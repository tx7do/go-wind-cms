<script setup lang="ts">
import type {contentservicev1_Category} from "@/api/generated/app/service/v1";
import {useCategoryStore} from "@/stores";

interface Props {
  category: contentservicev1_Category | null
  clickable?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  clickable: true
})

const emit = defineEmits<{
  click: [id: number]
}>()

const categoryStore = useCategoryStore()

function handleClick() {
  if (!props.category?.id || !props.clickable) return
  emit('click', props.category.id)
}
</script>

<template>
  <div
    class="category-card"
    :class="{ 'category-card--clickable': props.clickable }"
    @click="handleClick"
  >
    <div class="category-card-image">
      <img
        :src="categoryStore.getCategoryThumbnail(category)"
        :alt="categoryStore.getCategoryName(category)"
      />
      <div class="image-overlay"/>
    </div>
    <div class="category-card-content">
      <h3>{{ categoryStore.getCategoryName(category) }}</h3>
      <p>{{ categoryStore.getCategoryDescription(category) }}</p>
      <div class="category-card-meta">
        <span class="meta-icon">
          <span class="i-carbon:document"/>
        </span>
        <span class="meta-text">{{
            category?.postCount || 0
          }} {{ $t('page.posts.articles') }}</span>
      </div>
    </div>
  </div>
</template>

<style scoped lang="less">
.category-card {
  background: var(--color-bg);
  border-radius: 12px;
  overflow: hidden;
  cursor: pointer;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  border: 1px solid var(--color-border);
  height: 100%;
  display: flex;
  flex-direction: column;

  &--clickable {
    &:hover {
      transform: translateY(-6px);
      box-shadow: 0 8px 24px rgba(0, 0, 0, 0.1);
      border-color: var(--color-brand);

      .category-card-image {
        .image-overlay {
          opacity: 0.5;
        }

        img {
          transform: scale(1.08);
        }
      }

      h3 {
        color: var(--color-brand);
      }
    }
  }

  .category-card-image {
    position: relative;
    width: 100%;
    height: 160px;
    overflow: hidden;
    background: var(--color-bg);

    img {
      width: 100%;
      height: 100%;
      object-fit: cover;
      transition: transform 0.5s cubic-bezier(0.4, 0, 0.2, 1);
    }

    .image-overlay {
      position: absolute;
      top: 0;
      left: 0;
      right: 0;
      bottom: 0;
      background: rgba(0, 0, 0, 0.15);
      transition: opacity 0.3s;
      opacity: 0;
    }
  }

  .category-card-content {
    padding: 16px;

    h3 {
      font-size: 16px;
      font-weight: 700;
      margin: 0 0 8px 0;
      color: var(--color-text-primary);
      line-height: 1.4;
      display: -webkit-box;
      -webkit-line-clamp: 2;
      -webkit-box-orient: vertical;
      overflow: hidden;
      transition: color 0.3s;
    }

    p {
      color: var(--color-text-secondary);
      font-size: 13px;
      line-height: 1.6;
      margin: 0 0 12px 0;
      display: -webkit-box;
      -webkit-line-clamp: 2;
      -webkit-box-orient: vertical;
      overflow: hidden;
    }

    .category-card-meta {
      display: flex;
      align-items: center;
      gap: 6px;
      font-size: 12px;
      color: var(--color-text-secondary);
      font-weight: 500;

      .meta-icon {
        display: flex;
        align-items: center;
        font-size: 14px;
        opacity: 0.8;
      }
    }
  }
}

// 非 clickable 模式下移除 cursor
.category-card:not(.category-card--clickable) {
  cursor: default;

  &:hover {
    transform: none;
    box-shadow: 0 2px 8px rgba(0, 0, 0, 0.06);
    border-color: var(--color-border);
  }

  .category-card-image {
    .image-overlay {
      opacity: 0 !important;
    }

    img {
      transform: none !important;
    }
  }

  h3 {
    color: var(--color-text-primary) !important;
  }
}
</style>
