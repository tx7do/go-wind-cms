<script setup lang="ts">
import { type HTMLAttributes } from 'vue'
import { XIcon } from '@/plugins/xicon'
import { cn } from '@/lib/utils'
import { UiButton } from '@/components/ui/button'
import { useCarousel } from '.'

const props = withDefaults(defineProps<{
  variant?: 'default' | 'destructive' | 'outline' | 'secondary' | 'ghost' | 'link'
  size?: 'default' | 'sm' | 'lg' | 'icon'
  class?: HTMLAttributes['class']
}>(), {
  variant: 'outline',
  size: 'icon',
})

const { orientation, scrollNext, canScrollNext } = useCarousel()
</script>

<template>
  <UiButton
    :variant="variant"
    :size="size"
    :class="cn(
      'absolute h-8 w-8 rounded-full',
      orientation === 'horizontal'
        ? '-right-12 top-1/2 -translate-y-1/2'
        : '-bottom-12 left-1/2 -translate-x-1/2 rotate-90',
      props.class,
    )"
    :disabled="!canScrollNext"
    @click="scrollNext"
  >
    <XIcon icon="lucide:arrow-right" :size="16" />
    <span class="sr-only">Next slide</span>
  </UiButton>
</template>
