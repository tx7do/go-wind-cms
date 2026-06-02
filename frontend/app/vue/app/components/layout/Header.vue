<script setup lang="ts">
import { XIcon } from '@/plugins/xicon'
import { usePreferences } from '@/core/preferences/use-preferences'
import { useAccessStore } from '@/stores/modules/core/access.state'
import { useAuthStore } from '@/stores/modules/app/auth.state'

const { t } = useI18n()
const { themePreferences: themePref, setTheme: setThemeMode } = usePreferences()
const currentMode = computed(() => themePref.value.mode)

const localePath = useLocalePath()
const switchLocalePath = useSwitchLocalePath()

const changeLocale = (code: 'zh-CN' | 'en-US') => {
  navigateTo(switchLocalePath(code))
}
const accessStore = useAccessStore()
const authStore = useAuthStore()

const isLogin = computed(() => {
  const token = accessStore.accessToken
  return !!token?.value && !accessStore.loginExpired
})

const handleClickLogo = () => navigateTo(localePath('/'))
const handleClickSettings = () => navigateTo(localePath('/settings'))
const handleClickUserHomepage = () => navigateTo(localePath('/user'))
const handleClickLogin = () => navigateTo(localePath('/login'))
const handleClickRegister = () => navigateTo(localePath('/register'))
const handleClickLogout = async () => {
  if (isLogin.value) await authStore.logout()
}
</script>

<template>
  <header class="fixed top-0 left-0 right-0 z-1000 flex justify-center bg-background/80 backdrop-blur-md border-b border-border/50 dark:border-border/30 dark:bg-background/60">
    <div class="flex h-(--layout-header-height) w-full max-w-300 items-center gap-6 px-4 max-md:gap-3 max-md:px-3">
      <!-- Logo -->
      <button
        type="button"
        class="flex shrink-0 items-center gap-2 rounded-lg p-1 transition-colors hover:bg-primary/10"
        aria-label="Go to homepage"
        @click="handleClickLogo"
      >
        <img src="/logo.png" alt="Logo" width="36" height="36" class="h-9 w-9 shrink-0 max-md:h-8 max-md:w-8" />
        <span class="text-lg font-bold text-primary whitespace-nowrap max-md:hidden">
          {{ t('app.title') }}
        </span>
      </button>

      <!-- 桌面端导航区 -->
      <div class="min-w-0 flex-1 max-md:hidden">
        <LayoutTopNavbar />
      </div>

      <!-- 桌面端功能按钮区 -->
      <div class="flex shrink-0 items-center gap-1 max-md:hidden">
        <!-- 用户菜单 -->
        <UiDropdownMenu :modal="false">
          <UiDropdownMenuTrigger as-child>
            <UiButton variant="ghost" size="icon" aria-label="User menu">
              <XIcon icon="lucide:user" width="16" height="16" />
            </UiButton>
          </UiDropdownMenuTrigger>
          <UiDropdownMenuContent align="end">
            <template v-if="isLogin">
              <UiDropdownMenuItem @click="handleClickUserHomepage">
                <XIcon icon="lucide:home" width="16" height="16" />
                {{ t('menu.homepage') }}
              </UiDropdownMenuItem>
              <UiDropdownMenuItem @click="handleClickSettings">
                <XIcon icon="lucide:user" width="16" height="16" />
                {{ t('menu.my_profile') }}
              </UiDropdownMenuItem>
              <UiDropdownMenuSeparator />
              <UiDropdownMenuItem class="text-destructive" @click="handleClickLogout">
                <XIcon icon="lucide:log-out" width="16" height="16" />
                {{ t('menu.logout') }}
              </UiDropdownMenuItem>
            </template>
            <template v-else>
              <UiDropdownMenuItem @click="handleClickLogin">
                <XIcon icon="lucide:user" width="16" height="16" />
                {{ t('navbar.user.login') }}
              </UiDropdownMenuItem>
              <UiDropdownMenuItem @click="handleClickRegister">
                <XIcon icon="lucide:user" width="16" height="16" />
                {{ t('navbar.user.register') }}
              </UiDropdownMenuItem>
            </template>
          </UiDropdownMenuContent>
        </UiDropdownMenu>

        <!-- 语言菜单 -->
        <UiDropdownMenu :modal="false">
          <UiDropdownMenuTrigger as-child>
            <UiButton variant="ghost" size="icon" aria-label="Language">
              <XIcon icon="lucide:globe" width="16" height="16" />
            </UiButton>
          </UiDropdownMenuTrigger>
          <UiDropdownMenuContent align="end">
            <UiDropdownMenuItem @click="changeLocale('zh-CN')">
              简体中文
            </UiDropdownMenuItem>
            <UiDropdownMenuItem @click="changeLocale('en-US')">
              English
            </UiDropdownMenuItem>
          </UiDropdownMenuContent>
        </UiDropdownMenu>

        <!-- 主题菜单 -->
        <UiDropdownMenu :modal="false">
          <UiDropdownMenuTrigger as-child>
            <UiButton variant="ghost" size="icon" aria-label="Toggle theme">
              <XIcon                 :icon="currentMode === 'dark' ? 'lucide:moon' : currentMode === 'light' ? 'lucide:sun' : 'lucide:monitor'"
                width="16" height="16"
                class="theme-icon-animate"
              />
            </UiButton>
          </UiDropdownMenuTrigger>
          <UiDropdownMenuContent align="end">
            <UiDropdownMenuItem @click="setThemeMode('dark')">
              <XIcon icon="lucide:moon" width="16" height="16" />
              {{ t('navbar.theme.dark') }}
            </UiDropdownMenuItem>
            <UiDropdownMenuItem @click="setThemeMode('light')">
              <XIcon icon="lucide:sun" width="16" height="16" />
              {{ t('navbar.theme.light') }}
            </UiDropdownMenuItem>
            <UiDropdownMenuItem @click="setThemeMode('auto')">
              <XIcon icon="lucide:monitor" width="16" height="16" />
              {{ t('navbar.theme.system') }}
            </UiDropdownMenuItem>
          </UiDropdownMenuContent>
        </UiDropdownMenu>
      </div>

      <!-- 手机端：主题切换 + 汉堡选单 -->
      <div class="flex shrink-0 items-center gap-1 md:hidden">
        <UiDropdownMenu :modal="false">
          <UiDropdownMenuTrigger as-child>
            <UiButton variant="ghost" size="icon" aria-label="Toggle theme">
              <XIcon                 :icon="currentMode === 'dark' ? 'lucide:moon' : currentMode === 'light' ? 'lucide:sun' : 'lucide:monitor'"
                width="16" height="16"
                class="theme-icon-animate"
              />
            </UiButton>
          </UiDropdownMenuTrigger>
          <UiDropdownMenuContent align="end">
            <UiDropdownMenuItem @click="setThemeMode('dark')">
              <XIcon icon="lucide:moon" width="16" height="16" />
              {{ t('navbar.theme.dark') }}
            </UiDropdownMenuItem>
            <UiDropdownMenuItem @click="setThemeMode('light')">
              <XIcon icon="lucide:sun" width="16" height="16" />
              {{ t('navbar.theme.light') }}
            </UiDropdownMenuItem>
            <UiDropdownMenuItem @click="setThemeMode('auto')">
              <XIcon icon="lucide:monitor" width="16" height="16" />
              {{ t('navbar.theme.system') }}
            </UiDropdownMenuItem>
          </UiDropdownMenuContent>
        </UiDropdownMenu>

        <!-- 汉堡选单 -->
        <LayoutMobileNav />
      </div>
    </div>
  </header>
</template>
