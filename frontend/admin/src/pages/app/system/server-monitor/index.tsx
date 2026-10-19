import { Card, Descriptions, Spin, Tag, Typography, Empty } from 'antd';
import { useTranslation } from 'react-i18next';
import type { server_monitorservicev1_ServerMonitorInfo } from '@/api/generated/admin/service/v1';
import { useServerMonitorInfo } from '@/api/hooks/monitor-dashboards';
import { formatDateTime } from '@/utils/date';
import ContentContainer from '@/layouts/components/PageContainer/ContentContainer';

const { Text, Title } = Typography;

/**
 * 服务监控页面（只读）：Go 运行时 / 数据库连接池 / 主机信息。
 * 数据来自 GET /admin/v1/server-monitor，无任何写操作。
 */
const ServerMonitor = () => {
  const { t } = useTranslation('server-monitor');

  const { data: info, isLoading } = useServerMonitorInfo({
    refetchInterval: 10_000, // 每 10 秒自动刷新
  });

  if (isLoading) {
    return (
      <ContentContainer heightMode="auto" scrollable padding="16px">
        <div style={{ display: 'flex', justifyContent: 'center', padding: 48 }}>
          <Spin size="large" />
        </div>
      </ContentContainer>
    );
  }

  const infoTyped = info as server_monitorservicev1_ServerMonitorInfo | undefined;
  const go = infoTyped?.go;
  const db = infoTyped?.database;
  const host = infoTyped?.host;

  // proto3 的 64 位整数按 JSON 惯例序列化为字符串（int64-as-string），先归一为 number
  const fmtBytes = (bytes?: number | string) => {
    const n = Number(bytes);
    if (bytes === undefined || bytes === null || Number.isNaN(n)) return '-';
    if (n >= 1024 ** 3) return `${(n / 1024 ** 3).toFixed(2)} GB`;
    if (n >= 1024 ** 2) return `${(n / 1024 ** 2).toFixed(2)} MB`;
    if (n >= 1024) return `${(n / 1024).toFixed(1)} KB`;
    return `${n} B`;
  };

  const fmtUptime = (seconds?: number | string) => {
    const secs = Number(seconds);
    if (seconds === undefined || seconds === null || Number.isNaN(secs)) return '-';
    const d = Math.floor(secs / 86400);
    const h = Math.floor((secs % 86400) / 3600);
    const m = Math.floor((secs % 3600) / 60);
    if (d > 0) return `${d}d ${h}h ${m}m`;
    if (h > 0) return `${h}h ${m}m`;
    return `${m}m ${Math.floor(secs % 60)}s`;
  };

  const renderRow = (label: string, value?: React.ReactNode) => (
    <Descriptions.Item key={label} label={label}>
      {value ?? '-'}
    </Descriptions.Item>
  );

  return (
    <ContentContainer heightMode="fixed" scrollable padding="16px">
      {/* 数据库 */}
      <Card size="small" style={{ marginBottom: 16 }} title={<Text strong>{t('database')}</Text>} extra={
        db?.pingOk === true ? <Tag color="success">{t('pingOk')}</Tag>
          : db?.pingOk === false ? <Tag color="error">{t('pingFail')}</Tag>
            : null
      }>
        {db ? (
          <Descriptions column={2} size="small" bordered colon>
            {renderRow(t('dbDriver'), db.driver)}
            {renderRow(t('maxOpenConnections'), String(db.maxOpenConnections ?? '-'))}
            {renderRow(t('openConnections'), String(db.openConnections ?? '-'))}
            {renderRow(t('inUseConnections'), String(db.inUseConnections ?? '-'))}
            {renderRow(t('idleConnections'), String(db.idleConnections ?? '-'))}
            {db.pingError ? renderRow(t('pingError'), <span style={{ color: 'var(--ant-color-error)' }}>{db.pingError}</span>) : null}
          </Descriptions>
        ) : (
          <Empty description={t('noData')} />
        )}
      </Card>

      {/* Go 运行时 */}
      <Card size="small" style={{ marginBottom: 16 }} title={<Text strong>{t('goRuntime')}</Text>}>
        {go ? (
          <Descriptions column={2} size="small" bordered colon>
            {renderRow(t('goVersion'), go.version)}
            {renderRow(t('numGoroutine'), String(go.numGoroutine ?? '-'))}
            {renderRow(t('memAlloc'), fmtBytes(go.memAllocBytes))}
            {renderRow(t('memSys'), fmtBytes(go.memSysBytes))}
            {renderRow(t('gcCycles'), String(go.gcCycles ?? '-'))}
            {renderRow(t('uptime'), `${fmtUptime(go.uptimeSeconds)} (${formatDateTime(go.startedAt)})`)}
          </Descriptions>
        ) : (
          <Empty description={t('noData')} />
        )}
      </Card>

      {/* 主机 */}
      <Card size="small" style={{ marginBottom: 16 }} title={<Text strong>{t('host')}</Text>}>
        {host ? (
          <Descriptions column={2} size="small" bordered colon>
            {renderRow(t('hostname'), host.hostname)}
            {renderRow(t('os'), `${host.os ?? '-'} / ${host.arch ?? '-'}`)}
            {renderRow(t('numCpu'), String(host.numCpu ?? '-'))}
          </Descriptions>
        ) : (
          <Empty description={t('noData')} />
        )}
      </Card>

      <Title level={5} type="secondary" style={{ marginTop: 16 }}>
        {t('disclaimer')}
      </Title>
    </ContentContainer>
  );
};

export default ServerMonitor;
