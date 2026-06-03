"use client";

import React, {useEffect} from 'react';
import {useRouter} from 'next/navigation';
import {env} from "@/config";

const getDefaultLocale = () => {
    return env.defaultLocale;
};

const HomePage: React.FC = () => {
    const router = useRouter();

    useEffect(() => {
        const locale = getDefaultLocale();
        router.replace(`/${locale}`);
    }, [router]);

    return (
        <div className="flex h-screen w-full flex-col items-center justify-center bg-background">
            {/* 呼吸灯光晕 */}
            <div className="relative flex flex-col items-center gap-6">
                <div className="relative">
                    {/* 外圈脉冲光环 */}
                    <div className="absolute inset-0 -m-4 rounded-full bg-primary/10 animate-ping"
                         style={{animationDuration: '2s', animationIterationCount: 'infinite'}}/>
                    {/* 内圈光晕 */}
                    <div className="absolute inset-0 -m-2 rounded-full bg-primary/20 animate-pulse"/>
                    {/* Logo 圆 */}
                    <div
                        className="relative flex h-16 w-16 items-center justify-center rounded-2xl bg-primary/10 backdrop-blur-sm border border-primary/20">
                        <svg width="28" height="28" viewBox="0 0 24 24" fill="none" className="text-primary">
                            <path d="M12 2L2 7l10 5 10-5-10-5z" stroke="currentColor" strokeWidth="1.5"
                                  strokeLinecap="round" strokeLinejoin="round"/>
                            <path d="M2 17l10 5 10-5" stroke="currentColor" strokeWidth="1.5" strokeLinecap="round"
                                  strokeLinejoin="round"/>
                            <path d="M2 12l10 5 10-5" stroke="currentColor" strokeWidth="1.5" strokeLinecap="round"
                                  strokeLinejoin="round"/>
                        </svg>
                    </div>
                </div>
                <div className="flex flex-col items-center gap-2">
                    <p className="text-sm font-medium tracking-wide text-foreground/80">
                        GoWind CMS
                    </p>
                    <p className="text-xs text-muted-foreground/60 tracking-wider">
                        Loading...
                    </p>
                </div>
            </div>
        </div>
    );
};

export default HomePage;
