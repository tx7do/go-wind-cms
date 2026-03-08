<script setup lang="ts">
import {ref, onMounted} from 'vue'

import {$t} from '@/locales'
import {usePostStore} from '@/stores'
import {useLanguageChangeEffect} from '@/hooks/useLanguageChangeEffect';
import PostCard from '@/components/PostCard';
import type {contentservicev1_Post} from "@/api/generated/app/service/v1";

interface Props {
  initialPageSize?: number
  pageSizes?: number[]
  categoryId?: number | null
}

const props = withDefaults(defineProps<Props>(), {
  initialPageSize: 12,
  pageSizes: () => [12, 24, 36, 48],
  categoryId: null
})

interface Emits {
  (e: 'category-change', categoryId: number | null): void

  (e: 'page-change', page: number): void

  (e: 'page-size-change', pageSize: number): void
}

const emit = defineEmits<Emits>()

const postStore = usePostStore()

const loading = ref(false)
const posts = ref<contentservicev1_Post[]>([])
const total = ref(0)
const pagination = ref({
  page: 1,
  pageSize: props.initialPageSize,
})

async function loadPosts() {
  loading.value = true
  try {
    const filters: any = {status: 'POST_STATUS_PUBLISHED'}

    if (props.categoryId) {
      filters.categoryIds = [props.categoryId]
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

function handlePageChange(page: number) {
  pagination.value.page = page
  emit('page-change', page)
  loadPosts()
  window.scrollTo({top: 0, behavior: 'smooth'})
}

function handlePageSizeChange(pageSize: number) {
  pagination.value.pageSize = pageSize
  pagination.value.page = 1
  emit('page-size-change', pageSize)
  loadPosts()
}

// 监听语言切换，自动重新加载数据
useLanguageChangeEffect(async () => {
  await loadPosts();
}, {
  immediate: false,
  autoCleanup: true,
});

onMounted(async () => {
  await loadPosts()
})

// 暴露方法给父组件
defineExpose({
  refresh: loadPosts
})
</script>

<template>
  <div class="post-list-with-pagination">
    <!-- Main Content -->
    <main class="main-content">
      <!-- Loading Skeleton -->
      <div v-if="loading" class="posts-grid">
        <div v-for="i in pagination.pageSize" :key="i" class="post-card">
          <n-skeleton height="380px"/>
        </div>
      </div>

      <!-- Loaded Content -->
      <div v-else>
        <!-- Results Count -->
        <div v-if="posts.length > 0" class="results-info">
          <span>{{ $t('page.posts.found') }} <strong>{{
              total
            }}</strong> {{ $t('page.posts.articles') }}</span>
        </div>

        <!-- Posts Grid -->
        <div v-if="posts.length > 0" class="posts-grid">
          <PostCard
            v-for="post in posts"
            :key="post.id"
            :post="post"
          />
        </div>

        <!-- Empty State -->
        <n-empty v-else :description="$t('page.posts.no_results')" style="margin: 80px 0;">
          <template #icon>
            <span class="i-carbon:document-blank" style="font-size: 64px;"/>
          </template>
        </n-empty>

        <!-- Pagination -->
        <div v-if="total > pagination.pageSize" class="pagination-wrapper">
          <n-pagination
            v-model:page="pagination.page"
            :page-count="Math.ceil(total / pagination.pageSize)"
            :page-slot="7"
            show-size-picker
            :page-sizes="pageSizes"
            @update:page="handlePageChange"
            @update:page-size="handlePageSizeChange"
          />
        </div>
      </div>
    </main>
  </div>
</template>

<style scoped lang="less">
.post-list-with-pagination {
  width: 100%;
}

// Main Content
.main-content {
  // Posts Grid
  .posts-grid {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(320px, 1fr));
    gap: 32px;
    margin-bottom: 48px;
  }

  // Results Info
  .results-info {
    margin-bottom: 24px;
    padding: 16px 20px;
    background: var(--color-surface);
    border-radius: 12px;
    border: 1px solid var(--color-border);
    font-size: 14px;
    color: var(--color-text-secondary);

    strong {
      color: var(--color-brand);
      font-weight: 700;
    }
  }

  // Pagination
  .pagination-wrapper {
    display: flex;
    justify-content: center;
    padding: 40px 0;
  }
}

// 响应式
@media (max-width: 768px) {
  .category-section {
    margin-bottom: 32px;

    .category-tabs {
      padding: 16px;
      gap: 8px;

      :deep(.n-button) {
        padding: 0 16px;
        font-size: 14px;
      }
    }
  }

  .main-content {
    .posts-grid {
      grid-template-columns: 1fr;
      gap: 24px;
      margin-bottom: 32px;
    }

    .results-info {
      padding: 12px 16px;
      font-size: 13px;
    }

    .pagination-wrapper {
      padding: 32px 0;
    }
  }
}
</style>
