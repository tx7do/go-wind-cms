<script setup lang="ts">
import {useQuery} from '@tanstack/vue-query'
import {apiClient} from '@/api/client'

/**
 * 自定义字段渲染（内容建模前台读侧）：
 * 值约定为字符串（admin 编辑器 Record<string,string> 口径）：
 *   - RICHTEXT 按富文本注入；IMAGE/FILE 值形似 URL 时按图片/链接渲染；
 *   - RELATION 引用解析为公开摘要链接（缺失回退原样字符串）；
 *   - COMPONENT/DYNAMIC_ZONE 值为组件实例 JSON（模型响应携带组件 schema），
 *     解析后渲染键值卡（动态区按 _component 分发），解析失败回退原样字符串；
 *   - 其余原样展示。模型未绑定或加载失败静默不渲染。
 */
const props = defineProps<{
  postId?: number | null
  customFields?: Record<string, string>
}>()

const {t, locale} = useI18n()
const localePath = useLocalePath()

/** RELATION 引用值形态（与 core contentRefPattern 对齐） */
const RELATION_REF_RE = /^(post|page):\d+$/

interface FieldDef {
  name?: string
  label?: string
  type?: string
  componentConfig?: {
    componentCode?: string
    repeatable?: boolean
    allowedComponentCodes?: string[]
  }
}

interface InstanceData {
  _component?: string
  [key: string]: unknown
}

function parseInstances(value: string, expectArray: boolean): InstanceData[] | null {
  if (!value.trim()) return []
  try {
    const parsed: unknown = JSON.parse(value)
    if (expectArray) {
      if (!Array.isArray(parsed)) return null
      return parsed.filter(
        (x): x is InstanceData => Boolean(x) && typeof x === 'object' && !Array.isArray(x),
      )
    }
    if (!parsed || typeof parsed !== 'object' || Array.isArray(parsed)) return null
    return [parsed as InstanceData]
  } catch {
    return null
  }
}

function isImageUrl(v: string) {
  return /^https?:\/\//.test(v)
}

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

const componentMap = computed(() => {
  const map: Record<string, { name: string; fields: FieldDef[] }> = {}
  for (const c of modelQuery.data.value?.components ?? []) {
    if (c.definition?.code) {
      map[c.definition.code] = {name: c.definition.name ?? '', fields: c.fields ?? []}
    }
  }
  return map
})

const isUrl = (v: string) => /^https?:\/\//.test(v)

interface PairRow {
  label: string
  text: string
  image: boolean
}

function instanceRows(schemaFields: FieldDef[], inst: InstanceData): PairRow[] {
  const out: PairRow[] = []
  for (const f of schemaFields) {
    const key = f.name ?? ''
    const raw = inst[key]
    if (raw == null || typeof raw !== 'string' || raw === '') continue
    const isImg = f.type === 'FIELD_TYPE_IMAGE' && isImageUrl(raw)
    out.push({label: f.label || key, text: raw, image: isImg})
  }
  return out
}

const componentValue = (field: FieldDef, value: string) => {
  const cc = field.componentConfig
  if (!cc) return null
  const expectArray = !!cc.repeatable || field.type === 'FIELD_TYPE_DYNAMIC_ZONE'
  const instances = parseInstances(value, expectArray)
  if (instances === null) return null
  return instances.map((inst, i) => {
    const code = field.type === 'FIELD_TYPE_DYNAMIC_ZONE'
      ? String(inst._component ?? '')
      : cc.componentCode ?? ''
    const schema = componentMap.value[code]
    return {
      name: schema?.name ?? code,
      pairs: instanceRows(schema?.fields ?? [], inst),
      key: i,
    }
  })
}

const fieldByName = (name: string): FieldDef | undefined =>
  (modelQuery.data.value?.fields ?? []).find((f) => f.name === name)
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
          <div
            v-else-if="
              (row.type === 'FIELD_TYPE_COMPONENT' || row.type === 'FIELD_TYPE_DYNAMIC_ZONE') &&
                componentValue(fieldByName(row.name) ?? {}, row.value) !== null
            "
            class="flex flex-col gap-2"
          >
            <div
              v-for="card in componentValue(fieldByName(row.name) ?? {}, row.value)!"
              :key="card.key"
              class="rounded-lg border border-border p-3"
            >
              <div class="mb-1 text-xs font-medium text-muted-foreground">{{ card.name }}</div>
              <span v-if="card.pairs.length === 0" class="text-xs text-muted-foreground">—</span>
              <div v-for="(p, j) in card.pairs" :key="j" class="flex flex-col py-0.5">
                <span class="text-xs text-muted-foreground">{{ p.label }}</span>
                <img
                  v-if="p.image"
                  :src="p.text"
                  :alt="p.label"
                  class="max-h-48 rounded-md object-cover"
                />
                <span v-else class="break-all text-sm">{{ p.text }}</span>
              </div>
            </div>
          </div>
          <span v-else class="break-all">{{ row.value }}</span>
        </dd>
      </template>
    </dl>
  </section>
</template>
