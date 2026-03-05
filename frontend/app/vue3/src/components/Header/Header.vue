<script setup lang="ts">
</script>
<template>
  <div class="fixed-top">
    <div class="header-inner">
      <TopSearchBar/>
      <TopNavbar/>
    </div>
  </div>
</template>

<style scoped lang="less">
.fixed-top {
  width: 100%;
  z-index: 1000;
  background: linear-gradient(
    to bottom,
    var(--header-bg) 0%,
    var(--header-bg) 90%,
    rgba(255, 255, 255, 0)
  );
  backdrop-filter: blur(12px);
  border-bottom: 1px solid var(--color-border);
  color: var(--color-text-primary);
  position: sticky;
  top: 0;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.04);
  height: 100%;
  display: flex;
  justify-content: center;

  html.dark & {
    box-shadow: 0 2px 8px rgba(0, 0, 0, 0.2);
  }
}

.header-inner {
  display: flex;
  flex-direction: column;
  gap: 0;
  padding: 0;
  height: 100%;
  width: min(var(--layout-max-width), calc(100% - 2 * var(--space-6)));
  max-width: var(--layout-max-width);

  // 每行之间添加分割线
  & > :not(:last-child) {
    border-bottom: 1px solid var(--color-border);
  }
}

@media (max-width: 767px) {
  .header-inner {
    width: calc(100% - 2 * var(--space-4));
  }
}
</style>
