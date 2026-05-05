<script setup lang="ts">
import {definePage} from 'unplugin-vue-router/runtime'
import {$t} from '@/locales'

definePage({
  name: 'privacy',
  meta: {
    level: 2,
  },
})
</script>

<template>
  <div class="info-page">
    <!-- Hero Section -->
    <section class="hero">
      <div class="hero-content">
        <h1 class="hero-title">{{ $t('ui.button.privacy_agreement') }}</h1>
        <p class="hero-subtitle">{{ $t('page.legal.privacy.description') }}</p>
      </div>
    </section>

    <!-- Content Section -->
    <div class="info-card stick-content section-shell">
      <ul class="list">
        <li>{{ $t('page.legal.privacy.item_1') }}</li>
        <li>{{ $t('page.legal.privacy.item_2') }}</li>
        <li>{{ $t('page.legal.privacy.item_3') }}</li>
      </ul>
    </div>
  </div>
</template>

<style scoped lang="less">
.info-page {
  min-height: 100vh;
  background: var(--color-bg);
}

// Content Section
.info-card {
  max-width: 1200px;
  margin: 0 auto;
  padding: 40px 32px;
}

.list {
  margin: 0;
  padding-left: 18px;
  color: var(--color-text-primary);
  line-height: 1.8;
  font-size: 16px;

  li {
    margin-bottom: 12px;
  }
}

// Responsive Design
@media (max-width: 1024px) {
  .info-card {
    padding: 32px 24px;
  }
}

@media (max-width: 768px) {
  .info-card {
    padding: 24px 20px;
  }

  .list {
    font-size: 15px;
  }
}

@media (max-width: 480px) {
  .info-card {
    padding: 20px 16px;
  }

  .list {
    font-size: 14px;
    padding-left: 16px;
  }
}
</style>

<style lang="less">
// Hero Section - 非 scoped 样式，使用全局 hero.less 中的定义
// 只需要定义页面特定的 margin
.info-page .hero {
  margin-bottom: 40px;
}
</style>

