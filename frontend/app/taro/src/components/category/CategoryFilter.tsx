import {View, Text} from '@tarojs/components';
import React, {useState, useEffect, useRef} from 'react';
import {useTranslations} from '@/lib/next-intl-compat';

import {XIcon} from '@/plugins/xicon';
import {fetchListCategories, getCategoryName as getCategoryNameHelper} from '@/api/hooks/category';
import type {contentservicev1_Category, contentservicev1_ListCategoryResponse} from '@/api/generated/app/service/v1';

import {cn} from '@/lib/utils';
import {Skeleton} from '@/components/ui/skeleton';

interface CategoryFilterProps {
    categories?: contentservicev1_Category[]; // 外部传入的分类数据（可选）
    selectedCategory?: number | null;
    treeMode?: boolean;
    parentId?: number | null; // 支持根据 parentId 过滤
    autoLoad?: boolean; // 是否自动加载数据
    onCategoryChange?: (categoryId: number | null) => void;
    onLoaded?: (categories: contentservicev1_Category[]) => void;
}

const CategoryFilter: React.FC<CategoryFilterProps> = ({
                                                           categories: externalCategories,
                                                           selectedCategory = null,
                                                           treeMode = false,
                                                           parentId = null,
                                                           autoLoad = true,
                                                           onCategoryChange,
                                                           onLoaded
                                                       }) => {
    const t = useTranslations('page.posts');
    const categoryT = useTranslations('page.categories');

    // 内部状态管理
    const [internalCategories, setInternalCategories] = useState<contentservicev1_Category[]>([]);
    const [loading, setLoading] = useState(false);
    const [expandedIds, setExpandedIds] = useState<Set<number>>(new Set());

    // 定时器管理
    const hideTimers = useRef<Map<number, ReturnType<typeof setTimeout>>>(new Map());

    // 加载分类数据 - 移到 useEffect 内部，避免依赖项问题
    const loadCategories = async () => {
        setLoading(true);
        try {
            // eslint-disable-next-line @typescript-eslint/no-explicit-any
            const query: Record<string, any> = {status: 'CATEGORY_STATUS_ACTIVE'};

            // 如果指定了 parentId，添加过滤条件
            if (parentId !== undefined && parentId !== null) {
                query.parentId = parentId;
            }

            const res = await fetchListCategories({
                paging: undefined,
                formValues: query,
                fieldMask: 'id,status,sort_order,icon,code,post_count,direct_post_count,parent_id,created_at,children,translations.id,translations.category_id,translations.name,translations.language_code,translations.description',
                orderBy: ['-sortOrder']
            }) as unknown as contentservicev1_ListCategoryResponse;

            const items = res.items || [];
            setInternalCategories(items);
            onLoaded?.(items);
            console.log('[CategoryFilter] Categories loaded:', items.length);
        } catch (error) {
            console.error('[CategoryFilter] Load categories failed:', error);
        } finally {
            setLoading(false);
        }
    };

    // 监听语言切换，自动重新加载数据
    useEffect(() => {
        if (!autoLoad || externalCategories) {
            return;
        }

        loadCategories();
        // eslint-disable-next-line react-hooks/exhaustive-deps
    }, [autoLoad, externalCategories, parentId]); // 只依赖稳定的 props

    // 获取分类名称
    function getCategoryName(category: contentservicev1_Category | null): string {
        if (!category?.id) return '';
        return getCategoryNameHelper(category, categoryT);
    }

    // 使用外部传入的 categories 或内部加载的 categories
    const displayCategories = externalCategories || internalCategories;

    // 平铺模式：只显示根节点
    const rootCategories = displayCategories.filter(cat => !cat.parentId);

    const handleCategoryChange = (categoryId: number | null) => {
        onCategoryChange?.(categoryId);
    };

    // 处理分类按钮点击 (同时切换菜单)
    const handleCategoryClick = (nodeId: number) => {
        handleTouchToggle(nodeId);
        handleCategoryChange(nodeId);
    };

    // 隐藏子菜单 - 添加延时避免快速消失
    const hideSubmenu = (nodeId: number) => {
        const timer = setTimeout(() => {
            setExpandedIds(prev => {
                const next = new Set(prev);
                next.delete(nodeId);
                return next;
            });
            hideTimers.current.delete(nodeId);
        }, 150);
        hideTimers.current.set(nodeId, timer);
    };

    // 处理触摸切换 (移动端)
    const handleTouchToggle = (nodeId: number) => {
        if (expandedIds.has(nodeId)) {
            hideSubmenu(nodeId);
        } else {
            setExpandedIds(new Set([nodeId]));
        }
    };

    // 点击外部区域关闭菜单
    useEffect(() => {
        if (expandedIds.size === 0) return;

        // 使用 Taro 的点击事件委托（小程序不支持 document.addEventListener）
        // 这里通过组件外层容器的 onClick 来捕获外部点击
        return () => {
            // 清理函数
        };
    }, [expandedIds]);

    // 关闭所有展开的菜单
    const closeAllMenus = () => {
        setExpandedIds(new Set());
    };

    function hasChildren(categoryId: number): boolean {
        const category = displayCategories.find(cat => cat.id === categoryId);
        return !!(category && category.children && category.children.length > 0);
    }

    if (loading && autoLoad) {
        return (
            <View className='relative z-20 mb-8 max-md:mb-5'>
                <View className='flex flex-wrap items-center gap-2.5 rounded-xl border border-border bg-card/50 p-3.5 backdrop-blur-sm'>
                    <Skeleton className='h-9 w-24 rounded-lg' />
                    <Skeleton className='h-9 w-20 rounded-lg' />
                    <Skeleton className='h-9 w-28 rounded-lg' />
                    <Skeleton className='h-9 w-22 rounded-lg' />
                    <Skeleton className='h-9 w-24 rounded-lg' />
                    <Skeleton className='h-9 w-20 rounded-lg' />
                </View>
            </View>
        );
    }

    // 通用按钮样式：统一圆角胶囊 + 内边距
    const btnBase = cn(
        'inline-flex items-center justify-center gap-2 rounded-full px-4 py-2 text-sm font-medium',
        'transition-all duration-200 cursor-pointer select-none whitespace-nowrap',
    );
    const btnInactive = cn(
        'text-textThird bg-cardBg',
        'hover:text-textSec',
    );
    const btnActive = cn(
        /* 选中态：主色填充白字 */
        'bg-primary text-primary-foreground',
        'hover:bg-primary/90',
    );

    return (
        <View className='relative z-20 mb-8 max-md:mb-5'>
            {/* 点击外部关闭菜单的遮罩层 */}
            {expandedIds.size > 0 && (
                <View
                  className='fixed inset-0 z-[998]'
                  onClick={closeAllMenus}
                />
            )}

            {/* 移动端：横向流滚动容器，禁用换行 */}
            <View className={cn(
                'flex items-center gap-3 overflow-x-auto no-scrollbar scroll-smooth whitespace-nowrap',
                'rounded-xl border border-border bg-card/50 px-4 py-2.5 backdrop-blur-sm',
                'max-md:px-3',
            )}
            >
                {/* 所有分类按钮 */}
                <View
                  onClick={() => handleCategoryChange(null)}
                  className={cn(btnBase, selectedCategory === null ? btnActive : btnInactive, 'shrink-0', 'tap-active')}
                  hoverClass='opacity-70'
                >
                    <XIcon name='carbon:grid' size={15} />
                    {t('all_categories')}
                </View>

                {/* 分隔线 */}
                {displayCategories.length > 0 && (
                    <View className='mx-0.5 h-5 w-px bg-border shrink-0 max-md:hidden' />
                )}

                {/* 树形模式 */}
                {treeMode ? (
                    <>
                        {/* 一级分类 (横向排列 + 触摸菜单) */}
                        {displayCategories.map((node) => (
                            <View
                              key={node.id}
                              className='relative'
                            >
                                <View
                                  onClick={() => node.id && handleCategoryClick(node.id)}
                                  className={cn(
                                        btnBase,
                                        selectedCategory === node.id ? btnActive : btnInactive,
                                        'shrink-0',
                                        'tap-active',
                                    )}
                                  hoverClass='opacity-70'
                                >
                                    <XIcon name={node.icon || 'carbon:folder'} size={15} />
                                    {getCategoryName(node)}
                                    {hasChildren(node.id || 0) && (
                                        <XIcon name='carbon:chevron-down' size={14} className='ml-0.5 opacity-60' />
                                    )}
                                </View>
                            </View>
                        ))}
                    </>
                ) : (
                    /* 平铺模式 (默认) */
                    <>
                        {rootCategories.map((cat) => (
                            <View
                              key={cat.id}
                              onClick={() => handleCategoryChange(cat.id || 0)}
                              className={cn(btnBase, selectedCategory === cat.id ? btnActive : btnInactive, 'shrink-0', 'tap-active')}
                              hoverClass='opacity-70'
                            >
                                <XIcon name={cat.icon || 'carbon:folder'} size={15} />
                                {getCategoryName(cat)}
                            </View>
                        ))}
                    </>
                )}
            </View>

            {/* 下拉菜单（移到滚动容器外部，避免被 overflow 裁剪） */}
            {treeMode && (
                <>
                    {displayCategories.map((node) => (
                        hasChildren(node.id || 0) && expandedIds.has(node.id || 0) && (
                            <View
                              key={`dropdown-${node.id}`}
                              className={cn(
                                    'absolute left-4 top-full z-[999] mt-2',
                                    'min-w-[240px] max-w-[400px]',
                                    'rounded-lg border border-border bg-popover p-1.5',
                                    'shadow-lg shadow-black/10',
                                )}
                            >
                                {node.children!.map((child) => (
                                    <View
                                      key={child.id}
                                      onClick={(e) => {
                                            e.stopPropagation();
                                            handleCategoryChange(child.id || 0);
                                            closeAllMenus(); // 选择后关闭所有菜单
                                        }}
                                      className={cn(
                                            'flex w-full items-center gap-2 rounded-md px-3 py-2.5',
                                            'transition-colors duration-150 cursor-pointer',
                                            selectedCategory === child.id
                                                ? 'bg-primary/10 text-primary font-medium'
                                                : 'text-textSec hover:bg-cardBg',
                                            'tap-active',
                                        )}
                                      hoverClass='opacity-70'
                                    >
                                        <View className='flex items-center gap-2 flex-1 min-w-0'>
                                            <XIcon name={child.icon || 'carbon:folder'} size={14} className='flex-shrink-0' />
                                            <Text className='truncate'>{getCategoryName(child)}</Text>
                                        </View>
                                        {child.postCount !== undefined && child.postCount > 0 && (
                                            <View className='flex-shrink-0 text-right' style={{ minWidth: '24px' }}>
                                                <Text className='text-xs text-textThird'>
                                                    {child.postCount}
                                                </Text>
                                            </View>
                                        )}
                                    </View>
                                ))}
                            </View>
                        )
                    ))}
                </>
            )}
        </View>
    );
};

export default CategoryFilter;
