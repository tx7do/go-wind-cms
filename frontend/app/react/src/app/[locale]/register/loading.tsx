import {Skeleton} from '@/components/ui/skeleton';

/**
 * /register 注册页加载骨架
 */
export default function RegisterLoading() {
    return (
        <div className="flex min-h-screen w-full items-center justify-center bg-background p-4">
            <div className="w-full max-w-md space-y-6 rounded-2xl border border-border bg-card p-8 shadow-lg">
                <div className="flex flex-col items-center gap-2">
                    <Skeleton className="h-12 w-12 rounded-full"/>
                    <Skeleton className="h-7 w-32"/>
                    <Skeleton className="h-4 w-48"/>
                </div>
                <div className="space-y-4">
                    {Array.from({length: 4}).map((_, i) => (
                        <div key={i} className="space-y-2">
                            <Skeleton className="h-4 w-24"/>
                            <Skeleton className="h-10 w-full rounded-lg"/>
                        </div>
                    ))}
                    <Skeleton className="h-10 w-full rounded-lg"/>
                </div>
            </div>
        </div>
    );
}
