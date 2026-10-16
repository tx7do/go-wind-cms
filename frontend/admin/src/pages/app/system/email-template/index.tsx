import { useRef, useState } from 'react';
import { App, Button, Input, Popconfirm, Switch } from 'antd';
import { DeleteOutlined, EditOutlined, PlusOutlined } from '@ant-design/icons';
import { useTranslation } from 'react-i18next';
import ListTable from '@/components/common/ListTable';
import ContentContainer from '@/layouts/components/PageContainer/ContentContainer';
import { PaginationQuery } from '@/core';
import type { ActionType } from '@ant-design/pro-components';
import {
  useCreateEmailTemplate,
  useDeleteEmailTemplate,
  useEmailTemplatesList,
  useUpdateEmailTemplate,
} from '@/api/hooks/email-template';

/**
 * 邮件模板管理页：系统通知出站的可视化模板。
 * code 唯一；主题/正文支持 text/template 占位符
 * （{{.Title}}/{{.Content}}/{{.SiteURL}}）；渲染顺序：租户启用 → 平台启用 → 内置。
 */

interface TemplateRow {
  id?: number;
  updatedAt?: string;
  code?: string;
  name?: string;
  subject?: string;
  bodyHtml?: string;
  enabled?: boolean;
  description?: string;
}

const EmailTemplateManagement = () => {
  const { t } = useTranslation('email-template');
  const { message } = App.useApp();
  const actionRef = useRef<ActionType>(null);

  const [items, setItems] = useState<TemplateRow[]>([]);
  const [drawerOpen, setDrawerOpen] = useState(false);
  const [editing, setEditing] = useState<TemplateRow | undefined>();
  const [form, setForm] = useState<TemplateRow>({});

  const listMutation = useEmailTemplatesList();
  const createMutation = useCreateEmailTemplate();
  const updateMutation = useUpdateEmailTemplate();
  const deleteMutation = useDeleteEmailTemplate();

  const loadItems = async () => {
    const q = new PaginationQuery({
      paging: { page: 1, pageSize: 100 },
      orderBy: ['-id'],
    });
    const resp = await listMutation.mutateAsync(q);
    const rows = (resp?.items ?? []) as TemplateRow[];
    setItems(rows);
    return { data: rows, total: Number(resp?.total) || rows.length, success: true };
  };

  const openCreate = () => {
    setEditing(undefined);
    setForm({ enabled: true });
    setDrawerOpen(true);
  };

  const openEdit = (record: TemplateRow) => {
    setEditing(record);
    setForm({ ...record });
    setDrawerOpen(true);
  };

  const handleSave = async () => {
    if (!form.code || !form.name || !form.subject) {
      message.warning(t('requiredFields'));
      return;
    }
    try {
      if (editing?.id) {
        await updateMutation.mutateAsync({
          id: editing.id,
          data: { ...form, id: undefined } as never,
          updateMask: 'code,name,subject,bodyHtml,enabled,description',
        });
        message.success(t('updateSuccess'));
      } else {
        await createMutation.mutateAsync({ data: form as never });
        message.success(t('createSuccess'));
      }
      setDrawerOpen(false);
      actionRef.current?.reload();
    } catch (error) {
      console.error('Save email template failed:', error);
    }
  };

  const handleDelete = async (record: TemplateRow) => {
    if (!record.id) return;
    try {
      await deleteMutation.mutateAsync({ id: record.id });
      message.success(t('deleteSuccess'));
      actionRef.current?.reload();
    } catch (error) {
      console.error('Delete email template failed:', error);
    }
  };

  const columns = [
    { title: 'ID', dataIndex: 'id', width: 70, search: false },
    { title: t('code'), dataIndex: 'code', width: 160, ellipsis: true },
    { title: t('name'), dataIndex: 'name', width: 180, ellipsis: true },
    {
      title: t('enabled'),
      dataIndex: 'enabled',
      width: 90,
      search: false,
      render: (_: unknown, record: TemplateRow) => (
        <Switch checked={record.enabled ?? false} disabled size="small" />
      ),
    },
    {
      title: t('subject'),
      dataIndex: 'subject',
      search: false,
      ellipsis: true,
    },
    {
      title: t('updatedAt'),
      dataIndex: 'updatedAt',
      search: false,
      width: 160,
      render: (_: unknown, record: TemplateRow) =>
        record.updatedAt ? new Date(String(record.updatedAt)).toLocaleString() : '-',
    },
    {
      title: t('action'),
      valueType: 'option' as const,
      width: 140,
      fixed: 'right' as const,
      render: (_: unknown, record: TemplateRow) => (
        <span>
          <Button type="link" size="small" icon={<EditOutlined />} onClick={() => openEdit(record)}>
            {t('edit')}
          </Button>
          <Popconfirm title={t('deleteConfirm')} onConfirm={() => handleDelete(record)}>
            <Button type="link" size="small" danger icon={<DeleteOutlined />}>
              {t('delete')}
            </Button>
          </Popconfirm>
        </span>
      ),
    },
  ];

  return (
    <ContentContainer>
      <ListTable
        headerTitle={t('title')}
        actionRef={actionRef}
        loading={listMutation.isPending}
        rowKey="id"
        columns={columns}
        request={loadItems}
        dataSource={items}
        toolBarRender={() => [
          <Button key="create" type="primary" icon={<PlusOutlined />} onClick={openCreate}>
            {t('create')}
          </Button>,
        ]}
        pagination={{ pageSize: 20 }}
      />

      {drawerOpen && (
        <div
          style={{
            position: 'fixed',
            inset: 0,
            zIndex: 1000,
            background: 'rgba(0,0,0,0.45)',
            display: 'flex',
            justifyContent: 'flex-end',
          }}
          onClick={() => setDrawerOpen(false)}
        >
          <div
            style={{
              width: 640,
              maxWidth: '90vw',
              height: '100%',
              background: 'var(--ant-color-bg-container)',
              padding: 24,
              overflow: 'auto',
            }}
            onClick={(e) => e.stopPropagation()}
          >
            <h3 style={{ marginTop: 0 }}>{editing ? `${t('edit')} #${editing.id}` : t('create')}</h3>
            <div style={{ display: 'grid', gap: 12 }}>
              <label>
                <div>{t('code')}</div>
                <Input
                  value={form.code}
                  disabled={!!editing}
                  onChange={(e) => setForm({ ...form, code: e.target.value })}
                  placeholder="notify/system"
                />
              </label>
              <label>
                <div>{t('name')}</div>
                <Input
                  value={form.name}
                  onChange={(e) => setForm({ ...form, name: e.target.value })}
                />
              </label>
              <label>
                <div>{t('subject')}</div>
                <Input
                  value={form.subject}
                  onChange={(e) => setForm({ ...form, subject: e.target.value })}
                  placeholder="{{.Title}}"
                />
              </label>
              <label>
                <div>{t('bodyHtml')}</div>
                <Input.TextArea
                  rows={8}
                  value={form.bodyHtml}
                  onChange={(e) => setForm({ ...form, bodyHtml: e.target.value })}
                  placeholder='<div><h2>{{.Title}}</h2><div>{{.Content}}</div></div>'
                />
              </label>
              <div>
                {t('enabled')}
                <Switch
                  checked={form.enabled ?? true}
                  onChange={(v) => setForm({ ...form, enabled: v })}
                  style={{ marginLeft: 8 }}
                />
              </div>
              <label>
                <div>{t('description')}</div>
                <Input.TextArea
                  rows={2}
                  value={form.description}
                  onChange={(e: React.ChangeEvent<HTMLTextAreaElement>) => setForm({ ...form, description: e.target.value })}
                />
              </label>
            </div>
            <div style={{ marginTop: 16, display: 'flex', gap: 8, justifyContent: 'flex-end' }}>
              <Button onClick={() => setDrawerOpen(false)}>{t('common:button.cancel')}</Button>
              <Button
                type="primary"
                loading={createMutation.isPending || updateMutation.isPending}
                onClick={handleSave}
              >
                {t('save')}
              </Button>
            </div>
          </div>
        </div>
      )}
    </ContentContainer>
  );
};

export default EmailTemplateManagement;
