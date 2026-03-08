<script setup lang="ts">
import {computed, ref} from "vue";
import {definePage} from "unplugin-vue-router/runtime";
import {useRoute, useRouter} from 'vue-router'
import {$t} from '@/locales';

import AccountRegisterPage from "./components/AccountRegisterPage.vue";
import EmailRegisterPage from "./components/EmailRegisterPage.vue";
import PhoneRegisterPage from "./components/PhoneRegisterPage.vue";
import OtherRegisterPage from "./components/OtherRegisterPage.vue";
import {MoonOutline, SunnyOutline, LanguageOutline} from '@vicons/ionicons5'
import {isDark, toggleDark} from '@/composables/dark'
import {
  languageColumns,
  languageShorts,
  i18n,
  loadLocaleMessages,
  type SupportedLanguagesType
} from '@/locales'
import {updatePreferences} from "@/preferences";

definePage({
  name: 'register',
  meta: {
    level: 1,
    hideLayout: true, // 隐藏全局 Header 和 Footer
    forceLight: false, // 遵循全局主题设置
  },
})

const router = useRouter()
const route = useRoute()
const activeTab = ref('account');

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

const handleLoginClick = () => {
  router.push('/login')
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
  <div class="register-page" :class="{ 'force-light': forceLight }">
    <div class="register-controls">
      <n-dropdown trigger="hover" size="huge" :options="languageColumns"
                  @select="handleSelectLanguage">
        <n-button round class="control-btn" :aria-label="$t('navbar.top.language')">
          <template #icon>
            <n-icon>
              <LanguageOutline/>
            </n-icon>
          </template>
          <span class="lang-text">{{ getLanguageLabel() }}</span>
        </n-button>
      </n-dropdown>
      <n-button round class="control-btn"
                :aria-label="isDark ? $t('navbar.top.light_mode') : $t('navbar.top.dark_mode')"
                @click="handleToggleTheme">
        <template #icon>
          <n-icon>
            <component :is="isDark ? SunnyOutline : MoonOutline"/>
          </n-icon>
        </template>
      </n-button>
    </div>

    <div class="register-left">
      <div class="brand">
        <img src="/logo.png" :alt="$t('authentication.login.logo_alt')" class="brand-logo"/>
        <h1 class="brand-title">{{ $t('authentication.login.brand_title') }}</h1>
        <p class="brand-subtitle">{{ $t('authentication.login.brand_subtitle') }}</p>
      </div>

      <div class="benefits-list">
        <div class="benefit-item">
          <span>✓</span>
          <span>{{ $t('authentication.login.feature_projects') }}</span>
        </div>
        <div class="benefit-item">
          <span>✓</span>
          <span>{{ $t('authentication.login.feature_isolation') }}</span>
        </div>
        <div class="benefit-item">
          <span>✓</span>
          <span>{{ $t('authentication.login.feature_permissions') }}</span>
        </div>
        <div class="benefit-item">
          <span>✓</span>
          <span>{{ $t('authentication.login.feature_analytics') }}</span>
        </div>
      </div>
    </div>

    <div class="register-right">
      <n-card :bordered="false" class="register-card">
        <template #header>
          <div class="card-header">
            <h2>{{ $t('authentication.register.title') }}</h2>
            <p>{{ $t('authentication.register.register_with') }}</p>
          </div>
        </template>

        <n-tabs v-model:value="activeTab" type="segment" animated class="register-tabs">
          <n-tab-pane name="account" :tab="$t('authentication.login.tab_account')">
            <AccountRegisterPage/>
          </n-tab-pane>

          <n-tab-pane name="email" :tab="$t('authentication.login.tab_email')">
            <EmailRegisterPage/>
          </n-tab-pane>

          <n-tab-pane name="phone" :tab="$t('authentication.login.tab_phone')">
            <PhoneRegisterPage/>
          </n-tab-pane>

          <n-tab-pane name="other" :tab="$t('authentication.login.tab_other')">
            <OtherRegisterPage/>
          </n-tab-pane>
        </n-tabs>

        <div class="login-section">
          <p>{{ $t('authentication.register.already_have_account') }}
            <n-button text type="primary" @click="handleLoginClick">
              {{ $t('authentication.register.login_now') }}
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
            <n-button text type="primary" @click="handleTermsClick">
              {{ $t('authentication.login.terms_of_service') }}
            </n-button>
            {{ $t('authentication.login.terms_and') }}
            <n-button text type="primary" @click="handlePrivacyClick">
              {{ $t('authentication.login.privacy_policy') }}
            </n-button>
          </small>
        </div>
      </n-card>
    </div>
  </div>
</template>

<style scoped lang="less">
.register-page {
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

  .register-left {
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

    .benefits-list {
      display: flex;
      flex-direction: column;
      gap: 1rem;

      .benefit-item {
        display: flex;
        align-items: center;
        gap: 1rem;
        padding: 1rem 0;
        font-size: 1.05rem;
      }
    }
  }

  .register-right {
    flex: 1;
    display: flex;
    align-items: center;
    justify-content: center;
    padding: 2rem;
    background: var(--auth-right-bg);

    .register-card {
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

.register-controls {
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
    --n-color: rgba(255, 255, 255, 0.5);
    --n-color-hover: rgba(102, 126, 234, 0.15);
    --n-color-pressed: rgba(102, 126, 234, 0.2);
    --n-border: 1px solid rgba(102, 126, 234, 0.15);
    --n-border-hover: 1px solid rgba(102, 126, 234, 0.3);
    --n-text-color: #4a5568;
    --n-text-color-hover: var(--color-brand);
    min-width: auto;
    padding: 6px 12px;
    font-size: 13px;
    font-weight: 500;
    transition: all 0.2s ease;
    display: flex;
    align-items: center;
    gap: 6px;

    &:hover {
      transform: translateY(-1px);
      box-shadow: 0 2px 8px rgba(102, 126, 234, 0.15);
    }

    :deep(.n-icon) {
      font-size: 18px;
      color: #4a5568;
    }

    &:hover :deep(.n-icon) {
      color: var(--color-brand);
    }

    .lang-text {
      font-size: 13px;
      white-space: nowrap;
      color: #4a5568;
    }

    &:hover .lang-text {
      color: var(--color-brand);
    }
  }
}

.register-tabs {
  margin-bottom: 1.5rem;

  :deep(.n-tabs-nav--segment-type) {
    background: var(--auth-tab-bg);
  }

  :deep(.n-tabs-tab.n-tabs-tab--active) {
    background: var(--auth-tab-active-bg);
    color: var(--auth-tab-active-text);
  }
}

.login-section {
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

:deep(.login-section .n-button--text-type),
:deep(.back-home .n-button--text-type),
:deep(.terms .n-button--text-type) {
  --n-text-color-text: var(--color-brand);
  --n-text-color-text-hover: var(--auth-tab-active-text);
  --n-text-color-text-pressed: var(--auth-tab-active-text);
}

:deep(.login-section .n-button--text-type:focus-visible),
:deep(.back-home .n-button--text-type:focus-visible),
:deep(.terms .n-button--text-type:focus-visible) {
  box-shadow: var(--auth-focus-ring);
  border-radius: 6px;
}

html.dark {
  .register-controls {
    background: rgba(27, 31, 39, 0.95);
    box-shadow: 0 2px 12px rgba(0, 0, 0, 0.4);

    &:hover {
      background: rgba(27, 31, 39, 0.98);
      box-shadow: 0 4px 16px rgba(0, 0, 0, 0.5);
    }

    .control-btn {
      --n-color: rgba(42, 49, 64, 0.5);
      --n-color-hover: rgba(102, 126, 234, 0.2);
      --n-color-pressed: rgba(102, 126, 234, 0.25);
      --n-border: 1px solid rgba(102, 126, 234, 0.2);
      --n-border-hover: 1px solid rgba(102, 126, 234, 0.4);
      --n-text-color: #aab3c2;
      --n-text-color-hover: #b8c5ff;

      :deep(.n-icon) {
        color: #aab3c2;
      }

      &:hover :deep(.n-icon) {
        color: #b8c5ff;
      }

      .lang-text {
        color: #aab3c2;
      }

      &:hover .lang-text {
        color: #b8c5ff;
      }

      &:hover {
        box-shadow: 0 2px 8px rgba(102, 126, 234, 0.25);
      }
    }
  }

  .register-page:not(.force-light) {
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

@media (max-width: 1024px) {
  .register-page {
    flex-direction: column;

    .register-left {
      padding: 3rem 2rem;
      min-height: auto;

      .brand {
        margin-bottom: 2rem;

        .brand-logo {
          width: 60px;
          height: 60px;
        }

        .brand-title {
          font-size: 2rem;
        }

        .brand-subtitle {
          font-size: 1rem;
        }
      }

      .benefits-list {
        gap: 0.75rem;

        .benefit-item {
          padding: 0.75rem 0;
          font-size: 1rem;
        }
      }
    }

    .register-right {
      flex: 1;
      padding: 1.5rem;

      .register-card {
        max-width: 100%;
      }
    }
  }

  .register-controls {
    top: 16px;
    right: 16px;
    padding: 6px;
    gap: 6px;

    .control-btn {
      padding: 6px 10px;
      font-size: 12px;

      :deep(.n-icon) {
        font-size: 14px;
      }
    }
  }
}

@media (max-width: 768px) {
  .register-page {
    .register-left {
      padding: 2rem 1.5rem;
      min-height: auto;

      .brand {
        margin-bottom: 1.5rem;

        .brand-logo {
          width: 50px;
          height: 50px;
          margin-bottom: 0.75rem;
        }

        .brand-title {
          font-size: 1.75rem;
          margin-bottom: 0.4rem;
        }

        .brand-subtitle {
          font-size: 0.95rem;
        }
      }

      .benefits-list {
        gap: 0.5rem;

        .benefit-item {
          padding: 0.5rem 0;
          font-size: 0.95rem;
          gap: 0.75rem;
        }
      }
    }

    .register-right {
      padding: 1rem;

      .register-card {
        border: none;
        box-shadow: none;

        .card-header {
          margin-bottom: 1.25rem;

          h2 {
            font-size: 1.5rem;
            margin-bottom: 0.4rem;
          }

          p {
            font-size: 0.9rem;
          }
        }
      }
    }
  }

  .register-tabs {
    margin-bottom: 1.25rem;

    :deep(.n-tabs-nav) {
      font-size: 13px;
    }

    :deep(.n-tabs-tab) {
      padding: 8px 12px !important;
    }
  }

  .register-controls {
    top: 12px;
    right: 12px;
    padding: 6px;
    gap: 4px;

    .control-btn {
      padding: 6px 10px;
      font-size: 12px;

      :deep(.n-icon) {
        font-size: 14px;
      }
    }
  }
}

@media (max-width: 480px) {
  .register-page {
    .register-left {
      padding: 1.5rem 1rem;
      min-height: auto;

      .brand {
        margin-bottom: 1rem;

        .brand-logo {
          width: 40px;
          height: 40px;
          margin-bottom: 0.5rem;
        }

        .brand-title {
          font-size: 1.5rem;
          margin-bottom: 0.3rem;
        }

        .brand-subtitle {
          font-size: 0.85rem;
        }
      }

      .benefits-list {
        gap: 0.4rem;

        .benefit-item {
          padding: 0.4rem 0;
          font-size: 0.85rem;
          gap: 0.5rem;

          span:first-child {
            font-size: 1rem;
          }
        }
      }
    }

    .register-right {
      padding: 0.75rem;

      .register-card {
        padding: 0;

        :deep(.n-card__content) {
          padding: 1rem;
        }

        .card-header {
          margin-bottom: 1rem;

          h2 {
            font-size: 1.3rem;
            margin-bottom: 0.3rem;
          }

          p {
            font-size: 0.85rem;
          }
        }
      }
    }
  }

  .register-tabs {
    margin-bottom: 1rem;

    :deep(.n-tabs-nav) {
      font-size: 12px;
    }

    :deep(.n-tabs-tab) {
      padding: 6px 8px !important;
      font-size: 11px;
    }
  }

  .login-section {
    padding: 1rem 0;

    p {
      font-size: 0.85rem;
      line-height: 1.6;
    }

    :deep(.n-button) {
      font-size: 0.85rem;
    }
  }

  .back-home {
    padding: 0.75rem 0;

    :deep(.n-button) {
      font-size: 0.85rem;
    }
  }

  .terms {
    margin-top: 0.75rem;

    small {
      font-size: 0.75rem;
      line-height: 1.5;
    }

    :deep(.n-button) {
      font-size: 0.75rem;
    }
  }

  .register-controls {
    top: 8px;
    right: 8px;
    padding: 4px;
    gap: 4px;
    border-radius: 8px;

    .control-btn {
      padding: 4px 8px;
      font-size: 11px;
      min-width: 32px;
      height: 32px;

      :deep(.n-icon) {
        font-size: 14px;
      }

      .lang-text {
        display: none;
      }
    }
  }
}
</style>
