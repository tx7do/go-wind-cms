import {useState, useEffect, useMemo} from 'react';
import {useTranslation} from 'react-i18next';
import {View, Text, Image} from '@tarojs/components';
import XIcon from '@/plugins/xicon';
import {formatDateTime} from '@/utils';
import {type contentservicev1_ListPostResponse, type contentservicev1_Post, identityservicev1_User} from '@/api/generated/app/service/v1';
import {fetchListPosts, getPostTitle, getPostSummary} from '@/api/hooks/post';
import {fetchUserProfile} from '@/api/hooks/user-profile';
import {useI18nRouter} from '@/i18n/helpers';

export default function UserProfilePage() {
    const {t} = useTranslation();
    const router = useI18nRouter();
    const [loading, setLoading] = useState(false);
    const [user, setUser] = useState<identityservicev1_User | null>(null);
    const [posts, setPosts] = useState<contentservicev1_Post[]>([]);
    const [activeTab, setActiveTab] = useState<'posts' | 'activities' | 'collections'>('posts');

    const stats = useMemo(() => ({
        followers: user?.followers ?? 0,
        following: user?.following ?? 0,
        posts: user?.postCount ?? 0,
        likes: user?.likeCount ?? 0,
    }), [user]);

    useEffect(() => {
        (async () => {
            setLoading(true);
            try {
                const result = await fetchUserProfile() as identityservicev1_User;
                setUser(result || null);
            } catch (error) {
                console.error('Load user profile failed:', error);
            } finally {
                setLoading(false);
            }
        })();
    }, []);

    if (loading) {
        return (
            <View className='w-full bg-pageBg min-h-screen flex items-center justify-center'>
                <Text className='text-body text-textThird'>{t('page.user.loading') || 'Loading...'}</Text>
            </View>
        );
    }

    if (!user) {
        return (
            <View className='w-full bg-pageBg min-h-screen flex items-center justify-center'>
                <Text className='text-body text-textThird'>{t('page.user.please_login')}</Text>
            </View>
        );
    }

    const tabs = [
        {key: 'posts' as const, label: t('page.user.tab_posts')},
        {key: 'activities' as const, label: t('page.user.tab_activities')},
        {key: 'collections' as const, label: t('page.user.tab_collections')},
    ];

    return (
        <View className='w-full bg-pageBg'>
            {/* 用户信息头部 */}
            <View className='bg-cardBg px-[24rpx] py-[32rpx]'>
                <View className='flex items-center gap-[24rpx] mb-[24rpx]'>
                    <View className='w-[96rpx] h-[96rpx] rounded-full bg-primary/10 flex items-center justify-center'>
                        {user?.avatar ? (
                            <Image src={user.avatar} className='w-full h-full rounded-full' mode='aspectFill' />
                        ) : (
                            <Text className='text-title font-bold text-primary'>
                                {user?.nickname?.charAt(0) || user?.username?.charAt(0) || 'U'}
                            </Text>
                        )}
                    </View>
                    <View className='flex-1'>
                        <Text className='text-card-title font-bold text-textMain block'>{user?.nickname || user?.username}</Text>
                        {user?.description && (
                            <Text className='text-desc text-textSec block mt-[8rpx]'>{user.description}</Text>
                        )}
                    </View>
                    <View
                        className='px-[24rpx] h-[56rpx] rounded flex items-center border-[1rpx] border-primary'
                        hoverClass='tap-active'
                        onClick={() => router.push('/settings')}
                    >
                        <Text className='text-tips text-primary'>{t('page.user.edit_profile')}</Text>
                    </View>
                </View>

                {/* 统计数据 */}
                <View className='grid grid-cols-4 gap-[16rpx]'>
                    {[
                        {value: stats.following, label: t('page.user.following')},
                        {value: stats.followers, label: t('page.user.followers')},
                        {value: stats.posts, label: t('page.user.posts')},
                        {value: stats.likes, label: t('page.user.likes_received')},
                    ].map((stat) => (
                        <View key={stat.label} className='text-center'>
                            <Text className='text-card-title font-bold text-textMain block'>{stat.value}</Text>
                            <Text className='text-tips text-textThird block'>{stat.label}</Text>
                        </View>
                    ))}
                </View>
            </View>

            {/* Tab 切换 */}
            <View className='flex bg-cardBg border-b-[1rpx] border-splitLine mt-[16rpx]'>
                {tabs.map((tab) => (
                    <View
                        key={tab.key}
                        className={`flex-1 text-center py-[24rpx] min-h-touch ${activeTab === tab.key ? 'border-b-[4rpx] border-primary' : ''}`}
                        onClick={() => setActiveTab(tab.key)}
                    >
                        <Text className={activeTab === tab.key ? 'text-body text-primary font-bold' : 'text-body text-textSec'}>
                            {tab.label}
                        </Text>
                    </View>
                ))}
            </View>

            {/* Tab 内容 */}
            <View className='px-[24rpx] py-[24rpx]'>
                {activeTab === 'posts' && (
                    <View className='text-center py-[64rpx]'>
                        <XIcon name='carbon:document' size={48} className='text-textWeak' />
                        <Text className='text-body text-textThird block mt-[16rpx]'>{t('page.user.no_posts')}</Text>
                    </View>
                )}
                {activeTab === 'activities' && (
                    <View className='text-center py-[64rpx]'>
                        <XIcon name='carbon:target' size={48} className='text-textWeak' />
                        <Text className='text-body text-textThird block mt-[16rpx]'>{t('page.user.no_activities')}</Text>
                    </View>
                )}
                {activeTab === 'collections' && (
                    <View className='text-center py-[64rpx]'>
                        <XIcon name='carbon:bookmark' size={48} className='text-textWeak' />
                        <Text className='text-body text-textThird block mt-[16rpx]'>{t('page.user.no_collections')}</Text>
                    </View>
                )}
            </View>
        </View>
    );
}
