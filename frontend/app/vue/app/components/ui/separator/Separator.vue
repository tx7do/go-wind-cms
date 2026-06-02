<script setup lang="ts">
import { type HTMLAttributes, computed } from 'vue'
import { Separator as SeparatorPrimitive } from 'reka-ui'
import { cn } from '@/lib/utils'

const props = defineProps<{
  class?: HTMLAttributes['class']
  orientation?: 'horizontal' | 'vertical'
  decorative?: boolean
}>()

const delegatedProps = computed(() => {
  const { class: _, ...delegated } = props
  return delegated
})
</script>

<template>
  <SeparatorPrimitive
    v-bind="delegatedProps"
    :class="cn(
      'shrink-0 bg-border',
      orientation === 'vertical' ? 'h-full w-px' : 'h-px w-full',
      props.class,
    )"
  />
</template>
