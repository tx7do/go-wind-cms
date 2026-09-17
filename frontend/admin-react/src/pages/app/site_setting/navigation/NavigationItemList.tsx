import { useEffect, useRef, useState } from 'react';
import type { ProColumns, ActionType } from '@ant-design/pro-components';
import TableExportButton from '@/components/common/TableExportButton';
import { ProTable } from '@ant-design/pro-components';
import { Button, Popconfirm, Tag, App, Empty } from 'antd';
import { EditOutlined, DeleteOutlined, PlusOutlined } from '@ant-design/icons';
import { Icon } from '@iconify/react';
import { useQueryClient } from '@tanstack/react-query';
import { useTranslation } from 'react-i18next';
import dayjs from 'dayjs';

import { PaginationQuery } from '@/core';
import { TABLE } from '@/config/constants';
import { useProTableScrollY } from '@/hooks/useProTableScrollY';
import {
  fetchListNavigationItems,
  useDeleteNavigationItem,
} from '@/api/hooks/navigation-item';
import {
  getBoolLabel,
  getLinkTypeColor,
  getLinkTypeLabel,
  linkTypeOptions,
} from './constants';
import NavigationItemDrawer from './NavigationItemDrawer';

interface NavigationItemListProps {
  navigationId: number | null;
}

// 图标缺前缀时补 carbon:（对齐 vben 侧约定）
function getIconName(icon: string): string {
  if (!icon) return '';
  if (icon.includes(':')) return icon;
  return `carbon:${icon}`;
}

/**
 * 导航项列表（从表）：随左侧选中导航联动
 */
const NavigationItemList: React.FC<NavigationItemListProps> = ({
  navigationId,
}) => {
  const { t } = useTranslation('navigation-item');
  const actionRef = useRef<ActionType>(null);
  const containerRef = useRef<HTMLDivElement>(null);
  const tableScrollY = useProTableScrollY(containerRef);
  const queryClient = useQueryClient();
  const { message } = App.useApp();

  const [drawerOpen, setDrawerOpen] = useState(false);
  const [drawerMode, setDrawerMode] = useState<'create' | 'edit'>('create');
  const [editingRow, setEditingRow] = useState<any>(null);

  const deleteMutation = useDeleteNavigationItem({
    onSuccess: () => {
      message.success(t('deleteSuccess'));
      actionRef.current?.reload();
      queryClient.invalidateQueries({ queryKey: ['listNavigationItems'] });
    },
    onError: (error: Error) => {
      message.error(error.message || t('deleteFailed'));
    },
  });

  // 当选中导航变化时自动刷新列表
  useEffect(() => {
    actionRef.current?.reload();
  }, [navigationId]);

  const columns: ProColumns<any>[] = [
    {
      title: t('title'),
      dataIndex: 'title',
      fixed: 'left',
    },
    {
      title: t('description'),
      dataIndex: 'description',
      ellipsis: true,
    },
    {
      title: t('icon'),
      dataIndex: 'icon',
      width: 80,
      render: (_, record) =>
        record.icon ? (
          <Icon icon={getIconName(record.icon)} width="1.2em" height="1.2em" />
        ) : (
          ''
        ),
    },
    {
      title: t('url'),
      dataIndex: 'url',
      ellipsis: true,
    },
    {
      title: t('linkType'),
      dataIndex: 'linkType',
      valueType: 'select',
      fieldProps: {
        options: linkTypeOptions(t),
      },
      render: (_, record) => (
        <Tag color={getLinkTypeColor(record.linkType)}>
          {getLinkTypeLabel(t, record.linkType)}
        </Tag>
      ),
    },
    {
      title: t('isOpenNewTab'),
      dataIndex: 'isOpenNewTab',
      width: 90,
      render: (_, record) => getBoolLabel(t, record.isOpenNewTab),
    },
    {
      title: t('isInvalid'),
      dataIndex: 'isInvalid',
      width: 90,
      render: (_, record) => getBoolLabel(t, record.isInvalid),
    },
    {
      title: t('createdAt'),
      dataIndex: 'createdAt',
      width: 150,
      render: (_, record) =>
        record.createdAt ? dayjs(record.createdAt).format('YYYY-MM-DD HH:mm') : '',
    },
    {
      title: t('action'),
      valueType: 'option',
      width: 90,
      render: (_, record) => [
        <a
          key="edit"
          onClick={() => {
            setEditingRow(record);
            setDrawerMode('edit');
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
      <div
        ref={containerRef}
        className="page-container-content"
        style={{ padding: '0 4px', height: '100%', display: 'flex', flexDirection: 'column' }}
      >
        {navigationId ? (
          <ProTable<any>
            actionRef={actionRef}
            columns={columns}
            headerTitle={t('itemList')}
            params={{ navigationId }}
            request={async (params) => {
              try {
                const query = new PaginationQuery({
                  paging: {
                    page: params.current || 1,
                    pageSize: params.pageSize || TABLE.DEFAULT_PAGE_SIZE,
                  },
                  formValues: Object.fromEntries(
                    Object.entries(params).filter(
                      ([key]) =>
                        !['current', 'pageSize', 'navigationId'].includes(key),
                    ),
                  ),
                  orderBy: ['-created_at'],
                });

                const response = await fetchListNavigationItems(query);

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
                fetcher={fetchListNavigationItems}
                columns={columns}
                filename="navigation-items"
              />,
              <Button
                key="create"
                type="primary"
                icon={<PlusOutlined />}
                size="small"
                onClick={() => {
                  setEditingRow(null);
                  setDrawerMode('create');
                  setDrawerOpen(true);
                }}
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
        ) : (
          <div
            style={{
              flex: 1,
              display: 'flex',
              alignItems: 'center',
              justifyContent: 'center',
            }}
          >
            <Empty description={t('selectNavigationFirst')} />
          </div>
        )}
      </div>

      <NavigationItemDrawer
        open={drawerOpen}
        mode={drawerMode}
        data={editingRow}
        navigationId={navigationId}
        onClose={() => {
          setDrawerOpen(false);
          setEditingRow(null);
        }}
        onSuccess={() => {
          actionRef.current?.reload();
        }}
      />
    </>
  );
};

export default NavigationItemList;
