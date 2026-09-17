import { createBrowserRouter, type RouteObject } from 'react-router-dom';
import { injectRedirects } from './utils/inject-redirect';
import { sortRoutes } from './utils/sort-routes';
import { transformRoutesWithHandle } from './utils/transform-meta-to-handle';
import type { GenerateMenuAndRoutesOptions, AppRoute, AppRouteObject } from './types';
import { generateRoutesByBackend, generateRoutesByFrontend } from '@/core/router/generators';
import type { AccessModeType } from '@/core/preferences';

/**
 * 从路由列表中分离出：
 * - layoutRoutes: 包含 MainLayout/AuthGuard 的根路由（path='/'）
 * - staticRoutes: 不受 AuthGuard 保护的静态路由（auth/login/error 等）
 */
function separateRoutes(routes: AppRouteObject[]) {
  const layoutRoutes: AppRouteObject[] = [];  // path='/' 的布局路由
  const otherRoutes: AppRouteObject[] = [];    // 其他静态路由（auth/error等）

  for (const route of routes) {
    if (route.path === '/' && route.children) {
      layoutRoutes.push(route);
    } else {
      otherRoutes.push(route);
    }
  }

  return { layoutRoutes, otherRoutes };
}

export const createAccessibleRouter = async (
  mode: AccessModeType,
  options: GenerateMenuAndRoutesOptions,
): Promise<{ router: ReturnType<typeof createBrowserRouter>; routes: AppRouteObject[] }> => {
  let routes: AppRouteObject[] = [...options.routes];

  // 根据模式生成路由
  switch (mode) {
    case 'backend': {
      // 后端模式：从 API 获取路由树，动态转换组件
      if (!options.fetchMenuListAsync) {
        console.warn('[Router] Backend mode requires fetchMenuListAsync, falling back to frontend mode');
        routes = await generateRoutesByFrontend(
          routes,
          options.permissions ?? [],
          options.forbiddenElement,
        );
      } else {
        // 分离布局路由与静态路由（auth/error 等不受 AuthGuard 保护）
        const { otherRoutes } = separateRoutes(routes);

        // 后端返回的路由树（根节点 component="BasicLayout"，已包含 Layout）。
        // 注意：后端生成器只消费 fetchMenuListAsync/layoutMap/pageMap
        const backendRoutes = await generateRoutesByBackend({
          fetchMenuListAsync: options.fetchMenuListAsync,
          layoutMap: options.layoutMap,
          pageMap: options.pageMap,
        });

        // 合并：后端路由 + 静态路由（auth/error）
        routes = [...backendRoutes, ...otherRoutes];
      }
      break;
    }
    case 'frontend':
    default: {
      // 前端模式：基于静态路由 + 权限过滤
      routes = await generateRoutesByFrontend(
        routes,
        options.permissions ?? [],
        options.forbiddenElement,
      );
      break;
    }
  }

  if (options.autoInjectRedirect !== false)
    routes = injectRedirects(routes as unknown as AppRoute[]) as unknown as AppRouteObject[];
  if (options.autoSort !== false)
    routes = sortRoutes(routes as unknown as AppRoute[]) as unknown as AppRouteObject[];

  // 将 meta 转换为 handle，使 useMatches() 能获取路由元数据
  routes = transformRoutesWithHandle(routes);

  const router = createBrowserRouter(routes as RouteObject[], {
    future: {
      v7_relativeSplatPath: true,
    },
  });
  // 路由树随 router 一并返回：侧栏（MainLayout→useMenuData）需镜像实际挂载的路由
  return { router, routes };
};
