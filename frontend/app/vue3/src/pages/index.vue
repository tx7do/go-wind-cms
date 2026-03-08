<script setup lang="ts">
import {definePage} from 'unplugin-vue-router/runtime'
import {useRouter} from 'vue-router'
import {ref, onMounted, computed, onUnmounted} from 'vue'

import {usePostStore, useCategoryStore, useTagStore} from '@/stores/modules/app'
import {$t, i18n} from '@/locales'
import {XIcon} from "@/plugins/xicon";
import {formatDate} from "@/utils/date";
import {useLanguageChangeEffect} from '@/hooks/useLanguageChangeEffect';
import type {
  contentservicev1_Category,
  contentservicev1_Post,
  contentservicev1_Tag
} from "@/api/generated/app/service/v1";

definePage({
  name: 'home',
  meta: {
    level: 1,
  },
})

const router = useRouter()
const postStore = usePostStore()
const categoryStore = useCategoryStore()
const tagStore = useTagStore()

const latestPosts = ref<contentservicev1_Post[]>([])
const featuredPosts = ref<contentservicev1_Post[]>([])
const categories = ref<contentservicev1_Category[]>([])
const popularTags = ref<contentservicev1_Tag[]>([])
const loading = ref(false)

// 滚动动画相关
const observerRef = ref<IntersectionObserver | null>(null)

// 初始化滚动观察器
function initScrollObserver() {
  const options = {
    root: null,
    threshold: 0.1,
    rootMargin: '0px 0px -100px 0px'
  }

  observerRef.value = new IntersectionObserver((entries) => {
    entries.forEach((entry) => {
      if (entry.isIntersecting) {
        entry.target.classList.add('is-visible')
        // 一旦动画触发，停止观察该元素
        observerRef.value?.unobserve(entry.target)
      }
    })
  }, options)

  // 观察所有需要动画的元素
  const elements = document.querySelectorAll('.scroll-reveal-item')
  elements.forEach((el) => {
    observerRef.value?.observe(el)
  })
}

// 销毁观察器
function destroyScrollObserver() {
  if (observerRef.value) {
    observerRef.value.disconnect()
    observerRef.value = null
  }
}

// 加载最新文章
async function loadLatestPosts() {
  try {
    const res = await postStore.listPost(
      {page: 1, pageSize: 6},
      {status: 'POST_STATUS_PUBLISHED'}
    )
    latestPosts.value = res.items || []
  } catch (error) {
    console.error('Load latest posts failed:', error)
  }
}

// 加载推荐文章
async function loadFeaturedPosts() {
  try {
    const res = await postStore.listPost(
      {page: 1, pageSize: 3},
      {status: 'POST_STATUS_PUBLISHED', isFeatured: true}
    )
    featuredPosts.value = res.items || []
  } catch (error) {
    console.error('Load featured posts failed:', error)
  }
}

// 加载分类
async function loadCategories() {
  try {
    const res = await categoryStore.listCategory(
      {page: 1, pageSize: 8},
      {status: 'CATEGORY_STATUS_ACTIVE'}
    )
    categories.value = res.items || []
  } catch (error) {
    console.error('Load categories failed:', error)
  }
}

// 加载热门标签
async function loadPopularTags() {
  try {
    const res = await tagStore.listTag(
      {page: 1, pageSize: 6},
      {status: 'TAG_STATUS_ACTIVE', isFeatured: true}
    )
    // 将标签数据转换为显示格式
    popularTags.value = res.items?.map((tag: contentservicev1_Tag, index: number) => ({
      id: tag.id,
      name: tag.translations?.[0]?.name || tag.name || $t('page.home.tag_default'),
      color: tag.color || `hsl(${index * 60}, 100%, 50%)`,
      slug: tag.slug,
    })) || []
  } catch (error) {
    console.error('Load popular tags failed:', error)
  }
}

function getPostTitle(post: contentservicev1_Post) {
  const translation = postStore.getTranslation(post)
  return translation?.title || $t('page.home.untitled')
}

function getPostSummary(post: contentservicev1_Post) {
  const translation = postStore.getTranslation(post)
  return translation?.summary || ''
}

function getPostThumbnail(post: contentservicev1_Post) {
  const translation = postStore.getTranslation(post)
  return translation?.thumbnail || '/placeholder.jpg'
}

function getCategoryName(category: contentservicev1_Category) {
  const translation = categoryStore.getTranslation(category)
  return translation?.name || $t('page.home.category_default')
}

function handleViewPost(id: number) {
  router.push({
    path: `/post/${id}`,
    query: {
      from: 'home'
    }
  })
}

function handleViewCategory(id: number) {
  router.push(`/category/${id}`)
}

function getUpdatedDaysAgo(days: number = 3) {
  return $t('page.home.updated_days_ago', {days})
}

// Smooth scroll to categories section
function scrollToCategories() {
  const categoriesSection = document.querySelector('.categories-section')
  if (categoriesSection) {
    categoriesSection.scrollIntoView({
      behavior: 'smooth',
      block: 'start'
    })
  }
}

const features = computed(() => {
  // Track locale so this list recomputes immediately on language switch.
  void i18n.global.locale.value

  return [
    {
      icon: 'carbon:document',
      title: $t('page.home.flexible_content_management'),
      description: $t('page.home.content_management_desc'),
    },
    {
      icon: 'carbon:cloud',
      title: $t('page.home.multi_tenant_architecture'),
      description: $t('page.home.multi_tenant_desc'),
    },
    {
      icon: 'carbon:security',
      title: $t('page.home.enterprise_security'),
      description: $t('page.home.security_desc'),
    },
    {
      icon: 'carbon:analytics',
      title: $t('page.home.advanced_analytics'),
      description: $t('page.home.analytics_desc'),
    },
    {
      icon: 'carbon:api',
      title: $t('page.home.api_integration'),
      description: $t('page.home.api_integration_desc'),
    },
    {
      icon: 'carbon:collaborate',
      title: $t('page.home.real_time_collaboration'),
      description: $t('page.home.real_time_collaboration_desc'),
    },
  ]
})

onMounted(async () => {
  loading.value = true
  await Promise.all([
    loadLatestPosts(),
    loadFeaturedPosts(),
    loadCategories(),
    loadPopularTags(),
  ])
  loading.value = false

  // 数据加载完成后初始化滚动动画
  setTimeout(() => {
    initScrollObserver()
  }, 100)
})

// 监听语言切换，自动重新加载数据
useLanguageChangeEffect(async () => {
  loading.value = true
  await Promise.all([
    loadLatestPosts(),
    loadFeaturedPosts(),
    loadCategories(),
    loadPopularTags(),
  ])
  loading.value = false

  // 数据更新后，需要重新初始化滚动观察器
  // 先销毁旧的观察器
  destroyScrollObserver()
  // 等待 DOM 更新完成后重新初始化
  setTimeout(() => {
    initScrollObserver()
  }, 100)
}, {
  immediate: false,    // 已经在 onMounted 中加载过，不需要立即执行
  autoCleanup: true,  // 组件卸载时自动取消订阅
});

onUnmounted(() => {
  destroyScrollObserver()
})
</script>

<template>
  <div class="home-page">
    <!-- Hero Section -->
    <section class="hero">
      <!-- Animated Background Elements -->
      <div class="hero-bg-wrapper">
        <div class="hero-gradient-bg"></div>

        <!-- Grid Background -->
        <div class="hero-grid-bg"></div>

        <!-- Animated Shapes -->
        <div class="hero-animated-shapes">
          <div class="shape shape-1"></div>
          <div class="shape shape-2"></div>
          <div class="shape shape-3"></div>
        </div>

        <!-- Code Snippet Decorations -->
        <div class="hero-code-snippets">
          <div class="code-snippet snippet-1">
            <div class="code-line"><span class="code-keyword">func</span> <span
              class="code-function">GetPost</span>() {
            </div>
            <div class="code-line"><span class="code-keyword">return</span> post</div>
            <div class="code-line">}</div>
          </div>
          <div class="code-snippet snippet-2">
            <div class="code-line"><span class="code-tag">&lt;template&gt;</span></div>
            <div class="code-line"><span class="code-tag">&lt;div</span> <span class="code-attr">class</span>=<span
              class="code-string">"cms"</span><span class="code-tag">&gt;</span></div>
            <div class="code-line"><span class="code-tag">&lt;/div&gt;</span></div>
          </div>
          <div class="code-snippet snippet-3">
            <div class="code-line"><span class="code-comment">// CMS API</span></div>
            <div class="code-line"><span class="code-keyword">const</span> api = <span
              class="code-string">'v1'</span></div>
          </div>
        </div>

        <svg class="hero-waves" viewBox="0 0 1440 320" preserveAspectRatio="none">
          <path fill="rgba(255,255,255,0.1)"
                d="M0,96L48,112C96,128,192,160,288,160C384,160,480,128,576,122.7C672,117,768,139,864,138.7C960,139,1056,117,1152,101.3C1248,85,1344,75,1392,69.3L1440,64L1440,320L1392,320C1344,320,1248,320,1152,320C1056,320,960,320,864,320C768,320,672,320,576,320C480,320,384,320,288,320C192,320,96,320,48,320L0,320Z"></path>
        </svg>
      </div>

      <div class="hero-content">
        <h1 class="hero-title">{{ $t('authentication.login.brand_title') }}</h1>
        <p class="hero-subtitle">{{ $t('authentication.login.brand_subtitle') }}</p>
        <p class="hero-description">{{ $t('page.home.hero_description') }}</p>
        <div class="hero-actions">
          <n-button type="primary" size="large" class="btn-primary" @click="router.push('/post')">
            <template #icon>
              <XIcon name="carbon:document-export" :size="20"/>
            </template>
            {{ $t('page.home.browse_posts') }}
          </n-button>
          <n-button type="default" size="large" class="btn-secondary"
                    @click="router.push('/about')">
            <template #icon>
              <XIcon name="carbon:information" :size="20"/>
            </template>
            {{ $t('page.home.learn_more') }}
          </n-button>
        </div>
      </div>
    </section>

    <!-- Featured Posts Section (推荐阅读) - 上移到 Categories 之前 -->
    <section v-if="featuredPosts.length > 0 || loading" class="featured-section scroll-reveal">
      <div class="section-header">
        <h2>
          <XIcon name="carbon:star-filled" :size="28"
                 style="color: var(--color-brand); margin-right: 8px;"/>
          {{ $t('page.home.featured_posts') }}
        </h2>
        <n-button text type="primary" @click="router.push('/post')">
          {{ $t('page.home.view_all') }} →
        </n-button>
      </div>
      <!-- Loading Skeleton -->
      <div v-if="loading" class="featured-grid">
        <div v-for="i in 3" :key="i" class="featured-card">
          <n-skeleton height="220px"/>
          <div class="featured-content" style="padding: 24px;">
            <n-skeleton :width="'80%'" size="medium" style="margin-bottom: 12px;"/>
            <n-skeleton :width="'90%'" size="small" style="margin-bottom: 16px;"/>
            <div class="featured-meta" style="display: flex; gap: 12px;">
              <n-skeleton :width="60" size="small"/>
              <n-skeleton :width="60" size="small"/>
              <n-skeleton :width="60" size="small"/>
            </div>
          </div>
        </div>
      </div>
      <!-- Loaded Content -->
      <div v-else class="featured-grid">
        <div
          v-for="post in featuredPosts"
          :key="post.id"
          class="featured-card scroll-reveal-item"
          @click="handleViewPost(post.id)"
        >
          <div class="featured-image">
            <img :src="getPostThumbnail(post)" :alt="getPostTitle(post)"/>
            <div class="featured-overlay">
                <span class="featured-badge">
                  <XIcon name="carbon:star-filled" :size="14"/>
                  推荐
                </span>
            </div>
          </div>
          <div class="featured-content">
            <h3>{{ getPostTitle(post) }}</h3>
            <p>{{ getPostSummary(post) }}</p>
            <div class="featured-meta">
                <span>
                  <XIcon name="carbon:user" :size="14"/>
                  {{ post.authorName || '匿名' }}
                </span>
              <span>
                  <XIcon name="carbon:time" :size="14"/>
                  {{ formatDate(post.createdAt) }}
                </span>
              <span>
                  <XIcon name="carbon:view" :size="14"/>
                  {{ post.visits || 0 }}
                </span>
            </div>
          </div>
        </div>
      </div>

      <!-- Scroll Indicator -->
      <div class="scroll-indicator">
        <div class="scroll-indicator-content">
          <span class="scroll-indicator-text">{{ $t('page.home.explore_more_categories') }}</span>
          <div class="scroll-arrow" @click="scrollToCategories">
            <XIcon name="carbon:arrow-down" :size="24" class="arrow-icon"/>
          </div>
        </div>
      </div>
    </section>

    <!-- Categories Section -->
    <section class="categories-section scroll-reveal">
      <div class="section-header">
        <h2>{{ $t('page.home.categories') }}</h2>
        <n-button text type="primary" @click="router.push('/category')">
          {{ $t('page.home.view_all') }} →
        </n-button>
      </div>
      <!-- Loading Skeleton -->
      <div v-if="loading" class="categories-grid desktop-grid">
        <div v-for="i in 8" :key="i" class="category-card">
          <n-skeleton height="140px"/>
        </div>
      </div>
      <!-- Loaded Content -->
      <div v-else>
        <!-- Desktop: Grid Layout (4 columns) -->
        <div class="categories-grid desktop-grid">
          <div
            v-for="category in categories"
            :key="category.id"
            class="category-card scroll-reveal-item"
            :style="{ '--card-hue': (category.id % 6) * 60 }"
            @click="handleViewCategory(category.id)"
          >
            <!-- Card Background Gradient -->
            <div class="category-card-bg"></div>

            <!-- Content -->
            <div class="category-card-content">
              <div class="category-card-header">
                <div class="category-icon">
                  <XIcon
                    :name="category.icon?.includes(':') ? category.icon : `carbon:${category.icon || 'folder'}`"
                    :size="48"/>
                </div>
                <div class="category-info">
                  <h3>{{ getCategoryName(category) }}</h3>
                  <span class="post-count">
                    {{ $t('page.home.article_count', {count: category.postCount || 0}) }}
                  </span>
                </div>
              </div>

              <!-- Updated Badge -->
              <div class="category-badge">
                <XIcon name="carbon:time" :size="14"/>
                <span>{{ getUpdatedDaysAgo(3) }}</span>
              </div>
            </div>
          </div>
        </div>

        <!-- Mobile: Carousel (走马灯) -->
        <div class="categories-carousel mobile-carousel">
          <n-carousel
            :autoplay="true"
            :interval="5000"
            :show-arrow="true"
            :mouse-wheel="true"
            :touchable="true"
            :slides-per-view="1.5"
            :centered-slides="true"
            :scroll-eventListener="false"
            class="carousel-container"
          >
            <template
              v-for="category in categories"
              :key="category.id"
            >
              <div
                class="carousel-item"
                @click="handleViewCategory(category.id)"
              >
                <div
                  class="category-card carousel-card"
                  :style="{ '--card-hue': (category.id % 6) * 60 }"
                >
                  <!-- Card Background Gradient -->
                  <div class="category-card-bg"></div>

                  <!-- Content -->
                  <div class="category-card-content">
                    <div class="category-card-header">
                      <div class="category-icon">
                        <XIcon
                          :name="category.icon?.includes(':') ? category.icon : `carbon:${category.icon || 'folder'}`"
                          :size="48"/>
                      </div>
                      <div class="category-info">
                        <h3>{{ getCategoryName(category) }}</h3>
                        <span class="post-count">
                          {{ $t('page.home.article_count', {count: category.postCount || 0}) }}
                        </span>
                      </div>
                    </div>

                    <!-- Updated Badge -->
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
      </div>

      <!-- Popular Tags Section -->
      <div class="popular-tags-container">
        <h3 class="tags-title">{{ $t('page.home.popular_tags') }}</h3>
        <div class="tags-grid">
          <div
            v-for="tag in popularTags"
            :key="tag.id"
            class="tag-item"
            :style="{ '--tag-color': tag.color }"
          >
            <span class="tag-label">{{ tag.name }}</span>
          </div>
        </div>
      </div>
    </section>


    <!-- Latest Posts Section -->
    <section class="latest-section scroll-reveal">
      <div class="section-header">
        <h2>{{ $t('page.home.latest_posts') }}</h2>
        <n-button text type="primary" @click="router.push('/post')">
          {{ $t('page.home.view_all') }} →
        </n-button>
      </div>
      <!-- Loading Skeleton -->
      <div v-if="loading" class="posts-grid">
        <div v-for="i in 6" :key="i" class="post-card">
          <n-skeleton height="220px"/>
          <div style="padding: 24px;">
            <n-skeleton :width="'80%'" size="medium" style="margin-bottom: 12px;"/>
            <n-skeleton :width="'90%'" size="small" style="margin-bottom: 16px;"/>
            <div style="display: flex; gap: 12px;">
              <n-skeleton :width="60" size="small"/>
              <n-skeleton :width="60" size="small"/>
              <n-skeleton :width="60" size="small"/>
            </div>
          </div>
        </div>
      </div>
      <!-- Loaded Content -->
      <div v-else class="posts-grid">
        <div
          v-for="post in latestPosts"
          :key="post.id"
          class="post-card scroll-reveal-item"
          @click="handleViewPost(post.id)"
        >
          <div class="post-image">
            <img :src="getPostThumbnail(post)" :alt="getPostTitle(post)"/>
          </div>
          <div class="post-content">
            <h3>{{ getPostTitle(post) }}</h3>
            <p>{{ getPostSummary(post) }}</p>
            <div class="post-meta">
              <span class="author">{{ post.authorName }}</span>
              <span class="date">{{ formatDate(post.createdAt) }}</span>
              <span class="views">{{
                  $t('page.home.views_count', {count: post.visits || 0})
                }}</span>
            </div>
          </div>
        </div>
      </div>
    </section>

    <!-- Features Section -->
    <section class="features">
      <div class="section-header">
        <h2>{{ $t('page.home.platform_features') }}</h2>
      </div>
      <div class="features-grid">
        <div v-for="feature in features" :key="feature.title" class="feature-card">
          <div class="feature-icon">
            <XIcon :name="feature.icon" :size="48"/>
          </div>
          <h3>{{ feature.title }}</h3>
          <p>{{ feature.description }}</p>
        </div>
      </div>
    </section>
  </div>
</template>

<style scoped lang="less">
.home-page {
  min-height: 100vh;
  width: 100%;
  background: var(--color-bg);
}

// ========== Scroll Reveal Animations ==========
.scroll-reveal-item {
  opacity: 0;
  transform: translateY(30px);
  transition: opacity 0.6s cubic-bezier(0.4, 0, 0.2, 1),
  transform 0.6s cubic-bezier(0.4, 0, 0.2, 1);

  &.is-visible {
    opacity: 1;
    transform: translateY(0);
  }

  // 为每个卡片添加延迟，创造波浪效果
  &:nth-child(1) {
    transition-delay: 0.08s;
  }

  &:nth-child(2) {
    transition-delay: 0.16s;
  }

  &:nth-child(3) {
    transition-delay: 0.24s;
  }

  &:nth-child(4) {
    transition-delay: 0.32s;
  }

  &:nth-child(5) {
    transition-delay: 0.4s;
  }

  &:nth-child(6) {
    transition-delay: 0.48s;
  }

  &:nth-child(7) {
    transition-delay: 0.56s;
  }

  &:nth-child(8) {
    transition-delay: 0.64s;
  }

  &:nth-child(9) {
    transition-delay: 0.72s;
  }

  &:nth-child(10) {
    transition-delay: 0.8s;
  }

  &:nth-child(11) {
    transition-delay: 0.88s;
  }

  &:nth-child(12) {
    transition-delay: 0.96s;
  }
}

.scroll-reveal {
  .section-header {
    opacity: 0;
    transform: translateY(20px);
    animation: fadeInUp 0.6s ease-out 0.2s forwards;
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

// Hero Section
.hero {
  padding: 3.5rem 2rem 2.5rem;
  text-align: center;
  max-width: 100%;
  background: linear-gradient(135deg, #6366f1 0%, #8b5cf6 50%, #a855f7 100%);
  color: white;
  position: relative;
  overflow: hidden;
  min-height: 280px;
  height: auto;
  display: flex;
  align-items: center;
  justify-content: center;

  .hero-bg-wrapper {
    position: absolute;
    top: 0;
    left: 0;
    right: 0;
    bottom: 0;
    overflow: hidden;
    z-index: 0;
    pointer-events: none; // 防止背景元素阻挡按钮点击

    .hero-gradient-bg {
      position: absolute;
      top: 0;
      left: 0;
      right: 0;
      bottom: 0;
      background: radial-gradient(ellipse at 20% 50%, rgba(100, 200, 255, 0.15) 0%, transparent 50%),
      radial-gradient(ellipse at 80% 80%, rgba(200, 100, 255, 0.12) 0%, transparent 50%),
      linear-gradient(135deg, #6366f1 0%, #8b5cf6 50%, #a855f7 100%);
      animation: gradientShift 15s ease-in-out infinite;
    }

    .hero-grid-bg {
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
    }

    // Code Snippet Decorations
    .hero-code-snippets {
      position: absolute;
      width: 100%;
      height: 100%;
      pointer-events: none;
      perspective: 1000px; // 添加 3D 透视

      .code-snippet {
        position: absolute;
        background: rgba(0, 0, 0, 0.25);
        backdrop-filter: blur(12px);
        border: 1px solid rgba(255, 255, 255, 0.15);
        border-radius: 12px;
        padding: 14px 18px;
        font-family: 'Fira Code', 'Consolas', monospace;
        font-size: 13px;
        line-height: 1.6;
        box-shadow: 0 10px 40px rgba(0, 0, 0, 0.4),
        0 0 20px rgba(var(--color-primary-purple-rgb), 0.2),
        inset 0 1px 0 rgba(255, 255, 255, 0.1);
        opacity: 0.9;
        transform-style: preserve-3d;
        transition: all 0.3s ease;

        .code-line {
          white-space: nowrap;
          color: rgba(255, 255, 255, 0.95);
          text-shadow: 0 1px 2px rgba(0, 0, 0, 0.5);
        }

        .code-keyword {
          color: #fb7185;
          font-weight: 600;
        }

        .code-function {
          color: #fbbf24;
        }

        .code-tag {
          color: #60a5fa;
        }

        .code-attr {
          color: #a78bfa;
        }

        .code-string {
          color: #34d399;
        }

        .code-comment {
          color: rgba(255, 255, 255, 0.5);
          font-style: italic;
        }

        &.snippet-1 {
          top: 15%;
          right: 8%;
          animation: floatRotate3D 8s ease-in-out infinite,
          glowPulseSnippet 4s ease-in-out infinite;
        }

        &.snippet-2 {
          top: 55%;
          left: 5%;
          animation: floatRotate3D 10s ease-in-out infinite reverse,
          glowPulseSnippet 5s ease-in-out infinite 1s;
        }

        &.snippet-3 {
          bottom: 20%;
          right: 12%;
          animation: floatRotate3D 12s ease-in-out infinite,
          glowPulseSnippet 6s ease-in-out infinite 2s;
        }
      }
    }

    .hero-animated-shapes {
      position: absolute;
      width: 100%;
      height: 100%;

      .shape {
        position: absolute;
        border-radius: 50%;
        opacity: 0.15;
        filter: blur(60px);

        &.shape-1 {
          width: 300px;
          height: 300px;
          background: rgba(255, 255, 255, 0.5);
          top: -150px;
          right: -100px;
          animation: float 8s ease-in-out infinite;
        }

        &.shape-2 {
          width: 400px;
          height: 400px;
          background: rgba(255, 255, 255, 0.3);
          bottom: -200px;
          left: -150px;
          animation: float 10s ease-in-out infinite reverse;
        }

        &.shape-3 {
          width: 250px;
          height: 250px;
          background: rgba(255, 255, 255, 0.4);
          top: 50%;
          left: 50%;
          animation: float 12s ease-in-out infinite;
        }
      }
    }

    .hero-waves {
      position: absolute;
      bottom: 0;
      left: 0;
      width: 100%;
      height: 100px;
      animation: wave 10s linear infinite;
    }
  }

  .hero-content {
    animation: fadeInUp 0.8s ease-out;
    position: relative;
    z-index: 100 !important; // 使用最高 z-index
    max-width: 1200px;
    margin: 0 auto;
    pointer-events: auto !important; // 确保内容区域可交互
  }

  .hero-title {
    font-size: 4.5rem;
    font-weight: 900;
    margin-bottom: 1.5rem;
    line-height: 1.1;
    letter-spacing: -2px;
    text-shadow: 0 4px 16px rgba(0, 0, 0, 0.3),
    0 2px 8px rgba(0, 0, 0, 0.2);
    animation: slideDown 0.8s cubic-bezier(0.34, 1.56, 0.64, 1);
    background: linear-gradient(135deg, #ffffff 0%, rgba(255, 255, 255, 0.95) 100%);
    -webkit-background-clip: text;
    -webkit-text-fill-color: transparent;
    background-clip: text;
  }

  .hero-subtitle {
    font-size: 1.75rem;
    margin-bottom: 1rem;
    max-width: 700px;
    margin-left: auto;
    margin-right: auto;
    opacity: 0.95;
    font-weight: 300;
    letter-spacing: 0.5px;
    text-shadow: 0 2px 10px rgba(0, 0, 0, 0.25);
    animation: fadeInUp 0.8s ease-out 0.2s both;
    line-height: 1.4;
  }

  .hero-description {
    font-size: 1.15rem;
    margin-bottom: 2.5rem;
    max-width: 750px;
    margin-left: auto;
    margin-right: auto;
    opacity: 0.85;
    font-weight: 300;
    letter-spacing: 0.3px;
    text-shadow: 0 1px 8px rgba(0, 0, 0, 0.2);
    animation: fadeInUp 0.8s ease-out 0.4s both;
    line-height: 1.6;
  }

  .hero-actions {
    display: flex;
    gap: 1.5rem;
    justify-content: center;
    flex-wrap: wrap;
    animation: fadeInUp 0.8s ease-out 0.6s both;
    position: relative;
    z-index: 100 !important; // 使用最高 z-index
    pointer-events: auto !important; // 确保按钮区域可点击

    :deep(.n-button) {
      position: relative;
      z-index: 100 !important; // 确保按钮在最上层
      pointer-events: auto !important; // 确保按钮本身可点击
    }

    .btn-primary,
    .btn-secondary {
      pointer-events: auto !important; // 确保按钮可点击
      position: relative;
      z-index: 100 !important;
    }

    .btn-primary {
      :deep(.n-button) {
        background: linear-gradient(135deg, #10b981 0%, #059669 100%) !important;
        color: white !important;
        font-size: 1.05rem;
        font-weight: 700;
        padding: 16px 40px !important;
        border-radius: 12px !important;
        border: none !important;
        box-shadow: 0 8px 24px rgba(16, 185, 129, 0.35),
        0 4px 12px rgba(16, 185, 129, 0.2),
        0 0 0 1px rgba(255, 255, 255, 0.15) inset;
        transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
        text-shadow: 0 1px 2px rgba(0, 0, 0, 0.2);
        position: relative;
        overflow: hidden;

        &::before {
          content: '';
          position: absolute;
          top: 0;
          left: -100%;
          width: 100%;
          height: 100%;
          background: linear-gradient(90deg, transparent, rgba(255, 255, 255, 0.2), transparent);
          transition: left 0.5s;
          pointer-events: none; // 防止伪元素阻挡点击
        }

        &:hover {
          transform: translateY(-2px) scale(1.02);
          box-shadow: 0 12px 32px rgba(16, 185, 129, 0.45),
          0 6px 16px rgba(16, 185, 129, 0.25),
          0 0 0 1px rgba(255, 255, 255, 0.2) inset;
          background: linear-gradient(135deg, #059669 0%, #047857 100%) !important;

          &::before {
            left: 100%;
          }
        }

        &:active {
          transform: translateY(0) scale(0.98);
          box-shadow: 0 4px 16px rgba(16, 185, 129, 0.4),
          0 2px 8px rgba(16, 185, 129, 0.2);
        }
      }
    }

    .btn-secondary {
      :deep(.n-button) {
        background: rgba(255, 255, 255, 0.12) !important;
        color: white !important;
        font-size: 1.05rem;
        font-weight: 600;
        padding: 16px 40px !important;
        border-radius: 12px !important;
        border: 2px solid rgba(255, 255, 255, 0.3) !important;
        backdrop-filter: blur(12px);
        box-shadow: 0 4px 16px rgba(0, 0, 0, 0.15),
        0 0 0 1px rgba(255, 255, 255, 0.1) inset;
        transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
        text-shadow: 0 1px 3px rgba(0, 0, 0, 0.3);
        position: relative;
        overflow: hidden;

        &::before {
          content: '';
          position: absolute;
          top: 0;
          left: -100%;
          width: 100%;
          height: 100%;
          background: linear-gradient(90deg, transparent, rgba(255, 255, 255, 0.1), transparent);
          transition: left 0.5s;
          pointer-events: none; // 防止伪元素阻挡点击
        }

        &:hover {
          background: rgba(255, 255, 255, 0.2) !important;
          border-color: rgba(255, 255, 255, 0.5) !important;
          transform: translateY(-2px) scale(1.02);
          box-shadow: 0 8px 24px rgba(0, 0, 0, 0.2),
          0 0 0 1px rgba(255, 255, 255, 0.15) inset;

          &::before {
            left: 100%;
          }
        }

        &:active {
          transform: translateY(0) scale(0.98);
          background: rgba(255, 255, 255, 0.15) !important;
          box-shadow: 0 2px 12px rgba(0, 0, 0, 0.15);
        }
      }
    }
  }

  .hero-content {
    animation: fadeInUp 0.8s ease-out;
    position: relative;
    z-index: 1;
    max-width: 1200px;
    margin: 0 auto;
  }

  .hero-title {
    font-size: 4.5rem;
    font-weight: 900;
    margin-bottom: 1.5rem;
    line-height: 1.1;
    letter-spacing: -2px;
    position: relative;

    // 使用渐变文字效果
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

    // 多层外发光和阴影，创造强烈的视觉冲击
    text-shadow: 0 0 40px rgba(255, 255, 255, 0.9),
    0 0 80px rgba(var(--color-primary-purple-rgb), 0.7),
    0 0 120px rgba(99, 102, 241, 0.5),
    0 6px 24px rgba(0, 0, 0, 0.5),
    0 3px 12px rgba(0, 0, 0, 0.4);

    // 添加多层滤镜外发光
    filter: drop-shadow(0 0 30px rgba(var(--color-primary-purple-rgb), 0.6)) drop-shadow(0 0 60px rgba(99, 102, 241, 0.4)) drop-shadow(0 5px 15px rgba(0, 0, 0, 0.3));

    // 添加渐变动画和发光脉冲动画
    animation-name: slideDown, glowPulseTitle, gradientShine;
    animation-duration: 0.8s, 3s, 6s;
    animation-timing-function: cubic-bezier(0.34, 1.56, 0.64, 1), ease-in-out, linear;
    animation-iteration-count: 1, infinite, infinite;
    animation-delay: 0s, 0.5s, 0s;

    // 在标题下方添加光晕效果
    &::after {
      content: '';
      position: absolute;
      bottom: -20px;
      left: 50%;
      transform: translateX(-50%);
      width: 80%;
      height: 40px;
      background: radial-gradient(ellipse at center,
      rgba(var(--color-primary-purple-rgb), 0.3) 0%,
      rgba(99, 102, 241, 0.2) 50%,
      transparent 100%);
      filter: blur(20px);
      opacity: 0.8;
      animation: pulseGlow 3s ease-in-out infinite;
    }
  }

  .hero-subtitle {
    font-size: 1.75rem;
    margin-bottom: 1rem;
    max-width: 700px;
    margin-left: auto;
    margin-right: auto;
    opacity: 0.95;
    font-weight: 300;
    letter-spacing: 0.5px;
    text-shadow: 0 2px 10px rgba(0, 0, 0, 0.25);
    animation: fadeInUp 0.8s ease-out 0.2s both;
    line-height: 1.4;
  }

  .hero-description {
    font-size: 1.15rem;
    margin-bottom: 2.5rem;
    max-width: 750px;
    margin-left: auto;
    margin-right: auto;
    opacity: 0.85;
    font-weight: 300;
    letter-spacing: 0.3px;
    text-shadow: 0 1px 8px rgba(0, 0, 0, 0.2);
    animation: fadeInUp 0.8s ease-out 0.4s both;
    line-height: 1.6;
  }

  .hero-actions {
    display: flex;
    gap: 1.5rem;
    justify-content: center;
    flex-wrap: wrap;
    animation: fadeInUp 0.8s ease-out 0.6s both;
    position: relative;
    z-index: 100 !important; // 使用最高 z-index
    pointer-events: auto !important; // 确保按钮区域可点击

    :deep(.n-button) {
      position: relative;
      z-index: 100 !important; // 确保按钮在最上层
      pointer-events: auto !important; // 确保按钮本身可点击
    }

    .btn-primary,
    .btn-secondary {
      pointer-events: auto !important; // 确保按钮可点击
      position: relative;
      z-index: 100 !important;
    }

    .btn-primary {
      :deep(.n-button) {
        background: linear-gradient(135deg, #10b981 0%, #059669 100%) !important;
        color: white !important;
        font-size: 1.05rem;
        font-weight: 700;
        padding: 16px 40px !important;
        border-radius: 12px !important;
        border: none !important;
        box-shadow: 0 8px 24px rgba(16, 185, 129, 0.35),
        0 4px 12px rgba(16, 185, 129, 0.2),
        0 0 0 1px rgba(255, 255, 255, 0.15) inset;
        transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
        text-shadow: 0 1px 2px rgba(0, 0, 0, 0.2);
        position: relative;
        overflow: hidden;

        &::before {
          content: '';
          position: absolute;
          top: 0;
          left: -100%;
          width: 100%;
          height: 100%;
          background: linear-gradient(90deg, transparent, rgba(255, 255, 255, 0.2), transparent);
          transition: left 0.5s;
          pointer-events: none; // 防止伪元素阻挡点击
        }

        &:hover {
          transform: translateY(-2px) scale(1.02);
          box-shadow: 0 12px 32px rgba(16, 185, 129, 0.45),
          0 6px 16px rgba(16, 185, 129, 0.25),
          0 0 0 1px rgba(255, 255, 255, 0.2) inset;
          background: linear-gradient(135deg, #059669 0%, #047857 100%) !important;

          &::before {
            left: 100%;
          }
        }

        &:active {
          transform: translateY(0) scale(0.98);
          box-shadow: 0 4px 16px rgba(16, 185, 129, 0.4),
          0 2px 8px rgba(16, 185, 129, 0.2);
        }
      }
    }

    .btn-secondary {
      :deep(.n-button) {
        background: rgba(255, 255, 255, 0.12) !important;
        color: white !important;
        font-size: 1.05rem;
        font-weight: 600;
        padding: 16px 40px !important;
        border-radius: 12px !important;
        border: 2px solid rgba(255, 255, 255, 0.3) !important;
        backdrop-filter: blur(12px);
        box-shadow: 0 4px 16px rgba(0, 0, 0, 0.15),
        0 0 0 1px rgba(255, 255, 255, 0.1) inset;
        transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
        text-shadow: 0 1px 3px rgba(0, 0, 0, 0.3);
        position: relative;
        overflow: hidden;

        &::before {
          content: '';
          position: absolute;
          top: 0;
          left: -100%;
          width: 100%;
          height: 100%;
          background: linear-gradient(90deg, transparent, rgba(255, 255, 255, 0.1), transparent);
          transition: left 0.5s;
          pointer-events: none; // 防止伪元素阻挡点击
        }

        &:hover {
          background: rgba(255, 255, 255, 0.2) !important;
          border-color: rgba(255, 255, 255, 0.5) !important;
          transform: translateY(-2px) scale(1.02);
          box-shadow: 0 8px 24px rgba(0, 0, 0, 0.2),
          0 0 0 1px rgba(255, 255, 255, 0.15) inset;

          &::before {
            left: 100%;
          }
        }

        &:active {
          transform: translateY(0) scale(0.98);
          background: rgba(255, 255, 255, 0.15) !important;
          box-shadow: 0 2px 12px rgba(0, 0, 0, 0.15);
        }
      }
    }
  }

  .hero-content {
    animation: fadeInUp 0.8s ease-out;
    position: relative;
    z-index: 1;
    max-width: 1200px;
    margin: 0 auto;
  }

  .hero-title {
    font-size: 4.5rem;
    font-weight: 900;
    margin-bottom: 1.5rem;
    line-height: 1.1;
    letter-spacing: -2px;
    position: relative;

    // 使用渐变文字效果
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

    // 多层外发光和阴影，创造强烈的视觉冲击
    text-shadow: 0 0 40px rgba(255, 255, 255, 0.9),
    0 0 80px rgba(var(--color-primary-purple-rgb), 0.7),
    0 0 120px rgba(99, 102, 241, 0.5),
    0 6px 24px rgba(0, 0, 0, 0.5),
    0 3px 12px rgba(0, 0, 0, 0.4);

    // 添加多层滤镜外发光
    filter: drop-shadow(0 0 30px rgba(var(--color-primary-purple-rgb), 0.6)) drop-shadow(0 0 60px rgba(99, 102, 241, 0.4)) drop-shadow(0 6px 20px rgba(0, 0, 0, 0.4));

    // 添加渐变动画和发光脉冲动画
    animation-name: slideDown, glowPulseTitle, gradientShine;
    animation-duration: 0.8s, 3s, 6s;
    animation-timing-function: cubic-bezier(0.34, 1.56, 0.64, 1), ease-in-out, linear;
    animation-iteration-count: 1, infinite, infinite;
    animation-delay: 0s, 0.5s, 0s;

    // 在标题下方添加光晕效果
    &::after {
      content: '';
      position: absolute;
      bottom: -20px;
      left: 50%;
      transform: translateX(-50%);
      width: 80%;
      height: 40px;
      background: radial-gradient(ellipse at center,
      rgba(var(--color-primary-purple-rgb), 0.3) 0%,
      rgba(99, 102, 241, 0.2) 50%,
      transparent 100%);
      filter: blur(20px);
      opacity: 0.8;
      animation: pulseGlow 3s ease-in-out infinite;
    }
  }

  .hero-subtitle {
    font-size: 1.75rem;
    margin-bottom: 1rem;
    max-width: 700px;
    margin-left: auto;
    margin-right: auto;
    opacity: 0.95;
    font-weight: 300;
    letter-spacing: 0.5px;
    text-shadow: 0 2px 10px rgba(0, 0, 0, 0.25);
    animation: fadeInUp 0.8s ease-out 0.2s both;
    line-height: 1.4;
  }

  .hero-description {
    font-size: 1.15rem;
    margin-bottom: 2.5rem;
    max-width: 750px;
    margin-left: auto;
    margin-right: auto;
    opacity: 0.85;
    font-weight: 300;
    letter-spacing: 0.3px;
    text-shadow: 0 1px 8px rgba(0, 0, 0, 0.2);
    animation: fadeInUp 0.8s ease-out 0.4s both;
    line-height: 1.6;
  }

  .hero-actions {
    display: flex;
    gap: 1.5rem;
    justify-content: center;
    flex-wrap: wrap;
    animation: fadeInUp 0.8s ease-out 0.6s both;
    position: relative;
    z-index: 100 !important; // 使用最高 z-index
    pointer-events: auto !important; // 确保按钮区域可点击

    :deep(.n-button) {
      position: relative;
      z-index: 100 !important; // 确保按钮在最上层
      pointer-events: auto !important; // 确保按钮本身可点击
    }

    .btn-primary,
    .btn-secondary {
      pointer-events: auto !important; // 确保按钮可点击
      position: relative;
      z-index: 100 !important;
    }

    .btn-primary {
      :deep(.n-button) {
        background: linear-gradient(135deg, #10b981 0%, #059669 100%) !important;
        color: white !important;
        font-size: 1.05rem;
        font-weight: 700;
        padding: 16px 40px !important;
        border-radius: 12px !important;
        border: none !important;
        box-shadow: 0 8px 24px rgba(16, 185, 129, 0.35),
        0 4px 12px rgba(16, 185, 129, 0.2),
        0 0 0 1px rgba(255, 255, 255, 0.15) inset;
        transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
        text-shadow: 0 1px 2px rgba(0, 0, 0, 0.2);
        position: relative;
        overflow: hidden;

        &::before {
          content: '';
          position: absolute;
          top: 0;
          left: -100%;
          width: 100%;
          height: 100%;
          background: linear-gradient(90deg, transparent, rgba(255, 255, 255, 0.2), transparent);
          transition: left 0.5s;
          pointer-events: none; // 防止伪元素阻挡点击
        }

        &:hover {
          transform: translateY(-2px) scale(1.02);
          box-shadow: 0 12px 32px rgba(16, 185, 129, 0.45),
          0 6px 16px rgba(16, 185, 129, 0.25),
          0 0 0 1px rgba(255, 255, 255, 0.2) inset;
          background: linear-gradient(135deg, #059669 0%, #047857 100%) !important;

          &::before {
            left: 100%;
          }
        }

        &:active {
          transform: translateY(0) scale(0.98);
          box-shadow: 0 4px 16px rgba(16, 185, 129, 0.4),
          0 2px 8px rgba(16, 185, 129, 0.2);
        }
      }
    }

    .btn-secondary {
      :deep(.n-button) {
        background: rgba(255, 255, 255, 0.12) !important;
        color: white !important;
        font-size: 1.05rem;
        font-weight: 600;
        padding: 16px 40px !important;
        border-radius: 12px !important;
        border: 2px solid rgba(255, 255, 255, 0.3) !important;
        backdrop-filter: blur(12px);
        box-shadow: 0 4px 16px rgba(0, 0, 0, 0.15),
        0 0 0 1px rgba(255, 255, 255, 0.1) inset;
        transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
        text-shadow: 0 1px 3px rgba(0, 0, 0, 0.3);
        position: relative;
        overflow: hidden;

        &::before {
          content: '';
          position: absolute;
          top: 0;
          left: -100%;
          width: 100%;
          height: 100%;
          background: linear-gradient(90deg, transparent, rgba(255, 255, 255, 0.1), transparent);
          transition: left 0.5s;
          pointer-events: none; // 防止伪元素阻挡点击
        }

        &:hover {
          background: rgba(255, 255, 255, 0.2) !important;
          border-color: rgba(255, 255, 255, 0.5) !important;
          transform: translateY(-2px) scale(1.02);
          box-shadow: 0 8px 24px rgba(0, 0, 0, 0.2),
          0 0 0 1px rgba(255, 255, 255, 0.15) inset;

          &::before {
            left: 100%;
          }
        }

        &:active {
          transform: translateY(0) scale(0.98);
          background: rgba(255, 255, 255, 0.15) !important;
          box-shadow: 0 2px 12px rgba(0, 0, 0, 0.15);
        }
      }
    }
  }

  .hero-content {
    animation: fadeInUp 0.8s ease-out;
    position: relative;
    z-index: 1;
    max-width: 1200px;
    margin: 0 auto;
  }

  .hero-title {
    font-size: 4.5rem;
    font-weight: 900;
    margin-bottom: 1.5rem;
    line-height: 1.1;
    letter-spacing: -2px;
    position: relative;

    // 使用渐变文字效果
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

    // 多层外发光和阴影，创造强烈的视觉冲击
    text-shadow: 0 0 40px rgba(255, 255, 255, 0.9),
    0 0 80px rgba(var(--color-primary-purple-rgb), 0.7),
    0 0 120px rgba(99, 102, 241, 0.5),
    0 6px 24px rgba(0, 0, 0, 0.5),
    0 3px 12px rgba(0, 0, 0, 0.4);

    // 添加多层滤镜外发光
    filter: drop-shadow(0 0 30px rgba(var(--color-primary-purple-rgb), 0.6)) drop-shadow(0 0 60px rgba(99, 102, 241, 0.4)) drop-shadow(0 6px 20px rgba(0, 0, 0, 0.4));

    // 添加渐变动画和发光脉冲动画
    animation-name: slideDown, glowPulseTitle, gradientShine;
    animation-duration: 0.8s, 3s, 6s;
    animation-timing-function: cubic-bezier(0.34, 1.56, 0.64, 1), ease-in-out, linear;
    animation-iteration-count: 1, infinite, infinite;
    animation-delay: 0s, 0.5s, 0s;

    // 在标题下方添加光晕效果
    &::after {
      content: '';
      position: absolute;
      bottom: -20px;
      left: 50%;
      transform: translateX(-50%);
      width: 80%;
      height: 40px;
      background: radial-gradient(ellipse at center,
      rgba(var(--color-primary-purple-rgb), 0.3) 0%,
      rgba(99, 102, 241, 0.2) 50%,
      transparent 100%);
      filter: blur(20px);
      opacity: 0.8;
      animation: pulseGlow 3s ease-in-out infinite;
    }
  }

  .hero-subtitle {
    font-size: 1.75rem;
    margin-bottom: 1rem;
    max-width: 700px;
    margin-left: auto;
    margin-right: auto;
    opacity: 0.95;
    font-weight: 300;
    letter-spacing: 0.5px;
    text-shadow: 0 2px 10px rgba(0, 0, 0, 0.25);
    animation: fadeInUp 0.8s ease-out 0.2s both;
    line-height: 1.4;
  }

  .hero-description {
    font-size: 1.15rem;
    margin-bottom: 2.5rem;
    max-width: 750px;
    margin-left: auto;
    margin-right: auto;
    opacity: 0.85;
    font-weight: 300;
    letter-spacing: 0.3px;
    text-shadow: 0 1px 8px rgba(0, 0, 0, 0.2);
    animation: fadeInUp 0.8s ease-out 0.4s both;
    line-height: 1.6;
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

@keyframes float {
  0%, 100% {
    transform: translate(0, 0) rotate(0deg);
  }
  25% {
    transform: translate(30px, -20px) rotate(90deg);
  }
  50% {
    transform: translate(0, -40px) rotate(180deg);
  }
  75% {
    transform: translate(-30px, -20px) rotate(270deg);
  }
}

@keyframes wave {
  0% {
    transform: translateX(0);
  }
  100% {
    transform: translateX(100%);
  }
}

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

@keyframes gridMove {
  0% {
    transform: translate(0, 0);
  }
  100% {
    transform: translate(50px, 50px);
  }
}

@keyframes fadeInOut {
  0%, 100% {
    opacity: 0.5;
  }
  50% {
    opacity: 1;
  }
}

@keyframes glowPulse {
  0%, 100% {
    filter: drop-shadow(0 0 10px rgba(var(--color-primary-purple-rgb), 0.3));
    text-shadow: 0 0 20px rgba(255, 255, 255, 0.6),
    0 0 40px rgba(var(--color-primary-purple-rgb), 0.4),
    0 4px 15px rgba(0, 0, 0, 0.3),
    0 2px 5px rgba(0, 0, 0, 0.2);
  }
  50% {
    filter: drop-shadow(0 0 30px rgba(var(--color-primary-purple-rgb), 0.7));
    text-shadow: 0 0 40px rgba(255, 255, 255, 0.9),
    0 0 70px rgba(var(--color-primary-purple-rgb), 0.8),
    0 4px 20px rgba(0, 0, 0, 0.4),
    0 2px 8px rgba(0, 0, 0, 0.3);
  }
}

// 代码卡片 3D 旋转悬浮动画
@keyframes floatRotate3D {
  0%, 100% {
    transform: translate(0, 0) rotateX(0deg) rotateY(0deg) rotateZ(0deg);
  }
  25% {
    transform: translate(15px, -25px) rotateX(5deg) rotateY(-8deg) rotateZ(3deg);
  }
  50% {
    transform: translate(0, -50px) rotateX(-5deg) rotateY(8deg) rotateZ(-3deg);
  }
  75% {
    transform: translate(-15px, -25px) rotateX(5deg) rotateY(-5deg) rotateZ(5deg);
  }
}

// 代码卡片发光脉冲动画
@keyframes glowPulseSnippet {
  0%, 100% {
    box-shadow: 0 10px 40px rgba(0, 0, 0, 0.3),
    0 0 15px rgba(var(--color-primary-purple-rgb), 0.15),
    inset 0 1px 0 rgba(255, 255, 255, 0.1);
    opacity: 0.85;
  }
  50% {
    box-shadow: 0 15px 50px rgba(0, 0, 0, 0.4),
    0 0 30px rgba(var(--color-primary-purple-rgb), 0.4),
    0 0 60px rgba(99, 102, 241, 0.2),
    inset 0 1px 0 rgba(255, 255, 255, 0.15);
    opacity: 1;
  }
}

// 标题增强发光脉冲动画
@keyframes glowPulseTitle {
  0%, 100% {
    filter: drop-shadow(0 0 25px rgba(var(--color-primary-purple-rgb), 0.5)) drop-shadow(0 0 50px rgba(99, 102, 241, 0.3)) drop-shadow(0 5px 15px rgba(0, 0, 0, 0.3));
  }
  50% {
    filter: drop-shadow(0 0 40px rgba(var(--color-primary-purple-rgb), 0.8)) drop-shadow(0 0 80px rgba(99, 102, 241, 0.6)) drop-shadow(0 8px 20px rgba(0, 0, 0, 0.4));
  }
}

// 渐变光泽扫过动画
@keyframes gradientShine {
  0% {
    background-position: 0 center;
  }
  100% {
    background-position: 200% center;
  }
}

// 底部光晕脉冲动画
@keyframes pulseGlow {
  0%, 100% {
    opacity: 0.6;
    transform: translateX(-50%) scale(1);
  }
  50% {
    opacity: 1;
    transform: translateX(-50%) scale(1.1);
  }
}

// 图标弹跳动画
@keyframes iconBounce {
  0% {
    transform: scale(1) rotate(0deg);
  }
  30% {
    transform: scale(1.25) rotate(15deg);
  }
  50% {
    transform: scale(1.15) rotate(10deg);
  }
  70% {
    transform: scale(1.22) rotate(13deg);
  }
  100% {
    transform: scale(1.2) rotate(12deg);
  }
}

// 卡片轻微呼吸动画
@keyframes cardBreath {
  0%, 100% {
    box-shadow: 0 2px 8px rgba(0, 0, 0, 0.06),
    0 1px 3px rgba(0, 0, 0, 0.04);
  }
  50% {
    box-shadow: 0 4px 12px rgba(0, 0, 0, 0.08),
    0 2px 5px rgba(0, 0, 0, 0.05);
  }
}

// Bounce Down Animation for Scroll Indicator
@keyframes bounceDown {
  0%, 100% {
    transform: translateY(0);
  }
  50% {
    transform: translateY(8px);
  }
}

// Fade In Scale Animation
@keyframes fadeInScale {
  0% {
    opacity: 0;
    transform: scale(0.9);
  }
  100% {
    opacity: 1;
    transform: scale(1);
  }
}

// Section Common Styles
.section-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 2rem;
  max-width: 1200px;
  margin-left: auto;
  margin-right: auto;
  padding: 0 2rem;

  h2 {
    font-size: 2rem;
    font-weight: 600;
    color: var(--color-text-primary);
    margin: 0;
  }
}

// Categories Section
.categories-section {
  padding: 3rem 0;
  background: var(--color-surface);

  .section-header {
    max-width: 1200px;
    margin-left: auto;
    margin-right: auto;
    padding: 0 2rem;
    margin-bottom: 2rem;
  }

  .categories-grid {
    display: grid;
    grid-template-columns: repeat(4, 1fr);
    gap: 1.5rem;
    max-width: 1200px;
    padding: 0 2rem;
    margin: 0 auto 3rem;
  }

  // 走马灯轮播样式
  .categories-carousel {
    display: none; // 默认隐藏，手机端显示
    margin-bottom: 3rem;

    :deep(.n-carousel) {
      .n-carousel__viewport {
        border-radius: 16px;
      }

      .n-carousel__arrow {
        color: var(--color-brand);
        background: rgba(255, 255, 255, 0.9);
        border-radius: 50%;
        transition: all 0.3s;

        &:hover {
          background: var(--color-brand);
          color: white;
          transform: scale(1.1);
        }
      }

      .n-carousel__dots {
        bottom: -40px;
      }

      .n-carousel__dot {
        background: rgba(0, 0, 0, 0.2);
        border-radius: 50%;
        transition: all 0.3s;

        &.n-carousel__dot--active {
          background: var(--color-brand);
        }

        &:hover {
          background: var(--color-brand);
        }
      }
    }

    .carousel-item {
      width: calc(100vw - 2rem);
      padding: 0 0.5rem;
      display: flex;
      align-items: center;
    }

    .carousel-card {
      margin: 0;
      width: 100%;
    }
  }

  // 桌面网格显示
  .desktop-grid {
    display: grid !important;
  }

  .category-card {
    display: flex;
    flex-direction: column;
    padding: 1.5rem;
    background: var(--color-card-bg, linear-gradient(135deg, #ffffff 0%, #f9fafb 100%));
    border: 1px solid var(--color-card-border, #e5e7eb);
    border-radius: 16px;
    cursor: pointer;
    transition: all 0.4s cubic-bezier(0.4, 0, 0.2, 1);
    position: relative;
    overflow: visible;
    box-shadow: 0 2px 8px var(--color-card-shadow, rgba(0, 0, 0, 0.06)),
    0 1px 3px rgba(0, 0, 0, 0.04);
    min-height: 200px;
    height: 100%;

    // 添加轻微的呼吸动画
    animation: cardBreath 4s ease-in-out infinite;

    &::before {
      content: '';
      position: absolute;
      top: 0;
      left: 0;
      right: 0;
      bottom: 0;
      background: linear-gradient(135deg,
      hsl(var(--card-hue), 70%, 60%) 0%,
      hsl(var(--card-hue), 60%, 70%) 100%);
      opacity: 0;
      transition: opacity 0.4s ease;
      z-index: 0;
      border-radius: 16px;
      pointer-events: none;
    }

    &::after {
      content: '';
      position: absolute;
      top: 0;
      left: -100%;
      width: 100%;
      height: 100%;
      background: linear-gradient(90deg, transparent, rgba(255, 255, 255, 0.3), transparent);
      transition: left 0.6s;
      z-index: 0;
      pointer-events: none;
    }

    &:hover {
      transform: translateY(-12px) scale(1.02);
      box-shadow: 0 20px 40px rgba(0, 0, 0, 0.15),
      0 8px 16px rgba(0, 0, 0, 0.1),
      0 0 0 1px hsl(var(--card-hue), 70%, 60%);
      border-color: hsl(var(--card-hue), 70%, 60%);

      // 背景加深效果
      background: linear-gradient(135deg,
      rgba(250, 250, 255, 1) 0%,
      rgba(248, 248, 253, 1) 100%);

      &::before {
        opacity: 0.2;
      }

      &::after {
        left: 100%;
      }

      .category-icon {
        transform: scale(1.2) rotate(12deg);
        color: white;
        background: linear-gradient(135deg,
        hsl(var(--card-hue), 70%, 60%) 0%,
        hsl(var(--card-hue), 60%, 70%) 100%);
        box-shadow: 0 10px 20px hsl(var(--card-hue), 70%, 60%, 0.4),
        0 0 30px hsl(var(--card-hue), 70%, 60%, 0.3);
        animation: iconBounce 0.6s cubic-bezier(0.68, -0.55, 0.265, 1.55);
      }

      .category-info h3 {
        color: hsl(var(--card-hue), 70%, 45%);
        transform: translateX(3px);
      }

      .post-count {
        color: hsl(var(--card-hue), 60%, 50%);
      }

      .category-badge {
        transform: scale(1.05);
        box-shadow: 0 6px 12px hsla(var(--card-hue), 70%, 60%, 0.3);
        background: linear-gradient(135deg,
        hsla(var(--card-hue), 70%, 60%, 0.2) 0%,
        hsla(var(--card-hue), 60%, 70%, 0.15) 100%);
      }
    }

    .category-card-bg {
      position: absolute;
      top: 0;
      left: 0;
      right: 0;
      bottom: 0;
      background: linear-gradient(135deg,
      hsla(var(--card-hue), 70%, 60%, 0.05) 0%,
      hsla(var(--card-hue), 60%, 70%, 0.05) 100%);
      z-index: 0;
      pointer-events: none;
      border-radius: 16px;
    }

    .category-card-content {
      display: flex;
      flex-direction: column;
      align-items: flex-start;
      gap: 0.75rem;
      position: relative;
      z-index: 1;
      width: 100%;
      height: 100%;
    }

    .category-card-header {
      display: flex;
      align-items: flex-start;
      gap: 1rem;
      width: 100%;
    }

    .category-icon {
      font-size: 2.5rem;
      width: 70px;
      height: 70px;
      display: flex;
      align-items: center;
      justify-content: center;
      background: linear-gradient(135deg,
      hsla(var(--card-hue), 70%, 60%, 0.15) 0%,
      hsla(var(--card-hue), 60%, 70%, 0.1) 100%);
      border-radius: 14px;
      color: hsl(var(--card-hue), 70%, 60%);
      flex-shrink: 0;
      transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
      box-shadow: 0 2px 8px hsla(var(--card-hue), 70%, 60%, 0.1);
    }

    .category-info {
      flex: 1;
      width: 100%;

      h3 {
        font-size: 1.1rem;
        font-weight: 700;
        margin: 0 0 0.5rem 0;
        color: var(--color-text-primary);
        letter-spacing: 0.3px;
        transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
        word-break: break-word;

        // 浅色模式下优化文字颜色，提高对比度
        @media (prefers-color-scheme: light) {
          color: #111827 !important;
          font-weight: 800;
        }
      }

      .post-count {
        font-size: 0.95rem;
        color: var(--color-text-secondary);
        font-weight: 600;
        letter-spacing: 0.2px;
        transition: all 0.3s ease;

        // 浅色模式下优化文字颜色
        @media (prefers-color-scheme: light) {
          color: #4b5563 !important;
          font-weight: 700;
        }
      }
    }

    .category-badge {
      position: relative;
      z-index: 2;
      display: flex;
      align-items: center;
      gap: 0.5rem;
      margin-top: auto;
      padding: 0.5rem 0.875rem;
      background: linear-gradient(135deg,
      hsla(var(--card-hue), 70%, 60%, 0.1) 0%,
      hsla(var(--card-hue), 60%, 70%, 0.08) 100%);
      border-radius: 8px;
      border-left: 3px solid hsl(var(--card-hue), 70%, 60%);
      color: hsl(var(--card-hue), 70%, 50%);
      font-size: 0.8rem;

      // 浅色模式下优化文字颜色和背景
      @media (prefers-color-scheme: light) {
        color: hsl(var(--card-hue), 70%, 40%) !important;
        font-weight: 700;
        background: linear-gradient(135deg,
        hsla(var(--card-hue), 70%, 60%, 0.18) 0%,
        hsla(var(--card-hue), 60%, 70%, 0.14) 100%) !important;
        border-left-width: 4px;
      }
      font-weight: 600;
      width: fit-content;
      white-space: nowrap;
      transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
      opacity: 0.9;
    }

    &:hover .category-badge {
      opacity: 1;
    }
  }
}

// Popular Tags Section
.popular-tags-container {
  max-width: 1200px;
  margin: 0 auto;
  padding: 0 2rem;

  .tags-title {
    font-size: 1.25rem;
    font-weight: 600;
    margin: 0 0 1.5rem 0;
    color: var(--color-text-primary);

    @media (prefers-color-scheme: light) {
      color: #1f2937 !important;
    }
  }

  .tags-grid {
    display: flex;
    flex-wrap: wrap;
    gap: 1rem;
  }

  .tag-item {
    display: inline-flex;
    align-items: center;
    padding: 0.5rem 1rem;
    background: linear-gradient(135deg,
    rgba(255, 255, 255, 0.8) 0%,
    rgba(255, 255, 255, 0.6) 100%);
    border: 1.5px solid;
    border-color: var(--tag-color);
    border-radius: 24px;
    cursor: pointer;
    transition: all 0.3s ease;

    &:hover {
      background: var(--tag-color);
      color: white;
      transform: translateY(-2px);
      box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);

      .tag-label {
        color: white !important;
      }
    }
  }

  .tag-label {
    font-size: 0.9rem;
    font-weight: 500;
    color: var(--tag-color);
    transition: color 0.3s ease;

    @media (prefers-color-scheme: light) {
      color: var(--tag-color);
    }
  }
}

// Featured Posts Section
.featured-section {
  padding: 2.5rem 0;
  background: var(--color-bg);

  .featured-grid {
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(340px, 1fr));
    gap: 2rem;
    max-width: 1200px;
    margin: 0 auto;
    padding: 0 2rem;
  }

  .featured-card {
    background: var(--color-card-bg, #ffffff);
    border-radius: 16px;
    overflow: hidden;
    cursor: pointer;
    transition: all 0.4s cubic-bezier(0.4, 0, 0.2, 1);
    border: 1px solid var(--color-card-border, #e5e7eb);
    box-shadow: 0 2px 8px var(--color-card-shadow, rgba(0, 0, 0, 0.08)), 0 1px 3px rgba(0, 0, 0, 0.06);
    display: flex;
    flex-direction: column;
    height: 100%;

    &:hover {
      transform: translateY(-8px);
      box-shadow: 0 16px 32px var(--color-hover-shadow, rgba(0, 0, 0, 0.15)), 0 8px 16px var(--color-hover-shadow, rgba(0, 0, 0, 0.08));
      border-color: var(--color-brand);

      .featured-image {
        img {
          transform: scale(1.08);
        }

        .featured-overlay {
          opacity: 1;
        }
      }

      .featured-content h3 {
        color: var(--color-brand);
      }
    }

    .featured-image {
      position: relative;
      width: 100%;
      height: 220px;
      overflow: hidden;
      background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);

      img {
        width: 100%;
        height: 100%;
        object-fit: cover;
        transition: transform 0.6s cubic-bezier(0.4, 0, 0.2, 1);
      }

      .featured-overlay {
        position: absolute;
        top: 0;
        left: 0;
        right: 0;
        bottom: 0;
        background: linear-gradient(180deg,
        rgba(0, 0, 0, 0) 0%,
        rgba(0, 0, 0, 0.3) 100%);
        opacity: 0;
        transition: opacity 0.3s ease;
        display: flex;
        align-items: flex-start;
        justify-content: flex-end;
        padding: 12px;
      }

      .featured-badge {
        background: linear-gradient(135deg, #f59e0b 0%, #f97316 100%);
        color: white;
        padding: 6px 12px;
        border-radius: 20px;
        font-size: 12px;
        font-weight: 600;
        display: flex;
        align-items: center;
        gap: 4px;
        box-shadow: 0 2px 8px rgba(245, 158, 11, 0.4);
        backdrop-filter: blur(8px);

        :deep(svg) {
          filter: drop-shadow(0 0 2px rgba(255, 255, 255, 0.8));
        }
      }
    }

    .featured-content {
      padding: 24px;
      flex: 1;
      display: flex;
      flex-direction: column;

      h3 {
        font-size: 20px;
        font-weight: 600;
        margin: 0 0 12px 0;
        color: var(--color-text-primary);
        line-height: 1.4;
        transition: color 0.3s ease;
        display: -webkit-box;
        -webkit-line-clamp: 2;
        -webkit-box-orient: vertical;
        overflow: hidden;
        text-overflow: ellipsis;
      }

      p {
        font-size: 15px;
        line-height: 1.7;
        color: var(--color-text-secondary);
        margin: 0 0 16px 0;
        flex: 1;
        display: -webkit-box;
        -webkit-line-clamp: 3;
        -webkit-box-orient: vertical;
        overflow: hidden;
        text-overflow: ellipsis;
      }

      .featured-meta {
        display: flex;
        align-items: center;
        gap: 16px;
        padding-top: 16px;
        border-top: 1px solid var(--color-border);
        font-size: 13px;
        color: var(--color-text-secondary);

        span {
          display: flex;
          align-items: center;
          gap: 4px;

          :deep(svg) {
            opacity: 0.7;
          }
        }
      }
    }
  }

  // Scroll Indicator
  .scroll-indicator {
    display: flex;
    justify-content: center;
    align-items: center;
    padding: 3rem 0 2rem 0;
    margin: 0 auto;
    max-width: 1200px;

    .scroll-indicator-content {
      display: flex;
      flex-direction: column;
      align-items: center;
      gap: 1rem;
      animation: fadeInScale 0.8s ease-out;
    }

    .scroll-indicator-text {
      font-size: 1.1rem;
      font-weight: 500;
      color: var(--color-text-secondary);
      text-align: center;
      letter-spacing: 0.5px;
      transition: color 0.3s ease;
    }

    .scroll-arrow {
      width: 48px;
      height: 48px;
      border-radius: 50%;
      background: linear-gradient(135deg,
      hsla(var(--color-primary-purple-hue), 85%, 65%, 0.15) 0%,
      hsla(var(--color-primary-purple-hue), 85%, 65%, 0.08) 100%);
      border: 2px solid hsla(var(--color-primary-purple-hue), 85%, 65%, 0.3);
      display: flex;
      align-items: center;
      justify-content: center;
      cursor: pointer;
      transition: all 0.4s cubic-bezier(0.4, 0, 0.2, 1);
      animation: bounceDown 2s ease-in-out infinite;

      .arrow-icon {
        color: var(--color-brand);
        transition: transform 0.3s ease;
      }

      &:hover {
        background: linear-gradient(135deg,
        hsla(var(--color-primary-purple-hue), 85%, 65%, 0.25) 0%,
        hsla(var(--color-primary-purple-hue), 85%, 65%, 0.15) 100%);
        border-color: var(--color-brand);
        transform: translateY(4px) scale(1.05);
        box-shadow: 0 8px 20px hsla(var(--color-primary-purple-hue), 85%, 65%, 0.3);

        .arrow-icon {
          transform: translateY(3px);
        }
      }

      &:active {
        transform: translateY(6px) scale(0.98);
      }
    }

    @media (prefers-color-scheme: dark) {
      .scroll-indicator-text {
        color: var(--color-text-secondary);
      }

      .scroll-arrow {
        background: linear-gradient(135deg,
        hsla(var(--color-primary-purple-hue), 85%, 65%, 0.2) 0%,
        hsla(var(--color-primary-purple-hue), 85%, 65%, 0.1) 100%);
        border-color: hsla(var(--color-primary-purple-hue), 85%, 65%, 0.4);

        &:hover {
          background: linear-gradient(135deg,
          hsla(var(--color-primary-purple-hue), 85%, 65%, 0.3) 0%,
          hsla(var(--color-primary-purple-hue), 85%, 65%, 0.18) 100%);
          box-shadow: 0 8px 20px hsla(var(--color-primary-purple-hue), 85%, 65%, 0.4);
        }
      }
    }

    @media (max-width: 768px) {
      padding: 2rem 0 1.5rem 0;

      .scroll-indicator-text {
        font-size: 1rem;
      }

      .scroll-arrow {
        width: 42px;
        height: 42px;
      }
    }
  }
}

// Latest Posts Section
.latest-section {
  padding: 3rem 0;
  background: var(--color-surface);

  .posts-grid {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(320px, 1fr));
    gap: 2rem;
    max-width: 1200px;
    margin: 0 auto;
    padding: 0 2rem;
  }

  .post-card {
    background: var(--color-card-bg, #ffffff);
    border-radius: 12px;
    overflow: hidden;
    cursor: pointer;
    transition: all 0.4s cubic-bezier(0.4, 0, 0.2, 1);
    border: 1px solid var(--color-card-border, #e5e7eb);
    height: 100%;
    display: flex;
    flex-direction: column;
    box-shadow: 0 1px 3px var(--color-card-shadow, rgba(0, 0, 0, 0.08)), 0 1px 2px rgba(0, 0, 0, 0.06);

    &:hover {
      transform: translateY(-8px);
      box-shadow: 0 12px 24px var(--color-hover-shadow, rgba(0, 0, 0, 0.12)), 0 4px 8px var(--color-hover-shadow, rgba(0, 0, 0, 0.08));
      border-color: var(--color-brand);

      .post-image {
        &::before {
          opacity: 1;
        }

        img {
          transform: scale(1.08);
        }
      }

      .post-content h3 {
        color: var(--color-brand);
      }
    }

    .post-image {
      position: relative;
      width: 100%;
      height: 200px;
      overflow: hidden;
      background: var(--color-border);

      &::before {
        content: '';
        position: absolute;
        top: 0;
        left: 0;
        right: 0;
        bottom: 0;
        background: radial-gradient(ellipse at center, transparent 40%, rgba(255, 255, 255, 0.05) 100%);
        opacity: 0;
        transition: opacity 0.3s;
        mix-blend-mode: multiply;
      }

      img {
        width: 100%;
        height: 100%;
        object-fit: cover;
        transition: transform 0.5s cubic-bezier(0.4, 0, 0.2, 1);
      }
    }

    .post-content {
      padding: 1.5rem;
      flex: 1;
      display: flex;
      flex-direction: column;

      h3 {
        font-size: 1.1rem;
        font-weight: 700;
        margin: 0 0 0.5rem 0;
        color: var(--color-text-primary);
        display: -webkit-box;
        -webkit-line-clamp: 2;
        -webkit-box-orient: vertical;
        overflow: hidden;
        line-height: 1.4;
        transition: color 0.3s;
      }

      p {
        color: var(--color-text-secondary);
        font-size: 0.9rem;
        line-height: 1.6;
        margin: 0 0 1rem 0;
        display: -webkit-box;
        -webkit-line-clamp: 2;
        -webkit-box-orient: vertical;
        overflow: hidden;
        flex: 1;
      }

      .post-meta {
        display: flex;
        gap: 1rem;
        font-size: 0.85rem;
        color: var(--color-text-secondary);
        padding-top: 1rem;
        border-top: 1px solid var(--color-border);
        flex-wrap: wrap;

        span {
          display: flex;
          align-items: center;
          gap: 0.3rem;
          font-weight: 500;
        }
      }
    }
  }
}

// Features Section
.features {
  padding: 4rem 2rem;
  background: var(--color-bg);

  .section-header {
    text-align: center;
    display: block;
    margin-bottom: 3rem;
    max-width: 100%;
    padding: 0;
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

// 图片加载状态优化
.featured-image,
.post-image {
  background: linear-gradient(90deg, var(--color-border) 25%, transparent 25%, transparent 50%, var(--color-border) 50%, var(--color-border) 75%, transparent 75%);
  background-size: 60px 60px;
  animation: loading 1.5s infinite;

  img {
    &[loading="lazy"] {
      background: var(--color-border);
    }
  }
}

@keyframes loading {
  0% {
    background-position: 0 0;
  }
  100% {
    background-position: 60px 60px;
  }
}

// 暗黑模式优化 - 创造层次感和视觉节奏
html.dark {
  // 页面背景 - 优化配色更生动
  background: linear-gradient(135deg, #0f1420 0%, #1a1f35 50%, #0f1420 100%);
  background-attachment: fixed;

  // 页面背景装饰效果
  &::before {
    content: '';
    position: fixed;
    top: 0;
    left: 0;
    right: 0;
    bottom: 0;
    background-image: radial-gradient(ellipse at 20% 30%, rgba(99, 102, 241, 0.08) 0%, transparent 50%),
    radial-gradient(ellipse at 80% 70%, rgba(var(--color-primary-purple-rgb), 0.08) 0%, transparent 50%);
    pointer-events: none;
    z-index: -1;
  }

  // 页面主要内容区域
  .home-page {
    background: transparent; // 取消默认背景

    // Hero 区域暗黑模式优化
    .hero {
      background: linear-gradient(135deg, #6366f1 0%, #8b5cf6 50%, #a855f7 100%);
      box-shadow: 0 20px 60px rgba(99, 102, 241, 0.3);

      .hero-grid-bg {
        background-image: linear-gradient(rgba(255, 255, 255, 0.04) 1px, transparent 1px),
        linear-gradient(90deg, rgba(255, 255, 255, 0.04) 1px, transparent 1px);
        opacity: 0.5;
      }

      .hero-code-snippets .code-snippet {
        background: rgba(0, 0, 0, 0.4);
        border: 1px solid rgba(255, 255, 255, 0.2);
        box-shadow: 0 12px 48px rgba(0, 0, 0, 0.6),
        0 0 30px rgba(var(--color-primary-purple-rgb), 0.3),
        inset 0 1px 0 rgba(255, 255, 255, 0.15);
      }

      .hero-title {
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

      .hero-subtitle {
        color: rgba(255, 255, 255, 0.95);
        text-shadow: 0 2px 12px rgba(0, 0, 0, 0.4);
      }

      .hero-description {
        color: rgba(255, 255, 255, 0.85);
        text-shadow: 0 1px 10px rgba(0, 0, 0, 0.3);
      }
    }

    // 卡片基础样式
    .category-card {
      backdrop-filter: blur(10px);

      // 增强背景遮罩在暗黑模式下的可见度
      .category-card-bg {
        background: linear-gradient(135deg,
        hsla(var(--card-hue), 70%, 55%, 0.12) 0%,
        hsla(var(--card-hue), 60%, 65%, 0.08) 100%) !important;
        opacity: 1;
      }

      &:hover {
        background: linear-gradient(135deg, rgba(80, 95, 120, 0.95) 0%, rgba(65, 75, 95, 0.95) 100%);
        border-color: hsl(var(--card-hue), 70%, 60%);
        box-shadow: 0 12px 32px rgba(0, 0, 0, 0.5), 0 0 0 1px hsl(var(--card-hue), 70%, 60%);

        .category-card-bg {
          background: linear-gradient(135deg,
          hsla(var(--card-hue), 70%, 60%, 0.18) 0%,
          hsla(var(--card-hue), 60%, 70%, 0.12) 100%) !important;
        }
      }
    }

    .post-card {
      background: linear-gradient(135deg, rgba(75, 90, 120, 0.9) 0%, rgba(55, 70, 100, 0.9) 100%);
      border-color: rgba(150, 180, 255, 0.2);
      box-shadow: 0 2px 12px rgba(0, 0, 0, 0.4);
      backdrop-filter: blur(10px);

      &:hover {
        background: linear-gradient(135deg, rgba(99, 102, 241, 0.3) 0%, rgba(139, 92, 246, 0.3) 100%),
        linear-gradient(135deg, rgba(100, 120, 160, 0.95) 0%, rgba(80, 100, 140, 0.95) 100%);
        border-color: rgba(var(--color-primary-purple-rgb), 0.6) !important;
        box-shadow: 0 12px 32px rgba(99, 102, 241, 0.3),
        0 0 0 1px rgba(var(--color-primary-purple-rgb), 0.5),
        0 0 20px rgba(var(--color-primary-purple-rgb), 0.15);
      }
    }

    .featured-card {
      background: linear-gradient(135deg, rgba(75, 90, 120, 0.9) 0%, rgba(55, 70, 100, 0.9) 100%);
      border-color: rgba(150, 180, 255, 0.2);
      box-shadow: 0 2px 12px rgba(0, 0, 0, 0.4);
      backdrop-filter: blur(10px);

      &:hover {
        background: linear-gradient(135deg, rgba(99, 102, 241, 0.3) 0%, rgba(139, 92, 246, 0.3) 100%),
        linear-gradient(135deg, rgba(100, 120, 160, 0.95) 0%, rgba(80, 100, 140, 0.95) 100%);
        border-color: rgba(var(--color-primary-purple-rgb), 0.6) !important;
        box-shadow: 0 12px 32px rgba(99, 102, 241, 0.3),
        0 0 0 1px rgba(var(--color-primary-purple-rgb), 0.5),
        0 0 20px rgba(var(--color-primary-purple-rgb), 0.15);
      }
    }

    // Section 背景交替：创造视觉节奏 - 优化配色更生动
    .categories-section {
      background: linear-gradient(180deg, #1a1d26 0%, #222836 100%);
    }

    .featured-section {
      background: linear-gradient(180deg, #252b38 0%, #1f2633 100%);
    }

    .latest-section {
      background: linear-gradient(180deg, #1a1d26 0%, #222836 100%);
    }

    .features {
      background: linear-gradient(180deg, #252b38 0%, #1f2633 100%);
    }

    // Section header 文字优化
    .section-header h2 {
      text-shadow: 0 2px 8px rgba(0, 0, 0, 0.3);
      color: #f0f1f3 !important;
      font-weight: 700;
    }

    // 卡片内部元素优化
    .category-card {
      .category-icon {
        background: linear-gradient(135deg,
        hsla(var(--card-hue), 70%, 60%, 0.25) 0%,
        hsla(var(--card-hue), 60%, 70%, 0.15) 100%);
        box-shadow: 0 2px 12px hsla(var(--card-hue), 70%, 60%, 0.2);
        color: hsl(var(--card-hue), 75%, 70%) !important;
      }

      .category-info {
        h3 {
          color: #f3f4f6 !important;
          font-weight: 600;
          text-shadow: 0 2px 4px rgba(0, 0, 0, 0.6),
          0 1px 2px rgba(0, 0, 0, 0.4);
        }

        .post-count {
          color: #e5e7eb !important;
          font-weight: 600;
          text-shadow: 0 1px 3px rgba(0, 0, 0, 0.5),
          0 1px 2px rgba(0, 0, 0, 0.3);
        }
      }

      .category-badge {
        background: linear-gradient(135deg,
        hsla(var(--card-hue), 70%, 50%, 0.4) 0%,
        hsla(var(--card-hue), 60%, 55%, 0.3) 100%) !important;
        border-left: 3px solid hsla(var(--card-hue), 75%, 70%, 0.9) !important;
        backdrop-filter: blur(8px);
        box-shadow: 0 2px 8px hsla(var(--card-hue), 70%, 50%, 0.3),
        inset 0 1px 0 hsla(var(--card-hue), 80%, 80%, 0.15);
        opacity: 1 !important;

        // 图标和文字使用纯白色
        svg,
        span {
          color: #ffffff !important;
          text-shadow: 0 1px 2px rgba(0, 0, 0, 0.5);
        }
      }

      &:hover {
        .category-icon {
          color: hsl(var(--card-hue), 80%, 75%) !important;
        }

        .category-info {
          h3 {
            color: #ffffff !important;
            text-shadow: 0 2px 6px rgba(0, 0, 0, 0.7),
            0 1px 3px rgba(0, 0, 0, 0.5);
          }

          .post-count {
            color: #f3f4f6 !important;
            text-shadow: 0 1px 4px rgba(0, 0, 0, 0.6),
            0 1px 2px rgba(0, 0, 0, 0.4);
          }
        }

        .category-badge {
          background: linear-gradient(135deg,
          hsla(var(--card-hue), 75%, 55%, 0.5) 0%,
          hsla(var(--card-hue), 65%, 60%, 0.4) 100%) !important;
          border-left-color: hsla(var(--card-hue), 80%, 75%, 1) !important;
          box-shadow: 0 4px 12px hsla(var(--card-hue), 70%, 55%, 0.4),
          inset 0 1px 0 hsla(var(--card-hue), 85%, 85%, 0.2);
          transform: scale(1.05);
        }
      }
    }
  }

  // Feature卡片在暗黑模式下的优化
  .feature-card {
    background: linear-gradient(135deg, rgba(75, 90, 120, 0.9) 0%, rgba(55, 70, 100, 0.9) 100%);
    border-color: rgba(150, 180, 255, 0.2);
    box-shadow: 0 2px 12px rgba(0, 0, 0, 0.4);
    backdrop-filter: blur(10px);

    h3 {
      color: #e5e7eb !important;
    }

    p {
      color: #9ca3af !important;
    }

    &:hover {
      background: linear-gradient(135deg, rgba(99, 102, 241, 0.3) 0%, rgba(139, 92, 246, 0.3) 100%),
      linear-gradient(135deg, rgba(100, 120, 160, 0.95) 0%, rgba(80, 100, 140, 0.95) 100%);
      border-color: rgba(var(--color-primary-purple-rgb), 0.6) !important;
      box-shadow: 0 12px 32px rgba(99, 102, 241, 0.3),
      0 0 0 1px rgba(var(--color-primary-purple-rgb), 0.5),
      0 0 20px rgba(var(--color-primary-purple-rgb), 0.15);

      h3 {
        color: #f9fafb !important;
      }
    }
  }

  // 文章卡片文字优化
  .post-card,
  .featured-card {
    .post-content,
    .featured-content {
      h3 {
        color: #e5e7eb !important;
      }

      p {
        color: #9ca3af !important;
      }

      .post-meta,
      .featured-meta {
        color: #9ca3af !important;
      }
    }

    &:hover {
      .post-content h3,
      .featured-content h3 {
        color: #f9fafb !important;
      }
    }
  }

  // 热门标签区域暗黑模式优化
  .popular-tags-container {
    .tags-title {
      color: #f0f1f3 !important;
    }

    .tag-item {
      background: linear-gradient(135deg,
      rgba(75, 90, 120, 0.5) 0%,
      rgba(55, 70, 100, 0.5) 100%);
      border-color: var(--tag-color);

      &:hover {
        background: var(--tag-color);
        box-shadow: 0 4px 16px rgba(0, 0, 0, 0.3);
      }
    }

    .tag-label {
      color: var(--tag-color);
      filter: brightness(1.5);
    }
  }
}

// Responsive Design
@media (max-width: 1024px) {
  .section-header {
    padding: 0 1.5rem;

    h2 {
      font-size: 1.875rem;
    }
  }

  .categories-grid {
    grid-template-columns: repeat(2, 1fr);
    gap: 1.5rem;
    padding: 0 1.5rem;
  }

  .featured-grid,
  .posts-grid,
  .features-grid {
    padding: 0 1.5rem;
    gap: 1.5rem;
  }

  .featured-card {
    .featured-image {
      height: 200px;
    }

    .featured-content {
      padding: 20px;

      h3 {
        font-size: 18px;
      }

      p {
        font-size: 14px;
        -webkit-line-clamp: 3;
      }
    }
  }

  .hero {
    padding: 3rem 1.5rem 2.5rem;
    min-height: 300px;

    .hero-title {
      font-size: 3rem;
      letter-spacing: -1.2px;
      margin-bottom: 1.25rem;
    }

    .hero-subtitle {
      font-size: 1.35rem;
      margin-bottom: 1rem;
    }

    .hero-description {
      font-size: 1rem;
      margin-bottom: 2rem;
    }

    .hero-actions {
      gap: 1rem;

      :deep(.n-button) {
        padding: 12px 32px;
        font-size: 15px;
      }
    }

    // 隐藏代码片段在小屏幕上
    .hero-code-snippets {
      display: none;
    }
  }

  .category-card {
    padding: 1.5rem 1.25rem;
    min-height: 180px;

    .category-icon {
      width: 56px;
      height: 56px;
    }
  }

  .popular-tags-container {
    padding: 0 1.5rem;
  }
}

@media (max-width: 768px) {
  // 隐藏桌面网格，显示手机走马灯
  .categories-grid.desktop-grid {
    display: none !important;
  }

  .categories-carousel.mobile-carousel {
    display: block !important;
    padding: 0 1rem;
  }

  .hero {
    padding: 2.5rem 1.5rem 2rem;
    min-height: 280px;

    .hero-title {
      font-size: 2.2rem;
      margin-bottom: 1rem;
      letter-spacing: -0.8px;
    }

    .hero-subtitle {
      font-size: 1.15rem;
      margin-bottom: 0.875rem;
    }

    .hero-description {
      font-size: 0.95rem;
      margin-bottom: 1.75rem;
      line-height: 1.6;
    }

    // 隐藏装饰元素在移动端
    .hero-code-snippets,
    .hero-grid-bg,
    .hero-animated-shapes {
      display: none;
    }

    .hero-actions {
      flex-direction: column;
      gap: 0.875rem;
      width: 100%;
      max-width: 400px;
      margin: 0 auto;

      :deep(.n-button) {
        width: 100%;
        padding: 12px 24px;
        font-size: 15px;
      }
    }
  }

  .section-header {
    flex-direction: column;
    align-items: flex-start;
    gap: 1rem;
    padding: 0 1.5rem;

    h2 {
      font-size: 1.75rem;
    }

    :deep(.n-button) {
      font-size: 14px;
      padding: 8px 16px;
    }
  }

  .categories-section {
    padding: 3rem 0;
  }

  .featured-section,
  .latest-section,
  .features {
    padding: 3rem 0;
  }

  // Scroll Indicator responsive - hide on mobile
  .featured-section {
    .scroll-indicator {
      display: none; // 移动端隐藏滚动指示器
    }
  }

  .categories-grid {
    grid-template-columns: repeat(2, 1fr);
    gap: 1rem;
    padding: 0 1.5rem;
  }

  .featured-grid,
  .posts-grid {
    grid-template-columns: 1fr;
    gap: 1.5rem;
    padding: 0 1.5rem;
  }

  .featured-card {
    .featured-image {
      height: 200px;
    }

    .featured-content {
      padding: 20px;

      h3 {
        font-size: 18px;
        margin-bottom: 10px;
      }

      p {
        font-size: 14px;
        -webkit-line-clamp: 2;
        margin-bottom: 12px;
      }

      .featured-meta {
        padding-top: 12px;
        gap: 14px;
        font-size: 12px;

        span {
          :deep(svg) {
            width: 14px;
            height: 14px;
          }
        }
      }
    }
  }

  .popular-tags-container {
    padding: 0 1.5rem;

    .tags-title {
      font-size: 1.125rem;
      margin-bottom: 1.25rem;
    }

    .tags-grid {
      gap: 0.75rem;
    }

    .tag-item {
      padding: 0.4rem 0.875rem;
      font-size: 0.85rem;
    }
  }

  .category-card {
    padding: 1.5rem;

    .category-icon {
      width: 50px;
      height: 50px;
      font-size: 1.75rem;
    }

    .category-info h3 {
      font-size: 1rem;
    }

    .post-count {
      font-size: 0.8rem;
    }

    .update-badge {
      font-size: 0.7rem;
      padding: 0.3rem 0.7rem;
    }
  }

  .featured-card,
  .post-card {
    .featured-image,
    .post-image {
      height: 160px;
    }

    .featured-content,
    .post-content {
      padding: 1rem;

      h3 {
        font-size: 1rem;
        margin-bottom: 0.5rem;
        line-height: 1.4;
      }

      p {
        font-size: 0.85rem;
        line-height: 1.5;
      }

      .featured-meta,
      .post-meta {
        font-size: 0.75rem;
        gap: 0.75rem;
        margin-top: 0.75rem;
      }

      .category-badge {
        font-size: 0.7rem;
        padding: 0.25rem 0.625rem;
      }
    }
  }


  .features-grid {
    grid-template-columns: 1fr;
    gap: 1rem;
  }

  .feature-card {
    padding: 1.5rem 1rem;

    &:hover {
      transform: translateY(-4px);
    }

    .feature-icon {
      width: 60px;
      height: 60px;
      font-size: 2rem;
      margin-bottom: 1rem;
    }

    h3 {
      font-size: 1rem;
      margin-bottom: 0.625rem;
    }

    p {
      font-size: 0.875rem;
      line-height: 1.5;
    }
  }
}
</style>

