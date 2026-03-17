import React from 'react';
import {View, Text, Image} from '@tarojs/components';

import {usePostStore} from '@/store/slices/post/hooks';
import {useI18nRouter} from '@/i18n/helpers/useI18nRouter';

import type {contentservicev1_Post} from '@/api/generated/app/service/v1';
import {formatDate} from "@/utils";

import './index.scss';

interface PostCardProps {
  post: contentservicev1_Post;
  from?: string;
  categoryId?: number;
}

const PostCard: React.FC<PostCardProps> = ({
                                             post,
                                             from = 'post-list',
                                             categoryId
                                           }) => {
  const router = useI18nRouter();
  const postStore = usePostStore();

  const handleViewPost = () => {
    const query: string[] = [`from=${from}`];
    if (categoryId) {
      query.push(`categoryId=${categoryId}`);
    }

    router.push(`/post/${post.id}?${query.join('&')}`);

    // Taro 中滚动到顶部
    // TODO: 使用 Taro.pageScrollTo API
  };

  return (
    <View className='post-card' onClick={handleViewPost}>
      <View className='post-image'>
        <Image
          src={postStore.getPostThumbnail(post)}
          mode='aspectFill'
          className='post-image-img'
        />
        <View className='image-overlay'/>
      </View>
      <View className='post-content'>
        <Text className='post-title'>{postStore.getPostTitle(post)}</Text>
        <Text className='post-summary'>{postStore.getPostSummary(post)}</Text>
        <View className='post-meta'>
          <View className='meta-item'>
            <Text className='meta-icon'>👤</Text>
            <Text>{post.authorName}</Text>
          </View>
          <View className='meta-item'>
            <Text className='meta-icon'>📅</Text>
            <Text>{formatDate(post.createdAt)}</Text>
          </View>
          <View className='meta-item'>
            <Text className='meta-icon'>👁️</Text>
            <Text>{post.visits || 0}</Text>
          </View>
          <View className='meta-item'>
            <Text className='meta-icon'>👍</Text>
            <Text>{post.likes || 0}</Text>
          </View>
        </View>
      </View>
    </View>
  );
};

export default PostCard;
