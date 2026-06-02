<script setup lang="ts">
import { Toaster } from '@/components/ui/sonner'

const isDev = import.meta.dev
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
