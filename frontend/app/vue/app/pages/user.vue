<script setup lang="ts">
import { XIcon } from '@/plugins/xicon'
import { useUserStore } from '@/stores/modules/core/user.state'
import { useAccessStore } from '@/stores/modules/core/access.state'
import { useGetMe } from '@/api/composables/user-profile'

const { t } = useI18n()
const localePath = useLocalePath()
const userStore = useUserStore()
const accessStore = useAccessStore()

const isLogin = computed(() => {
  const token = accessStore.accessToken
  return !!token?.value && !accessStore.loginExpired
})

const { data: me, isLoading, error } = useGetMe({ enabled: isLogin })

const userInfo = computed(() => me.value || userStore.user)
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
        <!-- Not logged in -->
        <div v-if="!isLogin" class="flex flex-col items-center gap-6 py-20">
          <XIcon icon="carbon:locked" :size="48" class="text-muted-foreground" />
          <p class="text-lg text-muted-foreground">{{ t('authentication.login.please_login') }}</p>
          <UiButton @click="navigateTo(localePath('/login'))">
            {{ t('navbar.user.login') }}
          </UiButton>
        </div>

        <!-- Loading -->
        <div v-else-if="isLoading" class="rounded-2xl border border-border bg-card p-8">
          <div class="flex items-center gap-6">
            <UiSkeleton class="h-20 w-20 rounded-full" />
            <div class="space-y-2">
              <UiSkeleton class="h-6 w-32" />
              <UiSkeleton class="h-4 w-48" />
            </div>
          </div>
        </div>

        <!-- User Info -->
        <div v-else-if="userInfo" class="rounded-2xl border border-border bg-card p-8">
          <div class="flex items-center gap-6">
            <div class="flex h-20 w-20 items-center justify-center rounded-full bg-primary/10 text-primary">
              <XIcon icon="carbon:user" :size="40" />
            </div>
            <div>
              <h2 class="text-xl font-bold text-foreground">{{ userInfo.realname || userInfo.nickname || userInfo.username }}</h2>
              <p class="text-sm text-muted-foreground">{{ userInfo.email || userInfo.username }}</p>
            </div>
          </div>

          <div class="mt-8 space-y-4">
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
      </div>
    </section>
  </div>
</template>
