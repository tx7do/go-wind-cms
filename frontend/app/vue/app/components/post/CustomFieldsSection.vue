<script setup lang="ts">
import {useQuery} from '@tanstack/vue-query'
import {apiClient} from '@/api/client'

/**
 * 自定义字段渲染（内容建模前台读侧）：
 * 按文章分类链绑定的内容模型字段定义，渲染 post.customFields。
 * 值约定为字符串（admin 编辑器 Record<string,string> 口径）：
 *   - RICHTEXT 按富文本注入；IMAGE/FILE 值形似 URL 时按图片/链接渲染；
 *   - 其余（TEXT/NUMBER/RELATION 等）原样展示。
 * 模型未绑定或加载失败静默不渲染（自定义字段是增强信息，不阻断正文）。
 */
const props = defineProps<{
  postId?: number | null
  customFields?: Record<string, string>
}>()

const {t} = useI18n()

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
          <span v-else class="break-all">{{ row.value }}</span>
        </dd>
      </template>
    </dl>
  </section>
</template>
