import React, {useState} from 'react';
import {View, Text, Input} from '@tarojs/components';

import './index.scss';

interface VerificationCodeInputProps {
  /** 验证码长度，默认为 6 */
  length?: number;
  /** 受控模式的值 */
  value?: string;
  /** 值变化时的回调 */
  onChange?: (value: string) => void;
  /** 完成输入时的回调 */
  onComplete?: (value: string) => void;
  /** 是否禁用 */
  disabled?: boolean;
  /** 自定义类名 */
  className?: string;
  /** 输入框类型 */
  type?: 'number' | 'text';
}

/**
 * 验证码输入组件
 * 基于 Taro Input 实现
 */
export const VerificationCodeInput: React.FC<VerificationCodeInputProps> = ({
                                                                              length = 6,
                                                                              value = '',
                                                                              onChange,
                                                                              onComplete,
                                                                              className = '',
                                                                              type = 'number',
                                                                            }) => {
  const [focusedIndex, setFocusedIndex] = useState(-1);

  const handleChange = (e: any) => {
    const newValue = e.detail.value;
    // 只取最后一个字符
    if (newValue.length > 0) {
      const lastChar = newValue.slice(-1);
      const currentIndex = Math.min(newValue.length - 1, length - 1);

      // 更新值
      const values = Array(length).fill('');
      values[currentIndex] = lastChar;
      const result = values.join('');

      onChange?.(result);

      // 当输入完成时触发回调
      if (result.length === length && onComplete) {
        onComplete(result);
      }

      // 自动聚焦到下一个输入框
      if (currentIndex < length - 1) {
        setFocusedIndex(currentIndex + 1);
      }
    }
  };

  const handleFocus = () => {
    setFocusedIndex(0);
  };

  const handleBlur = () => {
    setFocusedIndex(-1);
  };

  // 生成数组用于渲染
  const slots = Array.from({length}, (_, i) => value[i] || '');

  return (
    <View className={`verification-code-input ${className}`}>
      <View className='otp-slots'>
        {slots.map((slot, index) => (
          <View
            key={index}
            className={`otp-slot ${focusedIndex === index ? 'focused' : ''} ${slot ? 'filled' : ''}`}
            onClick={() => {
              setFocusedIndex(index);
              // 触发隐藏 Input 的聚焦
              const hiddenInput = document.querySelector('.hidden-otp-input') as HTMLInputElement;
              if (hiddenInput) {
                hiddenInput.focus();
              }
            }}
          >
            <Text>{slot}</Text>
          </View>
        ))}
      </View>
      {/* 隐藏的 Input 用于接收输入 */}
      <Input
        className='hidden-otp-input'
        type={type === 'number' ? 'number' : 'text'}
        maxlength={length}
        onInput={handleChange}
        onFocus={handleFocus}
        onBlur={handleBlur}
        style={{
          position: 'absolute',
          opacity: 0,
          width: '1px',
          height: '1px',
          pointerEvents: 'none'
        }}
      />
    </View>
  );
};

export default VerificationCodeInput;
