import { Skeleton } from '../Skeleton';
import type { SkeletonPresetProps } from '../Skeleton/types';

/**
 * 页面级骨架屏
 * 适用：通用页面加载状态，包含标题 + 内容区域 + 操作区
 */
export const PageSkeleton = ({
  active = true,
  delay = 200,
  className = '',
}: SkeletonPresetProps) => {
  return (
    <div className={`page-skeleton ${className}`.trim()}>
      {/* 标题区域 */}
      <div className="flex items-center justify-between mb-4">
        <Skeleton
          active={active}
          delay={delay}
          preset="paragraph"
          rows={1}
          rowWidths={['40%']}
          className="h-6"
        />
        <div className="flex gap-2">
          <Skeleton active={active} delay={delay} width={60} height={32} />
          <Skeleton active={active} delay={delay} width={60} height={32} />
        </div>
      </div>

      {/* 内容区域 */}
      <div className="space-y-4">
        {/* 统计卡片行 */}
        <div className="grid grid-cols-1 md:grid-cols-4 gap-4">
          {Array.from({ length: 4 }).map((_, i) => (
            <div key={i} className="p-4 border rounded-lg">
              <Skeleton
                active={active}
                delay={delay + i * 50}
                preset="paragraph"
                rows={1}
                rowWidths={['60%']}
                className="mb-2"
              />
              <Skeleton active={active} delay={delay + i * 50} width="80%" height="2em" />
            </div>
          ))}
        </div>

        {/* 表格/列表区域 */}
        <div className="border rounded-lg overflow-hidden">
          {/* 表头 */}
          <div className="flex p-3 border-b bg-[var(--ant-color-fill-tertiary)]">
            {Array.from({ length: 5 }).map((_, i) => (
              <Skeleton
                key={i}
                active={active}
                delay={delay + 200 + i * 30}
                width={`${100 / 5}%`}
                height="1em"
                className="mx-1"
              />
            ))}
          </div>

          {/* 表体 */}
          {Array.from({ length: 5 }).map((_, rowIndex) => (
            <div key={rowIndex} className="flex p-3 border-b last:border-0">
              {Array.from({ length: 5 }).map((_, colIndex) => (
                <Skeleton
                  key={colIndex}
                  active={active}
                  delay={delay + 300 + rowIndex * 50 + colIndex * 20}
                  width={`${100 / 5}%`}
                  height="1em"
                  className="mx-1"
                />
              ))}
            </div>
          ))}
        </div>
      </div>
    </div>
  );
};

export default PageSkeleton;
