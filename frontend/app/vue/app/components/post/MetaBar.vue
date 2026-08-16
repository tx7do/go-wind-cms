<script setup lang="ts">
/**
 * 文章元数据条 — 作者 / 日期 / 点赞数 / 收藏数
 *
 * 用于 post/[id] 详情页 header 区域
 */
import { XIcon } from '@/plugins/xicon'
import { formatDate } from '@/utils/date'

defineProps<{
  authorName?: string
  createdAt?: string
  likes?: number | string
  watchCount?: number | string
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
                <XIcon icon="carbon:thumbs-up" :size="16" />
                <span>{{ likes || 0 }}</span>
            </div>
            <div class="flex items-center gap-1.5">
                <XIcon icon="carbon:bookmark" :size="16" />
                <span>{{ watchCount || 0 }}</span>
            </div>
  </div>
</template>
