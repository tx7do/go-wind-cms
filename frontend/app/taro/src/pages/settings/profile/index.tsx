import {useState} from 'react';
import Taro, {useDidShow} from '@tarojs/taro';
import {useTranslation} from 'react-i18next';
import {View, Text, Input, Textarea, Image} from '@tarojs/components';
import XIcon from '@/plugins/xicon';
import {usePageTitle} from '@/hooks/usePageTitle';
import {fetchUserProfile, updateMe} from '@/api/hooks/user-profile';
import {useI18nRouter} from '@/i18n/helpers';
import type {identityservicev1_User} from '@/api/generated/app/service/v1';

/**
 * 编辑资料（独立页）：渐变横幅 + 大头像 + 昵称/简介表单 + 底部悬浮保存条。
 * 保存成功后返回上一页；未变更时保存条为中性灰不可点。
 */
export default function ProfileEditPage() {
    const {t} = useTranslation();
    usePageTitle('page.user.edit_profile');
    const router = useI18nRouter();

    const [profile, setProfile] = useState<identityservicev1_User | null>(null);
    const [nickname, setNickname] = useState('');
    const [description, setDescription] = useState('');
    const [saving, setSaving] = useState(false);

    useDidShow(() => {
        fetchUserProfile()
            .then((result) => {
                const u = (result as identityservicev1_User) || null;
                setProfile(u);
                setNickname(u?.nickname || '');
                setDescription(u?.description || '');
            })
            .catch((error) => console.error('Load profile failed:', error));
    });

    const dirty = profile
        ? (nickname.trim() !== (profile.nickname || '')) || (description.trim() !== (profile.description || ''))
        : false;

    const save = async () => {
        if (!dirty || saving) return;
        setSaving(true);
        try {
            await updateMe({
                nickname: nickname.trim(),
                description: description.trim(),
            });
            Taro.showToast({title: t('common.success'), icon: 'success'});
            setTimeout(() => router.back(), 600);
        } catch (error) {
            console.error('Save profile failed:', error);
            Taro.showToast({title: t('common.error'), icon: 'none'});
        } finally {
            setSaving(false);
        }
    };

    return (
        <View className='au-page w-full bg-pageBg min-h-screen pb-[220rpx]'>
            {/* 渐变横幅：居中大头像 + 用户名 */}
            <View
              className='relative overflow-hidden pt-[56rpx] pb-[124rpx]'
              style={{background: 'linear-gradient(128deg, #57a3ff 0%, #1677ff 52%, #0b55d4 100%)'}}
            >
                <View
                  className='absolute rounded-full'
                  style={{
                      width: '380rpx',
                      height: '380rpx',
                      right: '-120rpx',
                      top: '-160rpx',
                      background: 'radial-gradient(circle, rgba(255,255,255,0.26), transparent 68%)',
                  }}
                />
                <View className='relative flex flex-col items-center'>
                    <View className='w-[160rpx] h-[160rpx] rounded-full bg-white/18 border-[length:3rpx] border-white/40 flex items-center justify-center overflow-hidden shadow-[0_16rpx_40rpx_rgba(3,20,60,0.35)]'>
                        {profile?.avatar ? (
                            <Image src={profile.avatar} className='w-full h-full' mode='aspectFill' />
                        ) : (
                            <Text className='text-[length:64rpx] font-bold text-white'>
                                {profile?.nickname?.charAt(0) || profile?.username?.charAt(0) || 'U'}
                            </Text>
                        )}
                    </View>
                    <Text className='text-desc text-white/85 mt-[16rpx]'>
                        @{profile?.username || ''}
                    </Text>
                </View>
            </View>

            {/* 表单卡 */}
            <View className='px-[24rpx] -mt-[76rpx]'>
                <View className='au-card au-card--glow p-[32rpx] au-anim'>
                    <View className='mb-[28rpx]'>
                        <Text className='text-tips text-textThird block mb-[12rpx]'>
                            {t('page.user.nickname')}
                        </Text>
                        <View className='au-input'>
                            <XIcon name='carbon:user' size={26} className='text-textThird shrink-0' />
                            <Input
                              className='au-input-field'
                              value={nickname}
                              placeholder={t('page.user.nickname')}
                              placeholderClass='text-textWeak'
                              maxlength={30}
                              onInput={(e) => setNickname(e.detail.value)}
                            />
                        </View>
                    </View>

                    <View>
                        <View className='flex items-center justify-between mb-[12rpx]'>
                            <Text className='text-tips text-textThird'>
                                {t('page.user.description')}
                            </Text>
                            <Text className='text-tips text-textWeak'>
                                {description.length}/120
                            </Text>
                        </View>
                        <View className='au-textarea-box rounded-[16rpx] bg-pageBg border-[length:1rpx] border-splitLine px-[20rpx] py-[16rpx] box-border'>
                            <Textarea
                              className='au-textarea'
                              autoHeight
                              value={description}
                              placeholder={t('page.user.description')}
                              placeholderClass='text-textWeak'
                              maxlength={120}
                              onInput={(e) => setDescription(e.detail.value)}
                            />
                        </View>
                    </View>
                </View>
            </View>

            {/* 底部悬浮保存条 */}
            <View className='au-bottom-bar au-glass'>
                <View
                  className={`${dirty ? 'btn-hero' : 'au-btn-disabled'} w-full h-[92rpx] rounded-full flex items-center justify-center gap-[12rpx]`}
                  hoverClass={dirty ? 'tap-active' : 'none'}
                  onClick={() => dirty && save()}
                >
                    {saving && <View className='au-btn-spin' />}
                    <Text className={`text-body font-bold ${dirty || saving ? 'text-white' : 'text-textThird'}`}>
                        {saving ? t('common.loading') : t('common.save')}
                    </Text>
                </View>
            </View>
        </View>
    );
}
