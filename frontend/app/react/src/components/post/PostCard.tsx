'use client';

import React from 'react';

import {XIcon} from '@/plugins/xicon';
import {
    getPostTitle,
    getPostSummary,
    getPostThumbnail,
} from '@/api/hooks/post';
import {useI18nRouter} from '@/i18n/helpers/useI18nRouter';

import type {contentservicev1_Post} from '@/api/generated/app/service/v1';

import {formatDate} from "@/utils";
import {cn} from '@/lib/utils';

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

    const handleViewPost = () => {
        const query: string[] = [`from=${from}`];
        if (categoryId) {
            query.push(`categoryId=${categoryId}`);
        }

        router.push(`/post/${post.id}?${query.join('&')}`);

        // 滚动到顶部
        window.scrollTo({top: 0, behavior: 'smooth'});
    };

    return (
        <article
            className={cn(
                'group flex h-full cursor-pointer flex-col overflow-hidden',
                'rounded-2xl border border-border bg-card shadow-sm',
                /* 空气动力学悬浮：纯 primary 辉光（无暗影底），更空灵 */
                'transition-all duration-500 ease-out',
                'hover:-translate-y-1.5 hover:border-primary/40',
                'hover:shadow-[0_20px_40px_-8px_hsl(var(--primary)/0.15)]',
            )}
            onClick={handleViewPost}
            style={{ willChange: 'transform, box-shadow' }}
        >
            <div className="relative h-[240px] w-full flex-shrink-0 overflow-hidden bg-background max-md:h-[200px]">
                <img
                    src={getPostThumbnail(post)}
                    alt={getPostTitle(post)}
                    className="h-full w-full object-cover transition-transform duration-700 ease-out group-hover:scale-[1.12]"
                />
                {/* hover 风迹渐变蒙层 */}
                <div className="absolute inset-0 bg-gradient-to-tr from-primary/20 via-transparent to-sky-400/20 opacity-0 transition-opacity duration-500 group-hover:opacity-100"/>
            </div>
            <div className="flex flex-1 flex-col gap-3 p-6 max-md:p-4 max-md:gap-2.5">
                {/* 标题：固定 min-h 确保两行高度 */}
                <h3 className={cn(
                    'line-clamp-2 min-h-[3.4em] text-lg font-bold leading-[1.7] text-foreground transition-colors duration-300',
                    'group-hover:text-primary',
                    'max-md:min-h-[3em] max-md:text-[17px]',
                )}>
                    {getPostTitle(post)}
                </h3>
                {/* 摘要：固定 min-h 确保三行高度，底部的数据列永退贴底 */}
                <p className={cn(
                    'line-clamp-3 min-h-[4.4em] flex-1 text-sm leading-[1.55] text-muted-foreground',
                    'max-md:min-h-[4em] max-md:text-[13px]',
                )}>
                    {getPostSummary(post)}
                </p>
                <div className={cn(
                    'flex flex-wrap gap-4 border-t border-border pt-3 text-[13px] font-medium text-muted-foreground',
                    'max-md:text-xs max-md:gap-3',
                )}>
                    <div className="flex items-center gap-1.5 whitespace-nowrap">
                        <XIcon name="carbon:user" size={16}/>
                        <span>{post.authorName}</span>
                    </div>
                    <div className="flex items-center gap-1.5 whitespace-nowrap">
                        <XIcon name="carbon:calendar" size={16}/>
                        <span>{formatDate(post.createdAt)}</span>
                    </div>
                    <div className="flex items-center gap-1.5 whitespace-nowrap">
                        <XIcon name="carbon:view" size={16}/>
                        <span>{post.visits || 0}</span>
                    </div>
                    <div className="flex items-center gap-1.5 whitespace-nowrap">
                        <XIcon name="carbon:thumbs-up" size={16}/>
                        <span>{post.likes || 0}</span>
                    </div>
                </div>
            </div>
        </article>
    );
};

export default PostCard;
