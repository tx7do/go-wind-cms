<script setup lang="ts">
import { useAuthStore } from '@/stores/modules/app/auth.state'

const { t } = useI18n()
const localePath = useLocalePath()
const authStore = useAuthStore()

const email = ref('')
const password = ref('')
const rememberMe = ref(false)
const loading = ref(false)
const errorMsg = ref('')

const handleLogin = async () => {
  errorMsg.value = ''
  if (!email.value || !password.value) return

  loading.value = true
  try {
    await authStore.authLogin({
      email: email.value,
      password: password.value,
    }, async () => {
      await navigateTo(localePath('/'))
    })
  } catch (e: any) {
    errorMsg.value = e?.message || t('authentication.login.login_failed')
  } finally {
    loading.value = false
  }
}

const inputBase = 'w-full rounded-lg border border-border bg-background px-4 py-2.5 text-sm text-foreground transition-colors hover:border-primary focus:border-primary focus:outline-none focus:ring-[3px] focus:ring-primary/15'
</script>

<template>
  <div class="space-y-4">
    <div v-if="errorMsg" class="rounded-lg border border-destructive/20 bg-destructive/5 px-4 py-3 text-sm text-destructive">
      {{ errorMsg }}
    </div>

    <div class="space-y-2">
      <label for="login-email" class="block text-sm font-medium text-foreground">
        {{ t('authentication.register.email') }}
      </label>
      <input
        id="login-email"
        v-model="email"
        type="email"
        :class="inputBase"
        :placeholder="t('authentication.login.placeholder_email')"
        autocomplete="email"
        autocapitalize="none"
        autocorrect="off"
        spellcheck="false"
      />
    </div>

    <div class="space-y-2">
      <label for="login-password-email" class="block text-sm font-medium text-foreground">
        {{ t('authentication.login.password') }}
      </label>
      <input
        id="login-password-email"
        v-model="password"
        type="password"
        :class="inputBase"
        :placeholder="t('authentication.login.placeholder_password')"
        autocomplete="current-password"
        @keydown.enter="handleLogin"
      />
    </div>

    <div class="flex w-full items-center justify-between">
      <label class="flex cursor-pointer items-center gap-2 text-sm text-muted-foreground">
        <input
          v-model="rememberMe"
          type="checkbox"
          class="h-4 w-4 rounded border-border accent-primary"
        />
        <span>{{ t('authentication.login.remember_me') }}</span>
      </label>
      <button
        class="cursor-pointer border-none bg-transparent text-sm text-primary transition-colors hover:text-primary/80 hover:underline"
      >
        {{ t('authentication.login.forgot_password') }}
      </button>
    </div>

    <button
      class="w-full cursor-pointer rounded-lg bg-primary px-4 py-2.5 text-sm font-medium text-primary-foreground transition-all hover:bg-primary/90 active:scale-[0.98]"
      :disabled="loading"
      @click="handleLogin"
    >
      {{ loading ? (t('authentication.login.logging_in') || 'Loading...') : t('authentication.login.login') }}
    </button>
  </div>
</template>
