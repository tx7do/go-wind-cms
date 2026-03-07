<script setup lang="ts">
import {ref} from 'vue'
import {useRouter} from 'vue-router'

interface Category {
  id: number
  translations?: Array<{
    name: string
    description: string
    thumbnail: string
  }>
  postCount: number
  children?: Category[]
  depth: number
}

interface Props {
  categories: Category[]
  level?: number
}

const props = withDefaults(defineProps<Props>(), {
  level: 0
})

const router = useRouter()

function getCategoryName(category: Category) {
  return category.translations?.[0]?.name || '未命名分类'
}

function getCategoryDescription(category: Category) {
  return category.translations?.[0]?.description || ''
}

function getCategoryThumbnail(category: Category) {
  return category.translations?.[0]?.thumbnail || '/placeholder.jpg'
}

function handleViewCategory(id: number) {
  router.push(`/category/${id}`)
}

const expandedCategories = ref<Set<number>>(new Set())

function toggleExpand(category: Category) {
  if (category.children && category.children.length > 0) {
    if (expandedCategories.value.has(category.id)) {
      expandedCategories.value.delete(category.id)
    } else {
      expandedCategories.value.add(category.id)
    }
  }
}

function isExpanded(category: Category) {
  return expandedCategories.value.has(category.id)
}
</script>

<template>
  <div class="category-tree">
    <div
      v-for="category in categories"
      :key="category.id"
      class="category-node"
      :class="[`level-${level}`]"
    >
      <div class="category-item" @click="handleViewCategory(category.id)">
        <div class="category-info">
          <div class="category-image">
            <img :src="getCategoryThumbnail(category)" :alt="getCategoryName(category)" />
            <div class="image-overlay" />
          </div>
          <div class="category-content">
            <h3 :class="{'has-children': category.children && category.children.length > 0}">
              {{ getCategoryName(category) }}
            </h3>
            <p class="description">{{ getCategoryDescription(category) }}</p>
            <div class="category-meta">
              <span class="meta-icon">
                <span class="i-carbon:document" />
              </span>
              <span class="meta-text">{{ category.postCount || 0 }} 篇文章</span>
            </div>
          </div>
        </div>
        
        <!-- 展开/收起按钮 (如果有子分类) -->
        <button
          v-if="category.children && category.children.length > 0"
          class="expand-btn"
          @click.stop="toggleExpand(category)"
          :class="{ expanded: isExpanded(category) }"
        >
          <span :class="isExpanded(category) ? 'i-carbon:chevron-down' : 'i-carbon:chevron-right'" />
        </button>
      </div>

      <!-- 递归渲染子分类 -->
      <transition name="slide">
        <CategoryTree
          v-if="category.children && category.children.length > 0 && isExpanded(category)"
          :categories="category.children"
          :level="level + 1"
        />
      </transition>
    </div>
  </div>
</template>

<style scoped lang="less">
.category-tree {
  display: flex;
  flex-direction: column;
  gap: 24px;
}

.category-node {
  .category-item {
    background: var(--color-surface);
    border-radius: 16px;
    overflow: hidden;
    cursor: pointer;
    transition: all 0.4s cubic-bezier(0.4, 0, 0.2, 1);
    border: 1px solid var(--color-border);
    display: flex;
    align-items: center;
    justify-content: space-between;
    box-shadow: 0 4px 12px rgba(0, 0, 0, 0.05);

    &:hover {
      transform: translateY(-4px);
      box-shadow: 0 12px 24px rgba(0, 0, 0, 0.1);
      border-color: var(--color-brand);

      .category-image {
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

    .category-info {
      display: flex;
      gap: 20px;
      flex: 1;
      padding: 20px;
    }

    .category-image {
      position: relative;
      width: 160px;
      height: 120px;
      flex-shrink: 0;
      overflow: hidden;
      border-radius: 12px;
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

    .category-content {
      flex: 1;
      display: flex;
      flex-direction: column;
      gap: 8px;

      h3 {
        font-size: 18px;
        font-weight: 700;
        margin: 0;
        color: var(--color-text-primary);
        line-height: 1.4;
        transition: color 0.3s;
        display: flex;
        align-items: center;
        gap: 8px;

        &.has-children {
          padding-left: 0;
        }
      }

      .description {
        color: var(--color-text-secondary);
        font-size: 14px;
        line-height: 1.6;
        margin: 0;
        display: -webkit-box;
        -webkit-line-clamp: 2;
        -webkit-box-orient: vertical;
        overflow: hidden;
      }

      .category-meta {
        display: flex;
        align-items: center;
        gap: 8px;
        font-size: 13px;
        color: var(--color-text-secondary);
        font-weight: 500;

        .meta-icon {
          display: flex;
          align-items: center;
          font-size: 16px;
          opacity: 0.8;
        }
      }
    }

    .expand-btn {
      padding: 12px;
      margin-right: 12px;
      background: transparent;
      border: 1px solid var(--color-border);
      border-radius: 8px;
      cursor: pointer;
      transition: all 0.3s;
      display: flex;
      align-items: center;
      justify-content: center;
      color: var(--color-text-secondary);

      &:hover {
        background: rgba(168, 85, 247, 0.1);
        border-color: var(--color-brand);
        color: var(--color-brand);
      }

      span[class^="i-"] {
        font-size: 20px;
        transition: transform 0.3s;
      }
    }
  }

  // 不同层级的缩进
  &.level-1 {
    margin-left: 40px;
  }

  &.level-2 {
    margin-left: 80px;
  }

  &.level-3 {
    margin-left: 120px;
  }
}

// Slide 动画
.slide-enter-active,
.slide-leave-active {
  transition: all 0.3s ease;
  overflow: hidden;
}

.slide-enter-from {
  opacity: 0;
  transform: translateY(-10px);
  max-height: 0;
}

.slide-enter-to {
  opacity: 1;
  transform: translateY(0);
  max-height: 1000px;
}

.slide-leave-from {
  opacity: 1;
  transform: translateY(0);
  max-height: 1000px;
}

.slide-leave-to {
  opacity: 0;
  transform: translateY(-10px);
  max-height: 0;
}

// 响应式
@media (max-width: 768px) {
  .category-node {
    .category-item {
      flex-direction: column;
      position: relative; // 为绝对定位的按钮提供参考

      .category-info {
        flex-direction: column;
        padding: 16px;
        gap: 12px;
        width: 100%;
      }

      .category-image {
        width: 100%;
        height: 200px;
        border-radius: 12px 12px 0 0; // 上圆角
      }

      .category-content {
        width: 100%;
        padding: 16px;
        gap: 8px;

        h3 {
          font-size: 17px;
          line-height: 1.3;
        }

        .description {
          font-size: 13px;
          line-height: 1.5;
        }

        .category-meta {
          font-size: 12px;
        }
      }

      .expand-btn {
        position: absolute;
        top: 12px;
        right: 12px;
        margin: 0;
        z-index: 10;
        background: var(--color-surface); // 使用主题色
        backdrop-filter: blur(8px);
        box-shadow: 0 2px 12px rgba(0, 0, 0, 0.15);
        border: 1px solid var(--color-border);
        color: var(--color-text-secondary);
        width: 36px;
        height: 36px;
        border-radius: 50%; // 圆形按钮
        padding: 0;
        display: flex;
        align-items: center;
        justify-content: center;
        transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
        
        &:hover {
          background: var(--color-brand);
          border-color: var(--color-brand);
          color: #fff;
          transform: scale(1.1);
          box-shadow: 0 4px 16px rgba(102, 126, 234, 0.4);
        }

        span[class^="i-"] {
          font-size: 18px;
        }
      }
    }

    &.level-1 {
      margin-left: 16px;
    }

    &.level-2 {
      margin-left: 32px;
    }

    &.level-3 {
      margin-left: 48px;
    }
  }

  // 优化子分类的间距
  .category-tree {
    gap: 16px;
  }
}
</style>
