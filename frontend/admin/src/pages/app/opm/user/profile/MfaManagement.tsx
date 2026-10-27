import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import { useTranslation } from 'react-i18next';
import { App, Button, Input, List, Modal, QRCode, Typography } from 'antd';
import { useState } from 'react';
import { apiClient } from '@/api/client';

/**
 * 两步验证（TOTP）管理区块（个人中心 → 安全设置）：
 *   - 未绑定：展示开启入口 → 弹窗展示绑定二维码 + secret → 输入验证器动态码确认；
 *   - 已绑定：展示关闭入口 → 弹窗输入当前动态码解除绑定（会话劫持者无动态码即无法关闭）。
 */
const MfaManagement = () => {
  const { t } = useTranslation('profile');
  const { message } = App.useApp();
  const queryClient = useQueryClient();

  const [bindOpen, setBindOpen] = useState(false);
  const [unbindOpen, setUnbindOpen] = useState(false);
  const [binding, setBinding] = useState<{
    secret: string;
    otpAuthUrl: string;
    operationId: string;
  } | null>(null);
  const [totpCode, setTotpCode] = useState('');
  const [confirmLoading, setConfirmLoading] = useState(false);

  const statusQuery = useQuery({
    queryKey: ['myMfaStatus'],
    queryFn: () => apiClient.mFAService.GetMFAStatus({}),
    retry: 0,
  });
  const enabled = !!statusQuery.data?.enabled;

  const invalidate = () => queryClient.invalidateQueries({ queryKey: ['myMfaStatus'] });

  // 开启：第一步拿 secret + 绑定 URI
  const startEnrollMutation = useMutation({
    mutationFn: () =>
      apiClient.mFAService.StartEnrollMethod({ method: 'TOTP' }),
    onSuccess: (resp) => {
      if (resp.totp) {
        setBinding({
          secret: resp.totp.secret || '',
          otpAuthUrl: resp.totp.otpAuthUrl || '',
          operationId: resp.operationId || '',
        });
      }
    },
    onError: (error) => {
      console.error('[MfaManagement] start enroll failed:', error);
      message.error(t('mfa.startFailed'));
    },
  });

  // 开启：第二步输码确认
  const confirmEnrollMutation = useMutation({
    mutationFn: (code: string) =>
      apiClient.mFAService.ConfirmEnrollMethod({
        method: 'TOTP',
        operationId: binding?.operationId || '',
        totpCode: code,
      }),
    onSuccess: (resp) => {
      if (!resp.success) {
        message.error(t('mfa.codeInvalid'));
        return;
      }
      message.success(t('mfa.enabled'));
      setBindOpen(false);
      setBinding(null);
      setTotpCode('');
      invalidate();
    },
    onError: (error) => {
      console.error('[MfaManagement] confirm enroll failed:', error);
      message.error(t('mfa.codeInvalid'));
    },
  });

  // 关闭
  const disableMutation = useMutation({
    mutationFn: (code: string) =>
      apiClient.mFAService.DisableMFA({ method: 'TOTP', totpCode: code }),
    onSuccess: () => {
      message.success(t('mfa.disabled'));
      setUnbindOpen(false);
      setTotpCode('');
      invalidate();
    },
    onError: (error) => {
      console.error('[MfaManagement] disable failed:', error);
      message.error(t('mfa.codeInvalid'));
    },
  });

  const openBind = () => {
    setTotpCode('');
    setBinding(null);
    setBindOpen(true);
    startEnrollMutation.mutate();
  };

  const submitCode = () => {
    if (!/^\d{6}$/.test(totpCode)) {
      message.error(t('mfa.codeFormatError'));
      return;
    }
    setConfirmLoading(true);
    (bindOpen ? confirmEnrollMutation : disableMutation)
      .mutateAsync(totpCode)
      .finally(() => setConfirmLoading(false));
  };

  return (
    <>
      <List
        dataSource={[
          {
            key: 'mfa',
            title: t('mfa.title'),
            description: enabled
              ? t('mfa.enabledDesc')
              : t('mfa.disabledDesc'),
          },
        ]}
        renderItem={(item) => (
          <List.Item
            actions={[
              <Button
                key="action"
                type="link"
                danger={enabled}
                loading={statusQuery.isLoading}
                onClick={enabled ? () => setUnbindOpen(true) : openBind}
              >
                {enabled ? t('mfa.disable') : t('mfa.enable')}
              </Button>,
            ]}
          >
            <List.Item.Meta title={item.title} description={item.description} />
          </List.Item>
        )}
      />

      {/* 绑定弹窗：二维码 + secret + 动态码确认 */}
      <Modal
        title={t('mfa.bindTitle')}
        open={bindOpen}
        onCancel={() => {
          setBindOpen(false);
          setBinding(null);
        }}
        onOk={submitCode}
        okText={t('mfa.confirm')}
        confirmLoading={confirmLoading || confirmEnrollMutation.isPending}
        okButtonProps={{ disabled: !binding }}
        destroyOnHidden
      >
        {binding ? (
          <div className="flex flex-col items-center gap-3 py-2">
            <QRCode value={binding.otpAuthUrl} size={192} />
            <Typography.Text type="secondary">
              {t('mfa.scanHint')}
            </Typography.Text>
            <Typography.Paragraph copyable={{ text: binding.secret }}>
              <Typography.Text code>
                {binding.secret.replace(/(.{4})/g, '$1 ')}
              </Typography.Text>
            </Typography.Paragraph>
            <Input
              value={totpCode}
              onChange={(e) =>
                setTotpCode(e.target.value.replace(/\D/g, '').slice(0, 6))
              }
              placeholder={t('mfa.codePlaceholder')}
              autoComplete="one-time-code"
              inputMode="numeric"
              maxLength={6}
              className="w-44 text-center tracking-[0.5em]"
              onPressEnter={submitCode}
            />
          </div>
        ) : (
          <div className="py-6 text-center">{t('mfa.loading')}</div>
        )}
      </Modal>

      {/* 解绑弹窗：输入当前动态码 */}
      <Modal
        title={t('mfa.unbindTitle')}
        open={unbindOpen}
        onCancel={() => setUnbindOpen(false)}
        onOk={submitCode}
        okText={t('mfa.confirmDisable')}
        okButtonProps={{ danger: true }}
        confirmLoading={confirmLoading || disableMutation.isPending}
        destroyOnHidden
      >
        <p className="mb-3 text-sm text-gray-500">{t('mfa.unbindHint')}</p>
        <Input
          value={totpCode}
          onChange={(e) =>
            setTotpCode(e.target.value.replace(/\D/g, '').slice(0, 6))
          }
          placeholder={t('mfa.codePlaceholder')}
          autoComplete="one-time-code"
          inputMode="numeric"
          maxLength={6}
          className="w-44 text-center tracking-[0.5em]"
          onPressEnter={submitCode}
        />
      </Modal>
    </>
  );
};

export default MfaManagement;
