<script setup lang="ts">
import { useAuthStore } from '@/stores/modules/app/auth.state'
import { apiClient, declareCaptchaSupport, setCaptchaValues } from '@/api/client'

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

// 渐进式验证码：初始不要求；后端对失败达阈值的登录投递 CAPTCHA_REQUIRED
// reason 后点亮输入与图片，图片拉取随之懒加载。
const captchaRequired = ref(false)
const captchaId = ref('')
const captchaImage = ref('')
const captchaValue = ref('')

// 本端有验证码 UI：声明支持（后端仅对声明方启用验证码 gate）；
// 卸载时清掉验证码 header，避免残留值污染后续请求。
onMounted(() => {
  declareCaptchaSupport()
})

onUnmounted(() => {
  setCaptchaValues(undefined, undefined)
})

// 点亮后懒拉验证码图
watch(captchaRequired, (required) => {
  if (!required) return
  refreshCaptcha()
})

// captcha id/value 变化同步到请求 header（id/value 齐备才携带）
watch([captchaId, captchaValue], ([id, value]) => {
  setCaptchaValues(id, value)
})

const refreshCaptcha = async () => {
  try {
    const resp = await apiClient.authenticationService.GenerateCaptcha({})
    captchaId.value = resp.captchaId ?? ''
    captchaImage.value = resp.imageBase64 ?? ''
  } catch (e) {
    // 验证码获取失败不阻断页面（登录时会再次校验），原始错误留控制台
    console.warn('fetch captcha failed', e)
  }
}

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
    // 渐进式验证码：CAPTCHA_REQUIRED 点亮输入与图片；已处于要求态的失败刷新图片。
    if (e?.reason === 'CAPTCHA_REQUIRED') {
      if (!captchaRequired.value) {
        captchaRequired.value = true
      } else {
        refreshCaptcha()
      }
    } else if (captchaRequired.value) {
      refreshCaptcha()
    }
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

    <div v-if="captchaRequired" class="space-y-2">
      <label for="login-captcha" class="block text-sm font-medium text-foreground">
        {{ t('authentication.login.captcha_label') }}
      </label>
      <div class="flex w-full items-center gap-2">
        <input
          id="login-captcha"
          v-model="captchaValue"
          type="text"
          :class="inputBase + ' min-w-0 flex-1'"
          :placeholder="t('authentication.login.placeholder_code')"
          autocomplete="off"
          maxlength="6"
          @keydown.enter="handleLogin"
        />
        <!-- 图片盒：点击刷新；加载失败退化为可点占位 -->
        <div
          class="flex h-[42px] w-[126px] shrink-0 cursor-pointer items-center justify-center overflow-hidden rounded-lg border border-border bg-muted"
          :title="t('authentication.login.captcha_refresh')"
          @click="refreshCaptcha"
        >
          <img
            v-if="captchaImage"
            :src="captchaImage"
            alt="captcha"
            class="h-full w-full object-cover"
          >
          <span v-else class="text-xs text-muted-foreground">
            {{ t('authentication.login.captcha_refresh') }}
          </span>
        </div>
      </div>
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
