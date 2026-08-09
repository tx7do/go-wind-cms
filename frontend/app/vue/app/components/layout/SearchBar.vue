<script setup lang="ts">
const { t } = useI18n()
const searchQuery = ref('')

const handleSearch = () => {
  // 搜索功能尚未接入后端，暂为占位；避免在每次按键时向控制台输出用户输入
}
</script>

<template>
  <div class="mx-2 hidden h-11 max-w-80 flex-1 items-center md:flex lg:max-w-80">
    <div class="relative w-full">
      <XIcon icon="lucide:search" width="16" height="16" class="absolute start-2 top-1/2 -translate-y-1/2 text-muted-foreground" />
      <UiInput
        class="h-full w-full ps-8"
        v-model="searchQuery"
        @keyup="handleSearch"
        :placeholder="t('navbar.top.search_placeholder')"
      />
    </div>
  </div>
</template>
