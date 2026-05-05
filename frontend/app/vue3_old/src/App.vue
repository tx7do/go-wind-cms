<script setup lang="ts">
import {storeToRefs} from 'pinia'
import {computed, onMounted} from 'vue'
import {darkTheme, dateZhCN, zhCN} from 'naive-ui'
import {useHead} from "@unhead/vue";
import {useRoute} from 'vue-router'

import {useAppStore, useRouteCacheStore} from '@/stores'

import useAutoThemeSwitcher from '@/hooks/useAutoThemeSwitcher'

import {$t, i18n} from '@/locales'

import {ApplicationProvider} from "@/components/ApplicationProvider";
import {useLanguageChangeStore} from "@/stores/modules/app/language-change.state";

useHead({
  title: $t('app.title'),
  meta: [
    {
      name: 'description',
      content: import.meta.env.VITE_APP_DESCRIPTION,
    },
    {
      name: 'keywords',
      content: import.meta.env.VITE_APP_KEYWORDS,
    },
  ],
  link: [
    {
      rel: 'icon',
      type: 'image/png',
      href: '/logo.png',
    },
  ],
})

const route = useRoute()
const appStore = useAppStore();
const routeCacheStore = useRouteCacheStore();

const {mode} = storeToRefs(appStore)

const {initializeThemeSwitcher} = useAutoThemeSwitcher(appStore)

const keepAliveRouteNames = computed(() => {
  return routeCacheStore.routeCaches as string[]
})

// 判断是否显示全局布局（Header 和 Footer）
const showLayout = computed(() => {
  return !route.meta.hideLayout
})

// 某些页面（如登录/注册）强制使用亮色主题，忽略全局 dark mode
const forceLightTheme = computed(() => {
  return route.meta.forceLight === true
})

const language = computed(() => {
  return i18n.global.locale.value === 'zhCN' ? zhCN : null
})
const dateLanguage = computed(() => {
  return i18n.global.locale.value === 'zhCN' ? dateZhCN : null
})
const theme = computed(() => {
  if (forceLightTheme.value)
    return null

  return mode.value === 'dark' ? darkTheme : null
})

onMounted(() => {
  initializeThemeSwitcher()
  
  // 初始化语言切换监听器
  const languageChangeStore = useLanguageChangeStore()
  languageChangeStore.initializeLanguageWatcher()
  
  console.log('[App] Language change watcher initialized')
})
</script>

<template>
  <n-config-provider
    :locale="language"
    :theme="theme"
    :date-locale="dateLanguage"
  >
    <n-global-style/>
    <ApplicationProvider>
      <Header v-if="showLayout" class="header header-color"/>
      <div class="content" :class="{ 'no-header': !showLayout }">
        <router-view v-slot="{ Component, route }">
          <transition name="fade" mode="out-in">
            <keep-alive :include="keepAliveRouteNames">
              <component :is="Component" :key="route.name"/>
            </keep-alive>
          </transition>
        </router-view>
      </div>
      <Footer v-if="showLayout" class="footer-wrapper"/>
    </ApplicationProvider>
  </n-config-provider>

  <n-back-top :right="100"/>
</template>

<style scoped lang="less">
@import "@/styles/app.less";

.header {
  position: fixed;
  top: 0;
  width: 100%;
  z-index: 1000;
  height: var(--layout-header-height);
  background-color: @header-color;
}

.content {
  min-height: 100vh;
  padding-top: var(--layout-header-height);
  width: 100%;
  background: var(--color-bg);

  &.no-header {
    padding-top: 0;
  }
}


.footer-wrapper {
  width: 100%;
  background: var(--color-surface);
}
</style>
