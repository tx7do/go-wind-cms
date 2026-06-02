import {Skeleton} from '@/components/ui/skeleton';

/**
 * /settings 设置页加载骨架
 */
export default function SettingsLoading() {
    return (
        <div className="w-full max-w-[1200px] mx-auto px-8 py-12 max-md:px-4">
            <div className="mb-8">
                <Skeleton className="h-9 w-40"/>
                <Skeleton className="mt-2 h-5 w-72"/>
            </div>
            <div className="grid grid-cols-[280px_1fr] gap-6 max-md:grid-cols-1">
                {/* 侧边栏 */}
                <aside>
                    <div className="space-y-2 rounded-xl border border-border bg-card p-4">
                        {Array.from({length: 5}).map((_, i) => (
                            <div key={i} className="flex items-center gap-3 rounded-lg p-3">
                                <Skeleton className="h-5 w-5 rounded"/>
                                <Skeleton className="h-5 flex-1"/>
                            </div>
                        ))}
                    </div>
                </aside>
                {/* 内容区 */}
                <main className="space-y-6 rounded-xl border border-border bg-card p-6">
                    {Array.from({length: 4}).map((_, i) => (
                        <div key={i} className="space-y-3 border-b border-border pb-6 last:border-b-0 last:pb-0">
                            <Skeleton className="h-6 w-32"/>
                            <Skeleton className="h-4 w-full"/>
                            <Skeleton className="h-4 w-3/4"/>
                            <div className="flex gap-3 pt-2">
                                <Skeleton className="h-10 w-28 rounded-lg"/>
                                <Skeleton className="h-10 w-28 rounded-lg"/>
                            </div>
                        </div>
                    ))}
                </main>
            </div>
        </div>
    );
}
