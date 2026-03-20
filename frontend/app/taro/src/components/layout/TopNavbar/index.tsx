import React, {useState, useEffect} from 'react';
import {View, Text} from '@tarojs/components';

import {useI18nRouter} from '@/i18n/helpers/useI18nRouter';
import {useNavigationStore} from '@/store/slices/navigation/hooks';
import {useLanguageChangeEffect} from '@/hooks/useLanguageChangeEffect';
import XIcon from '@/plugins/xicon';

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
  const [isMenuOpen, setIsMenuOpen] = useState(false); // 移动端菜单开关
  const [expandedSubMenus, setExpandedSubMenus] = useState<Set<string>>(new Set()); // 展开的子菜单
  const [submenuPosition, setSubmenuPosition] = useState<{top: number; left: number} | null>(null); // 子菜单位置

  // 切换菜单显示
  const toggleMenu = () => {
    setIsMenuOpen(!isMenuOpen);
  };

  // 切换子菜单展开/收起
  const toggleSubMenu = (parentId: string, event: any) => {
    const newExpanded = new Set(expandedSubMenus);
    if (newExpanded.has(parentId)) {
      newExpanded.delete(parentId);
      setSubmenuPosition(null);
    } else {
      newExpanded.add(parentId);
      // 计算子菜单位置
      if (event && event.currentTarget) {
        const rect = event.currentTarget.getBoundingClientRect();
        setSubmenuPosition({
          top: rect.bottom + window.scrollY + 8, // 在菜单项下方，留 8px 间隙
          left: rect.left + window.scrollX, // 与菜单项左对齐
        });
      }
    }
    setExpandedSubMenus(newExpanded);
  };

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
  const handleMenuClick = (item: siteservicev1_NavigationItem, key: string, event?: any) => {
    // 如果有子菜单，切换展开状态，不跳转
    if (item.children && item.children.length > 0) {
      toggleSubMenu(key, event);
    } else {
      setActiveKey(key);
      handleNavigate(item);
      onClick?.(Number(key));
      setIsMenuOpen(false); // 移动端点击后关闭菜单
    }
  };

  return (
    <View className='navbar-wrapper'>
      {/* 移动端菜单按钮 */}
      <View className='mobile-menu-toggle' onClick={toggleMenu}>
        <View className={`hamburger ${isMenuOpen ? 'active' : ''}`}>
          <View className='hamburger-line' />
          <View className='hamburger-line' />
          <View className='hamburger-line' />
        </View>
      </View>

      {/* PC 端导航菜单 */}
      <View className={`navbar-content ${isMenuOpen ? 'menu-open' : ''}`}>
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
                  onClick={(event: any) => handleMenuClick(item, item.id?.toString() || `nav-${index}`, event)}
                >
                  {item.icon && (
                    <View className='menu-icon'>
                      <XIcon name={`carbon:${item.icon}`} size={18} />
                    </View>
                  )}
                  <Text>{item.title}</Text>
                  {/* 子菜单箭头指示器 */}
                  {item.children && item.children.length > 0 && (
                    <View className='submenu-arrow'>
                      <XIcon 
                        name={`carbon:chevron-${expandedSubMenus.has(item.id?.toString() || '') ? 'up' : 'down'}`} 
                        size={16} 
                      />
                    </View>
                  )}
                </View>

                {/* 子菜单 */}
                {item.children && item.children.length > 0 && (
                  <View 
                    className={`submenu ${expandedSubMenus.has(item.id?.toString() || '') ? 'submenu-open' : ''}`}
                    style={
                      submenuPosition && expandedSubMenus.has(item.id?.toString() || '')
                        ? {
                            top: `${submenuPosition.top}px`,
                            left: `${submenuPosition.left}px`,
                          }
                        : undefined
                    }
                  >
                    {item.children.map((child, childIndex) => (
                      <View
                        key={child.id}
                        className='submenu-item'
                        onClick={() => handleMenuClick(child, child.id?.toString() || `nav-child-${childIndex}`)}
                      >
                        {child.icon && (
                          <View className='submenu-icon'>
                            <XIcon name={`carbon:${child.icon}`} size={16} />
                          </View>
                        )}
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
