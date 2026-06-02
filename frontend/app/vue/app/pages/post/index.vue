<script setup lang="ts">
const { t } = useI18n()

useHead({ title: t('page.posts.posts_list') })

const selectedCategoryId = ref<number | null>(null)

const handleCategoryChange = (categoryId: number | null) => {
  selectedCategoryId.value = categoryId
}

const queryParams = computed(() => {
  if (selectedCategoryId.value) {
    return { category_ids__in: [selectedCategoryId.value] }
  }
  return {}
})
</script>

<template>
  <div class="w-full">
    <LayoutPageHero
      :title="t('page.posts.posts_list')"
      :description="t('page.posts.explore_latest')"
      icon="carbon:document"
      size="md"
    />

    <LayoutSectionContainer>
      <CategoryFilter
        :selected-category="selectedCategoryId"
        :tree-mode="true"
        :auto-load="true"
        @category-change="handleCategoryChange"
      />

      <PostList
        :key="selectedCategoryId || 'all'"
        :query-params="queryParams"
        :initial-page-size="12"
      />
    </LayoutSectionContainer>
  </div>
</template>
