<script setup lang="ts">
import { cn } from '@/lib/utils'

const props = withDefaults(defineProps<{
  categories: any[]
  loading?: boolean
  showSkeleton?: boolean
  columns?: number
  gap?: number
  onCategoryClick?: (id: number) => void
}>(), {
  loading: false,
  showSkeleton: true,
  columns: 3,
  gap: 20,
})

function handleCategoryClick(id: number) {
  props.onCategoryClick?.(id)
}
</script>

<template>
  <div class="w-full">
    <!-- Loading Skeleton -->
    <div
      v-if="loading && showSkeleton"
      class="grid max-md:grid-cols-1"
      :style="{ gridTemplateColumns: `repeat(${columns}, 1fr)`, gap: `${gap}px` }"
    >
      <div v-for="index in columns" :key="index" class="overflow-hidden rounded-xl border border-border bg-background">
        <UiSkeleton class="h-40 w-full" />
        <div class="p-4">
          <UiSkeleton class="h-4 w-full" />
          <UiSkeleton class="h-4 w-3/4" />
          <UiSkeleton class="h-6 w-[60px]" />
        </div>
      </div>
    </div>

    <!-- Loaded Content -->
    <div
      v-if="!loading && categories.length > 0"
      class="grid grid-cols-[repeat(auto-fill,minmax(280px,1fr))] max-md:grid-cols-1"
      :style="{ gap: `${gap}px` }"
    >
      <CategoryCard
        v-for="category in categories"
        :key="category.id"
        :category="category"
        @click="handleCategoryClick"
      />
    </div>

    <!-- Empty State -->
    <div v-if="!loading && categories.length === 0">
      <UiAppEmpty :description="$t('page.categories.no_categories')" />
    </div>
  </div>
</template>
