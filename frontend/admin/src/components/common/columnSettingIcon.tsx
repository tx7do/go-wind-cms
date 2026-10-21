import type { HTMLAttributes } from 'react';

import React from 'react';

type ColumnSlidersIconProps = HTMLAttributes<HTMLSpanElement>;

/**
 * 列设置按钮的调节器图标（三滑杆，lucide sliders-horizontal 内联复刻；vben 同款造型）。
 * antd 图标集无此裸滑杆造型（ControlOutlined 带方框）。
 * 必须 forwardRef：pro-table ColumnSetting 用 Tooltip 包裹本图标，
 * Tooltip 靠 ref 定位触发器——函数组件挂不上 ref 时 Tooltip 静默失效。
 * 外层 .anticon span 借用 antd 的垂直对齐样式——裸 svg 会比相邻 antd 图标高 ~2px。
 * 用法：options={{ setting: { settingIcon: <ColumnSlidersIcon /> } }}
 */
export const ColumnSlidersIcon = React.forwardRef<
  HTMLSpanElement,
  ColumnSlidersIconProps
>((props, ref) => (
  <span
    ref={ref}
    {...props}
    className="anticon anticon-column-sliders"
    role="img"
    aria-hidden="true"
  >
    <svg
      viewBox="0 0 24 24"
      width="1em"
      height="1em"
      fill="none"
      stroke="currentColor"
      strokeWidth={2}
      strokeLinecap="round"
      aria-hidden="true"
    >
      <path d="M21 4h-7M10 4H3M14 2v4" />
      <path d="M21 12h-9M8 12H3M8 10v4" />
      <path d="M21 20h-5M12 20H3M16 18v4" />
    </svg>
  </span>
));

ColumnSlidersIcon.displayName = 'ColumnSlidersIcon';
