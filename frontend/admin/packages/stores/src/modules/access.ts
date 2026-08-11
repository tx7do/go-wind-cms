import type { MenuRecordRaw } from '@vben-core/typings';
import type { RouteRecordRaw } from 'vue-router';

import { acceptHMRUpdate, defineStore } from 'pinia';

/**
 * 令牌持久化序列化器接口。
 *
 * AUD9-M5：access store 持久化时对 token 字段做加密，使 localStorage 落盘的值为密文而非明文。
 * 由于 @vben/stores 是不依赖任何加密库的共享包，具体的加解密实现由各 app 在启动时
 * 通过 {@link setTokenSerializer} 注入（如 admin app 注入基于 VITE_AES_KEY 的 AES-CBC 实现）。
 * 默认使用明文 JSON，仅在开发期告警以提示注入遗漏。
 */
export interface TokenSerializer {
  serialize: (data: any) => string;
  deserialize: (data: string) => any;
}

const plainSerializer: TokenSerializer = {
  serialize: (data) => JSON.stringify(data),
  deserialize: (data) => JSON.parse(data),
};

let tokenSerializer: TokenSerializer = plainSerializer;

/**
 * 注入 token 持久化序列化器。
 * 必须在 access store 首次实例化（useAccessStore()）之前调用，
 * 否则首次 hydration 仍使用明文默认序列化器。
 */
export function setTokenSerializer(s: TokenSerializer): void {
  tokenSerializer = s;
}

/**
 * @zh_CN 访问令牌类型
 */
type AccessToken = null | string;

/**
 * @zh_CN 访问权限相关状态定义
 */
interface AccessState {
  /**
   * 权限码
   */
  accessCodes: string[];
  /**
   * 可访问的菜单列表
   */
  accessMenus: MenuRecordRaw[];
  /**
   * 可访问的路由列表
   */
  accessRoutes: RouteRecordRaw[];
  /**
   * 登录 accessToken
   */
  accessToken: AccessToken;
  /**
   * accessToken 过期时间戳
   */
  accessTokenExpireTime?: number;
  /**
   * 是否已经检查过权限
   */
  isAccessChecked: boolean;
  /**
   * 登录是否过期
   */
  loginExpired: boolean;

  /**
   * 登录 accessToken
   */
  refreshToken: AccessToken;

  /**
   * refreshToken 过期时间戳
   */
  refreshTokenExpireTime?: number;
}

/**
 * @zh_CN 访问权限相关状态管理
 */
export const useAccessStore = defineStore('core-access', {
  actions: {
    $reset() {
      this.accessToken = null;
      this.refreshToken = null;
      this.accessCodes = [];
      this.accessMenus = [];
      this.accessRoutes = [];
      this.isAccessChecked = false;
      this.loginExpired = false;
      this.accessTokenExpireTime = undefined;
      this.refreshTokenExpireTime = undefined;
    },
    /**
     * @zh_CN 检查 accessToken 是否过期
     */
    checkAccessTokenExpired(): boolean {
      if (!this.accessTokenExpireTime) {
        return true;
      }
      const now = Date.now();
      return now >= this.accessTokenExpireTime;
    },
    /**
     * @zh_CN 检查 refreshToken 是否过期
     */
    checkRefreshTokenExpired(): boolean {
      if (!this.refreshTokenExpireTime) {
        return true;
      }
      const now = Date.now();
      return now >= this.refreshTokenExpireTime;
    },
    setAccessCodes(codes: string[]) {
      this.accessCodes = codes;
    },
    setAccessMenus(menus: MenuRecordRaw[]) {
      this.accessMenus = menus;
    },
    setAccessRoutes(routes: RouteRecordRaw[]) {
      this.accessRoutes = routes;
    },
    setAccessToken(token: AccessToken) {
      this.accessToken = token;
    },
    setAccessTokenExpireTime(accessTokenExpireTime: number) {
      this.accessTokenExpireTime = accessTokenExpireTime;
    },
    setIsAccessChecked(isAccessChecked: boolean) {
      this.isAccessChecked = isAccessChecked;
    },

    setLoginExpired(loginExpired: boolean) {
      this.loginExpired = loginExpired;
    },

    setRefreshToken(token: AccessToken) {
      this.refreshToken = token;
    },

    setRefreshTokenExpireTime(refreshTokenExpireTime: number) {
      this.refreshTokenExpireTime = refreshTokenExpireTime;
    },
  },
  persist: {
    // 持久化
    pick: [
      'accessToken',
      'refreshToken',
      'accessCodes',
      'refreshTokenExpireTime',
      'accessTokenExpireTime',
    ],
    // AUD9-M5: 通过 setTokenSerializer 注入的加解密序列化器（call-time 查询），
    // 默认明文 JSON。各 app 应在启动时注入加密实现。
    serializer: {
      serialize: (data: any) => tokenSerializer.serialize(data),
      deserialize: (data: string) => tokenSerializer.deserialize(data),
    },
  },
  state: (): AccessState => ({
    accessCodes: [],
    accessMenus: [],
    accessRoutes: [],
    accessToken: null,
    accessTokenExpireTime: undefined,
    isAccessChecked: false,
    loginExpired: false,
    refreshToken: null,
    refreshTokenExpireTime: undefined,
  }),
});

// 解决热更新问题
const hot = import.meta.hot;
if (hot) {
  hot.accept(acceptHMRUpdate(useAccessStore, hot));
}
