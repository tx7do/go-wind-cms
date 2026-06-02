<script setup lang="ts">
const { t } = useI18n()
const searchQuery = ref('')

const handleSearch = () => {
  console.log('Searching for:', searchQuery.value)
}
</script>

<template>
  <div class="mx-2 hidden h-11 max-w-80 flex-1 items-center md:flex lg:max-w-80">
    <div class="relative w-full">
      <iconify-icon icon="lucide:search" width="16" height="16" class="absolute left-2 top-1/2 -translate-y-1/2 text-muted-foreground" />
      <UIInput
        class="h-full w-full pl-8"
        v-model="searchQuery"
        @keyup="handleSearch"
        :placeholder="t('navbar.top.search_placeholder')"
      />
    </div>
  </div>
</template>
