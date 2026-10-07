import type {AppRouteObject, RouteMeta, RouteModuleRecord, RouteModuleType} from "@/core/router";


/**
 * 合并动态路由模块的默认导出
 *
 * @param routeModules - Vite import.meta.glob 返回的模块对象
 * @param options - 可选配置
 * @returns 合并后的路由配置数组（AppRouteObject[]）
 *
 * @example
 *
 */
export function mergeRouteModules(
    routeModules: RouteModuleRecord,
    options?: {
        /** 是否过滤空模块 */
        filterEmpty?: boolean;
        /** 是否按路径排序 */
        sortByPath?: boolean;
    }
): AppRouteObject[] {
    const {filterEmpty = true, sortByPath = false} = options || {};

    const mergedRoutes: AppRouteObject[] = [];

    for (const [modulePath, routeModule] of Object.entries(routeModules)) {
        try {
            // 兼容多种导出形式
            let moduleRoutes: AppRouteObject[] = [];

            if (typeof routeModule === 'function') {
                // 异步模块：跳过（需用 mergeRouteModulesAsync）
                continue;
            } else if (Array.isArray(routeModule)) {
                // 直接导出数组
                moduleRoutes = routeModule;
            } else if (routeModule && typeof routeModule === 'object') {
                // 默认导出对象
                moduleRoutes = (routeModule as RouteModuleType).default || [];
            }

            // 过滤空模块
            if (filterEmpty && moduleRoutes.length === 0) {
                continue;
            }

            // 可选：注入模块路径信息（便于调试）
            if (import.meta.env.DEV) {
                moduleRoutes.forEach((route) => {
                    if (!route.meta) {
                        route.meta = {} as RouteMeta;
                    }
                    (route.meta as any).__modulePath = modulePath;
                });
            }

            mergedRoutes.push(...moduleRoutes);
        } catch (error) {
            console.warn(`Failed to load route module: ${modulePath}`, error);
        }
    }

    // 可选：按路径字典序排序
    if (sortByPath) {
        mergedRoutes.sort((a, b) => (a.path || '').localeCompare(b.path || ''));
    }

    return mergedRoutes;
}

/**
 * 异步合并动态路由模块（用于 lazy 加载）
 * @param routeModules - Vite import.meta.glob 返回的异步模块对象
 * @returns Promise<AppRouteObject[]>
 */
export async function mergeRouteModulesAsync(
    routeModules: Record<string, () => Promise<RouteModuleType | AppRouteObject[]>>
): Promise<AppRouteObject[]> {
    const results = await Promise.all(
        Object.values(routeModules).map(async (loader) => {
            try {
                const module = await loader();
                if (Array.isArray(module)) {
                    return module;
                }
                return (module as RouteModuleType).default || [];
            } catch (error) {
                console.warn('Failed to async load route module:', error);
                return [];
            }
        })
    );

    return results.flat();
}
