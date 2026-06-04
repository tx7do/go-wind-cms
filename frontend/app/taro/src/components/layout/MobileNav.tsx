import React, {useState, useEffect} from 'react';
import {View, Image} from '@tarojs/components';
import {useTranslations, useLocale} from '@/lib/next-intl-compat';
import logoImage from '@/assets/images/logo.png';

import {Button} from '@/components/ui/button';
import {
    Sheet,
    SheetContent,
    SheetHeader,
    SheetTitle,
} from '@/components/ui/sheet';

import {XIcon} from '@/plugins/xicon';
import {useI18nRouter} from '@/i18n/helpers/useI18nRouter';
import {useI18n} from '@/i18n';
import {usePreferences} from '@/core/preferences';
import {fetchListNavigations} from '@/api/hooks/navigation';
import {useLanguageChangeEffect} from '@/hooks/useLanguageChangeEffect';
import {cn} from '@/lib/utils';

import {useAccessStore} from '@/store/core/access/store';
import {useAuth} from '@/api/hooks/auth';

import type {siteservicev1_Navigation, siteservicev1_NavigationItem} from '@/api/generated/app/service/v1';

export default function MobileNav() {
    const t = useTranslations('navbar');
    const appT = useTranslations('app');
    const menuT = useTranslations('menu');
    const brandTitle = appT('title');

    const router = useI18nRouter();
    const locale = useLocale();
    const {theme: themePref, setThemeMode} = usePreferences();
    const {changeLocale} = useI18n();
    const accessStore = useAccessStore();
    const auth = useAuth();
    const currentMode = themePref.mode;

    const accessToken = accessStore.accessToken;
    const isLogin = !!accessToken && !accessStore.loginExpired;

    const [open, setOpen] = useState(false);
    const [navigationItems, setNavigationItems] = useState<siteservicev1_NavigationItem[]>([]);
    const [expandedId, setExpandedId] = useState<number | null>(null);

    useEffect(() => {
        let cancelled = false;
        (async () => {
            try {
                const res = await fetchListNavigations({
                    paging: {page: 1, pageSize: 10}
                }) as unknown as { items: siteservicev1_Navigation[]; total: number };

                if (cancelled) return;

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
                console.error('[MobileNav] 加载导航失败:', error);
            }
        })();
        return () => { cancelled = true; };
    }, []);

    useLanguageChangeEffect(() => {
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
            .catch(error => console.error('[MobileNav] 重新加载导航失败:', error));
    }, {immediate: false, autoCleanup: true});

    const handleNavigate = (item: siteservicev1_NavigationItem) => {
        if (item.isOpenNewTab) {
            // Taro 中打开新页面
            // 小程序不支持 window.open，使用 navigateTo
            if (item.url) router.push(item.url);
        } else if (item.url != null) {
            router.push(item.url);
        }
        setOpen(false);
    };

    const handleAction = (action: () => void) => {
        action();
        setOpen(false);
    };

    const themeIconName = currentMode === 'dark' ? 'carbon:moon' :
        currentMode === 'light' ? 'carbon:sun' : 'carbon:monitor';

    return (
        <>
            {/* 汉堡按钮 */}
            <Button
              variant='ghost'
              size='icon'
              onClick={() => setOpen(true)}
              className='md:hidden'
            >
                <XIcon name='carbon:menu' size={20} />
            </Button>

            {/* 抽屉 */}
            <Sheet open={open} onOpenChange={setOpen}>
                <SheetContent side='right' className='flex w-full flex-col p-0 sm:max-w-sm'>
                    {/* 顶部品牌区 */}
                    <SheetHeader className='border-b border-border px-5 py-4'>
                        <View className='flex items-center gap-2'>
                            <Image src={logoImage} mode='aspectFit' className='h-8 w-8' />
                            <SheetTitle className='text-base font-bold text-primary'>
                                {brandTitle}
                            </SheetTitle>
                        </View>
                    </SheetHeader>

                    {/* 滚动内容区 */}
                    <View className='flex-1 overflow-y-auto px-3 py-3'>
                        {/* 主导航 */}
                        {navigationItems.length > 0 && (
                            <View className='mb-3 space-y-0.5'>
                                {navigationItems.map((item) => {
                                    const itemId = item.id ?? 0;
                                    const hasChildren = (item.children?.length ?? 0) > 0;
                                    const isExpanded = expandedId === itemId;

                                    return (
                                        <View key={itemId}>
                                            <View
                                              className={cn(
                                                    'flex w-full items-center justify-between gap-2 rounded-lg px-3 py-2.5 text-sm font-medium',
                                                    'text-foreground/80 transition-colors',
                                                    'hover:bg-primary/10 hover:text-primary',
                                                )}
                                              onClick={() => {
                                                    if (hasChildren) {
                                                        setExpandedId(isExpanded ? null : itemId);
                                                    } else {
                                                        handleNavigate(item);
                                                    }
                                                }}
                                            >
                                                <View className='flex items-center gap-2'>
                                                    {item.icon && <XIcon name={`carbon:${item.icon}`} size={16} />}
                                                    <View>{item.title}</View>
                                                </View>
                                                {hasChildren && (
                                                    <XIcon name='carbon:chevron-down' size={16}
                                                      className={cn('transition-transform', isExpanded && 'rotate-180')}
                                                    />
                                                )}
                                            </View>

                                            {/* 子菜单 */}
                                            {hasChildren && isExpanded && (
                                                <View className='ml-3 mt-0.5 space-y-0.5 border-l border-border pl-3'>
                                                    {item.children!.map((child: siteservicev1_NavigationItem) => (
                                                        <View
                                                          key={child.id?.toString()}
                                                          className='flex w-full items-center gap-2 rounded-md px-3 py-2 text-sm text-foreground/70 transition-colors hover:bg-primary/10 hover:text-primary'
                                                          onClick={() => handleNavigate(child)}
                                                        >
                                                            {child.icon && <XIcon name={`carbon:${child.icon}`} size={14} />}
                                                            <View>{child.title}</View>
                                                        </View>
                                                    ))}
                                                </View>
                                            )}
                                        </View>
                                    );
                                })}
                            </View>
                        )}

                        {/* 分隔线 */}
                        <View className='my-2 border-t border-border' />

                        {/* 用户操作 */}
                        <View className='space-y-0.5'>
                            {isLogin ? (
                                <>
                                    <SheetNavItem
                                      icon={<XIcon name='carbon:home' size={16} />}
                                      label={menuT('homepage')}
                                      onClick={() => handleAction(() => router.push('/user'))}
                                    />
                                    <SheetNavItem
                                      icon={<XIcon name='carbon:user' size={16} />}
                                      label={menuT('my_profile')}
                                      onClick={() => handleAction(() => router.push('/settings'))}
                                    />
                                    <SheetNavItem
                                      icon={<XIcon name='carbon:logout' size={16} />}
                                      label={menuT('logout')}
                                      onClick={() => handleAction(async () => {
                                            await auth.logout();
                                        })}
                                      destructive
                                    />
                                </>
                            ) : (
                                <>
                                    <SheetNavItem
                                      icon={<XIcon name='carbon:user' size={16} />}
                                      label={t('user.login')}
                                      onClick={() => handleAction(() => router.push('/login'))}
                                    />
                                    <SheetNavItem
                                      icon={<XIcon name='carbon:user' size={16} />}
                                      label={t('user.register')}
                                      onClick={() => handleAction(() => router.push('/register'))}
                                    />
                                </>
                            )}
                        </View>

                        {/* 分隔线 */}
                        <View className='my-2 border-t border-border' />

                        {/* 语言切换 */}
                        <View className='mb-2 px-3 py-1 text-xs font-semibold text-muted-foreground'>
                            <XIcon name='carbon:earth' size={12} className='mr-1' />
                            {t('language.title')}
                        </View>
                        <View className='space-y-0.5'>
                            <SheetNavItem
                              label='简体中文'
                              onClick={() => handleAction(() => changeLocale('zh-CN'))}
                              active={locale === 'zh-CN'}
                            />
                            <SheetNavItem
                              label='English'
                              onClick={() => handleAction(() => changeLocale('en-US'))}
                              active={locale === 'en-US'}
                            />
                        </View>

                        {/* 分隔线 */}
                        <View className='my-2 border-t border-border' />

                        {/* 主题切换 */}
                        <View className='mb-2 px-3 py-1 text-xs font-semibold text-muted-foreground'>
                            <XIcon name={themeIconName} size={12} />
                            <View className='ml-1'>{t('theme.title') || 'Theme'}</View>
                        </View>
                        <View className='space-y-0.5'>
                            <SheetNavItem
                              icon={<XIcon name='carbon:sun' size={16} />}
                              label={t('theme.light')}
                              onClick={() => handleAction(() => setThemeMode('light'))}
                              active={currentMode === 'light'}
                            />
                            <SheetNavItem
                              icon={<XIcon name='carbon:moon' size={16} />}
                              label={t('theme.dark')}
                              onClick={() => handleAction(() => setThemeMode('dark'))}
                              active={currentMode === 'dark'}
                            />
                            <SheetNavItem
                              icon={<XIcon name='carbon:monitor' size={16} />}
                              label={t('theme.system')}
                              onClick={() => handleAction(() => setThemeMode('auto'))}
                              active={currentMode === 'auto'}
                            />
                        </View>
                    </View>
                </SheetContent>
            </Sheet>
        </>
    );
}

/** 抽屉内的导航项 */
function SheetNavItem({
    icon,
    label,
    onClick,
    active,
    destructive,
}: {
    icon?: React.ReactNode;
    label: React.ReactNode;
    onClick: () => void;
    active?: boolean;
    destructive?: boolean;
}) {
    return (
        <View
          className={cn(
                'flex w-full items-center gap-2 rounded-lg px-3 py-2.5 text-sm font-medium transition-colors',
                active
                    ? 'bg-primary/10 text-primary'
                    : destructive
                        ? 'text-destructive hover:bg-destructive/10'
                        : 'text-foreground/80 hover:bg-primary/10 hover:text-primary',
            )}
          onClick={onClick}
        >
            {icon}
            <View>{label}</View>
        </View>
    );
}
