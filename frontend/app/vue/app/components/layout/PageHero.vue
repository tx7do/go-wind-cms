<script setup lang="ts">
import { XIcon } from '@/plugins/xicon'
import { cn } from '@/lib/utils'

const props = defineProps<{
  title: string
  subtitle?: string
  description?: string
  icon?: string
  iconSize?: number
  accentColor?: string
  size?: 'sm' | 'md' | 'lg'
}>()

const minHeight = computed(() =>
  props.size === 'sm' ? 'min-h-[280px]' : props.size === 'lg' ? 'min-h-[380px]' : 'min-h-[340px]'
)
const padding = computed(() =>
  props.size === 'sm' ? 'py-16' : props.size === 'lg' ? 'py-28' : 'py-20'
)
</script>

<template>
  <section :class="cn(
    'relative w-full overflow-hidden',
    'bg-gradient-to-b from-[hsl(142,70%,96%)] via-[hsl(142,50%,98%)] to-background',
    'dark:from-slate-950 dark:via-slate-900/80 dark:to-background',
    minHeight, padding,
    'flex items-center justify-center text-center',
  )">
    <!-- 装饰线条层 -->
    <div class="absolute inset-0 opacity-[0.03] pointer-events-none">
      <div
        class="absolute inset-0"
        :style="{ backgroundImage: `repeating-linear-gradient(0deg, hsl(var(--primary)) 0, hsl(var(--primary)) 1px, transparent 1px, transparent 60px)` }"
      />
      <div
        class="absolute inset-0"
        :style="{ backgroundImage: `repeating-linear-gradient(135deg, hsl(var(--primary)) 0, hsl(var(--primary)) 1px, transparent 1px, transparent 80px)` }"
      />
    </div>

    <!-- 内容层 -->
    <div :class="cn('relative z-10 w-full px-8 max-md:px-4', 'max-w-3xl')">
      <!-- 图标 -->
      <div
        v-if="icon"
        class="mb-4 flex items-center justify-center"
        :style="{ color: accentColor || 'hsl(var(--primary))' }"
      >
        <XIcon :icon="icon" :size="iconSize ?? 48" />
      </div>

      <!-- 副标题 -->
      <p
        v-if="subtitle"
        :class="cn('mb-2 text-sm font-medium uppercase tracking-[0.2em]', size === 'sm' && 'text-xs')"
        :style="{ color: accentColor || 'hsl(var(--primary))' }"
      >
        {{ subtitle }}
      </p>

      <!-- 主标题 -->
      <h1 :class="cn(
        'mb-4 font-bold leading-tight text-foreground dark:text-white',
        size === 'lg' ? 'text-4xl max-md:text-2xl' : 'text-3xl max-md:text-2xl',
        size === 'sm' && 'text-2xl max-md:text-xl',
      )">
        {{ title }}
      </h1>

      <!-- 描述 -->
      <p
        v-if="description"
        :class="cn(
          'mx-auto font-light text-muted-foreground dark:text-slate-300',
          size === 'lg' ? 'max-w-2xl text-lg max-md:text-base' : 'max-w-xl text-base max-md:text-sm',
          size === 'sm' && 'text-sm max-md:text-xs',
        )"
      >
        {{ description }}
      </p>

      <!-- 自定义内容 slot -->
      <div v-if="$slots.default" :class="cn(
        'mt-5 flex flex-wrap items-center justify-center gap-4',
        'text-sm text-muted-foreground dark:text-slate-400',
      )">
        <slot />
      </div>
    </div>
  </section>
</template>
