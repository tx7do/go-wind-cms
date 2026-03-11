"use client";
import React, {useEffect} from 'react';
import {useRouter} from 'next/navigation';

import {env} from "@/config";

const getDefaultLocale = () => {
    return env.defaultLocale; // 或根据浏览器、配置等动态获取
};

const HomePage: React.FC = () => {
    const router = useRouter();

    useEffect(() => {
        const locale = getDefaultLocale();
        router.replace(`/${locale}`);
    }, [router]);

    return (
        <div>
            {/* 跳转前可显示 loading 或空白 */}
            <main>
                <h1>跳转中...</h1>
            </main>
        </div>
    );
};

export default HomePage;
