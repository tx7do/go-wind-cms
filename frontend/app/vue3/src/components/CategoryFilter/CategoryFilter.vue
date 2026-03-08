<script setup lang="ts">
import {computed, ref, onMounted} from 'vue'

import {$t} from '@/locales'
import {useCategoryStore} from '@/stores'
import {XIcon} from '@/plugins/xicon'
import {useLanguageChangeEffect} from "@/hooks/useLanguageChangeEffect";
import type {
  contentservicev1_Category
} from "@/api/generated/app/service/v1";

interface Props {
  categories?: any[] // 外部传入的分类数据（可选）
  selectedCategory?: number | null
  treeMode?: boolean
  parentId?: number | null // 支持根据 parentId 过滤
  autoLoad?: boolean // 是否自动加载数据
}

interface Emits {
  (e: 'category-change', categoryId: number | null): void

  (e: 'loaded', categories: any[]): void // 数据加载完成事件
}

const props = defineProps<Props>()
const emit = defineEmits<Emits>()
const categoryStore = useCategoryStore()

// 内部状态管理
const categories = ref<any[]>(props.categories || [])
const loading = ref(false)

// 加载分类数据
async function loadCategories() {
  loading.value = true
  try {
    const query: any = {status: 'CATEGORY_STATUS_ACTIVE'}

    // 如果指定了 parentId，添加过滤条件
    if (props.parentId !== undefined && props.parentId !== null) {
      query.parentId = props.parentId
    }

    const res = await categoryStore.listCategory(
      undefined,
      query
    )

    categories.value = res.items || []
    emit('loaded', categories.value)
    console.log('[CategoryFilter] Categories loaded:', categories.value.length)
  } catch (error) {
    console.error('[CategoryFilter] Load categories failed:', error)
  } finally {
    loading.value = false
  }
}

// 监听语言切换，自动重新加载数据
useLanguageChangeEffect(loadCategories, {
  immediate: true,      // 是否立即执行一次
  autoCleanup: true,    // 是否自动清理
});

// 创建一个响应式的分类名称映射
const categoryNamesMap = computed(() => {
  const map = new Map()
  const currentCategories = props.categories || categories.value

  currentCategories.forEach(cat => {
    if (cat?.id) {
      map.set(cat.id, categoryStore.getCategoryName(cat))
    }
  })
  // 处理子分类
  currentCategories.forEach(cat => {
    if (cat?.children) {
      cat.children.forEach((child: contentservicev1_Category) => {
        if (child?.id) {
          map.set(child.id, categoryStore.getCategoryName(child))
        }
      })
    }
  })
  return map
})

// 获取分类名称 - 使用响应式的 Map
function getCategoryName(category: any) {
  if (!category?.id) return ''
  return categoryNamesMap.value.get(category.id) || categoryStore.getCategoryName(category)
}

// 使用外部传入的 categories 或内部加载的 categories
const displayCategories = computed(() => {
  return props.categories || categories.value
})

// 定时器管理
const hideTimers = new Map<number, ReturnType<typeof setTimeout>>()

// 展开/收起状态
const expandedIds = ref<Set<number>>(new Set()) // 默认全部收起

// 平铺模式：只显示根节点
const rootCategories = computed(() => {
  const currentCategories = displayCategories.value
  return currentCategories.filter(cat => !cat.parentId)
})

function handleCategoryChange(categoryId: number | null) {
  emit('category-change', categoryId)
}

// 处理分类按钮点击 (同时切换菜单)
function handleCategoryClick(nodeId: number) {
  // 先切换菜单
  handleTouchToggle(nodeId)
  // 再触发分类选择
  handleCategoryChange(nodeId)
}

// 显示子菜单
function showSubmenu(nodeId: number) {
  if (hasChildren(nodeId)) {
    // 清除定时器
    if (hideTimers.has(nodeId)) {
      clearTimeout(hideTimers.get(nodeId))
      hideTimers.delete(nodeId)
    }
    expandedIds.value.add(nodeId)
  }
}

// 保持子菜单打开 (鼠标在菜单上时)
function keepSubmenuOpen(nodeId: number) {
  // 清除定时器，防止菜单消失
  if (hideTimers.has(nodeId)) {
    clearTimeout(hideTimers.get(nodeId))
    hideTimers.delete(nodeId)
  }
}

// 隐藏子菜单 - 添加延时避免快速消失
function hideSubmenu(nodeId: number) {
  // 设置定时器
  const timer = setTimeout(() => {
    expandedIds.value.delete(nodeId)
    hideTimers.delete(nodeId)
  }, 150)
  hideTimers.set(nodeId, timer)
}

// 处理触摸切换 (移动端)
function handleTouchToggle(nodeId: number) {
  // 如果已经有子菜单展开，且点击的是同一个，则收起
  if (expandedIds.value.has(nodeId)) {
    hideSubmenu(nodeId)
  } else {
    // 关闭其他展开的菜单
    expandedIds.value.clear()
    showSubmenu(nodeId)
  }
}

function hasChildren(categoryId: number): boolean {
  const currentCategories = displayCategories.value
  const category = currentCategories.find(cat => cat.id === categoryId)
  return !!(category && category.children && category.children.length > 0)
}

// 生命周期：自动加载数据
onMounted(async () => {
  // 如果没有外部传入 categories 且 autoLoad 未禁用，则自动加载
  if (!props.categories && props.autoLoad !== false) {
    await loadCategories()
  }
})
</script>

<template>
  <div class="category-filter">
    <div class="category-tabs">
      <!-- 所有分类按钮 -->
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

      <!-- 树形模式 -->
      <template v-if="treeMode">
        <!-- 一级分类 (横向排列 + 悬浮菜单) -->
        <div
          v-for="node in displayCategories"
          :key="node.id"
          class="category-item-wrapper"
          @mouseenter="showSubmenu(node.id)"
          @mouseleave="hideSubmenu(node.id)"
        >
          <n-button
            :type="selectedCategory === node.id ? 'primary' : 'default'"
            :ghost="selectedCategory !== node.id"
            size="large"
            @click="handleCategoryClick(node.id)"
            @touchstart="handleTouchToggle(node.id)"
          >
            <template #icon>
              <XIcon :name="node.icon || 'carbon:folder'"/>
            </template>
            {{ getCategoryName(node) }}
          </n-button>

          <!-- 子分类菜单 -->
          <div
            v-if="hasChildren(node.id)"
            class="category-submenu"
            v-show="expandedIds.has(node.id)"
            @mouseenter="keepSubmenuOpen(node.id)"
            @mouseleave="hideSubmenu(node.id)"
            @touchstart.stop
          >
            <n-button
              v-for="child in node.children"
              :key="child.id"
              :type="selectedCategory === child.id ? 'primary' : 'default'"
              :ghost="selectedCategory !== child.id"
              size="medium"
              @click.stop="handleCategoryChange(child.id)"
            >
              <template #icon>
                <XIcon :name="child.icon || 'carbon:folder'"/>
              </template>
              {{ getCategoryName(child) }}
            </n-button>
          </div>
        </div>
      </template>

      <!-- 平铺模式 (默认) -->
      <template v-else>
        <n-button
          v-for="cat in rootCategories"
          :key="cat.id"
          :type="selectedCategory === cat.id ? 'primary' : 'default'"
          :ghost="selectedCategory !== cat.id"
          size="large"
          @click="handleCategoryChange(cat.id)"
        >
          <template #icon>
            <XIcon :name="cat.icon || 'carbon:folder'"/>
          </template>
          {{ getCategoryName(cat) }}
        </n-button>
      </template>
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

// 悬浮菜单样式
.category-item-wrapper {
  position: relative;
  display: inline-block;

  .category-submenu {
    position: absolute;
    top: 100%;
    left: 0;
    z-index: 1000;
    display: flex;
    flex-wrap: wrap;
    gap: 6px;
    min-width: min-content;
    max-width: 280px;
    padding: 8px;
    margin-top: 8px;
    background: var(--color-surface);
    border-radius: 12px;
    box-shadow: 0 4px 24px rgba(0, 0, 0, 0.15);
    backdrop-filter: blur(10px);
    // 添加过渡动画
    animation: slideDown 0.2s ease-out;

    :deep(.n-button) {
      padding: 0 12px;
      font-size: 13px;
      height: 32px;
      flex-shrink: 0;
    }
  }
}

// 下滑动画
@keyframes slideDown {
  from {
    opacity: 0;
    transform: translateY(-8px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
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
