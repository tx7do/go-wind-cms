<script setup lang="ts">
import { DialogTitle, type DialogTitleProps, useForwardProps } from "reka-ui"
import { cn } from "@/lib/utils"
import type { HTMLAttributes } from "vue"

const props = defineProps<DialogTitleProps & { class?: HTMLAttributes["class"] }>()
const forwarded = useForwardProps(props)
</script>

<template>
  <DialogTitle
    v-bind="forwarded"
    :class="cn('text-lg font-semibold text-foreground', props.class)"
  >
    <slot />
  </DialogTitle>
</template>
