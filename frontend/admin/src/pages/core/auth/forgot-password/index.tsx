import { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { Button, Form, Input, App } from 'antd';
import { useTranslation } from 'react-i18next';
import { apiClient } from '@/api/client';
import { encryptPassword } from '@/utils/crypto';
import AuthLayout from '@/components/bussiness/AuthLayout';

// 从未知错误对象取本地化 message（request-client 拦截器已把 reason 译文写进 message）
const errMessage = (error: unknown, fallback: string): string => {
  if (error && typeof error === 'object' && 'message' in error) {
    const m = (error as { message?: unknown }).message;
    if (typeof m === 'string' && m) return m;
  }
  return fallback;
};

/**
 * 忘记密码页面（两步）：
 * 步骤 1 输入绑定邮箱 → 发送重置验证码（邮件走通知域 EMAIL 渠道）；
 * 步骤 2 输入验证码 + 新密码 → 重置成功跳登录页。
 * 新密码与登录同规：AES 密文传输（后端 NeedDecrypt 解密）。
 */
const ForgotPassword = () => {
  const { t } = useTranslation('forgot-password');
  const navigate = useNavigate();
  const { message } = App.useApp();
  const [form] = Form.useForm();
  const [step, setStep] = useState<1 | 2>(1);
  const [identifier, setIdentifier] = useState('');
  const [sending, setSending] = useState(false);
  const [resetting, setResetting] = useState(false);

  const sendCode = async () => {
    const values = await form.validateFields(['identifier']);
    setSending(true);
    try {
      await apiClient.authenticationService.ForgotPassword({
        identifier: values.identifier,
      });
      setIdentifier(values.identifier);
      setStep(2);
      message.success(t('codeSent'));
    } catch (error) {
      // 原始错误必须留在控制台：用户可见的那句翻译不包含服务端的原因
      console.error('send reset code failed', error);
      message.error(errMessage(error, t('sendFailed')));
    } finally {
      setSending(false);
    }
  };

  const resetByCode = async (values: { code: string; newPassword: string }) => {
    setResetting(true);
    try {
      await apiClient.authenticationService.ResetPasswordByCode({
        identifier,
        code: values.code,
        newPassword: encryptPassword(values.newPassword),
      });
      message.success(t('resetSuccess'));
      navigate('/auth/login');
    } catch (error) {
      console.error('reset password by code failed', error);
      message.error(errMessage(error, t('resetFailed')));
    } finally {
      setResetting(false);
    }
  };

  return (
    <AuthLayout
      title={t('title')}
      description={t('subtitle')}
      footerLink={{
        text: t('backToLogin'),
        linkText: t('backToLogin'),
        href: '/auth/login',
      }}
    >
      {step === 1 ? (
        <Form form={form} layout="vertical" onFinish={sendCode}>
          <Form.Item
            name="identifier"
            label={t('identifier')}
            rules={[
              { required: true, message: t('requiredIdentifier') },
              { type: 'email' as const, message: t('requiredIdentifier') },
            ]}
          >
            <Input placeholder={t('identifierPlaceholder')} size="large" />
          </Form.Item>
          <Button type="primary" size="large" block loading={sending} onClick={sendCode}>
            {t('sendCode')}
          </Button>
        </Form>
      ) : (
        <Form layout="vertical" onFinish={resetByCode}>
          <Form.Item
            name="code"
            label={t('code')}
            rules={[{ required: true, message: t('requiredCode') }]}
          >
            <Input placeholder={t('codePlaceholder')} size="large" />
          </Form.Item>
          <Form.Item
            name="newPassword"
            label={t('newPassword')}
            rules={[{ required: true, message: t('requiredNewPassword') }]}
          >
            <Input.Password placeholder={t('newPasswordPlaceholder')} size="large" />
          </Form.Item>
          <Button type="primary" size="large" block loading={resetting} onClick={() => form.submit()}>
            {t('reset')}
          </Button>
        </Form>
      )}
    </AuthLayout>
  );
};

export default ForgotPassword;
