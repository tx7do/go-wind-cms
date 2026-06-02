<script setup lang="ts">
import { type HTMLAttributes } from 'vue'
import { cn } from '@/lib/utils'
import { buttonVariants } from '@/components/ui/button'

const props = withDefaults(defineProps<{
  isActive?: boolean
  size?: 'default' | 'sm' | 'lg' | 'icon'
  class?: HTMLAttributes['class']
}>(), {
  size: 'icon',
})
</script>

<template>
  <button
    :aria-current="isActive ? 'page' : undefined"
    :class="cn(
      buttonVariants({ variant: isActive ? 'outline' : 'ghost', size: props.size }),
      'h-9 w-9',
      props.class,
    )"
    v-bind="$attrs"
  >
    <slot />
  </button>
</template>
