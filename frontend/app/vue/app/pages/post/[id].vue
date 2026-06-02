<script setup lang="ts">
const { t } = useI18n()
const route = useRoute()
const localePath = useLocalePath()
const id = computed(() => route.params.id)
</script>

<template>
  <div class="w-full">
    <LayoutPageHero
      :title="`${t('page.home.latest_posts')} #${id}`"
      description=""
      icon="carbon:document"
      size="md"
    />

    <LayoutSectionContainer>
      <div class="mb-8">
        <LayoutBackButton :label="t('navbar.top.home')" :on-click="() => navigateTo(localePath('/'))" />
      </div>

      <!-- Post detail placeholder -->
      <article class="mx-auto max-w-3xl">
        <div class="rounded-2xl border border-border bg-card p-8">
          <Skeleton class="h-8 w-3/4 mb-4" />
          <Skeleton class="h-4 w-full mb-2" />
          <Skeleton class="h-4 w-full mb-2" />
          <Skeleton class="h-4 w-2/3 mb-8" />
          <Skeleton class="h-64 w-full rounded-lg" />
          <div class="mt-8 space-y-4">
            <Skeleton class="h-4 w-full" />
            <Skeleton class="h-4 w-full" />
            <Skeleton class="h-4 w-full" />
            <Skeleton class="h-4 w-3/4" />
          </div>
        </div>
      </article>
    </LayoutSectionContainer>
  </div>
</template>
