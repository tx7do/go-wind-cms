<script setup lang="ts">
import {definePage} from 'unplugin-vue-router/runtime'
import {useRouter} from 'vue-router'
import {ref, onMounted, computed} from 'vue'
import {usePostStore, useCategoryStore} from '@/stores/modules/app'
import {$t, i18n} from '@/locales'
import {XIcon} from "@/plugins/xicon";

definePage({
  name: 'home',
  meta: {
    level: 1,
  },
})

const router = useRouter()
const postStore = usePostStore()
const categoryStore = useCategoryStore()

const latestPosts = ref<any[]>([])
const featuredPosts = ref<any[]>([])
const categories = ref<any[]>([])
const loading = ref(false)

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

function getPostTitle(post: any) {
  return post.translations?.[0]?.title || $t('page.home.untitled')
}

function getPostSummary(post: any) {
  return post.translations?.[0]?.summary || ''
}

function getPostThumbnail(post: any) {
  return post.translations?.[0]?.thumbnail || '/placeholder.jpg'
}

function getCategoryName(category: any) {
  return category.translations?.[0]?.name || $t('page.home.category_default')
}

function formatDate(dateString: string) {
  if (!dateString) return ''
  return new Date(dateString).toLocaleDateString()
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
  ])
  loading.value = false
})
</script>

<template>
  <div class="home-page">
    <!-- Hero Section -->
    <section class="hero">
      <!-- Animated Background Elements -->
      <div class="hero-bg-wrapper">
        <div class="hero-gradient-bg"></div>
        <div class="hero-animated-shapes">
          <div class="shape shape-1"></div>
          <div class="shape shape-2"></div>
          <div class="shape shape-3"></div>
        </div>
        <svg class="hero-waves" viewBox="0 0 1440 320" preserveAspectRatio="none">
          <path fill="rgba(255,255,255,0.1)" d="M0,96L48,112C96,128,192,160,288,160C384,160,480,128,576,122.7C672,117,768,139,864,138.7C960,139,1056,117,1152,101.3C1248,85,1344,75,1392,69.3L1440,64L1440,320L1392,320C1344,320,1248,320,1152,320C1056,320,960,320,864,320C768,320,672,320,576,320C480,320,384,320,288,320C192,320,96,320,48,320L0,320Z"></path>
        </svg>
      </div>

      <div class="hero-content">
        <h1 class="hero-title">{{ $t('authentication.login.brand_title') }}</h1>
        <p class="hero-subtitle">{{ $t('authentication.login.brand_subtitle') }}</p>
        <p class="hero-tagline">基于 Go + Vue3 的高性能、多语言 CMS 系统</p>
        <div class="hero-actions">
          <n-button type="primary" size="large" class="btn-primary" @click="router.push('/post')">
            <template #icon>
              <XIcon name="carbon:document-export" :size="20"/>
            </template>
            {{ $t('page.home.browse_posts') }}
          </n-button>
          <n-button type="default" size="large" class="btn-secondary" @click="router.push('/about')">
            <template #icon>
              <XIcon name="carbon:information" :size="20"/>
            </template>
            {{ $t('page.home.learn_more') }}
          </n-button>
        </div>
      </div>
    </section>

    <!-- Categories Section -->
    <section class="categories-section">
      <div class="section-header">
        <h2>{{ $t('page.home.categories') }}</h2>
        <n-button text type="primary" @click="router.push('/category')">
          {{ $t('page.home.view_all') }} →
        </n-button>
      </div>
      <n-spin :show="loading">
        <div class="categories-grid">
          <div
            v-for="category in categories"
            :key="category.id"
            class="category-card"
            :style="{ '--card-hue': (category.id % 6) * 60 }"
            @click="handleViewCategory(category.id)"
          >
            <!-- Card Background Gradient -->
            <div class="category-card-bg"></div>

            <!-- Content -->
            <div class="category-card-content">
              <div class="category-card-header">
                <div class="category-icon">
                  <XIcon :name="category.icon?.includes(':') ? category.icon : `carbon:${category.icon || 'folder'}`" :size="48"/>
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
                <span>3 天前更新</span>
              </div>
            </div>
          </div>
        </div>
      </n-spin>
    </section>

    <!-- Featured Posts Section -->
    <section v-if="featuredPosts.length > 0" class="featured-section">
      <div class="section-header">
        <h2>{{ $t('page.home.featured_posts') }}</h2>
      </div>
      <n-spin :show="loading">
        <div class="featured-grid">
          <div
            v-for="post in featuredPosts"
            :key="post.id"
            class="featured-card"
            @click="handleViewPost(post.id)"
          >
            <div class="featured-image">
              <img :src="getPostThumbnail(post)" :alt="getPostTitle(post)"/>
            </div>
            <div class="featured-content">
              <h3>{{ getPostTitle(post) }}</h3>
              <p>{{ getPostSummary(post) }}</p>
              <div class="featured-meta">
                <span>{{ post.authorName }}</span>
                <span>{{ formatDate(post.createdAt) }}</span>
              </div>
            </div>
          </div>
        </div>
      </n-spin>
    </section>

    <!-- Latest Posts Section -->
    <section class="latest-section">
      <div class="section-header">
        <h2>{{ $t('page.home.latest_posts') }}</h2>
        <n-button text type="primary" @click="router.push('/post')">
          {{ $t('page.home.view_all') }} →
        </n-button>
      </div>
      <n-spin :show="loading">
        <div class="posts-grid">
          <div
            v-for="post in latestPosts"
            :key="post.id"
            class="post-card"
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
      </n-spin>
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

// Hero Section
.hero {
  padding: 6rem 2rem;
  text-align: center;
  max-width: 100%;
  background: linear-gradient(135deg, var(--color-brand) 0%, #764ba2 100%);
  color: white;
  position: relative;
  overflow: hidden;
  min-height: 600px;
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

    .hero-gradient-bg {
      position: absolute;
      top: 0;
      left: 0;
      right: 0;
      bottom: 0;
      background:
        radial-gradient(ellipse at 20% 50%, rgba(255, 255, 255, 0.1) 0%, transparent 50%),
        radial-gradient(ellipse at 80% 80%, rgba(255, 255, 255, 0.08) 0%, transparent 50%),
        linear-gradient(135deg, var(--color-brand) 0%, #764ba2 100%);
      animation: gradientShift 15s ease-in-out infinite;
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
    z-index: 1;
    max-width: 1200px;
    margin: 0 auto;
  }

  .hero-title {
    font-size: 3.5rem;
    font-weight: 800;
    margin-bottom: 1rem;
    line-height: 1.2;
    letter-spacing: -1px;
    text-shadow: 0 2px 12px rgba(0, 0, 0, 0.25);
    animation: slideDown 0.8s cubic-bezier(0.34, 1.56, 0.64, 1);
  }

  .hero-subtitle {
    font-size: 1.5rem;
    margin-bottom: 0.75rem;
    max-width: 600px;
    margin-left: auto;
    margin-right: auto;
    opacity: 0.95;
    font-weight: 500;
    text-shadow: 0 1px 8px rgba(0, 0, 0, 0.2);
    animation: fadeInUp 0.8s ease-out 0.2s both;
  }

  .hero-tagline {
    font-size: 1.1rem;
    margin-bottom: 2.5rem;
    max-width: 700px;
    margin-left: auto;
    margin-right: auto;
    opacity: 0.9;
    font-weight: 400;
    text-shadow: 0 1px 6px rgba(0, 0, 0, 0.2);
    animation: fadeInUp 0.8s ease-out 0.4s both;
  }

  .hero-actions {
    display: flex;
    gap: 1.5rem;
    justify-content: center;
    flex-wrap: wrap;
    animation: fadeInUp 0.8s ease-out 0.6s both;

    .btn-primary {
      :deep(.n-button) {
        background: linear-gradient(135deg, #10b981 0%, #059669 100%) !important;
        color: white !important;
        font-size: 1rem;
        font-weight: 700;
        padding: 14px 36px !important;
        border-radius: 50px !important;
        border: none !important;
        box-shadow: 0 8px 24px rgba(16, 185, 129, 0.4), 0 1px 3px rgba(0, 0, 0, 0.1);
        transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
        text-shadow: 0 1px 3px rgba(0, 0, 0, 0.15);

        &:hover {
          transform: translateY(-3px);
          box-shadow: 0 12px 32px rgba(16, 185, 129, 0.5), 0 4px 8px rgba(0, 0, 0, 0.15);
          background: linear-gradient(135deg, #059669 0%, #047857 100%) !important;
        }

        &:active {
          transform: translateY(-1px);
          box-shadow: 0 4px 12px rgba(16, 185, 129, 0.4);
        }
      }
    }

    .btn-secondary {
      :deep(.n-button) {
        background: rgba(255, 255, 255, 0.15) !important;
        color: white !important;
        font-size: 1rem;
        font-weight: 700;
        padding: 14px 36px !important;
        border-radius: 50px !important;
        border: 2px solid rgba(255, 255, 255, 0.4) !important;
        backdrop-filter: blur(10px);
        box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1), 0 0 0 1px rgba(255, 255, 255, 0.2) inset;
        transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);

        &:hover {
          background: rgba(255, 255, 255, 0.25) !important;
          border-color: rgba(255, 255, 255, 0.6) !important;
          transform: translateY(-3px);
          box-shadow: 0 12px 32px rgba(0, 0, 0, 0.15), 0 0 0 1px rgba(255, 255, 255, 0.3) inset;
        }

        &:active {
          transform: translateY(-1px);
          background: rgba(255, 255, 255, 0.2) !important;
        }
      }
    }
  }
}

@keyframes gradientShift {
  0%, 100% {
    background:
      radial-gradient(ellipse at 20% 50%, rgba(255, 255, 255, 0.1) 0%, transparent 50%),
      radial-gradient(ellipse at 80% 80%, rgba(255, 255, 255, 0.08) 0%, transparent 50%),
      linear-gradient(135deg, var(--color-brand) 0%, #764ba2 100%);
  }
  50% {
    background:
      radial-gradient(ellipse at 80% 50%, rgba(255, 255, 255, 0.1) 0%, transparent 50%),
      radial-gradient(ellipse at 20% 80%, rgba(255, 255, 255, 0.08) 0%, transparent 50%),
      linear-gradient(135deg, #764ba2 0%, var(--color-brand) 100%);
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
  padding: 4rem 0;
  background: var(--color-surface);

  .categories-grid {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(280px, 1fr));
    gap: 2rem;
    max-width: 1200px;
    margin: 0 auto;
    padding: 0 2rem;
  }

  .category-card {
    display: flex;
    flex-direction: column;
    padding: 1.75rem;
    background: var(--color-card-bg, linear-gradient(135deg, #ffffff 0%, #f9fafb 100%));
    border: 1px solid var(--color-card-border, #e5e7eb);
    border-radius: 16px;
    cursor: pointer;
    transition: all 0.4s cubic-bezier(0.4, 0, 0.2, 1);
    position: relative;
    overflow: visible;
    box-shadow: 0 2px 8px var(--color-card-shadow, rgba(0, 0, 0, 0.06)), 0 1px 3px rgba(0, 0, 0, 0.04);
    min-height: 240px;
    height: 100%;

    &::before {
      content: '';
      position: absolute;
      top: 0;
      left: 0;
      right: 0;
      bottom: 0;
      background: linear-gradient(
        135deg,
        hsl(var(--card-hue), 70%, 60%) 0%,
        hsl(var(--card-hue), 60%, 70%) 100%
      );
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
      transform: translateY(-8px);
      box-shadow: 0 16px 32px rgba(0, 0, 0, 0.12), 0 4px 12px rgba(0, 0, 0, 0.08);
      border-color: hsl(var(--card-hue), 70%, 60%);

      &::before {
        opacity: 0.15;
      }

      &::after {
        left: 100%;
      }

      .category-icon {
        transform: scale(1.15) rotate(8deg);
        color: white;
        background: linear-gradient(
          135deg,
          hsl(var(--card-hue), 70%, 60%) 0%,
          hsl(var(--card-hue), 60%, 70%) 100%
        );
        box-shadow: 0 8px 16px hsl(var(--card-hue), 70%, 60%, 0.3);
      }

      .category-info h3 {
        color: hsl(var(--card-hue), 70%, 50%);
      }
    }

    .category-card-bg {
      position: absolute;
      top: 0;
      left: 0;
      right: 0;
      bottom: 0;
      background: linear-gradient(
        135deg,
        hsla(var(--card-hue), 70%, 60%, 0.05) 0%,
        hsla(var(--card-hue), 60%, 70%, 0.05) 100%
      );
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
      background: linear-gradient(
        135deg,
        hsla(var(--card-hue), 70%, 60%, 0.15) 0%,
        hsla(var(--card-hue), 60%, 70%, 0.1) 100%
      );
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
        transition: color 0.3s;
        word-break: break-word;
      }

      .post-count {
        font-size: 0.95rem;
        color: var(--color-text-secondary);
        font-weight: 600;
        letter-spacing: 0.2px;
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
      background: linear-gradient(
        135deg,
        hsla(var(--card-hue), 70%, 60%, 0.1) 0%,
        hsla(var(--card-hue), 60%, 70%, 0.08) 100%
      );
      border-radius: 8px;
      border-left: 3px solid hsl(var(--card-hue), 70%, 60%);
      color: hsl(var(--card-hue), 70%, 50%);
      font-size: 0.8rem;
      font-weight: 600;
      width: fit-content;
      white-space: nowrap;
      transition: all 0.3s ease;
      opacity: 0.85;
    }

    &:hover .category-badge {
      opacity: 1;
      transform: scale(1.05);
      box-shadow: 0 4px 8px hsla(var(--card-hue), 70%, 60%, 0.2);
    }
  }
}

// Featured Posts Section
.featured-section {
  padding: 4rem 0;
  background: var(--color-bg);

  .featured-grid {
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(320px, 1fr));
    gap: 2.5rem;
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

    &:hover {
      transform: translateY(-12px);
      box-shadow: 0 20px 40px var(--color-hover-shadow, rgba(0, 0, 0, 0.12)), 0 8px 16px var(--color-hover-shadow, rgba(0, 0, 0, 0.08));
      border-color: var(--color-brand);

      .featured-image {
        &::after {
          opacity: 1;
        }

        img {
          transform: scale(1.1);
        }
      }

      .featured-content h3 {
        color: var(--color-brand);
      }
    }

    .featured-image {
      position: relative;
      width: 100%;
      height: 240px;
      overflow: hidden;
      background: var(--color-bg);

      &::after {
        content: '';
        position: absolute;
        top: 0;
        left: 0;
        right: 0;
        bottom: 0;
        background: radial-gradient(ellipse at center, transparent 30%, rgba(102, 126, 234, 0.1) 100%);
        opacity: 0;
        transition: opacity 0.3s;
        mix-blend-mode: screen;
      }

      img {
        width: 100%;
        height: 100%;
        object-fit: cover;
        transition: transform 0.5s cubic-bezier(0.4, 0, 0.2, 1);
      }
    }

    .featured-content {
      padding: 2rem;

      h3 {
        font-size: 1.35rem;
        font-weight: 700;
        margin: 0 0 0.75rem 0;
        color: var(--color-text-primary);
        line-height: 1.4;
        transition: color 0.3s;
      }

      p {
        color: var(--color-text-secondary);
        line-height: 1.8;
        margin: 0 0 1.25rem 0;
        display: -webkit-box;
        -webkit-line-clamp: 2;
        -webkit-box-orient: vertical;
        overflow: hidden;
        font-size: 0.95rem;
      }

      .featured-meta {
        display: flex;
        gap: 1.5rem;
        font-size: 0.85rem;
        color: var(--color-text-secondary);
        padding-top: 1rem;
        border-top: 1px solid var(--color-border);
        font-weight: 500;
      }
    }
  }
}

// Latest Posts Section
.latest-section {
  padding: 4rem 0;
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
  // 定义暗黑模式的CSS变量 - 增强背景颜色深度
  --color-card-bg: linear-gradient(135deg, rgba(65, 75, 95, 0.85) 0%, rgba(48, 56, 72, 0.85) 100%);
  --color-card-border: rgba(255, 255, 255, 0.12);
  --color-card-shadow: rgba(0, 0, 0, 0.5);
  --color-hover-shadow: rgba(0, 0, 0, 0.6);

  // 卡片基础样式
  .category-card {
    backdrop-filter: blur(10px);

    // 增强背景遮罩在暗黑模式下的可见度
    .category-card-bg {
      background: linear-gradient(
        135deg,
        hsla(var(--card-hue), 70%, 55%, 0.12) 0%,
        hsla(var(--card-hue), 60%, 65%, 0.08) 100%
      ) !important;
      opacity: 1;
    }

    &:hover {
      background: linear-gradient(135deg, rgba(80, 95, 120, 0.95) 0%, rgba(65, 75, 95, 0.95) 100%);
      border-color: hsl(var(--card-hue), 70%, 60%);
      box-shadow: 0 12px 32px rgba(0, 0, 0, 0.5), 0 0 0 1px hsl(var(--card-hue), 70%, 60%);

      .category-card-bg {
        background: linear-gradient(
          135deg,
          hsla(var(--card-hue), 70%, 60%, 0.18) 0%,
          hsla(var(--card-hue), 60%, 70%, 0.12) 100%
        ) !important;
      }
    }
  }

  .post-card {
    background: linear-gradient(135deg, rgba(65, 75, 95, 0.85) 0%, rgba(48, 56, 72, 0.85) 100%);
    border-color: rgba(255, 255, 255, 0.12);
    box-shadow: 0 2px 12px rgba(0, 0, 0, 0.5);
    backdrop-filter: blur(10px);

    &:hover {
      background: linear-gradient(135deg, rgba(80, 95, 120, 0.95) 0%, rgba(65, 75, 95, 0.95) 100%);
      border-color: var(--color-brand);
      box-shadow: 0 12px 32px rgba(0, 0, 0, 0.6), 0 0 0 1px var(--color-brand);
    }
  }

  .featured-card {
    background: linear-gradient(135deg, rgba(65, 75, 95, 0.85) 0%, rgba(48, 56, 72, 0.85) 100%);
    border-color: rgba(255, 255, 255, 0.12);
    box-shadow: 0 2px 12px rgba(0, 0, 0, 0.5);
    backdrop-filter: blur(10px);

    &:hover {
      background: linear-gradient(135deg, rgba(80, 95, 120, 0.95) 0%, rgba(65, 75, 95, 0.95) 100%);
      border-color: var(--color-brand);
      box-shadow: 0 12px 32px rgba(0, 0, 0, 0.6), 0 0 0 1px var(--color-brand);
    }
  }

  // Section 背景交替：创造视觉节奏
  .categories-section {
    background: #18191c;
  }

  .featured-section {
    background: #1e2024;
  }

  .latest-section {
    background: #18191c;
  }

  .features {
    background: #1e2024;
  }

  // Section header 文字优化
  .section-header h2 {
    text-shadow: 0 2px 8px rgba(0, 0, 0, 0.3);
  }

  // 卡片内部元素优化
  .category-card {
    .category-icon {
      background: linear-gradient(
        135deg,
        hsla(var(--card-hue), 70%, 60%, 0.25) 0%,
        hsla(var(--card-hue), 60%, 70%, 0.15) 100%
      );
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
      background: linear-gradient(
        135deg,
        hsla(var(--card-hue), 70%, 50%, 0.4) 0%,
        hsla(var(--card-hue), 60%, 55%, 0.3) 100%
      ) !important;
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
        background: linear-gradient(
          135deg,
          hsla(var(--card-hue), 75%, 55%, 0.5) 0%,
          hsla(var(--card-hue), 65%, 60%, 0.4) 100%
        ) !important;
        border-left-color: hsla(var(--card-hue), 80%, 75%, 1) !important;
        box-shadow: 0 4px 12px hsla(var(--card-hue), 70%, 55%, 0.4),
                    inset 0 1px 0 hsla(var(--card-hue), 85%, 85%, 0.2);
        transform: scale(1.05);
      }
    }
  }

  // Feature卡片在暗黑模式下的优化
  .feature-card {
    background: linear-gradient(135deg, rgba(65, 75, 95, 0.85) 0%, rgba(48, 56, 72, 0.85) 100%);
    border-color: rgba(255, 255, 255, 0.12);
    box-shadow: 0 2px 12px rgba(0, 0, 0, 0.5);
    backdrop-filter: blur(10px);

    h3 {
      color: #e5e7eb !important;
    }

    p {
      color: #9ca3af !important;
    }

    &:hover {
      background: linear-gradient(135deg, rgba(80, 95, 120, 0.95) 0%, rgba(65, 75, 95, 0.95) 100%);
      box-shadow: 0 12px 32px rgba(0, 0, 0, 0.6);

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
}

// Responsive Design
@media (max-width: 1024px) {
  .section-header {
    padding: 0 1.5rem;
  }

  .categories-grid,
  .featured-grid,
  .posts-grid,
  .features-grid {
    padding: 0 1.5rem;
  }

  .hero-title {
    font-size: 2.8rem;
  }

  .hero-subtitle {
    font-size: 1.25rem;
  }
}

@media (max-width: 768px) {
  .hero {
    padding: 4rem 1.5rem;

    .hero-title {
      font-size: 2.2rem;
      margin-bottom: 0.75rem;
    }

    .hero-subtitle {
      font-size: 1rem;
      margin-bottom: 1.5rem;
    }

    .hero-actions {
      flex-direction: column;
      gap: 1rem;

      :deep(.n-button) {
        width: 100%;
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
  }

  .categories-grid,
  .featured-grid,
  .posts-grid {
    grid-template-columns: 1fr;
    gap: 1.5rem;
    padding: 0 1.5rem;
  }

  .category-card {
    padding: 1.5rem;

    .category-icon {
      width: 50px;
      height: 50px;
    }
  }

  .featured-card,
  .post-card {
    .featured-image,
    .post-image {
      height: 180px;
    }
  }

  .features-grid {
    grid-template-columns: 1fr;
    gap: 1.5rem;
    padding: 0 1.5rem;
  }

  .feature-card {
    padding: 2rem 1.5rem;

    .feature-icon {
      width: 70px;
      height: 70px;
      font-size: 2.5rem;
    }

    h3 {
      font-size: 1.1rem;
    }
  }
}

@media (max-width: 480px) {
  .hero {
    padding: 3rem 1rem;

    .hero-title {
      font-size: 1.75rem;
    }

    .hero-subtitle {
      font-size: 0.95rem;
    }
  }

  .section-header {
    padding: 0 1rem;

    h2 {
      font-size: 1.5rem;
    }
  }

  .categories-grid,
  .featured-grid,
  .posts-grid,
  .features-grid {
    gap: 1rem;
    padding: 0 1rem;
  }

  .category-card {
    padding: 1rem;
    gap: 1rem;

    .category-icon {
      width: 45px;
      height: 45px;
      font-size: 1.5rem;
    }

    .category-info {
      h3 {
        font-size: 1rem;
      }

      .post-count {
        font-size: 0.8rem;
      }
    }
  }

  .featured-card,
  .post-card {
    .featured-image,
    .post-image {
      height: 150px;
    }

    .featured-content,
    .post-content {
      padding: 1rem;

      h3 {
        font-size: 1rem;
      }

      p {
        font-size: 0.85rem;
      }

      .featured-meta,
      .post-meta {
        font-size: 0.75rem;
        gap: 0.75rem;
      }
    }
  }

  .features-grid {
    grid-template-columns: 1fr;
  }

  .feature-card {
    padding: 1.5rem 1rem;

    .feature-icon {
      width: 60px;
      height: 60px;
      font-size: 2rem;
      margin-bottom: 1rem;
    }

    h3 {
      font-size: 1rem;
    }

    p {
      font-size: 0.9rem;
    }
  }
}
</style>

