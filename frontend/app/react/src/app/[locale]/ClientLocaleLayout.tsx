"use client";

import React from "react";
import {NextIntlClientProvider} from "next-intl";

import Header from "@/components/layout/Header";
import Footer from "@/components/layout/Footer";

import ReduxProvider from "@/store/ReduxProvider";

import styles from './layout.module.css';

interface ClientLocaleLayoutProps {
    locale: string;
      messages: Record<string, unknown>;
    children: React.ReactNode;
}

const ClientLocaleLayout: React.FC<ClientLocaleLayoutProps> = ({locale, messages, children}) => {
    return (
        <NextIntlClientProvider locale={locale} messages={messages}>
            <ReduxProvider>
                <div className={styles.appContainer}>
                    <Header/>
                    <main className={styles.content}>
                        {children}
                    </main>
                    <Footer/>
                </div>
            </ReduxProvider>
        </NextIntlClientProvider>
    );
};

export default ClientLocaleLayout;
