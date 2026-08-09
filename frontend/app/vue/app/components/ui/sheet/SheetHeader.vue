<script setup lang="ts">
import { cn } from "@/lib/utils"
import type { HTMLAttributes } from "vue"

const props = defineProps<{ class?: HTMLAttributes["class"] }>()
</script>

<template>
  <div :class="cn('flex flex-col space-y-2 text-center sm:text-start', props.class)">
    <slot />
  </div>
</template>
