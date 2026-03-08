<script setup lang="ts">
import {definePage} from 'unplugin-vue-router/runtime'
import {computed, nextTick, onBeforeUnmount, onMounted, ref, watch} from 'vue'
import {useRoute, useRouter} from 'vue-router'
import {useMessage} from 'naive-ui'

import {usePostStore} from '@/stores/modules/app'
import {$t} from '@/locales'
import {useLanguageChangeEffect} from '@/hooks/useLanguageChangeEffect';

import {ContentViewer} from '@/components/ContentViewer'
import CommentSection from '@/components/CommentSection'

import type {contentservicev1_Post} from "@/api/generated/app/service/v1"
import {formatDate} from "@/utils/date";

// --- 常量定义 ---
const SCROLL_THRESHOLD = 500
const HEADING_OFFSET = 150
const THROTTLE_DELAY = 200

// --- 类型定义 ---
interface TocItem {
  id: string
  level: number
  text: string
  element: HTMLElement
}

// --- 路由定义 ---
definePage({
  name: 'post-detail',
  meta: {
    title: 'Post Detail',
    level: 3,
  },
})

// --- 状态初始化 ---
const route = useRoute()
const router = useRouter()
const postStore = usePostStore()
const message = useMessage()

// 修复类型定义，允许 null
const post = ref<contentservicev1_Post | null>(null)
const relatedPosts = ref<contentservicev1_Post[]>([])
const tableOfContents = ref<TocItem[]>([])
const activeHeading = ref<string>('')

const loading = ref(false)
const showBackToTop = ref(false)
const isLiked = ref(false)
const isBookmarked = ref(false)
const isTocExpanded = ref(true)

// --- 计算属性 ---
const postId = computed(() => {
  const id = route.params.id
  return id ? parseInt(id as string) : null
})

const currentTranslation = computed(() => {
  if (!post.value) return null;
  return postStore.getTranslation(post.value);
})

const displayTitle = computed(() => currentTranslation.value?.title || $t('page.post_detail.untitled'))
const displayContent = computed(() => currentTranslation.value?.content || '')
const displayThumbnail = computed(() => currentTranslation.value?.thumbnail || '/placeholder.jpg')

// --- 工具函数 ---
// 移到这里，避免 hoisting 依赖
function throttle(fn: Function, delay: number) {
  let timer: ReturnType<typeof setTimeout> | null = null
  return function (this: any, ...args: any[]) {
    if (!timer) {
      timer = setTimeout(() => {
        fn.apply(this, args)
        timer = null
      }, delay)
    }
  }
}

// --- 数据加载 ---
async function loadPost() {
  if (!postId.value) return
  loading.value = true
  try {
    post.value = await postStore.getPost(postId.value)
    // SEO 优化
    document.title = `${displayTitle.value} - ${import.meta.env.VITE_APP_TITLE || 'GoWind CMS'}`
    // 可以在 onBeforeUnmount 中恢复原标题，但通常单页应用不需要
  } catch (error) {
    console.error('Load post failed:', error)
    message.error($t('page.post_detail.load_failed'))
  } finally {
    loading.value = false
  }
}

async function loadRelatedPosts() {
  if (!post.value) return
  try {
    const categoryIds = post.value.categoryIds || []
    if (categoryIds.length > 0) {
      const res = await postStore.listPost(
        {page: 1, pageSize: 3},
        {
          status: 'POST_STATUS_PUBLISHED',
          categoryIds: categoryIds
        }
      )
      relatedPosts.value = (res.items || []).filter((p) => p.id !== postId.value)
    }
  } catch (error) {
    console.error('Load related posts failed:', error)
  }
}

// 加载所有数据
async function loadAllData() {
  await Promise.all([
    loadPost(),
    loadRelatedPosts(),
  ]);
}

// 监听语言切换，自动重新加载数据
useLanguageChangeEffect(loadAllData, {
  immediate: false,    // 已经在 onMounted 中加载，不需要立即执行
  autoCleanup: true,   // 组件卸载时自动取消订阅
});

// --- 交互逻辑 ---
function handleViewRelatedPost(id: number) {
  router.push(`/post/${id}`)
  window.scrollTo({top: 0, behavior: 'smooth'})
}

function handleBack() {
  const from = route.query.from as string
  const categoryId = route.query.categoryId as string
  if (from === 'category' && categoryId) {
    router.push(`/category/${categoryId}`)
  } else if (from === 'user') {
    router.push('/user')
  } else if (from === 'post') {
    router.push('/post')
  } else if (from === 'home') {
    router.push('/')
  } else if (window.history.length > 2) {
    router.back()
  } else {
    router.push('/post')
  }
}

function handleLike() {
  isLiked.value = !isLiked.value
  if (isLiked.value && post.value) {
    post.value.likes = (post.value.likes || 0) + 1
    message.success($t('page.post_detail.liked'))
  } else if (post.value && post.value.likes > 0) {
    post.value.likes -= 1
  }
}

function handleBookmark() {
  isBookmarked.value = !isBookmarked.value
  message.success(isBookmarked.value ? $t('page.post_detail.bookmarked') : $t('page.post_detail.unbookmarked'))
}

function handleShare() {
  const url = window.location.href
  if (navigator.share) {
    navigator.share({title: displayTitle.value, url}).then(() => {
      message.success($t('page.post_detail.shared'))
    }).catch(() => {
      copyToClipboard(url)
    })
  } else {
    copyToClipboard(url)
  }
}

function copyToClipboard(text: string) {
  navigator.clipboard.writeText(text).then(() => {
    message.success($t('page.post_detail.link_copied'))
  }).catch(() => {
    message.error($t('page.post_detail.copy_failed'))
  })
}

function scrollToTop() {
  window.scrollTo({top: 0, behavior: 'smooth'})
}

// --- 目录与滚动 ---
function generateTableOfContents() {
  nextTick(() => {
    // 在 .content-viewer 中查找所有标题 (ContentViewer 的根元素)
    const contentEl = document.querySelector('.content-viewer')
    if (!contentEl) return

    const headings = contentEl.querySelectorAll('h2, h3')
    const toc: TocItem[] = []

    headings.forEach((heading, index) => {
      const level = heading.tagName === 'H2' ? 2 : 3
      const id = `heading-${index}`
      if (!heading.id) heading.id = id
      heading.setAttribute('id', id)
      toc.push({
        id,
        level,
        text: heading.textContent || '',
        element: heading as HTMLElement // 缓存元素引用
      })
    })
    tableOfContents.value = toc

    // 调试输出
    console.log('Table of contents generated:', toc.length, 'items')
  })
}

function scrollToHeading(id: string) {
  const tocItem = tableOfContents.value.find(item => item.id === id)
  const element = tocItem?.element || document.getElementById(id)

  if (element) {
    const offset = 100 // 根据实际导航栏高度调整
    const elementPosition = element.getBoundingClientRect().top + window.pageYOffset
    window.scrollTo({
      top: elementPosition - offset,
      behavior: 'smooth'
    })
    activeHeading.value = id

    // 保留现有的 history state，只更新 hash
    if (history.pushState) {
      const currentState = history.state || {}
      history.pushState(currentState, '', `#${id}`)
    } else {
      window.location.hash = `#${id}`
    }
  }
}

function handleScroll() {
  showBackToTop.value = window.scrollY > SCROLL_THRESHOLD

  if (tableOfContents.value.length > 0) {
    let currentActive = activeHeading.value
    // 优化：直接使用缓存的 element，避免重复查询 DOM
    for (const heading of tableOfContents.value) {
      if (heading.element && heading.element.getBoundingClientRect().top < HEADING_OFFSET) {
        currentActive = heading.id
      }
    }
    activeHeading.value = currentActive
  }
}

const throttledScroll = throttle(handleScroll, THROTTLE_DELAY)

// --- 生命周期 ---
onMounted(async () => {
  await loadPost()

  // 等待内容完全加载后生成目录
  // 增加延迟确保 ContentViewer 已经渲染完成
  setTimeout(() => {
    generateTableOfContents()
    // 页面加载时检查 URL hash，自动滚动到对应位置
    if (window.location.hash) {
      const id = window.location.hash.slice(1)
      setTimeout(() => {
        scrollToHeading(id)
      }, 300)
    }
  }, 500)

  // 加载相关文章 (确保 post 已加载)
  await loadRelatedPosts()

  window.addEventListener('scroll', throttledScroll)
})

onBeforeUnmount(() => {
  window.removeEventListener('scroll', throttledScroll)
})

watch(() => displayContent.value, () => {
  generateTableOfContents()
})
</script>

<template>
  <div class="post-detail-page">
    <!-- Loading State - Skeleton Screen -->
    <div v-if="loading" class="post-loading">
      <div class="back-navigation">
        <n-skeleton :width="100" size="medium" />
      </div>
      <article class="post-article">
        <div class="post-banner">
          <n-skeleton height="300px" />
        </div>
        <div class="post-wrapper">
          <aside class="toc-sidebar">
            <div class="toc-container">
              <n-skeleton :width="200" size="large" style="margin-bottom: 16px;" />
              <n-skeleton :width="180" size="medium" style="margin-bottom: 8px;" />
              <n-skeleton :width="160" size="medium" style="margin-bottom: 8px;" />
              <n-skeleton :width="140" size="medium" style="margin-bottom: 8px;" />
              <n-skeleton :width="120" size="medium" />
            </div>
          </aside>
          <div class="article-content">
            <header class="post-header">
              <n-skeleton :width="'80%'" size="huge" style="margin-bottom: 16px;" />
              <n-skeleton :width="'60%'" size="large" style="margin-bottom: 24px;" />
              <div class="post-meta">
                <n-skeleton :width="100" size="medium" />
                <n-skeleton :width="100" size="medium" />
                <n-skeleton :width="100" size="medium" />
                <n-skeleton :width="100" size="medium" />
              </div>
            </header>
            <div class="post-content">
              <n-skeleton text :repeat="8" />
            </div>
          </div>
        </div>
      </article>
    </div>

    <!-- Loaded State -->
    <div v-else-if="post" class="post-container">
        <!-- Back Button -->
        <div class="back-navigation">
          <n-button text @click="handleBack()" aria-label="Back">
            <template #icon>
              <span class="i-carbon:arrow-left"/>
            </template>
            {{ $t('page.post_detail.back') }}
          </n-button>
        </div>

        <!-- Post Article -->
        <article class="post-article">
          <!-- Post Thumbnail Banner -->
          <div v-if="displayThumbnail" class="post-banner">
            <img
              :src="displayThumbnail"
              :alt="displayTitle"
              loading="lazy"
            />
            <div class="banner-overlay"/>
          </div>

          <div class="post-wrapper" :class="{ 'toc-collapsed': !isTocExpanded }">
            <!-- Table of Contents Sidebar -->
            <aside v-if="tableOfContents.length > 0 && isTocExpanded" class="toc-sidebar">
              <div class="toc-container">
                <div class="toc-header">
                  <h3 class="toc-title">
                    <span class="i-carbon:list"/>
                    {{ $t('page.post_detail.table_of_contents') }}
                  </h3>
                  <n-button
                    text
                    size="small"
                    @click="isTocExpanded = false"
                    class="toc-collapse-btn"
                    aria-label="Collapse TOC"
                  >
                    <span class="i-carbon:chevron-left"/>
                  </n-button>
                </div>
                <nav class="toc-list">
                  <a
                    v-for="item in tableOfContents"
                    :key="item.id"
                    href="javascript:void(0)"
                    :class="['toc-item', `level-${item.level}`, { active: activeHeading === item.id }]"
                    @click.prevent="scrollToHeading(item.id)"
                  >
                    {{ item.text }}
                  </a>
                </nav>
              </div>
            </aside>

            <div class="article-content">
              <!-- TOC Expand Trigger (when collapsed) -->
              <div v-if="tableOfContents.length > 0 && !isTocExpanded" class="toc-expand-trigger">
                <n-button
                  text
                  @click="isTocExpanded = true"
                  aria-label="Expand TOC"
                >
                  <template #icon>
                    <span class="i-carbon:list"/>
                  </template>
                  {{ $t('page.post_detail.table_of_contents') }}
                  <span class="i-carbon:chevron-right"/>
                </n-button>
              </div>

              <!-- Post Header -->
              <header class="post-header">
                <h1 class="post-title">{{ displayTitle }}</h1>
                <div class="post-meta">
                  <div class="meta-item">
                    <span class="i-carbon:user-avatar"/>
                    <span>{{ post.authorName }}</span>
                  </div>
                  <div class="meta-item">
                    <span class="i-carbon:calendar"/>
                    <span>{{ formatDate(post.createdAt) }}</span>
                  </div>
                  <div class="meta-item">
                    <span class="i-carbon:view"/>
                    <span>{{ post.visits || 0 }}</span>
                  </div>
                  <div class="meta-item">
                    <span class="i-carbon:thumbs-up"/>
                    <span>{{ post.likes || 0 }}</span>
                  </div>
                </div>
              </header>

              <!-- Post Content -->
              <div class="post-content">
                <ContentViewer :content="displayContent" type="markdown"/>
              </div>

              <!-- Post Actions -->
              <footer class="post-actions">
                <n-space size="large">
                  <n-button
                    size="large"
                    circle
                    :type="isLiked ? 'primary' : 'default'"
                    @click="handleLike"
                    :aria-label="$t('page.post_detail.likes')"
                  >
                    <template #icon>
                      <span :class="isLiked ? 'i-carbon:thumbs-up-filled' : 'i-carbon:thumbs-up'"/>
                    </template>
                  </n-button>
                  <n-button
                    size="large"
                    circle
                    :type="isBookmarked ? 'primary' : 'default'"
                    @click="handleBookmark"
                    :aria-label="$t('page.post_detail.bookmark')"
                  >
                    <template #icon>
                      <span
                        :class="isBookmarked ? 'i-carbon:bookmark-filled' : 'i-carbon:bookmark'"/>
                    </template>
                  </n-button>
                  <n-button
                    size="large"
                    circle
                    @click="handleShare"
                    :aria-label="$t('page.post_detail.share')"
                  >
                    <template #icon>
                      <span class="i-carbon:share"/>
                    </template>
                  </n-button>
                </n-space>
              </footer>
            </div>
          </div>
        </article>

        <!-- Comments Section -->
        <CommentSection
          :object-id="postId"
          content-type="CONTENT_TYPE_POST"
          @update:comments="() => {}"
        />

        <!-- Related Posts -->
        <section v-if="relatedPosts.length > 0" class="related-section">
          <div class="section-header">
            <h2>
              <span class="i-carbon:book"/>
              {{ $t('page.post_detail.related_posts') }}
            </h2>
          </div>
          <div class="related-grid">
            <div
              v-for="relatedPost in relatedPosts"
              :key="relatedPost.id"
              class="related-card"
              @click="handleViewRelatedPost(relatedPost.id)"
              role="button"
              tabindex="0"
              @keyup.enter="handleViewRelatedPost(relatedPost.id)"
            >
              <div class="related-image">
                <img
                  :src="relatedPost.translations?.[0]?.thumbnail || '/placeholder.jpg'"
                  :alt="relatedPost.translations?.[0]?.title"
                  loading="lazy"
                />
                <div class="image-overlay"/>
              </div>
              <div class="related-content">
                <h3>{{ relatedPost.translations?.[0]?.title }}</h3>
                <p>{{ relatedPost.translations?.[0]?.summary }}</p>
                <div class="related-meta">
                  <span><span class="i-carbon:view"/> {{ relatedPost.visits || 0 }}</span>
                  <span><span class="i-carbon:thumbs-up"/> {{ relatedPost.likes || 0 }}</span>
                </div>
              </div>
            </div>
          </div>
        </section>
      </div>
      <n-empty v-else :description="$t('page.post_detail.post_not_found')"/>

    <!-- Back to Top Button -->
    <transition name="fade">
      <n-button
        v-show="showBackToTop"
        class="back-to-top"
        circle
        size="large"
        type="primary"
        @click="scrollToTop"
        aria-label="Back to top"
      >
        <template #icon>
          <span class="i-carbon:arrow-up"/>
        </template>
      </n-button>
    </transition>
  </div>
</template>

<style scoped lang="less">
.post-detail-page {
  min-height: 100vh;
  background: var(--color-bg);
  padding-bottom: 80px;
}

// Loading Skeleton Styles
.post-loading {
  .post-article {
    max-width: 1200px;
    margin: 0 auto 40px;
    background: var(--color-surface);
    border-radius: 16px;
    overflow: hidden;
    box-shadow: 0 2px 8px rgba(0, 0, 0, 0.04);
  }

  .post-banner {
    position: relative;
    width: 100%;
    height: 300px;
    overflow: hidden;
  }

  .post-wrapper {
    display: flex;
    gap: 24px;
    padding: 40px;

    .toc-sidebar {
      flex: 0 0 200px;
      padding: 20px;
      background: var(--color-surface);
      border-radius: 12px;
    }

    .article-content {
      flex: 1;
      max-width: calc(100% - 224px);
    }
  }

  .post-header {
    margin-bottom: 32px;

    .post-meta {
      display: flex;
      gap: 16px;
      margin-top: 16px;
    }
  }

  .post-content {
    :deep(.n-skeleton) {
      margin-bottom: 12px;

      &:last-child {
        margin-bottom: 0;
      }
    }
  }
}

// Back Navigation
.back-navigation {
  max-width: 1200px;
  margin: 0 auto;
  padding: 20px 32px;

  :deep(.n-button) {
    color: var(--color-text-secondary);
    font-weight: 500;
    transition: all 0.3s;

    &:hover {
      color: var(--color-brand);
      transform: translateX(-4px);
    }
  }
}

// Post Article
.post-article {
  max-width: 1200px;
  margin: 0 auto 40px;
  background: var(--color-surface);
  border-radius: 16px;
  overflow: hidden;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.04);
}

// Post Wrapper (包含 TOC 和内容)
.post-wrapper {
  display: flex;
  gap: 24px;
  background: var(--color-surface);
  position: relative;
  transition: all 0.3s ease;

  &.toc-collapsed {
    .toc-sidebar {
      width: 0 !important;
      min-width: 0 !important;
      max-width: 0 !important;
      padding: 0;
      margin: 0;
      overflow: hidden;
      opacity: 0;
      border: none;
    }

    &::after {
      opacity: 0;
    }

    .article-content {
      flex: 1;
      max-width: 100%;
    }
  }

  &::after {
    content: '';
    position: absolute;
    left: 200px;
    top: 0;
    bottom: 0;
    width: 1px;
    background: var(--color-border);
    transition: opacity 0.3s ease;
  }
}

// TOC Expand Trigger
.toc-expand-trigger {
  position: absolute;
  top: 0;
  left: 0;
  z-index: 10;
  padding: 12px 16px;
  background: var(--color-bg);
  border-radius: 8px;
  border: 1px solid var(--color-border);
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  transition: all 0.3s;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.04);

  &:hover {
    border-color: var(--color-brand);
    background: rgba(168, 85, 247, 0.05);
    transform: translateX(2px);
  }

  :deep(.n-button) {
    color: var(--color-text-secondary);
    font-weight: 500;
    transition: all 0.3s;
    background: transparent;

    &:hover {
      color: var(--color-brand);
      background: transparent;
    }

    span[class^="i-"] {
      font-size: 16px;
      margin-left: 8px;
    }
  }
}

// Table of Contents Sidebar
.toc-sidebar {
  flex: 0 0 200px;
  width: 200px;
  padding: 40px 16px;
  background:
    var(--color-surface),
    linear-gradient(135deg, rgba(168, 85, 247, 0.06), rgba(102, 126, 234, 0.04)); // 渐变紫色调，增强层次
  backdrop-filter: blur(10px);
  box-shadow:
    2px 0 16px rgba(0, 0, 0, 0.1),
    -2px 0 12px rgba(0, 0, 0, 0.06),
    inset 0 0 0 1px rgba(168, 85, 247, 0.12); // 增强边框和阴影，提升层次感
  max-height: calc(100vh - 120px);
  position: sticky;
  top: 60px; // 距离顶部 60px 开始固定
  overflow-y: auto;
  transition: all 0.3s ease;
  border-radius: 0 12px 12px 0;
  align-self: flex-start; // 确保从顶部开始

  // 如果内容较短，目录框高度根据内容自适应
  &:not(:has(.toc-list:only-child)) {
    min-height: auto;
  }

  .toc-container {
    .toc-header {
      display: flex;
      justify-content: space-between;
      align-items: center;
      margin-bottom: 20px;
      padding-bottom: 12px;
      border-bottom: 1px solid rgba(168, 85, 247, 0.15);

      .toc-title {
        display: flex;
        align-items: center;
        gap: 8px;
        font-size: 15px;
        font-weight: 600;
        margin: 0;
        color: var(--color-text-primary);

        span[class^="i-"] {
          font-size: 17px;
          color: var(--color-brand);
        }
      }

      .toc-collapse-btn {
        padding: 4px;
        color: var(--color-text-secondary);
        transition: all 0.3s;

        &:hover {
          color: var(--color-brand);
          transform: translateX(-2px);
        }
      }
    }

    .toc-list {
      display: flex;
      flex-direction: column;
      gap: 6px;

      .toc-item {
        padding: 10px 12px;
        border-radius: 8px;
        color: var(--color-text-secondary);
        font-size: 13px;
        text-decoration: none;
        transition: all 0.25s cubic-bezier(0.4, 0, 0.2, 1);
        border-left: 2px solid transparent;
        overflow: hidden;
        text-overflow: ellipsis;
        white-space: nowrap;
        display: block;
        position: relative;

        &.level-2 {
          padding-left: 12px;
          font-weight: 500;
        }

        &.level-3 {
          padding-left: 26px;
          font-size: 12px;
          color: var(--color-text-secondary);
        }

        &:hover {
          background: rgba(168, 85, 247, 0.12); // 强化 hover 背景色
          color: var(--color-brand);
          border-left-color: var(--color-brand);
          transform: translateX(3px); // 添加位移效果
          box-shadow: 0 2px 8px rgba(168, 85, 247, 0.15); // 添加阴影效果
        }

        &.active {
          background: linear-gradient(90deg, rgba(168, 85, 247, 0.15), rgba(168, 85, 247, 0.08)); // 渐变背景
          color: var(--color-brand);
          border-left-color: var(--color-brand);
          font-weight: 600;
          box-shadow: 0 2px 8px rgba(168, 85, 247, 0.2); // 激活状态也有阴影
        }
      }
    }
  }

  // 美化滚动条
  &::-webkit-scrollbar {
    width: 4px;
  }

  &::-webkit-scrollbar-track {
    background: transparent;
  }

  &::-webkit-scrollbar-thumb {
    background: rgba(168, 85, 247, 0.3);
    border-radius: 2px;

    &:hover {
      background: rgba(168, 85, 247, 0.6);
    }
  }
}

// Post Banner
.post-banner {
  position: relative;
  width: 100%;
  padding-top: 56.25%; // 16:9 比例（优化过的）
  overflow: hidden;
  background: linear-gradient(135deg,
    var(--color-bg) 0%,
    var(--color-surface) 100%); // 使用 CSS 变量支持明暗主题

  img {
    position: absolute;
    top: 0;
    left: 0;
    width: 100%;
    height: 100%;
    object-fit: cover;
    display: block;
  }

  // ✅ 改进：多层过渡遮罩
  .banner-overlay {
    position: absolute;
    bottom: 0;
    left: 0;
    right: 0;
    height: 280px;
    background: linear-gradient(to bottom,
    transparent 0%,
    rgba(0, 0, 0, 0.15) 20%,
    rgba(0, 0, 0, 0.35) 50%,
    rgba(0, 0, 0, 0.6) 80%,
    rgba(0, 0, 0, 0.75) 100%);
    pointer-events: none;
    backdrop-filter: blur(0.5px);
  }
}

// Article Content
.article-content {
  padding: 48px 56px 40px;
  background: var(--color-surface);
  position: relative;
  flex: 1;
  transition: all 0.3s ease;

  .post-wrapper.toc-collapsed & {
    max-width: 100%;
  }

  &::before {
    content: '';
    position: absolute;
    top: 0;
    left: 0;
    right: 0;
    height: 12px;
    background: linear-gradient(to bottom,
      rgba(0, 0, 0, 0.08),
      rgba(0, 0, 0, 0.04) 50%,
      transparent);
    pointer-events: none;
  }
}

// Post Header
.post-header {
  margin-bottom: 36px;

  .post-title {
    font-size: 42px;
    font-weight: 800;
    margin: 0 0 24px 0;
    color: var(--color-text-primary);
    line-height: 1.3;
    letter-spacing: -0.5px;
  }

  .post-meta {
    display: flex;
    flex-wrap: wrap;
    gap: 24px;

    .meta-item {
      display: flex;
      align-items: center;
      gap: 8px;
      color: var(--color-text-secondary);
      font-size: 15px;
      font-weight: 500;

      span[class^="i-"] {
        font-size: 18px;
        opacity: 0.8;
      }
    }
  }
}

// Post Content
.post-content {
  font-size: 17px;
  line-height: 2.05; // 黄金行高比例，提升长文阅读舒适度
  letter-spacing: 0.3px; // 微调字间距，增强可读性
  color: var(--color-text-primary);
  margin-bottom: 40px;
  text-rendering: optimizeLegibility;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;

  :deep(h2) {
    font-size: 26px;
    font-weight: 700;
    margin: 56px 0 28px; // 增加标题上下间距，章节更清晰
    padding-bottom: 12px; // 增加下划线与文字距离
    border-bottom: 2px solid var(--color-brand);
    color: var(--color-text-primary);
    position: relative;
    line-height: 1.4; // 稍微增加标题行高
    letter-spacing: -0.3px;

    &::after {
      content: '';
      position: absolute;
      bottom: -2px;
      left: 0;
      width: 45px;
      height: 2px;
      background: linear-gradient(90deg, #a855f7, transparent);
    }
  }

  :deep(h3) {
    font-size: 21px;
    font-weight: 600;
    margin: 44px 0 24px; // 增加三级标题间距
    color: var(--color-text-primary);
    padding-left: 10px;
    border-left: 3px solid var(--color-brand);
    line-height: 1.45; // 优化标题行高
  }

  :deep(p) {
    margin: 26px 0; // 增加段间距到 26px，呼吸感更强
    text-align: justify;
    line-height: 2.15; // 优化行高到 2.15
    letter-spacing: 0.32px; // 微调字间距
    word-spacing: 0.08em; // 增加词间距，提升可读性

    & + p {
      margin-top: 24px; // 相邻段落间距同步增加
    }

    &:first-of-type::first-letter {
      font-size: 1.6em;
      font-weight: 600;
      color: var(--color-brand);
      margin-right: 6px;
      float: left;
      line-height: 1;
    }
  }

  :deep(img) {
    max-width: 100%;
    border-radius: 10px;
    margin: 36px 0; // 增加图片上下间距
    box-shadow: 0 4px 16px rgba(0, 0, 0, 0.08);
    transition: all 0.3s;
    cursor: zoom-in;
    display: block;

    &:hover {
      box-shadow: 0 8px 24px rgba(0, 0, 0, 0.12);
      transform: translateY(-2px);
    }
  }

  :deep(code) {
    background: rgba(150, 150, 150, 0.15); // 加深背景
    padding: 4px 10px;
    border-radius: 6px;
    font-family: 'Fira Code', 'Courier New', monospace;
    font-size: 0.88em;
    color: var(--color-brand);
    border: 1px solid rgba(168, 85, 247, 0.2); // 加强边框
    letter-spacing: 0.1px;

    // 行内代码不换行，但代码块中的代码可以换行
    &:not(pre code) {
      white-space: nowrap;
    }
  }

  :deep(pre) {
    background: #282c34; // One Dark 主题色
    padding: 24px; // 增加内边距
    border-radius: 12px;
    overflow-x: auto;
    border: 1px solid rgba(255, 255, 255, 0.12); // 加强边框
    margin: 40px 0; // 增加代码块整体间距
    line-height: 1.65; // 优化代码行高
    letter-spacing: 0.18px; // 微调字间距
    box-shadow: 0 6px 20px rgba(0, 0, 0, 0.18), 0 2px 8px rgba(0, 0, 0, 0.08); // 增强阴影
    position: relative;
    white-space: pre;
    word-wrap: normal;
    overflow-wrap: normal;

    &::before {
      content: attr(data-lang);
      position: absolute;
      top: 10px;
      right: 14px;
      font-size: 12px;
      color: rgba(255, 255, 255, 0.45);
      text-transform: uppercase;
      letter-spacing: 0.5px;
      font-weight: 600;
    }

    code {
      background: none;
      padding: 0;
      color: inherit;
      border: none;
      font-size: 13.5px;
      font-family: inherit;
      line-height: inherit;
      white-space: inherit;
      word-wrap: inherit;
      overflow-wrap: inherit;
    }
  }

  :deep(blockquote) {
    border-left: 5px solid var(--color-brand); // 加粗左边框
    padding: 26px 30px; // 进一步增加内边距
    margin: 40px 0; // 增加外边距
    background: linear-gradient(135deg, rgba(168, 85, 247, 0.08), rgba(168, 85, 247, 0.05)); // 渐变背景
    border-radius: 0 10px 10px 0;
    color: var(--color-text-secondary);
    font-style: italic;
    position: relative;
    letter-spacing: 0.25px; // 微调字间距
    line-height: 2.0; // 增加引用块行高
    box-shadow: 0 2px 8px rgba(168, 85, 247, 0.1); // 添加阴影

    &::before {
      content: '"';
      position: absolute;
      top: -8px;
      left: 14px;
      font-size: 64px; // 更大的引号
      color: var(--color-brand);
      opacity: 0.18;
      font-family: Georgia, serif;
      line-height: 1;
    }

    p {
      margin: 18px 0; // 增加引用段落间距
      text-indent: 0;
      text-align: left;
      position: relative;
      z-index: 1;
    }

    &:last-child {
      margin-bottom: 40px; // 增加底部间距
    }
  }

  :deep(ul), :deep(ol) {
    padding-left: 32px;
    margin: 32px 0; // 增加列表整体间距

    li {
      margin: 14px 0; // 增加列表项间距
      line-height: 2.0; // 优化列表项行高
      letter-spacing: 0.25px; // 微调字间距
      position: relative;

      &::marker {
        color: var(--color-brand);
        font-weight: 600;
      }
    }
  }

  :deep(a) {
    color: var(--color-brand);
    text-decoration: none;
    border-bottom: 1.5px solid rgba(168, 85, 247, 0.3);
    transition: all 0.25s;
    font-weight: 500;
    letter-spacing: 0.1px;

    &:hover {
      border-bottom-color: var(--color-brand);
      background: rgba(168, 85, 247, 0.08);
      padding: 2px 4px;
      margin: -2px -4px;
      border-radius: 4px;
    }
  }

  :deep(strong) {
    color: var(--color-brand);
    font-weight: 700;
    letter-spacing: 0.2px;
  }

  :deep(em) {
    color: var(--color-text-secondary);
    font-style: italic;
    letter-spacing: 0.15px;
  }

  :deep(hr) {
    border: none;
    height: 2px;
    background: linear-gradient(to right, transparent, var(--color-brand), transparent);
    margin: 48px 0;
  }

  :deep(table) {
    width: 100%;
    border-collapse: collapse;
    margin: 32px 0;
    border-radius: 8px;
    overflow: hidden;
    box-shadow: 0 2px 12px rgba(0, 0, 0, 0.08);
    font-size: 15px;
    letter-spacing: 0.15px;

    th {
      background: var(--color-brand);
      color: white;
      padding: 14px 16px;
      text-align: left;
      font-weight: 600;
      line-height: 1.5;
      letter-spacing: 0.3px;
    }

    td {
      padding: 12px 16px;
      border-bottom: 1px solid var(--color-border);
      line-height: 1.7;
    }

    tr:hover {
      background: rgba(168, 85, 247, 0.05);
    }
  }
}

// Back to Top Button
.back-to-top {
  position: fixed;
  bottom: 40px;
  right: 40px;
  z-index: 999;
  box-shadow: 0 8px 24px rgba(168, 85, 247, 0.3);

  &:hover {
    transform: translateY(-4px);
    box-shadow: 0 12px 32px rgba(168, 85, 247, 0.4);
  }
}

// Post Actions
.post-actions {
  margin-top: 48px;
  padding-top: 32px;
  border-top: 1px solid var(--color-border);
  display: flex;
  justify-content: center;

  :deep(.n-space) {
    display: flex;
    justify-content: center;
  }

  :deep(.n-button) {
    transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);

    &:hover {
      transform: translateY(-2px);
      box-shadow: 0 4px 12px rgba(0, 0, 0, 0.08);
    }
  }
}

// Related Posts Section
.related-section {
  max-width: 1200px;
  margin: 0 auto 40px;
  background: var(--color-surface);
  border-radius: 16px;
  padding: 48px 64px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.04);

  .section-header {
    margin-bottom: 32px;

    h2 {
      display: flex;
      align-items: center;
      gap: 12px;
      font-size: 28px;
      font-weight: 700;
      margin: 0;
      color: var(--color-text-primary);

      span[class^="i-"] {
        font-size: 32px;
        color: var(--color-brand);
      }
    }
  }

  .related-grid {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(320px, 1fr));
    gap: 32px;
  }

  .related-card {
    background: var(--color-bg);
    border-radius: 12px;
    overflow: hidden;
    border: 1px solid var(--color-border);
    transition: all 0.3s;
    cursor: pointer;

    &:hover {
      border-color: var(--color-brand);
      box-shadow: 0 8px 24px rgba(0, 0, 0, 0.08);
      transform: translateY(-4px);
    }

    .related-image {
      position: relative;
      width: 100%;
      padding-top: 56.25%; // 16:9
      overflow: hidden;
      background: linear-gradient(135deg, var(--color-bg) 0%, var(--color-surface) 100%);

      img {
        position: absolute;
        top: 0;
        left: 0;
        width: 100%;
        height: 100%;
        object-fit: cover;
        transition: all 0.3s;
      }

      .image-overlay {
        position: absolute;
        bottom: 0;
        left: 0;
        right: 0;
        height: 120px;
        background: linear-gradient(to bottom,
          transparent 0%,
          rgba(0, 0, 0, 0.15) 30%,
          rgba(0, 0, 0, 0.4) 70%,
          rgba(0, 0, 0, 0.6) 100%);
        pointer-events: none;
      }
    }

    .related-content {
      padding: 24px;

      h3 {
        font-size: 18px;
        font-weight: 600;
        margin: 0 0 12px 0;
        color: var(--color-text-primary);
        line-height: 1.4;
        display: -webkit-box;
        -webkit-line-clamp: 2;
        -webkit-box-orient: vertical;
        overflow: hidden;
      }

      p {
        font-size: 14px;
        color: var(--color-text-secondary);
        margin: 0 0 16px 0;
        line-height: 1.6;
        display: -webkit-box;
        -webkit-line-clamp: 2;
        -webkit-box-orient: vertical;
        overflow: hidden;
      }

      .related-meta {
        display: flex;
        gap: 16px;
        font-size: 13px;
        color: var(--color-text-secondary);

        span {
          display: flex;
          align-items: center;
          gap: 6px;

          span[class^="i-"] {
            font-size: 16px;
          }
        }
      }
    }
  }
}

// Fade Transition
.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.3s, transform 0.3s;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
  transform: translateY(20px);
}

// Responsive Design
@media (max-width: 1024px) {
  .back-navigation,
  .post-article,
  .comments-section,
  .related-section {
    margin-left: 24px;
    margin-right: 24px;
  }

  .article-content {
    padding: 40px 40px 36px;

    &::before {
      height: 6px;
    }
  }

  .post-wrapper {
    gap: 20px;

    &::after {
      left: 180px;
    }
  }

  .toc-sidebar {
    flex: 0 0 180px;
    width: 180px;
    padding: 36px 16px;

    .toc-container {
      .toc-header {
        .toc-title {
          font-size: 14px;
        }
      }

      .toc-list {
        gap: 5px;

        .toc-item {
          padding: 7px 9px;
          font-size: 12px;
        }
      }
    }
  }

  .post-header {
    margin-bottom: 36px;

    .post-title {
      font-size: 36px;
      margin-bottom: 20px;
      letter-spacing: -0.4px;
    }

    .post-meta {
      gap: 20px;

      .meta-item {
        font-size: 14px;

        span[class^="i-"] {
          font-size: 16px;
        }
      }
    }
  }

  .post-content {
    font-size: 16px;
    line-height: 1.95;
    margin-bottom: 36px;

    :deep(h2) {
      font-size: 26px;
      margin: 44px 0 22px;
      padding-bottom: 12px;

      &::after {
        width: 45px;
      }
    }

    :deep(h3) {
      font-size: 21px;
      margin: 32px 0 18px;
      padding-left: 12px;
    }

    :deep(p) {
      margin: 18px 0;
      line-height: 2.05;
    }

    :deep(img) {
      margin: 28px 0;
    }
  }

  .post-banner {
    padding-top: 50%; // 调整为 2:1 比例
  }

  .post-wrapper {
    gap: 24px;

    &::after {
      left: 200px;
    }
  }

  .related-grid {
    grid-template-columns: repeat(auto-fit, minmax(280px, 1fr));
    gap: 20px;
  }
}

@media (max-width: 768px) {
  .back-navigation {
    padding: 16px 20px;

    :deep(.n-button) {
      font-size: 14px;
    }
  }

  .post-article,
  .related-section {
    margin-left: 20px;
    margin-right: 20px;
    border-radius: 14px;
  }

  .post-banner {
    padding-top: 56.25%; // 16:9 比例
  }

  .post-wrapper {
    gap: 0;
    flex-direction: column;

    &::after {
      display: none;
    }
  }

  // 只隐藏 TOC 侧边栏
  .toc-sidebar {
    display: none;
  }

  .toc-expand-trigger {
    display: none;
  }

  .article-content {
    padding: 40px 32px 36px;
    width: 100%;
  }

  .post-header {
    margin-bottom: 32px;

    .post-title {
      font-size: 30px;
      margin-bottom: 18px;
      letter-spacing: -0.3px;
      line-height: 1.35;
    }

    .post-meta {
      gap: 16px;

      .meta-item {
        font-size: 13px;

        span[class^="i-"] {
          font-size: 15px;
        }
      }
    }
  }

  .post-content {
    font-size: 16px;
    line-height: 1.9;
    margin-bottom: 36px;

    :deep(h2) {
      font-size: 26px;
      margin: 40px 0 20px;
      padding-bottom: 12px;
      border-bottom-width: 2px;

      &::after {
        width: 45px;
        height: 2px;
        bottom: -2px;
      }
    }

    :deep(h3) {
      font-size: 21px;
      margin: 32px 0 18px;
      padding-left: 12px;
      border-left-width: 3px;
    }

    :deep(p) {
      margin: 18px 0;
      line-height: 2;
      text-indent: 1.5em;

      & + p {
        margin-top: 16px;
      }
    }

    :deep(img) {
      margin: 32px 0;
      border-radius: 10px;
    }

    :deep(code) {
      padding: 3px 8px;
      font-size: 0.88em;
    }

    :deep(pre) {
      padding: 20px;
      margin: 28px 0;
      border-radius: 10px;
    }

    :deep(blockquote) {
      padding: 16px 20px;
      margin: 28px 0;
      border-left-width: 4px;

      &::before {
        font-size: 50px;
        top: -8px;
      }
    }

    :deep(ul), :deep(ol) {
      padding-left: 24px;
      margin: 18px 0;

      li {
        margin: 14px 0;
        line-height: 1.9;
      }
    }

    :deep(a) {
      border-bottom-width: 1px;

      &:hover {
        padding: 1px 3px;
        margin: -1px -3px;
      }
    }

    :deep(hr) {
      margin: 32px 0;
    }

    :deep(table) {
      font-size: 12px;
      margin: 20px -4px; // 让表格更宽
      border-radius: 6px;

      th {
        padding: 8px 10px;
        font-size: 12px;
      }

      td {
        padding: 8px 10px;
        font-size: 12px;
      }
    }
  }
}
</style>
