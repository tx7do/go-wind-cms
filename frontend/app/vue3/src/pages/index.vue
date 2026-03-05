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
  padding: 4rem 2rem;
  text-align: center;
  max-width: 1200px;
  background: linear-gradient(135deg, var(--color-brand) 0%, var(--color-brand-accent) 100%);
  color: white;
  border-radius: var(--radius-lg);
  margin: 2rem auto 0;

  .hero-content {
    animation: fadeInUp 0.8s ease-out;
  }

  .hero-title {
    font-size: 3rem;
    font-weight: 700;
    margin-bottom: 1rem;
    line-height: 1.2;
  }

  .hero-subtitle {
    font-size: 1.25rem;
    margin-bottom: 2rem;
    max-width: 600px;
    margin-left: auto;
    margin-right: auto;
    opacity: 0.95;
  }

  .hero-actions {
    display: flex;
    gap: 1rem;
    justify-content: center;

    :deep(.n-button) {
      font-size: 1rem;
      font-weight: 500;
      transition: all 0.3s;
    }

    :deep(.n-button--primary) {
      background: white !important;
      color: var(--color-brand) !important;

      &:hover {
        opacity: 0.9;
      }
    }

    :deep(.n-button--default) {
      background: transparent !important;
      color: white !important;
      border: 2px solid white !important;

      &:hover {
        background: rgba(255, 255, 255, 0.1) !important;
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
  padding: 3rem 0;
  background: var(--color-surface);

  .categories-grid {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(250px, 1fr));
    gap: 1.5rem;
    max-width: 1200px;
    margin: 0 auto;
    padding: 0 2rem;
  }

  .category-card {
    display: flex;
    align-items: center;
    gap: 1rem;
    padding: 1.5rem;
    background: var(--color-bg);
    border: 1px solid var(--color-border);
    border-radius: var(--radius-md);
    cursor: pointer;
    transition: all 0.3s;

    &:hover {
      transform: translateY(-4px);
      box-shadow: 0 8px 16px rgba(0, 0, 0, 0.1);
      border-color: var(--color-brand);
    }

    .category-icon {
      font-size: 2rem;
      color: var(--color-brand);
      width: 48px;
      height: 48px;
      display: flex;
      align-items: center;
      justify-content: center;
      background: rgba(102, 126, 234, 0.1);
      border-radius: var(--radius-md);
    }

    .category-info {
      flex: 1;

      h3 {
        font-size: 1.1rem;
        font-weight: 600;
        margin: 0 0 0.25rem 0;
        color: var(--color-text-primary);
      }

      .post-count {
        font-size: 0.85rem;
        color: var(--color-text-secondary);
      }
    }
  }
}

// Featured Posts Section
.featured-section {
  padding: 3rem 0;
  background: var(--color-bg);

  .featured-grid {
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(320px, 1fr));
    gap: 2rem;
    max-width: 1200px;
    margin: 0 auto;
    padding: 0 2rem;
  }

  .featured-card {
    background: var(--color-surface);
    border-radius: var(--radius-md);
    overflow: hidden;
    cursor: pointer;
    transition: all 0.3s;
    border: 1px solid var(--color-border);

    &:hover {
      transform: translateY(-8px);
      box-shadow: 0 12px 24px rgba(0, 0, 0, 0.15);
    }

    .featured-image {
      width: 100%;
      height: 200px;
      overflow: hidden;

      img {
        width: 100%;
        height: 100%;
        object-fit: cover;
        transition: transform 0.3s;
      }

      &:hover img {
        transform: scale(1.05);
      }
    }

    .featured-content {
      padding: 1.5rem;

      h3 {
        font-size: 1.25rem;
        font-weight: 600;
        margin: 0 0 0.5rem 0;
        color: var(--color-text-primary);
      }

      p {
        color: var(--color-text-secondary);
        line-height: 1.6;
        margin: 0 0 1rem 0;
        display: -webkit-box;
        -webkit-line-clamp: 2;
        -webkit-box-orient: vertical;
        overflow: hidden;
      }

      .featured-meta {
        display: flex;
        gap: 1rem;
        font-size: 0.85rem;
        color: var(--color-text-secondary);
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
    grid-template-columns: repeat(auto-fill, minmax(350px, 1fr));
    gap: 2rem;
    max-width: 1200px;
    margin: 0 auto;
    padding: 0 2rem;
  }

  .post-card {
    background: var(--color-bg);
    border-radius: var(--radius-md);
    overflow: hidden;
    cursor: pointer;
    transition: all 0.3s;
    border: 1px solid var(--color-border);

    &:hover {
      transform: translateY(-4px);
      box-shadow: 0 8px 16px rgba(0, 0, 0, 0.1);
    }

    .post-image {
      width: 100%;
      height: 180px;
      overflow: hidden;

      img {
        width: 100%;
        height: 100%;
        object-fit: cover;
        transition: transform 0.3s;
      }

      &:hover img {
        transform: scale(1.05);
      }
    }

    .post-content {
      padding: 1.5rem;

      h3 {
        font-size: 1.15rem;
        font-weight: 600;
        margin: 0 0 0.5rem 0;
        color: var(--color-text-primary);
        display: -webkit-box;
        -webkit-line-clamp: 2;
        -webkit-box-orient: vertical;
        overflow: hidden;
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
      }

      .post-meta {
        display: flex;
        gap: 1rem;
        font-size: 0.85rem;
        color: var(--color-text-secondary);

        span {
          display: flex;
          align-items: center;
          gap: 0.25rem;
        }
      }
    }
  }
}

// Features Section
.features {
  padding: 3rem 2rem;
  background: var(--color-bg);

  .section-header {
    text-align: center;
    display: block;
    margin-bottom: 3rem;
  }

  .features-grid {
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(280px, 1fr));
    gap: 2rem;
    max-width: 1200px;
    margin: 0 auto;
  }

  .feature-card {
    padding: 2rem;
    border-radius: var(--radius-md);
    border: 1px solid var(--color-border);
    background: var(--color-surface);
    transition: all 0.3s;
    text-align: center;

    &:hover {
      transform: translateY(-4px);
      box-shadow: 0 12px 24px rgba(0, 0, 0, 0.1);
    }

    .feature-icon {
      font-size: 3rem;
      color: var(--color-brand);
      margin-bottom: 1rem;
    }

    h3 {
      font-size: 1.25rem;
      margin-bottom: 0.5rem;
      color: var(--color-text-primary);
    }

    p {
      color: var(--color-text-secondary);
      font-size: 0.95rem;
      line-height: 1.6;
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

// Responsive
@media (max-width: 768px) {
  .hero {
    padding: 2rem 1rem;
    margin-top: 1rem;

    .hero-title {
      font-size: 2rem;
    }

    .hero-subtitle {
      font-size: 1rem;
    }

    .hero-actions {
      flex-direction: column;
    }
  }

  .section-header {
    padding: 0 1rem;
  }

  .categories-grid,
  .posts-grid,
  .featured-grid {
    padding: 0 1rem;
    grid-template-columns: 1fr;
  }
}
</style>

