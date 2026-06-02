<script setup lang="ts">
import { XIcon } from '@/plugins/xicon'

defineProps<{
  title: string
  subtitle: string
  activeTab: string
  tabs: { key: string; label: string }[]
}>()

const emit = defineEmits<{
  'update:activeTab': [tab: string]
}>()

const { t } = useI18n()
const localePath = useLocalePath()

const tabBase = 'flex-1 cursor-pointer border-b-2 px-2 py-3 text-center text-sm font-medium transition-colors bg-transparent border-border text-muted-foreground hover:text-foreground'
const tabActive = 'border-b-primary text-primary'
const linkBtn = 'text-sm text-primary font-medium underline underline-offset-4 tracking-wide transition-colors hover:text-primary/80 cursor-pointer bg-transparent border-none'
</script>

<template>
  <div class="flex min-h-screen items-center justify-center bg-gradient-to-br from-primary/5 via-background to-background pb-12 max-md:pb-16">
    <LayoutControlPanel />

    <div class="w-full max-w-5xl mx-auto px-6">
      <div class="grid grid-cols-1 gap-12 items-center lg:grid-cols-2 lg:gap-16">
        <!-- Left: Brand -->
        <div class="hidden lg:flex relative flex-col justify-center">
          <div
            class="absolute inset-0 -z-10 opacity-60 pointer-events-none"
            :style="{
              background: 'radial-gradient(ellipse at 30% 40%, hsl(var(--primary) / 0.15) 0%, transparent 60%), radial-gradient(ellipse at 70% 70%, hsl(200 95% 60% / 0.06) 0%, transparent 55%)',
            }"
          />

          <div class="flex flex-col items-start gap-4">
            <img src="/logo.png" :alt="t('authentication.login.logo_alt')" class="h-16 w-auto" />
            <h1 class="text-4xl font-bold leading-tight">
              <span class="bg-gradient-to-r from-foreground via-primary to-sky-400 bg-clip-text text-transparent">
                {{ t('authentication.login.brand_title') }}
              </span>
            </h1>
            <p class="text-lg text-muted-foreground">{{ t('authentication.login.brand_subtitle') }}</p>
          </div>

          <div class="mt-12 space-y-4">
            <div v-for="key in ['feature_projects', 'feature_isolation', 'feature_permissions', 'feature_analytics']" :key="key" class="flex items-center gap-3 text-sm text-muted-foreground">
              <span class="flex h-5 w-5 shrink-0 items-center justify-center rounded-full bg-primary/15 text-primary text-xs font-bold">✓</span>
              <span>{{ t(`authentication.login.${key}`) }}</span>
            </div>
          </div>
        </div>

        <!-- Right: Auth Card -->
        <div class="flex justify-center lg:justify-end">
          <div class="w-full max-w-[420px] rounded-2xl border border-border bg-card p-8 shadow-lg max-md:p-6">
            <div class="mb-6">
              <h2 class="text-2xl font-bold text-foreground">{{ title }}</h2>
              <p class="mt-2 text-sm text-muted-foreground">{{ subtitle }}</p>
            </div>

            <!-- Tab Switch -->
            <div class="mb-6 flex border-b border-border">
              <button
                v-for="tab in tabs"
                :key="tab.key"
                :class="[tabBase, activeTab === tab.key ? tabActive : '']"
                @click="emit('update:activeTab', tab.key)"
              >
                {{ tab.label }}
              </button>
            </div>

            <!-- Form Content -->
            <div class="mb-6">
              <slot />
            </div>

            <!-- Switch Link -->
            <div v-if="$slots.switchLink" class="mb-4 text-center text-sm text-muted-foreground">
              <slot name="switchLink" />
            </div>

            <!-- Back Home -->
            <div class="mb-4 text-center">
              <button :class="linkBtn" @click="navigateTo(localePath('/'))">
                ← {{ t('authentication.login.back_home') }}
              </button>
            </div>

            <!-- Terms -->
            <div class="text-center leading-relaxed">
              <small class="text-xs text-muted-foreground/80">
                {{ t('authentication.login.terms_prefix') }}
                <button :class="[linkBtn, 'mx-1 text-xs']" @click="navigateTo(localePath('/terms'))">
                  {{ t('authentication.login.terms_of_service') }}
                </button>
                {{ t('authentication.login.terms_and') }}
                <button :class="[linkBtn, 'ml-1 text-xs']" @click="navigateTo(localePath('/privacy'))">
                  {{ t('authentication.login.privacy_policy') }}
                </button>
              </small>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
