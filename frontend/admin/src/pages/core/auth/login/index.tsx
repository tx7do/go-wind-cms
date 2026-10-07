import React, { useCallback, useEffect, useState } from 'react';
import { useTranslation } from 'react-i18next';
import { Form, Input, Button, Checkbox, App } from 'antd';
import { UserOutlined, LockOutlined, SafetyOutlined } from '@ant-design/icons';
import { useAuthStore } from '@/stores';
import { useNavigate, useSearchParams } from 'react-router-dom';
import { fetchGenerateCaptcha } from '@/api';

/**
 * 验证码控件：输入框与图片盒同排，组成单一受控组件挂进 Form.Item。
 * value/onChange 由 Form.Item 注入后透传给 Input；校验错误渲染在整行下方，
 * 验证码图不会被错误行挤偏（旧结构嵌套 Form.Item + items-center 会在报错时下坠）。
 */
const CaptchaField: React.FC<{
  value?: string;
  onChange?: (...event: any[]) => void;
  placeholder?: string;
  imageEl?: React.ReactNode;
}> = ({ value, onChange, placeholder, imageEl }) => (
  <div className="flex items-center gap-2">
    <Input
      prefix={<SafetyOutlined />}
      placeholder={placeholder}
      autoComplete="off"
      value={value}
      onChange={onChange}
    />
    {imageEl}
  </div>
);

const Login: React.FC = () => {
  const { t } = useTranslation('auth');
  const { login, loginLoading } = useAuthStore();
  const navigate = useNavigate();
  const [searchParams] = useSearchParams();
  const { message } = App.useApp();

  // 验证码状态
  const [captchaId, setCaptchaId] = useState<string>('');
  const [captchaImage, setCaptchaImage] = useState<string>('');
  const [captchaLoading, setCaptchaLoading] = useState(false);

  // 获取验证码
  const refreshCaptcha = useCallback(async () => {
    setCaptchaLoading(true);
    try {
      const resp = await fetchGenerateCaptcha();
      setCaptchaId(resp.captchaId ?? '');
      setCaptchaImage(resp.imageBase64 ?? '');
    } catch {
      // 验证码获取失败不阻断页面，登录时会再次校验
    } finally {
      setCaptchaLoading(false);
    }
  }, []);

  useEffect(() => {
    // React 18 StrictMode 下 effect 会执行两次，导致两个并发的
    // fetchGenerateCaptcha 请求；后到的响应会覆盖先到的 captchaId，
    // 而后端通常一次性消费 captcha，登录时用的 captchaId 可能对应已被
    // 先到请求作废的 captcha。这里用 cancelled 标志位丢弃首次（被 double-invoke
    // 的第一次）请求的结果，确保 state 始终是最后一次请求的值。
    let cancelled = false;
    const run = async () => {
      setCaptchaLoading(true);
      try {
        const resp = await fetchGenerateCaptcha();
        if (cancelled) return;
        setCaptchaId(resp.captchaId ?? '');
        setCaptchaImage(resp.imageBase64 ?? '');
      } catch {
        // 验证码获取失败不阻断页面，登录时会再次校验
      } finally {
        if (!cancelled) setCaptchaLoading(false);
      }
    };
    run();
    return () => {
      cancelled = true;
    };
  }, []);

  const handleSubmit = async (values: {
    username: string;
    password: string;
    tenant_code?: string;
    remember?: boolean;
    captcha?: string;
  }) => {
    // 跳转目标：redirect 参数 → 用户 homePath → 首页。
    // 校验必须为同源相对路径，防止开放重定向（如 ?redirect=https://evil.com 或 //evil.com）
    const resolveSafeRedirect = () => {
      const rawRedirect =
        searchParams.get('redirect') || useAuthStore.getState().userInfo?.homePath || '/';
      return typeof rawRedirect === 'string' &&
        rawRedirect.startsWith('/') &&
        !rawRedirect.startsWith('//')
        ? rawRedirect
        : '/';
    };
    try {
      await login(
        {
          username: values.username,
          password: values.password,
          tenant_code: values.tenant_code,
          grant_type: 'password',
        },
        // onSuccess：store 检测到回调即跳过默认的 window.location.href 整页跳转，
        // 改由这里立即经 react-router 跳转。必须赶在 AppRouter 重建 router
        // （createBrowserRouter 以 window.location 为初始地址）之前完成。
        // 旧实现是 store 整页跳 homePath 与本页 300ms 定时器 navigate(redirect)
        // 赛跑：location.href 同步改写 window.location 使新 router 落在 homePath，
        // 定时器的 pushState 又把地址栏推回 redirect，两者各说各话——
        // 复现为"URL /content/posts 渲染着仪表盘"，手动刷新才恢复。
        () => navigate(resolveSafeRedirect()),
        { id: captchaId, value: values.captcha ?? '' },
      );

      message.success(t('loginSuccess'));
    } catch (error: any) {
      // 登录失败后刷新验证码
      refreshCaptcha();
      // 弹出错误提示（与 CRUD 页面统一模式：优先用后端 message，兜底走 i18n）
      message.error(error?.message || t('loginFailed'));
    }
  };

  /**
   * 验证码图片盒：与输入框等高（h-11=44px），宽 132px 与后端 240×80（3:1）等比，
   * img object-cover 满高填充无留白。点击刷新验证码。
   */
  const captchaImageEl = (
    <div
      className="flex items-center justify-center overflow-hidden w-[132px] h-11 shrink-0 rounded-lg cursor-pointer border border-solid bg-white light:border-black/10 light:bg-black/[0.03]"
      title={t('captchaRefresh')}
      onClick={() => !captchaLoading && refreshCaptcha()}
    >
      {captchaImage ? (
        <img
          src={captchaImage}
          alt="captcha"
          className="h-full w-full object-cover"
        />
      ) : (
        <span className="text-slate-400 text-xs">
          {captchaLoading ? '...' : t('captchaRefresh')}
        </span>
      )}
    </div>
  );

  return (
    <div className="w-full max-w-[420px]">
      {/* 标题 */}
      <div className="mb-11">
        <h2 className="text-[34px] font-extrabold tracking-[-0.5px] mb-2.5 text-[color:var(--ant-color-text)]">
          {t('welcomeBack')}
        </h2>
        <p className="text-[15px] leading-relaxed text-[color:var(--ant-color-text-tertiary)]">
          {t('loginDescription')}
        </p>
      </div>

      {/* 登录表单卡片 —— 实底表面色 + 24px 大圆角 + 主色柔影（对齐 vben 认证面板） */}
      <div className="rounded-3xl border border-[color:var(--ant-color-border-secondary)] bg-[color:var(--ant-color-bg-container)] p-8 shadow-[0_12px_40px_-8px_rgba(0,107,230,0.18)]">
        <Form
          name="login"
          onFinish={handleSubmit}
          size="large"
          initialValues={{ remember: true }}
          className="login-form"
        >
          <Form.Item name="tenant_code" className="login-form-item">
            <Input
              prefix={<UserOutlined />}
              placeholder={t('tenantCodePlaceholder')}
              autoComplete="off"
            />
          </Form.Item>

          <Form.Item
            name="username"
            className="login-form-item"
            rules={[
              {
                required: true,
                message: t('usernameRequired'),
              },
            ]}
          >
            <Input
              prefix={<UserOutlined />}
              placeholder={t('usernamePlaceholder')}
              autoComplete="username"
            />
          </Form.Item>

          <Form.Item
            name="password"
            className="login-form-item"
            rules={[
              {
                required: true,
                message: t('passwordRequired'),
              },
            ]}
          >
            <Input.Password
              prefix={<LockOutlined />}
              placeholder={t('passwordPlaceholder')}
              autoComplete="current-password"
            />
          </Form.Item>

          {/* 验证码 —— 输入框与图片盒组成单一受控组件（CaptchaField），
              校验错误渲染在整行下方，图片盒与输入框始终保持等高对齐 */}
          <Form.Item
            name="captcha"
            rules={[
              {
                required: true,
                message: t('captchaRequired'),
              },
            ]}
            className="login-form-item"
          >
            <CaptchaField
              placeholder={t('captchaPlaceholder')}
              imageEl={captchaImageEl}
            />
          </Form.Item>

          <Form.Item className="login-remember-item">
            <Form.Item name="remember" valuePropName="checked" noStyle>
              <Checkbox>{t('rememberAccount')}</Checkbox>
            </Form.Item>
          </Form.Item>

          <Form.Item className="login-form-item">
            <Button
              type="primary"
              htmlType="submit"
              loading={loginLoading}
              block
              className="login-submit-btn"
            >
              {loginLoading ? t('loggingIn') : t('loginButton')}
            </Button>
          </Form.Item>
        </Form>
      </div>

      {/* 底部链接 */}
      <div className="mt-6 text-center text-[13px]">
        <span className="text-[color:var(--ant-color-text-secondary)]">
          {t('noAccount')}{' '}
        </span>
        {/* CMS 后端暂无找回密码 RPC，入口对齐 vben 侧改为禁用态（9e131069） */}
        <span
          className="cursor-not-allowed text-[color:var(--ant-color-text-quaternary)]"
          style={{ fontSize: 13 }}
          title={t('forgotPasswordDisabled')}
        >
          {t('forgotPassword')}
        </span>
        
      </div>
    </div>
  );
};

export default Login;
