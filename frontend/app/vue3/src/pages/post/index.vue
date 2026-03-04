<script setup lang="ts">
import { definePage } from 'unplugin-vue-router/runtime'
import { ref, onMounted, watch } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { usePostStore, useCategoryStore } from '@/stores/modules/app'
import { $t } from '@/locales'

definePage({
  name: 'post-list',
  meta: {
    title: 'Posts',
    level: 2,
  },
})

const router = useRouter()
const route = useRoute()
const postStore = usePostStore()
const categoryStore = useCategoryStore()

const loading = ref(false)
const posts = ref<any[]>([])
const categories = ref<any[]>([])
const total = ref(0)
const pagination = ref({
  page: 1,
  pageSize: 12,
})

const selectedCategory = ref<number | null>(null)
const searchQuery = ref('')

async function loadPosts() {
  loading.value = true
  try {
    const filters: any = { status: 'POST_STATUS_PUBLISHED' }

    if (selectedCategory.value) {
      filters.categoryIds = [selectedCategory.value]
    }

    if (searchQuery.value) {
      filters.title___icontains = searchQuery.value
    }

    const res = await postStore.listPost(pagination.value, filters)
    posts.value = res.items || []
    total.value = res.total || 0
  } catch (error) {
    console.error('Load posts failed:', error)
  } finally {
    loading.value = false
  }
}

async function loadCategories() {
  try {
    const res = await categoryStore.listCategory(
      undefined,
      { status: 'CATEGORY_STATUS_ACTIVE' }
    )
    categories.value = res.items || []
  } catch (error) {
    console.error('Load categories failed:', error)
  }
}

function getPostTitle(post: any) {
  return post.translations?.[0]?.title || $t('page.home.untitled')
}

function getPostSummary(post: any) {
  return post.translations?.[0]?.summary || ''
}

function getPostThumbnail(post: any) {
  return post.translations?.[0]?.thumbnail || '/placeholder.jpg'
}

function getCategoryName(category: any) {
  return category.translations?.[0]?.name || $t('page.home.category_default')
}

function formatDate(dateString: string) {
  if (!dateString) return ''
  return new Date(dateString).toLocaleDateString()
}

function handleViewPost(id: number) {
  router.push(`/post/${id}`)
}

function handleCategoryChange(categoryId: number | null) {
  selectedCategory.value = categoryId
  pagination.value.page = 1
  loadPosts()
}

function handlePageChange(page: number) {
  pagination.value.page = page
  loadPosts()
  window.scrollTo({ top: 0, behavior: 'smooth' })
}

function handleSearch() {
  pagination.value.page = 1
  loadPosts()
}

onMounted(async () => {
  await loadCategories()

  // 检查 URL 参数
  const categoryId = route.query.category
  if (categoryId) {
    selectedCategory.value = parseInt(categoryId as string)
  }

  await loadPosts()
})

watch(() => route.query.category, (newVal) => {
  if (newVal) {
    selectedCategory.value = parseInt(newVal as string)
    loadPosts()
  }
})
</script>

<template>
  <div class="post-list-page">
    <div class="page-header">
      <h1>{{ $t('page.posts.posts_list') }}</h1>
      <p>{{ $t('page.posts.explore_latest') }}</p>
    </div>

    <div class="page-content">
      <!-- Sidebar -->
      <aside class="sidebar">
        <n-card :title="$t('page.posts.search')" class="search-card">
          <n-input
            v-model:value="searchQuery"
            :placeholder="$t('page.posts.search_posts')"
            clearable
            @keyup.enter="handleSearch"
          >
            <template #suffix>
              <n-button text @click="handleSearch">
                <div class="i-carbon:search" />
              </n-button>
            </template>
          </n-input>
        </n-card>

        <n-card :title="$t('page.posts.all_categories')" class="filter-card">
          <n-menu
            :value="selectedCategory"
            :options="[
              { label: $t('page.posts.all_categories'), key: null },
              ...categories.map(cat => ({
                label: getCategoryName(cat),
                key: cat.id,
              }))
            ]"
            @update:value="handleCategoryChange"
          />
        </n-card>
      </aside>

      <!-- Main Content -->
      <main class="main-content">
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
                  <span class="views">{{ $t('page.posts.views_count', { count: post.visits || 0 }) }}</span>
                </div>
              </div>
            </div>
          </div>

          <n-empty v-else :description="$t('page.posts.no_results')" />

          <div v-if="total > pagination.pageSize" class="pagination-wrapper">
            <n-pagination
              v-model:page="pagination.page"
              :page-count="Math.ceil(total / pagination.pageSize)"
              @update:page="handlePageChange"
            />
          </div>
        </n-spin>
      </main>
    </div>
  </div>
</template>

<style scoped lang="less">
.post-list-page {
  max-width: 1400px;
  margin: 0 auto;
  padding: 2rem;
  background: var(--color-bg);
  color: var(--color-text-primary);
}

.page-header {
  text-align: center;
  margin-bottom: 3rem;

  h1 {
    font-size: 2.5rem;
    font-weight: 700;
    margin: 0 0 0.5rem 0;
    color: var(--color-text-primary);
  }

  p {
    font-size: 1.1rem;
    color: var(--color-text-secondary);
    margin: 0;
  }
}

.page-content {
  display: grid;
  grid-template-columns: 280px 1fr;
  gap: 2rem;
}

.sidebar {
  display: flex;
  flex-direction: column;
  gap: 1.5rem;

  .search-card,
  .filter-card {
    background: var(--color-surface);
    border: 1px solid var(--color-border);

    :deep(.n-card__content) {
      padding: 1rem;
    }

    :deep(.n-card-header) {
      border-bottom: 1px solid var(--color-border);
      padding: 1rem;
    }
  }
}

.main-content {
  min-height: 400px;
}

.posts-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(320px, 1fr));
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
    box-shadow: var(--post-card-hover-shadow);
  }

  .post-image {
    width: 100%;
    height: 200px;
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
      -webkit-line-clamp: 3;
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

.pagination-wrapper {
  display: flex;
  justify-content: center;
  margin-top: 2rem;
}

@media (max-width: 968px) {
  .page-content {
    grid-template-columns: 1fr;
  }

  .sidebar {
    order: 2;
  }

  .main-content {
    order: 1;
  }

  .posts-grid {
    grid-template-columns: 1fr;
  }
}
</style>

