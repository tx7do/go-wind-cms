import Taro from '@tarojs/taro';

/**
 * tabBar 页面（与 src/app.config.ts 的 tabBar.list 保持一致）
 * 小程序端跳转 tabBar 页必须使用 switchTab，navigateTo 会失败
 */
const TAB_PAGES = new Set([
    '/pages/index/index',
    '/pages/post/index',
    '/pages/favorites/index',
    '/pages/user/index',
]);

/**
 * 将短路径转为 Taro 内部页面路径
 *
 * 支持的模式：
 * /                       → /pages/index/index
 * /posts                  → /pages/post/index
 * /post/detail?id=1       → /pages/post/detail/index?id=1
 * /post/1                 → /pages/post/detail/index?id=1
 * /pages/...              → 原样返回
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

    // 带数字 ID 的短路径：/post/1 → /pages/post/detail/index?id=1
    const detailRoutes = ['post', 'tag', 'category', 'page'];
    if (segments.length === 2 && detailRoutes.includes(segments[0]) && /^\d+$/.test(segments[1])) {
        const fullPath = `/pages/${segments[0]}/detail/index`;
        return queryString ? `${fullPath}?id=${segments[1]}&${queryString}` : `${fullPath}?id=${segments[1]}`;
    }

    // 通用转换：/about → /pages/about/index
    const fullPath = `/pages/${segments.join('/')}/index`;
    return queryString ? `${fullPath}?${queryString}` : fullPath;
}

/**
 * tabBar 页参数中转：
 * switchTab 无法携带 query，先将参数写入 Storage，目标页 onShow 时读取并清除
 */
const TAB_QUERY_PREFIX = 'tabQuery:';

export function stashTabQuery(fullPath: string) {
    const [path, query] = fullPath.split('?');
    if (query) {
        Taro.setStorageSync(`${TAB_QUERY_PREFIX}${path}`, query);
    }
}

/** 读取并清除 switchTab 携带的页面参数（非 tab 参数场景返回空对象） */
export function consumeTabQuery(pagePath: string): Record<string, string> {
    const key = `${TAB_QUERY_PREFIX}${pagePath}`;
    let query = '';
    try {
        query = Taro.getStorageSync(key) || '';
    } catch {
        return {};
    }
    if (query) {
        try {
            Taro.removeStorageSync(key);
        } catch {
            // ignore
        }
    }
    const result: Record<string, string> = {};
    if (query) {
        query.split('&').forEach((pair) => {
            const [k, v] = pair.split('=');
            if (k) result[decodeURIComponent(k)] = decodeURIComponent(v || '');
        });
    }
    return result;
}

function isTabPage(fullPath: string): boolean {
    const [path] = fullPath.split('?');
    return TAB_PAGES.has(path);
}

/**
 * 页面栈深度。栈底（=1）时 navigateBack 在两端都是坏路径：
 * H5 走 history.go(0) → 整页重载，小程序直接 navigateBack:fail。
 */
export function canGoBack(): boolean {
    return Taro.getCurrentPages().length > 1;
}

/**
 * Taro 路由 Hook - 替代 React 版本的 useI18nRouter
 * Taro 不使用 locale 前缀路由，直接使用路径
 */
export function useI18nRouter() {
    const push = (path: string) => {
        const full = resolveTaroPath(path);
        if (isTabPage(full)) {
            // tabBar 页：switchTab 跳转，参数经 Storage 中转
            stashTabQuery(full);
            Taro.switchTab({url: full.split('?')[0]});
        } else {
            Taro.navigateTo({url: full});
        }
    };

    const replace = (path: string) => {
        const full = resolveTaroPath(path);
        if (isTabPage(full)) {
            // redirectTo 不能打开 tabBar 页，用 reLaunch 兜底
            Taro.reLaunch({url: full.split('?')[0]});
        } else {
            Taro.redirectTo({url: full});
        }
    };

    const reLaunch = (path: string) => {
        if (typeof window !== 'undefined') {
            // 加时间戳强制浏览器绕过缓存，必定整页重载
            const separator = path.includes('?') ? '&' : '?';
            const url = window.location.origin + path + separator + '_t=' + Date.now();
            window.location.href = url;
        } else {
            Taro.reLaunch({url: resolveTaroPath(path)});
        }
    };

    const back = () => {
        if (canGoBack()) {
            Taro.navigateBack();
        } else {
            // 栈底（分享链接直达、刷新后）退回首页；'/' 是 tab 页，push 内部走 switchTab
            push('/');
        }
    };

    return {
        push,
        replace,
        reLaunch,
        back,
        localizedPath: (path: string) => path,
    };
}
