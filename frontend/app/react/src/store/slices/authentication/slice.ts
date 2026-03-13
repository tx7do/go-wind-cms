import {createSlice, createAsyncThunk, PayloadAction} from '@reduxjs/toolkit';

import {createAuthenticationServiceClient, createUserProfileServiceClient} from "@/api/generated/app/service/v1";

import {requestApi} from "@/transport/rest";

import {setAccessToken, setLoginExpired} from "@/store/core/access/slice";
import {useAccessStore} from "@/store/core/access/hooks";

interface AuthState {
    loginLoading: boolean;
    loading: boolean;
    error: string | null;
}

const initialState: AuthState = {
    loginLoading: false,
    loading: false,
    error: null,
};

// 登录
export const login = createAsyncThunk(
    'authentication/login',
    async (
        params: { username?: string; email?: string; mobile?: string; password: string },
        {dispatch}
    ) => {
        const authnService = createAuthenticationServiceClient(requestApi);
        const userProfileService = createUserProfileServiceClient(requestApi);

        dispatch(setLoginLoading(true));
        try {
            const {access_token} = await authnService.Login({
                username: params.username,
                email: params.email,
                mobile: params.mobile,
                password: params.password, // 建议加密处理
                grant_type: 'password',
            });

            if (access_token) {
                // dispatch(setAccessToken(access_token));
                // 获取用户信息
                const userInfo = await userProfileService.GetUser({});
                // dispatch(setUserInfo(userInfo));
                // dispatch(setLoginExpired(false));
                return {userInfo};
            }
            return {userInfo: null};
        } finally {
            dispatch(setLoginLoading(false));
        }
    }
);

// 登出
export const logout = createAsyncThunk('authentication/logout', async (_, {dispatch}) => {
    const authnService = createAuthenticationServiceClient(requestApi);

    try {
        await authnService.Logout({});
    } catch {
    } finally {
        dispatch(setAccessToken(null));
        dispatch(setLoginExpired(false));

        dispatch(setLoginLoading(false));
    }
});

// 刷新token
export const refreshToken = createAsyncThunk(
    'authentication/refreshToken',
    async (refresh_token: string, {dispatch}) => {
        const authnService = createAuthenticationServiceClient(requestApi);

        const resp = await authnService.RefreshToken({
            grant_type: 'password',
            refresh_token,
        });
        // dispatch(setAccessToken(resp.access_token));
        return resp.access_token;
    }
);

export const reauthenticate = createAsyncThunk(
    'authentication/reauthenticate',
    async (_, {dispatch}) => {
        console.warn('Access token or refresh token is invalid or expired. ')

        setAccessToken(null)

        const accessStore = useAccessStore();
        const refresh_token = accessStore.access.refreshToken;
        if (refresh_token && refresh_token.value) {
            const access_token = await dispatch(refreshToken(refresh_token.value));
            if (access_token) {
                // dispatch(setAccessToken(access_token));
                dispatch(setLoginExpired(false));
            }
        }
    }
)

const authenticationSlice = createSlice({
    name: 'authentication',
    initialState,
    reducers: {
        setLoginLoading(state, action: PayloadAction<boolean>) {
            state.loginLoading = action.payload;
        },
        setLoading(state, action: PayloadAction<boolean>) {
            state.loading = action.payload;
        },
        setError(state, action: PayloadAction<string | null>) {
            state.error = action.payload;
        },
        clearError(state) {
            state.error = null;
        },
        resetAuth(state) {
            state.loginLoading = false;
            state.loading = false;
            state.error = null;
        },
    },
    extraReducers: (builder) => {
        // 可扩展异步 action 状态
    },
});

export const {setLoginLoading, resetAuth} = authenticationSlice.actions;
export default authenticationSlice.reducer;
