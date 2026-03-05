<script setup lang="ts">
import {ref, computed} from "vue";
import {$t} from '@/locales';

import EmailRegisterEnterCodePage from "@/pages/register/components/EmailRegisterEnterCodePage.vue";

const email = ref<string>('');
const visibleEnter = ref<boolean>(false);
const emailInputId = 'register-email-address';

// 简单的邮箱格式验证
const isValidEmail = computed(() => {
  if (!email.value) return false;
  const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
  return emailRegex.test(email.value);
});

function handleButtonNext() {
  if (!isValidEmail.value) {
    return;
  }
  visibleEnter.value = true;
}
</script>

<template>
  <div v-show="!visibleEnter" class="register-form">
    <!-- Email Input Group -->
    <div class="form-group">
      <label :for="emailInputId">{{ $t('authentication.register.email') }}</label>
      <n-input
        :id="emailInputId"
        v-model:value="email"
        :placeholder="$t('authentication.register.input_email')"
        clearable
        type="text"
        autocomplete="email"
        :status="email && !isValidEmail ? 'error' : undefined"
      />
      <span v-if="email && !isValidEmail" class="error-hint">
        {{ $t('authentication.register.invalid_email') }}
      </span>
    </div>

    <!-- Next Button -->
    <n-button
      type="primary"
      block
      size="large"
      class="register-button"
      :disabled="!isValidEmail"
      @click="handleButtonNext"
    >
      {{ $t('authentication.register.next_step') }}
    </n-button>
  </div>

  <EmailRegisterEnterCodePage v-show="visibleEnter" :email="email" />
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
