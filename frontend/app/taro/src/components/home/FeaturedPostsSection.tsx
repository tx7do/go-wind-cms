import {View, Text} from '@tarojs/components';
import {useTranslations} from '@/lib/next-intl-compat';
import {XIcon} from '@/plugins/xicon';
import PostList from '@/components/post/PostList';
import {useI18nRouter} from '@/i18n/helpers/useI18nRouter';

/**
 * 首页「精选文章」区块
 * - 精选文章使用 compact 模式（横排卡片）节省空间
 * - 显示 3 篇
 */
export default function FeaturedPostsSection() {
    const t = useTranslations('page.home');
    const router = useI18nRouter();

    return (
        <View className='w-full'>
            {/* 标题行 */}
            <View className='flex items-center justify-between mb-[32rpx]'>
                <View className='flex items-center gap-[8rpx]'>
                    <XIcon name='carbon:star-filled' size={20} className='text-primary' />
                    <Text className='text-card-title font-bold text-textMain'>
                        {t('featured_posts')}
                    </Text>
                </View>
                <View
                  className='px-[24rpx] py-[12rpx] rounded-full flex items-center justify-center'
                  style={{
                      backgroundColor: 'rgba(22,119,255,0.08)',
                  }}
                  onClick={() => router.push('/post')}
                  hoverClass='tap-active'
                >
                    <Text className='text-tips font-medium text-primary'>{t('view_all')} →</Text>
                </View>
            </View>

            {/* 文章列表：compact 横排卡片 */}
            <PostList
              queryParams={{status: 'POST_STATUS_PUBLISHED', isFeatured: true}}
              fieldMask='id,status,sortOrder,isFeatured,visits,likes,commentCount,authorName,availableLanguages,createdAt,translations.id,translations.postId,translations.languageCode,translations.title,translations.summary,translations.thumbnail'
              orderBy={['-sortOrder']}
              page={1}
              pageSize={3}
              showSkeleton
              from='home'
              showPagination={false}
              compact
            />
        </View>
    );
}
