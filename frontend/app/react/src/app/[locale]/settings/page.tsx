'use client';

import {useState, useMemo} from 'react';
import {useTranslations} from 'next-intl';
import XIcon from '@/plugins/xicon';

import '../../globals.css'; // 导入全局 CSS，确保 CSS 变量可用
import styles from './settings.module.css';

interface MenuItem {
    key: string;
    icon: string;
    label: string;
}

interface SettingItem {
    label: string;
    status?: string;
    type: 'button' | 'toggle' | 'select';
}

interface Section {
    title: string;
    description?: string;
    items: SettingItem[];
}

interface SettingsData {
    title: string;
    subtitle: string;
    sections: Section[];
}

export default function SettingsPage() {
    const t = useTranslations('settings');

    const [activeMenu, setActiveMenu] = useState<'account' | 'message' | 'preference'>('account');

    // 菜单项
    const menuItems: MenuItem[] = useMemo(() => [
        {key: 'account', icon: 'carbon:user', label: t('menu.account')},
        {key: 'message', icon: 'carbon:email', label: t('menu.message')},
        {key: 'preference', icon: 'carbon:settings', label: t('menu.preference')}
    ], [t]);

    // 账号设置数据
    const accountSettings: SettingsData = useMemo(() => ({
        title: t('account.title'),
        subtitle: t('account.subtitle'),
        sections: [
            {
                title: t('account.section_title'),
                description: t('account.section_desc'),
                items: [
                    {label: t('account.password'), status: t('account.password_not_set'), type: 'button'},
                    {label: t('account.bind_phone'), status: t('account.password_not_set'), type: 'button'},
                    {label: t('account.bind_email'), status: t('account.email_not_bound'), type: 'button'}
                ]
            },
            {
                title: t('account.third_party_title'),
                items: []
            }
        ]
    }), [t]);

    // 消息设置数据
    const messageSettings: SettingsData = useMemo(() => ({
        title: t('message.title'),
        subtitle: t('message.subtitle'),
        sections: [
            {
                title: t('message.email_notifications'),
                description: t('message.email_notifications_desc'),
                items: [
                    {label: t('message.system_messages'), status: t('message.enabled'), type: 'toggle'},
                    {label: t('message.comment_notifications'), status: t('message.enabled'), type: 'toggle'},
                    {label: t('message.activity_updates'), status: t('message.enabled'), type: 'toggle'},
                    {label: t('message.recommended_content'), status: t('message.enabled'), type: 'toggle'}
                ]
            },
            {
                title: t('message.email_frequency'),
                items: [
                    {label: t('message.frequency_desc'), type: 'select'}
                ]
            }
        ]
    }), [t]);

    // 偏好设置数据
    const preferenceSettings: SettingsData = useMemo(() => ({
        title: t('preference.title'),
        subtitle: t('preference.subtitle'),
        sections: [
            {
                title: t('preference.theme_settings'),
                description: t('preference.theme_desc'),
                items: [
                    {label: t('preference.theme'), type: 'select'}
                ]
            },
            {
                title: t('preference.language_settings'),
                description: t('preference.language_desc'),
                items: [
                    {label: t('preference.language'), type: 'select'}
                ]
            },
            {
                title: t('preference.content_preferences'),
                description: t('preference.content_desc'),
                items: [
                    {label: t('preference.hide_sensitive_content'), type: 'toggle'},
                    {label: t('preference.compact_mode'), type: 'toggle'},
                    {label: t('preference.show_recommendations'), type: 'toggle'}
                ]
            }
        ]
    }), [t]);

    // 当前设置数据
    const currentSettings: SettingsData = useMemo(() => {
        switch (activeMenu) {
            case 'account':
                return accountSettings;
            case 'message':
                return messageSettings;
            case 'preference':
                return preferenceSettings;
            default:
                return accountSettings;
        }
    }, [activeMenu, accountSettings, messageSettings, preferenceSettings]);

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
                    {/* 页面头部 */}
                    <div className={styles.contentHeader}>
                        <h1 className={styles.title}>{currentSettings.title}</h1>
                        <p className={styles.subtitle}>{currentSettings.subtitle}</p>
                    </div>

                    {/* 动态渲染各 section */}
                    {currentSettings.sections.map((section, sectionIndex) => (
                        <div key={sectionIndex} className={styles.section}>
                            <h2 className={styles.sectionTitle}>{section.title}</h2>
                            {section.description && (
                                <p className={styles.sectionDesc}>{section.description}</p>
                            )}

                            {/* 账户特殊处理 - 第三方账号 */}
                            {section.title === t('account.third_party_title') ? (
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
                            ) : (
                                /* 设置项列表 */
                                section.items.map((item, itemIndex) => (
                                    <div key={itemIndex} className={styles.settingItem}>
                                        <div className={styles.settingLabel}>
                                            <span className={styles.labelText}>{item.label}</span>
                                            {item.status && (
                                                <span className={styles.labelStatus}>{item.status}</span>
                                            )}
                                        </div>

                                        {/* 根据 type 渲染不同的控件 */}
                                        <div className={styles.settingControl}>
                                            {/* 编辑按钮 */}
                                            {item.type === 'button' && (
                                                <button className={styles.editBtn}>
                                                    {t('account.edit')}
                                                </button>
                                            )}

                                            {/* 切换开关 */}
                                            {item.type === 'toggle' && (
                                                <label className={styles.toggleSwitch}>
                                                    <input type="checkbox"/>
                                                    <span className={styles.slider}></span>
                                                </label>
                                            )}

                                            {/* 选择框 */}
                                            {item.type === 'select' && (
                                                <select className={styles.selectControl}>
                                                    {section.title === t('preference.theme') && (
                                                        <>
                                                            <option>{t('preference.theme_light')}</option>
                                                            <option>{t('preference.theme_dark')}</option>
                                                            <option>{t('preference.theme_auto')}</option>
                                                        </>
                                                    )}
                                                    {section.title === t('preference.language') && (
                                                        <>
                                                            <option>{t('preference.language_chinese')}</option>
                                                            <option>{t('preference.language_english')}</option>
                                                        </>
                                                    )}
                                                    {section.title === t('message.email_frequency') && (
                                                        <>
                                                            <option>{t('message.frequency_immediately')}</option>
                                                            <option>{t('message.frequency_daily')}</option>
                                                            <option>{t('message.frequency_weekly')}</option>
                                                            <option>{t('message.frequency_monthly')}</option>
                                                        </>
                                                    )}
                                                </select>
                                            )}
                                        </div>
                                    </div>
                                ))
                            )}
                        </div>
                    ))}
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
