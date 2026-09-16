<script setup lang="ts">
import { requestApi } from '@/core/transport/rest/request-api'
import { encryptByAES } from '@/utils'
import { useAppConfig } from '@/hooks/use-app-config'

const { t } = useI18n()
const { resolveApiError } = useApiError()
const localePath = useLocalePath()
const config = useAppConfig()

const username = ref('')
const password = ref('')
const confirmPassword = ref('')
const loading = ref(false)
const errorMsg = ref('')
const successMsg = ref('')

const isValidUsername = computed(() => {
  if (!username.value) return false
  return /^[a-zA-Z0-9_]{3,20}$/.test(username.value)
})

const isValidPassword = computed(() => password.value.length >= 6)

const isPasswordMatch = computed(() =>
  password.value === confirmPassword.value && confirmPassword.value.length > 0
)

const isFormValid = computed(() =>
  isValidUsername.value && isValidPassword.value && isPasswordMatch.value
)

async function handleRegister() {
  errorMsg.value = ''
  successMsg.value = ''
  if (!isFormValid.value || loading.value) return

  loading.value = true
  try {
    // 密码与登录口径一致:AES 加密后提交,租户归属由服务端按 Host 解析
    await requestApi({
      path: '/app/v1/register',
      method: 'POST',
      body: JSON.stringify({
        username: username.value,
        password: encryptByAES(password.value, config.aesKey),
      }),
    })
    successMsg.value = t('authentication.register.register_success') || 'Registration successful'
    setTimeout(() => {
      navigateTo(localePath('/login'))
    }, 1200)
  } catch (e: any) {
    errorMsg.value = resolveApiError(e)
  } finally {
    loading.value = false
  }
}

const inputBase = 'w-full rounded-lg border border-border bg-background px-4 py-2.5 text-sm text-foreground transition-colors hover:border-primary focus:border-primary focus:outline-none focus:ring-[3px] focus:ring-primary/15'
</script>

<template>
  <div class="space-y-4">
    <!-- Username -->
    <div class="space-y-2">
      <label for="register-account-username" class="block text-sm font-medium text-foreground">
        {{ t('authentication.register.username') }}
      </label>
      <input
        id="register-account-username"
        v-model="username"
        type="text"
        :class="[inputBase, username && !isValidUsername ? 'border-destructive focus:border-destructive focus:ring-destructive/15' : '']"
        :placeholder="t('authentication.register.input_username')"
        autocomplete="username"
        autocapitalize="none"
        autocorrect="off"
        spellcheck="false"
      />
      <span v-if="username && !isValidUsername" class="text-xs text-destructive">
        {{ t('authentication.register.invalid_username') }}
      </span>
    </div>

    <!-- Password -->
    <div class="space-y-2">
      <label for="register-account-password" class="block text-sm font-medium text-foreground">
        {{ t('authentication.register.password') }}
      </label>
      <input
        id="register-account-password"
        v-model="password"
        type="password"
        :class="[inputBase, password && !isValidPassword ? 'border-destructive focus:border-destructive focus:ring-destructive/15' : '']"
        :placeholder="t('authentication.register.input_password')"
        autocomplete="new-password"
      />
      <span v-if="password && !isValidPassword" class="text-xs text-destructive">
        {{ t('authentication.register.invalid_password') }}
      </span>
    </div>

    <!-- Confirm Password -->
    <div class="space-y-2">
      <label for="register-account-confirm-password" class="block text-sm font-medium text-foreground">
        {{ t('authentication.register.confirm_password') }}
      </label>
      <input
        id="register-account-confirm-password"
        v-model="confirmPassword"
        type="password"
        :class="[inputBase, confirmPassword && !isPasswordMatch ? 'border-destructive focus:border-destructive focus:ring-destructive/15' : '']"
        :placeholder="t('authentication.register.input_confirm_password')"
        autocomplete="new-password"
      />
      <span v-if="confirmPassword && !isPasswordMatch" class="text-xs text-destructive">
        {{ t('authentication.register.password_not_match') }}
      </span>
    </div>

    <!-- Register Button -->
    <button
      type="button"
      class="w-full cursor-pointer rounded-lg bg-primary px-4 py-2.5 text-sm font-medium text-primary-foreground transition-all hover:bg-primary/90 active:scale-[0.98] disabled:cursor-not-allowed disabled:opacity-60"
      :disabled="!isFormValid || loading"
      @click="handleRegister"
    >
      {{ loading ? '...' : t('authentication.register.register') }}
    </button>

    <p v-if="errorMsg" class="rounded-lg border border-destructive/20 bg-destructive/5 px-4 py-3 text-sm text-destructive">
      {{ errorMsg }}
    </p>
    <p v-if="successMsg" class="rounded-lg border border-primary/20 bg-primary/5 px-4 py-3 text-sm text-primary">
      {{ successMsg }}
    </p>
  </div>
</template>
