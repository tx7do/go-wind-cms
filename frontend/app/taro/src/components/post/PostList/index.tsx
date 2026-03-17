import React, {useState, useEffect, useCallback} from 'react';
import {useTranslation} from 'react-i18next';
import {View, Text} from '@tarojs/components';

import {AppEmpty} from '@/components/ui';

import {usePostStore} from '@/store/slices/post/hooks';
import type {
  contentservicev1_ListPostResponse,
  contentservicev1_Post
} from '@/api/generated/app/service/v1';

import PostCard from '../PostCard';

import './index.scss';

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
  columns?: number;
  showPagination?: boolean;
  pageSizes?: number[];
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
                                             columns = 3,
                                             showPagination = false,
                                           }) => {
  const {t} = useTranslation('page.posts');
  const postStore = usePostStore();
  const [posts, setPosts] = useState<contentservicev1_Post[]>([]);
  const [loading, setLoading] = useState(false);
  const [total, setTotal] = useState(0);
  const [currentPage, setCurrentPage] = useState<number>(page);
  const [currentPageSize, setCurrentPageSize] = useState<number>(pageSize ?? initialPageSize);

  // 同步外部 page 变化
  useEffect(() => {
    if (page !== undefined && page !== currentPage) {
      setCurrentPage(page);
    }
  }, [currentPage, page]);

  // 同步外部 pageSize 变化
  useEffect(() => {
    if (pageSize !== undefined && pageSize !== currentPageSize) {
      setCurrentPageSize(pageSize);
    }
  }, [currentPageSize, pageSize]);

  // 监听 categoryId 和 tagId 变化，重置页码并重新查询
  useEffect(() => {
    setCurrentPage(1);
    fetchPosts(1, currentPageSize);
  }, [categoryId, tagId]);

  const fetchPosts = useCallback(async (page: number, pageSize: number) => {
    setLoading(true);
    try {
      const res = await postStore.listPost({
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
  }, [postStore, queryParams, categoryId, tagId, fieldMask, orderBy]);

  // 监听页面变化
  useEffect(() => {
    fetchPosts(currentPage, currentPageSize);
  }, [currentPage, currentPageSize]);

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
    <View className='post-list-container'>
      {/* Loading Skeleton */}
      {loading && showSkeleton && (
        <View
          className='posts-grid'
          style={{gridTemplateColumns: `repeat(${columns}, 1fr)`}}
        >
          {Array.from({length: currentPageSize}).map((_, index) => (
            <View key={index} className='post-card-skeleton'>
              <View className='skeleton-image-placeholder'/>
              <View className='skeleton-content'>
                <View className='skeleton-title-placeholder'/>
                <View className='skeleton-paragraph-placeholder'/>
                <View className='skeleton-meta'>
                  <View className='skeleton-btn-small'/>
                  <View className='skeleton-btn-small'/>
                  <View className='skeleton-btn-small'/>
                </View>
              </View>
            </View>
          ))}
        </View>
      )}

      {/* Loaded Content */}
      {!loading && posts.length > 0 && (
        <>
          {showPagination && total > currentPageSize && (
            <View className='results-info'>
              <Text>{t('total_articles', {total})}</Text>
            </View>
          )}

          <View
            className='posts-grid'
            style={{gridTemplateColumns: `repeat(${columns}, 1fr)`}}
          >
            {posts.map((post) => (
              <PostCard
                key={post.id}
                post={post}
                from={from}
                categoryId={categoryId}
              />
            ))}
          </View>

          {showPagination && total > currentPageSize && (
            <View className='pagination-wrapper'>
              {/* TODO: 实现 Taro 兼容的分页组件 */}
              <View className='pagination-simple'>
                <Text
                  className={`page-btn ${currentPage === 1 ? 'disabled' : ''}`}
                  onClick={() => currentPage > 1 && handlePageChange(currentPage - 1)}
                >
                  上一页
                </Text>
                <Text className='page-info'>{currentPage}</Text>
                <Text
                  className={`page-btn ${currentPage * currentPageSize >= total ? 'disabled' : ''}`}
                  onClick={() => currentPage * currentPageSize < total && handlePageChange(currentPage + 1)}
                >
                  下一页
                </Text>
              </View>
            </View>
          )}
        </>
      )}

      {/* Empty State */}
      {!loading && posts.length === 0 && (
        <AppEmpty description={t('no_results')}/>
      )}
    </View>
  );
};

export default PostList;
