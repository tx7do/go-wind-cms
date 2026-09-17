import { Skeleton } from '../Skeleton';
import type { SkeletonPresetProps } from '../Skeleton/types';

export interface ListSkeletonProps extends SkeletonPresetProps {
  /** 列表项数 */
  itemCount?: number;

  /** 每项是否显示头像 */
  withAvatar?: boolean;

  /** 每项是否显示描述 */
  withDescription?: boolean;

  /** 每项是否显示操作按钮 */
  withActions?: boolean;
}

/**
 * 列表骨架屏
 * 适用：Antd List / 自定义列表组件
 */
export const ListSkeleton = ({
  active = true,
  delay = 100,
  className = '',
  itemCount = 5,
  withAvatar = true,
  withDescription = true,
  withActions = false,
}: ListSkeletonProps) => {
  return (
    <div className={`list-skeleton space-y-3 ${className}`.trim()}>
      {Array.from({ length: itemCount }).map((_, index) => (
        <div key={index} className="flex items-start gap-3 p-3 border rounded-lg">
          {/* 头像 */}
          {withAvatar && (
            <Skeleton
              active={active}
              delay={delay + index * 50}
              width={40}
              height={40}
              shape="circle"
            />
          )}

          {/* 内容区域 */}
          <div className="flex-1 min-w-0">
            {/* 标题 */}
            <Skeleton
              active={active}
              delay={delay + index * 50 + 20}
              width="60%"
              height="1.2em"
              className="mb-2"
            />

            {/* 描述 */}
            {withDescription && (
              <div className="space-y-1.5">
                <Skeleton
                  active={active}
                  delay={delay + index * 50 + 40}
                  width="90%"
                  height="1em"
                />
                <Skeleton
                  active={active}
                  delay={delay + index * 50 + 50}
                  width="70%"
                  height="1em"
                />
              </div>
            )}

            {/* 操作按钮 */}
            {withActions && (
              <div className="flex gap-2 mt-3">
                <Skeleton active={active} delay={delay + index * 50 + 70} width={50} height={24} />
                <Skeleton active={active} delay={delay + index * 50 + 80} width={50} height={24} />
              </div>
            )}
          </div>
        </div>
      ))}
    </div>
  );
};

export default ListSkeleton;
