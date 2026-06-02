import React from "react";
import type {Metadata} from "next";

import {env} from "@/config";
import {DEFAULT_LANGUAGE} from "@/i18n";
import StoreProvider from '@/store/StoreProvider';
import ThemeClientProvider from '@/components/layout/ThemeClientProvider';
import initThemeScript from '@/utils/init-theme-script';

export const metadata: Metadata = {
    title: env.appTitle,
    description: env.appDescription,
};


export default function RootLayout({children}: { children: React.ReactNode }) {
    return (
        <html
            lang={DEFAULT_LANGUAGE}
            suppressHydrationWarning
            className="m-0 p-0"
        >
        <head>
            <script
                dangerouslySetInnerHTML={{
                    __html: `(${initThemeScript.toString()})()`,
                }}
            />
        </head>
        <body className="m-0 p-0">
            <StoreProvider>
                <ThemeClientProvider>{children}</ThemeClientProvider>
            </StoreProvider>
        </body>
        </html>
    );
}
