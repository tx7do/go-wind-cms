<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'
import { XIcon } from '@/plugins/xicon'
import { fetchListTag } from '@/api/composables/tag'

const { t } = useI18n()
const localePath = useLocalePath()

interface TagItem {
  id: number
  name: string
  color: string
  postCount: number
}

const loading = ref(true)
const displayTags = ref<TagItem[]>([])
const hoveredId = ref<number | null>(null)

let abortController: AbortController | null = null

onMounted(async () => {
  loading.value = true
  try {
    const res = await fetchListTag({
      paging: { page: 1, pageSize: 6 },
      formValues: { status: 'TAG_STATUS_ACTIVE', isFeatured: true },
    }) as any

    const tagItems = res.items || []
    displayTags.value = tagItems
      .filter((tag: any) => tag.id !== undefined)
      .map((tag: any, index: number) => ({
        id: tag.id,
        name: tag.translations?.[0]?.name || t('page.tags.tag_untitled'),
        color: tag.color || `hsl(${(index * 67) % 360}, 70%, 55%)`,
        postCount: tag.postCount || 0,
      }))
  } catch (error) {
    console.error('Failed to load tags:', error)
  } finally {
    loading.value = false
  }
})

onUnmounted(() => { abortController?.abort() })
</script>

<template>
  <section class="w-full max-w-300 mx-auto px-8 py-12 max-md:px-4">
    <div class="mb-8 flex items-center justify-between">
      <h2 class="flex items-center gap-2 text-2xl font-extrabold tracking-tight text-foreground max-md:text-xl">
        <XIcon icon="carbon:fire" :size="28" class="mr-2 text-primary" />
        {{ t('page.tags.popular_tags') }}
      </h2>
      <UiButton variant="ghost" @click="navigateTo(localePath('/tag'))">
        {{ t('page.tags.view_all') }} →
      </UiButton>
    </div>
    <div class="w-full">
      <!-- Loading -->
      <div v-if="loading" class="flex flex-wrap justify-center gap-3">
        <div v-for="i in 6" :key="i" class="h-10 w-28">
          <UiSkeleton class="h-full w-full rounded-full" />
        </div>
      </div>
      <!-- Tags -->
      <div v-else class="flex flex-wrap justify-center gap-3 max-md:gap-2">
        <button
          v-for="tag in displayTags"
          :key="tag.id"
          class="group flex cursor-pointer items-center gap-2 rounded-full border-2 px-5 py-2.5 text-sm font-semibold transition-all duration-300 hover:-translate-y-0.5 hover:shadow-lg max-md:px-4 max-md:py-2 max-md:text-xs"
          :style="hoveredId === tag.id ? {
            backgroundColor: tag.color,
            color: '#ffffff',
            borderColor: tag.color,
            boxShadow: `0 8px 24px -8px ${tag.color}80`,
          } : {
            backgroundColor: `${tag.color}12`,
            color: tag.color,
            borderColor: `${tag.color}40`,
          }"
          @mouseenter="hoveredId = tag.id"
          @mouseleave="hoveredId = null"
          @click="navigateTo(localePath(`/tag/${tag.id}`))"
        >
          <span class="h-2 w-2 rounded-full transition-all duration-300" :style="{ backgroundColor: hoveredId === tag.id ? '#fff' : tag.color }" />
          <span class="truncate tracking-tight">{{ tag.name }}</span>
          <span
            v-if="tag.postCount > 0"
            class="rounded-full px-2 py-0.5 text-[10px] font-bold leading-none transition-colors duration-300"
            :style="hoveredId === tag.id ? { backgroundColor: 'rgba(255,255,255,0.25)', color: '#fff' } : { backgroundColor: `${tag.color}25` }"
          >
            {{ tag.postCount }}
          </span>
        </button>
      </div>
    </div>
  </section>
</template>
