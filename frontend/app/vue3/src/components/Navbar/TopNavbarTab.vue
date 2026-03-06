<script setup lang="ts">

import {$t} from "@/locales";

import type {TopNavBarTabItem} from "./types";
import {navigateTo} from "@/router";
import {useNavbarStore} from "@/stores";

const navbarStore = useNavbarStore();

const props = defineProps<{ data: TopNavBarTabItem }>();

/**
 * 隐藏所有的悬浮层
 */
function hideOverlay() {
  navbarStore.setActiveOverlay(null);
}

/**
 * 点击Tab
 * @param path
 */
function handleClickTab(path: string) {
  console.log('handleClickTab', path);

  navigateTo(path);

  hideOverlay();
}
</script>

<template>
  <div class="tab-title" @click="() => handleClickTab(props.data.path)">
    <n-icon v-show="props.data.icon != null" class="tab-icon">
      <component :is="props.data.icon"/>
    </n-icon>
    <span class="tab-label">{{ $t(props.data.name) }}</span>
  </div>
</template>

<style scoped lang="less">
.tab-title {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  gap: 4px; // 进一步减小间距从6px到4px
  padding: 8px 16px;
  border-radius: var(--radius-sm);
  font-size: 14px;
  font-weight: 500;
  color: var(--color-text-primary);
  cursor: pointer;
  transition: all 0.2s ease;
  white-space: nowrap;
  user-select: none;

  &:hover {
    color: var(--color-brand);
    background-color: rgba(102, 126, 234, 0.05);
  }

  &:active {
    background-color: rgba(102, 126, 234, 0.1);
  }
}

.tab-icon {
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 20px; // 增大图标尺寸从18px到20px
}

.tab-label {
  white-space: nowrap;
  line-height: 1.5;
}
</style>
