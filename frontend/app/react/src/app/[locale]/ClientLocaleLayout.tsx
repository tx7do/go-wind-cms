"use client";

import React from "react";
import {NextIntlClientProvider} from "next-intl";
import {usePathname} from 'next/navigation';

import Header from "@/components/layout/Header";
import Footer from "@/components/layout/Footer";
import NavigationProgress from "@/components/layout/NavigationProgress";

import {DEFAULT_TIME_ZONE} from "@/i18n";

import {cn} from "@/lib/utils";

interface ClientLocaleLayoutProps {
    locale: string;
    messages: Record<string, unknown>;
    children: React.ReactNode;
}

const ClientLocaleLayout: React.FC<ClientLocaleLayoutProps> = ({locale, messages, children}) => {
    const pathname = usePathname();

    // 检查是否是认证页面（register, login 等）
    const isAuthPage = pathname?.includes('/register') || pathname?.includes('/login');

    // 缓存清理由 useI18n.changeLocale 在导航前执行 queryClient.clear()，
    // 配合 key={locale} 强制重挂载内容区，确保子组件用新语言重新请求。
    return (
        <NextIntlClientProvider timeZone={DEFAULT_TIME_ZONE} locale={locale} messages={messages}>
            <div className="flex min-h-screen w-full flex-col">
                {!isAuthPage && <Header/>}
                {/* key={locale}：语言切换时强制重挂载内容区，触发所有 useEffect 重新加载数据 */}
                <main key={locale} className={cn(
                    'flex w-full flex-1 flex-col bg-background',
                    !isAuthPage && 'pt-(--layout-header-height) min-h-screen',
                )}>
                    {children}
                </main>
                {!isAuthPage && <Footer/>}
                <NavigationProgress/>
            </div>
        </NextIntlClientProvider>
    );
};

export default ClientLocaleLayout;
