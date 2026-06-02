<script setup lang="ts">
import { XIcon } from '@/plugins/xicon'

const { t } = useI18n()
const localePath = useLocalePath()

// Placeholder - will be replaced with actual user data
</script>

<template>
  <div class="w-full">
    <LayoutPageHero
      :title="t('menu.my_profile')"
      icon="carbon:user"
      size="md"
    />

    <section class="w-full py-12 max-md:py-8">
      <div class="mx-auto max-w-3xl px-4">
        <div class="rounded-2xl border border-border bg-card p-8">
          <div class="flex items-center gap-6">
            <div class="flex h-20 w-20 items-center justify-center rounded-full bg-primary/10 text-primary">
              <XIcon name="carbon:user" :size="40" />
            </div>
            <div>
              <h2 class="text-xl font-bold text-foreground">{{ t('menu.my_profile') }}</h2>
              <p class="text-sm text-muted-foreground">user@gowind.dev</p>
            </div>
          </div>

          <div class="mt-8 space-y-4">
            <UISettingRow :label="t('settings.account.password')" :description="t('settings.account.password_not_set')">
              <UIButton variant="outline" size="sm">{{ t('settings.account.edit') }}</UIButton>
            </UISettingRow>
            <UISettingRow :label="t('settings.account.bind_phone')" :description="t('settings.account.password_not_set')">
              <UIButton variant="outline" size="sm">{{ t('settings.account.edit') }}</UIButton>
            </UISettingRow>
            <UISettingRow :label="t('settings.account.bind_email')" :description="t('settings.account.email_not_bound')">
              <UIButton variant="outline" size="sm">{{ t('settings.account.edit') }}</UIButton>
            </UISettingRow>
          </div>
        </div>
      </div>
    </section>
  </div>
</template>
