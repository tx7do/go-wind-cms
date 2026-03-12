'use client';

import React, {useState, useEffect} from 'react';
import {Menu} from 'antd';
import type {MenuProps} from 'antd';

import {XIcon} from '@/plugins/xicon';
import {useI18nRouter} from '@/i18n/helpers/useI18nRouter';
import {useNavigationStore} from '@/store/slices/navigation/hooks';
import {useLanguageChangeEffect} from '@/hooks/useLanguageChangeEffect';

import type {siteservicev1_Navigation, siteservicev1_NavigationItem} from '@/api/generated/app/service/v1';
import styles from './TopNavbar.module.css';

export interface TopNavbarItem {
    key: string;
    label: string;
    icon?: string;
    children?: TopNavbarItem[];
}

interface TopNavbarProps {
    items?: TopNavbarItem[];
    onClick?: (key: number) => void;
}

export default function TopNavbar({onClick}: TopNavbarProps) {
    const navigationStore = useNavigationStore();
    const router = useI18nRouter();
    const [navigationItems, setNavigationItems] = useState<siteservicev1_NavigationItem[]>([]);
    const [isLoading, setIsLoading] = useState(true);

    // 将数据转换为 Ant Design Menu 所需的格式
    const getMenuOptions = (navItems: siteservicev1_NavigationItem[]): MenuProps['items'] => {
        return navItems.map(item => ({
            key: item.id?.toString() || `nav-${item.id}`,
            label: item.title,
            icon: item.icon ? <XIcon name={`carbon:${item.icon}`} size={22}/> : undefined,
            children: item.children && item.children.length > 0
                ? item.children.map((child: siteservicev1_NavigationItem) => ({
                    key: child.id?.toString() || `nav-${child.id}`,
                    label: child.title,
                    icon: child.icon ? <XIcon name={`carbon:${child.icon}`} size={18}/> : undefined,
                }))
                : undefined,
        })) as MenuProps['items'];
    };

    // 加载导航数据
    useEffect(() => {
        async function loadNavigation() {
            try {
                setIsLoading(true);
                const res = await navigationStore.listNavigation({
                    // @ts-expect-error - listNavigation 参数类型推断问题
                    paging: {page: 1, pageSize: 10}
                }) as unknown as { items: siteservicev1_Navigation[]; total: number };

                if (res.items && res.items.length > 0) {
                    const headerNav = res.items.find(nav =>
                        nav.location === 'HEADER' && nav.isActive === true
                    );

                    if (headerNav && headerNav.items && headerNav.items.length > 0) {
                        setNavigationItems(headerNav.items);
                    }
                }
            } catch (error) {
                console.error('[TopNavbar] 加载导航失败:', error);
            } finally {
                setIsLoading(false);
            }
        }

        loadNavigation();
        // eslint-disable-next-line react-hooks/exhaustive-deps
    }, []);

    // 监听语言变化，自动重新加载导航数据
    useLanguageChangeEffect(() => {
        setIsLoading(true);
        navigationStore.listNavigation({
            // @ts-expect-error - listNavigation 参数类型推断问题
            paging: {page: 1, pageSize: 10}
        })
            .then(res => {
                const navRes = res as unknown as { items: siteservicev1_Navigation[]; total: number };
                if (navRes.items && navRes.items.length > 0) {
                    const headerNav = navRes.items.find(nav =>
                        nav.location === 'HEADER' && nav.isActive === true
                    );
                    if (headerNav && headerNav.items && headerNav.items.length > 0) {
                        setNavigationItems(headerNav.items);
                    }
                }
            })
            .catch(error => {
                console.error('[TopNavbar] 重新加载导航失败:', error);
            })
            .finally(() => {
                setIsLoading(false);
            });
    }, {immediate: false, autoCleanup: true});

    // 处理导航点击
    const handleNavigate = (item: siteservicev1_NavigationItem) => {
        if (item.isOpenNewTab) {
            window.open(item.url, '_blank');
        } else {
            if (item.url != null) {
                router.push(item.url);
            }
        }
    };

    const menuItems = getMenuOptions(navigationItems) || [];

    // 根据主题设置 Menu 的 theme
    const [isDark, setIsDark] = useState(false);
    useEffect(() => {
        const updateTheme = () => {
            setIsDark(document.documentElement.classList.contains('dark'));
        };
        updateTheme();
        const observer = new MutationObserver(updateTheme);
        observer.observe(document.documentElement, {attributes: true, attributeFilter: ['class']});
        return () => observer.disconnect();
    }, []);

    return (
        <div className={styles.navbarWrapper} style={{width: '100%', minWidth: 0, flex: 1}}>
            <div className={styles.navbarContent} style={{width: '100%', minWidth: 0, flex: 1}}>
                {/* 导航菜单 */}
                {isLoading ? (
                    <div className={styles.loadingState}>
                        Loading navigation...
                    </div>
                ) : menuItems.length > 0 ? (
                    <Menu
                        mode="horizontal"
                        className={styles.navbarMenu}
                        theme={isDark ? 'dark' : 'light'}
                        items={menuItems}
                        style={{borderBottom: 'none', boxShadow: 'none', width: '100%', minWidth: 0, flex: 1}}
                        onClick={({key, keyPath}) => {
                            // keyPath: [childKey, parentKey] for submenu
                            let item: siteservicev1_NavigationItem | null;
                            if (keyPath.length > 1) {
                                // 子菜单项
                                item = navigationStore.findNavItem(navigationItems, Number(keyPath[0]));
                            } else {
                                item = navigationStore.findNavItem(navigationItems, Number(key));
                            }
                            if (item) handleNavigate(item);
                            onClick?.(Number(key));
                        }}
                    />
                ) : (
                    <div className={styles.loadingState}>
                        No navigation items
                    </div>
                )}
            </div>
        </div>
    );
}
