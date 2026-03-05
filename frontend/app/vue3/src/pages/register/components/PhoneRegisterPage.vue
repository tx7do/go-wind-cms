<script setup lang="ts">
import {onBeforeUnmount, ref} from "vue";
import {$t} from '@/locales';

const phone = ref('');
const verificationCode = ref('');
const codeSent = ref(false);
const codeCountdown = ref(0);
const phoneInputId = 'register-phone-number'
const codeInputId = 'register-phone-code'
let codeTimer: number | null = null

function clearCodeTimer() {
  if (codeTimer !== null) {
    clearInterval(codeTimer)
    codeTimer = null
  }
}

/**
 * 发送验证码
 */
function handleButtonSendVerifyCode() {
  if (!phone.value) {
    return;
  }

  codeSent.value = true;
  codeCountdown.value = 60;

  clearCodeTimer()
  codeTimer = window.setInterval(() => {
    codeCountdown.value--;
    if (codeCountdown.value <= 0) {
      clearCodeTimer()
      codeSent.value = false;
    }
  }, 1000);
}

/**
 * 登录或者注册
 */
function handleButtonRegisterOrLogin() {
  if (!phone.value || !verificationCode.value) {
    return;
  }
}

onBeforeUnmount(() => {
  clearCodeTimer()
})
</script>

<template>
  <div class="register-form">
    <!-- Phone Number Group -->
    <div class="form-group">
      <label :for="phoneInputId">{{ $t('authentication.register.phone') }}</label>
      <n-input
        :id="phoneInputId"
        v-model:value="phone"
        :placeholder="$t('authentication.register.input_phone')"
        clearable
        type="text"
        autocomplete="tel"
      />
    </div>

    <!-- Verification Code Group -->
    <div class="form-group">
      <label :for="codeInputId">{{ $t('authentication.register.code') }}</label>
      <div class="code-input-row">
        <n-input
          :id="codeInputId"
          v-model:value="verificationCode"
          :placeholder="$t('authentication.register.input_code')"
          maxlength="6"
          type="text"
          autocomplete="one-time-code"
        />
        <n-button
          :disabled="codeSent"
          :type="codeSent ? 'default' : 'primary'"
          class="send-code-btn"
          aria-live="polite"
          @click="handleButtonSendVerifyCode"
        >
          {{ codeSent ? `${codeCountdown}s` : $t('authentication.register.send_code') }}
        </n-button>
      </div>
    </div>

    <!-- Register Button -->
    <n-button
      type="primary"
      block
      size="large"
      class="register-button"
      @click="handleButtonRegisterOrLogin"
    >
      {{ $t('authentication.register.register_or_login') }}
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
}

.code-input-row {
  display: flex;
  gap: 0.75rem;
  align-items: center;

  :deep(.n-input) {
    flex: 1;
  }

  .send-code-btn {
    flex-shrink: 0;
    width: 120px;
    height: 40px;
    font-size: 0.9rem;
    white-space: nowrap;
  }
}

.register-button {
  font-weight: 600;
  height: 44px;
  margin-top: 0.5rem;
}
</style>
