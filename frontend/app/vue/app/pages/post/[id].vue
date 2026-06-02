<script setup lang="ts">
import { XIcon } from '@/plugins/xicon'
import { fetchPost, getPostTitle, getPostContent, getPostThumbnail } from '@/api/composables/post'

const { t } = useI18n()
const localePath = useLocalePath()
const route = useRoute()

const post = ref<any>(null)
const postLoading = ref(true)
const error = ref('')
const isLiked = ref(false)
const isBookmarked = ref(false)
const isTocExpanded = ref(true)
const contentRef = ref<HTMLElement | null>(null)

const postId = computed(() => Number(route.params.id))

const displayTitle = computed(() => post.value ? getPostTitle(post.value) : '')
const displayContent = computed(() => post.value ? getPostContent(post.value) : '')
const displayThumbnail = computed(() => post.value ? getPostThumbnail(post.value) : '')

const relatedPostsQuery = computed(() => {
  if (!post.value?.categoryIds) return null
  return {
    status: 'POST_STATUS_PUBLISHED',
    id__not: postId.value,
    category_ids__in: post.value.categoryIds,
  }
})

async function loadPost() {
  if (!postId.value) return
  postLoading.value = true
  try {
    post.value = await fetchPost(postId.value)
    if (post.value) {
      document.title = `${getPostTitle(post.value)} - GoWind Content Hub`
    }
  } catch (e) {
    console.error('[Post Detail] Load failed:', e)
    error.value = t('page.post_detail.load_failed')
  } finally {
    postLoading.value = false
  }
}

const isLoading = computed(() => postLoading.value && !post.value)

function handleBack() {
  const from = route.query.from as string
  const categoryId = route.query.categoryId as string

  if (from === 'category' && categoryId) {
    navigateTo(localePath(`/category/${categoryId}`))
    return
  }
  if (from === 'tag') { navigateTo(localePath('/tag')); return }
  if (from === 'user') { navigateTo(localePath('/user')); return }
  if (from === 'home') { navigateTo(localePath('/')); return }
  if (from === 'post-list') { navigateTo(localePath('/post')); return }

  if (import.meta.client && window.history.length > 2) {
    router.back()
    return
  }
  navigateTo(localePath('/post'))
}

function handleShare() {
  const url = window.location.href
  if (navigator.share) {
    navigator.share({ title: displayTitle.value, url }).catch(() => {
      navigator.clipboard.writeText(url)
    })
  } else {
    navigator.clipboard.writeText(url)
  }
}

onMounted(() => {
  loadPost()
})
</script>

<template>
  <div class="w-full">
    <!-- Loading -->
    <div v-if="isLoading" class="w-full">
      <div class="w-full max-w-[1200px] mx-auto px-8 py-8 max-md:px-4">
        <div class="mb-6 h-9 w-32 animate-pulse rounded bg-muted" />
        <article class="overflow-hidden rounded-xl border border-border bg-card shadow-sm">
          <div class="h-[300px] animate-pulse bg-muted max-md:h-[200px]" />
          <div class="flex gap-6 p-8 max-md:flex-col max-md:p-4">
            <aside class="w-[240px] shrink-0 max-md:hidden">
              <div class="sticky top-24 space-y-3 rounded-lg border border-border bg-background p-4">
                <div class="h-6 w-[200px] animate-pulse rounded bg-muted" />
                <div class="mt-4 h-5 w-[180px] animate-pulse rounded bg-muted" />
                <div class="mt-2 h-5 w-[160px] animate-pulse rounded bg-muted" />
                <div class="mt-2 h-5 w-[140px] animate-pulse rounded bg-muted" />
              </div>
            </aside>
            <div class="flex-1">
              <div class="mb-4 h-12 w-[80%] animate-pulse rounded bg-muted" />
              <div class="mb-4 h-8 w-[60%] animate-pulse rounded bg-muted" />
              <div class="mb-6 flex gap-4">
                <div class="h-5 w-[100px] animate-pulse rounded bg-muted" />
                <div class="h-5 w-[100px] animate-pulse rounded bg-muted" />
                <div class="h-5 w-[100px] animate-pulse rounded bg-muted" />
                <div class="h-5 w-[100px] animate-pulse rounded bg-muted" />
              </div>
              <div class="space-y-3">
                <div class="h-4 w-full animate-pulse rounded bg-muted" />
                <div class="h-4 w-[90%] animate-pulse rounded bg-muted" />
                <div class="h-4 w-[95%] animate-pulse rounded bg-muted" />
              </div>
            </div>
          </div>
        </article>
      </div>
    </div>

    <!-- Error -->
    <div v-else-if="error" class="w-full py-20 text-center">
      <LayoutSectionContainer>
        <div class="flex flex-col items-center gap-4">
          <XIcon icon="carbon:warning-alt" :size="48" class="text-muted-foreground" />
          <p class="text-lg text-muted-foreground">{{ error }}</p>
          <UiButton variant="outline" @click="navigateTo(localePath('/post'))">
            ← {{ t('page.post_detail.back') }}
          </UiButton>
        </div>
      </LayoutSectionContainer>
    </div>

    <!-- Not found -->
    <div v-else-if="!post" class="flex w-full items-center justify-center py-20">
      <div class="text-lg text-muted-foreground">Post not found</div>
    </div>

    <!-- Post Detail -->
    <template v-else>
      <!-- Back Navigation -->
      <LayoutSectionContainer no-padding class="!py-6">
        <LayoutBackButton :label="t('page.post_detail.back')" :on-click="handleBack" />
      </LayoutSectionContainer>

      <!-- Post Article -->
      <article class="w-full max-w-[1200px] mx-auto px-8 max-md:px-4">
        <!-- Thumbnail Banner -->
        <div v-if="displayThumbnail" class="relative mb-8 h-[300px] overflow-hidden rounded-xl max-md:h-[200px]">
          <img :src="displayThumbnail" :alt="displayTitle" class="h-full w-full object-cover" />
          <div class="absolute inset-0 bg-gradient-to-t from-black/30 to-transparent" />
        </div>

        <div :class="[!isTocExpanded ? 'lg:pl-12' : '', 'flex gap-6 max-md:flex-col']">
          <!-- Table of Contents -->
          <PostToc
            :content-ref="contentRef"
            :content-key="displayContent"
            :is-expanded="isTocExpanded"
            :title="t('page.post_detail.table_of_contents')"
            @collapse="isTocExpanded = false"
            @expand="isTocExpanded = true"
          />

          <div class="flex-1 min-w-0">
            <!-- Post Header -->
            <header class="mb-8">
              <h1 class="mb-5 text-3xl font-bold leading-tight text-foreground max-md:text-2xl">
                {{ displayTitle }}
              </h1>
              <PostMetaBar
                :author-name="post.authorName"
                :created-at="post.createdAt"
                :visits="post.visits"
                :likes="post.likes"
              />
            </header>

            <!-- Post Content -->
            <div ref="contentRef" class="mb-8">
              <ContentViewer :content="displayContent" type="markdown" />
            </div>

            <!-- Floating Actions -->
            <PostFloatingActions
              :is-liked="isLiked"
              :is-bookmarked="isBookmarked"
              :labels="{
                likes: t('page.post_detail.likes'),
                bookmark: t('page.post_detail.bookmark'),
                share: t('page.post_detail.share'),
              }"
              @like="isLiked = !isLiked"
              @bookmark="isBookmarked = !isBookmarked"
              @share="handleShare"
            />
          </div>
        </div>
      </article>

      <!-- Comments Section -->
      <LayoutSectionContainer no-padding>
        <CommentSection
          :object-id="postId"
          content-type="CONTENT_TYPE_POST"
        />
      </LayoutSectionContainer>

      <!-- Related Posts -->
      <LayoutSectionContainer>
        <div class="mb-6">
          <h2 class="flex items-center gap-2 text-2xl font-bold text-foreground">
            <XIcon icon="carbon:book" />
            <span>{{ t('page.post_detail.related_posts') }}</span>
          </h2>
        </div>
        <PostList
          v-if="relatedPostsQuery"
          :query-params="relatedPostsQuery"
          field-mask="id,status,sort_order,is_featured,visits,likes,comment_count,author_name,available_languages,created_at,translations.id,translations.post_id,translations.language_code,translations.title,translations.summary,translations.thumbnail"
          :order-by="['-sortOrder']"
          :page="1"
          :page-size="3"
          :show-skeleton="false"
        />
      </LayoutSectionContainer>

      <LayoutBackToTop :scroll-threshold="500" />
    </template>
  </div>
</template>
