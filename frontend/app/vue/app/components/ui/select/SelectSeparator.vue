<script setup lang="ts">
import { SelectSeparator, type SelectSeparatorProps, useForwardProps } from "reka-ui"
import { cn } from "@/lib/utils"

const props = defineProps<SelectSeparatorProps & { class?: string }>()
const forwarded = useForwardProps(props)
</script>

<template>
  <SelectSeparator
    v-bind="forwarded"
    :class="cn('-mx-1 my-1 h-px bg-muted', props.class)"
  />
</template>
