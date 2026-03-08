<script setup lang="ts">
import {$t} from '@/locales'
import {useCategoryStore} from '@/stores/modules/app'

interface Props {
  categories: any[]
  selectedCategory: number | null
}

interface Emits {
  (e: 'category-change', categoryId: number | null): void
}

const props = defineProps<Props>()
const emit = defineEmits<Emits>()
const categoryStore = useCategoryStore()

function handleCategoryChange(categoryId: number | null) {
  emit('category-change', categoryId)
}
</script>

<template>
  <div class="category-filter">
    <div class="category-tabs">
      <n-button
        :type="selectedCategory === null ? 'primary' : 'default'"
        :ghost="selectedCategory !== null"
        size="large"
        @click="handleCategoryChange(null)"
      >
        <template #icon>
          <span class="i-carbon:grid"/>
        </template>
        {{ $t('page.posts.all_categories') }}
      </n-button>
      <n-button
        v-for="cat in categories"
        :key="cat.id"
        :type="selectedCategory === cat.id ? 'primary' : 'default'"
        :ghost="selectedCategory !== cat.id"
        size="large"
        @click="handleCategoryChange(cat.id)"
      >
        <template #icon>
          <span class="i-carbon:folder"/>
        </template>
        {{ categoryStore.getCategoryName(cat) }}
      </n-button>
    </div>
  </div>
</template>

<style scoped lang="less">
.category-filter {
  margin-bottom: 48px;

  .category-tabs {
    display: flex;
    flex-wrap: wrap;
    gap: 12px;
    padding: 24px;
    background: var(--color-surface);
    border-radius: 16px;
    box-shadow: 0 2px 8px rgba(0, 0, 0, 0.04);

    :deep(.n-button) {
      border-radius: 50px;
      padding: 0 24px;
      font-weight: 500;
      transition: all 0.3s;

      &.n-button--primary-type:not(.n-button--ghost) {
        background: linear-gradient(135deg, var(--color-brand) 0%, #764ba2 100%);
        box-shadow: 0 4px 12px rgba(102, 126, 234, 0.3);
      }

      &.n-button--ghost {
        &:hover {
          color: var(--color-brand);
          border-color: var(--color-brand);
        }
      }
    }
  }
}

// 响应式
@media (max-width: 768px) {
  .category-filter {
    margin-bottom: 32px;

    .category-tabs {
      padding: 16px;
      gap: 8px;
      overflow-x: auto;
      overflow-y: hidden;
      -webkit-overflow-scrolling: touch;

      scrollbar-width: thin;
      scrollbar-color: rgba(102, 126, 234, 0.2) transparent;

      &::-webkit-scrollbar {
        height: 3px;
      }

      &::-webkit-scrollbar-track {
        background: transparent;
      }

      &::-webkit-scrollbar-thumb {
        background: rgba(102, 126, 234, 0.2);
        border-radius: 10px;
      }

      :deep(.n-button) {
        flex-shrink: 0;
        white-space: nowrap;
        font-size: 14px;
        padding: 0 18px;
      }
    }
  }
}

@media (max-width: 640px) {
  .category-filter {
    .category-tabs {
      padding: 12px;
      gap: 8px;
      flex-wrap: nowrap;
      justify-content: flex-start;

      :deep(.n-button) {
        font-size: 13px;
        padding: 0 16px;
        height: 36px;

        .n-button__icon {
          font-size: 16px;
        }
      }
    }
  }
}

@media (max-width: 480px) {
  .category-filter {
    margin-bottom: 24px;

    .category-tabs {
      padding: 10px;
      gap: 6px;

      :deep(.n-button) {
        font-size: 12px;
        padding: 0 14px;
        height: 34px;

        .n-button__icon {
          font-size: 14px;
          margin-right: 4px;
        }
      }
    }
  }
}
</style>
