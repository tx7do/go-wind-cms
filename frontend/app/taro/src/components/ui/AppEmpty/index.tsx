import React from 'react';
import {View, Text} from '@tarojs/components';

import XIcon from '@/plugins/xicon';

import './index.scss';

interface AppEmptyProps {
  /**
   * 空状态类型
   * - default: 默认空状态
   * - error: 错误状态 (如无效的 ID)
   * - noData: 无数据状态
   */
  variant?: 'default' | 'error' | 'noData';
  /**
   * 是否显示在页面容器中 (添加额外的 margin)
   */
  inContainer?: boolean;
  /**
   * 自定义图标
   */
  image?: React.ReactNode;
  /**
   * 描述文字
   */
  description?: string;
  /**
   * 自定义类名
   */
  className?: string;
  /**
   * 自定义样式
   */
  style?: React.CSSProperties;
}

/**
 * App 级别的 Empty 组件
 *
 * 特性:
 * - 统一的暗色模式适配
 * - 预设的样式变体
 * - 一致的视觉风格
 * - 支持自定义图标和样式
 */
export const AppEmpty: React.FC<AppEmptyProps> = ({
  variant = 'default',
  inContainer = false,
  image,
  description,
  className = '',
  style,
}) => {
  // 根据 variant 设置默认图标
  const defaultImage = image || (
    variant === 'error'
      ? <XIcon name='carbon:warning' size={48} className='empty-icon-svg' />
      : variant === 'noData'
        ? <XIcon name='carbon:chart-line' size={48} className='empty-icon-svg' />
        : <XIcon name='carbon:document' size={48} className='empty-icon-svg' />
  );

  return (
    <View
      className={`empty-wrapper ${inContainer ? 'in-container' : ''} ${className}`}
      style={style}
    >
      <View className='app-empty'>
        <View className='empty-image'>
          {defaultImage}
        </View>
        {description && (
          <Text className='empty-description'>{description}</Text>
        )}
      </View>
    </View>
  );
};

export default AppEmpty;
