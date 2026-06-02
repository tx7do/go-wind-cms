<script setup lang="ts">
import { type HTMLAttributes } from 'vue'
import { XIcon } from '@/plugins/xicon'
import { cn } from '@/lib/utils'

const props = defineProps<{
  class?: HTMLAttributes['class']
}>()
</script>

<template>
  <span aria-hidden :class="cn('flex h-9 w-9 items-center justify-center', props.class)">
    <XIcon icon="lucide:more-horizontal" :size="16" />
    <span class="sr-only">More pages</span>
  </span>
</template>
