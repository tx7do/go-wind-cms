<script setup lang="ts">
import {definePage} from 'unplugin-vue-router/runtime'
import {useRouter} from 'vue-router'
import {ref, onMounted, computed, onUnmounted} from 'vue'

import {$t, i18n} from '@/locales'
import {XIcon} from "@/plugins/xicon";
import {useLanguageChangeEffect} from '@/hooks/useLanguageChangeEffect';

import FeaturesSection from '@/components/FeaturesSection.vue';
import CategoryListSection from '@/components/CategoryListSection.vue';
import PopularTagsSection from '@/components/PopularTagsSection.vue';


definePage({
  name: 'home',
  meta: {
    level: 1,
  },
})

const router = useRouter()

const loading = ref(false)

const featuredPostList = ref();
const latestPostList = ref();
const categoryListSection = ref();
const popularTags = ref();

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

// 用于Featured Posts的props，确保对象引用稳定，避免重复API调用
const featuredQueryParams = ref({
  status: 'POST_STATUS_PUBLISHED',
  isFeatured: true
})
const featuredOrderBy = ref(['-sortOrder'])
const featuredFieldMask = ref('id,status,sort_order,is_featured,visits,likes,comment_count,author_name,available_languages,created_at,translations.id,translations.post_id,translations.language_code,translations.title,translations.summary,translations.thumbnail')

onMounted(async () => {
  loading.value = true
  await Promise.all([
    categoryListSection.value?.reload(),
    popularTags.value?.reload(),
    featuredPostList.value?.reload(),
    latestPostList.value?.reload(),
  ])

  loading.value = false

  // 数据加载完成后初始化滚动动画
  setTimeout(() => {
    initScrollObserver()
  }, 100)
})

// 监听语言切换，自动重新加载数据
useLanguageChangeEffect(async () => {
  // 数据更新后，需要重新初始化滚动观察器
  // 先销毁旧的观察器
  destroyScrollObserver()
  // 等待 DOM 更新完成后重新初始化
  setTimeout(() => {
    initScrollObserver()
  }, 100)
}, {
  immediate: false,      // 是否立即执行一次
  autoCleanup: true,    // 是否自动清理
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
              class="code-function">GetPost</span>() {'{'}
            </div>
            <div class="code-line"><span class="code-keyword">return</span> post</div>
            <div class="code-line">{'}'}</div>
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
    <section class="featured-section scroll-reveal">
      <div class="section-header">
        <h2 class="section-title">
          <XIcon name="carbon:star-filled" :size="28"
                 style="color: #6366f1; margin-right: 8px;"/>
          {{ $t('page.home.featured_posts') }}
        </h2>
        <n-button text type="primary" @click="router.push('/post')">
          {{ $t('page.home.view_all') }} →
        </n-button>
      </div>
      <!-- 直接渲染 PostList，由 PostList 控制 skeleton -->
      <div class="featured-grid">
        <PostList
          ref="featuredPostList"
          :query-params="featuredQueryParams"
          :field-mask="featuredFieldMask"
          :order-by="featuredOrderBy"
          :page="1"
          :page-size="3"
          :columns="3"
          :show-skeleton="true"
          from="home"
        ></PostList>
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
    <section class="categories-section-wrapper scroll-reveal">
      <CategoryListSection
        :title="$t('page.home.categories')"
        :showViewAll="true"
        :showCarousel="true"
        :pageSize="8"
      >
      </CategoryListSection>
    </section>

    <!-- Popular Tags Section -->
    <section class="popular-section scroll-reveal">
      <div class="section-header">
        <h2 class="section-title">
          <XIcon name="carbon:fire" :size="28"
                 style="color: #6366f1; margin-right: 8px;"/>
          {{ $t('page.home.popular_tags') }}
        </h2>
        <n-button text type="primary" @click="router.push('/tag')">
          {{ $t('page.home.view_all') }} →
        </n-button>
      </div>
      <div class="popular-grid">
        <PopularTagsSection
          :title="''"
        />
      </div>
    </section>

    <!-- Latest Posts Section -->
    <section class="latest-section scroll-reveal">
      <div class="section-header">
        <h2 class="section-title">
          <XIcon name="carbon:document" :size="28"
                 style="color: #6366f1; margin-right: 8px;"/>
          {{ $t('page.home.latest_posts') }}
        </h2>
        <n-button text type="primary" @click="router.push('/post')">
          {{ $t('page.home.view_all') }} →
        </n-button>
      </div>
      <div class="latest-grid">
        <PostList
          ref="latestPostList"
          :query-params="{status: 'POST_STATUS_PUBLISHED'}"
          :field-mask="'id,status,sort_order,is_featured,visits,likes,comment_count,author_name,available_languages,created_at,translations.id,translations.post_id,translations.language_code,translations.title,translations.summary,translations.thumbnail'"
          :order-by="['-createdAt']"
          :page="1"
          :page-size="6"
          :columns="3"
          :show-skeleton="true"
          from="home"
        ></PostList>
      </div>
    </section>

    <!-- Features Section -->
    <FeaturesSection :title="$t('page.home.platform_features')" :features="features"/>
  </div>
</template>

<style scoped lang="less">
// 主体样式
.home-page {
  min-height: 100vh;
  width: 100%;
  background: var(--color-bg);
}

// 只保留一份核心样式定义，移除所有重复
.hero {
  padding: 3.5rem 2rem 2.5rem;
  text-align: center;
  background: linear-gradient(135deg, #6366f1 0%, #8b5cf6 50%, #a855f7 100%);
  color: white;
  position: relative;
  overflow: hidden;
  min-height: 280px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.hero-content {
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
  background: linear-gradient(135deg, #fff 0%, #f0f0ff 100%);
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
  line-height: 1.6;
}

.hero-actions {
  display: flex;
  gap: 1.5rem;
  justify-content: center;
  flex-wrap: wrap;
}

.section-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 2rem;
  max-width: 1200px;
  margin-left: auto;
  margin-right: auto;
  padding: 0 2rem;
}

.section-title {
  font-size: 2.25rem;
  font-weight: 800;
  letter-spacing: -1px;
  background: linear-gradient(135deg, #6366f1 0%, #8b5cf6 50%, #a855f7 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
  text-shadow: 0 2px 8px rgba(99, 102, 241, 0.08), 0 1px 2px rgba(0, 0, 0, 0.08);
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
  background: linear-gradient(135deg, #6366f1 0%, #8b5cf6 50%, #a855f7 100%);
  color: white;
  position: relative;
  overflow: hidden;
  min-height: 280px;
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
    text-shadow: 0 2px 10px rgba(255, 255, 255, 0.7);
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
    text-shadow: 0 1px 8px rgba(255, 255, 255, 0.6);
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

  .section-title {
    font-size: 2.25rem;
    font-weight: 800;
    letter-spacing: -1px;
    color: var(--color-text-primary);
    margin: 0;
    line-height: 1.1;
    background: linear-gradient(135deg, #6366f1 0%, #8b5cf6 50%, #a855f7 100%);
    -webkit-background-clip: text;
    -webkit-text-fill-color: transparent;
    background-clip: text;
    /* text-fill-color: transparent; */
    text-shadow: 0 2px 8px rgba(99, 102, 241, 0.08), 0 1px 2px rgba(0, 0, 0, 0.08);
  }

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
      // 移除原有 h2 的单独样式
    }
  }

  // Categories Section
  .categories-section {
    background: #f8fafc;
    padding: 3rem 0;

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
      background: transparent;
    }

    .category-card {
      display: flex;
      flex-direction: column;
      padding: 1.5rem;
      background: linear-gradient(135deg, #fff 0%, #f3f4f6 100%);
      border: 1.5px solid #e0e7ef;
      box-shadow: 0 2px 8px #e0e7ff;
      color: #3730a3;
      border-radius: 16px;
      cursor: pointer;
      position: relative;
      min-height: 200px;
      height: 100%;
      transition: box-shadow 0.3s, border-color 0.3s, background 0.3s;
      overflow: visible;

      .category-icon {
        font-size: 2.5rem;
        width: 70px;
        height: 70px;
        display: flex;
        align-items: center;
        justify-content: center;
        background: #e0e7ff;
        color: #6366f1;
        border-radius: 14px;
        margin-bottom: 1rem;
        box-shadow: 0 2px 8px #e0e7ff44;
        transition: background 0.3s, color 0.3s, transform 0.3s;
      }

      .category-info {
        flex: 1;
        width: 100%;

        h3 {
          font-size: 1.1rem;
          font-weight: 800;
          margin: 0 0 0.5rem 0;
          color: #3730a3;
          letter-spacing: 0.3px;
          word-break: break-word;
        }

        .post-count {
          font-size: 0.95rem;
          color: #6366f1;
          font-weight: 700;
          letter-spacing: 0.2px;
        }
      }

      .category-badge {
        margin-top: auto;
        padding: 0.5rem 0.875rem;
        background: #e0e7ff;
        color: #6366f1;
        border-radius: 8px;
        border-left: 4px solid #6366f1;
        font-size: 0.8rem;
        font-weight: 700;
        width: fit-content;
        white-space: nowrap;
        display: flex;
        align-items: center;
        gap: 0.5rem;
        box-shadow: 0 2px 8px #e0e7ff44;
      }

      &:hover {
        background: #f1f5f9;
        border-color: #6366f1;
        box-shadow: 0 8px 24px #a5b4fc44;

        .category-icon {
          background: #6366f1;
          color: #fff;
          transform: scale(1.1) rotate(8deg);
        }

        .category-info h3 {
          color: #4338ca;
        }

        .post-count {
          color: #4338ca;
        }

        .category-badge {
          background: #6366f1;
          color: #fff;
          border-left-color: #4338ca;
        }
      }
    }
  }
}

html.dark {
  .categories-section {
    background: #181a20;
    padding: 3rem 0;

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
      background: transparent;
    }

    .category-card {
      display: flex;
      flex-direction: column;
      padding: 1.5rem;
      background: linear-gradient(135deg, #23263a 0%, #181a20 100%);
      border: 1.5px solid #23263a;
      box-shadow: 0 2px 8px #23263a;
      color: #e0e7ef;
      border-radius: 16px;
      cursor: pointer;
      position: relative;
      min-height: 200px;
      height: 100%;
      transition: box-shadow 0.3s, border-color 0.3s, background 0.3s;
      overflow: visible;

      .category-icon {
        font-size: 2.5rem;
        width: 70px;
        height: 70px;
        display: flex;
        align-items: center;
        justify-content: center;
        background: #23263a;
        color: #a5b4fc;
        border-radius: 14px;
        margin-bottom: 1rem;
        box-shadow: 0 2px 8px #23263a44;
        transition: background 0.3s, color 0.3s, transform 0.3s;
      }

      .category-info {
        flex: 1;
        width: 100%;

        h3 {
          font-size: 1.1rem;
          font-weight: 800;
          margin: 0 0 0.5rem 0;
          color: #e0e7ef;
          letter-spacing: 0.3px;
          word-break: break-word;
        }

        .post-count {
          font-size: 0.95rem;
          color: #a5b4fc;
          font-weight: 700;
          letter-spacing: 0.2px;
        }
      }

      .category-badge {
        margin-top: auto;
        padding: 0.5rem 0.875rem;
        background: #23263a;
        color: #a5b4fc;
        border-radius: 8px;
        border-left: 4px solid #6366f1;
        font-size: 0.8rem;
        font-weight: 700;
        width: fit-content;
        white-space: nowrap;
        display: flex;
        align-items: center;
        gap: 0.5rem;
        box-shadow: 0 2px 8px #23263a44;
      }

      &:hover {
        background: #23263a;
        border-color: #6366f1;
        box-shadow: 0 8px 24px #6366f144;

        .category-icon {
          background: #6366f1;
          color: #fff;
          transform: scale(1.1) rotate(8deg);
        }

        .category-info h3 {
          color: #a5b4fc;
        }

        .post-count {
          color: #a5b4fc;
        }

        .category-badge {
          background: #6366f1;
          color: #fff;
          border-left-color: #a5b4fc;
        }
      }
    }
  }

  .popular-tags-container {
    max-width: 1200px;
    margin: 0 auto;
    padding: 0 2rem;

    .tags-title {
      font-size: 1.25rem;
      font-weight: 600;
      margin: 0 0 1.5rem 0;
      color: #e0e7ef;
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
      background: linear-gradient(135deg, #23263a 0%, #181a20 100%);
      border: 1.5px solid var(--tag-color, #a5b4fc);
      border-radius: 24px;
      cursor: pointer;
      transition: all 0.3s;
      color: var(--tag-color, #a5b4fc);
      font-weight: 500;
      box-shadow: 0 2px 8px #23263a44;

      &:hover {
        background: var(--tag-color, #6366f1);
        color: #fff;
        box-shadow: 0 4px 12px #6366f144;

        .tag-label {
          color: #fff !important;
        }
      }
    }

    .tag-label {
      font-size: 0.95rem;
      color: var(--tag-color, #a5b4fc);
      transition: color 0.3s;
      font-weight: 600;
    }
  }
}

@media (prefers-color-scheme: light) {
  .popular-tags-container {
    max-width: 1200px;
    margin: 0 auto;
    padding: 0 2rem;

    .tags-title {
      font-size: 1.25rem;
      font-weight: 600;
      margin: 0 0 1.5rem 0;
      color: #1f2937;
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
      background: linear-gradient(135deg, #fff 0%, #f3f4f6 100%);
      border: 1.5px solid var(--tag-color, #6366f1);
      border-radius: 24px;
      cursor: pointer;
      transition: all 0.3s;
      color: var(--tag-color, #6366f1);
      font-weight: 500;
      box-shadow: 0 2px 8px #e0e7ff44;

      &:hover {
        background: var(--tag-color, #6366f1);
        color: #fff;
        box-shadow: 0 4px 12px #a5b4fc44;

        .tag-label {
          color: #fff !important;
        }
      }
    }

    .tag-label {
      font-size: 0.95rem;
      color: var(--tag-color, #6366f1);
      transition: color 0.3s;
      font-weight: 600;
    }
  }
}

// Scroll Indicator Styles
.scroll-indicator {
  margin-top: 2rem;
  display: flex;
  justify-content: center;
  align-items: center;
  width: 100%;

  .scroll-indicator-content {
    display: flex;
    align-items: center;
    gap: 1rem;
    background: none;
    padding: 0;

    .scroll-indicator-text {
      font-size: 2.25rem;
      font-weight: 800;
      letter-spacing: -1px;
      background: linear-gradient(135deg, #6366f1 0%, #8b5cf6 50%, #a855f7 100%);
      -webkit-background-clip: text;
      -webkit-text-fill-color: transparent;
      background-clip: text;
      text-shadow: 0 2px 8px rgba(99, 102, 241, 0.08), 0 1px 2px rgba(0, 0, 0, 0.08);
      margin: 0;
      line-height: 1.1;
    }

    .scroll-arrow {
      display: flex;
      align-items: center;
      justify-content: center;
      width: 36px;
      height: 36px;
      border-radius: 50%;
      background: #e0e7ff;
      box-shadow: 0 2px 8px #e0e7ff44;
      cursor: pointer;
      transition: background 0.3s;

      .arrow-icon {
        color: #6366f1;
        transition: color 0.3s;
      }

      &:hover {
        background: #6366f1;

        .arrow-icon {
          color: #fff;
        }
      }
    }
  }
}

// 移动端隐藏滚动指示器
@media (max-width: 900px) {
  .scroll-indicator {
    display: none;
  }
}

html.dark .scroll-indicator-content .scroll-indicator-text {
  background: linear-gradient(135deg, #a5b4fc 0%, #6366f1 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
  text-shadow: 0 2px 8px #23263a44;
  color: #e0e7ef;
}

.featured-section,
.latest-section,
.popular-section {
  margin-top: 3rem;
  margin-bottom: 3rem;

  .section-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 2rem;
    max-width: 1200px;
    margin-left: auto;
    margin-right: auto;
    padding: 0 2rem;

    .section-title {
      font-size: 2.25rem;
      font-weight: 800;
      letter-spacing: -1px;
      background: linear-gradient(135deg, #6366f1 0%, #8b5cf6 50%, #a855f7 100%);
      -webkit-background-clip: text;
      -webkit-text-fill-color: transparent;
      background-clip: text;
      text-shadow: 0 2px 8px rgba(99, 102, 241, 0.08), 0 1px 2px rgba(0, 0, 0, 0.08);
      margin: 0;
      line-height: 1.1;
      display: flex;
      align-items: center;
      gap: 0.5rem;
    }

    .n-button {
      font-size: 1rem;
      font-weight: 700;
      color: #6366f1;
      background: none;
      border: none;
      padding: 0.5rem 1rem;
      transition: color 0.3s;

      &:hover {
        color: #a855f7;
      }
    }
  }

  .featured-grid,
  .latest-grid,
  .popular-grid {
    max-width: 1200px;
    margin: 0 auto;
    padding: 0 2rem;
  }
}

html.dark .featured-section .section-header .section-title,
html.dark .latest-section .section-header .section-title,
html.dark .popular-section .section-header .section-title {
  background: linear-gradient(135deg, #a5b4fc 0%, #6366f1 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
  text-shadow: 0 2px 8px #23263a44;
  color: #e0e7ef;
}
</style>
