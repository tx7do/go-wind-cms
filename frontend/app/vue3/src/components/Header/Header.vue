<script setup lang="ts">
</script>

<template>
  <div class="fixed-top">
    <div class="stick-content header-inner">
      <TopSearchBar/>
      <TopNavbar/>
    </div>
  </div>
</template>

<style scoped lang="less">
.fixed-top {
  width: 100%;
  z-index: 1000;
  background: var(--header-bg);
  backdrop-filter: blur(8px);
  border-bottom: 1px solid var(--color-border);
  color: var(--color-text-primary);
}

.header-inner {
  padding: var(--space-4) 0 var(--space-2);
}
</style>
