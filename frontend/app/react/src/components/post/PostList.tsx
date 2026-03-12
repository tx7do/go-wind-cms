'use client';

import React, {useState, useEffect, useCallback} from 'react';
import {Pagination, Empty, Skeleton} from 'antd';
import {useTranslations} from 'next-intl';

import {usePostStore} from '@/store/slices/post/hooks';
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
    showSkeleton?: boolean;
    from?: string;
    categoryId?: number;
    columns?: number; // 控制列数
    showPagination?: boolean; // 是否显示分页
    pageSizes?: number[]; // 每页条数选项
}

const PostList: React.FC<PostListProps> = ({
                                               queryParams = {},
                                               fieldMask,
                                               orderBy,
                                               page = 1,
                                               pageSize = 6,
                                               showSkeleton = true,
                                               from = 'post-list',
                                               categoryId,
                                               columns = 3,
                                               showPagination = false,
                                               pageSizes = [6, 12, 24, 36]
                                           }) => {
    const t = useTranslations('page.posts');
    const postStore = usePostStore();
    const [posts, setPosts] = useState<contentservicev1_Post[]>([]);
    const [loading, setLoading] = useState(false);
    const [total, setTotal] = useState(0);
    const [currentPage, setCurrentPage] = useState<number>(page);
    const [currentPageSize, setCurrentPageSize] = useState<number>(pageSize);

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

    const fetchPosts = useCallback(async (page: number, pageSize: number) => {
        setLoading(true);
        try {
            const res = await postStore.listPost({
                // @ts-expect-error - 参数类型推断问题
                paging: {
                    page: page,
                    pageSize: pageSize,
                },
                formValues: queryParams,
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
    }, [postStore, queryParams, fieldMask, orderBy]);

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
                <div className={styles.emptyState}>
                    <Empty description={t('no_results')}/>
                </div>
            )}
        </div>
    );
};

export default PostList;
