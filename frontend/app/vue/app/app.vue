<script setup lang="ts">
import { Toaster } from '@/components/ui/sonner'

const isDev = import.meta.dev
const { t } = useI18n()
const config = useRuntimeConfig()

// @nuxtjs/i18n：把当前 locale 的 htmlAttrs（lang、dir）应用到 <html>。
// 阿拉伯语（ar）在 nuxt.config.ts 的 locale 声明里带了 dir:'rtl'，
// 这里会把 dir="rtl" 写入 <html>，使布局整体镜像；其余 locale 为 ltr。
// 非 strict-SEO 模式下需显式调用 useLocaleHead + useHead 才会注入。
// head 是 Ref，useHead 会自动响应 locale 变化（如切到 ar 时刷新 <html dir>）。
const head = useLocaleHead({ dir: true, lang: true, seo: false })
useHead({
  titleTemplate: (title) => title ? `${title} - ${config.public.appTitle}` : config.public.appTitle,
  htmlAttrs: head.value.htmlAttrs,
})
</script>

<template>
  <div>
    <NuxtLayout>
      <NuxtPage />
    </NuxtLayout>
    <Toaster />
    <ClientOnly>
      <VueQueryDevtools v-if="isDev" />
    </ClientOnly>
  </div>
</template>
