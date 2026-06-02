<script setup lang="ts">
import { ref, onMounted, onUnmounted, watch } from 'vue'
import { XIcon } from '@/plugins/xicon'
import { useLanguageChangeEffect } from '@/hooks/use-language-change-effect'
import { usePreferences } from '@/core/preferences/use-preferences'
import { useAccessStore } from '@/stores/modules/core/access.state'
import { useAuthStore } from '@/stores/modules/app/auth.state'
import { fetchListNavigation } from '@/api/composables/navigation'
import { cn } from '@/lib/utils'

// @ts-ignore
import type { siteservicev1_NavigationItem } from '@/api/generated/app/service/v1'

const { t } = useI18n()
const localePath = useLocalePath()
const { locale } = useI18n()
const { theme: themePref, setThemeMode } = usePreferences()
const currentMode = computed(() => themePref.value.mode)
const accessStore = useAccessStore()
const authStore = useAuthStore()

const isLogin = computed(() => {
  const token = accessStore.accessToken
  return !!token?.value && !accessStore.loginExpired
})

const open = ref(false)
const navigationItems = ref<siteservicev1_NavigationItem[]>([])
const expandedId = ref<number | null>(null)

const loadNav = async () => {
  try {
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
    console.error('[MobileNav] 加载导航失败:', error)
  }
}

onMounted(loadNav)

useLanguageChangeEffect(() => {
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
    .catch((error: any) => console.error('[MobileNav] 重新加载导航失败:', error))
}, { immediate: false, autoCleanup: true })

const handleNavigate = (item: siteservicev1_NavigationItem) => {
  if (item.isOpenNewTab) {
    window.open(item.url, '_blank')
  } else if (item.url != null) {
    navigateTo(localePath(item.url))
  }
  open.value = false
}

const handleAction = (action: () => void) => {
  action()
  open.value = false
}
</script>

<template>
  <!-- 汉堡按钮 -->
  <UiButton variant="ghost" size="icon" aria-label="Open menu" class="md:hidden" @click="open = true">
    <XIcon icon="lucide:menu" width="20" height="20" />
  </UiButton>

  <!-- 抽屉 -->
  <UiSheet :open="open" @update:open="open = $event">
    <UiSheetContent side="right" class="flex w-full flex-col p-0 sm:max-w-sm">
      <!-- 顶部品牌区 -->
      <UiSheetHeader class="border-b border-border px-5 py-4">
        <div class="flex items-center gap-2">
          <img src="/logo.png" alt="Logo" width="32" height="32" class="h-8 w-8" />
          <UiSheetTitle class="text-base font-bold text-primary">
            {{ t('app.title') }}
          </UiSheetTitle>
        </div>
      </UiSheetHeader>

      <!-- 滚动内容区 -->
      <div class="flex-1 overflow-y-auto px-3 py-3">
        <!-- 主导航 -->
        <nav v-if="navigationItems.length > 0" class="mb-3 space-y-0.5">
          <div v-for="item in navigationItems" :key="item.id ?? 0">
            <button
              type="button"
              :class="cn(
                'flex w-full items-center justify-between gap-2 rounded-lg px-3 py-2.5 text-left text-sm font-medium',
                'text-foreground/80 transition-colors',
                'hover:bg-primary/10 hover:text-primary',
              )"
              @click="() => {
                const hasChildren = (item.children?.length ?? 0) > 0
                if (hasChildren) {
                  expandedId = expandedId === (item.id ?? 0) ? null : (item.id ?? 0)
                } else {
                  handleNavigate(item)
                }
              }"
            >
              <span class="flex items-center gap-2">
                <XIcon v-if="item.icon" :name="`carbon:${item.icon}`" :size="16" />
                <span>{{ item.title }}</span>
              </span>
              <XIcon                 v-if="(item.children?.length ?? 0) > 0"
                icon="lucide:chevron-down"
                width="16" height="16"
                :class="cn('transition-transform', expandedId === (item.id ?? 0) && 'rotate-180')"
              />
            </button>

            <!-- 子菜单 -->
            <div
              v-if="(item.children?.length ?? 0) > 0 && expandedId === (item.id ?? 0)"
              class="ml-3 mt-0.5 space-y-0.5 border-l border-border pl-3"
            >
              <button
                v-for="child in item.children"
                :key="child.id?.toString()"
                type="button"
                class="flex w-full items-center gap-2 rounded-md px-3 py-2 text-left text-sm text-foreground/70 transition-colors hover:bg-primary/10 hover:text-primary"
                @click="handleNavigate(child)"
              >
                <XIcon v-if="child.icon" :name="`carbon:${child.icon}`" :size="14" />
                <span>{{ child.title }}</span>
              </button>
            </div>
          </div>
        </nav>

        <!-- 分隔线 -->
        <div class="my-2 border-t border-border" />

        <!-- 用户操作 -->
        <div class="space-y-0.5">
          <template v-if="isLogin">
            <button
              type="button"
              class="flex w-full items-center gap-2 rounded-lg px-3 py-2.5 text-left text-sm font-medium text-foreground/80 transition-colors hover:bg-primary/10 hover:text-primary"
              @click="handleAction(() => navigateTo(localePath('/user')))"
            >
              <XIcon icon="lucide:home" width="16" height="16" />
              <span>{{ t('menu.homepage') }}</span>
            </button>
            <button
              type="button"
              class="flex w-full items-center gap-2 rounded-lg px-3 py-2.5 text-left text-sm font-medium text-foreground/80 transition-colors hover:bg-primary/10 hover:text-primary"
              @click="handleAction(() => navigateTo(localePath('/settings')))"
            >
              <XIcon icon="lucide:user" width="16" height="16" />
              <span>{{ t('menu.my_profile') }}</span>
            </button>
            <button
              type="button"
              class="flex w-full items-center gap-2 rounded-lg px-3 py-2.5 text-left text-sm font-medium text-destructive hover:bg-destructive/10"
              @click="handleAction(async () => { await authStore.logout() })"
            >
              <XIcon icon="lucide:log-out" width="16" height="16" />
              <span>{{ t('menu.logout') }}</span>
            </button>
          </template>
          <template v-else>
            <button
              type="button"
              class="flex w-full items-center gap-2 rounded-lg px-3 py-2.5 text-left text-sm font-medium text-foreground/80 transition-colors hover:bg-primary/10 hover:text-primary"
              @click="handleAction(() => navigateTo(localePath('/login')))"
            >
              <XIcon icon="lucide:user" width="16" height="16" />
              <span>{{ t('navbar.user.login') }}</span>
            </button>
            <button
              type="button"
              class="flex w-full items-center gap-2 rounded-lg px-3 py-2.5 text-left text-sm font-medium text-foreground/80 transition-colors hover:bg-primary/10 hover:text-primary"
              @click="handleAction(() => navigateTo(localePath('/register')))"
            >
              <XIcon icon="lucide:user" width="16" height="16" />
              <span>{{ t('navbar.user.register') }}</span>
            </button>
          </template>
        </div>

        <!-- 分隔线 -->
        <div class="my-2 border-t border-border" />

        <!-- 语言切换 -->
        <div class="mb-2 px-3 py-1 text-xs font-semibold text-muted-foreground">
          <XIcon icon="lucide:globe" width="12" height="12" class="mr-1 inline" />
          {{ t('navbar.language.title') }}
        </div>
        <div class="space-y-0.5">
          <button
            type="button"
            :class="cn(
              'flex w-full items-center gap-2 rounded-lg px-3 py-2.5 text-left text-sm font-medium transition-colors',
              locale === 'zh' ? 'bg-primary/10 text-primary' : 'text-foreground/80 hover:bg-primary/10 hover:text-primary',
            )"
            @click="handleAction(() => navigateTo(switchLocalePath('zh')))"
          >
            <span>简体中文</span>
          </button>
          <button
            type="button"
            :class="cn(
              'flex w-full items-center gap-2 rounded-lg px-3 py-2.5 text-left text-sm font-medium transition-colors',
              locale === 'en' ? 'bg-primary/10 text-primary' : 'text-foreground/80 hover:bg-primary/10 hover:text-primary',
            )"
            @click="handleAction(() => navigateTo(switchLocalePath('en')))"
          >
            <span>English</span>
          </button>
        </div>

        <!-- 分隔线 -->
        <div class="my-2 border-t border-border" />

        <!-- 主题切换 -->
        <div class="mb-2 px-3 py-1 text-xs font-semibold text-muted-foreground">
          <XIcon             :icon="currentMode === 'dark' ? 'lucide:moon' : currentMode === 'light' ? 'lucide:sun' : 'lucide:monitor'"
            width="16" height="16"
          />
          <span class="ml-1">{{ t('navbar.theme.title') || 'Theme' }}</span>
        </div>
        <div class="space-y-0.5">
          <button
            v-for="mode in [
              { key: 'light', icon: 'lucide:sun', label: t('navbar.theme.light') },
              { key: 'dark', icon: 'lucide:moon', label: t('navbar.theme.dark') },
              { key: 'auto', icon: 'lucide:monitor', label: t('navbar.theme.system') },
            ]"
            :key="mode.key"
            type="button"
            :class="cn(
              'flex w-full items-center gap-2 rounded-lg px-3 py-2.5 text-left text-sm font-medium transition-colors',
              currentMode === mode.key ? 'bg-primary/10 text-primary' : 'text-foreground/80 hover:bg-primary/10 hover:text-primary',
            )"
            @click="handleAction(() => setThemeMode(mode.key as any))"
          >
            <XIcon :icon="mode.icon" width="16" height="16" />
            <span>{{ mode.label }}</span>
          </button>
        </div>
      </div>
    </UiSheetContent>
  </UiSheet>
</template>
