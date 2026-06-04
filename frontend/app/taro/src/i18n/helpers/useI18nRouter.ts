import Taro from '@tarojs/taro';

/**
 * 将短路径转为 Taro 内部页面路径
 *
 * /tag/detail?id=2       → /pages/tag/detail/index?id=2
 * /category/detail?id=5  → /pages/category/detail/index?id=5
 * /tag                   → /pages/tag/index
 *
 * Taro H5 router 会根据 customRoutes 自动将内部路径映射为浏览器 URL：
 * /pages/tag/detail/index?id=2 → 浏览器地址栏显示 /tag/2
 */
function resolveTaroPath(shortPath: string): string {
    const [path, queryString] = shortPath.split('?');
    const segments = path.split('/').filter(Boolean);
    // 根路径 '/' → /pages/index/index
    if (segments.length === 0) {
        const fullPath = '/pages/index/index';
        return queryString ? `${fullPath}?${queryString}` : fullPath;
    }
    const fullPath = `/pages/${segments.join('/')}/index`;
    return queryString ? `${fullPath}?${queryString}` : fullPath;
}

/**
 * Taro 路由 Hook - 替代 React 版本的 useI18nRouter
 * Taro 不使用 locale 前缀路由，直接使用路径
 */
export function useI18nRouter() {
    const push = (path: string) => {
        Taro.navigateTo({url: resolveTaroPath(path)});
    };

    const replace = (path: string) => {
        Taro.redirectTo({url: resolveTaroPath(path)});
    };

    const back = () => {
        Taro.navigateBack();
    };

    return {
        push,
        replace,
        back,
        localizedPath: (path: string) => path,
    };
}
