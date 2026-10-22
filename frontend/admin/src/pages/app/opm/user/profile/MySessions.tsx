import { useMemo } from 'react';
import { Table, Tag, Typography, App } from 'antd';
import type { TableColumnsType } from 'antd';
import { LogoutOutlined } from '@ant-design/icons';
import { useTranslation } from 'react-i18next';
import TableActions from '@/components/common/TableActions';
import {
  useListMyOnlineSessions,
  useRevokeMyOnlineSession,
} from '@/api/hooks/online-session';
import { useUserStore } from '@/stores/user';
import { useAuthStore } from '@/stores';
import type { online_sessionservicev1_OnlineSession as OnlineSession } from '@/api/generated/admin/service/v1';
import { formatDateTime } from '@/utils/date';

const { Text } = Typography;

// 从未知错误对象取本地化 message（request-client 拦截器已把 reason 译文写进 message）
const errMessage = (error: unknown, fallback: string): string => {
  if (error && typeof error === 'object' && 'message' in error) {
    const m = (error as { message?: unknown }).message;
    if (typeof m === 'string' && m) return m;
  }
  return fallback;
};

/**
 * 个人中心「活跃会话」：当前登录用户自己的全部在线会话自助管理。
 * 当前会话仅展示不可踢除（踢除会立即登出）；其它会话可单独下线。
 */
const MySessions = () => {
  const { t } = useTranslation('online-session');
  const { message } = App.useApp();

  // 当前 jti 从持久化的 access token payload 解出（与工作台同源）；
  // CMS 的 ListMy 不透传身份，admin BFF 要求显式传 user_id + current_jti。
  // userId 读 useUserStore（持久化、刷新后可用）；auth store 的 userInfo 仅登录期存在
  const accessToken = useAuthStore((s) => s.accessToken);
  const userId = useUserStore((s) => s.userInfo?.id);
  const currentJti = useMemo(() => {
    try {
      const payload = JSON.parse(
        atob((accessToken ?? '').split('.')[1].replace(/-/g, '+').replace(/_/g, '/')),
      );
      return typeof payload.jti === 'string' ? payload.jti : undefined;
    } catch {
      return undefined;
    }
  }, [accessToken]);

  const { data, isLoading } = useListMyOnlineSessions(currentJti, userId);
  const revokeMutation = useRevokeMyOnlineSession({
    onSuccess: () => {
      message.success(t('revokeSuccess'));
    },
  });

  const handleRevoke = async (record: OnlineSession) => {
    if (!record.jti) return;
    try {
      await revokeMutation.mutateAsync({
        clientType: record.clientType,
        jti: record.jti,
      });
    } catch (error) {
      console.error('Failed to revoke session:', error);
      message.error(errMessage(error, t('revokeFailed')));
    }
  };

  const columns: TableColumnsType<OnlineSession> = [
    {
      title: t('clientType'),
      dataIndex: 'clientType',
      width: 110,
      render: (_, record) =>
        record.clientType === 'app' ? (
          <Tag color="purple">{t('clientApp')}</Tag>
        ) : (
          <Tag color="blue">{t('clientAdmin')}</Tag>
        ),
    },
    {
      title: t('ipAddress'),
      dataIndex: 'ipAddress',
      width: 140,
      render: (val) => val || '-',
    },
    {
      title: t('userAgent'),
      dataIndex: 'userAgent',
      ellipsis: true,
      render: (val) => val || '-',
    },
    {
      title: t('loginAt'),
      dataIndex: 'loginAt',
      width: 170,
      render: (val) => formatDateTime(val),
    },
    {
      title: t('status'),
      width: 140,
      render: (_, record) =>
        record.current ? (
          <Tag color="green">{t('currentSession')}</Tag>
        ) : (
          <TableActions
            record={record}
            actions={[
              {
                key: 'revoke',
                label: t('revoke'),
                icon: <LogoutOutlined />,
                danger: true,
                confirm: { title: t('revokeConfirm'), okLoading: revokeMutation.isPending },
                onClick: () => handleRevoke(record),
              },
            ]}
          />
        ),
    },
  ];

  return (
    <div>
      <Text type="secondary" style={{ fontSize: 12, display: 'block', marginBottom: 8 }}>
        {t('mySessionsHint')}
      </Text>
      <Table<OnlineSession>
        rowKey={(record) => `${record.clientType}-${record.jti}`}
        dataSource={data?.items ?? []}
        loading={isLoading}
        pagination={false}
        size="small"
        columns={columns}
        locale={{ emptyText: t('noSessions') }}
      />
    </div>
  );
};

export default MySessions;
