import {useState} from 'react';
import {useDidShow} from '@tarojs/taro';
import {useTranslation} from 'react-i18next';
import {View, Text} from '@tarojs/components';
import XIcon from '@/plugins/xicon';
import {usePageTitle} from '@/hooks/usePageTitle';
import {usePreferences} from '@/core/preferences';
import type {SupportedLanguagesType, ThemeModeType} from '@/core/preferences';
import {OptionSheet, type OptionSheetOption} from '@/components/ui';
import {fetchUserProfile} from '@/api/hooks/user-profile';
import type {identityservicev1_User} from '@/api/generated/app/service/v1';
import {useIsLogin} from '@/store/core/access/store';
import {useI18nRouter} from '@/i18n/helpers';

type SheetKey = 'theme' | 'language' | null;

interface SettingRow {
    label: string;
    status: string;
    icon: string;
    onClick?: () => void;
}

/**
 * 设置页：菜单行 + 偏好（主题/语言 OptionSheet 即时预览）。
 * 资料编辑在独立页 /settings/profile（页面底部悬浮保存条），本页只留入口行。
 */
export default function SettingsPage() {
    const {t, i18n} = useTranslation();
    usePageTitle('page.title.settings');
    const {theme, app, setThemeMode, setLanguage} = usePreferences();
    const router = useI18nRouter();
    const isLogin = useIsLogin();
    const [sheet, setSheet] = useState<SheetKey>(null);
    // 「编辑资料」行展示当前昵称状态（从设置页返回时刷新）
    const [profileNickname, setProfileNickname] = useState('');

    useDidShow(() => {
        if (!isLogin) return;
        fetchUserProfile()
            .then((result) => setProfileNickname((result as identityservicev1_User | null)?.nickname || ''))
            .catch((error) => console.error('Load profile for settings failed:', error));
    });

    const themeLabels: Record<ThemeModeType, string> = {
        light: t('settings.preference.theme_light'),
        dark: t('settings.preference.theme_dark'),
        auto: t('settings.preference.theme_auto'),
    };
    const localeLabels: Record<SupportedLanguagesType, string> = {
        'zh-CN': t('settings.preference.language_chinese'),
        'en-US': t('settings.preference.language_english'),
    };

    const themeOptions: OptionSheetOption[] = [
        {value: 'light', label: themeLabels.light, icon: 'carbon:sun'},
        {value: 'dark', label: themeLabels.dark, icon: 'carbon:moon'},
        {value: 'auto', label: themeLabels.auto, icon: 'carbon:screen'},
    ];
    const languageOptions: OptionSheetOption[] = [
        {value: 'zh-CN', label: localeLabels['zh-CN']},
        {value: 'en-US', label: localeLabels['en-US']},
    ];

    const handleSelect = (value: string) => {
        if (sheet === 'theme') {
            setThemeMode(value as ThemeModeType);
        } else if (sheet === 'language') {
            const locale = value as SupportedLanguagesType;
            setLanguage(locale);
            void i18n.changeLanguage(locale);
        }
        setSheet(null);
    };

    const accountItems: SettingRow[] = [
        {label: t('page.user.edit_profile'), status: profileNickname, icon: 'carbon:user-role', onClick: () => router.push(isLogin ? '/settings/profile' : '/login')},
        // 用户资料接口不返回密码是否已设置（每账号必有口令），故不留假状态
        {label: t('settings.account.password'), status: '', icon: 'carbon:locked', onClick: () => router.push(isLogin ? '/settings/password' : '/login')},
    ];

    const preferenceItems: SettingRow[] = [
        {
            label: t('settings.preference.theme'),
            status: themeLabels[theme.mode],
            icon: 'carbon:asleep',
            onClick: () => setSheet('theme'),
        },
        {
            label: t('settings.preference.language'),
            status: localeLabels[app.locale],
            icon: 'carbon:language',
            onClick: () => setSheet('language'),
        },
    ];

    const renderRows = (items: SettingRow[]) => (
        <View className='au-card au-card--glow overflow-hidden au-anim'>
            {items.map((item, index) => (
                <View
                  key={index}
                  className={`flex items-center justify-between px-[28rpx] h-[104rpx] ${index < items.length - 1 ? 'border-b-[length:1rpx] border-splitLine' : ''}`}
                  hoverClass={item.onClick ? 'tap-active' : 'none'}
                  onClick={item.onClick}
                >
                    <View className='flex items-center gap-[20rpx]'>
                        <View className='w-[60rpx] h-[60rpx] rounded-[16rpx] bg-primary-tint flex items-center justify-center'>
                            <XIcon name={item.icon} size={26} className='text-primary' />
                        </View>
                        <Text className='text-body text-textMain'>{item.label}</Text>
                    </View>
                    <View className='flex items-center gap-[10rpx]'>
                        {item.status && <Text className='text-desc text-textThird'>{item.status}</Text>}
                        {item.onClick && <XIcon name='carbon:chevron-right' size={16} className='text-textWeak' />}
                    </View>
                </View>
            ))}
        </View>
    );

    return (
        <View className='au-page w-full bg-pageBg min-h-screen pb-[60rpx]'>
            <View className='px-[24rpx] py-[32rpx]'>
                {/* 账号设置 */}
                <Text className='au-section-title block mb-[24rpx]'>{t('settings.account.title')}</Text>
                {renderRows(accountItems)}

                {/* 偏好设置 */}
                <Text className='au-section-title block mt-[40rpx] mb-[24rpx]'>{t('settings.preference.title')}</Text>
                {renderRows(preferenceItems)}

                {/* 关于 / 法务 */}
                <Text className='au-section-title block mt-[40rpx] mb-[24rpx]'>{t('page.about.title')}</Text>
                <View className='au-card overflow-hidden au-anim au-anim-2'>
                    {[
                        {icon: 'carbon:information', label: t('page.title.about'), path: '/about'},
                        {icon: 'carbon:document', label: t('page.title.terms'), path: '/terms'},
                        {icon: 'carbon:shield-alert', label: t('page.title.privacy'), path: '/privacy'},
                        {icon: 'carbon:warning-alt', label: t('page.title.disclaimer'), path: '/disclaimer'},
                    ].map((item, index, arr) => (
                        <View
                          key={item.path}
                          className={`flex items-center justify-between px-[28rpx] h-[104rpx] ${index < arr.length - 1 ? 'border-b-[length:1rpx] border-splitLine' : ''}`}
                          hoverClass='tap-active'
                          onClick={() => router.push(item.path)}
                        >
                            <View className='flex items-center gap-[20rpx]'>
                                <View className='w-[60rpx] h-[60rpx] rounded-[16rpx] bg-primary-tint flex items-center justify-center'>
                                    <XIcon name={item.icon} size={26} className='text-primary' />
                                </View>
                                <Text className='text-body text-textMain'>{item.label}</Text>
                            </View>
                            <XIcon name='carbon:chevron-right' size={16} className='text-textWeak' />
                        </View>
                    ))}
                </View>
            </View>

            <OptionSheet
              visible={sheet !== null}
              title={sheet === 'language' ? t('settings.preference.language') : t('settings.preference.theme')}
              options={sheet === 'language' ? languageOptions : themeOptions}
              value={sheet === 'language' ? app.locale : theme.mode}
              onSelect={handleSelect}
              onClose={() => setSheet(null)}
            />
        </View>
    );
}
