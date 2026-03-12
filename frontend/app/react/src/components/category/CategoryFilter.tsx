'use client';

import React, {useState, useEffect, useRef} from 'react';
import {Button} from 'antd';
import {useTranslations} from 'next-intl';

import {XIcon} from '@/plugins/xicon';
import {useCategoryStore} from '@/store/slices/category/hooks';
import type {contentservicev1_Category, contentservicev1_ListCategoryResponse} from '@/api/generated/app/service/v1';

import styles from './CategoryFilter.module.css';

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
    const categoryStore = useCategoryStore();
    const categoryT = useTranslations('page.categories');

    // 内部状态管理
    const [internalCategories, setInternalCategories] = useState<contentservicev1_Category[]>([]);
    const [loading, setLoading] = useState(false);
    const [expandedIds, setExpandedIds] = useState<Set<number>>(new Set());

    // 定时器管理
    const hideTimers = useRef<Map<number, ReturnType<typeof setTimeout>>>(new Map());

    // 加载分类数据
    async function loadCategories() {
        setLoading(true);
        try {
            // eslint-disable-next-line @typescript-eslint/no-explicit-any
            const query: Record<string, any> = {status: 'CATEGORY_STATUS_ACTIVE'};

            // 如果指定了 parentId，添加过滤条件
            if (parentId !== undefined && parentId !== null) {
                query.parentId = parentId;
            }

            const res = await categoryStore.listCategory({
                // @ts-expect-error - listCategory 参数类型推断问题
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
    }

    // 监听语言切换，自动重新加载数据
    useEffect(() => {
        if (autoLoad && !externalCategories) {
            loadCategories();
        }
    }, [autoLoad, externalCategories, loadCategories]);

    // 获取分类名称
    function getCategoryName(category: contentservicev1_Category | null): string {
        if (!category?.id) return '';
        return categoryStore.getCategoryName(category, categoryT);
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

    // 显示子菜单
    const showSubmenu = (nodeId: number) => {
        if (hasChildren(nodeId)) {
            if (hideTimers.current.has(nodeId)) {
                clearTimeout(hideTimers.current.get(nodeId));
                hideTimers.current.delete(nodeId);
            }
            setExpandedIds(prev => new Set(prev).add(nodeId));
        }
    };

    // 保持子菜单打开 (鼠标在菜单上时)
    const keepSubmenuOpen = (nodeId: number) => {
        if (hideTimers.current.has(nodeId)) {
            clearTimeout(hideTimers.current.get(nodeId));
            hideTimers.current.delete(nodeId);
        }
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

    function hasChildren(categoryId: number): boolean {
        const category = displayCategories.find(cat => cat.id === categoryId);
        return !!(category && category.children && category.children.length > 0);
    }

    if (loading && autoLoad) {
        return <div className={styles.loading}>加载中...</div>;
    }

    return (
        <div className={styles.categoryFilter}>
            <div className={styles.categoryTabs}>
                {/* 所有分类按钮 */}
                <Button
                    type={selectedCategory === null ? 'primary' : 'default'}
                    ghost={selectedCategory !== null}
                    size="large"
                    onClick={() => handleCategoryChange(null)}
                    icon={<XIcon name="carbon:grid" size={16}/>}
                >
                    {t('all_categories')}
                </Button>

                {/* 树形模式 */}
                {treeMode ? (
                    <>
                        {/* 一级分类 (横向排列 + 悬浮菜单) */}
                        {displayCategories.map((node) => (
                            <div
                                key={node.id}
                                className={styles.categoryItemWrapper}
                                onMouseEnter={() => node.id && showSubmenu(node.id)}
                                onMouseLeave={() => node.id && hideSubmenu(node.id)}
                            >
                                <Button
                                    type={selectedCategory === node.id ? 'primary' : 'default'}
                                    ghost={selectedCategory !== node.id}
                                    size="large"
                                    onClick={() => node.id && handleCategoryClick(node.id)}
                                    icon={<XIcon name={node.icon || 'carbon:folder'} size={16}/>}
                                >
                                    {getCategoryName(node)}
                                </Button>

                                {/* 子分类菜单 */}
                                {hasChildren(node.id || 0) && expandedIds.has(node.id || 0) && (
                                    <div
                                        className={styles.categorySubmenu}
                                        onMouseEnter={() => keepSubmenuOpen(node.id || 0)}
                                        onMouseLeave={() => hideSubmenu(node.id || 0)}
                                    >
                                        {node.children!.map((child) => (
                                            <Button
                                                key={child.id}
                                                type={selectedCategory === child.id ? 'primary' : 'default'}
                                                ghost={selectedCategory !== child.id}
                                                size="medium"
                                                onClick={(e) => {
                                                    e.stopPropagation();
                                                    handleCategoryChange(child.id || 0);
                                                }}
                                                icon={<XIcon name={child.icon || 'carbon:folder'} size={14}/>}
                                            >
                                                {getCategoryName(child)}
                                            </Button>
                                        ))}
                                    </div>
                                )}
                            </div>
                        ))}
                    </>
                ) : (
                    /* 平铺模式 (默认) */
                    <>
                        {rootCategories.map((cat) => (
                            <Button
                                key={cat.id}
                                type={selectedCategory === cat.id ? 'primary' : 'default'}
                                ghost={selectedCategory !== cat.id}
                                size="large"
                                onClick={() => handleCategoryChange(cat.id || 0)}
                                icon={<XIcon name={cat.icon || 'carbon:folder'} size={16}/>}
                            >
                                {getCategoryName(cat)}
                            </Button>
                        ))}
                    </>
                )}
            </div>
        </div>
    );
};

export default CategoryFilter;
