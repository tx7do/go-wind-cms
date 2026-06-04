import {View, Text} from '@tarojs/components';
import React from 'react';

import {XIcon} from '@/plugins/xicon';
import {
    getPostTitle,
    getPostSummary,
    getPostThumbnail,
} from '@/api/hooks/post';
import {useI18nRouter} from '@/i18n/helpers/useI18nRouter';

import type {contentservicev1_Post} from '@/api/generated/app/service/v1';

import Image from '@/components/ui/image';
import {formatDate} from '@/utils';
import Taro from '@tarojs/taro';

interface PostCardProps {
    post: contentservicev1_Post;
    from?: string;
    categoryId?: number;
    /** 紧凑模式：横向卡片（缩略图左+文字右） */
    compact?: boolean;
}

const PostCard: React.FC<PostCardProps> = ({
    post,
    from = 'post-list',
    categoryId,
    compact = false,
}) => {
    const router = useI18nRouter();

    const handleViewPost = () => {
        const query: string[] = [`from=${from}`];
        if (categoryId) query.push(`categoryId=${categoryId}`);
        router.push(`/post/detail?id=${post.id}&${query.join('&')}`);
        Taro.pageScrollTo({scrollTop: 0, duration: 300});
    };

    const thumbnail = getPostThumbnail(post);
    const title = getPostTitle(post);
    const summary = getPostSummary(post);

    // ---- 紧凑横排卡片 ----
    if (compact) {
        return (
            <View
              className='flex flex-row rounded bg-cardBg overflow-hidden min-h-touch'
              onClick={handleViewPost}
              hoverClass='tap-active'
            >
                {/* 缩略图 / 图标占位 - 固定 16:9 比例 */}
                <View className='w-[240rpx] h-[136rpx] flex-shrink-0 overflow-hidden bg-pageBg'>
                    {thumbnail ? (
                        <Image src={thumbnail} mode='aspectFill' className='w-full h-full' />
                    ) : (
                        <View className='w-full h-full flex items-center justify-center'>
                            <XIcon name='carbon:document' size={40} className='text-textWeak' />
                        </View>
                    )}
                </View>
                <View className='flex-1 flex flex-col justify-between p-[20rpx] min-w-0'>
                    <Text
                      className='text-desc font-bold text-textMain leading-[1.4]'
                      style={{
                            overflow: 'hidden',
                            textOverflow: 'ellipsis',
                            display: '-webkit-box',
                            WebkitLineClamp: 2,
                            WebkitBoxOrient: 'vertical',
                        }}
                    >
                        {title}
                    </Text>
                    <View className='flex items-center gap-[20rpx] text-tips text-textSec'>
                        <View className='flex items-center gap-[6rpx]'>
                            <XIcon name='carbon:calendar' size={14} className='text-textSec' />
                            <Text className='text-tips text-textSec'>{formatDate(post.createdAt)}</Text>
                        </View>
                        <View className='flex items-center gap-[6rpx]'>
                            <XIcon name='carbon:view' size={14} className='text-textSec' />
                            <Text className='text-tips text-textSec'>{post.visits || 0}</Text>
                        </View>
                    </View>
                </View>
            </View>
        );
    }

    // ---- 默认竖排卡片（封面+标题+摘要） ----
    return (
        <View
          className='flex flex-col rounded-[16rpx] bg-cardBg overflow-hidden min-h-touch'
          onClick={handleViewPost}
          hoverClass='tap-active'
        >
            {/* 封面图 / 渐变背景 + 图标占位 */}
            <View className='w-full h-[280rpx] overflow-hidden flex items-center justify-center' style={{
                background: thumbnail ? 'none' : 'linear-gradient(135deg, rgba(22,119,255,0.06) 0%, rgba(114,46,209,0.06) 100%)',
            }}>
                {thumbnail ? (
                    <Image src={thumbnail} mode='aspectFill' className='w-full h-full' />
                ) : (
                    <XIcon name='carbon:document' size={64} className='text-primary' style={{opacity: 0.6}} />
                )}
            </View>

            {/* 内容区 */}
            <View className='flex flex-col gap-[12rpx] p-[24rpx]'>
                {/* 标题 */}
                {title && title !== 'title' && (
                    <Text
                      className='text-body font-bold text-textMain leading-[1.5]'
                      style={{
                            overflow: 'hidden',
                            textOverflow: 'ellipsis',
                            display: '-webkit-box',
                            WebkitLineClamp: 2,
                            WebkitBoxOrient: 'vertical',
                        }}
                    >
                        {title}
                    </Text>
                )}

                {/* 摘要 */}
                {summary && summary !== 'content' && (
                    <Text
                      className='text-desc text-textSec leading-[1.6]'
                      style={{
                            overflow: 'hidden',
                            textOverflow: 'ellipsis',
                            display: '-webkit-box',
                            WebkitLineClamp: 2,
                            WebkitBoxOrient: 'vertical',
                        }}
                    >
                        {summary}
                    </Text>
                )}

                {/* 元数据 */}
                <View className='flex items-center gap-[16rpx] pt-[12rpx] border-t-[1rpx] border-splitLine text-tips text-textThird'>
                    <View className='flex items-center gap-[4rpx]'>
                        <XIcon name='carbon:user' size={12} className='text-textThird' />
                        <Text className='text-tips text-textThird'>{post.authorName || '—'}</Text>
                    </View>
                    <View className='flex items-center gap-[4rpx]'>
                        <XIcon name='carbon:calendar' size={12} className='text-textThird' />
                        <Text className='text-tips text-textThird'>{formatDate(post.createdAt)}</Text>
                    </View>
                    <View className='flex items-center gap-[4rpx]'>
                        <XIcon name='carbon:view' size={12} className='text-textThird' />
                        <Text className='text-tips text-textThird'>{post.visits || 0}</Text>
                    </View>
                    <View className='flex items-center gap-[4rpx]'>
                        <XIcon name='carbon:thumbs-up' size={12} className='text-textThird' />
                        <Text className='text-tips text-textThird'>{post.likes || 0}</Text>
                    </View>
                </View>
            </View>
        </View>
    );
};

export default PostCard;
