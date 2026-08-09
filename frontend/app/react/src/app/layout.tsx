import React from "react";
import type {Metadata} from "next";

import './globals.css';

import {env} from "@/config";
import StoreProvider from '@/store/StoreProvider';
import ThemeClientProvider from '@/components/layout/ThemeClientProvider';
import DevToolsClient from '@/components/dev/DevToolsClient';

export const metadata: Metadata = {
    title: env.appTitle,
    description: env.appDescription,
}

// 注意：<html>/<head>/<body> 由 [locale]/layout.tsx 渲染。
// 静态导出 + [locale] 段下，locale 信息只在 [locale]/layout.tsx 可用（那里的
// setRequestLocale 已执行，next-intl 的 locale API 可用）。若在根布局调用
// getLocale()，会因 setRequestLocale 尚未执行而回落到 headers()，触发静态
// 导出报错。因此根布局只保留全局 providers，不再输出 <html>。
export default function RootLayout({children}: { children: React.ReactNode }) {
    return (
        <StoreProvider>
            <ThemeClientProvider>{children}</ThemeClientProvider>
        </StoreProvider>
    );
}
