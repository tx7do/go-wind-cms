import React from "react";
import {notFound} from 'next/navigation';
import {setRequestLocale} from 'next-intl/server';

import {isSupportedLocale, routing} from './routing';
import ClientLocaleLayout from './ClientLocaleLayout';
import DevToolsClient from '@/components/dev/DevToolsClient';
import initThemeScript from '@/utils/init-theme-script';
import getRequestConfig from '@/i18n/request';

type SupportedLocale = (typeof routing.locales)[number];

// 告诉 Next.js 在构建时为每个语言预生成静态页面
export function generateStaticParams() {
    return routing.locales.map((locale) => ({locale}));
}

/**
 * 将 locale 解析为 HTML 文档方向。
 *
 * 说明：next-intl 不提供 locale→direction 映射，这里用标准 `Intl.Locale` 的
 * `textInfo.direction`（返回 'ltr' | 'rtl'）。中英文为 'ltr'；若未来加入阿
 * 拉伯语（ar）等 RTL 语言，该 API 会返回 'rtl'，使 <html dir> 自动翻转，
 * 无需在此维护手写白名单。
 *
 * 兜底：解析失败时回落到 'ltr'，与历史行为一致。
 */
function resolveDirection(locale: string): 'ltr' | 'rtl' {
    try {
        // eslint-disable-next-line @typescript-eslint/no-explicit-any
        const dir = new (Intl as any).Locale(locale).textInfo?.direction;
        return dir === 'rtl' ? 'rtl' : 'ltr';
    } catch {
        return 'ltr';
    }
}

export default async function LocaleLayout({children, params}: {
    children: React.ReactNode;
    params: Promise<{ locale: string }>;
}) {
    const unwrappedParams = await params;
    const {locale} = unwrappedParams;

    const isValidLocale = isSupportedLocale(locale);

    // 无效语言直接返回 404
    if (!isValidLocale) {
        console.log('Invalid locale', locale)
        notFound();
    }

    const validLocale = locale as SupportedLocale;

    // 告诉 next-intl 当前 locale，避免内部调用 headers()（静态导出必需）
    setRequestLocale(validLocale);

    // 用当前 locale 设置 <html lang> 与 <html dir>。
    // 历史上 <html lang> 绑定到常量 DEFAULT_LANGUAGE='zh-CN'，导致 en-US 页面
    // 也被标成 zh-CN。这里改为按实际 locale 设置，并对 RTL 语言翻转方向。
    const dir = resolveDirection(validLocale);

    const requestLocale = Promise.resolve(validLocale);
    const {messages} = await getRequestConfig({requestLocale});
    return (
        <html
            lang={validLocale}
            dir={dir}
            suppressHydrationWarning
            data-scroll-behavior="smooth"
        >
        <head>
            <script
                dangerouslySetInnerHTML={{
                    __html: `(${initThemeScript.toString()})()`,
                }}
            />
        </head>
        <body>
            <ClientLocaleLayout locale={validLocale} messages={messages ?? {}}>
                {children}
            </ClientLocaleLayout>
            <DevToolsClient/>
        </body>
        </html>
    );
}
