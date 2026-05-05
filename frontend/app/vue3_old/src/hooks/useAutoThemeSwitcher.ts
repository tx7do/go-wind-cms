import type { useAppStore } from '@/stores'

export default function useAutoThemeSwitcher(appStore: ReturnType<typeof useAppStore>) {
  const handleAttributeChange = () => {
    const rootElement = document.documentElement
    if (rootElement.classList.contains('dark'))
      appStore.switchMode('dark')
    else
      appStore.switchMode('light')
  }

  const observerOptions = {
    attributes: true,
    attributeFilter: ['class'],
  }

  const observer = new MutationObserver(handleAttributeChange)

  const targetElement = document.querySelector('html')

  const initializeThemeSwitcher = () => {
    observer.observe(targetElement, observerOptions)
  }

  return { initializeThemeSwitcher }
}
