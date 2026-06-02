<script setup lang="ts">
import { fetchListPost, getPostTitle, getPostSummary, getPostThumbnail } from '@/api/composables/post'
import type { ListPostParams } from '@/api/composables/post'

const { t } = useI18n()
const localePath = useLocalePath()
const route = useRoute()

const posts = ref<any[]>([])
const loading = ref(true)

onMounted(async () => {
  try {
    const params: ListPostParams = {
      paging: { page: 1, pageSize: 20 },
      isTenantUser: true,
    }
    const res = await fetchListPost(params)
    posts.value = res?.items || res?.data || []
  } catch (e) {
    console.error('[Post List] 加载失败:', e)
  } finally {
    loading.value = false
  }
})
</script>

<template>
  <div class="w-full">
    <LayoutPageHero
      :title="t('page.posts.posts_list')"
      icon="carbon:document"
      size="md"
    />

    <section class="w-full py-12 max-md:py-8">
      <LayoutSectionContainer>
        <!-- Loading -->
        <div v-if="loading" class="grid gap-6 sm:grid-cols-2 lg:grid-cols-3">
          <div v-for="i in 6" :key="i" class="rounded-xl border border-border p-4">
            <UiSkeleton class="h-35 w-full" />
            <div class="mt-4 space-y-2">
              <UiSkeleton class="h-5 w-3/4" />
              <UiSkeleton class="h-4 w-full" />
            </div>
          </div>
        </div>

        <!-- Empty -->
        <div v-else-if="posts.length === 0" class="py-20 text-center">
          <UiAppEmpty :message="t('page.posts.no_results')" />
        </div>

        <!-- Post Grid -->
        <div v-else class="grid gap-6 sm:grid-cols-2 lg:grid-cols-3">
          <NuxtLink
            v-for="post in posts"
            :key="post.id"
            :to="localePath(`/post/${post.id}`)"
            class="group rounded-xl border border-border bg-card p-4 transition-colors hover:border-primary/30 hover:bg-primary/5"
          >
            <div class="aspect-video overflow-hidden rounded-lg bg-muted">
              <img
                v-if="getPostThumbnail(post)"
                :src="getPostThumbnail(post)"
                :alt="getPostTitle(post)"
                class="h-full w-full object-cover transition-transform group-hover:scale-105"
              />
            </div>
            <h3 class="mt-3 text-base font-semibold text-foreground group-hover:text-primary">
              {{ getPostTitle(post, t('page.post_detail.untitled')) }}
            </h3>
            <p class="mt-1 line-clamp-2 text-sm text-muted-foreground">
              {{ getPostSummary(post) }}
            </p>
          </NuxtLink>
        </div>
      </LayoutSectionContainer>
    </section>
  </div>
</template>
