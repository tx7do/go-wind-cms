import {PageHeroSkeleton, CategoryListSkeleton} from '@/components/ui/skeleton-presets';

/**
 * /category 分类列表页加载骨架
 */
export default function CategoryListLoading() {
    return (
        <div className="w-full">
            <PageHeroSkeleton size="md"/>
            <div className="w-full max-w-[1200px] mx-auto px-8 py-12 max-md:px-4">
                <CategoryListSkeleton count={6}/>
            </div>
        </div>
    );
}
