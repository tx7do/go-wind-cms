<script setup lang="ts">
import {definePage} from 'unplugin-vue-router/runtime'
import {computed, nextTick, onBeforeUnmount, onMounted, ref, watch} from 'vue'
import {useRoute, useRouter} from 'vue-router'
import {useMessage} from 'naive-ui'
import {useCommentStore, usePostStore} from '@/stores/modules/app'
import {$t, currentLocaleLanguageCode} from '@/locales'
import {ContentViewer} from '@/components/ContentViewer'
import type {commentservicev1_Comment, contentservicev1_Post} from "@/api/generated/app/service/v1"

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

interface CommentForm {
  content: string
  authorName: string
  email: string
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
const commentStore = useCommentStore()
const message = useMessage()

// 修复类型定义，允许 null
const post = ref<contentservicev1_Post | null>(null)
const comments = ref<commentservicev1_Comment[]>([])
const relatedPosts = ref<contentservicev1_Post[]>([])
const tableOfContents = ref<TocItem[]>([])
const activeHeading = ref<string>('')
const submitting = ref(false)
const loading = ref(false)
const showBackToTop = ref(false)
const isLiked = ref(false)
const isBookmarked = ref(false)
const isTocExpanded = ref(true)

const newComment = ref<CommentForm>({
  content: '',
  authorName: '',
  email: '',
})

// --- 计算属性 ---
const postId = computed(() => {
  const id = route.params.id
  return id ? parseInt(id as string) : null
})

const currentTranslation = computed(() => {
  if (!post.value?.translations) return null
  const locale = currentLocaleLanguageCode()
  return post.value.translations.find((t) => t.languageCode === locale) || post.value.translations[0]
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

function formatDate(dateString: string) {
  if (!dateString) return ''
  return new Date(dateString).toLocaleString()
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

async function loadComments() {
  if (!postId.value) return
  try {
    const res = await commentStore.listComment(
      {page: 1, pageSize: 50},
      {postId: postId.value, status: 'COMMENT_STATUS_APPROVED'}
    )
    comments.value = res.items || []
  } catch (error) {
    console.error('Load comments failed:', error)
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

// --- 交互逻辑 ---
async function handleSubmitComment() {
  if (!newComment.value.content || !newComment.value.authorName || !newComment.value.email) {
    message.warning($t('page.post_detail.fill_form_info'))
    return
  }
  if (!postId.value || submitting.value) return

  submitting.value = true
  try {
    await commentStore.createComment({
      postId: postId.value,
      content: newComment.value.content,
      authorName: newComment.value.authorName,
      email: newComment.value.email,
      status: 'COMMENT_STATUS_PENDING',
    })
    message.success($t('page.post_detail.comment_submitted'))
    newComment.value = {content: '', authorName: '', email: ''}
    await loadComments()
    // 用户体验优化：提交后滚动到评论列表
    setTimeout(() => {
      document.querySelector('.comments-list')?.scrollIntoView({behavior: 'smooth'})
    }, 100)
  } catch (error) {
    console.error('Submit comment failed:', error)
    message.error($t('page.post_detail.submit_comment_failed'))
  } finally {
    submitting.value = false
  }
}

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
    const contentEl = document.querySelector('.post-content')
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

    if (history.pushState) {
      history.pushState(null, '', `#${id}`)
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
  await Promise.all([loadComments(), loadRelatedPosts()])

  // 等待内容完全加载后生成目录
  setTimeout(() => {
    generateTableOfContents()
    // 页面加载时检查 URL hash，自动滚动到对应位置
    if (window.location.hash) {
      const id = window.location.hash.slice(1)
      setTimeout(() => {
        scrollToHeading(id)
      }, 300)
    }
  }, 100)

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
    <n-spin :show="loading">
      <div v-if="post" class="post-container">
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
        <section v-if="!post.disallowComment" class="comments-section">
          <div class="section-header">
            <h2>
              <span class="i-carbon:chat"/>
              {{ $t('page.post_detail.comments_count', {count: comments.length}) }}
            </h2>
          </div>

          <!-- Comment Form -->
          <div class="comment-form">
            <h3>{{ $t('page.post_detail.write_comment') }}</h3>
            <n-form>
              <n-grid :cols="2" :x-gap="16">
                <n-form-item-gi :label="$t('page.post_detail.nickname')">
                  <n-input
                    v-model:value="newComment.authorName"
                    :placeholder="$t('page.post_detail.enter_nickname')"
                    size="large"
                    :disabled="submitting"
                  />
                </n-form-item-gi>
                <n-form-item-gi :label="$t('page.post_detail.email')">
                  <n-input
                    v-model:value="newComment.email"
                    :placeholder="$t('page.post_detail.enter_email')"
                    type="text"
                    size="large"
                    :disabled="submitting"
                  />
                </n-form-item-gi>
              </n-grid>
              <n-form-item :label="$t('page.post_detail.comment_content')">
                <n-input
                  v-model:value="newComment.content"
                  type="textarea"
                  :rows="5"
                  :placeholder="$t('page.post_detail.write_comment')"
                  size="large"
                  :disabled="submitting"
                />
              </n-form-item>
              <n-form-item>
                <n-button
                  type="primary"
                  size="large"
                  @click="handleSubmitComment"
                  :loading="submitting"
                >
                  <template #icon>
                    <span class="i-carbon:send-alt"/>
                  </template>
                  {{ $t('page.post_detail.submit_comment') }}
                </n-button>
              </n-form-item>
            </n-form>
          </div>

          <!-- Comments List -->
          <div v-if="comments.length > 0" class="comments-list">
            <div v-for="comment in comments" :key="comment.id" class="comment-item">
              <div class="comment-avatar">
                <n-avatar :size="48" round>
                  {{ comment.authorName?.charAt(0) || 'U' }}
                </n-avatar>
              </div>
              <div class="comment-body">
                <div class="comment-header">
                  <strong class="comment-author">{{ comment.authorName }}</strong>
                  <span class="comment-date">{{ formatDate(comment.createdAt) }}</span>
                </div>
                <div class="comment-content">
                  <ContentViewer :content="comment.content" type="text"/>
                </div>
              </div>
            </div>
          </div>
          <n-empty v-else :description="$t('page.post_detail.no_comments')"
                   style="margin-top: 40px;"/>
        </section>

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
    </n-spin>

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
  background: var(--color-surface);
  box-shadow: 2px 0 8px rgba(0, 0, 0, 0.02);
  max-height: calc(100vh - 120px);
  position: sticky;
  top: 60px;
  overflow-y: auto;
  transition: all 0.3s ease;

  .toc-container {
    .toc-header {
      display: flex;
      justify-content: space-between;
      align-items: center;
      margin-bottom: 20px;

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
        padding: 8px 10px;
        border-radius: 6px;
        color: var(--color-text-secondary);
        font-size: 13px;
        text-decoration: none;
        transition: all 0.3s;
        border-left: 2px solid transparent;
        overflow: hidden;
        text-overflow: ellipsis;
        white-space: nowrap;
        display: block;

        &.level-2 {
          padding-left: 10px;
          font-weight: 500;
        }

        &.level-3 {
          padding-left: 24px;
          font-size: 12px;
          color: var(--color-text-secondary);
        }

        &:hover {
          background: rgba(168, 85, 247, 0.08);
          color: var(--color-text-primary);
          border-left-color: var(--color-brand);
        }

        &.active {
          background: rgba(168, 85, 247, 0.12);
          color: var(--color-brand);
          border-left-color: var(--color-brand);
          font-weight: 600;
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
  background: linear-gradient(135deg, #f5f5f5 0%, #e8e8e8 100%);

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
    rgba(0, 0, 0, 0.06),
    rgba(0, 0, 0, 0.03) 50%,
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
  line-height: 2;
  letter-spacing: 0.3px;
  color: var(--color-text-primary);
  margin-bottom: 40px;
  text-rendering: optimizeLegibility;

  :deep(h2) {
    font-size: 28px;
    font-weight: 700;
    margin: 56px 0 28px;
    padding-bottom: 12px;
    border-bottom: 2px solid var(--color-brand);
    color: var(--color-text-primary);
    position: relative;
    line-height: 1.4;
    letter-spacing: -0.2px;

    &::after {
      content: '';
      position: absolute;
      bottom: -2px;
      left: 0;
      width: 50px;
      height: 2px;
      background: linear-gradient(90deg, #a855f7, transparent);
    }
  }

  :deep(h3) {
    font-size: 22px;
    font-weight: 600;
    margin: 40px 0 24px;
    color: var(--color-text-primary);
    padding-left: 12px;
    border-left: 3px solid var(--color-brand);
    line-height: 1.4;
  }

  :deep(p) {
    margin: 24px 0;
    text-align: justify;
    line-height: 2.1;
    letter-spacing: 0.3px;
    word-spacing: 0.1em;

    & + p {
      margin-top: 20px;
    }

    &:first-of-type::first-letter {
      font-size: 1.8em;
      font-weight: 700;
      color: var(--color-brand);
      margin-right: 4px;
    }
  }

  :deep(img) {
    max-width: 100%;
    border-radius: 10px;
    margin: 32px 0;
    box-shadow: 0 6px 20px rgba(0, 0, 0, 0.1);
    transition: all 0.3s;
    cursor: zoom-in;
    display: block;

    &:hover {
      box-shadow: 0 10px 28px rgba(0, 0, 0, 0.15);
      transform: translateY(-3px);
    }
  }

  :deep(code) {
    background: var(--color-bg);
    padding: 4px 10px;
    border-radius: 6px;
    font-family: 'Fira Code', 'Courier New', monospace;
    font-size: 0.9em;
    color: var(--color-brand);
    border: 1px solid rgba(168, 85, 247, 0.15);
    letter-spacing: 0.1px;
  }

  :deep(pre) {
    background: var(--color-bg);
    padding: 20px;
    border-radius: 10px;
    overflow-x: auto;
    border: 1px solid var(--color-border);
    margin: 28px 0;
    line-height: 1.65;
    letter-spacing: 0.2px;
    box-shadow: inset 0 2px 8px rgba(0, 0, 0, 0.03);

    code {
      background: none;
      padding: 0;
      color: var(--color-text-primary);
      border: none;
      font-size: 0.92em;
    }
  }

  :deep(blockquote) {
    border-left: 4px solid var(--color-brand);
    padding: 18px 20px;
    margin: 28px 0;
    background: rgba(168, 85, 247, 0.04);
    border-radius: 0 6px 6px 0;
    color: var(--color-text-secondary);
    font-style: italic;
    position: relative;
    letter-spacing: 0.2px;
    line-height: 1.9;

    &::before {
      content: '"';
      position: absolute;
      top: -8px;
      left: 8px;
      font-size: 48px;
      color: var(--color-brand);
      opacity: 0.25;
      font-family: Georgia, serif;
    }

    p {
      margin: 12px 0;
      text-indent: 0;
    }
  }

  :deep(ul), :deep(ol) {
    padding-left: 36px;
    margin: 24px 0;

    li {
      margin: 14px 0;
      line-height: 2.05;
      letter-spacing: 0.25px;

      &::marker {
        color: var(--color-brand);
        font-weight: 600;
      }
    }
  }

  :deep(a) {
    color: var(--color-brand);
    text-decoration: none;
    border-bottom: 1.5px solid rgba(168, 85, 247, 0.25);
    transition: all 0.25s;
    font-weight: 500;
    letter-spacing: 0.15px;

    &:hover {
      border-bottom-color: var(--color-brand);
      background: rgba(168, 85, 247, 0.06);
      padding: 1px 3px;
      margin: -1px -3px;
      border-radius: 3px;
    }
  }

  :deep(strong) {
    color: var(--color-brand);
    font-weight: 700;
    letter-spacing: 0.25px;
  }

  :deep(em) {
    color: var(--color-text-secondary);
    font-style: italic;
    letter-spacing: 0.15px;
  }

  :deep(hr) {
    border: none;
    height: 1.5px;
    background: linear-gradient(to right, transparent, var(--color-brand), transparent);
    margin: 40px 0;
  }

  :deep(table) {
    width: 100%;
    border-collapse: collapse;
    margin: 28px 0;
    border-radius: 6px;
    overflow: hidden;
    box-shadow: 0 1px 6px rgba(0, 0, 0, 0.06);
    font-size: 14px;
    letter-spacing: 0.2px;

    th {
      background: var(--color-brand);
      color: white;
      padding: 14px;
      text-align: left;
      font-weight: 600;
      line-height: 1.55;
    }

    td {
      padding: 12px 14px;
      border-bottom: 1px solid var(--color-border);
      line-height: 1.75;
    }

    tr:hover {
      background: rgba(168, 85, 247, 0.04);
    }
  }
}

// Post Actions
.post-actions {
  padding: 32px 0;
  border-top: 1px solid var(--color-border);
  display: flex;
  justify-content: center;

  :deep(.n-button) {
    width: 56px;
    height: 56px;
    border: 2px solid var(--color-border);
    background: var(--color-surface);
    transition: all 0.3s;

    &:hover {
      border-color: var(--color-brand);
      color: var(--color-brand);
      transform: translateY(-4px);
      box-shadow: 0 8px 16px rgba(0, 0, 0, 0.1);
    }
  }
}

// Section Header
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

// Comments Section
.comments-section {
  max-width: 1200px;
  margin: 0 auto 40px;
  background: var(--color-surface);
  border-radius: 16px;
  padding: 48px 64px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.04);
}

// Comment Form
.comment-form {
  background: var(--color-bg);
  border-radius: 12px;
  padding: 32px;
  margin-bottom: 40px;
  border: 1px solid var(--color-border);

  h3 {
    font-size: 20px;
    font-weight: 600;
    margin: 0 0 24px 0;
    color: var(--color-text-primary);
  }

  :deep(.n-form-item) {
    margin-bottom: 24px;
  }

  :deep(.n-button) {
    padding: 0 32px;
  }
}

// Comments List
.comments-list {
  display: flex;
  flex-direction: column;
  gap: 24px;
}

.comment-item {
  display: flex;
  gap: 16px;
  padding: 24px;
  background: var(--color-bg);
  border-radius: 12px;
  border: 1px solid var(--color-border);
  transition: all 0.3s;

  &:hover {
    border-color: var(--color-brand);
    box-shadow: 0 4px 12px rgba(0, 0, 0, 0.06);
  }

  .comment-avatar {
    flex-shrink: 0;

    :deep(.n-avatar) {
      background: linear-gradient(135deg, var(--color-brand) 0%, #764ba2 100%);
      color: #fff;
      font-weight: 600;
    }
  }

  .comment-body {
    flex: 1;
    min-width: 0;

    .comment-header {
      display: flex;
      justify-content: space-between;
      align-items: center;
      margin-bottom: 12px;

      .comment-author {
        font-size: 16px;
        font-weight: 600;
        color: var(--color-text-primary);
      }

      .comment-date {
        font-size: 14px;
        color: var(--color-text-secondary);
      }
    }

    .comment-content {
      font-size: 15px;
      line-height: 1.7;
      color: var(--color-text-primary);
    }
  }
}

// Related Section
.related-section {
  max-width: 1200px;
  margin: 0 auto;
  background: var(--color-surface);
  border-radius: 16px;
  padding: 48px 64px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.04);
}

.related-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
  gap: 24px;
}

.related-card {
  cursor: pointer;
  transition: all 0.3s;
  border-radius: 12px;
  overflow: hidden;
  border: 1px solid var(--color-border);
  background: var(--color-bg);

  &:hover {
    transform: translateY(-8px);
    box-shadow: 0 12px 24px rgba(0, 0, 0, 0.12);
    border-color: var(--color-brand);

    .related-image {
      .image-overlay {
        background: rgba(0, 0, 0, 0.3);
      }

      img {
        transform: scale(1.1);
      }
    }
  }

  .related-image {
    position: relative;
    width: 100%;
    height: 200px;
    overflow: hidden;

    img {
      width: 100%;
      height: 100%;
      object-fit: cover;
      transition: transform 0.5s;
    }

    .image-overlay {
      position: absolute;
      top: 0;
      left: 0;
      right: 0;
      bottom: 0;
      background: rgba(0, 0, 0, 0.1);
      transition: all 0.3s;
    }
  }

  .related-content {
    padding: 20px;

    h3 {
      font-size: 18px;
      font-weight: 600;
      margin: 0 0 12px 0;
      color: var(--color-text-primary);
      display: -webkit-box;
      -webkit-line-clamp: 2;
      -webkit-box-orient: vertical;
      overflow: hidden;
      line-height: 1.4;
    }

    p {
      font-size: 14px;
      color: var(--color-text-secondary);
      margin: 0 0 16px 0;
      display: -webkit-box;
      -webkit-line-clamp: 2;
      -webkit-box-orient: vertical;
      overflow: hidden;
      line-height: 1.6;
    }

    .related-meta {
      display: flex;
      gap: 20px;
      font-size: 13px;
      color: var(--color-text-secondary);

      span {
        display: flex;
        align-items: center;
        gap: 4px;

        span[class^="i-"] {
          font-size: 16px;
        }
      }
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

  .comments-section,
  .related-section {
    padding: 40px 48px;
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
  .comments-section,
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
      padding-left: 32px;
      margin: 24px 0;

      li {
        margin: 14px 0;
        line-height: 1.9;
      }
    }

    :deep(table) {
      margin: 28px 0;

      th {
        padding: 14px;
        font-size: 14px;
      }

      td {
        padding: 12px 14px;
        font-size: 14px;
      }
    }
  }

  .post-actions {
    padding: 28px 0;

    :deep(.n-button) {
      width: 52px;
      height: 52px;
      border-width: 1.5px;
    }

    :deep(.n-space) {
      gap: 16px !important;
    }
  }

  .comments-section,
  .related-section {
    padding: 36px 32px;
  }

  .section-header {
    margin-bottom: 28px;

    h2 {
      font-size: 24px;

      span[class^="i-"] {
        font-size: 28px;
      }
    }
  }

  .comment-form {
    padding: 24px;
    margin-bottom: 32px;

    h3 {
      font-size: 18px;
      margin-bottom: 20px;
    }

    :deep(.n-form-item) {
      margin-bottom: 20px;
    }

    :deep(.n-grid) {
      gap: 12px !important;
    }

    :deep(.n-button) {
      padding: 0 28px;
      height: 44px;
    }
  }

  .comment-item {
    padding: 20px;
    gap: 14px;

    .comment-avatar {
      :deep(.n-avatar) {
        --n-size: 44px !important;
      }
    }

    .comment-body {
      .comment-header {
        flex-direction: column;
        align-items: flex-start;
        gap: 6px;
        margin-bottom: 10px;

        .comment-author {
          font-size: 15px;
        }

        .comment-date {
          font-size: 13px;
        }
      }

      .comment-content {
        font-size: 14px;
        line-height: 1.65;
      }
    }
  }

  .related-grid {
    grid-template-columns: 1fr;
    gap: 18px;
  }

  .related-card {
    .related-image {
      height: 180px;
    }

    .related-content {
      padding: 18px;

      h3 {
        font-size: 17px;
        margin-bottom: 10px;
      }

      p {
        font-size: 13px;
        margin-bottom: 14px;
      }

      .related-meta {
        gap: 16px;
        font-size: 12px;

        span span[class^="i-"] {
          font-size: 15px;
        }
      }
    }
  }

  .back-to-top {
    bottom: 28px;
    right: 28px;
    width: 50px !important;
    height: 50px !important;
  }
}

@media (max-width: 640px) {
  .back-navigation {
    padding: 14px 16px;

    :deep(.n-button) {
      font-size: 13px;
      padding: 0 8px;
    }
  }

  .post-article,
  .comments-section,
  .related-section {
    margin-left: 16px;
    margin-right: 16px;
    border-radius: 12px;
  }

  .post-banner {
    padding-top: 60%; // 5:3 比例
  }

  .article-content {
    padding: 32px 24px 28px;
  }

  .post-header {
    margin-bottom: 28px;

    .post-title {
      font-size: 26px;
      margin-bottom: 16px;
      line-height: 1.4;
    }

    .post-meta {
      gap: 14px;
      font-size: 12px;

      .meta-item {
        font-size: 12px;

        span[class^="i-"] {
          font-size: 14px;
        }
      }
    }
  }

  .post-wrapper {
    &::after {
      display: none;
    }
  }

  .post-content {
    font-size: 15px;
    line-height: 1.85;
    margin-bottom: 32px;

    :deep(h2) {
      font-size: 22px;
      margin: 36px 0 18px;
      padding-bottom: 10px;

      &::after {
        width: 40px;
      }
    }

    :deep(h3) {
      font-size: 19px;
      margin: 28px 0 16px;
      padding-left: 10px;
    }

    :deep(p) {
      margin: 16px 0;
      line-height: 1.9;
      text-indent: 1em;

      & + p {
        margin-top: 14px;
      }
    }

    :deep(img) {
      margin: 28px 0;
      border-radius: 8px;
    }

    :deep(code) {
      padding: 2px 6px;
      font-size: 0.86em;
    }

    :deep(pre) {
      padding: 16px;
      margin: 24px 0;
      border-radius: 8px;
      font-size: 13px;
    }

    :deep(blockquote) {
      padding: 14px 18px;
      margin: 24px 0;
      font-size: 14px;

      &::before {
        font-size: 45px;
      }
    }

    :deep(ul), :deep(ol) {
      padding-left: 28px;
      margin: 20px 0;

      li {
        margin: 12px 0;
      }
    }

    :deep(table) {
      font-size: 13px;

      th, td {
        padding: 10px 12px;
      }
    }
  }

  .post-actions {
    padding: 24px 0;

    :deep(.n-button) {
      width: 48px;
      height: 48px;
    }

    :deep(.n-space) {
      gap: 14px !important;
    }
  }

  .comments-section,
  .related-section {
    padding: 28px 24px;
  }

  .section-header {
    margin-bottom: 24px;

    h2 {
      font-size: 20px;

      span[class^="i-"] {
        font-size: 24px;
      }
    }
  }

  .comment-form {
    padding: 20px;
    margin-bottom: 28px;
    border-radius: 10px;

    h3 {
      font-size: 17px;
      margin-bottom: 18px;
    }

    :deep(.n-grid) {
      grid-template-columns: 1fr !important;
      gap: 0 !important;
    }

    :deep(.n-form-item) {
      margin-bottom: 18px;
    }

    :deep(.n-input) {
      font-size: 14px;
    }

    :deep(.n-button) {
      width: 100%;
      padding: 0 24px;
      height: 42px;
      font-size: 14px;
    }
  }

  .comments-list {
    gap: 18px;
  }

  .comment-item {
    padding: 16px;
    gap: 12px;
    border-radius: 10px;

    .comment-avatar {
      :deep(.n-avatar) {
        --n-size: 40px !important;
      }
    }

    .comment-body {
      .comment-header {
        .comment-author {
          font-size: 14px;
        }

        .comment-date {
          font-size: 12px;
        }
      }

      .comment-content {
        font-size: 13px;
        line-height: 1.6;
      }
    }
  }

  .related-card {
    border-radius: 10px;

    &:hover {
      transform: translateY(-6px);
    }

    .related-image {
      height: 160px;
    }

    .related-content {
      padding: 16px;

      h3 {
        font-size: 16px;
        margin-bottom: 8px;
        -webkit-line-clamp: 2;
      }

      p {
        font-size: 13px;
        margin-bottom: 12px;
        -webkit-line-clamp: 2;
      }

      .related-meta {
        gap: 14px;
        font-size: 11px;
      }
    }
  }

  .back-to-top {
    bottom: 24px;
    right: 24px;
    width: 48px !important;
    height: 48px !important;
  }
}

@media (max-width: 480px) {
  .back-navigation {
    padding: 12px 12px;

    :deep(.n-button) {
      font-size: 12px;
      padding: 0 6px;

      span[class^="i-"] {
        font-size: 16px;
      }
    }
  }

  .post-article,
  .comments-section,
  .related-section {
    margin-left: 12px;
    margin-right: 12px;
    border-radius: 10px;
  }

  .post-banner {
    padding-top: 66.67%; // 3:2 比例，更适合小屏
  }

  .article-content {
    padding: 28px 20px 24px;

    &::before {
      height: 5px;
    }
  }

  .post-header {
    margin-bottom: 24px;

    .post-title {
      font-size: 22px;
      margin-bottom: 14px;
      letter-spacing: -0.2px;
      line-height: 1.45;
    }

    .post-meta {
      gap: 12px;
      flex-wrap: wrap;

      .meta-item {
        font-size: 11px;

        span[class^="i-"] {
          font-size: 13px;
        }
      }
    }
  }

  .post-wrapper {
    &::after {
      display: none;
    }
  }

  .post-content {
    font-size: 14px;
    line-height: 1.8;
    letter-spacing: 0.2px;
    margin-bottom: 28px;

    :deep(h2) {
      font-size: 20px;
      margin: 32px 0 16px;
      padding-bottom: 8px;
      letter-spacing: -0.2px;

      &::after {
        width: 35px;
      }
    }

    :deep(h3) {
      font-size: 17px;
      margin: 24px 0 14px;
      padding-left: 8px;
      border-left-width: 3px;
    }

    :deep(p) {
      margin: 14px 0;
      line-height: 1.85;
      text-indent: 0; // 移动端取消首行缩进
      text-align: left; // 左对齐更适合小屏

      & + p {
        margin-top: 12px;
      }
    }

    :deep(img) {
      margin: 24px 0;
      border-radius: 6px;
      box-shadow: 0 6px 18px rgba(0, 0, 0, 0.1);

      &:hover {
        box-shadow: 0 8px 24px rgba(0, 0, 0, 0.15);
        transform: none; // 移动端取消 hover 位移
      }
    }

    :deep(code) {
      padding: 2px 6px;
      font-size: 0.85em;
      border-radius: 4px;
    }

    :deep(pre) {
      padding: 14px;
      margin: 20px -4px; // 负边距让代码块更宽
      border-radius: 6px;
      font-size: 12px;
      overflow-x: auto;
      -webkit-overflow-scrolling: touch;
    }

    :deep(blockquote) {
      padding: 12px 16px;
      margin: 20px 0;
      border-left-width: 3px;
      font-size: 13px;
      border-radius: 0 6px 6px 0;

      &::before {
        font-size: 40px;
        top: -6px;
        left: 8px;
      }
    }

    :deep(ul), :deep(ol) {
      padding-left: 24px;
      margin: 18px 0;

      li {
        margin: 10px 0;
        font-size: 14px;
        line-height: 1.8;
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

  .post-actions {
    padding: 20px 0;

    :deep(.n-button) {
      width: 44px;
      height: 44px;
      border-width: 1px;

      &:hover {
        transform: translateY(-2px);
      }
    }

    :deep(.n-space) {
      gap: 12px !important;
    }
  }

  .comments-section,
  .related-section {
    padding: 24px 20px;
  }

  .section-header {
    margin-bottom: 20px;

    h2 {
      font-size: 18px;
      gap: 10px;

      span[class^="i-"] {
        font-size: 22px;
      }
    }
  }

  .comment-form {
    padding: 16px;
    margin-bottom: 24px;

    h3 {
      font-size: 16px;
      margin-bottom: 16px;
    }

    :deep(.n-form-item) {
      margin-bottom: 16px;
    }

    :deep(.n-form-item-label) {
      font-size: 13px;
    }

    :deep(.n-input) {
      font-size: 13px;

      &.n-input--textarea {
        .n-input__textarea-el {
          min-height: 100px !important;
        }
      }
    }

    :deep(.n-button) {
      height: 40px;
      font-size: 13px;
      padding: 0 20px;
    }
  }

  .comments-list {
    gap: 16px;
  }

  .comment-item {
    padding: 14px;
    gap: 10px;

    .comment-avatar {
      :deep(.n-avatar) {
        --n-size: 36px !important;
        font-size: 14px;
      }
    }

    .comment-body {
      .comment-header {
        gap: 5px;
        margin-bottom: 8px;

        .comment-author {
          font-size: 13px;
        }

        .comment-date {
          font-size: 11px;
        }
      }

      .comment-content {
        font-size: 12px;
        line-height: 1.6;
      }
    }
  }

  .related-grid {
    gap: 14px;
  }

  .related-card {
    &:hover {
      transform: translateY(-4px);
      box-shadow: 0 8px 20px rgba(0, 0, 0, 0.1);
    }

    .related-image {
      height: 140px;
    }

    .related-content {
      padding: 14px;

      h3 {
        font-size: 15px;
        margin-bottom: 8px;
      }

      p {
        font-size: 12px;
        margin-bottom: 10px;
      }

      .related-meta {
        gap: 12px;
        font-size: 10px;

        span span[class^="i-"] {
          font-size: 13px;
        }
      }
    }
  }

  .back-to-top {
    bottom: 20px;
    right: 20px;
    width: 46px !important;
    height: 46px !important;
    box-shadow: 0 6px 18px rgba(168, 85, 247, 0.25);

    &:hover {
      transform: translateY(-2px);
      box-shadow: 0 8px 24px rgba(168, 85, 247, 0.35);
    }
  }

  // Empty state 优化
  :deep(.n-empty) {
    margin: 40px 0 !important;

    .n-empty__icon {
      font-size: 48px !important;
    }

    .n-empty__description {
      font-size: 13px;
    }
  }

  // Spin 加载优化
  :deep(.n-spin-container) {
    min-height: 300px;
  }
}
</style>