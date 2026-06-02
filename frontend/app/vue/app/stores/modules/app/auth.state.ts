import {ref} from 'vue'
import {defineStore} from 'pinia'
import CryptoJS from 'crypto-js';

import {DEFAULT_HOME_PATH, LOGIN_PATH} from '@/constants'
import type {Recordable} from '@/typings/helper'
import type {UserInfo} from "@/typings";

import {useAccessStore, useUserStore} from '@/stores'
import {resetAllStores} from "@/stores/setup";

// API 调用从 composables 层引入
import {login as apiLogin, logout as apiLogout, refreshToken as apiRefreshToken} from '@/api/composables/auth';
import {getMe} from '@/api/composables/user-profile';
import { toast } from 'vue-sonner';

export const useAuthStore = defineStore('auth', () => {
    const accessStore = useAccessStore()
    const userStore = useUserStore()
    const { t } = useI18n()

    const loginLoading = ref(false)

    /**
     * 异步处理登录操作
     * Asynchronously handle the login process
     * @param params 登录表单数据
     * @param onSuccess
     */
    async function authLogin(
        params: Recordable<any>,
        onSuccess?: () => Promise<void> | void,
    ) {
        // 异步处理用户登录操作并获取 accessToken
        let userInfo: null | UserInfo = null
        try {
            loginLoading.value = true

            const {access_token} = await apiLogin({
                username: params.username,
                email: params.email,
                mobile: params.mobile,
                password: encryptPassword(params.password),
                grant_type: 'password',
            })

            // 如果成功获取到 accessToken
            if (access_token) {
                accessStore.setAccessToken(access_token)

                // 获取用户信息并存储到 accessStore 中
                userInfo = await fetchUserInfo()

                userStore.setUserInfo(userInfo)
                // accessStore.setAccessCodes(accessCodes);

                if (accessStore.loginExpired) {
                    accessStore.setLoginExpired(false)
                } else {
                    onSuccess
                        ? await onSuccess?.()
                        : await navigateTo(userInfo.homePath || DEFAULT_HOME_PATH)
                }

                if (userInfo?.realname) {
                    toast.success(t('authentication.login.login_success'))
                }
            }
        } finally {
            loginLoading.value = false
        }

        return {
            userInfo,
        }
    }

    /**
     * 用户登出
     * @param redirect
     */
    async function logout(redirect: boolean = true) {
        try {
            await apiLogout()
        } catch {
            // 不做任何处理
        } finally {
            resetAllStores()
            accessStore.setLoginExpired(false)

            // 回登录页带上当前路由地址
            const currentPath = typeof window !== 'undefined' ? window.location.pathname : '/'
            await navigateTo({
                path: LOGIN_PATH,
                query: redirect
                    ? {
                        redirect: encodeURIComponent(currentPath),
                    }
                    : {},
            })
        }
    }

    /**
     * 重新认证
     */
    async function reauthenticate() {
        console.warn('Access token or refresh token is invalid or expired. ')

        accessStore.setAccessToken('')
        if (
            accessStore.isAccessChecked
        )
            accessStore.setLoginExpired(true)
        else {
            // await logout()
        }
    }

    /**
     * 刷新访问令牌
     */
    async function refreshToken() {
        const refreshTokenValue = typeof accessStore.refreshToken === 'string'
            ? accessStore.refreshToken
            : accessStore.refreshToken?.value ?? '';
        const resp = await apiRefreshToken(refreshTokenValue)
        const newToken = resp.access_token

        accessStore.setAccessToken(newToken || '')

        return newToken
    }

    /**
     * 拉取用户信息
     */
    async function fetchUserInfo() {
        return (await getMe()) as unknown as UserInfo;
    }

    /**
     * 加密数据
     * @param data 待加密数据
     * @param key 密钥
     * @param iv 初始向量
     */
    function encryptData(data: string, key: string, iv: string): string {
        const keyHex = CryptoJS.enc.Utf8.parse(key);
        const ivHex = CryptoJS.enc.Utf8.parse(iv);
        const encrypted = CryptoJS.AES.encrypt(data, keyHex, {
            iv: ivHex,
            mode: CryptoJS.mode.CBC,
            padding: CryptoJS.pad.Pkcs7,
        });
        return encrypted.toString();
    }

    /**
     * 加密密码
     * @param password 明文密码
     */
    function encryptPassword(password: string): string {
        const config = useRuntimeConfig()
        const key = config.public.aesKey as string;
        if (!key) {
            throw new Error('AES_KEY is not set in runtime config');
        }
        return encryptData(password, key, key);
    }

    function $reset() {
        loginLoading.value = false
    }

    return {
        $reset,
        authLogin,
        fetchUserInfo,
        loginLoading,
        logout,
        reauthenticate,
        refreshToken
    }
})
