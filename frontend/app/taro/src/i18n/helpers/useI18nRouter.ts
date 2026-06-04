import Taro from '@tarojs/taro';

/**
 * Taro 路由 Hook - 替代 React 版本的 useI18nRouter
 * Taro 不使用 locale 前缀路由，直接使用路径
 */
export function useI18nRouter() {
    const push = (path: string) => {
        Taro.navigateTo({url: path});
    };

    const replace = (path: string) => {
        Taro.redirectTo({url: path});
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
