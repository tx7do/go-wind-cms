<script setup lang="ts">
import { cn } from '@/lib/utils'

const props = withDefaults(defineProps<{
  count?: number
  class?: string
}>(), {
  count: 6,
})
</script>

<template>
  <div :class="cn('grid grid-cols-[repeat(auto-fill,minmax(260px,1fr))] gap-6', props.class)">
    <UiCategoryCardSkeleton v-for="(_, i) in count" :key="i" />
  </div>
</template>
