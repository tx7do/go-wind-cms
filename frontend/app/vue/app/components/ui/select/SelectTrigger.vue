<script setup lang="ts">
import { SelectIcon, SelectTrigger, type SelectTriggerProps, useForwardProps } from "reka-ui"
import { ChevronDown } from "lucide-vue-next"
import { cn } from "@/lib/utils"

const props = defineProps<SelectTriggerProps & { class?: string }>()
const forwarded = useForwardProps(props)
</script>

<template>
  <SelectTrigger
    v-bind="forwarded"
    :class="cn(
      'flex h-9 w-full items-center justify-between whitespace-nowrap rounded-md border border-input bg-transparent px-3 py-2 text-sm shadow-xs ring-offset-background placeholder:text-muted-foreground focus:outline-hidden focus:ring-1 focus:ring-ring disabled:cursor-not-allowed disabled:opacity-50 [&>span]:line-clamp-1',
      props.class,
    )"
  >
    <slot />
    <SelectIcon as-child>
      <ChevronDown class="size-4 opacity-50" />
    </SelectIcon>
  </SelectTrigger>
</template>
