import {useState, useMemo} from 'react';
import {useTranslation} from 'react-i18next';
import {View, Text} from '@tarojs/components';
import XIcon from '@/plugins/xicon';

export default function SettingsPage() {
    const {t} = useTranslation();
    const [activeMenu, setActiveMenu] = useState<'account' | 'message' | 'preference'>('account');

    const menuItems = useMemo(() => [
        {key: 'account' as const, icon: 'carbon:user', label: t('settings.menu.account')},
        {key: 'message' as const, icon: 'carbon:email', label: t('settings.menu.message')},
        {key: 'preference' as const, icon: 'carbon:settings', label: t('settings.menu.preference')},
    ], [t]);

    const sections = useMemo(() => {
        switch (activeMenu) {
            case 'account': return {
                title: t('settings.account.title'),
                items: [
                    {label: t('settings.account.password'), status: t('settings.account.password_not_set')},
                    {label: t('settings.account.bind_phone'), status: t('settings.account.password_not_set')},
                    {label: t('settings.account.bind_email'), status: t('settings.account.email_not_bound')},
                ],
            };
            case 'message': return {
                title: t('settings.message.title'),
                items: [
                    {label: t('settings.message.system_messages'), status: t('settings.message.enabled')},
                    {label: t('settings.message.comment_notifications'), status: t('settings.message.enabled')},
                    {label: t('settings.message.activity_updates'), status: t('settings.message.enabled')},
                ],
            };
            case 'preference': return {
                title: t('settings.preference.title'),
                items: [
                    {label: t('settings.preference.theme'), status: ''},
                    {label: t('settings.preference.language'), status: ''},
                    {label: t('settings.preference.hide_sensitive_content'), status: ''},
                ],
            };
        }
    }, [activeMenu, t]);

    return (
        <View className='w-full bg-pageBg'>
            {/* Tab 切换 */}
            <View className='flex bg-cardBg border-b-[1rpx] border-splitLine'>
                {menuItems.map((item) => (
                    <View
                        key={item.key}
                        className={`flex-1 flex items-center justify-center gap-[8rpx] py-[24rpx] min-h-touch ${activeMenu === item.key ? 'border-b-[4rpx] border-primary' : ''}`}
                        onClick={() => setActiveMenu(item.key)}
                    >
                        <XIcon name={item.icon} size={16} className={activeMenu === item.key ? 'text-primary' : 'text-textThird'} />
                        <Text className={`text-desc ${activeMenu === item.key ? 'text-primary font-bold' : 'text-textSec'}`}>
                            {item.label}
                        </Text>
                    </View>
                ))}
            </View>

            {/* 设置项列表 */}
            <View className='px-[24rpx] py-[32rpx]'>
                <Text className='text-card-title font-bold text-textMain block mb-[24rpx]'>{sections?.title}</Text>
                <View className='rounded bg-cardBg overflow-hidden'>
                    {sections?.items.map((item, index) => (
                        <View
                            key={index}
                            className={`flex items-center justify-between px-[24rpx] h-[88rpx] min-h-touch ${index < sections.items.length - 1 ? 'border-b-[1rpx] border-splitLine' : ''}`}
                            hoverClass='tap-active'
                        >
                            <Text className='text-body text-textMain'>{item.label}</Text>
                            <View className='flex items-center gap-[8rpx]'>
                                {item.status && <Text className='text-desc text-textThird'>{item.status}</Text>}
                                <XIcon name='carbon:chevron-right' size={16} className='text-textWeak' />
                            </View>
                        </View>
                    ))}
                </View>
            </View>
        </View>
    );
}
