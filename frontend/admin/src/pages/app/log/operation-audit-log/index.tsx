import { useRef } from 'react';
import type { ProColumns, ActionType } from '@ant-design/pro-components';
import dayjs, { type Dayjs } from 'dayjs';
import ListTable from '@/components/common/ListTable';
import { DownloadOutlined } from '@ant-design/icons';
import { Button, Dropdown, Tag, App } from 'antd';
import { useTranslation } from 'react-i18next';
import type { auditservicev1_OperationAuditLog as OperationAuditLog } from '@/api/generated/admin/service/v1';
import { PaginationQuery } from '@/core';
import { TABLE } from '@/config/constants';
import { fetchListOperationAuditLogs } from '@/api/hooks/operation-audit-log';
import { exportAuditLogsServer } from '@/api/hooks/audit-export';
import { exportAuditLogs, type AuditExportFormat } from '@/utils/csv';
import { useProTableScrollY } from '@/hooks/useProTableScrollY';
import ContentContainer from '@/layouts/components/PageContainer/ContentContainer';
import {
  getActionMap,
  getActionOptions,
  getSuccessStatusList,
  successToColor,
  successToName,
} from './constants';

/**
 * 操作审计日志页面
 */
const OperationAuditLogPage = () => {
  const { t } = useTranslation('operation-audit-log');
  const actionRef = useRef<ActionType>(null);
  // 最近一次列表查询参数：导出 CSV 时沿用当前搜索/排序条件
  const latestParamsRef = useRef<Record<string, any>>({});
  const { message } = App.useApp();

  const containerRef = useRef<HTMLDivElement>(null);
  const tableScrollY = useProTableScrollY(containerRef);

  const actionMap = getActionMap(t);
  const actionOptions = getActionOptions(t);

  const columns: ProColumns<OperationAuditLog>[] = [
    {
      title: t('serial'),
      dataIndex: 'id',
      width: 60,
      hideInSearch: true,
      render: (_, _record, index) => {
        const pagination = actionRef.current?.pageInfo;
        const page = pagination?.current || 1;
        const pageSize = pagination?.pageSize || 20;
        return (page - 1) * pageSize + index + 1;
      },
    },
    {
      title: t('createdAt'),
      dataIndex: 'createdAt',
      width: 170,
      valueType: 'dateTime',
      hideInSearch: true,
      sorter: true,
    },

    {
      title: t('createdAt'),
      dataIndex: 'createdAt',
      valueType: 'dateTimeRange',
      hideInTable: true,
      search: {
        transform: (values?: Array<string | Dayjs>) => {
          const fmt = (v?: string | Dayjs) =>
            typeof v === 'string' ? v : v ? dayjs(v).format('YYYY-MM-DD HH:mm:ss') : undefined;
          return {
            created_at__gte: fmt(values?.[0]),
            created_at__lte: fmt(values?.[1]),
          };
        },
      },
      fieldProps: {
        showTime: true,
        presets: [
          {
            label: t('common:dateRange.today'),
            value: [dayjs().startOf('day'), dayjs().endOf('day')],
          },
          {
            label: t('common:dateRange.yesterday'),
            value: [
              dayjs().subtract(1, 'day').startOf('day'),
              dayjs().subtract(1, 'day').endOf('day'),
            ],
          },
          {
            label: t('common:dateRange.thisWeek'),
            value: [dayjs().startOf('week'), dayjs().endOf('week')],
          },
          {
            label: t('common:dateRange.lastWeek'),
            value: [
              dayjs().subtract(1, 'week').startOf('week'),
              dayjs().subtract(1, 'week').endOf('week'),
            ],
          },
          {
            label: t('common:dateRange.thisMonth'),
            value: [dayjs().startOf('month'), dayjs().endOf('month')],
          },
          {
            label: t('common:dateRange.lastMonth'),
            value: [
              dayjs().subtract(1, 'month').startOf('month'),
              dayjs().subtract(1, 'month').endOf('month'),
            ],
          },
        ],
      },
    },
    {
      title: t('success'),
      dataIndex: 'success',
      width: 100,
      valueType: 'select',
      valueEnum: Object.fromEntries(
        getSuccessStatusList(t).map((item) => [
          item.value,
          { text: item.label, status: item.value === 'true' ? 'Success' : 'Error' },
        ]),
      ),
      render: (_, record) => (
        <Tag color={successToColor(record.success)}>{successToName(t, record.success)}</Tag>
      ),
    },
    {
      title: t('action'),
      dataIndex: 'action',
      width: 100,
      valueType: 'select',
      fieldProps: {
        options: actionOptions,
      },
      render: (_, record) => {
        const config = actionMap[record.action as keyof typeof actionMap];
        if (!config) return record.action || '-';
        return <Tag color={config.color}>{config.text}</Tag>;
      },
    },
    {
      title: t('resourceType'),
      dataIndex: 'resourceType',
      width: 150,
      ellipsis: true,
    },
    {
      title: t('resourceId'),
      dataIndex: 'resourceId',
      width: 200,
      hideInSearch: true,
      ellipsis: true,
    },
    {
      title: t('requestId'),
      dataIndex: 'requestId',
      width: 200,
      hideInSearch: true,
      ellipsis: true,
    },
    {
      title: t('username'),
      dataIndex: 'username',
      width: 120,
    },
    {
      title: t('geoLocation'),
      dataIndex: 'geoLocation',
      width: 150,
      hideInSearch: true,
      render: (_, record) => {
        const geo = record.geoLocation as any;
        if (!geo) return '-';
        return `${geo.province || ''} ${geo.city || ''}`.trim() || '-';
      },
    },
    {
      title: t('ipAddress'),
      dataIndex: 'ipAddress',
      width: 140,
    },
  ];

  // 按当前搜索条件导出 CSV（客户端分页聚合，上限 1 万行；全量归档走后端 JSONL 任务）
  const handleExport = async (format: AuditExportFormat) => {
    const exportColumns = columns
      .filter((c) => c.dataIndex && !c.hideInTable)
      .map((c) => ({
        key: String(c.dataIndex),
        title: typeof c.title === 'string' ? c.title : String(c.dataIndex),
      }));
    try {
      await exportAuditLogs({
        fetcher: (q) => fetchListOperationAuditLogs(q),
        filename: `operation-audit-logs-${Date.now()}.csv`,
        columns: exportColumns,
        params: latestParamsRef.current,
      }, format);
    } catch (error: any) {
      message.error(error?.message || 'Export failed');
    }
  };

  return (
    <ContentContainer heightMode="fixed" padding="16px" bottomMargin={0}>
      <div ref={containerRef} className="page-container-content">
        <ListTable<OperationAuditLog>
          actionRef={actionRef}
          columns={columns}
          request={async (params, sorter) => {
            const query = new PaginationQuery({
              paging: {
                page: params.current || 1,
                pageSize: params.pageSize || 20,
              },
              formValues: Object.fromEntries(
                Object.entries(params).filter(([key]) => !['current', 'pageSize'].includes(key)),
              ),
              orderBy:
                sorter && Object.keys(sorter).length > 0
                  ? Object.entries(sorter).map(([key, value]) =>
                      value === 'ascend' ? key : `-${key}`,
                    )
                  : undefined,
            });

            const response = await fetchListOperationAuditLogs(query);

            return {
              data: response.items || [],
              total: response.total || 0,
              success: true,
            };
          
}}
          rowKey="id"
          search={{
            labelWidth: 'auto',
            defaultCollapsed: false,
          }}
          pagination={{
            defaultPageSize: TABLE.DEFAULT_PAGE_SIZE,
            showSizeChanger: true,
            showQuickJumper: true,
          }}
          toolBarRender={() => [
				<Button
					key="server-export"
					type="default"
					onClick={() => {
						// 服务端全量导出（XLSX，上限 50 万行）；不携带搜索条件=全量归档语义
						exportAuditLogsServer('operation', new PaginationQuery({}))
							.then(() => undefined)
							.catch((error: Error) => console.error('server export failed:', error));
					}}
				>
					服务端全量
				</Button>,
            <Dropdown
              key="export"
              menu={{
                items: [
                  { key: 'csv', label: t('exportCsv') },
                  { key: 'xlsx', label: t('exportXlsx') },
                ],
                onClick: ({ key }) => handleExport(key as AuditExportFormat),
              }}
            >
              <Button icon={<DownloadOutlined />}>{t('export')}</Button>
            </Dropdown>,
          ]}

          options={{
            density: true,
            fullScreen: true,
            setting: true,
            reload: true,
          }}
          size="middle"
          bordered
          cardBordered={false}
          scroll={{
            y: tableScrollY,
            x: 1500,
          }}
        />
      </div>
    </ContentContainer>
  );
};

export default OperationAuditLogPage;
