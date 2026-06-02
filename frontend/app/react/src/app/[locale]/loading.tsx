import {Skeleton} from '@/components/ui/skeleton';

/**
 * [locale] 段的根加载页面
 * 在 layout/页面切换时显示
 */
export default function LocaleLoading() {
    return (
        <div className="w-full">
            {/* Hero Skeleton */}
            <div className="relative w-full overflow-hidden border-b border-border min-h-[420px] py-24 flex items-center justify-center">
                <div className="flex flex-col items-center gap-5">
                    <Skeleton className="h-14 w-72 rounded-xl"/>
                    <Skeleton className="h-6 w-96 max-md:w-72"/>
                    <Skeleton className="h-6 w-72 max-md:w-56"/>
                    <div className="mt-3 flex gap-4">
                        <Skeleton className="h-11 w-32 rounded-lg"/>
                        <Skeleton className="h-11 w-32 rounded-lg"/>
                    </div>
                </div>
            </div>

            {/* Featured Posts Section Skeleton */}
            <div className="w-full max-w-300 mx-auto px-8 py-12 max-md:px-4">
                <div className="mb-8 flex items-center justify-between">
                    <Skeleton className="h-8 w-40"/>
                    <Skeleton className="h-9 w-24 rounded-lg"/>
                </div>
                <div className="grid grid-cols-[repeat(auto-fill,minmax(300px,1fr))] gap-6 max-md:grid-cols-1">
                    {Array.from({length: 3}).map((_, i) => (
                        <div key={i} className="overflow-hidden rounded-2xl border border-border bg-card">
                            <Skeleton className="h-52 w-full"/>
                            <div className="space-y-3 p-5">
                                <Skeleton className="h-5 w-4/5"/>
                                <Skeleton className="h-4 w-full"/>
                                <Skeleton className="h-4 w-3/4"/>
                            </div>
                        </div>
                    ))}
                </div>
            </div>

            {/* Categories Section Skeleton */}
            <div className="w-full max-w-300 mx-auto px-8 py-12 max-md:px-4">
                <div className="mb-8 flex items-center justify-between">
                    <Skeleton className="h-8 w-32"/>
                    <Skeleton className="h-9 w-24 rounded-lg"/>
                </div>
                <div className="grid grid-cols-[repeat(auto-fill,minmax(260px,1fr))] gap-6 max-md:grid-cols-1">
                    {Array.from({length: 4}).map((_, i) => (
                        <div key={i} className="rounded-2xl border border-border bg-card p-6">
                            <Skeleton className="h-35 w-full rounded-xl"/>
                            <div className="mt-4 space-y-2">
                                <Skeleton className="h-5 w-3/5"/>
                                <Skeleton className="h-4 w-full"/>
                                <Skeleton className="h-4 w-4/5"/>
                            </div>
                        </div>
                    ))}
                </div>
            </div>
        </div>
    );
}
