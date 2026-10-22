import { type ClassValue, clsx } from "clsx"
import { extendTailwindMerge } from "tailwind-merge"

/**
 * tailwind-merge 默认只把 t-shirt 字号（text-base/text-lg）归入 font-size 组，
 * 其余 text-* 一律按 text-color 处理。本项目的字号被换成了
 * title / card-title / body / desc / tips（tailwind.config.js），于是
 * `text-white text-card-title` 变成同组互斥，只有最后一个存活——
 * 而 cva 的拼接顺序是 base → variant → size，主按钮的 text-white 永远排在 size 前，
 * 于是被静默丢弃，蓝底按钮退回继承的 --color-text-main（黑字压蓝底）。
 * 把这五个字号登记进 theme.text，它们才回到自己的组。
 */
const mergeTailwind = extendTailwindMerge({
    extend: {
        theme: {
            text: ['title', 'card-title', 'body', 'desc', 'tips'],
        },
    },
})

export function cn(...inputs: ClassValue[]) {
  return mergeTailwind(clsx(inputs))
}
