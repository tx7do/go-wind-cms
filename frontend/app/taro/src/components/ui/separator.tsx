import * as React from 'react';
import {View} from '@tarojs/components';
import {cn} from '@/lib/utils';

interface SeparatorProps {
    className?: string;
    orientation?: 'horizontal' | 'vertical';
    decorative?: boolean;
}

const Separator = React.forwardRef<any, SeparatorProps>(
    ({className, orientation = 'horizontal', decorative = true}, ref) => (
        <View
          ref={ref}
          className={cn(
                'shrink-0 bg-border',
                orientation === 'horizontal' ? 'h-px w-full' : 'h-full w-px',
                className,
            )}
        />
    ),
);
Separator.displayName = 'Separator';

export {Separator};
