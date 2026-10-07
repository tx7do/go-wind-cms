import type {AppRoute} from "@/core/router/types";

/**
 * 递归排序路由树（按 meta.order 升序）
 */
export const sortRoutes = (routes: AppRoute[]): AppRoute[] => {
    return routes
        .map((route) => ({
            ...route,
            children: route.children?.length ? sortRoutes(route.children) : undefined,
        }))
        .sort((a, b) => {
            const orderA = a.meta?.order ?? Number.MAX_SAFE_INTEGER;
            const orderB = b.meta?.order ?? Number.MAX_SAFE_INTEGER;
            return orderA - orderB;
        });
};
