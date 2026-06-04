import {View, Text} from '@tarojs/components';
import {useState, useEffect} from 'react';
import {XIcon} from '@/plugins/xicon';
import {useI18nRouter} from '@/i18n/helpers/useI18nRouter';
import {fetchListNavigations} from '@/api/hooks/navigation';
import {useLanguageChangeEffect} from '@/hooks/useLanguageChangeEffect';
import {Skeleton} from '@/components/ui/skeleton';
import type {siteservicev1_Navigation, siteservicev1_NavigationItem} from '@/api/generated/app/service/v1';

export default function TopNavbar() {
    const router = useI18nRouter();
    const [navigationItems, setNavigationItems] = useState<siteservicev1_NavigationItem[]>([]);
    const [isLoading, setIsLoading] = useState(true);
    const [expandedId, setExpandedId] = useState<number | null>(null);

    const loadNav = async () => {
        try {
            setIsLoading(true);
            const res = await fetchListNavigations({
                paging: {page: 1, pageSize: 10}
            }) as unknown as { items: siteservicev1_Navigation[]; total: number };
            if (res.items?.length) {
                const headerNav = res.items.find(nav => nav.location === 'HEADER' && nav.isActive === true);
                if (headerNav?.items?.length) setNavigationItems(headerNav.items);
            }
        } catch (error) {
            console.error('[TopNavbar] 加载导航失败:', error);
        } finally {
            setIsLoading(false);
        }
    };

    useEffect(() => { loadNav(); }, []);

    useLanguageChangeEffect(() => {
        fetchListNavigations({paging: {page: 1, pageSize: 10}})
            .then(res => {
                const navRes = res as unknown as { items: siteservicev1_Navigation[]; total: number };
                if (navRes.items?.length) {
                    const headerNav = navRes.items.find(nav => nav.location === 'HEADER' && nav.isActive === true);
                    if (headerNav?.items?.length) setNavigationItems(headerNav.items);
                }
            })
            .catch(error => console.error('[TopNavbar] 重新加载导航失败:', error))
            .finally(() => setIsLoading(false));
    }, {immediate: false, autoCleanup: true});

    const handleNavigate = (item: siteservicev1_NavigationItem) => {
        if (item.url != null) router.push(item.url);
    };

    if (isLoading) {
        return (
            <View className='flex gap-[16rpx] py-[16rpx]'>
                {Array.from({length: 4}).map((_, i) => (
                    <Skeleton key={i} className='h-[64rpx] w-[120rpx] rounded' />
                ))}
            </View>
        );
    }

    if (navigationItems.length === 0) return null;

    return (
        <View className='flex gap-[8rpx]'>
            {navigationItems.map((item) => {
                const itemId = item.id ?? 0;
                const hasChildren = (item.children?.length ?? 0) > 0;
                const isExpanded = expandedId === itemId;

                return (
                    <View key={itemId}>
                        <View
                          className='flex items-center gap-[8rpx] px-[24rpx] h-[64rpx] min-w-touch rounded text-desc text-textMain'
                          onClick={() => {
                                if (hasChildren) {
                                    setExpandedId(isExpanded ? null : itemId);
                                } else {
                                    handleNavigate(item);
                                }
                            }}
                          hoverClass='tap-active'
                        >
                            {item.icon && <XIcon name={`carbon:${item.icon}`} size={16} className='text-primary' />}
                            <Text className='whitespace-nowrap'>{item.title}</Text>
                            {hasChildren && (
                                <XIcon name='carbon:chevron-down' size={12} className={`text-textThird ${isExpanded ? 'rotate-180' : ''}`} />
                            )}
                        </View>

                        {hasChildren && isExpanded && (
                            <View className='ml-[24rpx] border-l-[2rpx] border-splitLine'>
                                {item.children!.map((child: siteservicev1_NavigationItem) => (
                                    <View
                                      key={child.id?.toString()}
                                      className='flex items-center gap-[8rpx] px-[24rpx] h-[72rpx] min-w-touch text-desc text-textSec'
                                      onClick={() => {
                                            handleNavigate(child);
                                            setExpandedId(null);
                                        }}
                                      hoverClass='tap-active'
                                    >
                                        {child.icon && <XIcon name={`carbon:${child.icon}`} size={14} className='text-textThird' />}
                                        <Text>{child.title}</Text>
                                    </View>
                                ))}
                            </View>
                        )}
                    </View>
                );
            })}
        </View>
    );
}
