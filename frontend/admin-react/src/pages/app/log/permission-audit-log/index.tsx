import { useRef } from 'react';
import type { ProColumns, ActionType } from '@ant-design/pro-components';
import { ProTable } from '@ant-design/pro-components';
import { DownloadOutlined } from '@ant-design/icons';
import { Button, Dropdown, Tag, App } from 'antd';
import { useTranslation } from 'react-i18next';
import type { auditservicev1_PermissionAuditLog as PermissionAuditLog } from '@/api/generated/admin/service/v1';
import { PaginationQuery } from '@/core';
import { TABLE } from '@/config/constants';
import { fetchListPermissionAuditLogs } from '@/api/hooks/permission-audit-log';
import { exportAuditLogs, type AuditExportFormat } from '@/utils/csv';
import { useProTableScrollY } from '@/hooks/useProTableScrollY';
import ContentContainer from '@/layouts/components/PageContainer/ContentContainer';
import { getActionMap, getActionOptions } from './constants';

/**
 * 权限审计日志页面
 */
const PermissionAuditLogPage = () => {
  const { t } = useTranslation('permission-audit-log');
  const actionRef = useRef<ActionType>(null);
  // 最近一次列表查询参数：导出 CSV 时沿用当前搜索/排序条件
  const latestParamsRef = useRef<Record<string, any>>({});
  const { message } = App.useApp();

  const containerRef = useRef<HTMLDivElement>(null);
  const tableScrollY = useProTableScrollY(containerRef);

  const actionMap = getActionMap(t);
  const actionOptions = getActionOptions(t);

  const columns: ProColumns<PermissionAuditLog>[] = [
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
      title: t('targetType'),
      dataIndex: 'targetType',
      width: 120,
    },
    {
      title: t('targetName'),
      dataIndex: 'targetName',
      width: 150,
      hideInSearch: true,
      ellipsis: true,
    },
    {
      title: t('reason'),
      dataIndex: 'reason',
      width: 200,
      hideInSearch: true,
      ellipsis: true,
    },
    {
      title: t('operatorName'),
      dataIndex: 'operatorName',
      width: 120,
      hideInSearch: false,
      render: (_, record) => {
        const info = record.operatorName;
        if (!info) return '-';
        return info || '-';
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
        fetcher: (q) => fetchListPermissionAuditLogs(q),
        filename: `permission-audit-logs-${Date.now()}.csv`,
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
        <ProTable<PermissionAuditLog>
          actionRef={actionRef}
          columns={columns}
          request={async (params, sorter) => {
            latestParamsRef.current = params;
            try {
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

              const response = await fetchListPermissionAuditLogs(query);

              return {
                data: response.items || [],
                total: response.total || 0,
                success: true,
              };
            } catch (error: any) {
              message.error(error.message || t('fetchFailed'));
              return {
                data: [],
                total: 0,
                success: false,
              };
            }
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
            x: 1300,
          }}
        />
      </div>
    </ContentContainer>
  );
};

export default PermissionAuditLogPage;
