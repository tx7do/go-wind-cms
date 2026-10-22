import {PropsWithChildren} from 'react';
import {View} from '@tarojs/components';
import Header from './Header';
import NavigationProgress from './NavigationProgress';
import ThemeClientProvider from './ThemeClientProvider';

// 移动端 App 布局：仅保留顶栏与页面内容。
// 基座带来的桌面网页 Footer/BackToTop（关于我们/免责条款等站点页脚）
// 不适合学生端 H5，已移除；法律声明类页面仍可通过路由直达。
//
// 顶栏只在 H5 渲染：小程序自带原生导航栏，再叠一层自定义栏就是两条栏，
// 且容器为顶栏预留的偏移会在原生栏之下再空出一截。
// 偏移量交给 .has-app-header，避免 96rpx 这个数在 TSX 与 CSS 各写一遍。
const isH5 = process.env.TARO_ENV === 'h5';

function Layout({children}: PropsWithChildren) {
  return (
    <ThemeClientProvider>
      <NavigationProgress />
      {isH5 && <Header />}
      <View className={isH5 ? 'app-page-container has-app-header' : 'app-page-container'}>
        {children}
      </View>
    </ThemeClientProvider>
  );
}

export default Layout;
