<script setup lang="ts">
import {ref, computed, watch} from 'vue'
import {
  NFlex, NInput
} from 'naive-ui'
import type {InputInst} from 'naive-ui'

interface Props {
  length?: number
  code?: string
}

interface InputTarget {
  key: string
  ref: InputInst | null
  value: string
}

const props = defineProps<Props>()
const emit = defineEmits(['update:code', 'complete'])

const length = props.length ?? 6

const inputs = ref<InputTarget[]>(Array.from({length}, (_, index) => ({
  key: `input-${index}`,
  ref: null,
  value: ''
})))

const code = computed(() => inputs.value.map(input => input.value).join(''))

const handleInput = (value: string, index: number) => {
  console.log('handleInput', value, index);

  if (value.length === 1 && /^[a-zA-Z0-9]$/.test(value)) {
    inputs.value[index].value = value
    if (index < length - 1) {
      (inputs.value[index + 1].ref as InputInst)?.focus()
    } else {
      (inputs.value[index].ref as InputInst)?.blur()
      emit('complete', code.value)
    }
  } else if (value === '') {
    if (index > 0) {
      (inputs.value[index - 1].ref as InputInst)?.focus()
    }
  } else {
    inputs.value[index].value = value.slice(0, 1)
  }
}

const handleFocus = (index: number) => {
  const inputRef = inputs.value[index].ref as InputInst
  inputRef?.select()
}

watch(code, (newCode) => {
  emit('update:code', newCode)
})
</script>

<template>
  <n-flex justify="space-around">
    <n-input
      v-for="(input, index) in inputs"
      :key="input.key"
      style="text-align: center; font-size: 38px; width: 60px; height: 60px;"
      :maxlength="1"
      size="large"
      :show-button="false"
      placeholder=""
      v-model:value="input.value"
      @input="value => handleInput(value, index)"
      @focus="() => handleFocus(index)"
      :ref="el => inputs[index].ref = el as unknown as InputInst"
      :allow-input="value => /^[a-zA-Z0-9]$/.test(value)"/>
  </n-flex>
</template>
