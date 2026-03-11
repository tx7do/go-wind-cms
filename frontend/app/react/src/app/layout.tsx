import type {Metadata} from "next";

import {env} from "@/config";

export const metadata: Metadata = {
    title: env.appTitle,
    description: env.appDescription,
};

export default function RootLayout({children}: {children: React.ReactNode}) {
    return (
        <html lang="zh-CN" suppressHydrationWarning>
        <body>
        {children}
        </body>
        </html>
    );
}
