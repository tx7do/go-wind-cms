'use client';

import React, {useState, useEffect} from 'react';

import {XIcon} from '@/plugins/xicon';
import {useI18nRouter} from '@/i18n/helpers/useI18nRouter';
import {fetchListNavigations, findNavItem} from '@/api/hooks/navigation';
import {useLanguageChangeEffect} from '@/hooks/useLanguageChangeEffect';
import {usePreferences} from '@/core/preferences';

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
    const router = useI18nRouter();
    const [navigationItems, setNavigationItems] = useState<siteservicev1_NavigationItem[]>([]);
    const [isLoading, setIsLoading] = useState(true);
    const [hoveredItem, setHoveredItem] = useState<number | null>(null);

    useEffect(() => {
        async function loadNavigation() {
            try {
                setIsLoading(true);
                const res = await fetchListNavigations({
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

    useLanguageChangeEffect(() => {
        setIsLoading(true);
        fetchListNavigations({paging: {page: 1, pageSize: 10}})
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
            .catch(error => console.error('[TopNavbar] 重新加载导航失败:', error))
            .finally(() => setIsLoading(false));
    }, {immediate: false, autoCleanup: true});

    const handleNavigate = (item: siteservicev1_NavigationItem) => {
        if (item.isOpenNewTab) {
            window.open(item.url, '_blank');
        } else {
            if (item.url != null) {
                router.push(item.url);
            }
        }
    };

    usePreferences(); // keep theme reactivity

    return (
        <div className={styles.navbarWrapper} style={{width: '100%', minWidth: 0, flex: 1}}>
            <div className={styles.navbarContent} style={{width: '100%', minWidth: 0, flex: 1}}>
                {isLoading ? (
                    <div className={styles.loadingState}>Loading navigation...</div>
                ) : navigationItems.length > 0 ? (
                    <nav className="flex w-full items-center gap-1">
                        {navigationItems.map((item) => {
                            const itemId = item.id?.toString() || '';
                            const hasChildren = item.children && item.children.length > 0;
                            return (
                                <div
                                    key={itemId}
                                    className="relative"
                                    onMouseEnter={() => setHoveredItem(Number(itemId))}
                                    onMouseLeave={() => setHoveredItem(null)}
                                >
                                    <button
                                        className="inline-flex h-9 items-center gap-1 rounded-md px-3 py-2 text-sm font-medium transition-colors hover:bg-accent hover:text-accent-foreground"
                                        onClick={() => {
                                            if (!hasChildren) {
                                                handleNavigate(item);
                                                onClick?.(Number(itemId));
                                            }
                                        }}
                                    >
                                        {item.icon && <XIcon name={`carbon:${item.icon}`} size={18}/>}
                                        <span>{item.title}</span>
                                    </button>
                                    {hasChildren && hoveredItem === Number(itemId) && (
                                        <div className="absolute left-0 top-full z-50 min-w-[8rem] rounded-md border bg-popover p-1 text-popover-foreground shadow-md">
                                            {item.children!.map((child: siteservicev1_NavigationItem) => (
                                                <button
                                                    key={child.id?.toString()}
                                                    className="flex w-full cursor-pointer select-none items-center gap-2 rounded-sm px-2 py-1.5 text-sm outline-none transition-colors hover:bg-accent hover:text-accent-foreground"
                                                    onClick={() => {
                                                        handleNavigate(child);
                                                        onClick?.(Number(child.id));
                                                    }}
                                                >
                                                    {child.icon && <XIcon name={`carbon:${child.icon}`} size={16}/>}
                                                    <span>{child.title}</span>
                                                </button>
                                            ))}
                                        </div>
                                    )}
                                </div>
                            );
                        })}
                    </nav>
                ) : (
                    <div className={styles.loadingState}>No navigation items</div>
                )}
            </div>
        </div>
    );
}
