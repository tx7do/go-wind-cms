<script setup lang="ts">
import { fetchListTag, getTagTranslation } from '@/api/composables/tag'
import type { ListTagParams } from '@/api/composables/tag'
import { XIcon } from '@/plugins/xicon'

const { t } = useI18n()
const localePath = useLocalePath()

const tags = ref<any[]>([])
const loading = ref(true)

function getTagName(tag: any) {
  const translation = getTagTranslation(tag)
  return translation?.name || t('page.tags.tag_untitled')
}

onMounted(async () => {
  try {
    const params: ListTagParams = {
      paging: { page: 1, pageSize: 100 },
      isTenantUser: true,
    }
    const res = await fetchListTag(params)
    tags.value = res?.items || res?.data || []
  } catch (e) {
    console.error('[Tag List] 加载失败:', e)
  } finally {
    loading.value = false
  }
})
</script>

<template>
  <div class="w-full">
    <LayoutPageHero
      :title="t('page.tags.tags_list')"
      icon="carbon:tag"
      size="md"
    />

    <section class="w-full py-12 max-md:py-8">
      <LayoutSectionContainer>
        <!-- Loading -->
        <div v-if="loading" class="flex flex-wrap gap-3">
          <UiSkeleton v-for="i in 12" :key="i" class="h-8 w-20 rounded-full" />
        </div>

        <!-- Empty -->
        <div v-else-if="tags.length === 0" class="py-20 text-center">
          <UiAppEmpty :message="t('page.tags.no_tags')" />
        </div>

        <!-- Tag Cloud -->
        <div v-else class="flex flex-wrap gap-3">
          <NuxtLink
            v-for="tag in tags"
            :key="tag.id"
            :to="localePath(`/tag/${tag.id}`)"
            class="inline-flex items-center gap-1.5 rounded-full border border-border bg-card px-4 py-2 text-sm font-medium text-foreground transition-colors hover:border-primary/30 hover:bg-primary/10 hover:text-primary"
          >
            <XIcon icon="carbon:tag" :size="14" />
            {{ getTagName(tag) }}
          </NuxtLink>
        </div>
      </LayoutSectionContainer>
    </section>
  </div>
</template>
