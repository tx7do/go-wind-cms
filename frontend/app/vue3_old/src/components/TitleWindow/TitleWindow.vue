<script setup lang="ts">
import {ArrowForward} from "@vicons/ionicons5";

import {useNavbarStore} from '@/stores';

export interface Props {
  title?: string;
  path?: string;
  icon?: any;
  showArrow?: boolean;
  enableClick?: boolean;
}

const props = withDefaults(
  defineProps<Props>(),
  {
    title: '',
    path: '',
    icon: null,
    showArrow: true,
    enableClick: true,
  }
)

const navbarStore = useNavbarStore();

function handleHeaderClicked() {
  if (props.enableClick === false) {
    return;
  }

  console.log('handleHeaderClicked', props.path);
  navbarStore.navigateToAndHideOverlay(props.path);
}
</script>

<template>
  <n-card :bordered="false" class="container">
    <template #header>
      <n-space
        @click.stop="handleHeaderClicked"
        justify="start"
        align="center"
        :class="props.enableClick ? 'header-click' : ''"
      >
        <n-icon size="25" v-if="props.icon!=null">
          <component :is="props.icon"/>
        </n-icon>

        <n-text>{{ props.title }}</n-text>

        <n-icon class="arrow-forward" v-show="props.showArrow">
          <ArrowForward/>
        </n-icon>
      </n-space>
    </template>

    <slot name="content"></slot>

    <template #header-extra>
      <slot name="header-extra"></slot>
    </template>
  </n-card>
</template>

<style scoped lang="less">
@import "@/styles/app.less";

.container {
  background-color: @content-color;
}

.header-click {
  cursor: pointer;
}

.arrow-forward {
  width: 40px;
  height: 20px;
  background-color: @home-card-color;
}
</style>
