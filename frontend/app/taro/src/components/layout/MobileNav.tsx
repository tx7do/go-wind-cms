import React, {useState, useEffect, useCallback, useRef, createContext, useContext} from 'react';
import {View, Text, Image} from '@tarojs/components';
import {useTranslations, useLocale} from '@/lib/next-intl-compat';
import logoImage from '@/assets/images/logo.png';
import {XIcon} from '@/plugins/xicon';
import {useI18nRouter} from '@/i18n/helpers/useI18nRouter';
import {useI18n} from '@/i18n';
import {usePreferences} from '@/core/preferences';
import {fetchListNavigations} from '@/api/hooks/navigation';
import type {siteservicev1_Navigation, siteservicev1_NavigationItem} from '@/api/generated/app/service/v1';

/** 全局抽屉状态上下文 */
const DrawerCtx = createContext<{
    open: boolean;
    setOpen: (v: boolean) => void;
}>({open: false, setOpen: () => {}});

/** 供 Layout 使用的 Provider + 抽屉渲染 */
export function MobileNavProvider({children}: {children: React.ReactNode}) {
    const [open, setOpen] = useState(false);
    return (
        <DrawerCtx.Provider value={{open, setOpen}}>
            {children}
            {open && <MobileNavDrawer />}
        </DrawerCtx.Provider>
    );
}

/** Header 中的汉堡触发按钮 */
export function MobileNavTrigger() {
    const {open, setOpen} = useContext(DrawerCtx);
    return (
        <View
            className='flex items-center justify-center w-[64rpx] h-[64rpx] rounded'
            onClick={() => setOpen(!open)}
            hoverClass='tap-active'
        >
            <XIcon name={open ? 'carbon:close' : 'carbon:menu'} size={20} className='text-textSec' />
        </View>
    );
}

/** 抽屉面板（在 Layout 顶层渲染，不受 Header fixed 嵌套影响） */
function MobileNavDrawer() {
    const {open, setOpen} = useContext(DrawerCtx);
    const t = useTranslations('navbar');
    const menuT = useTranslations('menu');
    const appT = useTranslations('app');
    const brandTitle = appT('title');

    const router = useI18nRouter();
    const locale = useLocale();
    const {theme: themePref, setThemeMode} = usePreferences();
    const {changeLocale} = useI18n();
    const currentMode = themePref.mode;

    const [navigationItems, setNavigationItems] = useState<siteservicev1_NavigationItem[]>([]);
    const [expandedId, setExpandedId] = useState<number | null>(null);

    const close = () => setOpen(false);

    const loadNavigations = useCallback(async () => {
        try {
            const res = await fetchListNavigations({paging: {page: 1, pageSize: 10}}) as unknown as { items: siteservicev1_Navigation[]; total: number };
            if (res.items?.length) {
                const headerNav = res.items.find(nav => nav.location === 'HEADER' && nav.isActive === true);
                if (headerNav?.items?.length) setNavigationItems(headerNav.items);
            }
        } catch (error) {
            console.error('[MobileNav] 加载导航失败:', error);
        }
    }, []);

    // 首次加载
    useEffect(() => {
        loadNavigations();
    }, [loadNavigations]);

    // locale 变化时重新加载导航（用 ref 跳过首次）
    const prevLocaleRef = useRef(locale);
    useEffect(() => {
        if (prevLocaleRef.current !== locale) {
            prevLocaleRef.current = locale;
            loadNavigations();
        }
    }, [locale, loadNavigations]);

    const handleNavigate = (item: siteservicev1_NavigationItem) => {
        if (item.url) router.push(item.url);
        close();
    };

    const handleAction = (action: () => void) => {
        action();
        close();
    };

    return (
        <View
            className='fixed top-0 left-0 right-0 bottom-0 z-[2000]'
            style={{width: '100vw', height: '100vh'}}
        >
            {/* 半透明遮罩 */}
            <View
                className='fixed top-0 left-0 right-0 bottom-0'
                style={{backgroundColor: 'rgba(0,0,0,0.5)', zIndex: 1}}
                onClick={close}
            />

            {/* 右侧抽屉面板 */}
            <View
                className='fixed top-0 bottom-0 bg-cardBg flex flex-col'
                style={{right: 0, width: '560rpx', zIndex: 2}}
            >
                {/* 顶部关闭按钮区 */}
                <View
                    className='flex items-center justify-end border-b-[1rpx] border-splitLine'
                    style={{height: '128rpx', padding: '16rpx 32rpx 16rpx 32rpx'}}
                >
                    <View
                        className='flex items-center justify-center rounded-full'
                        style={{width: '88rpx', height: '88rpx', backgroundColor: 'rgba(0,0,0,0.04)'}}
                        onClick={close}
                        hoverClass='tap-active'
                    >
                        <XIcon name='carbon:close' size={22} className='text-textMain' />
                    </View>
                </View>

                {/* 可滚动内容区 */}
                <View className='flex-1 overflow-y-auto'>
                    {/* 主导航 */}
                    {navigationItems.length > 0 && (
                        <View className='py-[8rpx]'>
                            {navigationItems.map((item) => {
                                const itemId = item.id ?? 0;
                                const hasChildren = (item.children?.length ?? 0) > 0;
                                const isExpanded = expandedId === itemId;

                                return (
                                    <View key={itemId}>
                                        <View
                                            className='flex items-center justify-between px-[32rpx]'
                                            style={{height: '96rpx'}}
                                            onClick={() => {
                                                if (hasChildren) {
                                                    setExpandedId(isExpanded ? null : itemId);
                                                } else {
                                                    handleNavigate(item);
                                                }
                                            }}
                                            hoverClass='tap-active'
                                        >
                                            <View className='flex items-center'>
                                                {item.icon && (
                                                    <View style={{marginRight: '16rpx'}}>
                                                        <XIcon name={`carbon:${item.icon}`} size={18} className='text-primary' />
                                                    </View>
                                                )}
                                                <Text className='text-body text-textMain'>{item.title}</Text>
                                            </View>
                                            {hasChildren && (
                                                <XIcon
                                                    name='carbon:chevron-down'
                                                    size={16}
                                                    className='text-textThird'
                                                    style={{transform: isExpanded ? 'rotate(180deg)' : 'rotate(0deg)', transition: 'transform 0.2s'}}
                                                />
                                            )}
                                        </View>

                                        {/* 子导航 */}
                                        {hasChildren && isExpanded && (
                                            <View
                                                className='border-l-[2rpx] border-splitLine'
                                                style={{marginLeft: '32rpx'}}
                                            >
                                                {item.children!.map((child: siteservicev1_NavigationItem) => (
                                                    <View
                                                        key={child.id?.toString()}
                                                        className='flex items-center px-[24rpx]'
                                                        style={{height: '84rpx'}}
                                                        onClick={() => handleNavigate(child)}
                                                        hoverClass='tap-active'
                                                    >
                                                        {child.icon && (
                                                            <View style={{marginRight: '16rpx'}}>
                                                                <XIcon name={`carbon:${child.icon}`} size={16} className='text-textThird' />
                                                            </View>
                                                        )}
                                                        <Text className='text-desc text-textSec'>{child.title}</Text>
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
                    <View className='bg-splitLine' style={{height: '1rpx', margin: '0 32rpx'}} />

                    {/* 用户操作 */}
                    <View className='py-[8rpx]'>
                        <DrawerItem icon='carbon:user' label={menuT('homepage')} onClick={() => handleAction(() => router.push('/user'))} />
                        <DrawerItem icon='carbon:settings' label={menuT('my_profile')} onClick={() => handleAction(() => router.push('/settings'))} />
                    </View>

                    {/* 分隔线 */}
                    <View className='bg-splitLine' style={{height: '1rpx', margin: '0 32rpx'}} />

                    {/* 语言切换 */}
                    <View className='py-[8rpx]'>
                        <View className='flex items-center px-[32rpx]' style={{height: '56rpx'}}>
                            <XIcon name='carbon:earth' size={12} className='text-textThird' />
                            <Text className='text-xs text-textThird' style={{marginLeft: '6rpx'}}>语言</Text>
                        </View>
                        <DrawerItem label='简体中文' onClick={() => handleAction(() => changeLocale('zh-CN'))} active={locale === 'zh-CN'} />
                        <DrawerItem label='English' onClick={() => handleAction(() => changeLocale('en-US'))} active={locale === 'en-US'} />
                    </View>

                    {/* 分隔线 */}
                    <View className='bg-splitLine' style={{height: '1rpx', margin: '0 32rpx'}} />

                    {/* 主题切换 */}
                    <View className='py-[8rpx]'>
                        <View className='flex items-center px-[32rpx]' style={{height: '56rpx'}}>
                            <XIcon
                                name={currentMode === 'dark' ? 'carbon:moon' : currentMode === 'light' ? 'carbon:sun' : 'carbon:monitor'}
                                size={12}
                                className='text-textThird'
                            />
                            <Text className='text-xs text-textThird' style={{marginLeft: '6rpx'}}>主题模式</Text>
                        </View>
                        <DrawerItem icon='carbon:sun' label={t('theme.light')} onClick={() => handleAction(() => setThemeMode('light'))} active={currentMode === 'light'} />
                        <DrawerItem icon='carbon:moon' label={t('theme.dark')} onClick={() => handleAction(() => setThemeMode('dark'))} active={currentMode === 'dark'} />
                        <DrawerItem icon='carbon:monitor' label={t('theme.system')} onClick={() => handleAction(() => setThemeMode('auto'))} active={currentMode === 'auto'} />
                    </View>
                </View>
            </View>
        </View>
    );
}

/** 抽屉内菜单项 - 增强选中状态视觉反馈 */
function DrawerItem({icon, label, onClick, active, destructive}: {
    icon?: string;
    label: React.ReactNode;
    onClick: () => void;
    active?: boolean;
    destructive?: boolean;
}) {
    // 选中状态：浅蓝色背景 + 主色文字 + 加粗
    const bgStyle = active ? {backgroundColor: 'rgba(22,119,255,0.08)'} : {};

    return (
        <View
            className='flex items-center px-[32rpx]'
            style={{height: '88rpx', ...bgStyle}}
            onClick={onClick}
            hoverClass='tap-active'
        >
            {icon && (
                <View style={{marginRight: '16rpx'}}>
                    <XIcon
                        name={icon}
                        size={18}
                        className={destructive ? 'text-danger' : active ? 'text-primary' : 'text-textSec'}
                    />
                </View>
            )}
            <Text
                className={`text-body flex-1 ${destructive ? 'text-danger' : active ? 'text-primary font-bold' : 'text-textMain'}`}
            >
                {label}
            </Text>
            {/* 选中状态：右侧显示勾选图标 */}
            {active && !destructive && (
                <XIcon name='carbon:checkmark' size={18} className='text-primary' />
            )}
        </View>
    );
}

/** 默认导出保持向后兼容 */
export default MobileNavTrigger;
