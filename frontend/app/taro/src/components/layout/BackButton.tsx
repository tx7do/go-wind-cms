import React from 'react';
import {View} from '@tarojs/components';
import {cn} from '@/lib/utils';
import {XIcon} from '@/plugins/xicon';

export interface BackButtonProps {
    label: string;
    onClick: () => void;
    className?: string;
}

const BackButton: React.FC<BackButtonProps> = ({label, onClick, className}) => {
    return (
        <View
          onClick={onClick}
          className={cn(
                'group flex cursor-pointer items-center gap-2 rounded-lg px-3 py-1.5',
                'text-sm text-muted-foreground transition-colors',
                'hover:bg-muted/60 hover:text-foreground',
                className,
            )}
        >
            <XIcon name='carbon:arrow-left' size={16} className='transition-transform duration-300 group-hover:-translate-x-1' />
            <View>{label}</View>
        </View>
    );
};

export default BackButton;
