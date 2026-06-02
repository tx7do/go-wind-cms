<script setup lang="ts">
import { XIcon } from '@/plugins/xicon'

const { t } = useI18n()
const localePath = useLocalePath()
</script>

<template>
  <section class="w-full max-w-300 mx-auto scroll-reveal px-8 py-12 max-md:px-4">
    <div class="mb-8 flex items-center justify-between">
      <h2 class="flex items-center gap-2 text-2xl font-extrabold tracking-tight text-foreground max-md:text-xl">
        <XIcon name="carbon:document" :size="28" class="mr-2 text-primary" />
        {{ t('page.home.latest_posts') }}
      </h2>
      <UiButton variant="ghost" @click="navigateTo(localePath('/post'))">
        {{ t('page.home.view_all') }} →
      </UiButton>
    </div>
    <div class="w-full">
      <!-- PostList placeholder -->
      <div class="grid grid-cols-1 gap-6 md:grid-cols-2 lg:grid-cols-3">
        <div v-for="i in 6" :key="i" class="flex min-h-50 flex-col rounded-2xl border border-border bg-card p-6">
          <UiSkeleton class="h-35 w-full" />
        </div>
      </div>
    </div>
  </section>
</template>
