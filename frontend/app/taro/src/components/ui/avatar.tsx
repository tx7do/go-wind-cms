import * as React from 'react';
import {View, Image as TaroImage} from '@tarojs/components';
import {cn} from '@/lib/utils';

const Avatar = React.forwardRef<any, {className?: string; children?: React.ReactNode}>(
    ({className, children, ...props}, ref) => (
        <View
          ref={ref}
          className={cn('relative flex h-10 w-10 shrink-0 overflow-hidden rounded-full', className)}
          {...props}
        >
            {children}
        </View>
    ),
);
Avatar.displayName = 'Avatar';

const AvatarImage = React.forwardRef<any, {className?: string; src?: string; alt?: string}>(
    ({className, src, alt, ...props}, ref) => (
        <TaroImage
          ref={ref}
          src={src || ''}
          className={cn('aspect-square h-full w-full', className)}
          mode='aspectFill'
          {...props}
        />
    ),
);
AvatarImage.displayName = 'AvatarImage';

const AvatarFallback = React.forwardRef<any, {className?: string; children?: React.ReactNode}>(
    ({className, children, ...props}, ref) => (
        <View
          ref={ref}
          className={cn('flex h-full w-full items-center justify-center rounded-full bg-muted', className)}
          {...props}
        >
            {children}
        </View>
    ),
);
AvatarFallback.displayName = 'AvatarFallback';

export {Avatar, AvatarImage, AvatarFallback};
