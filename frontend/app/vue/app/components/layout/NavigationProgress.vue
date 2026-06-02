<script setup lang="ts">
import { ref, watch, onUnmounted } from 'vue'
import { cn } from '@/lib/utils'

const route = useRoute()

const progress = ref(0)
const visible = ref(false)

let progressVal = 0
let intervalId: ReturnType<typeof setInterval> | null = null
let timeoutId: ReturnType<typeof setTimeout> | null = null
let prevPath: string | null = null

const clearTimers = () => {
  if (intervalId) { clearInterval(intervalId); intervalId = null }
  if (timeoutId) { clearTimeout(timeoutId); timeoutId = null }
}

const startProgress = () => {
  clearTimers()
  visible.value = true
  progressVal = 0
  progress.value = 0

  intervalId = setInterval(() => {
    if (progressVal < 30) {
      progressVal += Math.random() * 8 + 2
    } else if (progressVal < 60) {
      progressVal += Math.random() * 4 + 1
    } else if (progressVal < 85) {
      progressVal += Math.random() * 2
    } else {
      if (intervalId) { clearInterval(intervalId); intervalId = null }
    }
    progress.value = progressVal
  }, 200)
}

const finishProgress = () => {
  clearTimers()
  progressVal = 100
  progress.value = 100
  timeoutId = setTimeout(() => {
    visible.value = false
    setTimeout(() => { progressVal = 0; progress.value = 0 }, 100)
  }, 250)
}

// Watch route changes
watch(() => route.path, (newPath) => {
  if (prevPath === null) {
    prevPath = newPath
    return
  }
  if (prevPath !== newPath) {
    prevPath = newPath
    startProgress()
    timeoutId = setTimeout(() => finishProgress(), 600)
  }
})

onUnmounted(() => clearTimers())
</script>

<template>
  <div
    v-if="visible || progress > 0"
    class="fixed left-0 right-0 top-0 z-9999 pointer-events-none"
    role="progressbar"
    :aria-valuenow="Math.round(progress)"
    aria-valuemin="0"
    aria-valuemax="100"
  >
    <div class="h-[3px] w-full bg-transparent">
      <div
        :class="cn(
          'h-full bg-primary',
          'shadow-[0_0_10px_hsl(var(--primary)/0.6),0_0_4px_hsl(var(--primary)/0.8)]',
          'transition-[width,opacity] duration-300 ease-out',
        )"
        :style="{
          width: `${progress}%`,
          opacity: visible ? 1 : 0,
        }"
      />
    </div>
  </div>
</template>
