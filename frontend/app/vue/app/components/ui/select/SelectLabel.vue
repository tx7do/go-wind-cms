<script setup lang="ts">
import { SelectLabel, type SelectLabelProps, useForwardProps } from "reka-ui"
import { cn } from "@/lib/utils"

const props = defineProps<SelectLabelProps & { class?: string }>()
const forwarded = useForwardProps(props)
</script>

<template>
  <SelectLabel
    v-bind="forwarded"
    :class="cn('px-2 py-1.5 text-xs font-semibold text-muted-foreground', props.class)"
  >
    <slot />
  </SelectLabel>
</template>
