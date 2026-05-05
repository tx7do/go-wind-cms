import type { App } from 'vue'
import {createHead} from "@unhead/vue";

/**
 * 注册全局方法
 */
export function setupGlobalMethods(app: App) {
  const head = createHead()
  app.use(head)
}
