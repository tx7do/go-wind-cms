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

import TopNavbar from './TopNavbar';

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

    const handleClickLogo = () => router.push('/');
    const handleClickSettings = () => router.push('/settings');
    const handleClickUserHomepage = () => router.push('/user');
    const handleClickLogin = () => router.push('/login');
    const handleClickRegister = () => router.push('/register');
    const handleClickLogout = async () => {
        if (isLogin) await auth.logout();
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
        <header className="fixed top-0 left-0 right-0 z-1000 flex justify-center bg-background/95 backdrop-blur-md border-b border-border shadow-sm dark:shadow-lg">
            <div className="flex h-(--layout-header-height) w-full max-w-300 items-center gap-6 px-4 max-md:gap-3 max-md:px-3">
                {/* Logo */}
                <button
                    type="button"
                    onClick={handleClickLogo}
                    className="flex shrink-0 items-center gap-2 rounded-lg p-1 transition-colors hover:bg-primary/10"
                    aria-label="Go to homepage"
                >
                    {/* eslint-disable-next-line @next/next/no-img-element */}
                    <img
                        src="/logo.png"
                        alt="Logo"
                        width={36}
                        height={36}
                        className="h-9 w-9 shrink-0 max-md:h-8 max-md:w-8"
                    />
                    <span className="text-lg font-bold text-primary whitespace-nowrap max-md:hidden">
                        {brandTitle}
                    </span>
                </button>

                {/* 导航区 */}
                <div className="min-w-0 flex-1">
                    <TopNavbar/>
                </div>

                {/* 功能按钮区 */}
                <div className="flex shrink-0 items-center gap-1">
                    {/* 用户菜单 */}
                    <DropdownMenu modal={false}>
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
                    <DropdownMenu modal={false}>
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
                    <DropdownMenu modal={false}>
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
        </header>
    );
}
