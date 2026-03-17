import {useEffect} from 'react';
import Taro, {useLoad} from '@tarojs/taro';
import {View, Text} from '@tarojs/components';

import {DEFAULT_LANGUAGE} from '@/i18n';
import './index.scss';

const getDefaultLocale = () => {
    // 优先从本地存储获取用户上次选择的语言
    try {
        const savedLocale = Taro.getStorageSync('locale');
        if (savedLocale) {
            return savedLocale;
        }
    } catch (e) {
        // 忽略存储读取错误
    }
    
    // 其次根据浏览器语言判断
    if (typeof navigator !== 'undefined') {
        const browserLanguage = navigator.language;
        if (browserLanguage.startsWith('en')) {
            return 'en-US';
        }
        if (browserLanguage.startsWith('zh')) {
            return 'zh-CN';
        }
    }
    
    // 最后使用默认语言
    return DEFAULT_LANGUAGE;
};

export default function Index() {
    useLoad(() => {
        console.log('Page loaded.');
    });

    useEffect(() => {
        const locale = getDefaultLocale();
        // Taro 中使用 navigateTo 进行路由跳转
        Taro.navigateTo({
            url: `/${locale}/index`,
            fail: () => {
                // 如果 navigateTo 失败（可能是 tabBar 页面），使用 switchTab
                Taro.switchTab({
                    url: `/${locale}/index`,
                }).catch(() => {
                    console.error('Navigation failed');
                });
            }
        });
    }, []);

    return (
        <View className='index'>
            <View className='loading-container'>
                <Text className='loading-text'>Redirecting to home page...</Text>
            </View>
        </View>
    );
}
