import {useState, useMemo} from 'react';
import {useTranslation} from 'react-i18next';
import {View, Text} from '@tarojs/components';

import CategoryFilter from '@/components/category/CategoryFilter';
import PostList from '@/components/post/PostList';

export default function PostListPage() {
  const {t} = useTranslation();
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
          <Text className='hero-title'>{t('page.posts.posts_list')}</Text>
          <Text className='hero-subtitle'>{t('page.posts.explore_latest')}</Text>
        </View>
      </View>

      <View className='page-container'>
        <View className='category-filter-container'>
          <CategoryFilter
            selectedCategory={selectedCategoryId}
            treeMode
            autoLoad
            onCategoryChange={handleCategoryChange}
          />
        </View>

        <PostList
          key={selectedCategoryId || 'all'}  // 使用 key 强制刷新组件
          queryParams={queryParams}
          initialPageSize={12}
          pageSizes={[12, 24, 36, 48]}
        />
      </View>
    </View>
  );
}
