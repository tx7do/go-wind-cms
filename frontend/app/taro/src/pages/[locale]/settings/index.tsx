'use client';

import {useState, useMemo} from 'react';
import {useTranslation} from 'react-i18next';
import {View, Text} from '@tarojs/components';

import '../../globals.css'; // 导入全局 CSS，确保 CSS 变量可用
import styles from './settings.scss';
import XIcon from "@/plugins/xicon";

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
  const {t} = useTranslation('settings');

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
    <View className={styles.settingsPage}>
      <View className={styles.settingsContainer}>
        {/* 左侧导航 */}
        <View className={styles.sidebar}>
          <View className={styles.menu}>
            {menuItems.map((item) => (
              <View
                key={item.key}
                className={`${styles.menuItem} ${activeMenu === item.key ? styles.active : ''}`}
                onClick={() => setActiveMenu(item.key as 'account' | 'message' | 'preference')}
              >
                <View className={styles.menuIcon}>
                  <Text>{item.icon}</Text>
                </View>
                <Text className={styles.menuLabel}>{item.label}</Text>
              </View>
            ))}
          </View>
        </View>

        {/* 中间内容区 */}
        <View className={styles.content}>
          {/* 页面头部 */}
          <View className={styles.contentHeader}>
            <Text className={styles.title}>{currentSettings.title}</Text>
            <Text className={styles.subtitle}>{currentSettings.subtitle}</Text>
          </View>

          {/* 动态渲染各 section */}
          {currentSettings.sections.map((section, sectionIndex) => (
            <View key={sectionIndex} className={styles.section}>
              <Text className={styles.sectionTitle}>{section.title}</Text>
              {section.description && (
                <Text className={styles.sectionDesc}>{section.description}</Text>
              )}

              {/* 账户特殊处理 - 第三方账号 */}
              {section.title === t('account.third_party_title') && (
                <div className={styles.thirdPartyList}>
                  <div className={styles.thirdPartyItem}>
                    <div className={`${styles.platformIcon} ${styles.wechat}`}>
                      <XIcon name="fa:wechat" size={24} color="white"/>
                    </div>
                    <Text className={styles.platformLink}>
                      {t('account.bind_wechat')}
                    </Text>
                  </div>
                  <div className={styles.thirdPartyItem}>
                    <div className={`${styles.platformIcon} ${styles.weibo}`}>
                      <XIcon name="fa:weibo" size={24} color="white"/>
                    </div>
                    <Text className={styles.platformLink}>
                      {t('account.bind_weibo')}
                    </Text>
                  </div>
                  <div className={styles.thirdPartyItem}>
                    <div className={`${styles.platformIcon} ${styles.qq}`}>
                      <XIcon name="fa:qq" size={24} color="white"/>
                    </div>
                    <Text className={styles.platformLink}>
                      {t('account.bind_qq')}
                    </Text>
                  </div>
                </div>
              )}

              {/* 设置项列表 */}
              {section.title !== t('account.third_party_title') && section.items.map((item, itemIndex) => (
                <View key={itemIndex} className={styles.settingItem}>
                  <View className={styles.settingLabel}>
                    <Text className={styles.labelText}>{item.label}</Text>
                    {item.status && (
                      <Text className={styles.labelStatus}>{item.status}</Text>
                    )}
                  </View>

                  {/* 根据 type 渲染不同的控件 */}
                  <View className={styles.settingControl}>
                    {/* 编辑按钮 */}
                    {item.type === 'button' && (
                      <View className={styles.editBtn}>
                        <Text>{t('account.edit')}</Text>
                      </View>
                    )}

                    {/* 切换开关 */}
                    {item.type === 'toggle' && (
                      <View className={styles.toggleSwitch}>
                        {/* TODO: Taro Switch 组件 */}
                      </View>
                    )}

                    {/* 选择框 */}
                    {item.type === 'select' && (
                      <View className={styles.selectControl}>
                        {/* TODO: Taro Picker 组件 */}
                      </View>
                    )}
                  </View>
                </View>
              ))}
            </View>
          ))}
        </View>

        {/* 右侧帮助区 */}
        <View className={styles.helpSidebar}>
          <View className={styles.helpSection}>
            <Text className={styles.helpTitle}>{t('help.title')}</Text>
            <Text className={styles.helpSubtitle}>{t('help.account_password')}</Text>
            <View className={styles.helpList}>
              <Text>1. {t('help.q1')}</Text>
              <Text>2. {t('help.q2')}</Text>
              <Text>3. {t('help.q3')}</Text>
              <Text>4. {t('help.q4')}</Text>
              <Text>5. {t('help.q5')}</Text>
            </View>
            <Text className={styles.helpSubtitle}>{t('help.other_issues')}</Text>
            <View className={styles.helpList}>
              <Text>6. <Text>{t('help.q6')} {t('help.q6_link')}</Text></Text>
              <Text>7. <Text>{t('help.q7')} {t('help.q7_link')}</Text></Text>
            </View>
          </View>
        </View>
      </View>
    </View>
  );
}
