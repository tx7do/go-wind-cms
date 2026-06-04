import * as React from 'react';
import {View} from '@tarojs/components';
import {cn} from '@/lib/utils';

/* =================== DropdownMenu =================== */
interface DropdownMenuProps {
    open?: boolean;
    onOpenChange?: (open: boolean) => void;
    modal?: boolean;
    children: React.ReactNode;
}

const DropdownMenuContext = React.createContext<{
    open: boolean;
    setOpen: (open: boolean) => void;
}>({open: false, setOpen: () => {}});

function DropdownMenu({open: controlledOpen, onOpenChange, children}: DropdownMenuProps) {
    const [internalOpen, setInternalOpen] = React.useState(false);
    const open = controlledOpen ?? internalOpen;
    const setOpen = (v: boolean) => {
        setInternalOpen(v);
        onOpenChange?.(v);
    };

    return (
        <DropdownMenuContext.Provider value={{open, setOpen}}>
            {children}
        </DropdownMenuContext.Provider>
    );
}

function DropdownMenuTrigger({children, className}: {
    children: React.ReactNode;
    className?: string;
}) {
    const {open, setOpen} = React.useContext(DropdownMenuContext);
    return (
        <View className={className} onClick={() => setOpen(!open)}>
            {children}
        </View>
    );
}

function DropdownMenuContent({children, className}: {
    children: React.ReactNode;
    className?: string;
}) {
    const {open, setOpen} = React.useContext(DropdownMenuContext);

    if (!open) return null;

    return (
        <View className='fixed inset-0 z-1001' onClick={() => setOpen(false)}>
            <View
              className={cn(
                    'absolute z-1002 min-w-32 overflow-hidden rounded-md border bg-popover p-1 text-popover-foreground shadow-md',
                    className,
                )}
              onClick={(e: any) => e.stopPropagation()}
            >
                {children}
            </View>
        </View>
    );
}

function DropdownMenuItem({children, className, onClick, inset}: {
    children: React.ReactNode;
    className?: string;
    onClick?: () => void;
    inset?: boolean;
}) {
    const {setOpen} = React.useContext(DropdownMenuContext);
    return (
        <View
          className={cn(
                'relative flex cursor-default select-none items-center gap-2 rounded-sm px-2 py-1.5 text-sm outline-none transition-colors hover:bg-muted hover:text-primary',
                inset && 'pl-8',
                className,
            )}
          onClick={() => {
                onClick?.();
                setOpen(false);
            }}
        >
            {children}
        </View>
    );
}

function DropdownMenuSeparator({className}: {className?: string}) {
    return <View className={cn('-mx-1 my-1 h-px bg-muted', className)} />;
}

function DropdownMenuLabel({children, className, inset}: {
    children: React.ReactNode;
    className?: string;
    inset?: boolean;
}) {
    return (
        <View className={cn('px-2 py-1.5 text-sm font-semibold', inset && 'pl-8', className)}>
            {children}
        </View>
    );
}

function DropdownMenuGroup({children}: {children: React.ReactNode}) {
    return <View>{children}</View>;
}

function DropdownMenuPortal({children}: {children: React.ReactNode}) {
    return <>{children}</>;
}

// 不常用但需要导出占位
const DropdownMenuSub = ({children}: {children: React.ReactNode}) => <>{children}</>;
const DropdownMenuSubTrigger = ({children}: {children: React.ReactNode}) => <View>{children}</View>;
const DropdownMenuSubContent = ({children}: {children: React.ReactNode}) => <View>{children}</View>;
const DropdownMenuRadioGroup = ({children}: {children: React.ReactNode}) => <View>{children}</View>;
const DropdownMenuCheckboxItem = ({children}: {children: React.ReactNode}) => <View>{children}</View>;
const DropdownMenuRadioItem = ({children}: {children: React.ReactNode}) => <View>{children}</View>;
const DropdownMenuShortcut = ({children, className}: {children: React.ReactNode; className?: string}) => (
    <View className={cn('ml-auto text-xs tracking-widest opacity-60', className)}>{children}</View>
);

export {
    DropdownMenu,
    DropdownMenuTrigger,
    DropdownMenuContent,
    DropdownMenuItem,
    DropdownMenuCheckboxItem,
    DropdownMenuRadioItem,
    DropdownMenuLabel,
    DropdownMenuSeparator,
    DropdownMenuShortcut,
    DropdownMenuGroup,
    DropdownMenuPortal,
    DropdownMenuSub,
    DropdownMenuSubContent,
    DropdownMenuSubTrigger,
    DropdownMenuRadioGroup,
};
