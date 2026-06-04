import React, {useState, useEffect} from 'react';
import {View, Text, Image} from '@tarojs/components';
import {useTranslations, useLocale} from '@/lib/next-intl-compat';
import logoImage from '@/assets/images/logo.png';
import {XIcon} from '@/plugins/xicon';
import {useI18nRouter} from '@/i18n/helpers/useI18nRouter';
import {useI18n} from '@/i18n';
import {usePreferences} from '@/core/preferences';
import {fetchListNavigations} from '@/api/hooks/navigation';
import {useLanguageChangeEffect} from '@/hooks/useLanguageChangeEffect';
import type {siteservicev1_Navigation, siteservicev1_NavigationItem} from '@/api/generated/app/service/v1';

export default function MobileNav() {
    const t = useTranslations('navbar');
    const menuT = useTranslations('menu');
    const appT = useTranslations('app');
    const brandTitle = appT('title');

    const router = useI18nRouter();
    const locale = useLocale();
    const {theme: themePref, setThemeMode} = usePreferences();
    const {changeLocale} = useI18n();
    const currentMode = themePref.mode;

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
                    const headerNav = res.items.find(nav => nav.location === 'HEADER' && nav.isActive === true);
                    if (headerNav?.items?.length) setNavigationItems(headerNav.items);
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
                    const headerNav = navRes.items.find(nav => nav.location === 'HEADER' && nav.isActive === true);
                    if (headerNav?.items?.length) setNavigationItems(headerNav.items);
                }
            })
            .catch(error => console.error('[MobileNav] 重新加载导航失败:', error));
    }, {immediate: false, autoCleanup: true});

    const handleNavigate = (item: siteservicev1_NavigationItem) => {
        if (item.url) router.push(item.url);
        setOpen(false);
    };

    const handleAction = (action: () => void) => {
        action();
        setOpen(false);
    };

    // ---- 汉堡按钮 ----
    if (!open) {
        return (
            <View
              className='flex items-center justify-center w-[64rpx] h-[64rpx] rounded'
              onClick={() => setOpen(true)}
              hoverClass='tap-active'
            >
                <XIcon name='carbon:menu' size={20} className='text-textSec' />
            </View>
        );
    }

    // ---- 全屏抽屉 ----
    return (
        <View
          className='fixed top-0 left-0 right-0 bottom-0 z-[2000]'
          style={{width: '100vw', height: '100vh'}}
        >
            {/* 半透明遮罩 */}
            <View
              className='absolute top-0 left-0 right-0 bottom-0'
              style={{backgroundColor: 'rgba(0,0,0,0.4)'}}
              onClick={() => setOpen(false)}
            />

            {/* 右侧抽屉面板 */}
            <View
              className='absolute top-0 right-0 bottom-0 bg-cardBg flex flex-col'
              style={{width: '560rpx'}}
            >
                {/* 顶部品牌区 */}
                <View className='flex items-center px-[32rpx] border-b-[1rpx] border-splitLine'
                  style={{height: '120rpx'}}
                >
                    <Image src={logoImage} mode='aspectFit' className='h-[48rpx] w-[48rpx]' />
                    <Text className='text-card-title font-bold text-primary ml-[16rpx]'>{brandTitle}</Text>
                    <View
                      className='flex items-center justify-center w-[64rpx] h-[64rpx] ml-auto'
                      onClick={() => setOpen(false)}
                      hoverClass='tap-active'
                    >
                        <XIcon name='carbon:close' size={20} className='text-textThird' />
                    </View>
                </View>

                {/* 可滚动内容区 */}
                <View className='flex-1 overflow-y-auto py-[16rpx]'>
                    {/* 主导航 */}
                    {navigationItems.length > 0 && (
                        <View className='mb-[16rpx]'>
                            {navigationItems.map((item) => {
                                const itemId = item.id ?? 0;
                                const hasChildren = (item.children?.length ?? 0) > 0;
                                const isExpanded = expandedId === itemId;

                                return (
                                    <View key={itemId}>
                                        <View
                                          className='flex items-center justify-between px-[32rpx]'
                                          style={{height: '88rpx'}}
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
                                                    <View className='mr-[16rpx]'>
                                                        <XIcon name={`carbon:${item.icon}`} size={16} className='text-primary' />
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
                                                      style={{height: '80rpx'}}
                                                      onClick={() => handleNavigate(child)}
                                                      hoverClass='tap-active'
                                                    >
                                                        {child.icon && (
                                                            <View className='mr-[16rpx]'>
                                                                <XIcon name={`carbon:${child.icon}`} size={14} className='text-textThird' />
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
                    <View className='bg-splitLine mx-[32rpx] my-[16rpx]' style={{height: '1rpx'}} />

                    {/* 用户操作 */}
                    <DrawerItem icon='carbon:user' label={menuT('homepage')} onClick={() => handleAction(() => router.push('/user'))} />
                    <DrawerItem icon='carbon:settings' label={menuT('my_profile')} onClick={() => handleAction(() => router.push('/settings'))} />

                    {/* 分隔线 */}
                    <View className='bg-splitLine mx-[32rpx] my-[16rpx]' style={{height: '1rpx'}} />

                    {/* 语言切换 */}
                    <View className='px-[32rpx] py-[8rpx]'>
                        <View className='flex items-center'>
                            <XIcon name='carbon:earth' size={12} className='text-textThird mr-[8rpx]' />
                            <Text className='text-tips text-textThird'>{t('language.title')}</Text>
                        </View>
                    </View>
                    <DrawerItem label='简体中文' onClick={() => handleAction(() => changeLocale('zh-CN'))} active={locale === 'zh-CN'} />
                    <DrawerItem label='English' onClick={() => handleAction(() => changeLocale('en-US'))} active={locale === 'en-US'} />

                    {/* 分隔线 */}
                    <View className='bg-splitLine mx-[32rpx] my-[16rpx]' style={{height: '1rpx'}} />

                    {/* 主题切换 */}
                    <View className='px-[32rpx] py-[8rpx]'>
                        <View className='flex items-center'>
                            <XIcon name={currentMode === 'dark' ? 'carbon:moon' : currentMode === 'light' ? 'carbon:sun' : 'carbon:monitor'} size={12} className='text-textThird mr-[8rpx]' />
                            <Text className='text-tips text-textThird'>{t('theme.title') || 'Theme'}</Text>
                        </View>
                    </View>
                    <DrawerItem icon='carbon:sun' label={t('theme.light')} onClick={() => handleAction(() => setThemeMode('light'))} active={currentMode === 'light'} />
                    <DrawerItem icon='carbon:moon' label={t('theme.dark')} onClick={() => handleAction(() => setThemeMode('dark'))} active={currentMode === 'dark'} />
                    <DrawerItem icon='carbon:monitor' label={t('theme.system')} onClick={() => handleAction(() => setThemeMode('auto'))} active={currentMode === 'auto'} />
                </View>
            </View>
        </View>
    );
}

/** 抽屉内菜单项 */
function DrawerItem({icon, label, onClick, active, destructive}: {
    icon?: string;
    label: React.ReactNode;
    onClick: () => void;
    active?: boolean;
    destructive?: boolean;
}) {
    const bgStyle = active ? {backgroundColor: 'rgba(22,119,255,0.08)'} : {};

    return (
        <View
          className='flex items-center px-[32rpx]'
          style={{height: '88rpx', ...bgStyle}}
          onClick={onClick}
          hoverClass='tap-active'
        >
            {icon && (
                <View className='mr-[16rpx]'>
                    <XIcon
                      name={icon}
                      size={16}
                      className={destructive ? 'text-danger' : active ? 'text-primary' : 'text-textSec'}
                    />
                </View>
            )}
            <Text
              className={`text-body ${destructive ? 'text-danger' : active ? 'text-primary font-bold' : 'text-textMain'}`}
            >
                {label}
            </Text>
        </View>
    );
}
