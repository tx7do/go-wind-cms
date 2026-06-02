<script setup lang="ts">
import { cn } from '@/lib/utils'

const props = defineProps<{
  size?: 'sm' | 'md' | 'lg'
  class?: string
}>()

const minHeight = props.size === 'sm' ? 'min-h-[220px]' : props.size === 'lg' ? 'min-h-[340px]' : 'min-h-[280px]'
const padding = props.size === 'sm' ? 'py-14' : props.size === 'lg' ? 'py-24' : 'py-18'
</script>

<template>
  <section :class="cn('relative w-full overflow-hidden border-b border-border', minHeight, padding, 'flex items-center justify-center', props.class)">
    <div class="flex flex-col items-center gap-4">
      <UiSkeleton v-if="size !== 'sm'" class="h-12 w-12 rounded-full" />
      <UiSkeleton class="h-10 w-64 rounded-lg" />
      <UiSkeleton class="h-5 w-80 rounded-lg" />
      <div class="mt-2 flex gap-3">
        <UiSkeleton class="h-5 w-20 rounded-full" />
        <UiSkeleton class="h-5 w-20 rounded-full" />
      </div>
    </div>
  </section>
</template>
