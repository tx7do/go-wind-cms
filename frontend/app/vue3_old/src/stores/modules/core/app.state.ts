import {defineStore} from 'pinia'
import {ref} from 'vue'

export interface AppState {
  switchMode: (val: string) => void
}

const prefersDark
  = window.matchMedia
  && window.matchMedia('(prefers-color-scheme: dark)').matches

export const useAppStore = defineStore('app', () => {
  const theme = prefersDark ? 'dark' : 'light'
  const mode = ref(theme)

  const switchMode = (val: string) => {
    mode.value = val
  }

  function $reset() {
  }

  return {
    mode,
    switchMode,
    $reset
  }
}, {
  persist: true,
})

