import {useState, useEffect, useMemo} from 'react';
import {useTranslation} from 'react-i18next';
import {View, Text} from '@tarojs/components';
import Taro from '@tarojs/taro';

import {AppEmpty} from '@/components/ui';
import XIcon from '@/plugins/xicon';
import PostListWithPagination from '@/components/post/PostList';
import {useI18nRouter} from '@/i18n/helpers';

import {fetchTag, getTagTranslation} from '@/api/hooks/tag';
import type {contentservicev1_Tag} from '@/api/generated/app/service/v1';

export default function TagDetailPage() {
  const {t} = useTranslation();
  const router = useI18nRouter();

  const [tag, setTag] = useState<contentservicev1_Tag | null>(null);
  const [_loading, setLoading] = useState(false);

  const tagId = useMemo(() => {
    let id: string | null = null;
    if (typeof Taro.getCurrentInstance === 'function') {
      const instance = Taro.getCurrentInstance();
      const routeId = instance?.router?.params?.id;
      id = typeof routeId === 'string' ? routeId : null;
    } else {
      const pages = Taro.getCurrentPages();
      const currentPage = pages[pages.length - 1];
      const pageId = (currentPage as any).options?.id;
      id = typeof pageId === 'string' ? pageId : null;
    }
    return id ? parseInt(id) : null;
  }, []);

  async function loadTag() {
    if (!tagId) return;
    setLoading(true);
    try {
      const tagData = await fetchTag(tagId) as contentservicev1_Tag;
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

  const translation = getTagTranslation(tag);

  return (
    <View className='min-h-screen w-full bg-pageBg pb-[160rpx]'>
      {/* 页面标题 */}
      <View
        className='bg-cardBg px-[24rpx] pt-[40rpx] pb-[24rpx] border-b-[1rpx] border-splitLine'
        style={{borderTop: `3px solid ${tag?.color || '#1677ff'}`}}
      >
        <View className='flex items-center gap-[16rpx]'>
          <View
            className='w-[48rpx] h-[48rpx] rounded-full flex items-center justify-center'
            style={{backgroundColor: `${tag?.color || '#1677ff'}20`}}
          >
            <XIcon name='carbon:tag' size={20} style={{color: tag?.color || '#1677ff'}} />
          </View>
          <View className='flex-1 min-w-0'>
            <Text className='text-title font-bold text-textMain'>
              {translation?.name || t('page.tags.tag_untitled')}
            </Text>
            {translation?.description && (
              <Text className='text-desc text-textSec block mt-[8rpx]'>
                {translation.description}
              </Text>
            )}
          </View>
        </View>
        <View className='flex items-center gap-[16rpx] mt-[16rpx] text-tips text-textThird'>
          <View className='flex items-center gap-[6rpx]'>
            <XIcon name='carbon:document' size={14} className='text-textThird' />
            <Text className='text-tips text-textThird'>
              {tag?.postCount || 0} {t('page.posts.articles')}
            </Text>
          </View>
        </View>
      </View>

      {/* 返回按钮 */}
      <View className='px-[24rpx] pt-[24rpx]'>
        <View
          className='inline-flex items-center gap-[8rpx] rounded-full bg-cardBg px-[24rpx] py-[12rpx] border-[1rpx] border-splitLine tap-active'
          onClick={handleBack}
          hoverClass='opacity-80'
        >
          <XIcon name='carbon:arrow-left' size={14} className='text-textSec' />
          <Text className='text-desc text-textSec'>{t('page.categories.back_to_list')}</Text>
        </View>
      </View>

      {/* 文章列表 */}
      <View className='px-[24rpx] pt-[24rpx]'>
        {tagId && (
          <PostListWithPagination
            key={tagId}
            initialPageSize={10}
            pageSizes={[10, 20, 30, 40]}
            tagId={tagId}
            from='tag'
            showPagination={true}
          />
        )}
      </View>
    </View>
  );
}
