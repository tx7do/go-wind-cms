'use client';

import React from 'react';
import {Inbox} from 'lucide-react';

import {cn} from '@/lib/utils';

interface AppEmptyProps {
    /**
     * 空状态类型
     * - default: 默认空状态
     * - error: 错误状态 (如无效的 ID)
     * - noData: 无数据状态
     */
    variant?: 'default' | 'error' | 'noData';
    /**
     * 是否显示在页面容器中 (添加额外的 margin)
     */
    inContainer?: boolean;
    /**
     * 自定义图标
     */
    image?: React.ReactNode;
    /**
     * 描述文本
     */
    description?: React.ReactNode;
    /**
     * 自定义类名
     */
    className?: string;
    /**
     * 自定义样式
     */
    style?: React.CSSProperties;
}

/**
 * App 级别的 Empty 组件（shadcn/ui 风格，无 antd 依赖）
 */
export const AppEmpty: React.FC<AppEmptyProps> = ({
    variant = 'default',
    inContainer = false,
    image,
    description,
    className,
    style,
}) => {
    const defaultIcon = image || (
        variant === 'error'
            ? <Inbox className="h-16 w-16 text-muted-foreground" strokeWidth={1}/>
            : variant === 'noData'
                ? <Inbox className="h-16 w-16 text-muted-foreground" strokeWidth={1}/>
                : <Inbox className="h-16 w-16 text-muted-foreground" strokeWidth={1}/>
    );

    return (
        <div
            className={cn(
                'flex w-full items-center justify-center gap-4 py-12 px-5',
                inContainer && 'my-20',
                className,
            )}
            style={style}
        >
            <div className="opacity-50">{defaultIcon}</div>
            {description && (
                <p className="text-sm text-muted-foreground">{description}</p>
            )}
        </div>
    );
};

export default AppEmpty;
