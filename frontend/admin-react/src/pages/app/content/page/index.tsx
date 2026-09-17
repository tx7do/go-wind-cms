import { useRef, useState } from 'react';
import type { ProColumns, ActionType } from '@ant-design/pro-components';
import TableExportButton from '@/components/common/TableExportButton';
import { ProTable } from '@ant-design/pro-components';
import { Button, Popconfirm, Tag, App } from 'antd';
import { EditOutlined, DeleteOutlined, PlusOutlined } from '@ant-design/icons';
import { useQueryClient } from '@tanstack/react-query';
import { useTranslation } from 'react-i18next';
import { useNavigate } from 'react-router-dom';
import dayjs from 'dayjs';

import { PaginationQuery } from '@/core';
import { TABLE } from '@/config/constants';
import { useProTableScrollY } from '@/hooks/useProTableScrollY';
import { fetchListPages, useDeletePage } from '@/api/hooks/page';
import {
  getBoolLabel,
  getEditorTypeColor,
  getEditorTypeLabel,
  getPageStatusColor,
  getPageStatusLabel,
  getPageTypeColor,
  getPageTypeLabel,
  pageStatusOptions,
  pickTranslation,
} from './constants';

/**
 * 页面管理列表页（树形表格：子页面挂在 children 下，对齐 vben 侧 treeConfig）
 */
const PageList = () => {
  const { t, i18n } = useTranslation('page');
  const actionRef = useRef<ActionType>(null);
  const containerRef = useRef<HTMLDivElement>(null);
  const tableScrollY = useProTableScrollY(containerRef);
  const queryClient = useQueryClient();
  const navigate = useNavigate();
  const { message } = App.useApp();

  // 树形展开受控：数据异步加载后全展开
  const [expandedKeys, setExpandedKeys] = useState<readonly any[]>([]);

  const deleteMutation = useDeletePage({
    onSuccess: () => {
      message.success(t('deleteSuccess'));
      actionRef.current?.reload();
      queryClient.invalidateQueries({ queryKey: ['listPages'] });
    },
    onError: (error: Error) => {
      message.error(error.message || t('deleteFailed'));
    },
  });

  const columns: ProColumns<any>[] = [
    {
      title: t('slug'),
      dataIndex: 'slug',
      fixed: 'left',
      render: (_, record) => record.slug || pickTranslation(record, i18n.language)?.slug || '',
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
      title: t('type'),
      dataIndex: 'type',
      valueType: 'select',
      width: 110,
      fieldProps: {
        options: [
          'PAGE_TYPE_DEFAULT',
          'PAGE_TYPE_HOME',
          'PAGE_TYPE_ERROR_404',
          'PAGE_TYPE_ERROR_500',
          'PAGE_TYPE_CUSTOM',
        ].map((value) => ({ label: t(`typeMap.${value}`), value })),
        showSearch: true,
      },
      render: (_, record) => (
        <Tag color={getPageTypeColor(record.type)}>
          {getPageTypeLabel(t, record.type)}
        </Tag>
      ),
    },
    {
      title: t('showInNavigation'),
      dataIndex: 'showInNavigation',
      search: false,
      width: 130,
      render: (_, record) => getBoolLabel(t, record.showInNavigation),
    },
    {
      title: t('disallowComment'),
      dataIndex: 'disallowComment',
      search: false,
      width: 120,
      render: (_, record) => getBoolLabel(t, record.disallowComment),
    },
    {
      title: t('authorName'),
      dataIndex: 'authorName',
      search: false,
      ellipsis: true,
    },
    {
      title: t('status'),
      dataIndex: 'status',
      valueType: 'select',
      width: 100,
      fieldProps: {
        options: pageStatusOptions(t),
        showSearch: true,
      },
      render: (_, record) => (
        <Tag color={getPageStatusColor(record.status)}>
          {getPageStatusLabel(t, record.status)}
        </Tag>
      ),
    },
    {
      title: t('sortOrder'),
      dataIndex: 'sortOrder',
      search: false,
      width: 90,
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
            navigate(`/content/pages/edit/${record.id}?lang=${i18n.language}`)
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

  return (
    <div ref={containerRef} className="page-container-content" style={{ padding: '0 8px', height: '100%' }}>
      <ProTable<any>
        actionRef={actionRef}
        columns={columns}
        headerTitle={t('pageList')}
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
            });

            const response = await fetchListPages(query);

            // 树形展开受控：收集所有含 children 的父页面 id
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
        }}
        toolBarRender={() => [
          <TableExportButton
            key="export"
            fetcher={fetchListPages}
            columns={columns}
            filename="pages"
          />,
          <Button
            key="create"
            type="primary"
            icon={<PlusOutlined />}
            size="small"
            onClick={() => navigate(`/content/pages/create?lang=${i18n.language}`)}
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
        expandable={{
          expandedRowKeys: expandedKeys,
          onExpandedRowsChange: (keys: readonly any[]) => setExpandedKeys(keys),
        }}
      />
    </div>
  );
};

export default PageList;
