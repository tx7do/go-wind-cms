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
    const minHeight = size === 'sm' ? 'min-h-[220px]' : size === 'lg' ? 'min-h-[340px]' : 'min-h-[280px]';
    const padding = size === 'sm' ? 'py-14' : size === 'lg' ? 'py-24' : 'py-18';

    return (
        <section className={cn(
            'relative w-full overflow-hidden border-b border-border',
            minHeight,
            padding,
            'flex items-center justify-center text-center',
        )}>
            {/* 背景层：风的流体渐变 */}
            <div className="absolute inset-0 -z-10">
                {/* 基础渐变：倾斜风的轨迹 */}
                <div
                    className={cn(
                        'absolute inset-0',
                    )}
                    style={{
                        background: accentColor
                            ? `linear-gradient(135deg, ${accentColor}22 0%, hsl(var(--background)) 45%, ${accentColor}10 100%)`
                            : 'linear-gradient(135deg, hsl(var(--primary) / 0.10) 0%, hsl(var(--background)) 45%, hsl(var(--primary) / 0.05) 100%)',
                    }}
                />

                {/* 风的流体光晕：双径向层 */}
                <div
                    className="absolute inset-0 opacity-70"
                    style={{
                        background: `radial-gradient(ellipse at 25% 35%, ${accentColor || 'hsl(var(--primary))'}20 0%, transparent 55%), radial-gradient(ellipse at 75% 65%, rgba(56, 189, 248, 0.08) 0%, transparent 55%)`,
                    }}
                />

                {/* 风的流向斜线纹理 */}
                <div
                    className={cn(
                        'absolute inset-0 opacity-[0.05]',
                        'dark:opacity-[0.08]',
                    )}
                    style={{
                        backgroundImage: `repeating-linear-gradient(135deg, ${accentColor || 'hsl(var(--primary))'} 0, ${accentColor || 'hsl(var(--primary))'} 1px, transparent 1px, transparent 48px)`,
                        maskImage: 'radial-gradient(ellipse at center, black 20%, transparent 75%)',
                        WebkitMaskImage: 'radial-gradient(ellipse at center, black 20%, transparent 75%)',
                    }}
                />

                {/* 点阵装饰背景：科技门户层次感 */}
                <div
                    className={cn(
                        'absolute inset-0 opacity-[0.08]',
                        'dark:opacity-[0.12]',
                    )}
                    style={{
                        backgroundImage: `radial-gradient(${accentColor || 'hsl(var(--primary))'} 1px, transparent 1px)`,
                        backgroundSize: '24px 24px',
                        maskImage: 'radial-gradient(ellipse at center, black 10%, transparent 65%)',
                        WebkitMaskImage: 'radial-gradient(ellipse at center, black 10%, transparent 65%)',
                    }}
                />

                {/* 底部风迹光带：变宽为渐变带，更流体 */}
                <div
                    className="absolute bottom-0 left-0 right-0 h-0.5 opacity-40"
                    style={{
                        background: `linear-gradient(90deg, transparent, ${accentColor || 'hsl(var(--primary))'} 30%, #38bdf8 70%, transparent)`,
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

                {/* 主标题 */}
                <h1 className={cn(
                    'mb-4 font-bold leading-tight text-foreground',
                    size === 'lg' ? 'text-4xl max-md:text-2xl' : 'text-3xl max-md:text-2xl',
                    size === 'sm' && 'text-2xl max-md:text-xl',
                )}>
                    {title}
                </h1>

                {/* 描述：弱化一阶，拉开与主标题的字重对比 */}
                {description && (
                    <p className={cn(
                        'mx-auto text-muted-foreground/80 font-light',
                        size === 'lg' ? 'max-w-2xl text-lg max-md:text-base' : 'max-w-xl text-base max-md:text-sm',
                        size === 'sm' && 'text-sm max-md:text-xs',
                    )}>
                        {description}
                    </p>
                )}

                {/* 元信息 */}
                {meta && (
                    <div className={cn(
                        'mt-5 flex flex-wrap items-center justify-center gap-4',
                        'text-sm text-muted-foreground',
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
