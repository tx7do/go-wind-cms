import { useRef, useState } from 'react';
import type { ProColumns, ActionType } from '@ant-design/pro-components';
import TableExportButton from '@/components/common/TableExportButton';
import { ProTable } from '@ant-design/pro-components';
import { Button, Popconfirm, Tag, App } from 'antd';
import { EditOutlined, DeleteOutlined, PlusOutlined } from '@ant-design/icons';
import { useQueryClient } from '@tanstack/react-query';
import { useTranslation } from 'react-i18next';
import type { internal_messageservicev1_InternalMessage as InternalMessage } from '@/api/generated/admin/service/v1';
import { PaginationQuery } from '@/core';
import {
  fetchListInternalMessages,
  useDeleteInternalMessage,
} from '@/api/hooks/internal-message';
import { useProTableScrollY } from '@/hooks/useProTableScrollY';
import ContentContainer from '@/layouts/components/PageContainer/ContentContainer';
import { getStatusMap, getStatusOptions, getTypeMap, getTypeOptions } from './constants';
import MessageDrawer from './components/MessageDrawer';

/**
 * 内部消息列表页面
 */
const InternalMessageList = () => {
  const { t } = useTranslation('internal-message');
  const actionRef = useRef<ActionType>(null);
  const queryClient = useQueryClient();
  const { message } = App.useApp();

  const containerRef = useRef<HTMLDivElement>(null);
  const tableScrollY = useProTableScrollY(containerRef);

  // Drawer 状态管理
  const [drawerOpen, setDrawerOpen] = useState(false);
  const [drawerMode, setDrawerMode] = useState<'create' | 'edit'>('create');
  const [selectedMessage, setSelectedMessage] = useState<InternalMessage | undefined>();

  const statusMap = getStatusMap(t);
  const typeMap = getTypeMap(t);

  // 删除操作
  const deleteMutation = useDeleteInternalMessage({
    onSuccess: () => {
      message.success(t('deleteSuccess'));
      actionRef.current?.reload();
      queryClient.invalidateQueries({ queryKey: ['listInternalMessages'] });
    },
    onError: (error: Error) => {
      message.error(error.message || t('deleteFailed'));
    },
  });

  // 列配置
  const columns: ProColumns<InternalMessage>[] = [
    {
      title: t('title'),
      dataIndex: 'title',
    },
    {
      title: t('categoryName'),
      dataIndex: 'categoryName',
      hideInSearch: true,
    },
    {
      title: t('status'),
      dataIndex: 'status',
      valueType: 'select',
      fieldProps: {
        options: getStatusOptions(t),
      },
      render: (_, record) => {
        const status = record.status as keyof typeof statusMap;
        const config = statusMap[status] || { text: status, color: 'default' };
        return <Tag color={config.color}>{config.text}</Tag>;
      },
    },
    {
      title: t('type'),
      dataIndex: 'type',
      valueType: 'select',
      fieldProps: {
        options: getTypeOptions(t),
      },
      render: (_, record) => {
        const type = (record as any).type as keyof typeof typeMap;
        const config = typeMap[type] || { text: type, color: 'default' };
        return <Tag color={config.color}>{config.text}</Tag>;
      },
    },
    {
      title: t('senderName'),
      dataIndex: 'senderName',
      hideInSearch: true,
    },
    {
      title: t('createdAt'),
      dataIndex: 'createdAt',
      width: 180,
      valueType: 'dateTime',
      hideInSearch: true,
    },
    {
      title: t('action'),
      valueType: 'option',
      width: 90,
      render: (_, record) => [
        <a
          key="edit"
          onClick={() => {
            setDrawerMode('edit');
            setSelectedMessage(record);
            setDrawerOpen(true);
          }}
        >
          <EditOutlined />
        </a>,
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
          <ProTable<InternalMessage>
            actionRef={actionRef}
            columns={columns}
            request={async (params, _sorter, _filter) => {
              try {
                const formValues: Record<string, any> = {};
                // 过滤搜索参数
                Object.entries(params).forEach(([key, value]) => {
                  if (!['current', 'pageSize'].includes(key) && value !== undefined) {
                    formValues[key] = value;
                  }
                });

                const query = new PaginationQuery({
                  formValues,
                });

                const response = await fetchListInternalMessages(query);

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
              <TableExportButton key="export" fetcher={fetchListInternalMessages} columns={columns} filename="im-messages" />,
              <Button
                key="create"
                type="primary"
                icon={<PlusOutlined />}
                onClick={() => {
                  setDrawerMode('create');
                  setSelectedMessage(undefined);
                  setDrawerOpen(true);
                }}
              >
                {t('create')}
              </Button>,
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

      {/* 消息编辑/创建 Drawer */}
      <MessageDrawer
        open={drawerOpen}
        mode={drawerMode}
        data={selectedMessage}
        onClose={() => {
          setDrawerOpen(false);
          setSelectedMessage(undefined);
        }}
        onSuccess={() => {
          actionRef.current?.reload();
        }}
      />
    </>
  );
};

export default InternalMessageList;
