import { Skeleton } from '../Skeleton';
import type { SkeletonPresetProps } from '../Skeleton/types';

export interface FormSkeletonProps extends SkeletonPresetProps {
  /** 表单项数 */
  fieldCount?: number;

  /** 是否显示标签 */
  withLabels?: boolean;

  /** 是否显示表单项描述 */
  withHelp?: boolean;

  /** 是否显示提交按钮 */
  withSubmit?: boolean;
}

/**
 * 表单骨架屏
 * 适用：表单加载/配置动态表单场景
 */
export const FormSkeleton = ({
  active = true,
  delay = 100,
  className = '',
  fieldCount = 4,
  withLabels = true,
  withHelp = false,
  withSubmit = true,
}: FormSkeletonProps) => {
  return (
    <div className={`form-skeleton space-y-4 ${className}`.trim()}>
      {Array.from({ length: fieldCount }).map((_, index) => (
        <div key={index} className="space-y-2">
          {/* 标签 */}
          {withLabels && (
            <Skeleton active={active} delay={delay + index * 60} width={80} height="1em" />
          )}

          {/* 输入框 */}
          <Skeleton
            active={active}
            delay={delay + index * 60 + 20}
            width="100%"
            height={40}
            shape="round"
          />

          {/* 帮助文本 */}
          {withHelp && (
            <Skeleton
              active={active}
              delay={delay + index * 60 + 40}
              width="70%"
              height="0.9em"
              className="text-gray-400"
            />
          )}
        </div>
      ))}

      {/* 提交按钮区域 */}
      {withSubmit && (
        <div className="flex justify-end gap-2 pt-4 border-t">
          <Skeleton active={active} delay={delay + 300} width={80} height={36} />
          <Skeleton active={active} delay={delay + 320} width={80} height={36} />
        </div>
      )}
    </div>
  );
};

export default FormSkeleton;
