'use client';

import React from 'react';
import XIcon from '@/plugins/xicon';

export interface PostFloatingActionsProps {
    isLiked: boolean;
    isBookmarked: boolean;
    onLike: () => void;
    onBookmark: () => void;
    onShare: () => void;
    labels: {
        likes: string;
        bookmark: string;
        share: string;
    };
}

/**
 * 文章浮动操作按钮 — 点赞 / 收藏 / 分享
 *
 * 固定定位在页面右侧垂直居中，仅 lg 以上显示
 */
const PostFloatingActions: React.FC<PostFloatingActionsProps> = ({
    isLiked,
    isBookmarked,
    onLike,
    onBookmark,
    onShare,
    labels,
}) => {
    const btnBase = 'flex h-12 w-12 cursor-pointer items-center justify-center rounded-full border border-border bg-card shadow-md transition-all hover:border-primary hover:bg-primary/5 hover:shadow-lg';

    return (
        <div className="fixed top-1/2 right-8 z-10 hidden -translate-y-1/2 flex-col gap-3 lg:flex">
            <button
                onClick={onLike}
                className={`${btnBase} ${isLiked ? 'border-primary bg-primary/10 text-primary' : 'text-muted-foreground'}`}
                aria-label={labels.likes}
            >
                <XIcon name={isLiked ? 'carbon:thumbs-up-filled' : 'carbon:thumbs-up'} size={20}/>
            </button>
            <button
                onClick={onBookmark}
                className={`${btnBase} ${isBookmarked ? 'border-primary bg-primary/10 text-primary' : 'text-muted-foreground'}`}
                aria-label={labels.bookmark}
            >
                <XIcon name={isBookmarked ? 'carbon:bookmark-filled' : 'carbon:bookmark'} size={20}/>
            </button>
            <button
                onClick={onShare}
                className={`${btnBase} text-muted-foreground hover:text-primary`}
                aria-label={labels.share}
            >
                <XIcon name="carbon:share" size={20}/>
            </button>
        </div>
    );
};

export default PostFloatingActions;
