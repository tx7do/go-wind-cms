<template>
  <section class="features">
    <div class="section-header">
      <h2 class="section-title">
        <XIcon name="carbon:rocket" :size="28"
               style="color: #6366f1; margin-right: 8px;"/>
        {{ title }}
      </h2>
    </div>
    <div class="features-grid">
      <div v-for="feature in features" :key="feature.title" class="feature-card">
        <div class="feature-icon">
          <XIcon :name="feature.icon" :size="48" />
        </div>
        <h3>{{ feature.title }}</h3>
        <p>{{ feature.description }}</p>
      </div>
    </div>
  </section>
</template>

<script setup lang="ts">
import { XIcon } from '@/plugins/xicon';
import { defineProps } from 'vue';

defineProps<{
  title: string;
  features: Array<{
    icon: string;
    title: string;
    description: string;
  }>;
}>();
</script>

<style scoped lang="less">
.features {
  padding: 4rem 2rem;
  background: var(--color-bg);

  .section-header {
    text-align: center;
    display: block;
    margin-bottom: 3rem;
    max-width: 100%;
    padding: 0;

    .section-title {
      font-size: 2.25rem;
      font-weight: 800;
      letter-spacing: -1px;
      margin: 0 0 2rem 0;
      line-height: 1.1;
      background: linear-gradient(135deg, #6366f1 0%, #8b5cf6 50%, #a855f7 100%);
      -webkit-background-clip: text;
      -webkit-text-fill-color: transparent;
      background-clip: text;
      text-shadow: 0 2px 8px rgba(99,102,241,0.08), 0 1px 2px rgba(0,0,0,0.08);
    }
  }

  .features-grid {
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
    gap: 2rem;
    max-width: 1200px;
    margin: 0 auto;
  }

  .feature-card {
    padding: 2.5rem 2rem;
    border-radius: 12px;
    border: 1px solid var(--color-card-border, var(--color-border));
    background: var(--color-surface);
    transition: all 0.4s cubic-bezier(0.4, 0, 0.2, 1);
    text-align: center;
    position: relative;
    overflow: hidden;
    box-shadow: 0 2px 8px var(--color-card-shadow, rgba(0, 0, 0, 0.06));

    &::before {
      content: '';
      position: absolute;
      top: -100%;
      left: 50%;
      width: 100%;
      height: 100%;
      background: linear-gradient(135deg, rgba(102, 126, 234, 0.1) 0%, transparent 100%);
      transform: translateX(-50%);
      transition: top 0.5s cubic-bezier(0.4, 0, 0.2, 1);
    }

    &:hover {
      transform: translateY(-8px);
      box-shadow: 0 12px 32px var(--color-hover-shadow, rgba(0, 0, 0, 0.12));
      border-color: var(--color-brand);

      &::before {
        top: 0;
      }

      .feature-icon {
        transform: scale(1.15) rotate(-10deg);
        color: #fff;
        background: linear-gradient(135deg, var(--color-brand) 0%, #764ba2 100%);
      }

      h3 {
        color: var(--color-brand);
      }
    }

    .feature-icon {
      font-size: 3rem;
      color: var(--color-brand);
      margin-bottom: 1.5rem;
      width: 80px;
      height: 80px;
      display: flex;
      align-items: center;
      justify-content: center;
      margin-left: auto;
      margin-right: auto;
      background: rgba(102, 126, 234, 0.12);
      border-radius: 12px;
      transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
    }

    h3 {
      font-size: 1.25rem;
      margin-bottom: 0.75rem;
      color: var(--color-text-primary);
      font-weight: 700;
      transition: color 0.3s;
    }

    p {
      color: var(--color-text-secondary);
      font-size: 0.95rem;
      line-height: 1.7;
    }
  }
}

html.dark .features {
  .section-title {
    background: linear-gradient(135deg, #a5b4fc 0%, #6366f1 100%);
    -webkit-background-clip: text;
    -webkit-text-fill-color: transparent;
    background-clip: text;
    text-shadow: 0 2px 8px #23263a44;
    color: #e0e7ef;
  }
}
</style>
