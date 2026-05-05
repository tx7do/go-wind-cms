import {defineStore} from 'pinia'
import {ref} from 'vue'
import type {RouteRecordName} from 'vue-router'
import type {EnhancedRouteLocation} from '@/router/types'

export const useRouteCacheStore = defineStore('route-cache', () => {
  const routeCaches = ref<RouteRecordName[]>([])

  const addRoute = (route: EnhancedRouteLocation) => {
    if (routeCaches.value.includes(route.name))
      return

    if (route?.meta?.keepAlive)
      routeCaches.value.push(route.name)
  }

  function $reset() {
  }

  return {
    routeCaches,
    addRoute,
    $reset,
  }
})

