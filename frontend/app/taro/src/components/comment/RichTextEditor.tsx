import React from 'react';
import {View} from '@tarojs/components';
import {cn} from '@/lib/utils';
import {XIcon} from '@/plugins/xicon';

interface RichTextEditorProps {
    value: string;
    onChange: (content: string) => void;
    onSubmit: () => void;
    submitting?: boolean;
    placeholder?: string;
    maxLength?: number;
    submitLabel?: string;
}

/**
 * 简化版评论编辑器（替代 Tiptap）
 * 使用 textarea 而非富文本编辑器，兼容小程序
 */
export default function RichTextEditor({
    value,
    onChange,
    onSubmit,
    submitting = false,
    placeholder = '写下你的评论...',
    maxLength = 1000,
    submitLabel = '提交评论',
}: RichTextEditorProps) {
    const textLength = value.length;
    const isEmpty = textLength === 0;

    return (
        <View className='w-full overflow-hidden rounded-xl border border-border bg-card'>
            {/* 编辑区域 */}
            <View className='p-3'>
                <textarea
                  value={value}
                  onChange={(e) => onChange(e.target.value)}
                  placeholder={placeholder}
                  maxLength={maxLength}
                  className={cn(
                        'w-full min-h-[120px] max-h-[300px] text-sm bg-transparent',
                        'text-foreground placeholder:text-muted-foreground',
                        'focus:outline-none resize-none',
                    )}
                />
            </View>

            {/* 底部操作栏 */}
            <View className='flex items-center justify-between border-t border-border/50 bg-muted/20 p-2'>
                <View className={cn('text-xs', textLength > maxLength ? 'text-destructive' : 'text-muted-foreground')}>
                    {textLength} / {maxLength}
                </View>
                <View
                  onClick={() => {
                        if (!isEmpty && !submitting && textLength <= maxLength) {
                            onSubmit();
                        }
                    }}
                  className={cn(
                        'inline-flex items-center gap-2 rounded-md px-4 py-1.5 text-sm font-medium',
                        'bg-primary text-primary-foreground shadow-sm',
                        (submitting || isEmpty || textLength > maxLength) && 'opacity-50 pointer-events-none',
                    )}
                >
                    {submitting ? (
                        <XIcon name='carbon:circle-dash' className='animate-spin' size={16} />
                    ) : (
                        <XIcon name='carbon:send-alt' size={16} />
                    )}
                    {submitLabel}
                </View>
            </View>
        </View>
    );
}
