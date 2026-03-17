import React from 'react';
import {Icon} from '@iconify/react';

interface XIconProps {
    name: string;
    size?: number | string;
    color?: string;
    className?: string;
    style?: React.CSSProperties;
}

/**
 * 通用图标组件
 * 使用 Iconify 实现，支持 carbon、mdi 等所有 Iconify 图标库
 */
export const XIcon: React.FC<XIconProps> = ({name, size = 24, color, className, style}) => {
    return (
        <Icon
            icon={name}
            width={typeof size === 'string' ? size : size}
            height={typeof size === 'string' ? size : size}
            color={color}
            className={className}
            style={{
                display: 'inline-flex',
                alignItems: 'center',
                justifyContent: 'center',
                ...style,
            }}
        />
    );
};

// 默认导出
export default XIcon;
