<script setup lang="ts">
import { fetchListCategory } from '@/api/composables/category'

const { t } = useI18n()

useHead({ title: t('page.categories.categories') })
const localePath = useLocalePath()

const categories = ref<any[]>([])
const loading = ref(true)

onMounted(async () => {
  try {
    const res = await fetchListCategory({
      formValues: { status: 'CATEGORY_STATUS_ACTIVE' },
      fieldMask: 'id,status,sort_order,icon,code,post_count,direct_post_count,parent_id,created_at,children,translations.id,translations.category_id,translations.name,translations.language_code,translations.description,translations.thumbnail,translations.cover_image',
      orderBy: ['-sortOrder'],
    })
    categories.value = res?.items || []
  } catch (e) {
    console.error('[Category List] 加载失败:', e)
  } finally {
    loading.value = false
  }
})

function handleCategoryClick(id: number) {
  navigateTo(localePath(`/category/${id}`))
}
</script>

<template>
  <div class="w-full">
    <LayoutPageHero
      :title="t('page.categories.categories')"
      :description="t('page.categories.browse_all')"
      icon="carbon:folder"
      size="md"
    />

    <LayoutSectionContainer>
      <!-- Loading Skeleton -->
      <div v-if="loading" class="grid grid-cols-[repeat(auto-fill,minmax(280px,1fr))] gap-6">
        <div v-for="i in 6" :key="i">
          <UiSkeleton class="h-40 w-full" />
        </div>
      </div>

      <template v-else>
        <!-- Category Tree -->
        <CategoryTree
          v-if="categories.length > 0"
          :categories="categories"
          :on-category-click="handleCategoryClick"
        />

        <!-- Empty -->
        <UiAppEmpty
          v-else
          :description="t('page.categories.no_categories')"
          in-container
        />
      </template>
    </LayoutSectionContainer>
  </div>
</template>
