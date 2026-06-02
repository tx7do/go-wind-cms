<script setup lang="ts">
import { fetchListCategory, getCategoryName, getCategoryDescription, getCategoryThumbnail } from '@/api/composables/category'
import type { ListCategoryParams } from '@/api/composables/category'
import { XIcon } from '@/plugins/xicon'

const { t } = useI18n()
const localePath = useLocalePath()

const categories = ref<any[]>([])
const loading = ref(true)

function getIconName(icon: string | undefined) {
  return icon ? `carbon:${icon}` : 'carbon:folder'
}

onMounted(async () => {
  try {
    const params: ListCategoryParams = {
      paging: { page: 1, pageSize: 50 },
      isTenantUser: true,
    }
    const res = await fetchListCategory(params)
    categories.value = res?.items || res?.data || []
  } catch (e) {
    console.error('[Category List] 加载失败:', e)
  } finally {
    loading.value = false
  }
})
</script>

<template>
  <div class="w-full">
    <LayoutPageHero
      :title="t('page.categories.categories')"
      icon="carbon:folder-details"
      size="md"
    />

    <section class="w-full py-12 max-md:py-8">
      <LayoutSectionContainer>
        <!-- Loading -->
        <div v-if="loading" class="grid gap-6 sm:grid-cols-2 lg:grid-cols-3">
          <div v-for="i in 6" :key="i" class="rounded-xl border border-border p-4">
            <UiSkeleton class="h-35 w-full" />
            <div class="mt-4 space-y-2">
              <UiSkeleton class="h-5 w-3/4" />
              <UiSkeleton class="h-4 w-full" />
            </div>
          </div>
        </div>

        <!-- Empty -->
        <div v-else-if="categories.length === 0" class="py-20 text-center">
          <UiAppEmpty :message="t('page.categories.no_categories')" />
        </div>

        <!-- Category Grid -->
        <div v-else class="grid gap-6 sm:grid-cols-2 lg:grid-cols-3">
          <NuxtLink
            v-for="cat in categories"
            :key="cat.id"
            :to="localePath(`/category/${cat.id}`)"
            class="group rounded-xl border border-border bg-card p-6 transition-colors hover:border-primary/30 hover:bg-primary/5"
          >
            <div class="flex items-center gap-3">
              <div class="flex h-12 w-12 items-center justify-center rounded-lg bg-primary/10 text-primary">
                <XIcon :icon="getIconName(cat.icon)" :size="24" />
              </div>
              <div>
                <h3 class="text-base font-semibold text-foreground group-hover:text-primary">
                  {{ getCategoryName(cat, t('page.categories.category_untitled')) }}
                </h3>
                <p class="mt-1 line-clamp-2 text-sm text-muted-foreground">
                  {{ getCategoryDescription(cat) }}
                </p>
              </div>
            </div>
          </NuxtLink>
        </div>
      </LayoutSectionContainer>
    </section>
  </div>
</template>
