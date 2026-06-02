<script setup lang="ts">
import { SelectGroup, type SelectGroupProps, useForwardProps } from "reka-ui"

const props = defineProps<SelectGroupProps>()
const forwarded = useForwardProps(props)
</script>

<template>
  <SelectGroup v-bind="forwarded">
    <slot />
  </SelectGroup>
</template>
