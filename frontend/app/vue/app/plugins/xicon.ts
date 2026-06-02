import { Icon } from '@iconify/vue'
import { defineNuxtPlugin } from '#app'

export interface XIconProps {
    name: string
    size?: number | string
    color?: string
    className?: string
    style?: Record<string, string>
}

/**
 * 通用图标组件
 * 使用 Iconify 实现，支持 carbon、mdi 等所有 Iconify 图标库
 */
export const XIcon = Icon

export default defineNuxtPlugin((nuxtApp) => {
    nuxtApp.vueApp.component('XIcon', Icon)
})
