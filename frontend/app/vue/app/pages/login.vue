<script setup lang="ts">
import { useAuthStore } from '@/stores/modules/app/auth.state'

definePageMeta({ layout: 'auth' })

const { t } = useI18n()
const localePath = useLocalePath()
const authStore = useAuthStore()

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
  <div class="flex min-h-screen w-full items-center justify-center p-4">
    <div class="w-full max-w-md rounded-2xl border border-border bg-card p-8 shadow-lg">
      <div class="mb-8 text-center">
        <NuxtLink :to="localePath('/')" class="mb-4 inline-flex items-center gap-2">
          <img src="/logo.png" alt="Logo" width="40" height="40" class="h-10 w-10" />
          <span class="text-xl font-bold text-primary">{{ $t('app.title') }}</span>
        </NuxtLink>
        <h1 class="text-2xl font-bold text-foreground">{{ $t('authentication.login.title') }}</h1>
        <p class="mt-2 text-sm text-muted-foreground">{{ $t('authentication.login.login_with') }}</p>
      </div>

      <form class="space-y-4" @submit.prevent="handleLogin">
        <!-- Error message -->
        <div v-if="errorMsg" class="rounded-lg border border-destructive/20 bg-destructive/5 px-4 py-3 text-sm text-destructive">
          {{ errorMsg }}
        </div>

        <div>
          <UiLabel class="mb-2 block text-sm font-medium">{{ $t('authentication.login.username') }}</UiLabel>
          <UiInput v-model="form.username" type="text" :placeholder="$t('authentication.login.placeholder_email')" />
        </div>
        <div>
          <UiLabel class="mb-2 block text-sm font-medium">{{ $t('authentication.login.password') }}</UiLabel>
          <UiInput v-model="form.password" type="password" :placeholder="$t('authentication.login.placeholder_password')" />
        </div>
        <div class="flex items-center justify-between text-sm">
          <label class="flex items-center gap-2 text-muted-foreground">
            <input v-model="form.rememberMe" type="checkbox" class="rounded" />
            {{ $t('authentication.login.remember_me') }}
          </label>
          <NuxtLink :to="localePath('/')" class="text-primary hover:underline">{{ $t('authentication.login.forgot_password') }}</NuxtLink>
        </div>
        <UiButton class="w-full" type="submit" :disabled="loading">
          <span v-if="loading">{{ $t('authentication.login.logging_in') || 'Loading...' }}</span>
          <span v-else>{{ $t('authentication.login.login') }}</span>
        </UiButton>
      </form>

      <div class="mt-6 text-center text-sm text-muted-foreground">
        {{ $t('authentication.login.no_account') }}
        <NuxtLink :to="localePath('/register')" class="text-primary hover:underline">{{ $t('authentication.login.register_now') }}</NuxtLink>
      </div>
    </div>
  </div>
</template>
