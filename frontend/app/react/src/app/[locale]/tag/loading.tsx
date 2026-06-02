import {PageHeroSkeleton, TagListSkeleton} from '@/components/ui/skeleton-presets';

/**
 * /tag 标签列表页加载骨架
 */
export default function TagListLoading() {
    return (
        <div className="w-full">
            <PageHeroSkeleton size="md"/>
            <div className="w-full max-w-[1200px] mx-auto px-8 py-12 max-md:px-4">
                <TagListSkeleton count={12}/>
            </div>
        </div>
    );
}
