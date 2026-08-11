<script setup lang="ts">
import { fetchSearchPosts } from '@/api/composables/post'

const { t } = useI18n()
const route = useRoute()
const localePath = useLocalePath()

const query = computed(() => (route.query.q as string)?.trim() ?? '')

useHead({ title: () => query.value ? `${t('page.posts.search')}: ${query.value}` : t('page.posts.search') })

// 依赖 query 触发重新搜索
const { data, pending, error } = await useAsyncData(
  () => `search-${query.value}`,
  () => query.value ? fetchSearchPosts({ query: query.value, pageSize: 20 }) : Promise.resolve(null),
  { watch: [query] },
)

const hits = computed(() => data.value?.items ?? [])
const total = computed(() => data.value?.total ?? 0)

function openPost(postId?: number) {
  if (!postId) return
  navigateTo(localePath(`/post/${postId}`))
  window.scrollTo({ top: 0, behavior: 'smooth' })
}
</script>

<template>
  <div class="w-full">
    <LayoutPageHero
      :title="t('page.posts.search')"
      :description="query ? `"${query}"` : ''"
      icon="lucide:search"
      size="md"
    />

    <LayoutSectionContainer>
      <!-- 加载中 -->
      <div v-if="pending" class="flex justify-center py-20">
        <UiSpinner size="lg" class="text-primary" />
      </div>

      <!-- 出错 -->
      <div v-else-if="error" class="py-20 text-center text-muted-foreground">
        {{ t('page.posts.search_failed') }}
      </div>

      <!-- 无关键词 -->
      <div v-else-if="!query" class="py-20 text-center text-muted-foreground">
        {{ t('page.posts.search_placeholder') }}
      </div>

      <!-- 无结果 -->
      <div v-else-if="hits.length === 0" class="flex flex-col items-center py-20 text-muted-foreground">
        <XIcon icon="lucide:search-x" :size="48" class="mb-4 opacity-50" />
        <p>{{ t('page.posts.no_search_results', { query }) }}</p>
      </div>

      <!-- 结果列表 -->
      <div v-else>
        <p class="mb-6 text-sm text-muted-foreground">
          {{ t('page.posts.search_results_count', { count: total }) }}
        </p>
        <div class="grid gap-3">
          <article
            v-for="hit in hits"
            :key="hit.postId"
            class="group flex cursor-pointer items-center gap-3 rounded-xl border border-border bg-card p-5 shadow-sm transition-all duration-300 hover:-translate-y-0.5 hover:border-primary/40 hover:shadow-md"
            @click="openPost(hit.postId)"
          >
            <XIcon icon="carbon:document" :size="20" class="shrink-0 text-muted-foreground transition-colors group-hover:text-primary" />
            <h3 class="line-clamp-2 flex-1 text-base font-semibold text-foreground transition-colors group-hover:text-primary">
              {{ hit.title }}
            </h3>
            <XIcon icon="lucide:chevron-right" :size="18" class="shrink-0 text-muted-foreground opacity-0 transition-opacity group-hover:opacity-100" />
          </article>
        </div>
      </div>
    </LayoutSectionContainer>
  </div>
</template>
