<script setup lang="ts">
import { XIcon } from '@/plugins/xicon'
import { cn } from '@/lib/utils'
import { getPostTitle, getPostSummary, getPostThumbnail } from '@/api/composables/post'
import { formatDate } from '@/utils/date'

const props = withDefaults(defineProps<{
  post: any
  from?: string
  categoryId?: number
}>(), {
  from: 'post-list',
})

const localePath = useLocalePath()

function handleViewPost() {
  const query: string[] = [`from=${props.from}`]
  if (props.categoryId) {
    query.push(`categoryId=${props.categoryId}`)
  }
  navigateTo(`${localePath(`/post/${props.post.id}`)}?${query.join('&')}`)
  window.scrollTo({ top: 0, behavior: 'smooth' })
}
</script>

<template>
  <article
    :class="cn(
      'group flex h-full cursor-pointer flex-col overflow-hidden',
      'rounded-2xl border border-border bg-card shadow-sm',
      'transition-all duration-500 ease-out',
      'hover:-translate-y-1.5 hover:border-primary/40',
      'hover:shadow-[0_20px_40px_-8px_hsl(var(--primary)/0.15)]',
    )"
    style="will-change: transform, box-shadow"
    @click="handleViewPost"
  >
    <div class="relative h-[240px] w-full flex-shrink-0 overflow-hidden bg-background max-md:h-[200px]">
      <UiImage
        :src="getPostThumbnail(post)"
        :alt="getPostTitle(post)"
        class="h-full w-full object-cover transition-transform duration-700 ease-out group-hover:scale-[1.12]"
      />
      <div class="absolute inset-0 bg-gradient-to-tr from-primary/20 via-transparent to-sky-400/20 opacity-0 transition-opacity duration-500 group-hover:opacity-100" />
    </div>
    <div class="flex flex-1 flex-col gap-3 p-6 max-md:p-4 max-md:gap-2.5">
      <h3 :class="cn(
        'line-clamp-2 min-h-[3.4em] text-lg font-bold leading-[1.7] text-foreground transition-colors duration-300',
        'group-hover:text-primary',
        'max-md:min-h-[3em] max-md:text-[17px]',
      )">
        {{ getPostTitle(post) }}
      </h3>
      <p :class="cn(
        'line-clamp-3 min-h-[4.4em] flex-1 text-sm leading-relaxed text-muted-foreground',
        'max-md:min-h-[4em] max-md:text-xs max-md:leading-relaxed',
      )">
        {{ getPostSummary(post) }}
      </p>
      <div :class="cn(
        'border-t border-border pt-3 text-[13px] font-medium text-muted-foreground',
        'flex flex-col gap-1.5',
        'max-md:text-xs',
      )">
        <div class="flex flex-wrap gap-4">
          <div class="flex items-center gap-1.5 whitespace-nowrap">
            <XIcon icon="carbon:user" :size="16" />
            <span>{{ post.authorName || '—' }}</span>
          </div>
          <div class="flex items-center gap-1.5 whitespace-nowrap">
            <XIcon icon="carbon:calendar" :size="16" />
            <span>{{ formatDate(post.createdAt) }}</span>
          </div>
        </div>
        <div class="flex flex-wrap gap-4">
          <div class="flex items-center gap-1.5 whitespace-nowrap">
            <XIcon icon="carbon:view" :size="16" />
            <span>{{ post.visits || 0 }}</span>
          </div>
          <div class="flex items-center gap-1.5 whitespace-nowrap">
            <XIcon icon="carbon:thumbs-up" :size="16" />
            <span>{{ post.likes || 0 }}</span>
          </div>
        </div>
      </div>
    </div>
  </article>
</template>
