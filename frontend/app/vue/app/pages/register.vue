<script setup lang="ts">
definePageMeta({ layout: 'auth' })
const localePath = useLocalePath()
</script>

<template>
  <div class="flex min-h-screen w-full items-center justify-center p-4">
    <div class="w-full max-w-md rounded-2xl border border-border bg-card p-8 shadow-lg">
      <div class="mb-8 text-center">
        <NuxtLink :to="localePath('/')" class="mb-4 inline-flex items-center gap-2">
          <img src="/logo.png" alt="Logo" width="40" height="40" class="h-10 w-10" />
          <span class="text-xl font-bold text-primary">{{ $t('app.title') }}</span>
        </NuxtLink>
        <h1 class="text-2xl font-bold text-foreground">{{ $t('authentication.register.title') }}</h1>
        <p class="mt-2 text-sm text-muted-foreground">{{ $t('authentication.register.register_with') }}</p>
      </div>

      <form @submit.prevent class="space-y-4">
        <div>
          <UiLabel class="mb-2 block text-sm font-medium">{{ $t('authentication.register.username') }}</UiLabel>
          <UiInput type="text" :placeholder="$t('authentication.register.input_username')" />
        </div>
        <div>
          <UiLabel class="mb-2 block text-sm font-medium">{{ $t('authentication.register.email') }}</UiLabel>
          <UiInput type="email" :placeholder="$t('authentication.register.input_email')" />
        </div>
        <div>
          <UiLabel class="mb-2 block text-sm font-medium">{{ $t('authentication.register.password') }}</UiLabel>
          <UiInput type="password" :placeholder="$t('authentication.register.input_password')" />
        </div>
        <div>
          <UiLabel class="mb-2 block text-sm font-medium">{{ $t('authentication.register.confirm_password') }}</UiLabel>
          <UiInput type="password" :placeholder="$t('authentication.register.input_confirm_password')" />
        </div>
        <UiButton class="w-full" type="submit">{{ $t('authentication.register.register') }}</UiButton>
      </form>

      <div class="mt-6 text-center text-sm text-muted-foreground">
        {{ $t('authentication.register.already_have_account') }}
        <NuxtLink :to="localePath('/login')" class="text-primary hover:underline">{{ $t('authentication.register.login_now') }}</NuxtLink>
      </div>
    </div>
  </div>
</template>
