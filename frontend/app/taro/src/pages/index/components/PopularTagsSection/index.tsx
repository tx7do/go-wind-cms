import {useState, useEffect, useCallback} from 'react';
import {useTranslation} from 'react-i18next';
import {View, Text} from '@tarojs/components';
import XIcon from '@/plugins/xicon';

import {fetchListTags} from '@/api/hooks/tag';
import type {contentservicev1_ListTagResponse, contentservicev1_Tag} from '@/api/generated/app/service/v1';
import {useI18nRouter} from "@/i18n/helpers";interface TagItem {
  id: number;
  name: string;
  color: string;
}

// 将十六进制颜色转换为 RGB
function hexToRgb(hex: string): string {
  // 处理十六进制颜色
  if (hex.startsWith('#')) {
    const result = /^#?([a-f\d]{2})([a-f\d]{2})([a-f\d]{2})$/i.exec(hex);
    if (result) {
      return `${parseInt(result[1], 16)}, ${parseInt(result[2], 16)}, ${parseInt(result[3], 16)}`;
    }
  }
  // 处理 HSL 或其他格式，返回默认值
  return '139, 92, 246'; // 默认紫色
}

export default function PopularTagsSection() {
  const {t} = useTranslation();
  const tagStore = {listTag: fetchListTags}; // 使用 API hooks 替代 store

  const [_tags, setTags] = useState<contentservicev1_Tag[]>([]);
  const [loading, setLoading] = useState(false);
  const [displayTags, setDisplayTags] = useState<TagItem[]>([]);

  const loadPopularTags = useCallback(async () => {
    setLoading(true);
    try {
      const res = await tagStore.listTag({
        paging: {page: 1, pageSize: 6},
        formValues: {status: 'TAG_STATUS_ACTIVE', isFeatured: true},
        fieldMask: undefined,
        orderBy: undefined,
      }) as unknown as contentservicev1_ListTagResponse;

      const tagItems = res.items || [];
      setTags(tagItems);

      // 生成带颜色的标签
      const taggedItems: TagItem[] = tagItems
        .filter(tag => tag.id !== undefined)
        .map((tag, index) => ({
          id: tag.id!,
          name: tag.translations?.[0]?.name || t('page.tags.tag_untitled'),
          color: tag.color || `hsl(${index * 60}, 100%, 50%)`,
        }));

      setDisplayTags(taggedItems);
    } catch (error) {
      console.error('Failed to load tags:', error);
      setTags([]);
      setDisplayTags([]);
    } finally {
      setLoading(false);
    }
  }, []); // 移除依赖，避免无限循环

  useEffect(() => {
    loadPopularTags();
  }, []); // 空依赖数组，只在首次渲染时执行

  const router = useI18nRouter();

  const handleViewTag = (tag: TagItem) => {
    router.push(`/tag/${tag.id}`);
  };

  return (
    <View className='popular-section'>
      <View className='section-header'>
        <Text className='section-title'>
          <XIcon name='carbon:fire' size={24} /> {t('page.tags.popular_tags')}
        </Text>
        <View
          className='view-all-btn'
          onClick={() => router.push('/tag')}
        >
          <Text>{t('page.tags.view_all')} →</Text>
        </View>
      </View>
      <View className='tags-content'>
        {loading ? (
          <View className='tags-grid'>
            {Array.from({length: 6}).map((_, i) => (
              <View key={i} className='tag-item skeleton'>
                <View className='skeleton-placeholder' />
              </View>
            ))}
          </View>
        ) : (
          <View className='tags-grid'>
            {displayTags.map((tag) => (
              <View
                key={tag.id}
                className='tag-item'
                onClick={() => handleViewTag(tag)}
                style={{
                  background: `rgba(${hexToRgb(tag.color), 0.08})`,
                  borderColor: `rgba(${hexToRgb(tag.color), 0.2})`,
                  color: tag.color,
                }}
              >
                <Text className='tag-label'>{tag.name}</Text>
              </View>
            ))}
          </View>
        )}
      </View>
    </View>
  );
}
