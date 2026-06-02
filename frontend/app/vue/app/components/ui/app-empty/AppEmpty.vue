<script setup lang="ts">
import { cn } from '@/lib/utils'

const props = defineProps<{
  variant?: 'default' | 'error' | 'noData'
  inContainer?: boolean
  description?: string
  class?: string
}>()
</script>

<template>
  <div :class="cn('flex w-full items-center justify-center gap-4 py-12 px-5', props.inContainer && 'my-20', props.class)">
    <div class="opacity-50">
      <slot name="image">
        <XIcon icon="lucide:inbox" width="64" height="64" class="text-muted-foreground" />
      </slot>
    </div>
    <p v-if="description" class="text-sm text-muted-foreground">{{ description }}</p>
    <p v-else-if="$slots.default" class="text-sm text-muted-foreground">
      <slot />
    </p>
  </div>
</template>
