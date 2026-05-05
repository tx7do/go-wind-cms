<script setup lang="ts">
import {definePage} from 'unplugin-vue-router/runtime'
import {ref, onMounted, computed, watch} from 'vue'
import {useRoute, useRouter} from 'vue-router'
import {useMessage} from 'naive-ui'

import {useCategoryStore} from '@/stores'
import {$t} from '@/locales'
import type {
  contentservicev1_Category,
} from "@/api/generated/app/service/v1";
import {useLanguageChangeEffect} from '@/hooks/useLanguageChangeEffect';
import {CategoryList} from '@/components/CategoryList';

definePage({
  name: 'category-detail',
  meta: {
    title: 'Category Detail',
    level: 3,
  },
})

const route = useRoute()
const router = useRouter()
const categoryStore = useCategoryStore()
const message = useMessage()

const loading = ref(false)
const category = ref<contentservicev1_Category>(null)
const childCategories = ref<contentservicev1_Category[]>([])

const categoryId = computed(() => {
  const id = route.params.id
  return id ? parseInt(id as string) : null
})

// 获取父分类 ID
const parentCategoryId = computed(() => {
  if (!category.value?.parentId) return null
  return category.value.parentId
})

async function loadCategory() {
  if (!categoryId.value) return

  loading.value = true
  try {
    const categoryData = await categoryStore.getCategory(
      categoryId.value,
      'id,status,sort_order,icon,code,post_count,direct_post_count,parent_id,created_at,children,translations.id,translations.category_id,translations.name,translations.language_code,translations.description,translations.thumbnail,translations.cover_image',
    )
    category.value = categoryData
    // 提取子分类 - 使用新的数组引用确保响应式更新
    if (categoryData?.children && categoryData.children.length > 0) {
      // 创建新的数组引用，触发响应式更新
      childCategories.value = [...categoryData.children]
    } else {
      childCategories.value = []
    }
  } catch (error) {
    console.error('Load category failed:', error)
    message.error($t('page.post_detail.load_failed'))
  } finally {
    loading.value = false
  }
}

function handleViewChildCategory(id: number) {
  router.push(`/category/${id}`)
}

function handleBackToParent() {
  if (parentCategoryId.value) {
    router.push(`/category/${parentCategoryId.value}`)
  } else {
    // 如果没有父分类，返回分类列表页
    router.push('/category')
  }
}

onMounted(async () => {
  await loadCategory()
})

// 监听语言切换，自动重新加载分类数据
useLanguageChangeEffect(async () => {
  await loadCategory();
}, {
  immediate: false,      // 是否立即执行一次
  autoCleanup: true,    // 是否自动清理
});

// 监听路由参数变化（id 变化时重新加载）
watch(
  () => route.params.id,
  async (newId, oldId) => {
    if (newId !== oldId) {
      console.log('[Category Detail] Route param changed, reloading...')
      await loadCategory()
    }
  }
)
</script>

<template>
  <div class="category-detail-page">
    <!-- Hero Section -->
    <div class="hero-section">
      <div class="hero-content">
        <h1>{{ categoryStore.getCategoryName(category) }}</h1>
        <p v-if="categoryStore.getCategoryDescription(category)" class="category-description">{{
            categoryStore.getCategoryDescription(category)
          }}</p>
        <div class="category-stats">
          <div class="stat-item">
            <span class="i-carbon:document"/>
            <span>{{ category?.postCount || 0 }} {{ $t('page.posts.articles') }}</span>
          </div>
        </div>
      </div>
    </div>

    <!-- Posts Section -->
    <div class="page-container">
      <!-- Back to Parent Button -->
      <div class="back-button-container">
        <n-button
          size="small"
          @click="handleBackToParent"
        >
          <template #icon>
            <span class="i-carbon:arrow-left"/>
          </template>
          {{
            parentCategoryId ? $t('page.categories.back_to_parent') : $t('page.categories.back_to_list')
          }}
        </n-button>
      </div>

      <!-- Sub Categories List -->
      <div class="sub-categories-wrapper" v-if="childCategories.length > 0">
        <CategoryList
          :categories="childCategories"
          :loading="false"
          :show-skeleton="false"
          @category-click="handleViewChildCategory"
        />
      </div>

      <!-- Posts List with Pagination -->
      <PostListWithPagination
        :key="categoryId"
        :initial-page-size="10"
        :page-sizes="[10, 20, 30, 40]"
        :category-id="categoryId"
        from="category"
      />
    </div>
  </div>
</template>

<style scoped lang="less">
.category-detail-page {
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

    .category-description {
      font-size: 20px;
      color: rgba(255, 255, 255, 0.95);
      font-weight: 500;
      line-height: 1.5;
      max-width: 700px;
      margin: 0 auto 24px;
      text-shadow: 0 2px 10px rgba(0, 0, 0, 0.25);
      animation: fadeInUp 0.8s ease-out 0.2s both;
    }

    .category-stats {
      display: flex;
      justify-content: center;
      gap: 32px;
      flex-wrap: wrap;

      .stat-item {
        display: flex;
        align-items: center;
        gap: 8px;
        font-size: 16px;
        font-weight: 600;
        color: rgba(255, 255, 255, 0.9);
        text-shadow: 0 1px 4px rgba(0, 0, 0, 0.2);
        animation: fadeInUp 0.8s ease-out 0.3s both;

        span[class^="i-"] {
          font-size: 20px;
        }
      }
    }
  }
}

// Page Container
.page-container {
  max-width: 1200px;
  margin: 0 auto;
  padding: 0 32px 80px;
}

// Sub Categories Section
.sub-categories-section {
  margin-bottom: 48px;
  padding: 32px;
  background: var(--color-surface);
  border-radius: 16px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.05);

  .section-header {
    display: flex;
    align-items: center;
    gap: 12px;
    margin-bottom: 24px;

    h2 {
      font-size: 24px;
      font-weight: 700;
      margin: 0;
      color: var(--color-text-primary);
      display: flex;
      align-items: center;
      gap: 8px;

      span[class^="i-"] {
        font-size: 28px;
        color: var(--color-brand);
      }
    }
  }
}

// Sub Categories Wrapper
.sub-categories-wrapper {
  margin-bottom: 48px;
}

// Filter Controls
.filter-controls {
  display: flex;
  justify-content: flex-start;
  margin-bottom: 24px;
  padding-bottom: 24px;
  border-bottom: 1px solid var(--color-border);

  :deep(.n-button-group) {
    display: flex;
    gap: 12px; // 按钮间间距

    .n-button {
      min-width: 160px;
      font-weight: 500;
      border-radius: 8px;
      padding: 12px 20px;
      transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);

      span[class^="i-"] {
        font-size: 16px;
      }

      // 非激活状态按钮样式优化
      &:not(.n-button--primary) {
        background: var(--color-surface);
        border-color: var(--color-border);
        color: var(--color-text-secondary);

        &:hover {
          border-color: var(--color-brand);
          color: var(--color-brand);
          background: rgba(102, 126, 234, 0.05);
        }
      }
    }
  }
}

// Back Button
.back-button-container {
  margin-bottom: 24px;

  :deep(.n-button) {
    border-radius: 12px;
    padding: 12px 20px;
    font-weight: 600;
    font-size: 14px;
    transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
    background: var(--color-surface);
    border: 1.5px solid var(--color-border);
    color: var(--color-text-primary);
    box-shadow: 0 2px 8px rgba(0, 0, 0, 0.06);

    span[class^="i-"] {
      font-size: 16px;
      margin-right: 8px;
      transition: transform 0.3s cubic-bezier(0.4, 0, 0.2, 1);
    }

    &:hover {
      border-color: var(--color-brand);
      color: var(--color-brand);
      background: rgba(102, 126, 234, 0.08);
      transform: translateX(-4px);
      box-shadow: 0 4px 12px rgba(102, 126, 234, 0.15);

      span[class^="i-"] {
        transform: translateX(-4px);
      }
    }

    &:active {
      transform: translateX(-2px);
      box-shadow: 0 2px 6px rgba(102, 126, 234, 0.1);
    }
  }
}

// Responsive Design
@media (max-width: 1024px) {
  .sub-categories-section {
    padding: 24px;
    margin-bottom: 40px;

    .section-header {
      h2 {
        font-size: 22px;
      }
    }
  }

  .sub-categories-wrapper {
    margin-bottom: 40px;
  }

  .filter-controls {
    :deep(.n-button-group) {
      .n-button {
        min-width: 140px;
        font-size: 13px;
      }
    }
  }

  .posts-grid {
    grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
    gap: 24px;
  }
}

@media (max-width: 768px) {
  .back-button-container {
    margin-bottom: 20px;

    :deep(.n-button) {
      width: 100%;
      padding: 12px 16px;
      font-size: 14px;

      span[class^="i-"] {
        margin-right: 8px;
      }
    }
  }

  .sub-categories-section {
    padding: 20px;
    margin-bottom: 32px;

    .section-header {
      h2 {
        font-size: 20px;
      }
    }
  }

  .sub-categories-wrapper {
    margin-bottom: 32px;
  }

  .filter-controls {
    flex-direction: column;
    gap: 12px;

    :deep(.n-button-group) {
      width: 100%;
      flex-direction: column;
      gap: 12px;

      .n-button {
        width: 100%;
        padding: 14px 20px;
        border-radius: 10px;
        font-size: 14px;
        font-weight: 600;
        transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);

        // 激活状态 - 渐变背景和阴影
        &.n-button--primary {
          background: linear-gradient(135deg, var(--color-brand) 0%, rgba(102, 126, 234, 0.9) 100%);
          box-shadow: 0 4px 12px rgba(102, 126, 234, 0.3);
          border: none;

          &:hover {
            box-shadow: 0 6px 16px rgba(102, 126, 234, 0.4);
            transform: translateY(-1px);
          }

          &:active {
            transform: translateY(0);
            box-shadow: 0 2px 8px rgba(102, 126, 234, 0.25);
          }
        }

        // 非激活状态 - 增强视觉效果
        &:not(.n-button--primary) {
          background: var(--color-surface);
          border: 1.5px solid var(--color-border);
          color: var(--color-text-primary);
          box-shadow: 0 2px 4px rgba(0, 0, 0, 0.05);

          &:hover {
            border-color: var(--color-brand);
            color: var(--color-brand);
            background: rgba(102, 126, 234, 0.08);
            box-shadow: 0 4px 8px rgba(102, 126, 234, 0.15);
            transform: translateY(-2px);
          }

          &:active {
            transform: translateY(0);
            box-shadow: 0 1px 2px rgba(0, 0, 0, 0.05);
          }
        }
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
    &:hover {
      transform: translateY(-6px);
    }

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
        margin-bottom: 14px;
      }

      .post-meta {
        gap: 14px;
        padding-top: 14px;
        font-size: 12px;

        .meta-item {
          gap: 5px;

          span[class^="i-"] {
            font-size: 15px;
          }
        }
      }
    }
  }
}

@media (max-width: 640px) {
  .back-button-container {
    margin-bottom: 18px;

    :deep(.n-button) {
      padding: 11px 14px;
      font-size: 13px;

      span[class^="i-"] {
        font-size: 15px;
        margin-right: 6px;
      }
    }
  }

  .sub-categories-section {
    padding: 16px;
    margin-bottom: 28px;

    .section-header {
      h2 {
        font-size: 18px;
      }
    }
  }

  .sub-categories-wrapper {
    margin-bottom: 28px;
  }

  .filter-controls {
    :deep(.n-button-group) {
      .n-button {
        width: 100%;
        padding: 12px 16px;
        font-size: 13px;

        &.n-button--primary {
          box-shadow: 0 3px 10px rgba(102, 126, 234, 0.25);
        }
      }
    }
  }

  .results-info {
    font-size: 13px;
    margin-bottom: 18px;
  }

  .posts-grid {
    grid-template-columns: 1fr;
    gap: 18px;
    margin-bottom: 36px;
  }

  .post-card {
    border-radius: 14px;

    &:hover {
      transform: translateY(-5px);
    }

    .post-image {
      height: 180px;
    }

    .post-content {
      padding: 18px;

      .post-title {
        font-size: 16px;
        min-height: auto;
        margin-bottom: 10px;
        -webkit-line-clamp: 2;
      }

      .post-summary {
        font-size: 13px;
        margin-bottom: 12px;
        -webkit-line-clamp: 2;
      }

      .post-meta {
        gap: 12px;
        padding-top: 12px;
        font-size: 11px;

        .meta-item {
          span[class^="i-"] {
            font-size: 14px;
          }
        }
      }
    }
  }
}

@media (max-width: 480px) {
  .back-button-container {
    margin-bottom: 16px;

    :deep(.n-button) {
      padding: 12px 16px;
      font-size: 13px;
      border-radius: 12px;
      box-shadow: 0 2px 6px rgba(0, 0, 0, 0.05);

      span[class^="i-"] {
        font-size: 15px;
        margin-right: 8px;
      }
    }
  }

  .sub-categories-section {
    padding: 14px;
    margin-bottom: 24px;

    .section-header {
      h2 {
        font-size: 17px;
      }
    }
  }

  .sub-categories-wrapper {
    margin-bottom: 24px;
  }

  .filter-controls {
    :deep(.n-button-group) {
      .n-button {
        width: 100%;
        padding: 11px 14px;
        font-size: 12px;
        border-radius: 8px;

        &.n-button--primary {
          box-shadow: 0 2px 8px rgba(102, 126, 234, 0.2);
        }
      }
    }
  }

  .results-info {
    font-size: 12px;
    margin-bottom: 16px;
  }

  .posts-grid {
    grid-template-columns: 1fr;
    gap: 14px;
    margin-bottom: 32px;
  }

  .post-card {
    border-radius: 12px;

    &:hover {
      transform: translateY(-4px);
      box-shadow: 0 10px 28px rgba(0, 0, 0, 0.12);
    }

    .post-image {
      height: 160px;

      img {
        &:hover {
          transform: scale(1.05);
        }
      }
    }

    .post-content {
      padding: 16px;

      .post-title {
        font-size: 15px;
        margin-bottom: 8px;
        min-height: auto;
        line-height: 1.4;
      }

      .post-summary {
        font-size: 12px;
        margin-bottom: 10px;
        line-height: 1.6;
        -webkit-line-clamp: 2;
      }

      .post-meta {
        gap: 10px;
        padding-top: 10px;
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
        font-size: 40px;
        letter-spacing: -0.8px;
        margin-bottom: 14px;
      }

      .category-description {
        font-size: 18px;
        margin-bottom: 20px;
      }

      .category-stats {
        gap: 24px;

        .stat-item {
          font-size: 15px;

          span[class^="i-"] {
            font-size: 18px;
          }
        }
      }
    }
  }

  .page-container {
    padding: 0 24px 60px;
  }

  // Sub Categories Responsive
  .sub-categories-section {
    padding: 24px;
    margin-bottom: 40px;

    .section-header {
      h2 {
        font-size: 22px;
      }
    }

    .sub-categories-grid {
      grid-template-columns: repeat(auto-fill, minmax(260px, 1fr));
      gap: 16px;
    }
  }

  .filter-controls {
    :deep(.n-button-group) {
      .n-button {
        min-width: 140px;
        font-size: 13px;
      }
    }
  }

  .posts-grid {
    grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
    gap: 24px;
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
        font-size: 34px;
        margin-bottom: 12px;
        line-height: 1.2;
        letter-spacing: -0.6px;
      }

      .category-description {
        font-size: 17px;
        margin-bottom: 18px;
        line-height: 1.6;
      }

      .category-stats {
        gap: 20px;

        .stat-item {
          font-size: 14px;

          span[class^="i-"] {
            font-size: 17px;
          }
        }
      }
    }
  }

  .page-container {
    padding: 0 20px 50px;
  }

  // Sub Categories Responsive
  .sub-categories-section {
    padding: 20px;
    margin-bottom: 32px;

    .section-header {
      h2 {
        font-size: 20px;
      }
    }

    .sub-categories-grid {
      grid-template-columns: repeat(auto-fill, minmax(240px, 1fr));
      gap: 14px;
    }
  }

  .filter-controls {
    flex-direction: column;
    gap: 12px;

    :deep(.n-button-group) {
      width: 100%;

      .n-button {
        min-width: 100%;
        font-size: 13px;
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
    &:hover {
      transform: translateY(-6px);
    }

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
        margin-bottom: 14px;
      }

      .post-meta {
        gap: 14px;
        padding-top: 14px;
        font-size: 12px;

        .meta-item {
          gap: 5px;

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
    margin-bottom: 32px;

    &::after {
      background-size: 25px 25px;
    }

    .hero-content {
      h1 {
        font-size: 28px;
        margin-bottom: 10px;
        letter-spacing: -0.4px;
      }

      .category-description {
        font-size: 16px;
        margin-bottom: 16px;
      }

      .category-stats {
        gap: 18px;

        .stat-item {
          font-size: 13px;

          span[class^="i-"] {
            font-size: 16px;
          }
        }
      }
    }
  }

  .page-container {
    padding: 0 16px 40px;
  }

  .results-info {
    font-size: 13px;
    margin-bottom: 18px;
  }

  .posts-grid {
    grid-template-columns: 1fr;
    gap: 18px;
    margin-bottom: 36px;
  }

  .post-card {
    border-radius: 14px;

    &:hover {
      transform: translateY(-5px);
    }

    .post-image {
      height: 180px;
    }

    .post-content {
      padding: 18px;

      .post-title {
        font-size: 16px;
        min-height: auto;
        margin-bottom: 10px;
        -webkit-line-clamp: 2;
      }

      .post-summary {
        font-size: 13px;
        margin-bottom: 12px;
        -webkit-line-clamp: 2;
      }

      .post-meta {
        gap: 12px;
        padding-top: 12px;
        font-size: 11px;

        .meta-item {
          span[class^="i-"] {
            font-size: 14px;
          }
        }
      }
    }
  }

  .pagination-wrapper {
    padding: 28px 0;

    :deep(.n-pagination) {
      .n-pagination-item {
        min-width: 34px;
        height: 34px;
        font-size: 13px;
      }

      .n-pagination-prefix,
      .n-pagination-suffix {
        font-size: 13px;
      }
    }
  }
}

@media (max-width: 480px) {
  .hero-section {
    padding: 1.75rem 1rem 1.5rem;
    min-height: 220px;
    margin-bottom: 28px;

    &::before {
      animation: gradientShift 20s ease-in-out infinite;
    }

    &::after {
      background-size: 20px 20px;
      animation: gridMove 30s linear infinite;
    }

    .hero-content {
      padding: 0;

      h1 {
        font-size: 24px;
        margin-bottom: 8px;
        letter-spacing: -0.3px;
        line-height: 1.3;

        // 移动端简化发光效果
        text-shadow: 0 0 30px rgba(255, 255, 255, 0.8),
        0 0 60px rgba(var(--color-primary-purple-rgb), 0.6),
        0 4px 16px rgba(0, 0, 0, 0.4);

        filter: drop-shadow(0 0 20px rgba(var(--color-primary-purple-rgb), 0.5)) drop-shadow(0 3px 10px rgba(0, 0, 0, 0.3));
      }

      .category-description {
        font-size: 14px;
        margin-bottom: 14px;
        line-height: 1.5;
        max-width: 100%;
      }

      .category-stats {
        gap: 16px;
        justify-content: center;

        .stat-item {
          font-size: 12px;
          gap: 6px;

          span[class^="i-"] {
            font-size: 15px;
          }
        }
      }
    }
  }

  .page-container {
    padding: 0 12px 32px;
  }

  .results-info {
    font-size: 12px;
    margin-bottom: 16px;
  }

  .posts-grid {
    grid-template-columns: 1fr;
    gap: 14px;
    margin-bottom: 32px;
  }

  .post-card {
    border-radius: 12px;

    &:hover {
      transform: translateY(-4px);
      box-shadow: 0 10px 28px rgba(0, 0, 0, 0.12);
    }

    .post-image {
      height: 160px;

      img {
        &:hover {
          transform: scale(1.05);
        }
      }
    }

    .post-content {
      padding: 16px;

      .post-title {
        font-size: 15px;
        margin-bottom: 8px;
        min-height: auto;
        line-height: 1.4;
      }

      .post-summary {
        font-size: 12px;
        margin-bottom: 10px;
        line-height: 1.6;
        -webkit-line-clamp: 2;
      }

      .post-meta {
        gap: 10px;
        padding-top: 10px;
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
    padding: 24px 0;

    :deep(.n-pagination) {
      flex-wrap: wrap;
      justify-content: center;
      gap: 6px;

      .n-pagination-item {
        min-width: 32px;
        height: 32px;
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
    margin: 50px 0 !important;

    .n-empty__icon {
      span[class^="i-"] {
        font-size: 48px !important;
      }
    }

    .n-empty__description {
      font-size: 13px;
    }
  }

  // Spin 加载优化
  :deep(.n-spin-container) {
    min-height: 300px;
  }
}

// ...existing code...
// Hero 相关的动画关键帧已在 hero.less 中定义

// Dark Mode Enhancements
// Hero 的暗黑模式样式已在 hero.less 中定义
html.dark {

  .post-card {
    &:hover {
      box-shadow: 0 12px 32px rgba(0, 0, 0, 0.4);
    }
  }
}

</style>
