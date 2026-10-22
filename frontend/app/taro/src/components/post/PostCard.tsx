import {View, Text, Image} from '@tarojs/components';
import XIcon from '@/plugins/xicon';
import {cn} from '@/lib/utils';
import {formatDate} from '@/utils/date';

export interface PostCardProps {
    title: string;
    summary?: string;
    thumbnail?: string;
    author?: string;
    publishTime?: string;
    views?: number;
    className?: string;
    onClick?: () => void;
}

/**
 * 文章卡片（信息流通用）：左侧标题/摘要/元信息，右侧封面。
 * 无封面时右侧退化为品牌色渐变占位块，保持栅格节奏不塌。
 */
export default function PostCard({
    title,
    summary,
    thumbnail,
    author,
    publishTime,
    views,
    className,
    onClick,
}: PostCardProps) {
    return (
        <View
          className={cn('au-card au-card--glow overflow-hidden p-[22rpx] flex gap-[22rpx]', className)}
          hoverClass='tap-active'
          onClick={onClick}
        >
            {/* 文字主体 */}
            <View className='flex-1 min-w-0 flex flex-col'>
                <Text className='text-body font-bold text-textMain leading-[1.4]' numberOfLines={2}>
                    {title}
                </Text>
                {summary && (
                    <Text className='text-tips text-textThird leading-[1.6] mt-[10rpx]' numberOfLines={2}>
                        {summary}
                    </Text>
                )}
                <View className='flex items-center gap-[18rpx] mt-auto pt-[16rpx]'>
                    {author && (
                        <View className='flex items-center gap-[6rpx] min-w-0'>
                            <XIcon name='carbon:user' size={18} className='text-textWeak shrink-0' />
                            <Text className='text-tips text-textThird' numberOfLines={1}>{author}</Text>
                        </View>
                    )}
                    {publishTime && (
                        <View className='flex items-center gap-[6rpx] shrink-0'>
                            <XIcon name='carbon:calendar' size={18} className='text-textWeak' />
                            <Text className='text-tips text-textThird'>{formatDate(publishTime)}</Text>
                        </View>
                    )}
                    {typeof views === 'number' && views >= 0 && (
                        <View className='flex items-center gap-[6rpx] shrink-0'>
                            <XIcon name='carbon:view' size={18} className='text-textWeak' />
                            <Text className='text-tips text-textThird'>{views}</Text>
                        </View>
                    )}
                </View>
            </View>

            {/* 封面 / 渐变占位 */}
            <View className='w-[200rpx] h-[150rpx] rounded-[18rpx] overflow-hidden shrink-0'>
                {thumbnail ? (
                    <Image
                      src={thumbnail}
                      mode='aspectFill'
                      className='w-full h-full'
                      lazyLoad
                    />
                ) : (
                    <View
                      className='w-full h-full flex items-center justify-center'
                      style={{background: 'linear-gradient(128deg, #57a3ff 0%, #1677ff 52%, #0b55d4 100%)'}}
                    >
                        <XIcon name='carbon:document-blank' size={44} className='text-white/70' />
                    </View>
                )}
            </View>
        </View>
    );
}
