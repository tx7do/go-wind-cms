'use client';

import React from 'react';
import {ArrowLeft} from 'lucide-react';

import {cn} from '@/lib/utils';

export interface BackButtonProps {
    /** 按钮文字 */
    label: string;
    /** 点击回调 */
    onClick: () => void;
    /** 自定义 className */
    className?: string;
}

/**
 * 返回按钮 — 用于详情页面的返回导航
 *
 * 统一了 post/[id]、category/[id]、tag/[id] 的返回按钮样式
 */
const BackButton: React.FC<BackButtonProps> = ({label, onClick, className}) => {
    return (
        <button
            onClick={onClick}
            className={cn(
                'group flex cursor-pointer items-center gap-2 rounded-lg px-3 py-1.5',
                'text-sm text-muted-foreground transition-colors',
                'hover:bg-muted/60 hover:text-foreground',
                className,
            )}
        >
            <ArrowLeft className="h-4 w-4 transition-transform duration-300 group-hover:-translate-x-1"/>
            <span>{label}</span>
        </button>
    );
};

export default BackButton;
