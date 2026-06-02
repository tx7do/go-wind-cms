<script setup lang="ts">
import { XIcon } from '@/plugins/xicon'
import { fetchCategory, getCategoryName, getCategoryDescription } from '@/api/composables/category'

const { t } = useI18n()
const localePath = useLocalePath()
const route = useRoute()

const category = ref<any>(null)
const childCategories = ref<any[]>([])
const loading = ref(true)
const error = ref('')

const categoryId = computed(() => Number(route.params.id))

// 动态标题
const categoryTitle = computed(() =>
  category.value ? getCategoryName(category.value, t('page.categories.category_untitled')) : t('page.categories.categories')
)
useHead({ title: () => categoryTitle.value })

const parentCategoryId = computed(() => {
  if (!category.value?.parentId) return null
  return category.value.parentId
})

async function loadCategory() {
  if (!categoryId.value) return
  loading.value = true
  try {
    const catData = await fetchCategory(
      categoryId.value,
      'id,status,sort_order,icon,code,post_count,direct_post_count,parent_id,created_at,children,translations.id,translations.category_id,translations.name,translations.language_code,translations.description,translations.thumbnail,translations.cover_image'
    ) as any
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

function handleViewChildCategory(id: number) {
  navigateTo(localePath(`/category/${id}`))
}

function handleBackToParent() {
  if (parentCategoryId.value) {
    navigateTo(localePath(`/category/${parentCategoryId.value}`))
  } else {
    navigateTo(localePath('/category'))
  }
}

onMounted(loadCategory)
</script>

<template>
  <div class="w-full">
    <!-- Error -->
    <div v-if="error" class="w-full py-20 text-center">
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
        :description="getCategoryDescription(category) || undefined"
        :icon="category.icon || 'carbon:folder'"
        size="lg"
      >
        <template #default>
          <div class="flex items-center gap-2">
            <XIcon icon="carbon:document" :size="16" />
            <span>{{ category.postCount || 0 }} {{ t('page.posts.articles') }}</span>
          </div>
          <div v-if="childCategories.length > 0" class="flex items-center gap-2">
            <XIcon icon="carbon:folder" :size="16" />
            <span>{{ childCategories.length }} {{ t('page.categories.categories') }}</span>
          </div>
        </template>
      </LayoutPageHero>

      <LayoutSectionContainer>
        <!-- Back Button -->
        <div class="mb-8">
          <LayoutBackButton
            :label="parentCategoryId ? t('page.categories.back_to_parent') : t('page.categories.back_to_list')"
            :on-click="handleBackToParent"
          />
        </div>

        <!-- Sub Categories List -->
        <div v-if="childCategories.length > 0" class="mb-8">
          <CategoryList
            :categories="childCategories"
            :loading="false"
            :show-skeleton="false"
            :on-category-click="handleViewChildCategory"
          />
        </div>

        <!-- Posts List with Pagination -->
        <PostList
          :key="categoryId"
          :initial-page-size="10"
          :category-id="categoryId"
          from="category"
        />
      </LayoutSectionContainer>
    </template>
  </div>
</template>
