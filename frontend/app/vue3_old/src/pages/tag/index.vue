<script setup lang="ts">
import {definePage} from 'unplugin-vue-router/runtime'
import {ref, onMounted} from 'vue'
import {useRouter} from 'vue-router'
import {useMessage} from 'naive-ui'

import {$t} from '@/locales'
import {useTagStore} from '@/stores'
import type {
  contentservicev1_Tag,
} from "@/api/generated/app/service/v1";
import {useLanguageChangeEffect} from '@/hooks/useLanguageChangeEffect';

definePage({
  name: 'tag-list',
  meta: {
    title: 'Tags',
    level: 2,
  },
})

const router = useRouter()
const tagStore = useTagStore()
const message = useMessage()

const loading = ref(false)
const tags = ref<contentservicev1_Tag[]>([])
const page = ref(1)
const pageSize = ref(20)
const total = ref(0)

async function loadTags() {
  loading.value = true
  try {
    const res = await tagStore.listTag(
      {page: page.value, pageSize: pageSize.value},
      {status: 'TAG_STATUS_ACTIVE'}
    )
    tags.value = res.items || []
    total.value = res.total || 0
  } catch (error) {
    console.error('Load tags failed:', error)
    message.error($t('page.post_detail.load_failed'))
  } finally {
    loading.value = false
  }
}

function handleTagClick(id: number) {
  router.push(`/tag/${id}`)
}

function handlePageChange(newPage: number) {
  page.value = newPage
  loadTags()
}

function handlePageSizeChange(newSize: number) {
  pageSize.value = newSize
  page.value = 1
  loadTags()
}

onMounted(async () => {
  await loadTags()
})

// 监听语言切换，自动重新加载数据
useLanguageChangeEffect(async () => {
  await loadTags();
}, {
  immediate: false,      // 是否立即执行一次
  autoCleanup: true,    // 是否自动清理
});
</script>

<template>
  <div class="tag-list-page">
    <!-- Hero Section -->
    <div class="hero-section">
      <div class="hero-content">
        <h1>{{ $t('page.tags.tags_list') }}</h1>
        <p>{{ $t('page.tags.explore_all') }}</p>
        <div class="tag-stats">
          <div class="stat-item">
            <span class="i-carbon:tag"/>
            <span>{{ tags.length }} {{ $t('page.tags.total_tags') }}</span>
          </div>
        </div>
      </div>
    </div>

    <!-- Tags Grid -->
    <div class="page-container">
      <!-- Loading Skeleton -->
      <div v-if="loading" class="tags-grid">
        <div v-for="i in 12" :key="i" class="tag-card-skeleton">
          <n-skeleton height="120px"/>
        </div>
      </div>

      <!-- Tags List -->
      <div v-else-if="tags.length > 0 || total > 0" class="tags-grid">
        <div
          v-for="tag in tags"
          :key="tag.id"
          class="tag-card"
          :style="{
            borderColor: tag.color || 'var(--color-border)',
            background: `linear-gradient(135deg, ${tag.color}10 0%, var(--color-surface) 100%)`
          }"
          @click="handleTagClick(tag.id)"
        >
          <div class="tag-icon" :style="{color: tag.color || 'var(--color-brand)'}">
            <span :class="`i-carbon:${tag.icon || 'tag'}`"/>
          </div>
          <div class="tag-content">
            <h3>{{ tagStore.getTranslation(tag)?.name || $t('page.tags.tag_untitled') }}</h3>
            <p class="tag-description">
              {{ tagStore.getTranslation(tag)?.description || '' }}
            </p>
            <div class="tag-meta">
              <span class="meta-icon">
                <span class="i-carbon:document"/>
              </span>
              <span class="meta-text">{{ tag.postCount || 0 }} {{
                  $t('page.posts.articles')
                }}</span>
            </div>
          </div>
        </div>
        <div v-if="tags.length === 0 && total > 0" class="empty-page">
          <n-empty :description="$t('page.tags.no_tags_in_page')" />
        </div>
      </div>
      <!-- Empty State -->
      <div v-else class="empty-state">
        <n-empty :description="$t('page.tags.no_tags')"/>
      </div>
      <n-pagination
        v-if="total > pageSize"
        :page="page"
        :page-size="pageSize"
        :page-count="Math.ceil(total / pageSize)"
        :page-sizes="[10, 20, 50, 100]"
        :page-slot="7"
        :show-size-picker="true"
        @update:page="handlePageChange"
        @update:page-size="handlePageSizeChange"
        style="margin: 32px auto 0; display: flex; justify-content: center;"
      />
    </div>
  </div>
</template>

<style scoped lang="less">
.tag-list-page {
  min-height: 100vh;
  background: var(--color-bg);
}

// Hero Section
.hero-section {
  background: linear-gradient(135deg, #6366f1 0%, #8b5cf6 50%, #a855f7 100%);
  padding: 3rem 2rem;
  margin-bottom: 40px;
  position: relative;
  overflow: hidden;
  min-height: 300px;
  height: auto;
  display: flex;
  align-items: center;
  justify-content: center;
  text-align: center;
  color: white;

  &::before {
    content: '';
    position: absolute;
    top: 0;
    left: 0;
    right: 0;
    bottom: 0;
    background: radial-gradient(ellipse at 20% 50%, rgba(100, 200, 255, 0.15) 0%, transparent 50%),
    radial-gradient(ellipse at 80% 80%, rgba(200, 100, 255, 0.12) 0%, transparent 50%),
    linear-gradient(135deg, #6366f1 0%, #8b5cf6 50%, #a855f7 100%);
    animation: gradientShift 15s ease-in-out infinite;
    z-index: 0;
  }

  &::after {
    content: '';
    position: absolute;
    top: 0;
    left: 0;
    right: 0;
    bottom: 0;
    background-image: linear-gradient(rgba(255, 255, 255, 0.05) 1px, transparent 1px),
    linear-gradient(90deg, rgba(255, 255, 255, 0.05) 1px, transparent 1px);
    background-size: 50px 50px;
    animation: gridMove 20s linear infinite;
    opacity: 0.6;
    z-index: 1;
  }

  .hero-content {
    max-width: 1200px;
    margin: 0 auto;
    text-align: center;
    position: relative;
    z-index: 2;
    animation: fadeInUp 0.8s ease-out;

    h1 {
      font-size: 48px;
      font-weight: 800;
      margin: 0 0 16px 0;
      color: #ffffff;
      letter-spacing: -1px;
      line-height: 1.1;

      background: linear-gradient(135deg,
      #ffffff 0%,
      #f0f0ff 25%,
      #e0e0ff 50%,
      #f0f0ff 75%,
      #ffffff 100%);
      background-size: 200% auto;
      -webkit-background-clip: text;
      -webkit-text-fill-color: transparent;
      background-clip: text;

      text-shadow: 0 0 40px rgba(255, 255, 255, 0.9),
      0 0 80px rgba(var(--color-primary-purple-rgb), 0.7),
      0 0 120px rgba(99, 102, 241, 0.5),
      0 6px 24px rgba(0, 0, 0, 0.5),
      0 3px 12px rgba(0, 0, 0, 0.4);

      filter: drop-shadow(0 0 30px rgba(var(--color-primary-purple-rgb), 0.6)) drop-shadow(0 0 60px rgba(99, 102, 241, 0.4)) drop-shadow(0 5px 15px rgba(0, 0, 0, 0.3));

      animation-name: slideDown, glowPulseTitle, gradientShine;
      animation-duration: 0.8s, 3s, 6s;
      animation-timing-function: cubic-bezier(0.34, 1.56, 0.64, 1), ease-in-out, linear;
      animation-iteration-count: 1, infinite, infinite;
      animation-delay: 0s, 0.5s, 0s;
    }

    p {
      font-size: 20px;
      color: rgba(255, 255, 255, 0.95);
      margin: 0 0 24px 0;
      font-weight: 500;
      line-height: 1.5;
      text-shadow: 0 2px 10px rgba(0, 0, 0, 0.25);
      animation: fadeInUp 0.8s ease-out 0.2s both;
    }

    .tag-stats {
      display: flex;
      justify-content: center;
      gap: 32px;
      flex-wrap: wrap;

      .stat-item {
        display: flex;
        align-items: center;
        gap: 8px;
        font-size: 16px;
        font-weight: 600;
        color: rgba(255, 255, 255, 0.9);
        text-shadow: 0 1px 4px rgba(0, 0, 0, 0.2);
        animation: fadeInUp 0.8s ease-out 0.3s both;

        span[class^="i-"] {
          font-size: 20px;
        }
      }
    }
  }
}

// Page Container
.page-container {
  max-width: 1400px;
  margin: 0 auto;
  padding: 0 32px 80px;
}

// Tags Grid
.tags-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
  gap: 24px;
}

// Tag Card
.tag-card {
  background: var(--color-surface);
  border-radius: 16px;
  padding: 24px;
  cursor: pointer;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  border: 2px solid var(--color-border);
  height: 100%;
  display: flex;
  flex-direction: column;
  position: relative;
  overflow: hidden;

  &::before {
    content: '';
    position: absolute;
    top: 0;
    left: 0;
    right: 0;
    height: 4px;
    background: linear-gradient(90deg, var(--color-brand), var(--color-brand-transparent));
    opacity: 0;
    transition: opacity 0.3s;
  }

  &:hover {
    transform: translateY(-6px);
    box-shadow: 0 12px 32px rgba(0, 0, 0, 0.12);
    border-color: var(--color-brand);

    &::before {
      opacity: 1;
    }

    .tag-icon {
      transform: scale(1.15) rotate(5deg);
    }

    h3 {
      color: var(--color-brand);
    }
  }

  .tag-icon {
    font-size: 48px;
    margin-bottom: 16px;
    transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
    display: inline-flex;
    align-items: center;
    justify-content: center;
  }

  .tag-content {
    flex: 1;
    display: flex;
    flex-direction: column;

    h3 {
      font-size: 18px;
      font-weight: 700;
      margin: 0 0 10px 0;
      color: var(--color-text-primary);
      line-height: 1.4;
      transition: color 0.3s;
    }

    .tag-description {
      color: var(--color-text-secondary);
      font-size: 14px;
      line-height: 1.6;
      margin: 0 0 16px 0;
      display: -webkit-box;
      -webkit-line-clamp: 2;
      -webkit-box-orient: vertical;
      overflow: hidden;
      flex: 1;
    }

    .tag-meta {
      display: flex;
      align-items: center;
      gap: 6px;
      font-size: 13px;
      color: var(--color-text-secondary);
      font-weight: 500;
      padding-top: 12px;
      border-top: 1px solid var(--color-border);

      .meta-icon {
        display: flex;
        align-items: center;
        font-size: 14px;
        opacity: 0.8;
      }
    }
  }
}

// Skeleton
.tag-card-skeleton {
  background: var(--color-surface);
  border-radius: 16px;
  padding: 24px;
  height: 180px;
}

// Empty State
.empty-state {
  padding: 80px 20px;
  text-align: center;
}

// Responsive Design
@media (max-width: 1024px) {
  .hero-section {
    padding: 3rem 1.5rem 2.5rem;
    min-height: 280px;

    .hero-content {
      h1 {
        font-size: 40px;
        margin-bottom: 14px;
      }

      p {
        font-size: 18px;
        margin-bottom: 20px;
      }

      .tag-stats {
        gap: 24px;

        .stat-item {
          font-size: 15px;

          span[class^="i-"] {
            font-size: 18px;
          }
        }
      }
    }
  }

  .page-container {
    padding: 0 24px 60px;
  }

  .tags-grid {
    grid-template-columns: repeat(auto-fill, minmax(280px, 1fr));
    gap: 20px;
  }
}

@media (max-width: 768px) {
  .hero-section {
    padding: 2.5rem 1.5rem 2rem;
    min-height: 260px;

    &::after {
      background-size: 30px 30px;
    }

    .hero-content {
      h1 {
        font-size: 34px;
        margin-bottom: 12px;
        line-height: 1.2;
      }

      p {
        font-size: 17px;
        margin-bottom: 18px;
      }

      .tag-stats {
        gap: 20px;

        .stat-item {
          font-size: 14px;

          span[class^="i-"] {
            font-size: 17px;
          }
        }
      }
    }
  }

  .page-container {
    padding: 0 20px 50px;
  }

  .tags-grid {
    grid-template-columns: repeat(auto-fill, minmax(260px, 1fr));
    gap: 18px;
  }

  .tag-card {
    &:hover {
      transform: translateY(-5px);
    }

    .tag-icon {
      font-size: 42px;
    }

    .tag-content {
      h3 {
        font-size: 17px;
      }

      .tag-description {
        font-size: 13px;
      }
    }
  }
}

@media (max-width: 640px) {
  .hero-section {
    padding: 2rem 1rem;
    min-height: 240px;
    margin-bottom: 32px;

    &::after {
      background-size: 25px 25px;
    }

    .hero-content {
      h1 {
        font-size: 28px;
        margin-bottom: 10px;
      }

      p {
        font-size: 16px;
        margin-bottom: 16px;
      }
    }
  }

  .page-container {
    padding: 0 16px 40px;
  }

  .tags-grid {
    grid-template-columns: 1fr;
    gap: 16px;
  }

  .tag-card {
    .tag-icon {
      font-size: 40px;
    }

    .tag-content {
      h3 {
        font-size: 16px;
      }

      .tag-description {
        font-size: 13px;
      }
    }
  }
}

@media (max-width: 480px) {
  .hero-section {
    padding: 1.75rem 1rem 1.5rem;
    min-height: 220px;
    margin-bottom: 28px;

    &::before {
      animation: gradientShift 20s ease-in-out infinite;
    }

    &::after {
      background-size: 20px 20px;
      animation: gridMove 30s linear infinite;
    }

    .hero-content {
      h1 {
        font-size: 24px;
        margin-bottom: 8px;
        line-height: 1.3;

        text-shadow: 0 0 30px rgba(255, 255, 255, 0.8),
        0 0 60px rgba(var(--color-primary-purple-rgb), 0.6),
        0 4px 16px rgba(0, 0, 0, 0.4);

        filter: drop-shadow(0 0 20px rgba(var(--color-primary-purple-rgb), 0.5)) drop-shadow(0 3px 10px rgba(0, 0, 0, 0.3));
      }

      p {
        font-size: 14px;
        margin-bottom: 14px;
      }

      .tag-stats {
        gap: 16px;

        .stat-item {
          font-size: 12px;
          gap: 6px;

          span[class^="i-"] {
            font-size: 15px;
          }
        }
      }
    }
  }

  .page-container {
    padding: 0 12px 32px;
  }

  .tags-grid {
    grid-template-columns: 1fr;
    gap: 14px;
  }

  .tag-card {
    padding: 20px;

    .tag-icon {
      font-size: 36px;
      margin-bottom: 12px;
    }

    .tag-content {
      h3 {
        font-size: 15px;
        margin-bottom: 8px;
      }

      .tag-description {
        font-size: 12px;
        margin-bottom: 12px;
      }

      .tag-meta {
        font-size: 11px;
        gap: 5px;
      }
    }
  }
}

// Dark Mode Enhancements
html.dark {
  .tag-card {
    &:hover {
      box-shadow: 0 12px 32px rgba(0, 0, 0, 0.4);
    }
  }
}
</style>
