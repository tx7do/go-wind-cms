import {useEffect} from 'react';
import Taro from '@tarojs/taro';
import {useTranslation} from 'react-i18next';

/**
 * 根据当前语言动态设置页面标题（导航栏 + document.title）
 *
 * 用法：
 *   usePageTitle('page.home');        // i18n key
 *   usePageTitle('My Custom Title');  // 固定字符串
 *
 * 语言切换后自动更新。
 */
export function usePageTitle(titleKey: string) {
    const {t} = useTranslation();

    useEffect(() => {
        const title = t(titleKey);
        Taro.setNavigationBarTitle({title});
    }, [t, titleKey]);
}
