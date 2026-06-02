<script setup lang="ts">
import { ref, type HTMLAttributes } from 'vue'
import { cn } from '@/lib/utils'

const props = withDefaults(defineProps<{
  length?: number
  modelValue?: string
  disabled?: boolean
  type?: 'number' | 'text'
  class?: HTMLAttributes['class']
}>(), {
  length: 6,
  modelValue: '',
  disabled: false,
  type: 'number',
})

const emits = defineEmits<{
  'update:modelValue': [value: string]
  'complete': [value: string]
}>()

const inputRefs = ref<(HTMLInputElement | null)[]>([])

function focusInput(index: number) {
  const input = inputRefs.value[index]
  if (input) {
    input.focus()
    input.select()
  }
}

function handleInputChange(index: number, inputValue: string) {
  const char = inputValue.slice(-1)
  if (props.type === 'number' && char && !/^\d$/.test(char)) return

  const newValue = props.modelValue.split('')
  newValue[index] = char
  const joined = newValue.join('').slice(0, props.length)
  emits('update:modelValue', joined)

  if (char && index < props.length - 1) {
    focusInput(index + 1)
  }

  if (joined.length === props.length) {
    emits('complete', joined)
  }
}

function handleKeyDown(index: number, e: KeyboardEvent) {
  if (e.key === 'Backspace') {
    e.preventDefault()
    if (props.modelValue[index]) {
      const newValue = props.modelValue.split('')
      newValue[index] = ''
      emits('update:modelValue', newValue.join(''))
    } else if (index > 0) {
      focusInput(index - 1)
      const newValue = props.modelValue.split('')
      newValue[index - 1] = ''
      emits('update:modelValue', newValue.join(''))
    }
  } else if (e.key === 'ArrowLeft' && index > 0) {
    e.preventDefault()
    focusInput(index - 1)
  } else if (e.key === 'ArrowRight' && index < props.length - 1) {
    e.preventDefault()
    focusInput(index + 1)
  }
}

function handlePaste(e: ClipboardEvent) {
  e.preventDefault()
  const pasted = (e.clipboardData?.getData('text') || '').slice(0, props.length)
  if (props.type === 'number' && !/^\d+$/.test(pasted)) return
  emits('update:modelValue', pasted)
  if (pasted.length === props.length) {
    emits('complete', pasted)
  }
  focusInput(Math.min(pasted.length, props.length - 1))
}
</script>

<template>
  <div :class="cn('flex items-center gap-2', props.class)">
    <input
      v-for="(_, index) in length"
      :key="index"
      :ref="(el) => { inputRefs[index] = el as HTMLInputElement | null }"
      type="text"
      :inputMode="type === 'number' ? 'numeric' : 'text'"
      :maxlength="1"
      :value="modelValue[index] || ''"
      :disabled="disabled"
      class="flex h-10 w-10 items-center justify-center rounded-md border border-input bg-transparent text-center text-sm shadow-sm transition-colors placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-1 focus-visible:ring-ring disabled:cursor-not-allowed disabled:opacity-50 md:h-12 md:w-12 md:text-base"
      @input="handleInputChange(index, ($event.target as HTMLInputElement).value)"
      @keydown="handleKeyDown(index, $event)"
      @paste="handlePaste"
      @focus="($event.target as HTMLInputElement).select()"
    />
  </div>
</template>
