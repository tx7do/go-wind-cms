<script setup lang="ts">
import { XIcon } from '@/plugins/xicon'

const { t } = useI18n()

const handleGoogleRegister = () => {
  console.log('Google 注册')
}

const handleGithubRegister = () => {
  console.log('GitHub 注册')
}
</script>

<template>
  <div class="space-y-3">
    <button
      type="button"
      class="flex w-full cursor-pointer items-center justify-center gap-3 rounded-lg border border-border bg-card px-4 py-3 text-sm font-medium text-foreground transition-all hover:border-primary hover:bg-primary/5 active:scale-[0.98]"
      @click="handleGoogleRegister"
    >
      <XIcon icon="logos:google-icon" :size="18" />
      {{ t('authentication.login.social_google') }}
    </button>

    <button
      type="button"
      class="flex w-full cursor-pointer items-center justify-center gap-3 rounded-lg border border-border bg-card px-4 py-3 text-sm font-medium text-foreground transition-all hover:border-primary hover:bg-primary/5 active:scale-[0.98]"
      @click="handleGithubRegister"
    >
      <XIcon icon="logos:github-icon" :size="18" />
      {{ t('authentication.login.social_github') }}
    </button>
  </div>
</template>
