<script setup lang="ts">
import { useRoute, useRouter } from 'vue-router'
import { XIcon } from '@/plugins/xicon'
import { useAuthStore } from '@/stores/modules/app/auth.state'
import { apiClient } from '@/api/client'

const { t } = useI18n()
const route = useRoute()
const router = useRouter()
const authStore = useAuthStore()

const ssoLoading = ref(false)
// 回调只允许消费一次（HMR/重复挂载防抖）
const ssoHandled = ref(false)

/**
 * 企业 SSO（OIDC 授权码模式）登录入口：
 *  - 按钮点击 → 后台取授权地址（state 由服务端生成暂存）→ 整页跳 IdP；
 *  - IdP 回调携带 code+state 回到登录页 → 此处消费并换发本仓令牌。
 * OIDC 未配置时点击有后端错误提示，不影响本地账密登录。
 */
onMounted(async () => {
  const code = route.query.code as string | undefined
  const state = route.query.state as string | undefined
  if (!code || !state || ssoHandled.value) return
  ssoHandled.value = true

  // 消费掉地址栏 code/state（服务端 state 一次性）
  router.replace({ path: route.path })

  ssoLoading.value = true
  try {
    await authStore.authLogin({
      grant_type: 'authorization_code',
      code,
      oidc_state: state,
    })
  } catch (error) {
    console.error('OIDC login failed:', error)
  } finally {
    ssoLoading.value = false
  }
})

const handleSsoLogin = async () => {
  ssoLoading.value = true
  try {
    const resp = await apiClient.authenticationService.GetOidcAuthorizeUrl({
      // 回调落回本页：/login?code=&state= 由 onMounted 消费
      redirect_uri: `${window.location.origin}${route.path}`,
    })
    if (resp.authorizeUrl) {
      window.location.href = resp.authorizeUrl
      return
    }
    ssoLoading.value = false
  } catch (error) {
    console.error('Get OIDC authorize url failed:', error)
    ssoLoading.value = false
  }
}
</script>

<template>
  <div class="space-y-3">
    <button
      class="flex w-full cursor-pointer items-center justify-center gap-3 rounded-lg border border-border bg-card px-4 py-3 text-sm font-medium text-foreground transition-all hover:border-primary hover:bg-primary/5 active:scale-[0.98] disabled:cursor-not-allowed disabled:opacity-60"
      :disabled="ssoLoading"
      @click="handleSsoLogin"
    >
      <XIcon icon="mdi:shield-key-outline" :size="20" />
      <span>{{ ssoLoading ? t('authentication.login.sso_redirecting') : t('authentication.login.sso_enterprise') }}</span>
    </button>
  </div>
</template>
