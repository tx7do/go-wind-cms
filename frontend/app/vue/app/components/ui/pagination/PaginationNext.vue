<script setup lang="ts">
import { type HTMLAttributes } from 'vue'
import { XIcon } from '@/plugins/xicon'
import { cn } from '@/lib/utils'
import UiPaginationLink from './PaginationLink.vue'

const props = defineProps<{
  class?: HTMLAttributes['class']
}>()
</script>

<template>
  <UiPaginationLink
    aria-label="Go to next page"
    size="default"
    :class="cn('gap-1 px-2.5', props.class)"
    v-bind="$attrs"
  >
    <span>Next</span>
    <XIcon icon="lucide:chevron-right" :size="16" />
  </UiPaginationLink>
</template>
