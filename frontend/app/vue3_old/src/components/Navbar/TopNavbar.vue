<script setup lang="ts">
import {h, ref, onMounted} from 'vue';
import {useRouter} from 'vue-router';

import {useNavbarStore, useNavigationStore} from '@/stores';
import {XIcon} from "@/plugins/xicon";
import {useLanguageChangeEffect} from '@/hooks/useLanguageChangeEffect';
import type {
  siteservicev1_NavigationItem
} from "@/api/generated/app/service/v1";

import type {TopNavBarTabItem} from "./types";
import TopNavbarTab from './TopNavbarTab.vue';

const router = useRouter();
const navbarStore = useNavbarStore();
const navigationStore = useNavigationStore();

const leftTabList: TopNavBarTabItem[] = []
const rightTabList: TopNavBarTabItem[] = [];
const navigationItems = ref<siteservicev1_NavigationItem[]>([]);

function getMenuOptions(items: siteservicev1_NavigationItem[]) {
  return items.map(item => ({
    key: item.id,
    label: item.title,
    icon: item.icon ? () => h(XIcon, {name: `carbon:${item.icon}`, size: 22}) : undefined,
    children: item.children && item.children.length > 0
      ? item.children.map((child: siteservicev1_NavigationItem) => ({
        key: child.id,
        label: child.title,
        icon: child.icon ? () => h(XIcon, {name: `carbon:${child.icon}`, size: 18}) : undefined,
      }))
      : undefined,
  }));
}

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
      {location: 'HEADER', isActive: true}
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
function handleNavigate(item: siteservicev1_NavigationItem) {
  if (item.isOpenNewTab) {
    window.open(item.url, '_blank');
  } else {
    // 使用 router.push 时保留 history state
    router.push({
      path: item.url,
      // 保留必要的 state 信息
      state: {...history.state, from: 'navbar'}
    });
  }
}

onMounted(() => {
  loadNavigation();
});

// 监听语言切换，自动重新加载导航数据
useLanguageChangeEffect(loadNavigation, {
  immediate: false,      // 是否立即执行一次
  autoCleanup: true,    // 是否自动清理
});
</script>

<template>
  <div class="navbar-wrapper">
    <div class="navbar-left">
      <!-- 左侧导航菜单 -->
      <n-menu
        v-if="navigationItems.length > 0"
        mode="horizontal"
        :options="getMenuOptions(navigationItems)"
        @update:value="(key: number) => {
            const item = navigationStore.findNavItem(navigationItems, key);
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
    </div>

    <div class="navbar-right">
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
    </div>
  </div>
</template>

<style scoped lang="less">
@import '@/styles/app.less';

.navbar-wrapper {
  display: flex;
  align-items: center;
  justify-content: space-between;
  width: 100%;
  padding: 0 var(--space-4);
  height: 65px;
  overflow: hidden;
  gap: 4px;
}

.navbar-left,
.navbar-right {
  display: flex;
  align-items: center;
  height: 100%;
  min-width: 0;
}

.navbar-left {
  flex: 1;
  gap: 4px;
  overflow-x: auto;
  overflow-y: hidden;

  // 隐藏滚动条
  scrollbar-width: none;

  &::-webkit-scrollbar {
    display: none;
  }
}

.navbar-right {
  flex-shrink: 0;
  gap: 4px;
}

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
  overflow: visible !important;
  display: flex !important;
  align-items: center;
  gap: 0 !important;
}

:deep(.n-menu--horizontal) {
  height: 100%;
  padding: 0 !important;
  overflow: visible !important;
  display: flex !important;
  align-items: center;
  gap: 0 !important;
}

:deep(.n-menu-item) {
  color: var(--color-text-primary) !important;
  padding: 8px 14px !important;
  border-radius: var(--radius-sm);
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  font-size: 15px;
  font-weight: 500;
  height: 100%;
  line-height: 1.2;
  display: flex;
  align-items: center;
  white-space: nowrap;
  flex-shrink: 0;
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
  gap: 4px;
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
  overflow: visible !important;
  display: flex !important;
  align-items: center;
  flex-shrink: 0;
}

:deep(.n-tabs--horizontal) {
  padding: 0 !important;
  height: 100%;
  overflow: visible !important;
  display: flex !important;
  align-items: center;
  flex-shrink: 0;
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
  white-space: nowrap;
  flex-shrink: 0;
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

@media (max-width: 1024px) {
  .navbar-wrapper {
    padding: 0 var(--space-3);
    height: 60px;
  }

  :deep(.n-menu-item),
  :deep(.n-submenu-title),
  :deep(.n-tabs-tab) {
    padding: 8px 12px !important;
    font-size: 13px;
    min-height: 32px;
    flex-shrink: 0;
  }

  :deep(.n-submenu-children .n-menu-item) {
    padding: 8px 12px !important;
    font-size: 12px;
  }

  :deep(.n-menu-item-content) {
    gap: 3px;
  }

  :deep(.n-tabs-bar) {
    height: 2px;
  }
}

@media (max-width: 768px) {
  .navbar-wrapper {
    padding: 0 var(--space-2);
    height: 50px;
    overflow-x: auto;
    overflow-y: hidden;
    gap: 4px;
    -webkit-overflow-scrolling: touch;
    scroll-behavior: smooth;

    // 显示细滚动条
    scrollbar-width: thin;
    scrollbar-color: rgba(102, 126, 234, 0.2) transparent;

    &::-webkit-scrollbar {
      height: 3px;
      display: block;
    }

    &::-webkit-scrollbar-track {
      background: transparent;
    }

    &::-webkit-scrollbar-thumb {
      background: rgba(102, 126, 234, 0.2);
      border-radius: 10px;
    }

    &::-webkit-scrollbar-thumb:hover {
      background: rgba(102, 126, 234, 0.4);
    }
  }

  .navbar-left {
    gap: 4px;
    flex: 1;
    overflow-x: auto;
    overflow-y: hidden;
    -webkit-overflow-scrolling: touch;

    scrollbar-width: none;

    &::-webkit-scrollbar {
      display: none;
    }
  }

  .navbar-right {
    gap: 4px;
    flex-shrink: 0;
  }

  :deep(.n-menu-item),
  :deep(.n-submenu-title),
  :deep(.n-tabs-tab) {
    padding: 8px 12px !important;
    font-size: 13px;
    min-height: 32px;
    height: 32px;
    flex-shrink: 0;
    white-space: nowrap;
  }

  :deep(.n-menu-item-content) {
    gap: 3px;
  }

  :deep(.n-menu-item-content__icon) {
    width: 18px;
    height: 18px;
    font-size: 18px;
    flex-shrink: 0;
  }

  :deep(.n-menu-item-content__label) {
    white-space: nowrap;
    font-size: 13px;
    flex-shrink: 0;
    font-weight: 500;
  }

  :deep(.n-submenu-children .n-menu-item) {
    padding: 8px 12px !important;
    font-size: 13px;
    min-height: 32px;
  }

  :deep(.n-tabs-bar) {
    height: 2px;
  }
}

@media (max-width: 480px) {
  .navbar-wrapper {
    padding: 0 var(--space-2);
    height: 48px;
    gap: 4px;
    overflow-x: auto;
    overflow-y: hidden;
    -webkit-overflow-scrolling: touch;
    scroll-behavior: smooth;

    // 显示滚动条，让用户知道可以滚动
    scrollbar-width: thin;
    scrollbar-color: rgba(102, 126, 234, 0.3) transparent;

    &::-webkit-scrollbar {
      height: 3px;
      display: block;
    }

    &::-webkit-scrollbar-track {
      background: transparent;
    }

    &::-webkit-scrollbar-thumb {
      background: rgba(102, 126, 234, 0.3);
      border-radius: 10px;
    }

    &::-webkit-scrollbar-thumb:hover {
      background: rgba(102, 126, 234, 0.5);
    }

    // 添加渐变指示器提示可以滚动
    position: relative;

    &::after {
      content: '';
      position: absolute;
      top: 0;
      right: 0;
      bottom: 3px;
      width: 30px;
      background: linear-gradient(to left, var(--header-bg) 0%, transparent 100%);
      pointer-events: none;
      z-index: 1;
    }
  }

  .navbar-left {
    gap: 4px;
    flex-shrink: 0;
    min-width: min-content;
  }

  .navbar-right {
    display: none;
  }

  :deep(.n-menu-item),
  :deep(.n-submenu-title),
  :deep(.n-tabs-tab) {
    padding: 8px 12px !important;
    font-size: 13px;
    min-height: 36px;
    height: 36px;
    flex-shrink: 0;
    border-radius: 18px !important;
    white-space: nowrap;
  }

  :deep(.n-menu-item-content) {
    gap: 3px;
  }

  :deep(.n-menu-item-content__label) {
    display: inline-block;
    font-size: 13px;
    font-weight: 500;
    white-space: nowrap;
  }

  :deep(.n-menu-item-content__icon) {
    font-size: 18px;
    flex-shrink: 0;
  }

  :deep(.n-submenu-children .n-menu-item) {
    padding: 8px 12px !important;
    font-size: 13px;
    min-height: 36px;
  }

  :deep(.n-submenu-children .n-menu-item-content__label) {
    display: inline-block !important;
  }

  :deep(.n-tabs-bar) {
    height: 2px;
  }
}
</style>
