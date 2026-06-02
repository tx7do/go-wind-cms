<script setup lang="ts">
import { cn } from '@/lib/utils'

const props = defineProps<{
  label: string
  onClick: () => void
  class?: string
}>()
</script>

<template>
  <button
    :class="cn(
      'group flex cursor-pointer items-center gap-2 rounded-lg px-3 py-1.5',
      'text-sm text-muted-foreground transition-colors',
      'hover:bg-muted/60 hover:text-foreground',
      props.class,
    )"
    @click="onClick"
  >
    <XIcon icon="lucide:arrow-left" width="16" height="16" class="transition-transform duration-300 group-hover:-translate-x-1" />
    <span>{{ label }}</span>
  </button>
</template>
