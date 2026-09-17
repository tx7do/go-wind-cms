import {mapTree} from "@/utils";
import type {AppRoute} from "@/core/router/types";

/**
 * 自动注入重定向
 * 规则：父路由有子路由且未配置 redirect 时，指向第一个子路由
 */
export const injectRedirects = (routes: AppRoute[]): AppRoute[] => {
    return mapTree(routes, (route) => {
        // 已有 redirect 或无子路由，直接返回
        if (route.redirect || !route.children?.length) {
            return route;
        }

        const firstChild = route.children[0];

        // 子路由路径无效（相对路径/含参数），跳过
        if (!firstChild?.path || !firstChild.path.startsWith('/')) {
            return route;
        }

        // 注入重定向（不可变更新）
        return {
            ...route,
            redirect: firstChild.path,
        };
    });
};
