<script setup lang="ts">
import type {contentservicev1_Category} from "@/api/generated/app/service/v1";
import CategoryCard from "@/components/CategoryCard";

interface Props {
  categories: contentservicev1_Category[]
  loading?: boolean
  showSkeleton?: boolean
  columns?: number
  gap?: number
}

withDefaults(defineProps<Props>(), {
  loading: false,
  showSkeleton: true,
  columns: 3,
  gap: 20
})

const emit = defineEmits<{
  categoryClick: [id: number]
}>()

function handleCategoryClick(id: number) {
  emit('categoryClick', id)
}
</script>

<template>
  <div class="category-list-container">
    <!-- Loading Skeleton -->
    <div v-if="loading && showSkeleton" class="category-grid" :style="{
      gridTemplateColumns: `repeat(auto-fill, minmax(${280}px, 1fr))`,
      gap: `${gap}px`
    }">
      <div v-for="i in columns" :key="i" class="category-card-skeleton">
        <n-skeleton height="160px"/>
        <div style="padding: 16px;">
          <n-skeleton :width="'80%'" size="medium" style="margin-bottom: 8px;"/>
          <n-skeleton :width="'90%'" size="small" style="margin-bottom: 12px;"/>
          <n-skeleton :width="60" size="small"/>
        </div>
      </div>
    </div>

    <!-- Loaded Content -->
    <div v-else-if="!loading && categories.length > 0" class="category-grid" :style="{
      gridTemplateColumns: `repeat(auto-fill, minmax(${280}px, 1fr))`,
      gap: `${gap}px`
    }">
      <CategoryCard
        v-for="category in categories"
        :key="category.id"
        :category="category"
        @click="handleCategoryClick(category.id)"
      />
    </div>

    <!-- Empty State -->
    <div v-else-if="!loading && categories.length === 0" class="empty-state">
      <n-empty :description="$t('page.categories.no_categories')"/>
    </div>
  </div>
</template>

<style scoped lang="less">
.category-list-container {
  width: 100%;
}

.category-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(280px, 1fr));
  gap: 20px;
}

.category-card-skeleton {
  background: var(--color-bg);
  border-radius: 12px;
  overflow: hidden;
  border: 1px solid var(--color-border);
}

.empty-state {
  padding: 60px 20px;
  text-align: center;
}

// Responsive Design
@media (max-width: 1024px) {
  .category-grid {
    grid-template-columns: repeat(auto-fill, minmax(260px, 1fr));
    gap: 16px;
  }
}

@media (max-width: 768px) {
  .category-grid {
    grid-template-columns: repeat(auto-fill, minmax(240px, 1fr));
    gap: 14px;
  }
}

@media (max-width: 640px) {
  .category-grid {
    grid-template-columns: 1fr;
    gap: 12px;
  }
}
</style>
