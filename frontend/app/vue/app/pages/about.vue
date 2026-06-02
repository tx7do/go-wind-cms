<script setup lang="ts">
const { t } = useI18n()
</script>

<template>
  <div class="w-full">
    <LayoutPageHero
      :title="t('page.about.title')"
      :subtitle="t('page.about.subtitle')"
      :description="t('page.about.description')"
      icon="carbon:information"
      :icon-size="56"
      size="lg"
    />

    <!-- About Section -->
    <section class="w-full bg-background py-20">
      <div class="w-full max-w-[1200px] mx-auto px-8 max-md:px-4">
        <div class="grid grid-cols-1 gap-12 lg:grid-cols-2 lg:gap-16">
          <div class="antialiased">
            <h2 class="mb-6 text-3xl font-bold text-foreground">{{ t('page.about.about_us') }}</h2>
            <p class="mb-4 text-[15px] leading-loose tracking-normal text-muted-foreground">{{ t('page.about.about_us_desc_1') }}</p>
            <p class="mb-4 text-[15px] leading-loose tracking-normal text-muted-foreground">{{ t('page.about.about_us_desc_2') }}</p>
            <p class="text-[15px] leading-loose tracking-normal text-muted-foreground">{{ t('page.about.about_us_desc_3') }}</p>
          </div>
          <div class="grid grid-cols-2 gap-4">
            <div v-for="stat in [
              { num: '10K+', label: t('page.about.stat_users') },
              { num: '500+', label: t('page.about.stat_projects') },
              { num: '99.9%', label: t('page.about.stat_uptime') },
              { num: '24/7', label: t('page.about.stat_support') },
            ]" :key="stat.label" class="rounded-xl border border-border bg-card p-6 text-center shadow-sm">
              <div class="mb-3 text-3xl font-bold text-primary">{{ stat.num }}</div>
              <div class="text-[15px] font-medium text-muted-foreground">{{ stat.label }}</div>
            </div>
          </div>
        </div>
      </div>
    </section>

    <!-- Features -->
    <section class="w-full bg-muted/30 py-20">
      <div class="w-full max-w-[1200px] mx-auto px-8 max-md:px-4">
        <div class="mb-12 text-center">
          <h2 class="mb-4 text-3xl font-bold text-foreground">{{ t('page.about.features') }}</h2>
          <p class="text-muted-foreground">{{ t('page.about.features_desc') }}</p>
        </div>
        <div class="grid grid-cols-[repeat(auto-fit,minmax(280px,1fr))] gap-6">
          <div
            v-for="feature in [
              { icon: 'carbon:document-add', title: t('page.about.feature_content'), description: t('page.about.feature_content_desc') },
              { icon: 'carbon:cloud-upload', title: t('page.about.feature_multi_tenant'), description: t('page.about.feature_multi_tenant_desc') },
              { icon: 'carbon:security', title: t('page.about.feature_security'), description: t('page.about.feature_security_desc') },
              { icon: 'carbon:api', title: t('page.about.feature_api'), description: t('page.about.feature_api_desc') },
              { icon: 'carbon:collaborate', title: t('page.about.feature_collaboration'), description: t('page.about.feature_collaboration_desc') },
              { icon: 'carbon:analytics', title: t('page.about.feature_analytics'), description: t('page.about.feature_analytics_desc') },
            ]"
            :key="feature.title"
            class="group rounded-xl border border-border bg-card p-8 text-center shadow-sm transition-all duration-300 hover:-translate-y-1 hover:border-primary/40 hover:shadow-lg"
          >
            <div class="mx-auto mb-4 flex h-16 w-16 items-center justify-center rounded-xl bg-primary/10 text-primary transition-all duration-300 group-hover:bg-primary group-hover:text-white group-hover:shadow-md">
              <XIcon :icon="feature.icon" :size="32" />
            </div>
            <h3 class="mb-2 text-lg font-bold text-foreground transition-colors group-hover:text-primary">{{ feature.title }}</h3>
            <p class="text-sm leading-relaxed text-muted-foreground">{{ feature.description }}</p>
          </div>
        </div>
      </div>
    </section>

    <!-- CTA -->
    <section class="w-full bg-primary/5 py-20">
      <div class="w-full max-w-[1200px] mx-auto px-8 max-md:px-4">
        <div class="mx-auto max-w-[800px] text-center">
          <h2 class="mb-4 text-3xl font-bold text-foreground">{{ t('page.about.cta_title') }}</h2>
          <p class="mb-8 text-lg text-muted-foreground">{{ t('page.about.cta_desc') }}</p>
          <div class="flex flex-wrap justify-center gap-4">
            <UiButton size="lg" class="rounded-lg px-6 py-2.5 font-medium">{{ t('page.about.cta_explore') }}</UiButton>
            <UiButton variant="outline" size="lg" class="rounded-lg px-6 py-2.5 font-medium">{{ t('page.about.cta_contact') }}</UiButton>
          </div>
        </div>
      </div>
    </section>
  </div>
</template>
