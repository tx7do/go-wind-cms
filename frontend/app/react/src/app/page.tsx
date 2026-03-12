"use client";

import React, {useEffect} from 'react';
import {useRouter} from 'next/navigation'; // 使用原始 Next.js router
import {Spin} from 'antd';

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
        <div style={{
            display: 'flex',
            flexDirection: 'column',
            justifyContent: 'center',
            alignItems: 'center',
            height: '100vh'
        }}>
            <Spin size="large"/>
            <main>
                <h1 style={{marginTop: 24}}>Redirecting to home page...</h1>
            </main>
        </div>
    );
};

export default HomePage;
