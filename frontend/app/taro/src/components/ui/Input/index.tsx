import {View, Input as TaroInput, type InputProps as TaroInputProps} from '@tarojs/components';
import './index.scss';

interface InputProps {
  value?: string;
  placeholder?: string;
  type?: TaroInputProps['type'];
  disabled?: boolean;
  error?: boolean;
  onChange?: (value: string) => void;
  onBlur?: () => void;
  className?: string;
}

export function Input({
                        value = '',
                        placeholder = '',
                        type = 'text',
                        disabled = false,
                        error = false,
                        onChange,
                        onBlur,
                        className = '',
                      }: InputProps) {
  const handleChange = (e: any) => {
    if (onChange) {
      onChange(e.detail.value);
    }
  };

  const handleBlur = () => {
    if (onBlur) {
      onBlur();
    }
  };

  return (
    <View className={`taro-input ${error ? 'taro-input--error' : ''} ${className}`}>
      <TaroInput
        className='taro-input__field'
        value={value}
        placeholder={placeholder}
        type={type}
        disabled={disabled}
        onInput={handleChange}
        onBlur={handleBlur}
      />
    </View>
  );
}
