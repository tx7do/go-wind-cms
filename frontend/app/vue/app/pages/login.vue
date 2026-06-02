<script setup lang="ts">
import { useAuthStore } from '@/stores/modules/app/auth.state'

definePageMeta({ layout: 'auth' })

const { t } = useI18n()

useHead({ title: t('authentication.login.title') })
const localePath = useLocalePath()
const authStore = useAuthStore()

const activeTab = ref('account')

const tabs = computed(() => [
  { key: 'account', label: t('authentication.login.tab_account') },
  { key: 'email', label: t('authentication.login.tab_email') },
  { key: 'phone', label: t('authentication.login.tab_phone') },
  { key: 'other', label: t('authentication.login.tab_other') },
])

const form = reactive({
  username: '',
  password: '',
  rememberMe: false,
})

const loading = ref(false)
const errorMsg = ref('')

const handleLogin = async () => {
  errorMsg.value = ''
  if (!form.username || !form.password) {
    errorMsg.value = t('authentication.login.fill_required')
    return
  }
  loading.value = true
  try {
    await authStore.authLogin({
      username: form.username,
      password: form.password,
    }, async () => {
      await navigateTo(localePath('/'))
    })
  } catch (e: any) {
    errorMsg.value = e?.message || t('authentication.login.login_failed')
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <AuthLayout
    :title="t('authentication.login.title')"
    :subtitle="t('authentication.login.login_with')"
    :active-tab="activeTab"
    :tabs="tabs"
    @update:active-tab="activeTab = $event"
  >
    <!-- Account login form -->
    <template v-if="activeTab === 'account'">
      <form class="space-y-4" @submit.prevent="handleLogin">
        <div v-if="errorMsg" class="rounded-lg border border-destructive/20 bg-destructive/5 px-4 py-3 text-sm text-destructive">
          {{ errorMsg }}
        </div>
        <div>
          <UiLabel class="mb-2 block text-sm font-medium">{{ t('authentication.login.username') }}</UiLabel>
          <UiInput v-model="form.username" type="text" :placeholder="t('authentication.login.placeholder_email')" />
        </div>
        <div>
          <UiLabel class="mb-2 block text-sm font-medium">{{ t('authentication.login.password') }}</UiLabel>
          <UiInput v-model="form.password" type="password" :placeholder="t('authentication.login.placeholder_password')" />
        </div>
        <div class="flex items-center justify-between text-sm">
          <label class="flex items-center gap-2 text-muted-foreground">
            <input v-model="form.rememberMe" type="checkbox" class="rounded" />
            {{ t('authentication.login.remember_me') }}
          </label>
          <button class="text-primary hover:underline bg-transparent border-none cursor-pointer text-sm">{{ t('authentication.login.forgot_password') }}</button>
        </div>
        <UiButton class="w-full" type="submit" :disabled="loading">
          <span v-if="loading">{{ t('authentication.login.logging_in') || 'Loading...' }}</span>
          <span v-else>{{ t('authentication.login.login') }}</span>
        </UiButton>
      </form>
    </template>

    <!-- Other tabs: placeholder -->
    <template v-else>
      <div class="py-8 text-center text-sm text-muted-foreground">
        {{ t('authentication.login.coming_soon') || 'Coming soon...' }}
      </div>
    </template>

    <template #switchLink>
      {{ t('authentication.login.no_account') }}
      <NuxtLink :to="localePath('/register')" class="text-primary hover:underline">{{ t('authentication.login.register_now') }}</NuxtLink>
    </template>
  </AuthLayout>
</template>
