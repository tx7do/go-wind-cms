import {useState, useEffect} from 'react';
import {useTranslation} from 'react-i18next';
import {View, Text} from '@tarojs/components';
import {AppEmpty} from '@/components/ui';

import {useCategoryStore} from '@/store/slices/category/hooks';
import CategoryTree from '@/components/category/CategoryTree';

import {contentservicev1_Category, contentservicev1_ListCategoryResponse} from "@/api/generated/app/service/v1";
import {useI18nRouter} from "@/i18n/helpers";

import './category.scss';

export default function CategoryListPage() {
  const {t} = useTranslation();
  const router = useI18nRouter();
  const categoryStore = useCategoryStore();

  const [loading, setLoading] = useState(false);
  const [categories, setCategories] = useState<contentservicev1_Category[]>([]);

  async function loadCategories() {
    setLoading(true);
    try {
      const res = (await categoryStore.listCategory({
        paging: undefined,
        formValues: {status: 'CATEGORY_STATUS_ACTIVE'},
        fieldMask: 'id,status,sort_order,icon,code,post_count,direct_post_count,parent_id,created_at,children,translations.id,translations.category_id,translations.name,translations.language_code,translations.description,translations.thumbnail,translations.cover_image',
        orderBy: ['-sortOrder']
      })) as unknown as contentservicev1_ListCategoryResponse;
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
    router.push(`/category/${id}`);
  };

  return (
    <View className="category-page">
      {/* Hero Section */}
      <View className="hero-section">
        <View className="hero-content">
          <Text className="hero-title">内容分类</Text>
          <Text className="hero-subtitle">浏览所有内容分类</Text>
        </View>
      </View>

      {/* Content Section */}
      <View className="page-container">
        {/* Loading Skeleton - TODO: 实现 Taro 的加载骨架屏 */}
        {loading ? (
          <View className="categories-loading">
            <Text>{t('page.common.loading')}</Text>
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
                image={<span className="i-carbon:folder-blank" style={{fontSize: '64px'}}/>}
              />
            )}
          </>
        )}
      </View>
    </View>
  );
}
