<script setup lang="ts">
const { t } = useI18n()

const phone = ref('')
const verificationCode = ref('')
const codeSent = ref(false)
const codeCountdown = ref(0)

let codeTimer: ReturnType<typeof setInterval> | null = null

watch([codeSent, codeCountdown], () => {
  if (codeSent.value && codeCountdown.value > 0) {
    codeTimer = setInterval(() => {
      codeCountdown.value--
      if (codeCountdown.value <= 0) {
        codeSent.value = false
        if (codeTimer) clearInterval(codeTimer)
      }
    }, 1000)
  }
})

onUnmounted(() => {
  if (codeTimer) clearInterval(codeTimer)
})

function handleSendCode() {
  if (!phone.value) return
  codeSent.value = true
  codeCountdown.value = 60
}

function handleRegister() {
  if (!phone.value || !verificationCode.value) return
  console.log('手机号注册/登录:', { phone: phone.value, verificationCode: verificationCode.value })
}

const inputBase = 'w-full rounded-lg border border-border bg-background px-4 py-2.5 text-sm text-foreground transition-colors hover:border-primary focus:border-primary focus:outline-none focus:ring-[3px] focus:ring-primary/15'
</script>

<template>
  <div class="space-y-4">
    <!-- Phone Number -->
    <div class="space-y-2">
      <label for="register-phone-number" class="block text-sm font-medium text-foreground">
        {{ t('authentication.register.phone') }}
      </label>
      <input
        id="register-phone-number"
        v-model="phone"
        type="text"
        :class="inputBase"
        :placeholder="t('authentication.register.input_phone')"
        autocomplete="tel"
        autocapitalize="none"
        autocorrect="off"
        spellcheck="false"
      />
    </div>

    <!-- Verification Code -->
    <div class="space-y-2">
      <label for="register-phone-code" class="block text-sm font-medium text-foreground">
        {{ t('authentication.register.code') }}
      </label>
      <div class="flex gap-2">
        <input
          id="register-phone-code"
          v-model="verificationCode"
          type="text"
          class="flex-1 rounded-lg border border-border bg-background px-4 py-2.5 text-sm text-foreground transition-colors hover:border-primary focus:border-primary focus:outline-none focus:ring-[3px] focus:ring-primary/15"
          :placeholder="t('authentication.register.input_code')"
          maxlength="6"
          autocomplete="one-time-code"
          autocapitalize="none"
          autocorrect="off"
          spellcheck="false"
        />
        <button
          type="button"
          :disabled="codeSent"
          :class="[
            'shrink-0 cursor-pointer rounded-lg border border-border bg-card px-4 py-2.5 text-sm font-medium text-foreground transition-all hover:border-primary hover:bg-primary/5 active:scale-[0.98]',
            codeSent ? 'cursor-not-allowed opacity-60' : '',
          ]"
          @click="handleSendCode"
        >
          {{ codeSent ? `${codeCountdown}s` : t('authentication.register.send_code') }}
        </button>
      </div>
    </div>

    <!-- Register Button -->
    <button
      type="button"
      class="w-full cursor-pointer rounded-lg bg-primary px-4 py-2.5 text-sm font-medium text-primary-foreground transition-all hover:bg-primary/90 active:scale-[0.98]"
      @click="handleRegister"
    >
      {{ t('authentication.register.register_or_login') }}
    </button>
  </div>
</template>
