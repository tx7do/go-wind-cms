<script setup lang="ts">
import {h, ref, onMounted} from 'vue';
import {useRouter} from 'vue-router';

import {useNavbarStore, useNavigationStore} from "@/stores";
import {XIcon} from '@/plugins/xicon';

import type {TopNavBarTabItem} from "./types";
import TopNavbarTab from './TopNavbarTab.vue';

const router = useRouter();
const navbarStore = useNavbarStore();
const navigationStore = useNavigationStore();

const leftTabList: TopNavBarTabItem[] = []
const rightTabList: TopNavBarTabItem[] = [];
const navigationItems = ref<any[]>([]);

/**
 * 隐藏所有的悬浮层
 */
function hideOverlay() {
  navbarStore.setActiveOverlay(null);
}

/**
 * 加载导航数据
 */
async function loadNavigation() {
  try {
    const res = await navigationStore.listNavigation(
      {page: 1, pageSize: 10},
      {location: 'header', isActive: true}
    );
    if (res.items && res.items.length > 0) {
      navigationItems.value = res.items[0].items || [];
    }
  } catch (error) {
    console.error('Load navigation failed:', error);
  }
}

/**
 * 处理导航点击
 */
function handleNavigate(item: any) {
  if (item.isOpenNewTab) {
    window.open(item.url, '_blank');
  } else {
    router.push(item.url);
  }
}

onMounted(() => {
  loadNavigation();
});
</script>

<template>
  <n-space justify="space-between" align="center" style="width: 100%; padding: 0 var(--space-4)">
    <n-space align="center" :size="4">
      <!-- 左侧导航菜单 -->
      <n-menu
        v-if="navigationItems.length > 0"
        mode="horizontal"
        :options="navigationItems.map(item => ({
          key: item.id,
          label: item.title,
          icon: item.icon ? () => h(XIcon, { name: `carbon:${item.icon}`, size: 18 }) : undefined,
          children: item.children && item.children.length > 0 ? item.children.map((child: any) => ({
            key: child.id,
            label: child.title,
            icon: child.icon ? () => h(XIcon, { name: `carbon:${child.icon}`, size: 16 }) : undefined,
          })) : undefined,
        }))"
        @update:value="(key: number) => {
          const findItem = (items: any[]): any => {
            for (const item of items) {
              if (item.id === key) return item;
              if (item.children) {
                const found = findItem(item.children);
                if (found) return found;
              }
            }
            return null;
          };
          const item = findItem(navigationItems);
          if (item) handleNavigate(item);
        }"
      />

      <!-- 原有的 Tabs -->
      <n-tabs
        v-model:value="navbarStore.activeOverlay"
        type="bar"
        trigger="hover"
        animated
        @mouseleave="hideOverlay"
      >
        <n-tab-pane
          v-for="item in leftTabList" :key="item.key" :name="item.key"
          :tab="() => h(TopNavbarTab, { data: item })"
        >
          <div class="overlay-mask">
            <div class="overlay">
              <component :is="item.component" class="stick-content"/>
            </div>
          </div>
        </n-tab-pane>
      </n-tabs>
    </n-space>

    <n-space align="center" :size="4">
      <n-tabs
        v-model:value="navbarStore.activeOverlay"
        type="bar"
        trigger="hover"
        animated
        @mouseleave="hideOverlay"
      >
        <n-tab-pane
          v-for="item in rightTabList" :key="item.key" :name="item.key"
          :tab="() => h(TopNavbarTab, { data: item })"
        >
          <div class="overlay-mask">
            <div class="overlay">
              <transition name="fade" mode="out-in">
                <component :is="item.component" class="stick-content"/>
              </transition>
            </div>
          </div>
        </n-tab-pane>
      </n-tabs>
    </n-space>
  </n-space>
</template>

<style scoped lang="less">
@import '@/styles/app.less';

.overlay {
  background-color: var(--color-surface);
  border-top: 1px solid var(--color-border);
  z-index: 9999;
  pointer-events: auto;
  min-height: 15vh;
  color: var(--color-text-primary);
}

.overlay-mask {
  position: fixed;
  left: 0;
  width: 100%;
  height: 100%;
  background-color: rgba(0, 0, 0, 0.3);
  backdrop-filter: blur(4px);
  z-index: 9998;
  pointer-events: none;
}

// Menu 样式
:deep(.n-menu) {
  background: transparent !important;
  color: var(--color-text-primary);
  border: none !important;
}

:deep(.n-menu--horizontal) {
  height: auto;
  padding: 0 !important;
}

:deep(.n-menu-item) {
  color: var(--color-text-primary) !important;
  padding: 10px 14px !important;
  border-radius: var(--radius-sm);
  transition: all 0.2s;
  font-size: 15px;
  height: auto;
  line-height: 1.4;
  min-height: 36px;
  display: flex;
  align-items: center;
}

:deep(.n-menu-item:hover) {
  color: var(--color-brand) !important;
  background: rgba(102, 126, 234, 0.08) !important;
}

:deep(.n-menu-item--selected) {
  color: var(--color-brand) !important;
  background: rgba(102, 126, 234, 0.12) !important;
}

:deep(.n-menu-item-content) {
  display: flex;
  align-items: center;
  gap: 8px;
}

:deep(.n-menu-item-content__icon) {
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

:deep(.n-menu-item-content__label) {
  flex: 1;
  white-space: nowrap;
}

// Submenu 样式
:deep(.n-submenu) {
  color: var(--color-text-primary);
  padding: 0 !important;
}

:deep(.n-submenu__children) {
  background: transparent;
}

:deep(.n-submenu-children) {
  background: var(--color-surface) !important;
  border: 1px solid var(--color-border) !important;
  border-radius: var(--radius-md);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
  margin-top: 4px;
  padding: 8px 0 !important;
}

:deep(.n-submenu-children .n-menu-item) {
  color: var(--color-text-primary) !important;
  padding: 10px 16px !important;
  font-size: 14px;
}

:deep(.n-submenu-children .n-menu-item:hover) {
  color: var(--color-brand) !important;
  background: rgba(102, 126, 234, 0.08) !important;
}

:deep(.n-submenu-children .n-menu-item--selected) {
  color: var(--color-brand) !important;
  background: rgba(102, 126, 234, 0.12) !important;
}

// Submenu 展开箭头样式
:deep(.n-submenu-title) {
  padding: 10px 14px !important;
  font-size: 15px;
  min-height: 36px;
}

// Tabs 样式
:deep(.n-tabs) {
  color: var(--color-text-primary);
}

:deep(.n-tabs--horizontal) {
  padding: 0 !important;
}

// Tab 项样式
:deep(.n-tabs-tab) {
  color: var(--color-text-primary) !important;
  padding: 10px 14px !important;
  font-size: 15px;
  min-height: 36px;
  display: flex;
  align-items: center;
}

:deep(.n-tabs-tab:hover) {
  color: var(--color-brand) !important;
  background: rgba(102, 126, 234, 0.08);
}

// 激活的 Tab
:deep(.n-tabs-tab--active) {
  color: var(--color-brand) !important;
}

// Tab 下划线
:deep(.n-tabs-bar) {
  background-color: var(--color-brand);
  height: 2px;
}

// Tab 内容
:deep(.n-tabs-tab__label) {
  color: inherit;
}

// Space 组件
:deep(.n-space) {
  color: var(--color-text-primary);
}

.default-enter-active {
  transition: all 0.3s ease-in;
}

.default-leave-active {
  transition: all 0.3s ease-out;
}

.default-enter,
.default-leave-to {
  max-height: 0;
}
</style>
