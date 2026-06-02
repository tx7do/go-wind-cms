<script setup lang="ts">
import { cn } from '@/lib/utils'

const props = defineProps<{
  icon: string
  title: string
  description: string
  items: string[]
}>()
</script>

<template>
  <div class="w-full">
    <LayoutPageHero
      :title="title"
      :description="description"
      :icon="icon"
      size="sm"
    />

    <section class="w-full py-12 max-md:py-8">
      <div class="mx-auto max-w-4xl px-4">
        <div :class="cn(
          'rounded-2xl border border-border bg-card p-8 md:p-10',
          'shadow-[0_20px_50px_rgba(0,0,0,0.15)]',
          'transition-all duration-300',
          'hover:border-primary/20',
        )">
          <ul class="space-y-5">
            <li
              v-for="(item, index) in items"
              :key="index"
              class="flex items-start gap-4 text-base leading-relaxed text-foreground max-md:text-sm"
            >
              <span :class="cn(
                'shrink-0 font-mono text-sm font-bold tabular-nums',
                'text-primary',
                'mt-0.5 min-w-[2rem]',
              )">
                {{ String(index + 1).padStart(2, '0') }}.
              </span>
              <span>{{ item }}</span>
            </li>
          </ul>
        </div>
      </div>
    </section>
  </div>
</template>
