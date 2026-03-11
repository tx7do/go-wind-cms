"use client";
import React from "react";
import {NextIntlClientProvider} from "next-intl";
import Header from "@/components/layout/Header";
import Footer from "@/components/layout/Footer";
import ReduxProvider from "@/store/ReduxProvider";

interface ClientLocaleLayoutProps {
  messages: Record<string, unknown>;
  locale: string;
  children: React.ReactNode;
}

const ClientLocaleLayout: React.FC<ClientLocaleLayoutProps> = ({messages, locale, children}) => {
  return (
    <NextIntlClientProvider messages={messages} locale={locale}>
      <ReduxProvider>
        <div style={{display: 'flex', flexDirection: 'column', minHeight: '100vh'}}>
          <Header/>
          <main style={{flex: 1}}>
            {children}
          </main>
          <Footer/>
        </div>
      </ReduxProvider>
    </NextIntlClientProvider>
  );
};

export default ClientLocaleLayout;

