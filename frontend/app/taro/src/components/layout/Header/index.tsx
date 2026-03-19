import {useState, useEffect} from 'react';
import {useTranslation} from 'react-i18next';
import {View, Text, Image} from '@tarojs/components';

import {useI18n} from '@/i18n';
import {useI18nRouter} from '@/i18n/helpers/useI18nRouter';
import {useThemeStore, useThemeMode} from '@/store/core/theme/hooks';
import {ThemeMode} from "@/store/types";
import XIcon from '@/plugins/xicon';

import TopNavbar from '../TopNavbar';

import './index.scss';
import {useAccessStore} from "@/store/core/access/hooks";
import {useAuthenticationStore} from "@/store/slices/authentication/hooks";

export default function Header() {
  const {t: navbarT} = useTranslation('navbar');
  const {t: appT} = useTranslation('app');
  const {t: menuT} = useTranslation('menu');
  const brandTitle = appT('title');

  const themeStore = useThemeStore();
  const currentMode = useThemeMode();
  const {changeLocale} = useI18n();
  const router = useI18nRouter();
  const accessStore = useAccessStore();
  const authenticationStore = useAuthenticationStore();

  const accessToken = accessStore.access.accessToken;
  const isLogin = !!accessToken && !accessStore.access.loginExpired;

  const handleClickLogo = () => {
    router.push('/');
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
      await authenticationStore.logout();
    }
  };

  // 根据登录状态动态生成用户菜单项
  const userMenuItems = isLogin
    ? [
      {
        key: 'homepage',
        label: menuT('homepage'),
        icon: 'carbon:home',
        onClick: handleClickUserHomepage
      },
      {
        key: 'profile',
        label: menuT('my_profile'),
        icon: 'carbon:user',
        onClick: handleClickSettings
      },
      {
        type: 'divider' as const
      },
      {
        key: 'logout',
        label: menuT('logout'),
        icon: 'carbon:logout',
        danger: true,
        onClick: handleClickLogout
      },
    ]
    : [
      {
        key: 'login',
        label: navbarT('user.login'),
        icon: 'carbon:login',
        onClick: handleClickLogin
      },
      {
        key: 'register',
        label: navbarT('user.register'),
        icon: 'carbon:register',
        onClick: handleClickRegister
      }
    ];

  // 语言菜单
  const languageMenuItems = [
    {
      key: 'zh-CN',
      label: '简体中文',
      onClick: () => changeLocale('zh-CN')
    },
    {
      key: 'en-US',
      label: 'English',
      onClick: () => changeLocale('en-US')
    }
  ];

  const themeMenuItems = [
    {
      key: 'dark',
      label: navbarT('theme.dark'),
      icon: 'carbon:moon',
      onClick: () => themeStore.setMode('dark')
    },
    {
      key: 'light',
      label: navbarT('theme.light'),
      icon: 'carbon:sun',
      onClick: () => themeStore.setMode('light')
    },
    {
      key: 'system',
      label: navbarT('theme.system'),
      icon: 'carbon:ibm-watsonx-orchestrate',
      onClick: () => themeStore.setMode('system')
    },
  ];

  const themeIconMap: Record<ThemeMode, string> = {
    dark: 'carbon:moon',
    light: 'carbon:sun',
    system: 'carbon:ibm-watsonx-orchestrate'
  };

  const validMode: ThemeMode = (currentMode && ['dark', 'light', 'system'].includes(currentMode))
    ? currentMode
    : 'system';
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
              onClick={handleClickLogo}
            >
              <Image
                src='/logo.png'
                style={{width: 55, height: 55}}
                className='logo'
              />
              <Text className='site-name'>{brandTitle}</Text>
            </View>
            {/* 主导航菜单 */}
            <View className='navbar-menu-wrap'>
              <TopNavbar/>
            </View>
          </View>
          {/* 功能按钮区 */}
          <View className='actions'>
            <View className='space'>
              <View className='dropdown-wrapper'>
                <View
                  className='icon-btn'
                  onClick={() => console.log('User menu clicked')}
                >
                  <XIcon name='carbon:user' size={20} />
                </View>
                {/* TODO: 实现 Dropdown 组件 */}
              </View>
              <View className='dropdown-wrapper'>
                <View
                  className='icon-btn'
                  onClick={() => console.log('Language clicked')}
                >
                  <XIcon name='carbon:language' size={20} className='lang-icon' />
                </View>
              </View>
              <View className='dropdown-wrapper'>
                <View
                  className='icon-btn'
                  onClick={() => console.log('Theme clicked')}
                >
                  <XIcon name={displayIcon} size={20} className='theme-icon' />
                </View>
              </View>
            </View>
          </View>
        </View>
      </View>
    </View>
  );
}
