import React from 'react';
import { Icon } from '@iconify/react';
import * as Icons from '@ant-design/icons';

/**
 * 根据图标名称字符串解析为图标组件
 * 支持两种格式：
 * 1. Iconify 格式: "lucide:building-2", "mdi:home", "carbon:dashboard"
 * 2. Ant Design 格式: "DashboardOutlined", "HomeOutlined"（向后兼容）
 */
export const getIconFromName = (iconName?: string): React.ReactNode => {
  if (!iconName) return null;

  // 判断是否为 Iconify 格式（包含冒号）
  if (iconName.includes(':')) {
    const [prefix, name] = iconName.split(':');
    return <Icon icon={`${prefix}:${name}`} width="1em" height="1em" />;
  }

  // 向后兼容：Ant Design 图标
  const IconComponent = (Icons as any)[iconName];
  if (IconComponent) {
    return React.createElement(IconComponent);
  }

  return null;
};
