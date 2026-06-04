import * as React from 'react';
import {View} from '@tarojs/components';
import {cn} from '@/lib/utils';
import {XIcon} from '@/plugins/xicon';

/* =================== Select =================== */
interface SelectProps {
    value?: string;
    onValueChange?: (value: string) => void;
    children: React.ReactNode;
}

const SelectContext = React.createContext<{
    value?: string;
    onChange?: (value: string) => void;
    open: boolean;
    setOpen: (v: boolean) => void;
}>({open: false, setOpen: () => {}});

function Select({value, onValueChange, children}: SelectProps) {
    const [open, setOpen] = React.useState(false);
    return (
        <SelectContext.Provider value={{value, onChange: onValueChange, open, setOpen}}>
            {children}
        </SelectContext.Provider>
    );
}

function SelectGroup({children}: {children: React.ReactNode}) {
    return <View>{children}</View>;
}

function SelectValue({placeholder}: {placeholder?: string}) {
    const {value} = React.useContext(SelectContext);
    return <View>{value || placeholder}</View>;
}

function SelectTrigger({children, className}: {children: React.ReactNode; className?: string}) {
    const {open, setOpen} = React.useContext(SelectContext);
    return (
        <View
          className={cn(
                'flex h-9 w-full items-center justify-between rounded-md border border-input bg-transparent px-3 py-2 text-sm shadow-sm',
                className,
            )}
          onClick={() => setOpen(!open)}
        >
            {children}
            <XIcon name='carbon:chevron-down' size={16} className='opacity-50' />
        </View>
    );
}

function SelectContent({children, className, position}: {
    children: React.ReactNode;
    className?: string;
    position?: string;
}) {
    const {open, setOpen} = React.useContext(SelectContext);
    if (!open) return null;

    // position 参数保留用于将来扩展（如 'popper' | 'item-aligned'）
    void position;

    return (
        <View className='fixed inset-0 z-50' onClick={() => setOpen(false)}>
            <View
              className={cn(
                    'absolute z-50 max-h-96 min-w-32 overflow-hidden rounded-md border bg-popover text-popover-foreground shadow-md',
                    className,
                )}
              onClick={(e: any) => e.stopPropagation()}
            >
                <View className='p-1'>
                    {children}
                </View>
            </View>
        </View>
    );
}

function SelectItem({children, value, className}: {
    children: React.ReactNode;
    value: string;
    className?: string;
}) {
    const {value: selectedValue, onChange, setOpen} = React.useContext(SelectContext);
    const isSelected = selectedValue === value;

    return (
        <View
          className={cn(
                'relative flex w-full cursor-default select-none items-center rounded-sm py-1.5 pl-2 pr-8 text-sm outline-none hover:bg-muted hover:text-primary',
                isSelected && 'bg-muted text-primary',
                className,
            )}
          onClick={() => {
                onChange?.(value);
                setOpen(false);
            }}
        >
            {isSelected && (
                <View className='absolute right-2 flex h-3.5 w-3.5 items-center justify-center'>
                    <XIcon name='carbon:checkmark' size={16} />
                </View>
            )}
            {children}
        </View>
    );
}

function SelectLabel({children, className}: {children: React.ReactNode; className?: string}) {
    return <View className={cn('px-2 py-1.5 text-sm font-semibold', className)}>{children}</View>;
}

function SelectSeparator({className}: {className?: string}) {
    return <View className={cn('-mx-1 my-1 h-px bg-muted', className)} />;
}

function SelectScrollUpButton() {
    return null;
}

function SelectScrollDownButton() {
    return null;
}

export {
    Select,
    SelectGroup,
    SelectValue,
    SelectTrigger,
    SelectContent,
    SelectLabel,
    SelectItem,
    SelectSeparator,
    SelectScrollUpButton,
    SelectScrollDownButton,
};
