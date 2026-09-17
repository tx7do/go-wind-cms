/**
 * ContentContainer 组件
 * 统一页面内容区域的布局和样式
 */

import React from 'react';
import './ContentContainer.style.css';

export interface ContentContainerProps {
  /** 子组件 */
  children: React.ReactNode;
  /** 内边距，默认 16px */
  padding?: string | number;
  /** 是否显示滚动条，默认 false */
  scrollable?: boolean;
  /** 
   * 高度模式：
   * - 'auto': 自适应高度，内容有多高就撑多高
   * - 'fixed': 固定高度，占满可视区域（默认）
   */
  heightMode?: 'auto' | 'fixed';
  /** 底部外边距，仅在 heightMode='fixed' 时生效，默认 0 */
  bottomMargin?: string | number;
  /** 自定义样式 */
  style?: React.CSSProperties;
  /** 自定义类名 */
  className?: string;
}

/**
 * 内容容器组件
 * 
 * 功能：
 * - 统一页面内容区域的布局
 * - 支持两种高度模式：自适应（auto）和固定（fixed）
 * - 提供一致的内边距
 * - 可选的滚动支持
 * - 动态高度计算
 * 
 * 使用示例：
 * ```tsx
 * // 方式 1：固定高度（默认）
 * <ContentContainer heightMode="fixed" bottomMargin={16}>
 *   <Card>内容</Card>
 * </ContentContainer>
 * 
 * // 方式 2：自适应高度
 * <ContentContainer heightMode="auto">
 *   <Card>内容</Card>
 * </ContentContainer>
 * ```
 */
const ContentContainer: React.FC<ContentContainerProps> = ({
  children,
  padding = '16px',
  scrollable = false,
  heightMode = 'fixed',
  bottomMargin = 0,
  style,
  className,
}) => {
  const containerClassName = [
    'content-container',
    heightMode === 'auto' && 'content-container--auto',
    heightMode === 'fixed' && 'content-container--fixed',
    scrollable && 'content-container--scrollable',
    className,
  ]
    .filter(Boolean)
    .join(' ');

  const bottomMarginValue = typeof bottomMargin === 'number' ? `${bottomMargin}px` : bottomMargin;

  return (
    <div
      className={containerClassName}
      style={{
        // 使用内联样式确保不被外部 CSS 覆盖
        display: 'flex',
        flexDirection: 'column',
        flex: heightMode === 'fixed' ? 1 : undefined,
        height: heightMode === 'auto' ? 'auto' : undefined,
        minHeight: 0,
        // scrollable 必须在内联样式里生效：内联优先级高于 CSS class，
        // 此前 overflow 只看 heightMode，导致 --scrollable 类的 overflow:auto
        // 永远被内联覆盖，scrollable prop 实际是死功能（页面长内容无滚动条）。
        overflow: scrollable
          ? 'auto'
          : heightMode === 'fixed'
            ? 'hidden'
            : 'visible',
        // padding 和 margin 由 props 控制
        padding: typeof padding === 'number' ? `${padding}px` : padding,
        marginBottom: heightMode === 'fixed' ? bottomMarginValue : undefined,
        boxSizing: 'border-box',
        ...style,
      }}
    >
      {children}
    </div>
  );
};

export default ContentContainer;
