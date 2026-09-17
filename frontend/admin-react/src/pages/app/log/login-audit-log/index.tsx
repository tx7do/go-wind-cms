import { useRef, useState } from 'react';
import type { ProColumns, ActionType } from '@ant-design/pro-components';
import { ProTable } from '@ant-design/pro-components';
import { Button, Dropdown, Tag, App } from 'antd';
import { DownloadOutlined, InfoCircleOutlined } from '@ant-design/icons';
import { useTranslation } from 'react-i18next';
import type { auditservicev1_LoginAuditLog as LoginAuditLog } from '@/api/generated/admin/service/v1';
import { PaginationQuery } from '@/core';
import { TABLE } from '@/config/constants';
import { fetchListLoginAuditLogs } from '@/api/hooks/login-audit-log';
import { exportAuditLogs, type AuditExportFormat } from '@/utils/csv';
import { useProTableScrollY } from '@/hooks/useProTableScrollY';
import ContentContainer from '@/layouts/components/PageContainer/ContentContainer';
import {
  getStatusMap,
  getStatusOptions,
  getActionTypeMap,
  getActionTypeOptions,
  getRiskLevelMap,
  getRiskLevelOptions,
} from './constants';
import LoginAuditLogDetailDrawer from './DetailDrawer';

/**
 * 登录审计日志页面
 */
const LoginAuditLogPage = () => {
  const { t } = useTranslation('login-audit-log');
  const actionRef = useRef<ActionType>(null);
  // 最近一次列表查询参数：导出 CSV 时沿用当前搜索/排序条件
  const latestParamsRef = useRef<Record<string, any>>({});
  const { message } = App.useApp();

  const containerRef = useRef<HTMLDivElement>(null);
  const tableScrollY = useProTableScrollY(containerRef);

  const statusMap = getStatusMap(t);
  const actionTypeMap = getActionTypeMap(t);
  const riskLevelMap = getRiskLevelMap(t);

  const [detailRecord, setDetailRecord] = useState<LoginAuditLog | null>(null);

  const columns: ProColumns<LoginAuditLog>[] = [
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
      title: t('status'),
      dataIndex: 'status',
      width: 100,
      valueType: 'select',
      fieldProps: {
        options: getStatusOptions(t),
      },
      render: (_, record) => {
        const config = statusMap[record.status as keyof typeof statusMap];
        if (!config) return record.status || '-';
        return <Tag color={config.color}>{config.text}</Tag>;
      },
    },
    {
      title: t('username'),
      dataIndex: 'username',
      width: 120,
    },
    {
      title: t('actionType'),
      dataIndex: 'actionType',
      width: 120,
      valueType: 'select',
      fieldProps: {
        options: getActionTypeOptions(t),
      },
      render: (_, record) => {
        const config = actionTypeMap[record.actionType as keyof typeof actionTypeMap];
        if (!config) return record.actionType || '-';
        return <Tag color={config.color}>{config.text}</Tag>;
      },
    },
    {
      title: t('riskLevel'),
      dataIndex: 'riskLevel',
      width: 100,
      valueType: 'select',
      fieldProps: {
        options: getRiskLevelOptions(t),
      },
      render: (_, record) => {
        const config = riskLevelMap[record.riskLevel as keyof typeof riskLevelMap];
        if (!config) return record.riskLevel || '-';
        return <Tag color={config.color}>{config.text}</Tag>;
      },
    },
    {
      title: t('platform'),
      dataIndex: 'deviceInfo',
      width: 180,
      hideInSearch: true,
      render: (_, record) => {
        const info = record.deviceInfo as any;
        if (!info) return '-';
        return `${info.osName || ''} ${info.browserName || ''}`.trim() || '-';
      },
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
    {
      title: t('action'),
      valueType: 'option',
      width: 80,
      render: (_, record) => [
        <a
          key="detail"
          title={t('detail')}
          onClick={() => setDetailRecord(record)}
        >
          <InfoCircleOutlined />
        </a>,
      ],
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
        fetcher: (q) => fetchListLoginAuditLogs(q),
        filename: `login-audit-logs-${Date.now()}.csv`,
        columns: exportColumns,
        params: latestParamsRef.current,
      }, format);
    } catch (error: any) {
      message.error(error?.message || 'Export failed');
    }
  };

  return (
    <>
      <ContentContainer heightMode="fixed" padding="16px" bottomMargin={0}>
        <div ref={containerRef} className="page-container-content">
          <ProTable<LoginAuditLog>
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
                    Object.entries(params).filter(
                      ([key]) => !['current', 'pageSize'].includes(key),
                    ),
                  ),
                  orderBy:
                    sorter && Object.keys(sorter).length > 0
                      ? Object.entries(sorter).map(([key, value]) =>
                          value === 'ascend' ? key : `-${key}`,
                        )
                      : undefined,
                });

                const response = await fetchListLoginAuditLogs(query);

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

      <LoginAuditLogDetailDrawer
        open={detailRecord !== null}
        data={detailRecord}
        onClose={() => setDetailRecord(null)}
      />
    </>
  );
};

export default LoginAuditLogPage;
