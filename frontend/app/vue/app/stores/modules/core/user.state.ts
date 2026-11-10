import {acceptHMRUpdate, defineStore} from 'pinia'

interface IUser {
    [key: string]: any

    avatar: string
    // ID 已升级为 uint64（proto3 JSON 下为字符串）
    id: string
    nickname: string
    realname: string
    roles?: string[]
    username: string
}

interface IUserState {
    user: IUser | null
    roles: string[]
}

export const useUserStore = defineStore('core-user', {
    state: (): IUserState => ({
        user: null,
        roles: [],
    }),
    getters: {
        tenantId: (state): string | null => state.user?.tenantId ?? null,
        isLoggedIn: (state): boolean => state.user !== null,
        isTenantUser: (state): boolean => {
            // tenantId 为 uint64 字符串；空串与 '0' 均视为无租户
            const tenantId = state.user?.tenantId;
            return !!tenantId && tenantId !== '0';
        },
    },
    actions: {
        setUserInfo(userInfo: IUser | null) {
            this.user = userInfo;
            this.roles = userInfo?.roles ?? [];
        },

        setUserRoles(newRoles: string[]) {
            this.roles = newRoles;
            if (this.user) this.user.roles = newRoles;
        },

        $reset() {
            this.user = null;
            this.roles = [];
        },
    },
})

const hot = import.meta.hot
if (hot)
    hot.accept(acceptHMRUpdate(useUserStore, hot))
