import * as React from 'react';
import {View} from '@tarojs/components';
import {cn} from '@/lib/utils';
import {XIcon} from '@/plugins/xicon';

/* =================== Sheet Context =================== */
const SheetContext = React.createContext<{
    open: boolean;
    setOpen: (v: boolean) => void;
}>({open: false, setOpen: () => {}});

/* =================== Sheet (侧滑面板) =================== */
interface SheetProps {
    open?: boolean;
    onOpenChange?: (open: boolean) => void;
    children: React.ReactNode;
}

function Sheet({open: controlledOpen, onOpenChange, children}: SheetProps) {
    const [internalOpen, setInternalOpen] = React.useState(false);
    const open = controlledOpen ?? internalOpen;
    const setOpen = (v: boolean) => {
        setInternalOpen(v);
        onOpenChange?.(v);
    };

    if (!open) return null;

    return (
        <SheetContext.Provider value={{open, setOpen}}>
            {children}
        </SheetContext.Provider>
    );
}

function SheetTrigger({children, asChild, className}: {
    children: React.ReactNode;
    asChild?: boolean;
    className?: string;
}) {
    return <View className={className}>{children}</View>;
}

function SheetClose({children}: {children: React.ReactNode}) {
    return <>{children}</>;
}

function SheetPortal({children}: {children: React.ReactNode}) {
    return <>{children}</>;
}

function SheetOverlay({className}: {className?: string}) {
    const {setOpen} = React.useContext(SheetContext);
    return (
        <View
          className={cn(
                'fixed inset-0 z-2000 bg-black/60',
                className,
            )}
          onClick={() => setOpen(false)}
        />
    );
}

interface SheetContentProps {
    side?: 'top' | 'right' | 'bottom' | 'left';
    hideCloseButton?: boolean;
    className?: string;
    children: React.ReactNode;
}

function SheetContent({side = 'right', hideCloseButton, className, children}: SheetContentProps) {
    const {setOpen} = React.useContext(SheetContext);
    const positionClass = {
        top: 'inset-x-0 top-0 border-b',
        right: 'inset-y-0 right-0 h-full w-full border-l sm:max-w-sm',
        bottom: 'inset-x-0 bottom-0 border-t',
        left: 'inset-y-0 left-0 h-full w-full border-r sm:max-w-sm',
    }[side];

    return (
        <>
            <SheetOverlay />
            <View
              className={cn(
                    'fixed z-2001 gap-4 bg-background p-6 shadow-lg',
                    positionClass,
                    className,
                )}
            >
                {children}
                {!hideCloseButton && (
                    <View
                      className={cn(
                            'absolute right-4 top-4 rounded-sm opacity-70 cursor-pointer',
                        )}
                      onClick={() => setOpen(false)}
                    >
                        <XIcon name='carbon:close' size={16} />
                    </View>
                )}
            </View>
        </>
    );
}

function SheetHeader({className, ...props}: {className?: string; children?: React.ReactNode; [key: string]: any}) {
    return (
        <View className={cn('flex flex-col space-y-2 text-center sm:text-left', className)} {...props} />
    );
}
SheetHeader.displayName = 'SheetHeader';

function SheetFooter({className, ...props}: {className?: string; children?: React.ReactNode; [key: string]: any}) {
    return (
        <View className={cn('flex flex-col-reverse sm:flex-row sm:justify-end sm:space-x-2', className)} {...props} />
    );
}
SheetFooter.displayName = 'SheetFooter';

function SheetTitle({className, children}: {className?: string; children?: React.ReactNode}) {
    return (
        <View className={cn('text-lg font-semibold text-foreground', className)}>
            {children}
        </View>
    );
}
SheetTitle.displayName = 'SheetTitle';

function SheetDescription({className, children}: {className?: string; children?: React.ReactNode}) {
    return (
        <View className={cn('text-sm text-muted-foreground', className)}>
            {children}
        </View>
    );
}
SheetDescription.displayName = 'SheetDescription';

export {
    Sheet,
    SheetTrigger,
    SheetClose,
    SheetContent,
    SheetHeader,
    SheetFooter,
    SheetTitle,
    SheetDescription,
    SheetPortal,
    SheetOverlay,
};
