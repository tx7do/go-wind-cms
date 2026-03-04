<script setup lang="ts">
import {ref} from 'vue';
import {SearchOutline} from "@vicons/ionicons5";

import {$t} from "@/locales";

const props = withDefaults(
  defineProps<{
    round?: boolean
  }>(),
  {
    round: false,
  }
);

const searchQuery = ref('');

function handleSearch() {
  console.log('Searching for:', searchQuery.value);
}
</script>

<template>
  <n-input
    :round="props.round"
    class="search-bar"
    v-model:value="searchQuery"
    :placeholder="$t('navbar.top.search_placeholder')"
    @keyup="handleSearch"
  >
    <template #prefix>
      <n-icon :component="SearchOutline"/>
    </template>
  </n-input>
</template>

<style scoped lang="less">
.search-bar {
  display: flex;
  align-items: center;
  flex: 1;
  margin: 0 20px;

  --n-color: var(--header-control-bg);
  --n-color-focus: var(--header-control-bg);
  --n-text-color: var(--header-control-text);
  --n-placeholder-color: var(--header-control-text-muted);
  --n-border: 1px solid var(--header-control-border);
  --n-border-hover: 1px solid var(--color-brand);
  --n-border-focus: 1px solid var(--color-brand);
  --n-border-radius: 8px;
  --n-box-shadow-focus: var(--header-focus-ring);

  :deep(.n-input-wrapper) {
    background: var(--header-control-bg);
    border-radius: 8px !important;
    box-shadow: inset 0 0 0 1px var(--header-control-border) !important;
  }

  :deep(.n-input-wrapper:hover) {
    border-radius: 8px !important;
    box-shadow: inset 0 0 0 1px var(--color-brand) !important;
  }

  :deep(.n-input.n-input--focus .n-input-wrapper) {
    border-radius: 8px !important;
    box-shadow: inset 0 0 0 1px var(--color-brand), var(--header-focus-ring) !important;
  }

  :deep(.n-input__input-el) {
    color: var(--header-control-text) !important;
  }

  :deep(.n-input__placeholder) {
    color: var(--header-control-text-muted) !important;
  }

  :deep(.n-input__prefix .n-icon) {
    color: var(--header-control-text-muted) !important;
  }
}
</style>
