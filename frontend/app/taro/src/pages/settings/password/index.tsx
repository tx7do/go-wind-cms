import {useMemo, useState} from 'react';
import Taro from '@tarojs/taro';
import {useTranslation} from 'react-i18next';
import {View, Text, Input} from '@tarojs/components';
import XIcon from '@/plugins/xicon';
import {usePageTitle} from '@/hooks/usePageTitle';
import {changePassword} from '@/api/hooks/user-profile';
import {useI18nRouter} from '@/i18n/helpers';

/**
 * 修改密码：旧密码 + 新密码 + 确认。
 * 新旧密码均 AES 加密传输（与登录同款）。强度提示纯前端提示，最终以服务端校验为准。
 */
export default function ChangePasswordPage() {
    const {t} = useTranslation();
    usePageTitle('page.title.settings');
    const router = useI18nRouter();

    const [oldPassword, setOldPassword] = useState('');
    const [newPassword, setNewPassword] = useState('');
    const [confirmPassword, setConfirmPassword] = useState('');
    const [submitting, setSubmitting] = useState(false);

    const strength = useMemo(() => {
        let score = 0;
        if (newPassword.length >= 8) score += 1;
        if (/[a-z]/.test(newPassword) && /[A-Z]/.test(newPassword)) score += 1;
        if (/\d/.test(newPassword)) score += 1;
        if (/[^A-Za-z0-9]/.test(newPassword)) score += 1;
        return score; // 0..4
    }, [newPassword]);

    const valid = oldPassword.length > 0 && newPassword.length >= 6 && newPassword === confirmPassword;

    const submit = async () => {
        if (!valid || submitting) return;
        setSubmitting(true);
        try {
            await changePassword(oldPassword, newPassword);
            Taro.showToast({title: t('common.success'), icon: 'success'});
            setTimeout(() => router.back(), 800);
        } catch (error) {
            console.error('Change password failed:', error);
            Taro.showToast({title: t('common.error'), icon: 'none'});
        } finally {
            setSubmitting(false);
        }
    };

    const strengthLabels = [t('common.warning'), t('common.warning'), t('common.info'), t('common.success'), t('common.success')];
    const strengthColors = ['#ff4d4f', '#ff9f4b', '#ff9f4b', '#22c3d8', '#35c26a'];

    const fields = [
        {
            label: t('settings.account.old_password'),
            value: oldPassword,
            setValue: setOldPassword,
            placeholder: t('settings.account.old_password'),
            key: 'old' as const,
        },
        {
            label: t('settings.account.new_password'),
            value: newPassword,
            setValue: setNewPassword,
            placeholder: t('settings.account.new_password'),
            key: 'new' as const,
        },
        {
            label: t('authentication.register.confirm_password'),
            value: confirmPassword,
            setValue: setConfirmPassword,
            placeholder: t('authentication.register.confirm_password'),
            key: 'confirm' as const,
        },
    ];

    return (
        <View className='au-page w-full bg-pageBg min-h-screen pb-[60rpx]'>
            <View className='px-[24rpx] pt-[36rpx]'>
                <View className='au-anim au-card au-card--glow p-[32rpx]'>
                    {fields.map((field) => (
                        <View key={field.key} className='mb-[24rpx]'>
                            <Text className='login-field-label'>{field.label}</Text>
                            <View className='au-input'>
                                <XIcon name='carbon:locked' size={26} className='text-textThird shrink-0' />
                                <Input
                                  className='au-input-field'
                                  password
                                  value={field.value}
                                  placeholder={field.placeholder}
                                  placeholderClass='text-textWeak'
                                  onInput={(e) => field.setValue(e.detail.value)}
                                />
                            </View>
                        </View>
                    ))}

                    {/* 强度指示（新密码非空时展示） */}
                    {newPassword.length > 0 && (
                        <View className='mb-[24rpx]'>
                            <View className='flex gap-[10rpx] mb-[10rpx]'>
                                {[0, 1, 2, 3].map((i) => (
                                    <View
                                      key={i}
                                      className='flex-1 h-[8rpx] rounded-full'
                                      style={{background: i < strength ? strengthColors[strength] : 'var(--color-split-line)'}}
                                    />
                                ))}
                            </View>
                            <Text className='text-tips' style={{color: strengthColors[strength]}}>
                                {strengthLabels[strength]}
                            </Text>
                        </View>
                    )}

                    {confirmPassword.length > 0 && newPassword !== confirmPassword && (
                        <Text className='text-tips text-danger block mb-[20rpx]'>
                            {t('authentication.register.password_mismatch')}
                        </Text>
                    )}

                    <View
                      className={`${valid ? 'btn-hero' : 'au-btn-disabled bg-splitLine'} h-[92rpx] rounded-full flex items-center justify-center mt-[8rpx]`}
                      hoverClass={valid ? 'tap-active' : 'none'}
                      onClick={() => valid && submit()}
                    >
                        <Text className={`text-body font-bold ${valid ? 'text-white' : 'text-textThird'}`}>
                            {submitting ? t('common.loading') : t('common.save')}
                        </Text>
                    </View>
                </View>
            </View>
        </View>
    );
}
