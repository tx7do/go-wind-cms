<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'
import { XIcon } from '@/plugins/xicon'
import { cn } from '@/lib/utils'
import { fetchListCategory } from '@/api/composables/category'

const { t } = useI18n()
const localePath = useLocalePath()

const categories = ref<any[]>([])
const loading = ref(true)

let abortController: AbortController | null = null

const loadCategories = async () => {
  if (abortController) abortController.abort()
  abortController = new AbortController()

  loading.value = true
  try {
    const res = await fetchListCategory({
      paging: { page: 1, pageSize: 8 },
      formValues: { status: 'CATEGORY_STATUS_ACTIVE' },
      fieldMask: 'id,status,sortOrder,icon,code,postCount,directPostCount,parent_id,createdAt,translations.id,translations.categoryId,translations.name,translations.languageCode,translations.description',
      orderBy: ['-sortOrder', '-postCount'],
    }) as any

    categories.value = res.items || []
  } catch (error) {
    console.error('Failed to load categories:', error)
    categories.value = []
  } finally {
    loading.value = false
  }
}

onMounted(loadCategories)
onUnmounted(() => { abortController?.abort() })

const getIconName = (icon?: string): string => {
  if (!icon) return 'carbon:folder'
  return icon.includes(':') ? icon : `carbon:${icon}`
}

const getCategoryName = (category: any): string => {
  return category.translations?.[0]?.name || ''
}
</script>

<template>
  <section class="categories-section w-full max-w-300 mx-auto scroll-reveal px-8 py-12 max-md:px-4">
    <!-- Section Header -->
    <div class="mb-8 flex items-center justify-between">
      <h2 class="flex items-center gap-2 text-2xl font-extrabold tracking-tight text-foreground max-md:text-xl">
        <XIcon icon="carbon:folder-details" :size="28" class="mr-2 text-primary" />
        {{ t('page.home.categories') }}
      </h2>
      <UiButton variant="ghost" @click="navigateTo(localePath('/category'))">
        {{ t('page.home.view_all') }} →
      </UiButton>
    </div>

    <!-- Loading Skeleton -->
    <div v-if="loading" class="grid grid-cols-[repeat(auto-fill,minmax(260px,1fr))] gap-6">
      <div v-for="i in 8" :key="i" class="flex min-h-50 flex-col rounded-2xl border border-border bg-card p-6">
        <UiSkeleton class="h-35 w-full" />
      </div>
    </div>
    <template v-else>
      <!-- Desktop Grid -->
      <div class="grid grid-cols-[repeat(auto-fill,minmax(260px,1fr))] gap-6 max-md:hidden">
        <div
          v-for="category in categories"
          :key="category.id"
          class="group scroll-reveal-item relative flex min-h-50 h-full cursor-pointer flex-col rounded-2xl border border-border bg-card p-5 transition-all duration-300 ease-out hover:-translate-y-1 hover:border-primary/30 hover:bg-accent/5 hover:shadow-[0_12px_30px_rgba(34,197,94,0.06)]"
          @click="navigateTo(localePath(`/category/${category.id}`))"
        >
          <div class="pointer-events-none absolute inset-0 bg-linear-to-tr from-primary/5 via-transparent to-sky-400/5 opacity-0 transition-opacity duration-500 group-hover:opacity-100 rounded-2xl" />
          <div class="relative z-1 flex h-full flex-col justify-between">
            <div class="mb-5 flex gap-5">
              <div class="mb-4 flex h-17.5 w-17.5 shrink-0 items-center justify-center rounded-xl bg-primary/10 text-primary border border-primary/20 text-3xl transition-all duration-500 ease-out group-hover:scale-110 group-hover:bg-emerald-500 group-hover:text-white group-hover:border-emerald-500 group-hover:shadow-[0_8px_24px_-4px_hsl(var(--primary)/0.4)]">
                <XIcon :icon="getIconName(category.icon)" :size="48" />
              </div>
              <div class="flex-1 w-full min-w-0">
                <h3 class="mb-2 text-lg font-extrabold leading-tight tracking-wide text-foreground transition-colors">
                  {{ getCategoryName(category) }}
                </h3>
                <p class="text-xs text-muted-foreground mt-1">
                  <span class="mr-1 font-medium text-primary">{{ category.postCount || 0 }}</span>
                  {{ t('page.home.articles_unit') }}
                </p>
              </div>
            </div>
            <div class="mt-auto w-fit inline-flex items-center gap-1.5 rounded-md px-2.5 py-1 bg-muted border border-border text-xs text-muted-foreground transition-colors duration-300 group-hover:text-foreground/70 group-hover:border-primary/20">
              <XIcon icon="carbon:time" :size="14" />
              <span>{{ t('page.home.updated_days_ago', { days: 3 }) }}</span>
            </div>
          </div>
        </div>
      </div>

      <!-- Mobile carousel -->
      <div class="hidden max-md:flex w-full gap-4 overflow-x-auto no-scrollbar scroll-smooth px-4 pb-4 pt-2 -mx-4">
        <div
          v-for="category in categories"
          :key="category.id"
          class="group scroll-reveal-item w-[170px] shrink-0 flex cursor-pointer flex-col items-center rounded-2xl border border-border bg-card p-4 text-center transition-all duration-300 ease-out hover:-translate-y-1 hover:border-primary/30"
          @click="navigateTo(localePath(`/category/${category.id}`))"
        >
          <div class="relative z-1 flex h-full w-full flex-col items-center justify-between gap-3">
            <div class="flex h-14 w-14 shrink-0 items-center justify-center rounded-xl bg-primary/10 text-primary border border-primary/20 transition-all duration-500 ease-out group-hover:scale-110 group-hover:bg-emerald-500 group-hover:text-white">
              <XIcon :icon="getIconName(category.icon)" :size="32" />
            </div>
            <div class="flex w-full min-w-0 flex-1 flex-col items-center gap-1">
              <h3 class="line-clamp-2 text-sm font-bold leading-tight tracking-wide text-foreground">
                {{ getCategoryName(category) }}
              </h3>
              <p class="text-xs text-muted-foreground">
                <span class="mr-1 font-medium text-primary">{{ category.postCount || 0 }}</span>
                {{ t('page.home.articles_unit') }}
              </p>
            </div>
            <div class="mt-auto inline-flex w-fit items-center gap-1.5 rounded-md px-2.5 py-1 bg-muted border border-border text-xs text-muted-foreground transition-colors duration-300 group-hover:text-foreground/70 group-hover:border-primary/20">
              <XIcon icon="carbon:time" :size="14" />
              <span>{{ t('page.home.updated_days_ago', { days: 3 }) }}</span>
            </div>
          </div>
        </div>
      </div>
    </template>
  </section>
</template>
