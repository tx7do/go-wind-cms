import {useState, useMemo} from 'react';
import {useTranslation} from 'react-i18next';
import {View, Text} from '@tarojs/components';

import CategoryFilter from '@/components/category/CategoryFilter';
import PostList from '@/components/post/PostList';

import './index.scss';

export default function PostListPage() {
  const {t} = useTranslation('page');
  const [selectedCategoryId, setSelectedCategoryId] = useState<number | null>(null);

  const handleCategoryChange = (categoryId: number | null) => {
    console.log('[PostListPage] Category changed:', categoryId);
    setSelectedCategoryId(categoryId);
  };

  // 构建查询参数 - 使用 useMemo 确保稳定性
  const queryParams = useMemo(() => {
    if (selectedCategoryId) {
      return {category_ids__in: [selectedCategoryId]};
    }
    return {};
  }, [selectedCategoryId]);

  return (
    <View className='post-list-page'>
      {/* Hero Section */}
      <View className='hero-section'>
        <View className='hero-content'>
          <Text>{t('posts.posts_list')}</Text>
          <Text>{t('posts.explore_latest')}</Text>
        </View>
      </View>

      <View className='page-container'>
        <CategoryFilter
          selectedCategory={selectedCategoryId}
          treeMode={true}
          autoLoad={true}
          onCategoryChange={handleCategoryChange}
        />

        <PostList
          key={selectedCategoryId || 'all'}  // 使用 key 强制重新渲染
          queryParams={queryParams}
          initialPageSize={12}
          pageSizes={[12, 24, 36, 48]}
        />
      </View>
    </View>
  );
}
