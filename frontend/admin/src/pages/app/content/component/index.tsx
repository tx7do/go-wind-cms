import { useRef, useState } from 'react';
import type { ActionType } from '@ant-design/pro-components';
import {
  App,
  Button,
  Collapse,
  Drawer,
  Input,
  InputNumber,
  Popconfirm,
  Select,
  Space,
  Switch,
} from 'antd';
import { DeleteOutlined, EditOutlined, PlusOutlined } from '@ant-design/icons';
import { useTranslation } from 'react-i18next';
import ListTable from '@/components/common/ListTable';
import ContentContainer from '@/layouts/components/PageContainer/ContentContainer';
import type {
  contentservicev1_ComponentDefinition as ComponentDefinition,
  contentservicev1_FieldDefinition,
} from '@/api/generated/admin/service/v1';
import {
  useCreateComponentDefinition,
  useDeleteComponentDefinition,
  fetchComponentFieldDefinitions,
  useListComponentDefinitions,
  useSaveComponentFieldDefinitions,
  useUpdateComponentDefinition,
} from '@/api/hooks/component-definition';

/** 组件字段面板（与模型字段同 schema；不支持嵌套组件——后端校验拒绝） */
interface FieldPanel {
  clientKey: string;
  name: string;
  type: string;
  label: string;
  isRequired: boolean;
  sortOrder: number;
}

let fieldCounter = 0;

function makeEmptyField(): FieldPanel {
  fieldCounter += 1;
  return {
    clientKey: `cfield-${fieldCounter}`,
    name: '',
    type: 'FIELD_TYPE_TEXT',
    label: '',
    isRequired: false,
    sortOrder: 0,
  };
}

/**
 * 组件定义管理页（内容建模 Components）：组件 CRUD + 字段构建器
 * （字段按 componentId 整体替换保存；组件被模型字段引用时后端拒绝删除）。
 */
const ComponentDefinitionManagement = () => {
  const { t } = useTranslation('component-definition');
  const { message } = App.useApp();
  const actionRef = useRef<ActionType>(null);

  const [formOpen, setFormOpen] = useState(false);
  const [saving, setSaving] = useState(false);
  const [form, setForm] = useState<{
    id: number;
    code: string;
    name: string;
    description: string;
    fields: FieldPanel[];
  }>({ id: 0, code: '', name: '', description: '', fields: [] });
  const [activeKeys, setActiveKeys] = useState<string[]>([]);

  const listQuery = useListComponentDefinitions();
  const createMutation = useCreateComponentDefinition();
  const updateMutation = useUpdateComponentDefinition();
  const deleteMutation = useDeleteComponentDefinition();
  const saveFieldsMutation = useSaveComponentFieldDefinitions();

  const dataSource = (listQuery.data?.items ?? []) as ComponentDefinition[];

  const setField = (clientKey: string, patch: Partial<FieldPanel>) => {
    setForm((prev) => ({
      ...prev,
      fields: prev.fields.map((f) =>
        f.clientKey === clientKey ? { ...f, ...patch } : f,
      ),
    }));
  };

  const addField = () => {
    const f = makeEmptyField();
    setForm((prev) => ({ ...prev, fields: [...prev.fields, f] }));
    setActiveKeys((prev) => [...prev, f.clientKey]);
  };

  const removeField = (clientKey: string) => {
    setForm((prev) => ({
      ...prev,
      fields: prev.fields.filter((f) => f.clientKey !== clientKey),
    }));
    setActiveKeys((prev) => prev.filter((k) => k !== clientKey));
  };

  const openCreate = () => {
    setForm({ id: 0, code: '', name: '', description: '', fields: [] });
    setFormOpen(true);
  };

  const openEdit = (record: ComponentDefinition) => {
    const componentId = record.id ?? 0;
    setForm({
      id: componentId,
      code: record.code ?? '',
      name: record.name ?? '',
      description: record.description ?? '',
      fields: [],
    });
    setFormOpen(true);
    // 异步拉字段（setState 在异步回调内，不在 effect 同步链路）
    fetchComponentFieldDefinitions(componentId)
      .then((resp) => {
        setForm((prev) => {
          if (prev.id !== componentId) return prev;
          const fields: FieldPanel[] = (resp.items ?? []).map((f) => {
            fieldCounter += 1;
            return {
              clientKey: `cfield-${fieldCounter}`,
              name: f.name ?? '',
              type: f.type ?? 'FIELD_TYPE_TEXT',
              label: f.label ?? '',
              isRequired: f.isRequired ?? false,
              sortOrder: f.sortOrder ?? 0,
            } as FieldPanel;
          });
          return { ...prev, fields };
        });
      })
      .catch((error: unknown) => {
        console.error('Load component fields failed:', error);
      });
  };

  const handleSubmit = async () => {
    if (!form.code.trim() || !form.name.trim()) {
      message.warning(t('nameRequired'));
      return;
    }
    setSaving(true);
    try {
      const fieldsPayload = form.fields.map((f) => ({
        name: f.name,
        type: f.type,
        label: f.label,
        isRequired: f.isRequired,
        sortOrder: f.sortOrder,
      })) as Array<Partial<contentservicev1_FieldDefinition>>;
      if (form.id) {
        await updateMutation.mutateAsync({
          id: form.id,
          values: { name: form.name.trim(), description: form.description },
        });
        await saveFieldsMutation.mutateAsync({
          componentId: form.id,
          fields: fieldsPayload,
        });
      } else {
        const created = await createMutation.mutateAsync({
          code: form.code.trim(),
          name: form.name.trim(),
          description: form.description,
        });
        if (created.id && fieldsPayload.length > 0) {
          await saveFieldsMutation.mutateAsync({
            componentId: created.id,
            fields: fieldsPayload,
          });
        }
      }
      message.success(t('saveSuccess'));
      setFormOpen(false);
      listQuery.refetch();
    } catch (error) {
      // 错误消息由 hooks 的 onError 统一提示
      console.error('Save component definition failed:', error);
    } finally {
      setSaving(false);
    }
  };

  const handleDelete = async (record: ComponentDefinition) => {
    if (!record.id) return;
    try {
      await deleteMutation.mutateAsync({ id: record.id });
      message.success(t('deleteSuccess'));
      listQuery.refetch();
    } catch (error) {
      console.error('Delete component definition failed:', error);
    }
  };

  const columns = [
    { title: 'ID', dataIndex: 'id', width: 60 },
    { title: t('code'), dataIndex: 'code', width: 200 },
    { title: t('name'), dataIndex: 'name', ellipsis: true },
    { title: t('description'), dataIndex: 'description', ellipsis: true },
    {
      title: t('fieldCount'),
      dataIndex: 'fieldCount',
      width: 90,
      render: (v: unknown) => Number(v || 0),
    },
    {
      title: t('action'),
      width: 150,
      render: (_: unknown, record: ComponentDefinition) => (
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
        loading={listQuery.isPending}
        rowKey="id"
        columns={columns}
        request={async () => {
          const items = dataSource;
          return {
            data: items,
            total: Number(listQuery.data?.total) || items.length,
            success: true,
          };
        }}
        dataSource={dataSource}
        toolBarRender={() => [
          <Button key="create" type="primary" icon={<PlusOutlined />} onClick={openCreate}>
            {t('create')}
          </Button>,
        ]}
        pagination={{ pageSize: 10 }}
      />

      <Drawer
        title={form.id ? t('editTitle') : t('createTitle')}
        width={620}
        open={formOpen}
        onClose={() => setFormOpen(false)}
        destroyOnHidden
        footer={
          <div style={{ display: 'flex', justifyContent: 'flex-end' }}>
            <Button type="primary" loading={saving} onClick={handleSubmit}>
              {t('save')}
            </Button>
          </div>
        }
      >
        <div style={{ display: 'flex', flexDirection: 'column', gap: 12, padding: 4 }}>
          <div style={{ display: 'flex', flexDirection: 'column', gap: 8 }}>
            <label style={{ fontSize: 14 }}>{t('code')}</label>
            <Input
              value={form.code}
              disabled={!!form.id}
              placeholder={t('codePlaceholder')}
              onChange={(e) => setForm((f) => ({ ...f, code: e.target.value }))}
            />
          </div>
          <div style={{ display: 'flex', flexDirection: 'column', gap: 8 }}>
            <label style={{ fontSize: 14 }}>{t('name')}</label>
            <Input
              value={form.name}
              onChange={(e) => setForm((f) => ({ ...f, name: e.target.value }))}
            />
          </div>
          <div style={{ display: 'flex', flexDirection: 'column', gap: 8 }}>
            <label style={{ fontSize: 14 }}>{t('description')}</label>
            <Input.TextArea
              rows={2}
              value={form.description}
              onChange={(e) => setForm((f) => ({ ...f, description: e.target.value }))}
            />
          </div>

          <hr style={{ border: 'none', borderTop: '1px solid var(--ant-color-border)' }} />

          <div style={{ display: 'flex', alignItems: 'center', justifyContent: 'space-between' }}>
            <span style={{ fontSize: 15, fontWeight: 500 }}>{t('fieldsTitle')}</span>
            <Button size="small" icon={<PlusOutlined />} onClick={addField}>
              {t('addField')}
            </Button>
          </div>

          {form.fields.length > 0 ? (
            <Collapse
              activeKey={activeKeys}
              onChange={(keys) => setActiveKeys(Array.isArray(keys) ? keys : [keys])}
              items={form.fields.map((field) => ({
                key: field.clientKey,
                label: field.name || t('unnamedField'),
                children: (
                  <div style={{ display: 'flex', flexDirection: 'column', gap: 12, padding: 8 }}>
                    <div style={{ display: 'flex', justifyContent: 'flex-end' }}>
                      <Button danger size="small" onClick={() => removeField(field.clientKey)}>
                        {t('removeField')}
                      </Button>
                    </div>
                    <div style={{ display: 'flex', flexDirection: 'column', gap: 8 }}>
                      <label style={{ fontSize: 14 }}>{t('fieldName')}</label>
                      <Input
                        value={field.name}
                        onChange={(e) => setField(field.clientKey, { name: e.target.value })}
                      />
                    </div>
                    <div style={{ display: 'flex', flexDirection: 'column', gap: 8 }}>
                      <label style={{ fontSize: 14 }}>{t('fieldType')}</label>
                      <Select
                        value={field.type}
                        options={[
                          { value: 'FIELD_TYPE_TEXT', label: t('fieldTypeText') },
                          { value: 'FIELD_TYPE_NUMBER', label: t('fieldTypeNumber') },
                          { value: 'FIELD_TYPE_IMAGE', label: t('fieldTypeImage') },
                          { value: 'FIELD_TYPE_FILE', label: t('fieldTypeFile') },
                        ]}
                        onChange={(v) => setField(field.clientKey, { type: v })}
                      />
                    </div>
                    <div style={{ display: 'flex', flexDirection: 'column', gap: 8 }}>
                      <label style={{ fontSize: 14 }}>{t('fieldLabel')}</label>
                      <Input
                        value={field.label}
                        onChange={(e) => setField(field.clientKey, { label: e.target.value })}
                      />
                    </div>
                    <div style={{ display: 'flex', flexDirection: 'column', gap: 8 }}>
                      <label style={{ fontSize: 14 }}>{t('sortOrder')}</label>
                      <InputNumber
                        min={0}
                        style={{ width: '100%' }}
                        value={field.sortOrder}
                        onChange={(v) => setField(field.clientKey, { sortOrder: v ?? 0 })}
                      />
                    </div>
                    <div style={{ display: 'flex', alignItems: 'center', gap: 8 }}>
                      <Switch
                        checked={field.isRequired}
                        onChange={(v) => setField(field.clientKey, { isRequired: v })}
                      />
                      <span style={{ fontSize: 14 }}>{t('fieldIsRequired')}</span>
                    </div>
                  </div>
                ),
              }))}
            />
          ) : (
            <div
              style={{
                textAlign: 'center',
                color: 'var(--ant-color-text-quaternary)',
                padding: '32px 0',
              }}
            >
              {t('emptyFields')}
            </div>
          )}
        </div>
      </Drawer>
    </ContentContainer>
  );
};

export default ComponentDefinitionManagement;
