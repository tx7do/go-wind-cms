<script setup lang="ts">
import {
  LogoGoogle,
  LogoGithub
} from '@vicons/ionicons5'
import {$t} from '@/locales';

function handleGoogleLogin() {
  // TODO: 实现 Google 登录逻辑
  console.log('Google 登录');
}

function handleGithubLogin() {
  // TODO: 实现 GitHub 登录逻辑
  console.log('GitHub 登录');
}
</script>

<template>
  <div class="other-login">
    <n-button
      block
      size="large"
      class="social-button"
      @click="handleGithubLogin"
    >
      <template #icon>
        <n-icon :component="LogoGithub"/>
      </template>
      {{ $t('authentication.login.social_github') }}
    </n-button>

    <n-button
      block
      size="large"
      class="social-button"
      @click="handleGoogleLogin"
    >
      <template #icon>
        <n-icon :component="LogoGoogle"/>
      </template>
      {{ $t('authentication.login.social_google') }}
    </n-button>
  </div>
</template>

<style scoped lang="less">
.other-login {
  display: flex;
  flex-direction: column;
  gap: 1rem;
}

.social-button {
  font-weight: 500;
  height: 44px;
  background: var(--auth-card-bg);
  color: var(--auth-text-primary);
  border: 1px solid var(--auth-border);
  transition: all 0.3s ease;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 0.5rem;
  font-size: 0.95rem;

  :deep(.n-icon) {
    font-size: 18px;
  }

  &:hover {
    border-color: var(--color-brand);
    color: var(--color-brand);
    background: var(--auth-tab-bg);
    box-shadow: 0 2px 8px rgba(102, 126, 234, 0.12);
  }

  &:focus-visible {
    box-shadow: var(--auth-focus-ring);
  }

  &:active {
    transform: scale(0.98);
  }
}
</style>
