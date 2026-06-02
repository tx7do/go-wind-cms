import React from 'react';
import {Skeleton} from './skeleton';
import {cn} from '@/lib/utils';

/* ========== Page Hero Skeleton ========== */

interface PageHeroSkeletonProps {
    size?: 'sm' | 'md' | 'lg';
    className?: string;
}

export function PageHeroSkeleton({size = 'md', className}: PageHeroSkeletonProps) {
    const minHeight = size === 'sm' ? 'min-h-[220px]' : size === 'lg' ? 'min-h-[340px]' : 'min-h-[280px]';
    const padding = size === 'sm' ? 'py-14' : size === 'lg' ? 'py-24' : 'py-18';

    return (
        <section className={cn(
            'relative w-full overflow-hidden border-b border-border',
            minHeight,
            padding,
            'flex items-center justify-center',
            className,
        )}>
            <div className="flex flex-col items-center gap-4">
                {size !== 'sm' && <Skeleton className="h-12 w-12 rounded-full"/>}
                <Skeleton className="h-10 w-64 rounded-lg"/>
                <Skeleton className="h-5 w-80 rounded-lg"/>
                <div className="mt-2 flex gap-3">
                    <Skeleton className="h-5 w-20 rounded-full"/>
                    <Skeleton className="h-5 w-20 rounded-full"/>
                </div>
            </div>
        </section>
    );
}

/* ========== Post Card Skeleton ========== */

interface PostCardSkeletonProps {
    className?: string;
}

export function PostCardSkeleton({className}: PostCardSkeletonProps) {
    return (
        <div className={cn('overflow-hidden rounded-2xl border border-border bg-card', className)}>
            <Skeleton className="h-52 w-full"/>
            <div className="space-y-3 p-5">
                <Skeleton className="h-5 w-4/5"/>
                <Skeleton className="h-4 w-full"/>
                <Skeleton className="h-4 w-3/4"/>
                <div className="flex items-center gap-3 pt-2">
                    <Skeleton className="h-5 w-16 rounded-full"/>
                    <Skeleton className="h-5 w-16 rounded-full"/>
                    <Skeleton className="h-4 w-20"/>
                </div>
            </div>
        </div>
    );
}

/* ========== Post List Skeleton ========== */

interface PostListSkeletonProps {
    columns?: number;
    count?: number;
    className?: string;
}

export function PostListSkeleton({columns = 3, count = 6, className}: PostListSkeletonProps) {
    return (
        <div
            className={cn('grid gap-6 max-md:!grid-cols-1 max-md:gap-4', className)}
            style={{gridTemplateColumns: `repeat(${columns}, 1fr)`}}
        >
            {Array.from({length: count}).map((_, i) => (
                <PostCardSkeleton key={i}/>
            ))}
        </div>
    );
}

/* ========== Category Card Skeleton ========== */

interface CategoryCardSkeletonProps {
    className?: string;
}

export function CategoryCardSkeleton({className}: CategoryCardSkeletonProps) {
    return (
        <div className={cn('flex min-h-50 flex-col rounded-2xl border border-border bg-card p-6', className)}>
            <Skeleton className="h-35 w-full rounded-xl"/>
            <div className="mt-4 space-y-2">
                <Skeleton className="h-5 w-3/5"/>
                <Skeleton className="h-4 w-full"/>
                <Skeleton className="h-4 w-4/5"/>
            </div>
        </div>
    );
}

/* ========== Category List Skeleton ========== */

interface CategoryListSkeletonProps {
    count?: number;
    className?: string;
}

export function CategoryListSkeleton({count = 6, className}: CategoryListSkeletonProps) {
    return (
        <div
            className={cn('grid grid-cols-[repeat(auto-fill,minmax(260px,1fr))] gap-6', className)}
        >
            {Array.from({length: count}).map((_, i) => (
                <CategoryCardSkeleton key={i}/>
            ))}
        </div>
    );
}

/* ========== Tag Card Skeleton ========== */

interface TagCardSkeletonProps {
    className?: string;
}

export function TagCardSkeleton({className}: TagCardSkeletonProps) {
    return (
        <div className={cn(
            'flex flex-col items-center gap-3 rounded-2xl border border-border bg-card p-6',
            className,
        )}>
            <Skeleton className="h-10 w-10 rounded-full"/>
            <Skeleton className="h-5 w-20"/>
            <Skeleton className="h-4 w-24"/>
        </div>
    );
}

/* ========== Tag List Skeleton ========== */

interface TagListSkeletonProps {
    count?: number;
    className?: string;
}

export function TagListSkeleton({count = 12, className}: TagListSkeletonProps) {
    return (
        <div
            className={cn('grid grid-cols-[repeat(auto-fill,minmax(260px,1fr))] gap-6', className)}
        >
            {Array.from({length: count}).map((_, i) => (
                <TagCardSkeleton key={i}/>
            ))}
        </div>
    );
}

/* ========== Article Detail Skeleton ========== */

interface ArticleDetailSkeletonProps {
    className?: string;
}

export function ArticleDetailSkeleton({className}: ArticleDetailSkeletonProps) {
    return (
        <div className={cn('w-full max-w-[1200px] mx-auto px-8 py-8 max-md:px-4', className)}>
            {/* Back button */}
            <Skeleton className="mb-6 h-9 w-32 rounded-lg"/>
            <article className="overflow-hidden rounded-xl border border-border bg-card shadow-sm">
                {/* Cover image */}
                <Skeleton className="h-[300px] w-full max-md:h-[200px]"/>
                <div className="flex gap-6 p-8 max-md:flex-col max-md:p-4">
                    {/* Sidebar (TOC) */}
                    <aside className="w-[240px] shrink-0 max-md:hidden">
                        <div className="sticky top-24 space-y-3 rounded-lg border border-border bg-background p-4">
                            <Skeleton className="h-6 w-[200px]"/>
                            <Skeleton className="mt-3 h-5 w-[180px]"/>
                            <Skeleton className="mt-2 h-5 w-[160px]"/>
                            <Skeleton className="mt-2 h-5 w-[140px]"/>
                            <Skeleton className="mt-2 h-5 w-[170px]"/>
                        </div>
                    </aside>
                    {/* Content */}
                    <div className="flex-1">
                        <Skeleton className="mb-4 h-12 w-[80%]"/>
                        <Skeleton className="mb-4 h-6 w-[60%]"/>
                        <div className="mb-6 flex gap-3">
                            <Skeleton className="h-6 w-[80px] rounded-full"/>
                            <Skeleton className="h-6 w-[80px] rounded-full"/>
                            <Skeleton className="h-6 w-[80px] rounded-full"/>
                            <Skeleton className="h-6 w-[80px] rounded-full"/>
                        </div>
                        <div className="space-y-3">
                            <Skeleton className="h-4 w-full"/>
                            <Skeleton className="h-4 w-[95%]"/>
                            <Skeleton className="h-4 w-[90%]"/>
                            <Skeleton className="h-4 w-full"/>
                            <Skeleton className="h-4 w-[85%]"/>
                            <Skeleton className="h-4 w-[92%]"/>
                            <Skeleton className="h-4 w-full"/>
                            <Skeleton className="h-4 w-[88%]"/>
                        </div>
                    </div>
                </div>
            </article>
        </div>
    );
}

/* ========== User Page Skeleton ========== */

interface UserPageSkeletonProps {
    className?: string;
}

export function UserPageSkeleton({className}: UserPageSkeletonProps) {
    return (
        <div className={cn('w-full', className)}>
            {/* Top gradient */}
            <div className="h-[200px] animate-pulse bg-gradient-to-r from-primary/10 to-primary/5"/>
            <div className="w-full max-w-[1200px] mx-auto grid grid-cols-[300px_1fr] gap-6 px-8 py-8 max-md:grid-cols-1 max-md:px-4">
                {/* Sidebar */}
                <aside>
                    <div className="space-y-4 rounded-xl border border-border bg-card p-6">
                        <Skeleton className="mx-auto h-20 w-20 rounded-full"/>
                        <Skeleton className="mx-auto h-6 w-32"/>
                        <Skeleton className="mx-auto h-4 w-24"/>
                        <div className="space-y-3 pt-4">
                            <Skeleton className="h-4 w-3/4"/>
                            <Skeleton className="h-4 w-1/2"/>
                            <Skeleton className="h-4 w-2/3"/>
                        </div>
                    </div>
                </aside>
                {/* Main */}
                <main>
                    <div className="space-y-4 rounded-xl border border-border bg-card p-6">
                        <Skeleton className="h-8 w-40"/>
                        <div className="space-y-3">
                            <Skeleton className="h-4 w-full"/>
                            <Skeleton className="h-4 w-3/4"/>
                            <Skeleton className="h-4 w-5/6"/>
                            <Skeleton className="h-4 w-2/3"/>
                        </div>
                    </div>
                </main>
            </div>
        </div>
    );
}

/* ========== Generic Content Page Skeleton ========== */

interface ContentPageSkeletonProps {
    hasHero?: boolean;
    heroSize?: 'sm' | 'md' | 'lg';
    className?: string;
}

export function ContentPageSkeleton({hasHero = true, heroSize = 'md', className}: ContentPageSkeletonProps) {
    return (
        <div className={cn('w-full', className)}>
            {hasHero && <PageHeroSkeleton size={heroSize}/>}
            <div className="w-full max-w-[1200px] mx-auto px-8 py-12 max-md:px-4">
                <div className="space-y-4">
                    <Skeleton className="h-6 w-48"/>
                    <div className="grid grid-cols-[repeat(auto-fill,minmax(280px,1fr))] gap-6">
                        {Array.from({length: 6}).map((_, i) => (
                            <div key={i} className="space-y-3 rounded-xl border border-border bg-card p-5">
                                <Skeleton className="h-32 w-full rounded-lg"/>
                                <Skeleton className="h-5 w-3/5"/>
                                <Skeleton className="h-4 w-full"/>
                                <Skeleton className="h-4 w-4/5"/>
                            </div>
                        ))}
                    </div>
                </div>
            </div>
        </div>
    );
}
