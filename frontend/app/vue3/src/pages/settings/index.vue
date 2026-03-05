<script setup lang="ts">
import { ref, computed } from 'vue'
import { $t } from '@/locales'
import { definePage } from 'unplugin-vue-router/runtime'

definePage({
  name: 'settings',
  meta: {
    level: 1,
  },
})

const activeMenu = ref('account')

const menuItems = computed(() => [
  { key: 'account', icon: '👤', label: $t('settings.menu.account') },
  { key: 'message', icon: '📧', label: $t('settings.menu.message') },
  { key: 'preference', icon: '⚙️', label: $t('settings.menu.preference') }
])
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
            <span class="menu-icon">{{ item.icon }}</span>
            <span class="menu-label">{{ item.label }}</span>
          </div>
        </nav>
      </aside>

      <!-- 右侧内容区 -->
      <main class="content">
        <!-- 账号与密码 Tab -->
        <div v-show="activeMenu === 'account'">
          <div class="content-header">
            <h1 class="title">{{ $t('settings.account.title') }}</h1>
            <p class="subtitle">{{ $t('settings.account.subtitle') }}</p>
          </div>

          <!-- 账号设置说明 -->
          <div class="section">
            <h2 class="section-title">{{ $t('settings.account.section_title') }}</h2>
            <p class="section-desc">{{ $t('settings.account.section_desc') }}</p>
          </div>

          <!-- 密码 -->
          <div class="setting-item">
            <div class="setting-label">
              <span class="label-text">{{ $t('settings.account.password') }}</span>
              <span class="label-status">{{ $t('settings.account.password_not_set') }}</span>
            </div>
            <button class="edit-btn">{{ $t('settings.account.edit') }}</button>
          </div>

          <!-- 绑定手机 -->
          <div class="setting-item">
            <div class="setting-label">
              <span class="label-text">{{ $t('settings.account.bind_phone') }}</span>
              <span class="label-status">{{ $t('settings.account.password_not_set') }}</span>
            </div>
            <button class="edit-btn">{{ $t('settings.account.edit') }}</button>
          </div>

          <!-- 绑定邮箱 -->
          <div class="setting-item">
            <div class="setting-label">
              <span class="label-text">{{ $t('settings.account.bind_email') }}</span>
              <span class="label-status">{{ $t('settings.account.email_not_bound') }}</span>
            </div>
            <button class="edit-btn">{{ $t('settings.account.edit') }}</button>
          </div>

          <!-- 绑定第三方账号 -->
          <div class="section">
            <h2 class="section-title">{{ $t('settings.account.third_party_title') }}</h2>
            <div class="third-party-list">
              <div class="third-party-item">
                <span class="platform-icon wechat">微</span>
                <span class="platform-name">
                  {{ $t('settings.account.wechat_account') }}<br />
                  (WeChat)
                </span>
              </div>
              <div class="third-party-item">
                <span class="platform-icon weibo">微</span>
                <span class="platform-link">{{ $t('settings.account.bind_weibo') }}</span>
              </div>
              <div class="third-party-item">
                <span class="platform-icon qq">QQ</span>
                <span class="platform-link">{{ $t('settings.account.bind_qq') }}</span>
              </div>
            </div>
          </div>
        </div>

        <!-- 消息与邮件 Tab -->
        <div v-show="activeMenu === 'message'">
          <div class="content-header">
            <h1 class="title">{{ $t('settings.message.title') }}</h1>
            <p class="subtitle">{{ $t('settings.message.subtitle') }}</p>
          </div>

          <!-- 邮件通知 -->
          <div class="section">
            <h2 class="section-title">{{ $t('settings.message.email_notifications') }}</h2>
            <p class="section-desc">{{ $t('settings.message.email_notifications_desc') }}</p>
          </div>

          <!-- 系统消息 -->
          <div class="setting-item">
            <div class="setting-label">
              <span class="label-text">{{ $t('settings.message.system_messages') }}</span>
              <span class="label-status">{{ $t('settings.message.enabled') }}</span>
            </div>
            <label class="toggle-switch">
              <input type="checkbox" checked />
              <span class="slider"></span>
            </label>
          </div>

          <!-- 评论通知 -->
          <div class="setting-item">
            <div class="setting-label">
              <span class="label-text">{{ $t('settings.message.comment_notifications') }}</span>
              <span class="label-status">{{ $t('settings.message.enabled') }}</span>
            </div>
            <label class="toggle-switch">
              <input type="checkbox" checked />
              <span class="slider"></span>
            </label>
          </div>

          <!-- 活动更新 -->
          <div class="setting-item">
            <div class="setting-label">
              <span class="label-text">{{ $t('settings.message.activity_updates') }}</span>
              <span class="label-status">{{ $t('settings.message.enabled') }}</span>
            </div>
            <label class="toggle-switch">
              <input type="checkbox" checked />
              <span class="slider"></span>
            </label>
          </div>

          <!-- 推荐内容 -->
          <div class="setting-item">
            <div class="setting-label">
              <span class="label-text">{{ $t('settings.message.recommended_content') }}</span>
              <span class="label-status">{{ $t('settings.message.enabled') }}</span>
            </div>
            <label class="toggle-switch">
              <input type="checkbox" />
              <span class="slider"></span>
            </label>
          </div>

          <!-- 邮件频率 -->
          <div class="section" style="margin-top: 30px;">
            <h2 class="section-title">{{ $t('settings.message.email_frequency') }}</h2>
          </div>

          <div class="setting-item">
            <div class="setting-label">
              <span class="label-text">{{ $t('settings.message.frequency_desc') }}</span>
            </div>
            <select class="select-control">
              <option>{{ $t('settings.message.frequency_immediately') }}</option>
              <option>{{ $t('settings.message.frequency_daily') }}</option>
              <option>{{ $t('settings.message.frequency_weekly') }}</option>
              <option>{{ $t('settings.message.frequency_monthly') }}</option>
            </select>
          </div>
        </div>

        <!-- 偏好设置 Tab -->
        <div v-show="activeMenu === 'preference'">
          <div class="content-header">
            <h1 class="title">{{ $t('settings.preference.title') }}</h1>
            <p class="subtitle">{{ $t('settings.preference.subtitle') }}</p>
          </div>

          <!-- 主题设置 -->
          <div class="section">
            <h2 class="section-title">{{ $t('settings.preference.theme_settings') }}</h2>
            <p class="section-desc">{{ $t('settings.preference.theme_desc') }}</p>
          </div>

          <div class="setting-item">
            <div class="setting-label">
              <span class="label-text">{{ $t('settings.preference.theme') }}</span>
            </div>
            <select class="select-control">
              <option>{{ $t('settings.preference.theme_light') }}</option>
              <option>{{ $t('settings.preference.theme_dark') }}</option>
              <option>{{ $t('settings.preference.theme_auto') }}</option>
            </select>
          </div>

          <!-- 语言设置 -->
          <div class="section" style="margin-top: 30px;">
            <h2 class="section-title">{{ $t('settings.preference.language_settings') }}</h2>
            <p class="section-desc">{{ $t('settings.preference.language_desc') }}</p>
          </div>

          <div class="setting-item">
            <div class="setting-label">
              <span class="label-text">{{ $t('settings.preference.language') }}</span>
            </div>
            <select class="select-control">
              <option>{{ $t('settings.preference.language_chinese') }}</option>
              <option>{{ $t('settings.preference.language_english') }}</option>
            </select>
          </div>

          <!-- 内容偏好 -->
          <div class="section" style="margin-top: 30px;">
            <h2 class="section-title">{{ $t('settings.preference.content_preferences') }}</h2>
            <p class="section-desc">{{ $t('settings.preference.content_desc') }}</p>
          </div>

          <div class="setting-item">
            <div class="setting-label">
              <span class="label-text">{{ $t('settings.preference.hide_sensitive_content') }}</span>
            </div>
            <label class="toggle-switch">
              <input type="checkbox" />
              <span class="slider"></span>
            </label>
          </div>

          <div class="setting-item">
            <div class="setting-label">
              <span class="label-text">{{ $t('settings.preference.compact_mode') }}</span>
            </div>
            <label class="toggle-switch">
              <input type="checkbox" />
              <span class="slider"></span>
            </label>
          </div>

          <div class="setting-item">
            <div class="setting-label">
              <span class="label-text">{{ $t('settings.preference.show_recommendations') }}</span>
            </div>
            <label class="toggle-switch">
              <input type="checkbox" checked />
              <span class="slider"></span>
            </label>
          </div>
        </div>
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
  padding: 20px 0;
}

.settings-container {
  max-width: 1200px;
  margin: 0 auto;
  display: flex;
  gap: 20px;
  padding: 0 20px;
}

.sidebar {
  width: 180px;
  background: var(--color-surface);
  border: 1px solid var(--color-border);
  border-radius: 8px;
  padding: 10px 0;
  height: fit-content;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.05);
}

.menu-item {
  padding: 12px 20px;
  cursor: pointer;
  display: flex;
  align-items: center;
  gap: 10px;
  color: var(--color-text-primary);
  transition: all 0.2s;

  &:hover {
    background: var(--post-accent-bg-hover);
  }

  &.active {
    background: var(--post-accent-bg-active);
    color: var(--color-brand);
  }
}

.menu-icon {
  font-size: 18px;
}

.menu-label {
  font-size: 14px;
  font-weight: 500;
}

.content {
  flex: 1;
  background: var(--color-surface);
  border: 1px solid var(--color-border);
  border-radius: 8px;
  padding: 30px;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.05);
}

.content-header {
  margin-bottom: 30px;
  padding-bottom: 20px;
  border-bottom: 1px solid var(--color-border);
}

.title {
  font-size: 24px;
  font-weight: 600;
  margin: 0 0 8px 0;
  color: var(--color-text-primary);
}

.subtitle {
  font-size: 14px;
  color: var(--color-text-secondary);
  margin: 0;
}

.section {
  margin-bottom: 30px;
}

.section-title {
  font-size: 18px;
  font-weight: 600;
  margin: 0 0 10px 0;
  color: var(--color-text-primary);
}

.section-desc {
  font-size: 14px;
  color: var(--color-text-secondary);
  line-height: 1.6;
  margin: 0;
}

.setting-item {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  padding: 20px 0;
  border-bottom: 1px solid var(--color-border);

  &:last-child {
    border-bottom: none;
  }
}

.setting-label {
  flex: 1;
}

.label-text {
  font-size: 16px;
  font-weight: 500;
  color: var(--color-text-primary);
  display: block;
  margin-bottom: 8px;
}

.label-status,
.label-value {
  font-size: 14px;
  color: var(--color-text-secondary);
}

.domain-info {
  margin-top: 8px;
}

.domain-desc {
  font-size: 14px;
  color: var(--color-text-secondary);
  margin: 0 0 4px 0;
}

.domain-url {
  font-size: 14px;
  color: var(--color-brand);
  margin: 0;
}

.edit-btn {
  padding: 6px 16px;
  background: transparent;
  border: 1px solid var(--color-brand);
  color: var(--color-brand);
  border-radius: 6px;
  cursor: pointer;
  font-size: 14px;
  font-weight: 500;
  transition: all 0.2s;
  white-space: nowrap;

  &:hover {
    background: var(--color-brand);
    color: #ffffff;
    box-shadow: 0 2px 8px rgba(102, 126, 234, 0.3);
  }

  &:active {
    transform: translateY(1px);
  }
}

.toggle-switch {
  position: relative;
  display: inline-flex;
  align-items: center;
  width: 50px;
  height: 28px;
  cursor: pointer;

  input {
    opacity: 0;
    width: 0;
    height: 0;

    &:checked + .slider {
      background-color: var(--color-brand);

      &::before {
        transform: translateX(22px);
      }
    }

    &:disabled + .slider {
      opacity: 0.6;
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
    background-color: #ccc;
    transition: 0.3s;
    border-radius: 28px;

    &::before {
      position: absolute;
      content: '';
      height: 22px;
      width: 22px;
      left: 3px;
      bottom: 3px;
      background-color: white;
      transition: 0.3s;
      border-radius: 50%;
    }
  }
}

.select-control {
  padding: 8px 12px;
  background: var(--color-surface);
  color: var(--color-text-primary);
  border: 1px solid var(--color-border);
  border-radius: 6px;
  font-size: 14px;
  cursor: pointer;
  transition: all 0.2s;
  min-width: 150px;

  &:hover {
    border-color: var(--color-brand);
  }

  &:focus {
    outline: none;
    border-color: var(--color-brand);
    box-shadow: 0 0 0 2px var(--post-accent-bg-hover);
  }

  option {
    background: var(--color-surface);
    color: var(--color-text-primary);
  }
}

.third-party-list {
  display: flex;
  gap: 20px;
  margin-top: 15px;
  flex-wrap: wrap;
}

.third-party-item {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 10px 15px;
  background: var(--color-bg);
  border: 1px solid var(--color-border);
  border-radius: 8px;
  transition: all 0.2s;

  &:hover {
    border-color: var(--color-brand);
    box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  }
}

.platform-icon {
  width: 40px;
  height: 40px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
  font-size: 20px;
  font-weight: bold;
  flex-shrink: 0;

  &.wechat {
    background: #09bb07;
  }

  &.weibo {
    background: #e6162d;
  }

  &.qq {
    background: #56b6e7;
  }
}

.platform-name {
  font-size: 13px;
  color: var(--color-text-secondary);
  line-height: 1.4;
}

.platform-link {
  font-size: 14px;
  color: var(--color-brand);
  cursor: pointer;
  font-weight: 500;
  transition: all 0.2s;

  &:hover {
    text-decoration: underline;
    opacity: 0.8;
  }
}

.help-sidebar {
  width: 260px;
  background: var(--color-surface);
  border: 1px solid var(--color-border);
  border-radius: 8px;
  padding: 20px;
  height: fit-content;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.05);
}

.help-title {
  font-size: 16px;
  font-weight: 600;
  margin: 0 0 15px 0;
  color: var(--color-text-primary);
  padding-bottom: 10px;
  border-bottom: 1px solid var(--color-border);
}

.help-subtitle {
  font-size: 14px;
  font-weight: 600;
  margin: 15px 0 10px 0;
  color: var(--color-text-primary);
}

.help-list {
  list-style: none;
  padding: 0;
  margin: 0;

  li {
    font-size: 13px;
    color: var(--color-text-secondary);
    line-height: 2;
    cursor: pointer;
    transition: color 0.2s;

    &:hover {
      color: var(--color-brand);
    }

    a {
      color: var(--color-brand);
      text-decoration: none;
      transition: all 0.2s;

      &:hover {
        text-decoration: underline;
        opacity: 0.8;
      }
    }
  }
}

// 响应式设计
@media (max-width: 1024px) {
  .settings-container {
    flex-direction: column;
  }

  .sidebar {
    width: 100%;
  }

  .menu {
    display: flex;
    overflow-x: auto;
    padding: 0;
  }

  .menu-item {
    flex-direction: column;
    min-width: 100px;
    text-align: center;
    padding: 15px;
  }

  .help-sidebar {
    width: 100%;
  }
}

@media (max-width: 768px) {
  .settings-page {
    padding: 10px 0;
  }

  .settings-container {
    padding: 0 10px;
    gap: 10px;
  }

  .content {
    padding: 20px;
  }

  .content-header {
    margin-bottom: 20px;
  }

  .title {
    font-size: 20px;
  }

  .setting-item {
    flex-direction: column;
    gap: 15px;
  }

  .edit-btn {
    align-self: flex-start;
  }

  .third-party-list {
    flex-direction: column;
    gap: 10px;
  }
}
</style>
