import type { HTMLAttributes } from 'react';

import React from 'react';

type SettingBoltIconProps = HTMLAttributes<HTMLSpanElement>;

/**
 * 偏好设置按钮图标（六角螺母/bolt 造型，Remix settings-3 内联复刻；ele 顶栏 SvgIcon "setting" 同款）。
 * antd 图标集无此造型（SettingOutlined 是齿轮）。
 * 必须 forwardRef：Tooltip 包裹本图标时靠 ref 定位触发器，
 * 函数组件挂不上 ref 时 Tooltip 静默失效。
 * 外层 .anticon span 借用 antd 的垂直对齐样式——裸 svg 会比相邻 antd 图标高 ~2px。
 */
export const SettingBoltIcon = React.forwardRef<
  HTMLSpanElement,
  SettingBoltIconProps
>((props, ref) => (
  <span
    ref={ref}
    {...props}
    className="anticon anticon-setting-bolt"
    role="img"
    aria-hidden="true"
  >
    <svg
      viewBox="0 0 24 24"
      width="1em"
      height="1em"
      fill="currentColor"
      aria-hidden="true"
    >
      <path d="m12 1 9.5 5.5v11L12 23l-9.5-5.5v-11L12 1zm0 2.311L4.5 7.653v8.694l7.5 4.342 7.5-4.342V7.653L12 3.311zM12 16a4 4 0 1 1 0-8 4 4 0 0 1 0 8zm0-2a2 2 0 1 0 0-4 2 2 0 0 0 0 4z" />
    </svg>
  </span>
));

SettingBoltIcon.displayName = 'SettingBoltIcon';
