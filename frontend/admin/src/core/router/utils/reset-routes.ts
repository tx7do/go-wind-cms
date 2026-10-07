import {traverseTreeValues} from '@/utils';
import type {AppRouteObject} from "@/core/router/types";

/**
 * 重置可访问路由（保留白名单）
 *
 * ⚠️ React Router v6+ 不支持运行时修改路由，
 * 此函数返回过滤后的路由树，需配合 createBrowserRouter 重新创建 router
 *
 * @param allRoutes - 当前完整路由树（含动态路由）
 * @param staticRouteNames - 静态路由白名单（保留这些路由）
 * @returns 过滤后的路由树（仅含白名单路由）
 */
export function resetAccessibleRoutes(
    allRoutes: AppRouteObject[],
    staticRouteNames: string[]
): AppRouteObject[] {
    // 1. 构建白名单 Set（O(1) 查找）
    const whitelist = new Set(staticRouteNames.filter(Boolean));

    // 2. 递归过滤路由树
    return filterRoutesByWhitelist(allRoutes, whitelist);
}

/**
 * 递归过滤路由树，仅保留白名单中的路由
 */
function filterRoutesByWhitelist(
    routes: AppRouteObject[],
    whitelist: Set<string>
): AppRouteObject[] {
    return routes.reduce<AppRouteObject[]>((acc, route) => {
        // 路由必须有 name 才能被过滤
        if (!route.name) {
            console.warn(`Route with path "${route.path}" is missing a "name" field and will be skipped.`);
            return acc;
        }

        // 如果当前路由在白名单中，保留它并递归处理子路由
        if (whitelist.has(route.name)) {
            const filtered = {...route};

            if (route.children?.length) {
                const childFiltered = filterRoutesByWhitelist(route.children, whitelist);
                if (childFiltered.length > 0) {
                    filtered.children = childFiltered;
                } else {
                    // 子路由全被过滤，但父路由在白名单中：保留父路由但清空 children
                    filtered.children = undefined;
                }
            }

            acc.push(filtered);
        }
        // 如果当前路由不在白名单，直接跳过（及其所有子路由）

        return acc;
    }, []);
}

/**
 * 提取路由树中所有有效的 name（用于生成白名单）
 * @param routes - 路由树
 * @returns 有效的路由 name 数组
 */
export function extractRouteNames(routes: AppRouteObject[]): string[] {
    return traverseTreeValues<AppRouteObject, string | undefined>(
        routes,
        (route) => {
            if (!route.name) {
                console.warn(`Route with path "${route.path}" is missing a "name" field.`);
                return undefined;
            }
            return route.name;
        }
    ).filter((name): name is string => Boolean(name));
}
