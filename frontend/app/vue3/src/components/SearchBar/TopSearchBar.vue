<script setup lang="ts">
import {ref, computed} from 'vue';
import {
  AlbumsOutline,
  CardOutline,
  HomeOutline,
  LayersOutline,
  LogOutOutline,
  PersonCircleOutline,
  ShieldCheckmarkOutline,
  MoonOutline,
  SunnyOutline,
} from '@vicons/ionicons5';
import {CrownOutlined} from '@vicons/antd';
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
import { isDark, toggleDark } from '@/composables/dark';

import logoImage from '@/assets/images/logo.png';

const isLogin = ref(false);

const userOptions = computed(() => {
  return [
    {
      label: $t('menu.my_profile'),
      key: 'profile',
      icon: renderIcon(PersonCircleOutline),
    },
    {
      label: $t('menu.my_business'),
      key: 'business',
      icon: renderIcon(LayersOutline),
    },
    {
      label: $t('menu.my_cards'),
      key: 'cards',
      icon: renderIcon(CardOutline),
    },
    {
      label: $t('menu.my_events'),
      key: 'events',
      icon: renderIcon(AlbumsOutline),
    },
    {
      label: $t('menu.my_account_security'),
      key: 'accountSecurity',
      icon: renderIcon(ShieldCheckmarkOutline),
    },
    {
      label: $t('menu.homepage'),
      key: 'homepage',
      icon: renderIcon(HomeOutline),
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

function handleSelectUserItem(key: string | number) {
  console.log('handleSelectUserItem:', key)
}

function handleClickButtonVip() {
  console.log('handleClickButtonVip:')
}

function handleToggleTheme() {
  toggleDark();
}
</script>

<template>
  <div class="top-bar">
    <n-image :src="logoImage" class="logo" alt="Logo" preview-disabled/>
    <SearchBar />
    <div class="actions">
      <n-space>
        <n-popover trigger="hover">
          <template #trigger>
            <n-button
              v-show="isLogin"
              class="vip-btn"
              @click="handleClickButtonVip"
            >
              <template #icon>
                <CrownOutlined/>
              </template>
              VIP
            </n-button>
          </template>
        </n-popover>
        <n-dropdown
          trigger="hover" size="huge"
          :options="userOptions" @select="handleSelectUserItem"
        >
          <n-button v-show="isLogin" text class="icon-btn" @click="navigateTo('/user')">
            <n-icon>
              <UserCircle/>
            </n-icon>
          </n-button>
        </n-dropdown>
        <n-divider :vertical="true"/>
        <n-button v-show="isLogin" text class="icon-btn" @click="navigateTo('/notifications')">
          1
        </n-button>
        <n-button v-show="!isLogin" type="info" class="header-login-btn" @click="navigateTo('/login')">
          {{ $t('navbar.top.login') }}
        </n-button>
        <n-button v-show="!isLogin" type="primary" class="header-register-btn" @click="navigateTo('/register')">
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
              <component :is="isDark ? SunnyOutline : MoonOutline" />
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

.logo {
  height: 40px;
}

.actions {
  display: flex;
  align-items: center;

  // 分割线
  :deep(.n-divider) {
    background-color: var(--color-border);
  }
}

.vip-btn {
  background: linear-gradient(90deg, #ffd700, #ffa500) !important;
  color: #111 !important;
  border: none !important;
}

.vip-btn:hover {
  opacity: 0.9;
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
