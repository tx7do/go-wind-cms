<script setup lang="ts">
import { computed, ref } from 'vue'
import { RotateCcw, Monitor, Sun, Moon } from 'lucide-vue-next'
import {
  COLOR_PRESETS,
  preferencesManager,
  updatePreferences,
} from '@/core/preferences'
import { usePreferences } from '@/core/preferences/use-preferences'
import type { Preferences, ThemeModeType } from '@/core/preferences/types'

import {
  Sheet,
  SheetContent,
  SheetHeader,
  SheetTitle,
  SheetDescription,
} from '@/components/ui/sheet'
import { Label } from '@/components/ui/label'
import { Button } from '@/components/ui/button'
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select'

const props = defineProps<{
  open: boolean
}>()

const emit = defineEmits<{
  'update:open': [value: boolean]
}>()

const isOpen = computed({
  get: () => props.open,
  set: (val) => emit('update:open', val),
})

const { locale } = usePreferences()
const preferences = preferencesManager.getPreferences()

// 主题模式选项
const themeModes: { label: string; value: ThemeModeType; icon: typeof Monitor }[] = [
  { label: '自动', value: 'auto', icon: Monitor },
  { label: '浅色', value: 'light', icon: Sun },
  { label: '深色', value: 'dark', icon: Moon },
]

// 语言选项
const locales = [
  { label: '中文', value: 'zh-CN' },
  { label: 'English', value: 'en-US' },
]

// 页面切换动画选项
const transitionOptions = [
  { label: '淡入淡出', value: 'fade' },
  { label: '向下滑入', value: 'fade-down' },
  { label: '滑入', value: 'fade-slide' },
  { label: '向上滑入', value: 'fade-up' },
]

// 响应式状态
const currentMode = ref<ThemeModeType>(preferences.theme.mode)
const selectedColorIdx = ref(
  COLOR_PRESETS.findIndex((p) => p.type === preferences.theme.builtinType),
)
const currentTransition = ref(preferences.transition.name)

function applyThemeMode(mode: ThemeModeType) {
  currentMode.value = mode
  updatePreferences({ theme: { mode } })
}

function applyColorPreset(idx: number) {
  selectedColorIdx.value = idx
  const preset = COLOR_PRESETS[idx]
  if (preset) {
    updatePreferences({
      theme: {
        builtinType: preset.type,
        colorPrimary: preset.color,
      },
    })
  }
}

function applyLocale(value: string) {
  updatePreferences({ app: { locale: value as Preferences['app']['locale'] } })
}

function applyTransition(name: string) {
  currentTransition.value = name
  updatePreferences({ transition: { name } })
}

function handleReset() {
  preferencesManager.resetPreferences()
  const prefs = preferencesManager.getPreferences()
  currentMode.value = prefs.theme.mode
  currentTransition.value = prefs.transition.name
  selectedColorIdx.value = COLOR_PRESETS.findIndex(
    (p) => p.type === prefs.theme.builtinType,
  )
}
</script>

<template>
  <Sheet v-model:open="isOpen">
    <SheetContent side="right" class="w-full sm:max-w-md overflow-y-auto">
      <SheetHeader>
        <SheetTitle>偏好设置</SheetTitle>
        <SheetDescription>自定义你的浏览体验</SheetDescription>
      </SheetHeader>

      <div class="mt-6 space-y-8">
        <!-- 主题模式 -->
        <div class="space-y-2">
          <Label>主题模式</Label>
          <div class="grid grid-cols-3 gap-3">
            <button
              v-for="item in themeModes"
              :key="item.value"
              class="flex flex-col items-center gap-2 rounded-lg border-2 p-4 transition-colors hover:bg-accent"
              :class="
                currentMode === item.value
                  ? 'border-primary bg-primary/10'
                  : 'border-border'
              "
              @click="applyThemeMode(item.value)"
            >
              <component :is="item.icon" class="size-5" />
              <span class="text-xs">{{ item.label }}</span>
            </button>
          </div>
        </div>

        <!-- 主题色 -->
        <div class="space-y-2">
          <Label>主题色</Label>
          <div class="flex flex-wrap gap-3">
            <button
              v-for="(preset, idx) in COLOR_PRESETS"
              :key="preset.type"
              class="size-8 rounded-full border-2 transition-transform hover:scale-110"
              :class="
                selectedColorIdx === idx
                  ? 'border-primary ring-2 ring-primary/30'
                  : 'border-border'
              "
              :style="
                preset.color
                  ? { backgroundColor: preset.color }
                  : {
                      background:
                        'conic-gradient(red, orange, yellow, green, blue, indigo, violet, red)',
                    }
              "
              :title="preset.type"
              @click="applyColorPreset(idx)"
            />
          </div>
        </div>

        <!-- 语言 -->
        <div class="space-y-2">
          <Label>语言</Label>
          <div class="grid grid-cols-2 gap-3">
            <button
              v-for="item in locales"
              :key="item.value"
              class="rounded-lg border-2 px-4 py-2 text-sm transition-colors hover:bg-accent"
              :class="
                locale === item.value
                  ? 'border-primary bg-primary/10'
                  : 'border-border'
              "
              @click="applyLocale(item.value)"
            >
              {{ item.label }}
            </button>
          </div>
        </div>

        <!-- 页面切换动画 -->
        <div class="space-y-2">
          <Label>页面切换动画</Label>
          <Select
            :model-value="currentTransition"
            @update:model-value="applyTransition"
          >
            <SelectTrigger>
              <SelectValue placeholder="选择动画效果" />
            </SelectTrigger>
            <SelectContent>
              <SelectItem
                v-for="item in transitionOptions"
                :key="item.value"
                :value="item.value"
              >
                {{ item.label }}
              </SelectItem>
            </SelectContent>
          </Select>
        </div>

        <!-- 重置按钮 -->
        <div class="pt-4">
          <Button
            variant="outline"
            class="w-full"
            @click="handleReset"
          >
            <RotateCcw class="size-4" />
            恢复默认设置
          </Button>
        </div>
      </div>
    </SheetContent>
  </Sheet>
</template>
