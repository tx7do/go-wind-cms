import {ref} from 'vue'
import {useRouter} from 'vue-router'
import {defineStore} from 'pinia'
import CryptoJS from 'crypto-js';

import {DEFAULT_HOME_PATH, LOGIN_PATH} from '@/constants'
import type {Recordable} from '@/typings/helper'
import type {UserInfo} from "@/typings";

import {useAccessStore, useUserStore} from '@/stores'
import {resetAllStores} from "@/stores/setup";

import {$t} from '@/locales';
import {
  createAuthenticationServiceClient,
  createUserProfileServiceClient
} from "@/api/generated/app/service/v1";
import {requestClientRequestHandler} from "@/transport/rpc/request";

export const useAuthStore = defineStore('auth', () => {
  const accessStore = useAccessStore()
  const userStore = useUserStore()
  const authnService = createAuthenticationServiceClient(
    requestClientRequestHandler,
  );
  const userProfileService = createUserProfileServiceClient(
    requestClientRequestHandler,
  );

  const router = useRouter()

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

      const {access_token} = await authnService.Login({
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
            : await router.push(userInfo.homePath || DEFAULT_HOME_PATH)
        }

        if (userInfo?.realname) {
          window.$message.success($t('authentication.login.login_success'))
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
      await authnService.Logout({})
    } catch {
      // 不做任何处理
    } finally {
      resetAllStores()
      accessStore.setLoginExpired(false)

      // 回登录页带上当前路由地址
      await router.replace({
        path: LOGIN_PATH,
        query: redirect
          ? {
            redirect: encodeURIComponent(router.currentRoute.value.fullPath),
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

    accessStore.setAccessToken(null)
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
    const resp = await authnService.RefreshToken({
      grant_type: 'password',
      refresh_token: accessStore.refreshToken ?? '',
    })
    const newToken = resp.access_token

    accessStore.setAccessToken(newToken)

    return newToken
  }

  /**
   * 拉取用户信息
   */
  async function fetchUserInfo() {
    return (await userProfileService.GetUser({})) as unknown as UserInfo;
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
    const key = import.meta.env.VITE_AES_KEY;
    if (!key) {
      throw new Error('VITE_AES_KEY is not set in environment');
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
