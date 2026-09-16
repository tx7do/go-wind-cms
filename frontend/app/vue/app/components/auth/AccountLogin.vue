<script setup lang="ts">
import { useAuthStore } from '@/stores/modules/app/auth.state'

const { t } = useI18n()
const { resolveApiError } = useApiError()
const localePath = useLocalePath()
const route = useRoute()
const authStore = useAuthStore()

const username = ref('')
const password = ref('')
const rememberMe = ref(false)
const loading = ref(false)
const errorMsg = ref('')

const handleLogin = async () => {
  errorMsg.value = ''
  if (!username.value || !password.value) return

  loading.value = true
  try {
    await authStore.authLogin({
      username: username.value,
      password: password.value,
    }, async () => {
      // 登录前被 401 拦截到本页时会带 redirect 参数，登录成功后应回到原页面。
      // 历史链接可能存在双重编码（%252F），防御性解码一次。
      const redirect = route.query.redirect
      let target = typeof redirect === 'string' ? redirect : ''
      if (/%[0-9a-fA-F]{2}/.test(target)) {
        try {
          target = decodeURIComponent(target)
        } catch {
          // 保留原值
        }
      }
      await navigateTo(target.startsWith('/') ? target : localePath('/'))
    })
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
    <div v-if="errorMsg" class="rounded-lg border border-destructive/20 bg-destructive/5 px-4 py-3 text-sm text-destructive">
      {{ errorMsg }}
    </div>

    <div class="space-y-2">
      <label for="login-username" class="block text-sm font-medium text-foreground">
        {{ t('authentication.login.username') }}
      </label>
      <input
        id="login-username"
        v-model="username"
        type="text"
        :class="inputBase"
        :placeholder="t('authentication.register.input_username')"
        autocomplete="username"
        autocapitalize="none"
        autocorrect="off"
        spellcheck="false"
      />
    </div>

    <div class="space-y-2">
      <label for="login-password" class="block text-sm font-medium text-foreground">
        {{ t('authentication.login.password') }}
      </label>
      <input
        id="login-password"
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
      <!-- 忘记密码依赖邮件重置链路(后端暂未实现),先以禁用态展示 -->
      <button
        class="cursor-not-allowed border-none bg-transparent text-sm text-muted-foreground opacity-60"
        disabled
        :title="t('authentication.login.coming_soon')"
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
