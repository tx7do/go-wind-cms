import {PageHeroSkeleton, PostListSkeleton} from '@/components/ui/skeleton-presets';

/**
 * /post 列表页加载骨架
 */
export default function PostListLoading() {
    return (
        <div className="w-full">
            <PageHeroSkeleton size="md"/>
            <div className="w-full max-w-[1200px] mx-auto px-8 py-12 max-md:px-4">
                {/* 分类筛选骨架 */}
                <div className="mb-8 flex flex-wrap gap-3">
                    {Array.from({length: 6}).map((_, i) => (
                        <div key={i} className="h-9 w-20 rounded-lg bg-muted animate-pulse"/>
                    ))}
                </div>
                <PostListSkeleton columns={3} count={9}/>
            </div>
        </div>
    );
}
