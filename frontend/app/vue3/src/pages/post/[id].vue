<script setup lang="ts">
import { definePage } from 'unplugin-vue-router/runtime'
import { ref, onMounted, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { usePostStore, useCommentStore } from '@/stores/modules/app'
import { useMessage } from 'naive-ui'
import { $t } from '@/locales'
import { ContentViewer } from '@/components/ContentViewer'

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
      { page: 1, pageSize: 50 },
      { postId: postId.value, status: 'COMMENT_STATUS_APPROVED' }
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
        { page: 1, pageSize: 3 },
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

function getTitle() {
  return post.value?.translations?.[0]?.title || $t('page.post_detail.untitled')
}

function getContent() {
  return post.value?.translations?.[0]?.content || ''
}

function getThumbnail() {
  return post.value?.translations?.[0]?.thumbnail || '/placeholder.jpg'
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
  window.scrollTo({ top: 0, behavior: 'smooth' })
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

onMounted(async () => {
  await loadPost()
  await Promise.all([
    loadComments(),
    loadRelatedPosts(),
  ])
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
              <span class="i-carbon:arrow-left" />
            </template>
            {{ $t('page.post_detail.back') }}
          </n-button>
        </div>

        <!-- Post Header with Thumbnail -->
        <article class="post-article">
          <!-- Post Thumbnail Banner -->
          <div v-if="getThumbnail()" class="post-banner">
            <img :src="getThumbnail()" :alt="getTitle()" />
            <div class="banner-overlay" />
          </div>

          <div class="article-content">
            <!-- Post Header -->
            <header class="post-header">
              <h1 class="post-title">{{ getTitle() }}</h1>
              <div class="post-meta">
                <div class="meta-item">
                  <span class="i-carbon:user-avatar" />
                  <span>{{ post.authorName }}</span>
                </div>
                <div class="meta-item">
                  <span class="i-carbon:calendar" />
                  <span>{{ formatDate(post.createdAt) }}</span>
                </div>
                <div class="meta-item">
                  <span class="i-carbon:view" />
                  <span>{{ post.visits || 0 }}</span>
                </div>
                <div class="meta-item">
                  <span class="i-carbon:thumbs-up" />
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
                <n-button size="large" circle>
                  <template #icon>
                    <span class="i-carbon:thumbs-up" />
                  </template>
                </n-button>
                <n-button size="large" circle>
                  <template #icon>
                    <span class="i-carbon:bookmark" />
                  </template>
                </n-button>
                <n-button size="large" circle>
                  <template #icon>
                    <span class="i-carbon:share" />
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
              <span class="i-carbon:chat" />
              {{ $t('page.post_detail.comments_count', { count: comments.length }) }}
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
                    <span class="i-carbon:send-alt" />
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
          <n-empty v-else :description="$t('page.post_detail.no_comments')" style="margin-top: 40px;" />
        </section>

        <!-- Related Posts -->
        <section v-if="relatedPosts.length > 0" class="related-section">
          <div class="section-header">
            <h2>
              <span class="i-carbon:book" />
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
                <div class="image-overlay" />
              </div>
              <div class="related-content">
                <h3>{{ relatedPost.translations?.[0]?.title }}</h3>
                <p>{{ relatedPost.translations?.[0]?.summary }}</p>
                <div class="related-meta">
                  <span><span class="i-carbon:view" /> {{ relatedPost.visits || 0 }}</span>
                  <span><span class="i-carbon:thumbs-up" /> {{ relatedPost.likes || 0 }}</span>
                </div>
              </div>
            </div>
          </div>
        </section>
      </div>

      <n-empty v-else :description="$t('page.post_detail.post_not_found')" />
    </n-spin>
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
  height: 500px;
  overflow: hidden;

  img {
    width: 100%;
    height: 100%;
    object-fit: cover;
  }

  .banner-overlay {
    position: absolute;
    bottom: 0;
    left: 0;
    right: 0;
    height: 150px;
    background: linear-gradient(to bottom, transparent, rgba(0, 0, 0, 0.6));
  }
}

// Article Content
.article-content {
  padding: 48px 64px;
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
  line-height: 1.8;
  color: var(--color-text-primary);
  margin-bottom: 48px;

  :deep(h2) {
    font-size: 32px;
    font-weight: 700;
    margin: 48px 0 24px;
    padding-bottom: 12px;
    border-bottom: 2px solid var(--color-border);
  }

  :deep(h3) {
    font-size: 26px;
    font-weight: 600;
    margin: 36px 0 20px;
  }

  :deep(p) {
    margin: 20px 0;
    text-align: justify;
  }

  :deep(img) {
    max-width: 100%;
    border-radius: 12px;
    margin: 32px 0;
    box-shadow: 0 4px 12px rgba(0, 0, 0, 0.08);
  }

  :deep(code) {
    background: var(--color-bg);
    padding: 4px 8px;
    border-radius: 6px;
    font-family: 'Fira Code', 'Courier New', monospace;
    font-size: 0.9em;
    color: var(--color-brand);
  }

  :deep(pre) {
    background: var(--color-bg);
    padding: 24px;
    border-radius: 12px;
    overflow-x: auto;
    border: 1px solid var(--color-border);

    code {
      background: none;
      padding: 0;
      color: var(--color-text-primary);
    }
  }

  :deep(blockquote) {
    border-left: 4px solid var(--color-brand);
    padding-left: 24px;
    margin: 24px 0;
    color: var(--color-text-secondary);
    font-style: italic;
  }

  :deep(ul), :deep(ol) {
    padding-left: 32px;
    margin: 20px 0;

    li {
      margin: 12px 0;
    }
  }

  :deep(a) {
    color: var(--color-brand);
    text-decoration: none;
    border-bottom: 1px solid transparent;
    transition: all 0.3s;

    &:hover {
      border-bottom-color: var(--color-brand);
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

// Responsive Design
@media (max-width: 1024px) {
  .back-navigation,
  .post-article,
  .comments-section,
  .related-section {
    margin-left: 20px;
    margin-right: 20px;
  }

  .article-content,
  .comments-section,
  .related-section {
    padding: 32px 40px;
  }

  .post-banner {
    height: 400px;
  }
}

@media (max-width: 768px) {
  .back-navigation {
    padding: 16px 20px;
  }

  .post-banner {
    height: 300px;
  }

  .article-content,
  .comments-section,
  .related-section {
    padding: 24px 20px;
  }

  .post-header {
    .post-title {
      font-size: 28px;
    }

    .post-meta {
      gap: 16px;

      .meta-item {
        font-size: 13px;
      }
    }
  }

  .post-content {
    font-size: 16px;

    :deep(h2) {
      font-size: 24px;
      margin: 32px 0 16px;
    }

    :deep(h3) {
      font-size: 20px;
      margin: 24px 0 12px;
    }
  }

  .post-actions {
    :deep(.n-button) {
      width: 48px;
      height: 48px;
    }
  }

  .section-header h2 {
    font-size: 22px;
  }

  .comment-form {
    padding: 20px;

    h3 {
      font-size: 18px;
    }
  }

  .comment-item {
    padding: 16px;

    .comment-avatar {
      :deep(.n-avatar) {
        --n-size: 40px !important;
      }
    }
  }

  .related-grid {
    grid-template-columns: 1fr;
  }
}

@media (max-width: 480px) {
  .back-navigation,
  .post-article,
  .comments-section,
  .related-section {
    margin-left: 12px;
    margin-right: 12px;
    border-radius: 12px;
  }

  .article-content,
  .comments-section,
  .related-section {
    padding: 20px 16px;
  }

  .post-banner {
    height: 250px;
  }

  .post-header .post-title {
    font-size: 24px;
  }

  .post-content {
    font-size: 15px;
  }
}
</style>

