<script setup lang="ts">
import { XIcon } from '@/plugins/xicon'
import { cn } from '@/lib/utils'

const props = defineProps<{
  icon?: string
  iconSize?: number
  class?: string
}>()
</script>

<template>
  <span :class="cn('flex items-center gap-1.5 text-muted-foreground', props.class)">
    <XIcon v-if="icon" :icon="icon" :size="iconSize ?? 16" class="text-muted-foreground" />
    <slot />
  </span>
</template>
