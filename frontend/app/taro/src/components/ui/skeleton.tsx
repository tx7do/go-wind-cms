import {View, ViewProps} from '@tarojs/components';
import {cn} from '@/lib/utils';

/**
 * 骨架屏占位块
 * 使用 splitLine 色作为 shimmer 背景
 */
function Skeleton({className, ...props}: ViewProps) {
    return (
        <View
          className={cn('animate-pulse rounded bg-splitLine', className)}
          {...props}
        />
    );
}

export {Skeleton};
