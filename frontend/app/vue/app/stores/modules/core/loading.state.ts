import { defineStore } from 'pinia'
import { ref } from 'vue'

export const useLoadingStore = defineStore('loading', () => {
  const isLoading = ref(false)
  const error = ref<boolean | null>(null)

  const start = () => {
    isLoading.value = true
    error.value = null
  }

  const finish = () => {
    isLoading.value = false
  }

  const error_ = () => {
    isLoading.value = false
    error.value = true
  }

  return {
    isLoading,
    error,
    start,
    finish,
    error_,
  }
})
