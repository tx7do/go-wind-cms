import {View, Text, Textarea} from '@tarojs/components';
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
 * 评论编辑器 - Taro 跨端兼容版
 * 使用 Taro Textarea 组件替代 HTML textarea
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
    const isOverLimit = textLength > maxLength;

    return (
        <View className='w-full overflow-hidden rounded-[16rpx] border-[1rpx] border-splitLine bg-cardBg'>
            {/* 编辑区域 */}
            <View className='p-[16rpx]'>
                <Textarea
                  value={value}
                  onInput={(e) => onChange(e.detail.value ?? '')}
                  placeholder={placeholder}
                  maxlength={maxLength}
                  className='w-full min-h-[160rpx] text-body text-textMain bg-transparent placeholder:text-textWeak'
                  style={{lineHeight: 1.6}}
                  disabled={submitting}
                />
            </View>

            {/* 底部操作栏 */}
            <View className='flex items-center justify-between border-t-[1rpx] border-splitLine/50 bg-pageBg/20 px-[16rpx] py-[12rpx]'>
                <Text className={cn('text-tips', isOverLimit ? 'text-danger' : 'text-textSec')}>
                    {textLength} / {maxLength}
                </Text>
                <View
                  className={cn(
                    'flex items-center gap-[8rpx] rounded-[8rpx] px-[24rpx] py-[12rpx] text-desc font-medium',
                    'bg-primary text-white',
                    (submitting || isEmpty || isOverLimit) && 'opacity-50',
                  )}
                  onClick={() => {
                    if (!isEmpty && !submitting && !isOverLimit) {
                        onSubmit();
                    }
                  }}
                >
                    {submitting ? (
                        <XIcon name='carbon:circle-dash' className='animate-spin' size={14} />
                    ) : (
                        <XIcon name='carbon:send-alt' size={14} />
                    )}
                    <Text className='text-white'>{submitLabel}</Text>
                </View>
            </View>
        </View>
    );
}
