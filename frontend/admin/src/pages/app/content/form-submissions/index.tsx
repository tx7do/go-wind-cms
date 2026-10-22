import { useRef, useState } from 'react';
import type { ActionType } from '@ant-design/pro-components';
import { App, Button, Drawer, Space } from 'antd';
import { DeleteOutlined, DownloadOutlined, EyeOutlined } from '@ant-design/icons';
import { useTranslation } from 'react-i18next';
import ListTable from '@/components/common/ListTable';
import TableActions from '@/components/common/TableActions';
import ContentContainer from '@/layouts/components/PageContainer/ContentContainer';
import { PaginationQuery } from '@/core';
import type { contentservicev1_FormSubmission as FormSubmission } from '@/api/generated/admin/service/v1';
import {
  submissionsToCsv,
  useDeleteFormSubmission,
  useFormSubmissionsList,
} from '@/api/hooks/form-submission';

/**
 * 表单提交管理页：页面 FORM 区块的访客提交数据。
 * 列表按时间倒序；支持按页面/区块过滤（query 传 pageId/sectionId）；
 * 详情抽屉展开完整提交数据；CSV 由前端按当前数据拼装导出。
 */
const FormSubmissionManagement = () => {
  const { t } = useTranslation('form-submission');
  const { message } = App.useApp();
  const actionRef = useRef<ActionType>(null);

  const [detail, setDetail] = useState<FormSubmission | undefined>();
  const [items, setItems] = useState<FormSubmission[]>([]);
  const lastParamsRef = useRef<{ pageId?: string; sectionId?: string }>({});

  const listMutation = useFormSubmissionsList();
  const deleteMutation = useDeleteFormSubmission();

  // ListTable request：首载与重试都走这里；增删后用 actionRef.reload() 重跑
  const loadItems = async () => {
    const q = new PaginationQuery({
      paging: { page: 1, pageSize: 200 },
      formValues:
        lastParamsRef.current.pageId || lastParamsRef.current.sectionId
          ? lastParamsRef.current
          : undefined,
      orderBy: ['-id'],
    });
    const resp = await listMutation.mutateAsync(q);
    const rows = (resp?.items ?? []) as FormSubmission[];
    setItems(rows);
    return { data: rows, total: Number(resp?.total) || rows.length, success: true };
  };

  const handleDelete = async (record: FormSubmission) => {
    if (!record.id) return;
    try {
      await deleteMutation.mutateAsync({ id: record.id });
      message.success(t('deleteSuccess'));
      actionRef.current?.reload();
    } catch (error) {
      console.error('Delete form submission failed:', error);
    }
  };

  const handleExportCsv = () => {
    if (items.length === 0) {
      message.warning(t('nothingToExport'));
      return;
    }
    const blob = new Blob([submissionsToCsv(items)], { type: 'text/csv;charset=utf-8' });
    const url = URL.createObjectURL(blob);
    const a = document.createElement('a');
    a.href = url;
    a.download = `form-submissions-${new Date().toISOString().slice(0, 10)}.csv`;
    a.click();
    URL.revokeObjectURL(url);
  };

  const columns = [
    { title: 'ID', dataIndex: 'id', width: 70 },
    { title: t('pageId'), dataIndex: 'pageId', width: 90 },
    { title: t('sectionId'), dataIndex: 'sectionId', width: 90 },
    {
      title: t('dataSummary'),
      render: (_: unknown, record: FormSubmission) => {
        const data = record.data ?? {};
        const keys = Object.keys(data);
        if (keys.length === 0) return '-';
        return keys
          .slice(0, 3)
          .map((k) => `${k}: ${data[k]}`)
          .join('　') + (keys.length > 3 ? `…(+${keys.length - 3})` : '');
      },
    },
    {
      title: t('submittedAt'),
      dataIndex: 'createdAt',
      width: 170,
      render: (v: unknown) => (v ? new Date(String(v)).toLocaleString() : '-'),
    },
    {
      title: t('action'),
      width: 150,
      render: (_: unknown, record: FormSubmission) => (
        <TableActions
          record={record}
          actions={[
            {
              key: 'detail',
              label: t('detail'),
              icon: <EyeOutlined />,
              onClick: () => setDetail(record),
            },
            {
              key: 'delete',
              label: t('common:button.delete'),
              icon: <DeleteOutlined />,
              danger: true,
              confirm: { title: t('deleteConfirm') },
              onClick: () => handleDelete(record),
            },
          ]}
        />
      ),
    },
  ];

  return (
    <ContentContainer>
      <ListTable
        headerTitle={t('title')}
        actionRef={actionRef}
        search={false}
        loading={listMutation.isPending}
        rowKey="id"
        columns={columns}
        request={loadItems}
        dataSource={items}
        toolBarRender={() => [
          <Button key="export" icon={<DownloadOutlined />} onClick={handleExportCsv}>
            {t('exportCsv')}
          </Button>,
        ]}
        pagination={{ pageSize: 20 }}
      />

      <Drawer
        open={!!detail}
        onClose={() => setDetail(undefined)}
        title={`${t('detail')} #${detail?.id ?? ''}`}
        width={480}
        destroyOnHidden
      >
        <Space direction="vertical" style={{ width: '100%' }} size={8}>
          {Object.entries(detail?.data ?? {}).map(([key, value]) => (
            <div key={key} style={{ borderBottom: '1px solid var(--ant-color-split)', paddingBottom: 8 }}>
              <div style={{ fontWeight: 500 }}>{key}</div>
              <div style={{ whiteSpace: 'pre-wrap', wordBreak: 'break-all', color: 'var(--ant-color-text-secondary)' }}>
                {value || '-'}
              </div>
            </div>
          ))}
          {Object.keys(detail?.data ?? {}).length === 0 && (
            <div style={{ color: 'var(--ant-color-text-tertiary)' }}>{t('noData')}</div>
          )}
        </Space>
      </Drawer>
    </ContentContainer>
  );
};

export default FormSubmissionManagement;
