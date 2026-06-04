import {useState, useEffect} from 'react';
import {useTranslation} from 'react-i18next';
import {View, Text} from '@tarojs/components';

import {AppEmpty} from '@/components/ui';
import {Skeleton} from '@/components/ui/skeleton';
import XIcon from '@/plugins/xicon';
import {usePageTitle} from '@/hooks/usePageTitle';

import CategoryTree from '@/components/category/CategoryTree';

import {fetchListCategories} from '@/api/hooks/category';
import {useI18nRouter} from '@/i18n/helpers';

import {contentservicev1_Category} from '@/api/generated/app/service/v1';

export default function CategoryListPage() {
  const {t} = useTranslation();
  usePageTitle('page.title.categories');
  const router = useI18nRouter();

  const [loading, setLoading] = useState(false);
  const [categories, setCategories] = useState<contentservicev1_Category[]>([]);

  async function loadCategories() {
    setLoading(true);
    try {
      const res = await fetchListCategories({
        paging: undefined,
        formValues: {status: 'CATEGORY_STATUS_ACTIVE'},
        fieldMask: 'id,status,sort_order,icon,code,post_count,direct_post_count,parent_id,created_at,children,translations.id,translations.category_id,translations.name,translations.language_code,translations.description,translations.thumbnail,translations.cover_image',
        orderBy: ['-sortOrder']
      });
      setCategories(res.items || []);
    } catch (error) {
      console.error('Load categories failed:', error);
    } finally {
      setLoading(false);
    }
  }

  useEffect(() => {
    loadCategories();
  }, []);

  const handleCategoryClick = (id: number) => {
    router.push(`/category/detail?id=${id}`);
  };

  // 统计分类总数和文章总数
  const totalPosts = categories.reduce((sum, cat) => sum + (cat.postCount || 0), 0);

  return (
    <View className='min-h-screen w-full bg-pageBg pb-[240rpx]'>
      {/* 页面标题 - px-[24rpx] 与 Header 对齐 */}
      <View className='bg-cardBg px-[24rpx] pt-[40rpx] pb-[24rpx] border-b-[1rpx] border-splitLine'>
        <Text className='text-title font-bold text-textMain'>{t('page.categories.categories_list')}</Text>
        <Text className='text-desc text-textSec block mt-[8rpx]'>{t('page.categories.explore_all')}</Text>
        <View className='flex items-center gap-[16rpx] mt-[16rpx]'>
          <View className='flex items-center gap-[6rpx]'>
            <XIcon name='carbon:folder' size={14} className='text-textSec' />
            <Text className='text-tips text-textSec'>{categories.length} {t('page.categories.categories_count')}</Text>
          </View>
          <View className='w-[1rpx] h-[20rpx] bg-splitLine' />
          <View className='flex items-center gap-[6rpx]'>
            <XIcon name='carbon:document' size={14} className='text-textSec' />
            <Text className='text-tips text-textSec'>{totalPosts} {t('page.posts.articles')}</Text>
          </View>
        </View>
      </View>

      {/* 分类列表 */}
      <View className='px-[24rpx] pt-[24rpx]'>
        {loading ? (
          <View className='flex flex-col gap-[24rpx]'>
            {Array.from({length: 4}).map((_, i) => (
              <View key={i} className='rounded-[16rpx] bg-cardBg overflow-hidden'>
                <View className='flex items-center gap-[24rpx] p-[24rpx]'>
                  <Skeleton className='w-[180rpx] h-[120rpx] rounded-[12rpx]' />
                  <View className='flex-1 flex flex-col gap-[16rpx]'>
                    <Skeleton className='h-[32rpx] w-[60%] rounded' />
                    <Skeleton className='h-[24rpx] w-full rounded' />
                    <Skeleton className='h-[24rpx] w-[40%] rounded' />
                  </View>
                </View>
              </View>
            ))}
          </View>
        ) : (
          <>
            {categories.length > 0 ? (
              <CategoryTree
                categories={categories}
                onCategoryClick={handleCategoryClick}
              />
            ) : (
              <AppEmpty
                description={t('page.categories.no_categories')}
                inContainer
                image={<View style={{fontSize: '48px'}}><XIcon name='carbon:folder-blank' size={48} /></View>}
              />
            )}
          </>
        )}
      </View>
    </View>
  );
}
