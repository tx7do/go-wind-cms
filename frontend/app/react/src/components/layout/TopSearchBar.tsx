'use client';

import React from 'react';
import {Image, Button, Divider, Space, Dropdown} from 'antd';
import {
    UserOutlined,
    MoonOutlined,
    SunOutlined,
    GlobalOutlined,
    HomeOutlined,
    LogoutOutlined
} from '@ant-design/icons';
import {useTranslations} from 'next-intl';
import {useRouter} from 'next/navigation';

import styles from './TopSearchBar.module.css';

import {useLanguageStore} from "@/store/core/language/hooks";
import {useThemeStore} from '@/store/core/theme/hooks';

import {useI18n} from '@/i18n';

import SearchBar from './SearchBar';

interface TopSearchBarProps {
    brandTitle: string;
    logoSrc?: string;
    onLogoClick?: () => void;
}

export default function TopSearchBar({brandTitle, logoSrc = '/logo.png', onLogoClick}: TopSearchBarProps) {
    const t = useTranslations('navbar.top');
    const menuT = useTranslations('menu');
    const themeStore = useThemeStore();
    const languageStore = useLanguageStore();
    const {changeLocale} = useI18n();
    const router = useRouter();

    // 模拟登录状态（实际应该从 store 获取）
    const isLogin = false; // TODO: 从 accessStore 获取

    const handleClickLogo = () => {
        if (onLogoClick) {
            onLogoClick();
        } else {
            router.push('/');
        }
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

    // 用户菜单选项
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

    // 语言菜单选项
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
        const newMode = themeStore.theme.mode === 'dark' ? 'light' : 'dark';
        themeStore.setMode(newMode);
    };

    return (
        <div className={styles.topBar}>
            {/* Logo Section */}
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
                    src={logoSrc}
                    alt="Logo"
                    width={55}
                    height={55}
                    className={styles.logo}
                    preview={false}
                />
                <span className={styles.siteName}>{brandTitle}</span>
            </div>

            {/* Search Bar */}
            <SearchBar/>

            {/* Actions */}
            <div className={styles.actions}>
                <Space size="middle">
                    {/* User Menu */}
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
                                type="text"
                                className={styles.headerLoginBtn}
                                aria-label={t('login')}
                                onClick={handleClickLogin}
                            >
                                {t('login')}
                            </Button>
                            <Button
                                type="text"
                                className={styles.headerRegisterBtn}
                                aria-label={t('register')}
                                onClick={handleClickRegister}
                            >
                                {t('register')}
                            </Button>
                        </>
                    )}

                    {/* Language Switcher */}
                    <Dropdown
                        menu={{items: languageMenuItems, onClick: handleLanguageChange}}
                        trigger={['hover']}
                    >
                        <Button
                            shape="round"
                            className={styles.langBtn}
                            aria-label="Language"
                            icon={<GlobalOutlined/>}
                        >
                            <span
                                className={styles.langText}>
                                {languageStore.language.locale === 'zh-CN' ? '中文' : 'EN'}
                            </span>
                        </Button>
                    </Dropdown>

                    {/* Theme Toggle */}
                    <Button
                        shape="round"
                        className={styles.themeBtn}
                        aria-label="Toggle theme"
                        onClick={toggleDarkMode}
                        icon={themeStore.theme.mode === 'dark' || themeStore.theme.mode === 'system' ? <SunOutlined/> :
                            <MoonOutlined/>}
                    >
                      <span className={styles.themeText}>
                       {themeStore.theme.mode === 'dark' ? t('light_mode') : t('dark_mode')}
                      </span>
                    </Button>
                </Space>
            </div>
        </div>
    );
}
