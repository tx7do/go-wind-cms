<script setup lang="ts">
import { XIcon } from '@/plugins/xicon'
import { usePreferences } from '@/core/preferences/use-preferences'
import { requestApi } from '@/core/transport/rest/request-api'
import { encryptByAES } from '@/utils'
import { useAppConfig } from '@/hooks/use-app-config'
import { useAccessStore } from '@/stores/modules/core/access.state'

const { t } = useI18n()
const { locale } = useI18n()
const { resolveApiError } = useApiError()

useHead({ title: t('settings.account.title') })
const localePath = useLocalePath()
const switchLocalePath = useSwitchLocalePath()

const config = useAppConfig()
const accessStore = useAccessStore()

const { themePreferences: themePref, setTheme: setThemeMode } = usePreferences()

const activeMenu = ref<'account' | 'message' | 'preference'>('account')

const menuItems = [
  { key: 'account', icon: 'carbon:user', label: t('settings.menu.account') },
  { key: 'message', icon: 'carbon:email', label: t('settings.menu.message') },
  { key: 'preference', icon: 'carbon:settings', label: t('settings.menu.preference') },
]

// ── 修改密码 ──
const pwdEditing = ref(false)
const pwdLoading = ref(false)
const pwdError = ref('')
const pwdSuccess = ref(false)
const oldPassword = ref('')
const newPassword = ref('')
const confirmPassword = ref('')

const pwdFormValid = computed(() =>
  oldPassword.value.length > 0
  && newPassword.value.length >= 6
  && newPassword.value === confirmPassword.value
)

const inputBase = 'w-full rounded-lg border border-border bg-background px-4 py-2.5 text-sm text-foreground transition-colors hover:border-primary focus:border-primary focus:outline-none focus:ring-[3px] focus:ring-primary/15'

function openPwdForm() {
  pwdEditing.value = true
  pwdError.value = ''
  pwdSuccess.value = false
  oldPassword.value = ''
  newPassword.value = ''
  confirmPassword.value = ''
}

function cancelPwdForm() {
  pwdEditing.value = false
  pwdError.value = ''
}

async function submitChangePassword() {
  pwdError.value = ''
  pwdSuccess.value = false
  if (!accessStore.accessToken?.value) {
    pwdError.value = t('settings.account.login_required')
    return
  }
  if (newPassword.value.length < 6) {
    pwdError.value = t('settings.account.password_too_short')
    return
  }
  if (newPassword.value !== confirmPassword.value) {
    pwdError.value = t('settings.account.password_mismatch')
    return
  }

  pwdLoading.value = true
  try {
    // 与登录/注册口径一致:新旧密码均 AES 加密后提交,服务端解密校验/入库
    await requestApi({
      path: '/app/v1/me/password',
      method: 'POST',
      body: JSON.stringify({
        oldPassword: encryptByAES(oldPassword.value, config.aesKey),
        newPassword: encryptByAES(newPassword.value, config.aesKey),
      }),
    })
    pwdSuccess.value = true
    pwdEditing.value = false
    oldPassword.value = ''
    newPassword.value = ''
    confirmPassword.value = ''
  } catch (e: any) {
    pwdError.value = resolveApiError(e)
  } finally {
    pwdLoading.value = false
  }
}
</script>

<template>
  <div class="w-full py-8 max-md:py-4">
    <div class="w-full max-w-[1200px] mx-auto grid grid-cols-[220px_1fr_280px] gap-6 px-8 max-md:grid-cols-1 max-md:px-4">
      <!-- Left nav -->
      <aside class="max-md:hidden">
        <nav class="space-y-1">
          <div
            v-for="item in menuItems"
            :key="item.key"
            class="flex cursor-pointer items-center gap-3 rounded-lg border px-3 py-2.5 text-sm font-medium transition-all"
            :class="activeMenu === item.key
              ? 'border-primary/20 bg-primary/5 text-primary'
              : 'border-transparent text-muted-foreground hover:bg-muted hover:text-foreground'"
            @click="activeMenu = item.key as any"
          >
            <div class="flex h-8 w-8 items-center justify-center rounded-md bg-primary/10 text-primary">
              <XIcon :icon="item.icon" :size="18" />
            </div>
            <span>{{ item.label }}</span>
          </div>
        </nav>
      </aside>

      <!-- Main content -->
      <main class="min-w-0">
        <!-- Account settings -->
        <template v-if="activeMenu === 'account'">
          <div class="mb-8">
            <h1 class="mb-2 text-2xl font-bold text-foreground">{{ t('settings.account.title') }}</h1>
            <p class="text-sm text-muted-foreground">{{ t('settings.account.subtitle') }}</p>
          </div>
          <div class="mb-8">
            <h2 class="mb-2 text-lg font-semibold text-foreground">{{ t('settings.account.section_title') }}</h2>
            <p class="mb-4 text-sm text-muted-foreground">{{ t('settings.account.section_desc') }}</p>
            <div class="space-y-3">
              <div class="rounded-lg border border-border bg-cardBg">
                <UiSettingRow :label="t('settings.account.password')" :description="t('settings.account.password_not_set')">
                  <UiButton v-if="!pwdEditing" variant="outline" size="sm" @click="openPwdForm">
                    {{ t('settings.account.edit') }}
                  </UiButton>
                  <UiButton v-else variant="outline" size="sm" @click="cancelPwdForm">
                    {{ t('settings.account.cancel') }}
                  </UiButton>
                </UiSettingRow>
                <div v-if="pwdEditing" class="space-y-3 border-t border-border px-4 py-4">
                  <input
                    v-model="oldPassword"
                    type="password"
                    :placeholder="t('settings.account.input_old_password')"
                    autocomplete="current-password"
                    :class="inputBase"
                  />
                  <input
                    v-model="newPassword"
                    type="password"
                    :placeholder="t('settings.account.input_new_password')"
                    autocomplete="new-password"
                    :class="inputBase"
                  />
                  <input
                    v-model="confirmPassword"
                    type="password"
                    :placeholder="t('settings.account.input_confirm_new_password')"
                    autocomplete="new-password"
                    :class="inputBase"
                  />
                  <p v-if="pwdError" class="text-sm text-destructive">{{ pwdError }}</p>
                  <div class="flex justify-end gap-2">
                    <UiButton variant="outline" size="sm" @click="cancelPwdForm">
                      {{ t('settings.account.cancel') }}
                    </UiButton>
                    <UiButton size="sm" :disabled="!pwdFormValid || pwdLoading" @click="submitChangePassword">
                      {{ pwdLoading ? t('settings.account.saving') : t('settings.account.save') }}
                    </UiButton>
                  </div>
                </div>
              </div>
              <UiSettingRow :label="t('settings.account.bind_phone')" :description="t('settings.account.password_not_set')">
                <UiButton variant="outline" size="sm" disabled :title="t('settings.account.not_available')">
                  {{ t('settings.account.not_available') }}
                </UiButton>
              </UiSettingRow>
              <UiSettingRow :label="t('settings.account.bind_email')" :description="t('settings.account.email_not_bound')">
                <UiButton variant="outline" size="sm" disabled :title="t('settings.account.not_available')">
                  {{ t('settings.account.not_available') }}
                </UiButton>
              </UiSettingRow>
              <p v-if="pwdSuccess" class="rounded-lg border border-primary/20 bg-primary/5 px-4 py-3 text-sm text-primary">
                {{ t('settings.account.change_success') }}
              </p>
            </div>
          </div>
        </template>

        <!-- Message settings -->
        <template v-if="activeMenu === 'message'">
          <div class="mb-8">
            <h1 class="mb-2 text-2xl font-bold text-foreground">{{ t('settings.message.title') }}</h1>
            <p class="text-sm text-muted-foreground">{{ t('settings.message.subtitle') }}</p>
          </div>
          <div class="mb-8">
            <h2 class="mb-2 text-lg font-semibold text-foreground">{{ t('settings.message.email_notifications') }}</h2>
            <div class="space-y-3">
              <UiSettingRow :label="t('settings.message.system_messages')" />
              <UiSettingRow :label="t('settings.message.comment_notifications')" />
              <UiSettingRow :label="t('settings.message.activity_updates')" />
              <UiSettingRow :label="t('settings.message.recommended_content')" />
            </div>
          </div>
        </template>

        <!-- Preference settings -->
        <template v-if="activeMenu === 'preference'">
          <div class="mb-8">
            <h1 class="mb-2 text-2xl font-bold text-foreground">{{ t('settings.preference.title') }}</h1>
            <p class="text-sm text-muted-foreground">{{ t('settings.preference.subtitle') }}</p>
          </div>
          <div class="mb-8">
            <h2 class="mb-2 text-lg font-semibold text-foreground">{{ t('settings.preference.theme_settings') }}</h2>
            <p class="mb-4 text-sm text-muted-foreground">{{ t('settings.preference.theme_desc') }}</p>
            <UiSettingRow :label="t('settings.preference.theme')">
              <UiSelect v-model="themePref.mode" @update:model-value="(v: any) => setThemeMode(v)">
                <UiSelectTrigger class="w-[180px] h-8">
                  <UiSelectValue />
                </UiSelectTrigger>
                <UiSelectContent>
                  <UiSelectItem value="light">{{ t('settings.preference.theme_light') }}</UiSelectItem>
                  <UiSelectItem value="dark">{{ t('settings.preference.theme_dark') }}</UiSelectItem>
                  <UiSelectItem value="auto">{{ t('settings.preference.theme_auto') }}</UiSelectItem>
                </UiSelectContent>
              </UiSelect>
            </UiSettingRow>
          </div>
        </template>
      </main>

      <!-- Right help -->
      <aside class="max-md:hidden">
        <div class="sticky top-24 rounded-lg border border-border bg-card p-5">
          <h3 class="mb-3 text-base font-bold text-foreground">{{ t('settings.help.title') }}</h3>
          <h4 class="mb-2 text-sm font-semibold text-foreground">{{ t('settings.help.account_password') }}</h4>
          <ul class="mb-4 space-y-2 text-xs">
            <li class="text-muted-foreground hover:text-foreground">1. {{ t('settings.help.q1') }}</li>
            <li class="text-muted-foreground hover:text-foreground">2. {{ t('settings.help.q2') }}</li>
            <li class="text-muted-foreground hover:text-foreground">3. {{ t('settings.help.q3') }}</li>
            <li class="text-muted-foreground hover:text-foreground">4. {{ t('settings.help.q4') }}</li>
          </ul>
        </div>
      </aside>
    </div>
  </div>
</template>
