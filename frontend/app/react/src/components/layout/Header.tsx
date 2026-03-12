'use client';

import React from 'react';
import {useTranslations} from 'next-intl';
import {Image, Button, Divider, Space, Dropdown} from 'antd';
import {
    UserOutlined,
    HomeOutlined,
    LogoutOutlined
} from '@ant-design/icons';

import {useI18n} from '@/i18n';
import {useI18nRouter} from '@/i18n/helpers/useI18nRouter';
import {useThemeStore} from '@/store/core/theme/hooks';

import TopNavbar from './TopNavbar';

import styles from './Header.module.css';


export default function Header() {
    const t = useTranslations('navbar.top');
    const appT = useTranslations('app');
    const brandTitle = appT('title');
    const menuT = useTranslations('menu');
    const themeStore = useThemeStore();
    const {changeLocale} = useI18n();
    const router = useI18nRouter();
    const isLogin = false; // TODO: 从 accessStore 获取

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
        // TODO: 实现登出逻辑
        console.log('Logout');
    };

    const userMenuItems = [
        {
            key: 'homepage',
            icon: <HomeOutlined/>,
            label: menuT('homepage'),
            onClick: handleClickUserHomepage,
        },
        {
            key: 'profile',
            icon: <UserOutlined/>,
            label: menuT('my_profile'),
            onClick: handleClickSettings,
        },
        {
            type: 'divider' as const,
        },
        {
            key: 'logout',
            icon: <LogoutOutlined/>,
            label: menuT('logout'),
            danger: true,
            onClick: handleClickLogout,
        },
    ];

    const languageMenuItems = [
        {
            key: 'zh-CN',
            label: t('language') === '语言' ? '简体中文' : 'Chinese',
        },
        {
            key: 'en-US',
            label: t('language') === '语言' ? 'English' : 'English',
        },
    ];
    const handleLanguageChange = ({key}: { key: string }) => {
        changeLocale(key);
    };
    const toggleDarkMode = () => {
        let currentEffectiveMode: 'dark' | 'light';
        if (themeStore.theme.mode === 'system') {
            currentEffectiveMode = window.matchMedia('(prefers-color-scheme: dark)').matches ? 'dark' : 'light';
        } else {
            currentEffectiveMode = themeStore.theme.mode;
        }
        const newMode: 'dark' | 'light' = currentEffectiveMode === 'dark' ? 'light' : 'dark';
        themeStore.setMode(newMode);
        if (typeof window !== 'undefined') {
            localStorage.setItem('themeMode', newMode);
        }
    };

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
                            {isLogin ? (
                                <>
                                    <Dropdown
                                        menu={{items: userMenuItems}}
                                        trigger={['hover']}
                                    >
                                        <Button
                                            type="text"
                                            className={styles.iconBtn}
                                            aria-label={menuT('my_profile')}
                                            icon={<UserOutlined/>}
                                        />
                                    </Dropdown>
                                </>
                            ) : (
                                <>
                                    <Divider orientation="vertical"/>
                                    <Button
                                        type="primary"
                                        className={styles.headerLoginBtn}
                                        aria-label={t('login')}
                                        onClick={handleClickLogin}
                                    >
                                        {t('login')}
                                    </Button>
                                    <Button
                                        type="primary"
                                        className={styles.headerRegisterBtn}
                                        aria-label={t('register')}
                                        onClick={handleClickRegister}
                                    >
                                        {t('register')}
                                    </Button>
                                </>
                            )}
                            <Dropdown
                                menu={{items: languageMenuItems, onClick: handleLanguageChange}}
                                trigger={['click']}
                            >
                                <Button
                                    type="text"
                                    className={styles.langBtn}
                                    aria-label="Language"
                                    icon={<span className={styles.langIcon}>{'🌐'}</span>}
                                />
                            </Dropdown>
                            <Dropdown
                                menu={{items: [
                                        {key: 'dark', label: '暗色', icon: <span>🌙</span>, onClick: () => themeStore.setMode('dark')},
                                        {key: 'light', label: '亮色', icon: <span>☀️</span>, onClick: () => themeStore.setMode('light')},
                                        {key: 'system', label: '跟随系统', icon: <span>🖥️</span>, onClick: () => themeStore.setMode('system')},
                                    ]}}
                                trigger={['click']}
                            >
                                <Button
                                    type="text"
                                    className={styles.themeBtn}
                                    aria-label="Toggle theme"
                                    icon={<span className={styles.themeIcon}>{themeStore.theme.mode === 'dark' ? '🌙' : themeStore.theme.mode === 'light' ? '☀️' : '🖥️'}</span>}
                                />
                            </Dropdown>
                        </Space>
                    </div>
                </div>
            </div>
        </div>
    );
}
