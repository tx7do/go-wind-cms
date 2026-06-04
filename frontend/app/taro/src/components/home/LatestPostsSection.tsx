import {View, Text} from '@tarojs/components';
import {useTranslations} from '@/lib/next-intl-compat';
import {XIcon} from '@/plugins/xicon';
import {useI18nRouter} from '@/i18n/helpers/useI18nRouter';
import PostList from '@/components/post/PostList';

/**
 * 首页「最新文章」区块
 * - 默认竖排卡片（封面+标题+摘要）
 * - 显示 4 篇
 */
export default function LatestPostsSection() {
    const t = useTranslations('page.home');
    const router = useI18nRouter();

    return (
        <View className='w-full'>
            {/* 标题行 */}
            <View className='flex items-center justify-between mb-[32rpx]'>
                <View className='flex items-center gap-[8rpx]'>
                    <XIcon name='carbon:document' size={20} className='text-primary' />
                    <Text className='text-card-title font-bold text-textMain'>
                        {t('latest_posts')}
                    </Text>
                </View>
                <View
                  className='px-[16rpx] py-[8rpx] min-w-touch min-h-touch flex items-center justify-center'
                  onClick={() => router.push('/post')}
                  hoverClass='tap-active'
                >
                    <Text className='text-desc text-primary'>{t('view_all')} →</Text>
                </View>
            </View>

            {/* 文章列表：默认竖排卡片 */}
            <PostList
              queryParams={{status: 'POST_STATUS_PUBLISHED'}}
              fieldMask='id,status,sortOrder,isFeatured,visits,likes,commentCount,authorName,availableLanguages,createdAt,translations.id,translations.postId,translations.languageCode,translations.title,translations.summary,translations.thumbnail'
              orderBy={['-createdAt']}
              page={1}
              pageSize={4}
              showSkeleton
              from='home'
              showPagination={false}
            />
        </View>
    );
}
