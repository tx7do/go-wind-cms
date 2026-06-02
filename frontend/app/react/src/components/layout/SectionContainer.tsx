import React from 'react';
import {cn} from '@/lib/utils';

export interface SectionContainerProps {
    children: React.ReactNode;
    className?: string;
    /** 最大宽度 */
    maxWidth?: 'default' | 'wide' | 'narrow';
    /** 是否保留垂直 padding */
    noPadding?: boolean;
}

/**
 * 页面内容区通用容器
 *
 * 替代重复的 `w-full max-w-[1200px] mx-auto px-8 py-12 max-md:px-4` 模式
 *
 * - default: max-w-[1200px] (标准页面)
 * - wide: max-w-300 (首页 sections)
 * - narrow: max-w-3xl (法定页面、窄文章)
 */
const SectionContainer: React.FC<SectionContainerProps> = ({
    children,
    className,
    maxWidth = 'default',
    noPadding = false,
}) => {
    const widthClass = maxWidth === 'wide'
        ? 'max-w-300'
        : maxWidth === 'narrow'
            ? 'max-w-3xl'
            : 'max-w-[1200px]';

    return (
        <div className={cn(
            'w-full mx-auto px-8 max-md:px-4',
            !noPadding && 'py-12',
            widthClass,
            className,
        )}>
            {children}
        </div>
    );
};

export default SectionContainer;
