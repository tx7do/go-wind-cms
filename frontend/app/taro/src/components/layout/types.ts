import {ReactNode, ComponentType, ElementType} from 'react';

/**
 * 导航项基础配置
 */
export interface TopNavBarTabItem<P = Record<string, unknown>> {
    /** 唯一标识，用于路由匹配和状态管理 */
    key: string;

    /** 显示文本，支持字符串或 ReactNode（可渲染自定义内容） */
    label: ReactNode;

    /** 图标，支持：图标名称字符串 / 图标组件 / ReactNode */
    icon?: ReactNode | ElementType<{ className?: string }>;

    /** 关联的页面组件，使用泛型支持组件 props 类型检查 */
    component?: ComponentType<P>;

    /** 路由路径（可选，用于自动路由跳转） */
    path?: string;

    /** 是否禁用该标签项 */
    disabled?: boolean;

    /** 是否隐藏该标签项（不渲染） */
    hidden?: boolean;

    /** 角标配置，支持数字/文字/自定义内容 */
    badge?: {
        content?: ReactNode;
        color?: string;
        dot?: boolean;  // 是否显示为小红点
        offset?: [number, number];  // 偏移量 [x, y]
    };

    /** 子菜单项（支持多级导航） */
    children?: TopNavBarTabItem<P>[];

    /** 权限标识，用于权限控制显示 */
    permission?: string | string[];

    /** 自定义数据，透传给组件使用 */
    meta?: Record<string, unknown>;

    /** 组件初始化时传递的固定 props（与 component 泛型 P 匹配） */
    componentProps?: Partial<P>;
}
