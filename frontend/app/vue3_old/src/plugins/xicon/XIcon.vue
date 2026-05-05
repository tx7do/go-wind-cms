<template>
  <IconifyIcon v-if="useIconifyDirect" :icon="name" :width="size" :height="size" :color="iconColor" />
  <Icon v-else class="icon" v-bind="$attrs" tag="span" :size="size" :color="iconColor">
    <IconComponent/>
  </Icon>
</template>

<script lang="ts" setup>
import {computed} from 'vue'
import type {Component} from 'vue'
import {Icon} from '@vicons/utils'
import {Icon as IconifyIcon} from '@iconify/vue'
import {Help} from '@vicons/carbon'
import icons from './icons'
import {preferences} from "@/preferences";

const props = withDefaults(
  defineProps<{
    name: string
    size?: number
    color?: string
  }>(),
  {
    size: 25,
    color: ''
  }
)

const iconColor = computed<string>(() => (props.color ? props.color : preferences.theme.colorPrimary))

// 判断是否使用 Iconify 直接加载
const useIconifyDirect = computed<boolean>(() => {
  // 如果图标名称包含冒号(如 carbon:folder),则使用 Iconify 直接加载
  return props.name.includes(':')
})

// 兼容原有的图标查表方式
const IconComponent = computed<Component>(() => {
  console.log('icon name', props.name)
  return icons[props.name] ? icons[props.name] : Help
})
</script>

<style lang="less" scoped>
.icon {
@apply h-full inline-flex justify-center items-center;
}
</style>
