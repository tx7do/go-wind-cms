<script setup lang="ts">
import { cn } from '@/lib/utils'

const props = withDefaults(defineProps<{
  columns?: number
  count?: number
  class?: string
}>(), {
  columns: 3,
  count: 6,
})
</script>

<template>
  <div
    :class="cn('grid gap-6 max-md:!grid-cols-1 max-md:gap-4', props.class)"
    :style="{ gridTemplateColumns: `repeat(${columns}, 1fr)` }"
  >
    <UiPostCardSkeleton v-for="(_, i) in count" :key="i" />
  </div>
</template>
