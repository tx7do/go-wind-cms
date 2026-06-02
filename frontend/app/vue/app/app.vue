<script setup lang="ts">
import { Toaster } from '@/components/ui/sonner'

const isDev = import.meta.dev
const { t } = useI18n()
const config = useRuntimeConfig()

useHead({
  titleTemplate: (title) => title ? `${title} - ${config.public.appTitle}` : config.public.appTitle,
})
</script>

<template>
  <div>
    <NuxtLayout>
      <NuxtPage />
    </NuxtLayout>
    <Toaster />
    <ClientOnly>
      <VueQueryDevtools v-if="isDev" />
    </ClientOnly>
  </div>
</template>
