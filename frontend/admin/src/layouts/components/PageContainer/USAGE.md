# 页面刷新功能使用说明

## 概述

顶部栏的刷新按钮现在只刷新当前激活的 PageContainer，而不是整个页面。这是通过全局的页面刷新管理器实现的。

## 工作原理

1. **PageContainer 注册刷新回调**：每个 PageContainer 组件在挂载时，会将自己的 `onRefresh` 回调注册到全局 store 中
2. **顶部栏触发刷新**：点击顶部栏的刷新按钮时，会调用所有已注册的刷新回调
3. **组件卸载时清理**：PageContainer 卸载时会自动移除自己的刷新回调

## 使用方法

### 基础用法

```tsx
import { PageContainer } from '@/layouts/components/PageContainer';

const MyPage = () => {
  // 定义刷新逻辑
  const handleRefresh = async () => {
    console.log('刷新页面数据...');
    // 可以执行 API 请求等操作
    await fetchData();
  };

  return (
    <PageContainer
      title="我的页面"
      showRefresh={true}        // 显示刷新按钮（可选）
      onRefresh={handleRefresh}  // 刷新回调
    >
      {/* 页面内容 */}
    </PageContainer>
  );
};
```

### 完整示例

```tsx
import { useState, useEffect } from 'react';
import { PageContainer } from '@/layouts/components/PageContainer';
import { Table, Card } from 'antd';

const UserListPage = () => {
  const [data, setData] = useState([]);
  const [loading, setLoading] = useState(false);

  // 加载数据
  const loadData = async () => {
    setLoading(true);
    try {
      const response = await fetch('/api/users');
      const result = await response.json();
      setData(result.data);
    } finally {
      setLoading(false);
    }
  };

  // 初始加载
  useEffect(() => {
    loadData();
  }, []);

  // 刷新回调
  const handleRefresh = async () => {
    console.log('正在刷新用户列表...');
    await loadData();
    console.log('刷新完成');
  };

  return (
    <PageContainer
      title="用户列表"
      loading={loading}
      showRefresh={true}
      onRefresh={handleRefresh}
    >
      <Card>
        <Table 
          dataSource={data} 
          columns={[
            { title: '姓名', dataIndex: 'name' },
            { title: '邮箱', dataIndex: 'email' },
          ]}
        />
      </Card>
    </PageContainer>
  );
};

export default UserListPage;
```

### Dashboard 示例

```tsx
import { PageContainer } from '@/layouts/components/PageContainer';
import { StatsCard, TrendChart } from './components';

const Dashboard = () => {
  // 刷新统计数据
  const handleRefresh = async () => {
    console.log('刷新 Dashboard 数据...');
    // 可以在这里重新获取统计数据、图表数据等
    await Promise.all([
      fetchStats(),
      fetchTrendData(),
      fetchChartData(),
    ]);
  };

  return (
    <PageContainer
      title="工作台"
      showRefresh={true}
      onRefresh={handleRefresh}
    >
      <Row gutter={[16, 16]}>
        <Col span={6}>
          <StatsCard title="用户量" value="2,000" />
        </Col>
        {/* 其他统计卡片 */}
      </Row>
      <TrendChart />
    </PageContainer>
  );
};
```

## API 说明

### PageContainer Props

| 属性 | 类型 | 说明 |
|------|------|------|
| `showRefresh` | `boolean` | 是否在 PageContainer 内部显示刷新按钮 |
| `onRefresh` | `() => void \| Promise<void>` | 刷新回调函数，支持异步操作 |
| `pageKey` | `string` | 页面唯一标识，默认为路由路径 |

### 工作流程

```
用户点击顶部栏刷新按钮
         ↓
triggerPageRefresh() 被调用
         ↓
遍历所有注册的 refreshCallbacks
         ↓
依次执行每个页面的 onRefresh 回调
         ↓
页面数据更新，UI 自动刷新
```

## 注意事项

1. **多个 PageContainer**：如果同时有多个 PageContainer 存在（例如 Tab 页），点击刷新会触发所有页面的刷新回调

2. **异步操作**：`onRefresh` 可以是异步函数，会在执行完成后自动结束

3. **错误处理**：建议在 `onRefresh` 中自行处理错误，store 会捕获并记录错误但不会中断其他页面的刷新

4. **性能考虑**：如果只需要刷新特定页面，可以在 `onRefresh` 中添加判断逻辑：

```tsx
const handleRefresh = async () => {
  // 只在页面可见时才刷新
  if (!document.hidden) {
    await loadData();
  }
};
```

## 与全屏功能的配合

刷新和全屏功能可以同时使用：

```tsx
<PageContainer
  title="我的页面"
  showRefresh={true}      // 显示刷新按钮
  showFullscreen={true}   // 显示全屏按钮
  onRefresh={handleRefresh}
>
  {/* 页面内容 */}
</PageContainer>
```

在全屏状态下，刷新按钮仍然可用，只会刷新内容区域。
