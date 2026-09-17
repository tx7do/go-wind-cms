import type { ReactNode } from 'react';

export type SkeletonAnimation = 'wave' | 'pulse' | 'none';

export type SkeletonShape = 'round' | 'square' | 'circle';

export interface SkeletonBlockProps {
  /** 宽度（支持像素/百分比） */
  width?: string | number;

  /** 高度（支持像素/百分比） */
  height?: string | number;

  /** 形状 */
  shape?: SkeletonShape;

  /** 动画类型 */
  animation?: SkeletonAnimation;

  /** 自定义类名 */
  className?: string;

  /** 自定义样式 */
  style?: React.CSSProperties;

  /** 子内容（用于嵌套） */
  children?: ReactNode;
}

export interface SkeletonProps extends SkeletonBlockProps {
  /** 是否显示（控制显隐，避免卸载重挂载） */
  active?: boolean;

  /** 延迟显示时间（避免闪烁，单位：ms） */
  delay?: number;

  /** 骨架块配置数组 */
  blocks?: SkeletonBlockProps[];

  /** 预设模板名称 */
  preset?: 'paragraph' | 'avatar' | 'button' | 'input' | 'image';

  /** 段落行数（preset='paragraph' 时生效） */
  rows?: number;

  /** 段落宽度分布（百分比数组，总和建议 <= 100） */
  rowWidths?: (string | number)[];
}

export interface SkeletonPresetProps {
  /** 是否激活动画 */
  active?: boolean;

  /** 延迟显示（避免闪烁） */
  delay?: number;

  /** 自定义类名 */
  className?: string;
}
