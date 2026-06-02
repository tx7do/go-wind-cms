'use client';

import React from 'react';
import {XIcon} from '@/plugins/xicon';
import {cn} from '@/lib/utils';

export interface MetaDataItemProps {
    /** 图标名（XIcon 格式） */
    icon?: string;
    /** 图标尺寸 */
    iconSize?: number;
    /** 文字内容 */
    text: React.ReactNode;
    /** 自定义 className */
    className?: string;
}

/**
 * 元数据项 — icon + text 的行内组合
 *
 * 用于文章卡片、用户页、PageHero meta 等处统一的元数据展示
 */
const MetaDataItem: React.FC<MetaDataItemProps> = ({icon, iconSize = 16, text, className}) => {
    return (
        <span className={cn('flex items-center gap-1.5 text-muted-foreground', className)}>
            {icon && (
                <XIcon name={icon} size={iconSize} className="text-muted-foreground"/>
            )}
            {text}
        </span>
    );
};

export default MetaDataItem;
