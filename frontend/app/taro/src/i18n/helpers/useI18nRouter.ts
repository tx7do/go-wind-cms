import Taro from '@tarojs/taro';

/**
 * 将短路径转为 Taro 内部页面路径
 *
 * 支持的模式：
 * /                       → /pages/index/index
 * /about                  → /pages/about/index
 * /category               → /pages/category/index
 * /category/1             → /pages/category/detail/index?id=1
 * /category/detail?id=1   → /pages/category/detail/index?id=1
 * /tag/2                  → /pages/tag/detail/index?id=2
 * /post/3                 → /pages/post/detail/index?id=3
 */
function resolveTaroPath(shortPath: string): string {
    const [path, queryString] = shortPath.split('?');
    const segments = path.split('/').filter(Boolean);

    // 根路径 '/' → /pages/index/index
    if (segments.length === 0) {
        const fullPath = '/pages/index/index';
        return queryString ? `${fullPath}?${queryString}` : fullPath;
    }

    // 已是内部路径 /pages/... 直接返回
    if (segments[0] === 'pages') {
        return queryString ? `${path}?${queryString}` : path;
    }

    // 带数字 ID 的短路径：/category/123 → /pages/category/detail/index?id=123
    const detailRoutes = ['category', 'tag', 'post'];
    if (segments.length === 2 && detailRoutes.includes(segments[0]) && /^\d+$/.test(segments[1])) {
        const fullPath = `/pages/${segments[0]}/detail/index`;
        return queryString ? `${fullPath}?id=${segments[1]}&${queryString}` : `${fullPath}?id=${segments[1]}`;
    }

    // 通用转换：/about → /pages/about/index
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
