import { useRef } from 'react';
import type { ProColumns, ActionType } from '@ant-design/pro-components';
import TableExportButton from '@/components/common/TableExportButton';
import { ProTable } from '@ant-design/pro-components';
import { Button, Popconfirm, Tag, App } from 'antd';
import { EditOutlined, DeleteOutlined, PlusOutlined } from '@ant-design/icons';
import { Icon } from '@iconify/react';
import { useQueryClient } from '@tanstack/react-query';
import { useTranslation } from 'react-i18next';
import { useNavigate } from 'react-router-dom';
import dayjs from 'dayjs';

import { PaginationQuery } from '@/core';
import { TABLE } from '@/config/constants';
import { useProTableScrollY } from '@/hooks/useProTableScrollY';
import { useDeleteTag, fetchListTags } from '@/api/hooks/tag';
import {
  getBoolLabel,
  getIconName,
  getTagStatusColor,
  getTagStatusLabel,
  pickTranslation,
  tagStatusOptions,
} from './constants';

/**
 * 标签管理列表页
 */
const TagList = () => {
  const { t, i18n } = useTranslation('tag');
  const actionRef = useRef<ActionType>(null);
  const containerRef = useRef<HTMLDivElement>(null);
  const tableScrollY = useProTableScrollY(containerRef);
  const queryClient = useQueryClient();
  const navigate = useNavigate();
  const { message } = App.useApp();

  const deleteMutation = useDeleteTag({
    onSuccess: () => {
      message.success(t('deleteSuccess'));
      actionRef.current?.reload();
      queryClient.invalidateQueries({ queryKey: ['listTags'] });
    },
    onError: (error: Error) => {
      message.error(error.message || t('deleteFailed'));
    },
  });

  const handleCreate = () => {
    navigate(`/content/tags/create?lang=${i18n.language}`);
  };

  const handleEdit = (row: any) => {
    navigate(`/content/tags/edit/${row.id}?lang=${i18n.language}`);
  };

  const columns: ProColumns<any>[] = [
    {
      title: t('name'),
      dataIndex: 'translations.name',
      fixed: 'left',
      render: (_, record) => pickTranslation(record, i18n.language)?.name || '',
    },
    {
      title: t('description'),
      dataIndex: 'translations.description',
      ellipsis: true,
      render: (_, record) =>
        pickTranslation(record, i18n.language)?.description || '',
    },
    {
      title: t('color'),
      dataIndex: 'color',
      search: false,
      width: 110,
      render: (_, record) =>
        record.color ? <Tag color={record.color}>{record.color}</Tag> : '',
    },
    {
      title: t('icon'),
      dataIndex: 'icon',
      search: false,
      width: 70,
      render: (_, record) =>
        record.icon ? <Icon icon={getIconName(record.icon)} width="1.2em" height="1.2em" /> : '',
    },
    {
      title: t('group'),
      dataIndex: 'group',
      search: false,
    },
    {
      title: t('isFeatured'),
      dataIndex: 'isFeatured',
      search: false,
      width: 90,
      render: (_, record) => getBoolLabel(t, record.isFeatured),
    },
    {
      title: t('postCount'),
      dataIndex: 'postCount',
      search: false,
      width: 100,
    },
    {
      title: t('status'),
      dataIndex: 'status',
      valueType: 'select',
      width: 100,
      fieldProps: {
        options: tagStatusOptions(t),
        showSearch: true,
      },
      render: (_, record) => (
        <Tag color={getTagStatusColor(record.status)}>
          {getTagStatusLabel(t, record.status)}
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
        <a key="edit" onClick={() => handleEdit(record)}>
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
        headerTitle={t('tagList')}
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

            const response = await fetchListTags(query);

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
            fetcher={fetchListTags}
            columns={columns}
            filename="tags"
          />,
          <Button
            key="create"
            type="primary"
            icon={<PlusOutlined />}
            size="small"
            onClick={handleCreate}
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
        scroll={{ y: tableScrollY }}
      />
    </div>
  );
};

export default TagList;
