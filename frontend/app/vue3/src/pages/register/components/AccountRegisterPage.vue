<script setup lang="ts">
import {ref, computed} from "vue";
import {$t} from '@/locales';

const username = ref('');
const password = ref('');
const confirmPassword = ref('');
const usernameInputId = 'register-account-username'
const passwordInputId = 'register-account-password'
const confirmPasswordInputId = 'register-account-confirm-password'

// 用户名验证（3-20个字符，只能包含字母、数字、下划线）
const isValidUsername = computed(() => {
  if (!username.value) return false;
  const usernameRegex = /^[a-zA-Z0-9_]{3,20}$/;
  return usernameRegex.test(username.value);
});

// 密码强度验证（至少6个字符）
const isValidPassword = computed(() => {
  return password.value.length >= 6;
});

// 确认密码验证
const isPasswordMatch = computed(() => {
  return password.value === confirmPassword.value && confirmPassword.value.length > 0;
});

// 表单是否可提交
const isFormValid = computed(() => {
  return isValidUsername.value && isValidPassword.value && isPasswordMatch.value;
});

function handleButtonRegister() {
  if (!isFormValid.value) {
    return;
  }

  console.log('注册信息：', {
    username: username.value,
    password: password.value,
  });
}
</script>

<template>
  <div class="register-form">
    <!-- Username -->
    <div class="form-group">
      <label :for="usernameInputId">{{ $t('authentication.register.username') }}</label>
      <n-input
        :id="usernameInputId"
        v-model:value="username"
        :placeholder="$t('authentication.register.input_username')"
        clearable
        type="text"
        autocomplete="username"
        :status="username && !isValidUsername ? 'error' : undefined"
      />
      <span v-if="username && !isValidUsername" class="error-hint">
        {{ $t('authentication.register.invalid_username') }}
      </span>
    </div>

    <!-- Password -->
    <div class="form-group">
      <label :for="passwordInputId">{{ $t('authentication.register.password') }}</label>
      <n-input
        :id="passwordInputId"
        v-model:value="password"
        :placeholder="$t('authentication.register.input_password')"
        clearable
        type="password"
        show-password-on="click"
        autocomplete="new-password"
        :status="password && !isValidPassword ? 'error' : undefined"
      />
      <span v-if="password && !isValidPassword" class="error-hint">
        {{ $t('authentication.register.invalid_password') }}
      </span>
    </div>

    <!-- Confirm Password -->
    <div class="form-group">
      <label :for="confirmPasswordInputId">{{ $t('authentication.register.confirm_password') }}</label>
      <n-input
        :id="confirmPasswordInputId"
        v-model:value="confirmPassword"
        :placeholder="$t('authentication.register.input_confirm_password')"
        clearable
        type="password"
        show-password-on="click"
        autocomplete="new-password"
        :status="confirmPassword && !isPasswordMatch ? 'error' : undefined"
      />
      <span v-if="confirmPassword && !isPasswordMatch" class="error-hint">
        {{ $t('authentication.register.password_not_match') }}
      </span>
    </div>

    <!-- Register Button -->
    <n-button
      type="primary"
      block
      size="large"
      class="register-button"
      :disabled="!isFormValid"
      @click="handleButtonRegister"
    >
      {{ $t('authentication.register.register') }}
    </n-button>
  </div>
</template>

<style scoped lang="less">
.register-form {
  display: flex;
  flex-direction: column;
  gap: 1.5rem;

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
}

.form-group {
  display: flex;
  flex-direction: column;
  gap: 0.75rem;

  label {
    display: block;
    font-weight: 500;
    font-size: 0.95rem;
    color: var(--auth-text-primary);
  }

  .error-hint {
    font-size: 0.85rem;
    color: #d03050;
    margin-top: -0.5rem;
  }
}

.register-button {
  font-weight: 600;
  height: 44px;
  margin-top: 0.5rem;

  &:disabled {
    cursor: not-allowed;
    opacity: 0.5;
  }
}
</style>
