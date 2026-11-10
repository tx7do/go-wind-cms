import { useMemo } from 'react';
import { useAuthStore } from '@/stores/auth';
import { useUserStore } from '@/stores/user';

/** 光标调色板（按用户 id 取色，稳定区分协作者） */
const CURSOR_PALETTE = [
  '#f78fb3', '#30b8db', '#7bd389', '#f5a623',
  '#9b8cff', '#5fb4ea', '#e5a17e', '#68c7a0',
];

/**
 * 把同源相对 ws 路径（如 VITE_COLLAB_WS_URL=/collab）按当前页面来源
 * 绝对化为 ws://（http 页）或 wss://（https 页）。y-websocket 的
 * WebsocketProvider 不接受相对 URL；绝对 ws(s):// 形态则原样返回，
 * 兼容既有按域名配置的部署。SSR/无 window 场景原样返回（admin 为纯 SPA）。
 */
function absolutizeWsUrl(url: string): string {
  if (/^wss?:\/\//i.test(url) || typeof window === 'undefined') return url;
  const proto = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
  return `${proto}//${window.location.host}${url.startsWith('/') ? url : `/${url}`}`;
}

/**
 * 实时协同选项（富文本编辑器）：
 * 房间 = {targetType}:{id}:{locale}；token 走 ws query 参数（浏览器 WS 不能带头）。
 * 未登录/无 id（创建模式）返回 null = 关闭协同。
 */
export function useCollaboration(
  targetType: 'post' | 'page',
  targetId?: string | null,
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
    // userId 是 uint64 字符串；取色需要数值下标，用 Number 折算（< 2^53 安全）
    const id = Number(userId ?? 0) || 0;
    const result = {
      room: `${targetType}:${targetId}:${locale || 'zh-CN'}`,
      wsUrl: absolutizeWsUrl(import.meta.env.VITE_COLLAB_WS_URL || '/collab'),
      token,
      userName: userName || `用户${userId ?? ''}`,
      userColor: CURSOR_PALETTE[id % CURSOR_PALETTE.length],
    };
    return result;
  }, [targetType, targetId, locale, token, userId, userName]);
}
