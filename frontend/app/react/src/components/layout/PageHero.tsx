'use client';

import React from 'react';
import {XIcon} from '@/plugins/xicon';
import {cn} from '@/lib/utils';

interface PageHeroProps {
    title: string;
    subtitle?: string;
    description?: string;
    icon?: string;
    iconSize?: number;
    meta?: React.ReactNode;
    accentColor?: string;
    children?: React.ReactNode;
    size?: 'sm' | 'md' | 'lg';
}

/**
 * 通用页面 Hero 组件
 * - sm: 紧凑型，适用于法律页面（隐私、条款等）
 * - md: 标准型，适用于列表页（Post、Category、Tag）
 * - lg: 大型，适用于详情页（分类详情、标签详情）
 */
const PageHero: React.FC<PageHeroProps> = ({
    title,
    subtitle,
    description,
    icon,
    iconSize = 48,
    meta,
    accentColor,
    children,
    size = 'md',
}) => {
    const minHeight = size === 'sm' ? 'min-h-[280px]' : size === 'lg' ? 'min-h-[380px]' : 'min-h-[340px]';
    const padding = size === 'sm' ? 'py-16' : size === 'lg' ? 'py-28' : 'py-20';

    return (
        <section className={cn(
            'relative w-full overflow-hidden',
            'dark:bg-gradient-to-b dark:from-slate-950 dark:via-slate-900/80 dark:to-background',
            'bg-gradient-to-b from-[hsl(142,70%,96%)] via-[hsl(142,50%,98%)] to-background',
            minHeight,
            padding,
            'flex items-center justify-center text-center',
        )}>
            {/* 装饰线条层：低透明度风轨纹理 */}
            <div className="absolute inset-0 opacity-[0.03] pointer-events-none">
                {/* 水平流动线条 */}
                <div
                    className="absolute inset-0"
                    style={{
                        backgroundImage: `repeating-linear-gradient(0deg, hsl(var(--primary)) 0, hsl(var(--primary)) 1px, transparent 1px, transparent 60px)`,
                    }}
                />
                {/* 斜向风轨线条 */}
                <div
                    className="absolute inset-0"
                    style={{
                        backgroundImage: `repeating-linear-gradient(135deg, hsl(var(--primary)) 0, hsl(var(--primary)) 1px, transparent 1px, transparent 80px)`,
                    }}
                />
            </div>


            {/* 内容层 */}
            <div className={cn(
                'relative z-10 w-full px-8 max-md:px-4',
                size === 'lg' ? 'max-w-3xl' : 'max-w-3xl',
            )}>
                {/* 图标 */}
                {icon && (
                    <div
                        className="mb-4 flex items-center justify-center"
                        style={{color: accentColor || 'hsl(var(--primary))'}}
                    >
                        <XIcon name={icon} size={iconSize}/>
                    </div>
                )}

                {/* 副标题 */}
                {subtitle && (
                    <p className={cn(
                        'mb-2 text-sm font-medium uppercase tracking-[0.2em]',
                        size === 'sm' && 'text-xs',
                    )} style={{color: accentColor || 'hsl(var(--primary))'}}>
                        {subtitle}
                    </p>
                )}

                {/* 主标题：深色模式白字，亮色模式深色字 */}
                <h1 className={cn(
                    'mb-4 font-bold leading-tight dark:text-white text-foreground',
                    size === 'lg' ? 'text-4xl max-md:text-2xl' : 'text-3xl max-md:text-2xl',
                    size === 'sm' && 'text-2xl max-md:text-xl',
                )}>
                    {title}
                </h1>

                {/* 描述：深色模式浅灰，亮色模式 muted */}
                {description && (
                    <p className={cn(
                        'mx-auto font-light dark:text-slate-300 text-muted-foreground',
                        size === 'lg' ? 'max-w-2xl text-lg max-md:text-base' : 'max-w-xl text-base max-md:text-sm',
                        size === 'sm' && 'text-sm max-md:text-xs',
                    )}>
                        {description}
                    </p>
                )}

                {/* 元信息：深色模式浅灰，亮色模式 muted */}
                {meta && (
                    <div className={cn(
                        'mt-5 flex flex-wrap items-center justify-center gap-4',
                        'text-sm dark:text-slate-400 text-muted-foreground',
                    )}>
                        {meta}
                    </div>
                )}

                {/* 自定义内容 */}
                {children}
            </div>
        </section>
    );
};

export default PageHero;
