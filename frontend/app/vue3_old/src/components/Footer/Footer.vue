<template>
  <footer class="footer-outer">
    <div class="footer-container stick-content">
      <nav class="footer-links" aria-label="Footer links">
        <n-button
          v-for="link in footerLinks"
          :key="link.key"
          text
          class="footer-link"
          @click="handleFooterLinkClick(link.key)"
        >
          {{ $t(link.labelKey) }}
        </n-button>
      </nav>

      <div class="footer-meta">
        <span class="copyright">{{ $t('ui.copyright') }}</span>
        <n-divider vertical class="meta-divider"/>
        <div class="social-list" aria-label="Social links">
          <n-button
            v-for="item in socialLinks"
            :key="item.key"
            text
            class="social-btn"
            :aria-label="item.name"
            @click="handleSocialClick(item.name)"
          >
            <template #icon>
              <component :is="item.icon"/>
            </template>
          </n-button>
        </div>
      </div>
    </div>
  </footer>
</template>

<script setup lang="ts">
import type {Component} from 'vue'
import {FacebookSquare, Telegram, InstagramSquare, TwitterSquare, Weixin, Weibo, Qq} from '@vicons/fa'
import {$t} from '@/locales'
import {navigateTo} from '@/router'

type FooterLinkKey =
  | 'about_us'
  | 'contact_us'
  | 'non_responsibility'
  | 'privacy_agreement'
  | 'terms_of_service'

interface FooterLinkItem {
  key: FooterLinkKey
  labelKey: `ui.button.${FooterLinkKey}`
}

interface SocialItem {
  key: 'twitter' | 'facebook' | 'instagram' | 'telegram' | 'wechat' | 'weibo' | 'qq'
  name: 'Twitter' | 'Facebook' | 'Instagram' | 'Telegram' | 'WeChat' | 'Weibo' | 'QQ'
  icon: Component
}

const footerLinks: FooterLinkItem[] = [
  {key: 'about_us', labelKey: 'ui.button.about_us'},
  {key: 'contact_us', labelKey: 'ui.button.contact_us'},
  {key: 'non_responsibility', labelKey: 'ui.button.non_responsibility'},
  {key: 'privacy_agreement', labelKey: 'ui.button.privacy_agreement'},
  {key: 'terms_of_service', labelKey: 'ui.button.terms_of_service'},
]

const footerRouteMap: Record<FooterLinkKey, string> = {
  about_us: '/about',
  contact_us: '/contact',
  non_responsibility: '/disclaimer',
  privacy_agreement: '/privacy',
  terms_of_service: '/terms',
}

const socialLinks: SocialItem[] = [
  {key: 'twitter', name: 'Twitter', icon: TwitterSquare},
  {key: 'facebook', name: 'Facebook', icon: FacebookSquare},
  {key: 'instagram', name: 'Instagram', icon: InstagramSquare},
  {key: 'telegram', name: 'Telegram', icon: Telegram},
  {key: 'wechat', name: 'WeChat', icon: Weixin},
  {key: 'weibo', name: 'Weibo', icon: Weibo},
  {key: 'qq', name: 'QQ', icon: Qq},
]

function handleFooterLinkClick(key: FooterLinkKey) {
  navigateTo(footerRouteMap[key])
}

function handleSocialClick(name: SocialItem['name']) {
  console.log(name)
}
</script>

<style scoped lang="less">
.footer-outer {
  width: 100%;
  background: var(--color-surface);
  display: flex;
  justify-content: center;
}

.footer-container {
  min-height: 72px;
  padding: var(--space-4) 0;
  color: var(--color-text-secondary);
  border-top: 1px solid var(--color-border);
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: var(--space-4);
}

.footer-links {
  display: flex;
  align-items: center;
  flex-wrap: wrap;
  gap: 4px;
}

.footer-meta {
  display: flex;
  align-items: center;
  gap: var(--space-3);
  flex-shrink: 0;
}

.footer-link {
  --n-text-color-text: var(--color-text-secondary);
  --n-text-color-text-hover: var(--color-text-primary);
  --n-text-color-text-pressed: var(--color-text-primary);
  padding: 0 4px;
}

.footer-link:deep(.n-button__content) {
  color: var(--color-text-secondary) !important;
}

.footer-link:hover:deep(.n-button__content) {
  color: var(--color-text-primary) !important;
}

.social-list {
  display: flex;
  align-items: center;
  gap: 2px;
}

.social-btn {
  font-size: 24px;
  color: var(--color-text-secondary) !important;
  transition: all 0.2s ease;
  --n-text-color-text: var(--color-text-secondary);
  --n-text-color-text-hover: var(--color-brand);
  --n-text-color-text-pressed: var(--color-brand);
}

.social-btn:hover {
  transform: translateY(-1px);
}

.copyright {
  color: var(--color-text-secondary);
  font-size: 14px;
  white-space: nowrap;
}

.meta-divider:deep(.n-divider) {
  background-color: var(--color-border);
}

@media (max-width: 960px) {
  .footer-container {
    flex-direction: column;
    align-items: flex-start;
    gap: var(--space-3);
    padding: var(--space-4) 0 var(--space-5);
  }

  .footer-meta {
    width: 100%;
    justify-content: space-between;
  }

  .meta-divider {
    display: none;
  }
}
</style>
