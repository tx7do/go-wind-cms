<script setup lang="ts">
const { t } = useI18n()

const email = ref('')
const visibleEnter = ref(false)

const isValidEmail = computed(() => {
  if (!email.value) return false
  return /^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(email.value)
})

function handleNext() {
  if (!isValidEmail.value) return
  visibleEnter.value = true
}

const inputBase = 'w-full rounded-lg border border-border bg-background px-4 py-2.5 text-sm text-foreground transition-colors hover:border-primary focus:border-primary focus:outline-none focus:ring-[3px] focus:ring-primary/15'
</script>

<template>
  <div>
    <div v-if="!visibleEnter" class="space-y-4">
      <!-- Email Input Group -->
      <div class="space-y-2">
        <label for="register-email-address" class="block text-sm font-medium text-foreground">
          {{ t('authentication.register.email') }}
        </label>
        <input
          id="register-email-address"
          v-model="email"
          type="text"
          :class="[inputBase, email && !isValidEmail ? 'border-destructive focus:border-destructive focus:ring-destructive/15' : '']"
          :placeholder="t('authentication.register.input_email')"
          autocomplete="email"
          autocapitalize="none"
          autocorrect="off"
          spellcheck="false"
        />
        <span v-if="email && !isValidEmail" class="text-xs text-destructive">
          {{ t('authentication.register.invalid_email') }}
        </span>
      </div>

      <!-- Next Button -->
      <button
        type="button"
        class="w-full cursor-pointer rounded-lg bg-primary px-4 py-2.5 text-sm font-medium text-primary-foreground transition-all hover:bg-primary/90 active:scale-[0.98] disabled:cursor-not-allowed disabled:opacity-60"
        :disabled="!isValidEmail"
        @click="handleNext"
      >
        {{ t('authentication.register.next_step') }}
      </button>
    </div>

    <!-- Enter Code Step -->
    <AuthEmailRegisterEnter v-else :email="email" />
  </div>
</template>
