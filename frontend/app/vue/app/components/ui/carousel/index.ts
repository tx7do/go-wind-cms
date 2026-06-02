import emblaCarouselVue from 'embla-carousel-vue'
import { type HTMLAttributes, computed, provide, inject, ref, type Ref, watch, onMounted, onUnmounted } from 'vue'
import { cn } from '@/lib/utils'

type CarouselApi = ReturnType<typeof emblaCarouselVue>[1]
type CarouselOptions = Parameters<typeof emblaCarouselVue>[0]
type CarouselPlugin = Parameters<typeof emblaCarouselVue>[1]

interface CarouselProps {
  opts?: CarouselOptions
  plugins?: CarouselPlugin
  orientation?: 'horizontal' | 'vertical'
  class?: HTMLAttributes['class']
}

interface CarouselContext {
  carouselRef: Ref<HTMLElement | null>
  api: CarouselApi
  orientation: Ref<'horizontal' | 'vertical'>
  scrollPrev: () => void
  scrollNext: () => void
  canScrollPrev: Ref<boolean>
  canScrollNext: Ref<boolean>
}

const CAROUSEL_KEY = Symbol('carousel')

function useCarousel() {
  const ctx = inject<CarouselContext>(CAROUSEL_KEY)
  if (!ctx) throw new Error('useCarousel must be used within a <UiCarousel />')
  return ctx
}

export { useCarousel, type CarouselApi }
export { default as Carousel } from './Carousel.vue'
export { default as CarouselContent } from './CarouselContent.vue'
export { default as CarouselItem } from './CarouselItem.vue'
export { default as CarouselPrevious } from './CarouselPrevious.vue'
export { default as CarouselNext } from './CarouselNext.vue'

export { CAROUSEL_KEY, type CarouselContext, type CarouselProps }
