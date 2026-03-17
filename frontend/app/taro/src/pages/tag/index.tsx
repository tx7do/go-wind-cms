import {useState, useEffect} from 'react';
import {useTranslation} from 'react-i18next';
import {View, Text} from '@tarojs/components';
import Taro from '@tarojs/taro';
import {AppEmpty} from '@/components/ui';

import {useTagStore} from '@/store/slices/tag/hooks';
import {contentservicev1_ListTagResponse, contentservicev1_Tag} from '@/api/generated/app/service/v1';

import '../../globals.css'; // 导入全局 CSS，确保 CSS 变量可用
import styles from './tag-list.scss';

export default function TagListPage() {
  const {t} = useTranslation('page');
  const tagStore = useTagStore();

  const [loading, setLoading] = useState(false);
  const [tags, setTags] = useState<contentservicev1_Tag[]>([]);
  const [page, setPage] = useState(1);
  const [pageSize, setPageSize] = useState(20);
  const [total, setTotal] = useState(0);

  async function loadTags() {
    setLoading(true);
    try {
      const res = await tagStore.listTag({
        paging: {
          page: page,
          pageSize: pageSize,
        },
        formValues: {
          status: 'TAG_STATUS_ACTIVE'
        },
        fieldMask: null,
        orderBy: null,
      }) as unknown as contentservicev1_ListTagResponse;
      setTags(res.items || []);
      setTotal(res.total || 0);
    } catch (error) {
      console.error('Load tags failed:', error);
    } finally {
      setLoading(false);
    }
  }

  function handleTagClick(id: number) {
    Taro.navigateTo({url: `/pages/[locale]/tag/${id}/page`});
  }

  function handlePageChange(newPage: number) {
    setPage(newPage);
    loadTags();
  }

  function handlePageSizeChange(newSize: number) {
    setPageSize(newSize);
    setPage(1);
    loadTags();
  }

  useEffect(() => {
    loadTags();
  }, [page, pageSize]);

  return (
    <View className={styles['tag-list-page']}>
      {/* Hero Section */}
      <View className={styles['hero-section']}>
        <View className={styles['hero-content']}>
          <Text>{t('tags.tags_list')}</Text>
          <Text>{t('tags.explore_all')}</Text>
          <View className={styles['tag-stats']}>
            <View className={styles['stat-item']}>
              <Text>🏷️</Text>
              <Text>{tags.length} {t('tags.total_tags')}</Text>
            </View>
          </View>
        </View>
      </View>

      {/* Tags Grid */}
      <View className={styles['page-container']}>
        {/* Loading Skeleton - TODO: 实现 Taro 骨架屏 */}
        {loading ? (
          <View className={styles['tags-grid']}>
            <Text>加载中...</Text>
          </View>
        ) : (
          <>
            {tags.length > 0 && (
              <View className={styles['tags-grid']}>
                {tags.map((tag) => (
                  <View
                    key={tag.id}
                    className={styles['tag-card']}
                    style={{
                      borderColor: tag.color || 'var(--color-border)',
                      background: `linear-gradient(135deg, ${tag.color}10 0%, var(--color-surface) 100%)`
                    }}
                    onClick={() => handleTagClick(tag.id || 0)}
                  >
                    <View
                      className={styles['tag-icon']}
                      style={{color: tag.color || 'var(--color-brand)'}}
                    >
                      <Text>🏷️</Text>
                    </View>
                    <View className={styles['tag-content']}>
                      <Text>{tagStore.getTranslation(tag)?.name || t('tags.tag_untitled')}</Text>
                      <Text className={styles['tag-description']}>
                        {tagStore.getTranslation(tag)?.description || ''}
                      </Text>
                      <View className={styles['tag-meta']}>
                        <Text>📄</Text>
                        <Text>
                          {tag.postCount || 0} {t('posts.articles')}
                        </Text>
                      </View>
                    </View>
                  </View>
                ))}
                {tags.length === 0 && total > 0 && (
                  <AppEmpty description={t('tags.no_tags_in_page')}/>
                )}
              </View>
            )}

            {!loading && tags.length === 0 && total === 0 && (
              <AppEmpty inContainer description={t('tags.no_tags')}/>
            )}

            {total > pageSize && (
              <View style={{margin: '32px auto 0', display: 'flex', justifyContent: 'center'}}>
                {/* TODO: 实现 Taro 分页组件 */}
                <Text>第 {page} 页，共 {Math.ceil(total / pageSize)} 页</Text>
              </View>
            )}
          </>
        )}
      </View>
    </View>
  );
}
