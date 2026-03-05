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
  <n-space justify="space-between" align="center" style="width: 100%; padding: 0 var(--space-4); height: 56px; overflow: hidden;">
    <n-space align="center" :size="4" style="height: 100%; overflow: hidden; display: flex; align-items: center;">
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

    <n-space align="center" :size="4" style="height: 100%; overflow: hidden; display: flex; align-items: center;">
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
  padding: 0 !important;
  height: 100%;
  overflow: hidden;
  display: flex;
  align-items: center;
}

:deep(.n-menu--horizontal) {
  height: 100%;
  padding: 0 !important;
  overflow: hidden;
  display: flex;
  align-items: center;
}

:deep(.n-menu-item) {
  color: var(--color-text-primary) !important;
  padding: 8px 12px !important;
  border-radius: var(--radius-sm);
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  font-size: 15px;
  font-weight: 500;
  height: 100%;
  line-height: 1.2;
  display: flex;
  align-items: center;
}

:deep(.n-menu-item:hover) {
  color: var(--color-brand) !important;
  background: rgba(102, 126, 234, 0.1) !important;
  transform: translateX(2px);
}

:deep(.n-menu-item--selected) {
  color: var(--color-brand) !important;
  background: rgba(102, 126, 234, 0.15) !important;
  font-weight: 600;
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
  transition: all 0.3s;
}

:deep(.n-menu-item:hover .n-menu-item-content__icon) {
  transform: scale(1.1);
  color: var(--color-brand);
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
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.12);
  margin-top: 8px;
  padding: 8px 0 !important;
  animation: slideDown 0.3s cubic-bezier(0.4, 0, 0.2, 1);
}

:deep(.n-submenu-children .n-menu-item) {
  color: var(--color-text-primary) !important;
  padding: 10px 16px !important;
  font-size: 14px;
  font-weight: 500;
  transition: all 0.3s;
}

:deep(.n-submenu-children .n-menu-item:hover) {
  color: var(--color-brand) !important;
  background: rgba(102, 126, 234, 0.1) !important;
  padding-left: 20px !important;
}

:deep(.n-submenu-children .n-menu-item--selected) {
  color: var(--color-brand) !important;
  background: rgba(102, 126, 234, 0.15) !important;
}

// Submenu 展开箭头样式
:deep(.n-submenu-title) {
  padding: 8px 12px !important;
  font-size: 15px;
  font-weight: 500;
  min-height: 32px;
  transition: all 0.3s;
  border-radius: var(--radius-sm);
}

:deep(.n-submenu--selected > .n-submenu-title) {
  color: var(--color-brand) !important;
  background: rgba(102, 126, 234, 0.1) !important;
}

// Tabs 样式
:deep(.n-tabs) {
  color: var(--color-text-primary);
  padding: 0 !important;
  height: 100%;
  overflow: hidden;
  display: flex;
  align-items: center;
}

:deep(.n-tabs--horizontal) {
  padding: 0 !important;
  height: 100%;
  overflow: hidden;
  display: flex;
  align-items: center;
}

// Tab 项样式
:deep(.n-tabs-tab) {
  color: var(--color-text-primary) !important;
  padding: 8px 12px !important;
  font-size: 15px;
  font-weight: 500;
  height: 100%;
  display: flex;
  align-items: center;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  border-radius: var(--radius-sm) var(--radius-sm) 0 0;
}

:deep(.n-tabs-tab:hover) {
  color: var(--color-brand) !important;
  background: rgba(102, 126, 234, 0.08);
}

// 激活的 Tab
:deep(.n-tabs-tab--active) {
  color: var(--color-brand) !important;
  background: rgba(102, 126, 234, 0.1) !important;
  font-weight: 600;
}

// Tab 下划线
:deep(.n-tabs-bar) {
  background: linear-gradient(90deg, var(--color-brand) 0%, var(--color-brand) 100%);
  height: 3px;
}

// Tab 内容
:deep(.n-tabs-tab__label) {
  color: inherit;
  transition: all 0.3s;
}

// Space 组件
:deep(.n-space) {
  color: var(--color-text-primary);
}

// 动画
@keyframes slideDown {
  from {
    opacity: 0;
    transform: translateY(-10px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

.default-enter-active {
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
}

.default-leave-active {
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
}

.default-enter,
.default-leave-to {
  max-height: 0;
  opacity: 0;
}

// Responsive
@media (max-width: 768px) {
  :deep(.n-menu-item),
  :deep(.n-submenu-title),
  :deep(.n-tabs-tab) {
    padding: 8px 10px !important;
    font-size: 13px;
    min-height: 32px;
  }

  :deep(.n-submenu-children .n-menu-item) {
    padding: 8px 12px !important;
    font-size: 12px;
  }

  :deep(.n-menu-item-content) {
    gap: 4px;
  }

  :deep(.n-tabs-bar) {
    height: 2px;
  }
}
</style>
