import { useCallback, useEffect, useMemo, useRef, useState } from 'react';
import type { MutableRefObject } from 'react';
import { ProTable } from '@ant-design/pro-components';
import type { ActionType, ProTableProps } from '@ant-design/pro-components';
import { Alert, Button, Empty } from 'antd';
import { ReloadOutlined } from '@ant-design/icons';
import { usePreferencesStore } from '@/core/preferences';
import { useI18n } from '@/core/i18n';
import { getDefaultErrorMsg } from '@/core/transport/rest/utils';
import { TABLE } from '@/config/constants';
import { TableSkeleton } from './TableSkeleton';

/**
 * 加载态出场阈值（ms）。本地接口常在 100ms 内返回，任何早于阈值就出现的加载态（骨架屏或 Spin）
 * 都比「什么都不显示」更闪，所以首屏骨架与表内刷新的 Spin 共用这一个阈值。
 */
const LOADING_DELAY = 250;

/** pagination=false 的页面没有 pageSize 可依据，骨架行数取一个接近满屏的常量 */
const FALLBACK_SKELETON_ROWS = 10;

type ListTableProps<
  RecordType extends Record<string, any>,
  Params extends Record<string, any> = Record<string, any>,
> = ProTableProps<RecordType, Params>;

type LoadPhase = 'initial' | 'ready' | 'failed';

/**
 * 列表页统一表格：在 ProTable 之上补三件页面各自都没有的事——首屏骨架屏（按页面 columns 生成，
 * 受偏好设置 transition.loading 控制）、快响应不闪的 Spin（delay）、请求失败的错误态与重试。
 *
 * 页面把 `<ProTable>` 换成 `<ListTable>` 即可，request / columns / actionRef 用法不变。
 */
export function ListTable<
  RecordType extends Record<string, any>,
  Params extends Record<string, any> = Record<string, any>,
>(props: ListTableProps<RecordType, Params> & { striped?: boolean }) {
  const {
    request,
    columns,
    actionRef,
    pagination,
    size,
    bordered,
    scroll,
    loading,
    rowClassName,
    striped = true,
    ...restProps
  } = props;
  const { t } = useI18n('common');
  const skeletonEnabled = usePreferencesStore((state) => state.preferences.transition.loading);

  const [phase, setPhase] = useState<LoadPhase>('initial');
  /** 请求抛错时的原始文案；页面自己 catch 后返回 success:false 的情况拿不到原因，退化为通用文案 */
  const [errorReason, setErrorReason] = useState<string | null>(null);
  const [showSkeleton, setShowSkeleton] = useState(false);

  const loadedOnceRef = useRef(false);
  const skeletonTimerRef = useRef<number | null>(null);
  const ownActionRef = useRef<ActionType | undefined>(undefined);

  // 页面传进来的 request 是内联箭头函数（每次渲染都是新引用）。走 ref 取值，给 ProTable 的
  // request 才能保持恒定引用，否则包装层自己会变成重新取数的原因。
  const requestRef = useRef(request);
  requestRef.current = request;

  // ProTable 只接受一个 actionRef：这里收下它，再把同一个 action 对象转交给页面的 ref。
  // 子组件 effect 先于父组件 effect 运行，所以本页的 actionRef 在页面首次 reload 前就已就位。
  useEffect(() => {
    if (!actionRef) return;
    if (typeof actionRef === 'function') {
      actionRef(ownActionRef.current);
      return;
    }
    (actionRef as MutableRefObject<ActionType | undefined>).current = ownActionRef.current;
  });

  // 斑马纹：默认开启，对齐 ele/vben 共享表格的 stripe 默认值（值表见 docs/design-language.md §4 表格条）。
  // 页面自带 rowClassName 时叠加而非覆盖；占位行（暂无数据）由 CSS 侧排除。
  const composedRowClassName = useCallback(
    (record: RecordType, index: number, indent: number) => {
      const pageCls =
        typeof rowClassName === 'function'
          ? (rowClassName(record, index, indent) ?? '')
          : (rowClassName ?? '');
      if (!striped || index % 2 !== 1) return pageCls;
      return pageCls ? `${pageCls} ant-table-row-striped` : 'ant-table-row-striped';
    },
    [rowClassName, striped],
  );

  const stopSkeleton = useCallback(() => {
    if (skeletonTimerRef.current !== null) {
      window.clearTimeout(skeletonTimerRef.current);
      skeletonTimerRef.current = null;
    }
    setShowSkeleton(false);
  }, []);

  useEffect(() => stopSkeleton, [stopSkeleton]);

  const wrappedRequest = useCallback<NonNullable<ListTableProps<RecordType, Params>['request']>>(
    async (params, sort, filter) => {
      setErrorReason(null);
      if (!loadedOnceRef.current && skeletonEnabled) {
        if (skeletonTimerRef.current !== null) window.clearTimeout(skeletonTimerRef.current);
        skeletonTimerRef.current = window.setTimeout(() => setShowSkeleton(true), LOADING_DELAY);
      }
      try {
        const result = await requestRef.current?.(params, sort, filter);
        if (!result || result.success === false) {
          setPhase('failed');
          return result ?? { data: [], success: false };
        }
        loadedOnceRef.current = true;
        setPhase('ready');
        return result;
      } catch (error) {
        console.error('[ListTable] 列表数据请求失败', error);
        setPhase('failed');
        // 本仓错误规范 reason 优先：axios error.message 只是 "Request failed with
        // status code 500" 这类通用文案，真正的业务原因在 response.data.reason
        // （getDefaultErrorMsg 即全仓统一的提取器，getErrorMsg 的非 i18n 版）。
        setErrorReason(getDefaultErrorMsg(error));
        return { data: [], total: 0, success: false };
      } finally {
        stopSkeleton();
      }
    },
    [skeletonEnabled, stopSkeleton],
  );

  const retry = useCallback(() => {
    setPhase('initial');
    ownActionRef.current?.reload?.();
  }, []);

  const skeletonRows = useMemo(() => {
    if (pagination === false) return FALLBACK_SKELETON_ROWS;
    const configured = pagination?.defaultPageSize ?? pagination?.pageSize;
    return configured && configured > 0 ? configured : TABLE.DEFAULT_PAGE_SIZE;
  }, [pagination]);

  // 骨架继承页面的两向 scroll：scroll.x 保证列宽与真实表格逐列对齐，scroll.y 保证两态都有
  // 同一层限高 .ant-table-body —— 分页器的落点由它决定，骨架缺这层就会在数据到位时整段跳动。
  const scrollX = typeof scroll === 'object' ? scroll?.x : undefined;
  const scrollY = typeof scroll === 'object' ? scroll?.y : undefined;

  const tableViewRender = useCallback<NonNullable<ListTableProps<RecordType, Params>['tableViewRender']>>(
    (_tableProps, defaultDom) => {
      // 首屏就失败：一行数据都没有，"暂无数据"会被读成查询成功但结果为空，必须换成错误态
      if (phase === 'failed' && !loadedOnceRef.current) {
        return (
          <div className="flex flex-col items-center justify-center gap-3 py-12">
            <Empty image={Empty.PRESENTED_IMAGE_SIMPLE} description={errorReason ?? t('loading.failed')} />
            <Button icon={<ReloadOutlined />} onClick={retry}>
              {t('button.retry')}
            </Button>
          </div>
        );
      }
      if (showSkeleton) {
        return (
          <TableSkeleton<RecordType>
            columns={columns}
            rows={skeletonRows}
            size={size}
            bordered={bordered}
            scrollX={scrollX}
            scrollY={scrollY}
            pagination={pagination}
          />
        );
      }
      // 已有数据时刷新失败：旧数据留着可读，只在表格上方补一条带重试入口的提示
      if (phase === 'failed') {
        return (
          <>
            <Alert
              className="mb-2"
              type="error"
              showIcon
              title={errorReason ?? t('loading.failed')}
              action={
                <Button size="small" icon={<ReloadOutlined />} onClick={retry}>
                  {t('button.retry')}
                </Button>
              }
            />
            {defaultDom}
          </>
        );
      }
      return defaultDom;
    },
    [phase, showSkeleton, errorReason, retry, columns, skeletonRows, size, bordered, scrollX, scrollY, pagination, t],
  );

  return (
    <ProTable<RecordType, Params>
      {...restProps}
      columns={columns}
      // 这四个 prop 上面为了喂骨架/算骨架行数而被解构出来，解构即从 restProps 里摘走，
      // 不在此显式转交就会被 ProTable 静默丢掉（scroll.y 丢 → 无粘性表头与限高 body，
      // scroll.x 丢 → 列被压到容器宽度，pagination 丢 → 页面关不掉分页器）。
      pagination={pagination}
      size={size}
      bordered={bordered}
      scroll={scroll}
      actionRef={ownActionRef}
      request={wrappedRequest}
      rowClassName={composedRowClassName}
      loading={{ ...(typeof loading === 'object' ? loading : {}), delay: LOADING_DELAY }}
      tableViewRender={tableViewRender}
    />
  );
}

export default ListTable;
