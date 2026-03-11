'use client';

import React, {useState, useEffect} from 'react';
import {Menu, Tabs} from 'antd';

import {XIcon} from '@/plugins/xicon';

import type {siteservicev1_Navigation, siteservicev1_NavigationItem} from '@/api/generated/app/service/v1';
import {useNavigationStore} from '@/store/slices/navigation/hooks';
import {useLanguageStore} from '@/store/core/language/hooks';
import {useLanguageChangeEffect} from '@/hooks/useLanguageChangeEffect';

import type {TopNavBarTabItem} from './types';
import TopNavbarTab from './TopNavbarTab';
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

// 模拟左侧和右侧的 Tab 列表（实际应该从 store 获取）
const leftTabList: TopNavBarTabItem[] = [];
const rightTabList: TopNavBarTabItem[] = [];

export default function TopNavbar({onClick}: TopNavbarProps) {
    const navigationStore = useNavigationStore();
    const languageStore = useLanguageStore();
    const [navigationItems, setNavigationItems] = useState<siteservicev1_NavigationItem[]>([]);
    const [activeOverlay, setActiveOverlay] = useState<string | null>(null);
    const [isLoading, setIsLoading] = useState(true);

    // 将数据转换为 Ant Design Menu 所需的格式
    const getMenuOptions = (navItems: siteservicev1_NavigationItem[]) => {
        return navItems.map(item => ({
            key: item.id,
            label: item.title,
            icon: item.icon ? <XIcon name={`carbon:${item.icon}`} size={22}/> : undefined,
            children: item.children && item.children.length > 0
                ? item.children.map((child: siteservicev1_NavigationItem) => ({
                    key: child.id,
                    label: child.title,
                    icon: child.icon ? <XIcon name={`carbon:${child.icon}`} size={18}/> : undefined,
                }))
                : undefined,
        }));
    };

    // 隐藏悬浮层
    const hideOverlay = () => {
        setActiveOverlay(null);
    };

    // 加载导航数据
    useEffect(() => {
        async function loadNavigation() {
            try {
                setIsLoading(true);

                // 调用 API 获取所有导航数据
                const res = await navigationStore.listNavigation(
                    {page: 1, pageSize: 10}
                ) as unknown as { items: siteservicev1_Navigation[]; total: number };

                if (res.items && res.items.length > 0) {
                    res.items.forEach((nav, index) => {
                        console.log(`[${index}] location: ${nav.location}, isActive: ${nav.isActive}, items count: ${nav.items?.length}`);
                        console.log(`[${index}] items:`, nav.items);
                    });

                    // 在前端过滤出 HEADER 位置且激活的导航
                    const headerNav = res.items.find(nav =>
                        nav.location === 'HEADER' && nav.isActive === true
                    );

                    if (headerNav && headerNav.items && headerNav.items.length > 0) {
                        setNavigationItems(headerNav.items);
                    }
                }
            } catch (error) {
                console.error('❌ 加载导航失败:', error);
            } finally {
                setIsLoading(false);
            }
        }

        loadNavigation();
        // eslint-disable-next-line react-hooks/exhaustive-deps
    }, []); // 移除 languageStore.language.locale 依赖，改用 useLanguageChangeEffect

    // 监听语言变化，自动重新加载导航数据
    useLanguageChangeEffect(() => {
        console.log('[TopNavbar] 语言变化，重新加载导航数据');
        // 重新加载导航数据的逻辑已经在 useEffect 中，这里不需要额外操作
        // 因为 useEffect 的依赖是空的，所以不会自动重新执行
        // 我们需要手动触发重新加载
        setIsLoading(true);
        navigationStore.listNavigation({page: 1, pageSize: 10})
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
                console.error('❌ 重新加载导航失败:', error);
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
                // TODO: 使用 router.push 时保留 history state
                window.location.href = item.url;
            }
        }
    };

    const menuItems = getMenuOptions(navigationItems);

    console.log('=== Menu 渲染 ===');
    console.log('navigationItems:', navigationItems);
    console.log('menuItems:', menuItems);
    console.log('isLoading:', isLoading);

    return (
        <div className={styles.navbarWrapper}>
            <div className={styles.navbarLeft}>
                {/* 左侧导航菜单 */}
                {isLoading ? (
                    <div className={styles.loadingState}>
                        Loading navigation...
                    </div>
                ) : navigationItems.length > 0 ? (
                    <Menu
                        mode="horizontal"
                        className={styles.navbarMenu}
                        items={menuItems}
                        onClick={({key}) => {
                            const item = navigationStore.findNavItem(navigationItems, Number(key));
                            if (item) handleNavigate(item);
                            onClick?.(Number(key));
                        }}
                    />
                ) : (
                    <div className={styles.loadingState}>
                        No navigation items
                    </div>
                )}

                {/* 原有的 Tabs */}
                {leftTabList.length > 0 && (
                    <Tabs
                        activeKey={activeOverlay || undefined}
                        type="card"
                        trigger="hover"
                        animated
                        onTabClick={(key) => setActiveOverlay(key)}
                        onMouseLeave={hideOverlay}
                        items={leftTabList.map(item => ({
                            key: item.key,
                            label: <TopNavbarTab data={item}/>,
                            children: (
                                <div className={styles.overlayMask}>
                                    <div className={styles.overlay}>
                                        <item.component/>
                                    </div>
                                </div>
                            ),
                        }))}
                    />
                )}
            </div>

            <div className={styles.navbarRight}>
                {rightTabList.length > 0 && (
                    <Tabs
                        activeKey={activeOverlay || undefined}
                        type="card"
                        trigger="hover"
                        animated
                        onTabClick={(key) => setActiveOverlay(key)}
                        onMouseLeave={hideOverlay}
                        items={rightTabList.map(item => ({
                            key: item.key,
                            label: <TopNavbarTab data={item}/>,
                            children: (
                                <div className={styles.overlayMask}>
                                    <div className={styles.overlay}>
                                        <item.component/>
                                    </div>
                                </div>
                            ),
                        }))}
                    />
                )}
            </div>
        </div>
    );
}
