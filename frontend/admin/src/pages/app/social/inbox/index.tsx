import { useMemo, useState } from 'react';
import { Badge, Card, List, Popconfirm, Segmented, Tag } from 'antd';
import { useQueryClient } from '@tanstack/react-query';
import { useTranslation } from 'react-i18next';
import dayjs from 'dayjs';

import { useListSocialInbox, useMarkSocialInboxRead } from '@/api/hooks/social';
import { getPlatformColor, getPlatformLabel } from '../account/constants';

/**
 * 社交收件箱：回拉平台侧的回复/提及（Bluesky/Mastodon）。
 * 周期任务每 15 分钟自动回拉；条目可标记已读。
 */
const SocialInboxList = () => {
  const { t } = useTranslation('social');
  const queryClient = useQueryClient();
  const [filter, setFilter] = useState<'all' | 'unread'>('all');

  const { data, isLoading } = useListSocialInbox({
    unreadOnly: filter === 'unread',
    limit: 100,
  });

  const invalidate = () => {
    queryClient.invalidateQueries({ queryKey: ['listSocialInbox'] });
  };

  const markRead = useMarkSocialInboxRead({
    onSuccess: () => invalidate(),
  });

  const items = useMemo(() => data?.items ?? [], [data]);
  const unreadCount = Number(data?.unreadCount) || 0;

  return (
    <Card bordered={false}>
      <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: 12 }}>
        <Segmented
          value={filter}
          onChange={(v) => setFilter(v as 'all' | 'unread')}
          options={[
            { label: t('inboxAll'), value: 'all' },
            { label: t('inboxUnread'), value: 'unread' },
          ]}
        />
        <Badge count={unreadCount} overflowCount={99} />
      </div>

      <List
        loading={isLoading}
        dataSource={items}
        locale={{ emptyText: t('inboxEmpty') }}
        renderItem={(item: any) => (
          <List.Item
            actions={[
              !item.readAt && (
                <Popconfirm
                  key="read"
                  title={t('markReadConfirmTitle')}
                  onConfirm={() => markRead.mutate({ id: item.id })}
                  okText={t('common:button.ok')}
                  cancelText={t('common:button.cancel')}
                >
                  <a>{t('markRead')}</a>
                </Popconfirm>
              ),
              item.url && (
                <a key="link" href={item.url} target="_blank" rel="noreferrer">
                  {t('view')}
                </a>
              ),
            ].filter(Boolean)}
          >
            <List.Item.Meta
              title={
                <span style={{ display: 'inline-flex', gap: 8, alignItems: 'center' }}>
                  {!item.readAt && <Badge status="processing" />}
                  <Tag color={getPlatformColor(item.platform)}>
                    {getPlatformLabel(t, item.platform)}
                  </Tag>
                  <span>{item.type}</span>
                  <span style={{ color: 'var(--ant-color-text-tertiary)', fontSize: 12 }}>
                    {item.authorName || item.authorHandle || '—'}
                  </span>
                </span>
              }
              description={
                <>
                  <div style={{ fontSize: 13 }}>{item.content || '—'}</div>
                  <div style={{ fontSize: 12, color: 'var(--ant-color-text-tertiary)', marginTop: 2 }}>
                    {item.occurredAt ? dayjs(item.occurredAt).format('YYYY-MM-DD HH:mm') : ''}
                  </div>
                </>
              }
            />
          </List.Item>
        )}
      />
    </Card>
  );
};

export default SocialInboxList;
