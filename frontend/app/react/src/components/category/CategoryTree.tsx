'use client';

import React, {useState} from 'react';
import {useTranslations} from 'next-intl';

import {XIcon} from '@/plugins/xicon';
import {useCategoryStore} from '@/store/slices/category/hooks';
import type {contentservicev1_Category} from '@/api/generated/app/service/v1';

import styles from './CategoryTree.module.css';

interface CategoryTreeProps {
    categories: contentservicev1_Category[];
    level?: number;
    onCategoryClick?: (id: number) => void;
}

const CategoryTree: React.FC<CategoryTreeProps> = ({
                                                       categories = [],
                                                       level = 0,
                                                       onCategoryClick
                                                   }) => {
    const t = useTranslations('page.categories');
    const categoryStore = useCategoryStore();
    const [expandedCategories, setExpandedCategories] = useState<Set<number>>(new Set());

    const handleViewCategory = (id: number) => {
        onCategoryClick?.(id);
    };

    const toggleExpand = (e: React.MouseEvent, category: contentservicev1_Category) => {
        e.stopPropagation();
        if (category.children && category.children.length > 0) {
            const newExpanded = new Set(expandedCategories);
            if (newExpanded.has(category.id || 0)) {
                newExpanded.delete(category.id || 0);
            } else {
                newExpanded.add(category.id || 0);
            }
            setExpandedCategories(newExpanded);
        }
    };

    const isExpanded = (category: contentservicev1_Category) => {
        return expandedCategories.has(category.id || 0);
    };

    if (!categories || categories.length === 0) return null;

    return (
        <div className={styles.categoryTree}>
            {categories.map((category) => (
                <div
                    key={category.id}
                    className={`${styles.categoryNode} ${styles[`level${level}`]}`}
                >
                    <div
                        className={styles.categoryItem}
                        onClick={() => handleViewCategory(category.id || 0)}
                    >
                        <div className={styles.categoryInfo}>
                            <div className={styles.categoryImage}>
                                <img
                                    src={categoryStore.getCategoryThumbnail(category)}
                                    alt={categoryStore.getCategoryName(category, t)}
                                />
                                <div className={styles.imageOverlay}/>
                            </div>
                            <div className={styles.categoryContent}>
                                <h3 className={category.children && category.children.length > 0 ? styles.hasChildren : ''}>
                                    {categoryStore.getCategoryName(category, t)}
                                </h3>
                                <p className={styles.description}>
                                    {categoryStore.getCategoryDescription(category)}
                                </p>
                                <div className={styles.categoryMeta}>
                                    <span className={styles.metaIcon}>
                                        <XIcon name="carbon:document" size={16}/>
                                    </span>
                                    <span className={styles.metaText}>
                                        {category.postCount || 0} {t('articles_count')}
                                    </span>
                                </div>
                            </div>
                        </div>

                        {/* 展开/收起按钮 (如果有子分类) */}
                        {category.children && category.children.length > 0 && (
                            <button
                                className={`${styles.expandBtn} ${isExpanded(category) ? styles.expanded : ''}`}
                                onClick={(e) => toggleExpand(e, category)}
                            >
                                <XIcon
                                    name={isExpanded(category) ? 'carbon:chevron-down' : 'carbon:chevron-right'}
                                    size={20}
                                />
                            </button>
                        )}
                    </div>

                    {/* 递归渲染子分类 */}
                    {category.children && category.children.length > 0 && isExpanded(category) && (
                        <div className={styles.slideContainer}>
                            <CategoryTree
                                categories={category.children}
                                level={level + 1}
                                onCategoryClick={onCategoryClick}
                            />
                        </div>
                    )}
                </div>
            ))}
        </div>
    );
};

export default CategoryTree;
