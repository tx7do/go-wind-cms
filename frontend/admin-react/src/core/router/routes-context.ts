import { createContext } from "react";

import type { AppRouteObject } from "./types";

/**
 * 当前实际挂载的路由树（createAccessibleRouter 的产物）：
 * 后端模式=后端下发路由；前端模式=权限过滤后的静态路由。
 * 侧栏等消费方必须镜像此树保证"菜单即路由"——后端模式下若回退静态全量表，
 * 会出现未授权菜单可见、点击 404 的越权观感。
 */
export const AccessibleRoutesContext = createContext<AppRouteObject[] | null>(null);
