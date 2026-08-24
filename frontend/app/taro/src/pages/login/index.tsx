import {useState} from 'react';
import Taro, {useRouter} from '@tarojs/taro';
import {useTranslation} from 'react-i18next';
import {View, Text, Input} from '@tarojs/components';
import XIcon from '@/plugins/xicon';
import {Button} from '@/components/ui/button';
import {useAuth} from '@/api/hooks/auth';
import {useI18nRouter} from '@/i18n/helpers';
import {usePageTitle} from '@/hooks/usePageTitle';

/**
 * 账号密码登录页
 *
 * 登录编排（AES 加密、token 持久化、profile 拉取）由 useAuth().login 内部处理，
 * 成功后 reLaunch 回首页（或 redirect 参数指定页），避免登录页残留在页面栈。
 */
export default function LoginPage() {
    const {t} = useTranslation();
    usePageTitle('authentication.login.title');
    const router = useI18nRouter();
    const taroRouter = useRouter();
    const {login} = useAuth();

    const [username, setUsername] = useState('');
    const [password, setPassword] = useState('');
    const [submitting, setSubmitting] = useState(false);

    async function handleLogin() {
        if (!username || !password) {
            Taro.showToast({title: t('authentication.login.login_failed'), icon: 'none'});
            return;
        }

        setSubmitting(true);
        try {
            await login({username, password}, async () => {
                Taro.showToast({title: t('authentication.login.login_success'), icon: 'none'});
                // redirect 参数优先（logout 流程带入），否则回首页
                const redirect = taroRouter.params.redirect;
                if (redirect) {
                    await Taro.reLaunch({url: decodeURIComponent(redirect)});
                } else {
                    await Taro.reLaunch({url: '/pages/index/index'});
                }
            });
        } catch (error) {
            console.error('Login failed:', error);
            Taro.showToast({title: t('authentication.login.login_failed'), icon: 'none'});
        } finally {
            setSubmitting(false);
        }
    }

    return (
        <View className='min-h-screen w-full bg-pageBg flex flex-col justify-center px-[48rpx]'>
            {/* 品牌区 */}
            <View className='flex flex-col items-center mb-[64rpx]'>
                <View className='w-[120rpx] h-[120rpx] rounded-[32rpx] bg-primary/10 flex items-center justify-center mb-[24rpx]'>
                    <XIcon name='carbon:wind-gusts' size={48} className='text-primary'/>
                </View>
                <Text className='text-title font-bold text-textMain block mb-[8rpx]'>
                    {t('authentication.login.title')}
                </Text>
                <Text className='text-desc text-textSec'>
                    {t('authentication.login.brand_subtitle')}
                </Text>
            </View>

            {/* 表单卡片 */}
            <View className='rounded-[32rpx] bg-cardBg p-[40rpx]'>
                <View className='mb-[32rpx]'>
                    <Text className='text-desc text-textSec block mb-[12rpx]'>
                        {t('authentication.login.username')}
                    </Text>
                    <Input
                        className='h-[88rpx] px-[24rpx] text-[28rpx] text-textMain rounded-[16rpx] bg-pageBg'
                        value={username}
                        placeholder={t('authentication.login.username')}
                        onInput={(e) => setUsername(e.detail.value)}
                    />
                </View>

                <View className='mb-[48rpx]'>
                    <Text className='text-desc text-textSec block mb-[12rpx]'>
                        {t('authentication.login.password')}
                    </Text>
                    <Input
                        className='h-[88rpx] px-[24rpx] text-[28rpx] text-textMain rounded-[16rpx] bg-pageBg'
                        value={password}
                        password
                        placeholder={t('authentication.login.placeholder_password')}
                        onInput={(e) => setPassword(e.detail.value)}
                    />
                </View>

                <Button
                    variant='default'
                    size='lg'
                    className='w-full'
                    disabled={submitting}
                    onClick={handleLogin}
                >
                    {t('authentication.login.login')}
                </Button>

                <View
                    className='flex justify-center mt-[32rpx]'
                    hoverClass='tap-active'
                    onClick={() => router.reLaunch('/')}
                >
                    <Text className='text-tips text-textSec'>
                        {t('authentication.login.back_home')}
                    </Text>
                </View>
            </View>
        </View>
    );
}
