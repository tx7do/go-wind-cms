# PageContainer & ContentContainer 使用指南

## 概述

### PageContainer
企业级页面容器组件，提供：
- 自动面包屑导航
- 动态页面标题
- 权限控制
- 加载状态

### ContentContainer
**页面布局限定容器**，在业务页面中使用，提供：
- 统一的高度管理（固定/自适应）
- 一致的内边距
- Flex 布局支持

## 组件层级

```
MainLayout
  → PageContainer（面包屑、标题、权限）
    → 业务页面（MenuManagement、ApiManagement 等）
      → ContentContainer（布局容器）
        → 业务组件（Card、ProTable 等）
```

## 使用示例

### 固定高度（列表页、表格页）

```tsx
import { Card, Empty } from 'antd';
import ContentContainer from '@/layouts/components/PageContainer/ContentContainer';

const MyPage = () => {
  return (
    <ContentContainer heightMode="fixed" bottomMargin={16}>
      <Card style={{ flex: 1, overflow: 'auto' }}>
        <Empty description="页面内容" />
      </Card>
    </ContentContainer>
  );
};
```

**特点**：
- 占满可视区域
- 内容超出时滚动
- 可配置底部 margin

### 自适应高度（详情页、表单页）

```tsx
import { Card, Empty } from 'antd';
import ContentContainer from '@/layouts/components/PageContainer/ContentContainer';

const MyPage = () => {
  return (
    <ContentContainer heightMode="auto">
      <Card>
        <Empty description="页面内容" />
      </Card>
    </ContentContainer>
  );
};
```

**特点**：
- 跟随内容高度
- 不限制高度

## ContentContainer Props

```typescript
interface ContentContainerProps {
  /** 子组件 */
  children: React.ReactNode;
  
  /** 内边距，默认 16px */
  padding?: string | number;
  
  /** 是否显示滚动条，默认 false */
  scrollable?: boolean;
  
  /** 
   * 高度模式：
   * - 'auto': 自适应高度
   * - 'fixed': 固定高度（默认）
   */
  heightMode?: 'auto' | 'fixed';
  
  /** 底部外边距，仅在 fixed 模式生效，默认 0 */
  bottomMargin?: string | number;
  
  /** 自定义样式 */
  style?: React.CSSProperties;
  
  /** 自定义类名 */
  className?: string;
}
```

## 样式规范

### 固定高度模式
```css
.content-container--fixed {
  flex: 1;
  min-height: 0;
  overflow: hidden;
}
```

### 自适应高度模式
```css
.content-container--auto {
  height: auto;
  min-height: 0;
  overflow: visible;
}
```

## 注意事项

1. **ContentContainer 在业务页面中使用**，不在 PageContainer 内部自动调用
2. 页面组件的子元素需要设置 `flex: 1` 才能占满空间（仅 fixed 模式）
3. 需要滚动时，在子元素上设置 `overflow: 'auto'`
4. 全屏模式下应使用自适应高度
