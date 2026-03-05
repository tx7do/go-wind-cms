<script setup lang="ts">
import {ref, computed} from 'vue';
import {useRoute} from 'vue-router';
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
const route = useRoute();

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

function handleClickLogo() {
  if (route.path === '/') {
    return;
  }
  navigateTo('/');
}

</script>

<template>
  <div class="top-bar">
    <div
      class="logo-section"
      role="button"
      tabindex="0"
      aria-label="Go to homepage"
      @click="handleClickLogo"
      @keydown.enter.prevent="handleClickLogo"
      @keydown.space.prevent="handleClickLogo"
    >
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
  gap: 16px;
  padding: 10px 32px;
  color: var(--color-text-primary);
  height: 65px;
  overflow: hidden;
  width: 100%;
}

.logo-section {
  display: flex;
  align-items: center;
  gap: 10px;
  cursor: pointer;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  border-radius: var(--radius-sm);
  padding: 8px 12px;
  flex-shrink: 0;

  &:hover {
    background-color: rgba(102, 126, 234, 0.12);

    .site-name {
      color: var(--color-brand);
    }
  }
}

.logo {
  height: 45px;
  flex-shrink: 0;
  transition: transform 0.3s;

  .logo-section:hover & {
    transform: scale(1.05);
  }
}

.site-name {
  font-size: 17px;
  font-weight: 700;
  color: var(--color-brand);
  white-space: nowrap;
  letter-spacing: 0.8px;
  transition: all 0.3s;
}

.actions {
  display: flex;
  align-items: center;

  // 分割线
  :deep(.n-divider) {
    background-color: var(--color-border);
    opacity: 0.5;
  }
}

.icon-btn {
  font-size: 24px;
  --n-text-color: var(--header-control-text-muted);
  --n-text-color-hover: var(--color-brand);
  --n-text-color-pressed: var(--color-brand);
  transition: all 0.3s;
  border-radius: 8px;

  &:hover {
    background: rgba(102, 126, 234, 0.08) !important;
  }
}

:deep(.icon-btn.n-button.n-button--text-type:hover) {
  background: rgba(102, 126, 234, 0.08) !important;
}

.lang-btn,
.theme-btn {
  --n-color: var(--header-control-bg);
  --n-color-hover: rgba(102, 126, 234, 0.08);
  --n-color-pressed: rgba(102, 126, 234, 0.12);
  --n-border: 1px solid var(--header-control-border);
  --n-border-hover: 1px solid var(--color-brand);
  --n-border-pressed: 1px solid var(--color-brand);
  --n-text-color: var(--header-control-text);
  --n-text-color-hover: var(--color-brand);
  --n-text-color-pressed: var(--color-brand);
  font-weight: 500;
  transition: all 0.3s;
  min-width: 80px;

  &:hover {
    box-shadow: 0 2px 8px rgba(102, 126, 234, 0.15);
  }
}

.header-login-btn,
.header-register-btn {
  border: 1px solid transparent;
  font-weight: 600;
  font-size: 14px;
  transition: all 0.3s;
}

.header-login-btn {
  --n-color: var(--header-login-bg);
  --n-color-hover: rgba(102, 126, 234, 0.1);
  --n-color-pressed: rgba(102, 126, 234, 0.15);
  --n-border: 1px solid var(--header-login-border);
  --n-border-hover: 1px solid var(--color-brand);
  --n-border-pressed: 1px solid var(--color-brand);
  --n-text-color: var(--header-login-text);
  --n-text-color-hover: var(--color-brand);
  --n-text-color-pressed: var(--color-brand);

  &:hover {
    box-shadow: 0 2px 8px rgba(102, 126, 234, 0.1);
  }
}

.header-register-btn {
  --n-color: var(--header-register-bg);
  --n-color-hover: var(--color-brand);
  --n-color-pressed: var(--color-brand);
  --n-border: 1px solid var(--header-register-border);
  --n-border-hover: 1px solid var(--color-brand);
  --n-border-pressed: 1px solid var(--color-brand);
  --n-text-color: #fff;
  --n-text-color-hover: #fff;
  --n-text-color-pressed: #fff;

  &:hover {
    box-shadow: 0 4px 12px rgba(102, 126, 234, 0.3);
  }
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

// Responsive
@media (max-width: 1024px) {
  .top-bar {
    padding: 12px 24px;
    gap: 12px;
    height: 60px;
  }

  .logo {
    height: 42px;
  }

  .site-name {
    font-size: 15px;
  }
}

@media (max-width: 768px) {
  .top-bar {
    padding: 10px 16px;
    gap: 8px;
    height: 56px;
  }

  .logo {
    height: 36px;
  }

  .site-name {
    font-size: 13px;
    display: none;
  }

  .lang-btn,
  .theme-btn {
    min-width: auto;
    padding: 0 8px !important;

    :deep(.n-button__content) {
      font-size: 12px;
    }
  }

  .header-login-btn,
  .header-register-btn {
    font-size: 12px;
    padding: 0 12px !important;
  }
}
</style>
