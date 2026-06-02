"use client";

import React from "react";
import {NextIntlClientProvider} from "next-intl";
import {usePathname} from 'next/navigation';

import Header from "@/components/layout/Header";
import Footer from "@/components/layout/Footer";
import GlobalLoading from "@/components/layout/GlobalLoading";

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
    
    return (
        <NextIntlClientProvider timeZone={DEFAULT_TIME_ZONE} locale={locale} messages={messages}>
            <div className="flex min-h-screen w-full flex-col">
                {!isAuthPage && <Header/>}
                <main className={cn(
                    'flex w-full flex-1 flex-col bg-background',
                    !isAuthPage && 'pt-[var(--layout-header-height)] min-h-screen',
                )}>
                    {children}
                </main>
                {!isAuthPage && <Footer/>}
                <GlobalLoading/>
            </div>
        </NextIntlClientProvider>
    );
};

export default ClientLocaleLayout;
