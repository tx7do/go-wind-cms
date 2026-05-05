import type {RouteLocationNormalized} from 'vue-router'

export interface EnhancedRouteLocation extends RouteLocationNormalized {
  meta: {
    level?: number | unknown
    name?: string
    keepAlive?: boolean
  }
}

export enum RouteMode {
  HISTORY = 'history',
  HASH = 'hash'
}
