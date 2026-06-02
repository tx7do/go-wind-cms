'use client';

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
        <div className="flex flex-wrap items-center gap-x-5 gap-y-2 pb-5 text-sm font-medium text-muted-foreground border-b border-border/40">
            {authorName && (
                <div className="flex items-center gap-1.5">
                    <XIcon name="carbon:user-avatar" size={16}/>
                    <span>{authorName}</span>
                </div>
            )}
            {createdAt && (
                <div className="flex items-center gap-1.5">
                    <XIcon name="carbon:calendar" size={16}/>
                    <span>{formatDate(createdAt)}</span>
                </div>
            )}
            <div className="flex items-center gap-1.5">
                <XIcon name="carbon:view" size={16}/>
                <span>{visits || 0}</span>
            </div>
            <div className="flex items-center gap-1.5">
                <XIcon name="carbon:thumbs-up" size={16}/>
                <span>{likes || 0}</span>
            </div>
        </div>
    );
};

export default PostMetaBar;
