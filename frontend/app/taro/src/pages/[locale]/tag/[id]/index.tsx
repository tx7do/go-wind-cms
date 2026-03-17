'use client';

import {useState, useEffect, useMemo} from 'react';
import {useTranslation} from 'react-i18next';
import {View, Text} from '@tarojs/components';
import Taro from '@tarojs/taro';
import {AppEmpty} from '@/components/ui';

import {XIcon} from '@/plugins/xicon';
import {useI18nRouter} from "@/i18n/helpers";

import {useTagStore} from '@/store/slices/tag/hooks';
import PostListWithPagination from '@/components/post/PostList';
import type {contentservicev1_Tag} from "@/api/generated/app/service/v1";

import '../../../globals.css'; // 导入全局 CSS，确保 CSS 变量可用
import styles from './tag-detail.scss';

export default function TagDetailPage() {
    const {t} = useTranslation('page');
    const tagId = useMemo(() => {
        const pages = Taro.getCurrentPages();
        const currentPage = pages[pages.length - 1];
        const id = (currentPage as any).options?.id;
        return id ? parseInt(id) : null;
    }, []);

    async function loadTag() {
        if (!tagId) return;

        setLoading(true);
        try {
            const tagData = await tagStore.getTag({
                // @ts-expect-error - 参数类型推断问题
                id: tagId
            }) as unknown as contentservicev1_Tag;
            setTag(tagData);
        } catch (error) {
            console.error('Load tag failed:', error);
        } finally {
            setLoading(false);
        }
    }

    function handleBack() {
        Taro.navigateBack();
    }

    useEffect(() => {
        loadTag();
    }, [tagId]);

    if (!tagId) {
        return <AppEmpty variant="error" description="Invalid tag ID"/>;
    }

    return (
        <View className={styles['tag-detail-page']}>
            {/* Hero Section */}
            <View
                className={styles['hero-section']}
                style={{
                    background: tag?.color
                        ? `linear-gradient(135deg, ${tag.color} 0%, ${tag.color}dd 50%, ${tag.color}aa 100%)`
                        : undefined
                }}
            >
                <View className={styles['hero-content']}>
                    <View
                        className={styles['tag-icon']}
                        style={{color: tag?.color || '#6366f1'}}
                    >
                        <Text>🏷️</Text>
                    </View>
                    <Text>{tagStore.getTranslation(tag)?.name || t('tags.tag_untitled')}</Text>
                    {tagStore.getTranslation(tag)?.description && (
                        <Text className={styles['tag-description']}>
                            {tagStore.getTranslation(tag)?.description}
                        </Text>
                    )}
                    <View className={styles['tag-stats']}>
                        <View className={styles['stat-item']}>
                            <Text>📄</Text>
                            <Text>{tag?.postCount || 0} {t('posts.articles')}</Text>
                        </View>
                    </View>
                </View>
            </View>

            {/* Posts Section */}
            <View className={styles['page-container']}>
                {/* Back Button */}
                <View className={styles['back-button-container']}>
                    <View className={styles['back-btn']} onClick={handleBack}>
                        <Text>← {t('categories.back_to_list')}</Text>
                    </View>
                </View>

                {/* Posts List with Pagination */}
                {tagId && (
                    <PostListWithPagination
                        key={tagId}
                        initialPageSize={10}
                        pageSizes={[10, 20, 30, 40]}
                        tagId={tagId}
                        from="tag"
                    />
                )}
            </View>
        </View>
    );
}
