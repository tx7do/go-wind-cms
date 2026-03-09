<script setup lang="ts">
import {computed} from 'vue';
import {useRoute} from 'vue-router';
import {
  HomeOutline,
  LogOutOutline,
  PersonCircleOutline,
  MoonOutline,
  SunnyOutline,
  LanguageOutline,
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
import {useAccessStore, useAuthStore} from '@/stores';

import logoImage from '@/assets/images/logo.png';

const accessStore = useAccessStore();
const authStore = useAuthStore();
const isLogin = computed(() => Boolean(accessStore.accessToken));
const route = useRoute();

type UserOptionKey = 'homepage' | 'profile' | 'logout'

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

const languageLabel = computed(() => {
  const lang = i18n.global.locale.value;
  const ret = languageShorts.find((item) => item.key === lang);
  return ret?.label ?? lang;
});

function navigateIfNeeded(path: string) {
  if (route.path === path) {
    return;
  }
  navigateTo(path);
}

async function handleSelectLanguage(key: string | number) {
  const locale = key as SupportedLanguagesType;
  if (locale === i18n.global.locale.value) {
    return;
  }

  console.log('handleSelectLanguage', key);

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
  const actionMap: Record<UserOptionKey, () => void | Promise<void>> = {
    homepage: handleClickUserHomepage,
    profile: handleClickSettings,
    logout: handleClickLogout,
  };
  void actionMap[key as UserOptionKey]?.();
}

function handleClickUserAvatar() {
  handleClickSettings();
}

function handleClickSettings() {
  navigateIfNeeded('/settings');
}

function handleClickUserHomepage() {
  navigateIfNeeded('/user');
}

function handleClickRegister() {
  navigateIfNeeded('/register')
}

function handleClickLogin() {
  navigateIfNeeded('/login')
}

async function handleClickLogout() {
  await authStore.logout(false);
}

function handleClickLogo() {
  navigateIfNeeded('/');
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
        <!-- User logged in: show user avatar -->
        <n-dropdown
          v-if="isLogin"
          trigger="hover"
          size="huge"
          :options="userOptions"
          @select="handleSelectUserItem"
        >
          <n-button
            text
            class="icon-btn"
            :aria-label="$t('menu.my_profile')"
            @click="handleClickUserAvatar"
          >
            <n-icon>
              <UserCircle/>
            </n-icon>
          </n-button>
        </n-dropdown>

        <!-- User not logged in: show login/register buttons -->
        <template v-if="!isLogin">
          <n-divider :vertical="true"/>
          <n-button text class="header-login-btn"
                    :aria-label="$t('navbar.top.login')"
                    @click="handleClickLogin">
            {{ $t('navbar.top.login') }}
          </n-button>
          <n-button text class="header-register-btn"
                    :aria-label="$t('navbar.top.register')"
                    @click="handleClickRegister">
            {{ $t('navbar.top.register') }}
          </n-button>
        </template>
        <n-dropdown trigger="hover" size="huge" :options="languageColumns"
                    @select="handleSelectLanguage">
          <n-button round class="lang-btn" aria-label="Language">
            <template #icon>
              <n-icon>
                <LanguageOutline/>
              </n-icon>
            </template>
            <span class="lang-text">{{ languageLabel }}</span>
          </n-button>
        </n-dropdown>
        <n-button
          round
          class="theme-btn"
          :aria-label="isDark ? $t('navbar.top.light_mode') : $t('navbar.top.dark_mode')"
          @click="handleToggleTheme"
        >
          <template #icon>
            <n-icon>
              <component :is="isDark ? SunnyOutline : MoonOutline"/>
            </n-icon>
          </template>
          <span class="theme-text">{{
              isDark ? $t('navbar.top.light_mode') : $t('navbar.top.dark_mode')
            }}</span>
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
  padding: 12px 32px;
  color: var(--color-text-primary);
  height: 75px;
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
  height: 55px;
  flex-shrink: 0;
  transition: transform 0.3s;

  .logo-section:hover & {
    transform: scale(1.05);
  }
}

.site-name {
  font-size: 20px;
  font-weight: 700;
  color: var(--color-brand);
  white-space: nowrap;
  letter-spacing: 0.8px;
  transition: all 0.3s;
}

.actions {
  display: flex;
  align-items: center;
  flex-shrink: 0;
  height: 100%;

  // 分割线
  :deep(.n-divider) {
    background-color: var(--color-border);
    opacity: 0.5;
  }

  :deep(.n-space) {
    display: flex;
    align-items: center;
    height: 100%;
    gap: 8px;
  }
}

.icon-btn {
  font-size: 24px;
  --n-text-color: var(--header-control-text-muted);
  --n-text-color-hover: var(--color-brand);
  --n-text-color-pressed: var(--color-brand);
  transition: all 0.3s;
  border-radius: 8px;
  display: flex !important;
  align-items: center !important;
  justify-content: center !important;
  height: 40px !important;
  min-width: 40px !important;

  &:hover {
    background: rgba(102, 126, 234, 0.08) !important;
  }

  :deep(.n-icon) {
    display: flex;
    align-items: center;
    justify-content: center;
  }
}

:deep(.icon-btn.n-button.n-button--text-type:hover) {
  background: rgba(102, 126, 234, 0.08) !important;
}

.lang-btn,
.theme-btn {
  font-weight: 500;
  transition: all 0.3s;
  display: flex !important;
  align-items: center !important;
  justify-content: center !important;
  height: 40px !important;
  gap: 4px !important;
  border: 1px solid var(--color-border);
  background: var(--header-control-bg);
  color: var(--header-control-text);

  &:hover {
    background: var(--header-control-hover-bg);
    color: var(--color-brand);
    border-color: var(--color-brand);
  }

  :deep(.n-button__content) {
    display: flex;
    align-items: center;
    justify-content: center;
    gap: 4px;
    line-height: 1;
  }

  :deep(.n-icon) {
    display: flex;
    align-items: center;
    justify-content: center;
  }
}

.lang-text,
.theme-text {
  font-size: 14px;
  font-weight: 500;
}

.header-login-btn,
.header-register-btn {
  font-weight: 500;
  font-size: 14px;
  transition: all 0.3s;
  display: flex !important;
  align-items: center !important;
  justify-content: center !important;
  height: 40px !important;
  padding: 0 16px !important;
  border-radius: 8px;
  border: 1px solid rgba(102, 126, 234, 0.3);
  background: linear-gradient(135deg, rgba(102, 126, 234, 0.05) 0%, rgba(139, 92, 246, 0.05) 100%);
  color: var(--color-brand);

  &:hover {
    border-color: var(--color-brand);
    background: linear-gradient(135deg, rgba(102, 126, 234, 0.12) 0%, rgba(139, 92, 246, 0.12) 100%);
    box-shadow: 0 2px 8px rgba(102, 126, 234, 0.15);
  }
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
    height: 70px;
  }

  .logo {
    height: 50px;
  }

  .site-name {
    font-size: 18px;
  }
}

@media (max-width: 768px) {
  .top-bar {
    padding: 10px 12px;
    gap: 8px;
    height: 60px;
  }

  .logo {
    height: 40px;
  }

  .site-name {
    font-size: 15px;
    display: block;
    flex-shrink: 0;
    max-width: 100px;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
  }

  // 搜索框已在SearchBar组件中隐藏
  .actions {
    flex-shrink: 0;
    display: flex;
    align-items: center;
    gap: 4px;

    :deep(.n-space) {
      gap: 4px !important;
      display: flex;
      align-items: center;
      height: 100%;
    }
  }

  .lang-btn,
  .theme-btn {
    height: 36px !important;
    padding: 0 12px !important;
    font-size: 13px !important;
  }

  .theme-btn {
    .theme-text {
      display: none;
    }
  }

  .lang-text,
  .theme-text {
    font-size: 13px;
  }

  .header-login-btn,
  .header-register-btn {
    font-size: 12px;
    padding: 4px 10px !important;
    min-width: 60px;
    flex-shrink: 0;
  }

  .icon-btn {
    font-size: 18px;
    padding: 0 8px !important;
    min-width: 36px !important;
    width: 36px !important;
    height: 36px !important;
    display: flex;
    align-items: center;
    justify-content: center;
    flex-shrink: 0;
  }

  :deep(.n-divider) {
    height: 20px !important;
  }
}

@media (max-width: 480px) {
  .top-bar {
    padding: 8px 8px;
    gap: 6px;
    height: 56px;
  }

  .logo {
    height: 36px;
  }

  .logo-section {
    padding: 4px 8px;
    gap: 8px;
  }

  .site-name {
    font-size: 14px;
    max-width: 80px;
  }

  .actions {
    gap: 2px !important;

    :deep(.n-space) {
      gap: 2px !important;
    }
  }

  :deep(.n-divider) {
    display: none;
  }

  .header-login-btn,
  .header-register-btn {
    display: none;
  }

  .lang-btn,
  .theme-btn {
    min-width: 32px !important;
    width: 32px !important;
    height: 32px !important;
    padding: 0 !important;

    :deep(.n-icon) {
      font-size: 16px;
    }
  }

  .lang-text,
  .theme-text {
    display: none;
  }

  .icon-btn {
    padding: 0 4px !important;
    font-size: 16px;
    min-width: 32px !important;
    width: 32px !important;
    height: 32px !important;
  }
}
</style>
