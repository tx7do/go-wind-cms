<script setup lang="ts">
import { cn } from '@/lib/utils'

const props = defineProps<{
  label: string
  description?: string
  class?: string
}>()
</script>

<template>
  <div :class="cn('flex items-center justify-between rounded-lg border border-border bg-card p-4', props.class)">
    <div v-if="description">
      <div class="text-sm font-medium text-foreground">{{ label }}</div>
      <div class="mt-1 text-xs text-muted-foreground">{{ description }}</div>
    </div>
    <div v-else class="text-sm font-medium text-foreground">{{ label }}</div>
    <slot />
  </div>
</template>
