import { Splitter } from 'antd';
import { useState } from 'react';

import ContentContainer from '@/layouts/components/PageContainer/ContentContainer';
import PermissionGroupList from './PermissionGroupList';
import PermissionList from './PermissionList';

/**
 * 权限点管理页面
 * 使用 Splitter 实现左右分栏布局
 */
const PermissionPointManagement = () => {
  const [currentGroupId, setCurrentGroupId] = useState<number | null>(null);

  return (
    <ContentContainer heightMode="fixed" padding="16px" bottomMargin={0}>
      <Splitter style={{ height: '100%', flex: 1, minHeight: 0 }}>
        <Splitter.Panel collapsible defaultSize="40%" min="25%" max="55%" style={{ display: 'flex', flexDirection: 'column' }}>
          <PermissionGroupList currentGroupId={currentGroupId} onGroupSelect={setCurrentGroupId} />
        </Splitter.Panel>
        <Splitter.Panel style={{ display: 'flex', flexDirection: 'column' }}>
          <PermissionList groupId={currentGroupId} />
        </Splitter.Panel>
      </Splitter>
    </ContentContainer>
  );
};

export default PermissionPointManagement;
