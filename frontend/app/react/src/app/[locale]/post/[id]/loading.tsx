import {ArticleDetailSkeleton} from '@/components/ui/skeleton-presets';

/**
 * /post/[id] 文章详情页加载骨架
 */
export default function PostDetailLoading() {
    return (
        <div className="w-full">
            <ArticleDetailSkeleton/>
        </div>
    );
}
