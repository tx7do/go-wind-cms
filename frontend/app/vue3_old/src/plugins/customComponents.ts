import type { App } from 'vue'

import XIcon from '@/plugins/xicon'

/**
 * 全局注册自定义组件
 * @param app
 */
export function setupCustomComponents(app: App) {
  app.use(XIcon)
}
