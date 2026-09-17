import { useEffect, useRef, useState } from 'react';
import type { ProColumns, ActionType } from '@ant-design/pro-components';
import TableExportButton from '@/components/common/TableExportButton';
import { ProTable } from '@ant-design/pro-components';
import { Button, Popconfirm, Tabs, Tag, App } from 'antd';
import { EditOutlined, DeleteOutlined, PlusOutlined } from '@ant-design/icons';
import { useQueryClient } from '@tanstack/react-query';
import { useTranslation } from 'react-i18next';
import { useNavigate } from 'react-router-dom';
import dayjs from 'dayjs';

import { PaginationQuery } from '@/core';
import { TABLE } from '@/config/constants';
import { useProTableScrollY } from '@/hooks/useProTableScrollY';
import { useDeletePost, fetchListPosts } from '@/api/hooks/post';
import { fetchFlattenedCategoryOptions } from '@/api/hooks/category';
import {
  getBoolLabel,
  getEditorTypeColor,
  getEditorTypeLabel,
  getPostStatusColor,
  getPostStatusLabel,
  postStatusOptions,
  pickTranslation,
} from './constants';

/** 状态 Tab：ALL 之外为 POST_STATUS_* 枚举名 */
const ALL_STATUS = 'ALL';
const COUNT_REFRESH_THROTTLE_MS = 3000;

// 列表字段掩码（对齐 vben 侧，避免大字段拖慢列表）
const LIST_FIELD_MASK =
  'id,status,sort_order,is_featured,author_name,available_languages,created_at,code,editor_type,disallow_comment,in_progress,auto_summary,is_featured,translations.id,translations.post_id,translations.language_code,translations.title,translations.summary,';

/**
 * 帖子管理列表页（状态 Tab 带计数；Tab 是唯一的状态筛选入口，重置不清 Tab 是有意行为）
 */
const PostList = () => {
  const { t, i18n } = useTranslation('post');
  const actionRef = useRef<ActionType>(null);
  const containerRef = useRef<HTMLDivElement>(null);
  const tableScrollY = useProTableScrollY(containerRef);
  const queryClient = useQueryClient();
  const navigate = useNavigate();
  const { message } = App.useApp();

  const [activeStatus, setActiveStatus] = useState<string>(ALL_STATUS);
  const [statusCounts, setStatusCounts] = useState<Record<string, number>>({});
  const [categoryOptions, setCategoryOptions] = useState<
    { label: string; value: number }[]
  >([]);

  const countThrottleRef = useRef({ lastAt: 0, inFlight: false, queued: false });

  const deleteMutation = useDeletePost({
    onSuccess: () => {
      message.success(t('deleteSuccess'));
      actionRef.current?.reload();
      queryClient.invalidateQueries({ queryKey: ['listPosts'] });
      refreshCounts(true);
    },
    onError: (error: Error) => {
      message.error(error.message || t('deleteFailed'));
    },
  });

  /**
   * 刷新各状态 Tab 的计数（每次一个 pageSize=1 的计数查询，取 total）
   * uint64 经 JSON 序列化是字符串，需强转数字
   */
  const refreshCounts = async (force = false) => {
    const now = Date.now();
    const gate = countThrottleRef.current;
    if (!force && now - gate.lastAt < COUNT_REFRESH_THROTTLE_MS) return;
    if (gate.inFlight) {
      gate.queued = true;
      return;
    }
    gate.inFlight = true;
    gate.lastAt = now;
    try {
      const keys = [ALL_STATUS, ...postStatusOptions(t).map((o) => o.value as string)];
      const entries = await Promise.all(
        keys.map(async (key) => {
          try {
            const resp = await fetchListPosts(
              new PaginationQuery({
                paging: { page: 1, pageSize: 1 },
                formValues: key === ALL_STATUS ? {} : { status: key },
                fieldMask: 'id',
              }),
            );
            return [key, resp.total] as const;
          } catch {
            return [key, undefined] as const;
          }
        }),
      );
      const next: Record<string, number> = {};
      for (const [key, total] of entries) {
        const n = Number(total);
        if (Number.isFinite(n)) {
          next[key] = n;
        }
      }
      setStatusCounts(next);
    } finally {
      gate.inFlight = false;
      if (gate.queued) {
        gate.queued = false;
        refreshCounts(true);
      }
    }
  };

  useEffect(() => {
    refreshCounts();
    fetchFlattenedCategoryOptions(i18n.language)
      .then(setCategoryOptions)
      .catch(() => setCategoryOptions([]));
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

  const columns: ProColumns<any>[] = [
    {
      title: t('postTitle'),
      dataIndex: 'translations.title',
      fixed: 'left',
      ellipsis: true,
      // 标题在翻译层，列表查询不带该过滤（对齐 vben 侧搜索字段）
      search: false,
      render: (_, record) => pickTranslation(record, i18n.language)?.title || '',
    },
    {
      title: t('slug'),
      dataIndex: 'code',
      ellipsis: true,
    },
    {
      // 仅搜索入口：分类筛选（查询时转 category_ids__in 数组走连接表）
      title: t('category'),
      dataIndex: 'categoryIds',
      valueType: 'select',
      hideInTable: true,
      fieldProps: {
        options: categoryOptions,
        showSearch: true,
        mode: 'multiple',
        allowClear: true,
        filterOption: (input: string, option: any) =>
          String(option?.label ?? '').toLowerCase().includes(input.toLowerCase()),
      },
    },
    {
      title: t('editorType'),
      dataIndex: 'editorType',
      search: false,
      width: 110,
      render: (_, record) => (
        <Tag color={getEditorTypeColor(record.editorType)}>
          {getEditorTypeLabel(t, record.editorType)}
        </Tag>
      ),
    },
    {
      title: t('authorName'),
      dataIndex: 'authorName',
      search: false,
      ellipsis: true,
      width: 120,
    },
    {
      title: t('status'),
      dataIndex: 'status',
      search: false,
      width: 100,
      render: (_, record) => (
        <Tag color={getPostStatusColor(record.status)}>
          {getPostStatusLabel(t, record.status)}
        </Tag>
      ),
    },
    {
      title: t('sortOrder'),
      dataIndex: 'sortOrder',
      search: false,
      width: 80,
    },
    {
      title: t('disallowComment'),
      dataIndex: 'disallowComment',
      search: false,
      width: 90,
      render: (_, record) => getBoolLabel(t, record.disallowComment),
    },
    {
      title: t('isFeatured'),
      dataIndex: 'isFeatured',
      search: false,
      width: 80,
      render: (_, record) => getBoolLabel(t, record.isFeatured),
    },
    {
      title: t('createdAt'),
      dataIndex: 'createdAt',
      search: false,
      width: 150,
      render: (_, record) =>
        record.createdAt ? dayjs(record.createdAt).format('YYYY-MM-DD HH:mm') : '',
    },
    {
      title: t('action'),
      valueType: 'option',
      width: 100,
      fixed: 'right',
      render: (_, record) => [
        <a
          key="edit"
          onClick={() =>
            navigate(`/content/posts/edit/${record.id}?lang=${i18n.language}`)
          }
        >
          <EditOutlined />
        </a>,
        <Popconfirm
          key="delete"
          title={t('deleteConfirmTitle')}
          description={t('deleteConfirmDesc', { moduleName: t('moduleName') })}
          onConfirm={() => {
            record.id && deleteMutation.mutate({ id: record.id });
          }}
          okText={t('common:button.ok')}
          cancelText={t('common:button.cancel')}
        >
          <a style={{ color: 'var(--ant-color-error)' }}>
            <DeleteOutlined />
          </a>
        </Popconfirm>,
      ],
    },
  ];

  const statusTabs = [
    { key: ALL_STATUS, label: t('statusAll') },
    ...postStatusOptions(t).map((o) => ({ key: o.value as string, label: o.label as string })),
  ];

  return (
    <div ref={containerRef} className="page-container-content" style={{ padding: '0 8px', height: '100%' }}>
      <ProTable<any>
        actionRef={actionRef}
        columns={columns}
        headerTitle={
          <div style={{ display: 'flex', alignItems: 'center', gap: 16, flexWrap: 'wrap' }}>
            <span>{t('postList')}</span>
            <Tabs
              size="small"
              activeKey={activeStatus}
              onChange={(key) => {
                setActiveStatus(key);
                actionRef.current?.reload();
              }}
              items={statusTabs.map((tab) => ({
                key: tab.key,
                label: (
                  <span>
                    {tab.label}
                    {statusCounts[tab.key] !== undefined && (
                      <span style={{ marginLeft: 4, fontSize: 12, opacity: 0.6 }}>
                        {statusCounts[tab.key]}
                      </span>
                    )}
                  </span>
                ),
              }))}
              style={{ marginBottom: 0 }}
            />
          </div>
        }
        // activeStatus 变化会触发重新请求（Tab 是唯一状态筛选入口）
        params={{ activeStatus }}
        request={async (params) => {
          try {
            const query: Record<string, unknown> = Object.fromEntries(
              Object.entries(params).filter(
                ([key]) => !['current', 'pageSize', 'activeStatus'].includes(key),
              ),
            );
            // 状态 Tab 是唯一的状态筛选入口，ALL 表示不过滤
            if (activeStatus === ALL_STATUS) {
              delete query.status;
            } else {
              query.status = activeStatus;
            }
            // 分类过滤走连接表：字段名必须带 __in 后缀且值为数组，
            // 通用转换器对单值生成 Value，而 post_repo 的连接表过滤只认 Values
            if (query.categoryIds != null) {
              const v = query.categoryIds;
              delete query.categoryIds;
              query.category_ids__in = Array.isArray(v) ? v : [v];
            }

            const response = await fetchListPosts(
              new PaginationQuery({
                paging: {
                  page: params.current || 1,
                  pageSize: params.pageSize || TABLE.DEFAULT_PAGE_SIZE,
                },
                formValues: query,
                fieldMask: LIST_FIELD_MASK,
              }),
            );

            return {
              data: response.items || [],
              total: Number(response.total) || 0,
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
          span: 12,
        }}
        pagination={{
          defaultPageSize: TABLE.DEFAULT_PAGE_SIZE,
          showSizeChanger: true,
        }}
        toolBarRender={() => [
          <TableExportButton
            key="export"
            fetcher={fetchListPosts}
            columns={columns}
            filename="posts"
          />,
          <Button
            key="create"
            type="primary"
            icon={<PlusOutlined />}
            size="small"
            onClick={() => navigate(`/content/posts/create?lang=${i18n.language}`)}
          >
            {t('create')}
          </Button>,
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
        scroll={{ x: 1300, y: tableScrollY }}
      />
    </div>
  );
};

export default PostList;
