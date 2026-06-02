<script setup lang="ts">
import { XIcon } from '@/plugins/xicon'

export interface TocItem {
  id: string
  level: number
  text: string
  index: number
}

const props = defineProps<{
  contentRef: HTMLElement | null
  contentKey: string
  isExpanded: boolean
  title: string
}>()

const emit = defineEmits<{
  collapse: []
  expand: []
}>()

const HEADING_OFFSET = 150
const THROTTLE_DELAY = 200

const tableOfContents = ref<TocItem[]>([])
const activeHeading = ref('')

let throttleTimer: ReturnType<typeof setTimeout> | null = null

function getLiveHeading(index: number): Element | null {
  if (!props.contentRef) return null
  const headings = props.contentRef.querySelectorAll('h2, h3')
  return headings[index] || null
}

function generateToc() {
  if (!props.contentRef) return
  const headings = props.contentRef.querySelectorAll('h2, h3')
  const toc: TocItem[] = []

  headings.forEach((heading, index) => {
    const id = `heading-${index}`
    heading.setAttribute('id', id)
    toc.push({
      id,
      level: heading.tagName === 'H2' ? 2 : 3,
      text: heading.textContent || '',
      index,
    })
  })

  tableOfContents.value = toc
}

function scrollToHeading(id: string) {
  const tocItem = tableOfContents.value.find(item => item.id === id)
  if (!tocItem) return

  const element = getLiveHeading(tocItem.index)
  if (!element) return

  element.setAttribute('id', id)

  const rect = element.getBoundingClientRect()
  const HEADER_OFFSET = 80
  const targetPosition = window.scrollY + rect.top - HEADER_OFFSET

  window.scrollTo({
    top: Math.max(0, targetPosition),
    behavior: 'smooth',
  })

  activeHeading.value = id

  if (window.history.pushState) {
    const currentState = window.history.state || {}
    window.history.pushState(currentState, '', `#${id}`)
  }
}

function handleScroll() {
  if (tableOfContents.value.length === 0) return
  let currentActive = ''
  for (const tocItem of tableOfContents.value) {
    const el = getLiveHeading(tocItem.index)
    if (el && el.getBoundingClientRect().top < HEADING_OFFSET) {
      currentActive = tocItem.id
    }
  }
  if (currentActive !== activeHeading.value) {
    activeHeading.value = currentActive
  }
}

function throttledScroll() {
  if (!throttleTimer) {
    throttleTimer = setTimeout(() => {
      handleScroll()
      throttleTimer = null
    }, THROTTLE_DELAY)
  }
}

// Generate TOC when content changes
watch(() => props.contentKey, (key) => {
  if (!key) return
  setTimeout(generateToc, 100)
  setTimeout(generateToc, 1500)
}, { immediate: true })

// Check URL hash on load
watch([() => props.contentKey, () => tableOfContents.value.length], () => {
  if (!props.contentKey || tableOfContents.value.length === 0 || !import.meta.client) return
  if (window.location.hash) {
    const hashId = window.location.hash.slice(1)
    setTimeout(() => scrollToHeading(hashId), 500)
  }
})

// Scroll listener
onMounted(() => {
  window.addEventListener('scroll', throttledScroll)
})

onUnmounted(() => {
  window.removeEventListener('scroll', throttledScroll)
  if (throttleTimer) clearTimeout(throttleTimer)
})
</script>

<template>
  <!-- Sidebar TOC -->
  <aside v-if="isExpanded && tableOfContents.length > 0" class="w-60 shrink-0 max-md:hidden">
    <nav class="sticky top-24 rounded-lg border border-border bg-card p-4 pr-5 border-r border-r-border/50">
      <div class="mb-3 flex items-center justify-between">
        <h3 class="flex items-center gap-2 text-sm font-semibold text-foreground">
          <XIcon icon="carbon:list" />
          <span>{{ title }}</span>
        </h3>
        <button
          class="cursor-pointer rounded p-1 text-muted-foreground transition-colors hover:bg-primary/5 hover:text-primary"
          @click="emit('collapse')"
        >
          <XIcon icon="carbon:chevron-left" />
        </button>
      </div>
      <div class="space-y-0.5">
        <a
          v-for="item in tableOfContents"
          :key="item.id"
          :href="`#${item.id}`"
          :class="[
            'block truncate rounded px-2 py-1.5 text-sm transition-colors',
            activeHeading === item.id
              ? 'bg-primary/10 font-medium text-primary'
              : 'text-muted-foreground hover:bg-muted/50 hover:text-foreground',
            item.level === 3 ? 'pl-6 text-[13px]' : '',
          ]"
          @click.prevent="scrollToHeading(item.id)"
        >
          {{ item.text }}
        </a>
      </div>
    </nav>
  </aside>

  <!-- Collapsed expand button -->
  <div v-if="!isExpanded && tableOfContents.length > 0" class="fixed top-1/2 left-4 z-10 -translate-y-1/2">
    <button
      class="flex cursor-pointer items-center gap-2 rounded-lg border border-border bg-card px-3 py-2 text-sm text-foreground shadow-md transition-all hover:border-primary hover:bg-primary/5"
      @click="emit('expand')"
    >
      <XIcon icon="carbon:list" />
      <span class="max-md:hidden">{{ title }}</span>
      <XIcon icon="carbon:chevron-right" />
    </button>
  </div>
</template>
