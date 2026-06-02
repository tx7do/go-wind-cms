<script setup lang="ts">
import { XIcon } from '@/plugins/xicon'
import { usePreferences } from '@/core/preferences/use-preferences'

const { t } = useI18n()
const { locale } = useI18n()
const localePath = useLocalePath()
const switchLocalePath = useSwitchLocalePath()

const { themePreferences: themePref, setTheme: setThemeMode } = usePreferences()

const activeMenu = ref<'account' | 'message' | 'preference'>('account')

const menuItems = [
  { key: 'account', icon: 'carbon:user', label: t('settings.menu.account') },
  { key: 'message', icon: 'carbon:email', label: t('settings.menu.message') },
  { key: 'preference', icon: 'carbon:settings', label: t('settings.menu.preference') },
]
</script>

<template>
  <div class="w-full py-8 max-md:py-4">
    <div class="w-full max-w-[1200px] mx-auto grid grid-cols-[220px_1fr_280px] gap-6 px-8 max-md:grid-cols-1 max-md:px-4">
      <!-- Left nav -->
      <aside class="max-md:hidden">
        <nav class="space-y-1">
          <div
            v-for="item in menuItems"
            :key="item.key"
            class="flex cursor-pointer items-center gap-3 rounded-lg border px-3 py-2.5 text-sm font-medium transition-all"
            :class="activeMenu === item.key
              ? 'border-primary/20 bg-primary/5 text-primary'
              : 'border-transparent text-muted-foreground hover:bg-muted hover:text-foreground'"
            @click="activeMenu = item.key as any"
          >
            <div class="flex h-8 w-8 items-center justify-center rounded-md bg-primary/10 text-primary">
              <XIcon :icon="item.icon" :size="18" />
            </div>
            <span>{{ item.label }}</span>
          </div>
        </nav>
      </aside>

      <!-- Main content -->
      <main class="min-w-0">
        <!-- Account settings -->
        <template v-if="activeMenu === 'account'">
          <div class="mb-8">
            <h1 class="mb-2 text-2xl font-bold text-foreground">{{ t('settings.account.title') }}</h1>
            <p class="text-sm text-muted-foreground">{{ t('settings.account.subtitle') }}</p>
          </div>
          <div class="mb-8">
            <h2 class="mb-2 text-lg font-semibold text-foreground">{{ t('settings.account.section_title') }}</h2>
            <p class="mb-4 text-sm text-muted-foreground">{{ t('settings.account.section_desc') }}</p>
            <div class="space-y-3">
              <UiSettingRow :label="t('settings.account.password')" :description="t('settings.account.password_not_set')">
                <UiButton variant="outline" size="sm">{{ t('settings.account.edit') }}</UiButton>
              </UiSettingRow>
              <UiSettingRow :label="t('settings.account.bind_phone')" :description="t('settings.account.password_not_set')">
                <UiButton variant="outline" size="sm">{{ t('settings.account.edit') }}</UiButton>
              </UiSettingRow>
              <UiSettingRow :label="t('settings.account.bind_email')" :description="t('settings.account.email_not_bound')">
                <UiButton variant="outline" size="sm">{{ t('settings.account.edit') }}</UiButton>
              </UiSettingRow>
            </div>
          </div>
        </template>

        <!-- Message settings -->
        <template v-if="activeMenu === 'message'">
          <div class="mb-8">
            <h1 class="mb-2 text-2xl font-bold text-foreground">{{ t('settings.message.title') }}</h1>
            <p class="text-sm text-muted-foreground">{{ t('settings.message.subtitle') }}</p>
          </div>
          <div class="mb-8">
            <h2 class="mb-2 text-lg font-semibold text-foreground">{{ t('settings.message.email_notifications') }}</h2>
            <div class="space-y-3">
              <UiSettingRow :label="t('settings.message.system_messages')" />
              <UiSettingRow :label="t('settings.message.comment_notifications')" />
              <UiSettingRow :label="t('settings.message.activity_updates')" />
              <UiSettingRow :label="t('settings.message.recommended_content')" />
            </div>
          </div>
        </template>

        <!-- Preference settings -->
        <template v-if="activeMenu === 'preference'">
          <div class="mb-8">
            <h1 class="mb-2 text-2xl font-bold text-foreground">{{ t('settings.preference.title') }}</h1>
            <p class="text-sm text-muted-foreground">{{ t('settings.preference.subtitle') }}</p>
          </div>
          <div class="mb-8">
            <h2 class="mb-2 text-lg font-semibold text-foreground">{{ t('settings.preference.theme_settings') }}</h2>
            <p class="mb-4 text-sm text-muted-foreground">{{ t('settings.preference.theme_desc') }}</p>
            <UiSettingRow :label="t('settings.preference.theme')">
              <UiSelect v-model="themePref.mode" @update:model-value="(v: any) => setThemeMode(v)">
                <UiSelectTrigger class="w-[180px] h-8">
                  <UiSelectValue />
                </UiSelectTrigger>
                <UiSelectContent>
                  <UiSelectItem value="light">{{ t('settings.preference.theme_light') }}</UiSelectItem>
                  <UiSelectItem value="dark">{{ t('settings.preference.theme_dark') }}</UiSelectItem>
                  <UiSelectItem value="auto">{{ t('settings.preference.theme_auto') }}</UiSelectItem>
                </UiSelectContent>
              </UiSelect>
            </UiSettingRow>
          </div>
        </template>
      </main>

      <!-- Right help -->
      <aside class="max-md:hidden">
        <div class="sticky top-24 rounded-lg border border-border bg-card p-5">
          <h3 class="mb-3 text-base font-bold text-foreground">{{ t('settings.help.title') }}</h3>
          <h4 class="mb-2 text-sm font-semibold text-foreground">{{ t('settings.help.account_password') }}</h4>
          <ul class="mb-4 space-y-2 text-xs">
            <li class="text-muted-foreground hover:text-foreground">1. {{ t('settings.help.q1') }}</li>
            <li class="text-muted-foreground hover:text-foreground">2. {{ t('settings.help.q2') }}</li>
            <li class="text-muted-foreground hover:text-foreground">3. {{ t('settings.help.q3') }}</li>
            <li class="text-muted-foreground hover:text-foreground">4. {{ t('settings.help.q4') }}</li>
          </ul>
        </div>
      </aside>
    </div>
  </div>
</template>
