'use client';

import {useState, useMemo} from 'react';
import {useTranslations, useLocale} from 'next-intl';
import {Switch} from '@/components/ui/switch';
import {Select, SelectContent, SelectItem, SelectTrigger, SelectValue} from '@/components/ui/select';
import {Button} from '@/components/ui/button';
import XIcon from '@/plugins/xicon';

import {usePreferences} from '@/core/preferences';
import {useI18n} from '@/i18n';
import type {ThemeModeType, SupportedLanguagesType} from '@/core/preferences';

interface MenuItem {
    key: string;
    icon: string;
    label: string;
}

export default function SettingsPage() {
    const t = useTranslations('settings');
    const locale = useLocale();

    const {
        theme: themePref,
        content: contentPref,
        setThemeMode,
        updateContent,
        setLanguage,
        resetPreferences,
    } = usePreferences();
    const {changeLocale} = useI18n();

    const [activeMenu, setActiveMenu] = useState<'account' | 'message' | 'preference'>('account');

    const menuItems: MenuItem[] = useMemo(() => [
        {key: 'account', icon: 'carbon:user', label: t('menu.account')},
        {key: 'message', icon: 'carbon:email', label: t('menu.message')},
        {key: 'preference', icon: 'carbon:settings', label: t('menu.preference')}
    ], [t]);

    const themeOptions = useMemo(() => [
        {value: 'light' as ThemeModeType, label: t('preference.theme_light')},
        {value: 'dark' as ThemeModeType, label: t('preference.theme_dark')},
        {value: 'auto' as ThemeModeType, label: t('preference.theme_auto')},
    ], [t]);

    const languageOptions = useMemo(() => [
        {value: 'zh-CN' as SupportedLanguagesType, label: t('preference.language_chinese')},
        {value: 'en-US' as SupportedLanguagesType, label: t('preference.language_english')},
    ], [t]);

    const handleThemeChange = (mode: ThemeModeType) => {
        setThemeMode(mode);
    };

    const handleLanguageChange = (newLocale: SupportedLanguagesType) => {
        setLanguage(newLocale);
        if (newLocale !== locale) {
            changeLocale(newLocale);
        }
    };

    return (
        <div className="w-full py-8 max-md:py-4">
            <div className="w-full max-w-[1200px] mx-auto grid grid-cols-[220px_1fr_280px] gap-6 px-8 max-md:grid-cols-1 max-md:px-4">
                {/* 左侧导航 */}
                <aside className="max-md:hidden">
                    <nav className="space-y-1">
                        {menuItems.map((item) => (
                            <div
                                key={item.key}
                                className={`flex cursor-pointer items-center gap-3 rounded-lg border border-transparent px-3 py-2.5 text-sm font-medium transition-all ${
                                    activeMenu === item.key
                                        ? 'border-primary/20 bg-primary/5 text-primary'
                                        : 'text-muted-foreground hover:bg-muted hover:text-foreground'
                                }`}
                                onClick={() => setActiveMenu(item.key as 'account' | 'message' | 'preference')}
                            >
                                <div className="flex h-8 w-8 items-center justify-center rounded-md bg-primary/10 text-primary">
                                    <XIcon name={item.icon} size={18}/>
                                </div>
                                <span>{item.label}</span>
                            </div>
                        ))}
                    </nav>
                </aside>

                {/* 中间内容区 */}
                <main className="min-w-0">
                    {/* 账号设置 */}
                    {activeMenu === 'account' && (
                        <>
                            <div className="mb-8">
                                <h1 className="mb-2 text-2xl font-bold text-foreground">{t('account.title')}</h1>
                                <p className="text-sm text-muted-foreground">{t('account.subtitle')}</p>
                            </div>
                            <div className="mb-8">
                                <h2 className="mb-2 text-lg font-semibold text-foreground">{t('account.section_title')}</h2>
                                <p className="mb-4 text-sm text-muted-foreground">{t('account.section_desc')}</p>
                                <div className="space-y-3">
                                    <div className="flex items-center justify-between rounded-lg border border-border bg-card p-4">
                                        <div>
                                            <div className="text-sm font-medium text-foreground">{t('account.password')}</div>
                                            <div className="mt-1 text-xs text-muted-foreground">{t('account.password_not_set')}</div>
                                        </div>
                                        <Button variant="outline" size="sm">{t('account.edit')}</Button>
                                    </div>
                                    <div className="flex items-center justify-between rounded-lg border border-border bg-card p-4">
                                        <div>
                                            <div className="text-sm font-medium text-foreground">{t('account.bind_phone')}</div>
                                            <div className="mt-1 text-xs text-muted-foreground">{t('account.password_not_set')}</div>
                                        </div>
                                        <Button variant="outline" size="sm">{t('account.edit')}</Button>
                                    </div>
                                    <div className="flex items-center justify-between rounded-lg border border-border bg-card p-4">
                                        <div>
                                            <div className="text-sm font-medium text-foreground">{t('account.bind_email')}</div>
                                            <div className="mt-1 text-xs text-muted-foreground">{t('account.email_not_bound')}</div>
                                        </div>
                                        <Button variant="outline" size="sm">{t('account.edit')}</Button>
                                    </div>
                                </div>
                            </div>
                            <div>
                                <h2 className="mb-4 text-lg font-semibold text-foreground">{t('account.third_party_title')}</h2>
                                <div className="grid grid-cols-3 gap-3 max-md:grid-cols-1">
                                    <div className="flex items-center gap-3 rounded-lg border border-border bg-card p-4">
                                        <div className="flex h-10 w-10 items-center justify-center rounded-lg bg-green-500">
                                            <XIcon name="fa:wechat" size={24} color="white"/>
                                        </div>
                                        <span className="text-sm font-medium text-foreground">
                                            {t('account.bind_wechat')}
                                        </span>
                                    </div>
                                    <div className="flex items-center gap-3 rounded-lg border border-border bg-card p-4">
                                        <div className="flex h-10 w-10 items-center justify-center rounded-lg bg-red-500">
                                            <XIcon name="fa:weibo" size={24} color="white"/>
                                        </div>
                                        <span className="text-sm font-medium text-foreground">
                                            {t('account.bind_weibo')}
                                        </span>
                                    </div>
                                    <div className="flex items-center gap-3 rounded-lg border border-border bg-card p-4">
                                        <div className="flex h-10 w-10 items-center justify-center rounded-lg bg-blue-500">
                                            <XIcon name="fa:qq" size={24} color="white"/>
                                        </div>
                                        <span className="text-sm font-medium text-foreground">
                                            {t('account.bind_qq')}
                                        </span>
                                    </div>
                                </div>
                            </div>
                        </>
                    )}

                    {/* 消息设置 */}
                    {activeMenu === 'message' && (
                        <>
                            <div className="mb-8">
                                <h1 className="mb-2 text-2xl font-bold text-foreground">{t('message.title')}</h1>
                                <p className="text-sm text-muted-foreground">{t('message.subtitle')}</p>
                            </div>
                            <div className="mb-8">
                                <h2 className="mb-2 text-lg font-semibold text-foreground">{t('message.email_notifications')}</h2>
                                <p className="mb-4 text-sm text-muted-foreground">{t('message.email_notifications_desc')}</p>
                                <div className="space-y-3">
                                    <div className="flex items-center justify-between rounded-lg border border-border bg-card p-4">
                                        <div className="text-sm font-medium text-foreground">{t('message.system_messages')}</div>
                                        <Switch defaultChecked/>
                                    </div>
                                    <div className="flex items-center justify-between rounded-lg border border-border bg-card p-4">
                                        <div className="text-sm font-medium text-foreground">{t('message.comment_notifications')}</div>
                                        <Switch defaultChecked/>
                                    </div>
                                    <div className="flex items-center justify-between rounded-lg border border-border bg-card p-4">
                                        <div className="text-sm font-medium text-foreground">{t('message.activity_updates')}</div>
                                        <Switch defaultChecked/>
                                    </div>
                                    <div className="flex items-center justify-between rounded-lg border border-border bg-card p-4">
                                        <div className="text-sm font-medium text-foreground">{t('message.recommended_content')}</div>
                                        <Switch defaultChecked/>
                                    </div>
                                </div>
                            </div>
                            <div>
                                <h2 className="mb-4 text-lg font-semibold text-foreground">{t('message.email_frequency')}</h2>
                                <div className="flex items-center justify-between rounded-lg border border-border bg-card p-4">
                                    <div className="text-sm font-medium text-foreground">{t('message.frequency_desc')}</div>
                                    <Select defaultValue="immediately">
                                        <SelectTrigger className="w-[180px] h-8">
                                            <SelectValue/>
                                        </SelectTrigger>
                                        <SelectContent>
                                            <SelectItem value="immediately">{t('message.frequency_immediately')}</SelectItem>
                                            <SelectItem value="daily">{t('message.frequency_daily')}</SelectItem>
                                            <SelectItem value="weekly">{t('message.frequency_weekly')}</SelectItem>
                                            <SelectItem value="monthly">{t('message.frequency_monthly')}</SelectItem>
                                        </SelectContent>
                                    </Select>
                                </div>
                            </div>
                        </>
                    )}

                    {/* 偏好设置 - 连接 usePreferences */}
                    {activeMenu === 'preference' && (
                        <>
                            <div className="mb-8">
                                <h1 className="mb-2 text-2xl font-bold text-foreground">{t('preference.title')}</h1>
                                <p className="text-sm text-muted-foreground">{t('preference.subtitle')}</p>
                            </div>
                            <div className="mb-8">
                                <h2 className="mb-2 text-lg font-semibold text-foreground">{t('preference.theme_settings')}</h2>
                                <p className="mb-4 text-sm text-muted-foreground">{t('preference.theme_desc')}</p>
                                <div className="flex items-center justify-between rounded-lg border border-border bg-card p-4">
                                    <div className="text-sm font-medium text-foreground">{t('preference.theme')}</div>
                                    <Select value={themePref.mode} onValueChange={(v) => handleThemeChange(v as ThemeModeType)}>
                                        <SelectTrigger className="w-[180px] h-8">
                                            <SelectValue/>
                                        </SelectTrigger>
                                        <SelectContent>
                                            {themeOptions.map(opt => (
                                                <SelectItem key={opt.value} value={opt.value}>{opt.label}</SelectItem>
                                            ))}
                                        </SelectContent>
                                    </Select>
                                </div>
                            </div>
                            <div className="mb-8">
                                <h2 className="mb-2 text-lg font-semibold text-foreground">{t('preference.language_settings')}</h2>
                                <p className="mb-4 text-sm text-muted-foreground">{t('preference.language_desc')}</p>
                                <div className="flex items-center justify-between rounded-lg border border-border bg-card p-4">
                                    <div className="text-sm font-medium text-foreground">{t('preference.language')}</div>
                                    <Select value={locale as SupportedLanguagesType} onValueChange={(v) => handleLanguageChange(v as SupportedLanguagesType)}>
                                        <SelectTrigger className="w-[180px] h-8">
                                            <SelectValue/>
                                        </SelectTrigger>
                                        <SelectContent>
                                            {languageOptions.map(opt => (
                                                <SelectItem key={opt.value} value={opt.value}>{opt.label}</SelectItem>
                                            ))}
                                        </SelectContent>
                                    </Select>
                                </div>
                            </div>
                            <div className="mb-8">
                                <h2 className="mb-2 text-lg font-semibold text-foreground">{t('preference.content_preferences')}</h2>
                                <p className="mb-4 text-sm text-muted-foreground">{t('preference.content_desc')}</p>
                                <div className="space-y-3">
                                    <div className="flex items-center justify-between rounded-lg border border-border bg-card p-4">
                                        <div className="text-sm font-medium text-foreground">{t('preference.hide_sensitive_content')}</div>
                                        <Switch
                                            checked={contentPref.hideSensitiveContent}
                                            onCheckedChange={(checked) => updateContent({hideSensitiveContent: checked})}
                                        />
                                    </div>
                                    <div className="flex items-center justify-between rounded-lg border border-border bg-card p-4">
                                        <div className="text-sm font-medium text-foreground">{t('preference.compact_mode')}</div>
                                        <Switch
                                            checked={contentPref.compactMode}
                                            onCheckedChange={(checked) => updateContent({compactMode: checked})}
                                        />
                                    </div>
                                    <div className="flex items-center justify-between rounded-lg border border-border bg-card p-4">
                                        <div className="text-sm font-medium text-foreground">{t('preference.show_recommendations')}</div>
                                        <Switch
                                            checked={contentPref.showRecommendations}
                                            onCheckedChange={(checked) => updateContent({showRecommendations: checked})}
                                        />
                                    </div>
                                </div>
                            </div>
                            <div>
                                <div className="flex items-center justify-between rounded-lg border border-border bg-card p-4">
                                    <div className="text-sm font-medium text-foreground">恢复默认设置</div>
                                    <Button variant="destructive" size="sm" onClick={() => resetPreferences()}>
                                        重置
                                    </Button>
                                </div>
                            </div>
                        </>
                    )}
                </main>

                {/* 右侧帮助区 */}
                <aside className="max-md:hidden">
                    <div className="sticky top-24 rounded-lg border border-border bg-card p-5">
                        <h3 className="mb-3 text-base font-bold text-foreground">{t('help.title')}</h3>
                        <h4 className="mb-2 text-sm font-semibold text-foreground">{t('help.account_password')}</h4>
                        <ul className="mb-4 space-y-1.5 text-xs text-muted-foreground">
                            <li>1. {t('help.q1')}</li>
                            <li>2. {t('help.q2')}</li>
                            <li>3. {t('help.q3')}</li>
                            <li>4. {t('help.q4')}</li>
                            <li>5. {t('help.q5')}</li>
                        </ul>
                        <h4 className="mb-2 text-sm font-semibold text-foreground">{t('help.other_issues')}</h4>
                        <ul className="space-y-1.5 text-xs text-muted-foreground">
                            <li>6. <a href="#" className="text-primary hover:underline">{t('help.q6')} {t('help.q6_link')}</a></li>
                            <li>7. <a href="#" className="text-primary hover:underline">{t('help.q7')} {t('help.q7_link')}</a></li>
                        </ul>
                    </div>
                </aside>
            </div>
        </div>
    );
}
