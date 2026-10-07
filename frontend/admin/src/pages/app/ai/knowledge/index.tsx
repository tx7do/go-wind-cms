import { useRef, useState } from 'react';
import type { ProColumns, ActionType } from '@ant-design/pro-components';
import ListTable from '@/components/common/ListTable';
import { Button, Popconfirm, Tag, App } from 'antd';
import { DeleteOutlined, EditOutlined, FileTextOutlined, PlusOutlined } from '@ant-design/icons';
import { useQueryClient } from '@tanstack/react-query';
import { useTranslation } from 'react-i18next';
import type { aiservicev1_AiKnowledgeBase as AiKnowledgeBase } from '@/api/generated/admin/service/v1';
import { PaginationQuery } from '@/core';
import { fetchListAiKnowledgeBases, useDeleteAiKnowledgeBase } from '@/api/hooks/ai-knowledge';
import { useProTableScrollY } from '@/hooks/useProTableScrollY';
import ContentContainer from '@/layouts/components/PageContainer/ContentContainer';
import { TABLE } from '@/config/constants';
import KnowledgeDrawer from './components/KnowledgeDrawer';
import DocsDrawer from './components/DocsDrawer';

export default function AiKnowledgePage() {
  const { t } = useTranslation('aiKnowledge');
  const { message } = App.useApp();
  const queryClient = useQueryClient();
  const actionRef = useRef<ActionType>(null);
  const containerRef = useRef<HTMLDivElement>(null);
  const tableScrollY = useProTableScrollY(containerRef);
  const [drawerOpen, setDrawerOpen] = useState(false);
  const [drawerMode, setDrawerMode] = useState<'create' | 'edit'>('create');
  const [selected, setSelected] = useState<AiKnowledgeBase | undefined>(undefined);
  const [docsTarget, setDocsTarget] = useState<AiKnowledgeBase | undefined>(undefined);
  const [docsOpen, setDocsOpen] = useState(false);

  const deleteMutation = useDeleteAiKnowledgeBase({
    onSuccess: () => {
      message.success(t('deleteSuccess'));
      actionRef.current?.reload();
      queryClient.invalidateQueries({ queryKey: ['listAiKnowledgeBases'] });
    },
    onError: (error: Error) => message.error(error.message || t('fetchFailed')),
  });

  const columns: ProColumns[] = [
    { title: t('name'), dataIndex: 'name', minWidth: 160 },
    { title: t('description'), dataIndex: 'description', ellipsis: true },
    { title: t('embeddingModel'), dataIndex: 'embeddingModel', width: 180 },
    {
      title: t('docCount'),
      dataIndex: 'docCount',
      width: 90,
      render: (_, record) => <Tag color="blue">{record.docCount ?? 0}</Tag>,
    },
    {
      title: t('action'),
      valueType: 'option',
      width: 220,
      fixed: 'right',
      render: (_, record) => [
        <Button
          key="docs"
          type="link"
          size="small"
          icon={<FileTextOutlined />}
          onClick={() => {
            setDocsTarget(record);
            setDocsOpen(true);
          }}
        >
          {t('docs')}
        </Button>,
        <Button
          key="edit"
          type="link"
          size="small"
          icon={<EditOutlined />}
          onClick={() => {
            setSelected(record);
            setDrawerMode('edit');
            setDrawerOpen(true);
          }}
        >
          {t('edit')}
        </Button>,
        <Popconfirm
          key="delete"
          title={t('deleteConfirmTitle')}
          description={t('deleteConfirmDesc')}
          onConfirm={() => deleteMutation.mutate({ id: record.id! })}
        >
          <Button type="link" size="small" danger icon={<DeleteOutlined />}>
            {t('deleteConfirmTitle')}
          </Button>
        </Popconfirm>,
      ],
    },
  ];

  return (
    <ContentContainer heightMode="fixed" padding="16px" bottomMargin={0}>
      <div ref={containerRef} className="page-container-content">
        <ListTable<AiKnowledgeBase>
          actionRef={actionRef}
          rowKey="id"
          search={false}
          scroll={{ y: tableScrollY, x: 900 }}
          request={async (params) => {
            const { current, pageSize } = params;
            const query = new PaginationQuery({
              paging: { page: current || 1, pageSize: pageSize || TABLE.DEFAULT_PAGE_SIZE },
            });
            const res = await fetchListAiKnowledgeBases(query);
            return { data: res.items || [], total: res.total || 0, success: true };
          }}
          toolBarRender={() => [
            <Button
              key="create"
              type="primary"
              icon={<PlusOutlined />}
              onClick={() => {
                setSelected(undefined);
                setDrawerMode('create');
                setDrawerOpen(true);
              }}
            >
              {t('create')}
            </Button>,
          ]}
          columns={columns}
        />
        <KnowledgeDrawer
          open={drawerOpen}
          mode={drawerMode}
          data={selected}
          onClose={() => setDrawerOpen(false)}
          onSuccess={() => {
            setDrawerOpen(false);
            actionRef.current?.reload();
            queryClient.invalidateQueries({ queryKey: ['listAiKnowledgeBases'] });
          }}
        />
        <DocsDrawer
          open={docsOpen}
          data={docsTarget}
          onClose={() => setDocsOpen(false)}
        />
      </div>
    </ContentContainer>
  );
}
