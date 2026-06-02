<script setup lang="ts">
import { XIcon } from '@/plugins/xicon'
import { formatDate } from '@/utils/date'

defineProps<{
  authorName?: string
  createdAt?: string
  visits?: number | string
  likes?: number | string
}>()
</script>

<template>
  <div class="flex flex-wrap items-center gap-x-5 gap-y-2 pb-5 text-sm font-medium text-muted-foreground border-b border-border/40">
    <div v-if="authorName" class="flex items-center gap-1.5">
      <XIcon icon="carbon:user-avatar" :size="16" />
      <span>{{ authorName }}</span>
    </div>
    <div v-if="createdAt" class="flex items-center gap-1.5">
      <XIcon icon="carbon:calendar" :size="16" />
      <span>{{ formatDate(createdAt) }}</span>
    </div>
    <div class="flex items-center gap-1.5">
      <XIcon icon="carbon:view" :size="16" />
      <span>{{ visits || 0 }}</span>
    </div>
    <div class="flex items-center gap-1.5">
      <XIcon icon="carbon:thumbs-up" :size="16" />
      <span>{{ likes || 0 }}</span>
    </div>
  </div>
</template>
