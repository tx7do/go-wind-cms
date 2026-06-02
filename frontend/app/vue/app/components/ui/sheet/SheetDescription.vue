<script setup lang="ts">
import { DialogDescription, type DialogDescriptionProps, useForwardProps } from "reka-ui"
import { cn } from "@/lib/utils"
import type { HTMLAttributes } from "vue"

const props = defineProps<DialogDescriptionProps & { class?: HTMLAttributes["class"] }>()
const forwarded = useForwardProps(props)
</script>

<template>
  <DialogDescription
    v-bind="forwarded"
    :class="cn('text-sm text-muted-foreground', props.class)"
  >
    <slot />
  </DialogDescription>
</template>
