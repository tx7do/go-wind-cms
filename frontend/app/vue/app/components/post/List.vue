<script setup lang="ts">
import { fetchListPost } from '@/api/composables/post'
import { useGetCounts, extractCount } from '@/api/composables/interaction'

const props = withDefaults(defineProps<{
  queryParams?: object
  fieldMask?: string
  orderBy?: string[]
  page?: number
  pageSize?: number
  initialPageSize?: number
  showSkeleton?: boolean
  from?: string
  categoryId?: number
  tagId?: number
  columns?: number
  showPagination?: boolean
}>(), {
  queryParams: () => ({}),
  page: 1,
  initialPageSize: 10,
  showSkeleton: true,
  from: 'post-list',
  columns: 3,
  showPagination: false,
})

const { t } = useI18n()

const posts = ref<any[]>([])
const loading = ref(false)
const total = ref(0)
const currentPage = ref(props.page)
const currentPageSize = ref(props.pageSize ?? props.initialPageSize)

async function fetchPosts(page: number, pageSize: number) {
  loading.value = true
  try {
    const res = await fetchListPost({
      paging: { page, pageSize },
      formValues: {
        ...props.queryParams,
        ...(props.categoryId && { category_ids__in: [props.categoryId] }),
        ...(props.tagId && { tag_ids__in: [props.tagId] }),
        status: 'POST_STATUS_PUBLISHED',
      },
      fieldMask: props.fieldMask,
      orderBy: props.orderBy,
    })
    posts.value = res?.items || []
    total.value = res?.total || 0
  } catch (error) {
    console.error('PostList fetch failed:', error)
    posts.value = []
  } finally {
    loading.value = false
  }
}

// 收集当前页 post 的 ID，用于批量查询点赞计数。
// 计数由 interaction_counter 表提供（post.likes 列已移除）。
const postIds = computed(() =>
  (posts.value ?? [])
    .map(p => p.id)
    .filter((id): id is number => typeof id === 'number'),
)
const countsResp = useGetCounts('TARGET_TYPE_POST', postIds, ['COUNTER_METRIC_LIKE'])

// 将 post 与其点赞计数配对，供渲染使用。
// extractCount 对未出现在响应中的 (target, metric) 返回 0。
const postsWithLikeCounts = computed(() =>
  (posts.value ?? []).map(p => ({
    post: p,
    likeCount: extractCount(countsResp.data?.value, p.id, 'COUNTER_METRIC_LIKE'),
  })),
)

function handlePageChange(newPage: number) {
  currentPage.value = newPage
}

watch([currentPage, currentPageSize], ([page, size]) => {
  fetchPosts(page, size)
})

watch([() => props.categoryId, () => props.tagId], () => {
  currentPage.value = 1
  fetchPosts(1, currentPageSize.value)
})

onMounted(() => {
  fetchPosts(currentPage.value, currentPageSize.value)
})
</script>

<template>
  <div class="w-full">
    <!-- Loading Skeleton -->
    <div
      v-if="loading && showSkeleton"
      class="grid gap-6 max-md:!grid-cols-1 max-md:gap-4"
      :style="{ gridTemplateColumns: `repeat(${columns}, 1fr)` }"
    >
      <div v-for="i in currentPageSize" :key="i" class="overflow-hidden rounded-2xl border border-border bg-card">
        <UiSkeleton class="h-60 w-full" />
        <div class="p-6">
          <UiSkeleton class="h-4 w-full" />
          <UiSkeleton class="h-4 w-3/4" />
          <div class="mt-4 flex gap-3 border-t border-border pt-4">
            <UiSkeleton class="h-6 w-[60px]" />
            <UiSkeleton class="h-6 w-[60px]" />
            <UiSkeleton class="h-6 w-[60px]" />
          </div>
        </div>
      </div>
    </div>

    <!-- Loaded Content -->
    <template v-if="!loading && posts.length > 0">
      <div
        v-if="showPagination && total > currentPageSize"
        class="mb-6 rounded-xl border border-border bg-card px-5 py-4 text-sm text-muted-foreground"
      >
        <span>{{ t('page.posts.total_articles', { total }) }}</span>
      </div>

      <div
        class="grid gap-6 max-md:!grid-cols-1 max-md:gap-4"
        :style="{ gridTemplateColumns: `repeat(${columns}, 1fr)` }"
      >
        <PostCard
          v-for="({ post, likeCount }, index) in postsWithLikeCounts"
          :key="`${post.id}-${index}`"
          :post="post"
          :from="from"
          :category-id="categoryId"
          :like-count="likeCount"
        />
      </div>

      <div v-if="showPagination && total > currentPageSize" class="flex justify-center py-10">
        <div class="flex items-center justify-center gap-2">
          <UiButton
            variant="outline"
            size="sm"
            :disabled="currentPage <= 1"
            @click="handlePageChange(currentPage - 1)"
          >
            Previous
          </UiButton>
          <span class="px-2 text-sm text-muted-foreground">
            {{ currentPage }} / {{ Math.ceil(total / currentPageSize) }}
          </span>
          <UiButton
            variant="outline"
            size="sm"
            :disabled="currentPage >= Math.ceil(total / currentPageSize)"
            @click="handlePageChange(currentPage + 1)"
          >
            Next
          </UiButton>
        </div>
      </div>
    </template>

    <!-- Empty State -->
    <div v-if="!loading && posts.length === 0">
      <slot name="empty">
        <UiAppEmpty :description="t('page.posts.no_results')" />
      </slot>
    </div>
  </div>
</template>
