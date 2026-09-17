import { useRef } from 'react';
import type { ProColumns, ActionType } from '@ant-design/pro-components';
import { ProTable } from '@ant-design/pro-components';
import { DownloadOutlined } from '@ant-design/icons';
import { Button, Dropdown, Tag, App } from 'antd';
import { useTranslation } from 'react-i18next';
import type { permissionservicev1_PolicyEvaluationLog as PolicyEvaluationLog } from '@/api/generated/admin/service/v1';
import { PaginationQuery } from '@/core';
import { TABLE, METHOD_LIST } from '@/config/constants';
import { fetchListPolicyEvaluationLogs } from '@/api/hooks/policy-evaluation-log';
import { exportAuditLogs, type AuditExportFormat } from '@/utils/csv';
import { useProTableScrollY } from '@/hooks/useProTableScrollY';
import ContentContainer from '@/layouts/components/PageContainer/ContentContainer';

/**
 * 评估结果筛选项
 */
const getResultStatusList = (t: (key: string) => string) => [
  { label: t('result.true'), value: 'true' },
  { label: t('result.false'), value: 'false' },
];

/**
 * 根据评估结果获取标签颜色
 */
function resultToColor(result: boolean | undefined): string {
  if (result === true) return 'success';
  if (result === false) return 'error';
  return 'default';
}

/**
 * 根据评估结果获取显示文本
 */
function resultToName(t: (key: string) => string, result: boolean | undefined): string {
  if (result === true) return t('result.true');
  if (result === false) return t('result.false');
  return '-';
}

/**
 * 策略评估日志页面
 */
const PolicyEvaluationLog = () => {
  const { t } = useTranslation('policy-evaluation-log');
  const actionRef = useRef<ActionType>(null);
  // 最近一次列表查询参数：导出 CSV 时沿用当前搜索/排序条件
  const latestParamsRef = useRef<Record<string, any>>({});
  const { message } = App.useApp();

  const containerRef = useRef<HTMLDivElement>(null);
  const tableScrollY = useProTableScrollY(containerRef);

  const columns: ProColumns<PolicyEvaluationLog>[] = [
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
      title: t('result'),
      dataIndex: 'result',
      width: 100,
      valueType: 'select',
      valueEnum: Object.fromEntries(
        getResultStatusList(t).map((item) => [
          item.value,
          { text: item.label, status: item.value === 'true' ? 'Success' : 'Error' },
        ]),
      ),
      render: (_, record) => (
        <Tag color={resultToColor(record.result)}>
          {resultToName(t, record.result)}
        </Tag>
      ),
    },
    {
      title: t('requestMethod'),
      dataIndex: 'requestMethod',
      width: 110,
      valueType: 'select',
      valueEnum: Object.fromEntries(
        METHOD_LIST.map((item) => [item.value, { text: item.label, status: 'Default' }]),
      ),
      render: (_, record) => {
        const colorMap: Record<string, string> = {
          GET: 'success',
          POST: 'processing',
          PUT: 'warning',
          DELETE: 'error',
          PATCH: 'default',
          HEAD: 'default',
          OPTIONS: 'default',
        };

  return (
          <span style={{ color: colorMap[record.requestMethod || ''] || '#666' }}>
            {record.requestMethod}
          </span>
        );
      },
    },
    {
      title: t('requestPath'),
      dataIndex: 'requestPath',
      width: 240,
      ellipsis: true,
    },
    {
      title: t('policyId'),
      dataIndex: 'policyId',
      width: 90,
    },
    {
      title: t('permissionId'),
      dataIndex: 'permissionId',
      width: 100,
      hideInSearch: true,
    },
    {
      title: t('userId'),
      dataIndex: 'userId',
      width: 90,
    },
    {
      title: t('membershipId'),
      dataIndex: 'membershipId',
      width: 110,
      hideInSearch: true,
    },
    {
      title: t('tenantId'),
      dataIndex: 'tenantId',
      width: 90,
      hideInSearch: true,
    },
    {
      title: t('ipAddress'),
      dataIndex: 'ipAddress',
      width: 140,
    },
    {
      title: t('traceId'),
      dataIndex: 'traceId',
      width: 180,
      hideInSearch: true,
      ellipsis: true,
    },
    {
      title: t('evaluationContext'),
      dataIndex: 'evaluationContext',
      width: 200,
      hideInSearch: true,
      ellipsis: true,
    },
    {
      title: t('effectDetails'),
      dataIndex: 'effectDetails',
      width: 180,
      hideInSearch: true,
      ellipsis: true,
    },
    {
      title: t('scopeSql'),
      dataIndex: 'scopeSql',
      width: 160,
      hideInSearch: true,
      ellipsis: true,
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
        fetcher: (q) => fetchListPolicyEvaluationLogs(q),
        filename: `policy-evaluation-logs-${Date.now()}.csv`,
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
        <ProTable<PolicyEvaluationLog>
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

              const response = await fetchListPolicyEvaluationLogs(query);

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
            x: 1600,
          }}
        />
      </div>
    </ContentContainer>
  );
};

export default PolicyEvaluationLog;
