<script setup lang="ts">
import { definePage } from 'unplugin-vue-router/runtime'
import { ref, onMounted, computed } from 'vue'
import { useRouter } from 'vue-router'
import { useUserProfileStore } from '@/stores/modules/app/user-profile.state'
import { $t } from '@/locales'
import type { identityservicev1_User } from '@/api/generated/app/service/v1'

definePage({
  name: 'user-profile',
  meta: {
    title: 'User Profile',
    level: 2,
  },
})

const router = useRouter()
const userProfileStore = useUserProfileStore()

const loading = ref(false)
const user = ref<identityservicev1_User | null>(null)
const activeTab = ref('posts')

// 模拟统计数据
const stats = computed(() => ({
  followers: 196,
  following: 1,
  posts: 2,
  likes: 947,
}))

// 格式化性别
const formatGender = (gender?: string) => {
  if (!gender) return $t('page.user.gender_secret')
  const genderMap: Record<string, string> = {
    'SECRET': $t('page.user.gender_secret'),
    'MALE': $t('page.user.gender_male'),
    'FEMALE': $t('page.user.gender_female'),
  }
  return genderMap[gender] || gender
}

// 格式化状态
const formatStatus = (status?: string) => {
  if (!status) return '-'
  const statusMap: Record<string, string> = {
    'NORMAL': $t('page.user.status_normal'),
    'DISABLED': $t('page.user.status_disabled'),
    'PENDING': $t('page.user.status_pending'),
    'LOCKED': $t('page.user.status_locked'),
    'EXPIRED': $t('page.user.status_expired'),
    'CLOSED': $t('page.user.status_closed'),
  }
  return statusMap[status] || status
}

// 格式化日期时间
const formatDateTime = (timestamp?: any) => {
  if (!timestamp) return '-'
  try {
    let date: Date
    if (timestamp.seconds) {
      date = new Date(timestamp.seconds * 1000)
    } else if (typeof timestamp === 'string') {
      date = new Date(timestamp)
    } else {
      return '-'
    }
    return date.toLocaleString('zh-CN', {
      year: 'numeric',
      month: '2-digit',
      day: '2-digit',
      hour: '2-digit',
      minute: '2-digit',
    })
  } catch (error) {
    return '-'
  }
}

// 状态标签类型
const statusType = computed(() => {
  if (!user.value?.status) return 'default'
  const typeMap: Record<string, any> = {
    'NORMAL': 'success',
    'DISABLED': 'error',
    'PENDING': 'warning',
    'LOCKED': 'error',
    'EXPIRED': 'warning',
    'CLOSED': 'error',
  }
  return typeMap[user.value.status] || 'default'
})

// 获取用户信息
async function loadUserProfile() {
  loading.value = true
  try {
    const result = await userProfileStore.getMe()
    if (result) {
      user.value = result
    } else {
      window.$message?.error($t('page.user.load_failed'))
    }
  } catch (error) {
    console.error('Load user profile failed:', error)
    window.$message?.error($t('page.user.load_failed'))
  } finally {
    loading.value = false
  }
}

// 编辑资料
function handleEdit() {
  router.push('/settings')
}

onMounted(() => {
  loadUserProfile()
})
</script>

<template>
  <div class="user-profile-page">
    <n-spin :show="loading">
      <template v-if="user">
        <div class="profile-container">
        <!-- 顶部背景和基本信息 -->
        <div class="profile-header">
          <div class="header-bg" />
          <div class="header-content">
            <div class="user-basic">
              <n-avatar
                :size="80"
                :src="user.avatar || undefined"
                class="user-avatar"
              >
                {{ user.nickname?.charAt(0) || user.username?.charAt(0) || 'U' }}
              </n-avatar>
              <div class="user-info">
                <div class="user-name-row">
                  <h1 class="user-name">{{ user.nickname || user.username || $t('page.user.not_login') }}</h1>
                  <div class="user-badges">
                    <n-tag v-if="user.gender" size="small" :bordered="false" class="gender-tag">
                      {{ formatGender(user.gender) }}
                    </n-tag>
                  </div>
                </div>
                <div v-if="user.description" class="user-bio">{{ user.description }}</div>
                <div class="user-meta-info">
                  <span v-if="user.region" class="meta-item">
                    <span class="i-carbon:location" />
                    {{ user.region }}
                  </span>
                  <span v-if="user.tenantName" class="meta-item">
                    <span class="i-carbon:enterprise" />
                    {{ user.tenantName }}
                  </span>
                  <span v-if="user.positionNames && user.positionNames.length > 0" class="meta-item">
                    <span class="i-carbon:user-role" />
                    {{ user.positionNames.join(', ') }}
                  </span>
                </div>
              </div>
              <div class="user-actions">
                <n-button type="primary" @click="handleEdit">
                  <template #icon>
                    <span class="i-carbon:edit" />
                  </template>
                  {{ $t('page.user.edit_profile') }}
                </n-button>
              </div>
            </div>

            <!-- 统计数据 -->
            <div class="user-stats">
              <div class="stat-item">
                <div class="stat-value">{{ stats.following }}</div>
                <div class="stat-label">{{ $t('page.user.following') }}</div>
              </div>
              <div class="stat-item">
                <div class="stat-value">{{ stats.followers }}</div>
                <div class="stat-label">{{ $t('page.user.followers') }}</div>
              </div>
              <div class="stat-item">
                <div class="stat-value">{{ stats.posts }}</div>
                <div class="stat-label">{{ $t('page.user.posts') }}</div>
              </div>
              <div class="stat-item">
                <div class="stat-value">{{ stats.likes }}</div>
                <div class="stat-label">{{ $t('page.user.likes_received') }}</div>
              </div>
            </div>
          </div>
        </div>

        <!-- 内容区域 -->
        <div class="profile-body">
          <!-- 左侧边栏 -->
          <aside class="profile-sidebar">
            <n-card :bordered="false" class="info-card">
              <div class="info-section">
                <h3 class="section-title">{{ $t('page.user.basic_info') }}</h3>
                <div class="info-list">
                  <div v-if="user.username" class="info-item">
                    <span class="info-label">{{ $t('page.user.username') }}</span>
                    <span class="info-value">{{ user.username }}</span>
                  </div>
                  <div v-if="user.realname" class="info-item">
                    <span class="info-label">{{ $t('page.user.realname') }}</span>
                    <span class="info-value">{{ user.realname }}</span>
                  </div>
                  <div v-if="user.email" class="info-item">
                    <span class="info-label">{{ $t('page.user.email') }}</span>
                    <span class="info-value">{{ user.email }}</span>
                  </div>
                  <div v-if="user.mobile" class="info-item">
                    <span class="info-label">{{ $t('page.user.mobile') }}</span>
                    <span class="info-value">{{ user.mobile }}</span>
                  </div>
                </div>
              </div>

              <n-divider style="margin: 16px 0" />

              <div class="info-section">
                <h3 class="section-title">{{ $t('page.user.account_info') }}</h3>
                <div class="info-list">
                  <div class="info-item">
                    <span class="info-label">{{ $t('page.user.status') }}</span>
                    <n-tag :type="statusType" size="small">
                      {{ formatStatus(user.status) }}
                    </n-tag>
                  </div>
                  <div v-if="user.roleNames && user.roleNames.length > 0" class="info-item">
                    <span class="info-label">{{ $t('page.user.roles') }}</span>
                    <span class="info-value">{{ user.roleNames.join(', ') }}</span>
                  </div>
                  <div v-if="user.createdAt" class="info-item">
                    <span class="info-label">{{ $t('page.user.created_at') }}</span>
                    <span class="info-value">{{ formatDateTime(user.createdAt) }}</span>
                  </div>
                  <div v-if="user.lastLoginAt" class="info-item">
                    <span class="info-label">{{ $t('page.user.last_login_at') }}</span>
                    <span class="info-value">{{ formatDateTime(user.lastLoginAt) }}</span>
                  </div>
                </div>
              </div>
            </n-card>
          </aside>

          <!-- 主内容区 -->
          <main class="profile-main">
            <n-card :bordered="false" class="content-card">
              <n-tabs v-model:value="activeTab" type="line" animated>
                <n-tab-pane name="posts" :tab="$t('page.user.tab_posts')">
                  <n-empty :description="$t('page.user.no_posts')" style="padding: 60px 0">
                    <template #icon>
                      <span class="i-carbon:document" style="font-size: 48px" />
                    </template>
                  </n-empty>
                </n-tab-pane>
                <n-tab-pane name="activities" :tab="$t('page.user.tab_activities')">
                  <n-empty :description="$t('page.user.no_activities')" style="padding: 60px 0">
                    <template #icon>
                      <span class="i-carbon:activity" style="font-size: 48px" />
                    </template>
                  </n-empty>
                </n-tab-pane>
                <n-tab-pane name="collections" :tab="$t('page.user.tab_collections')">
                  <n-empty :description="$t('page.user.no_collections')" style="padding: 60px 0">
                    <template #icon>
                      <span class="i-carbon:bookmark" style="font-size: 48px" />
                    </template>
                  </n-empty>
                </n-tab-pane>
              </n-tabs>
            </n-card>
          </main>
        </div>
      </div>
      </template>

      <template v-else>
        <div class="empty-container">
          <n-empty :description="$t('page.user.please_login')">
            <template #extra>
              <n-button type="primary" size="large" @click="router.push('/login')">
                {{ $t('authentication.login.title') }}
              </n-button>
            </template>
          </n-empty>
        </div>
      </template>
    </n-spin>
  </div>
</template>

<style scoped lang="less">
.user-profile-page {
  width: 100%;
  min-height: calc(100vh - 64px);
  background: var(--n-color-base);
}

.profile-container {
  width: 100%;
}

.empty-container {
  min-height: calc(100vh - 200px);
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 60px 20px;
}

// 顶部区域
.profile-header {
  position: relative;
  background: var(--n-card-color);
  border-bottom: 1px solid var(--n-border-color);

  .header-bg {
    width: 100%;
    height: 160px;
    background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
    position: relative;

    &::after {
      content: '';
      position: absolute;
      bottom: 0;
      left: 0;
      right: 0;
      height: 60px;
      background: linear-gradient(to bottom, transparent, rgba(0, 0, 0, 0.1));
    }
  }

  .header-content {
    max-width: 1200px;
    margin: 0 auto;
    padding: 0 32px;
  }

  .user-basic {
    display: flex;
    gap: 24px;
    align-items: flex-start;
    margin-top: -50px;
    padding-bottom: 24px;

    .user-avatar {
      flex-shrink: 0;
      border: 5px solid var(--n-card-color);
      box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
      font-size: 36px;
      font-weight: 600;
      background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
      color: #fff;
    }

    .user-info {
      flex: 1;
      padding-top: 20px;

      .user-name-row {
        display: flex;
        align-items: center;
        gap: 12px;
        margin-bottom: 10px;

        .user-name {
          font-size: 26px;
          font-weight: 700;
          color: var(--n-text-color);
          margin: 0;
          letter-spacing: -0.5px;
        }

        .user-badges {
          display: flex;
          gap: 8px;
          align-items: center;

          .gender-tag {
            background: var(--n-color-target);
            font-weight: 500;
          }
        }
      }

      .user-bio {
        color: var(--n-text-color-2);
        font-size: 15px;
        line-height: 1.6;
        margin-bottom: 14px;
        max-width: 600px;
      }

      .user-meta-info {
        display: flex;
        gap: 24px;
        flex-wrap: wrap;

        .meta-item {
          display: flex;
          align-items: center;
          gap: 6px;
          color: var(--n-text-color-3);
          font-size: 14px;
          transition: color 0.3s;

          &:hover {
            color: var(--n-text-color-2);
          }

          span[class^="i-"] {
            font-size: 18px;
            opacity: 0.8;
          }
        }
      }
    }

    .user-actions {
      padding-top: 20px;

      button {
        padding: 0 24px;
        height: 36px;
        font-weight: 500;
        box-shadow: 0 2px 8px rgba(102, 126, 234, 0.3);

        &:hover {
          box-shadow: 0 4px 12px rgba(102, 126, 234, 0.4);
        }
      }
    }
  }

  .user-stats {
    display: flex;
    gap: 48px;
    padding: 20px 0;
    border-top: 1px solid var(--n-border-color);

    .stat-item {
      display: flex;
      flex-direction: column;
      align-items: center;
      cursor: pointer;
      transition: all 0.3s ease;
      padding: 8px 16px;
      border-radius: 8px;
      min-width: 80px;

      &:hover {
        background: var(--n-color-target);
        transform: translateY(-3px);
      }

      .stat-value {
        font-size: 24px;
        font-weight: 700;
        color: var(--n-text-color);
        margin-bottom: 6px;
        background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
        -webkit-background-clip: text;
        -webkit-text-fill-color: transparent;
        background-clip: text;
      }

      .stat-label {
        font-size: 13px;
        color: var(--n-text-color-3);
        font-weight: 500;
      }
    }
  }
}

// 内容区域
.profile-body {
  max-width: 1200px;
  margin: 32px auto;
  padding: 0 32px;
  display: grid;
  grid-template-columns: 320px 1fr;
  gap: 32px;
  align-items: start;
}

// 左侧边栏
.profile-sidebar {
  position: sticky;
  top: 90px;

  .info-card {
    box-shadow: 0 2px 8px rgba(0, 0, 0, 0.04);
    border-radius: 12px;
    overflow: hidden;

    :deep(.n-card__content) {
      padding: 24px;
    }

    .info-section {
      .section-title {
        font-size: 16px;
        font-weight: 700;
        color: var(--n-text-color);
        margin: 0 0 20px 0;
        padding-bottom: 12px;
        border-bottom: 2px solid var(--n-border-color);
        letter-spacing: -0.3px;
      }

      .info-list {
        display: flex;
        flex-direction: column;
        gap: 16px;

        .info-item {
          display: flex;
          justify-content: space-between;
          align-items: flex-start;
          gap: 16px;
          font-size: 14px;
          line-height: 1.6;

          .info-label {
            color: var(--n-text-color-3);
            flex-shrink: 0;
            font-weight: 500;
            min-width: 70px;
          }

          .info-value {
            color: var(--n-text-color);
            text-align: right;
            word-break: break-word;
            flex: 1;
            font-weight: 400;
          }
        }
      }
    }
  }
}

// 主内容区
.profile-main {
  min-height: 500px;

  .content-card {
    box-shadow: 0 2px 8px rgba(0, 0, 0, 0.04);
    border-radius: 12px;
    overflow: hidden;

    :deep(.n-card__content) {
      padding: 0;
    }

    :deep(.n-tabs) {
      .n-tabs-nav {
        padding: 0 32px;
        background: var(--n-card-color);

        .n-tabs-tab {
          font-weight: 500;
          font-size: 15px;
          padding: 16px 0;
        }
      }

      .n-tab-pane {
        padding: 32px;
        min-height: 400px;
      }
    }
  }
}

// 响应式设计
@media (max-width: 968px) {
  .profile-body {
    grid-template-columns: 1fr;
    gap: 24px;
    margin: 24px auto;
  }

  .profile-sidebar {
    position: static;
  }

  .profile-header {
    .user-basic {
      flex-direction: column;
      align-items: center;
      text-align: center;
      margin-top: -60px;

      .user-info {
        display: flex;
        flex-direction: column;
        align-items: center;

        .user-name-row {
          justify-content: center;
        }

        .user-bio {
          text-align: center;
        }

        .user-meta-info {
          justify-content: center;
        }
      }

      .user-actions {
        padding-top: 0;
      }
    }

    .user-stats {
      justify-content: center;
      gap: 32px;
    }
  }
}

@media (max-width: 640px) {
  .profile-header {
    .header-bg {
      height: 140px;
    }

    .header-content {
      padding: 0 16px;
    }

    .user-basic {
      gap: 16px;
      margin-top: -50px;
      padding-bottom: 16px;

      .user-avatar {
        width: 70px !important;
        height: 70px !important;
        font-size: 28px;
        border-width: 4px;
      }

      .user-info {
        .user-name-row {
          .user-name {
            font-size: 22px;
          }
        }

        .user-bio {
          font-size: 14px;
        }

        .user-meta-info {
          gap: 16px;

          .meta-item {
            font-size: 13px;
          }
        }

        .user-actions {
          width: 100%;

          button {
            width: 100%;
          }
        }
      }
    }

    .user-stats {
      gap: 12px;
      flex-wrap: wrap;
      padding: 18px 0 16px;

      .stat-item {
        min-width: calc(50% - 6px);
        padding: 10px 16px;

        .stat-value {
          font-size: 22px;
          margin-bottom: 6px;
        }

        .stat-label {
          font-size: 13px;
        }
      }
    }
  }

  .profile-body {
    padding: 0 16px;
    margin: 16px auto;
    gap: 16px;
  }

  .profile-main {
    .content-card {
      :deep(.n-tabs) {
        .n-tabs-nav {
          padding: 0 16px;

          .n-tabs-tab {
            font-size: 14px;
            padding: 12px 0;
          }
        }

        .n-tab-pane {
          padding: 20px 16px;
          min-height: 300px;
        }
      }
    }
  }
}
</style>
