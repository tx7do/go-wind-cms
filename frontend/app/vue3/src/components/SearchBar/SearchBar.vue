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
  margin: 0 16px;
  max-width: 600px;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);

  --n-color: var(--header-control-bg);
  --n-color-focus: var(--header-control-bg);
  --n-text-color: var(--header-control-text);
  --n-placeholder-color: var(--header-control-text-muted);
  --n-border: 1px solid var(--header-control-border);
  --n-border-hover: 1px solid var(--color-brand);
  --n-border-focus: 1px solid var(--color-brand);
  --n-border-radius: 8px;
  --n-box-shadow-focus: 0 0 0 2px rgba(102, 126, 234, 0.1);

  :deep(.n-input-wrapper) {
    background: var(--header-control-bg);
    border-radius: 8px !important;
    box-shadow: inset 0 0 0 1px var(--header-control-border) !important;
    transition: all 0.3s;
  }

  :deep(.n-input-wrapper:hover) {
    border-radius: 8px !important;
    box-shadow: inset 0 0 0 1px var(--color-brand) !important;
  }

  :deep(.n-input.n-input--focus .n-input-wrapper) {
    border-radius: 8px !important;
    box-shadow: inset 0 0 0 1px var(--color-brand), 0 0 0 2px rgba(102, 126, 234, 0.1) !important;
  }

  :deep(.n-input__input-el) {
    color: var(--header-control-text) !important;
    font-weight: 500;
  }

  :deep(.n-input__placeholder) {
    color: var(--header-control-text-muted) !important;
  }

  :deep(.n-input__prefix .n-icon) {
    color: var(--header-control-text-muted) !important;
    transition: color 0.3s;
  }

  :deep(.n-input.n-input--focus .n-input__prefix .n-icon) {
    color: var(--color-brand) !important;
  }
}

// Responsive
@media (max-width: 1024px) {
  .search-bar {
    margin: 0 12px;
    max-width: 400px;
  }
}

@media (max-width: 768px) {
  .search-bar {
    margin: 0 8px;
    max-width: 300px;

    :deep(.n-input__input-el) {
      font-size: 14px;
    }
  }
}

@media (max-width: 480px) {
  .search-bar {
    margin: 0 4px;
    max-width: 200px;

    :deep(.n-input__input-el) {
      font-size: 12px;
    }

    :deep(.n-input__prefix .n-icon) {
      font-size: 14px !important;
    }
  }
}
</style>
