<script setup lang="ts">
import {definePage} from 'unplugin-vue-router/runtime'
import {ref, onMounted} from 'vue'

import {useCategoryStore} from '@/stores/modules/app'
import {$t} from '@/locales'
import CategoryFilter from '@/components/CategoryFilter';
import PostListWithPagination from '@/components/PostListWithPagination';

definePage({
  name: 'post-list',
  meta: {
    title: 'Posts',
    level: 2,
  },
})

const categoryStore = useCategoryStore()

const postListRef = ref<InstanceType<typeof PostListWithPagination> | null>(null)
const categories = ref<any[]>([])
const selectedCategoryId = ref<number | null>(null)

async function loadCategories() {
  try {
    const res = await categoryStore.listCategory(
      undefined,
      {status: 'CATEGORY_STATUS_ACTIVE'}
    )
    categories.value = res.items || []
  } catch (error) {
    console.error('Load categories failed:', error)
  }
}

function handleCategoryChange(categoryId: number | null) {
  selectedCategoryId.value = categoryId
}

onMounted(async () => {
  await loadCategories()
})
</script>

<template>
  <div class="post-list-page">
    <!-- Hero Section -->
    <div class="hero-section">
      <div class="hero-content">
        <h1>{{ $t('page.posts.posts_list') }}</h1>
        <p>{{ $t('page.posts.explore_latest') }}</p>
      </div>
    </div>

    <div class="page-container">
      <CategoryFilter
        :categories="categories"
        :selected-category="selectedCategoryId"
        @category-change="handleCategoryChange"
      />
      <PostListWithPagination
        ref="postListRef"
        :initial-page-size="12"
        :page-sizes="[12, 24, 36, 48]"
        :category-id="selectedCategoryId"
      />
    </div>
  </div>
</template>

<style scoped lang="less">
.post-list-page {
  min-height: 100vh;
  background: var(--color-bg);
}

// Hero Section
.hero-section {
  background: linear-gradient(135deg, #6366f1 0%, #8b5cf6 50%, #a855f7 100%);
  padding: 3rem 2rem;
  margin-bottom: 40px;
  position: relative;
  overflow: hidden;
  min-height: 350px;
  height: auto;
  display: flex;
  align-items: center;
  justify-content: center;
  text-align: center;
  color: white;

  // 背景装饰：渐变和网格
  &::before {
    content: '';
    position: absolute;
    top: 0;
    left: 0;
    right: 0;
    bottom: 0;
    background: radial-gradient(ellipse at 20% 50%, rgba(100, 200, 255, 0.15) 0%, transparent 50%),
    radial-gradient(ellipse at 80% 80%, rgba(200, 100, 255, 0.12) 0%, transparent 50%),
    linear-gradient(135deg, #6366f1 0%, #8b5cf6 50%, #a855f7 100%);
    animation: gradientShift 15s ease-in-out infinite;
    z-index: 0;
  }

  // 网格背景
  &::after {
    content: '';
    position: absolute;
    top: 0;
    left: 0;
    right: 0;
    bottom: 0;
    background-image: linear-gradient(rgba(255, 255, 255, 0.05) 1px, transparent 1px),
    linear-gradient(90deg, rgba(255, 255, 255, 0.05) 1px, transparent 1px);
    background-size: 50px 50px;
    animation: gridMove 20s linear infinite;
    opacity: 0.6;
    z-index: 1;
  }

  .hero-content {
    max-width: 1200px;
    margin: 0 auto;
    text-align: center;
    position: relative;
    z-index: 2;
    animation: fadeInUp 0.8s ease-out;

    h1 {
      font-size: 48px;
      font-weight: 800;
      margin: 0 0 16px 0;
      color: #ffffff;
      letter-spacing: -1px;
      line-height: 1.1;

      // 添加发光效果
      background: linear-gradient(135deg,
      #ffffff 0%,
      #f0f0ff 25%,
      #e0e0ff 50%,
      #f0f0ff 75%,
      #ffffff 100%);
      background-size: 200% auto;
      -webkit-background-clip: text;
      -webkit-text-fill-color: transparent;
      background-clip: text;

      text-shadow: 0 0 40px rgba(255, 255, 255, 0.9),
      0 0 80px rgba(var(--color-primary-purple-rgb), 0.7),
      0 0 120px rgba(99, 102, 241, 0.5),
      0 6px 24px rgba(0, 0, 0, 0.5),
      0 3px 12px rgba(0, 0, 0, 0.4);

      filter: drop-shadow(0 0 30px rgba(var(--color-primary-purple-rgb), 0.6)) drop-shadow(0 0 60px rgba(99, 102, 241, 0.4)) drop-shadow(0 5px 15px rgba(0, 0, 0, 0.3));

      animation-name: slideDown, glowPulseTitle, gradientShine;
      animation-duration: 0.8s, 3s, 6s;
      animation-timing-function: cubic-bezier(0.34, 1.56, 0.64, 1), ease-in-out, linear;
      animation-iteration-count: 1, infinite, infinite;
      animation-delay: 0s, 0.5s, 0s;
    }

    p {
      font-size: 20px;
      color: rgba(255, 255, 255, 0.95);
      margin: 0;
      font-weight: 500;
      text-shadow: 0 2px 10px rgba(0, 0, 0, 0.25);
      animation: fadeInUp 0.8s ease-out 0.2s both;
    }
  }
}

// Page Container
.page-container {
  max-width: 1400px;
  margin: 0 auto;
  padding: 0 32px 80px;
}

// Category Section
.category-section {
  margin-bottom: 48px;

  .category-tabs {
    display: flex;
    flex-wrap: wrap;
    gap: 12px;
    padding: 24px;
    background: var(--color-surface);
    border-radius: 16px;
    box-shadow: 0 2px 8px rgba(0, 0, 0, 0.04);

    :deep(.n-button) {
      border-radius: 50px;
      padding: 0 24px;
      font-weight: 500;
      transition: all 0.3s;

      &.n-button--primary-type:not(.n-button--ghost) {
        background: linear-gradient(135deg, var(--color-brand) 0%, #764ba2 100%);
        box-shadow: 0 4px 12px rgba(102, 126, 234, 0.3);
      }

      &.n-button--ghost {
        &:hover {
          color: var(--color-brand);
          border-color: var(--color-brand);
          background: var(--post-accent-bg-hover);
        }
      }
    }
  }
}

// Results Info
.results-info {
  margin-bottom: 24px;
  font-size: 15px;
  color: var(--color-text-secondary);
  font-weight: 500;

  strong {
    color: var(--color-brand);
    font-weight: 700;
  }
}

// Posts Grid
.posts-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(340px, 1fr));
  gap: 32px;
  margin-bottom: 48px;
}

.post-card {
  background: var(--color-surface);
  border-radius: 16px;
  overflow: hidden;
  cursor: pointer;
  transition: all 0.4s cubic-bezier(0.4, 0, 0.2, 1);
  border: 1px solid var(--color-border);
  height: 100%;
  display: flex;
  flex-direction: column;

  &:hover {
    transform: translateY(-8px);
    box-shadow: 0 12px 32px rgba(0, 0, 0, 0.12);
    border-color: var(--color-brand);

    .post-image {
      .image-overlay {
        background: rgba(102, 126, 234, 0.15);
      }

      img {
        transform: scale(1.08);
      }
    }

    .post-title {
      color: var(--color-brand);
    }
  }

  .post-image {
    position: relative;
    width: 100%;
    height: 220px;
    overflow: hidden;
    background: var(--color-bg);

    img {
      width: 100%;
      height: 100%;
      object-fit: cover;
      transition: transform 0.5s cubic-bezier(0.4, 0, 0.2, 1);
    }

    .image-overlay {
      position: absolute;
      top: 0;
      left: 0;
      right: 0;
      bottom: 0;
      background: rgba(0, 0, 0, 0.05);
      transition: all 0.4s;
    }
  }

  .post-content {
    padding: 24px;
    flex: 1;
    display: flex;
    flex-direction: column;

    .post-title {
      font-size: 19px;
      font-weight: 700;
      margin: 0 0 12px 0;
      color: var(--color-text-primary);
      display: -webkit-box;
      -webkit-line-clamp: 2;
      -webkit-box-orient: vertical;
      overflow: hidden;
      line-height: 1.4;
      transition: color 0.3s;
      min-height: 53px;
    }

    .post-summary {
      color: var(--color-text-secondary);
      font-size: 15px;
      line-height: 1.7;
      margin: 0 0 16px 0;
      display: -webkit-box;
      -webkit-line-clamp: 3;
      -webkit-box-orient: vertical;
      overflow: hidden;
      flex: 1;
    }

    .post-meta {
      display: flex;
      flex-wrap: wrap;
      gap: 16px;
      padding-top: 16px;
      border-top: 1px solid var(--color-border);
      font-size: 13px;
      color: var(--color-text-secondary);

      .meta-item {
        display: flex;
        align-items: center;
        gap: 6px;
        font-weight: 500;

        span[class^="i-"] {
          font-size: 16px;
          opacity: 0.8;
        }
      }
    }
  }
}

// Pagination
.pagination-wrapper {
  display: flex;
  justify-content: center;
  padding: 40px 0;

  :deep(.n-pagination) {
    .n-pagination-item {
      border-radius: 8px;
      transition: all 0.3s;

      &.n-pagination-item--active {
        background: var(--color-brand);
        color: #fff;
      }

      &:not(.n-pagination-item--disabled):hover {
        background: var(--post-accent-bg-hover);
        color: var(--color-brand);
      }
    }
  }
}

// Responsive Design
@media (max-width: 1024px) {
  .hero-section {
    padding: 3rem 1.5rem 2.5rem;
    min-height: 300px;

    .hero-content {
      h1 {
        font-size: 3rem;
        letter-spacing: -1.2px;
      }

      p {
        font-size: 1.25rem;
      }
    }
  }

  .page-container {
    padding: 0 24px 60px;
  }

  .category-section {
    margin-bottom: 40px;

    .category-tabs {
      padding: 20px;
      gap: 10px;

      :deep(.n-button) {
        padding: 0 20px;
        font-size: 15px;
      }
    }
  }

  .posts-grid {
    grid-template-columns: repeat(auto-fill, minmax(320px, 1fr));
    gap: 28px;
  }

  .post-card {
    .post-image {
      height: 200px;
    }

    .post-content {
      padding: 22px;
    }
  }
}

@media (max-width: 768px) {
  .hero-section {
    padding: 2.5rem 1.5rem 2rem;
    min-height: 280px;

    &::after {
      background-size: 30px 30px;
    }

    .hero-content {
      h1 {
        font-size: 2.2rem;
        margin-bottom: 1rem;
        letter-spacing: -0.8px;
      }

      p {
        font-size: 1.1rem;
      }
    }
  }

  .page-container {
    padding: 0 20px 50px;
  }

  .category-section {
    margin-bottom: 32px;

    .category-tabs {
      padding: 16px;
      gap: 8px;
      overflow-x: auto;
      overflow-y: hidden;
      -webkit-overflow-scrolling: touch;

      // 显示滚动条提示
      scrollbar-width: thin;
      scrollbar-color: rgba(102, 126, 234, 0.2) transparent;

      &::-webkit-scrollbar {
        height: 3px;
      }

      &::-webkit-scrollbar-track {
        background: transparent;
      }

      &::-webkit-scrollbar-thumb {
        background: rgba(102, 126, 234, 0.2);
        border-radius: 10px;
      }

      :deep(.n-button) {
        flex-shrink: 0;
        white-space: nowrap;
        font-size: 14px;
        padding: 0 18px;
      }
    }
  }

  .results-info {
    font-size: 14px;
    margin-bottom: 20px;
  }

  .posts-grid {
    grid-template-columns: repeat(auto-fill, minmax(280px, 1fr));
    gap: 20px;
    margin-bottom: 40px;
  }

  .post-card {
    .post-image {
      height: 200px;
    }

    .post-content {
      padding: 20px;

      .post-title {
        font-size: 17px;
        min-height: 48px;
      }

      .post-summary {
        font-size: 14px;
      }

      .post-meta {
        gap: 14px;
        font-size: 12px;

        .meta-item {
          span[class^="i-"] {
            font-size: 15px;
          }
        }
      }
    }
  }

  .pagination-wrapper {
    padding: 32px 0;
  }
}

@media (max-width: 640px) {
  .hero-section {
    padding: 2rem 1rem;
    min-height: 240px;

    &::after {
      background-size: 25px 25px;
    }

    .hero-content {
      h1 {
        font-size: 1.75rem;
        letter-spacing: -0.5px;
      }

      p {
        font-size: 1rem;
      }
    }
  }

  .page-container {
    padding: 0 16px 40px;
  }

  .category-section {
    margin-bottom: 28px;

    .category-tabs {
      padding: 12px;
      gap: 8px;
      flex-wrap: nowrap;
      justify-content: flex-start;

      :deep(.n-button) {
        font-size: 13px;
        padding: 0 16px;
        height: 36px;

        .n-button__icon {
          font-size: 16px;
        }
      }
    }
  }

  .results-info {
    font-size: 13px;
    margin-bottom: 16px;
  }

  .posts-grid {
    grid-template-columns: 1fr;
    gap: 16px;
    margin-bottom: 32px;
  }

  .post-card {
    &:hover {
      transform: translateY(-4px);
    }

    .post-image {
      height: 180px;
    }

    .post-content {
      padding: 16px;

      .post-title {
        font-size: 16px;
        min-height: auto;
        margin-bottom: 10px;
      }

      .post-summary {
        font-size: 13px;
        line-height: 1.6;
        -webkit-line-clamp: 2;
        margin-bottom: 14px;
      }

      .post-meta {
        gap: 12px;
        padding-top: 14px;
        font-size: 11px;

        .meta-item {
          gap: 5px;

          span[class^="i-"] {
            font-size: 14px;
          }
        }
      }
    }
  }

  .pagination-wrapper {
    padding: 24px 0;

    :deep(.n-pagination) {
      .n-pagination-item {
        min-width: 32px;
        height: 32px;
        font-size: 13px;
      }

      .n-pagination-prefix,
      .n-pagination-suffix {
        font-size: 13px;
      }

      .n-pagination-item__button {
        padding: 0 8px;
      }
    }
  }
}

@media (max-width: 480px) {
  .hero-section {
    padding: 1.75rem 1rem 1.5rem;
    min-height: 220px;

    &::after {
      background-size: 20px 20px;
    }

    .hero-content {
      padding: 0;

      h1 {
        font-size: 1.6rem;
        margin-bottom: 0.75rem;
        letter-spacing: -0.3px;
      }

      p {
        font-size: 0.95rem;
      }
    }
  }

  .page-container {
    padding: 0 12px 32px;
  }

  .category-section {
    margin-bottom: 24px;

    .category-tabs {
      padding: 10px;
      gap: 6px;

      :deep(.n-button) {
        font-size: 12px;
        padding: 0 14px;
        height: 34px;

        .n-button__icon {
          font-size: 14px;
          margin-right: 4px;
        }
      }
    }
  }

  .results-info {
    font-size: 12px;
    margin-bottom: 14px;
  }

  .posts-grid {
    gap: 14px;
    margin-bottom: 28px;
  }

  .post-card {
    border-radius: 12px;

    &:hover {
      transform: translateY(-3px);
      box-shadow: 0 8px 24px rgba(0, 0, 0, 0.1);
    }

    .post-image {
      height: 160px;
    }

    .post-content {
      padding: 14px;

      .post-title {
        font-size: 15px;
        margin-bottom: 8px;
        -webkit-line-clamp: 2;
      }

      .post-summary {
        font-size: 12px;
        line-height: 1.5;
        margin-bottom: 12px;
      }

      .post-meta {
        gap: 10px;
        padding-top: 12px;
        font-size: 10px;
        flex-wrap: wrap;

        .meta-item {
          gap: 4px;

          span[class^="i-"] {
            font-size: 13px;
          }
        }
      }
    }
  }

  .pagination-wrapper {
    padding: 20px 0;

    :deep(.n-pagination) {
      flex-wrap: wrap;
      justify-content: center;
      gap: 6px;

      .n-pagination-item {
        min-width: 30px;
        height: 30px;
        font-size: 12px;
      }

      .n-pagination-prefix,
      .n-pagination-suffix {
        font-size: 12px;
        margin: 0 4px;
      }

      .n-pagination-item__button {
        padding: 0 6px;
      }

      // 隐藏 size picker 在小屏幕
      .n-pagination-size-picker {
        display: none;
      }
    }
  }

  // Empty state 优化
  :deep(.n-empty) {
    margin: 60px 0 !important;

    .n-empty__icon {
      font-size: 48px !important;
    }

    .n-empty__description {
      font-size: 14px;
    }
  }
}

// 动画关键帧定义
@keyframes slideDown {
  from {
    opacity: 0;
    transform: translateY(-30px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

@keyframes fadeInUp {
  from {
    opacity: 0;
    transform: translateY(20px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

@keyframes gradientShift {
  0%, 100% {
    background: radial-gradient(ellipse at 20% 50%, rgba(100, 200, 255, 0.15) 0%, transparent 50%),
    radial-gradient(ellipse at 80% 80%, rgba(200, 100, 255, 0.12) 0%, transparent 50%),
    linear-gradient(135deg, #6366f1 0%, #8b5cf6 50%, #a855f7 100%);
  }
  50% {
    background: radial-gradient(ellipse at 80% 50%, rgba(100, 200, 255, 0.15) 0%, transparent 50%),
    radial-gradient(ellipse at 20% 80%, rgba(200, 100, 255, 0.12) 0%, transparent 50%),
    linear-gradient(135deg, #8b5cf6 0%, #a855f7 50%, #6366f1 100%);
  }
}

@keyframes gridMove {
  0% {
    transform: translate(0, 0);
  }
  100% {
    transform: translate(50px, 50px);
  }
}

@keyframes glowPulseTitle {
  0%, 100% {
    filter: drop-shadow(0 0 25px rgba(var(--color-primary-purple-rgb), 0.5)) drop-shadow(0 0 50px rgba(99, 102, 241, 0.3)) drop-shadow(0 5px 15px rgba(0, 0, 0, 0.3));
  }
  50% {
    filter: drop-shadow(0 0 40px rgba(var(--color-primary-purple-rgb), 0.8)) drop-shadow(0 0 80px rgba(99, 102, 241, 0.6)) drop-shadow(0 8px 20px rgba(0, 0, 0, 0.4));
  }
}

@keyframes gradientShine {
  0% {
    background-position: 0 center;
  }
  100% {
    background-position: 200% center;
  }
}

// Dark Mode Enhancements
html.dark {
  .hero-section {
    background: linear-gradient(135deg, #6366f1 0%, #8b5cf6 50%, #a855f7 100%);
    box-shadow: 0 20px 60px rgba(99, 102, 241, 0.3);

    &::before {
      opacity: 1;
    }

    .hero-content {
      h1 {
        background: linear-gradient(135deg,
        #ffffff 0%,
        #f5f5ff 25%,
        #ebebff 50%,
        #f5f5ff 75%,
        #ffffff 100%);
        background-size: 200% auto;
        -webkit-background-clip: text;
        -webkit-text-fill-color: transparent;
        background-clip: text;

        text-shadow: 0 0 50px rgba(255, 255, 255, 1),
        0 0 100px rgba(var(--color-primary-purple-rgb), 0.8),
        0 0 150px rgba(99, 102, 241, 0.6),
        0 8px 30px rgba(0, 0, 0, 0.6),
        0 4px 15px rgba(0, 0, 0, 0.4);

        filter: drop-shadow(0 0 35px rgba(var(--color-primary-purple-rgb), 0.7)) drop-shadow(0 0 70px rgba(99, 102, 241, 0.5)) drop-shadow(0 6px 20px rgba(0, 0, 0, 0.4));
      }
    }
  }

  .post-card {
    &:hover {
      box-shadow: 0 12px 32px rgba(0, 0, 0, 0.4);
    }
  }
}
</style>

