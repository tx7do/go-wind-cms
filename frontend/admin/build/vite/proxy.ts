import type {ProxyOptions} from 'vite';

type ProxyList = [string, string][];

type ProxyTargetList = Record<string, ProxyOptions>;

/**
 * 创建跨域代理
 *
 * 后端 HTTP 路由本身带 /admin/v1 前缀（见 api/gen/.../*_http.pb.go 的
 * r.POST("/admin/v1/...") 注册），故代理仅做 target 转发，不做 path rewrite。
 * 之前的 rewrite 会把 /admin 前缀剥掉，导致后端收到的路径（如 /v1/captcha）
 * 不匹配任何注册路由而 404。
 *
 * @param list - 二维数组参数：[prefix, target]
 */
export function createProxy(list: ProxyList = []) {
    const res: ProxyTargetList = {};

    for (const [prefix, target] of list) {
        res[`^${prefix}`] = {
            target,
            changeOrigin: true,
        };
    }

    return res;
}
