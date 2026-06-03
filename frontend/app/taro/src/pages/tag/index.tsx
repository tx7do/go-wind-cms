import {useState, useEffect} from 'react';
import {useTranslation} from 'react-i18next';
import {View, Text} from '@tarojs/components';

import {AppEmpty} from '@/components/ui';
import XIcon from '@/plugins/xicon';
import Pagination from '@/components/Pagination';
import {useI18nRouter} from "@/i18n/helpers";

import {fetchListTags, getTagTranslation} from '@/api/hooks/tag';
import {contentservicev1_ListTagResponse, contentservicev1_Tag} from '@/api/generated/app/service/v1';

import './tag-list.scss';

export default function TagListPage() {
  const {t} = useTranslation();

  const router = useI18nRouter();

  const [loading, setLoading] = useState(false);
  const [tags, setTags] = useState<contentservicev1_Tag[]>([]);
  const [page, setPage] = useState(1);
  const [pageSize, _setPageSize] = useState(20);
  const [total, setTotal] = useState(0);

  async function loadTags() {
    setLoading(true);
    try {
      const res = await fetchListTags({
        paging: {
          page: page,
          pageSize: pageSize,
        },
        formValues: {
          status: 'TAG_STATUS_ACTIVE'
        },
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
    router.push(`/tag/${id}`);
  }

  function handlePageChange(newPage: number) {
    setPage(newPage);
    loadTags();
  }

  useEffect(() => {
    loadTags();
  }, [page, pageSize]);

  useEffect(() => {
    if (page !== 1) {
      setPage(1);
    }
  }, [pageSize]);

  return (
    <View className='tag-list-page'>
      {/* Hero Section */}
      <View className='hero-section'>
        <View className='hero-content'>
          <Text className='hero-title'>{t('page.tags.tags_list')}</Text>
          <Text className='hero-subtitle'>{t('page.tags.explore_all')}</Text>
          <View className='tag-stats'>
            <View className='stat-item'>
              <XIcon name='carbon:tag' size={20} />
              <Text>{total} {t('page.tags.total_tags')}</Text>
            </View>
          </View>
        </View>
      </View>

      {/* Tags Grid */}
      <View className='page-container'>
        {/* Loading Skeleton - TODO: 实现 Taro 骨架屏 */}
        {loading ? (
          <View className='tags-grid'>
            <Text>{t('common.loading')}</Text>
          </View>
        ) : (
          <>
            {tags.length > 0 && (
              <View className='tags-grid'>
                {tags.map((tag) => (
                  <View
                    key={tag.id}
                    className='tag-card'
                    style={{
                      borderColor: tag.color || 'rgba(0, 0, 0, 0.08)',
                      background: `linear-gradient(135deg, ${tag.color || '#f0f0f0'}15 0%, #ffffff 100%)`
                    }}
                    onClick={() => handleTagClick(tag.id || 0)}
                  >
                    <View className='tag-content'>
                      <Text>{getTagTranslation(tag)?.name || t('page.tags.tag_untitled')}</Text>
                      <Text className='tag-description'>
                        {getTagTranslation(tag)?.description || ''}
                      </Text>
                      <View className='tag-meta'>
                        <XIcon name='carbon:document' size={16} />
                        <Text>
                          {tag.postCount || 0} {t('page.posts.articles')}
                        </Text>
                      </View>
                    </View>
                  </View>
                ))}
                {tags.length === 0 && total > 0 && (
                  <AppEmpty description={t('page.tags.no_tags_in_page')} />
                )}
              </View>
            )}

            {!loading && tags.length === 0 && total === 0 && (
              <AppEmpty inContainer description={t('page.tags.no_tags')} />
            )}

            {total > pageSize && (
              <View className='pagination-wrapper'>
                <Pagination
                  current={page}
                  total={total}
                  pageSize={pageSize}
                  onChange={handlePageChange}
                  showSizeChanger
                />
              </View>
            )}
          </>
        )}
      </View>
    </View>
  );
}
