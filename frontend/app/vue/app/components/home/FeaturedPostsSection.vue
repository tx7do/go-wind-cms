<script setup lang="ts">
import { XIcon } from '@/plugins/xicon'
import { fetchListPost, getPostTitle, getPostSummary, getPostThumbnail } from '@/api/composables/post'

const { t } = useI18n()
const localePath = useLocalePath()

const posts = ref<any[]>([])
const loading = ref(true)

const scrollToCategories = () => {
  const section = document.querySelector('.categories-section')
  if (section) {
    section.scrollIntoView({ behavior: 'smooth', block: 'start' })
  }
}

onMounted(async () => {
  try {
    const res = await fetchListPost({
      paging: { page: 1, pageSize: 3 },
      formValues: { isFeatured: true },
      isTenantUser: true,
    }) as any
    posts.value = res?.items || []
  } catch (e) {
    console.error('[FeaturedPosts] 加载失败:', e)
  } finally {
    loading.value = false
  }
})
</script>

<template>
  <section class="w-full max-w-300 mx-auto scroll-reveal px-8 py-12 max-md:px-4">
    <div class="mb-8 flex items-center justify-between">
      <h2 class="flex items-center gap-2 text-2xl font-extrabold tracking-tight text-foreground max-md:text-xl">
        <XIcon icon="carbon:star-filled" :size="28" class="mr-2 text-primary" />
        {{ t('page.home.featured_posts') }}
      </h2>
      <UiButton variant="ghost" @click="navigateTo(localePath('/post'))">
        {{ t('page.home.view_all') }} →
      </UiButton>
    </div>
    <div class="w-full">
      <!-- Loading -->
      <div v-if="loading" class="grid grid-cols-1 gap-6 md:grid-cols-2 lg:grid-cols-3">
        <div v-for="i in 3" :key="i" class="flex min-h-50 flex-col rounded-2xl border border-border bg-card p-6">
          <UiSkeleton class="h-35 w-full" />
        </div>
      </div>

      <!-- Empty -->
      <div v-else-if="posts.length === 0" class="py-16 text-center">
        <p class="text-muted-foreground">{{ t('page.posts.no_results') }}</p>
      </div>

      <!-- Posts Grid -->
      <div v-else class="grid grid-cols-1 gap-6 md:grid-cols-2 lg:grid-cols-3">
        <NuxtLink
          v-for="post in posts"
          :key="post.id"
          :to="localePath(`/post/${post.id}`)"
          class="group flex min-h-50 flex-col rounded-2xl border border-border bg-card p-6 transition-colors hover:border-primary/30 hover:bg-primary/5"
        >
          <div class="aspect-video overflow-hidden rounded-lg bg-muted">
            <UiImage
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
    </div>
    <div class="mt-8 flex w-full items-center justify-center max-md:hidden">
      <div class="flex items-center gap-4">
        <span class="text-2xl font-extrabold leading-tight tracking-tight text-primary max-md:text-xl">
          {{ t('page.home.explore_more_categories') }}
        </span>
        <div
          class="group flex h-9 w-9 cursor-pointer items-center justify-center rounded-full border border-primary/20 bg-primary/10 transition-all duration-300 hover:border-primary hover:bg-primary hover:shadow-[0_0_20px_hsl(var(--primary)/0.3)]"
          @click="scrollToCategories"
        >
          <XIcon icon="carbon:arrow-down" :size="24" class="text-primary transition-colors group-hover:text-white" />
        </div>
      </div>
    </div>
  </section>
</template>
