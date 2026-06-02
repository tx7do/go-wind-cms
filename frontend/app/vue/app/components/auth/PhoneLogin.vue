<script setup lang="ts">
const { t } = useI18n()

const phone = ref('')
const verificationCode = ref('')
const codeSent = ref(false)
const countdown = ref(0)

let timer: ReturnType<typeof setTimeout> | null = null

watch(countdown, (val) => {
  if (val > 0) {
    timer = setTimeout(() => countdown.value--, 1000)
  } else if (codeSent.value) {
    codeSent.value = false
  }
})

onUnmounted(() => {
  if (timer) clearTimeout(timer)
})

function handleSendCode() {
  if (!phone.value) return
  codeSent.value = true
  countdown.value = 60
}

function handleLogin() {
  if (!phone.value || !verificationCode.value) return
  console.log('登录信息:', { phone: phone.value, verificationCode: verificationCode.value })
}

const inputBase = 'w-full rounded-lg border border-border bg-background px-4 py-2.5 text-sm text-foreground transition-colors hover:border-primary focus:border-primary focus:outline-none focus:ring-[3px] focus:ring-primary/15'
</script>

<template>
  <div class="space-y-4">
    <div class="space-y-2">
      <label for="login-phone" class="block text-sm font-medium text-foreground">
        {{ t('authentication.register.phone') }}
      </label>
      <input
        id="login-phone"
        v-model="phone"
        type="tel"
        :class="inputBase"
        :placeholder="t('authentication.login.placeholder_phone')"
        autocomplete="tel"
        autocapitalize="none"
        autocorrect="off"
        spellcheck="false"
      />
    </div>

    <div class="space-y-2">
      <label for="login-code" class="block text-sm font-medium text-foreground">
        {{ t('authentication.register.code') }}
      </label>
      <div class="flex gap-2">
        <input
          id="login-code"
          v-model="verificationCode"
          type="text"
          class="flex-1 rounded-lg border border-border bg-background px-4 py-2.5 text-sm text-foreground transition-colors hover:border-primary focus:border-primary focus:outline-none focus:ring-[3px] focus:ring-primary/15"
          :placeholder="t('authentication.login.placeholder_code')"
          maxlength="6"
          autocomplete="one-time-code"
          autocapitalize="none"
          autocorrect="off"
          spellcheck="false"
          @keydown.enter="handleLogin"
        />
        <button
          :class="[
            'shrink-0 cursor-pointer rounded-lg border border-border bg-card px-4 py-2.5 text-sm font-medium text-foreground transition-all hover:border-primary hover:bg-primary/5 active:scale-[0.98]',
            codeSent ? 'cursor-not-allowed opacity-60' : '',
          ]"
          :disabled="codeSent"
          @click="handleSendCode"
        >
          {{ codeSent ? `${countdown}s` : t('authentication.register.send_code') }}
        </button>
      </div>
    </div>

    <button
      class="w-full cursor-pointer rounded-lg bg-primary px-4 py-2.5 text-sm font-medium text-primary-foreground transition-all hover:bg-primary/90 active:scale-[0.98]"
      @click="handleLogin"
    >
      {{ t('authentication.login.login') }}
    </button>
  </div>
</template>
