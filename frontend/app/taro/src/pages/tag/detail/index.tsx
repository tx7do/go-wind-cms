import {useState, useEffect, useMemo} from 'react';
import {useTranslation} from 'react-i18next';
import {View, Text} from '@tarojs/components';
import Taro from '@tarojs/taro';

import {AppEmpty} from '@/components/ui';
import XIcon from '@/plugins/xicon';
import PostListWithPagination from '@/components/post/PostList';
import {useI18nRouter} from "@/i18n/helpers";

import {useTagStore} from '@/store/slices/tag/hooks';
import type {contentservicev1_Tag} from "@/api/generated/app/service/v1";

import './tag-detail.scss';

export default function TagDetailPage() {
  const {t} = useTranslation();
  const router = useI18nRouter();

  const tagStore = useTagStore();
  const [tag, setTag] = useState<contentservicev1_Tag | null>(null);

  const [_loading, setLoading] = useState(false);
  const tagId = useMemo(() => {
    const pages = Taro.getCurrentPages();
    const currentPage = pages[pages.length - 1];
    const id = (currentPage as any).options?.id;
    return id ? parseInt(id) : null;
  }, []);

  async function loadTag() {
    if (!tagId) return;

    setLoading(true);
    try {
      const tagData = await tagStore.getTag({
        id: tagId
      }) as unknown as contentservicev1_Tag;
      setTag(tagData);
    } catch (error) {
      console.error('Load tag failed:', error);
    } finally {
      setLoading(false);
    }
  }

  function handleBack() {
    router.push(`/tag`);
  }

  useEffect(() => {
    loadTag();
  }, [tagId]);

  if (!tagId) {
    return <AppEmpty variant='error' description='Invalid tag ID' />;
  }

  const translation = tagStore.getTranslation(tag);

  return (
    <View className='tag-detail-page'>
      {/* Hero Section */}
      <View className='hero-section'>
        <View className='hero-content'>
          <Text className='tag-title'>{translation?.name || t('page.tags.tag_untitled')}</Text>
          {translation?.description && (
            <Text className='tag-description'>
              {translation.description}
            </Text>
          )}
          <View className='tag-stats'>
            <View className='stat-item'>
              <XIcon name='carbon:document' size={20} />
              <Text>{tag?.postCount || 0} {t('page.posts.articles')}</Text>
            </View>
          </View>
        </View>
      </View>

      {/* Posts Section */}
      <View className='page-container'>
        {/* Back Button */}
        <View className='back-button-container'>
          <View className='back-btn' onClick={handleBack}>
            <Text>← {t('page.categories.back_to_list')}</Text>
          </View>
        </View>

        {/* Posts List with Pagination */}
        {tagId && (
          <PostListWithPagination
            key={tagId}
            initialPageSize={10}
            pageSizes={[10, 20, 30, 40]}
            tagId={tagId}
            from='tag'
          />
        )}
      </View>
    </View>
  );
}
