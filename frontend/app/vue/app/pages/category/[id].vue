<script setup lang="ts">
import { XIcon } from '@/plugins/xicon'
import { fetchCategory, getCategoryName, getCategoryDescription, getCategoryThumbnail } from '@/api/composables/category'
import { fetchListPost, getPostTitle, getPostSummary, getPostThumbnail } from '@/api/composables/post'
import type { ListPostParams } from '@/api/composables/post'

const { t } = useI18n()
const localePath = useLocalePath()
const route = useRoute()

const category = ref<any>(null)
const posts = ref<any[]>([])
const loading = ref(true)
const error = ref('')

const categoryId = computed(() => Number(route.params.id))

onMounted(async () => {
  try {
    if (!categoryId.value) {
      error.value = t('page.categories.not_found')
      return
    }
    const [catRes, postRes] = await Promise.all([
      fetchCategory(categoryId.value),
      fetchListPost({
        paging: { page: 1, pageSize: 20 },
        formValues: { categoryId: categoryId.value },
        isTenantUser: true,
      }),
    ])
    category.value = catRes
    posts.value = (postRes as any)?.items || []
  } catch (e) {
    console.error('[Category Detail] 加载失败:', e)
    error.value = t('page.categories.load_failed')
  } finally {
    loading.value = false
  }
})
</script>

<template>
  <div class="w-full">
    <!-- Loading -->
    <div v-if="loading" class="w-full py-20">
      <LayoutSectionContainer>
        <div class="mx-auto max-w-3xl space-y-4">
          <UiSkeleton class="h-10 w-1/2" />
          <UiSkeleton class="h-5 w-3/4" />
        </div>
      </LayoutSectionContainer>
    </div>

    <!-- Error -->
    <div v-else-if="error" class="w-full py-20 text-center">
      <LayoutSectionContainer>
        <div class="flex flex-col items-center gap-4">
          <XIcon icon="carbon:warning-alt" :size="48" class="text-muted-foreground" />
          <p class="text-lg text-muted-foreground">{{ error }}</p>
          <UiButton variant="outline" @click="navigateTo(localePath('/category'))">
            {{ t('common.back') || '← Back' }}
          </UiButton>
        </div>
      </LayoutSectionContainer>
    </div>

    <!-- Category Detail -->
    <template v-else-if="category">
      <LayoutPageHero
        :title="getCategoryName(category, t('page.categories.category_untitled'))"
        :subtitle="getCategoryDescription(category)"
        icon="carbon:folder-details"
        size="sm"
      />

      <section class="w-full py-12 max-md:py-8">
        <LayoutSectionContainer>
          <!-- Empty -->
          <div v-if="posts.length === 0" class="py-20 text-center">
            <UiAppEmpty :message="t('page.posts.no_results')" />
          </div>

          <!-- Posts Grid -->
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
    </template>
  </div>
</template>
