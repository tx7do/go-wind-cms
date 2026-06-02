<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'
import { XIcon } from '@/plugins/xicon'
import { useLanguageChangeEffect } from '@/hooks/use-language-change-effect'
import { usePreferences } from '@/core/preferences/use-preferences'
import { cn } from '@/lib/utils'
import { fetchListNavigation } from '@/api/composables/navigation'

// @ts-ignore - generated types
import type { siteservicev1_NavigationItem } from '@/api/generated/app/service/v1'

const localePath = useLocalePath()
const { t } = useI18n()

const navigationItems = ref<siteservicev1_NavigationItem[]>([])
const isLoading = ref(true)
const openMenuId = ref<number | null>(null)
let closeTimer: ReturnType<typeof setTimeout> | null = null

const loadNav = async () => {
  try {
    isLoading.value = true
    const res = await fetchListNavigation({
      paging: { page: 1, pageSize: 10 },
    }) as any

    if (res.items?.length) {
      const headerNav = res.items.find((nav: any) =>
        nav.location === 'HEADER' && nav.isActive === true
      )
      const items = headerNav?.items
      if (items && items.length > 0) {
        navigationItems.value = items
      }
    }
  } catch (error) {
    console.error('[TopNavbar] 加载导航失败:', error)
  } finally {
    isLoading.value = false
  }
}

onMounted(() => {
  loadNav()
})

useLanguageChangeEffect(() => {
  isLoading.value = true
  fetchListNavigation({ paging: { page: 1, pageSize: 10 } })
    .then((res: any) => {
      if (res.items?.length) {
        const headerNav = res.items.find((nav: any) =>
          nav.location === 'HEADER' && nav.isActive === true
        )
        const items = headerNav?.items
        if (items && items.length > 0) {
          navigationItems.value = items
        }
      }
    })
    .catch((error: any) => console.error('[TopNavbar] 重新加载导航失败:', error))
    .finally(() => (isLoading.value = false))
}, { immediate: false, autoCleanup: true })

onUnmounted(() => {
  if (closeTimer) clearTimeout(closeTimer)
})

const handleNavigate = (item: siteservicev1_NavigationItem) => {
  if (item.isOpenNewTab) {
    window.open(item.url, '_blank')
  } else if (item.url != null) {
    navigateTo(localePath(item.url))
  }
}

const handleMouseEnter = (id: number) => {
  if (closeTimer) {
    clearTimeout(closeTimer)
    closeTimer = null
  }
  openMenuId.value = id
}

const handleMouseLeave = () => {
  closeTimer = setTimeout(() => (openMenuId.value = null), 150)
}

usePreferences() // keep theme reactivity
</script>

<template>
  <div v-if="isLoading" class="flex h-full items-center gap-2">
    <UiSkeleton v-for="i in 5" :key="i" class="h-7 w-16 rounded-lg" />
  </div>
  <nav v-else-if="navigationItems.length > 0" class="flex h-full items-center gap-0.5">
    <div
      v-for="item in navigationItems"
      :key="item.id ?? 0"
      class="relative"
      @mouseenter="handleMouseEnter(item.id ?? 0)"
      @mouseleave="handleMouseLeave"
    >
      <button
        type="button"
        :class="cn(
          'inline-flex h-9 items-center gap-1.5 rounded-lg px-3 text-sm font-medium',
          'text-foreground/80 transition-all duration-200',
          'hover:bg-primary/10 hover:text-primary',
          (item.children?.length ?? 0) > 0 && openMenuId === (item.id ?? 0) && 'bg-primary/10 text-primary',
          'max-md:px-2',
        )"
        @click="() => {
          if (!(item.children && item.children.length > 0)) {
            handleNavigate(item)
          }
        }"
      >
        <XIcon v-if="item.icon" :icon="`carbon:${item.icon}`" :size="16" />
        <span class="whitespace-nowrap">{{ item.title }}</span>
        <XIcon
          v-if="(item.children?.length ?? 0) > 0"
          icon="carbon:chevron-down"
          :size="12"
          :class="cn('transition-transform duration-200', openMenuId === (item.id ?? 0) && 'rotate-180')"
        />
      </button>

      <div
        v-if="(item.children?.length ?? 0) > 0 && openMenuId === (item.id ?? 0)"
        :class="cn(
          'absolute left-0 top-full z-1001 mt-1.5',
          'min-w-50 rounded-lg border border-border bg-popover p-1.5',
          'shadow-lg shadow-black/5',
          'animate-in fade-in-0 zoom-in-95 duration-150',
        )"
        @mouseenter="handleMouseEnter(item.id ?? 0)"
        @mouseleave="handleMouseLeave"
      >
        <button
          v-for="child in item.children"
          :key="child.id?.toString()"
          type="button"
          class="flex w-full cursor-pointer items-center gap-2 rounded-md px-2.5 py-2 text-left text-sm text-foreground/80 transition-colors duration-150 hover:bg-primary/10 hover:text-primary"
          @click="() => {
            handleNavigate(child)
            openMenuId = null
          }"
        >
          <XIcon v-if="child.icon" :icon="`carbon:${child.icon}`" :size="14" />
          <span class="truncate">{{ child.title }}</span>
        </button>
      </div>
    </div>
  </nav>
</template>
