import {useState} from 'react';
import {View} from '@tarojs/components';
import {Input} from '@/components/ui/input';
import {XIcon} from '@/plugins/xicon';
import {useTranslations} from '@/lib/next-intl-compat';

export default function SearchBar() {
    const t = useTranslations('navbar.top');
    const [searchQuery, setSearchQuery] = useState('');

    return (
        <View className='px-[24rpx] py-[16rpx]'>
            <View className='relative w-full'>
                <View className='absolute left-[16rpx] top-1/2 -translate-y-1/2'>
                    <XIcon name='carbon:search' size={16} className='text-textThird' />
                </View>
                <Input
                  className='w-full h-[80rpx] rounded pl-[56rpx] bg-cardBg border-[1rpx] border-splitLine'
                  value={searchQuery}
                  onInput={(e: any) => setSearchQuery(e.detail?.value ?? e.target?.value ?? '')}
                  placeholder={t('search_placeholder')}
                />
            </View>
        </View>
    );
}
