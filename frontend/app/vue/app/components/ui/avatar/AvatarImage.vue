<script setup lang="ts">
import { type HTMLAttributes } from 'vue'
import { cn } from '@/lib/utils'

const props = defineProps<{
  src?: string
  alt?: string
  class?: HTMLAttributes['class']
}>()
</script>

<template>
  <img :src="src" :alt="alt" :class="cn('aspect-square h-full w-full', props.class)" />
</template>
