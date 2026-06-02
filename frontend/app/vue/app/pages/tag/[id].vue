<script setup lang="ts">
import { fetchTag, getTagTranslation } from '@/api/composables/tag'
import { fetchListPost, getPostTitle, getPostSummary, getPostThumbnail } from '@/api/composables/post'
import { XIcon } from '@/plugins/xicon'

const { t } = useI18n()
const localePath = useLocalePath()
const route = useRoute()

const tag = ref<any>(null)
const loading = ref(true)
const posts = ref<any[]>([])
const postsLoading = ref(true)

// Pagination
const page = ref(1)
const pageSize = ref(10)
const total = ref(0)

const tagId = computed(() => Number(route.params.id))

function getTagName(tagData: any) {
  const translation = getTagTranslation(tagData)
  return translation?.name || t('page.tags.tag_untitled')
}

function getTagDescription(tagData: any) {
  const translation = getTagTranslation(tagData)
  return translation?.description || ''
}

async function loadTag() {
  if (!tagId.value) return
  loading.value = true
  try {
    tag.value = await fetchTag(tagId.value) as any
  } catch (e) {
    console.error('[Tag Detail] 加载失败:', e)
  } finally {
    loading.value = false
  }
}

async function loadPosts() {
  if (!tagId.value) return
  postsLoading.value = true
  try {
    const res = await fetchListPost({
      paging: { page: page.value, pageSize: pageSize.value },
      formValues: { tagIds: [tagId.value] },
      isTenantUser: true,
    }) as any
    posts.value = res?.items || []
    total.value = res?.total || 0
  } catch (e) {
    console.error('[Tag Posts] 加载失败:', e)
  } finally {
    postsLoading.value = false
  }
}

function handlePageChange(newPage: number) {
  page.value = newPage
  loadPosts()
}

onMounted(async () => {
  await loadTag()
  await loadPosts()
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

    <!-- Tag Detail -->
    <template v-else-if="tag">
      <LayoutPageHero
        :title="getTagName(tag)"
        :description="getTagDescription(tag) || undefined"
        :icon="`carbon:${tag.icon || 'tag'}`"
        :icon-size="56"
        size="lg"
        :accent-color="tag.color || undefined"
      >
        <template #default>
          <div class="flex items-center gap-2">
            <XIcon icon="carbon:document" :size="16" />
            <span>{{ tag.postCount || 0 }} {{ t('page.posts.articles') }}</span>
          </div>
        </template>
      </LayoutPageHero>

      <LayoutSectionContainer>
        <!-- Back Button -->
        <div class="mb-8">
          <UiButton variant="ghost" size="sm" class="gap-2" @click="navigateTo(localePath('/tag'))">
            <XIcon icon="carbon:arrow-left" :size="16" />
            {{ t('page.categories.back_to_list') }}
          </UiButton>
        </div>

        <!-- Posts Loading -->
        <div v-if="postsLoading" class="grid gap-6 sm:grid-cols-2 lg:grid-cols-3">
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
          <UiAppEmpty>
            {{ t('page.posts.no_results') }}
          </UiAppEmpty>
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

        <!-- Pagination -->
        <div v-if="total > pageSize" class="mt-8 flex items-center justify-center gap-2">
          <UiButton
            variant="outline"
            size="sm"
            :disabled="page <= 1"
            @click="handlePageChange(page - 1)"
          >
            {{ t('common.previous') || 'Previous' }}
          </UiButton>
          <span class="text-sm text-muted-foreground">
            {{ page }} / {{ Math.ceil(total / pageSize) }}
          </span>
          <UiButton
            variant="outline"
            size="sm"
            :disabled="page >= Math.ceil(total / pageSize)"
            @click="handlePageChange(page + 1)"
          >
            {{ t('common.next') || 'Next' }}
          </UiButton>
        </div>
      </LayoutSectionContainer>
    </template>

    <!-- Tag not found -->
    <div v-else class="w-full py-20 text-center">
      <LayoutSectionContainer>
        <div class="flex flex-col items-center gap-4">
          <XIcon icon="carbon:warning-alt" :size="48" class="text-muted-foreground" />
          <p class="text-lg text-muted-foreground">{{ t('page.tags.no_tags') }}</p>
          <UiButton variant="outline" @click="navigateTo(localePath('/tag'))">
            {{ t('common.back') || '← Back' }}
          </UiButton>
        </div>
      </LayoutSectionContainer>
    </div>
  </div>
</template>
