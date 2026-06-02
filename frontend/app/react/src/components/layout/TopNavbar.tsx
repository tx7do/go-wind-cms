'use client';

import React, {useState, useEffect, useRef} from 'react';

import {XIcon} from '@/plugins/xicon';
import {useI18nRouter} from '@/i18n/helpers/useI18nRouter';
import {fetchListNavigations} from '@/api/hooks/navigation';
import {useLanguageChangeEffect} from '@/hooks/useLanguageChangeEffect';
import {usePreferences} from '@/core/preferences';
import {cn} from '@/lib/utils';

import type {siteservicev1_Navigation, siteservicev1_NavigationItem} from '@/api/generated/app/service/v1';

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
    const [openMenuId, setOpenMenuId] = useState<number | null>(null);
    const closeTimerRef = useRef<ReturnType<typeof setTimeout> | null>(null);

    const loadNav = async () => {
        try {
            setIsLoading(true);
            const res = await fetchListNavigations({
                paging: {page: 1, pageSize: 10}
            }) as unknown as { items: siteservicev1_Navigation[]; total: number };

            if (res.items?.length) {
                const headerNav = res.items.find(nav =>
                    nav.location === 'HEADER' && nav.isActive === true
                );
                const items = headerNav?.items;
                if (items && items.length > 0) {
                    setNavigationItems(items);
                }
            }
        } catch (error) {
            console.error('[TopNavbar] 加载导航失败:', error);
        } finally {
            setIsLoading(false);
        }
    };

    useEffect(() => {
        loadNav();

    }, []);

    useLanguageChangeEffect(() => {
        setIsLoading(true);
        fetchListNavigations({paging: {page: 1, pageSize: 10}})
            .then(res => {
                const navRes = res as unknown as { items: siteservicev1_Navigation[]; total: number };
                if (navRes.items?.length) {
                    const headerNav = navRes.items.find(nav =>
                        nav.location === 'HEADER' && nav.isActive === true
                    );
                    const items = headerNav?.items;
                    if (items && items.length > 0) {
                        setNavigationItems(items);
                    }
                }
            })
            .catch(error => console.error('[TopNavbar] 重新加载导航失败:', error))
            .finally(() => setIsLoading(false));
    }, {immediate: false, autoCleanup: true});

    // 清理定时器
    useEffect(() => {
        return () => {
            if (closeTimerRef.current) clearTimeout(closeTimerRef.current);
        };
    }, []);

    const handleNavigate = (item: siteservicev1_NavigationItem) => {
        if (item.isOpenNewTab) {
            window.open(item.url, '_blank');
        } else if (item.url != null) {
            router.push(item.url);
        }
    };

    const handleMouseEnter = (id: number) => {
        if (closeTimerRef.current) {
            clearTimeout(closeTimerRef.current);
            closeTimerRef.current = null;
        }
        setOpenMenuId(id);
    };

    const handleMouseLeave = () => {
        closeTimerRef.current = setTimeout(() => setOpenMenuId(null), 150);
    };

    usePreferences(); // keep theme reactivity

    if (isLoading) {
        return (
            <div className="flex h-full items-center text-sm text-muted-foreground">
                Loading…
            </div>
        );
    }

    if (navigationItems.length === 0) {
        return null;
    }

    return (
        <nav className="flex h-full items-center gap-0.5">
            {navigationItems.map((item) => {
                const itemId = item.id ?? 0;
                const hasChildren = (item.children?.length ?? 0) > 0;
                const isOpen = openMenuId === itemId;

                return (
                    <div
                        key={itemId}
                        className="relative"
                        onMouseEnter={() => handleMouseEnter(itemId)}
                        onMouseLeave={handleMouseLeave}
                    >
                        <button
                            type="button"
                            className={cn(
                                'inline-flex h-9 items-center gap-1.5 rounded-lg px-3 text-sm font-medium',
                                'text-foreground/80 transition-all duration-200',
                                'hover:bg-primary/10 hover:text-primary',
                                hasChildren && isOpen && 'bg-primary/10 text-primary',
                                'max-md:px-2',
                            )}
                            onClick={() => {
                                if (!hasChildren) {
                                    handleNavigate(item);
                                    onClick?.(Number(itemId));
                                }
                            }}
                        >
                            {item.icon && <XIcon name={`carbon:${item.icon}`} size={16}/>}
                            <span className="whitespace-nowrap">{item.title}</span>
                            {hasChildren && (
                                <XIcon name="carbon:chevron-down" size={12} className={cn('transition-transform duration-200', isOpen && 'rotate-180')}/>
                            )}
                        </button>

                        {hasChildren && isOpen && (
                            <div
                                className={cn(
                                    'absolute left-0 top-full z-1001 mt-1.5',
                                    'min-w-50 rounded-lg border border-border bg-popover p-1.5',
                                    'shadow-lg shadow-black/5',
                                    'animate-in fade-in-0 zoom-in-95 duration-150',
                                )}
                                onMouseEnter={() => handleMouseEnter(itemId)}
                                onMouseLeave={handleMouseLeave}
                            >
                                {item.children!.map((child: siteservicev1_NavigationItem) => (
                                    <button
                                        key={child.id?.toString()}
                                        type="button"
                                        className="flex w-full cursor-pointer items-center gap-2 rounded-md px-2.5 py-2 text-left text-sm text-foreground/80 transition-colors duration-150 hover:bg-primary/10 hover:text-primary"
                                        onClick={() => {
                                            handleNavigate(child);
                                            onClick?.(Number(child.id));
                                            setOpenMenuId(null);
                                        }}
                                    >
                                        {child.icon && <XIcon name={`carbon:${child.icon}`} size={14}/>}
                                        <span className="truncate">{child.title}</span>
                                    </button>
                                ))}
                            </div>
                        )}
                    </div>
                );
            })}
        </nav>
    );
}
