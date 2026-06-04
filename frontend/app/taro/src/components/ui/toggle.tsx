import * as React from 'react';
import {View} from '@tarojs/components';
import {cva, type VariantProps} from 'class-variance-authority';
import {cn} from '@/lib/utils';

const toggleVariants = cva(
    'inline-flex items-center justify-center gap-2 rounded-md text-sm font-medium transition-colors hover:bg-muted hover:text-muted-foreground disabled:pointer-events-none disabled:opacity-50',
    {
        variants: {
            variant: {
                default: 'bg-transparent',
                outline: 'border border-input bg-transparent shadow-sm hover:bg-muted hover:text-primary'
            },
            size: {
                default: 'h-9 px-2 min-w-9',
                sm: 'h-8 px-1.5 min-w-8',
                lg: 'h-10 px-2.5 min-w-10',
            },
        },
        defaultVariants: {
            variant: 'default',
            size: 'default',
        },
    }
);

interface ToggleProps extends VariantProps<typeof toggleVariants> {
    className?: string;
    pressed?: boolean;
    onPressedChange?: (pressed: boolean) => void;
    disabled?: boolean;
    children?: React.ReactNode;
}

function Toggle({className, variant, size, pressed, onPressedChange, disabled, children}: ToggleProps) {
    return (
        <View
          className={cn(
                toggleVariants({variant, size, className}),
                pressed && 'bg-slate-900 text-white dark:bg-slate-100 dark:text-slate-950 dark:font-semibold',
            )}
          onClick={() => {
                if (!disabled) onPressedChange?.(!pressed);
            }}
        >
            {children}
        </View>
    );
}

Toggle.displayName = 'Toggle';

export {Toggle, toggleVariants};
