import { useRef, useState } from 'react';
import type { ProColumns, ActionType } from '@ant-design/pro-components';
import ListTable from '@/components/common/ListTable';
import {
  DrawerForm,
  ProFormText,
  ProFormTextArea,
  ProFormSwitch,
} from '@ant-design/pro-components';
import { Button, Popconfirm, Tag, App } from 'antd';
import { EditOutlined, DeleteOutlined, PlusOutlined, HistoryOutlined } from '@ant-design/icons';
import { useQueryClient } from '@tanstack/react-query';
import { useTranslation } from 'react-i18next';
import type { aiservicev1_AiAgentApp as AiAgentApp } from '@/api/generated/admin/service/v1';
import { PaginationQuery } from '@/core';
import { TABLE } from '@/config/constants';
import {
  fetchListAiAgentApps,
  useDeleteAiAgentApp,
  useCreateAiAgentApp,
  useUpdateAiAgentApp,
} from '@/api/hooks/ai-agent';
import { useProTableScrollY } from '@/hooks/useProTableScrollY';
import ContentContainer from '@/layouts/components/PageContainer/ContentContainer';
import AgentRunsDrawer from './agent-runs-drawer';

// 从未知错误对象取本地化 message（request-client 拦截器已把 reason 译文写进 message）
const errMessage = (error: unknown, fallback: string): string => {
  if (error && typeof error === 'object' && 'message' in error) {
    const m = (error as { message?: unknown }).message;
    if (typeof m === 'string' && m) return m;
  }
  return fallback;
};

/** 图 DSL 文本粗校验：非空时必须是合法 JSON 对象（结构强校验在保存时由后端执行引擎切片承担） */
const validateGraphJson = (_rule: unknown, value: string | undefined): Promise<void> => {
  if (!value || !value.trim()) return Promise.resolve();
  try {
    const parsed = JSON.parse(value);
    if (typeof parsed !== 'object' || parsed === null || Array.isArray(parsed)) {
      return Promise.reject(new Error('graph must be a JSON object'));
    }
    return Promise.resolve();
  } catch {
    return Promise.reject(new Error('invalid JSON'));
  }
};

/**
 * AI Agent 编排应用列表（切片①）：
 * 基本信息与图 DSL（JSON 文本）在这里维护；可视化画布（@xyflow/react）与
 * 试运行在后续切片接入——图的结构校验届时由后端保存路径强校验。
 */
const AiAgentAppManagement = () => {
  const { t } = useTranslation('ai-agent');
  const actionRef = useRef<ActionType>(null);
  const queryClient = useQueryClient();
  const { message } = App.useApp();
  const containerRef = useRef<HTMLDivElement>(null);
  const tableScrollY = useProTableScrollY(containerRef);

  const [formOpen, setFormOpen] = useState(false);
  const [editing, setEditing] = useState<AiAgentApp | undefined>();
  const [runsApp, setRunsApp] = useState<AiAgentApp | undefined>();

  const refresh = () => {
    actionRef.current?.reload();
    queryClient.invalidateQueries({ queryKey: ['listAiAgentApps'] });
  };

  const deleteMutation = useDeleteAiAgentApp();
  const createMutation = useCreateAiAgentApp();
  const updateMutation = useUpdateAiAgentApp();

  const handleDelete = async (record: AiAgentApp) => {
    if (!record.id) return;
    try {
      await deleteMutation.mutateAsync({ id: record.id });
      message.success(t('deleteSuccess'));
      refresh();
    } catch (error: any) {
      // 原始错误必须留在控制台：用户可见的那句翻译不包含服务端的原因
      console.error('delete ai agent app failed', error);
      message.error(errMessage(error, t('deleteFailed')));
    }
  };

  const handleSubmit = async (values: Record<string, any>) => {
    try {
      if (!editing?.id) {
        await createMutation.mutateAsync({ data: values as AiAgentApp });
        message.success(t('createSuccess'));
      } else {
        await updateMutation.mutateAsync({
          id: editing.id,
          data: { id: editing.id, ...values } as AiAgentApp,
          updateMask: 'name,description,enabled,graph',
        });
        message.success(t('updateSuccess'));
      }
      setFormOpen(false);
      refresh();
      return true;
    } catch (error: any) {
      console.error('save ai agent app failed', error);
      message.error(errMessage(error, t('saveFailed')));
      return false;
    }
  };

  const columns: ProColumns<AiAgentApp>[] = [
    {
      title: t('name'),
      dataIndex: 'name',
      width: 200,
      ellipsis: true,
    },
    {
      title: t('description'),
      dataIndex: 'description',
      ellipsis: true,
      render: (_, record) => record.description || '-',
    },
    {
      title: t('enabled'),
      dataIndex: 'enabled',
      width: 100,
      render: (_, record) =>
        record.enabled ? <Tag color="success">{t('enabledOn')}</Tag> : <Tag>{t('enabledOff')}</Tag>,
    },
    {
      title: t('updatedAt'),
      dataIndex: 'updatedAt',
      width: 170,
    },
    {
      title: t('actions'),
      valueType: 'option',
      width: 220,
      fixed: 'right',
      render: (_, record) => [
        <Button
          key="edit"
          type="link"
          size="small"
          icon={<EditOutlined />}
          onClick={() => {
            setEditing(record);
            setFormOpen(true);
          }}
        >
          {t('edit')}
        </Button>,
        <Button
          key="runs"
          type="link"
          size="small"
          icon={<HistoryOutlined />}
          onClick={() => setRunsApp(record)}
        >
          {t('openRuns')}
        </Button>,
        <Popconfirm key="delete" title={t('deleteConfirm')} onConfirm={() => handleDelete(record)}>
          <Button danger type="link" size="small" icon={<DeleteOutlined />}>
            {t('delete')}
          </Button>
        </Popconfirm>,
      ],
    },
  ];

  return (
    <ContentContainer heightMode="fixed" padding="16px" bottomMargin={0}>
      <div ref={containerRef} className="page-container-content">
        <ListTable<AiAgentApp>
          actionRef={actionRef}
          columns={columns}
          request={async (params) => {
            const query = new PaginationQuery({
              paging: { page: params.current || 1, pageSize: params.pageSize || 20 },
              formValues: Object.fromEntries(
                Object.entries(params).filter(([key]) => !['current', 'pageSize'].includes(key)),
              ),
            });
            const response = await fetchListAiAgentApps(query);
            return { data: response.items || [], total: Number(response.total || 0), success: true };
          }}
          rowKey="id"
          search={false}
          pagination={{ defaultPageSize: TABLE.DEFAULT_PAGE_SIZE, showSizeChanger: true }}
          options={{ density: true, fullScreen: true, setting: true, reload: true }}
          toolBarRender={() => [
            <Button
              key="create"
              type="primary"
              icon={<PlusOutlined />}
              onClick={() => {
                setEditing(undefined);
                setFormOpen(true);
              }}
            >
              {t('create')}
            </Button>,
          ]}
          size="middle"
          bordered
          scroll={{ y: tableScrollY }}
        />
      </div>

      <DrawerForm
        title={editing ? t('edit') : t('create')}
        width={560}
        open={formOpen}
        onOpenChange={setFormOpen}
        drawerProps={{ destroyOnHidden: true }}
        onFinish={handleSubmit}
        initialValues={
          editing
            ? { ...editing, graph: editing.graph ?? '' }
            : { enabled: true, graph: '' }
        }
      >
        <ProFormText
          name="name"
          label={t('name')}
          rules={[{ required: true, message: t('requiredName') }]}
        />
        <ProFormTextArea name="description" label={t('description')} fieldProps={{ rows: 2 }} />
        <ProFormSwitch name="enabled" label={t('enabled')} />
        <ProFormTextArea
          name="graph"
          label="Graph DSL (JSON)"
          fieldProps={{ rows: 10, style: { fontFamily: 'monospace', fontSize: 12 } }}
          extra={t('graphJsonHint')}
          rules={[{ validator: validateGraphJson }]}
        />
      </DrawerForm>

      <AgentRunsDrawer open={!!runsApp} app={runsApp} onClose={() => setRunsApp(undefined)} />
    </ContentContainer>
  );
};

export default AiAgentAppManagement;
