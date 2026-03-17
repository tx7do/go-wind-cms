import React, {useState, useEffect} from 'react';
import {View, Text} from '@tarojs/components';

import {useI18nRouter} from '@/i18n/helpers/useI18nRouter';
import {useNavigationStore} from '@/store/slices/navigation/hooks';
import {useLanguageChangeEffect} from '@/hooks/useLanguageChangeEffect';

import type {siteservicev1_Navigation, siteservicev1_NavigationItem} from '@/api/generated/app/service/v1';

import './index.scss';

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
  const [activeKey, setActiveKey] = useState<string | null>(null);

  // 加载导航数据
  useEffect(() => {
    async function loadNavigation() {
      try {
        setIsLoading(true);
        const res = await navigationStore.listNavigation({
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
  }, []);

  // 监听语言变化，自动重新加载导航数据
  useLanguageChangeEffect(() => {
    setIsLoading(true);
    navigationStore.listNavigation({
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
      // Taro 中打开外部链接需要使用 Taro.navigateToMiniProgram 或其他方式
      console.log('Open external link:', item.url);
    } else {
      if (item.url != null) {
        router.push(item.url);
      }
    }
  };

  // 处理菜单点击
  const handleMenuClick = (item: siteservicev1_NavigationItem, key: string) => {
    setActiveKey(key);
    handleNavigate(item);
    onClick?.(Number(key));
  };

  return (
    <View className='navbar-wrapper'>
      <View className='navbar-content'>
        {/* 导航菜单 */}
        {isLoading ? (
          <View className='loading-state'>
            Loading navigation...
          </View>
        ) : navigationItems.length > 0 ? (
          <View className='navbar-menu'>
            {navigationItems.map((item, index) => (
              <React.Fragment key={item.id}>
                {/* 一级菜单 */}
                <View
                  className={`menu-item ${activeKey === item.id?.toString() ? 'active' : ''}`}
                  onClick={() => handleMenuClick(item, item.id?.toString() || `nav-${index}`)}
                >
                  {item.icon && <Text className='menu-icon'>{item.icon}</Text>}
                  <Text>{item.title}</Text>
                </View>

                {/* 子菜单 */}
                {item.children && item.children.length > 0 && (
                  <View className='submenu'>
                    {item.children.map((child, childIndex) => (
                      <View
                        key={child.id}
                        className='submenu-item'
                        onClick={() => handleMenuClick(child, child.id?.toString() || `nav-child-${childIndex}`)}
                      >
                        {child.icon && <Text className='submenu-icon'>{child.icon}</Text>}
                        <Text>{child.title}</Text>
                      </View>
                    ))}
                  </View>
                )}
              </React.Fragment>
            ))}
          </View>
        ) : (
          <View className='loading-state'>
            No navigation items
          </View>
        )}
      </View>
    </View>
  );
}
