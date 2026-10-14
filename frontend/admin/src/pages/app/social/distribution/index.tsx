import { useCallback, useMemo, useRef, useState } from 'react';
import type { ProColumns, ActionType } from '@ant-design/pro-components';
import ListTable from '@/components/common/ListTable';
import { Calendar, Card, Popconfirm, Segmented, Tag, App, Tooltip } from 'antd';
import { EditOutlined, DeleteOutlined, SendOutlined, SyncOutlined } from '@ant-design/icons';
import { useQueryClient } from '@tanstack/react-query';
import { useTranslation } from 'react-i18next';
import dayjs from 'dayjs';

import { PaginationQuery } from '@/core';
import { TABLE } from '@/config/constants';
import { useProTableScrollY } from '@/hooks/useProTableScrollY';
import {
  fetchListSocialDistributions,
  useDeleteSocialDistribution,
  usePublishSocialDistribution,
  useRefreshSocialDistributionStats,
} from '@/api/hooks/social';
import {
  getPlatformColor,
  getPlatformLabel,
  getStatusColor,
  getStatusLabel,
  platformOptions,
  statusOptions,
} from './constants';
import SocialDistributionDrawer from './SocialDistributionDrawer';

/**
 * 内容分发列表页：CMS 文章 → 社交平台的发布记录（状态/统计/失败原因）
 */
const SocialDistributionList = () => {
  const { t } = useTranslation('social');
  const actionRef = useRef<ActionType>(null);
  const containerRef = useRef<HTMLDivElement>(null);
  const tableScrollY = useProTableScrollY(containerRef);
  const queryClient = useQueryClient();
  const { message } = App.useApp();

  const [drawerOpen, setDrawerOpen] = useState(false);
  const [editingRow, setEditingRow] = useState<any>(null);
  const [calendarPrefillDate, setCalendarPrefillDate] = useState<string | undefined>();

  // ── 日历视图 ──
  const [viewMode, setViewMode] = useState<'list' | 'calendar'>('list');
  const [monthItems, setMonthItems] = useState<any[]>([]);
  const [calendarMonth, setCalendarMonth] = useState<dayjs.Dayjs>(dayjs());

  const fetchMonthItems = useCallback(async (month: dayjs.Dayjs) => {
    const query = new PaginationQuery({
      paging: { page: 1, pageSize: 100 },
      formValues: {
        publish_time__gte: month.startOf('month').format('YYYY-MM-DD HH:mm:ss'),
        publish_time__lte: month.endOf('month').format('YYYY-MM-DD HH:mm:ss'),
      },
      orderBy: ['publish_time'],
    });
    const resp = await fetchListSocialDistributions(query);
    setMonthItems(resp.items || []);
  }, []);

  const handleViewChange = (mode: 'list' | 'calendar') => {
    setViewMode(mode);
    if (mode === 'calendar') {
      fetchMonthItems(calendarMonth);
    }
  };

  const itemsByDay = useMemo(() => {
    const map: Record<string, any[]> = {};
    for (const item of monthItems) {
      if (!item.publishTime) continue;
      const key = dayjs(item.publishTime).format('YYYY-MM-DD');
      (map[key] ??= []).push(item);
    }
    return map;
  }, [monthItems]);

  const renderDateCell = (current: dayjs.Dayjs, info: any) => {
    if (info.type !== 'date') return null;
    const dayItems = itemsByDay[current.format('YYYY-MM-DD')] || [];
    if (!dayItems.length) return null;
    return (
      <div style={{ display: 'flex', flexDirection: 'column', gap: 2 }}>
        {dayItems.slice(0, 3).map((item: any) => (
          <div
            key={item.id}
            style={{ fontSize: 12, lineHeight: '18px', overflow: 'hidden', textOverflow: 'ellipsis', whiteSpace: 'nowrap' }}
          >
            <Tag color={getStatusColor(item.status)} style={{ marginRight: 4 }}>
              {dayjs(item.publishTime).format('HH:mm')}
            </Tag>
            {getPlatformLabel(t, item.platform)}
          </div>
        ))}
        {dayItems.length > 3 && (
          <span style={{ fontSize: 12, color: 'var(--ant-color-text-tertiary)' }}>+{dayItems.length - 3}</span>
        )}
      </div>
    );
  };

  const invalidate = () => {
    queryClient.invalidateQueries({ queryKey: ['listSocialDistributions'] });
    actionRef.current?.reload();
  };

  const deleteMutation = useDeleteSocialDistribution({
    onSuccess: () => {
      message.success(t('deleteSuccessDist'));
      invalidate();
    },
    onError: (error: Error) => {
      message.error(error.message || t('deleteFailedDist'));
    },
  });

  const publishMutation = usePublishSocialDistribution({
    onSuccess: () => {
      message.success(t('publishQueued'));
      invalidate();
    },
    onError: (error: Error) => {
      message.error(error.message || t('publishFailed'));
    },
  });

  const refreshStatsMutation = useRefreshSocialDistributionStats({
    onSuccess: (row) => {
      message.success(t('statsRefreshed', {
        likes: row.likes ?? 0,
        replies: row.replies ?? 0,
        reposts: row.reposts ?? 0,
      }));
      invalidate();
    },
    onError: (error: Error) => {
      message.error(error.message || t('statsRefreshFailed'));
    },
  });

  const columns: ProColumns<any>[] = [
    {
      title: t('platform'),
      dataIndex: 'platform',
      valueType: 'select',
      width: 110,
      fieldProps: {
        options: platformOptions(t),
        showSearch: true,
      },
      render: (_, record) => (
        <Tag color={getPlatformColor(record.platform)}>
          {getPlatformLabel(t, record.platform)}
        </Tag>
      ),
    },
    {
      title: t('socialAccount'),
      dataIndex: 'accountName',
      search: false,
      width: 140,
      ellipsis: true,
      render: (_, record) => record.accountName || `#${record.accountId}`,
    },
    {
      title: t('sourcePost'),
      dataIndex: 'postTitle',
      search: false,
      width: 160,
      ellipsis: true,
      render: (_, record) => record.postTitle || (record.postId ? `#${record.postId}` : '—'),
    },
    {
      title: t('overrideText'),
      dataIndex: 'overrideText',
      search: false,
      ellipsis: true,
      width: 220,
    },
    {
      title: t('status'),
      dataIndex: 'status',
      valueType: 'select',
      width: 100,
      fieldProps: {
        options: statusOptions(t),
        showSearch: true,
      },
      render: (_, record) => (
        <Tooltip title={record.status === 'STATUS_FAILED' ? record.errorMessage : undefined}>
          <Tag color={getStatusColor(record.status)}>{getStatusLabel(t, record.status)}</Tag>
        </Tooltip>
      ),
    },
    {
      title: t('publishTime'),
      dataIndex: 'publishTime',
      search: false,
      width: 145,
      render: (_, record) =>
        record.publishTime ? dayjs(record.publishTime).format('MM-DD HH:mm') : '—',
    },
    {
      title: t('publishedAt'),
      dataIndex: 'publishedAt',
      search: false,
      width: 145,
      render: (_, record) =>
        record.publishedAt ? dayjs(record.publishedAt).format('MM-DD HH:mm') : '—',
    },
    {
      title: t('stats'),
      key: 'stats',
      search: false,
      width: 140,
      render: (_, record) =>
        record.status === 'STATUS_PUBLISHED' ? (
          <span>
            ❤ {record.likes ?? 0} · 💬 {record.replies ?? 0} · 🔁 {record.reposts ?? 0}
          </span>
        ) : (
          '—'
        ),
    },
    {
      title: t('externalLink'),
      dataIndex: 'externalUri',
      search: false,
      width: 100,
      render: (_, record) =>
        record.externalUri ? (
          <a href={record.externalUri} target="_blank" rel="noreferrer">
            {t('view')}
          </a>
        ) : (
          '—'
        ),
    },
    {
      title: t('action'),
      valueType: 'option',
      width: 190,
      fixed: 'right',
      render: (_, record) => {
        const actions: any[] = [];
        const publishable = ['STATUS_DRAFT', 'STATUS_SCHEDULED', 'STATUS_FAILED'].includes(
          record.status,
        );
        if (publishable) {
          actions.push(
            <a
              key="publish"
              onClick={() => record.id && publishMutation.mutate({ id: record.id })}
            >
              <SendOutlined /> {t('publishNow')}
            </a>,
            <a
              key="edit"
              onClick={() => {
                setEditingRow(record);
                setDrawerOpen(true);
              }}
            >
              <EditOutlined />
            </a>,
          );
        }
        if (record.status === 'STATUS_PUBLISHED' && record.externalId) {
          actions.push(
            <a
              key="stats"
              onClick={() => record.id && refreshStatsMutation.mutate({ id: record.id })}
            >
              <SyncOutlined /> {t('refreshStats')}
            </a>,
          );
        }
        actions.push(
          <Popconfirm
            key="delete"
            title={t('deleteConfirmTitle')}
            description={t('deleteDistributionConfirmDesc')}
            onConfirm={() => {
              if (record.id) deleteMutation.mutate({ id: record.id });
            }}
            okText={t('common:button.ok')}
            cancelText={t('common:button.cancel')}
          >
            <a style={{ color: 'var(--ant-color-error)' }}>
              <DeleteOutlined />
            </a>
          </Popconfirm>,
        );
        return actions;
      },
    },
  ];

  return (
    <>
      <div ref={containerRef} className="page-container-content" style={{ padding: '0 8px', height: '100%', display: 'flex', flexDirection: 'column' }}>
        <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: 8 }}>
          <Segmented
            value={viewMode}
            onChange={(v) => handleViewChange(v as 'list' | 'calendar')}
            options={[
              { label: t('viewList'), value: 'list' },
              { label: t('viewCalendar'), value: 'calendar' },
            ]}
          />
          {viewMode === 'calendar' && (
            <span style={{ color: 'var(--ant-color-text-tertiary)', fontSize: 12 }}>{t('calendarHint')}</span>
          )}
        </div>
        {viewMode === 'calendar' ? (
          <Card bordered={false} style={{ flex: 1, minHeight: 0, overflow: 'auto' }}>
            <Calendar
              cellRender={renderDateCell}
              onPanelChange={(date) => {
                setCalendarMonth(date);
                fetchMonthItems(date);
              }}
              onSelect={(date, info) => {
                if (info?.source === 'date') {
                  setEditingRow(null);
                  setCalendarPrefillDate(date.hour(10).minute(0).second(0).toISOString());
                  setDrawerOpen(true);
                }
              }}
            />
          </Card>
        ) : (
        <ListTable<any>
          actionRef={actionRef}
          columns={columns}
          headerTitle={t('distributionList')}
          request={async (params) => {
            const query = new PaginationQuery({
              paging: {
                page: params.current || 1,
                pageSize: params.pageSize || TABLE.DEFAULT_PAGE_SIZE,
              },
              formValues: Object.fromEntries(
                Object.entries(params).filter(
                  ([key]) => !['current', 'pageSize'].includes(key),
                ),
              ),
              orderBy: ['-created_at'],
            });

            const response = await fetchListSocialDistributions(query);

            return {
              data: response.items || [],
              total: Number(response.total) || 0,
              success: true,
            };
          }}
          rowKey="id"
          search={{
            labelWidth: 'auto',
            defaultCollapsed: false,
            span: 12,
          }}
          pagination={{
            defaultPageSize: TABLE.DEFAULT_PAGE_SIZE,
            showSizeChanger: true,
            showQuickJumper: true,
          }}
          toolBarRender={() => [
            <a
              key="create"
              onClick={() => {
                setEditingRow(null);
                setDrawerOpen(true);
              }}
            >
              + {t('createDistribution')}
            </a>,
          ]}
          options={{
            density: false,
            fullScreen: false,
            setting: false,
            reload: true,
          }}
          size="small"
          bordered
          cardBordered={false}
          scroll={{ x: 1500, y: tableScrollY }}
        />
        )}
      </div>

      <SocialDistributionDrawer
        open={drawerOpen}
        data={editingRow}
        prefill={editingRow ? undefined : { publishTime: calendarPrefillDate }}
        onClose={() => {
          setDrawerOpen(false);
          setEditingRow(null);
          setCalendarPrefillDate(undefined);
        }}
        onSuccess={() => {
          actionRef.current?.reload();
        }}
      />
    </>
  );
};

export default SocialDistributionList;
