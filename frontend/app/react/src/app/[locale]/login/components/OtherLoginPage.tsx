'use client';

import {useSearchParams} from 'next/navigation';
import {useTranslations} from 'next-intl';
import {useEffect, useRef, useState} from 'react';
import XIcon from '@/plugins/xicon';
import {apiClient} from '@/api/client';
import {useAuth} from '@/api/hooks/auth';

/**
 * 企业 SSO（OIDC 授权码模式）登录入口：
 *  - 按钮点击 → 后台取授权地址（state 由服务端生成暂存）→ 整页跳 IdP；
 *  - IdP 回调携带 code+state 回到登录页 → 本组件消费并换发本仓令牌。
 * OIDC 未配置时点击有后端错误提示，不影响本地账密登录。
 */
export default function OtherLoginPage() {
    const t = useTranslations('authentication');
    const searchParams = useSearchParams();
    const auth = useAuth();

    const [ssoLoading, setSsoLoading] = useState(false);
    // StrictMode 下 effect 双跑，回调只允许消费一次
    const ssoHandledRef = useRef(false);

    useEffect(() => {
        const code = searchParams.get('code');
        const state = searchParams.get('state');
        if (!code || !state || ssoHandledRef.current) return;
        ssoHandledRef.current = true;

        setSsoLoading(true);
        auth.login({
            grant_type: 'authorization_code',
            code,
            oidc_state: state,
            password: '',
        })
            .catch((error: unknown) => {
                console.error('OIDC login failed:', error);
            })
            .finally(() => setSsoLoading(false));
        // eslint-disable-next-line react-hooks/exhaustive-deps
    }, []);

    const handleSsoLogin = async () => {
        setSsoLoading(true);
        try {
            const resp = await apiClient.authenticationService.GetOidcAuthorizeUrl({
                // 回调落回本页：/login?code=&state= 由上方 effect 消费
                redirect_uri: `${window.location.origin}${window.location.pathname}`,
            });
            if (resp.authorizeUrl) {
                window.location.href = resp.authorizeUrl;
                return;
            }
            setSsoLoading(false);
        } catch (error) {
            console.error('Get OIDC authorize url failed:', error);
            setSsoLoading(false);
        }
    };

    return (
        <div className="space-y-3">
            <button
                className="flex w-full cursor-pointer items-center justify-center gap-3 rounded-lg border border-border bg-card px-4 py-3 text-sm font-medium text-foreground transition-all hover:border-primary hover:bg-primary/5 active:scale-[0.98] disabled:cursor-not-allowed disabled:opacity-60"
                disabled={ssoLoading}
                onClick={handleSsoLogin}
            >
                <XIcon name="mdi:shield-key-outline" size={20}/>
                <span>{ssoLoading ? t('login.sso_redirecting') : t('login.sso_enterprise')}</span>
            </button>
        </div>
    );
}
