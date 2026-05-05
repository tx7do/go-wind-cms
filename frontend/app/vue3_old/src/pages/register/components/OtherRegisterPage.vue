<script setup lang="ts">
import {
  LogoGoogle,
  LogoGithub
} from '@vicons/ionicons5'
import {$t} from '@/locales';

function handleButtonGoogleRegister() {
  // TODO: 实现 Google 登录逻辑
}

function handleButtonGithubRegister() {
  // TODO: 实现 GitHub 登录逻辑
}
</script>

<template>
  <div class="other-register">

    <!-- Social Buttons -->
    <n-button
      block
      size="large"
      class="social-button"
      @click="handleButtonGoogleRegister"
    >
      <template #icon>
        <n-icon :component="LogoGoogle"/>
      </template>
      {{ $t('authentication.login.social_google') }}
    </n-button>

    <n-button
      block
      size="large"
      class="social-button"
      @click="handleButtonGithubRegister"
    >
      <template #icon>
        <n-icon :component="LogoGithub"/>
      </template>
      {{ $t('authentication.login.social_github') }}
    </n-button>
  </div>
</template>

<style scoped lang="less">
.other-register {
  display: flex;
  flex-direction: column;
  gap: 1rem;
}

.social-button {
  font-weight: 600;
  height: 44px;
  background: var(--auth-card-bg);
  color: var(--auth-text-primary);
  border: 1px solid var(--auth-border);
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 0.5rem;
  font-size: 0.95rem;
  border-radius: 12px;

  :deep(.n-icon) {
    font-size: 18px;
  }

  &:hover {
    border-color: var(--color-brand);
    color: #fff;
    background: linear-gradient(135deg, var(--color-brand) 0%, #a855f7 100%);
    box-shadow: 0 4px 12px rgba(168, 85, 247, 0.3);
    transform: translateY(-2px);
  }

  &:focus-visible {
    box-shadow: var(--auth-focus-ring);
  }

  &:active {
    transform: translateY(0);
  }
}
</style>
