import React from "react";
import type {Metadata} from "next";
import Script from "next/script";

import {env} from "@/config";
import {DEFAULT_LANGUAGE} from "@/i18n";
import ThemeClientProvider from '@/components/layout/ThemeClientProvider';
import ReduxProvider from '@/store/ReduxProvider';
import initThemeScript from '@/utils/init-theme-script';

export const metadata: Metadata = {
    title: env.appTitle,
    description: env.appDescription,
};


export default function RootLayout({children}: {children: React.ReactNode}) {
    return (
        <html 
            lang={DEFAULT_LANGUAGE} 
            suppressHydrationWarning
            style={{margin: 0, padding: 0}}
        >
            <head>
                <script
                    dangerouslySetInnerHTML={{
                        __html: `(${initThemeScript.toString()})()`,
                    }}
                />
            </head>
            <body style={{margin: 0, padding: 0}}>
                <ReduxProvider>
                    <ThemeClientProvider>{children}</ThemeClientProvider>
                </ReduxProvider>
            </body>
        </html>
    );
}
