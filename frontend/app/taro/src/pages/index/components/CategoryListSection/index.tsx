import {useState, useEffect, useCallback, useMemo, useRef} from 'react';
import {useTranslation} from 'react-i18next';
import {View, Text} from '@tarojs/components';
import XIcon from '@/plugins/xicon';

import {useCategoryStore} from '@/store/slices/category/hooks';
import type {contentservicev1_Category, contentservicev1_ListCategoryResponse} from '@/api/generated/app/service/v1';

import HomeCategoryCard from '../HomeCategoryCard';

import './index.scss';
import {useI18nRouter} from "@/i18n/helpers";

interface CategoryListSectionProps {
  skeletonCount?: number;
  pageSize?: number;
  page?: number;
  filter?: Record<string, unknown>;
  orderBy?: string[];
  fieldMask?: string;
  showHeader?: boolean;
}

export default function CategoryListSection(props: CategoryListSectionProps) {
  const {
    skeletonCount = 8,
    pageSize = 8,
    page = 1,
    filter: filterProp,
    orderBy: orderByProp,
    fieldMask: fieldMaskProp,
    showHeader = true
  } = props;
  const {t} = useTranslation();
  const categoryStore = useCategoryStore();
  // 用 useRef 缓存 listCategory 方法，确保依赖稳定
  const listCategoryRef = useRef(categoryStore.listCategory);

  const [categories, setCategories] = useState<contentservicev1_Category[]>([]);
  const [loading, setLoading] = useState(false);

  // useMemo 缓存参数，避免每次渲染新对象/数组
  const filter = useMemo(() => filterProp ?? {status: 'CATEGORY_STATUS_ACTIVE'}, [filterProp]);
  const orderBy = useMemo(() => orderByProp ?? ['-sortOrder', '-postCount'], [orderByProp]);
  const fieldMask = useMemo(() => fieldMaskProp ?? 'id,status,sortOrder,icon,code,postCount,directPostCount,parent_id,createdAt,translations.id,translations.categoryId,translations.name,translations.languageCode,translations.description', [fieldMaskProp]);

  const loadCategories = useCallback(async () => {
    setLoading(true);
    try {
      console.log('Load categories:', {page, pageSize, filter, fieldMask, orderBy});
      const res = await listCategoryRef.current({
        paging: {page, pageSize},
        formValues: filter,
        fieldMask,
        orderBy,
      }) as unknown as contentservicev1_ListCategoryResponse;
      setCategories(res.items || []);
    } catch (error) {
      console.error('Failed to load categories:', error);
      setCategories([]);
    } finally {
      setLoading(false);
    }
  }, [filter, fieldMask, orderBy, page, pageSize]);

  useEffect(() => {
    loadCategories();
  }, [loadCategories]);

  const router = useI18nRouter();

  const handleViewCategoryDetail = (id: number) => {
    router.push(`/category/${id}`);
  };

  const handleViewCategory = () => {
    router.push(`/category`);
  };

  return (
    <View className='categories-section scroll-reveal'>
      {/* Section Header */}
      {showHeader && (
        <View className='section-header'>
          <Text className='section-title'>
            <XIcon name='carbon:folder' size={24} /> {t('page.home.categories')}
          </Text>
          <View
            className='view-all-btn'
            onClick={handleViewCategory}
          >
            <Text>{t('page.home.view_all')} →</Text>
          </View>
        </View>
      )}

      {/* Loading Skeleton */}
      {loading ? (
        <View className='categories-grid'>
          {Array.from({length: skeletonCount}).map((_, i) => (
            <View key={i} className='category-card-skeleton'/>
          ))}
        </View>
      ) : (
        <View className='categories-grid'>
          {categories.map((category) => (
            <HomeCategoryCard
              key={category.id}
              category={category}
              onClick={handleViewCategoryDetail}
            />
          ))}
        </View>
      )}
    </View>
  );
}
