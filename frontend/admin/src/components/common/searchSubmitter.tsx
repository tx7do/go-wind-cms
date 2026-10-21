import type { ReactElement, ReactNode } from 'react';

import { SearchOutlined, SyncOutlined } from '@ant-design/icons';
import { cloneElement, isValidElement } from 'react';

type SearchButtonElement = ReactElement<{
  icon?: ReactNode;
  type?: 'default' | 'primary' | 'dashed' | 'link' | 'text';
}>;

/**
 * 列表搜索栏提交区统一配置（三端 UX 对齐 ele）：
 * 「搜索」实心蓝 + 放大镜在前，「重置」文字样式 + 循环箭头在后。
 *
 * pro-form 的 Submitter 固定按 [重置, 提交] 顺序渲染且按钮不带图标；
 * SearchConfig 类型层又 Omit 了 submitter，故走 optionRender：
 * 翻转默认 dom 顺序并用 cloneElement 补图标/样式。
 * 用法：search={{ ...原有配置, optionRender: searchOptionRender }}
 */
export const searchOptionRender = (
  _searchConfig: unknown,
  _props: unknown,
  dom: ReactNode[],
): ReactNode[] => {
  const nodes = dom.filter((n): n is SearchButtonElement => isValidElement(n));
  if (nodes.length < 2) return dom;
  const [reset, submit] = nodes;
  return [
    cloneElement(submit, { icon: <SearchOutlined /> }),
    cloneElement(reset, { icon: <SyncOutlined />, type: 'text' }),
  ];
};
