<script setup lang="ts">
import {apiClient} from '@/api/client'

/**
 * 媒体引用图片（即时图片变换接入）：
 *   - IMAGE 值为 media_asset:<id> 引用时优先经免鉴权变换端点渲染
 *     （只缩小不放大，w=1280、q=80）；
 *   - gif/svg 等不可变换源（415）或资产缺失（404）回退元数据 url，再败回占位图；
 *   - http(s) 直链原样渲染（无资产 id，无法参与变换）。
 */
const props = defineProps<{
  value: string
  alt: string
}>()

const config = useAppConfig()

// 渲染链：变换 URL → 元数据 url（回退）→ 占位图（再败终态）
const id = computed(() => mediaRefId(props.value))
const fallbackUrl = ref<string | null>(null)
const failed = ref(false)

const src = computed(() => {
  if (failed.value) return '/placeholder.svg'
  if (id.value && !fallbackUrl.value) {
    return mediaTransformUrl(config.apiURL, id.value)
  }
  return fallbackUrl.value ?? props.value
})

function onError() {
  if (id.value && fallbackUrl.value === null && !failed.value) {
    // 变换端点不可用：回退元数据 url。取元数据失败仅告警——
    // 这是 best-effort 兜底链的中间环，终态仍会落占位图。
    apiClient.mediaAssetService
      .Get({id: Number(id.value)})
      .then((asset) => {
        if (asset.url) fallbackUrl.value = asset.url
        else failed.value = true
      })
      .catch((err: unknown) => {
        console.error('[MediaRefImage] 元数据回退失败:', err)
        failed.value = true
      })
    return
  }
  failed.value = true
}
</script>

<template>
  <img
    v-if="isMediaImageValue(value)"
    :src="src"
    :alt="alt"
    loading="lazy"
    @error="onError"
  >
</template>
