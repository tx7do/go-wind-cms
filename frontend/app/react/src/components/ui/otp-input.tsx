'use client';

import * as React from 'react';
import {cn} from '@/lib/utils';

interface OtpInputProps {
    length?: number;
    value?: string;
    onChange?: (value: string) => void;
    onComplete?: (value: string) => void;
    disabled?: boolean;
    type?: 'number' | 'text';
    className?: string;
}

/**
 * OTP Input component - shadcn/ui style replacement for antd Input.OTP
 */
const OtpInput = React.forwardRef<HTMLDivElement, OtpInputProps>(
    (
        {
            length = 6,
            value = '',
            onChange,
            onComplete,
            disabled = false,
            type = 'number',
            className,
        },
        ref,
    ) => {
        const inputRefs = React.useRef<(HTMLInputElement | null)[]>([]);

        const focusInput = (index: number) => {
            const input = inputRefs.current[index];
            if (input) {
                input.focus();
                input.select();
            }
        };

        const handleInputChange = (index: number, inputValue: string) => {
            const char = inputValue.slice(-1);
            // Validate input based on type
            if (type === 'number' && char && !/^\d$/.test(char)) return;

            const newValue = value.split('');
            newValue[index] = char;
            const joined = newValue.join('').slice(0, length);
            onChange?.(joined);

            // Auto-advance to next input
            if (char && index < length - 1) {
                focusInput(index + 1);
            }

            // Check if complete
            if (joined.length === length && onComplete) {
                onComplete(joined);
            }
        };

        const handleKeyDown = (index: number, e: React.KeyboardEvent<HTMLInputElement>) => {
            if (e.key === 'Backspace') {
                e.preventDefault();
                if (value[index]) {
                    // Clear current
                    const newValue = value.split('');
                    newValue[index] = '';
                    onChange?.(newValue.join(''));
                } else if (index > 0) {
                    // Move to previous and clear
                    focusInput(index - 1);
                    const newValue = value.split('');
                    newValue[index - 1] = '';
                    onChange?.(newValue.join(''));
                }
            } else if (e.key === 'ArrowLeft' && index > 0) {
                e.preventDefault();
                focusInput(index - 1);
            } else if (e.key === 'ArrowRight' && index < length - 1) {
                e.preventDefault();
                focusInput(index + 1);
            }
        };

        const handlePaste = (e: React.ClipboardEvent<HTMLInputElement>) => {
            e.preventDefault();
            const pasted = e.clipboardData.getData('text').slice(0, length);
            if (type === 'number' && !/^\d+$/.test(pasted)) return;
            onChange?.(pasted);
            if (pasted.length === length && onComplete) {
                onComplete(pasted);
            }
            focusInput(Math.min(pasted.length, length - 1));
        };

        return (
            <div ref={ref} className={cn('flex items-center gap-2', className)}>
                {Array.from({length}).map((_, index) => (
                    <input
                        key={index}
                        ref={(el) => { inputRefs.current[index] = el; }}
                        type={type === 'number' ? 'text' : 'text'}
                        inputMode={type === 'number' ? 'numeric' : 'text'}
                        maxLength={1}
                        value={value[index] || ''}
                        onChange={(e) => handleInputChange(index, e.target.value)}
                        onKeyDown={(e) => handleKeyDown(index, e)}
                        onPaste={handlePaste}
                        onFocus={(e) => e.target.select()}
                        disabled={disabled}
                        className="flex h-10 w-10 items-center justify-center rounded-md border border-input bg-transparent text-center text-sm shadow-sm transition-colors placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-1 focus-visible:ring-ring disabled:cursor-not-allowed disabled:opacity-50 md:h-12 md:w-12 md:text-base"
                    />
                ))}
            </div>
        );
    },
);
OtpInput.displayName = 'OtpInput';

export {OtpInput};
