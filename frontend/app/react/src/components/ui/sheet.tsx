'use client';

import * as React from 'react';
import * as DialogPrimitive from '@radix-ui/react-dialog';
import {X} from 'lucide-react';

import {cn} from '@/lib/utils';

const Sheet = DialogPrimitive.Root;
const SheetTrigger = DialogPrimitive.Trigger;
const SheetClose = DialogPrimitive.Close;
const SheetPortal = DialogPrimitive.Portal;

const SheetOverlay = React.forwardRef<
    React.ComponentRef<typeof DialogPrimitive.Overlay>,
    React.ComponentPropsWithoutRef<typeof DialogPrimitive.Overlay>
>(({className, ...props}, ref) => (
    <DialogPrimitive.Overlay
        ref={ref}
        className={cn(
            'fixed inset-0 z-2000 bg-black/60 backdrop-blur-sm',
            'data-[state=open]:animate-in data-[state=closed]:animate-out',
            'data-[state=closed]:fade-out-0 data-[state=open]:fade-in-0',
            className,
        )}
        {...props}
    />
));
SheetOverlay.displayName = DialogPrimitive.Overlay.displayName;

interface SheetContentProps extends React.ComponentPropsWithoutRef<typeof DialogPrimitive.Content> {
    /** 抽屉滑出方向 */
    side?: 'top' | 'right' | 'bottom' | 'left';
    /** 是否隐藏右上角关闭按钮 */
    hideCloseButton?: boolean;
}

const sideAnimClass = {
    top: 'data-[state=closed]:slide-out-to-top data-[state=open]:slide-in-from-top',
    right: 'data-[state=closed]:slide-out-to-right data-[state=open]:slide-in-from-right',
    bottom: 'data-[state=closed]:slide-out-to-bottom data-[state=open]:slide-in-from-bottom',
    left: 'data-[state=closed]:slide-out-to-left data-[state=open]:slide-in-from-left',
};

const SheetContent = React.forwardRef<
    React.ComponentRef<typeof DialogPrimitive.Content>,
    SheetContentProps
>(({side = 'right', className, children, hideCloseButton, ...props}, ref) => {
    const positionClass = {
        top: 'inset-x-0 top-0 border-b',
        right: 'inset-y-0 right-0 h-full w-full border-l sm:max-w-sm',
        bottom: 'inset-x-0 bottom-0 border-t',
        left: 'inset-y-0 left-0 h-full w-full border-r sm:max-w-sm',
    }[side];

    return (
        <SheetPortal>
            <SheetOverlay/>
            <DialogPrimitive.Content
                ref={ref}
                className={cn(
                    'fixed z-2000 gap-4 bg-background p-6 shadow-lg transition ease-in-out',
                    'data-[state=open]:animate-in data-[state=closed]:animate-out',
                    'data-[state=open]:duration-300 data-[state=closed]:duration-200',
                    sideAnimClass[side],
                    positionClass,
                    className,
                )}
                {...props}
            >
                {children}
                {!hideCloseButton && (
                    <DialogPrimitive.Close
                        className={cn(
                            'absolute right-4 top-4 rounded-sm opacity-70 ring-offset-background',
                            'transition-opacity hover:opacity-100 focus:outline-none focus:ring-2 focus:ring-ring focus:ring-offset-2',
                            'disabled:pointer-events-none',
                        )}
                    >
                        <X className="h-4 w-4"/>
                        <span className="sr-only">Close</span>
                    </DialogPrimitive.Close>
                )}
            </DialogPrimitive.Content>
        </SheetPortal>
    );
});
SheetContent.displayName = DialogPrimitive.Content.displayName;

const SheetHeader = ({className, ...props}: React.HTMLAttributes<HTMLDivElement>) => (
    <div className={cn('flex flex-col space-y-2 text-center sm:text-left', className)} {...props}/>
);
SheetHeader.displayName = 'SheetHeader';

const SheetFooter = ({className, ...props}: React.HTMLAttributes<HTMLDivElement>) => (
    <div className={cn('flex flex-col-reverse sm:flex-row sm:justify-end sm:space-x-2', className)} {...props}/>
);
SheetFooter.displayName = 'SheetFooter';

const SheetTitle = React.forwardRef<
    React.ComponentRef<typeof DialogPrimitive.Title>,
    React.ComponentPropsWithoutRef<typeof DialogPrimitive.Title>
>(({className, ...props}, ref) => (
    <DialogPrimitive.Title
        ref={ref}
        className={cn('text-lg font-semibold text-foreground', className)}
        {...props}
    />
));
SheetTitle.displayName = DialogPrimitive.Title.displayName;

const SheetDescription = React.forwardRef<
    React.ComponentRef<typeof DialogPrimitive.Description>,
    React.ComponentPropsWithoutRef<typeof DialogPrimitive.Description>
>(({className, ...props}, ref) => (
    <DialogPrimitive.Description
        ref={ref}
        className={cn('text-sm text-muted-foreground', className)}
        {...props}
    />
));
SheetDescription.displayName = DialogPrimitive.Description.displayName;

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
