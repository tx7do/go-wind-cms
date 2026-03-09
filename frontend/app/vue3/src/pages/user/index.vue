<script setup lang="ts">
import {definePage} from 'unplugin-vue-router/runtime'
import {ref, onMounted, computed} from 'vue'
import {useRouter} from 'vue-router'

import {useUserProfileStore} from '@/stores/modules/app/user-profile.state'
import {usePostStore} from '@/stores/modules/app/post.state'
import {$t} from '@/locales'
import type {contentservicev1_Post, identityservicev1_User} from '@/api/generated/app/service/v1'
import {useLanguageChangeEffect} from '@/hooks/useLanguageChangeEffect';

definePage({
  name: 'user-profile',
  meta: {
    title: 'User Profile',
    level: 2,
  },
})

const router = useRouter()
const userProfileStore = useUserProfileStore()
const postStore = usePostStore()

const loading = ref(false)
const postsLoading = ref(false)
const user = ref<identityservicev1_User | null>(null)
const activeTab = ref('posts')
const posts = ref<contentservicev1_Post[]>([])
const postsTotal = ref(0)

// 统计数据
const stats = computed(() => ({
  followers: user.value?.followers ?? 0,
  following: user.value?.following ?? 0,
  posts: user.value?.postCount ?? 0,
  likes: user.value?.likeCount ?? 0,
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

// 性别标签类型
const genderType = computed(() => {
  if (!user.value?.gender) return 'default'
  const typeMap: Record<string, any> = {
    'MALE': 'info',
    'FEMALE': 'default',
    'SECRET': 'default',
  }
  return typeMap[user.value.gender] || 'default'
})

// 获取用户信息
async function loadUserProfile() {
  loading.value = true
  try {
    const result = await userProfileStore.getMe()
    user.value = result || null

    // 如果成功获取用户信息且当前 tab 是 posts，自动加载帖子列表
    if (user.value && activeTab.value === 'posts') {
      await loadUserPosts()
    }
  } catch (error) {
    console.error('Load user profile failed:', error)
    user.value = null // 确保失败时设置为 null
  } finally {
    loading.value = false
  }
}

// 加载用户的帖子列表
async function loadUserPosts() {
  if (!user.value?.id) return

  postsLoading.value = true
  try {
    const result = await postStore.listPost(
      {page: 1, pageSize: 10},
      {author_id: user.value.id}, // 使用作者ID进行过滤
      null,
      ['-createdAt'] // 按创建时间倒序
    )

    if (result) {
      posts.value = result.items || []
      postsTotal.value = result.total || 0
    }
  } catch (error) {
    console.error('Load user posts failed:', error)
  } finally {
    postsLoading.value = false
  }
}

// 编辑资料
function handleEdit() {
  router.push('/settings')
}

// 查看帖子详情
function handleViewPost(postId: number) {
  router.push({
    path: `/post/${postId}`,
    query: {
      from: 'user'
    }
  })
}

// 监听 activeTab 变化，加载对应内容
async function handleTabChange(tabName: string) {
  if (tabName === 'posts' && posts.value.length === 0) {
    await loadUserPosts()
  }
}

onMounted(() => {
  loadUserProfile()
})

// 监听语言切换，自动重新加载用户数据和帖子
useLanguageChangeEffect(async () => {
  await Promise.all([
    loadUserProfile(),
    activeTab.value === 'posts' ? loadUserPosts() : Promise.resolve(),
  ]);
}, {
  immediate: false,      // 是否立即执行一次
  autoCleanup: true,    // 是否自动清理
});
</script>

<template>
  <div class="user-profile-page">
    <!-- Loading Skeleton -->
    <div v-if="loading" class="profile-loading">
      <div class="profile-header">
        <n-skeleton height="200px"/>
      </div>
      <div class="profile-main">
        <aside class="profile-sidebar">
          <div class="info-card">
            <n-skeleton text :repeat="8"/>
          </div>
        </aside>
        <main class="profile-content-area">
          <div class="content-card">
            <n-skeleton text :repeat="10"/>
          </div>
        </main>
      </div>
    </div>

    <!-- Loaded Content -->
    <template v-else>
      <!-- 顶部背景和基本信息 -->
      <div class="profile-header">
        <div class="header-bg"></div>
        <div class="profile-content">
          <!-- 用户基本信息 -->
          <div class="user-basic-info">
            <div class="avatar-section">
              <div class="avatar-wrapper">
                <n-avatar
                  :size="100"
                  :src="user?.avatar || undefined"
                  class="user-avatar"
                >
                  {{ user?.nickname?.charAt(0) || user?.username?.charAt(0) || 'U' }}
                </n-avatar>
              </div>
            </div>

            <div class="info-section">
              <div class="user-header">
                <div class="user-title">
                  <h1 class="user-name">{{ user?.nickname || user?.username }}</h1>
                  <n-tag v-if="user?.gender" size="small" :type="genderType">
                    {{ formatGender(user?.gender) }}
                  </n-tag>
                </div>
                <n-button type="primary" @click="handleEdit">
                  {{ $t('page.user.edit_profile') }}
                </n-button>
              </div>

              <p v-if="user?.description" class="user-bio">{{ user?.description }}</p>

              <div class="user-meta">
                  <span v-if="user?.region" class="meta-item">
                    {{ user?.region }}
                  </span>
                <span v-if="user?.tenantName" class="meta-item">
                    {{ user?.tenantName }}
                  </span>
                <span v-if="user?.positionNames?.length" class="meta-item">
                    {{ user?.positionNames.join(', ') }}
                  </span>
              </div>
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
      <div class="profile-main">
        <!-- 左侧面板 -->
        <aside class="profile-sidebar">
          <div class="info-card">
            <div class="card-section">
              <h3 class="section-title">{{ $t('page.user.basic_info') }}</h3>
              <div class="info-list">
                <div v-if="user?.username" class="info-row">
                  <span class="info-key">{{ $t('page.user.username') }}:</span>
                  <span class="info-val">{{ user.username }}</span>
                </div>
                <div v-if="user?.realname" class="info-row">
                  <span class="info-key">{{ $t('page.user.realname') }}:</span>
                  <span class="info-val">{{ user.realname }}</span>
                </div>
                <div v-if="user?.email" class="info-row">
                  <span class="info-key">{{ $t('page.user.email') }}:</span>
                  <span class="info-val">{{ user.email }}</span>
                </div>
                <div v-if="user?.mobile" class="info-row">
                  <span class="info-key">{{ $t('page.user.mobile') }}:</span>
                  <span class="info-val">{{ user.mobile }}</span>
                </div>
              </div>
            </div>

            <div class="divider"></div>

            <div class="card-section">
              <h3 class="section-title">{{ $t('page.user.account_info') }}</h3>
              <div class="info-list">
                <div class="info-row">
                  <span class="info-key">{{ $t('page.user.status') }}:</span>
                  <n-tag :type="statusType" size="small">
                    {{ formatStatus(user?.status) }}
                  </n-tag>
                </div>
                <div v-if="user?.roleNames?.length" class="info-row">
                  <span class="info-key">{{ $t('page.user.roles') }}:</span>
                  <span class="info-val">{{ user.roleNames.join(', ') }}</span>
                </div>
                <div v-if="user?.createdAt" class="info-row">
                  <span class="info-key">{{ $t('page.user.created_at') }}:</span>
                  <span class="info-val">{{ formatDateTime(user.createdAt) }}</span>
                </div>
                <div v-if="user?.lastLoginAt" class="info-row">
                  <span class="info-key">{{ $t('page.user.last_login_at') }}:</span>
                  <span class="info-val">{{ formatDateTime(user.lastLoginAt) }}</span>
                </div>
              </div>
            </div>
          </div>
        </aside>

        <!-- 主内容区 -->
        <main class="profile-content-area">
          <div class="content-card">
            <n-tabs v-model:value="activeTab" type="line" animated @update:value="handleTabChange">
              <n-tab-pane name="posts" :tab="$t('page.user.tab_posts')">
                <!-- Posts Loading Skeleton -->
                <div v-if="postsLoading" class="posts-list-loading">
                  <n-skeleton text :repeat="5"/>
                </div>
                <!-- Posts Loaded -->
                <div v-else>
                  <div v-if="posts.length > 0" class="posts-list">
                    <div v-for="post in posts" :key="post.id" class="post-item">
                      <div class="post-content">
                        <h3 class="post-title">{{ postStore.getPostTitle(post) }}</h3>
                        <p class="post-summary">{{ postStore.getPostSummary(post) }}</p>
                        <div class="post-meta">
                            <span class="meta-info">
                              <span class="i-carbon:view"/>
                              {{ post.visits || 0 }} {{ $t('page.user.views') }}
                            </span>
                          <span class="meta-info">
                              <span class="i-carbon:thumbs-up"/>
                              {{ post.likes || 0 }} {{ $t('page.user.likes') }}
                            </span>
                          <span class="meta-info">
                              <span class="i-carbon:chat"/>
                              {{ post.commentCount || 0 }} {{ $t('page.user.comments') }}
                            </span>
                          <span class="meta-info">
                              <span class="i-carbon:time"/>
                              {{ formatDateTime(post.createdAt) }}
                            </span>
                        </div>
                      </div>
                      <n-button text @click="handleViewPost(post.id)">
                        {{ $t('page.user.view_post') }} →
                      </n-button>
                    </div>
                  </div>
                  <n-empty v-else :description="$t('page.user.no_posts')">
                    <template #icon>
                      <span class="i-carbon:document"/>
                    </template>
                  </n-empty>
                </div>
              </n-tab-pane>
              <n-tab-pane name="activities" :tab="$t('page.user.tab_activities')">
                <n-empty :description="$t('page.user.no_activities')">
                  <template #icon>
                    <span class="i-carbon:activity"/>
                  </template>
                </n-empty>
              </n-tab-pane>
              <n-tab-pane name="collections" :tab="$t('page.user.tab_collections')">
                <n-empty :description="$t('page.user.no_collections')">
                  <template #icon>
                    <span class="i-carbon:bookmark"/>
                  </template>
                </n-empty>
              </n-tab-pane>
            </n-tabs>
          </div>
        </main>
      </div>
    </template>

    <!-- User Not Logged In -->
    <template v-if="!user">
      <div class="empty-state">
        <n-empty :description="$t('page.user.please_login')">
          <template #extra>
            <n-button type="primary" size="large" @click="router.push('/login')">
              {{ $t('authentication.login.title') }}
            </n-button>
          </template>
        </n-empty>
      </div>
    </template>
  </div>
</template>

<style scoped lang="less">
.user-profile-page {
  min-height: 100vh;
  background: var(--color-bg);
  overflow-x: hidden;
}

// Loading Skeleton Styles
.profile-loading {
  .profile-header {
    margin-bottom: 40px;

    :deep(.n-skeleton) {
      border-radius: 16px;
    }
  }

  .profile-main {
    display: flex;
    gap: 24px;
    padding: 24px;

    .profile-sidebar {
      width: 300px;
      flex-shrink: 0;

      .info-card {
        :deep(.n-skeleton) {
          margin-bottom: 12px;
        }
      }
    }

    .profile-content-area {
      flex: 1;

      .content-card {
        :deep(.n-skeleton) {
          margin-bottom: 12px;
        }
      }
    }
  }
}

.posts-list-loading {
  :deep(.n-skeleton) {
    margin-bottom: 16px;

    &:last-child {
      margin-bottom: 0;
    }
  }
}

.empty-state {
  min-height: calc(100vh - 120px);
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 60px 20px;
}

// ============ Header Section ============
.profile-header {
  position: relative;
  background: var(--color-surface);
  border-bottom: 1px solid var(--color-border);
  overflow: hidden;

  .header-bg {
    width: 100%;
    height: 200px;
    background: linear-gradient(135deg, #667eea 0%, #764ba2 50%, #a855f7 100%);
    position: relative;
    overflow: hidden;

    // 添加装饰性背景
    &::before {
      content: '';
      position: absolute;
      top: -50%;
      right: -10%;
      width: 60%;
      height: 200%;
      background: radial-gradient(circle, rgba(255, 255, 255, 0.1) 0%, transparent 70%);
      animation: float 20s ease-in-out infinite;
    }

    &::after {
      content: '';
      position: absolute;
      bottom: -30%;
      left: -5%;
      width: 50%;
      height: 150%;
      background: radial-gradient(circle, rgba(255, 255, 255, 0.08) 0%, transparent 70%);
      animation: float 25s ease-in-out infinite reverse;
    }
  }

  .profile-content {
    max-width: 1200px;
    margin: 0 auto;
    padding: 0 32px;
  }
}

// User Basic Info
.user-basic-info {
  display: flex;
  gap: 32px;
  align-items: flex-start;
  margin-top: -50px;
  padding-bottom: 24px;

  .avatar-section {
    flex-shrink: 0;

    .avatar-wrapper {
      position: relative;
      width: 100px;
      height: 100px;

      .user-avatar {
        width: 100%;
        height: 100%;
        border: 5px solid rgba(255, 255, 255, 0.2);
        box-shadow: 0 4px 16px rgba(0, 0, 0, 0.2), 0 0 0 1px rgba(102, 126, 234, 0.3);
        font-size: 40px;
        font-weight: 600;
        background: linear-gradient(135deg, var(--color-brand) 0%, #764ba2 100%);
        color: #fff;
      }
    }
  }

  // 深色模式下的avatar优化
  html.dark & {
    .avatar-section {
      .avatar-wrapper {
        .user-avatar {
          border: 5px solid rgba(255, 255, 255, 0.15);
          box-shadow: 0 6px 20px rgba(0, 0, 0, 0.6), 0 0 0 2px rgba(129, 147, 255, 0.4);
        }
      }
    }
  }

  .info-section {
    flex: 1;
    background: rgba(255, 255, 255, 0.98);
    padding: 24px 28px;
    border-radius: 16px;
    box-shadow: 0 8px 24px rgba(0, 0, 0, 0.12), 0 2px 8px rgba(0, 0, 0, 0.08);
    backdrop-filter: blur(12px);
    border: 1px solid rgba(255, 255, 255, 0.3);

    .user-header {
      display: flex;
      justify-content: space-between;
      align-items: flex-start;
      margin-bottom: 12px;
      gap: 20px;

      .user-title {
        display: flex;
        align-items: center;
        gap: 12px;

        .user-name {
          font-size: 32px;
          font-weight: 800;
          margin: 0;
          color: #1a1a1a;
          letter-spacing: -0.5px;
          text-shadow: none;
        }

        :deep(.n-tag) {
          margin-top: 6px;
          font-weight: 700;
          padding: 8px 20px;
          border-radius: 24px;
          border: none;
          box-shadow: 0 3px 10px rgba(0, 0, 0, 0.15);
          font-size: 13px;
          letter-spacing: 0.5px;
          transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
          position: relative;
          overflow: hidden;

          &::before {
            content: '';
            position: absolute;
            top: -50%;
            left: -50%;
            width: 200%;
            height: 200%;
            background: linear-gradient(
              45deg,
              transparent 30%,
              rgba(255, 255, 255, 0.3) 50%,
              transparent 70%
            );
            transform: rotate(45deg);
            animation: shimmer 3s infinite;
          }

          &.n-tag--info {
            background: linear-gradient(135deg, #1890ff 0%, #096dd9 50%, #0050b3 100%);
            color: #fff;
            box-shadow:
              0 4px 12px rgba(24, 144, 255, 0.4),
              inset 0 1px 0 rgba(255, 255, 255, 0.3);
          }

          &.n-tag--default {
            background: linear-gradient(135deg, #ff7875 0%, #f5222d 50%, #cf1322 100%);
            color: #fff;
            box-shadow:
              0 4px 12px rgba(255, 120, 117, 0.4),
              inset 0 1px 0 rgba(255, 255, 255, 0.3);
          }

          &:hover {
            transform: translateY(-2px) scale(1.05);
            box-shadow: 0 6px 16px rgba(0, 0, 0, 0.2);
          }
        }
      }

      button {
        flex-shrink: 0;
      }
    }

    .user-bio {
      font-size: 15px;
      color: #595959;
      line-height: 1.7;
      margin: 16px 0;
      max-width: 650px;
      padding-left: 4px;
    }

    .user-meta {
      display: flex;
      gap: 14px;
      flex-wrap: wrap;

      .meta-item {
        font-size: 14px;
        color: #595959;
        padding: 8px 16px;
        background: linear-gradient(135deg, #f5f5f5 0%, #fafafa 100%);
        border-radius: 10px;
        border: 1px solid #d9d9d9;
        transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
        font-weight: 600;
        box-shadow: 0 1px 3px rgba(0, 0, 0, 0.05);

        &:hover {
          color: #1890ff;
          background: linear-gradient(135deg, #e6f7ff 0%, #bae7ff 100%);
          border-color: #1890ff;
          transform: translateY(-2px);
          box-shadow: 0 4px 8px rgba(24, 144, 255, 0.2);
        }
      }
    }
  }

  // 深色模式下的 info-section
  html.dark & {
    .info-section {
      background: rgba(27, 31, 39, 0.85);
      box-shadow: 0 8px 24px rgba(0, 0, 0, 0.4), 0 2px 8px rgba(0, 0, 0, 0.3);
      border: 1px solid rgba(255, 255, 255, 0.1);

      .user-header {
        .user-title {
          .user-name {
            color: var(--color-text-primary);
            text-shadow: 0 1px 3px rgba(0, 0, 0, 0.3);
          }

          :deep(.n-tag) {
            box-shadow: 0 4px 16px rgba(0, 0, 0, 0.6), inset 0 1px 0 rgba(255, 255, 255, 0.2);

            &::before {
              background: linear-gradient(
                45deg,
                transparent 30%,
                rgba(255, 255, 255, 0.2) 50%,
                transparent 70%
              );
            }

            &.n-tag--info {
              background: linear-gradient(135deg, #40a9ff 0%, #1890ff 50%, #096dd9 100%);
              box-shadow:
                0 4px 16px rgba(64, 169, 255, 0.5),
                inset 0 1px 0 rgba(255, 255, 255, 0.25);
            }

            &.n-tag--default {
              background: linear-gradient(135deg, #ff7875 0%, #ff4d4f 50%, #f5222d 100%);
              box-shadow:
                0 4px 16px rgba(255, 120, 117, 0.5),
                inset 0 1px 0 rgba(255, 255, 255, 0.25);
            }

            &:hover {
              box-shadow:
                0 6px 20px rgba(0, 0, 0, 0.7),
                inset 0 1px 0 rgba(255, 255, 255, 0.3);
            }
          }
        }
      }

      .user-bio {
        color: var(--color-text-secondary);
      }

      .user-meta {
        .meta-item {
          color: var(--color-text-secondary);
          background: rgba(255, 255, 255, 0.05);
          border-color: var(--color-border);

          &:hover {
            color: var(--color-brand);
            background: var(--post-accent-bg-hover);
            border-color: var(--color-brand);
            transform: translateY(-2px);
            box-shadow: 0 4px 12px rgba(102, 126, 234, 0.3);
          }
        }
      }
    }
  }
}

// User Stats
.user-stats {
  display: flex;
  gap: 48px;
  padding: 28px 0;
  border-top: 1px solid var(--color-border);
  background: linear-gradient(135deg, rgba(102, 126, 234, 0.02) 0%, rgba(168, 85, 247, 0.02) 100%);

  .stat-item {
    flex: 1;
    text-align: center;
    cursor: pointer;
    padding: 16px;
    border-radius: 12px;
    transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
    position: relative;
    overflow: hidden;

    &::before {
      content: '';
      position: absolute;
      top: 0;
      left: 0;
      width: 100%;
      height: 100%;
      background: linear-gradient(135deg, rgba(102, 126, 234, 0.05) 0%, rgba(168, 85, 247, 0.05) 100%);
      opacity: 0;
      transition: opacity 0.3s;
    }

    &:hover {
      background: var(--color-bg);
      transform: translateY(-4px);
      box-shadow: 0 8px 20px rgba(0, 0, 0, 0.1);

      &::before {
        opacity: 1;
      }
    }

    .stat-value {
      font-size: 28px;
      font-weight: 800;
      margin-bottom: 10px;
      background: linear-gradient(135deg, #667eea 0%, #764ba2 50%, #a855f7 100%);
      -webkit-background-clip: text;
      -webkit-text-fill-color: transparent;
      background-clip: text;
      filter: drop-shadow(0 2px 4px rgba(102, 126, 234, 0.3));
    }

    .stat-label {
      font-size: 14px;
      color: var(--color-text-secondary);
      font-weight: 600;
      letter-spacing: 0.3px;
    }
  }
}

// ============ Main Content ============
.profile-main {
  max-width: 1200px;
  margin: 32px auto;
  padding: 0 32px;
  display: grid;
  grid-template-columns: 300px 1fr;
  gap: 24px;
  align-items: start;
}

// Sidebar
.profile-sidebar {
  .info-card {
    position: sticky;
    top: 100px;
    background: var(--color-surface);
    border: 1px solid var(--color-border);
    border-radius: 16px;
    padding: 24px;
    box-shadow: 0 2px 8px rgba(0, 0, 0, 0.06), 0 1px 3px rgba(0, 0, 0, 0.03);
    transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);

    &:hover {
      border-color: var(--post-accent-bg-active);
      box-shadow: 0 8px 24px rgba(0, 0, 0, 0.12), 0 4px 12px rgba(0, 0, 0, 0.08);
      transform: translateY(-2px);
    }

    .card-section {
      .section-title {
        font-size: 16px;
        font-weight: 800;
        color: var(--color-text-primary);
        margin: 0 0 18px 0;
        padding-bottom: 12px;
        border-bottom: 3px solid transparent;
        border-image-slice: 1;
        letter-spacing: 0.5px;
      }

      .info-list {
        display: flex;
        flex-direction: column;
        gap: 14px;

        .info-row {
          display: flex;
          justify-content: space-between;
          align-items: center;
          gap: 12px;
          font-size: 14px;
          line-height: 1.7;
          padding: 10px 0;
          border-radius: 8px;
          transition: all 0.25s cubic-bezier(0.4, 0, 0.2, 1);

          &:hover {
            background: linear-gradient(135deg, rgba(102, 126, 234, 0.05) 0%, rgba(168, 85, 247, 0.05) 100%);
            padding-left: 8px;
            padding-right: 8px;
            transform: translateX(2px);
          }

          .info-key {
            color: var(--color-text-secondary);
            flex-shrink: 0;
            font-weight: 700;
          }

          .info-val {
            color: var(--color-text-primary);
            text-align: right;
            word-break: break-word;
            flex: 1;
            font-weight: 600;
          }
        }
      }
    }

    .divider {
      height: 1px;
      background: var(--color-border);
      margin: 16px 0;
    }
  }
}

// Main Content Area
.profile-content-area {
  .content-card {
    background: var(--color-surface);
    border: 1px solid var(--color-border);
    border-radius: 16px;
    box-shadow: 0 2px 8px rgba(0, 0, 0, 0.06), 0 1px 3px rgba(0, 0, 0, 0.03);
    overflow: hidden;
    transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);

    &:hover {
      border-color: var(--post-accent-bg-active);
      box-shadow: 0 8px 24px rgba(0, 0, 0, 0.12), 0 4px 12px rgba(0, 0, 0, 0.08);
    }

    :deep(.n-tabs) {
      .n-tabs-nav {
        padding: 0 24px;
        background: linear-gradient(135deg, rgba(102, 126, 234, 0.03) 0%, rgba(168, 85, 247, 0.03) 100%);
        border-bottom: 1px solid var(--color-border);

        .n-tabs-tab {
          font-weight: 600;
          font-size: 15px;
          padding: 16px 0;
          margin-right: 24px;
          color: var(--color-text-secondary);
          transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
          position: relative;

          &::after {
            content: '';
            position: absolute;
            bottom: 0;
            left: 0;
            width: 0;
            height: 3px;
            background: linear-gradient(90deg, #667eea 0%, #764ba2 50%, #a855f7 100%);
            transition: width 0.3s ease;
          }

          &:hover {
            color: var(--color-text-primary);

            &::after {
              width: 100%;
            }
          }

          &.n-tabs-tab--active {
            color: var(--color-brand);
            font-weight: 700;

            &::after {
              width: 100%;
            }
          }
        }
      }

      .n-tab-pane {
        padding: 36px;
        min-height: 450px;
        background: var(--color-surface);
      }
    }

    :deep(.n-empty) {
      padding: 60px 20px;
    }
  }
}

// Posts List Styling
.posts-list {
  display: flex;
  flex-direction: column;
  gap: 18px;

  .post-item {
    display: flex;
    justify-content: space-between;
    align-items: flex-start;
    gap: 24px;
    padding: 24px;
    background: var(--color-bg);
    border: 1px solid var(--color-border);
    border-radius: 14px;
    transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
    box-shadow: 0 2px 6px rgba(0, 0, 0, 0.04);

    &:hover {
      border-color: var(--color-brand);
      box-shadow: 0 8px 24px rgba(102, 126, 234, 0.15), 0 4px 12px rgba(0, 0, 0, 0.08);
      transform: translateY(-4px) scale(1.01);
      background: linear-gradient(135deg, rgba(102, 126, 234, 0.03) 0%, rgba(168, 85, 247, 0.03) 100%);
    }

    .post-content {
      flex: 1;

      .post-title {
        font-size: 19px;
        font-weight: 700;
        color: var(--color-text-primary);
        margin: 0 0 10px 0;
        line-height: 1.5;
        cursor: pointer;
        transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
        letter-spacing: -0.3px;

        &:hover {
          color: var(--color-brand);
          transform: translateX(4px);
        }
      }

      .post-summary {
        font-size: 14px;
        color: var(--color-text-secondary);
        line-height: 1.8;
        margin: 0 0 14px 0;
        display: -webkit-box;
        -webkit-line-clamp: 2;
        -webkit-box-orient: vertical;
        overflow: hidden;
        padding-left: 2px;
      }

      .post-meta {
        display: flex;
        gap: 20px;
        flex-wrap: wrap;

        .meta-info {
          display: flex;
          align-items: center;
          gap: 4px;
          font-size: 13px;
          color: var(--color-text-secondary);

          span[class^="i-"] {
            font-size: 16px;
          }
        }
      }
    }

    :deep(.n-button) {
      color: var(--color-brand);
      font-weight: 500;

      &:hover {
        color: var(--color-brand);
        background: var(--post-accent-bg-hover);
      }
    }
  }
}

// ============ Responsive Design ============
@media (max-width: 1024px) {
  .profile-main {
    grid-template-columns: 1fr;
    gap: 20px;
    margin: 24px auto;
    padding: 0 24px;
  }

  .profile-sidebar {
    .info-card {
      position: static;
      display: grid;
      grid-template-columns: 1fr 1fr;
      gap: 0;

      .card-section {
        padding: 16px 12px;

        &:nth-child(2) {
          border-left: 1px solid var(--color-border);
        }

        .section-title {
          font-size: 14px;
          margin-bottom: 12px;
        }

        .info-list {
          gap: 10px;

          .info-row {
            font-size: 13px;
          }
        }
      }

      .divider {
        display: none;
      }
    }
  }

  .user-stats {
    gap: 24px;

    .stat-item {
      .stat-value {
        font-size: 22px;
      }
    }
  }
}

@media (max-width: 768px) {
  .profile-header {
    .header-bg {
      height: 140px;
    }

    .profile-content {
      padding: 0 20px;
    }
  }

  .user-basic-info {
    flex-direction: column;
    align-items: center;
    text-align: center;
    margin-top: -40px;
    gap: 20px;

    .avatar-section {
      .avatar-wrapper {
        width: 80px;
        height: 80px;

        .user-avatar {
          font-size: 32px;
        }
      }
    }

    .info-section {
      width: 100%;
      padding-top: 0;

      .user-header {
        flex-direction: column;
        align-items: center;
        margin-bottom: 12px;

        .user-title {
          justify-content: center;

          .user-name {
            font-size: 24px;
          }
        }

        button {
          width: 100%;
        }
      }

      .user-bio {
        text-align: center;
      }

      .user-meta {
        justify-content: center;
        gap: 12px;
      }
    }
  }

  .user-stats {
    gap: 12px;
    flex-wrap: wrap;
    padding: 16px 0;

    .stat-item {
      flex: 0 1 calc(50% - 6px);

      .stat-value {
        font-size: 20px;
      }

      .stat-label {
        font-size: 12px;
      }
    }
  }

  .profile-main {
    grid-template-columns: 1fr;
    margin: 20px auto;
    padding: 0 16px;
    gap: 16px;
  }

  .profile-sidebar {
    .info-card {
      grid-template-columns: 1fr;

      .card-section {
        padding: 12px 0;

        &:nth-child(2) {
          border-left: none;
          border-top: 1px solid var(--color-border);
          padding-top: 12px;
        }

        .section-title {
          font-size: 13px;
          margin-bottom: 10px;
        }

        .info-list {
          gap: 8px;

          .info-row {
            font-size: 12px;
            flex-direction: column;
            align-items: flex-start;

            .info-val {
              text-align: left;
            }
          }
        }
      }
    }
  }

  .profile-content-area {
    .content-card {
      :deep(.n-tabs) {
        .n-tabs-nav {
          padding: 0 16px;

          .n-tabs-tab {
            font-size: 14px;
            margin-right: 12px;
            padding: 12px 0;
          }
        }

        .n-tab-pane {
          padding: 20px;
          min-height: 300px;
        }
      }
    }
  }
}

@media (max-width: 480px) {
  .profile-header {
    .header-bg {
      height: 120px;
    }

    .profile-content {
      padding: 0 12px;
    }
  }

  .user-basic-info {
    margin-top: -35px;

    .avatar-section {
      .avatar-wrapper {
        width: 70px;
        height: 70px;

        .user-avatar {
          font-size: 28px;
          border-width: 3px;
        }
      }
    }

    .info-section {
      .user-header {
        .user-title {
          .user-name {
            font-size: 20px;
          }
        }
      }
    }
  }

  .user-stats {
    gap: 8px;

    .stat-item {
      min-width: calc(50% - 4px);
      padding: 8px;

      .stat-value {
        font-size: 18px;
      }
    }
  }

  .profile-main {
    padding: 0 12px;
    margin: 16px auto;
    gap: 12px;
  }

  .profile-content-area {
    .content-card {
      :deep(.n-tabs) {
        .n-tabs-nav {
          padding: 0 12px;

          .n-tabs-tab {
            font-size: 13px;
            margin-right: 8px;
          }
        }

        .n-tab-pane {
          padding: 16px;
          min-height: 250px;
        }
      }
    }
  }
}

// 动画关键帧定义
@keyframes float {
  0%, 100% {
    transform: translate(0, 0) rotate(0deg);
  }
  33% {
    transform: translate(15px, -20px) rotate(8deg);
  }
  66% {
    transform: translate(-15px, 15px) rotate(-8deg);
  }
}

@keyframes shimmer {
  0% {
    transform: translateX(-100%) rotate(45deg);
  }
  100% {
    transform: translateX(100%) rotate(45deg);
  }
}
</style>
