import {useTranslation} from 'react-i18next';
import {View, Text} from '@tarojs/components';
import XIcon from '@/plugins/xicon';

import PostList from '@/components/post/PostList';
import {useI18nRouter} from "@/i18n/helpers";

import './index.scss';

export default function LatestPostsSection() {
  const {t} = useTranslation();

  const router = useI18nRouter();

  return (
    <View className='latest-section scroll-reveal'>
      <View className='section-header'>
        <Text className='section-title'>
          <XIcon name='carbon:document' size={24} /> {t('page.home.latest_posts')}
        </Text>
        <View className='view-all-btn' onClick={() => router.push('/post')}>
          <Text>{t('page.home.view_all')} →</Text>
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
