import {acceptHMRUpdate, defineStore} from 'pinia'

interface IUser {
  [key: string]: any
  avatar: string
  id: number
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
    tenantId: (state): number | null => state.user?.tenantId ?? null,
    isLoggedIn: (state): boolean => state.user !== null,
    isTenantUser: (state): boolean => {
      const tenantId = state.user?.tenantId;
      return tenantId !== null && tenantId !== undefined && tenantId > 0;
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
