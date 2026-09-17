import { Navigate, useLocation } from 'react-router-dom';

import type { AppRouteObject } from '@/core/router/types';
import { useAuthStore } from '@/stores';

import BlankLayout from '@/layouts/BlankLayout';
import RouteErrorFallback from '@/layouts/components/ErrorFallback/RouteErrorFallback';
import {
  Unauthorized,
  Forbidden,
  InternalError,
  Offline,
  ComingSoon,
  NotFound,
} from '@/pages/core/error';

/**
 * 认证感知的 404：未登录时命中通配符（如刷新受保护深链 /content/posts），
 * 展示 404 是误导——真实原因是路由树尚未包含该页。跳登录并携带回跳地址；
 * 已登录时的未知路径仍展示 404（真实不存在）。
 */
function AuthAwareNotFound() {
  const location = useLocation();
  const accessToken = useAuthStore((s) => s.accessToken);

  if (!accessToken) {
    const redirect = encodeURIComponent(location.pathname + location.search);
    return <Navigate to={`/auth/login?redirect=${redirect}`} replace />;
  }
  return <NotFound />;
}

export const errorRoutes: AppRouteObject[] = [
  {
    name: 'error-pages',
    path: '/',
    element: <BlankLayout />,
    errorElement: <RouteErrorFallback />,
    meta: { title: 'routes:error-pages', hideInMenu: true, hideInTab: true },
    children: [
      {
        name: 'unauthorized',
        path: '401',
        element: <Unauthorized />,
        meta: { title: 'routes:unauthorized', ignoreAccess: true, hideInMenu: true, hideInTab: true },
      },
      {
        name: 'forbidden',
        path: '403',
        element: <Forbidden />,
        meta: { title: 'routes:forbidden', ignoreAccess: true, hideInMenu: true, hideInTab: true },
      },
      {
        name: 'server-error',
        path: '500',
        element: <InternalError />,
        meta: { title: 'routes:server-error', ignoreAccess: true, hideInMenu: true, hideInTab: true },
      },
      {
        name: 'offline',
        path: 'offline',
        element: <Offline />,
        meta: { title: 'routes:offline', ignoreAccess: true, hideInMenu: true, hideInTab: true },
      },
      {
        name: 'coming-soon',
        path: 'coming-soon',
        element: <ComingSoon />,
        meta: { title: 'routes:coming-soon', ignoreAccess: true, hideInMenu: true, hideInTab: true },
      },

      // ========== 404 通配符路由（必须放在最后！） ==========
      {
        name: 'not-found',
        path: '*',
        element: <AuthAwareNotFound />,
        meta: { title: 'routes:not-found', ignoreAccess: true, hideInMenu: true, hideInTab: true },
      },
    ],
  },
];

export default errorRoutes;
