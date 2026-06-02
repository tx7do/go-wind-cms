"use client";

import React from "react";
import {NextIntlClientProvider} from "next-intl";
import {usePathname} from 'next/navigation';

import Header from "@/components/layout/Header";
import Footer from "@/components/layout/Footer";
import GlobalLoading from "@/components/layout/GlobalLoading";

import {DEFAULT_TIME_ZONE} from "@/i18n";

import styles from './layout.module.css';

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
            <div className={styles.appContainer}>
                {!isAuthPage && <Header/>}
                <main className={`${styles.content} ${isAuthPage ? styles.noHeader : ''}`}>
                    {children}
                </main>
                {!isAuthPage && <Footer/>}
                <GlobalLoading/>
            </div>
        </NextIntlClientProvider>
    );
};

export default ClientLocaleLayout;
