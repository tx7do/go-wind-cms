<script setup lang="ts">
import { computed, ref } from 'vue'
import type {
  contentservicev1_Section,
  contentservicev1_SectionTranslation,
} from '@/api/generated/app/service/v1'
import { getCurrentLocale } from '@/utils/locale'
import { apiClient } from '@/api/client'

/**
 * Section 内容键契约（三端渲染器共用，与 react SectionRenderer.tsx 保持同步）
 *
 * 语言相关内容存 translations[].content（map），主内容键一律为 `value`；
 * 语言无关样式存 config（map）。缺译回退：匹配当前语言 → 第一条翻译。
 * 键明细见 react 端文件头表格：RICH_TEXT/MARKDOWN/TITLE/IMAGE/GALLERY/VIDEO/
 * BUTTON/DIVIDER/SPACER/CODE/HTML/FORM/CAROUSEL/CUSTOM。
 */

interface SectionRendererProps {
  sections?: contentservicev1_Section[]
  class?: string
}

const props = defineProps<SectionRendererProps>()

/** 取区块当前语言的翻译：匹配当前语言 → 回退第一条 */
function pickTranslation(
  section: contentservicev1_Section,
): contentservicev1_SectionTranslation | null {
  if (!section?.translations || section.translations.length === 0) return null
  const locale = getCurrentLocale()
  return (
    section.translations.find(
      (tr: contentservicev1_SectionTranslation) => tr.languageCode === locale,
    ) || section.translations[0]
  )
}

function getContent(section: contentservicev1_Section): string {
  return pickTranslation(section)?.content?.value || ''
}

/** GALLERY/CAROUSEL 图片列表：value 为 JSON 数组（URL 或 {src,caption}）或换行分隔 URL */
function getImageList(section: contentservicev1_Section): { src: string; caption?: string }[] {
  const raw = getContent(section).trim()
  if (!raw) return []
  if (raw.startsWith('[')) {
    try {
      const parsed = JSON.parse(raw)
      if (Array.isArray(parsed)) {
        return parsed
          .map((item: unknown) => {
            if (typeof item === 'string') return { src: item }
            if (item && typeof item === 'object') {
              const obj = item as { src?: string; url?: string; caption?: string }
              const src = obj.src || obj.url || ''
              return src ? { src, caption: obj.caption } : null
            }
            return null
          })
          .filter((x): x is { src: string; caption?: string } => x !== null)
      }
    } catch {
      // JSON 解析失败时按纯 URL 文本处理，不中断整页渲染
    }
  }
  return raw
    .split(/\r?\n/)
    .map((line) => line.trim())
    .filter(Boolean)
    .map((src) => ({ src }))
}

const sortedSections = computed(() =>
  [...(props.sections ?? [])].sort((a, b) => (a.sortOrder ?? 0) - (b.sortOrder ?? 0)),
)

const alignClass = (config?: Record<string, string>) =>
  config?.align === 'center' ? 'text-center' : config?.align === 'right' ? 'text-end' : ''

function titleLevel(section: contentservicev1_Section): number {
  return Math.min(Math.max(Number(section.config?.level) || 2, 1), 6)
}

function titleSizeClass(level: number): string {
  return ['text-4xl', 'text-3xl', 'text-2xl', 'text-xl', 'text-lg', 'text-base'][level - 1] || 'text-2xl'
}

function galleryColumns(section: contentservicev1_Section): string {
  const columns = Math.min(Math.max(Number(section.config?.columns) || 3, 1), 4)
  return `repeat(${columns}, minmax(0, 1fr))`
}

function spacerHeight(section: contentservicev1_Section): string {
  const height = Math.min(Math.max(Number(section.config?.height) || 48, 8), 400)
  return `${height}px`
}

/** YouTube / Bilibili 站点链接转嵌入地址；直链视频返回 null 走 <video> */
function videoEmbed(url: string): string | null {
  try {
    const u = new URL(url)
    if (u.hostname.includes('youtube.com')) {
      const id = u.searchParams.get('v')
      if (id) return `https://www.youtube.com/embed/${id}`
    }
    if (u.hostname === 'youtu.be') {
      return `https://www.youtube.com/embed${u.pathname}`
    }
    if (u.hostname.includes('bilibili.com')) {
      const id = u.searchParams.get('bvid')
      if (id) return `https://player.bilibili.com/player.html?bvid=${id}`
    }
  } catch {
    // 非 URL 按 <video> 直链处理
  }
  return null
}

function isExternal(href: string): boolean {
  return /^https?:\/\//.test(href)
}

function sectionKey(section: contentservicev1_Section, index: number): string {
  return section.id != null ? `s-${section.id}` : `s-${index}`
}

/** FORM 区块字段定义（config.fields JSON 数组元素；与 core 侧校验契约一致） */
interface FormFieldDef {
  name: string
  type?: string
  required?: boolean
}

// FORM 区块状态：按区块 id 隔离（页面可能多个表单区块）
const formValues = ref<Record<string, Record<string, string>>>({})
const formState = ref<Record<string, 'idle' | 'submitting' | 'done' | 'error'>>({})

function formFields(section: contentservicev1_Section): FormFieldDef[] {
  try {
    const parsed = JSON.parse(section.config?.fields || '[]')
    if (Array.isArray(parsed)) {
      return parsed.filter(
        (f: unknown): f is FormFieldDef => !!f && typeof (f as FormFieldDef).name === 'string',
      )
    }
  } catch {
    // 字段定义非法：渲染为空表单（仅标题），不中断整页
  }
  return []
}

function formStateOf(section: contentservicev1_Section): string {
  return formState.value[section.id ?? 0] ?? 'idle'
}

/** 表单值记录（惰性初始化，模板 v-model 直接绑定到返回对象） */
function formValuesOf(section: contentservicev1_Section): Record<string, string> {
  const sid = section.id ?? 0
  if (!formValues.value[sid]) formValues.value[sid] = {}
  return formValues.value[sid]
}

function fieldLabel(section: contentservicev1_Section, field: FormFieldDef): string {
  return pickTranslation(section)?.content?.[field.name] || field.name
}

async function submitForm(section: contentservicev1_Section) {
  const sid = section.id ?? 0
  if (!section.pageId || formState.value[sid] === 'submitting') return
  formState.value[sid] = 'submitting'
  try {
    await apiClient.formSubmissionService.Submit({
      pageId: section.pageId,
      sectionId: sid,
      data: formValues.value[sid] ?? {},
    })
    formState.value[sid] = 'done'
    formValues.value[sid] = {}
  } catch (e) {
    console.error('Form submit failed:', e)
    formState.value[sid] = 'error'
  }
}
</script>

<template>
  <div v-if="sortedSections.length > 0" class="space-y-8" :class="props.class">
    <div
      v-for="(section, index) in sortedSections"
      :key="sectionKey(section, index)"
      :class="section.config?.className"
    >
      <!-- 富文本 -->
      <ContentViewer
        v-if="section.type === 'SECTION_TYPE_RICH_TEXT'"
        :content="getContent(section)"
        type="html"
      />

      <!-- Markdown -->
      <ContentViewer
        v-else-if="section.type === 'SECTION_TYPE_MARKDOWN'"
        :content="getContent(section)"
        type="markdown"
      />

      <!-- 标题 -->
      <component
        :is="`h${titleLevel(section)}`"
        v-else-if="section.type === 'SECTION_TYPE_TITLE'"
        class="font-bold text-foreground"
        :class="[titleSizeClass(titleLevel(section)), alignClass(section.config)]"
      >
        {{ getContent(section) }}
      </component>

      <!-- 图片 -->
      <figure v-else-if="section.type === 'SECTION_TYPE_IMAGE' && getContent(section)" class="overflow-hidden">
        <a
          v-if="section.config?.href"
          :href="section.config.href"
          target="_blank"
          rel="noopener noreferrer"
        >
          <UiImage :src="getContent(section)" :alt="section.config?.alt || ''" class="w-full rounded-xl object-cover" />
        </a>
        <UiImage
          v-else
          :src="getContent(section)"
          :alt="section.config?.alt || ''"
          class="w-full rounded-xl object-cover"
        />
        <figcaption v-if="section.config?.caption" class="mt-2 text-center text-sm text-muted-foreground">
          {{ section.config.caption }}
        </figcaption>
      </figure>

      <!-- 图集 -->
      <div
        v-else-if="section.type === 'SECTION_TYPE_GALLERY' && getImageList(section).length > 0"
        class="grid gap-4"
        :style="{ gridTemplateColumns: galleryColumns(section) }"
      >
        <figure v-for="(img, i) in getImageList(section)" :key="i" class="overflow-hidden rounded-lg">
          <UiImage :src="img.src" :alt="img.caption || ''" class="aspect-[4/3] w-full object-cover" />
          <figcaption v-if="img.caption" class="mt-1 text-center text-xs text-muted-foreground">
            {{ img.caption }}
          </figcaption>
        </figure>
      </div>

      <!-- 轮播：横向滚动吸附 -->
      <div
        v-else-if="section.type === 'SECTION_TYPE_CAROUSEL' && getImageList(section).length > 0"
        class="-mx-2 flex snap-x snap-mandatory gap-4 overflow-x-auto px-2 pb-2"
      >
        <figure
          v-for="(img, i) in getImageList(section)"
          :key="i"
          class="w-[80%] shrink-0 snap-center overflow-hidden rounded-xl max-md:w-[85%]"
        >
          <UiImage :src="img.src" :alt="img.caption || ''" class="aspect-video w-full object-cover" />
          <figcaption v-if="img.caption" class="mt-1 text-center text-xs text-muted-foreground">
            {{ img.caption }}
          </figcaption>
        </figure>
      </div>

      <!-- 视频 -->
      <div v-else-if="section.type === 'SECTION_TYPE_VIDEO' && getContent(section)" class="overflow-hidden rounded-xl border border-border">
        <iframe
          v-if="videoEmbed(getContent(section))"
          :src="videoEmbed(getContent(section)) || ''"
          title="video"
          class="aspect-video w-full"
          allowfullscreen
          loading="lazy"
        />
        <video
          v-else
          :src="getContent(section)"
          :poster="section.config?.poster || undefined"
          controls
          class="aspect-video w-full"
        />
      </div>

      <!-- 按钮 -->
      <div v-else-if="section.type === 'SECTION_TYPE_BUTTON' && getContent(section) && (section.config?.href || section.config?.url || section.config?.link)" :class="alignClass(section.config)">
        <UiButton
          :as="'a'"
          :href="section.config?.href || section.config?.url || section.config?.link"
          :target="isExternal(section.config?.href || section.config?.url || section.config?.link || '') ? '_blank' : undefined"
          :rel="isExternal(section.config?.href || section.config?.url || section.config?.link || '') ? 'noopener noreferrer' : undefined"
          :variant="section.config?.variant === 'outline' ? 'outline' : 'default'"
        >
          {{ getContent(section) }}
        </UiButton>
      </div>

      <!-- 分割线 -->
      <hr v-else-if="section.type === 'SECTION_TYPE_DIVIDER'" class="border-border" />

      <!-- 间距 -->
      <div v-else-if="section.type === 'SECTION_TYPE_SPACER'" :style="{ height: spacerHeight(section) }" aria-hidden="true" />

      <!-- 代码：借道 Viewer 的 Markdown 管线获得 Shiki 高亮 -->
      <ContentViewer
        v-else-if="section.type === 'SECTION_TYPE_CODE' && getContent(section)"
        :content="'```' + (section.config?.language || '') + '\n' + getContent(section) + '\n```'"
        type="markdown"
      />

      <!-- 自定义 HTML（Viewer 内 DOMPurify 消毒） -->
      <ContentViewer
        v-else-if="section.type === 'SECTION_TYPE_HTML'"
        :content="getContent(section)"
        type="html"
      />

      <!-- 表单：字段定义在 config.fields，标签在翻译 content 的 <字段名> 键 -->
      <form
        v-else-if="section.type === 'SECTION_TYPE_FORM'"
        class="space-y-4"
        @submit.prevent="submitForm(section)"
      >
        <div
          v-if="pickTranslation(section)?.content?.title || section.name"
          class="text-lg font-semibold text-foreground"
        >
          {{ pickTranslation(section)?.content?.title || section.name }}
        </div>
        <template v-if="formStateOf(section) !== 'done'">
          <div v-for="field in formFields(section)" :key="field.name" class="space-y-1">
            <label :for="`f-${section.id}-${field.name}`" class="text-sm font-medium text-foreground">
              {{ fieldLabel(section, field) }}
              <span v-if="field.required" class="text-destructive"> *</span>
            </label>
            <textarea
              v-if="field.type === 'textarea'"
              :id="`f-${section.id}-${field.name}`"
              v-model="formValuesOf(section)[field.name]"
              :required="field.required"
              rows="4"
              class="w-full rounded-lg border border-border bg-background px-3 py-2 text-sm text-foreground outline-none focus:border-primary"
            />
            <input
              v-else
              :id="`f-${section.id}-${field.name}`"
              v-model="formValuesOf(section)[field.name]"
              :type="field.type === 'email' ? 'email' : field.type === 'tel' ? 'tel' : 'text'"
              :required="field.required"
              class="w-full rounded-lg border border-border bg-background px-3 py-2 text-sm text-foreground outline-none focus:border-primary"
            />
          </div>
          <div v-if="formStateOf(section) === 'error'" class="text-sm text-destructive">
            {{ pickTranslation(section)?.content?.errorText || 'Submit failed, please retry' }}
          </div>
          <UiButton type="submit" :disabled="formStateOf(section) === 'submitting'">
            {{ pickTranslation(section)?.content?.submitText || 'Submit' }}
          </UiButton>
        </template>
        <div
          v-else
          class="rounded-xl border border-border bg-muted/40 p-6 text-center text-sm text-muted-foreground"
        >
          {{ pickTranslation(section)?.content?.successText || pickTranslation(section)?.content?.title }}
        </div>
      </form>

      <!-- 自定义/未知类型：按 Markdown 兜底渲染 value -->
      <ContentViewer
        v-else-if="getContent(section)"
        :content="getContent(section)"
        type="markdown"
      />
    </div>
  </div>
</template>
