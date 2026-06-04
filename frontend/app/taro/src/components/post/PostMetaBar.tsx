import {View, Text} from '@tarojs/components';
import React from 'react';
import XIcon from '@/plugins/xicon';
import {formatDate} from '@/utils';

export interface PostMetaBarProps {
    authorName?: string;
    createdAt?: string;
    visits?: number | string;
    likes?: number | string;
}

/**
 * 文章元数据条 — 作者 / 日期 / 浏览量 / 点赞数
 *
 * 用于 post/[id] 详情页 header 区域
 */
const PostMetaBar: React.FC<PostMetaBarProps> = ({authorName, createdAt, visits, likes}) => {
    return (
        <View className='flex flex-wrap items-center gap-x-5 gap-y-2 pb-5 text-sm font-medium text-muted-foreground border-b border-border/40'>
            {authorName && (
                <View className='flex items-center gap-1.5'>
                    <XIcon name='carbon:user-avatar' size={16} />
                    <Text>{authorName}</Text>
                </View>
            )}
            {createdAt && (
                <View className='flex items-center gap-1.5'>
                    <XIcon name='carbon:calendar' size={16} />
                    <Text>{formatDate(createdAt)}</Text>
                </View>
            )}
            <View className='flex items-center gap-1.5'>
                <XIcon name='carbon:view' size={16} />
                <Text>{visits || 0}</Text>
            </View>
            <View className='flex items-center gap-1.5'>
                <XIcon name='carbon:thumbs-up' size={16} />
                <Text>{likes || 0}</Text>
            </View>
        </View>
    );
};

export default PostMetaBar;
