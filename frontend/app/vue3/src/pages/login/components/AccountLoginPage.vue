<script setup lang="ts">
import {ref} from 'vue';
import {$t} from '@/locales';

const username = ref('');
const password = ref('');
const rememberMe = ref(false);
const usernameInputId = 'login-account-username'
const passwordInputId = 'login-account-password'

function handleLogin() {
  if (!username.value || !password.value) {
    return;
  }

  console.log('登录信息：', {
    username: username.value,
    password: password.value,
    rememberMe: rememberMe.value,
  });
}

function handleForgotPassword() {
  // TODO: 跳转到忘记密码页面
  console.log('忘记密码');
}
</script>

<template>
  <div class="login-form">
    <div class="form-group">
      <label :for="usernameInputId">{{ $t('authentication.login.username') }}</label>
      <n-input
        :id="usernameInputId"
        v-model:value="username"
        :placeholder="$t('authentication.register.input_username')"
        type="text"
        autocomplete="username"
        clearable
      />
    </div>
    <div class="form-group">
      <label :for="passwordInputId">{{ $t('authentication.login.password') }}</label>
      <n-input
        :id="passwordInputId"
        v-model:value="password"
        type="password"
        :placeholder="$t('authentication.login.placeholder_password')"
        autocomplete="current-password"
        clearable
        show-password-on="click"
        @keyup.enter="handleLogin"
      />
    </div>
    <div class="form-options">
      <n-checkbox v-model:checked="rememberMe">{{ $t('authentication.login.remember_me') }}</n-checkbox>
      <n-button text type="primary" @click="handleForgotPassword">
        {{ $t('authentication.login.forgot_password') }}
      </n-button>
    </div>
    <n-button type="primary" block size="large" class="login-button" @click="handleLogin">
      {{ $t('authentication.login.login') }}
    </n-button>
  </div>
</template>

<style scoped lang="less">
.login-form {
  display: flex;
  flex-direction: column;

  :deep(.n-input) {
    --n-color: var(--auth-card-bg);
    --n-color-focus: var(--auth-card-bg);
    --n-text-color: var(--auth-text-primary);
    --n-placeholder-color: var(--auth-text-secondary);
    --n-border: 1px solid var(--auth-border);
    --n-border-hover: 1px solid var(--color-brand);
    --n-border-focus: 1px solid var(--color-brand);
    --n-box-shadow-focus: var(--auth-focus-ring);
  }

  :deep(.n-checkbox) {
    --n-text-color: var(--auth-text-secondary);
  }
}

.form-group {
  margin-bottom: 1.5rem;

  label {
    display: block;
    font-weight: 500;
    margin-bottom: 0.5rem;
    font-size: 0.95rem;
    color: var(--auth-text-primary);
  }
}

.form-options {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 1.5rem;
  font-size: 0.95rem;

  :deep(.n-button--text-type) {
    --n-text-color-text: var(--color-brand);
    --n-text-color-text-hover: var(--auth-tab-active-text);
  }
}

.login-button {
  font-weight: 600;
  height: 44px;
  margin-bottom: 1.5rem;
}
</style>
