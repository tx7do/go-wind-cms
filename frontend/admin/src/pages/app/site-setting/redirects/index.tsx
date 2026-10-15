import { useRef, useState } from 'react';
import type { ActionType } from '@ant-design/pro-components';
import { App, Button, Input, Modal, Popconfirm, Select, Space, Switch } from 'antd';
import {
  DeleteOutlined,
  DownloadOutlined,
  EditOutlined,
  FileTextOutlined,
  PlusOutlined,
} from '@ant-design/icons';
import { useTranslation } from 'react-i18next';
import { PaginationQuery } from '@/core';
import ListTable from '@/components/common/ListTable';
import ContentContainer from '@/layouts/components/PageContainer/ContentContainer';
import type { contentservicev1_Redirect as Redirect } from '@/api/generated/admin/service/v1';
import {
  redirectsToNginxMap,
  useCreateRedirect,
  useDeleteRedirect,
  useRedirectsList,
  useUpdateRedirect,
} from '@/api/hooks/redirect';

const STATUS_CODES = [301, 302, 307, 308];

/**
 * URL 重定向管理页：SEO 迁移/短链。
 * 消费面：vue SSR 路由中间件查 /app/v1/redirects 做 3xx；nginx/CDN 用导出的
 * map 配置做边缘层直跳（导出文件按注释 include 即可）。
 */
const RedirectManagement = () => {
  const { t } = useTranslation('redirect');
  const { message } = App.useApp();
  const actionRef = useRef<ActionType>(null);

  const [items, setItems] = useState<Redirect[]>([]);
  const [formOpen, setFormOpen] = useState(false);
  const [editing, setEditing] = useState<Redirect | undefined>();
  const [form, setForm] = useState<{
    fromPath: string;
    toUrl: string;
    statusCode: number;
    enabled: boolean;
    description: string;
  }>({ fromPath: '', toUrl: '', statusCode: 301, enabled: true, description: '' });

  const listMutation = useRedirectsList();
  const createMutation = useCreateRedirect();
  const updateMutation = useUpdateRedirect();
  const deleteMutation = useDeleteRedirect();

  // 列表拉取（供保存/删除后刷新与导出取数）
  const loadItems = async () => {
    const q = new PaginationQuery({
      paging: { page: 1, pageSize: 500 },
      orderBy: ['-id'],
    });
    const resp = await listMutation.mutateAsync(q);
    const items = (resp?.items ?? []) as Redirect[];
    setItems(items);
    return items;
  };

  const openCreate = () => {
    setEditing(undefined);
    setForm({ fromPath: '', toUrl: '', statusCode: 301, enabled: true, description: '' });
    setFormOpen(true);
  };

  const openEdit = (record: Redirect) => {
    setEditing(record);
    setForm({
      fromPath: record.fromPath ?? '',
      toUrl: record.toUrl ?? '',
      statusCode: record.statusCode ?? 301,
      enabled: record.enabled ?? true,
      description: record.description ?? '',
    });
    setFormOpen(true);
  };

  const handleSave = async () => {
    if (!form.fromPath.trim() || !form.toUrl.trim()) {
      message.warning(t('requiredFields'));
      return;
    }
    try {
      if (editing?.id) {
        await updateMutation.mutateAsync({
          id: editing.id,
          data: { ...form, fromPath: form.fromPath.trim(), toUrl: form.toUrl.trim() },
          updateMask: 'fromPath,toUrl,statusCode,enabled,description',
        });
      } else {
        await createMutation.mutateAsync({
          data: { ...form, fromPath: form.fromPath.trim(), toUrl: form.toUrl.trim() },
        });
      }
      message.success(t('saveSuccess'));
      setFormOpen(false);
      await loadItems();
    } catch (error) {
      console.error('Save redirect failed:', error);
    }
  };

  const handleDelete = async (record: Redirect) => {
    if (!record.id) return;
    try {
      await deleteMutation.mutateAsync({ id: record.id });
      message.success(t('deleteSuccess'));
      await loadItems();
    } catch (error) {
      console.error('Delete redirect failed:', error);
    }
  };

  const download = (content: string, filename: string, mime: string) => {
    const blob = new Blob([content], { type: mime });
    const url = URL.createObjectURL(blob);
    const a = document.createElement('a');
    a.href = url;
    a.download = filename;
    a.click();
    URL.revokeObjectURL(url);
  };

  const columns = [
    { title: 'ID', dataIndex: 'id', width: 60 },
    { title: t('fromPath'), dataIndex: 'fromPath', ellipsis: true, copyable: true },
    { title: t('toUrl'), dataIndex: 'toUrl', ellipsis: true },
    {
      title: t('statusCode'),
      dataIndex: 'statusCode',
      width: 90,
      render: (v: unknown) => <span style={{ fontFamily: 'monospace' }}>{String(v ?? 301)}</span>,
    },
    {
      title: t('enabled'),
      dataIndex: 'enabled',
      width: 80,
      render: (_: unknown, record: Redirect) => (
        <Switch
          checkedChildren="301"
          unCheckedChildren="off"
          checked={record.enabled ?? false}
          onChange={async (checked) => {
            if (!record.id) return;
            try {
              await updateMutation.mutateAsync({
                id: record.id,
                data: { enabled: checked },
                updateMask: 'enabled',
              });
              await loadItems();
            } catch (error) {
              console.error('Toggle redirect failed:', error);
            }
          }}
        />
      ),
    },
    { title: t('description'), dataIndex: 'description', ellipsis: true },
    {
      title: t('action'),
      width: 140,
      render: (_: unknown, record: Redirect) => (
        <Space>
          <Button type="link" size="small" icon={<EditOutlined />} onClick={() => openEdit(record)}>
            {t('edit')}
          </Button>
          <Popconfirm title={t('deleteConfirm')} onConfirm={() => handleDelete(record)}>
            <Button type="link" size="small" danger icon={<DeleteOutlined />}>
              {t('delete')}
            </Button>
          </Popconfirm>
        </Space>
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
        dataSource={items}
        request={async () => {
          try {
            const data = await loadItems();
            return { data, total: data.length, success: true };
          } catch (error) {
            console.error('Load redirects failed:', error);
            return { data: [], total: 0, success: false };
          }
        }}
        toolBarRender={() => [
          <Button
            key="nginx"
            size="small"
            icon={<DownloadOutlined />}
            onClick={() => {
              download(
                redirectsToNginxMap(items),
                'redirects-nginx.conf',
                'text/plain;charset=utf-8',
              );
            }}
          >
            {t('exportNginx')}
          </Button>,
          <Button
            key="json"
            size="small"
            icon={<FileTextOutlined />}
            onClick={() => {
              download(
                JSON.stringify({ version: 1, redirects: items }, null, 2),
                'redirects.json',
                'application/json',
              );
            }}
          >
            {t('exportJson')}
          </Button>,
          <Button key="create" type="primary" icon={<PlusOutlined />} onClick={openCreate}>
            {t('create')}
          </Button>,
        ]}
        pagination={{ pageSize: 20 }}
      />

      <Modal
        open={formOpen}
        title={editing ? t('edit') : t('create')}
        onCancel={() => setFormOpen(false)}
        onOk={handleSave}
        confirmLoading={createMutation.isPending || updateMutation.isPending}
        width={560}
        destroyOnHidden
      >
        <Space direction="vertical" style={{ width: '100%' }} size={12}>
          <Input
            placeholder={t('fromPlaceholder')}
            value={form.fromPath}
            onChange={(e) => setForm((f) => ({ ...f, fromPath: e.target.value }))}
          />
          <Input
            placeholder={t('toPlaceholder')}
            value={form.toUrl}
            onChange={(e) => setForm((f) => ({ ...f, toUrl: e.target.value }))}
          />
          <Select
            value={form.statusCode}
            options={STATUS_CODES.map((c) => ({ value: c, label: `${c}` }))}
            onChange={(v: number) => setForm((f) => ({ ...f, statusCode: v }))}
            style={{ width: 120 }}
          />
          <Input
            placeholder={t('descriptionPlaceholder')}
            value={form.description}
            onChange={(e) => setForm((f) => ({ ...f, description: e.target.value }))}
          />
          <Space>
            <span>{t('enabled')}</span>
            <Switch checked={form.enabled} onChange={(v) => setForm((f) => ({ ...f, enabled: v }))} />
          </Space>
        </Space>
      </Modal>
    </ContentContainer>
  );
};

export default RedirectManagement;
