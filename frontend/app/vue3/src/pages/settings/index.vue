<script setup lang="ts">
import { ref, computed } from 'vue'
import { $t } from '@/locales'
import { definePage } from 'unplugin-vue-router/runtime'
import { XIcon } from '@/plugins/xicon'

definePage({
  name: 'settings',
  meta: {
    level: 1,
  },
})

interface MenuItem {
  key: string
  icon: string
  label: string
}

const activeMenu = ref('account')

const menuItems = computed<MenuItem[]>(() => [
  { key: 'account', icon: 'carbon:user', label: $t('settings.menu.account') },
  { key: 'message', icon: 'carbon:email', label: $t('settings.menu.message') },
  { key: 'preference', icon: 'carbon:settings', label: $t('settings.menu.preference') }
])

// 账号设置数据
const accountSettings = computed(() => ({
  title: $t('settings.account.title'),
  subtitle: $t('settings.account.subtitle'),
  sections: [
    {
      title: $t('settings.account.section_title'),
      description: $t('settings.account.section_desc'),
      items: [
        { label: $t('settings.account.password'), status: $t('settings.account.password_not_set'), type: 'button' },
        { label: $t('settings.account.bind_phone'), status: $t('settings.account.password_not_set'), type: 'button' },
        { label: $t('settings.account.bind_email'), status: $t('settings.account.email_not_bound'), type: 'button' }
      ]
    },
    {
      title: $t('settings.account.third_party_title'),
      items: []
    }
  ]
}))

// 消息设置数据
const messageSettings = computed(() => ({
  title: $t('settings.message.title'),
  subtitle: $t('settings.message.subtitle'),
  sections: [
    {
      title: $t('settings.message.email_notifications'),
      description: $t('settings.message.email_notifications_desc'),
      items: [
        { label: $t('settings.message.system_messages'), status: $t('settings.message.enabled'), type: 'toggle' },
        { label: $t('settings.message.comment_notifications'), status: $t('settings.message.enabled'), type: 'toggle' },
        { label: $t('settings.message.activity_updates'), status: $t('settings.message.enabled'), type: 'toggle' },
        { label: $t('settings.message.recommended_content'), status: $t('settings.message.enabled'), type: 'toggle' }
      ]
    },
    {
      title: $t('settings.message.email_frequency'),
      items: [
        { label: $t('settings.message.frequency_desc'), type: 'select' }
      ]
    }
  ]
}))

// 偏好设置数据
const preferenceSettings = computed(() => ({
  title: $t('settings.preference.title'),
  subtitle: $t('settings.preference.subtitle'),
  sections: [
    {
      title: $t('settings.preference.theme_settings'),
      description: $t('settings.preference.theme_desc'),
      items: [
        { label: $t('settings.preference.theme'), type: 'select' }
      ]
    },
    {
      title: $t('settings.preference.language_settings'),
      description: $t('settings.preference.language_desc'),
      items: [
        { label: $t('settings.preference.language'), type: 'select' }
      ]
    },
    {
      title: $t('settings.preference.content_preferences'),
      description: $t('settings.preference.content_desc'),
      items: [
        { label: $t('settings.preference.hide_sensitive_content'), type: 'toggle' },
        { label: $t('settings.preference.compact_mode'), type: 'toggle' },
        { label: $t('settings.preference.show_recommendations'), type: 'toggle' }
      ]
    }
  ]
}))

const currentSettings = computed(() => {
  switch (activeMenu.value) {
    case 'account':
      return accountSettings.value
    case 'message':
      return messageSettings.value
    case 'preference':
      return preferenceSettings.value
    default:
      return accountSettings.value
  }
})

</script>

<template>
  <div class="settings-page">
    <div class="settings-container">
      <!-- 左侧导航 -->
      <aside class="sidebar">
        <nav class="menu">
          <div
            v-for="item in menuItems"
            :key="item.key"
            class="menu-item"
            :class="{ active: activeMenu === item.key }"
            @click="activeMenu = item.key"
          >
            <div class="menu-icon">
              <XIcon :name="item.icon" :size="18" />
            </div>
            <span class="menu-label">{{ item.label }}</span>
          </div>
        </nav>
      </aside>

      <!-- 右侧内容区 -->
      <main class="content">
        <!-- 页面头部 -->
        <div class="content-header">
          <h1 class="title">{{ currentSettings.title }}</h1>
          <p class="subtitle">{{ currentSettings.subtitle }}</p>
        </div>

        <!-- 动态渲染各section -->
        <template v-for="section in currentSettings.sections" :key="section.title">
          <div class="section">
            <h2 class="section-title">{{ section.title }}</h2>
            <p v-if="section.description" class="section-desc">{{ section.description }}</p>

            <!-- 账户特殊处理 - 第三方账号 -->
            <div v-if="section.title === $t('settings.account.third_party_title')" class="third-party-list">
              <div class="third-party-item">
                <div class="platform-icon wechat">
                  <XIcon name="fa:wechat" :size="24" color="white" />
                </div>
                <span class="platform-name">
                  {{ $t('settings.account.wechat_account') }}<br />
                  (WeChat)
                </span>
              </div>
              <div class="third-party-item">
                <div class="platform-icon weibo">
                  <XIcon name="fa:weibo" :size="24" color="white" />
                </div>
                <span class="platform-link">{{ $t('settings.account.bind_weibo') }}</span>
              </div>
              <div class="third-party-item">
                <div class="platform-icon qq">
                  <XIcon name="fa:qq" :size="24" color="white" />
                </div>
                <span class="platform-link">{{ $t('settings.account.bind_qq') }}</span>
              </div>
            </div>

            <!-- 设置项列表 -->
            <template v-else>
              <div
                v-for="(item, idx) in section.items"
                :key="idx"
                class="setting-item"
              >
                <div class="setting-label">
                  <span class="label-text">{{ item.label }}</span>
                  <span v-if="item.status" class="label-status">{{ item.status }}</span>
                </div>

                <!-- 根据type渲染不同的控件 -->
                <div class="setting-control">
                  <!-- 编辑按钮 -->
                  <button v-if="item.type === 'button'" class="edit-btn">
                    {{ $t('settings.account.edit') }}
                  </button>

                  <!-- 切换开关 -->
                  <label v-else-if="item.type === 'toggle'" class="toggle-switch">
                    <input type="checkbox" />
                    <span class="slider"></span>
                  </label>

                  <!-- 选择框 -->
                  <select v-else-if="item.type === 'select'" class="select-control">
                    <option v-if="section.title === $t('settings.preference.theme')">
                      {{ $t('settings.preference.theme_light') }}
                    </option>
                    <option v-if="section.title === $t('settings.preference.language')">
                      {{ $t('settings.preference.language_chinese') }}
                    </option>
                  </select>
                </div>
              </div>
            </template>
          </div>
        </template>
      </main>

      <!-- 右侧帮助区 -->
      <aside class="help-sidebar">
        <div class="help-section">
          <h3 class="help-title">{{ $t('settings.help.title') }}</h3>
          <h4 class="help-subtitle">{{ $t('settings.help.account_password') }}</h4>
          <ul class="help-list">
            <li>1. {{ $t('settings.help.q1') }}</li>
            <li>2. {{ $t('settings.help.q2') }}</li>
            <li>3. {{ $t('settings.help.q3') }}</li>
            <li>4. {{ $t('settings.help.q4') }}</li>
            <li>5. {{ $t('settings.help.q5') }}</li>
          </ul>
          <h4 class="help-subtitle">{{ $t('settings.help.other_issues') }}</h4>
          <ul class="help-list">
            <li>6. <a href="#">{{ $t('settings.help.q6') }} {{ $t('settings.help.q6_link') }}</a></li>
            <li>7. <a href="#">{{ $t('settings.help.q7') }} {{ $t('settings.help.q7_link') }}</a></li>
          </ul>
        </div>
      </aside>
    </div>
  </div>
</template>

<style scoped lang="less">
.settings-page {
  min-height: 100vh;
  background: var(--color-bg);
  padding: 30px 0;
}

.settings-container {
  max-width: 1400px;
  margin: 0 auto;
  display: flex;
  gap: 24px;
  padding: 0 24px;
}

.sidebar {
  width: 200px;
  background: var(--color-surface);
  border: 1px solid var(--color-border);
  border-radius: 12px;
  padding: 12px 0;
  height: fit-content;
  position: sticky;
  top: calc(var(--layout-header-height) + 24px);
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.04);
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);

  html.dark & {
    box-shadow: 0 2px 8px rgba(0, 0, 0, 0.2);
  }
}

.menu {
  display: flex;
  flex-direction: column;
  gap: 4px;
  padding: 8px;
}

.menu-item {
  padding: 12px 16px;
  cursor: pointer;
  display: flex;
  align-items: center;
  gap: 12px;
  color: var(--color-text-primary);
  border-radius: 8px;
  transition: all 0.2s cubic-bezier(0.4, 0, 0.2, 1);
  font-weight: 500;

  &:hover {
    background: var(--post-accent-bg-hover);
    transform: translateX(2px);
  }

  &.active {
    background: var(--post-accent-bg-active);
    color: var(--color-brand);
    box-shadow: 0 0 0 2px rgba(102, 126, 234, 0.1);
  }
}

.menu-icon {
  font-size: 18px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: currentColor;

  :deep(.icon) {
    display: flex;
    align-items: center;
    justify-content: center;
  }
}

.menu-label {
  font-size: 14px;
  font-weight: 500;
}

.content {
  flex: 1;
  background: var(--color-surface);
  border: 1px solid var(--color-border);
  border-radius: 12px;
  padding: 40px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.04);
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);

  html.dark & {
    box-shadow: 0 2px 8px rgba(0, 0, 0, 0.2);
  }
}

.content-header {
  margin-bottom: 40px;
  padding-bottom: 24px;
  border-bottom: 2px solid var(--color-border);
  animation: slideInDown 0.3s cubic-bezier(0.4, 0, 0.2, 1);
}

.title {
  font-size: 28px;
  font-weight: 700;
  margin: 0 0 10px 0;
  color: var(--color-text-primary);
  background: linear-gradient(135deg, var(--color-brand) 0%, var(--color-brand-accent) 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
}

.subtitle {
  font-size: 15px;
  color: var(--color-text-secondary);
  margin: 0;
  line-height: 1.6;
}

.section {
  margin-bottom: 36px;
  animation: slideInUp 0.3s cubic-bezier(0.4, 0, 0.2, 1);

  &:last-child {
    margin-bottom: 0;
  }
}

.section-title {
  font-size: 18px;
  font-weight: 600;
  margin: 0 0 12px 0;
  color: var(--color-text-primary);
  display: flex;
  align-items: center;
  gap: 8px;

  &::before {
    content: '';
    width: 4px;
    height: 20px;
    background: var(--color-brand);
    border-radius: 2px;
    transition: all 0.3s;
  }
}

.section-desc {
  font-size: 14px;
  color: var(--color-text-secondary);
  line-height: 1.6;
  margin: 8px 0 20px 0;
  padding: 0;
}

.setting-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 18px 0;
  border-bottom: 1px solid var(--color-border);
  transition: all 0.2s;

  &:last-child {
    border-bottom: none;
  }

  &:hover {
    background: rgba(102, 126, 234, 0.02);
    padding: 18px 12px;
    margin: 0 -12px;
    border-radius: 6px;
  }
}

.setting-label {
  flex: 1;
}

.setting-control {
  display: flex;
  align-items: center;
  gap: 12px;
  flex-shrink: 0;
}

.label-text {
  font-size: 16px;
  font-weight: 500;
  color: var(--color-text-primary);
  display: block;
  margin-bottom: 6px;
}

.label-status {
  font-size: 13px;
  color: var(--color-text-secondary);
  font-weight: 400;
}

.edit-btn {
  padding: 8px 20px;
  background: transparent;
  border: 2px solid var(--color-brand);
  color: var(--color-brand);
  border-radius: 6px;
  cursor: pointer;
  font-size: 14px;
  font-weight: 600;
  transition: all 0.2s cubic-bezier(0.4, 0, 0.2, 1);
  white-space: nowrap;
  position: relative;
  overflow: hidden;

  &::before {
    content: '';
    position: absolute;
    top: 50%;
    left: 50%;
    width: 0;
    height: 0;
    background: var(--color-brand);
    border-radius: 50%;
    transform: translate(-50%, -50%);
    transition: width 0.3s, height 0.3s;
    z-index: -1;
  }

  &:hover {
    color: #ffffff;
    box-shadow: 0 4px 12px rgba(102, 126, 234, 0.3);

    &::before {
      width: 300px;
      height: 300px;
    }
  }

  &:active {
    transform: scale(0.98);
  }
}

.toggle-switch {
  position: relative;
  display: inline-flex;
  align-items: center;
  width: 52px;
  height: 32px;
  cursor: pointer;

  input {
    opacity: 0;
    width: 0;
    height: 0;

    &:checked + .slider {
      background: linear-gradient(135deg, var(--color-brand) 0%, var(--color-brand-accent) 100%);
      box-shadow: 0 0 8px rgba(102, 126, 234, 0.3);

      &::before {
        transform: translateX(22px);
        box-shadow: 0 2px 4px rgba(0, 0, 0, 0.2);
      }
    }

    &:disabled + .slider {
      opacity: 0.5;
      cursor: not-allowed;
    }
  }

  .slider {
    position: absolute;
    cursor: pointer;
    top: 0;
    left: 0;
    right: 0;
    bottom: 0;
    background-color: #d0d0d0;
    transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
    border-radius: 32px;
    box-shadow: inset 0 1px 3px rgba(0, 0, 0, 0.1);

    &::before {
      position: absolute;
      content: '';
      height: 26px;
      width: 26px;
      left: 3px;
      bottom: 3px;
      background-color: white;
      transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
      border-radius: 50%;
      box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
    }
  }
}

.select-control {
  padding: 10px 14px;
  background: var(--color-surface);
  color: var(--color-text-primary);
  border: 2px solid var(--color-border);
  border-radius: 6px;
  font-size: 14px;
  cursor: pointer;
  transition: all 0.2s cubic-bezier(0.4, 0, 0.2, 1);
  min-width: 180px;
  font-weight: 500;

  &:hover {
    border-color: var(--color-brand);
    box-shadow: 0 0 0 2px rgba(102, 126, 234, 0.1);
  }

  &:focus {
    outline: none;
    border-color: var(--color-brand);
    box-shadow: 0 0 0 3px rgba(102, 126, 234, 0.15);
  }

  option {
    background: var(--color-surface);
    color: var(--color-text-primary);
  }
}

.third-party-list {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(200px, 1fr));
  gap: 16px;
  margin-top: 20px;
}

.third-party-item {
  display: flex;
  align-items: center;
  gap: 14px;
  padding: 14px 16px;
  background: var(--color-bg);
  border: 2px solid var(--color-border);
  border-radius: 10px;
  transition: all 0.2s cubic-bezier(0.4, 0, 0.2, 1);
  cursor: pointer;

  &:hover {
    border-color: var(--color-brand);
    box-shadow: 0 4px 12px rgba(102, 126, 234, 0.15);
    transform: translateY(-2px);
    background: var(--post-accent-bg-hover);
  }
}

.platform-icon {
  width: 44px;
  height: 44px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
  font-size: 22px;
  font-weight: bold;
  flex-shrink: 0;
  box-shadow: 0 2px 6px rgba(0, 0, 0, 0.15);
  transition: all 0.2s;

  .third-party-item:hover & {
    transform: scale(1.1);
  }

  &.wechat {
    background: linear-gradient(135deg, #09bb07 0%, #07a603 100%);
  }

  &.weibo {
    background: linear-gradient(135deg, #e6162d 0%, #d01b26 100%);
  }

  &.qq {
    background: linear-gradient(135deg, #56b6e7 0%, #4a9dd1 100%);
  }

  :deep(.icon) {
    display: flex;
    align-items: center;
    justify-content: center;
    height: 100%;
    width: 100%;
  }
}

.platform-name {
  font-size: 13px;
  color: var(--color-text-secondary);
  line-height: 1.4;
  font-weight: 500;
}

.platform-link {
  font-size: 14px;
  color: var(--color-brand);
  cursor: pointer;
  font-weight: 600;
  transition: all 0.2s;

  &:hover {
    text-decoration: underline;
    color: var(--color-brand-accent);
  }
}

.help-sidebar {
  width: 280px;
  background: var(--color-surface);
  border: 1px solid var(--color-border);
  border-radius: 12px;
  padding: 24px;
  height: fit-content;
  position: sticky;
  top: calc(var(--layout-header-height) + 24px);
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.04);
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);

  html.dark & {
    box-shadow: 0 2px 8px rgba(0, 0, 0, 0.2);
  }
}

.help-title {
  font-size: 16px;
  font-weight: 700;
  margin: 0 0 18px 0;
  color: var(--color-text-primary);
  padding-bottom: 12px;
  border-bottom: 2px solid var(--color-border);
  background: linear-gradient(135deg, var(--color-brand) 0%, var(--color-brand-accent) 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
}

.help-subtitle {
  font-size: 13px;
  font-weight: 700;
  margin: 16px 0 10px 0;
  color: var(--color-text-primary);
  text-transform: uppercase;
  letter-spacing: 0.5px;
}

.help-list {
  list-style: none;
  padding: 0;
  margin: 0;

  li {
    font-size: 13px;
    color: var(--color-text-secondary);
    line-height: 1.8;
    cursor: pointer;
    transition: all 0.2s;
    padding: 4px 0;

    &:hover {
      color: var(--color-brand);
      padding-left: 4px;
    }

    a {
      color: var(--color-brand);
      text-decoration: none;
      font-weight: 600;
      transition: all 0.2s;

      &:hover {
        color: var(--color-brand-accent);
        text-decoration: underline;
      }
    }
  }
}

// 动画
@keyframes slideInDown {
  from {
    opacity: 0;
    transform: translateY(-10px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

@keyframes slideInUp {
  from {
    opacity: 0;
    transform: translateY(10px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

// 响应式设计
@media (max-width: 1200px) {
  .settings-container {
    max-width: 100%;
    gap: 20px;
  }

  .sidebar {
    width: 180px;
  }

  .content {
    padding: 32px;
  }

  .help-sidebar {
    width: 260px;
  }
}

@media (max-width: 1024px) {
  .settings-container {
    flex-direction: column;
    gap: 20px;
  }

  .sidebar {
    width: 100%;
    position: static;
    top: auto;
  }

  .menu {
    flex-direction: row;
    overflow-x: auto;
    padding: 8px;
    gap: 8px;
  }

  .menu-item {
    flex-shrink: 0;
    min-width: 120px;
    justify-content: center;
    flex-direction: column;
  }

  .content {
    padding: 28px;
  }

  .help-sidebar {
    width: 100%;
    position: static;
    top: auto;
  }

  .title {
    font-size: 24px;
  }

  .section-title {
    font-size: 16px;
  }
}

@media (max-width: 768px) {
  .settings-page {
    padding: 16px 0;
  }

  .settings-container {
    padding: 0 16px;
    gap: 16px;
  }

  .content {
    padding: 20px;
    border-radius: 8px;
  }

  .content-header {
    margin-bottom: 28px;
    padding-bottom: 16px;
  }

  .title {
    font-size: 20px;
    -webkit-text-fill-color: unset;
    background: none;
    color: var(--color-text-primary);
  }

  .subtitle {
    font-size: 13px;
  }

  .section {
    margin-bottom: 24px;
  }

  .section-title {
    font-size: 15px;

    &::before {
      width: 3px;
      height: 16px;
    }
  }

  .section-desc {
    font-size: 13px;
    margin: 6px 0 16px 0;
  }

  .setting-item {
    flex-direction: column;
    gap: 12px;
    align-items: flex-start;
    padding: 16px 0;

    &:hover {
      padding: 16px 0;
      margin: 0;
    }
  }

  .setting-control {
    width: 100%;
    justify-content: flex-end;
  }

  .edit-btn {
    width: 100%;
    padding: 10px 16px;
  }

  .toggle-switch {
    width: 48px;
    height: 28px;
  }

  .select-control {
    width: 100%;
    min-width: unset;
  }

  .third-party-list {
    grid-template-columns: repeat(auto-fill, minmax(160px, 1fr));
    gap: 12px;
  }

  .third-party-item {
    flex-direction: column;
    text-align: center;
    padding: 12px;
    gap: 10px;
  }

  .platform-icon {
    width: 40px;
    height: 40px;
    font-size: 18px;
  }

  .platform-name,
  .platform-link {
    font-size: 12px;
  }

  .help-sidebar {
    padding: 20px;
  }

  .help-title {
    font-size: 14px;
    margin-bottom: 12px;
    -webkit-text-fill-color: unset;
    background: none;
    color: var(--color-text-primary);
  }

  .help-subtitle {
    font-size: 12px;
    margin: 12px 0 8px 0;
  }

  .help-list li {
    font-size: 12px;
    line-height: 1.6;
  }
}

@media (max-width: 480px) {
  .settings-page {
    padding: 12px 0;
  }

  .settings-container {
    padding: 0 12px;
    gap: 12px;
  }

  .sidebar {
    border-radius: 6px;
  }

  .menu-item {
    min-width: 80px;
    padding: 10px 12px;
    gap: 6px;
  }

  .menu-icon {
    font-size: 18px;
  }

  .menu-label {
    font-size: 12px;
  }

  .content {
    padding: 16px;
    border-radius: 8px;
  }

  .title {
    font-size: 18px;
  }

  .label-text {
    font-size: 14px;
  }

  .label-status {
    font-size: 12px;
  }

  .section-title {
    font-size: 14px;

    &::before {
      width: 3px;
      height: 14px;
    }
  }
}

// 暗色模式优化
html.dark {
  .sidebar {
    background: rgba(27, 31, 39, 0.8);
  }

  .content {
    background: rgba(27, 31, 39, 0.8);
  }

  .toggle-switch {
    .slider {
      background-color: #3a4555;

      &::before {
        background-color: #e8ecf3;
      }
    }

    input:checked + .slider {
      &::before {
        background-color: #ffffff;
      }
    }
  }

  .third-party-item {
    background: rgba(0, 0, 0, 0.3);
    border-color: rgba(255, 255, 255, 0.1);

    &:hover {
      background: rgba(0, 0, 0, 0.5);
      border-color: var(--color-brand);
    }
  }

  .help-sidebar {
    background: rgba(27, 31, 39, 0.8);
  }

  .help-list li {
    &:hover {
      color: #69c0ff;
    }

    a {
      color: #69c0ff;

      &:hover {
        color: #91d5ff;
      }
    }
  }
}
</style>
