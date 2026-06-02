'use client';

import React, {useState, useEffect} from 'react';
import {useTranslations} from 'next-intl';
import {Image, Button, Space, Dropdown} from 'antd';
import {
    UserOutlined,
    HomeOutlined,
    LogoutOutlined
} from '@ant-design/icons';

import {useI18n} from '@/i18n';
import {useI18nRouter} from '@/i18n/helpers/useI18nRouter';
import {usePreferences} from '@/core/preferences';
import type {ThemeModeType} from '@/core/preferences';

import TopNavbar from './TopNavbar';

import styles from './Header.module.css';
import {useAccessStore} from "@/store/core/access/store";
import {useAuth} from "@/api/hooks/auth";

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

    // 根据登录状态动态生成用户菜单项
    const userMenuItems = isLogin
        ? [
            {
                key: 'homepage',
                label: menuT('homepage'),
                icon: <HomeOutlined/>,
                onClick: handleClickUserHomepage
            },
            {
                key: 'profile',
                label: menuT('my_profile'),
                icon: <UserOutlined/>,
                onClick: handleClickSettings
            },
            {
                type: 'divider' as const
            },
            {
                key: 'logout',
                label: menuT('logout'),
                icon: <LogoutOutlined/>,
                danger: true,
                onClick: handleClickLogout
            },
        ]
        : [
            {
                key: 'login',
                label: t('user.login'),
                icon: <UserOutlined/>,
                onClick: handleClickLogin
            },
            {
                key: 'register',
                label: t('user.register'),
                icon: <UserOutlined/>,
                onClick: handleClickRegister
            }
        ];
    // 占位语言菜单
    const languageMenuItems = [
        {
            key: 'zh-CN',
            label: '简体中文',
            onClick: () => changeLocale('zh-CN')
        },
        {
            key: 'en-US',
            label: 'English',
            onClick: () => changeLocale('en-US')
        }
    ];
    // 占位语言切换
    const handleLanguageChange = ({key}: { key: string }) => {
        changeLocale(key);
    };

    const themeMenuItems = [
        {
            key: 'dark',
            label: t('theme.dark'),
            icon: <span>🌙</span>,
            onClick: () => setThemeMode('dark')
        },
        {
            key: 'light',
            label: t('theme.light'),
            icon: <span>☀️</span>,
            onClick: () => setThemeMode('light')
        },
        {
            key: 'auto',
            label: t('theme.system'),
            icon: <span>🖥️</span>,
            onClick: () => setThemeMode('auto')
        },
    ];

    // 主题图标映射
    const themeIconMap: Record<ThemeModeType, string> = {
        dark: '🌙',
        light: '☀️',
        auto: '🖥️'
    };

    const validMode: ThemeModeType = (currentMode && ['dark', 'light', 'auto'].includes(currentMode))
        ? currentMode as ThemeModeType
        : 'auto';

    // SSR 兼容：只在客户端渲染图标，避免 hydration 不匹配
    const [mounted, setMounted] = useState(false);
    useEffect(() => {
        // eslint-disable-next-line react-hooks/set-state-in-effect
        setMounted(true);
    }, []);

    const displayIcon = mounted ? themeIconMap[validMode] : '🖥️';

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
                            <Image
                                src={'/logo.png'}
                                alt="Logo"
                                width={55}
                                height={55}
                                className={styles.logo}
                                preview={false}
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
                        <Space size="middle">
                            <Dropdown
                                menu={{
                                    items: userMenuItems
                                }}
                                trigger={['click']}
                                popupRender={menu => (
                                    <div className={styles.dropdownMenu}>{menu}</div>
                                )}
                            >
                                <Button
                                    type="text"
                                    className={styles.iconBtn}
                                    aria-label="User menu"
                                    icon={<UserOutlined/>}
                                />
                            </Dropdown>
                            <Dropdown
                                menu={{items: languageMenuItems, onClick: handleLanguageChange}}
                                trigger={['click']}
                            >
                                <Button
                                    type="text"
                                    className={styles.iconBtn}
                                    aria-label="Language"
                                    icon={<span className={styles.langIcon}>{'🌐'}</span>}
                                />
                            </Dropdown>
                            <Dropdown
                                menu={{items: themeMenuItems}}
                                trigger={['click']}
                            >
                                <Button
                                    type="text"
                                    className={styles.iconBtn}
                                    aria-label="Toggle theme"
                                    icon={<span className={styles.themeIcon}>{displayIcon}</span>}
                                />
                            </Dropdown>
                        </Space>
                    </div>
                </div>
            </div>
        </div>
    );
}
