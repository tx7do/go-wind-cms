<script setup lang="ts">
import { DropdownMenuTrigger } from 'reka-ui'

withDefaults(defineProps<{ asChild?: boolean }>(), { asChild: true })
</script>

<template>
  <DropdownMenuTrigger :as-child="asChild">
    <slot />
  </DropdownMenuTrigger>
</template>
