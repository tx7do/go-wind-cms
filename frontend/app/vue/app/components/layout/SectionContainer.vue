<script setup lang="ts">
import { cn } from '@/lib/utils'

const props = defineProps<{
  class?: string
  maxWidth?: 'default' | 'wide' | 'narrow'
  noPadding?: boolean
}>()

const widthClass = computed(() =>
  props.maxWidth === 'wide' ? 'max-w-300'
    : props.maxWidth === 'narrow' ? 'max-w-3xl'
    : 'max-w-[1200px]'
)
</script>

<template>
  <div :class="cn('w-full mx-auto px-8 max-md:px-4', !noPadding && 'py-12', widthClass, props.class)">
    <slot />
  </div>
</template>
