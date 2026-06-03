<script setup lang="ts">
import { XIcon } from '@/plugins/xicon'
import { cn } from '@/lib/utils'
import { getCategoryName, getCategoryDescription, getCategoryThumbnail } from '@/api/composables/category'

const props = withDefaults(defineProps<{
  category: any
  clickable?: boolean
}>(), {
  clickable: true,
})

const emit = defineEmits<{
  click: [id: number]
}>()

const { t } = useI18n()

function handleClick() {
  if (!props.category?.id || !props.clickable) return
  emit('click', props.category.id)
}
</script>

<template>
  <div
    v-if="category"
    :class="cn(
      'group flex h-full flex-col overflow-hidden rounded-xl border border-border bg-background',
      'transition-all duration-500 ease-out',
      clickable && 'cursor-pointer hover:-translate-y-1.5 hover:border-primary/40 hover:shadow-[0_20px_40px_-8px_hsl(var(--primary)/0.15)]',
      !clickable && 'cursor-default',
    )"
    @click="handleClick"
  >
    <div class="relative h-[160px] w-full overflow-hidden bg-background">
      <UiImage
        :src="getCategoryThumbnail(category)"
        :alt="getCategoryName(category, t('page.categories.category_untitled'))"
        :class="cn(
          'h-full w-full object-cover transition-transform duration-700 ease-out',
          clickable && 'group-hover:scale-[1.12]',
        )"
      />
      <div :class="cn(
        'absolute inset-0 bg-gradient-to-tr from-primary/20 via-transparent to-sky-400/20 transition-opacity duration-500',
        clickable ? 'opacity-0 group-hover:opacity-100' : 'opacity-0',
      )" />
    </div>
    <div class="p-4">
      <h3 :class="cn(
        'mb-2 line-clamp-2 text-base font-bold leading-tight text-foreground transition-colors',
        clickable && 'group-hover:text-primary',
      )">
        {{ getCategoryName(category, t('page.categories.category_untitled')) }}
      </h3>
      <p class="mb-3 line-clamp-2 text-[13px] leading-relaxed text-muted-foreground">
        {{ getCategoryDescription(category) }}
      </p>
      <div class="flex items-center gap-1.5 text-xs font-medium text-muted-foreground">
        <span class="flex items-center">
          <XIcon icon="carbon:document" :size="14" />
        </span>
        <span>
          {{ category.postCount || 0 }} {{ t('page.categories.articles_count') }}
        </span>
      </div>
    </div>
  </div>
</template>
