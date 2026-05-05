import type {App} from "vue";

import router from "./router";

export function setupRouter(app: App<Element>) {
  app.use(router)
}

export * from './router'
export * from './utils'
