<script setup lang="ts">
import { ChevronDown } from "lucide-vue-next"
import { SelectScrollDownButton, type SelectScrollDownButtonProps, useForwardProps } from "reka-ui"
import { cn } from "@/lib/utils"

const props = defineProps<SelectScrollDownButtonProps & { class?: string }>()
const forwarded = useForwardProps(props)
</script>

<template>
  <SelectScrollDownButton
    v-bind="forwarded"
    :class="cn('flex cursor-default items-center justify-center py-1', props.class)"
  >
    <ChevronDown class="size-4" />
  </SelectScrollDownButton>
</template>
