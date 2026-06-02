<script setup lang="ts">
const props = defineProps<{
  email: string
}>()

const { t } = useI18n()

const code = ref('')
const isCodeComplete = ref(false)
const codeSent = ref(true)
const codeCountdown = ref(60)

let timer: ReturnType<typeof setInterval> | null = null

function startCountdown() {
  codeSent.value = true
  codeCountdown.value = 60

  timer = setInterval(() => {
    codeCountdown.value--
    if (codeCountdown.value <= 0) {
      if (timer) clearInterval(timer)
      codeSent.value = false
    }
  }, 1000)
}

onMounted(() => {
  startCountdown()
})

onUnmounted(() => {
  if (timer) clearInterval(timer)
})

function handleInputComplete() {
  if (code.value.length === 6) {
    isCodeComplete.value = true
  }
}

function handleResend() {
  if (codeSent.value) return
  startCountdown()
}

function handleConfirm() {
  if (!isCodeComplete.value) return
  console.log('确认注册')
}
</script>

<template>
  <div class="space-y-4">
    <!-- Email Hint -->
    <div class="rounded-lg bg-primary/5 p-4 text-center">
      <p class="text-sm font-medium text-foreground">{{ t('authentication.register.code_sent_title') }}</p>
      <p class="mt-1 text-base font-semibold text-primary">{{ props.email }}</p>
      <p class="mt-1 text-xs text-muted-foreground">{{ t('authentication.register.code_sent_subtitle') }}</p>
    </div>

    <!-- Verification Code Input -->
    <div>
      <input
        v-model="code"
        type="text"
        maxlength="6"
        class="w-full rounded-lg border border-border bg-background px-4 py-3 text-center text-lg tracking-widest text-foreground transition-colors hover:border-primary focus:border-primary focus:outline-none focus:ring-[3px] focus:ring-primary/15"
        autocapitalize="none"
        autocorrect="off"
        spellcheck="false"
        :placeholder="t('authentication.register.input_code')"
        @input="handleInputComplete"
      />
    </div>

    <!-- Resend Button -->
    <div class="text-center">
      <button
        class="cursor-pointer border-none bg-transparent text-sm text-primary transition-colors hover:text-primary/80 hover:underline disabled:cursor-not-allowed disabled:text-muted-foreground disabled:no-underline"
        :disabled="codeSent"
        @click="handleResend"
      >
        {{ codeSent ? `${codeCountdown}s ${t('authentication.register.resend_after')}` : t('authentication.register.resend') }}
      </button>
    </div>

    <!-- Confirm Button -->
    <button
      type="button"
      class="w-full cursor-pointer rounded-lg bg-primary px-4 py-2.5 text-sm font-medium text-primary-foreground transition-all hover:bg-primary/90 active:scale-[0.98] disabled:cursor-not-allowed disabled:opacity-60"
      :disabled="!isCodeComplete"
      @click="handleConfirm"
    >
      {{ t('authentication.register.confirm') }}
    </button>
  </div>
</template>
