<script setup lang="ts">
import { XIcon } from '@/plugins/xicon'
import { usePreferences } from '@/core/preferences/use-preferences'

const { isDark, toggleTheme } = usePreferences()
const { locale } = useI18n()
const switchLocalePath = useSwitchLocalePath()

const languageOptions: { key: 'zh-CN' | 'en-US'; label: string }[] = [
  { key: 'zh-CN', label: '中文' },
  { key: 'en-US', label: 'English' },
]

const handleSelectLanguage = (key: string) => {
  navigateTo(switchLocalePath(key as 'zh-CN' | 'en-US'))
}
</script>

<template>
  <div class="fixed top-8 right-8 z-100 flex gap-3 rounded-2xl border border-border bg-card p-2.5 px-3.5 backdrop-blur-md shadow-lg transition-all duration-300 hover:-translate-y-0.5 hover:shadow-xl max-lg:top-4 max-lg:right-4 max-lg:p-1.5 max-lg:gap-1.5 max-md:top-3 max-md:right-3 max-md:p-1 max-md:gap-1">
    <select
      class="rounded-lg border border-border bg-background px-3 py-2 text-sm font-medium text-foreground cursor-pointer outline-none transition-all duration-300 hover:border-primary hover:bg-card max-md:px-2 max-md:py-1.5 max-md:text-xs"
      :value="locale"
      @change="handleSelectLanguage(($event.target as HTMLSelectElement).value)"
    >
      <option v-for="opt in languageOptions" :key="opt.key" :value="opt.key">
        {{ opt.label }}
      </option>
    </select>
    <button
      class="flex min-w-9 h-9 items-center justify-center gap-1.5 rounded-lg bg-background border border-border text-foreground cursor-pointer transition-all duration-300 hover:-translate-y-0.5 hover:shadow-md hover:border-primary hover:bg-card max-md:min-w-8 max-md:h-8 max-md:text-sm"
      @click="toggleTheme"
      aria-label="Toggle theme"
    >
      <XIcon :icon="isDark ? 'carbon:sun' : 'carbon:moon'" :size="18" class="theme-icon-animate" />
    </button>
  </div>
</template>
