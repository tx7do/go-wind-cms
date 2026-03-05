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
  router.push(`/post/${id}`)
}

function handlePageChange(page: number) {
  pagination.value.page = page
  loadPosts()
}

onMounted(async () => {
  await loadCategory()
  await loadPosts()
})
</script>

<template>
  <div class="category-detail-page">
    <div class="category-header">
      <div class="header-content">
        <h1>{{ getCategoryName() }}</h1>
        <p v-if="getCategoryDescription()">{{ getCategoryDescription() }}</p>
        <p v-else class="no-description">{{ $t('page.home.category_default') }}</p>
        <div class="category-meta">
          <span>{{ pagination.itemCount }} {{ $t('page.posts.articles') }}</span>
        </div>
      </div>
    </div>

    <div class="posts-container">
      <n-spin :show="loading">
        <div v-if="posts.length > 0" class="posts-grid">
          <div
            v-for="post in posts"
            :key="post.id"
            class="post-card"
            @click="handleViewPost(post.id)"
          >
            <div class="post-image">
              <img :src="getPostThumbnail(post)" :alt="getPostTitle(post)" />
            </div>
            <div class="post-content">
              <h3>{{ getPostTitle(post) }}</h3>
              <p>{{ getPostSummary(post) }}</p>
              <div class="post-meta">
                <span class="author">{{ post.authorName }}</span>
                <span class="date">{{ formatDate(post.createdAt) }}</span>
                <span class="views">{{ $t('page.home.views_count', { count: post.visits || 0 }) }}</span>
              </div>
            </div>
          </div>
        </div>

        <n-empty
          v-else-if="!loading"
          :description="$t('page.posts.no_results')"
        />
      </n-spin>

      <!-- Pagination -->
      <div v-if="posts.length > 0" class="pagination">
        <n-pagination
          :page="pagination.page"
          :page-size="pagination.pageSize"
          :item-count="pagination.itemCount"
          @update:page="handlePageChange"
        />
      </div>
    </div>
  </div>
</template>

<style scoped lang="less">
.category-detail-page {
  min-height: 100vh;
  background: var(--color-bg);
}

.category-header {
  background: linear-gradient(135deg, var(--color-brand) 0%, var(--color-brand-accent) 100%);
  color: white;
  padding: 3rem 2rem;
  text-align: center;

  .header-content {
    max-width: 1200px;
    margin: 0 auto;

    h1 {
      font-size: 2.5rem;
      font-weight: 700;
      margin: 0 0 1rem 0;
    }

    p {
      font-size: 1.1rem;
      margin: 0 0 1rem 0;
      opacity: 0.95;
    }

    .no-description {
      opacity: 0.75;
    }

    .category-meta {
      display: flex;
      justify-content: center;
      gap: 2rem;
      font-size: 0.95rem;

      span {
        display: flex;
        align-items: center;
        gap: 0.5rem;
      }
    }
  }
}

.posts-container {
  max-width: 1200px;
  margin: 0 auto;
  padding: 3rem 2rem;
}

.posts-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(350px, 1fr));
  gap: 2rem;
  margin-bottom: 2rem;
}

.post-card {
  background: var(--color-surface);
  border-radius: var(--radius-md);
  overflow: hidden;
  cursor: pointer;
  transition: all 0.3s;
  border: 1px solid var(--color-border);

  &:hover {
    transform: translateY(-4px);
    box-shadow: 0 8px 16px rgba(0, 0, 0, 0.1);
  }

  .post-image {
    width: 100%;
    height: 180px;
    overflow: hidden;

    img {
      width: 100%;
      height: 100%;
      object-fit: cover;
      transition: transform 0.3s;
    }

    &:hover img {
      transform: scale(1.05);
    }
  }

  .post-content {
    padding: 1.5rem;

    h3 {
      font-size: 1.15rem;
      font-weight: 600;
      margin: 0 0 0.5rem 0;
      color: var(--color-text-primary);
      display: -webkit-box;
      -webkit-line-clamp: 2;
      -webkit-box-orient: vertical;
      overflow: hidden;
    }

    p {
      color: var(--color-text-secondary);
      font-size: 0.9rem;
      line-height: 1.6;
      margin: 0 0 1rem 0;
      display: -webkit-box;
      -webkit-line-clamp: 2;
      -webkit-box-orient: vertical;
      overflow: hidden;
    }

    .post-meta {
      display: flex;
      gap: 1rem;
      font-size: 0.85rem;
      color: var(--color-text-secondary);

      span {
        display: flex;
        align-items: center;
        gap: 0.25rem;
      }
    }
  }
}

.pagination {
  display: flex;
  justify-content: center;
  padding: 2rem 0;
}

@media (max-width: 768px) {
  .category-header {
    padding: 2rem 1rem;

    .header-content {
      h1 {
        font-size: 1.75rem;
      }

      p {
        font-size: 0.95rem;
      }

      .category-meta {
        flex-direction: column;
        gap: 0.5rem;
      }
    }
  }

  .posts-container {
    padding: 1.5rem 1rem;
  }

  .posts-grid {
    grid-template-columns: 1fr;
  }
}
</style>

