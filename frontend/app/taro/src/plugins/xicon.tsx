import React from 'react';
import {Image} from '@tarojs/components';

interface XIconProps {
  name: string; // H5: iconify 名称；小程序: 本地 svg 路径
  size?: number | string;
  color?: string;
  className?: string;
  style?: React.CSSProperties;
  svg?: string; // 小程序端可传 SVG 代码字符串
}

export const XIcon: React.FC<XIconProps> = ({name, size = 24, color, className = '', style, svg}) => {
  const env = process.env.TARO_ENV;
  // UnoCSS class
  const unoClass = [
    typeof size === 'number' ? `text-[${size}px]` : `text-[${size}]`,
    color ? `text-[${color}]` : '',
    className
  ].filter(Boolean).join(' ');

  if (env === 'h5') {
    // H5端用iconify web component（需在index.html引入iconify-icon.min.js）
    return (
      // @ts-ignore
      <iconify-icon
        icon={name}
        class={unoClass}
        style={{fontSize: size, color, ...style}}
      />
    );
  }
  // 小程序端
  if (svg) {
    // 直接渲染SVG代码
    // @ts-ignore
    return <svg style={{width: size, height: size, ...style}} className={unoClass}
                dangerouslySetInnerHTML={{__html: svg}}/>;
  }
  // 渲染本地SVG文件
  return (
    <Image
      src={name}
      style={{width: size, height: size, ...style}}
      className={unoClass}
      mode="aspectFit"
    />
  );
};

export default XIcon;
