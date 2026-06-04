import {View} from '@tarojs/components';
import {cn} from '@/lib/utils';

interface SwitchProps {
    checked?: boolean;
    onCheckedChange?: (checked: boolean) => void;
    disabled?: boolean;
    className?: string;
}

function Switch({checked, onCheckedChange, disabled, className}: SwitchProps) {
    return (
        <View
          className={cn(
                'inline-flex h-5 w-9 shrink-0 cursor-pointer items-center rounded-full border-2 border-transparent shadow-sm transition-colors',
                checked ? 'bg-primary' : 'bg-input',
                disabled && 'cursor-not-allowed opacity-50',
                className,
            )}
          onClick={() => {
                if (!disabled) onCheckedChange?.(!checked);
            }}
        >
            <View
              className={cn(
                    'pointer-events-none block h-4 w-4 rounded-full bg-background shadow-lg transition-transform',
                    checked ? 'translate-x-4' : 'translate-x-0',
                )}
            />
        </View>
    );
}

export {Switch};
