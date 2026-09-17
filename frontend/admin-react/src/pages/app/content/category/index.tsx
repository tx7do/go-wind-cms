import { useRef, useState } from 'react';
import type { ProColumns, ActionType } from '@ant-design/pro-components';
import TableExportButton from '@/components/common/TableExportButton';
import { ProTable } from '@ant-design/pro-components';
import { Button, Image, Popconfirm, Tag, App } from 'antd';
import { EditOutlined, DeleteOutlined, PlusOutlined } from '@ant-design/icons';
import { Icon } from '@iconify/react';
import { useQueryClient } from '@tanstack/react-query';
import { useTranslation } from 'react-i18next';
import { useNavigate } from 'react-router-dom';
import dayjs from 'dayjs';

import { PaginationQuery } from '@/core';
import { TABLE } from '@/config/constants';
import { useProTableScrollY } from '@/hooks/useProTableScrollY';
import { fetchListCategories, useDeleteCategory } from '@/api/hooks/category';
import {
  getCategoryStatusColor,
  getCategoryStatusLabel,
  getBoolLabel,
  pickTranslation,
  categoryStatusOptions,
} from './constants';

// 列表需要的字段掩码（对齐 vben 侧，保证树形 children 返回）
const LIST_FIELD_MASK =
  'id,status,sort_order,is_nav,icon,code,thumbnail,post_count,direct_post_count,available_languages,parent_id,children,created_by,created_at,translations,translations.id,translations.language_code,translations.name,translations.slug,translations.description,translations.cover_image';

/**
 * 分类管理列表页（树形表格：子分类挂在 children 下）
 */
const CategoryList = () => {
  const { t, i18n } = useTranslation('category');
  const actionRef = useRef<ActionType>(null);
  const containerRef = useRef<HTMLDivElement>(null);
  const tableScrollY = useProTableScrollY(containerRef);
  const queryClient = useQueryClient();
  const navigate = useNavigate();
  const { message } = App.useApp();

  // 树形展开受控：数据异步加载后全展开
  const [expandedKeys, setExpandedKeys] = useState<readonly any[]>([]);

  const deleteMutation = useDeleteCategory({
    onSuccess: () => {
      message.success(t('deleteSuccess'));
      actionRef.current?.reload();
      queryClient.invalidateQueries({ queryKey: ['listCategories'] });
    },
    onError: (error: Error) => {
      message.error(error.message || t('deleteFailed'));
    },
  });

  const columns: ProColumns<any>[] = [
    {
      title: t('name'),
      dataIndex: 'translations',
      fixed: 'left',
      render: (_, record) => pickTranslation(record, i18n.language)?.name || '',
    },
    {
      title: t('icon'),
      dataIndex: 'icon',
      search: false,
      width: 70,
      render: (_, record) =>
        record.icon ? <Icon icon={record.icon} width="1.2em" height="1.2em" /> : '',
    },
    {
      title: t('thumbnail'),
      dataIndex: 'thumbnail',
      search: false,
      width: 80,
      render: (_, record) =>
        record.thumbnail ? <Image src={record.thumbnail} width={50} /> : '',
    },
    {
      title: t('coverImage'),
      dataIndex: 'coverImage',
      search: false,
      width: 80,
      render: (_, record) => {
        const cover = pickTranslation(record, i18n.language)?.coverImage;
        return cover ? <Image src={cover} width={50} /> : '';
      },
    },
    {
      title: t('description'),
      dataIndex: 'description',
      search: false,
      ellipsis: true,
      render: (_, record) => pickTranslation(record, i18n.language)?.description || '',
    },
    {
      title: t('isNav'),
      dataIndex: 'isNav',
      search: false,
      width: 90,
      render: (_, record) => getBoolLabel(t, record.isNav),
    },
    {
      title: t('postCount'),
      dataIndex: 'postCount',
      search: false,
      width: 100,
    },
    {
      title: t('directPostCount'),
      dataIndex: 'directPostCount',
      search: false,
      width: 120,
    },
    {
      title: t('status'),
      dataIndex: 'status',
      valueType: 'select',
      width: 100,
      fieldProps: {
        options: categoryStatusOptions(t),
        showSearch: true,
      },
      render: (_, record) => (
        <Tag color={getCategoryStatusColor(record.status)}>
          {getCategoryStatusLabel(t, record.status)}
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
            navigate(`/content/categories/edit/${record.id}?lang=${i18n.language}`)
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
        headerTitle={t('categoryList')}
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
              fieldMask: LIST_FIELD_MASK,
            });

            const response = await fetchListCategories(query);

            // 树形展开受控：收集所有含 children 的父分类 id
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
            fetcher={fetchListCategories}
            columns={columns}
            filename="categories"
          />,
          <Button
            key="create"
            type="primary"
            icon={<PlusOutlined />}
            size="small"
            onClick={() => navigate(`/content/categories/create?lang=${i18n.language}`)}
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
        scroll={{ x: 1200, y: tableScrollY }}
        expandable={{
          expandedRowKeys: expandedKeys,
          onExpandedRowsChange: (keys: readonly any[]) => setExpandedKeys(keys),
        }}
      />
    </div>
  );
};

export default CategoryList;
