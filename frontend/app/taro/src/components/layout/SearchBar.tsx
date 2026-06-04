import {useState} from 'react';
import {View} from '@tarojs/components';
import {Input} from '@/components/ui/input';
import {XIcon} from '@/plugins/xicon';
import {useTranslations} from '@/lib/next-intl-compat';

export default function SearchBar() {
    const t = useTranslations('navbar.top');
    const [searchQuery, setSearchQuery] = useState('');

    const handleSearch = () => {
        console.log('Searching for:', searchQuery);
    };

    return (
        <View className='mx-2 hidden h-11 max-w-80 flex-1 md:flex lg:max-w-80'>
            <View className='relative w-full'>
                <View className='absolute left-2 top-1/2 -translate-y-1/2 text-muted-foreground'>
                    <XIcon name='carbon:search' size={16} />
                </View>
                <Input
                  className='h-full w-full pl-8'
                  value={searchQuery}
                  onInput={(e: any) => setSearchQuery(e.detail?.value ?? e.target?.value ?? '')}
                  placeholder={t('search_placeholder')}
                />
            </View>
        </View>
    );
}
