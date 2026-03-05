<script setup lang="ts">
import { definePage } from 'unplugin-vue-router/runtime'
import { ref, onMounted, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { usePostStore, useCategoryStore } from '@/stores/modules/app'
import { useMessage } from 'naive-ui'
import { $t } from '@/locales'

definePage({
  name: 'category-detail',
  meta: {
    title: 'Category Detail',
    level: 3,
  },
})

const route = useRoute()
const router = useRouter()
const postStore = usePostStore()
const categoryStore = useCategoryStore()
const message = useMessage()

const loading = ref(false)
const category = ref<any>(null)
const posts = ref<any[]>([])
const pagination = ref({
  page: 1,
  pageSize: 10,
  itemCount: 0,
})

const categoryId = computed(() => {
  const id = route.params.id
  return id ? parseInt(id as string) : null
})

async function loadCategory() {
  if (!categoryId.value) return

  loading.value = true
  try {
    category.value = await categoryStore.getCategory(categoryId.value)
  } catch (error) {
    console.error('Load category failed:', error)
    message.error($t('page.post_detail.load_failed'))
  } finally {
    loading.value = false
  }
}

async function loadPosts() {
  if (!categoryId.value) return

  loading.value = true
  try {
    const res = await postStore.listPost(
      {
        page: pagination.value.page,
        pageSize: pagination.value.pageSize,
      },
      {
        status: 'POST_STATUS_PUBLISHED',
        categoryIds: [categoryId.value],
      }
    )
    posts.value = res.items || []
    pagination.value.itemCount = res.total || 0
  } catch (error) {
    console.error('Load posts failed:', error)
    message.error($t('page.post_detail.load_failed'))
  } finally {
    loading.value = false
  }
}

function getCategoryName() {
  return category.value?.translations?.[0]?.name || $t('page.home.category_default')
}

function getCategoryDescription() {
  return category.value?.translations?.[0]?.description || ''
}

function getPostTitle(post: any) {
  return post.translations?.[0]?.title || $t('page.post_detail.untitled')
}

function getPostSummary(post: any) {
  return post.translations?.[0]?.summary || ''
}

function getPostThumbnail(post: any) {
  return post.translations?.[0]?.thumbnail || '/placeholder.jpg'
}

function formatDate(dateString: string) {
  if (!dateString) return ''
  return new Date(dateString).toLocaleDateString()
}

function handleViewPost(id: number) {
  router.push({
    path: `/post/${id}`,
    query: {
      from: 'category',
      categoryId: categoryId.value?.toString()
    }
  })
}

function handlePageChange(page: number) {
  pagination.value.page = page
  loadPosts()
}

function handlePageSizeChange(pageSize: number) {
  pagination.value.pageSize = pageSize
  pagination.value.page = 1
  loadPosts()
}

onMounted(async () => {
  await loadCategory()
  await loadPosts()
})
</script>

<template>
  <div class="category-detail-page">
    <!-- Hero Section -->
    <div class="hero-section">
      <div class="hero-content">
        <h1>{{ getCategoryName() }}</h1>
        <p v-if="getCategoryDescription()" class="category-description">{{ getCategoryDescription() }}</p>
        <div class="category-stats">
          <div class="stat-item">
            <span class="i-carbon:document" />
            <span>{{ pagination.itemCount }} {{ $t('page.posts.articles') }}</span>
          </div>
        </div>
      </div>
    </div>

    <!-- Posts Section -->
    <div class="page-container">
      <n-spin :show="loading">
        <!-- Results Count -->
        <div v-if="posts.length > 0" class="results-info">
          <span>找到 <strong>{{ pagination.itemCount }}</strong> 篇文章</span>
        </div>

        <!-- Posts Grid -->
        <div v-if="posts.length > 0" class="posts-grid">
          <article
            v-for="post in posts"
            :key="post.id"
            class="post-card"
            @click="handleViewPost(post.id)"
          >
            <div class="post-image">
              <img :src="getPostThumbnail(post)" :alt="getPostTitle(post)" />
              <div class="image-overlay" />
            </div>
            <div class="post-content">
              <h3 class="post-title">{{ getPostTitle(post) }}</h3>
              <p class="post-summary">{{ getPostSummary(post) }}</p>
              <div class="post-meta">
                <div class="meta-item">
                  <span class="i-carbon:user" />
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
            </div>
          </article>
        </div>

        <n-empty v-else :description="$t('page.posts.no_results')" style="margin: 80px 0;">
          <template #icon>
            <span class="i-carbon:document-blank" style="font-size: 64px;" />
          </template>
        </n-empty>

        <!-- Pagination -->
        <div v-if="posts.length > 0" class="pagination-wrapper">
          <n-pagination
            :page="pagination.page"
            :page-size="pagination.pageSize"
            :item-count="pagination.itemCount"
            :page-slot="7"
            show-size-picker
            :page-sizes="[10, 20, 30, 40]"
            @update:page="handlePageChange"
            @update:page-size="handlePageSizeChange"
          />
        </div>
      </n-spin>
    </div>
  </div>
</template>

<style scoped lang="less">
.category-detail-page {
  min-height: 100vh;
  background: var(--color-bg);
}

// Hero Section
.hero-section {
  background: linear-gradient(135deg, var(--color-brand) 0%, #764ba2 100%);
  padding: 60px 32px;
  margin-bottom: 40px;
  position: relative;
  overflow: hidden;

  &::before {
    content: '';
    position: absolute;
    top: 0;
    left: 0;
    right: 0;
    bottom: 0;
    background: url('data:image/svg+xml,<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 1440 320"><path fill="rgba(255,255,255,0.05)" d="M0,96L48,112C96,128,192,160,288,160C384,160,480,128,576,122.7C672,117,768,139,864,138.7C960,139,1056,117,1152,101.3C1248,85,1344,75,1392,69.3L1440,64L1440,320L1392,320C1344,320,1248,320,1152,320C1056,320,960,320,864,320C768,320,672,320,576,320C480,320,384,320,288,320C192,320,96,320,48,320L0,320Z"></path></svg>') no-repeat bottom;
    background-size: cover;
    opacity: 0.3;
  }

  .hero-content {
    max-width: 1200px;
    margin: 0 auto;
    text-align: center;
    position: relative;
    z-index: 1;

    h1 {
      font-size: 48px;
      font-weight: 800;
      margin: 0 0 16px 0;
      color: #ffffff;
      text-shadow: 0 2px 8px rgba(0, 0, 0, 0.2);
      letter-spacing: -1px;
    }

    .category-description {
      font-size: 20px;
      color: rgba(255, 255, 255, 0.95);
      margin: 0 0 24px 0;
      font-weight: 500;
      line-height: 1.5;
      max-width: 700px;
      margin-left: auto;
      margin-right: auto;
    }

    .category-stats {
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

// Results Info
.results-info {
  margin-bottom: 24px;
  font-size: 15px;
  color: var(--color-text-secondary);
  font-weight: 500;

  strong {
    color: var(--color-brand);
    font-weight: 700;
  }
}

// Posts Grid
.posts-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(320px, 1fr));
  gap: 32px;
  margin-bottom: 48px;
}

.post-card {
  background: var(--color-surface);
  border-radius: 16px;
  overflow: hidden;
  cursor: pointer;
  transition: all 0.4s cubic-bezier(0.4, 0, 0.2, 1);
  border: 1px solid var(--color-border);
  height: 100%;
  display: flex;
  flex-direction: column;

  &:hover {
    transform: translateY(-8px);
    box-shadow: 0 12px 32px rgba(0, 0, 0, 0.12);
    border-color: var(--color-brand);

    .post-image {
      .image-overlay {
        background: rgba(102, 126, 234, 0.15);
      }

      img {
        transform: scale(1.08);
      }
    }

    .post-title {
      color: var(--color-brand);
    }
  }

  .post-image {
    position: relative;
    width: 100%;
    height: 220px;
    overflow: hidden;
    background: var(--color-bg);

    img {
      width: 100%;
      height: 100%;
      object-fit: cover;
      transition: transform 0.5s cubic-bezier(0.4, 0, 0.2, 1);
    }

    .image-overlay {
      position: absolute;
      top: 0;
      left: 0;
      right: 0;
      bottom: 0;
      background: rgba(0, 0, 0, 0.05);
      transition: all 0.4s;
    }
  }

  .post-content {
    padding: 24px;
    flex: 1;
    display: flex;
    flex-direction: column;

    .post-title {
      font-size: 19px;
      font-weight: 700;
      margin: 0 0 12px 0;
      color: var(--color-text-primary);
      display: -webkit-box;
      -webkit-line-clamp: 2;
      -webkit-box-orient: vertical;
      overflow: hidden;
      line-height: 1.4;
      transition: color 0.3s;
      min-height: 53px;
    }

    .post-summary {
      color: var(--color-text-secondary);
      font-size: 15px;
      line-height: 1.7;
      margin: 0 0 16px 0;
      display: -webkit-box;
      -webkit-line-clamp: 3;
      -webkit-box-orient: vertical;
      overflow: hidden;
      flex: 1;
    }

    .post-meta {
      display: flex;
      flex-wrap: wrap;
      gap: 16px;
      padding-top: 16px;
      border-top: 1px solid var(--color-border);
      font-size: 13px;
      color: var(--color-text-secondary);

      .meta-item {
        display: flex;
        align-items: center;
        gap: 6px;
        font-weight: 500;

        span[class^="i-"] {
          font-size: 16px;
          opacity: 0.8;
        }
      }
    }
  }
}

// Pagination
.pagination-wrapper {
  display: flex;
  justify-content: center;
  padding: 40px 0;

  :deep(.n-pagination) {
    .n-pagination-item {
      border-radius: 8px;
      transition: all 0.3s;

      &.n-pagination-item--active {
        background: var(--color-brand);
        color: #fff;
      }

      &:not(.n-pagination-item--disabled):hover {
        background: var(--post-accent-bg-hover);
        color: var(--color-brand);
      }
    }
  }
}

// Responsive Design
@media (max-width: 1024px) {
  .hero-section {
    padding: 48px 24px;
  }

  .page-container {
    padding: 0 24px 60px;
  }

  .posts-grid {
    gap: 24px;
  }

  .hero-section .hero-content {
    h1 {
      font-size: 36px;
    }

    .category-description {
      font-size: 18px;
    }

    .category-stats {
      gap: 24px;
    }
  }
}

@media (max-width: 768px) {
  .hero-section {
    padding: 40px 20px;

    .hero-content {
      h1 {
        font-size: 32px;
      }

      .category-description {
        font-size: 16px;
      }

      .category-stats {
        gap: 16px;

        .stat-item {
          font-size: 14px;
        }
      }
    }
  }

  .page-container {
    padding: 0 20px 40px;
  }

  .posts-grid {
    grid-template-columns: repeat(auto-fill, minmax(280px, 1fr));
    gap: 20px;
  }

  .post-card {
    .post-image {
      height: 200px;
    }

    .post-content {
      padding: 20px;

      .post-title {
        font-size: 17px;
      }

      .post-summary {
        font-size: 14px;
      }
    }
  }
}

@media (max-width: 480px) {
  .hero-section {
    padding: 32px 16px;

    .hero-content {
      h1 {
        font-size: 28px;
      }

      .category-description {
        font-size: 15px;
      }

      .category-stats {
        flex-direction: column;
        gap: 12px;

        .stat-item {
          font-size: 13px;
        }
      }
    }
  }

  .page-container {
    padding: 0 16px 40px;
  }

  .posts-grid {
    grid-template-columns: 1fr;
    gap: 16px;
  }

  .post-card {
    .post-image {
      height: 180px;
    }

    .post-content {
      padding: 16px;

      .post-title {
        font-size: 16px;
        min-height: auto;
      }

      .post-summary {
        font-size: 13px;
        -webkit-line-clamp: 2;
      }

      .post-meta {
        gap: 12px;
        font-size: 12px;
      }
    }
  }

  .pagination-wrapper {
    padding: 24px 0;

    :deep(.n-pagination) {
      .n-pagination-item {
        min-width: 32px;
        height: 32px;
      }
    }
  }
}

// Dark Mode Enhancements
html.dark {
  .hero-section {
    &::before {
      opacity: 0.1;
    }
  }

  .post-card {
    &:hover {
      box-shadow: 0 12px 32px rgba(0, 0, 0, 0.4);
    }
  }
}
</style>

