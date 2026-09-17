import { Skeleton } from '../Skeleton';
import type { SkeletonPresetProps } from '../Skeleton/types';

export interface CardSkeletonProps extends SkeletonPresetProps {
  /** 是否显示封面图 */
  withCover?: boolean;

  /** 是否显示头像（配合封面） */
  withAvatar?: boolean;

  /** 是否显示操作按钮 */
  withActions?: boolean;
}

/**
 * 卡片骨架屏
 * 适用：Antd Card / 自定义卡片列表
 */
export const CardSkeleton = ({
  active = true,
  delay = 150,
  className = '',
  withCover = true,
  withAvatar = false,
  withActions = true,
}: CardSkeletonProps) => {
  return (
    <div className={`card-skeleton border rounded-lg overflow-hidden ${className}`.trim()}>
      {/* 封面图 */}
      {withCover && (
        <Skeleton
          active={active}
          delay={delay}
          width="100%"
          height={160}
          shape="square"
          className="block"
        />
      )}

      {/* 内容区域 */}
      <div className="p-4 space-y-3">
        {/* 头像 + 标题（可选） */}
        {withAvatar && (
          <div className="flex items-center gap-3">
            <Skeleton active={active} delay={delay + 50} width={40} height={40} shape="circle" />
            <Skeleton active={active} delay={delay + 70} width="60%" height="1.2em" />
          </div>
        )}

        {/* 标题 */}
        <Skeleton active={active} delay={delay + 100} width="80%" height="1.3em" />

        {/* 描述 */}
        <div className="space-y-1.5">
          <Skeleton active={active} delay={delay + 120} width="100%" height="1em" />
          <Skeleton active={active} delay={delay + 130} width="90%" height="1em" />
          <Skeleton active={active} delay={delay + 140} width="70%" height="1em" />
        </div>

        {/* 操作按钮 */}
        {withActions && (
          <div className="flex gap-2 pt-2">
            <Skeleton active={active} delay={delay + 180} width={60} height={28} />
            <Skeleton active={active} delay={delay + 190} width={60} height={28} />
          </div>
        )}
      </div>
    </div>
  );
};

export default CardSkeleton;
