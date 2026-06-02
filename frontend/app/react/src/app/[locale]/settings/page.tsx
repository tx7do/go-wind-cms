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

import '../../globals.css';
import styles from './settings.module.css';

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
        <div className={styles.settingsPage}>
            <div className={styles.settingsContainer}>
                {/* 左侧导航 */}
                <aside className={styles.sidebar}>
                    <nav className={styles.menu}>
                        {menuItems.map((item) => (
                            <div
                                key={item.key}
                                className={`${styles.menuItem} ${activeMenu === item.key ? styles.active : ''}`}
                                onClick={() => setActiveMenu(item.key as 'account' | 'message' | 'preference')}
                            >
                                <div className={styles.menuIcon}>
                                    <XIcon name={item.icon} size={18}/>
                                </div>
                                <span className={styles.menuLabel}>{item.label}</span>
                            </div>
                        ))}
                    </nav>
                </aside>

                {/* 中间内容区 */}
                <main className={styles.content}>
                    {/* 账号设置 */}
                    {activeMenu === 'account' && (
                        <>
                            <div className={styles.contentHeader}>
                                <h1 className={styles.title}>{t('account.title')}</h1>
                                <p className={styles.subtitle}>{t('account.subtitle')}</p>
                            </div>
                            <div className={styles.section}>
                                <h2 className={styles.sectionTitle}>{t('account.section_title')}</h2>
                                <p className={styles.sectionDesc}>{t('account.section_desc')}</p>
                                <div className={styles.settingItem}>
                                    <div className={styles.settingLabel}>
                                        <span className={styles.labelText}>{t('account.password')}</span>
                                        <span className={styles.labelStatus}>{t('account.password_not_set')}</span>
                                    </div>
                                    <div className={styles.settingControl}>
                                        <Button variant="outline" size="sm">{t('account.edit')}</Button>
                                    </div>
                                </div>
                                <div className={styles.settingItem}>
                                    <div className={styles.settingLabel}>
                                        <span className={styles.labelText}>{t('account.bind_phone')}</span>
                                        <span className={styles.labelStatus}>{t('account.password_not_set')}</span>
                                    </div>
                                    <div className={styles.settingControl}>
                                        <Button variant="outline" size="sm">{t('account.edit')}</Button>
                                    </div>
                                </div>
                                <div className={styles.settingItem}>
                                    <div className={styles.settingLabel}>
                                        <span className={styles.labelText}>{t('account.bind_email')}</span>
                                        <span className={styles.labelStatus}>{t('account.email_not_bound')}</span>
                                    </div>
                                    <div className={styles.settingControl}>
                                        <Button variant="outline" size="sm">{t('account.edit')}</Button>
                                    </div>
                                </div>
                            </div>
                            <div className={styles.section}>
                                <h2 className={styles.sectionTitle}>{t('account.third_party_title')}</h2>
                                <div className={styles.thirdPartyList}>
                                    <div className={styles.thirdPartyItem}>
                                        <div className={`${styles.platformIcon} ${styles.wechat}`}>
                                            <XIcon name="fa:wechat" size={24} color="white"/>
                                        </div>
                                        <span className={styles.platformLink}>
                                            {t('account.bind_wechat')}
                                        </span>
                                    </div>
                                    <div className={styles.thirdPartyItem}>
                                        <div className={`${styles.platformIcon} ${styles.weibo}`}>
                                            <XIcon name="fa:weibo" size={24} color="white"/>
                                        </div>
                                        <span className={styles.platformLink}>
                                            {t('account.bind_weibo')}
                                        </span>
                                    </div>
                                    <div className={styles.thirdPartyItem}>
                                        <div className={`${styles.platformIcon} ${styles.qq}`}>
                                            <XIcon name="fa:qq" size={24} color="white"/>
                                        </div>
                                        <span className={styles.platformLink}>
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
                            <div className={styles.contentHeader}>
                                <h1 className={styles.title}>{t('message.title')}</h1>
                                <p className={styles.subtitle}>{t('message.subtitle')}</p>
                            </div>
                            <div className={styles.section}>
                                <h2 className={styles.sectionTitle}>{t('message.email_notifications')}</h2>
                                <p className={styles.sectionDesc}>{t('message.email_notifications_desc')}</p>
                                <div className={styles.settingItem}>
                                    <div className={styles.settingLabel}>
                                        <span className={styles.labelText}>{t('message.system_messages')}</span>
                                    </div>
                                    <div className={styles.settingControl}>
                                        <Switch defaultChecked/>
                                    </div>
                                </div>
                                <div className={styles.settingItem}>
                                    <div className={styles.settingLabel}>
                                        <span className={styles.labelText}>{t('message.comment_notifications')}</span>
                                    </div>
                                    <div className={styles.settingControl}>
                                        <Switch defaultChecked/>
                                    </div>
                                </div>
                                <div className={styles.settingItem}>
                                    <div className={styles.settingLabel}>
                                        <span className={styles.labelText}>{t('message.activity_updates')}</span>
                                    </div>
                                    <div className={styles.settingControl}>
                                        <Switch defaultChecked/>
                                    </div>
                                </div>
                                <div className={styles.settingItem}>
                                    <div className={styles.settingLabel}>
                                        <span className={styles.labelText}>{t('message.recommended_content')}</span>
                                    </div>
                                    <div className={styles.settingControl}>
                                        <Switch defaultChecked/>
                                    </div>
                                </div>
                            </div>
                            <div className={styles.section}>
                                <h2 className={styles.sectionTitle}>{t('message.email_frequency')}</h2>
                                <div className={styles.settingItem}>
                                    <div className={styles.settingLabel}>
                                        <span className={styles.labelText}>{t('message.frequency_desc')}</span>
                                    </div>
                                    <div className={styles.settingControl}>
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
                            </div>
                        </>
                    )}

                    {/* 偏好设置 - 连接 usePreferences */}
                    {activeMenu === 'preference' && (
                        <>
                            <div className={styles.contentHeader}>
                                <h1 className={styles.title}>{t('preference.title')}</h1>
                                <p className={styles.subtitle}>{t('preference.subtitle')}</p>
                            </div>
                            <div className={styles.section}>
                                <h2 className={styles.sectionTitle}>{t('preference.theme_settings')}</h2>
                                <p className={styles.sectionDesc}>{t('preference.theme_desc')}</p>
                                <div className={styles.settingItem}>
                                    <div className={styles.settingLabel}>
                                        <span className={styles.labelText}>{t('preference.theme')}</span>
                                    </div>
                                    <div className={styles.settingControl}>
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
                            </div>
                            <div className={styles.section}>
                                <h2 className={styles.sectionTitle}>{t('preference.language_settings')}</h2>
                                <p className={styles.sectionDesc}>{t('preference.language_desc')}</p>
                                <div className={styles.settingItem}>
                                    <div className={styles.settingLabel}>
                                        <span className={styles.labelText}>{t('preference.language')}</span>
                                    </div>
                                    <div className={styles.settingControl}>
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
                            </div>
                            <div className={styles.section}>
                                <h2 className={styles.sectionTitle}>{t('preference.content_preferences')}</h2>
                                <p className={styles.sectionDesc}>{t('preference.content_desc')}</p>
                                <div className={styles.settingItem}>
                                    <div className={styles.settingLabel}>
                                        <span className={styles.labelText}>{t('preference.hide_sensitive_content')}</span>
                                    </div>
                                    <div className={styles.settingControl}>
                                        <Switch
                                            checked={contentPref.hideSensitiveContent}
                                            onCheckedChange={(checked) => updateContent({hideSensitiveContent: checked})}
                                        />
                                    </div>
                                </div>
                                <div className={styles.settingItem}>
                                    <div className={styles.settingLabel}>
                                        <span className={styles.labelText}>{t('preference.compact_mode')}</span>
                                    </div>
                                    <div className={styles.settingControl}>
                                        <Switch
                                            checked={contentPref.compactMode}
                                            onCheckedChange={(checked) => updateContent({compactMode: checked})}
                                        />
                                    </div>
                                </div>
                                <div className={styles.settingItem}>
                                    <div className={styles.settingLabel}>
                                        <span className={styles.labelText}>{t('preference.show_recommendations')}</span>
                                    </div>
                                    <div className={styles.settingControl}>
                                        <Switch
                                            checked={contentPref.showRecommendations}
                                            onCheckedChange={(checked) => updateContent({showRecommendations: checked})}
                                        />
                                    </div>
                                </div>
                            </div>
                            <div className={styles.section}>
                                <div className={styles.settingItem}>
                                    <div className={styles.settingLabel}>
                                        <span className={styles.labelText}>恢复默认设置</span>
                                    </div>
                                    <div className={styles.settingControl}>
                                        <Button variant="destructive" size="sm" onClick={() => resetPreferences()}>
                                            重置
                                        </Button>
                                    </div>
                                </div>
                            </div>
                        </>
                    )}
                </main>

                {/* 右侧帮助区 */}
                <aside className={styles.helpSidebar}>
                    <div className={styles.helpSection}>
                        <h3 className={styles.helpTitle}>{t('help.title')}</h3>
                        <h4 className={styles.helpSubtitle}>{t('help.account_password')}</h4>
                        <ul className={styles.helpList}>
                            <li>1. {t('help.q1')}</li>
                            <li>2. {t('help.q2')}</li>
                            <li>3. {t('help.q3')}</li>
                            <li>4. {t('help.q4')}</li>
                            <li>5. {t('help.q5')}</li>
                        </ul>
                        <h4 className={styles.helpSubtitle}>{t('help.other_issues')}</h4>
                        <ul className={styles.helpList}>
                            <li>6. <a href="#">{t('help.q6')} {t('help.q6_link')}</a></li>
                            <li>7. <a href="#">{t('help.q7')} {t('help.q7_link')}</a></li>
                        </ul>
                    </div>
                </aside>
            </div>
        </div>
    );
}
