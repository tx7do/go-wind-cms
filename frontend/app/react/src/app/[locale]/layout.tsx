import React from "react";
import {notFound} from 'next/navigation';

import {isSupportedLocale, routing} from './routing';
import ClientLocaleLayout from './ClientLocaleLayout';
import getRequestConfig from '@/i18n/request';

type SupportedLocale = (typeof routing.locales)[number];

// 告诉 Next.js 在构建时为每个语言预生成静态页面
export function generateStaticParams() {
    return routing.locales.map((locale) => ({locale}));
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
    const requestLocale = Promise.resolve(validLocale);

    const {messages} = await getRequestConfig({requestLocale});
    return (
        <ClientLocaleLayout locale={validLocale} messages={messages ?? {}}>
            {children}
        </ClientLocaleLayout>
    );
}
