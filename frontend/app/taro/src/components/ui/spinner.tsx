import {View} from '@tarojs/components';
import {cn} from '@/lib/utils';
import {XIcon} from '@/plugins/xicon';

interface SpinnerProps {
    className?: string;
    size?: 'sm' | 'md' | 'lg';
}

const sizeMap = {
    sm: 16,
    md: 24,
    lg: 32,
};

function Spinner({className, size = 'md', ...props}: SpinnerProps) {
    return (
        <View className={cn('flex items-center justify-center', className)} {...props}>
            <XIcon name='carbon:circle-dash' size={sizeMap[size]} className='animate-spin text-muted-foreground' />
        </View>
    );
}

export {Spinner};
