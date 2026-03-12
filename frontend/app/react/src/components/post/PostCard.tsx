'use client';

import React from 'react';

import {XIcon} from '@/plugins/xicon';
import {usePostStore} from '@/store/slices/post/hooks';
import {useI18nRouter} from '@/i18n/helpers/useI18nRouter';

import type {contentservicev1_Post} from '@/api/generated/app/service/v1';

import styles from './PostCard.module.css';
import {formatDate} from "@/utils";

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

        // 滚动到顶部
        window.scrollTo({top: 0, behavior: 'smooth'});
    };

    return (
        <article className={styles.postCard} onClick={handleViewPost}>
            <div className={styles.postImage}>
                <img
                    src={postStore.getPostThumbnail(post)}
                    alt={postStore.getPostTitle(post)}
                />
                <div className={styles.imageOverlay}/>
            </div>
            <div className={styles.postContent}>
                <h3 className={styles.postTitle}>{postStore.getPostTitle(post)}</h3>
                <p className={styles.postSummary}>{postStore.getPostSummary(post)}</p>
                <div className={styles.postMeta}>
                    <div className={styles.metaItem}>
                        <XIcon name="carbon:user" size={16}/>
                        <span>{post.authorName}</span>
                    </div>
                    <div className={styles.metaItem}>
                        <XIcon name="carbon:calendar" size={16}/>
                        <span>{formatDate(post.createdAt)}</span>
                    </div>
                    <div className={styles.metaItem}>
                        <XIcon name="carbon:view" size={16}/>
                        <span>{post.visits || 0}</span>
                    </div>
                    <div className={styles.metaItem}>
                        <XIcon name="carbon:thumbs-up" size={16}/>
                        <span>{post.likes || 0}</span>
                    </div>
                </div>
            </div>
        </article>
    );
};

export default PostCard;
