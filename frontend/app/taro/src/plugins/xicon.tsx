import React from 'react';
import {Text} from '@tarojs/components';
import {iconGlyphs} from './generated/icon-glyphs';
import './generated/iconfont.css';

/**
 * 图标一律走构建期裁出来的图标字体（node scripts/gen-icons.mjs），不走 svg：
 * 小程序端 svg 根本进不了渲染树（Taro 的模板选择只认固定节点名白名单），
 * 若 H5 留 svg、小程序走字体，预览面和出货面就不是同一套字形。
 * 字体的好处是颜色仍由 color 决定，调用点上的 text-primary 之类照旧生效。
 */
const isH5 = process.env.TARO_ENV === 'h5';

interface XIconProps {
    name: string;
    size?: number | string;
    color?: string;
    className?: string;
    style?: React.CSSProperties;
}

/**
 * 通用图标组件
 * name 用 '<库>:<图标名>'，目前只装了 carbon 的子集
 */
export const XIcon: React.FC<XIconProps> = ({name, size = 24, color, className, style}) => {
    // 图标名写错在生成阶段就会 fail，这里只兜「加了新图标没重跑生成脚本」
    const glyph = iconGlyphs[name];
    if (!glyph) return null;
    const px = Number(size) || 24;
    // 内联样式里的 rpx 不会被 postcss 改写，H5 只认 px；小程序按 750 设计稿 ×2
    const fontSize = isH5 ? `${px}px` : `${px * 2}rpx`;

    return (
        <Text
          className={`gw-icon ${className || ''}`}
          style={{fontSize, color, ...style}}
        >
            {glyph}
        </Text>
    );
};

// 默认导出
export default XIcon;
