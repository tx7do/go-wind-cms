import * as React from 'react';
import {View, Text} from '@tarojs/components';
import {cva, type VariantProps} from 'class-variance-authority';

import {cn} from '@/lib/utils';

/**
 * 小程序友好按钮组件
 * - 所有尺寸使用 rpx
 * - hover 反馈通过 hoverClass='tap-active' 实现
 * - 移除 hover:/focus-visible: 等伪类
 */
const buttonVariants = cva(
    'flex items-center justify-center whitespace-nowrap rounded font-medium',
    {
        variants: {
            variant: {
                default: 'bg-primary text-white',
                hero: 'btn-hero',
                glass: 'btn-glass',
                destructive: 'bg-danger text-white',
                outline: 'border-[length:2rpx] border-splitLine bg-cardBg text-textMain',
                secondary: 'bg-pageBg text-textSec',
                ghost: 'bg-transparent text-textSec',
                link: 'bg-transparent text-primary underline',
            },
            size: {
                default: 'min-h-touch px-[32rpx] py-[16rpx] text-body',
                sm: 'min-h-[64rpx] px-[24rpx] py-[8rpx] text-desc rounded-sm',
                lg: 'min-h-[96rpx] px-[48rpx] py-[20rpx] text-card-title rounded-lg',
                icon: 'min-w-touch min-h-touch',
            },
        },
        defaultVariants: {
            variant: 'default',
            size: 'default',
        },
    },
);

export interface ButtonProps extends VariantProps<typeof buttonVariants> {
    className?: string;
    children?: React.ReactNode;
    onClick?: (e: any) => void;
    disabled?: boolean;
    /** 提交中：保留本 variant 的底色并加转环，与 disabled 的「中性灰不可点」区分开 */
    loading?: boolean;
    type?: 'button' | 'submit' | 'reset';
    variant?: 'default' | 'hero' | 'glass' | 'destructive' | 'outline' | 'secondary' | 'ghost' | 'link';
    size?: 'default' | 'sm' | 'lg' | 'icon';
    hoverClass?: string;
}

const Button = React.forwardRef<any, ButtonProps>(
    ({className, variant, size, children, disabled, loading, hoverClass, ...props}, ref) => {
        const inert = disabled || loading;
        return (
            <View
              className={cn(
                    buttonVariants({variant, size, className}),
                    loading && 'gap-[12rpx]',
                    disabled && !loading && 'au-btn-disabled',
                )}
              ref={ref}
              hoverClass={inert ? '' : (hoverClass ?? 'tap-active')}
              onClick={inert ? undefined : props.onClick}
              {...props}
            >
                {loading && <View className='au-btn-spin' />}
                {typeof children === 'string' ? (
                    <Text className='text-inherit'>{children}</Text>
                ) : (
                    children
                )}
            </View>
        );
    },
);
Button.displayName = 'Button';

export {Button, buttonVariants};
