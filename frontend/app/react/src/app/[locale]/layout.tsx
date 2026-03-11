import {notFound} from 'next/navigation';

import {routing} from './routing';
import ClientLocaleLayout from './ClientLocaleLayout';
import getRequestConfig from '@/i18n/config';

export default async function LocaleLayout({children, params}: {
    children: React.ReactNode;
    params: Promise<{ locale: string }>;
}) {
    const {locale} = await params;
    if (!routing.locales.includes(locale as 'zh-CN' | 'en-US')) {
        notFound();
    }
    // 获取 messages
    const {messages} = await getRequestConfig({locale});
    return (
        <ClientLocaleLayout locale={locale} messages={messages}>
            {children}
        </ClientLocaleLayout>
    );
}
