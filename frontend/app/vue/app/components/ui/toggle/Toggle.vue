<script setup lang="ts">
import { type HTMLAttributes } from 'vue'
import { ToggleRoot } from 'reka-ui'
import { cn } from '@/lib/utils'
import { toggleVariants, type ToggleVariants } from '.'

const props = defineProps<{
  variant?: ToggleVariants['variant']
  size?: ToggleVariants['size']
  class?: HTMLAttributes['class']
  modelValue?: boolean
  defaultValue?: boolean
}>()

const emits = defineEmits<{
  'update:modelValue': [value: boolean]
}>()
</script>

<template>
  <ToggleRoot
    :class="cn(toggleVariants({ variant, size }), props.class)"
    :pressed="modelValue"
    :default-pressed="defaultValue"
    @update:pressed="emits('update:modelValue', $event)"
  >
    <slot />
  </ToggleRoot>
</template>
