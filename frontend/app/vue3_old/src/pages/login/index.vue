<script setup lang="ts">
import { definePage } from 'unplugin-vue-router/runtime'
import { computed, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { $t } from '@/locales'

import AccountLoginPage from './components/AccountLoginPage.vue'
import EmailLoginPage from './components/EmailLoginPage.vue'
import PhoneLoginPage from './components/PhoneLoginPage.vue'
import OtherLoginPage from './components/OtherLoginPage.vue'
import { MoonOutline, SunnyOutline, LanguageOutline } from '@vicons/ionicons5'
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
          <template #icon>
            <n-icon>
              <LanguageOutline/>
            </n-icon>
          </template>
          <span class="lang-text">{{ getLanguageLabel() }}</span>
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
  --auth-button-gradient: linear-gradient(135deg, var(--color-brand) 0%, #a855f7 100%);
  --auth-button-shadow: 0 4px 16px rgba(168, 85, 247, 0.3);

  display: flex;
  min-height: 100vh;
  background: var(--auth-page-bg);
  overflow-x: hidden;

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
    position: relative;
    overflow: hidden;

    // 添加装饰性背景元素
    &::before {
      content: '';
      position: absolute;
      top: -50%;
      right: -20%;
      width: 80%;
      height: 200%;
      background: radial-gradient(circle, rgba(255, 255, 255, 0.1) 0%, transparent 70%);
      animation: float 20s ease-in-out infinite;
    }

    &::after {
      content: '';
      position: absolute;
      bottom: -30%;
      left: -10%;
      width: 60%;
      height: 150%;
      background: radial-gradient(circle, rgba(255, 255, 255, 0.08) 0%, transparent 70%);
      animation: float 25s ease-in-out infinite reverse;
    }

    .brand {
      text-align: center;
      margin-bottom: 3rem;

      .brand-logo {
        width: 80px;
        height: 80px;
        margin-bottom: 1rem;
        animation: logoFloat 3s ease-in-out infinite;
        filter: drop-shadow(0 4px 12px rgba(0, 0, 0, 0.2));
      }

      .brand-title {
        font-size: 2.5rem;
        font-weight: 700;
        margin-bottom: 0.5rem;
        text-shadow: 0 2px 8px rgba(0, 0, 0, 0.2);
        animation: fadeInUp 0.8s ease-out;
      }

      .brand-subtitle {
        font-size: 1.1rem;
        opacity: 0.92;
        animation: fadeInUp 0.8s ease-out 0.2s both;
      }
    }

    .features-list {
      display: flex;
      flex-direction: column;
      gap: 1rem;
      position: relative;
      z-index: 1;

      .feature-item {
        display: flex;
        align-items: center;
        gap: 1rem;
        padding: 1rem 0;
        font-size: 1.05rem;
        animation: fadeInLeft 0.6s ease-out;
        animation-fill-mode: both;

        &:nth-child(1) { animation-delay: 0.1s; }
        &:nth-child(2) { animation-delay: 0.2s; }
        &:nth-child(3) { animation-delay: 0.3s; }
        &:nth-child(4) { animation-delay: 0.4s; }

        span:first-child {
          color: rgba(255, 255, 255, 0.9);
          font-size: 1.2rem;
          filter: drop-shadow(0 2px 4px rgba(0, 0, 0, 0.2));
        }
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
      border-radius: 20px;
      background: var(--auth-card-bg);
      box-shadow: 0 8px 32px rgba(0, 0, 0, 0.08),
                  0 2px 16px rgba(0, 0, 0, 0.04),
                  inset 0 1px 0 rgba(255, 255, 255, 0.1);
      backdrop-filter: blur(10px);
      transition: all 0.4s cubic-bezier(0.4, 0, 0.2, 1);

      &:hover {
        box-shadow: 0 12px 40px rgba(0, 0, 0, 0.12),
                    0 4px 20px rgba(0, 0, 0, 0.06),
                    inset 0 1px 0 rgba(255, 255, 255, 0.15);
        transform: translateY(-2px);
      }

      :deep(.n-card__content) {
        color: var(--auth-text-primary);
        padding: 32px;
      }

      .card-header {
        text-align: center;
        margin-bottom: 2rem;

        h2 {
          font-size: 1.875rem;
          margin-bottom: 0.5rem;
          color: var(--auth-text-primary);
          font-weight: 800;
          letter-spacing: -0.5px;
          background: var(--auth-button-gradient);
          -webkit-background-clip: text;
          -webkit-text-fill-color: transparent;
          background-clip: text;
          animation: fadeInDown 0.6s ease-out;
        }

        p {
          color: var(--auth-text-secondary);
          font-size: 0.975rem;
          animation: fadeInUp 0.6s ease-out 0.2s both;
        }
      }
    }
  }
}

.login-tabs {
  margin-bottom: 1.5rem;

  :deep(.n-tabs-nav--segment-type) {
    background: var(--auth-tab-bg);
    border-radius: 12px;
    padding: 4px;
  }

  :deep(.n-tabs-tab) {
    border-radius: 8px;
    transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
    font-weight: 600;
    padding: 10px 16px;

    &.n-tabs-tab--active {
      background: var(--auth-tab-active-bg);
      color: var(--auth-tab-active-text);
      box-shadow: 0 2px 8px rgba(102, 126, 234, 0.15);
      transform: translateY(-1px);
    }

    &:hover:not(.n-tabs-tab--active) {
      background: rgba(102, 126, 234, 0.1);
    }
  }
}

.register-section {
  text-align: center;
  padding: 1.5rem 0;
  border-top: 1px solid var(--auth-border);
  border-bottom: 1px solid var(--auth-border);
  background: linear-gradient(135deg, rgba(102, 126, 234, 0.02) 0%, rgba(168, 85, 247, 0.02) 100%);

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
  font-weight: 600;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);

  &:hover {
    transform: translateY(-1px);
  }
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
  border-radius: 16px;
  backdrop-filter: blur(12px);
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.1), 0 2px 8px rgba(0, 0, 0, 0.06);
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  border: 1px solid rgba(255, 255, 255, 0.3);

  &:hover {
    background: rgba(255, 255, 255, 0.95);
    box-shadow: 0 8px 24px rgba(0, 0, 0, 0.15), 0 4px 12px rgba(0, 0, 0, 0.1);
    transform: translateY(-2px);
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
    padding: 8px 14px;
    font-size: 13px;
    font-weight: 600;
    transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
    display: flex;
    align-items: center;
    gap: 6px;
    border-radius: 12px;
  
    &:hover {
      transform: translateY(-2px);
      box-shadow: 0 4px 12px rgba(102, 126, 234, 0.2);
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

html.dark {
  .login-controls {
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

@media (max-width: 480px) {
  .login-controls {
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

// 动画关键帧定义
@keyframes slideDown {
  from {
    opacity: 0;
    transform: translateY(-30px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

@keyframes fadeInUp {
  from {
    opacity: 0;
    transform: translateY(20px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

@keyframes fadeInDown {
  from {
    opacity: 0;
    transform: translateY(-20px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

@keyframes fadeInLeft {
  from {
    opacity: 0;
    transform: translateX(-20px);
  }
  to {
    opacity: 1;
    transform: translateX(0);
  }
}

@keyframes float {
  0%, 100% {
    transform: translate(0, 0) rotate(0deg);
  }
  33% {
    transform: translate(10px, -15px) rotate(5deg);
  }
  66% {
    transform: translate(-10px, 10px) rotate(-5deg);
  }
}

@keyframes logoFloat {
  0%, 100% {
    transform: translateY(0) scale(1);
  }
  50% {
    transform: translateY(-8px) scale(1.02);
  }
}
</style>
