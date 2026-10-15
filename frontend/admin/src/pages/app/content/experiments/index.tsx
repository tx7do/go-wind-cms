import { useEffect, useRef, useState } from 'react';
import type { ActionType, ProColumns } from '@ant-design/pro-components';
import {
  App,
  Button,
  Drawer,
  Form,
  Input,
  InputNumber,
  Modal,
  Popconfirm,
  Select,
  Space,
  Tag,
} from 'antd';
import {
  DeleteOutlined,
  EditOutlined,
  FieldTimeOutlined,
  PlusOutlined,
  ReloadOutlined,
} from '@ant-design/icons';
import { useTranslation } from 'react-i18next';
import ListTable from '@/components/common/ListTable';
import ContentContainer from '@/layouts/components/PageContainer/ContentContainer';
import { PaginationQuery } from '@/core';
import { apiClient } from '@/api/client';
import type {
  contentservicev1_Experiment as Experiment,
  contentservicev1_ExperimentVariant as ExperimentVariant,
} from '@/api/generated/admin/service/v1';
import {
  useCreateExperiment,
  useDeleteExperiment,
  useExperimentsList,
  useExperimentVariantStats,
  useSetExperimentStatus,
  useUpdateExperiment,
} from '@/api/hooks/experiment';

const STATUS_COLORS: Record<string, string> = {
  EXPERIMENT_STATUS_DRAFT: 'default',
  EXPERIMENT_STATUS_RUNNING: 'processing',
  EXPERIMENT_STATUS_PAUSED: 'warning',
  EXPERIMENT_STATUS_ENDED: 'success',
};

/**
 * 内容实验管理页：A/B 分流 + 受众定向。
 * 列表 request 模式（loadItems + 受控 dataSource）；抽屉内整组编辑变体，
 * sections 快照为 {"sections":[...]} protojson 形态（可从目标页面一键导入）。
 */
const ExperimentManagement = () => {
  const { t } = useTranslation('experiment');
  const { message } = App.useApp();
  const actionRef = useRef<ActionType>(null);

  const [items, setItems] = useState<Experiment[]>([]);
  const [drawerOpen, setDrawerOpen] = useState(false);
  const [editing, setEditing] = useState<Experiment | undefined>();
  const [statsTarget, setStatsTarget] = useState<Experiment | undefined>();
  const [importing, setImporting] = useState(false);

  const [form] = Form.useForm();

  const listMutation = useExperimentsList();
  const createMutation = useCreateExperiment();
  const updateMutation = useUpdateExperiment();
  const deleteMutation = useDeleteExperiment();
  const setStatusMutation = useSetExperimentStatus();

  const statsQuery = useExperimentVariantStats(
    statsTarget?.id ? { experimentId: statsTarget.id } : null,
  );

  const loadItems = async () => {
    const q = new PaginationQuery({
      paging: { page: 1, pageSize: 100 },
      orderBy: ['-id'],
    });
    const resp = await listMutation.mutateAsync(q);
    const rows = (resp?.items ?? []) as Experiment[];
    setItems(rows);
    return { data: rows, total: Number(resp?.total) || rows.length, success: true };
  };

  useEffect(() => {
    if (statsTarget?.id) {
      statsQuery.refetch().catch((error: Error) =>
        console.error('Refetch variant stats failed:', error),
      );
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [statsTarget?.id]);

  const openCreate = () => {
    setEditing(undefined);
    form.resetFields();
    form.setFieldsValue({ assignmentType: 'ASSIGNMENT_TYPE_RANDOM', variants: [] });
    setDrawerOpen(true);
  };

  const openEdit = (record: Experiment) => {
    setEditing(record);
    form.resetFields();
    form.setFieldsValue({
      name: record.name,
      pageId: record.pageId,
      assignmentType: record.assignmentType,
      description: record.description,
      variants: (record.variants ?? []).map((v) => ({
        name: v.name,
        weight: v.weight,
        conditions: v.conditions,
        sections: v.sections,
        sortOrder: v.sortOrder,
      })),
    });
    setDrawerOpen(true);
  };

  const handleSubmit = async () => {
    try {
      const values = await form.validateFields();
      const data: Experiment = {
        name: values.name,
        pageId: values.pageId,
        assignmentType: values.assignmentType,
        description: values.description,
        variants: (values.variants ?? []) as ExperimentVariant[],
      };
      if (editing?.id) {
        await updateMutation.mutateAsync({
          id: editing.id,
          data,
          updateMask: 'name,pageId,assignmentType,description,variants',
        });
        message.success(t('updateSuccess'));
      } else {
        await createMutation.mutateAsync({ data });
        message.success(t('createSuccess'));
      }
      setDrawerOpen(false);
      actionRef.current?.reload();
    } catch (error) {
      if (error && typeof error === 'object' && 'errorFields' in error) return;
      console.error('Save experiment failed:', error);
    }
  };

  const handleSetStatus = async (record: Experiment, status: string, successTip: string) => {
    if (!record.id) return;
    try {
      await setStatusMutation.mutateAsync({ id: record.id, status } as never);
      message.success(successTip);
      actionRef.current?.reload();
    } catch (error) {
      console.error('Set experiment status failed:', error);
    }
  };

  const handleDelete = async (record: Experiment) => {
    if (!record.id) return;
    try {
      await deleteMutation.mutateAsync({ id: record.id });
      message.success(t('deleteSuccess'));
      actionRef.current?.reload();
    } catch (error) {
      console.error('Delete experiment failed:', error);
    }
  };

  // 从目标页面导入当前区块作为变体 sections 快照
  const handleImportSections = async () => {
    const pageId = form.getFieldValue('pageId');
    if (!pageId) {
      message.warning(t('selectPageFirst'));
      return;
    }
    setImporting(true);
    try {
      const page = await apiClient.pageService.Get({ id: pageId } as never);
      const snapshot = JSON.stringify({ sections: page?.sections ?? [] });
      const variants: Array<{ sections?: string }> =
        (form.getFieldValue('variants') as Array<{ sections?: string }>) ?? [];
      const next = variants.map((v) => ({ ...v, sections: snapshot }));
      form.setFieldValue('variants', next);
      message.success(t('importDone'));
    } catch (error) {
      console.error('Import page sections failed:', error);
      message.error(t('importFailed'));
    } finally {
      setImporting(false);
    }
  };

  const statusLabel = (status?: string) =>
    status ? t(`statusMap.${status}`) : '-';
  const assignmentLabel = (v?: string) =>
    v ? t(`assignmentMap.${v}`) : '-';

  const renderStatusActions = (record: Experiment) => {
    const nodes: React.ReactNode[] = [];
    if (record.status === 'EXPERIMENT_STATUS_DRAFT' || record.status === 'EXPERIMENT_STATUS_PAUSED') {
      nodes.push(
        <Button
          key="run"
          type="link"
          size="small"
          onClick={() => handleSetStatus(record, 'EXPERIMENT_STATUS_RUNNING', t('startSuccess'))}
        >
          {t('start')}
        </Button>,
      );
    }
    if (record.status === 'EXPERIMENT_STATUS_RUNNING') {
      nodes.push(
        <Button
          key="pause"
          type="link"
          size="small"
          onClick={() => handleSetStatus(record, 'EXPERIMENT_STATUS_PAUSED', t('pauseSuccess'))}
        >
          {t('pause')}
        </Button>,
      );
    }
    if (record.status === 'EXPERIMENT_STATUS_RUNNING' || record.status === 'EXPERIMENT_STATUS_PAUSED') {
      nodes.push(
        <Popconfirm key="end" title={t('endConfirm')} onConfirm={() => handleSetStatus(record, 'EXPERIMENT_STATUS_ENDED', t('endSuccess'))}>
          <Button type="link" size="small" danger>
            {t('end')}
          </Button>
        </Popconfirm>,
      );
    }
    return nodes;
  };

  const columns: ProColumns<Experiment>[] = [
    { title: 'ID', dataIndex: 'id', width: 70, search: false },
    {
      title: t('name'),
      dataIndex: 'name',
      width: 180,
      ellipsis: true,
    },
    { title: t('pageId'), dataIndex: 'pageId', width: 90, search: false },
    {
      title: t('assignmentType'),
      dataIndex: 'assignmentType',
      valueType: 'select',
      width: 120,
      fieldProps: {
        options: ['ASSIGNMENT_TYPE_RANDOM', 'ASSIGNMENT_TYPE_RULE'].map((v) => ({
          label: t(`assignmentMap.${v}`),
          value: v,
        })),
      },
      render: (_, record) => <Tag>{assignmentLabel(record.assignmentType)}</Tag>,
    },
    {
      title: t('status'),
      dataIndex: 'status',
      valueType: 'select',
      width: 100,
      fieldProps: {
        options: [
          'EXPERIMENT_STATUS_DRAFT',
          'EXPERIMENT_STATUS_RUNNING',
          'EXPERIMENT_STATUS_PAUSED',
          'EXPERIMENT_STATUS_ENDED',
        ].map((v) => ({ label: t(`statusMap.${v}`), value: v })),
      },
      render: (_, record) => (
        <Tag color={STATUS_COLORS[record.status ?? ''] ?? 'default'}>
          {statusLabel(record.status)}
        </Tag>
      ),
    },
    {
      title: t('variants'),
      search: false,
      render: (_, record) =>
        (record.variants ?? [])
          .map((v) => `${v.name}(${v.weight ?? 0}%)`)
          .join(' / ') || '-',
    },
    {
      title: t('updatedAt'),
      dataIndex: 'updatedAt',
      search: false,
      width: 160,
      render: (_, record) =>
        record.updatedAt ? new Date(String(record.updatedAt)).toLocaleString() : '-',
    },
    {
      title: t('action'),
      valueType: 'option',
      width: 280,
      fixed: 'right',
      render: (_, record) => (
        <Space size={0} wrap>
          <Button type="link" size="small" icon={<EditOutlined />} onClick={() => openEdit(record)}>
            {t('edit')}
          </Button>
          {renderStatusActions(record)}
          <Button
            type="link"
            size="small"
            icon={<FieldTimeOutlined />}
            onClick={() => setStatsTarget(record)}
          >
            {t('stats')}
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
        loading={listMutation.isPending}
        rowKey="id"
        columns={columns}
        request={loadItems}
        dataSource={items}
        toolBarRender={() => [
          <Button key="refresh" icon={<ReloadOutlined />} onClick={() => actionRef.current?.reload()} />,
          <Button key="create" type="primary" icon={<PlusOutlined />} onClick={openCreate}>
            {t('create')}
          </Button>,
        ]}
        pagination={{ pageSize: 20 }}
      />

      <Drawer
        open={drawerOpen}
        onClose={() => setDrawerOpen(false)}
        title={editing ? `${t('edit')} #${editing.id}` : t('create')}
        width={720}
        destroyOnHidden
        extra={
          <Button type="primary" loading={createMutation.isPending || updateMutation.isPending} onClick={handleSubmit}>
            {t('save')}
          </Button>
        }
      >
        <Form form={form} layout="vertical">
          <Form.Item name="name" label={t('name')} rules={[{ required: true }]}>
            <Input placeholder={t('namePlaceholder')} />
          </Form.Item>
          <Space size={16} style={{ display: 'flex' }}>
            <Form.Item name="pageId" label={t('pageId')} rules={[{ required: true }]} style={{ width: 200 }}>
              <InputNumber min={1} style={{ width: '100%' }} />
            </Form.Item>
            <Form.Item name="assignmentType" label={t('assignmentType')} rules={[{ required: true }]} style={{ width: 220 }}>
              <Select
                options={['ASSIGNMENT_TYPE_RANDOM', 'ASSIGNMENT_TYPE_RULE'].map((v) => ({
                  label: t(`assignmentMap.${v}`),
                  value: v,
                }))}
              />
            </Form.Item>
          </Space>
          <Form.Item name="description" label={t('description')}>
            <Input.TextArea rows={2} />
          </Form.Item>

          <Form.List name="variants">
            {(fields, { add, remove }) => (
              <div style={{ marginBottom: 16 }}>
                <div style={{ display: 'flex', justifyContent: 'space-between', marginBottom: 8 }}>
                  <strong>{t('variants')}</strong>
                  <Space>
                    <Button size="small" loading={importing} onClick={handleImportSections}>
                      {t('importSections')}
                    </Button>
                    <Button size="small" icon={<PlusOutlined />} onClick={() => add({ weight: 50 })}>
                      {t('addVariant')}
                    </Button>
                  </Space>
                </div>
                {fields.map((field) => (
                  <div
                    key={field.key}
                    style={{
                      border: '1px solid var(--ant-color-border)',
                      borderRadius: 8,
                      padding: '12px 12px 4px',
                      marginBottom: 12,
                    }}
                  >
                    <Space size={12} style={{ display: 'flex' }} align="baseline">
                      <Form.Item
                        name={[field.name, 'name']}
                        label={t('variantName')}
                        rules={[{ required: true }]}
                        style={{ width: 140 }}
                      >
                        <Input placeholder="A / B" />
                      </Form.Item>
                      <Form.Item
                        name={[field.name, 'weight']}
                        label={t('weight')}
                        style={{ width: 120 }}
                      >
                        <InputNumber min={0} max={100} style={{ width: '100%' }} />
                      </Form.Item>
                      <Form.Item name={[field.name, 'sortOrder']} label={t('sortOrder')} style={{ width: 100 }}>
                        <InputNumber min={0} style={{ width: '100%' }} />
                      </Form.Item>
                      <Button type="link" danger size="small" onClick={() => remove(field.name)}>
                        {t('removeVariant')}
                      </Button>
                    </Space>
                    <Form.Item
                      name={[field.name, 'conditions']}
                      label={t('conditions')}
                      tooltip={t('conditionsTooltip')}
                    >
                      <Input.TextArea rows={2} placeholder='[{"field":"device","op":"eq","value":"mobile"}]' />
                    </Form.Item>
                    <Form.Item
                      name={[field.name, 'sections']}
                      label={t('sections')}
                      tooltip={t('sectionsTooltip')}
                    >
                      <Input.TextArea rows={4} placeholder='{"sections":[...]}' />
                    </Form.Item>
                  </div>
                ))}
              </div>
            )}
          </Form.List>
        </Form>
      </Drawer>

      <Modal
        open={!!statsTarget}
        onCancel={() => setStatsTarget(undefined)}
        title={`${t('stats')} - ${statsTarget?.name ?? ''}`}
        footer={null}
        width={420}
      >
        {statsQuery.isPending && <div>{t('loading')}</div>}
        {!statsQuery.isPending && (statsQuery.data ?? []).length === 0 && (
          <div style={{ color: 'var(--ant-color-text-tertiary)' }}>{t('noData')}</div>
        )}
        {(statsQuery.data ?? []).map((row) => (
          <div
            key={row.variantName}
            style={{ display: 'flex', justifyContent: 'space-between', padding: '8px 0', borderBottom: '1px solid var(--ant-color-split)' }}
          >
            <span>{row.variantName}</span>
            <strong>{row.participants}</strong>
          </div>
        ))}
      </Modal>
    </ContentContainer>
  );
};

export default ExperimentManagement;
