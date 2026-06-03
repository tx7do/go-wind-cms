<script setup lang="ts">
import { XIcon } from '@/plugins/xicon'
import {
  getCategoryName,
  getCategoryDescription,
  getCategoryThumbnail,
} from '@/api/composables/category'
import { cn } from '@/lib/utils'

const props = withDefaults(defineProps<{
  categories: any[]
  level?: number
  onCategoryClick?: (id: number) => void
}>(), {
  level: 0,
})

const { t } = useI18n()

const expandedCategories = ref<Set<number>>(new Set())

const levelMarginClass: Record<number, string> = {
  0: '',
  1: 'ml-8 max-md:ml-4',
  2: 'ml-16 max-md:ml-8',
  3: 'ml-24 max-md:ml-12',
}

function handleViewCategory(id: number) {
  props.onCategoryClick?.(id)
}

function toggleExpand(e: Event, category: any) {
  e.stopPropagation()
  if (category.children && category.children.length > 0) {
    const id = category.id || 0
    if (expandedCategories.value.has(id)) {
      expandedCategories.value.delete(id)
    } else {
      expandedCategories.value.add(id)
    }
  }
}

function isExpanded(category: any) {
  return expandedCategories.value.has(category.id || 0)
}
</script>

<template>
  <div v-if="categories && categories.length > 0" class="flex flex-col gap-4 max-md:gap-3">
    <div
      v-for="category in categories"
      :key="category.id"
      :class="cn('flex flex-col', levelMarginClass[props.level] || '')"
    >
      <div
        :class="cn(
          'group relative flex cursor-pointer items-center overflow-hidden',
          'rounded-xl border border-border/60 bg-card/50',
          'transition-all duration-300',
          'hover:bg-card hover:border-primary/30 hover:shadow-[0_8px_24px_-8px_hsl(var(--primary)/0.12)]',
          isExpanded(category) && 'border-primary/40 bg-card shadow-sm',
        )"
        @click="handleViewCategory(category.id || 0)"
      >
        <!-- 左侧色条指示器 -->
        <div :class="cn(
          'w-1 self-stretch flex-shrink-0 bg-primary/0 transition-colors duration-300',
          'group-hover:bg-primary',
          isExpanded(category) && 'bg-primary',
        )" />

        <div class="flex flex-1 items-center gap-4 p-4 max-md:flex-col max-md:items-start max-md:p-3 max-md:gap-3">
          <!-- 缩略图 -->
          <div :class="cn(
            'relative h-[90px] w-[130px] flex-shrink-0 overflow-hidden rounded-lg bg-muted',
            'max-md:h-[180px] max-md:w-full',
          )">
            <UiImage
              :src="getCategoryThumbnail(category)"
              :alt="getCategoryName(category, t('page.categories.category_untitled'))"
              class="h-full w-full object-cover transition-transform duration-500 group-hover:scale-105"
            />
            <div class="absolute inset-0 bg-gradient-to-t from-black/30 via-transparent to-transparent opacity-0 transition-opacity duration-300 group-hover:opacity-100" />
          </div>

          <!-- 文字区 -->
          <div class="flex flex-1 flex-col gap-1.5 max-md:w-full">
            <h3 :class="cn(
              'flex items-center gap-2 text-base font-semibold leading-tight text-foreground transition-colors',
              'group-hover:text-primary',
            )">
              {{ getCategoryName(category, t('page.categories.category_untitled')) }}
            </h3>
            <p class="line-clamp-2 text-sm leading-relaxed text-muted-foreground">
              {{ getCategoryDescription(category) }}
            </p>
            <div class="mt-1 flex items-center gap-3 text-xs text-muted-foreground">
              <span class="flex items-center gap-1">
                <XIcon icon="carbon:document" :size="14" />
                {{ category.postCount || 0 }} {{ t('page.categories.articles_count') }}
              </span>
              <span v-if="category.children && category.children.length > 0" class="flex items-center gap-1">
                <XIcon icon="carbon:folder" :size="14" />
                {{ category.children.length }}
              </span>
            </div>
          </div>

          <!-- 右侧箭头：无子分类时显示 -->
          <div
            v-if="!(category.children && category.children.length > 0)"
            class="flex-shrink-0 text-muted-foreground/50 transition-all duration-300 group-hover:text-primary group-hover:translate-x-1 max-md:hidden"
          >
            <XIcon icon="carbon:chevron-right" :size="20" />
          </div>
        </div>

        <!-- 展开/收起按钮 -->
        <button
          v-if="category.children && category.children.length > 0"
          :class="cn(
            'flex items-center justify-center px-3',
            'border-l border-border text-muted-foreground',
            'transition-colors duration-200',
            'hover:bg-primary/5 hover:text-primary',
            'max-md:absolute max-md:top-2 max-md:right-2 max-md:z-10',
            'max-md:h-8 max-md:w-8 max-md:rounded-full max-md:border max-md:border-border',
            'max-md:bg-card/90 max-md:shadow-sm max-md:border-l',
            'max-md:hover:bg-primary max-md:hover:text-white max-md:hover:scale-110',
          )"
          @click="toggleExpand($event, category)"
        >
          <XIcon
            :icon="isExpanded(category) ? 'carbon:chevron-down' : 'carbon:chevron-right'"
            :size="20"
            :class="cn('transition-transform duration-300', isExpanded(category) && 'rotate-180')"
          />
        </button>
      </div>

      <!-- 递归渲染子分类 -->
      <div
        v-if="category.children && category.children.length > 0 && isExpanded(category)"
        class="mt-2 overflow-hidden border-l-2 border-primary/20 ml-3"
      >
        <CategoryTree
          :categories="category.children"
          :level="props.level + 1"
          :on-category-click="props.onCategoryClick"
        />
      </div>
    </div>
  </div>
</template>
