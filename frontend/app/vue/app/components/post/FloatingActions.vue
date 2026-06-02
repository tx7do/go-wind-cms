<script setup lang="ts">
import { XIcon } from '@/plugins/xicon'

defineProps<{
  isLiked: boolean
  isBookmarked: boolean
  labels: {
    likes: string
    bookmark: string
    share: string
  }
}>()

const emit = defineEmits<{
  like: []
  bookmark: []
  share: []
}>()
</script>

<template>
  <div class="fixed top-1/2 right-8 z-10 hidden -translate-y-1/2 flex-col gap-3 lg:flex">
    <button
      :class="[
        'flex h-12 w-12 cursor-pointer items-center justify-center rounded-full border border-border bg-card shadow-md transition-all hover:border-primary hover:bg-primary/5 hover:shadow-lg',
        isLiked ? 'border-primary bg-primary/10 text-primary' : 'text-muted-foreground',
      ]"
      :aria-label="labels.likes"
      @click="emit('like')"
    >
      <XIcon :icon="isLiked ? 'carbon:thumbs-up-filled' : 'carbon:thumbs-up'" :size="20" />
    </button>
    <button
      :class="[
        'flex h-12 w-12 cursor-pointer items-center justify-center rounded-full border border-border bg-card shadow-md transition-all hover:border-primary hover:bg-primary/5 hover:shadow-lg',
        isBookmarked ? 'border-primary bg-primary/10 text-primary' : 'text-muted-foreground',
      ]"
      :aria-label="labels.bookmark"
      @click="emit('bookmark')"
    >
      <XIcon :icon="isBookmarked ? 'carbon:bookmark-filled' : 'carbon:bookmark'" :size="20" />
    </button>
    <button
      class="flex h-12 w-12 cursor-pointer items-center justify-center rounded-full border border-border bg-card text-muted-foreground shadow-md transition-all hover:border-primary hover:bg-primary/5 hover:shadow-lg hover:text-primary"
      :aria-label="labels.share"
      @click="emit('share')"
    >
      <XIcon icon="carbon:share" :size="20" />
    </button>
  </div>
</template>
