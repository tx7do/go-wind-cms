import {PropsWithChildren} from 'react';
import {View} from '@tarojs/components';
import Taro from '@tarojs/taro';

import Header from '../Header';
import Footer from '../Footer';

// 需要排除 Header 和 Footer 的页面路径
const EXCLUDED_PAGES = [
  '/pages/login/index',
  '/pages/register/index',
  '/pages/forgot-password/index', // 如果后续创建忘记密码页面
];

interface LayoutProps extends PropsWithChildren {}

/**
 * 全局布局组件
 * - Header 固定在顶部
 * - Footer 在页面底部（不固定）
 * - 自动排除登录、注册、忘记密码等页面
 */
export default function Layout({children}: LayoutProps) {
  // 获取当前页面路径
  const currentPagePath = Taro.getCurrentPages()[0]?.route;

  // 判断是否需要排除 Header 和 Footer
  const shouldExcludeLayout = currentPagePath && EXCLUDED_PAGES.some(path =>
    currentPagePath.startsWith(path.replace('/pages', ''))
  );

  // 如果是排除页面，直接渲染内容
  if (shouldExcludeLayout) {
    return <>{children}</>;
  }

  // 否则渲染完整布局
  return (
    <View className='global-layout'>
      {/* 固定顶部的 Header */}
      <Header />

      {/* 页面内容区域 */}
      <View className='page-content'>
        {children}
      </View>

      {/* 底部的 Footer（不固定） */}
      <Footer />
    </View>
  );
}
