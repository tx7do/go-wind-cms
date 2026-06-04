import * as React from 'react';
import {Input as TaroInput} from '@tarojs/components';
import {cn} from '@/lib/utils';

interface InputProps {
    className?: string;
    type?: 'text' | 'number' | 'idcard' | 'digit' | 'safe-password' | 'nickname';
    value?: string;
    placeholder?: string;
    disabled?: boolean;
    onInput?: (e: any) => void;
    onFocus?: (e: any) => void;
    onBlur?: (e: any) => void;
}

const Input = React.forwardRef<any, InputProps>(
    ({className, type, ...props}, ref) => {
        return (
            <TaroInput
              type={type}
              className={cn(
                    'flex h-9 w-full rounded-md border border-input bg-transparent px-3 py-1 text-base shadow-sm transition-colors placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-1 focus-visible:ring-ring disabled:cursor-not-allowed disabled:opacity-50 md:text-sm',
                    className,
                )}
              ref={ref}
              {...props}
            />
        );
    },
);
Input.displayName = 'Input';

export {Input};
