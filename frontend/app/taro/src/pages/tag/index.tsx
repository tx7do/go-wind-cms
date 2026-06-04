import {useState, useEffect} from 'react';
import {useTranslation} from 'react-i18next';
import {View, Text} from '@tarojs/components';

import {AppEmpty} from '@/components/ui';
import {Skeleton} from '@/components/ui/skeleton';
import {Button} from '@/components/ui/button';
import XIcon from '@/plugins/xicon';
import {usePageTitle} from '@/hooks/usePageTitle';

import {fetchListTags, getTagTranslation} from '@/api/hooks/tag';
import {contentservicev1_ListTagResponse, contentservicev1_Tag} from '@/api/generated/app/service/v1';
import {useTranslations} from '@/lib/next-intl-compat';
import {useI18nRouter} from '@/i18n/helpers';

export default function TagListPage() {
  const {t} = useTranslation();
  usePageTitle('page.title.tags');
  const pt = useTranslations('page.tags');
  const router = useI18nRouter();

  const [loading, setLoading] = useState(false);
  const [tags, setTags] = useState<contentservicev1_Tag[]>([]);
  const [page, setPage] = useState(1);
  const [pageSize] = useState(20);
  const [total, setTotal] = useState(0);

  async function loadTags() {
    setLoading(true);
    try {
      const res = await fetchListTags({
        paging: {page, pageSize},
        formValues: {status: 'TAG_STATUS_ACTIVE'},
        fieldMask: undefined,
        orderBy: undefined,
      }) as contentservicev1_ListTagResponse;
      setTags(res.items || []);
      setTotal(res.total || 0);
    } catch (error) {
      console.error('Load tags failed:', error);
    } finally {
      setLoading(false);
    }
  }

  function handleTagClick(id: number) {
    router.push(`/tag/detail?id=${id}`);
  }

  useEffect(() => {
    loadTags();
  }, [page, pageSize]);

  const totalPages = Math.ceil(total / pageSize);

  return (
    <View className='min-h-screen w-full bg-pageBg pb-[240rpx]'>
      {/* 页面标题 - px-[24rpx] 与 Header/卡片网格对齐 */}
      <View className='bg-cardBg px-[24rpx] pt-[40rpx] pb-[24rpx] border-b-[1rpx] border-splitLine'>
        <Text className='text-title font-bold text-textMain'>{t('page.tags.tags_list')}</Text>
        <Text className='text-desc text-textSec block mt-[8rpx]'>{t('page.tags.explore_all')}</Text>
        <View className='flex items-center gap-[12rpx] mt-[16rpx]'>
          <XIcon name='carbon:tag' size={14} className='text-textSec' />
          <Text className='text-tips text-textSec'>{total} {t('page.tags.total_tags')}</Text>
        </View>
      </View>

      {/* 标签网格 */}
      <View className='px-[24rpx] pt-[24rpx]'>
        {loading ? (
          <View className='grid grid-cols-2 gap-[24rpx]'>
            {Array.from({length: 6}).map((_, i) => (
              <View key={i} className='rounded-[16rpx] bg-cardBg overflow-hidden'>
                <Skeleton className='w-full h-[200rpx]' />
              </View>
            ))}
          </View>
        ) : (
          <>
            {tags.length > 0 && (
              <View className='grid grid-cols-2 gap-[24rpx]'>
                {tags.map((tag) => {
                  const tagName = getTagTranslation(tag)?.name || t('page.tags.tag_untitled');
                  const tagDesc = getTagTranslation(tag)?.description || '';

                  return (
                    <View
                      key={tag.id}
                      className='flex flex-col rounded-[16rpx] bg-cardBg overflow-hidden tap-active'
                      style={{
                        borderTop: `3px solid ${tag.color || '#1677ff'}`,
                      }}
                      onClick={() => handleTagClick(tag.id || 0)}
                      hoverClass='opacity-80'
                    >
                      <View className='flex flex-col p-[24rpx]'>
                        {/* 标题行 */}
                        <View className='flex items-center gap-[12rpx]'>
                          <View
                            className='w-[40rpx] h-[40rpx] rounded-full flex items-center justify-center flex-shrink-0'
                            style={{backgroundColor: `${tag.color || '#1677ff'}20`}}
                          >
                            <XIcon name='carbon:tag' size={16} style={{color: tag.color || '#1677ff'}} />
                          </View>
                          <Text
                            className='text-body font-bold text-textMain flex-1 min-w-0'
                            style={{
                              overflow: 'hidden',
                              textOverflow: 'ellipsis',
                              whiteSpace: 'nowrap',
                            }}
                          >
                            {tagName}
                          </Text>
                        </View>

                        {/* 描述文字 - 使用较亮的颜色提升暗色模式可读性 */}
                        {tagDesc && (
                          <Text
                            className='text-tips leading-[1.6] mt-[12rpx] mb-[16rpx]'
                            style={{
                              color: 'var(--color-text-sec)',
                              overflow: 'hidden',
                              textOverflow: 'ellipsis',
                              display: '-webkit-box',
                              WebkitLineClamp: 2,
                              WebkitBoxOrient: 'vertical',
                            }}
                          >
                            {tagDesc}
                          </Text>
                        )}

                        {/* 底部篇数 - 图标与文字间距充足 */}
                        <View className='flex items-center gap-[12rpx]'>
                          <XIcon name='carbon:document' size={12} className='text-textSec' />
                          <Text className='text-tips text-textSec'>
                            {tag.postCount || 0} {t('page.posts.articles')}
                          </Text>
                        </View>
                      </View>
                    </View>
                  );
                })}
              </View>
            )}

            {!loading && tags.length === 0 && total === 0 && (
              <AppEmpty description={t('page.tags.no_tags')} />
            )}

            {/* 分页 */}
            {totalPages > 1 && (
              <View className='flex items-center justify-center gap-[24rpx] py-[48rpx]'>
                <Button
                  variant='outline'
                  size='sm'
                  disabled={page <= 1}
                  onClick={() => setPage(page - 1)}
                >
                  {pt('previous') || '上一页'}
                </Button>
                <Text className='text-desc text-textSec'>
                  {page} / {totalPages}
                </Text>
                <Button
                  variant='outline'
                  size='sm'
                  disabled={page >= totalPages}
                  onClick={() => setPage(page + 1)}
                >
                  {pt('next') || '下一页'}
                </Button>
              </View>
            )}
          </>
        )}
      </View>
    </View>
  );
}
