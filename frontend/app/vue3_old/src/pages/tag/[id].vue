<script setup lang="ts">
import {definePage} from 'unplugin-vue-router/runtime'
import {computed, onMounted, ref, watch} from 'vue'
import {useRoute, useRouter} from 'vue-router'
import {useMessage} from 'naive-ui'

import {$t} from '@/locales'
import {useTagStore} from '@/stores'
import type {contentservicev1_Tag,} from "@/api/generated/app/service/v1";
import {useLanguageChangeEffect} from '@/hooks/useLanguageChangeEffect';

definePage({
  name: 'tag-detail',
  meta: {
    title: 'Tag Detail',
    level: 3,
  },
})

const route = useRoute()
const router = useRouter()
const tagStore = useTagStore()
const message = useMessage()

const loading = ref(false)
const tag = ref<contentservicev1_Tag>(null)

const tagId = computed(() => {
  const id = route.params.id
  return id ? parseInt(id as string) : null
})

async function loadTag() {
  if (!tagId.value) return

  loading.value = true
  try {
    tag.value = await tagStore.getTag(tagId.value)
  } catch (error) {
    console.error('Load tag failed:', error)
    message.error($t('page.post_detail.load_failed'))
  } finally {
    loading.value = false
  }
}

function handleBack() {
  router.push('/tag')
}

onMounted(async () => {
  await loadTag()
})

// 监听语言切换，自动重新加载数据
useLanguageChangeEffect(async () => {
  await loadTag();
}, {
  immediate: false,      // 是否立即执行一次
  autoCleanup: true,    // 是否自动清理
});

// 监听路由参数变化（id 变化时重新加载）
watch(
  () => route.params.id,
  async (newId, oldId) => {
    if (newId !== oldId) {
      console.log('[Tag Detail] Route param changed, reloading...')
      await loadTag()
    }
  }
)
</script>

<template>
  <div class="tag-detail-page">
    <!-- Hero Section -->
    <div class="hero-section" :style="{
      background: tag?.color ? `linear-gradient(135deg, ${tag.color} 0%, ${tag.color}dd 50%, ${tag.color}aa 100%)` : undefined
    }">
      <div class="hero-content">
        <div class="tag-icon" :style="{color: tag?.color || '#6366f1'}">
          <span :class="`i-carbon:${tag?.icon || 'tag'}`"/>
        </div>
        <h1>{{ tagStore.getTranslation(tag)?.name || $t('page.tags.tag_untitled') }}</h1>
        <p v-if="tagStore.getTranslation(tag)?.description" class="tag-description">
          {{ tagStore.getTranslation(tag)?.description }}
        </p>
        <div class="tag-stats">
          <div class="stat-item">
            <span class="i-carbon:document"/>
            <span>{{ tag?.postCount || 0 }} {{ $t('page.posts.articles') }}</span>
          </div>
        </div>
      </div>
    </div>

    <!-- Posts Section -->
    <div class="page-container">
      <!-- Back Button -->
      <div class="back-button-container">
        <n-button
          size="small"
          @click="handleBack"
        >
          <template #icon>
            <span class="i-carbon:arrow-left"/>
          </template>
          {{ $t('page.categories.back_to_list') }}
        </n-button>
      </div>

      <!-- Posts List with Pagination -->
      <PostListWithPagination
        v-if="tagId"
        :key="tagId"
        :initial-page-size="10"
        :page-sizes="[10, 20, 30, 40]"
        :tag-id="tagId"
        from="tag"
      />
    </div>
  </div>
</template>

<style scoped lang="less">
.tag-detail-page {
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
  min-height: 350px;
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
    radial-gradient(ellipse at 80% 80%, rgba(200, 100, 255, 0.12) 0%, transparent 50%);
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

    .tag-icon {
      font-size: 64px;
      margin-bottom: 20px;
      display: inline-flex;
      align-items: center;
      justify-content: center;
      animation: fadeInUp 0.8s ease-out 0.1s both;
    }

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

    .tag-description {
      font-size: 20px;
      color: rgba(255, 255, 255, 0.95);
      font-weight: 500;
      line-height: 1.5;
      max-width: 700px;
      margin: 0 auto 24px;
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
  max-width: 1200px;
  margin: 0 auto;
  padding: 0 32px 80px;
}

// Back Button
.back-button-container {
  margin-bottom: 24px;

  :deep(.n-button) {
    border-radius: 12px;
    padding: 12px 20px;
    font-weight: 600;
    font-size: 14px;
    transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
    background: var(--color-surface);
    border: 1.5px solid var(--color-border);
    color: var(--color-text-primary);
    box-shadow: 0 2px 8px rgba(0, 0, 0, 0.06);

    span[class^="i-"] {
      font-size: 16px;
      margin-right: 8px;
      transition: transform 0.3s cubic-bezier(0.4, 0, 0.2, 1);
    }

    &:hover {
      border-color: var(--color-brand);
      color: var(--color-brand);
      background: rgba(102, 126, 234, 0.08);
      transform: translateX(-4px);
      box-shadow: 0 4px 12px rgba(102, 126, 234, 0.15);

      span[class^="i-"] {
        transform: translateX(-4px);
      }
    }

    &:active {
      transform: translateX(-2px);
      box-shadow: 0 2px 6px rgba(102, 126, 234, 0.1);
    }
  }
}

// Responsive Design
@media (max-width: 1024px) {
  .hero-section {
    padding: 3rem 1.5rem 2.5rem;
    min-height: 300px;

    .hero-content {
      .tag-icon {
        font-size: 56px;
      }

      h1 {
        font-size: 40px;
        letter-spacing: -0.8px;
        margin-bottom: 14px;
      }

      .tag-description {
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
}

@media (max-width: 768px) {
  .hero-section {
    padding: 2.5rem 1.5rem 2rem;
    min-height: 280px;

    &::after {
      background-size: 30px 30px;
    }

    .hero-content {
      .tag-icon {
        font-size: 48px;
        margin-bottom: 16px;
      }

      h1 {
        font-size: 34px;
        margin-bottom: 12px;
        line-height: 1.2;
        letter-spacing: -0.6px;
      }

      .tag-description {
        font-size: 17px;
        margin-bottom: 18px;
        line-height: 1.6;
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

  .back-button-container {
    margin-bottom: 20px;

    :deep(.n-button) {
      width: 100%;
      padding: 12px 16px;
      font-size: 14px;

      span[class^="i-"] {
        margin-right: 8px;
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
      .tag-icon {
        font-size: 42px;
        margin-bottom: 14px;
      }

      h1 {
        font-size: 28px;
        margin-bottom: 10px;
        letter-spacing: -0.4px;
      }

      .tag-description {
        font-size: 16px;
        margin-bottom: 16px;
      }

      .tag-stats {
        gap: 18px;

        .stat-item {
          font-size: 13px;

          span[class^="i-"] {
            font-size: 16px;
          }
        }
      }
    }
  }

  .page-container {
    padding: 0 16px 40px;
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
      .tag-icon {
        font-size: 36px;
        margin-bottom: 12px;
      }

      h1 {
        font-size: 24px;
        margin-bottom: 8px;
        line-height: 1.3;

        text-shadow: 0 0 30px rgba(255, 255, 255, 0.8),
        0 0 60px rgba(var(--color-primary-purple-rgb), 0.6),
        0 4px 16px rgba(0, 0, 0, 0.4);

        filter: drop-shadow(0 0 20px rgba(var(--color-primary-purple-rgb), 0.5)) drop-shadow(0 3px 10px rgba(0, 0, 0, 0.3));
      }

      .tag-description {
        font-size: 14px;
        margin-bottom: 14px;
        line-height: 1.5;
        max-width: 100%;
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

  .back-button-container {
    margin-bottom: 16px;

    :deep(.n-button) {
      padding: 12px 16px;
      font-size: 13px;
      border-radius: 12px;
      box-shadow: 0 2px 6px rgba(0, 0, 0, 0.05);

      span[class^="i-"] {
        font-size: 15px;
        margin-right: 8px;
      }
    }
  }
}

// Dark Mode Enhancements
html.dark {
  .hero-section {
    &::before {
      opacity: 0.8;
    }
  }
}
</style>
