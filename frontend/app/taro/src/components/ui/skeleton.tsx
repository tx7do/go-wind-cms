import * as React from 'react';
import {View} from '@tarojs/components';
import {cn} from '@/lib/utils';

function Skeleton({className, ...props}: React.HTMLAttributes<HTMLDivElement>) {
    return <View className={cn('animate-pulse rounded-md bg-muted', className)} {...props} />;
}

export {Skeleton};
