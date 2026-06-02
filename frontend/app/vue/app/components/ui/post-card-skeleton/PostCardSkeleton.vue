<script setup lang="ts">
import { cn } from '@/lib/utils'

const props = defineProps<{
  class?: string
}>()
</script>

<template>
  <div :class="cn('overflow-hidden rounded-2xl border border-border bg-card', props.class)">
    <UiSkeleton class="h-52 w-full" />
    <div class="space-y-3 p-5">
      <UiSkeleton class="h-5 w-4/5" />
      <UiSkeleton class="h-4 w-full" />
      <UiSkeleton class="h-4 w-3/4" />
      <div class="flex items-center gap-3 pt-2">
        <UiSkeleton class="h-5 w-16 rounded-full" />
        <UiSkeleton class="h-5 w-16 rounded-full" />
        <UiSkeleton class="h-4 w-20" />
      </div>
    </div>
  </div>
</template>
