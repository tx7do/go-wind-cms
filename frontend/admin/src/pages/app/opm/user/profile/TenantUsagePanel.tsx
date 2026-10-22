import { Card, Progress, Spin } from 'antd';
import { useTranslation } from 'react-i18next';
import { useMyTenantUsage } from '@/api/hooks/tenant-quota';

// 64 位整数按 proto3 JSON 惯例为字符串，格式化前必须 Number() 归一
const num = (v: unknown): number => Number(v ?? 0);

const formatBytes = (bytes: number): string => {
  if (bytes >= 1024 ** 3) return `${(bytes / 1024 ** 3).toFixed(2)} GB`;
  if (bytes >= 1024 ** 2) return `${(bytes / 1024 ** 2).toFixed(2)} MB`;
  if (bytes >= 1024) return `${(bytes / 1024).toFixed(2)} KB`;
  return `${bytes} B`;
};

const quotaMeta: Record<
  string,
  { format: (v: number) => string }
> = {
  USER_LIMIT: { format: (v) => v.toLocaleString() },
  STORAGE: { format: formatBytes },
  API_CALL: { format: (v) => v.toLocaleString() },
  AI_TOKENS: { format: (v) => v.toLocaleString() },
};

/**
 * 租户用量面板（个人中心自助视图）：四类配额的水位条。
 * 无任何配额行时显示"未配置配额"提示（闸门 fail-open 的可见性面）。
 */
const TenantUsagePanel = () => {
  const { t } = useTranslation('tenant-quota');
  const { data: usage, isLoading } = useMyTenantUsage();

  if (isLoading) {
    return <Spin style={{ display: 'block', margin: '48px auto' }} />;
  }

  const quotas = usage?.quotas ?? [];
  const usedOf: Record<string, number> = {
    USER_LIMIT: num(usage?.userCount),
    STORAGE: num(usage?.storageUsedBytes),
    API_CALL: num(usage?.apiCallCount),
    AI_TOKENS: num(usage?.aiTokensMonth),
  };

  return (
    <div style={{ maxWidth: 640 }}>
      <p style={{ color: 'var(--ant-color-text-secondary)', marginBottom: 24 }}>
        {t('usage.hint')} {usage?.subscriptionPlan ? `（${t('usage.plan')}: ${usage.subscriptionPlan}）` : ''}
      </p>

      {quotas.length === 0 && (
        <Card size="small">{t('usage.noQuotas')}</Card>
      )}

      {quotas.map((q) => {
        const type = String(q.quotaType ?? '');
        const limit = num(q.quotaValue);
        const used = usedOf[type] ?? 0;
        const pct = limit > 0 ? Math.min(100, Math.round((used / limit) * 100)) : used > 0 ? 100 : 0;
        const fmt = quotaMeta[type]?.format ?? ((v: number) => v.toLocaleString());
        return (
          <div key={type} style={{ marginBottom: 20 }}>
            <div style={{ display: 'flex', justifyContent: 'space-between', marginBottom: 4 }}>
              <span>{t(`quotaTypeMap.${type}`, type)}</span>
              <span style={{ color: 'var(--ant-color-text-secondary)' }}>
                {fmt(used)} / {limit > 0 ? fmt(limit) : t('usage.unlimited')}（{pct}%）
              </span>
            </div>
            <Progress
              percent={pct}
              status={pct >= 100 ? 'exception' : pct >= 80 ? 'active' : 'normal'}
              size={['100%', 8]}
            />
          </div>
        );
      })}
    </div>
  );
};

export default TenantUsagePanel;
