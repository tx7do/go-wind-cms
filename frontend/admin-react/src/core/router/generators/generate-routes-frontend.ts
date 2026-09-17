import type {ReactNode} from "react";

import type {AppRouteObject, RouteMeta} from "@/core/router/types";

/**
 * 动态生成路由 - 前端权限过滤方式
 * @param routes - 静态路由表（AppRouteObject[]，含 meta）
 * @param permissions - 用户权限码列表
 * @param forbiddenElement - 无权限时渲染的元素（<Forbidden />）
 */
export async function generateRoutesByFrontend(
    routes: AppRouteObject[],
    permissions: string[],
    forbiddenElement?: ReactNode
): Promise<AppRouteObject[]> {
    // 纯变换：逐层构造新节点（element 等属性按引用共享），绝不回写传入的树。
    // 该树是模块级单例（router/modules/* 经 business-routes 装配），本函数在
    // 未登录时就会以空权限被调用——就地过滤会把带 authority 的节点从单例里
    // 永久剔除，登录后重建也找不回。vue 端以 accessible 层 cloneDeep 防线解决同类
    // 问题（实测 lodash 对嵌套函数按引用保留，cloneDeep 亦可行）；此处取纯过滤：
    // 未触碰子树零拷贝，且避免整树复制含冻结的 React 元素。
    const walk = (nodes: AppRouteObject[]): AppRouteObject[] => {
        const kept: AppRouteObject[] = [];
        for (const node of nodes) {
            const pass = hasPermission(node, permissions);
            const mvwf = menuVisibleWithForbidden(node);
            if (!pass && !mvwf) {
                continue;
            }
            const children = Array.isArray(node.children)
                ? walk(node.children)
                : undefined;
            const next: AppRouteObject =
                children === undefined ? node : {...node, children};
            if (mvwf && forbiddenElement) {
                // 对齐 vue-vben / vue-element：声明了 menuVisibleWithForbidden
                // 的路由保留在菜单中，但无论是否持码，元素一律换成 403
                kept.push({...next, element: forbiddenElement});
            } else {
                kept.push(next);
            }
        }
        return kept;
    };
    return walk(routes);
}

/**
 * 判断路由是否有权限访问
 * @param route - 路由对象
 * @param permissions - 用户权限码列表
 * @returns 是否允许访问
 */
function hasPermission(route: AppRouteObject, permissions: string[]): boolean {
    const meta = route.meta as RouteMeta | undefined;

    // 情况 1: 无权限要求 = 公开路由（登录页、404 等）
    if (!meta?.authority?.length) {
        return true;
    }

    // 情况 2: 配置了 ignoreAccess = 忽略权限检查
    if (meta?.ignoreAccess) {
        return true;
    }

    // 情况 3: 检查权限码数组（authority 中任一元素匹配用户权限即可）
    if (meta.authority.some((code) => permissions.includes(code))) {
        return true;
    }

    // 情况 4: 无权限，但可能是"菜单可见但禁止访问"，交给上层处理
    return false;
}

/**
 * 判断路由是否声明了"菜单可见，但访问返回 403"
 * @param route - 路由对象
 * @returns 是否为 menuVisibleWithForbidden 节点
 */
function menuVisibleWithForbidden(route: AppRouteObject): boolean {
    const meta = route.meta as RouteMeta | undefined;

    // 有权限要求且显式声明了 menuVisibleWithForbidden
    return (
        !!meta?.authority?.length &&
        meta?.menuVisibleWithForbidden === true
    );
}

export {hasPermission};
