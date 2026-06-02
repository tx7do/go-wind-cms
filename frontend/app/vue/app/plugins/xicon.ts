import { Icon } from '@iconify/vue'
import { defineNuxtPlugin } from '#app'
import { defineComponent, h, type PropType } from 'vue'

/**
 * 通用图标组件 — 对齐 React 版 XIcon 行为
 * - 默认 size=24（与 React 版一致）
 * - 添加 inline-flex 布局
 * - 同时支持 size / width+height 属性
 */
const XIconComponent = defineComponent({
    name: 'XIcon',
    props: {
        icon: { type: String, required: true },
        size: { type: [Number, String] as PropType<number | string>, default: 24 },
        width: { type: [Number, String] as PropType<number | string>, default: undefined },
        height: { type: [Number, String] as PropType<number | string>, default: undefined },
        color: { type: String, default: undefined },
    },
    setup(props, { attrs }) {
        return () => h(Icon, {
            icon: props.icon,
            width: props.width ?? props.size,
            height: props.height ?? props.size,
            color: props.color,
            class: 'inline-flex items-center justify-center',
            ...attrs,
        })
    },
})

export const XIcon = XIconComponent

export default defineNuxtPlugin((nuxtApp) => {
    nuxtApp.vueApp.component('XIcon', XIconComponent)
})
