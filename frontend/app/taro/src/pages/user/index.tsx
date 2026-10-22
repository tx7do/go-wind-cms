import {View, Text, Image} from '@tarojs/components';
import {useState} from 'react';
import {useDidShow} from '@tarojs/taro';
import {useTranslation} from 'react-i18next';
import XIcon from '@/plugins/xicon';
import {useI18nRouter} from '@/i18n/helpers';
import {usePageTitle} from '@/hooks/usePageTitle';
import {useAuth} from '@/api/hooks/auth';
import {fetchUserProfile} from '@/api/hooks/user-profile';
import {useAccessStore, useIsLogin} from '@/store/core/access/store';
import type {identityservicev1_User} from '@/api/generated/app/service/v1';

/**
 * 个人中心（tabBar「我的」）
 * 布局：渐变横幅（头像/昵称/签名）→ 资料卡 → 通用菜单（设置/改密/关于/法务）→ 退出登录。
 * 统计卡不再展示假数据：只展示真实可得的账号信息（用户名/角色/语言）。
 */
export default function UserProfilePage() {
    const {t} = useTranslation();
    usePageTitle('page.title.user');
    const router = useI18nRouter();
    const {logout} = useAuth();
    const isLogin = useIsLogin();
    const accessToken = useAccessStore((s) => s.accessToken);
    const [user, setUser] = useState<identityservicev1_User | null>(null);
    const [loading, setLoading] = useState(false);

    useDidShow(() => {
        if (!isLogin) {
            setUser(null);
            return;
        }
        setLoading(true);
        fetchUserProfile()
            .then((result) => setUser((result as identityservicev1_User) || null))
            .catch((error) => {
                console.error('Load user profile failed:', error);
                setUser(null);
            })
            .finally(() => setLoading(false));
    });

    // —— 未登录：登录引导 ——
    if (!isLogin) {
        return (
            <View className='au-page w-full bg-pageBg min-h-screen flex flex-col items-center justify-center px-[48rpx]'>
                <View className='au-card au-card--glow w-full flex flex-col items-center gap-[28rpx] py-[72rpx] px-[40rpx]'>
                    <View className='w-[112rpx] h-[112rpx] rounded-full bg-primary-tint flex items-center justify-center'>
                        <XIcon name='carbon:user-avatar' size={52} className='text-primary' />
                    </View>
                    <Text className='text-body text-textThird'>{t('page.user.please_login')}</Text>
                    <View
                      className='btn-hero px-[72rpx] h-[88rpx] rounded-full flex items-center justify-center'
                      hoverClass='tap-active'
                      onClick={() => router.push('/login')}
                    >
                        <Text className='text-body text-white font-bold'>{t('authentication.login.login')}</Text>
                    </View>
                </View>
            </View>
        );
    }

    const menuItems = [
        {icon: 'carbon:settings', label: t('page.title.settings'), path: '/settings'},
        {icon: 'carbon:password', label: t('settings.account.password'), path: '/settings/password'},
        {icon: 'carbon:star', label: t('page.title.favorites'), path: '/favorites'},
        {icon: 'carbon:information', label: t('page.title.about'), path: '/about'},
        {icon: 'carbon:email', label: t('page.title.contact'), path: '/contact'},
        {icon: 'carbon:document', label: t('page.title.terms'), path: '/terms'},
        {icon: 'carbon:shield-alert', label: t('page.title.privacy'), path: '/privacy'},
    ];

    return (
        <View className='au-page w-full bg-pageBg min-h-screen pb-[48rpx]'>
            {/* 资料头部：渐变横幅 + 悬浮头像 */}
            <View
              className='relative overflow-hidden px-[32rpx] pt-[44rpx] pb-[72rpx]'
              style={{background: 'linear-gradient(128deg, #57a3ff 0%, #1677ff 52%, #0b55d4 100%)'}}
            >
                <View
                  className='absolute rounded-full'
                  style={{
                      width: '360rpx',
                      height: '360rpx',
                      right: '-110rpx',
                      top: '-150rpx',
                      background: 'radial-gradient(circle, rgba(255,255,255,0.26), transparent 68%)',
                  }}
                />
                <View className='relative flex items-center gap-[28rpx]'>
                    <View className='w-[120rpx] h-[120rpx] rounded-full bg-white/18 border-[length:2rpx] border-white/35 flex items-center justify-center overflow-hidden shadow-[0_12rpx_32rpx_rgba(3,20,60,0.35)]'>
                        {user?.avatar ? (
                            <Image src={user.avatar} className='w-full h-full rounded-full' mode='aspectFill' />
                        ) : (
                            <Text className='text-title font-bold text-white'>
                                {user?.nickname?.charAt(0) || user?.username?.charAt(0) || 'U'}
                            </Text>
                        )}
                    </View>
                    <View className='flex-1 min-w-0'>
                        <Text className='text-card-title font-bold text-white block' numberOfLines={1}>
                            {user?.nickname || user?.username || t('page.user.not_login')}
                        </Text>
                        {user?.username && (
                            <Text className='text-tips text-white/80 block mt-[6rpx]'>@{user.username}</Text>
                        )}
                    </View>
                    <View
                      className='px-[22rpx] h-[56rpx] rounded-full flex items-center border-[length:1rpx] border-white/40 bg-white/12'
                      hoverClass='tap-active'
                      onClick={() => router.push('/settings/profile')}
                    >
                        <Text className='text-tips text-white'>{t('page.user.edit_profile')}</Text>
                    </View>
                </View>
            </View>

            <View className='px-[24rpx] -mt-[44rpx]'>
                {/* 账号信息卡（只展示真实字段） */}
                <View className='au-card au-card--glow p-[32rpx] au-anim'>
                    <View className='flex justify-around'>
                        <View className='items-center flex flex-col gap-[8rpx] min-w-0'>
                            <Text className='text-desc font-bold text-textMain' numberOfLines={1}>
                                {user?.nickname || user?.username || '-'}
                            </Text>
                            <Text className='text-tips text-textThird'>{t('page.user.nickname')}</Text>
                        </View>
                        <View className='w-[1rpx] bg-splitLine' />
                        <View className='items-center flex flex-col gap-[8rpx] min-w-0'>
                            <Text className='text-desc font-bold text-textMain' numberOfLines={1}>
                                {user?.email || user?.mobile || '-'}
                            </Text>
                            <Text className='text-tips text-textThird'>{t('page.user.contact_info')}</Text>
                        </View>
                        <View className='w-[1rpx] bg-splitLine' />
                        <View className='items-center flex flex-col gap-[8rpx] min-w-0'>
                            <Text className='text-desc font-bold text-textMain' numberOfLines={1}>
                                {user?.roles?.length ? user.roles.length : 0}
                            </Text>
                            <Text className='text-tips text-textThird'>{t('page.user.roles')}</Text>
                        </View>
                    </View>
                </View>

                {/* 通用菜单 */}
                <View className='au-card overflow-hidden mt-[24rpx] au-anim au-anim-1'>
                    {menuItems.map((item, index) => (
                        <View
                          key={item.path}
                          className={`flex items-center justify-between px-[28rpx] h-[104rpx] ${index < menuItems.length - 1 ? 'border-b-[length:1rpx] border-splitLine' : ''}`}
                          hoverClass='tap-active'
                          onClick={() => router.push(item.path)}
                        >
                            <View className='flex items-center gap-[20rpx]'>
                                <View className='w-[60rpx] h-[60rpx] rounded-[16rpx] bg-primary-tint flex items-center justify-center'>
                                    <XIcon name={item.icon} size={26} className='text-primary' />
                                </View>
                                <Text className='text-body text-textMain'>{item.label}</Text>
                            </View>
                            <XIcon name='carbon:chevron-right' size={16} className='text-textWeak' />
                        </View>
                    ))}
                </View>

                {/* 退出登录 */}
                <View
                  className='au-card mt-[24rpx] h-[96rpx] flex items-center justify-center au-anim au-anim-2'
                  hoverClass='tap-active'
                  onClick={() => logout()}
                >
                    <Text className='text-body text-danger'>{t('menu.logout')}</Text>
                </View>

                {loading && (
                    <View className='py-[24rpx] text-center'>
                        <Text className='text-tips text-textThird'>{t('page.user.loading')}</Text>
                    </View>
                )}

                {/* 令牌存在但资料为空：静默加载失败的兜底提示（不弹窗打断） */}
                {!loading && !user && accessToken?.value && (
                    <View className='py-[16rpx] text-center'>
                        <Text className='text-tips text-textWeak'>{t('page.user.load_failed')}</Text>
                    </View>
                )}
            </View>
        </View>
    );
}
