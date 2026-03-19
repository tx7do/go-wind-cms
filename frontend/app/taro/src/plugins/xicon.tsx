import React from 'react';
import { Image } from '@tarojs/components';
import { Icon } from '@iconify/react'; // 只在 H5 端有效，打包时 tree-shaking

interface XIconProps {
    name: string; // H5: iconify 名称；小程序: 本地 svg 路径
    size?: number | string;
    color?: string;
    className?: string;
    style?: React.CSSProperties;
    svg?: string; // 小程序端可传 SVG 代码字符串
}

/**
 * 通用图标组件
 * H5: 使用 Iconify 实现
 * 小程序: 支持本地 SVG 路径或 SVG 代码
 */
export const XIcon: React.FC<XIconProps> = ({name, size = 24, color, className, style, svg}) => {
    const env = process.env.TARO_ENV;

    if (env === 'h5') {
        // H5 端用 Iconify
        return (
            <Icon
                icon={name}
                width={size}
                height={size}
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
    }

    // 小程序端
    if (svg) {
        // 直接渲染 SVG 代码
        // @ts-ignore
        return <svg style={{width: size, height: size, ...style}} className={className} dangerouslySetInnerHTML={{__html: svg}} />;
    }
    // 渲染本地 SVG 文件
    return (
        <Image
            src={name}
            style={{width: size, height: size, ...style}}
            className={className}
            mode="aspectFit"
        />
    );
};

// 默认导出
export default XIcon;
