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
      <div class="hero-content">
        <h1 class="hero-title">{{ $t('authentication.login.brand_title') }}</h1>
        <p class="hero-subtitle">{{ $t('authentication.login.brand_subtitle') }}</p>
        <div class="hero-actions">
          <n-button type="primary" size="large" @click="router.push('/post')">
            {{ $t('page.home.browse_posts') }}
          </n-button>
          <n-button type="default" size="large" class="learn-more-btn" @click="router.push('/about')">
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
            @click="handleViewCategory(category.id)"
          >
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

  &::before {
    content: '';
    position: absolute;
    top: 0;
    left: 0;
    right: 0;
    bottom: 0;
    background: url('data:image/svg+xml,<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 1440 320"><path fill="rgba(255,255,255,0.05)" d="M0,96L48,112C96,128,192,160,288,160C384,160,480,128,576,122.7C672,117,768,139,864,138.7C960,139,1056,117,1152,101.3C1248,85,1344,75,1392,69.3L1440,64L1440,320L1392,320C1344,320,1248,320,1152,320C1056,320,960,320,864,320C768,320,672,320,576,320C480,320,384,320,288,320C192,320,96,320,48,320L0,320Z"></path></svg>') no-repeat bottom;
    background-size: cover;
    opacity: 0.3;
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
    text-shadow: 0 2px 8px rgba(0, 0, 0, 0.2);
  }

  .hero-subtitle {
    font-size: 1.5rem;
    margin-bottom: 2.5rem;
    max-width: 600px;
    margin-left: auto;
    margin-right: auto;
    opacity: 0.95;
    font-weight: 500;
  }

  .hero-actions {
    display: flex;
    gap: 1.5rem;
    justify-content: center;
    flex-wrap: wrap;

    :deep(.n-button) {
      font-size: 1rem;
      font-weight: 600;
      padding: 12px 32px;
      border-radius: 50px;
      transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
    }

    :deep(.n-button--primary) {
      background: white !important;
      color: var(--color-brand) !important;
      box-shadow: 0 8px 24px rgba(0, 0, 0, 0.2);

      &:hover {
        transform: translateY(-2px);
        box-shadow: 0 12px 32px rgba(0, 0, 0, 0.3);
      }
    }

    :deep(.n-button--default) {
      background: transparent !important;
      color: white !important;
      border: 2px solid white !important;

      &:hover {
        background: rgba(255, 255, 255, 0.15) !important;
        transform: translateY(-2px);
      }
    }

    :deep(.learn-more-btn) {
      color: white !important;
      border-color: white !important;

      &:hover {
        background: rgba(255, 255, 255, 0.15) !important;
        color: white !important;
      }

      &:active {
        background: rgba(255, 255, 255, 0.2) !important;
        color: white !important;
      }
    }
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
    align-items: flex-start;
    gap: 1.5rem;
    padding: 2rem;
    background: var(--color-bg);
    border: 1px solid var(--color-border);
    border-radius: 12px;
    cursor: pointer;
    transition: all 0.4s cubic-bezier(0.4, 0, 0.2, 1);
    position: relative;
    overflow: hidden;

    &::before {
      content: '';
      position: absolute;
      top: 0;
      left: -100%;
      width: 100%;
      height: 100%;
      background: linear-gradient(90deg, transparent, rgba(102, 126, 234, 0.1), transparent);
      transition: left 0.6s;
    }

    &:hover {
      transform: translateY(-6px);
      box-shadow: 0 12px 32px rgba(0, 0, 0, 0.12);
      border-color: var(--color-brand);

      &::before {
        left: 100%;
      }

      .category-icon {
        transform: scale(1.1) rotate(5deg);
        color: #fff;
        background: linear-gradient(135deg, var(--color-brand) 0%, #764ba2 100%);
      }
    }

    .category-icon {
      font-size: 2.5rem;
      width: 60px;
      height: 60px;
      display: flex;
      align-items: center;
      justify-content: center;
      background: rgba(102, 126, 234, 0.12);
      border-radius: 12px;
      color: var(--color-brand);
      flex-shrink: 0;
      transition: all 0.3s;
    }

    .category-info {
      flex: 1;

      h3 {
        font-size: 1.15rem;
        font-weight: 700;
        margin: 0 0 0.5rem 0;
        color: var(--color-text-primary);
        letter-spacing: 0.5px;
      }

      .post-count {
        font-size: 0.9rem;
        color: var(--color-text-secondary);
        font-weight: 500;
      }
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
    background: var(--color-surface);
    border-radius: 16px;
    overflow: hidden;
    cursor: pointer;
    transition: all 0.4s cubic-bezier(0.4, 0, 0.2, 1);
    border: 1px solid var(--color-border);
    box-shadow: 0 4px 12px rgba(0, 0, 0, 0.05);

    &:hover {
      transform: translateY(-12px);
      box-shadow: 0 20px 40px rgba(0, 0, 0, 0.15);
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
        background: linear-gradient(135deg, rgba(102, 126, 234, 0.2) 0%, rgba(118, 75, 162, 0.2) 100%);
        opacity: 0;
        transition: opacity 0.3s;
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
    background: var(--color-bg);
    border-radius: 12px;
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
        background: rgba(0, 0, 0, 0.1);
        opacity: 0;
        transition: opacity 0.3s;
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
    border: 1px solid var(--color-border);
    background: var(--color-surface);
    transition: all 0.4s cubic-bezier(0.4, 0, 0.2, 1);
    text-align: center;
    position: relative;
    overflow: hidden;

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
      box-shadow: 0 12px 32px rgba(0, 0, 0, 0.12);
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

