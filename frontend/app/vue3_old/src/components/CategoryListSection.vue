<template>
  <section class="categories-section scroll-reveal">
    <!-- Loading Skeleton -->
    <div v-if="loading" class="categories-grid desktop-grid">
      <div v-for="i in skeletonCount ?? 8" :key="i" class="category-card">
        <n-skeleton height="140px"/>
      </div>
    </div>
    <!-- Loaded Content -->
    <div v-else>
      <div class="categories-grid desktop-grid">
        <div
          v-for="category in categories"
          :key="category.id"
          class="category-card scroll-reveal-item"
          :style="{ '--card-hue': (category.id % 6) * 60 }"
          @click="handleViewCategory(category.id)"
        >
          <div class="category-card-bg"></div>
          <div class="category-card-content">
            <div class="category-card-header">
              <div class="category-icon">
                <XIcon
                  :name="category.icon?.includes(':') ? category.icon : `carbon:${category.icon || 'folder'}`"
                  :size="48"
                />
              </div>
              <div class="category-info">
                <h3>{{ categoryStore.getCategoryName(category) }}</h3>
                <span class="post-count">
                  {{ $t('page.home.article_count', {count: category.postCount || 0}) }}
                </span>
              </div>
            </div>
            <div class="category-badge">
              <XIcon name="carbon:time" :size="14"/>
              <span>{{ getUpdatedDaysAgo(3) }}</span>
            </div>
          </div>
        </div>
      </div>
      <div v-if="showCarousel" class="categories-carousel mobile-carousel">
        <n-carousel
          :autoplay="true"
          :interval="5000"
          :show-arrow="false"
          :mouse-wheel="true"
          :touchable="true"
          :slides-per-view="1.5"
          :centered-slides="true"
          :scroll-eventListener="false"
          :show-dots="false"
          class="carousel-container"
        >
          <template v-for="category in categories" :key="category.id">
            <div class="carousel-item" @click="handleViewCategory(category.id)">
              <div class="category-card carousel-card"
                   :style="{ '--card-hue': (category.id % 6) * 60 }">
                <div class="category-card-bg"></div>
                <div class="category-card-content">
                  <div class="category-card-header">
                    <div class="category-icon">
                      <XIcon
                        :name="category.icon?.includes(':') ? category.icon : `carbon:${category.icon || 'folder'}`"
                        :size="48"
                      />
                    </div>
                    <div class="category-info">
                      <h3>{{ categoryStore.getCategoryName(category) }}</h3>
                      <span class="post-count">
                        {{ $t('page.home.article_count', {count: category.postCount || 0}) }}
                      </span>
                    </div>
                  </div>
                  <div class="category-badge">
                    <XIcon name="carbon:time" :size="14"/>
                    <span>{{ getUpdatedDaysAgo(3) }}</span>
                  </div>
                </div>
              </div>
            </div>
          </template>
        </n-carousel>
      </div>
      <slot name="extra"></slot>
    </div>
  </section>
</template>

<script setup lang="ts">
import {ref, onMounted, defineExpose, watch} from 'vue';
import {XIcon} from '@/plugins/xicon';
import {useCategoryStore} from '@/stores';
import {useLanguageChangeEffect} from '@/hooks/useLanguageChangeEffect';
import type {contentservicev1_Category} from '@/api/generated/app/service/v1';
import {$t} from "@/locales";
import router from "@/router/router";

const props = defineProps<{
  skeletonCount?: number;
  showCarousel?: boolean;
  pageSize?: number;
  filter?: Record<string, any>;
  orderBy?: string[];
  fieldMask?: string;
}>();

const categoryStore = useCategoryStore();
const categories = ref<contentservicev1_Category[]>([]);
const loading = ref(false);

function getUpdatedDaysAgo(days: number = 3) {
  return $t('page.home.updated_days_ago', {days})
}

function handleViewCategory(id: number) {
  router.push(`/category/${id}`)
}

async function loadCategories() {
  loading.value = true;
  try {
    const res = await categoryStore.listCategory(
      {page: 1, pageSize: props.pageSize ?? 8},
      props.filter ?? {status: 'CATEGORY_STATUS_ACTIVE'},
      props.fieldMask ?? 'id,status,sort_order,icon,code,post_count,direct_post_count,parent_id,created_at,translations.id,translations.category_id,translations.name,translations.language_code,translations.description',
      props.orderBy ?? ['-sortOrder', '-postCount']
    );
    categories.value = res.items || [];
  } catch (e) {
    categories.value = [];
  }
  loading.value = false;
}

function reload() {
  loadCategories();
}

defineExpose({reload});

onMounted(() => {
  loadCategories();
});

// 支持外部 filter/orderBy/pageSize 变化时自动刷新
watch(() => [props.filter, props.orderBy, props.pageSize, props.fieldMask], () => {
  loadCategories();
}, {deep: true});

// 监听语言切换自动刷新分类
useLanguageChangeEffect(() => {
  loadCategories();
}, {immediate: false, autoCleanup: true});
</script>

<style scoped lang="less">
// 网格布局样式
.categories-grid {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 1.5rem;
  max-width: 1200px;
  margin: 0 auto;
  padding: 0 2rem;
}

.categories-carousel {
  display: none;
}

.category-card {
  display: flex;
  flex-direction: column;
  padding: 2rem 1.75rem;
  background: #fff;
  border: 1.5px solid #e0e7ef;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.06);
  color: #3730a3;
  border-radius: 16px;
  cursor: pointer;
  position: relative;
  min-height: 220px;
  height: 100%;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  overflow: visible;

  .category-icon {
    font-size: 2.5rem;
    width: 80px;
    height: 80px;
    display: flex;
    align-items: center;
    justify-content: center;
    background: #e0e7ff;
    color: #6366f1;
    border-radius: 16px;
    margin-bottom: 1.25rem;
    box-shadow: 0 4px 12px rgba(99, 102, 241, 0.15);
    transition: all 0.4s cubic-bezier(0.4, 0, 0.2, 1);
  }

  .category-info {
    flex: 1;
    width: 100%;

    h3 {
      font-size: 1.15rem;
      font-weight: 800;
      margin: 0 0 0.5rem 0;
      color: #1e293b;
      letter-spacing: 0.3px;
      word-break: break-word;
      line-height: 1.4;
      transition: color 0.3s;
    }

    .post-count {
      font-size: 0.95rem;
      color: #6366f1;
      font-weight: 600;
      letter-spacing: 0.2px;
      display: inline-block;
      padding: 0.35rem 0.875rem;
      background: rgba(99, 102, 241, 0.08);
      border-radius: 8px;
      border: 1px solid rgba(99, 102, 241, 0.15);
      transition: all 0.3s;
      margin-bottom: 0.75rem;
    }
  }

  .category-badge {
    margin-top: auto;
    padding: 0.625rem 1rem;
    background: #e0e7ff;
    color: #6366f1;
    border-radius: 10px;
    font-size: 0.825rem;
    font-weight: 600;
    width: fit-content;
    white-space: nowrap;
    display: flex;
    align-items: center;
    gap: 0.5rem;
    box-shadow: 0 2px 6px rgba(99, 102, 241, 0.12);
    transition: all 0.3s;
    border: 1px solid rgba(99, 102, 241, 0.12);
  }

  &:hover {
    transform: translateY(-8px);
    background: #fff;
    border-color: #6366f1;
    box-shadow: 0 12px 28px rgba(99, 102, 241, 0.2);

    .category-icon {
      background: linear-gradient(135deg, #6366f1 0%, #8b5cf6 100%);
      color: #fff;
      transform: scale(1.12) rotate(-8deg);
      box-shadow: 0 8px 20px rgba(99, 102, 241, 0.35);
    }

    .category-info h3 {
      color: #4338ca;
    }

    .post-count {
      color: #4338ca;
      background: rgba(67, 56, 202, 0.1);
      border-color: rgba(67, 56, 202, 0.2);
    }

    .category-badge {
      background: linear-gradient(135deg, #6366f1 0%, #8b5cf6 100%);
      color: #fff;
      border-left-color: #4338ca;
      box-shadow: 0 4px 12px rgba(99, 102, 241, 0.3);
      border-color: rgba(99, 102, 241, 0.2);
    }
  }
}

// ========== 深色模式 ==========
html.dark {
  .category-card {
    background: #1e293b;
    border: 1px solid rgba(51, 65, 85, 0.6);
    box-shadow: 0 2px 8px rgba(0, 0, 0, 0.15);
    color: #e0e7ef;

    .category-icon {
      background: rgba(99, 102, 241, 0.15);
      color: #a5b4fc;
      box-shadow: 0 4px 12px rgba(0, 0, 0, 0.2);
    }

    .category-info h3 {
      color: #e2e8f0;
    }

    .post-count {
      color: #a5b4fc;
      background: rgba(165, 180, 252, 0.08);
      border: 1px solid rgba(165, 180, 252, 0.1);
    }

    .category-badge {
      background: #1e293b;
      color: #a5b4fc;
      box-shadow: 0 2px 6px rgba(0, 0, 0, 0.15);
      border: 1px solid rgba(51, 65, 85, 0.5);
    }

    &:hover {
      transform: translateY(-8px);
      background: #1e293b;
      border-color: rgba(99, 102, 241, 0.6);
      box-shadow: 0 8px 24px rgba(99, 102, 241, 0.2);

      .category-icon {
        background: linear-gradient(135deg, #6366f1 0%, #8b5cf6 100%);
        color: #fff;
        box-shadow: 0 6px 16px rgba(99, 102, 241, 0.35);
      }

      .category-info h3 {
        color: #f1f5f9;
      }

      .post-count {
        color: #c7d2fe;
        background: rgba(199, 210, 254, 0.12);
        border-color: rgba(199, 210, 254, 0.18);
      }

      .category-badge {
        background: linear-gradient(135deg, #6366f1 0%, #8b5cf6 100%);
        color: #fff;
        box-shadow: 0 4px 12px rgba(99, 102, 241, 0.3);
        border-color: rgba(99, 102, 241, 0.5);
      }
    }
  }

  // 移动端适配
  @media (max-width: 900px) {
    .categories-grid {
      display: none;
    }

    .categories-carousel {
      display: block;
      width: 100%;
      margin: 0 auto;
      padding: 0 0.5rem 2rem 0.5rem;

      .carousel-container {
        width: 100%;
        padding: 1.5rem 0;
      }

      .carousel-item {
        padding: 0.75rem 0.5rem;
        display: flex;
        justify-content: center;
        align-items: stretch;
      }

      .category-card.carousel-card {
        min-width: 280px;
        max-width: 320px;
        margin: 0 auto;
        height: auto;
        display: flex;
        flex-direction: column;
      }
    }
  }
}
</style>
