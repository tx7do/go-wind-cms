import {PropsWithChildren} from 'react';
import {View} from '@tarojs/components';
import Header from './Header';
import Footer from './Footer';
import NavigationProgress from './NavigationProgress';
import BackToTop from './BackToTop';
import ThemeClientProvider from './ThemeClientProvider';

function Layout({children}: PropsWithChildren) {
  return (
    <ThemeClientProvider>
      <NavigationProgress />
      <Header />
      <View className='min-h-screen pt-[88rpx]'>
        {children}
      </View>
      <Footer />
      <BackToTop />
    </ThemeClientProvider>
  );
}

export default Layout;
