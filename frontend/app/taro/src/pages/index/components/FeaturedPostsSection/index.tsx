import {useTranslation} from 'react-i18next';
import {View, Text} from '@tarojs/components';
import {useMemo} from 'react';

import PostList from '@/components/post/PostList';

import './index.scss';
import {useI18nRouter} from "@/i18n/helpers";

export default function FeaturedPostsSection() {
  const {t} = useTranslation();
  // 用 useMemo 缓存 queryParams 和 orderBy，避免重复请求
  const queryParams = useMemo(() => ({status: 'POST_STATUS_PUBLISHED', isFeatured: true}), []);
  const orderBy = useMemo(() => ['-sortOrder'], []);

  const router = useI18nRouter();

  return (
    <View className='featured-section scroll-reveal'>
      <View className='section-header'>
        <Text className='section-title'>
          ⭐ {t('page.home.featured_posts')}
        </Text>
        <View className='view-all-btn' onClick={() => router.push('/post')}>
          <Text>{t('page.home.view_all')} →</Text>
        </View>
      </View>
      <View className='featured-grid'>
        <PostList
          queryParams={queryParams}
          fieldMask='id,status,sortOrder,isFeatured,visits,likes,commentCount,authorName,availableLanguages,createdAt,translations.id,translations.postId,translations.languageCode,translations.title,translations.summary,translations.thumbnail'
          orderBy={orderBy}
          page={1}
          pageSize={3}
          columns={1}
          showSkeleton={true}
          from='home'
          showPagination={false}
        />
      </View>
    </View>
  );
}
