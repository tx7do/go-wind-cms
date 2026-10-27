import { App, Button, Form, Input } from 'antd';
import { DeleteOutlined, PlusOutlined, ReloadOutlined } from '@ant-design/icons';
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import { useTranslation } from 'react-i18next';
import ListTable from '@/components/common/ListTable';
import ContentContainer from '@/layouts/components/PageContainer/ContentContainer';
import { apiClient } from '@/api/client';

/**
 * 搜索同义词管理页（租户级词库）：检索时查询侧按触发词扩展（OR 组），
 * 免索引重建、保存即时生效。
 */
export default function SearchSynonymPage() {
  const { t } = useTranslation('search-synonym');
  const { message } = App.useApp();
  const queryClient = useQueryClient();
  const [form] = Form.useForm();

  const listQuery = useQuery({
    queryKey: ['searchSynonyms'],
    queryFn: () => apiClient.searchSynonymService.List({}),
  });

  const refresh = () => queryClient.invalidateQueries({ queryKey: ['searchSynonyms'] });

  const createMutation = useMutation({
    mutationFn: (values: { term: string; synonymsText: string }) => {
      const synonyms = values.synonymsText
        .split(/[,，\n]+/)
        .map((x) => x.trim())
        .filter(Boolean);
      return apiClient.searchSynonymService.Create({
        data: { term: values.term.trim(), synonyms, enabled: true },
      });
    },
    onSuccess: () => {
      message.success(t('created'));
      form.resetFields();
      refresh();
    },
    onError: (error: Error) => {
      console.error('create search synonym failed:', error);
      message.error(error?.message || t('saveFailed'));
    },
  });

  const deleteMutation = useMutation({
    mutationFn: (id: number) => apiClient.searchSynonymService.Delete({ id }),
    onSuccess: () => {
      message.success(t('deleted'));
      refresh();
    },
    onError: (error: Error) => {
      console.error('delete search synonym failed:', error);
      message.error(error?.message || t('deleteFailed'));
    },
  });

  const handleSubmit = async () => {
    const values = await form.validateFields().catch(() => null);
    if (values === null) return;
    await createMutation.mutateAsync(values);
  };

  const columns = [
    { title: 'ID', dataIndex: 'id', width: 70 },
    {
      title: t('term'),
      dataIndex: 'term',
      width: 180,
    },
    {
      title: t('synonyms'),
      dataIndex: 'synonyms',
      ellipsis: true,
      render: (_: unknown, record: { synonyms?: string[] }) =>
        (record.synonyms ?? []).join(' / ') || '-',
    },
    {
      title: t('action'),
      width: 90,
      render: (_: unknown, record: { id?: number }) => (
        <Button
          size="small"
          danger
          icon={<DeleteOutlined />}
          onClick={() => record.id && deleteMutation.mutate(record.id)}
        />
      ),
    },
  ];

  return (
    <ContentContainer>
      <div className="mb-4 rounded-xl border border-solid border-gray-200 bg-white p-4 dark:border-gray-700 dark:bg-gray-900">
        <Form form={form} layout="inline" onFinish={handleSubmit}>
          <Form.Item
            name="term"
            rules={[{ required: true, message: t('termRequired') }]}
            style={{ marginBottom: 8 }}
          >
            <Input placeholder={t('termPlaceholder')} style={{ width: 200 }} />
          </Form.Item>
          <Form.Item
            name="synonymsText"
            rules={[{ required: true, message: t('synonymsRequired') }]}
            style={{ marginBottom: 8, flex: 1, minWidth: 260 }}
          >
            <Input placeholder={t('synonymsPlaceholder')} />
          </Form.Item>
          <Form.Item style={{ marginBottom: 8 }}>
            <Button type="primary" htmlType="submit" icon={<PlusOutlined />} loading={createMutation.isPending}>
              {t('add')}
            </Button>
          </Form.Item>
        </Form>
        <div className="text-xs text-[color:var(--ant-color-text-tertiary)]">{t('hint')}</div>
      </div>

      <ListTable
        headerTitle={t('title')}
        rowKey="id"
        search={false}
        loading={listQuery.isPending}
        columns={columns}
        request={async () => {
          const items = listQuery.data?.items ?? [];
          return { data: items, total: items.length, success: true };
        }}
        dataSource={listQuery.data?.items ?? []}
        toolBarRender={() => [
          <Button key="reload" icon={<ReloadOutlined />} onClick={() => listQuery.refetch()}>
            {t('reload')}
          </Button>,
        ]}
        pagination={{ pageSize: 15 }}
      />
    </ContentContainer>
  );
}
