'use client';

import {useState, useEffect, useMemo} from 'react';
import {useParams} from 'next/navigation';
import {useTranslations} from 'next-intl';
import {Button} from '@/components/ui/button';
import {ArrowLeft} from 'lucide-react';
import {AppEmpty} from '@/components/ui';

import {XIcon} from '@/plugins/xicon';
import {useI18nRouter} from "@/i18n/helpers";

import {fetchTag, getTranslation as getTagTranslation} from '@/api/hooks/tag';
import PostListWithPagination from '@/components/post/PostList';
import type {contentservicev1_Tag} from "@/api/generated/app/service/v1";

import '../../../globals.css'; // 导入全局 CSS，确保 CSS 变量可用
import styles from './tag-detail.module.css';

export default function TagDetailPage() {
    const t = useTranslations('page');
    const params = useParams();
    const router = useI18nRouter();

    const [loading, setLoading] = useState(false);
    const [tag, setTag] = useState<contentservicev1_Tag | null>(null);

    const tagId = useMemo(() => {
        const id = params?.id;
        return id ? parseInt(id as string) : null;
    }, [params?.id]);

    async function loadTag() {
        if (!tagId) return;

        setLoading(true);
        try {
            const tagData = await fetchTag(tagId!) as unknown as contentservicev1_Tag;
            setTag(tagData);
        } catch (error) {
            console.error('Load tag failed:', error);
        } finally {
            setLoading(false);
        }
    }

    function handleBack() {
        router.push('/tag');
    }

    useEffect(() => {
        loadTag();
    }, [tagId]);

    if (!tagId) {
        return <AppEmpty variant="error" description="Invalid tag ID"/>;
    }

    return (
        <div className={styles['tag-detail-page']}>
            {/* Hero Section */}
            <div
                className={styles['hero-section']}
                style={{
                    background: tag?.color
                        ? `linear-gradient(135deg, ${tag.color} 0%, ${tag.color}dd 50%, ${tag.color}aa 100%)`
                        : undefined
                }}
            >
                <div className={styles['hero-content']}>
                    <div
                        className={styles['tag-icon']}
                        style={{color: tag?.color || '#6366f1'}}
                    >
                        <XIcon name={`carbon:${tag?.icon || 'tag'}`} size={64}/>
                    </div>
                    <h1>{getTagTranslation(tag)?.name || t('tags.tag_untitled')}</h1>
                    {getTagTranslation(tag)?.description && (
                        <p className={styles['tag-description']}>
                            {getTagTranslation(tag)?.description}
                        </p>
                    )}
                    <div className={styles['tag-stats']}>
                        <div className={styles['stat-item']}>
                            <XIcon name="carbon:document" size={20}/>
                            <span>{tag?.postCount || 0} {t('posts.articles')}</span>
                        </div>
                    </div>
                </div>
            </div>

            {/* Posts Section */}
            <div className={styles['page-container']}>
                {/* Back Button */}
                <div className={styles['back-button-container']}>
                    <Button
                        variant="outline"
                        size="sm"
                        onClick={handleBack}
                        className="gap-2"
                    >
                        <ArrowLeft className="h-4 w-4"/>
                        {t('categories.back_to_list')}
                    </Button>
                </div>

                {/* Posts List with Pagination */}
                {tagId && (
                    <PostListWithPagination
                        key={tagId}
                        initialPageSize={10}
                        pageSizes={[10, 20, 30, 40]}
                        tagId={tagId}
                        from="tag"
                    />
                )}
            </div>
        </div>
    );
}
