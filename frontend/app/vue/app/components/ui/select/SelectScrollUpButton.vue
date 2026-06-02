<script setup lang="ts">
import { ChevronUp } from "lucide-vue-next"
import { SelectScrollUpButton, type SelectScrollUpButtonProps, useForwardProps } from "reka-ui"
import { cn } from "@/lib/utils"

const props = defineProps<SelectScrollUpButtonProps & { class?: string }>()
const forwarded = useForwardProps(props)
</script>

<template>
  <SelectScrollUpButton
    v-bind="forwarded"
    :class="cn('flex cursor-default items-center justify-center py-1', props.class)"
  >
    <ChevronUp class="size-4" />
  </SelectScrollUpButton>
</template>
