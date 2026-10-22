import {View, Text, Image} from '@tarojs/components';
import {useTranslations} from '@/lib/next-intl-compat';
import logoImage from '@/assets/images/logo.png';
import {XIcon} from '@/plugins/xicon';
import {canGoBack, useI18nRouter} from '@/i18n/helpers/useI18nRouter';

export default function Header() {
    const appT = useTranslations('app');
    const commonT = useTranslations('common');
    const router = useI18nRouter();

    // 回首页：直接用浏览器原生导航绕过 Taro SPA 路由
    const goHome = () => {
        if (typeof window !== 'undefined') {
            window.location.href = window.location.origin + '/';
        }
    };

    // 栈深在渲染期读取：Layout 重渲染由 children 变化驱动（与 NavigationProgress
    // 读 window.location.pathname 同一机制），每次路由切换都会重新取值。
    const showBack = canGoBack();

    // 状态栏留白交给 .app-header 的 padding-top: var(--safe-top)。
    // 原先用 getSystemInfoSync().statusBarHeight 计算：H5 恒为 0（白等一次异步），
    // 而本组件只在 H5 渲染，env(safe-area-inset-top) 才是浏览器给得出的那个值。
    return (
        <View className='app-header fixed top-0 left-0 right-0 z-[1000]'>
            <View className='flex h-[96rpx] items-center justify-between px-[28rpx]'>
                {/* 二级页给返回键，栈底（tab 页）给品牌 —— 对应小程序原生栏的同一分工 */}
                {showBack ? (
                    <View
                      className='flex items-center gap-[8rpx] min-w-touch min-h-touch'
                      aria-label={commonT('back')}
                      onClick={() => router.back()}
                      hoverClass='tap-active'
                    >
                        <XIcon name='carbon:chevron-left' size={26} className='text-textMain' />
                        <Text className='text-body text-textMain'>{commonT('back')}</Text>
                    </View>
                ) : (
                    <View
                      className='flex items-center gap-[16rpx] min-w-touch min-h-touch'
                      onClick={goHome}
                      hoverClass='tap-active'
                    >
                        <View className='h-[56rpx] w-[56rpx] rounded-[16rpx] overflow-hidden shadow-[0_6rpx_18rpx_rgba(10,108,255,0.35)]'>
                            <Image src={logoImage} mode='aspectFill' className='w-full h-full' />
                        </View>
                        <Text className='header-brand text-card-title font-bold whitespace-nowrap'>
                            {appT('title')}
                        </Text>
                    </View>
                )}

                {/* 右侧操作区：搜索 / 个人中心 */}
                <View className='flex items-center gap-[8rpx]'>
                    <View
                      className='flex items-center justify-center w-[68rpx] h-[68rpx] rounded-full bg-primary-tint'
                      onClick={() => router.push('/search')}
                      hoverClass='tap-active'
                    >
                        <XIcon name='carbon:search' size={24} className='text-primary' />
                    </View>
                    <View
                      className='flex items-center justify-center w-[68rpx] h-[68rpx] rounded-full bg-primary-tint'
                      onClick={() => router.push('/user')}
                      hoverClass='tap-active'
                    >
                        <XIcon name='carbon:user' size={24} className='text-primary' />
                    </View>
                </View>
            </View>
        </View>
    );
}
