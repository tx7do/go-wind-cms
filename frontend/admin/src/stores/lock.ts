import { create } from 'zustand';
import { persist } from 'zustand/middleware';

interface LockState {
  isLockScreen: boolean;
  lockScreenPassword?: string;
  lockScreen: (password: string) => void;
  unlockScreen: () => void;
}

/**
 * 锁屏密码明文持久化是有意的：锁屏只是"离开工位遮个屏"的界面遮罩，不是认证边界，
 * 真正的鉴权仍由后端令牌负责。auth-storage 被清掉时遮罩一并失效，不会绕过登录。
 */
export const useLockStore = create<LockState>()(
  persist(
    (set) => ({
      isLockScreen: false,
      lockScreenPassword: undefined,
      lockScreen: (password) => set({ isLockScreen: true, lockScreenPassword: password }),
      unlockScreen: () => set({ isLockScreen: false, lockScreenPassword: undefined }),
    }),
    {
      name: 'lock-storage',
      partialize: (state) => ({
        isLockScreen: state.isLockScreen,
        lockScreenPassword: state.lockScreenPassword,
      }),
    },
  ),
);
