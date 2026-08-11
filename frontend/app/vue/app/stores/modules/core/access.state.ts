import {acceptHMRUpdate, defineStore} from 'pinia'
import {encryptByAES, decryptByAES} from '@/utils/crypto'

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
 * 从 Nuxt runtime config 取得 AES key（与密码加密同源）。
 * 修复 AUD9-M5：access store 持久化时对 token 字段做 AES 加密，
 * 使 localStorage 落盘的值为密文而非明文。
 */
function getAesKey(): string {
    const config = useRuntimeConfig()
    const key = config.public.aesKey as string;
    if (!key) {
        throw new Error('AES_KEY is not set in runtime config');
    }
    return key;
}

/**
 * 自定义序列化器：落盘前对 token 字段加密，读盘后解密。
 * 仅 accessToken/refreshToken 的 value 字段加密，其余字段（accessCodes/Menus/Routes）明文。
 */
const tokenAwareSerializer = {
    serialize(data: any): string {
        const clone = JSON.parse(JSON.stringify(data));
        const key = getAesKey();
        if (clone?.accessToken?.value) {
            clone.accessToken.value = encryptByAES(clone.accessToken.value, key);
        }
        if (clone?.refreshToken?.value) {
            clone.refreshToken.value = encryptByAES(clone.refreshToken.value, key);
        }
        return JSON.stringify(clone);
    },
    deserialize(data: string): any {
        const parsed = JSON.parse(data);
        const key = getAesKey();
        if (parsed?.accessToken?.value) {
            try {
                parsed.accessToken.value = decryptByAES(parsed.accessToken.value, key);
            } catch {
                parsed.accessToken = null;
            }
        }
        if (parsed?.refreshToken?.value) {
            try {
                parsed.refreshToken.value = decryptByAES(parsed.refreshToken.value, key);
            } catch {
                parsed.refreshToken = null;
            }
        }
        return parsed;
    },
};

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
        // 持久化（v4 用 pick 替代 paths）
        // storage 由 nuxt.config.ts 的 piniaPluginPersistedstate.storage 统一设为 localStorage
        pick: ['accessToken', 'refreshToken', 'accessCodes', 'accessMenus', 'accessRoutes'],
        // AUD9-M5: 落盘 token 字段 AES 加密，读盘解密
        serializer: tokenAwareSerializer,
    },
})

// 解决热更新问题
const hot = import.meta.hot
if (hot)
    hot.accept(acceptHMRUpdate(useAccessStore, hot))
