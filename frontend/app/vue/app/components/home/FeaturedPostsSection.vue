<script setup lang="ts">
import { XIcon } from '@/plugins/xicon'

const { t } = useI18n()
const localePath = useLocalePath()

const scrollToCategories = () => {
  const section = document.querySelector('.categories-section')
  if (section) {
    section.scrollIntoView({ behavior: 'smooth', block: 'start' })
  }
}
</script>

<template>
  <section class="w-full max-w-300 mx-auto scroll-reveal px-8 py-12 max-md:px-4">
    <div class="mb-8 flex items-center justify-between">
      <h2 class="flex items-center gap-2 text-2xl font-extrabold tracking-tight text-foreground max-md:text-xl">
        <XIcon name="carbon:star-filled" :size="28" class="mr-2 text-primary" />
        {{ t('page.home.featured_posts') }}
      </h2>
      <UIButton variant="ghost" @click="navigateTo(localePath('/post'))">
        {{ t('page.home.view_all') }} →
      </UIButton>
    </div>
    <div class="w-full">
      <!-- PostList component placeholder - will be replaced with actual PostList when business components are created -->
      <div class="grid grid-cols-1 gap-6 md:grid-cols-2 lg:grid-cols-3">
        <div v-for="i in 3" :key="i" class="flex min-h-50 flex-col rounded-2xl border border-border bg-card p-6">
          <Skeleton class="h-35 w-full" />
        </div>
      </div>
    </div>
    <div class="mt-8 flex w-full items-center justify-center max-md:hidden">
      <div class="flex items-center gap-4">
        <span class="text-2xl font-extrabold leading-tight tracking-tight text-primary max-md:text-xl">
          {{ t('page.home.explore_more_categories') }}
        </span>
        <div
          class="group flex h-9 w-9 cursor-pointer items-center justify-center rounded-full border border-primary/20 bg-primary/10 transition-all duration-300 hover:border-primary hover:bg-primary hover:shadow-[0_0_20px_hsl(var(--primary)/0.3)]"
          @click="scrollToCategories"
        >
          <XIcon name="carbon:arrow-down" :size="24" class="text-primary transition-colors group-hover:text-white" />
        </div>
      </div>
    </div>
  </section>
</template>
