import { useCallback, useEffect, useMemo, useState } from 'react';
import { App, Button, Calendar, List, Modal, Spin, Tag } from 'antd';
import { ReloadOutlined } from '@ant-design/icons';
import dayjs, { type Dayjs } from 'dayjs';
import { useTranslation } from 'react-i18next';
import { useNavigate } from 'react-router-dom';
import ContentContainer from '@/layouts/components/PageContainer/ContentContainer';
import { PaginationQuery } from '@/core';
import { fetchListPosts } from '@/api/hooks/post';

/**
 * 内容日历（排期视图）：按发布时间把「定时排期」与「已发布」文章铺到月历上。
 * 数据走既有 ListPosts（状态过滤 + 大页拉取后前端按日分桶）——排期量级小，
 * 不加后端查询；点击条目跳文章编辑页。
 */

type PostRow = {
  id?: number;
  status?: string;
  publishTime?: string;
  translations?: { title?: string; languageCode?: string }[];
};

const STATUS_COLOR: Record<string, string> = {
  POST_STATUS_SCHEDULED: 'purple',
  POST_STATUS_PUBLISHED: 'green',
};

function statusKey(status?: string): string {
  // 兼容带/不带前缀的枚举形态
  if (!status) return '';
  return status.startsWith('POST_STATUS_') ? status : `POST_STATUS_${status}`;
}

function postTitle(p: PostRow): string {
  for (const tr of p.translations ?? []) {
    if (tr.title) return tr.title;
  }
  return `#${p.id ?? ''}`;
}

export default function ContentCalendarPage() {
  const { t } = useTranslation('content-calendar');
  const { message } = App.useApp();
  const navigate = useNavigate();

  const [loading, setLoading] = useState(false);
  const [posts, setPosts] = useState<PostRow[]>([]);
  const [dayOpen, setDayOpen] = useState<Dayjs | undefined>();

  const loadPosts = useCallback(async () => {
    setLoading(true);
    try {
      const scheduled = await fetchListPosts(new PaginationQuery({
        paging: { page: 1, pageSize: 200 },
        formValues: { status: 'POST_STATUS_SCHEDULED' },
        orderBy: ['publishTime'],
      }));
      const published = await fetchListPosts(new PaginationQuery({
        paging: { page: 1, pageSize: 100 },
        formValues: { status: 'POST_STATUS_PUBLISHED' },
        orderBy: ['-publishTime'],
      }));
      setPosts([...(scheduled.items ?? []), ...(published.items ?? [])] as PostRow[]);
    } catch (error) {
      console.error('Load calendar posts failed:', error);
      message.error(t('loadFailed'));
    } finally {
      setLoading(false);
    }
  }, [message, t]);

  useEffect(() => {
    // 经宏任务拉起：loadPosts 首行同步 setLoading 会被 set-state-in-effect 规则命中
    const timer = setTimeout(() => void loadPosts(), 0);
    return () => clearTimeout(timer);
  }, [loadPosts]);

  const byDay = useMemo(() => {
    const map = new Map<string, PostRow[]>();
    for (const p of posts) {
      if (!p.publishTime) continue;
      const key = dayjs(p.publishTime).format('YYYY-MM-DD');
      const list = map.get(key) ?? [];
      list.push(p);
      map.set(key, list);
    }
    return map;
  }, [posts]);

  const dayPosts = dayOpen ? byDay.get(dayOpen.format('YYYY-MM-DD')) ?? [] : [];

  const statusTag = (status?: string) => (
    <Tag color={STATUS_COLOR[statusKey(status)] ?? 'default'}>
      {t(statusKey(status) === 'POST_STATUS_SCHEDULED' ? 'scheduled' : 'published')}
    </Tag>
  );

  const cellRender = (value: Dayjs) => {
    const list = byDay.get(value.format('YYYY-MM-DD')) ?? [];
    if (list.length === 0) return null;
    return (
      <ul style={{ listStyle: 'none', margin: 0, padding: 0 }}>
        {list.slice(0, 3).map((p) => (
          <li key={p.id} style={{ marginBottom: 2 }}>
            {statusTag(p.status)}
            <span style={{ fontSize: 12 }}>{postTitle(p)}</span>
          </li>
        ))}
        {list.length > 3 && (
          <li style={{ fontSize: 12, color: 'var(--ant-color-text-secondary)' }}>
            +{list.length - 3}
          </li>
        )}
      </ul>
    );
  };

  return (
    <ContentContainer>
      <div style={{ marginBottom: 8 }}>
        <Button icon={<ReloadOutlined />} loading={loading} onClick={() => void loadPosts()}>
          {t('reload')}
        </Button>
      </div>
      <Spin spinning={loading}>
        <Calendar
          onSelect={(value) => {
            if ((byDay.get(value.format('YYYY-MM-DD')) ?? []).length > 0) {
              setDayOpen(value);
            }
          }}
          cellRender={(value, info) =>
            info.type === 'date' ? cellRender(value) : info.originNode
          }
        />
      </Spin>

      <Modal
        title={dayOpen ? dayOpen.format('YYYY-MM-DD') : ''}
        open={!!dayOpen}
        onCancel={() => setDayOpen(undefined)}
        footer={null}
        width={520}
        destroyOnHidden
      >
        <List
          dataSource={dayPosts}
          renderItem={(p) => (
            <List.Item
              style={{ cursor: 'pointer' }}
              onClick={() => p.id && navigate(`/content/posts/edit/${p.id}`)}
            >
              <List.Item.Meta
                title={postTitle(p)}
                description={statusTag(p.status)}
              />
              <span style={{ fontSize: 12, color: 'var(--ant-color-text-secondary)' }}>
                {p.publishTime ? dayjs(p.publishTime).format('HH:mm') : ''}
              </span>
            </List.Item>
          )}
        />
      </Modal>
    </ContentContainer>
  );
}
