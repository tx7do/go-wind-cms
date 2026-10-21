import { useRef } from 'react';
import type { PointerEvent as ReactPointerEvent, ThHTMLAttributes } from 'react';

/** 拖拽调整列宽的下限（px），防止把一列拖没 */
export const MIN_COLUMN_WIDTH = 60;

/**
 * ListTable 经 column.onHeaderCell 注入、随 rc-table 的 additionalProps 铺到表头 th 上的
 * 拖拽参数。只有这些键会被本组件摘走，其余属性原样透传给 <th>（含排序注入的 onClick）。
 */
export interface ResizableCellProps {
  'data-resize-key'?: string;
  onResizeStart?: (info: {
    key: string;
    width: number;
    allWidths: Record<string, number>;
    containerWidth: number;
  }) => void;
  onResize?: (key: string, width: number) => void;
  onResizeEnd?: () => void;
  onResizeReset?: () => void;
  resizing?: boolean;
  resizeTitle?: string;
}

type HeaderCellProps = ResizableCellProps & ThHTMLAttributes<HTMLTableCellElement>;

/**
 * 表头单元格的列宽拖拽把手。热区是列右缘的隐形 span：pointerdown 起 pointer capture 全程
 * 接管指针（移出表头也不丢事件），按下时把当前表头行所有列的渲染宽度上报给 ListTable 做
 * 「首拖锁定」，移动中持续回报新宽度，双击把手恢复自适应。
 *
 * 不引第三方库：react-resizable 年久失修不兼容 React 19，指针事件自持足够。
 */
export function ResizableHeaderCell(props: HeaderCellProps) {
  const {
    'data-resize-key': resizeKey,
    onResizeStart,
    onResize,
    onResizeEnd,
    onResizeReset,
    resizing,
    resizeTitle,
    children,
    ...thProps
  } = props;

  const dragRef = useRef<{ key: string; startX: number; startWidth: number } | null>(null);

  const handlePointerDown = (event: ReactPointerEvent<HTMLSpanElement>) => {
    if (event.button !== 0 || !resizeKey || !onResizeStart) return;
    event.preventDefault();
    const th = event.currentTarget.closest('th');
    if (!th) return;
    // 用 getBoundingClientRect 而非 offsetWidth：自适应布局下列宽普遍带小数，取整会让
    // 锁定前后总宽漂移，凭空多出一条横向滚动条
    const startWidth = th.getBoundingClientRect().width;
    const allWidths: Record<string, number> = {};
    th.closest('thead')
      ?.querySelectorAll<HTMLTableCellElement>('th[data-resize-key]')
      .forEach((el) => {
        const key = el.getAttribute('data-resize-key');
        if (key) allWidths[key] = el.getBoundingClientRect().width;
      });
    // 滚动容器（表头的直接父级）的可见宽度，ListTable 用它判定要不要留弹性列
    const containerWidth = th.closest('thead')?.parentElement?.parentElement?.clientWidth ?? 0;
    dragRef.current = { key: resizeKey, startX: event.clientX, startWidth };
    event.currentTarget.setPointerCapture(event.pointerId);
    onResizeStart({ key: resizeKey, width: startWidth, allWidths, containerWidth });
  };

  const handlePointerMove = (event: ReactPointerEvent<HTMLSpanElement>) => {
    const drag = dragRef.current;
    if (!drag || !onResize) return;
    onResize(drag.key, Math.max(MIN_COLUMN_WIDTH, drag.startWidth + event.clientX - drag.startX));
  };

  const handlePointerEnd = (event: ReactPointerEvent<HTMLSpanElement>) => {
    if (!dragRef.current) return;
    dragRef.current = null;
    try {
      event.currentTarget.releasePointerCapture(event.pointerId);
    } catch {
      // 指针可能已随拖出窗口而提前释放捕获；纯善后，无可处理
    }
    onResizeEnd?.();
  };

  return (
    <th
      {...thProps}
      data-resize-key={resizeKey}
      className={`${thProps.className ?? ''}${resizing ? ' ant-table-cell-resizing' : ''}`}
    >
      {children}
      {resizeKey && onResizeStart && (
        <span
          className="list-table-col-resizer"
          title={resizeTitle}
          onPointerDown={handlePointerDown}
          onPointerMove={handlePointerMove}
          onPointerUp={handlePointerEnd}
          onPointerCancel={handlePointerEnd}
          // 点击/双击都发生在把手上时不冒泡：避免拖完抬手被判成一次点击、误触发同列排序
          onClick={(event) => event.stopPropagation()}
          onDoubleClick={() => onResizeReset?.()}
        />
      )}
    </th>
  );
}

export default ResizableHeaderCell;
