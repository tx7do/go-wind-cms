import {useEffect, useMemo, useState} from 'react';
import {useTranslation} from 'react-i18next';
import {View, Text} from '@tarojs/components';
import Taro from '@tarojs/taro';

import {AppEmpty} from '@/components/ui';
import XIcon from '@/plugins/xicon';
import {useI18nRouter} from '@/i18n/helpers';

import CategoryList from '@/components/category/CategoryList';
import PostListWithPagination from '@/components/post/PostList';
import {
  fetchCategory,
  getCategoryName,
  getCategoryDescription,
} from '@/api/hooks/category';
import {contentservicev1_Category} from '@/api/generated/app/service/v1';

export default function CategoryDetailPage() {
  const {t} = useTranslation();
  const router = useI18nRouter();

  const [detail, setDetail] = useState<contentservicev1_Category | null>(null);
  const [childCategories, setChildCategories] = useState<contentservicev1_Category[]>([]);

  const categoryId = useMemo(() => {
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

  const parentCategoryId = useMemo(() => {
    if (!detail?.parentId) return null;
    return detail.parentId;
  }, [detail?.parentId]);

  async function loadCategory() {
    if (!categoryId) return;

    try {
      const loadedCategory = await fetchCategory({
        id: categoryId,
        fieldMask: 'id,status,sort_order,icon,code,post_count,direct_post_count,parent_id,created_at,children,translations.id,translations.category_id,translations.name,translations.language_code,translations.description,translations.thumbnail,translations.cover_image'
      });
      setDetail(loadedCategory);
      if (loadedCategory?.children && loadedCategory.children.length > 0) {
        setChildCategories([...loadedCategory.children]);
      } else {
        setChildCategories([]);
      }
    } catch (error) {
      console.error('Load category failed:', error);
    }
  }

  function handleViewChildCategory(id: number) {
    router.push(`/category/detail?id=${id}`);
  }

  function handleBack() {
    if (parentCategoryId) {
      router.push(`/category/detail?id=${parentCategoryId}`);
    } else {
      router.push('/category');
    }
  }

  useEffect(() => {
    loadCategory();
  }, [categoryId]);

  if (!categoryId) {
    return <AppEmpty variant='error' description='Invalid category ID' />;
  }

  const categoryName = getCategoryName(detail);
  const categoryDesc = getCategoryDescription(detail);
  const hasChildren = childCategories.length > 0;

  return (
    <View className='min-h-screen w-full bg-pageBg pb-[160rpx]'>
      {/* 页面标题 */}
      <View className='bg-cardBg px-[24rpx] pt-[40rpx] pb-[24rpx] border-b-[1rpx] border-splitLine'>
        <View className='flex items-center gap-[16rpx]'>
          {/* 分类图标 */}
          <View className='w-[48rpx] h-[48rpx] rounded-full flex items-center justify-center bg-primary/10 flex-shrink-0'>
            <XIcon name={detail?.icon || 'carbon:folder'} size={20} className='text-primary' />
          </View>
          <View className='flex-1 min-w-0'>
            <Text className='text-title font-bold text-textMain'>
              {categoryName}
            </Text>
            {categoryDesc && (
              <Text className='text-desc text-textSec block mt-[8rpx]'>
                {categoryDesc}
              </Text>
            )}
          </View>
        </View>
        {/* 统计信息 */}
        <View className='flex items-center gap-[16rpx] mt-[16rpx] text-tips text-textThird'>
          <View className='flex items-center gap-[6rpx]'>
            <XIcon name='carbon:document' size={14} className='text-textThird' />
            <Text className='text-tips text-textThird'>
              {detail?.postCount || 0} {t('page.posts.articles')}
            </Text>
          </View>
          {hasChildren && (
            <>
              <View className='w-[1rpx] h-[20rpx] bg-splitLine' />
              <View className='flex items-center gap-[6rpx]'>
                <XIcon name='carbon:folder' size={14} className='text-textThird' />
                <Text className='text-tips text-textThird'>
                  {childCategories.length} {t('page.categories.sub_categories')}
                </Text>
              </View>
            </>
          )}
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
          <Text className='text-desc text-textSec'>
            {parentCategoryId ? t('page.categories.back_to_parent') : t('page.categories.back_to_list')}
          </Text>
        </View>
      </View>

      {/* 子分类列表 */}
      {hasChildren && (
        <View className='px-[24rpx] pt-[24rpx]'>
          <View className='flex items-center gap-[12rpx] mb-[16rpx]'>
            <XIcon name='carbon:folder' size={16} className='text-textSec' />
            <Text className='text-body font-bold text-textMain'>{t('page.categories.sub_categories')}</Text>
          </View>
          <CategoryList
            categories={childCategories}
            loading={false}
            showSkeleton={false}
            columns={2}
            onCategoryClick={handleViewChildCategory}
          />
        </View>
      )}

      {/* 文章列表 */}
      <View className='px-[24rpx] pt-[24rpx]'>
        {categoryId && (
          <PostListWithPagination
            key={categoryId}
            initialPageSize={12}
            categoryId={categoryId}
            from='category'
            showPagination
          />
        )}
      </View>
    </View>
  );
}
