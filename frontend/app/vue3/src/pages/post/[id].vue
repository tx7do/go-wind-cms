<script setup lang="ts">
import {definePage} from 'unplugin-vue-router/runtime'
import {ref, onMounted, computed} from 'vue'
import {useRoute, useRouter} from 'vue-router'
import {usePostStore, useCommentStore} from '@/stores/modules/app'
import {useMessage} from 'naive-ui'
import {$t, currentLocaleLanguageCode} from '@/locales'
import {ContentViewer} from '@/components/ContentViewer'

definePage({
  name: 'post-detail',
  meta: {
    title: 'Post Detail',
    level: 3,
  },
})

const route = useRoute()
const router = useRouter()
const postStore = usePostStore()
const commentStore = useCommentStore()
const message = useMessage()

const loading = ref(false)
const post = ref<any>(null)
const comments = ref<any[]>([])
const relatedPosts = ref<any[]>([])
const showBackToTop = ref(false)
const isLiked = ref(false)
const isBookmarked = ref(false)

const newComment = ref({
  content: '',
  authorName: '',
  email: '',
})

const postId = computed(() => {
  const id = route.params.id
  return id ? parseInt(id as string) : null
})

async function loadPost() {
  if (!postId.value) return

  loading.value = true
  try {
    post.value = await postStore.getPost(postId.value)
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
      relatedPosts.value = (res.items || []).filter((p: any) => p.id !== postId.value)
    }
  } catch (error) {
    console.error('Load related posts failed:', error)
  }
}

function getTranslation() {
  const locale = currentLocaleLanguageCode();
  return post.value?.translations?.find((t: any) => t.languageCode === locale) || post.value?.translations?.[0]
}

function getTitle() {
  const translation = getTranslation();
  if (translation) {
    return translation.title || $t('page.post_detail.untitled')
  }

  return translation?.title || $t('page.post_detail.untitled')
}

function getContent() {
  const translation = getTranslation();
  if (translation) {
    return translation.content || ''
  }

  return ''
}

function getThumbnail() {
  const translation = getTranslation();
  if (translation && translation.thumbnail) {
    return translation.thumbnail
  }

  return '/placeholder.jpg'
}

function formatDate(dateString: string) {
  if (!dateString) return ''
  return new Date(dateString).toLocaleString()
}

async function handleSubmitComment() {
  if (!newComment.value.content || !newComment.value.authorName || !newComment.value.email) {
    message.warning($t('page.post_detail.fill_form_info'))
    return
  }

  if (!postId.value) return

  try {
    await commentStore.createComment({
      postId: postId.value,
      content: newComment.value.content,
      authorName: newComment.value.authorName,
      email: newComment.value.email,
      status: 'COMMENT_STATUS_PENDING',
    })

    message.success($t('page.post_detail.comment_submitted'))

    // 重置表单
    newComment.value = {
      content: '',
      authorName: '',
      email: '',
    }
  } catch (error) {
    console.error('Submit comment failed:', error)
    message.error($t('page.post_detail.submit_comment_failed'))
  }
}

function handleViewRelatedPost(id: number) {
  router.push(`/post/${id}`)
  window.scrollTo({top: 0, behavior: 'smooth'})
}

function handleBack() {
  // 检查是否有来源页面信息
  const from = route.query.from as string
  const categoryId = route.query.categoryId as string

  if (from === 'category' && categoryId) {
    // 从分类页面进入，返回到分类页面
    router.push(`/category/${categoryId}`)
  } else if (from === 'user') {
    // 从用户资料页面进入，返回到用户资料页面
    router.push('/user')
  } else if (from === 'post') {
    // 从文章列表页面进入，返回到文章列表页面
    router.push('/post')
  } else if (from === 'home') {
    // 从首页进入，返回到首页
    router.push('/')
  } else if (window.history.length > 2) {
    // 有浏览历史，使用浏览器返回
    router.back()
  } else {
    // 默认返回文章列表
    router.push('/post')
  }
}

function handleLike() {
  isLiked.value = !isLiked.value
  if (isLiked.value) {
    message.success($t('page.post_detail.liked'))
    if (post.value) {
      post.value.likes = (post.value.likes || 0) + 1
    }
  } else {
    if (post.value && post.value.likes > 0) {
      post.value.likes -= 1
    }
  }
}

function handleBookmark() {
  isBookmarked.value = !isBookmarked.value
  if (isBookmarked.value) {
    message.success($t('page.post_detail.bookmarked'))
  } else {
    message.info($t('page.post_detail.unbookmarked'))
  }
}

function handleShare() {
  const url = window.location.href
  if (navigator.share) {
    navigator.share({
      title: getTitle(),
      url: url
    }).then(() => {
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

function handleScroll() {
  showBackToTop.value = window.scrollY > 500
}

onMounted(async () => {
  await loadPost()
  await Promise.all([
    loadComments(),
    loadRelatedPosts(),
  ])

  // 监听滚动事件
  window.addEventListener('scroll', handleScroll)
})

// 清理事件监听
import {onBeforeUnmount} from 'vue'

onBeforeUnmount(() => {
  window.removeEventListener('scroll', handleScroll)
})
</script>

<template>
  <div class="post-detail-page">
    <n-spin :show="loading">
      <div v-if="post" class="post-container">
        <!-- Back Button -->
        <div class="back-navigation">
          <n-button text @click="handleBack()">
            <template #icon>
              <span class="i-carbon:arrow-left"/>
            </template>
            {{ $t('page.post_detail.back') }}
          </n-button>
        </div>

        <!-- Post Header with Thumbnail -->
        <article class="post-article">
          <!-- Post Thumbnail Banner -->
          <div v-if="getThumbnail()" class="post-banner">
            <img :src="getThumbnail()" :alt="getTitle()"/>
            <div class="banner-overlay"/>
          </div>

          <div class="article-content">
            <!-- Post Header -->
            <header class="post-header">
              <h1 class="post-title">{{ getTitle() }}</h1>
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
              <ContentViewer
                :content="getContent()"
                type="markdown"
              />
            </div>

            <!-- Post Actions -->
            <footer class="post-actions">
              <n-space size="large">
                <n-button
                  size="large"
                  circle
                  :type="isLiked ? 'primary' : 'default'"
                  @click="handleLike"
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
                >
                  <template #icon>
                    <span :class="isBookmarked ? 'i-carbon:bookmark-filled' : 'i-carbon:bookmark'"/>
                  </template>
                </n-button>
                <n-button
                  size="large"
                  circle
                  @click="handleShare"
                >
                  <template #icon>
                    <span class="i-carbon:share"/>
                  </template>
                </n-button>
              </n-space>
            </footer>
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
                  />
                </n-form-item-gi>
                <n-form-item-gi :label="$t('page.post_detail.email')">
                  <n-input
                    v-model:value="newComment.email"
                    :placeholder="$t('page.post_detail.enter_email')"
                    type="text"
                    size="large"
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
                />
              </n-form-item>
              <n-form-item>
                <n-button type="primary" size="large" @click="handleSubmitComment">
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
                  <ContentViewer
                    :content="comment.content"
                    type="text"
                  />
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
            >
              <div class="related-image">
                <img
                  :src="relatedPost.translations?.[0]?.thumbnail || '/placeholder.jpg'"
                  :alt="relatedPost.translations?.[0]?.title"
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

// Post Banner
.post-banner {
  position: relative;
  width: 100%;
  padding-top: 56.25%; // 16:9 比例
  overflow: hidden;
  background: var(--color-bg);

  img {
    position: absolute;
    top: 0;
    left: 0;
    width: 100%;
    height: 100%;
    object-fit: cover;
  }

  .banner-overlay {
    position: absolute;
    bottom: 0;
    left: 0;
    right: 0;
    height: 200px;
    background: linear-gradient(to bottom,
    transparent 0%,
    rgba(0, 0, 0, 0.3) 40%,
    rgba(0, 0, 0, 0.7) 80%,
    rgba(0, 0, 0, 0.85) 100%);
    pointer-events: none;
  }
}

// Article Content
.article-content {
  padding: 56px 64px 48px;
  background: var(--color-surface);
  position: relative;

  // 顶部留白过渡区域
  &::before {
    content: '';
    position: absolute;
    top: 0;
    left: 0;
    right: 0;
    height: 8px;
    background: linear-gradient(to bottom, rgba(0, 0, 0, 0.03), transparent);
  }
}

// Post Header
.post-header {
  margin-bottom: 40px;

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
  font-size: 18px;
  line-height: 2;
  letter-spacing: 0.3px;
  color: var(--color-text-primary);
  margin-bottom: 48px;

  :deep(h2) {
    font-size: 32px;
    font-weight: 700;
    margin: 56px 0 28px;
    padding-bottom: 16px;
    border-bottom: 3px solid var(--color-brand);
    color: var(--color-text-primary);
    position: relative;

    &::after {
      content: '';
      position: absolute;
      bottom: -3px;
      left: 0;
      width: 60px;
      height: 3px;
      background: linear-gradient(90deg, #a855f7, transparent);
    }
  }

  :deep(h3) {
    font-size: 26px;
    font-weight: 600;
    margin: 44px 0 24px;
    color: var(--color-text-primary);
    padding-left: 16px;
    border-left: 4px solid var(--color-brand);
  }

  :deep(p) {
    margin: 24px 0;
    text-align: justify;
    text-indent: 2em;
    line-height: 2.2;

    // 段落间距
    & + p {
      margin-top: 20px;
    }
  }

  :deep(img) {
    max-width: 100%;
    border-radius: 12px;
    margin: 40px 0;
    box-shadow: 0 8px 24px rgba(0, 0, 0, 0.12);
    transition: all 0.3s;
    cursor: zoom-in;

    &:hover {
      box-shadow: 0 12px 32px rgba(0, 0, 0, 0.18);
      transform: translateY(-4px);
    }
  }

  :deep(code) {
    background: var(--color-bg);
    padding: 4px 10px;
    border-radius: 6px;
    font-family: 'Fira Code', 'Courier New', monospace;
    font-size: 0.9em;
    color: var(--color-brand);
    border: 1px solid rgba(168, 85, 247, 0.2);
  }

  :deep(pre) {
    background: var(--color-bg);
    padding: 24px;
    border-radius: 12px;
    overflow-x: auto;
    border: 1px solid var(--color-border);
    margin: 32px 0;

    code {
      background: none;
      padding: 0;
      color: var(--color-text-primary);
      border: none;
    }
  }

  :deep(blockquote) {
    border-left: 5px solid var(--color-brand);
    padding: 20px 24px;
    margin: 32px 0;
    background: rgba(168, 85, 247, 0.05);
    border-radius: 0 8px 8px 0;
    color: var(--color-text-secondary);
    font-style: italic;
    position: relative;

    &::before {
      content: '"';
      position: absolute;
      top: -10px;
      left: 10px;
      font-size: 60px;
      color: var(--color-brand);
      opacity: 0.3;
      font-family: Georgia, serif;
    }
  }

  :deep(ul), :deep(ol) {
    padding-left: 40px;
    margin: 28px 0;

    li {
      margin: 16px 0;
      line-height: 2;

      &::marker {
        color: var(--color-brand);
        font-weight: 600;
      }
    }
  }

  :deep(a) {
    color: var(--color-brand);
    text-decoration: none;
    border-bottom: 2px solid rgba(168, 85, 247, 0.3);
    transition: all 0.3s;
    font-weight: 500;

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
  }

  :deep(em) {
    color: var(--color-text-secondary);
    font-style: italic;
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
    box-shadow: 0 2px 8px rgba(0, 0, 0, 0.08);

    th {
      background: var(--color-brand);
      color: white;
      padding: 16px;
      text-align: left;
      font-weight: 600;
    }

    td {
      padding: 14px 16px;
      border-bottom: 1px solid var(--color-border);
    }

    tr:hover {
      background: rgba(168, 85, 247, 0.05);
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
    padding: 48px 48px 40px;

    &::before {
      height: 6px;
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
    font-size: 17px;
    margin-bottom: 40px;

    :deep(h2) {
      font-size: 28px;
      margin: 48px 0 24px;
      padding-bottom: 14px;

      &::after {
        width: 50px;
      }
    }

    :deep(h3) {
      font-size: 23px;
      margin: 36px 0 20px;
      padding-left: 14px;
    }

    :deep(p) {
      margin: 20px 0;
      line-height: 2.1;
    }

    :deep(img) {
      margin: 36px 0;
    }
  }

  .comments-section,
  .related-section {
    padding: 40px 48px;
  }

  .post-banner {
    padding-top: 50%; // 调整为 2:1 比例
  }

  .back-to-top {
    bottom: 32px;
    right: 32px;
    width: 54px !important;
    height: 54px !important;
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

  .article-content {
    padding: 40px 32px 36px;
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

