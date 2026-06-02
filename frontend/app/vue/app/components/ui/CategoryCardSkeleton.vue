<script setup lang="ts">
import { cn } from '@/lib/utils'

const props = defineProps<{
  class?: string
}>()
</script>

<template>
  <div :class="cn('flex min-h-50 flex-col rounded-2xl border border-border bg-card p-6', props.class)">
    <UiSkeleton class="h-35 w-full rounded-xl" />
    <div class="mt-4 space-y-2">
      <UiSkeleton class="h-5 w-3/5" />
      <UiSkeleton class="h-4 w-full" />
      <UiSkeleton class="h-4 w-4/5" />
    </div>
  </div>
</template>
