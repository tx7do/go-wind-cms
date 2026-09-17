// 基础组件
export { Skeleton, SkeletonBlock, SkeletonCompat } from './Skeleton';
export type {
  SkeletonProps,
  SkeletonBlockProps,
  SkeletonAnimation,
  SkeletonShape,
} from './Skeleton/types';

// 预设模板
export { PageSkeleton } from './presets/PageSkeleton';
export { TableSkeleton, type TableSkeletonProps } from './presets/TableSkeleton';
export { ListSkeleton, type ListSkeletonProps } from './presets/ListSkeleton';
export { FormSkeleton, type FormSkeletonProps } from './presets/FormSkeleton';
export { CardSkeleton, type CardSkeletonProps } from './presets/CardSkeleton';

// Hooks
export { useSkeletonAnimation } from './hooks/useSkeletonAnimation';

// 默认导出（基础 Skeleton）
export { Skeleton as default } from './Skeleton';
