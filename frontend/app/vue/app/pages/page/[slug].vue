<script setup lang="ts">
import type { contentservicev1_Page } from '@/api/generated/app/service/v1'
import { fetchPage, fetchPageBySlug, fetchPagePreview, getPageTranslation } from '@/api/composables/page'
import { recordView } from '@/api/composables/interaction'

/**
 * CMS Page 公开渲染页：
 *   /page/[slug] —— slug 参数为纯数字时按 id 查询，否则按 slug。
 * 仅返回已发布页面（app BFF Get 过滤），草稿/归档按未找到处理。
 */
const route = useRoute()
const { t } = useI18n()

const slugParam = computed(() => String(route.params.slug || ''))
// 草稿预览：admin「预览」按钮带短时效令牌跳转，后端校验后放行非已发布内容
const previewToken = computed(() => String(route.query.preview_token || ''))

const page = ref<contentservicev1_Page | null>(null)
const loading = ref(true)

async function load() {
  if (!slugParam.value) {
    loading.value = false
    return
  }
  loading.value = true
  try {
    const numericId = Number(slugParam.value)
    const fetched = previewToken.value
      ? await fetchPagePreview(numericId > 0 ? numericId : 0, previewToken.value)
      : await fetchPageBySlug(slugParam.value)
    page.value = fetched ?? null
  } catch (e) {
    console.error('[CMS Page] Load failed:', e)
    page.value = null
  } finally {
    loading.value = false
  }
}

watch([slugParam, previewToken], load, { immediate: true })

// 浏览量上报（PV +1）：页面加载成功且拿到 id 后记一次
const viewReportedId = ref<number | null>(null)
watch(
  () => page.value?.id,
  (pid) => {
    if (!pid || viewReportedId.value === pid) return
    viewReportedId.value = pid
    recordView('TARGET_TYPE_PAGE', pid)
  },
)

// SSR/CSR 都生效的标题与 SEO meta
useHead(() => {
  const tr = page.value ? getPageTranslation(page.value) : null
  const seo = tr?.seo
  return {
    title: () => tr?.title || t('page.post_detail.untitled'),
    meta: seo?.metaDescription
      ? [{ name: 'description', content: seo.metaDescription }]
      : [],
  }
})
</script>

<template>
  <!-- Loading -->
  <div v-if="loading" class="mx-auto w-full max-w-[900px] px-8 py-12 max-md:px-4">
    <div class="mb-8 h-9 w-32 animate-pulse rounded bg-muted" />
    <div class="mb-6 h-10 w-2/3 animate-pulse rounded bg-muted" />
    <div class="space-y-4">
      <div class="h-4 w-full animate-pulse rounded bg-muted" />
      <div class="h-4 w-[90%] animate-pulse rounded bg-muted" />
      <div class="h-4 w-[95%] animate-pulse rounded bg-muted" />
    </div>
  </div>

  <!-- Not found -->
  <LayoutSectionContainer v-else-if="!page">
    <div class="flex w-full items-center justify-center py-20">
      <div class="text-lg text-muted-foreground">{{ t('page.post_detail.post_not_found') }}</div>
    </div>
  </LayoutSectionContainer>

  <!-- Page -->
  <PageCmsPageView v-else :page="page" :back-label="t('page.post_detail.back')" />
</template>
