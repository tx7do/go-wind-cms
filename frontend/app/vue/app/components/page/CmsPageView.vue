<script setup lang="ts">
import type { contentservicev1_Page } from '@/api/generated/app/service/v1'
import { getPageTranslation } from '@/api/composables/page'

/**
 * CMS Page 前台渲染视图：标题 + 封面 + SectionRenderer。
 * redirectUrl 非空时按页面配置跳转（对齐 Page.redirect_url 语义）。
 */
interface CmsPageViewProps {
  page: contentservicev1_Page
  backLabel?: string
}

const props = defineProps<CmsPageViewProps>()

const localePath = useLocalePath()

const title = computed(() => getPageTranslation(props.page)?.title || '')
const cover = computed(() => props.page.thumbnail || '')

function goBack() {
  if (import.meta.client && window.history.length > 2) {
    window.history.back()
    return
  }
  navigateTo(localePath('/'))
}

watch(
  () => props.page.redirectUrl,
  (url) => {
    if (!url || !import.meta.client) return
    if (/^https?:\/\//.test(url)) {
      window.location.href = url
    } else {
      navigateTo(localePath(url))
    }
  },
  { immediate: true },
)
</script>

<template>
  <div class="w-full">
    <LayoutSectionContainer v-if="props.backLabel" no-padding class="!py-6">
      <LayoutBackButton :label="props.backLabel" :on-click="goBack" />
    </LayoutSectionContainer>

    <article class="mx-auto w-full max-w-[900px] px-8 max-md:px-4">
      <div v-if="cover" class="relative mb-8 h-[300px] overflow-hidden rounded-xl max-md:h-[200px]">
        <UiImage :src="cover" :alt="title" class="h-full w-full object-cover" />
      </div>
      <h1
        v-if="title"
        class="mb-8 text-3xl font-bold leading-tight text-foreground max-md:text-2xl"
      >
        {{ title }}
      </h1>
      <PageSectionRenderer :sections="props.page.sections" />
    </article>
  </div>
</template>
