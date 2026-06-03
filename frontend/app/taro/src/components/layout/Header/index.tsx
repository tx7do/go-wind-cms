import {useState, useEffect} from 'react';
import {useTranslation} from 'react-i18next';
import {View, Text, Image} from '@tarojs/components';

import {useI18nRouter} from '@/i18n/helpers/useI18nRouter';
import XIcon from '@/plugins/xicon';

import {useAccessStore} from "@/store/core/access/store";
import {usePreferences} from '@/core/preferences/hooks/usePreferences';

import TopNavbar from '../TopNavbar';

import './index.scss';


export default function Header() {
  const {t} = useTranslation();

  const preferences = usePreferences();
  const router = useI18nRouter();
  const accessStore = useAccessStore();

  const accessToken = accessStore.accessToken;
  const isLogin = !!accessToken && !accessStore.loginExpired;

  // 下拉菜单状态
  const [openDropdown, setOpenDropdown] = useState<'user' | 'language' | 'theme' | null>(null);

  const handleClickLogo = () => {
    router.push('pages/index');
  };
  const handleClickSettings = () => {
    router.push('/settings');
  };
  const handleClickUserHomepage = () => {
    router.push('/user');
  };
  const handleClickLogin = () => {
    router.push('/login');
  };
  const handleClickRegister = () => {
    router.push('/register');
  };
  const handleClickLogout = async () => {
    console.log('Logout');

    if (isLogin) {
      // TODO: 调用登出API
      accessStore.clearTokens();
    }
    setOpenDropdown(null);
  };

  // 切换下拉菜单
  const toggleDropdown = (type: 'user' | 'language' | 'theme') => {
    setOpenDropdown(openDropdown === type ? null : type);
  };

  // 关闭所有下拉菜单
  const closeAllDropdowns = () => {
    setOpenDropdown(null);
  };

  // 根据登录状态动态生成用户菜单项
  const userMenuItems = isLogin
    ? [
      {
        key: 'homepage',
        label: t('menu.homepage'),
        icon: 'carbon:home',
        onClick: handleClickUserHomepage
      },
      {
        key: 'profile',
        label: t('menu.my_profile'),
        icon: 'carbon:user',
        onClick: handleClickSettings
      },
      {
        type: 'divider' as const
      },
      {
        key: 'logout',
        label: t('menu.logout'),
        icon: 'carbon:logout',
        danger: true,
        onClick: handleClickLogout
      },
    ]
    : [
      {
        key: 'login',
        label: t('navbar.user.login'),
        icon: 'carbon:user',
        onClick: handleClickLogin
      },
      {
        key: 'register',
        label: t('navbar.user.register'),
        icon: 'carbon:user',
        onClick: handleClickRegister
      }
    ];

  // 语言菜单
  const languageMenuItems = [
    {
      key: 'zh-CN',
      label: '简体中文',
      onClick: () => preferences.setLanguage('zh-CN')
    },
    {
      key: 'en-US',
      label: 'English',
      onClick: () => preferences.setLanguage('en-US')
    }
  ];

  const themeMenuItems = [
    {
      key: 'dark',
      label: t('navbar.theme.dark'),
      icon: 'carbon:moon',
      onClick: () => preferences.setThemeMode('dark')
    },
    {
      key: 'light',
      label: t('navbar.theme.light'),
      icon: 'carbon:sun',
      onClick: () => preferences.setThemeMode('light')
    },
    {
      key: 'auto',
      label: t('navbar.theme.system'),
      icon: 'carbon:ibm-watsonx-orchestrate',
      onClick: () => preferences.setThemeMode('auto')
    },
  ];

  const themeIconMap: Record<string, string> = {
    dark: 'carbon:moon',
    light: 'carbon:sun',
    auto: 'carbon:ibm-watsonx-orchestrate'
  };

  const validMode = ['dark', 'light', 'auto'].includes(preferences.theme.mode)
    ? preferences.theme.mode
    : 'auto';
  const iconValue = themeIconMap[validMode];

  // SSR 兼容
  const [mounted, setMounted] = useState(false);
  useEffect(() => {
    setMounted(true);
  }, []);

  const displayIcon = mounted ? iconValue : 'carbon:ibm-watsonx-orchestrate';

  return (
    <View className='fixed-top'>
      <View className='header-inner'>
        <View className='top-bar'>
          {/* Logo + 导航区 */}
          <View className='logo-nav-section'>
            <View
              className='logo-section'
              role='button'
              aria-label='Go to homepage'
              onClick={handleClickLogo}
            >
              <Image
                src='/assets/images/logo.png'
                className='logo'
              />
              <Text className='site-name'>{t('app.title')}</Text>
            </View>
            {/* 主导航菜单 */}
            <View className='navbar-menu-wrap'>
              <TopNavbar />
            </View>
          </View>
          {/* 功能按钮区 */}
          <View className='actions'>
            <View className='space'>
              {/* 用户菜单 */}
              <View className='dropdown-wrapper'>
                <View
                  className='icon-btn'
                  role='button'
                  aria-label='User menu'
                  onClick={() => toggleDropdown('user')}
                >
                  <XIcon name='carbon:user' size={24} />
                </View>
                {openDropdown === 'user' && (
                  <View className='dropdown-menu'>
                    {userMenuItems.map((item) => (
                      <View
                        key={item.key}
                        className={`dropdown-item ${item.danger ? 'danger' : ''}`}
                        onClick={() => {
                          if (item.onClick) item.onClick();
                          closeAllDropdowns();
                        }}
                      >
                        {item.icon && (
                          <XIcon name={item.icon} size={18} className='dropdown-item-icon' />
                        )}
                        <Text className='dropdown-item-label'>{item.label}</Text>
                      </View>
                    ))}
                  </View>
                )}
              </View>
              {/* 语言菜单 */}
              <View className='dropdown-wrapper'>
                <View
                  className='icon-btn'
                  role='button'
                  aria-label='Language'
                  onClick={() => toggleDropdown('language')}
                >
                  <XIcon name='carbon:language' size={24} className='lang-icon' />
                </View>
                {openDropdown === 'language' && (
                  <View className='dropdown-menu'>
                    {languageMenuItems.map((item) => (
                      <View
                        key={item.key}
                        className='dropdown-item'
                        onClick={() => {
                          if (item.onClick) {
                            item.onClick();
                          }
                          closeAllDropdowns();
                        }}
                      >
                        <Text className='dropdown-item-label'>{item.label}</Text>
                      </View>
                    ))}
                  </View>
                )}
              </View>
              {/* 主题菜单 */}
              <View className='dropdown-wrapper'>
                <View
                  className='icon-btn'
                  role='button'
                  aria-label='Toggle theme'
                  onClick={() => toggleDropdown('theme')}
                >
                  <XIcon name={displayIcon} size={24} className='theme-icon' />
                </View>
                {openDropdown === 'theme' && (
                  <View className='dropdown-menu'>
                    {themeMenuItems.map((item) => (
                      <View
                        key={item.key}
                        className='dropdown-item'
                        onClick={() => {
                          item.onClick();
                          closeAllDropdowns();
                        }}
                      >
                        {item.icon && (
                          <XIcon name={item.icon} size={18} className='dropdown-item-icon' />
                        )}
                        <Text className='dropdown-item-label'>{item.label}</Text>
                      </View>
                    ))}
                  </View>
                )}
              </View>
            </View>
          </View>
        </View>
      </View>
    </View>
  );
}
