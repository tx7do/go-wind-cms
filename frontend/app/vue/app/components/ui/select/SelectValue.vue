<script setup lang="ts">
import { SelectValue, type SelectValueProps, useForwardProps } from "reka-ui"

const props = defineProps<SelectValueProps>()
const forwarded = useForwardProps(props)
</script>

<template>
  <SelectValue v-bind="forwarded">
    <slot />
  </SelectValue>
</template>
