<script setup lang="ts">
import { fetchListTag, getTagTranslation } from '@/api/composables/tag'
import { XIcon } from '@/plugins/xicon'
import { cn } from '@/lib/utils'

const { t } = useI18n()
const localePath = useLocalePath()

const tags = ref<any[]>([])
const loading = ref(true)
const page = ref(1)
const pageSize = ref(20)
const total = ref(0)

function getTagName(tag: any) {
  const translation = getTagTranslation(tag)
  return translation?.name || t('page.tags.tag_untitled')
}

function getTagDescription(tag: any) {
  const translation = getTagTranslation(tag)
  return translation?.description || ''
}

async function loadTags() {
  loading.value = true
  try {
    const res = await fetchListTag({
      paging: { page: page.value, pageSize: pageSize.value },
      formValues: { status: 'TAG_STATUS_ACTIVE' },
    }) as any
    tags.value = res?.items || []
    total.value = res?.total || 0
  } catch (e) {
    console.error('[Tag List] 加载失败:', e)
  } finally {
    loading.value = false
  }
}

function handlePageChange(newPage: number) {
  page.value = newPage
  loadTags()
}

onMounted(loadTags)
</script>

<template>
  <div class="w-full">
    <LayoutPageHero
      :title="t('page.tags.tags_list')"
      :description="t('page.tags.explore_all')"
      icon="carbon:tag"
      size="md"
    >
      <template #default>
        <div class="flex items-center gap-2 text-sm text-muted-foreground">
          <span>{{ total || tags.length }} {{ t('page.tags.total_tags') }}</span>
        </div>
      </template>
    </LayoutPageHero>

    <LayoutSectionContainer>
      <!-- Loading Skeleton -->
      <div v-if="loading" class="grid grid-cols-[repeat(auto-fill,minmax(260px,1fr))] gap-6">
        <div v-for="i in 12" :key="i">
          <UiSkeleton class="h-[120px] w-full" />
        </div>
      </div>

      <template v-else>
        <!-- Tags Grid -->
        <div v-if="tags.length > 0 || total > 0" class="grid grid-cols-[repeat(auto-fill,minmax(260px,1fr))] gap-6">
          <div
            v-for="tag in tags"
            :key="tag.id"
            :class="cn(
              'group flex cursor-pointer flex-col items-center gap-4 rounded-2xl border bg-card p-6',
              'transition-all duration-300 ease-out',
              'hover:-translate-y-1.5',
              'hover:shadow-[0_10px_30px_-4px_hsl(var(--primary)/0.12)]',
            )"
            :style="{
              borderColor: tag.color || 'hsl(var(--border))',
              background: `linear-gradient(135deg, ${tag.color}10 0%, hsl(var(--card)) 100%)`,
            }"
            @click="navigateTo(localePath(`/tag/${tag.id}`))"
          >
            <div
              class="flex h-12 w-12 items-center justify-center rounded-lg transition-transform duration-300 group-hover:scale-110"
              :style="{ color: tag.color || 'hsl(var(--primary))' }"
            >
              <XIcon :icon="`carbon:${tag.icon || 'tag'}`" :size="48" />
            </div>
            <div class="flex flex-1 flex-col items-center gap-2 text-center">
              <h3 class="text-[15px] font-semibold text-foreground transition-colors group-hover:text-primary">
                {{ getTagName(tag) }}
              </h3>
              <p class="line-clamp-2 min-h-[2.6em] text-xs leading-relaxed text-muted-foreground">
                {{ getTagDescription(tag) }}
              </p>
              <div class="mt-auto flex items-center gap-1 pt-2 text-xs font-medium text-muted-foreground">
                <XIcon icon="carbon:document" :size="14" />
                <span>{{ tag.postCount || 0 }} {{ t('page.posts.articles') }}</span>
              </div>
            </div>
          </div>

          <div v-if="tags.length === 0 && total > 0" class="col-span-full">
            <UiAppEmpty>
              {{ t('page.tags.no_tags_in_page') }}
            </UiAppEmpty>
          </div>
        </div>

        <!-- Empty -->
        <div v-if="tags.length === 0 && total === 0" class="py-20 text-center">
          <UiAppEmpty :in-container="true">
            {{ t('page.tags.no_tags') }}
          </UiAppEmpty>
        </div>

        <!-- Pagination -->
        <div v-if="total > pageSize" class="mt-8 flex items-center justify-center gap-2">
          <UiButton
            variant="outline"
            size="sm"
            :disabled="page <= 1"
            @click="handlePageChange(page - 1)"
          >
            {{ t('common.previous') || 'Previous' }}
          </UiButton>
          <span class="text-sm text-muted-foreground">
            {{ page }} / {{ Math.ceil(total / pageSize) }}
          </span>
          <UiButton
            variant="outline"
            size="sm"
            :disabled="page >= Math.ceil(total / pageSize)"
            @click="handlePageChange(page + 1)"
          >
            {{ t('common.next') || 'Next' }}
          </UiButton>
        </div>
      </template>
    </LayoutSectionContainer>
  </div>
</template>
