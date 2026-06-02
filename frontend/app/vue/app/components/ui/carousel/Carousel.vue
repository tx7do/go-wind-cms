<script setup lang="ts">
import { type HTMLAttributes, provide, ref, watch, onMounted, onUnmounted } from 'vue'
import emblaCarouselVue from 'embla-carousel-vue'
import { cn } from '@/lib/utils'
import { CAROUSEL_KEY, type CarouselContext, type CarouselApi } from '.'

const props = withDefaults(defineProps<{
  opts?: Parameters<typeof emblaCarouselVue>[0]
  plugins?: Parameters<typeof emblaCarouselVue>[1]
  orientation?: 'horizontal' | 'vertical'
  class?: HTMLAttributes['class']
}>(), {
  orientation: 'horizontal',
})

const [emblaRef, emblaApi] = emblaCarouselVue(
  computed(() => ({
    ...props.opts,
    axis: props.orientation === 'horizontal' ? 'x' : 'y',
  })),
  computed(() => props.plugins),
)

const canScrollPrev = ref(false)
const canScrollNext = ref(false)

const scrollPrev = () => emblaApi.value?.scrollPrev()
const scrollNext = () => emblaApi.value?.scrollNext()

function onSelect() {
  if (!emblaApi.value) return
  canScrollPrev.value = emblaApi.value.canScrollPrev()
  canScrollNext.value = emblaApi.value.canScrollNext()
}

let inited = false
watch(emblaApi, (api) => {
  if (!api) return
  onSelect()
  api.on('select', onSelect)
  api.on('reInit', onSelect)
  inited = true
})

const handleKeyDown = (event: KeyboardEvent) => {
  if (event.key === 'ArrowLeft') {
    event.preventDefault()
    scrollPrev()
  } else if (event.key === 'ArrowRight') {
    event.preventDefault()
    scrollNext()
  }
}

provide(CAROUSEL_KEY, {
  carouselRef: emblaRef,
  api: emblaApi,
  orientation: computed(() => props.orientation),
  scrollPrev,
  scrollNext,
  canScrollPrev,
  canScrollNext,
} as CarouselContext)
</script>

<template>
  <div
    ref="emblaRef"
    :class="cn('relative', props.class)"
    role="region"
    aria-roledescription="carousel"
    @keydown.capture="handleKeyDown"
  >
    <slot />
  </div>
</template>
