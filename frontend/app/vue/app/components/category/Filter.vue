<script setup lang="ts">
import { XIcon } from '@/plugins/xicon'
import { fetchListCategory, getCategoryName as getCategoryNameHelper } from '@/api/composables/category'
import type { contentservicev1_Category } from '@/api/generated/app/service/v1'
import { cn } from '@/lib/utils'

const props = withDefaults(defineProps<{
  categories?: contentservicev1_Category[]
  selectedCategory?: number | null
  treeMode?: boolean
  parentId?: number | null
  autoLoad?: boolean
}>(), {
  selectedCategory: null,
  treeMode: false,
  parentId: null,
  autoLoad: true,
})

const emit = defineEmits<{
  categoryChange: [categoryId: number | null]
  loaded: [categories: contentservicev1_Category[]]
}>()

const { t } = useI18n()

const internalCategories = ref<contentservicev1_Category[]>([])
const loading = ref(false)
const expandedIds = ref<Set<number>>(new Set())
const hideTimers = new Map<number, ReturnType<typeof setTimeout>>()

const displayCategories = computed(() => props.categories || internalCategories.value)

async function loadCategories() {
  loading.value = true
  try {
    const query: Record<string, any> = { status: 'CATEGORY_STATUS_ACTIVE' }
    if (props.parentId !== undefined && props.parentId !== null) {
      query.parentId = props.parentId
    }

    const res = await fetchListCategory({
      formValues: query,
      fieldMask: 'id,status,sort_order,icon,code,post_count,direct_post_count,parent_id,created_at,children,translations.id,translations.category_id,translations.name,translations.language_code,translations.description',
      orderBy: ['-sortOrder'],
    })
    const items = res?.items || []
    internalCategories.value = items
    emit('loaded', items)
  } catch (error) {
    console.error('[CategoryFilter] Load categories failed:', error)
  } finally {
    loading.value = false
  }
}

function getCategoryName(category: contentservicev1_Category | null): string {
  if (!category?.id) return ''
  return getCategoryNameHelper(category)
}

function hasChildren(categoryId: number): boolean {
  const category = displayCategories.value.find(cat => cat.id === categoryId)
  return !!(category?.children && category.children.length > 0)
}

function handleCategoryChange(categoryId: number | null) {
  emit('categoryChange', categoryId)
}

function handleCategoryClick(nodeId: number) {
  handleTouchToggle(nodeId)
  handleCategoryChange(nodeId)
}

function showSubmenu(nodeId: number) {
  if (hasChildren(nodeId)) {
    if (hideTimers.has(nodeId)) {
      clearTimeout(hideTimers.get(nodeId)!)
      hideTimers.delete(nodeId)
    }
    expandedIds.value = new Set(expandedIds.value).add(nodeId)
  }
}

function keepSubmenuOpen(nodeId: number) {
  if (hideTimers.has(nodeId)) {
    clearTimeout(hideTimers.get(nodeId)!)
    hideTimers.delete(nodeId)
  }
}

function hideSubmenu(nodeId: number) {
  const timer = setTimeout(() => {
    const next = new Set(expandedIds.value)
    next.delete(nodeId)
    expandedIds.value = next
    hideTimers.delete(nodeId)
  }, 150)
  hideTimers.set(nodeId, timer)
}

function handleTouchToggle(nodeId: number) {
  if (expandedIds.value.has(nodeId)) {
    hideSubmenu(nodeId)
  } else {
    expandedIds.value = new Set([nodeId])
  }
}

onMounted(() => {
  if (props.autoLoad && !props.categories) {
    loadCategories()
  }
})

// 子菜单 fixed 定位计算
function submenuStyle(nodeId: number) {
  const el = document.querySelector(`[data-category-id="${nodeId}"]`)
  if (!el) return {}
  const rect = el.getBoundingClientRect()
  return {
    top: `${rect.bottom + 6}px`,
    left: `${rect.left}px`,
  }
}

// 按钮样式
const btnBase = cn(
  'inline-flex items-center gap-2 rounded-lg px-4 py-2.5 text-sm font-medium',
  'transition-all duration-200 cursor-pointer select-none',
  'border border-border/60 bg-transparent',
  'hover:bg-muted/60 hover:border-border',
)
const btnInactive = 'text-muted-foreground hover:text-foreground'
const btnActive = 'bg-primary text-primary-foreground border-primary hover:bg-primary/90 hover:border-primary'
</script>

<template>
  <!-- Loading -->
  <div v-if="loading && autoLoad" class="relative z-50 mb-10 max-md:mb-6">
    <div class="flex flex-wrap items-center gap-2.5 rounded-xl border border-border bg-card/50 p-3.5 backdrop-blur-sm">
      <UiSkeleton class="h-9 w-24 rounded-lg" />
      <UiSkeleton class="h-9 w-20 rounded-lg" />
      <UiSkeleton class="h-9 w-28 rounded-lg" />
      <UiSkeleton class="h-9 w-22 rounded-lg" />
      <UiSkeleton class="h-9 w-24 rounded-lg" />
      <UiSkeleton class="h-9 w-20 rounded-lg" />
    </div>
  </div>

  <!-- Content -->
  <div v-else class="relative z-50 mb-10 max-md:mb-6">
    <div :class="cn(
      'flex items-center gap-3 overflow-x-auto no-scrollbar scroll-smooth whitespace-nowrap',
      'rounded-xl border border-border bg-card/50 px-4 py-2.5 backdrop-blur-sm',
      'max-md:px-3',
    )">
      <!-- 所有分类按钮 -->
      <button
        type="button"
        :class="cn(btnBase, selectedCategory === null ? btnActive : btnInactive, 'shrink-0')"
        @click="handleCategoryChange(null)"
      >
        <XIcon icon="carbon:grid" :size="15" />
        {{ t('page.posts.all_categories') }}
      </button>

      <!-- 分隔线 -->
      <div v-if="displayCategories.length > 0" class="mx-0.5 h-5 w-px bg-border shrink-0 max-md:hidden" />

      <!-- 树形模式 -->
      <template v-if="treeMode">
        <div
          v-for="node in displayCategories"
          :key="node.id"
          class="relative"
          :data-category-id="node.id"
          @mouseenter="node.id && showSubmenu(node.id)"
          @mouseleave="node.id && hideSubmenu(node.id)"
        >
          <button
            type="button"
            :class="cn(btnBase, selectedCategory === node.id ? btnActive : btnInactive, 'shrink-0')"
            @click="node.id && handleCategoryClick(node.id)"
          >
            <XIcon :icon="node.icon || 'carbon:folder'" :size="15" />
            {{ getCategoryName(node) }}
            <XIcon
              v-if="hasChildren(node.id || 0)"
              icon="carbon:chevron-down"
              :size="14"
              class="ml-0.5 opacity-60"
            />
          </button>
        </div>
      </template>

      <!-- 平铺模式 -->
      <template v-else>
        <button
          v-for="cat in displayCategories.filter(c => !c.parentId)"
          :key="cat.id"
          type="button"
          :class="cn(btnBase, selectedCategory === cat.id ? btnActive : btnInactive, 'shrink-0')"
          @click="handleCategoryChange(cat.id || 0)"
        >
          <XIcon :icon="cat.icon || 'carbon:folder'" :size="15" />
          {{ getCategoryName(cat) }}
        </button>
      </template>
    </div>

    <!-- 子分类菜单（提升到 overflow 容器外部，使用 fixed 定位） -->
    <div
      v-for="node in displayCategories"
      :key="'submenu-' + node.id"
    >
      <div
        v-if="treeMode && hasChildren(node.id || 0) && expandedIds.has(node.id || 0)"
        :class="cn(
          'fixed z-[1001] mt-1.5',
          'min-w-[200px] max-w-[320px]',
          'rounded-lg border border-border bg-popover p-1.5',
          'shadow-lg shadow-black/5',
          'animate-in fade-in-0 zoom-in-95 duration-150',
        )"
        :style="submenuStyle(node.id || 0)"
        @mouseenter="keepSubmenuOpen(node.id || 0)"
        @mouseleave="hideSubmenu(node.id || 0)"
      >
        <button
          v-for="child in node.children"
          :key="child.id"
          type="button"
          :class="cn(
            'flex w-full items-center gap-2 rounded-md px-2.5 py-2 text-left text-sm',
            'transition-colors duration-150 cursor-pointer',
            selectedCategory === child.id
              ? 'bg-primary/10 text-primary font-medium'
              : 'text-foreground/80 hover:bg-muted hover:text-foreground',
          )"
          @click.stop="handleCategoryChange(child.id || 0)"
        >
          <XIcon :icon="child.icon || 'carbon:folder'" :size="14" />
          <span class="truncate">{{ getCategoryName(child) }}</span>
          <span
            v-if="child.postCount !== undefined && child.postCount > 0"
            class="ml-auto text-xs text-muted-foreground"
          >
            {{ child.postCount }}
          </span>
        </button>
      </div>
    </div>
  </div>
</template>
