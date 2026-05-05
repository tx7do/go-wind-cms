import * as NaiveUI from 'naive-ui'
import { computed } from 'vue'
import {preferences} from "@/preferences";
import {lighten} from "@/utils/utils";

/**
 * 挂载 Naive-ui 脱离上下文的 API
 * 如果你想在 setup 外使用 useDialog、useMessage、useNotification、useLoadingBar，可以通过 createDiscreteApi 来构建对应的 API。
 * https://www.naiveui.com/zh-CN/dark/components/discrete
 */

export function setupNaiveDiscreteApi() {
  const configProviderPropsRef = computed(() => ({
    theme: preferences.theme.mode === 'dark' ? NaiveUI.darkTheme : undefined,
    themeOverrides: {
      common: {
        primaryColor: preferences.theme.colorPrimary,
        primaryColorHover: lighten(preferences.theme.colorPrimary, 6),
        primaryColorPressed: lighten(preferences.theme.colorPrimary, 6)
      },
      LoadingBar: {
        colorLoading: preferences.theme.colorPrimary
      }
    }
  }))

  const { message, dialog, notification, loadingBar } = NaiveUI.createDiscreteApi(
    ['message', 'dialog', 'notification', 'loadingBar'],
    {
      configProviderProps: configProviderPropsRef
    }
  )

  window['$message'] = message
  window['$dialog'] = dialog
  window['$notification'] = notification
  window['$loading'] = loadingBar
}
