import { Splitter } from 'antd';
import { useState } from 'react';

import ContentContainer from '@/layouts/components/PageContainer/ContentContainer';
import NavigationItemList from './NavigationItemList';
import NavigationList from './NavigationList';

/**
 * 导航管理页面
 * 使用 Splitter 实现左右分栏布局：左侧导航列表，右侧导航项列表
 */
const NavigationManagement = () => {
  const [currentNavigationId, setCurrentNavigationId] = useState<number | null>(
    null,
  );

  return (
    <ContentContainer heightMode="fixed" padding="16px" bottomMargin={0}>
      <Splitter style={{ height: '100%', flex: 1, minHeight: 0 }}>
        <Splitter.Panel
          collapsible
          defaultSize="45%"
          min="25%"
          max="75%"
          style={{ display: 'flex', flexDirection: 'column' }}
        >
          <NavigationList
            currentNavigationId={currentNavigationId}
            onNavigationSelect={setCurrentNavigationId}
          />
        </Splitter.Panel>
        <Splitter.Panel style={{ display: 'flex', flexDirection: 'column' }}>
          <NavigationItemList navigationId={currentNavigationId} />
        </Splitter.Panel>
      </Splitter>
    </ContentContainer>
  );
};

export default NavigationManagement;
