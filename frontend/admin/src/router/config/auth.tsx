import { Navigate } from 'react-router-dom';

import { type AppRouteObject } from '@/core/router';
import { GuestGuard } from '@/router/guards';

import UserLayout from '@/layouts/UserLayout';
import RouteErrorFallback from '@/layouts/components/ErrorFallback/RouteErrorFallback.tsx';

import Login from '@/pages/core/auth/login';
import MfaChallenge from '@/pages/core/auth/mfa-challenge';
import ForgotPassword from '@/pages/core/auth/forgot-password';

/**
 * 认证相关路由配置
 * 登录等页面
 * 这些路由不受 AuthGuard 保护，使用 GuestGuard 防止已登录用户访问
 */
export const authRoutes: AppRouteObject[] = [
  {
    // 惯用入口别名：/login 直达登录页，避免落进 404 通配符后带出
    // redirect=%2Flogin 导致登录成功又跳回 404 页
    name: 'login-alias',
    path: '/login',
    element: <Navigate to="/auth/login" replace />,
    meta: { title: 'routes:login', ignoreAccess: true, hideInMenu: true, hideInTab: true },
  },
  {
    name: 'auth',
    path: '/auth',
    element: <UserLayout requireAuth={false} />,
    errorElement: <RouteErrorFallback />,
    meta: { title: 'routes:auth', ignoreAccess: true, hideInMenu: true, hideInTab: true },
    children: [
      {
        name: 'login',
        path: 'login',
        element: (
          <GuestGuard>
            <Login />
          </GuestGuard>
        ),
        meta: { title: 'routes:login', ignoreAccess: true },
      },
      {
        name: 'mfa-challenge',
        path: 'mfa-challenge',
        element: <MfaChallenge />,
        meta: { title: 'routes:mfaChallenge', ignoreAccess: true, hideInMenu: true },
      },
      {
        name: 'forgot-password',
        path: 'forgot-password',
        element: (
          <GuestGuard>
            <ForgotPassword />
          </GuestGuard>
        ),
        meta: { title: 'routes:forgot-password', ignoreAccess: true },
      },
    ],
  },
];

export default authRoutes;
