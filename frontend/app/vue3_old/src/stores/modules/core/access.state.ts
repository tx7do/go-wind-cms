import {acceptHMRUpdate, defineStore} from 'pinia'

/**
 * 令牌载荷
 */
interface TokenPayload {
  /**
   * 令牌值
   */
  value: string;
  /**
   * 令牌过期时间
   */
  expiresAt?: number;
}

interface AccessState {
  /**
   * 权限码
   */
  accessCodes: string[];
  /**
   * 权限菜单
   */
  accessMenus?: string[];
  /**
   * 权限路由
   */
  accessRoutes?: string[];
  /**
   * 登录 accessToken
   */
  accessToken: TokenPayload | null;
  /**
   * 登录 refreshToken
   */
  refreshToken: TokenPayload | null;
  /**
   * 是否已经检查过权限
   */
  isAccessChecked: boolean;
  /**
   * 登录是否过期
   */
  loginExpired: boolean;
}

/**
 * @zh_CN 访问权限相关
 */
export const useAccessStore = defineStore('access', {
  state: (): AccessState => ({
    accessCodes: [],
    accessMenus: [],
    accessRoutes: [],
    accessToken: null,
    refreshToken: null,
    isAccessChecked: false,
    loginExpired: false,
  }),
  actions: {
    setAccessCodes(codes: string[]) {
      this.accessCodes = codes
    },
    setAccessMenus(menus: string[]) {
      this.accessMenus = menus
    },
    setAccessRoutes(routes: string[]) {
      this.accessRoutes = routes
    },
    setAccessToken(token: string, expiresAt?: number) {
      this.accessToken = {value: token, expiresAt}
    },
    setRefreshToken(token: string, expiresAt?: number) {
      this.refreshToken = {value: token, expiresAt}
    },
    setLoginExpired(loginExpired: boolean) {
      this.loginExpired = loginExpired
    },
    setIsAccessChecked(isChecked: boolean) {
      this.isAccessChecked = isChecked
    },
  },
  persist: {
    // 持久化
    paths: ['accessToken', 'refreshToken', 'accessCodes', 'accessMenus', 'accessRoutes'],
  },
})

// 解决热更新问题
const hot = import.meta.hot
if (hot)
  hot.accept(acceptHMRUpdate(useAccessStore, hot))
