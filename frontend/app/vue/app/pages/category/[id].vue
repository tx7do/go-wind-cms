<script setup lang="ts">
import {XIcon} from '@/plugins/xicon'
import {fetchCategory, getCategoryName, getCategoryDescription} from '@/api/composables/category'
import {fetchListPost, getPostTitle, getPostSummary, getPostThumbnail} from '@/api/composables/post'

const {t} = useI18n()
const localePath = useLocalePath()
const route = useRoute()

const category = ref<any>(null)
const childCategories = ref<any[]>([])
const posts = ref<any[]>([])
const loading = ref(true)
const error = ref('')

// Pagination
const page = ref(1)
const pageSize = ref(10)
const total = ref(0)

const categoryId = computed(() => Number(route.params.id))

// 动态标题
const categoryTitle = computed(() =>
    category.value ? getCategoryName(category.value, t('page.categories.category_untitled')) : t('page.categories.categories')
)
useHead({title: () => categoryTitle.value})

const parentCategoryId = computed(() => {
  if (!category.value?.parentId) return null
  return category.value.parentId
})

function getIconName(icon: string | undefined) {
  return icon ? (icon.includes(':') ? icon : `carbon:${icon}`) : 'carbon:folder'
}

async function loadCategory() {
  if (!categoryId.value) return
  loading.value = true
  try {
    const catData = await fetchCategory(categoryId.value) as any
    category.value = catData
    if (catData?.children?.length > 0) {
      childCategories.value = [...catData.children]
    } else {
      childCategories.value = []
    }
  } catch (e) {
    console.error('[Category Detail] 加载失败:', e)
    error.value = t('page.categories.load_failed')
  } finally {
    loading.value = false
  }
}

async function loadPosts() {
  if (!categoryId.value) return
  try {
    const res = await fetchListPost({
      paging: {page: page.value, pageSize: pageSize.value},
      formValues: {categoryId: categoryId.value},
      isTenantUser: true,
    }) as any
    posts.value = res?.items || []
    total.value = res?.total || 0
  } catch (e) {
    console.error('[Category Posts] 加载失败:', e)
  }
}

function handlePageChange(newPage: number) {
  page.value = newPage
  loadPosts()
}

function handleBackToParent() {
  if (parentCategoryId.value) {
    navigateTo(localePath(`/category/${parentCategoryId.value}`))
  } else {
    navigateTo(localePath('/category'))
  }
}

onMounted(async () => {
  await loadCategory()
  await loadPosts()
})
</script>

<template>
  <div class="w-full">
    <!-- Loading -->
    <div v-if="loading" class="w-full py-20">
      <LayoutSectionContainer>
        <div class="mx-auto max-w-3xl space-y-4">
          <UiSkeleton class="h-10 w-1/2"/>
          <UiSkeleton class="h-5 w-3/4"/>
        </div>
      </LayoutSectionContainer>
    </div>

    <!-- Error -->
    <div v-else-if="error" class="w-full py-20 text-center">
      <LayoutSectionContainer>
        <div class="flex flex-col items-center gap-4">
          <XIcon icon="carbon:warning-alt" :size="48" class="text-muted-foreground"/>
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
          :description="getCategoryDescription(category) || undefined"
          :icon="getIconName(category.icon)"
          size="lg"
      >
        <template #default>
          <div class="flex items-center gap-4">
            <div class="flex items-center gap-2">
              <XIcon icon="carbon:document" :size="16"/>
              <span>{{ category.postCount || 0 }} {{ t('page.posts.articles') }}</span>
            </div>
            <div v-if="childCategories.length > 0" class="flex items-center gap-2">
              <XIcon icon="carbon:folder" :size="16"/>
              <span>{{ childCategories.length }} {{ t('page.categories.categories') }}</span>
            </div>
          </div>
        </template>
      </LayoutPageHero>

      <LayoutSectionContainer>
        <!-- Back Button -->
        <div class="mb-8">
          <UiButton variant="ghost" size="sm" class="gap-2" @click="handleBackToParent">
            <XIcon icon="carbon:arrow-left" :size="16"/>
            {{ parentCategoryId ? t('page.categories.back_to_parent') : t('page.categories.back_to_list') }}
          </UiButton>
        </div>

        <!-- Sub Categories -->
        <div v-if="childCategories.length > 0" class="mb-8">
          <div class="grid grid-cols-[repeat(auto-fill,minmax(200px,1fr))] gap-4">
            <div
                v-for="child in childCategories"
                :key="child.id"
                class="group flex cursor-pointer items-center gap-3 rounded-xl border border-border bg-card p-4 transition-all duration-300 hover:-translate-y-0.5 hover:border-primary/30 hover:bg-primary/5"
                @click="navigateTo(localePath(`/category/${child.id}`))"
            >
              <div class="flex h-10 w-10 shrink-0 items-center justify-center rounded-lg bg-primary/10 text-primary">
                <XIcon :icon="getIconName(child.icon)" :size="24"/>
              </div>
              <div class="min-w-0 flex-1">
                <h4 class="truncate text-sm font-semibold text-foreground group-hover:text-primary">
                  {{ getCategoryName(child, t('page.categories.category_untitled')) }}
                </h4>
                <p class="text-xs text-muted-foreground">
                  {{ child.postCount || 0 }} {{ t('page.posts.articles') }}
                </p>
              </div>
            </div>
          </div>
        </div>

        <!-- Empty -->
        <div v-if="posts.length === 0" class="py-20 text-center">
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
  </div>
</template>
