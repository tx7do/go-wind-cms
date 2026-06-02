'use client';

import React, {useState, useEffect} from 'react';
import {useTranslations} from 'next-intl';
import {User, Home, LogOut, Globe, Moon, Sun, Monitor} from 'lucide-react';

import {Button} from '@/components/ui/button';
import {
    DropdownMenu,
    DropdownMenuContent,
    DropdownMenuItem,
    DropdownMenuSeparator,
    DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu';

import {useI18n} from '@/i18n';
import {useI18nRouter} from '@/i18n/helpers/useI18nRouter';
import {usePreferences} from '@/core/preferences';
import type {ThemeModeType} from '@/core/preferences';

import TopNavbar from './TopNavbar';

import styles from './Header.module.css';
import {useAccessStore} from '@/store/core/access/store';
import {useAuth} from '@/api/hooks/auth';

export default function Header() {
    const t = useTranslations('navbar');
    const appT = useTranslations('app');
    const menuT = useTranslations('menu');
    const brandTitle = appT('title');

    const {theme: themePref, setThemeMode} = usePreferences();
    const currentMode = themePref.mode;
    const {changeLocale} = useI18n();
    const router = useI18nRouter();
    const accessStore = useAccessStore();
    const auth = useAuth();

    const accessToken = accessStore.accessToken;
    const isLogin = !!accessToken && !accessStore.loginExpired;

    const handleClickLogo = () => {
        router.push('/');
    };
    const handleClickSettings = () => {
        router.push('/settings');
    };
    const handleClickUserHomepage = () => {
        router.push('/user');
    };
    const handleClickLogin = () => {
        router.push('/login');
    };
    const handleClickRegister = () => {
        router.push('/register');
    };
    const handleClickLogout = async () => {
        console.log('Logout');
        if (isLogin) {
            await auth.logout();
        }
    };

    // SSR 兼容
    const [mounted, setMounted] = useState(false);
    useEffect(() => {
        // eslint-disable-next-line react-hooks/set-state-in-effect
        setMounted(true);
    }, []);

    const themeIcon = (() => {
        if (!mounted) return <Monitor className="h-4 w-4"/>;
        if (currentMode === 'dark') return <Moon className="h-4 w-4"/>;
        if (currentMode === 'light') return <Sun className="h-4 w-4"/>;
        return <Monitor className="h-4 w-4"/>;
    })();

    return (
        <div className={styles.fixedTop}>
            <div className={styles.headerInner}>
                <div className={styles.topBar}>
                    {/* Logo + 导航区 */}
                    <div className={styles.logoNavSection}>
                        <div
                            className={styles.logoSection}
                            role="button"
                            tabIndex={0}
                            aria-label="Go to homepage"
                            onClick={handleClickLogo}
                            onKeyDown={(e) => {
                                if (e.key === 'Enter' || e.key === ' ') {
                                    handleClickLogo();
                                }
                            }}
                        >
                            {/* eslint-disable-next-line @next/next/no-img-element */}
                            <img
                                src="/logo.png"
                                alt="Logo"
                                width={55}
                                height={55}
                                className={styles.logo}
                            />
                            <span className={styles.siteName}>{brandTitle}</span>
                        </div>
                        {/* 主导航菜单 */}
                        <div className={styles.navbarMenuWrap}>
                            <TopNavbar/>
                        </div>
                    </div>
                    {/* 功能按钮区 */}
                    <div className={styles.actions}>
                        <div className="flex items-center gap-2">
                            {/* 用户菜单 */}
                            <DropdownMenu>
                                <DropdownMenuTrigger asChild>
                                    <Button variant="ghost" size="icon" aria-label="User menu">
                                        <User className="h-4 w-4"/>
                                    </Button>
                                </DropdownMenuTrigger>
                                <DropdownMenuContent align="end">
                                    {isLogin ? (
                                        <>
                                            <DropdownMenuItem onClick={handleClickUserHomepage}>
                                                <Home className="h-4 w-4"/>
                                                {menuT('homepage')}
                                            </DropdownMenuItem>
                                            <DropdownMenuItem onClick={handleClickSettings}>
                                                <User className="h-4 w-4"/>
                                                {menuT('my_profile')}
                                            </DropdownMenuItem>
                                            <DropdownMenuSeparator/>
                                            <DropdownMenuItem
                                                onClick={handleClickLogout}
                                                className="text-destructive"
                                            >
                                                <LogOut className="h-4 w-4"/>
                                                {menuT('logout')}
                                            </DropdownMenuItem>
                                        </>
                                    ) : (
                                        <>
                                            <DropdownMenuItem onClick={handleClickLogin}>
                                                <User className="h-4 w-4"/>
                                                {t('user.login')}
                                            </DropdownMenuItem>
                                            <DropdownMenuItem onClick={handleClickRegister}>
                                                <User className="h-4 w-4"/>
                                                {t('user.register')}
                                            </DropdownMenuItem>
                                        </>
                                    )}
                                </DropdownMenuContent>
                            </DropdownMenu>

                            {/* 语言菜单 */}
                            <DropdownMenu>
                                <DropdownMenuTrigger asChild>
                                    <Button variant="ghost" size="icon" aria-label="Language">
                                        <Globe className="h-4 w-4"/>
                                    </Button>
                                </DropdownMenuTrigger>
                                <DropdownMenuContent align="end">
                                    <DropdownMenuItem onClick={() => changeLocale('zh-CN')}>
                                        简体中文
                                    </DropdownMenuItem>
                                    <DropdownMenuItem onClick={() => changeLocale('en-US')}>
                                        English
                                    </DropdownMenuItem>
                                </DropdownMenuContent>
                            </DropdownMenu>

                            {/* 主题菜单 */}
                            <DropdownMenu>
                                <DropdownMenuTrigger asChild>
                                    <Button variant="ghost" size="icon" aria-label="Toggle theme">
                                        {themeIcon}
                                    </Button>
                                </DropdownMenuTrigger>
                                <DropdownMenuContent align="end">
                                    <DropdownMenuItem onClick={() => setThemeMode('dark')}>
                                        <Moon className="h-4 w-4"/>
                                        {t('theme.dark')}
                                    </DropdownMenuItem>
                                    <DropdownMenuItem onClick={() => setThemeMode('light')}>
                                        <Sun className="h-4 w-4"/>
                                        {t('theme.light')}
                                    </DropdownMenuItem>
                                    <DropdownMenuItem onClick={() => setThemeMode('auto')}>
                                        <Monitor className="h-4 w-4"/>
                                        {t('theme.system')}
                                    </DropdownMenuItem>
                                </DropdownMenuContent>
                            </DropdownMenu>
                        </div>
                    </div>
                </div>
            </div>
        </div>
    );
}
