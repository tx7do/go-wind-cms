import {Skeleton} from '@/components/ui/skeleton';

/**
 * 通用静态页面 (about, contact, privacy, terms, disclaimer) 加载骨架
 */
export default function StaticPageLoading() {
    return (
        <div className="w-full">
            {/* Hero Skeleton */}
            <div className="relative w-full overflow-hidden border-b border-border min-h-[280px] py-18 flex items-center justify-center">
                <div className="flex flex-col items-center gap-4">
                    <Skeleton className="h-12 w-12 rounded-full"/>
                    <Skeleton className="h-10 w-56 rounded-lg"/>
                    <Skeleton className="h-5 w-80 max-md:w-64"/>
                </div>
            </div>
            {/* Content Skeleton */}
            <div className="w-full max-w-3xl mx-auto px-8 py-12 max-md:px-4">
                <div className="space-y-4">
                    {Array.from({length: 3}).map((_, i) => (
                        <div key={i} className="space-y-2">
                            <Skeleton className="h-6 w-48"/>
                            <Skeleton className="h-4 w-full"/>
                            <Skeleton className="h-4 w-[95%]"/>
                            <Skeleton className="h-4 w-[88%]"/>
                            <Skeleton className="h-4 w-[92%]"/>
                        </div>
                    ))}
                </div>
            </div>
        </div>
    );
}
