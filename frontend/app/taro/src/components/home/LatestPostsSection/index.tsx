import {useTranslation} from 'react-i18next';
import {View, Text} from '@tarojs/components';

import PostList from '@/components/post/PostList';

import './index.scss';

export default function LatestPostsSection() {
  const {t} = useTranslation('page.home');

  return (
    <View className='latest-section scroll-reveal'>
      <View className='section-header'>
        <Text className='section-title'>
          📄 {t('latest_posts')}
        </Text>
        <View className='view-all-btn' onClick={() => console.log('Navigate to /post')}>
          <Text>{t('view_all')} →</Text>
        </View>
      </View>
      <View className='latest-grid'>
        <PostList
          queryParams={{status: 'POST_STATUS_PUBLISHED'}}
          fieldMask='id,status,sortOrder,isFeatured,visits,likes,commentCount,authorName,availableLanguages,createdAt,translations.id,translations.postId,translations.languageCode,translations.title,translations.summary,translations.thumbnail'
          orderBy={['-createdAt']}
          page={1}
          pageSize={6}
          columns={3}
          showSkeleton={true}
          from='home'
          showPagination={false}
        />
      </View>
    </View>
  );
}
