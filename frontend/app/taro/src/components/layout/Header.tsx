import {useState, useEffect} from 'react';
import {View, Image} from '@tarojs/components';
import {useTranslations} from '@/lib/next-intl-compat';
import logoImage from '@/assets/images/logo.png';

import {Button} from '@/components/ui/button';
import {
    DropdownMenu,
    DropdownMenuContent,
    DropdownMenuItem,
    DropdownMenuSeparator,
    DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu';

import {XIcon} from '@/plugins/xicon';
import {useI18n} from '@/i18n';
import {useI18nRouter} from '@/i18n/helpers/useI18nRouter';
import {usePreferences} from '@/core/preferences';
import {useAccessStore} from '@/store/core/access/store';
import {useAuth} from '@/api/hooks/auth';

import TopNavbar from './TopNavbar';
import MobileNav from './MobileNav';

export default function Header() {
    const t = useTranslations('navbar');
    const appT = useTranslations('app');
    const menuT = useTranslations('menu');
    const brandTitle = appT('title');

    const {theme: themePref, setThemeMode} = usePreferences();
    const currentMode = themePref.mode;
    const {changeLocale} = useI18n();
    const router = useI18nRouter();
    const accessStore = useAccessStore();
    const auth = useAuth();

    const accessToken = accessStore.accessToken;
    const isLogin = !!accessToken && !accessStore.loginExpired;

    const handleClickLogo = () => router.push('/');
    const handleClickSettings = () => router.push('/settings');
    const handleClickUserHomepage = () => router.push('/user');
    const handleClickLogin = () => router.push('/login');
    const handleClickRegister = () => router.push('/register');
    const handleClickLogout = async () => {
        if (isLogin) await auth.logout();
    };

    const [mounted, setMounted] = useState(false);
    useEffect(() => {
        setMounted(true);
    }, []);

    const themeIconName = (() => {
        if (!mounted) return 'carbon:monitor';
        if (currentMode === 'dark') return 'carbon:moon';
        if (currentMode === 'light') return 'carbon:sun';
        return 'carbon:monitor';
    })();

    return (
        <View className='fixed top-0 left-0 right-0 z-1000 flex justify-center bg-background/80 backdrop-blur-md border-b border-border/50 dark:border-border/30 dark:bg-background/60'>
            <View className='flex h-(--layout-header-height) w-full max-w-300 items-center gap-6 px-4 max-md:gap-3 max-md:px-3'>
                {/* Logo */}
                <View
                  onClick={handleClickLogo}
                  className='flex shrink-0 items-center gap-2 rounded-lg p-1 transition-colors hover:bg-primary/10'
                >
                    <View className='h-9 w-9 shrink-0 max-md:h-8 max-md:w-8'>
                        <Image src={logoImage} mode='aspectFit' className='w-full h-full' />
                    </View>
                    <View className='text-lg font-bold text-primary whitespace-nowrap max-md:hidden'>
                        {brandTitle}
                    </View>
                </View>

                {/* 桌面端导航区 */}
                <View className='min-w-0 flex-1 max-md:hidden'>
                    <TopNavbar />
                </View>

                {/* 桌面端功能按钮区 */}
                <View className='flex shrink-0 items-center gap-1 max-md:hidden'>
                    {/* 用户菜单 */}
                    <DropdownMenu modal={false}>
                        <DropdownMenuTrigger>
                            <Button variant='ghost' size='icon'>
                                <XIcon name='carbon:user' size={16} />
                            </Button>
                        </DropdownMenuTrigger>
                        <DropdownMenuContent align='end'>
                            {isLogin ? (
                                <>
                                    <DropdownMenuItem onClick={handleClickUserHomepage}>
                                        <XIcon name='carbon:home' size={16} />
                                        {menuT('homepage')}
                                    </DropdownMenuItem>
                                    <DropdownMenuItem onClick={handleClickSettings}>
                                        <XIcon name='carbon:user' size={16} />
                                        {menuT('my_profile')}
                                    </DropdownMenuItem>
                                    <DropdownMenuSeparator />
                                    <DropdownMenuItem
                                      onClick={handleClickLogout}
                                      className='text-destructive'
                                    >
                                        <XIcon name='carbon:logout' size={16} />
                                        {menuT('logout')}
                                    </DropdownMenuItem>
                                </>
                            ) : (
                                <>
                                    <DropdownMenuItem onClick={handleClickLogin}>
                                        <XIcon name='carbon:user' size={16} />
                                        {t('user.login')}
                                    </DropdownMenuItem>
                                    <DropdownMenuItem onClick={handleClickRegister}>
                                        <XIcon name='carbon:user' size={16} />
                                        {t('user.register')}
                                    </DropdownMenuItem>
                                </>
                            )}
                        </DropdownMenuContent>
                    </DropdownMenu>

                    {/* 语言菜单 */}
                    <DropdownMenu modal={false}>
                        <DropdownMenuTrigger>
                            <Button variant='ghost' size='icon'>
                                <XIcon name='carbon:earth' size={16} />
                            </Button>
                        </DropdownMenuTrigger>
                        <DropdownMenuContent align='end'>
                            <DropdownMenuItem onClick={() => changeLocale('zh-CN')}>
                                简体中文
                            </DropdownMenuItem>
                            <DropdownMenuItem onClick={() => changeLocale('en-US')}>
                                English
                            </DropdownMenuItem>
                        </DropdownMenuContent>
                    </DropdownMenu>

                    {/* 主题菜单 */}
                    <DropdownMenu modal={false}>
                        <DropdownMenuTrigger>
                            <Button variant='ghost' size='icon'>
                                <XIcon name={themeIconName} size={16} />
                            </Button>
                        </DropdownMenuTrigger>
                        <DropdownMenuContent align='end'>
                            <DropdownMenuItem onClick={() => setThemeMode('dark')}>
                                <XIcon name='carbon:moon' size={16} />
                                {t('theme.dark')}
                            </DropdownMenuItem>
                            <DropdownMenuItem onClick={() => setThemeMode('light')}>
                                <XIcon name='carbon:sun' size={16} />
                                {t('theme.light')}
                            </DropdownMenuItem>
                            <DropdownMenuItem onClick={() => setThemeMode('auto')}>
                                <XIcon name='carbon:monitor' size={16} />
                                {t('theme.system')}
                            </DropdownMenuItem>
                        </DropdownMenuContent>
                    </DropdownMenu>
                </View>

                {/* 手机端：主题切换 + 汉堡选单 */}
                <View className='flex shrink-0 items-center gap-1 md:hidden'>
                    <DropdownMenu modal={false}>
                        <DropdownMenuTrigger>
                            <Button variant='ghost' size='icon'>
                                <XIcon name={themeIconName} size={16} />
                            </Button>
                        </DropdownMenuTrigger>
                        <DropdownMenuContent align='end'>
                            <DropdownMenuItem onClick={() => setThemeMode('dark')}>
                                <XIcon name='carbon:moon' size={16} />
                                {t('theme.dark')}
                            </DropdownMenuItem>
                            <DropdownMenuItem onClick={() => setThemeMode('light')}>
                                <XIcon name='carbon:sun' size={16} />
                                {t('theme.light')}
                            </DropdownMenuItem>
                            <DropdownMenuItem onClick={() => setThemeMode('auto')}>
                                <XIcon name='carbon:monitor' size={16} />
                                {t('theme.system')}
                            </DropdownMenuItem>
                        </DropdownMenuContent>
                    </DropdownMenu>

                    {/* 汉堡选单 */}
                    <MobileNav />
                </View>
            </View>
        </View>
    );
}
