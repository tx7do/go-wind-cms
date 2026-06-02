'use client';

import React, {useState, useEffect, useCallback} from 'react';
import {Pagination, Skeleton} from 'antd';
import {useTranslations} from 'next-intl';
import {AppEmpty} from '@/components/ui';

import {fetchListPosts} from '@/api/hooks/post';
import type {
    contentservicev1_ListPostResponse,
    contentservicev1_Post
} from '@/api/generated/app/service/v1';

import PostCard from './PostCard';
import styles from './PostList.module.css';

interface PostListProps {
    queryParams?: object;
    fieldMask?: string;
    orderBy?: string[];
    page?: number;
    pageSize?: number;
    initialPageSize?: number; // 初始每页条数
    showSkeleton?: boolean;
    from?: string;
    categoryId?: number;
    tagId?: number; // 新增 tagId 支持
    columns?: number; // 控制列数
    showPagination?: boolean; // 是否显示分页
    pageSizes?: number[]; // 每页条数选项
}

const PostList: React.FC<PostListProps> = ({
                                               queryParams = {},
                                               fieldMask,
                                               orderBy,
                                               page = 1,
                                               pageSize,
                                               initialPageSize = 10, // 默认 10 条
                                               showSkeleton = true,
                                               from = 'post-list',
                                               categoryId,
                                               tagId,
                                               columns = 3,
                                               showPagination = false,
                                               pageSizes = [10, 20, 30, 40]
                                           }) => {
    const t = useTranslations('page.posts');
    const [posts, setPosts] = useState<contentservicev1_Post[]>([]);
    const [loading, setLoading] = useState(false);
    const [total, setTotal] = useState(0);
    const [currentPage, setCurrentPage] = useState<number>(page);
    const [currentPageSize, setCurrentPageSize] = useState<number>(pageSize ?? initialPageSize);

    // 同步外部 page 变化（只在首次渲染或外部 page 真正变化时）
    useEffect(() => {
        if (page !== undefined && page !== currentPage) {
            setCurrentPage(page);
        }
    }, [currentPage, page]);

    // 同步外部 pageSize 变化（只在首次渲染或外部 pageSize 真正变化时）
    useEffect(() => {
        if (pageSize !== undefined && pageSize !== currentPageSize) {
            setCurrentPageSize(pageSize);
        }
    }, [currentPageSize, pageSize]);

    // 监听 categoryId 和 tagId 变化，重置页码并重新查询
    useEffect(() => {
        setCurrentPage(1);
        // 直接调用 fetchPosts，使用新的 categoryId/tagId
        fetchPosts(1, currentPageSize);
    }, [categoryId, tagId]);

    const fetchPosts = useCallback(async (page: number, pageSize: number) => {
        setLoading(true);
        try {
            const res = await fetchListPosts({
                paging: {
                    page: page,
                    pageSize: pageSize,
                },
                formValues: {
                    ...queryParams,
                    ...(categoryId && {category_ids__in: [categoryId]}),
                    ...(tagId && {tag_ids__in: [tagId]})
                },
                fieldMask: fieldMask,
                orderBy: orderBy
            }) as unknown as contentservicev1_ListPostResponse;
            setPosts(res.items || []);
            setTotal(res.total || 0);
        } catch (error) {
            console.error('PostList fetch failed:', error);
            setPosts([]);
        } finally {
            setLoading(false);
        }
    }, [queryParams, categoryId, tagId, fieldMask, orderBy]);

    // 监听页面变化
    useEffect(() => {
        fetchPosts(currentPage, currentPageSize);
    }, [currentPage, currentPageSize]); // 移除 fetchPosts 依赖，避免循环

    // 监听外部 page 变化（已在上面处理）

    // 处理页面变化
    const handlePageChange = (newPage: number) => {
        setCurrentPage(newPage);
    };

    // 处理每页条数变化
    const handlePageSizeChange = (size: number, newPage: number) => {
        setCurrentPageSize(size);
        setCurrentPage(newPage);
    };

    return (
        <div className={styles.postListContainer}>
            {/* Loading Skeleton */}
            {loading && showSkeleton && (
                <div className={styles.postsGrid} style={{gridTemplateColumns: `repeat(${columns}, 1fr)`}}>
                    {Array.from({length: currentPageSize}).map((_, index) => (
                        <div key={index} className={styles.postCardSkeleton}>
                            <Skeleton.Image style={{height: 240}}/>
                            <div className={styles.skeletonContent}>
                                <Skeleton paragraph={{rows: 2}} title={false}/>
                                <div className={styles.skeletonMeta}>
                                    <Skeleton.Button size="small" style={{width: 60}}/>
                                    <Skeleton.Button size="small" style={{width: 60}}/>
                                    <Skeleton.Button size="small" style={{width: 60}}/>
                                </div>
                            </div>
                        </div>
                    ))}
                </div>
            )}

            {/* Loaded Content */}
            {!loading && posts.length > 0 && (
                <>
                    {showPagination && total > currentPageSize && (
                        <div className={styles.resultsInfo}>
                            <span>{t('total_articles', {total})}</span>
                        </div>
                    )}

                    <div className={styles.postsGrid} style={{gridTemplateColumns: `repeat(${columns}, 1fr)`}}>
                        {posts.map((post) => (
                            <PostCard
                                key={post.id}
                                post={post}
                                from={from}
                                categoryId={categoryId}
                            />
                        ))}
                    </div>

                    {showPagination && total > currentPageSize && (
                        <div className={styles.paginationWrapper}>
                            <Pagination
                                current={currentPage}
                                pageSize={currentPageSize}
                                total={total}
                                showSizeChanger
                                showQuickJumper
                                pageSizeOptions={pageSizes.map(String)}
                                onChange={handlePageChange}
                                onShowSizeChange={handlePageSizeChange}
                                showTotal={(total) => t('total_articles', {total})}
                            />
                        </div>
                    )}
                </>
            )}

            {/* Empty State */}
            {!loading && posts.length === 0 && (
                <AppEmpty description={t('no_results')}/>
            )}
        </div>
    );
};

export default PostList;
