import {useState, useMemo} from 'react';
import {useTranslation} from 'react-i18next';
import {View, Text} from '@tarojs/components';

import XIcon from "@/plugins/xicon";

import './settings.scss';

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
  const {t} = useTranslation();

  const [activeMenu, setActiveMenu] = useState<'account' | 'message' | 'preference'>('account');

  // 菜单项
  const menuItems: MenuItem[] = useMemo(() => [
    {key: 'account', icon: 'carbon:user', label: t('settings.menu.account')},
    {key: 'message', icon: 'carbon:email', label: t('settings.menu.message')},
    {key: 'preference', icon: 'carbon:settings', label: t('settings.menu.preference')}
  ], [t]);

  // 账号设置数据
  const accountSettings: SettingsData = useMemo(() => ({
    title: t('settings.account.title'),
    subtitle: t('settings.account.subtitle'),
    sections: [
      {
        title: t('settings.account.section_title'),
        description: t('settings.account.section_desc'),
        items: [
          {label: t('settings.account.password'), status: t('settings.account.password_not_set'), type: 'button'},
          {label: t('settings.account.bind_phone'), status: t('settings.account.password_not_set'), type: 'button'},
          {label: t('settings.account.bind_email'), status: t('settings.account.email_not_bound'), type: 'button'}
        ]
      },
      {
        title: t('settings.account.third_party_title'),
        items: []
      }
    ]
  }), [t]);

  // 消息设置数据
  const messageSettings: SettingsData = useMemo(() => ({
    title: t('settings.message.title'),
    subtitle: t('settings.message.subtitle'),
    sections: [
      {
        title: t('settings.message.email_notifications'),
        description: t('settings.message.email_notifications_desc'),
        items: [
          {label: t('settings.message.system_messages'), status: t('settings.message.enabled'), type: 'toggle'},
          {label: t('settings.message.comment_notifications'), status: t('settings.message.enabled'), type: 'toggle'},
          {label: t('settings.message.activity_updates'), status: t('settings.message.enabled'), type: 'toggle'},
          {label: t('settings.message.recommended_content'), status: t('settings.message.enabled'), type: 'toggle'}
        ]
      },
      {
        title: t('settings.message.email_frequency'),
        items: [
          {label: t('settings.message.frequency_desc'), type: 'select'}
        ]
      }
    ]
  }), [t]);

  // 偏好设置数据
  const preferenceSettings: SettingsData = useMemo(() => ({
    title: t('settings.preference.title'),
    subtitle: t('settings.preference.subtitle'),
    sections: [
      {
        title: t('settings.preference.theme_settings'),
        description: t('settings.preference.theme_desc'),
        items: [
          {label: t('settings.preference.theme'), type: 'select'}
        ]
      },
      {
        title: t('settings.preference.language_settings'),
        description: t('settings.preference.language_desc'),
        items: [
          {label: t('settings.preference.language'), type: 'select'}
        ]
      },
      {
        title: t('settings.preference.content_preferences'),
        description: t('settings.preference.content_desc'),
        items: [
          {label: t('settings.preference.hide_sensitive_content'), type: 'toggle'},
          {label: t('settings.preference.compact_mode'), type: 'toggle'},
          {label: t('settings.preference.show_recommendations'), type: 'toggle'}
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
    <View className="settings-page">
      <View className="settings-container">
        {/* 左侧导航 */}
        <View className="sidebar">
          <View className="menu">
            {menuItems.map((item) => (
              <View
                key={item.key}
                className={`menu-item ${activeMenu === item.key ? 'active' : ''}`}
                onClick={() => setActiveMenu(item.key as 'account' | 'message' | 'preference')}
              >
                <View className="menu-icon">
                  <Text>{item.icon}</Text>
                </View>
                <Text className="menu-label">{item.label}</Text>
              </View>
            ))}
          </View>
        </View>

        {/* 中间内容区 */}
        <View className="content">
          {/* 页面头部 */}
          <View className="content-header">
            <Text className="title">{currentSettings.title}</Text>
            <Text className="subtitle">{currentSettings.subtitle}</Text>
          </View>

          {/* 动态渲染各 section */}
          {currentSettings.sections.map((section, sectionIndex) => (
            <View key={sectionIndex} className="section">
              <Text className="section-title">{section.title}</Text>
              {section.description && (
                <Text className="section-desc">{section.description}</Text>
              )}

              {/* 账户特殊处理 - 第三方账号 */}
              {section.title === t('settings.account.third_party_title') && (
                <div className="third-party-list">
                  <div className="third-party-item">
                    <div className={`platform-icon wechat`}>
                      <XIcon name="fa:wechat" size={24} color="white"/>
                    </div>
                    <Text className="platform-link">
                      {t('settings.account.bind_wechat')}
                    </Text>
                  </div>
                  <div className="third-party-item">
                    <div className={`platform-icon weibo`}>
                      <XIcon name="fa:weibo" size={24} color="white"/>
                    </div>
                    <Text className="platform-link">
                      {t('settings.account.bind_weibo')}
                    </Text>
                  </div>
                  <div className="third-party-item">
                    <div className={`platform-icon qq`}>
                      <XIcon name="fa:qq" size={24} color="white"/>
                    </div>
                    <Text className="platform-link">
                      {t('settings.account.bind_qq')}
                    </Text>
                  </div>
                </div>
              )}

              {/* 设置项列表 */}
              {section.title !== t('settings.account.third_party_title') && section.items.map((item, itemIndex) => (
                <View key={itemIndex} className="setting-item">
                  <View className="setting-label">
                    <Text className="label-text">{item.label}</Text>
                    {item.status && (
                      <Text className="label-status">{item.status}</Text>
                    )}
                  </View>

                  {/* 根据 type 渲染不同的控件 */}
                  <View className="setting-control">
                    {/* 编辑按钮 */}
                    {item.type === 'button' && (
                      <View className="edit-btn">
                        <Text>{t('settings.account.edit')}</Text>
                      </View>
                    )}

                    {/* 切换开关 */}
                    {item.type === 'toggle' && (
                      <View className="toggle-switch">
                        {/* TODO: Taro Switch 组件 */}
                      </View>
                    )}

                    {/* 选择框 */}
                    {item.type === 'select' && (
                      <View className="select-control">
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
        <View className="help-sidebar">
          <View className="help-section">
            <Text className="help-title">{t('settings.help.title')}</Text>
            <Text className="help-subtitle">{t('settings.help.account_password')}</Text>
            <View className="help-list">
              <Text>1. {t('settings.help.q1')}</Text>
              <Text>2. {t('settings.help.q2')}</Text>
              <Text>3. {t('settings.help.q3')}</Text>
              <Text>4. {t('settings.help.q4')}</Text>
              <Text>5. {t('settings.help.q5')}</Text>
            </View>
            <Text className="help-subtitle">{t('settings.help.other_issues')}</Text>
            <View className="help-list">
              <Text>6. <Text>{t('settings.help.q6')} {t('settings.help.q6_link')}</Text></Text>
              <Text>7. <Text>{t('settings.help.q7')} {t('settings.help.q7_link')}</Text></Text>
            </View>
          </View>
        </View>
      </View>
    </View>
  );
}
