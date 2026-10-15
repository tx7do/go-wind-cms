<script setup lang="ts">
import { XIcon } from '@/plugins/xicon'
import type { contentservicev1_Page } from '@/api/generated/app/service/v1'
import { fetchPageBySlug } from '@/api/composables/page'

const { t } = useI18n()

useHead({ title: t('page.legal.contact.title') })

// CMS 优先：后台有 slug=contact 的已发布页面（带区块）则整页走 CMS 渲染，否则回退静态内容
const cmsPage = ref<contentservicev1_Page | null>(null)
const cmsLoading = ref(true)

onMounted(async () => {
  try {
    const fetched = await fetchPageBySlug('contact')
    cmsPage.value = fetched?.sections?.length ? fetched : null
  } catch (e) {
    console.error('[Contact] Load cms page failed:', e)
    cmsPage.value = null
  } finally {
    cmsLoading.value = false
  }
})

const cards = [
  {
    icon: 'carbon:help',
    title: t('page.legal.contact.support_title'),
    desc: t('page.legal.contact.support_desc'),
    email: t('page.legal.contact.support_email'),
    accent: 'primary' as const,
  },
  {
    icon: 'mdi:briefcase-outline',
    title: t('page.legal.contact.business_title'),
    desc: t('page.legal.contact.business_desc'),
    email: t('page.legal.contact.business_email'),
    accent: 'blue' as const,
  },
  {
    icon: 'carbon:time',
    title: t('page.legal.contact.hours_title'),
    desc: t('page.legal.contact.hours_desc'),
    accent: 'slate' as const,
  },
]
</script>

<template>
  <PageCmsPageView v-if="!cmsLoading && cmsPage" :page="cmsPage" />
  <div v-else class="w-full">
    <LayoutPageHero
      :title="t('page.legal.contact.title')"
      :description="t('page.legal.contact.description')"
      icon="carbon:email"
      size="sm"
    />

    <section class="w-full py-12 max-md:py-8">
      <div class="mx-auto grid max-w-5xl grid-cols-1 gap-6 px-4 md:grid-cols-3">
        <div
          v-for="(card, index) in cards"
          :key="index"
          class="group relative flex flex-col items-center rounded-2xl border border-border bg-card p-6 text-center shadow-sm transition-all duration-300 ease-out hover:-translate-y-1 hover:shadow-lg"
          :class="card.accent === 'primary' ? 'hover:border-primary/40' : card.accent === 'blue' ? 'hover:border-blue-500/40' : 'hover:border-border'"
        >
          <div
            class="mb-4 flex h-12 w-12 items-center justify-center rounded-xl border"
            :class="card.accent === 'primary' ? 'bg-primary/10 text-primary border-primary/20' : card.accent === 'blue' ? 'bg-blue-500/10 text-blue-500 border-blue-500/20' : 'bg-muted text-muted-foreground border-border'"
          >
            <XIcon :icon="card.icon" :size="22" />
          </div>
          <h3 class="text-base font-semibold text-foreground">{{ card.title }}</h3>
          <p class="mt-2 text-xs leading-relaxed text-muted-foreground">{{ card.desc }}</p>
          <a
            v-if="card.email"
            :href="`mailto:${card.email}`"
            class="mt-3 inline-block rounded-md px-3 py-1.5 text-sm font-medium select-all transition-colors hover:underline bg-primary/5"
            :class="card.accent === 'primary' ? 'text-primary' : card.accent === 'blue' ? 'text-blue-500' : 'text-foreground'"
          >
            {{ card.email }}
          </a>
        </div>
      </div>
    </section>
  </div>
</template>
