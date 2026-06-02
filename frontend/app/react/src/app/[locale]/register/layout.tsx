import React from "react";
import {notFound} from 'next/navigation';
import {NextIntlClientProvider} from 'next-intl';
import GlobalLoading from '@/components/layout/GlobalLoading';
import {DEFAULT_TIME_ZONE} from '@/i18n';
import {isSupportedLocale, routing} from '../routing';
import getRequestConfig from '@/i18n/request';

type SupportedLocale = (typeof routing.locales)[number];

export default async function RegisterLayout({
    children,
    params,
}: {
    children: React.ReactNode;
    params: Promise<{ locale: string }>;
}) {
    const unwrappedParams = await params;
    const {locale} = unwrappedParams;

    const isValidLocale = isSupportedLocale(locale);

    // 无效语言直接返回 404
    if (!isValidLocale) {
        console.log('Invalid locale', locale);
        notFound();
    }

    const validLocale = locale as SupportedLocale;
    const requestLocale = Promise.resolve(validLocale);

    const {messages} = await getRequestConfig({requestLocale});

    return (
        <NextIntlClientProvider timeZone={DEFAULT_TIME_ZONE} locale={validLocale} messages={messages ?? {}}>
            <GlobalLoading/>
            <div className="min-h-screen">
                {children}
            </div>
        </NextIntlClientProvider>
    );
}
