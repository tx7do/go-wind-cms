import React from 'react';
import {View, Text} from '@tarojs/components';
import {useTranslation} from "react-i18next";

import './index.scss';


interface ButtonProps {
  type?: 'primary' | 'default' | 'danger';
  size?: 'small' | 'medium' | 'large';
  disabled?: boolean;
  loading?: boolean;
  onClick?: () => void;
  children: React.ReactNode;
  className?: string;
}

export function Button({
                         type = 'default',
                         size = 'medium',
                         disabled = false,
                         loading = false,
                         onClick,
                         children,
                         className = '',
                       }: ButtonProps) {
  const handleClick = () => {
    if (!disabled && !loading && onClick) {
      onClick();
    }
  };

  const {t} = useTranslation();

  return (
    <View
      className={`taro-button taro-button--${type} taro-button--${size} ${disabled ? 'taro-button--disabled' : ''} ${loading ? 'taro-button--loading' : ''} ${className}`}
      onClick={handleClick}
    >
      {loading && <Text className='taro-button__loading'>{t('common.loading')}</Text>}
      <Text>{children}</Text>
    </View>
  );
}
