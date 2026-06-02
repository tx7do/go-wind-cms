<script setup lang="ts">
import { XIcon } from '@/plugins/xicon'
import { usePreferences } from '@/core/preferences/use-preferences'
import { useAccessStore } from '@/stores/modules/core/access.state'
import { useAuthStore } from '@/stores/modules/app/auth.state'

const { t } = useI18n()
const { theme: themePref, setThemeMode } = usePreferences()
const currentMode = computed(() => themePref.value.mode)

const localePath = useLocalePath()
const switchLocalePath = useSwitchLocalePath()
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
        <UIDropdownMenu>
          <UIDropdownMenuTrigger as-child>
            <UIButton variant="ghost" size="icon" aria-label="User menu">
              <iconify-icon icon="lucide:user" width="16" height="16" />
            </UIButton>
          </UIDropdownMenuTrigger>
          <UIDropdownMenuContent align="end">
            <template v-if="isLogin">
              <UIDropdownMenuItem @click="handleClickUserHomepage">
                <iconify-icon icon="lucide:home" width="16" height="16" />
                {{ t('menu.homepage') }}
              </UIDropdownMenuItem>
              <UIDropdownMenuItem @click="handleClickSettings">
                <iconify-icon icon="lucide:user" width="16" height="16" />
                {{ t('menu.my_profile') }}
              </UIDropdownMenuItem>
              <UIDropdownMenuSeparator />
              <UIDropdownMenuItem class="text-destructive" @click="handleClickLogout">
                <iconify-icon icon="lucide:log-out" width="16" height="16" />
                {{ t('menu.logout') }}
              </UIDropdownMenuItem>
            </template>
            <template v-else>
              <UIDropdownMenuItem @click="handleClickLogin">
                <iconify-icon icon="lucide:user" width="16" height="16" />
                {{ t('navbar.user.login') }}
              </UIDropdownMenuItem>
              <UIDropdownMenuItem @click="handleClickRegister">
                <iconify-icon icon="lucide:user" width="16" height="16" />
                {{ t('navbar.user.register') }}
              </UIDropdownMenuItem>
            </template>
          </UIDropdownMenuContent>
        </UIDropdownMenu>

        <!-- 语言菜单 -->
        <UIDropdownMenu>
          <UIDropdownMenuTrigger as-child>
            <UIButton variant="ghost" size="icon" aria-label="Language">
              <iconify-icon icon="lucide:globe" width="16" height="16" />
            </UIButton>
          </UIDropdownMenuTrigger>
          <UIDropdownMenuContent align="end">
            <UIDropdownMenuItem @click="switchLocalePath('zh')">
              简体中文
            </UIDropdownMenuItem>
            <UIDropdownMenuItem @click="switchLocalePath('en')">
              English
            </UIDropdownMenuItem>
          </UIDropdownMenuContent>
        </UIDropdownMenu>

        <!-- 主题菜单 -->
        <UIDropdownMenu>
          <UIDropdownMenuTrigger as-child>
            <UIButton variant="ghost" size="icon" aria-label="Toggle theme">
              <iconify-icon
                :icon="currentMode === 'dark' ? 'lucide:moon' : currentMode === 'light' ? 'lucide:sun' : 'lucide:monitor'"
                width="16" height="16"
                class="theme-icon-animate"
              />
            </UIButton>
          </UIDropdownMenuTrigger>
          <UIDropdownMenuContent align="end">
            <UIDropdownMenuItem @click="setThemeMode('dark')">
              <iconify-icon icon="lucide:moon" width="16" height="16" />
              {{ t('navbar.theme.dark') }}
            </UIDropdownMenuItem>
            <UIDropdownMenuItem @click="setThemeMode('light')">
              <iconify-icon icon="lucide:sun" width="16" height="16" />
              {{ t('navbar.theme.light') }}
            </UIDropdownMenuItem>
            <UIDropdownMenuItem @click="setThemeMode('auto')">
              <iconify-icon icon="lucide:monitor" width="16" height="16" />
              {{ t('navbar.theme.system') }}
            </UIDropdownMenuItem>
          </UIDropdownMenuContent>
        </UIDropdownMenu>
      </div>

      <!-- 手机端：主题切换 + 汉堡选单 -->
      <div class="flex shrink-0 items-center gap-1 md:hidden">
        <UIDropdownMenu>
          <UIDropdownMenuTrigger as-child>
            <UIButton variant="ghost" size="icon" aria-label="Toggle theme">
              <iconify-icon
                :icon="currentMode === 'dark' ? 'lucide:moon' : currentMode === 'light' ? 'lucide:sun' : 'lucide:monitor'"
                width="16" height="16"
                class="theme-icon-animate"
              />
            </UIButton>
          </UIDropdownMenuTrigger>
          <UIDropdownMenuContent align="end">
            <UIDropdownMenuItem @click="setThemeMode('dark')">
              <iconify-icon icon="lucide:moon" width="16" height="16" />
              {{ t('navbar.theme.dark') }}
            </UIDropdownMenuItem>
            <UIDropdownMenuItem @click="setThemeMode('light')">
              <iconify-icon icon="lucide:sun" width="16" height="16" />
              {{ t('navbar.theme.light') }}
            </UIDropdownMenuItem>
            <UIDropdownMenuItem @click="setThemeMode('auto')">
              <iconify-icon icon="lucide:monitor" width="16" height="16" />
              {{ t('navbar.theme.system') }}
            </UIDropdownMenuItem>
          </UIDropdownMenuContent>
        </UIDropdownMenu>

        <!-- 汉堡选单 -->
        <LayoutMobileNav />
      </div>
    </div>
  </header>
</template>
