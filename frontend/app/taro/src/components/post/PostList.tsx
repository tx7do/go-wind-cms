import {View, Text} from '@tarojs/components';
import React, {useState, useEffect, useCallback} from 'react';
import {Skeleton} from '@/components/ui/skeleton';
import {Button} from '@/components/ui/button';
import {useTranslations} from '@/lib/next-intl-compat';
import {AppEmpty} from '@/components/ui';

import {fetchListPosts} from '@/api/hooks/post';
import type {
    contentservicev1_ListPostResponse,
    contentservicev1_Post
} from '@/api/generated/app/service/v1';

import PostCard from './PostCard';

interface PostListProps {
    queryParams?: object;
    fieldMask?: string;
    orderBy?: string[];
    page?: number;
    pageSize?: number;
    initialPageSize?: number;
    showSkeleton?: boolean;
    from?: string;
    categoryId?: number;
    tagId?: number;
    columns?: number; // 保留接口兼容，实际移动端固定单列
    showPagination?: boolean;
    pageSizes?: number[];
    /** 紧凑横排模式 */
    compact?: boolean;
}

const PostList: React.FC<PostListProps> = ({
    queryParams = {},
    fieldMask,
    orderBy,
    page = 1,
    pageSize,
    initialPageSize = 10,
    showSkeleton = true,
    from = 'post-list',
    categoryId,
    tagId,
    showPagination = false,
    pageSizes = [10, 20, 30, 40],
    compact = false,
}) => {
    const t = useTranslations('page.posts');
    const [posts, setPosts] = useState<contentservicev1_Post[]>([]);
    const [loading, setLoading] = useState(false);
    const [total, setTotal] = useState(0);
    const [currentPage, setCurrentPage] = useState<number>(page);
    const [currentPageSize, setCurrentPageSize] = useState<number>(pageSize ?? initialPageSize);

    useEffect(() => {
        if (page !== undefined && page !== currentPage) {
            setCurrentPage(page);
        }
    }, [currentPage, page]);

    useEffect(() => {
        if (pageSize !== undefined && pageSize !== currentPageSize) {
            setCurrentPageSize(pageSize);
        }
    }, [currentPageSize, pageSize]);

    useEffect(() => {
        setCurrentPage(1);
        fetchPosts(1, currentPageSize);
    }, [categoryId, tagId]);

    const fetchPosts = useCallback(async (page: number, pageSize: number) => {
        setLoading(true);
        try {
            const res = await fetchListPosts({
                paging: {page, pageSize},
                formValues: {
                    ...queryParams,
                    ...(categoryId && {category_ids__in: [categoryId]}),
                    ...(tagId && {tag_ids__in: [tagId]}),
                },
                fieldMask,
                orderBy,
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

    useEffect(() => {
        fetchPosts(currentPage, currentPageSize);
    }, [currentPage, currentPageSize]);

    const handlePageChange = (newPage: number) => {
        setCurrentPage(newPage);
    };

    const handlePageSizeChange = (size: number, newPage: number) => {
        setCurrentPageSize(size);
        setCurrentPage(newPage);
    };

    const totalPages = Math.ceil(total / currentPageSize);

    return (
        <View className='w-full'>
            {/* Loading Skeleton */}
            {loading && showSkeleton && (
                <View className='flex flex-col gap-[24rpx]'>
                    {Array.from({length: Math.min(currentPageSize, 3)}).map((_, index) => (
                        <View key={index} className='rounded bg-cardBg overflow-hidden'>
                            {!compact ? (
                                <>
                                    <Skeleton className='w-full h-[280rpx]' />
                                    <View className='p-[24rpx] flex flex-col gap-[16rpx]'>
                                        <Skeleton className='h-[32rpx] w-3/4 rounded' />
                                        <Skeleton className='h-[28rpx] w-full rounded' />
                                        <Skeleton className='h-[28rpx] w-1/2 rounded' />
                                    </View>
                                </>
                            ) : (
                                <View className='flex flex-row p-[20rpx] gap-[20rpx]'>
                                    <Skeleton className='w-[200rpx] h-[200rpx] rounded flex-shrink-0' />
                                    <View className='flex-1 flex flex-col gap-[12rpx] justify-center'>
                                        <Skeleton className='h-[28rpx] w-full rounded' />
                                        <Skeleton className='h-[28rpx] w-2/3 rounded' />
                                        <Skeleton className='h-[24rpx] w-1/2 rounded' />
                                    </View>
                                </View>
                            )}
                        </View>
                    ))}
                </View>
            )}

            {/* Loaded Content */}
            {!loading && posts.length > 0 && (
                <>
                    {showPagination && total > currentPageSize && (
                        <View className='mb-[24rpx] rounded bg-cardBg px-[24rpx] py-[20rpx]'>
                            <Text className='text-desc text-textSec'>{t('total_articles', {total})}</Text>
                        </View>
                    )}

                    <View className='flex flex-col gap-[24rpx]'>
                        {posts.map((post, index) => (
                            <PostCard
                              key={`${post.id}-${index}`}
                              post={post}
                              from={from}
                              categoryId={categoryId}
                              compact={compact}
                            />
                        ))}
                    </View>

                    {showPagination && totalPages > 1 && (
                        <View className='flex items-center justify-center gap-[24rpx] py-[48rpx]'>
                            <Button
                              variant='outline'
                              size='sm'
                              disabled={currentPage <= 1}
                              onClick={() => handlePageChange(currentPage - 1)}
                            >
                                {t('previous') || '上一页'}
                            </Button>
                            <Text className='text-desc text-textSec'>
                                {currentPage} / {totalPages}
                            </Text>
                            <Button
                              variant='outline'
                              size='sm'
                              disabled={currentPage >= totalPages}
                              onClick={() => handlePageChange(currentPage + 1)}
                            >
                                {t('next') || '下一页'}
                            </Button>
                        </View>
                    )}
                </>
            )}

            {/* Empty State */}
            {!loading && posts.length === 0 && (
                <AppEmpty description={t('no_results')} />
            )}
        </View>
    );
};

export default PostList;
