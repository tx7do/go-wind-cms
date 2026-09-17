import { useRef, useState } from 'react';
import type { ProColumns, ActionType } from '@ant-design/pro-components';
import TableExportButton from '@/components/common/TableExportButton';
import { ProTable } from '@ant-design/pro-components';
import { Button, Popconfirm, Tag, Switch, App } from 'antd';
import {
  EditOutlined,
  DeleteOutlined,
  PlusOutlined,
  CaretRightOutlined,
  PauseOutlined,
  ReloadOutlined,
} from '@ant-design/icons';
import { useQueryClient } from '@tanstack/react-query';
import { useTranslation } from 'react-i18next';
import type { taskservicev1_Task as Task } from '@/api/generated/admin/service/v1';
import { PaginationQuery } from '@/core';
import {
  fetchListTasks,
  useDeleteTask,
  useUpdateTask,
  useControlTask,
  useStartAllTasks,
  useStopAllTasks,
  useRestartAllTasks,
} from '@/api/hooks/task';
import { useProTableScrollY } from '@/hooks/useProTableScrollY';
import ContentContainer from '@/layouts/components/PageContainer/ContentContainer';
import { getTaskTypeMap, getTaskTypeOptions } from './constants';
import TaskDrawer from './components/TaskDrawer';

/**
 * 任务管理页面
 */
const TaskManagement = () => {
  const { t } = useTranslation('task');
  const actionRef = useRef<ActionType>(null);
  const queryClient = useQueryClient();
  const { message } = App.useApp();

  const containerRef = useRef<HTMLDivElement>(null);
  const tableScrollY = useProTableScrollY(containerRef);

  // Drawer 状态管理
  const [drawerOpen, setDrawerOpen] = useState(false);
  const [drawerMode, setDrawerMode] = useState<'create' | 'edit'>('create');
  const [selectedTask, setSelectedTask] = useState<Task | undefined>();

  const taskTypeMap = getTaskTypeMap(t);

  // 删除操作
  const deleteMutation = useDeleteTask({
    onSuccess: () => {
      message.success(t('deleteSuccess'));
      actionRef.current?.reload();
      queryClient.invalidateQueries({ queryKey: ['listTasks'] });
    },
    onError: (error: Error) => {
      message.error(error.message || t('deleteFailed'));
    },
  });

  // 更新操作（用于启用/禁用切换）
  const updateMutation = useUpdateTask();

  // 记录当前正在切换启用/禁用的任务 id，使该行的 Switch 显示 loading，
  // 而不是所有行共享 updateMutation.isPending 导致全部开关同时转圈。
  const [loadingTaskId, setLoadingTaskId] = useState<number | undefined>(undefined);

  // 控制单个任务
  const controlMutation = useControlTask({
    onSuccess: () => {
      message.success(t('operationSuccess'));
      actionRef.current?.reload();
    },
    onError: () => {
      message.error(t('operationFailed'));
    },
  });

  // 全部启动
  const startAllMutation = useStartAllTasks({
    onSuccess: () => {
      message.success(t('operationSuccess'));
      actionRef.current?.reload();
    },
    onError: () => {
      message.error(t('operationFailed'));
    },
  });

  // 全部停止
  const stopAllMutation = useStopAllTasks({
    onSuccess: () => {
      message.success(t('operationSuccess'));
      actionRef.current?.reload();
    },
    onError: () => {
      message.error(t('operationFailed'));
    },
  });

  // 全部重启
  const restartAllMutation = useRestartAllTasks({
    onSuccess: () => {
      message.success(t('operationSuccess'));
      actionRef.current?.reload();
    },
    onError: () => {
      message.error(t('operationFailed'));
    },
  });

  // 启用/禁用切换
  const handleEnableChange = async (record: Task, checked: boolean) => {
    setLoadingTaskId(record.id as number);
    try {
      await updateMutation.mutateAsync({
        id: record.id as number,
        values: { enable: checked },
      });
      await controlMutation.mutateAsync({
        typeName: record.typeName as string,
        controlType: checked ? 'Start' : 'Stop',
      });
      message.success(t('updateStatusSuccess'));
      actionRef.current?.reload();
    } catch {
      message.error(t('updateStatusFailed'));
    } finally {
      setLoadingTaskId(undefined);
    }
  };

  // 列配置
  const columns: ProColumns<Task>[] = [
    {
      title: t('serial'),
      dataIndex: 'id',
      width: 60,
      hideInSearch: true,
      render: (_, record, index) => {
        const pagination = record.id !== undefined ? actionRef.current?.pageInfo : undefined;
        const page = pagination?.current || 1;
        const pageSize = pagination?.pageSize || 10;
        return (page - 1) * pageSize + index + 1;
      },
    },
    {
      title: t('type'),
      dataIndex: 'type',
      width: 110,
      valueType: 'select',
      fieldProps: {
        options: getTaskTypeOptions(t),
      },
      render: (_, record) => {
        const type = record.type as keyof typeof taskTypeMap;
        const config = taskTypeMap[type] || { text: type, color: 'default' };
        return <Tag color={config.color}>{config.text}</Tag>;
      },
    },
    {
      title: t('typeName'),
      dataIndex: 'typeName',
      hideInSearch: true,
    },
    {
      title: t('taskPayload'),
      dataIndex: 'taskPayload',
      hideInSearch: true,
      ellipsis: true,
    },
    {
      title: t('cronSpec'),
      dataIndex: 'cronSpec',
      hideInSearch: true,
      width: 140,
    },
    {
      title: t('enable'),
      dataIndex: 'enable',
      width: 100,
      hideInSearch: true,
      render: (_, record) => (
        <Switch
          checked={record.enable as boolean}
          loading={loadingTaskId === record.id}
          onChange={(checked) => handleEnableChange(record, checked)}
        />
      ),
    },
    {
      title: t('createdAt'),
      dataIndex: 'createdAt',
      width: 180,
      valueType: 'dateTime',
      hideInSearch: true,
    },
    {
      title: t('remark'),
      dataIndex: 'remark',
      hideInSearch: true,
      ellipsis: true,
    },
    {
      title: t('action'),
      valueType: 'option',
      // 操作列需容纳 执行记录+多个控制按钮，过窄会导致按钮溢出且点击被单元格遮挡（同步 vben 81342f5d）
      width: 300,
      render: (_, record) => [
        <a
          key="edit"
          onClick={() => {
            setDrawerMode('edit');
            setSelectedTask(record);
            setDrawerOpen(true);
          }}
        >
          <EditOutlined />
        </a>,
        <Popconfirm
          key="start"
          title={t('startTaskConfirmTitle')}
          description={t('startTaskConfirmDesc', { moduleName: t('moduleName') })}
          onConfirm={() =>
            controlMutation.mutate({ typeName: record.typeName as string, controlType: 'Start' })
          }
          okText={t('common:button.ok')}
          cancelText={t('common:button.cancel')}
        >
          <a style={{ color: '#52c41a' }}><CaretRightOutlined /></a>
        </Popconfirm>,
        <Popconfirm
          key="stop"
          title={t('stopTaskConfirmTitle')}
          description={t('stopTaskConfirmDesc', { moduleName: t('moduleName') })}
          onConfirm={() =>
            controlMutation.mutate({ typeName: record.typeName as string, controlType: 'Stop' })
          }
          okText={t('common:button.ok')}
          cancelText={t('common:button.cancel')}
        >
          <a style={{ color: 'var(--ant-color-error)' }}><PauseOutlined /></a>
        </Popconfirm>,
        <Popconfirm
          key="restart"
          title={t('restartTaskConfirmTitle')}
          description={t('restartTaskConfirmDesc', { moduleName: t('moduleName') })}
          onConfirm={() =>
            controlMutation.mutate({ typeName: record.typeName as string, controlType: 'Restart' })
          }
          okText={t('common:button.ok')}
          cancelText={t('common:button.cancel')}
        >
          <a><ReloadOutlined /></a>
        </Popconfirm>,
        <Popconfirm
          key="delete"
          title={t('deleteConfirmTitle')}
          description={t('deleteConfirmDesc', { moduleName: t('moduleName') })}
          onConfirm={() => record.id && deleteMutation.mutate({ id: record.id })}
          okText={t('common:button.ok')}
          cancelText={t('common:button.cancel')}
        >
          <a style={{ color: 'var(--ant-color-error)' }}><DeleteOutlined /></a>
        </Popconfirm>,
      ],
    },
  ];

  return (
    <>
      <ContentContainer heightMode="fixed" padding="16px" bottomMargin={0}>
        <div ref={containerRef} className="page-container-content">
          <ProTable<Task>
            actionRef={actionRef}
            columns={columns}
            request={async (params, _sorter, _filter) => {
              try {
                const query = new PaginationQuery({
                  paging: {
                    page: params.current || 1,
                    pageSize: params.pageSize || 10,
                  },
                  formValues: Object.fromEntries(
                    Object.entries(params).filter(
                      ([key]) => !['current', 'pageSize'].includes(key),
                    ),
                  ),
                });

                const response = await fetchListTasks(query);

                return {
                  data: response.items || [],
                  total: response.total || 0,
                  success: true,
                };
              } catch (error: any) {
                message.error(error.message || t('fetchFailed'));
                return { data: [], total: 0, success: false };
              }
            }}
            rowKey="id"
            search={{
              labelWidth: 'auto',
              defaultCollapsed: false,
            }}
            pagination={false}
            toolBarRender={() => [
              <TableExportButton key="export" fetcher={fetchListTasks} columns={columns} filename="tasks" />,
              <Button
                key="create"
                type="primary"
                icon={<PlusOutlined />}
                onClick={() => {
                  setDrawerMode('create');
                  setSelectedTask(undefined);
                  setDrawerOpen(true);
                }}
              >
                {t('create')}
              </Button>,
              <Popconfirm
                key="startAll"
                title={t('startAllConfirmTitle')}
                description={t('startAllConfirmDesc', { moduleName: t('moduleName') })}
                onConfirm={() => startAllMutation.mutate()}
                okText={t('common:button.ok')}
                cancelText={t('common:button.cancel')}
              >
                <Button
                  style={{ backgroundColor: '#52c41a', borderColor: '#52c41a', color: '#fff' }}
                  icon={<CaretRightOutlined />}
                  loading={startAllMutation.isPending}
                >
                  {t('startAll')}
                </Button>
              </Popconfirm>,
              <Popconfirm
                key="stopAll"
                title={t('stopAllConfirmTitle')}
                description={t('stopAllConfirmDesc', { moduleName: t('moduleName') })}
                onConfirm={() => stopAllMutation.mutate()}
                okText={t('common:button.ok')}
                cancelText={t('common:button.cancel')}
              >
                <Button danger icon={<PauseOutlined />} loading={stopAllMutation.isPending}>
                  {t('stopAll')}
                </Button>
              </Popconfirm>,
              <Popconfirm
                key="restartAll"
                title={t('restartAllConfirmTitle')}
                description={t('restartAllConfirmDesc', { moduleName: t('moduleName') })}
                onConfirm={() => restartAllMutation.mutate()}
                okText={t('common:button.ok')}
                cancelText={t('common:button.cancel')}
              >
                <Button icon={<ReloadOutlined />} loading={restartAllMutation.isPending}>
                  {t('restartAll')}
                </Button>
              </Popconfirm>,
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
            scroll={{ y: tableScrollY }}
          />
        </div>
      </ContentContainer>

      {/* 任务编辑/创建 Drawer */}
      <TaskDrawer
        open={drawerOpen}
        mode={drawerMode}
        data={selectedTask}
        onClose={() => {
          setDrawerOpen(false);
          setSelectedTask(undefined);
        }}
        onSuccess={() => {
          actionRef.current?.reload();
        }}
      />
    </>
  );
};

export default TaskManagement;
