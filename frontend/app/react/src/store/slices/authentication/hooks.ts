import {useSelector, useDispatch} from 'react-redux';

import type {IAuthenticationState} from '../../types';
import type {AppDispatch, RootState} from '@/store';
import {authenticationservicev1_GrantType, identityservicev1_User} from '@/api/generated/app/service/v1';
import {setLoginLoading} from './slice';
import {
    createAuthenticationServiceClient,
    createUserProfileServiceClient,
} from '@/api/generated/app/service/v1';
import {requestClientRequestHandler} from '@/transport/rest';

export interface LoginParams {
    username?: string;
    email?: string;
    mobile?: string;
    password: string;
    grant_type?: string;
}

/**
 * 加密密码（需要实现 AES 加密）
 * @param password 明文密码
 */
function encryptPassword(password: string): string {
    // TODO: 实现 AES 加密逻辑
    // const key = import.meta.env.VITE_AES_KEY;
    // return encryptData(password, key, key);

    // 暂时返回明文，后续需要实现加密
    console.warn('密码未加密，需要实现 AES 加密逻辑');
    return password;
}

export function useAuthenticationStore() {
    const dispatch = useDispatch<AppDispatch>();
    const authentication = useSelector<RootState, IAuthenticationState>(
        (state) => state.authentication,
    );

    // 创建服务客户端
    const authService = createAuthenticationServiceClient(
        requestClientRequestHandler,
    );
    const userProfileService = createUserProfileServiceClient(
        requestClientRequestHandler,
    );

    /**
     * 异步处理用户登录操作并获取 accessToken
     */
    async function authLogin(
        params: LoginParams,
        onSuccess?: () => Promise<void> | void,
    ) {
        let userInfo: identityservicev1_User | null = null;

        try {
            dispatch(setLoginLoading(true));
            // 调用登录 API
            const response = await authService.Login({
                username: params.username,
                email: params.email,
                mobile: params.mobile,
                password: encryptPassword(params.password),
                grant_type: (params.grant_type || "password") as authenticationservicev1_GrantType,
            });
            if (response.access_token) {
                userInfo = await fetchUserInfo();
                // TODO: 更新 user store
                // userStore.setUserInfo(userInfo);
                // dispatch(setUserInfo(userInfo));
                if (onSuccess) {
                    await onSuccess();
                }
            }
        } catch (error) {
            throw error;
        } finally {
            dispatch(setLoginLoading(false));
        }
        return {
            userInfo,
        };
    }

    /**
     * 拉取用户信息
     */
    async function fetchUserInfo(): Promise<identityservicev1_User> {
        return await userProfileService.GetUser({});
    }

    return {
        authentication,
        dispatch,
        authLogin,
        fetchUserInfo,
    };
}

