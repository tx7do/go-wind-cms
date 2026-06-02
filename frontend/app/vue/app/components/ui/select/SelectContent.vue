<script setup lang="ts">
import {
  SelectContent,
  type SelectContentEmits,
  type SelectContentProps,
  SelectPortal,
  SelectViewport,
  useForwardPropsEmits,
} from "reka-ui"
import { cn } from "@/lib/utils"
import SelectScrollUpButton from "./SelectScrollUpButton.vue"
import SelectScrollDownButton from "./SelectScrollDownButton.vue"

const props = withDefaults(
  defineProps<SelectContentProps & { class?: string }>(),
  { position: "popper" },
)
const emits = defineEmits<SelectContentEmits>()
const forwarded = useForwardPropsEmits(props, emits)
</script>

<template>
  <SelectPortal>
    <SelectContent
      v-bind="{ ...forwarded, ...$attrs }"
      :class="cn(
        'relative z-50 max-h-(--reka-select-content-available-height) min-w-[8rem] overflow-x-hidden overflow-y-auto rounded-md border bg-popover text-popover-foreground shadow-md',
        position === 'popper' && 'data-[side=bottom]:translate-y-1 data-[side=left]:-translate-x-1 data-[side=right]:translate-x-1 data-[side=top]:-translate-y-1',
        props.class,
      )"
    >
      <SelectScrollUpButton />
      <SelectViewport
        :class="cn('p-1', position === 'popper' && 'h-(--reka-select-trigger-height) w-full min-w-(--reka-select-trigger-width)')"
      >
        <slot />
      </SelectViewport>
      <SelectScrollDownButton />
    </SelectContent>
  </SelectPortal>
</template>
