<script setup lang="ts">
import { definePage } from 'unplugin-vue-router/runtime'
import { computed, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { $t } from '@/locales'

import AccountLoginPage from './components/AccountLoginPage.vue'
import EmailLoginPage from './components/EmailLoginPage.vue'
import PhoneLoginPage from './components/PhoneLoginPage.vue'
import OtherLoginPage from './components/OtherLoginPage.vue'
import { MoonOutline, SunnyOutline } from '@vicons/ionicons5'
import { isDark, toggleDark } from '@/composables/dark'
import { languageColumns, languageShorts, i18n, loadLocaleMessages, type SupportedLanguagesType } from '@/locales'
import { updatePreferences } from "@/preferences";

definePage({
  name: 'login',
  meta: {
    level: 1,
    hideLayout: true, // 隐藏全局 Header 和 Footer
    forceLight: false, // 遵循全局主题设置
  },
})

const router = useRouter()
const route = useRoute()
const activeTab = ref('account')

const forceLight = computed(() => route.meta.forceLight === true)

function getLanguageLabel() {
  const lang = i18n.global.locale.value;
  const ret = languageShorts.find((item) => item.key === lang);
  return ret?.label ?? lang;
}

async function handleSelectLanguage(key: string | number) {
  const locale = key as SupportedLanguagesType;
  if (locale === i18n.global.locale.value) {
    return;
  }

  updatePreferences({
    app: {
      locale,
    },
  });

  await loadLocaleMessages(locale);
}

function handleToggleTheme() {
  toggleDark();
}

const handleRegisterClick = () => {
  router.push('/register')
}

const handleBackHome = () => {
  router.push('/')
}

const handleTermsClick = () => {
  router.push('/terms')
}

const handlePrivacyClick = () => {
  router.push('/privacy')
}
</script>

<template>
  <div class="login-page" :class="{ 'force-light': forceLight }">
    <div class="login-left">
      <div class="brand">
        <img src="/logo.png" :alt="$t('authentication.login.logo_alt')" class="brand-logo" />
        <h1 class="brand-title">{{ $t('authentication.login.brand_title') }}</h1>
        <p class="brand-subtitle">{{ $t('authentication.login.brand_subtitle') }}</p>
      </div>

      <div class="features-list">
        <div class="feature-item">
          <span>✓</span>
          <span>{{ $t('authentication.login.feature_projects') }}</span>
        </div>
        <div class="feature-item">
          <span>✓</span>
          <span>{{ $t('authentication.login.feature_isolation') }}</span>
        </div>
        <div class="feature-item">
          <span>✓</span>
          <span>{{ $t('authentication.login.feature_permissions') }}</span>
        </div>
        <div class="feature-item">
          <span>✓</span>
          <span>{{ $t('authentication.login.feature_analytics') }}</span>
        </div>
      </div>
    </div>

    <div class="login-controls">
      <n-dropdown trigger="hover" size="huge" :options="languageColumns" @select="handleSelectLanguage">
        <n-button round class="control-btn" :aria-label="$t('navbar.top.language')">
          {{ getLanguageLabel() }}
        </n-button>
      </n-dropdown>
      <n-button round class="control-btn" :aria-label="isDark ? $t('navbar.top.light_mode') : $t('navbar.top.dark_mode')" @click="handleToggleTheme">
        <template #icon>
          <n-icon>
            <component :is="isDark ? SunnyOutline : MoonOutline"/>
          </n-icon>
        </template>
      </n-button>
    </div>

    <div class="login-right">
      <n-card :bordered="false" class="login-card">
        <template #header>
          <div class="card-header">
            <h2>{{ $t('authentication.login.title') }}</h2>
            <p>{{ $t('authentication.login.login_with') }}</p>
          </div>
        </template>

        <n-tabs v-model:value="activeTab" type="segment" animated class="login-tabs">
          <n-tab-pane name="account" :tab="$t('authentication.login.tab_account')">
            <AccountLoginPage />
          </n-tab-pane>

          <n-tab-pane name="email" :tab="$t('authentication.login.tab_email')">
            <EmailLoginPage />
          </n-tab-pane>

          <n-tab-pane name="phone" :tab="$t('authentication.login.tab_phone')">
            <PhoneLoginPage />
          </n-tab-pane>

          <n-tab-pane name="other" :tab="$t('authentication.login.tab_other')">
            <OtherLoginPage />
          </n-tab-pane>
        </n-tabs>

        <div class="register-section">
          <p>{{ $t('authentication.login.no_account') }}
            <n-button text type="primary" @click="handleRegisterClick">
              {{ $t('authentication.login.register_now') }}
            </n-button>
          </p>
        </div>

        <div class="back-home">
          <n-button text type="primary" @click="handleBackHome">
            ← {{ $t('authentication.login.back_home') }}
          </n-button>
        </div>

        <div class="terms">
          <small>
            {{ $t('authentication.login.terms_prefix') }}
            <n-button text type="primary" @click="handleTermsClick">{{ $t('authentication.login.terms_of_service') }}</n-button>
            {{ $t('authentication.login.terms_and') }}
            <n-button text type="primary" @click="handlePrivacyClick">{{ $t('authentication.login.privacy_policy') }}</n-button>
          </small>
        </div>
      </n-card>
    </div>
  </div>
</template>

<style scoped lang="less">
.login-page {
  --auth-page-bg: var(--color-bg);
  --auth-right-bg: var(--color-surface);
  --auth-card-bg: var(--color-surface);
  --auth-border: var(--color-border);
  --auth-text-primary: var(--color-text-primary);
  --auth-text-secondary: var(--color-text-secondary);
  --auth-brand-start: var(--color-brand);
  --auth-brand-end: var(--color-brand-accent);
  --auth-brand-text: #ffffff;
  --auth-tab-bg: rgba(102, 126, 234, 0.08);
  --auth-tab-active-bg: rgba(102, 126, 234, 0.16);
  --auth-tab-active-text: var(--color-brand);
  --auth-focus-ring: 0 0 0 3px rgba(102, 126, 234, 0.18);

  display: flex;
  min-height: 100vh;
  background: var(--auth-page-bg);

  &.force-light {
    --auth-page-bg: #f5f7fa;
    --auth-right-bg: #ffffff;
    --auth-card-bg: #ffffff;
    --auth-border: #e8e8e8;
    --auth-text-primary: #333333;
    --auth-text-secondary: #666666;
    --auth-tab-bg: rgba(102, 126, 234, 0.08);
    --auth-tab-active-bg: rgba(102, 126, 234, 0.16);
    --auth-tab-active-text: #667eea;
  }

  .login-left {
    flex: 1;
    background: linear-gradient(135deg, var(--auth-brand-start) 0%, var(--auth-brand-end) 100%);
    color: var(--auth-brand-text);
    padding: 4rem 2rem;
    display: flex;
    flex-direction: column;
    justify-content: space-between;

    .brand {
      text-align: center;
      margin-bottom: 3rem;

      .brand-logo {
        width: 80px;
        height: 80px;
        margin-bottom: 1rem;
      }

      .brand-title {
        font-size: 2.5rem;
        font-weight: 700;
        margin-bottom: 0.5rem;
      }

      .brand-subtitle {
        font-size: 1.1rem;
        opacity: 0.92;
      }
    }

    .features-list {
      display: flex;
      flex-direction: column;
      gap: 1rem;

      .feature-item {
        display: flex;
        align-items: center;
        gap: 1rem;
        padding: 1rem 0;
        font-size: 1.05rem;
      }
    }
  }

  .login-right {
    flex: 1;
    display: flex;
    align-items: center;
    justify-content: center;
    padding: 2rem;
    background: var(--auth-right-bg);

    .login-card {
      width: 100%;
      max-width: 420px;
      border: 1px solid var(--auth-border);
      border-radius: var(--radius-md);
      background: var(--auth-card-bg);

      :deep(.n-card__content) {
        color: var(--auth-text-primary);
      }

      .card-header {
        text-align: center;
        margin-bottom: 1.5rem;

        h2 {
          font-size: 1.75rem;
          margin-bottom: 0.5rem;
          color: var(--auth-text-primary);
        }

        p {
          color: var(--auth-text-secondary);
          font-size: 0.95rem;
        }
      }
    }
  }
}

.login-tabs {
  margin-bottom: 1.5rem;

  :deep(.n-tabs-nav--segment-type) {
    background: var(--auth-tab-bg);
  }

  :deep(.n-tabs-tab.n-tabs-tab--active) {
    background: var(--auth-tab-active-bg);
    color: var(--auth-tab-active-text);
  }
}

.register-section {
  text-align: center;
  padding: 1.5rem 0;
  border-top: 1px solid var(--auth-border);
  border-bottom: 1px solid var(--auth-border);

  p {
    margin: 0;
    font-size: 0.95rem;
    color: var(--auth-text-secondary);
  }
}

.back-home {
  text-align: center;
  padding: 1rem 0;

  :deep(.n-button) {
    font-size: 0.9rem;
  }
}

.terms {
  text-align: center;
  margin-top: 1rem;
  color: var(--auth-text-secondary);
  line-height: 1.6;

  small {
    font-size: 0.85rem;
  }
}

:deep(.register-section .n-button--text-type),
:deep(.back-home .n-button--text-type),
:deep(.terms .n-button--text-type) {
  --n-text-color-text: var(--color-brand);
  --n-text-color-text-hover: var(--auth-tab-active-text);
  --n-text-color-text-pressed: var(--auth-tab-active-text);
}

:deep(.register-section .n-button--text-type:focus-visible),
:deep(.back-home .n-button--text-type:focus-visible),
:deep(.terms .n-button--text-type:focus-visible) {
  box-shadow: var(--auth-focus-ring);
  border-radius: 6px;
}

html.dark {
  .login-page:not(.force-light) {
    --auth-page-bg: #121317;
    --auth-right-bg: #1b1f27;
    --auth-card-bg: #1b1f27;
    --auth-border: #2a3140;
    --auth-text-primary: #e8ecf3;
    --auth-text-secondary: #aab3c2;
    --auth-tab-bg: rgba(129, 147, 255, 0.18);
    --auth-tab-active-bg: rgba(129, 147, 255, 0.26);
    --auth-tab-active-text: #b8c5ff;
    --auth-focus-ring: 0 0 0 3px rgba(129, 147, 255, 0.28);
  }
}

.login-controls {
  position: fixed;
  top: 20px;
  right: 20px;
  display: flex;
  gap: 8px;
  z-index: 100;
  background: rgba(255, 255, 255, 0.9);
  padding: 8px;
  border-radius: 12px;
  backdrop-filter: blur(8px);
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.1);
  transition: all 0.3s ease;

  &:hover {
    background: rgba(255, 255, 255, 0.95);
    box-shadow: 0 4px 16px rgba(0, 0, 0, 0.15);
  }

  .control-btn {
    --n-color: transparent;
    --n-color-hover: rgba(102, 126, 234, 0.08);
    --n-color-pressed: rgba(102, 126, 234, 0.12);
    --n-border: 1px solid transparent;
    --n-border-hover: 1px solid rgba(102, 126, 234, 0.2);
    --n-text-color: var(--color-text-secondary);
    --n-text-color-hover: var(--color-brand);
    min-width: auto;
    padding: 6px 12px;
    font-size: 13px;
    font-weight: 500;
    transition: all 0.2s ease;

    &:hover {
      transform: translateY(-1px);
    }

    :deep(.n-icon) {
      font-size: 16px;
    }
  }
}

html.dark {
  .login-controls {
    background: rgba(27, 31, 39, 0.9);
    box-shadow: 0 2px 12px rgba(0, 0, 0, 0.3);

    &:hover {
      background: rgba(27, 31, 39, 0.95);
      box-shadow: 0 4px 16px rgba(0, 0, 0, 0.4);
    }

    .control-btn {
      --n-text-color: var(--color-text-secondary);
      --n-text-color-hover: var(--color-brand);
      --n-color-hover: rgba(102, 126, 234, 0.12);
      --n-color-pressed: rgba(102, 126, 234, 0.18);
    }
  }
}

@media (max-width: 768px) {
  .login-controls {
    top: 12px;
    right: 12px;
    padding: 6px;
    gap: 6px;

    .control-btn {
      padding: 4px 8px;
      font-size: 12px;

      :deep(.n-icon) {
        font-size: 14px;
      }
    }
  }

  .login-page {
    flex-direction: column;

    .login-left {
      padding: 2rem;
      min-height: auto;
      display: none;
    }

    .login-right {
      flex: 1;

      .login-card {
        max-width: 100%;
      }
    }
  }
}
</style>
