<script setup lang="ts">
const { t } = useI18n()
const route = useRoute()
const localePath = useLocalePath()
const id = computed(() => route.params.id)
</script>

<template>
  <div class="w-full">
    <LayoutPageHero
      :title="`${t('page.tags.popular_tags')} #${id}`"
      description=""
      icon="carbon:tag"
      size="md"
    />

    <LayoutSectionContainer>
      <div class="mb-8">
        <LayoutBackButton :label="t('navbar.top.home')" :on-click="() => navigateTo(localePath('/'))" />
      </div>

      <!-- Tag posts placeholder -->
      <div class="grid grid-cols-1 gap-6 md:grid-cols-2 lg:grid-cols-3">
        <div v-for="i in 6" :key="i" class="flex min-h-50 flex-col rounded-2xl border border-border bg-card p-6">
          <Skeleton class="h-35 w-full" />
          <div class="mt-4 space-y-2">
            <Skeleton class="h-5 w-3/4" />
            <Skeleton class="h-4 w-full" />
          </div>
        </div>
      </div>
    </LayoutSectionContainer>
  </div>
</template>
