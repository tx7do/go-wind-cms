<script setup lang="ts">
import {ref, computed} from 'vue';
import {
  HomeOutline,
  LogOutOutline,
  PersonCircleOutline,
  MoonOutline,
  SunnyOutline,
} from '@vicons/ionicons5';
import {UserCircle} from '@vicons/fa';

import {
  languageColumns,
  languageShorts,
  i18n,
  $t,
  loadLocaleMessages,
  type SupportedLanguagesType
} from '@/locales';

import {navigateTo} from '@/router';
import {renderIcon} from '@/utils';
import {updatePreferences} from "@/preferences";
import {isDark, toggleDark} from '@/composables/dark';

import logoImage from '@/assets/images/logo.png';

const isLogin = ref(true);

const userOptions = computed(() => {
  return [
    {
      label: $t('menu.homepage'),
      key: 'homepage',
      icon: renderIcon(HomeOutline),
    },
    {
      label: $t('menu.my_profile'),
      key: 'profile',
      icon: renderIcon(PersonCircleOutline),
    },
    {
      type: 'divider',
      key: 'd1',
    },
    {
      label: $t('menu.logout'),
      key: 'logout',
      icon: renderIcon(LogOutOutline),
    },
  ]
});

function getLanguageLabel() {
  const lang = i18n.global.locale.value;
  const ret = languageShorts.find((item) => item.key === lang);
  return ret?.label ?? lang;
}

async function handleSelectLanguage(key: string | number) {
  console.log('handleSelectLanguage:', key)

  const locale = key as SupportedLanguagesType;
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

function handleSelectUserItem(key: string | number) {
  console.log('handleSelectUserItem:', key)
  switch (key) {
    case 'homepage':
      handleClickUserHomepage();
      break;
    case 'profile':
      handleClickSettings();
      break;
    case 'logout':
      handleClickLogout();
      break;
  }
}

function handleClickUserAvatar() {
  handleClickSettings();
}

function handleClickSettings() {
  navigateTo('/settings');
}

function handleClickUserHomepage() {
  navigateTo('/user');
}

function handleClickRegister() {
  navigateTo('/register')
}

function handleClickLogin() {
  navigateTo('/login')
}

function handleClickLogout() {
}

</script>

<template>
  <div class="top-bar">
    <div class="logo-section">
      <n-image :src="logoImage" class="logo" alt="Logo" preview-disabled/>
      <span class="site-name">{{ $t('authentication.login.brand_title') }}</span>
    </div>
    <SearchBar/>
    <div class="actions">
      <n-space>
        <n-dropdown
          trigger="hover"
          size="huge"
          :options="userOptions"
          @select="handleSelectUserItem"
        >
          <n-button v-show="isLogin" text class="icon-btn" @click="handleClickUserAvatar()">
            <n-icon>
              <UserCircle/>
            </n-icon>
          </n-button>
        </n-dropdown>
        <n-divider :vertical="true"/>
        <n-button v-show="!isLogin" type="info" class="header-login-btn"
                  @click="handleClickLogin()">
          {{ $t('navbar.top.login') }}
        </n-button>
        <n-button v-show="!isLogin" type="primary" class="header-register-btn"
                  @click="handleClickRegister()">
          {{ $t('navbar.top.register') }}
        </n-button>
        <n-dropdown trigger="hover" size="huge" :options="languageColumns"
                    @select="handleSelectLanguage">
          <n-button round class="lang-btn">
            {{ getLanguageLabel() }}
          </n-button>
        </n-dropdown>
        <n-button round class="theme-btn" @click="handleToggleTheme">
          <template #icon>
            <n-icon>
              <component :is="isDark ? SunnyOutline : MoonOutline"/>
            </n-icon>
          </template>
          {{ isDark ? $t('navbar.top.light_mode') : $t('navbar.top.dark_mode') }}
        </n-button>
      </n-space>
    </div>
  </div>
</template>

<style scoped lang="less">
.top-bar {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: var(--space-4);
  padding: 0;
  color: var(--color-text-primary);
}

.logo-section {
  display: flex;
  align-items: center;
  gap: 0.75rem;
  cursor: pointer;
  transition: all 0.2s;
  border-radius: var(--radius-sm);
  padding: 4px 8px;

  &:hover {
    background-color: rgba(102, 126, 234, 0.08);

    .site-name {
      color: var(--color-brand-accent);
    }
  }
}

.logo {
  height: 40px;
  flex-shrink: 0;
}

.site-name {
  font-size: 1.1rem;
  font-weight: 600;
  color: var(--color-brand);
  white-space: nowrap;
  letter-spacing: 0.5px;
  transition: color 0.2s;
}

.actions {
  display: flex;
  align-items: center;

  // 分割线
  :deep(.n-divider) {
    background-color: var(--color-border);
  }
}

.icon-btn {
  font-size: 28px;
  --n-text-color: var(--header-control-text-muted);
  --n-text-color-hover: var(--header-control-text);
  --n-text-color-pressed: var(--header-control-text);
}

:deep(.icon-btn.n-button.n-button--text-type:hover) {
  background: var(--header-control-hover-bg) !important;
}

.lang-btn,
.theme-btn {
  --n-color: var(--header-control-bg);
  --n-color-hover: var(--header-control-bg);
  --n-color-pressed: var(--header-control-bg);
  --n-border: 1px solid var(--header-control-border);
  --n-border-hover: 1px solid var(--color-brand);
  --n-border-pressed: 1px solid var(--color-brand);
  --n-text-color: var(--header-control-text);
  --n-text-color-hover: var(--color-brand);
  --n-text-color-pressed: var(--color-brand);
}

.header-login-btn,
.header-register-btn {
  border: 1px solid transparent;
}

.header-login-btn {
  --n-color: var(--header-login-bg);
  --n-color-hover: var(--header-login-bg);
  --n-color-pressed: var(--header-login-bg);
  --n-border: 1px solid var(--header-login-border);
  --n-border-hover: 1px solid var(--color-brand);
  --n-border-pressed: 1px solid var(--color-brand);
  --n-text-color: var(--header-login-text);
  --n-text-color-hover: var(--header-login-text);
  --n-text-color-pressed: var(--header-login-text);
}

.header-register-btn {
  --n-color: var(--header-register-bg);
  --n-color-hover: var(--header-register-bg);
  --n-color-pressed: var(--header-register-bg);
  --n-border: 1px solid var(--header-register-border);
  --n-border-hover: 1px solid var(--color-brand);
  --n-border-pressed: 1px solid var(--color-brand);
  --n-text-color: var(--header-register-text);
  --n-text-color-hover: var(--header-register-text);
  --n-text-color-pressed: var(--header-register-text);
}

:deep(.lang-btn .n-button__content),
:deep(.theme-btn .n-button__content),
:deep(.header-login-btn .n-button__content),
:deep(.header-register-btn .n-button__content) {
  color: inherit !important;
}

:deep(.lang-btn .n-icon),
:deep(.theme-btn .n-icon),
:deep(.icon-btn .n-icon) {
  color: currentColor;
}
</style>
