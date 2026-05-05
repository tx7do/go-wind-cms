<script setup lang="ts">
import {ref, computed} from "vue";
import {$t} from '@/locales';

import VerificationCodeInput from "@/components/VerificationCodeInput";

const props = defineProps<{
  email: string
}>();

const code = ref(Array(6).fill(''));
const isCodeComplete = ref(false);
const codeSent = ref(true); // 初始已发送
const codeCountdown = ref(60);

// 启动倒计时
const startCountdown = () => {
  codeSent.value = true;
  codeCountdown.value = 60;

  const timer = setInterval(() => {
    codeCountdown.value--;
    if (codeCountdown.value <= 0) {
      clearInterval(timer);
      codeSent.value = false;
    }
  }, 1000);
};

// 初始化倒计时
startCountdown();

// 确认按钮是否可用
const isConfirmEnabled = computed(() => {
  return isCodeComplete.value;
});

function handleInputComplete(value: string) {
  console.log('验证码输入完成：', value);
  isCodeComplete.value = true;
}

function handleButtonResend() {
  if (codeSent.value) return;

  // TODO: 实现重新发送验证码逻辑
  startCountdown();
}

function handleButtonConfirm() {
  if (!isConfirmEnabled.value) return;

  // TODO: 实现确认注册逻辑
}
</script>

<template>
  <div class="code-container">
    <!-- Email Hint -->
    <div class="email-sent-info">
      <p class="hint-title">{{ $t('authentication.register.code_sent_title') }}</p>
      <p class="email-display">{{ props.email }}</p>
      <p class="hint-subtitle">{{ $t('authentication.register.code_sent_subtitle') }}</p>
    </div>

    <!-- Verification Code Input -->
    <div class="code-input-wrapper">
      <VerificationCodeInput
        :update:code="code"
        :length="6"
        @complete="handleInputComplete"
      />
    </div>

    <!-- Resend Button -->
    <div class="resend-section">
      <n-button
        text
        type="primary"
        :disabled="codeSent"
        @click="handleButtonResend"
      >
        {{ codeSent ? `${codeCountdown}s ${$t('authentication.register.resend_after')}` : $t('authentication.register.resend') }}
      </n-button>
    </div>

    <!-- Confirm Button -->
    <n-button
      type="primary"
      block
      size="large"
      class="confirm-button"
      :disabled="!isConfirmEnabled"
      @click="handleButtonConfirm"
    >
      {{ $t('authentication.register.confirm') }}
    </n-button>
  </div>
</template>

<style scoped lang="less">
.code-container {
  display: flex;
  flex-direction: column;
  gap: 1.5rem;
  padding: 0.5rem 0;
}

.email-sent-info {
  text-align: center;
  margin-bottom: 0.5rem;

  .hint-title {
    font-size: 0.95rem;
    color: var(--color-text-primary);
    margin: 0 0 0.5rem 0;
    font-weight: 500;
  }

  .email-display {
    font-size: 1rem;
    color: var(--color-brand);
    margin: 0 0 0.5rem 0;
    font-weight: 600;
  }

  .hint-subtitle {
    font-size: 0.85rem;
    color: var(--color-text-secondary);
    margin: 0;
    line-height: 1.5;
  }
}

.code-input-wrapper {
  display: flex;
  justify-content: center;
  margin: 0.5rem 0;
}

.resend-section {
  text-align: center;
  margin: -0.5rem 0 0.5rem 0;
}

.confirm-button {
  font-weight: 600;
  height: 44px;
  margin-top: 0.5rem;
  background: linear-gradient(135deg, var(--color-brand) 0%, #a855f7 100%);
  border: none;
  box-shadow: 0 4px 16px rgba(168, 85, 247, 0.3);
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);

  &:hover:not(:disabled) {
    transform: translateY(-2px);
    box-shadow: 0 6px 20px rgba(168, 85, 247, 0.4);
  }

  &:active:not(:disabled) {
    transform: translateY(0);
  }

  &:disabled {
    cursor: not-allowed;
    opacity: 0.5;
  }
}
</style>
