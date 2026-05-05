/// <reference types="vite/client" />

// 引入naive对应的定义类型
import type { DialogApiInjection } from 'naive-ui/lib/dialog/src/DialogProvider'
import type { MessageApiInjection } from 'naive-ui/lib/message/src/MessageProvider'

declare module '*.vue' {
  import type { DefineComponent } from 'vue'

  // eslint-disable-next-line ts/ban-types
  const component: DefineComponent<{}, {}, any>
  export default component
}

declare global {
  interface Window {
    $message: MessageApiInjection
    $dialog: DialogApiInjection
  }
}
