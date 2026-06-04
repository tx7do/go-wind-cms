import * as React from 'react';
import {View} from '@tarojs/components';

import {cn} from '@/lib/utils';
import {buttonVariants} from '@/components/ui/button';
import {XIcon} from '@/plugins/xicon';

const Pagination = ({className, ...props}: {className?: string; children?: React.ReactNode; [key: string]: any}) => (
    <View
      className={cn('mx-auto flex w-full justify-center', className)}
      {...props}
    />
);
Pagination.displayName = 'Pagination';

const PaginationContent = ({className, ...props}: {className?: string; children?: React.ReactNode; [key: string]: any}) => (
    <View className={cn('flex flex-row items-center gap-1', className)} {...props} />
);
PaginationContent.displayName = 'PaginationContent';

const PaginationItem = ({className, ...props}: {className?: string; children?: React.ReactNode; [key: string]: any}) => (
    <View className={cn('', className)} {...props} />
);
PaginationItem.displayName = 'PaginationItem';

type PaginationLinkProps = {
    isActive?: boolean;
    size?: 'default' | 'sm' | 'lg' | 'icon';
    className?: string;
    children?: React.ReactNode;
    onClick?: () => void;
};

const PaginationLink = ({
    className,
    isActive,
    size = 'icon',
    children,
    onClick,
}: PaginationLinkProps) => (
    <View
      className={cn(
            buttonVariants({variant: isActive ? 'outline' : 'ghost', size}),
            'h-9 w-9',
            className,
        )}
      onClick={onClick}
    >
        {children}
    </View>
);
PaginationLink.displayName = 'PaginationLink';

const PaginationPrevious = ({className, ...props}: {className?: string; onClick?: () => void}) => (
    <PaginationLink
      size='default'
      className={cn('gap-1 px-2.5', className)}
      onClick={props.onClick}
    >
        <XIcon name='carbon:chevron-left' size={16} />
    </PaginationLink>
);
PaginationPrevious.displayName = 'PaginationPrevious';

const PaginationNext = ({className, ...props}: {className?: string; onClick?: () => void}) => (
    <PaginationLink
      size='default'
      className={cn('gap-1 px-2.5', className)}
      onClick={props.onClick}
    >
        <XIcon name='carbon:chevron-right' size={16} />
    </PaginationLink>
);
PaginationNext.displayName = 'PaginationNext';

const PaginationEllipsis = ({className}: {className?: string}) => (
    <View className={cn('flex h-9 w-9 items-center justify-center', className)}>
        <XIcon name='carbon:overflow-menu-horizontal' size={16} />
    </View>
);
PaginationEllipsis.displayName = 'PaginationEllipsis';

export {
    Pagination,
    PaginationContent,
    PaginationLink,
    PaginationItem,
    PaginationPrevious,
    PaginationNext,
    PaginationEllipsis,
};
