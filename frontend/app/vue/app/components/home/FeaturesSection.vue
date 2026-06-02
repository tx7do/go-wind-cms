<script setup lang="ts">
import { XIcon } from '@/plugins/xicon'
import { cn } from '@/lib/utils'

const { t } = useI18n()

const features = [
  { icon: 'carbon:document', title: t('page.home.flexible_content_management'), description: t('page.home.content_management_desc') },
  { icon: 'carbon:cloud', title: t('page.home.multi_tenant_architecture'), description: t('page.home.multi_tenant_desc') },
  { icon: 'carbon:security', title: t('page.home.enterprise_security'), description: t('page.home.security_desc') },
  { icon: 'carbon:analytics', title: t('page.home.advanced_analytics'), description: t('page.home.analytics_desc') },
  { icon: 'carbon:api', title: t('page.home.api_integration'), description: t('page.home.api_integration_desc') },
  { icon: 'carbon:collaborate', title: t('page.home.real_time_collaboration'), description: t('page.home.real_time_collaboration_desc') },
]
</script>

<template>
  <section class="w-full bg-background py-16">
    <div class="w-full max-w-300 mx-auto px-8 max-md:px-4">
      <div class="mb-12 text-center">
        <h2 class="inline-flex items-center gap-2 text-[2.25rem] font-extrabold leading-tight tracking-tight text-foreground">
          <XIcon icon="carbon:rocket" :size="28" class="mr-2 text-primary" />
          {{ t('page.home.platform_features') }}
        </h2>
      </div>
      <div class="grid grid-cols-[repeat(auto-fit,minmax(300px,1fr))] gap-8">
        <div
          v-for="(feature, index) in features"
          :key="index"
          :class="cn(
            'group relative overflow-hidden rounded-xl border border-border bg-card p-10 text-center',
            'shadow-sm',
            'transition-all duration-500 ease-out',
            'hover:-translate-y-1.5 hover:border-primary/40',
            'hover:shadow-[0_20px_40px_-8px_hsl(var(--primary)/0.15)]',
          )"
        >
          <div class="pointer-events-none absolute inset-0 bg-linear-to-tr from-primary/5 via-transparent to-sky-400/5 opacity-0 transition-opacity duration-500 group-hover:opacity-100" />
          <div :class="cn(
            'relative mx-auto mb-6 flex h-20 w-20 items-center justify-center rounded-xl',
            'bg-primary/10 text-3xl text-primary',
            'transition-all duration-500 ease-out',
            'group-hover:scale-110 group-hover:-rotate-6 group-hover:bg-primary group-hover:text-white',
            'group-hover:shadow-[0_8px_24px_-4px_hsl(var(--primary)/0.4)]',
          )">
            <XIcon :icon="feature.icon" :size="48" />
          </div>
          <h3 class="mb-3 text-xl font-bold text-foreground transition-colors group-hover:text-primary">
            {{ feature.title }}
          </h3>
          <p class="text-[15px] leading-relaxed text-muted-foreground">
            {{ feature.description }}
          </p>
        </div>
      </div>
    </div>
  </section>
</template>
