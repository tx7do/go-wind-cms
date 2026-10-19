import { useState } from 'react';
import { Card, Descriptions, Spin, Tag, Typography, Empty, Table } from 'antd';
import { useTranslation } from 'react-i18next';
import type {
  redis_cacheservicev1_RedisCacheMonitorInfo,
  redis_cacheservicev1_SlowLogEntry,
} from '@/api/generated/admin/service/v1';
import { useRedisCacheMonitorInfo } from '@/api/hooks/monitor-dashboards';
import { formatDateTime } from '@/utils/date';
import ContentContainer from '@/layouts/components/PageContainer/ContentContainer';

const { Text, Title } = Typography;

/**
 * Redis 缓存监控页面（只读）
 * 展示 Redis INFO（按 section 拆分的键值对）、DBSIZE、SLOWLOG。
 * 数据通过 GET /admin/v1/redis-cache-monitor 拉取，无任何写操作。
 */
const RedisCacheMonitor = () => {
  const { t } = useTranslation('redis-cache-monitor');

  const { data: info, isLoading } = useRedisCacheMonitorInfo();
  const [expandedSections, setExpandedSections] = useState<Record<string, boolean>>({});

  if (isLoading) {
    return (
      <ContentContainer heightMode="auto" scrollable padding="16px">
        <div style={{ display: 'flex', justifyContent: 'center', padding: 48 }}>
          <Spin size="large" />
        </div>
      </ContentContainer>
    );
  }

  const infoTyped = info as redis_cacheservicev1_RedisCacheMonitorInfo | undefined;
  const sections = infoTyped?.sections ?? [];
  const dbSize = infoTyped?.dbSize ?? 0;
  const slowlog = infoTyped?.slowlog ?? [];

  const slowLogColumns = [
    { title: t('slowlogId'), dataIndex: 'id', width: 80 },
    { title: t('slowlogTime'), dataIndex: 'createdAt', width: 170, render: (v?: string) => (v ? formatDateTime(v) : '-') },
    { title: t('slowlogDuration'), dataIndex: 'durationUsec', width: 120, render: (v?: number) => (v !== undefined ? `${v} µs` : '-') },
    { title: t('slowlogArgs'), dataIndex: 'args', ellipsis: true, render: (_: unknown, r: redis_cacheservicev1_SlowLogEntry) => (r.args ?? []).join(' ') },
    { title: t('slowlogClient'), dataIndex: 'clientAddr', width: 150, render: (v?: string) => v || '-' },
  ];

  return (
    // 长卡片流页面：fixed+scrollable 成对——十几个 INFO 卡片+慢日志表超出
    // 视口时必须能滚，默认 fixed 是 overflow:hidden 会直接裁掉内容。
    <ContentContainer heightMode="fixed" scrollable padding="16px">
      {/* 单值指标：DBSIZE */}
      <Card size="small" style={{ marginBottom: 16 }} title={<Text strong>{t('dbSize')}</Text>}>
        <Descriptions column={1} size="small" bordered colon>
          <Descriptions.Item label={t('dbSize')}>
            <Tag color={dbSize > 0 ? 'orange' : 'default'}>{String(dbSize)}</Tag>
          </Descriptions.Item>
        </Descriptions>
      </Card>

      {/* INFO sections：每个 section 一张可折叠卡片 */}
      {sections.length === 0 ? (
        <Card size="small" style={{ marginBottom: 16 }}>
          <Empty description={t('noInfo')} />
        </Card>
      ) : (
        sections.map((section, idx) => {
          const sectionName = section?.name ?? `${t('section')} ${idx}`;
          const key = `${sectionName}-${idx}`;
          const expanded = expandedSections[key] ?? false;
          return (
            <Card
              key={key}
              size="small"
              style={{ marginBottom: 16 }}
              title={
                <div
                  style={{ display: 'flex', alignItems: 'center', gap: 8, cursor: 'pointer' }}
                  onClick={() =>
                    setExpandedSections((prev) => ({ ...prev, [key]: !prev[key] }))
                  }
                >
                  <Text strong>{sectionName}</Text>
                  <Tag>{(section?.entries ?? []).length}</Tag>
                  <Text type="secondary" style={{ fontSize: 12 }}>
                    {expanded ? t('collapse') : t('expand')}
                  </Text>
                </div>
              }
            >
              {expanded ? (
                (section?.entries ?? []).length === 0 ? (
                  <Empty description={t('noEntries')} />
                ) : (
                  <Descriptions
                    column={1}
                    size="small"
                    bordered
                    colon
                    contentStyle={{ wordBreak: 'break-all' }}
                    labelStyle={{ width: 240 }}
                  >
                    {(section?.entries ?? []).map((entry, eIdx) => (
                      <Descriptions.Item key={eIdx} label={entry?.key ?? ''}>
                        <span style={{ wordBreak: 'break-all' }}>{entry?.value ?? ''}</span>
                      </Descriptions.Item>
                    ))}
                  </Descriptions>
                )
              ) : (
                // 收起态用轻量提示行，不用 Empty 插画——Empty 语义是"无数据"，
                // 大片收起时会让整页看起来像加载失败/空表
                <Text type="secondary" style={{ fontSize: 12 }}>
                  {t('collapsed')}
                </Text>
              )}
            </Card>
          );
        })
      )}

      {/* 慢日志 */}
      <Card size="small" style={{ marginBottom: 16 }} title={<Text strong>{t('slowlog')}</Text>}>
        {slowlog.length === 0 ? (
          <Empty description={t('noSlowlog')} />
        ) : (
          <Table
            rowKey="id"
            size="small"
            bordered
            columns={slowLogColumns}
            dataSource={slowlog}
            pagination={false}
          />
        )}
      </Card>

      <Title level={5} type="secondary" style={{ marginTop: 16 }}>
        {t('disclaimer')}
      </Title>
    </ContentContainer>
  );
};

export default RedisCacheMonitor;
