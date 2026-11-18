'use client';

import {useCallback, useEffect, useRef, useState} from 'react';
import {useTranslations} from 'next-intl';
import {useAuth} from "@/api/hooks/auth";
import {apiClient, declareCaptchaSupport, setCaptchaValues} from '@/api/client';

/**
 * 验证码输入行：挂载即拉取验证码图（父级按 CAPTCHA_REQUIRED 条件挂载，
 * 挂载=点亮），提供点击刷新。id/value 任意变化都经 onCaptchaChange 上抛，
 * 由父级同步到请求 header。
 */
function CaptchaField({
    t,
    refTick,
    onCaptchaChange,
}: {
    t: (key: string) => string;
    refTick: number;
    onCaptchaChange: (id: string, value: string) => void;
}) {
    const [captchaId, setCaptchaId] = useState('');
    const [captchaImage, setCaptchaImage] = useState('');
    const [captchaValue, setCaptchaValue] = useState('');
    // 点击刷新信号：递增驱动重新拉图（不重挂载组件，输入框内容保留）
    const [reloadTick, setReloadTick] = useState(0);
    // captchaValue 的 ref 镜像：拉图回调上报当前输入值不依赖闭包新鲜度
    const valueRef = useRef('');
    valueRef.current = captchaValue;

    const report = useCallback((id: string, value: string) => {
        onCaptchaChange(id, value);
    }, [onCaptchaChange]);

    // 挂载/重挂载/点击刷新即拉图；
    // 拉图的 setState 都在 await 之后（异步），不触发 set-state-in-effect。
    useEffect(() => {
        let cancelled = false;
        (async () => {
            try {
                const resp = await apiClient.authenticationService.GenerateCaptcha({});
                if (cancelled) return;
                setCaptchaId(resp.captchaId ?? '');
                setCaptchaImage(resp.imageBase64 ?? '');
                report(resp.captchaId ?? '', valueRef.current);
            } catch (e) {
                // 验证码获取失败不阻断页面（登录时会再次校验），原始错误留控制台
                console.warn('fetch captcha failed', e);
            }
        })();
        return () => {
            cancelled = true;
        };
        // eslint-disable-next-line react-hooks/exhaustive-deps -- 仅挂载/重挂载/点击刷新时拉图
    }, [refTick, reloadTick]);

    const handleValue = (value: string) => {
        setCaptchaValue(value);
        report(captchaId, value);
    };

    return (
        <div className="space-y-2">
            <label htmlFor="login-captcha" className="block text-sm font-medium text-foreground">
                {t('login.captcha_label')}
            </label>
            <div className="flex w-full items-center gap-2">
                <input
                    id="login-captcha"
                    type="text"
                    className="min-w-0 flex-1 rounded-lg border border-border bg-background px-4 py-2.5 text-sm text-foreground transition-colors hover:border-primary focus:border-primary focus:outline-none focus:ring-[3px] focus:ring-primary/15"
                    value={captchaValue}
                    onChange={(e) => handleValue(e.target.value)}
                    placeholder={t('login.placeholder_code')}
                    autoComplete="off"
                    maxLength={6}
                />
                {/* 图片盒：点击刷新；加载失败退化为可点占位 */}
                <div
                    className="flex h-[42px] w-[126px] shrink-0 cursor-pointer items-center justify-center overflow-hidden rounded-lg border border-border bg-muted"
                    title={t('login.captcha_refresh')}
                    onClick={() => setReloadTick((tick) => tick + 1)}
                >
                    {captchaImage ? (
                        // eslint-disable-next-line @next/next/no-img-element
                        <img
                            src={captchaImage}
                            alt="captcha"
                            className="h-full w-full object-cover"
                        />
                    ) : (
                        <span className="text-xs text-muted-foreground">{t('login.captcha_refresh')}</span>
                    )}
                </div>
            </div>
        </div>
    );
}

export default function AccountLoginPage() {
    const t = useTranslations('authentication');
    const [username, setUsername] = useState('');
    const [password, setPassword] = useState('');
    const [rememberMe, setRememberMe] = useState(false);

    // 渐进式验证码：初始不要求；后端对失败达阈值的登录投递 CAPTCHA_REQUIRED
    // reason 后点亮输入与图片（条件挂载 CaptchaField）。
    const [captchaRequired, setCaptchaRequired] = useState(false);
    // 失败刷新信号：递增驱动 CaptchaField 重挂载重新拉图
    const [captchaRefreshTick, setCaptchaRefreshTick] = useState(0);
    const [error, setError] = useState('');

    const auth = useAuth();

    // captcha id/value 上抛即同步到请求 header（组件卸载/未要求时清空）
    const handleCaptchaChange = useCallback((id: string, value: string) => {
        setCaptchaValues(id, value);
    }, []);

    // 本端有验证码 UI：声明支持（后端仅对声明方启用验证码 gate）；
    // 卸载时清掉验证码 header，避免残留值污染后续请求。
    useEffect(() => {
        declareCaptchaSupport();
        return () => setCaptchaValues(undefined, undefined);
    }, []);

    const handleLogin = async () => {
        if (!username || !password) {
            return;
        }
        setError('');
        try {
            await auth.login({
                username: username,
                password: password,
            });
        } catch (e: unknown) {
            // 渐进式验证码：CAPTCHA_REQUIRED 点亮输入与图片；已处于要求态的失败刷新图片。
            const reason = (e as {reason?: string} | null)?.reason;
            if (reason === 'CAPTCHA_REQUIRED') {
                if (!captchaRequired) {
                    setCaptchaRequired(true);
                } else {
                    setCaptchaRefreshTick((tick) => tick + 1);
                }
            } else if (captchaRequired) {
                setCaptchaRefreshTick((tick) => tick + 1);
            }
            // 展示错误（request-client 已按 reason 字典翻译 message；兜底走 i18n）
            setError((e as {message?: string} | null)?.message || t('login.login_failed'));
        }
    };

    return (
        <div className="space-y-4">
            <div className="space-y-2">
                <label htmlFor="login-username" className="block text-sm font-medium text-foreground">
                    {t('login.username')}
                </label>
                <input
                    id="login-username"
                    type="text"
                    className="w-full rounded-lg border border-border bg-background px-4 py-2.5 text-sm text-foreground transition-colors hover:border-primary focus:border-primary focus:outline-none focus:ring-[3px] focus:ring-primary/15"
                    value={username}
                    onChange={(e) => setUsername(e.target.value)}
                    placeholder={t('register.input_username')}
                    autoComplete="username"
                    autoCapitalize="none"
                    autoCorrect="off"
                    spellCheck="false"
                />
            </div>

            <div className="space-y-2">
                <label htmlFor="login-password" className="block text-sm font-medium text-foreground">
                    {t('login.password')}
                </label>
                <input
                    id="login-password"
                    type="password"
                    className="w-full rounded-lg border border-border bg-background px-4 py-2.5 text-sm text-foreground transition-colors hover:border-primary focus:border-primary focus:outline-none focus:ring-[3px] focus:ring-primary/15"
                    value={password}
                    onChange={(e) => setPassword(e.target.value)}
                    placeholder={t('login.placeholder_password')}
                    autoComplete="current-password"
                    onKeyDown={(e) => e.key === 'Enter' && handleLogin()}
                />
            </div>

            {captchaRequired && (
                <CaptchaField
                    key={captchaRefreshTick}
                    t={t as unknown as (key: string) => string}
                    refTick={captchaRefreshTick}
                    onCaptchaChange={handleCaptchaChange}
                />
            )}

            {error && (
                <p className="text-sm text-destructive" role="alert">{error}</p>
            )}

            <div className="flex w-full items-center justify-between">
                <label className="flex cursor-pointer items-center gap-2 text-sm text-muted-foreground">
                    <input
                        type="checkbox"
                        checked={rememberMe}
                        onChange={(e) => setRememberMe(e.target.checked)}
                        className="h-4 w-4 rounded border-border accent-primary"
                    />
                    <span>{t('login.remember_me')}</span>
                </label>
                {/* 忘记密码依赖邮件重置链路(后端暂未实现),先以禁用态展示 */}
                <button
                    className="cursor-not-allowed text-sm text-muted-foreground opacity-60 bg-transparent border-none"
                    disabled
                    title={t('login.coming_soon')}
                >
                    {t('login.forgot_password')}
                </button>
            </div>

            <button
                className="w-full cursor-pointer rounded-lg bg-primary px-4 py-2.5 text-sm font-medium text-primary-foreground transition-all hover:bg-primary/90 active:scale-[0.98]"
                onClick={handleLogin}
            >
                {t('login.login')}
            </button>
        </div>
    );
}
