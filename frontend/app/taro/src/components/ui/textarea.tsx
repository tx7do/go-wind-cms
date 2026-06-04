import * as React from 'react';
import {View} from '@tarojs/components';
import {cn} from '@/lib/utils';

interface TextareaProps {
    className?: string;
    value?: string;
    placeholder?: string;
    disabled?: boolean;
    onInput?: (e: any) => void;
    [key: string]: any;
}

const Textarea = React.forwardRef<any, TextareaProps>(
    ({className, ...props}, ref) => {
        return (
            <View
              className={cn(
                    'flex min-h-15 w-full rounded-md border border-input bg-transparent px-3 py-2 text-base shadow-sm placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-1 focus-visible:ring-ring disabled:cursor-not-allowed disabled:opacity-50 md:text-sm',
                    className,
                )}
              ref={ref}
              {...props}
            />
        );
    },
);
Textarea.displayName = 'Textarea';

export {Textarea};
