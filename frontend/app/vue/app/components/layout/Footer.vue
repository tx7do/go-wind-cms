<script setup lang="ts">
import { XIcon } from '@/plugins/xicon'

type FooterLinkKey = 'about_us' | 'contact_us' | 'non_responsibility' | 'privacy_agreement' | 'terms_of_service'

const { t } = useI18n()
const localePath = useLocalePath()

const footerLinks: { key: FooterLinkKey; labelKey: string }[] = [
  { key: 'about_us', labelKey: 'ui.button.about_us' },
  { key: 'contact_us', labelKey: 'ui.button.contact_us' },
  { key: 'non_responsibility', labelKey: 'ui.button.non_responsibility' },
  { key: 'privacy_agreement', labelKey: 'ui.button.privacy_agreement' },
  { key: 'terms_of_service', labelKey: 'ui.button.terms_of_service' },
]

const footerRouteMap: Record<FooterLinkKey, string> = {
  about_us: '/about',
  contact_us: '/contact',
  non_responsibility: '/disclaimer',
  privacy_agreement: '/privacy',
  terms_of_service: '/terms',
}

const socialLinks = [
  { key: 'twitter', name: 'Twitter', icon: 'fa:twitter' },
  { key: 'facebook', name: 'Facebook', icon: 'fa:facebook' },
  { key: 'instagram', name: 'Instagram', icon: 'fa:instagram' },
  { key: 'telegram', name: 'Telegram', icon: 'fa:telegram' },
  { key: 'wechat', name: 'WeChat', icon: 'fa:wechat' },
  { key: 'weibo', name: 'Weibo', icon: 'fa:weibo' },
  { key: 'qq', name: 'QQ', icon: 'fa:qq' },
]

const handleFooterLinkClick = (key: FooterLinkKey) => {
  navigateTo(localePath(footerRouteMap[key]))
}

const handleSocialClick = (name: string) => {
  console.log(name)
}
</script>

<template>
  <footer class="w-full bg-card flex justify-center">
    <div class="flex w-full max-w-300 items-center justify-between gap-4 border-t border-border px-6 py-4 text-muted-foreground min-h-18 max-md:flex-col max-md:items-start max-md:gap-3 max-md:py-4 max-md:pb-5">
      <nav class="flex flex-wrap items-center gap-1" aria-label="Footer links">
        <UiButton
          v-for="link in footerLinks"
          :key="link.key"
          variant="ghost"
          class="px-1 text-muted-foreground hover:text-foreground transition-colors"
          @click="handleFooterLinkClick(link.key)"
        >
          {{ t(link.labelKey) }}
        </UiButton>
      </nav>

      <div class="flex shrink-0 items-center gap-3 max-md:w-full max-md:justify-between">
        <span class="text-sm whitespace-nowrap">
          {{ t('ui.copyright') }}
        </span>
        <UiSeparator class="h-3.5 w-px bg-border max-md:hidden" />
        <div class="flex items-center gap-1" aria-label="Social links">
          <UiButton
            v-for="social in socialLinks"
            :key="social.key"
            variant="ghost"
            size="icon"
            class="text-2xl text-muted-foreground hover:text-primary transition-colors duration-300 hover:-translate-y-0.5"
            :aria-label="social.name"
            @click="handleSocialClick(social.name)"
          >
            <XIcon :icon="social.icon" />
          </UiButton>
        </div>
      </div>
    </div>
  </footer>
</template>
