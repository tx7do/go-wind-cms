<script setup lang="ts">
import {useQuery} from '@tanstack/vue-query'
import {apiClient} from '@/api/client'

/**
 * 自定义字段渲染（内容建模前台读侧）：
 * 按文章分类链绑定的内容模型字段定义，渲染 post.customFields。
 * 值约定为字符串（admin 编辑器 Record<string,string> 口径）：
 *   - RICHTEXT 按富文本注入；IMAGE/FILE 值形似 URL 时按图片/链接渲染；
 *   - RELATION 值为 "post:12"/"page:7" 引用时批量解析为公开摘要，
 *     渲染成指向目标详情页的链接（草稿/软删目标解析缺失，回退原样字符串）；
 *   - 其余（TEXT/NUMBER 等）原样展示。
 * 模型未绑定或加载失败静默不渲染（自定义字段是增强信息，不阻断正文）。
 */
const props = defineProps<{
  postId?: number | null
  customFields?: Record<string, string>
}>()

const {t, locale} = useI18n()
const localePath = useLocalePath()

/** RELATION 引用值形态（与 core contentRefPattern 对齐；category/tag v1 不解析） */
const RELATION_REF_RE = /^(post|page):\d+$/

const modelQuery = useQuery({
  queryKey: ['postContentModel', () => props.postId],
  queryFn: () => apiClient.postService.GetPostContentModel({postId: props.postId!}),
  enabled: computed(() => !!props.postId),
  retry: 0,
  staleTime: 5 * 60 * 1000,
})

const rows = computed(() => {
  const fields = modelQuery.data.value?.fields ?? []
  const cf = props.customFields ?? {}
  return fields
    .filter((f) => cf[f.name ?? ''] !== undefined && cf[f.name ?? ''] !== '')
    .map((f) => ({
      label: f.label || f.name || '',
      type: f.type ?? 'FIELD_TYPE_TEXT',
      name: f.name ?? '',
      value: cf[f.name ?? ''],
    }))
})

const relationRefs = computed(() =>
  Array.from(new Set(rows.value
    .filter((r) => r.type === 'FIELD_TYPE_RELATION' && RELATION_REF_RE.test(r.value.trim()))
    .map((r) => r.value.trim()))),
)

const resolveQuery = useQuery({
  queryKey: ['contentReferences', () => relationRefs.value.join(','), locale],
  queryFn: () =>
    apiClient.postService.ResolveContentReferences({
      refs: relationRefs.value,
      languageCode: locale.value,
    }),
  enabled: computed(() => relationRefs.value.length > 0),
  retry: 0,
  staleTime: 5 * 60 * 1000,
})

const resolvedByRef = computed(() => {
  const map = new Map<string, { title: string; href: string }>()
  for (const it of resolveQuery.data.value?.items ?? []) {
    if (!it.ref || !it.slug) continue
    // vue 端 post 详情按 id 路由、page 详情按 slug 路由
    const href = it.targetEntityType === 'page'
      ? localePath(`/page/${it.slug}`)
      : localePath(`/post/${it.targetId}`)
    map.set(it.ref, {title: it.title || it.ref, href})
  }
  return map
})

const relationTarget = (value: string) => {
  const ref = value.trim()
  return RELATION_REF_RE.test(ref) ? resolvedByRef.value.get(ref) : undefined
}

const isUrl = (v: string) => /^https?:\/\//.test(v)
</script>

<template>
  <section
    v-if="rows.length > 0"
    class="mb-8 rounded-xl border border-border bg-card p-6"
  >
    <h2 class="mb-4 text-lg font-semibold text-foreground">
      {{ modelQuery.data.value?.name || t('page.post_detail.custom_fields') }}
    </h2>
    <dl class="grid grid-cols-1 gap-x-8 gap-y-3 sm:grid-cols-[180px_1fr]">
      <template v-for="row in rows" :key="row.name">
        <dt class="text-sm font-medium text-muted-foreground">{{ row.label }}</dt>
        <dd class="min-w-0 text-sm text-foreground">
          <div
            v-if="row.type === 'FIELD_TYPE_RICHTEXT'"
            class="max-w-none text-sm [&_a]:text-primary [&_img]:max-w-full"
            v-html="row.value"
          />
          <img
            v-else-if="row.type === 'FIELD_TYPE_IMAGE' && isUrl(row.value)"
            :src="row.value"
            :alt="row.label"
            class="max-h-64 rounded-lg object-cover"
          />
          <a
            v-else-if="row.type === 'FIELD_TYPE_FILE' && isUrl(row.value)"
            :href="row.value"
            target="_blank"
            rel="noopener noreferrer"
            class="break-all text-primary hover:underline"
          >
            {{ row.value }}
          </a>
          <a
            v-else-if="relationTarget(row.value)"
            :href="relationTarget(row.value)!.href"
            class="break-all text-primary hover:underline"
          >
            {{ relationTarget(row.value)!.title }}
          </a>
          <span v-else class="break-all">{{ row.value }}</span>
        </dd>
      </template>
    </dl>
  </section>
</template>
