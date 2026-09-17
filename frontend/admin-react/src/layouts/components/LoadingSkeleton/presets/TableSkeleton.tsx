import { useMemo } from 'react';
import { Skeleton } from '../Skeleton';
import type { SkeletonPresetProps } from '../Skeleton/types';

export interface TableSkeletonProps extends SkeletonPresetProps {
  /** 列数 */
  columns?: number;

  /** 行数 */
  rows?: number;

  /** 列宽分布（百分比数组） */
  columnWidths?: (string | number)[];

  /** 是否显示操作列 */
  withActions?: boolean;

  /** 是否显示复选框列 */
  withSelection?: boolean;
}

/**
 * 表格骨架屏
 * 适用：ProTable / Antd Table 加载状态
 */
export const TableSkeleton = ({
  active = true,
  delay = 150,
  className = '',
  columns = 5,
  rows = 6,
  columnWidths,
  withActions = true,
  withSelection = true,
}: TableSkeletonProps) => {
  // 计算列配置
  const colConfigs = useMemo(() => {
    const configs: { width: string | number; isAction?: boolean; isSelection?: boolean }[] = [];

    if (withSelection) {
      configs.push({ width: 40, isSelection: true });
    }

    const contentCols = columns - (withSelection ? 1 : 0) - (withActions ? 1 : 0);
    const defaultWidth = `${100 / columns}%`;

    for (let i = 0; i < contentCols; i++) {
      configs.push({
        width: columnWidths?.[i] || defaultWidth,
      });
    }

    if (withActions) {
      configs.push({ width: 100, isAction: true });
    }

    return configs;
  }, [columns, columnWidths, withActions, withSelection]);

  return (
    <div className={`table-skeleton ${className}`.trim()}>
      {/* 表头 */}
      <div className="flex items-center p-3 border-b bg-[var(--ant-color-fill-tertiary)]">
        {colConfigs.map((col, i) => (
          <Skeleton
            key={`header-${i}`}
            active={active}
            delay={delay}
            width={col.width}
            height="1em"
            className={`mx-1 ${col.isAction ? 'text-right' : ''}`}
          />
        ))}
      </div>

      {/* 表体 */}
      {Array.from({ length: rows }).map((_, rowIndex) => (
        <div key={rowIndex} className="flex items-center p-3 border-b last:border-0">
          {colConfigs.map((col, colIndex) => (
            <Skeleton
              key={`row-${rowIndex}-col-${colIndex}`}
              active={active}
              delay={delay + 100 + rowIndex * 40 + colIndex * 20}
              width={col.width}
              height="1em"
              className={`mx-1 ${col.isSelection ? 'rounded-full' : ''} ${col.isAction ? 'text-right' : ''}`}
            />
          ))}
        </div>
      ))}

      {/* 分页区域 */}
      <div className="flex items-center justify-end p-3 gap-2">
        <Skeleton active={active} delay={delay + 500} width={100} height={32} />
        <Skeleton active={active} delay={delay + 520} width={60} height={32} />
        <Skeleton active={active} delay={delay + 540} width={60} height={32} />
      </div>
    </div>
  );
};

export default TableSkeleton;
