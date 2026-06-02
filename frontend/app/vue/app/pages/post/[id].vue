<script setup lang="ts">
import { XIcon } from '@/plugins/xicon'
import { fetchPost, getPostTitle, getPostSummary, getPostThumbnail, getPostContent } from '@/api/composables/post'

const { t } = useI18n()
const localePath = useLocalePath()
const route = useRoute()

const post = ref<any>(null)
const loading = ref(true)
const error = ref('')

const postId = computed(() => Number(route.params.id))

onMounted(async () => {
  try {
    if (!postId.value) {
      error.value = t('page.post_detail.not_found')
      return
    }
    post.value = await fetchPost(postId.value)
  } catch (e) {
    console.error('[Post Detail] 加载失败:', e)
    error.value = t('page.post_detail.load_failed')
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
          <UiSkeleton class="h-10 w-3/4" />
          <UiSkeleton class="h-5 w-1/3" />
          <div class="mt-8 space-y-3">
            <UiSkeleton class="h-4 w-full" />
            <UiSkeleton class="h-4 w-full" />
            <UiSkeleton class="h-4 w-2/3" />
          </div>
        </div>
      </LayoutSectionContainer>
    </div>

    <!-- Error -->
    <div v-else-if="error" class="w-full py-20 text-center">
      <LayoutSectionContainer>
        <div class="flex flex-col items-center gap-4">
          <XIcon icon="carbon:warning-alt" :size="48" class="text-muted-foreground" />
          <p class="text-lg text-muted-foreground">{{ error }}</p>
          <UiButton variant="outline" @click="navigateTo(localePath('/post'))">
            {{ t('common.back') || '← Back' }}
          </UiButton>
        </div>
      </LayoutSectionContainer>
    </div>

    <!-- Post Detail -->
    <template v-else-if="post">
      <LayoutPageHero
        :title="getPostTitle(post, t('page.post_detail.untitled'))"
        :subtitle="getPostSummary(post)"
        size="sm"
      />

      <section class="w-full py-12 max-md:py-8">
        <LayoutSectionContainer>
          <article class="mx-auto max-w-3xl">
            <!-- Thumbnail -->
            <div v-if="getPostThumbnail(post)" class="mb-8 overflow-hidden rounded-xl">
              <img
                :src="getPostThumbnail(post)"
                :alt="getPostTitle(post)"
                class="h-auto w-full object-cover"
              />
            </div>

            <!-- Content -->
            <div class="prose prose-neutral dark:prose-invert max-w-none" v-html="getPostContent(post)" />
          </article>
        </LayoutSectionContainer>
      </section>
    </template>
  </div>
</template>
