<script setup lang="ts">
import { DialogTrigger, type DialogTriggerProps, useForwardProps } from "reka-ui"

const props = defineProps<DialogTriggerProps & { class?: string }>()
const forwarded = useForwardProps(props)
</script>

<template>
  <DialogTrigger v-bind="forwarded" :class="props.class">
    <slot />
  </DialogTrigger>
</template>
