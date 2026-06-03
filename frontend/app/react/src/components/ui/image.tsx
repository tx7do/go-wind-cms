'use client';

import React, { useState } from 'react';

const DEFAULT_FALLBACK = '/placeholder.png';

export interface ImageProps extends React.ImgHTMLAttributes<HTMLImageElement> {
    /** 图片加载失败时的回退地址，默认 /placeholder.png */
    fallbackSrc?: string;
}

/**
 * 公共 Image 组件
 * - 透传所有原生 <img> 属性
 * - 内置 onError 回退至占位图，避免破图
 */
const Image: React.FC<ImageProps> = ({ fallbackSrc = DEFAULT_FALLBACK, onError, src, ...rest }) => {
    const [hasError, setHasError] = useState(false);

    const handleError = (e: React.SyntheticEvent<HTMLImageElement>) => {
        const img = e.currentTarget;
        // 防止死循环：如果回退图也失败则不再替换
        if (img.src !== fallbackSrc) {
            img.src = fallbackSrc;
            setHasError(true);
        }
        onError?.(e);
    };

    return (
        <img
            src={hasError ? fallbackSrc : src}
            onError={handleError}
            {...rest}
        />
    );
};

export default Image;
