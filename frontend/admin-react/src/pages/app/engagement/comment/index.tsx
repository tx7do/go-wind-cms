import { useRef, useState } from 'react';
import type { ProColumns, ActionType } from '@ant-design/pro-components';
import TableExportButton from '@/components/common/TableExportButton';
import { ProTable } from '@ant-design/pro-components';
import { Popconfirm, Tag, App, Modal, Input } from 'antd';
import { EditOutlined, DeleteOutlined, MessageOutlined } from '@ant-design/icons';
import { useQueryClient } from '@tanstack/react-query';
import { useTranslation } from 'react-i18next';
import dayjs from 'dayjs';

import { PaginationQuery } from '@/core';
import { TABLE } from '@/config/constants';
import { useProTableScrollY } from '@/hooks/useProTableScrollY';
import { useUserStore } from '@/stores/user';
import {
  fetchListComments,
  useCreateComment,
  useDeleteComment,
} from '@/api/hooks/comment';
import {
  getAuthorTypeColor,
  getAuthorTypeLabel,
  getContentTypeColor,
  getContentTypeLabel,
  getStatusColor,
  getStatusLabel,
  authorTypeOptions,
  contentTypeOptions,
  statusOptions,
} from './constants';
import CommentDrawer from './CommentDrawer';

/**
 * 评论管理列表页
 * 后端返回树形结构（回复挂在父评论 children 下），ProTable 原生渲染树形表格
 */
const CommentList = () => {
  const { t } = useTranslation('comment');
  const actionRef = useRef<ActionType>(null);
  const containerRef = useRef<HTMLDivElement>(null);
  const tableScrollY = useProTableScrollY(containerRef);
  const queryClient = useQueryClient();
  const { message } = App.useApp();
  const userInfo = useUserStore((s) => s.userInfo);

  const [drawerOpen, setDrawerOpen] = useState(false);
  const [editingRow, setEditingRow] = useState<any>(null);

  // 回复弹窗状态
  const [replyOpen, setReplyOpen] = useState(false);
  const [replySubmitting, setReplySubmitting] = useState(false);
  const [replyContent, setReplyContent] = useState('');
  const [replyRow, setReplyRow] = useState<any>(null);

  // 树形展开受控：数据异步加载后 defaultExpandAllRows 不生效，改为加载后全展开
  const [expandedKeys, setExpandedKeys] = useState<readonly any[]>([]);

  const deleteMutation = useDeleteComment({
    onSuccess: () => {
      message.success(t('deleteSuccess'));
      actionRef.current?.reload();
      queryClient.invalidateQueries({ queryKey: ['listComments'] });
    },
    onError: (error: Error) => {
      message.error(error.message || t('deleteFailed'));
    },
  });

  const createMutation = useCreateComment({
    onSuccess: () => {
      message.success(t('createSuccess'));
      queryClient.invalidateQueries({ queryKey: ['listComments'] });
      actionRef.current?.reload();
    },
    onError: (error: Error) => {
      message.error(error.message || t('createFailed'));
    },
  });

  /* 回复：以管理员身份对选中评论发表回复 */
  const handleReply = (row: any) => {
    setReplyRow(row);
    setReplyContent('');
    setReplyOpen(true);
  };

  const submitReply = async () => {
    const row = replyRow;
    if (!row || !replyContent.trim()) return;

    setReplySubmitting(true);
    try {
      await createMutation.mutateAsync({
        // 树形绑定 + 目标对象沿用父评论；replyToId 驱动前台"作者回复"徽章
        data: {
          parentId: row.id,
          replyToId: row.id,
          contentType: row.contentType,
          objectId: row.objectId,
          content: replyContent,
          // BFF 只盖 created_by，作者字段需前端提供；管理员回复直接以 APPROVED 发布
          authorType: 'AUTHOR_TYPE_ADMIN',
          authorId: userInfo?.id ?? 0,
          authorName: userInfo?.username,
          status: 'STATUS_APPROVED',
        } as any,
      });
      setReplyOpen(false);
      setReplyContent('');
    } finally {
      setReplySubmitting(false);
    }
  };

  const columns: ProColumns<any>[] = [
    {
      title: t('authorName'),
      dataIndex: 'authorName',
      width: 140,
    },
    {
      title: t('content'),
      dataIndex: 'content',
      ellipsis: true,
      render: (_, record) => {
        const text = String(record.content ?? '');
        return text.length > 50 ? `${text.slice(0, 50)}...` : text;
      },
    },
    {
      title: t('contentType'),
      dataIndex: 'contentType',
      valueType: 'select',
      width: 110,
      fieldProps: {
        options: contentTypeOptions(t),
        showSearch: true,
      },
      render: (_, record) => (
        <Tag color={getContentTypeColor(record.contentType)}>
          {getContentTypeLabel(t, record.contentType)}
        </Tag>
      ),
    },
    {
      title: t('authorType'),
      dataIndex: 'authorType',
      valueType: 'select',
      width: 110,
      fieldProps: {
        options: authorTypeOptions(t),
        showSearch: true,
      },
      render: (_, record) => (
        <Tag color={getAuthorTypeColor(record.authorType)}>
          {getAuthorTypeLabel(t, record.authorType)}
        </Tag>
      ),
    },
    {
      title: t('status'),
      dataIndex: 'status',
      valueType: 'select',
      width: 110,
      fieldProps: {
        options: statusOptions(t),
        showSearch: true,
      },
      render: (_, record) => (
        <Tag color={getStatusColor(record.status)}>
          {getStatusLabel(t, record.status)}
        </Tag>
      ),
    },
    {
      title: t('ipAddress'),
      dataIndex: 'ipAddress',
      search: false,
      width: 130,
    },
    {
      title: t('location'),
      dataIndex: 'location',
      search: false,
      width: 120,
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
      // 搜索专用：时间范围 → created_at__gte/lte
      title: t('createdAtRange'),
      dataIndex: 'created_at_range',
      valueType: 'dateTimeRange',
      hideInTable: true,
      fieldProps: { showTime: true, presets: [
        { label: t('dateRange.today'), value: [dayjs().startOf('day'), dayjs().endOf('day')] },
        { label: t('dateRange.yesterday'), value: [dayjs().subtract(1, 'day').startOf('day'), dayjs().subtract(1, 'day').endOf('day')] },
        { label: t('dateRange.thisWeek'), value: [dayjs().startOf('week'), dayjs().endOf('week')] },
        { label: t('dateRange.lastWeek'), value: [dayjs().subtract(1, 'week').startOf('week'), dayjs().subtract(1, 'week').endOf('week')] },
        { label: t('dateRange.thisMonth'), value: [dayjs().startOf('month'), dayjs().endOf('month')] },
        { label: t('dateRange.lastMonth'), value: [dayjs().subtract(1, 'month').startOf('month'), dayjs().subtract(1, 'month').endOf('month')] },
      ] },
      search: {
        transform: (value: any[]) => ({
          created_at__gte: value?.[0]
            ? dayjs(value[0]).format('YYYY-MM-DD HH:mm:ss')
            : undefined,
          created_at__lte: value?.[1]
            ? dayjs(value[1]).format('YYYY-MM-DD HH:mm:ss')
            : undefined,
        }),
      },
    },
    {
      title: t('action'),
      valueType: 'option',
      width: 150,
      fixed: 'right',
      render: (_, record) => [
        <a key="reply" onClick={() => handleReply(record)}>
          <MessageOutlined /> {t('reply')}
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

  return (
    <>
      <div ref={containerRef} className="page-container-content" style={{ padding: '0 8px', height: '100%' }}>
        <ProTable<any>
          actionRef={actionRef}
          columns={columns}
          headerTitle={t('commentList')}
          request={async (params) => {
            try {
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

              const response = await fetchListComments(query);

              // 全部展开：收集所有含 children 的父评论 id
              const collectParentIds = (items: any[]): any[] =>
                items.flatMap((item) =>
                  item.children?.length
                    ? [item.id, ...collectParentIds(item.children)]
                    : [],
                );
              setExpandedKeys(collectParentIds(response.items || []));

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
            span: 12,
          }}
          pagination={{
            defaultPageSize: TABLE.DEFAULT_PAGE_SIZE,
            showSizeChanger: true,
            showQuickJumper: true,
          }}
          toolBarRender={() => [
            <TableExportButton
              key="export"
              fetcher={fetchListComments}
              columns={columns}
              filename="comments"
            />,
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
          scroll={{ x: 1100, y: tableScrollY }}
          expandable={{
            expandedRowKeys: expandedKeys,
            onExpandedRowsChange: (keys: readonly any[]) => setExpandedKeys(keys),
          }}
        />
      </div>

      <CommentDrawer
        open={drawerOpen}
        data={editingRow}
        onClose={() => {
          setDrawerOpen(false);
          setEditingRow(null);
        }}
        onSuccess={() => {
          actionRef.current?.reload();
        }}
      />

      {/* 回复弹窗：以管理员身份回复选中评论 */}
      <Modal
        open={replyOpen}
        title={replyRow ? t('replyTo', { authorName: replyRow.authorName ?? '' }) : t('reply')}
        confirmLoading={replySubmitting}
        okText={t('reply')}
        cancelText={t('common:button.cancel')}
        onOk={submitReply}
        onCancel={() => setReplyOpen(false)}
        destroyOnHidden
      >
        <div style={{ display: 'flex', flexDirection: 'column', gap: 8 }}>
          <div>{t('content')}</div>
          <Input.TextArea
            rows={4}
            value={replyContent}
            onChange={(e) => setReplyContent(e.target.value)}
          />
        </div>
      </Modal>
    </>
  );
};

export default CommentList;
