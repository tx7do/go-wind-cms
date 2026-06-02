import {PageHeroSkeleton, PostListSkeleton} from '@/components/ui/skeleton-presets';

/**
 * /tag/[id] 标签详情页加载骨架
 */
export default function TagDetailLoading() {
    return (
        <div className="w-full">
            <PageHeroSkeleton size="lg"/>
            <div className="w-full max-w-[1200px] mx-auto px-8 py-12 max-md:px-4">
                <PostListSkeleton columns={3} count={6}/>
            </div>
        </div>
    );
}
