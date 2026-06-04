import {View, Text, Image} from '@tarojs/components';
import {useTranslations} from '@/lib/next-intl-compat';
import logoImage from '@/assets/images/logo.png';
import {XIcon} from '@/plugins/xicon';
import {useI18nRouter} from '@/i18n/helpers/useI18nRouter';
import MobileNav from './MobileNav';

export default function Header() {
    const appT = useTranslations('app');
    const brandTitle = appT('title');
    const router = useI18nRouter();

    return (
        <View className='fixed top-0 left-0 right-0 z-[1000] bg-cardBg border-b-[1rpx] border-splitLine'>
            <View className='flex h-[88rpx] items-center justify-between px-[24rpx]'>
                {/* Logo */}
                <View
                  className='flex items-center gap-[16rpx] min-w-touch min-h-touch justify-center'
                  onClick={() => router.push('/')}
                  hoverClass='tap-active'
                >
                    <View className='h-[48rpx] w-[48rpx]'>
                        <Image src={logoImage} mode='aspectFit' className='w-full h-full' />
                    </View>
                    <Text className='text-card-title font-bold text-primary whitespace-nowrap'>
                        {brandTitle}
                    </Text>
                </View>

                {/* 右侧操作区 */}
                <View className='flex items-center gap-[8rpx]'>
                    {/* 用户按钮 */}
                    <View
                      className='flex items-center justify-center w-[64rpx] h-[64rpx] rounded'
                      onClick={() => router.push('/user')}
                      hoverClass='tap-active'
                    >
                        <XIcon name='carbon:user' size={20} className='text-textSec' />
                    </View>

                    {/* 汉堡菜单 */}
                    <MobileNav />
                </View>
            </View>
        </View>
    );
}
