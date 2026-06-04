import * as React from 'react';
import {View} from '@tarojs/components';
import {cva} from 'class-variance-authority';
import {cn} from '@/lib/utils';
import {XIcon} from '@/plugins/xicon';

const navigationMenuTriggerStyle = cva(
    'inline-flex h-9 w-max items-center justify-center rounded-md bg-transparent px-4 py-2 text-sm font-medium transition-colors hover:bg-muted hover:text-primary focus:bg-muted focus:text-primary focus:outline-none disabled:pointer-events-none disabled:opacity-50'
);

function NavigationMenu({className, children}: {className?: string; children: React.ReactNode}) {
    return (
        <View className={cn('relative z-10 flex max-w-max flex-1 items-center justify-center', className)}>
            {children}
        </View>
    );
}

function NavigationMenuList({className, children}: {className?: string; children: React.ReactNode}) {
    return (
        <View className={cn('flex flex-1 items-center justify-center space-x-1', className)}>
            {children}
        </View>
    );
}

function NavigationMenuItem({children, value}: {children: React.ReactNode; value?: string}) {
    return <View className='relative'>{children}</View>;
}

const NavigationMenuContext = React.createContext<{activeValue: string | null; setActiveValue: (v: string | null) => void}>({
    activeValue: null,
    setActiveValue: () => {},
});

function NavigationMenuTrigger({className, children, value}: {
    className?: string;
    children: React.ReactNode;
    value?: string;
}) {
    const {activeValue, setActiveValue} = React.useContext(NavigationMenuContext);
    const isOpen = activeValue === value;
    return (
        <View          className={cn(navigationMenuTriggerStyle(), 'group', isOpen && 'bg-muted/50', className)}

          onClick={() => setActiveValue(isOpen ? null : (value ?? null))}
        >
            {children}
            <XIcon
              name='carbon:chevron-down'
              size={12}
              className={cn('ml-1 transition-transform duration-300', isOpen && 'rotate-180')}
            />
        </View>
    );
}

function NavigationMenuContent({className, children}: {className?: string; children: React.ReactNode}) {
    return (
        <View className={cn('left-0 top-full w-full md:absolute md:w-auto', className)}>
            {children}
        </View>
    );
}

function NavigationMenuLink({children, className, onClick}: {
    children: React.ReactNode;
    className?: string;
    onClick?: () => void;
}) {
    return (
        <View className={className} onClick={onClick}>
            {children}
        </View>
    );
}

function NavigationMenuViewport({className}: {className?: string}) {
    return null;
}

function NavigationMenuIndicator({className}: {className?: string}) {
    return null;
}

export {
    navigationMenuTriggerStyle,
    NavigationMenu,
    NavigationMenuList,
    NavigationMenuItem,
    NavigationMenuContent,
    NavigationMenuTrigger,
    NavigationMenuLink,
    NavigationMenuIndicator,
    NavigationMenuViewport,
    NavigationMenuContext,
};
