import { useMemo } from 'react';
import { useAuthStore } from '@/stores/auth';
import { useUserStore } from '@/stores/user';

/** 光标调色板（按用户 id 取色，稳定区分协作者） */
const CURSOR_PALETTE = [
  '#f78fb3', '#30b8db', '#7bd389', '#f5a623',
  '#9b8cff', '#5fb4ea', '#e5a17e', '#68c7a0',
];

/**
 * 实时协同选项（富文本编辑器）：
 * 房间 = {targetType}:{id}:{locale}；token 走 ws query 参数（浏览器 WS 不能带头）。
 * 未登录/无 id（创建模式）返回 null = 关闭协同。
 */
export function useCollaboration(
  targetType: 'post' | 'page',
  targetId?: number | null,
  locale?: string,
) {
  const token = useAuthStore((s) => s.accessToken);
  const userId = useUserStore((s) => s.userInfo?.id);
  const userName =
    useUserStore(
      (s) =>
        s.userInfo?.nickname ||
        s.userInfo?.realname ||
        s.userInfo?.username ||
        '',
    );

  return useMemo(() => {
    if (!targetId || !token) return null;
    const id = userId ?? 0;
    return {
      room: `${targetType}:${targetId}:${locale || 'zh-CN'}`,
      wsUrl: import.meta.env.VITE_COLLAB_WS_URL ?? 'ws://localhost:6602',
      token,
      userName: userName || `用户${id}`,
      userColor: CURSOR_PALETTE[id % CURSOR_PALETTE.length],
    };
  }, [targetType, targetId, locale, token, userId, userName]);
}
