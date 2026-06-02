<script setup lang="ts">
import { fetchTag, getTagTranslation } from '@/api/composables/tag'

const { t } = useI18n()
const localePath = useLocalePath()
const route = useRoute()

const tag = ref<any>(null)
const loading = ref(true)
const tagId = computed(() => Number(route.params.id))

// 动态标题
const tagTitle = computed(() =>
  tag.value ? getTagName(tag.value) : t('page.tags.tags_list')
)
useHead({ title: () => tagTitle.value })

function getTagName(tagData: any) {
  const translation = getTagTranslation(tagData)
  return translation?.name || t('page.tags.tag_untitled')
}

function getTagDescription(tagData: any) {
  const translation = getTagTranslation(tagData)
  return translation?.description || ''
}

async function loadTag() {
  if (!tagId.value) return
  loading.value = true
  try {
    tag.value = await fetchTag(tagId.value) as any
  } catch (e) {
    console.error('[Tag Detail] 加载失败:', e)
  } finally {
    loading.value = false
  }
}

function handleBack() {
  navigateTo(localePath('/tag'))
}

onMounted(loadTag)
</script>

<template>
  <div class="w-full">
    <!-- Loading -->
    <div v-if="loading" class="w-full py-20">
      <LayoutSectionContainer>
        <div class="mx-auto max-w-3xl space-y-4">
          <UiSkeleton class="h-10 w-1/2" />
          <UiSkeleton class="h-5 w-3/4" />
        </div>
      </LayoutSectionContainer>
    </div>

    <!-- Tag Detail -->
    <template v-else-if="tag">
      <LayoutPageHero
        :title="getTagName(tag)"
        :description="getTagDescription(tag) || undefined"
        :icon="`carbon:${tag.icon || 'tag'}`"
        :icon-size="56"
        size="lg"
        :accent-color="tag.color || undefined"
      >
        <template #default>
          <div class="flex items-center gap-2">
            <XIcon icon="carbon:document" :size="16" />
            <span>{{ tag.postCount || 0 }} {{ t('page.posts.articles') }}</span>
          </div>
        </template>
      </LayoutPageHero>

      <LayoutSectionContainer>
        <!-- Back Button -->
        <div class="mb-8">
          <LayoutBackButton
            :label="t('page.categories.back_to_list')"
            :on-click="handleBack"
          />
        </div>

        <!-- Posts List with Pagination -->
        <PostList
          v-if="tagId"
          :key="tagId"
          :initial-page-size="10"
          :tag-id="tagId"
          from="tag"
        />
      </LayoutSectionContainer>
    </template>

    <!-- Tag not found / Invalid ID -->
    <div v-else class="w-full py-20">
      <LayoutSectionContainer>
        <UiAppEmpty variant="error" :description="t('page.tags.no_tags')" />
      </LayoutSectionContainer>
    </div>
  </div>
</template>
