import {useEffect, useState} from 'react';
import Taro, {useRouter} from '@tarojs/taro';
import {useTranslation} from 'react-i18next';
import {View, Text, Input, Image} from '@tarojs/components';
import XIcon from '@/plugins/xicon';
import {useAuth} from '@/api/hooks/auth';
import {useI18nRouter} from '@/i18n/helpers';
import {usePageTitle} from '@/hooks/usePageTitle';
import logoImage from '@/assets/images/logo.png';
import './login.css';

const REMEMBER_USERNAME_KEY = 'login-remember-username';

/**
 * 账号密码登录页（深空蓝极光品牌风格，全屏沉浸布局）
 *
 * 极光 + 玻璃语言浅/暗两态共用，色档由 login.css 的 --lg-* 变量按主题切换，
 * 跟随应用主题设置（ThemeClientProvider 在 H5 给 <html>、小程序给主题根 View 挂 .dark）。
 * 登录编排（AES 加密、token 持久化、profile 拉取）由 useAuth().login 内部处理，
 * 成功后 reLaunch 回首页（或 redirect 参数指定页），避免登录页残留在页面栈。
 * 后端 app 面未开放注册 RPC，页面仅提供登录（注册引导移除）。
 */
export default function LoginPage() {
    const {t} = useTranslation();
    usePageTitle('authentication.login.title');
    const router = useI18nRouter();
    const taroRouter = useRouter();
    const {login} = useAuth();

    const [username, setUsername] = useState('');
    const [password, setPassword] = useState('');
    const [showPassword, setShowPassword] = useState(false);
    const [remember, setRemember] = useState(false);
    const [submitting, setSubmitting] = useState(false);

    // 登录页为全屏布局：挂载期间隐藏应用顶栏与页面预留的顶栏高度（仅 H5 有 body）
    useEffect(() => {
        if (typeof document === 'undefined') return;
        document.body.classList.add('fullscreen-auth');
        return () => document.body.classList.remove('fullscreen-auth');
    }, []);

    // 恢复“记住我”保存的用户名
    useEffect(() => {
        try {
            const saved = Taro.getStorageSync(REMEMBER_USERNAME_KEY);
            if (saved) {
                setUsername(saved);
                setRemember(true);
            }
        } catch {
            // 本地存储不可用时不影响登录主流程
        }
    }, []);

    async function handleLogin() {
        if (!username || !password) {
            Taro.showToast({title: t('authentication.login.login_failed'), icon: 'none'});
            return;
        }

        setSubmitting(true);
        try {
            await login({username, password}, async () => {
                try {
                    if (remember) {
                        Taro.setStorageSync(REMEMBER_USERNAME_KEY, username);
                    } else {
                        Taro.removeStorageSync(REMEMBER_USERNAME_KEY);
                    }
                } catch {
                    // 本地存储不可用不影响登录成功
                }
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

    function handleForgotPassword() {
        Taro.showToast({title: t('authentication.login.forgot_password_tip'), icon: 'none'});
    }

    const features = [
        {icon: 'carbon:application', label: t('authentication.login.feature_projects')},
        {icon: 'carbon:encryption', label: t('authentication.login.feature_isolation')},
        {icon: 'carbon:rule', label: t('authentication.login.feature_permissions')},
    ];

    return (
        <View className='login-page'>
            {/* 背景：网格底纹 + 极光光斑 */}
            <View className='login-grid' />
            <View className='login-orb login-orb--1' />
            <View className='login-orb login-orb--2' />
            <View className='login-orb login-orb--3' />

            <View className='login-content'>
                <View className='flex-1' />

                {/* 品牌区 */}
                <View className='flex flex-col items-center'>
                    <View className='login-logo-wrap login-anim login-anim-2'>
                        <View className='login-logo-ring' />
                        <View className='login-logo'>
                            <Image src={logoImage} mode='aspectFill' className='w-full h-full' />
                        </View>
                    </View>

                    <Text className='login-title login-anim login-anim-3'>
                        {t('authentication.login.brand_title')}
                    </Text>
                    <Text className='login-subtitle login-anim login-anim-3'>
                        {t('authentication.login.brand_subtitle')}
                    </Text>

                    <View className='flex items-center justify-center gap-[12rpx] mt-[16rpx] login-anim login-anim-4'>
                        {features.map((f) => (
                            <View key={f.label} className='login-chip login-chip--blue'>
                                <XIcon name={f.icon} size={22} />
                                <Text>{f.label}</Text>
                            </View>
                        ))}
                    </View>
                </View>

                {/* 玻璃表单卡 */}
                <View className='login-card login-anim login-anim-5 mt-[22rpx]'>
                    {/* 用户名 */}
                    <View className='mb-[22rpx]'>
                        <Text className='login-field-label'>
                            {t('authentication.login.username')}
                        </Text>
                        <View className='login-input'>
                            <XIcon name='carbon:user' size={32} className='login-field-icon' />
                            <Input
                              className='login-field'
                              value={username}
                              placeholder={t('authentication.login.username')}
                              placeholderClass='login-placeholder'
                              onInput={(e) => setUsername(e.detail.value)}
                            />
                        </View>
                    </View>

                    {/* 密码 */}
                    <View className='mb-[16rpx]'>
                        <Text className='login-field-label'>
                            {t('authentication.login.password')}
                        </Text>
                        <View className='login-input'>
                            <XIcon name='carbon:locked' size={32} className='login-field-icon' />
                            <Input
                              className='login-field'
                              value={password}
                              password={!showPassword}
                              placeholder={t('authentication.login.placeholder_password')}
                              placeholderClass='login-placeholder'
                              onInput={(e) => setPassword(e.detail.value)}
                            />
                            <View
                              className='login-eye'
                              onClick={() => setShowPassword((v) => !v)}
                            >
                                <XIcon
                                  name={showPassword ? 'carbon:view-off' : 'carbon:view'}
                                  size={30}
                                />
                            </View>
                        </View>
                    </View>

                    {/* 记住我 / 忘记密码 */}
                    <View className='flex items-center justify-between mb-[20rpx] min-h-[56rpx]'>
                        <View
                          className='flex items-center gap-[10rpx]'
                          onClick={() => setRemember((v) => !v)}
                        >
                            <XIcon
                              name={remember ? 'carbon:checkbox-checked' : 'carbon:checkbox'}
                              size={30}
                              className={remember ? 'login-check-on' : 'login-check-off'}
                            />
                            <Text className='login-check-label'>
                                {t('authentication.login.remember_me')}
                            </Text>
                        </View>
                        <View onClick={handleForgotPassword}>
                            <Text className='login-row-link'>
                                {t('authentication.login.forgot_password')}
                            </Text>
                        </View>
                    </View>

                    {/* 提交按钮 */}
                    <View
                      className={`login-submit ${submitting ? 'login-submit--disabled' : ''}`}
                      onClick={submitting ? undefined : handleLogin}
                    >
                        <Text>
                            {submitting
                                ? t('authentication.login.submitting')
                                : t('authentication.login.login')}
                        </Text>
                    </View>
                </View>

                <View className='flex-1' />

                {/* 底部：返回首页 + 条款 */}
                <View className='flex flex-col items-center pt-[14rpx] login-anim login-anim-6'>
                    <View
                      className='px-[24rpx] py-[8rpx]'
                      onClick={() => router.reLaunch('/')}
                    >
                        <Text className='login-back'>
                            {t('authentication.login.back_home')}
                        </Text>
                    </View>
                    <View className='flex items-center flex-wrap justify-center mt-[4rpx] px-[24rpx]'>
                        <Text className='login-footer-text'>
                            {t('authentication.login.terms_prefix')}
                        </Text>
                        <Text
                          className='login-footer-link'
                          onClick={() => router.push('/terms')}
                        >
                            {t('authentication.login.terms_of_service')}
                        </Text>
                        <Text className='login-footer-text'>
                            {t('authentication.login.terms_and')}
                        </Text>
                        <Text
                          className='login-footer-link'
                          onClick={() => router.push('/privacy')}
                        >
                            {t('authentication.login.privacy_policy')}
                        </Text>
                    </View>
                </View>
            </View>
        </View>
    );
}
