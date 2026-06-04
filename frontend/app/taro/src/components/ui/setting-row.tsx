import React from 'react';
import {View} from '@tarojs/components';
import {cn} from '@/lib/utils';

export interface SettingRowProps {
    /** 左侧标签文字 */
    label: string;
    /** 左侧描述文字（小字，可选） */
    description?: string;
    /** 右侧操作区域 */
    children: React.ReactNode;
    /** 自定义 className */
    className?: string;
}

/**
 * 设置项行 — 统一 settings 页面中重复的 setting row 样式
 *
 * 结构：[label + description] ←→ [action]
 */
const SettingRow: React.FC<SettingRowProps> = ({label, description, children, className}) => {
    return (
        <View className={cn(
            'flex items-center justify-between rounded-lg border border-border bg-card p-4',
            className,
        )}
        >
            {description ? (
                <View>
                    <View className='text-sm font-medium text-foreground'>{label}</View>
                    <View className='mt-1 text-xs text-muted-foreground'>{description}</View>
                </View>
            ) : (
                <View className='text-sm font-medium text-foreground'>{label}</View>
            )}
            {children}
        </View>
    );
};

export default SettingRow;
