<script setup lang="ts">
import { ref, onMounted, onUnmounted, watch } from 'vue'
import { cn } from '@/lib/utils'

const props = defineProps<{
  scrollThreshold?: number
  onClick?: () => void
  class?: string
  visible?: boolean
}>()

const showBackToTop = ref(false)

const handleScroll = () => {
  showBackToTop.value = window.scrollY > (props.scrollThreshold ?? 500)
}

onMounted(() => {
  window.addEventListener('scroll', handleScroll)
  handleScroll()
})

onUnmounted(() => {
  window.removeEventListener('scroll', handleScroll)
})

const isVisible = computed(() => props.visible !== undefined ? props.visible : showBackToTop.value)

const scrollToTop = () => {
  window.scrollTo({ top: 0, behavior: 'smooth' })
  props.onClick?.()
}
</script>

<template>
  <button
    v-if="isVisible"
    :class="cn(
      'fixed bottom-10 right-10 z-999 flex h-14 w-14 items-center justify-center',
      'rounded-full border-none bg-primary text-white cursor-pointer',
      'text-2xl shadow-lg transition-all duration-300',
      'hover:-translate-y-1 hover:shadow-xl hover:bg-primary/80',
      'active:-translate-y-0.5',
      props.class,
    )"
    aria-label="Back to top"
    @click="scrollToTop"
  >
    <XIcon icon="carbon:arrow-up" />
  </button>
</template>
