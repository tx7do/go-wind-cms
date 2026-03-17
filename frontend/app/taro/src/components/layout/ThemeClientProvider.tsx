import React, {useEffect, useRef} from 'react';
import {useThemeMode} from '@/store/core/theme/hooks';

export default function ThemeClientProvider({children}: { children: React.ReactNode }) {
    const mode = useThemeMode();
    const mqRef = useRef<any>(null);

    useEffect(() => {
        const html = document.documentElement;

        console.log('[Theme Provider] mode:', mode, 'current classes:', html.className);

        // 清理之前的监听器
        if (mqRef.current) {
            mqRef.current.onchange = null;
        }

        if (mode === 'system') {
            // 跟随系统
            const mq = window.matchMedia('(prefers-color-scheme: dark)');
            mqRef.current = mq;

            //  立即应用当前系统主题
            const currentTheme = mq.matches ? 'dark' : 'light';

            // 只有当当前类名与系统主题不一致时才更新
            if (!html.classList.contains(currentTheme)) {
                console.log('[Theme Provider] Applying system theme:', currentTheme);
                html.classList.remove('dark', 'light');
                html.classList.add(currentTheme);
            }

            mq.onchange = (e) => {
                const newTheme = e.matches ? 'dark' : 'light';
                if (!html.classList.contains(newTheme)) {
                    console.log('[Theme Provider] System theme changed to:', newTheme);
                    html.classList.remove('dark', 'light');
                    html.classList.add(newTheme);
                }
            };
        } else {
            // 直接应用指定主题
            if (!html.classList.contains(mode)) {
                console.log('[Theme Provider] Applying stored mode:', mode);
                html.classList.remove('dark', 'light');
                html.classList.add(mode);
            }
        }

        // 清理函数
        return () => {
            if (mqRef.current) {
                mqRef.current.onchange = null;
            }
        };
    }, [mode]);

    return children;
}
