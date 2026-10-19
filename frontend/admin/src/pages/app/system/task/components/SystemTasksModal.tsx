import { Modal, Table, Tag, Typography } from 'antd';
import type { ColumnsType } from 'antd/es/table';
import { useTranslation } from 'react-i18next';
import type {
  taskservicev1_InspectSystemTasksResponse,
  taskservicev1_SystemTaskFailure,
  taskservicev1_SystemTaskScheduleEntry,
  taskservicev1_SystemTaskStateSummary,
} from '@/api/generated/admin/service/v1';
import { useInspectSystemTasks } from '@/api/hooks/task';
import { formatDateTime } from '@/utils/date';

const { Text } = Typography;

type ScheduleEntry = taskservicev1_SystemTaskScheduleEntry;
type StateSummary = taskservicev1_SystemTaskStateSummary;

/**
 * 系统级常驻任务弹窗（方案 C'：只读 asynq Inspector，零新表）。
 *
 * 数据来自 asynq 本身：
 *  - 调度计划表：cron + 下次/上次入队时间（上次为空 = 注册后从未运行）；
 *  - 运行状态表：active/pending/retry/archived 计数 + 最近失败明细
 *    （错误原文/失败时间/已重试次数，失败时间倒序、每类型最多 10 条）。
 * asynq 队列无命名空间：后端只汇总调度条目覆盖到的任务类型，
 * 共库错开 URI 后看到的即本实例自己的任务。
 */
const SystemTasksModal = ({
  open,
  onClose,
}: {
  open: boolean;
  onClose: () => void;
}) => {
  const { t } = useTranslation('task');
  const inspectQuery = useInspectSystemTasks({ enabled: open });

  const data: taskservicev1_InspectSystemTasksResponse | undefined =
    inspectQuery.data;

  const scheduleColumns: ColumnsType<ScheduleEntry> = [
    {
      title: t('sysTaskType'),
      dataIndex: 'taskType',
      render: (v: string) => <Tag color="geekblue">{v}</Tag>,
    },
    { title: t('sysTaskCron'), dataIndex: 'cronSpec', width: 130 },
    {
      title: t('sysTaskNext'),
      dataIndex: 'nextEnqueueAt',
      width: 170,
      render: (v?: string) => (v ? formatDateTime(v) : '-'),
    },
    {
      title: t('sysTaskPrev'),
      dataIndex: 'prevEnqueueAt',
      width: 170,
      render: (v?: string) =>
        v ? (
          formatDateTime(v)
        ) : (
          <Text type="secondary">{t('sysTaskNeverRun')}</Text>
        ),
    },
  ];

  const summaryColumns: ColumnsType<StateSummary> = [
    {
      title: t('sysTaskType'),
      dataIndex: 'taskType',
      render: (v: string) => <Tag color="geekblue">{v}</Tag>,
    },
    { title: t('sysTaskActive'), dataIndex: 'active', width: 70, align: 'center' },
    { title: t('sysTaskPending'), dataIndex: 'pending', width: 80, align: 'center' },
    {
      title: t('sysTaskRetry'),
      dataIndex: 'retry',
      width: 70,
      align: 'center',
      render: (v: number) => (v > 0 ? <Tag color="warning">{v}</Tag> : v),
    },
    {
      title: t('sysTaskArchived'),
      dataIndex: 'archived',
      width: 80,
      align: 'center',
      render: (v: number) => (v > 0 ? <Tag color="error">{v}</Tag> : v),
    },
  ];

  const failureColumns: ColumnsType<taskservicev1_SystemTaskFailure> = [
    { title: t('sysTaskType'), dataIndex: 'taskType', width: 160 },
    {
      title: t('sysTaskState'),
      dataIndex: 'state',
      width: 90,
      render: (v: string) => (
        <Tag color={v === 'archived' ? 'error' : 'warning'}>{v}</Tag>
      ),
    },
    { title: t('sysTaskLastError'), dataIndex: 'lastError', ellipsis: true },
    {
      title: t('sysTaskFailedAt'),
      dataIndex: 'lastFailedAt',
      width: 170,
      render: (v?: string) => (v ? formatDateTime(v) : '-'),
    },
    {
      title: t('sysTaskRetried'),
      key: 'retried',
      width: 90,
      align: 'center',
      render: (_: unknown, r) => `${r.retried ?? 0}/${r.maxRetry ?? 0}`,
    },
  ];

  const summaries: StateSummary[] = data?.summaries ?? [];
  const failures: taskservicev1_SystemTaskFailure[] = summaries.flatMap(
    (s) => s.recentFailures ?? [],
  );

  return (
    <Modal
      title={t('sysTasksTitle')}
      open={open}
      onCancel={onClose}
      footer={null}
      width={860}
      destroyOnHidden
    >
      {inspectQuery.isFetching && (
        <div style={{ padding: 48, textAlign: 'center' }}>{t('sysTaskLoading')}</div>
      )}
      {inspectQuery.isError && (
        <Text type="danger">
          {inspectQuery.error?.message || t('sysTaskLoadFailed')}
        </Text>
      )}
      {data && (
        <>
          <Typography.Title level={5} style={{ marginBottom: 8 }}>
            {t('sysTaskSchedules')}
            <Text type="secondary" style={{ marginLeft: 8, fontSize: 12 }}>
              {t('sysTaskQueueHint', { queue: data.queue ?? '' })}
            </Text>
          </Typography.Title>
          <Table
            columns={scheduleColumns}
            dataSource={(data.schedules ?? []) as ScheduleEntry[]}
            pagination={false}
            rowKey="taskType"
            size="small"
            bordered
          />

          <Typography.Title level={5} style={{ margin: '20px 0 8px' }}>
            {t('sysTaskStates')}
          </Typography.Title>
          <Table
            columns={summaryColumns}
            dataSource={summaries}
            pagination={false}
            rowKey="taskType"
            size="small"
            bordered
          />

          <Typography.Title level={5} style={{ margin: '20px 0 8px' }}>
            {t('sysTaskFailures')}
          </Typography.Title>
          {failures.length === 0 ? (
            <Text type="secondary">{t('sysTaskNoFailures')}</Text>
          ) : (
            <Table
              columns={failureColumns}
              dataSource={failures}
              pagination={false}
              rowKey={(r) =>
                `${r.taskType}-${r.state}-${r.lastFailedAt ?? ''}-${r.retried ?? 0}`
              }
              size="small"
              bordered
            />
          )}
        </>
      )}
    </Modal>
  );
};

export default SystemTasksModal;
