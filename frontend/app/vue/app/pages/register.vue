<script setup lang="ts">
definePageMeta({ layout: 'auth' })

const { t } = useI18n()

useHead({ title: t('authentication.register.title') })
const localePath = useLocalePath()

const activeTab = ref('account')

const tabs = computed(() => [
  { key: 'account', label: t('authentication.login.tab_account') },
  { key: 'email', label: t('authentication.login.tab_email') },
  { key: 'phone', label: t('authentication.login.tab_phone') },
  { key: 'other', label: t('authentication.login.tab_other') },
])

const form = reactive({
  username: '',
  email: '',
  password: '',
  confirmPassword: '',
})

const loading = ref(false)
const errorMsg = ref('')

const handleRegister = async () => {
  errorMsg.value = ''
  if (!form.username || !form.email || !form.password || !form.confirmPassword) {
    errorMsg.value = t('authentication.register.fill_required')
    return
  }
  if (form.password !== form.confirmPassword) {
    errorMsg.value = t('authentication.register.password_mismatch')
    return
  }
  loading.value = true
  try {
    // TODO: 调用注册 API
    await navigateTo(localePath('/login'))
  } catch (e: any) {
    errorMsg.value = e?.message || t('authentication.register.register_failed')
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <AuthLayout
    :title="t('authentication.register.title')"
    :subtitle="t('authentication.register.register_with')"
    :active-tab="activeTab"
    :tabs="tabs"
    @update:active-tab="activeTab = $event"
  >
    <template v-if="activeTab === 'account'">
      <form class="space-y-4" @submit.prevent="handleRegister">
        <div v-if="errorMsg" class="rounded-lg border border-destructive/20 bg-destructive/5 px-4 py-3 text-sm text-destructive">
          {{ errorMsg }}
        </div>
        <div>
          <UiLabel class="mb-2 block text-sm font-medium">{{ t('authentication.register.username') }}</UiLabel>
          <UiInput v-model="form.username" type="text" :placeholder="t('authentication.register.input_username')" />
        </div>
        <div>
          <UiLabel class="mb-2 block text-sm font-medium">{{ t('authentication.register.email') }}</UiLabel>
          <UiInput v-model="form.email" type="email" :placeholder="t('authentication.register.input_email')" />
        </div>
        <div>
          <UiLabel class="mb-2 block text-sm font-medium">{{ t('authentication.register.password') }}</UiLabel>
          <UiInput v-model="form.password" type="password" :placeholder="t('authentication.register.input_password')" />
        </div>
        <div>
          <UiLabel class="mb-2 block text-sm font-medium">{{ t('authentication.register.confirm_password') }}</UiLabel>
          <UiInput v-model="form.confirmPassword" type="password" :placeholder="t('authentication.register.input_confirm_password')" />
        </div>
        <UiButton class="w-full" type="submit" :disabled="loading">
          <span v-if="loading">{{ t('authentication.register.registering') || 'Loading...' }}</span>
          <span v-else>{{ t('authentication.register.register') }}</span>
        </UiButton>
      </form>
    </template>

    <template v-else>
      <div class="py-8 text-center text-sm text-muted-foreground">
        {{ t('authentication.login.coming_soon') || 'Coming soon...' }}
      </div>
    </template>

    <template #switchLink>
      {{ t('authentication.register.already_have_account') }}
      <NuxtLink :to="localePath('/login')" class="text-primary hover:underline">{{ t('authentication.register.login_now') }}</NuxtLink>
    </template>
  </AuthLayout>
</template>
