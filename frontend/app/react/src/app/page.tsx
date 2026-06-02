"use client";

import React, {useEffect} from 'react';
import {useRouter} from 'next/navigation'; // 使用原始 Next.js router
import {Spinner} from '@/components/ui/spinner';

import {env} from "@/config";


const getDefaultLocale = () => {
    return env.defaultLocale; // 或根据浏览器、配置等动态获取
};

const HomePage: React.FC = () => {
    const router = useRouter(); // 使用原始 router，不自动添加 locale

    useEffect(() => {
        const locale = getDefaultLocale();
        router.replace(`/${locale}`); // 直接替换到带 locale 的路径
    }, [router]);

    return (
        <div className="flex h-screen flex-col items-center justify-center">
            <Spinner size="lg"/>
            <main className="mt-6">
                <h1 className="text-lg text-muted-foreground">Redirecting to home page...</h1>
            </main>
        </div>
    );
};

export default HomePage;
