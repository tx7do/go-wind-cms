'use client';

import React from 'react';
import {Input} from 'antd';

import styles from './VerificationCodeInput.module.css';

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
 * 基于 Ant Design Input.OTP 实现
 */
export const VerificationCodeInput: React.FC<VerificationCodeInputProps> = ({
                                                                                length = 6,
                                                                                value = '',
                                                                                onChange,
                                                                                onComplete,
                                                                                disabled = false,
                                                                                className = '',
                                                                                type = 'number',
                                                                            }) => {
    const handleChange = (newValue: string) => {
        onChange?.(newValue);

        // 当输入完成时触发回调
        if (newValue.length === length && onComplete) {
            onComplete(newValue);
        }
    };

    return (
        <div className={`${styles.verificationCodeInput} ${className}`}>
            <Input.OTP
                length={length}
                value={value}
                onChange={handleChange}
                disabled={disabled}
                type={type}
                className={styles.otpInput}
            />
        </div>
    );
};

export default VerificationCodeInput;
