<script setup lang="ts">
import { definePage } from 'unplugin-vue-router/runtime'
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useCategoryStore } from '@/stores/modules/app'
import { $t } from '@/locales'

definePage({
  name: 'category-list',
  meta: {
    title: 'Categories',
    level: 2,
  },
})

const router = useRouter()
const categoryStore = useCategoryStore()

const loading = ref(false)
const categories = ref<any[]>([])

async function loadCategories() {
  loading.value = true
  try {
    const res = await categoryStore.listCategory(
      undefined,
      { status: 'CATEGORY_STATUS_ACTIVE' }
    )
    categories.value = res.items || []
  } catch (error) {
    console.error('Load categories failed:', error)
  } finally {
    loading.value = false
  }
}

function getCategoryName(category: any) {
  return category.translations?.[0]?.name || $t('page.home.category_default')
}

function getCategoryDescription(category: any) {
  return category.translations?.[0]?.description || ''
}

function getCategoryThumbnail(category: any) {
  return category.translations?.[0]?.thumbnail || '/placeholder.jpg'
}

function handleViewCategory(id: number) {
  router.push({
    path: '/post',
    query: { category: id }
  })
}

onMounted(() => {
  loadCategories()
})
</script>

<template>
  <div class="category-page">
    <!-- Hero Section -->
    <div class="hero-section">
      <div class="hero-content">
        <h1>{{ $t('page.categories.categories') }}</h1>
        <p>{{ $t('page.categories.browse_all') }}</p>
      </div>
    </div>

    <!-- Content Section -->
    <div class="page-container">
      <n-spin :show="loading">
        <div v-if="categories.length > 0" class="categories-grid">
          <div
            v-for="category in categories"
            :key="category.id"
            class="category-card"
            @click="handleViewCategory(category.id)"
          >
            <div class="category-image">
              <img :src="getCategoryThumbnail(category)" :alt="getCategoryName(category)" />
              <div class="image-overlay" />
              <div class="view-button">
                <span>{{ $t('page.home.browse_posts') }}</span>
                <span class="i-carbon:arrow-right" />
              </div>
            </div>
            <div class="category-content">
              <h3>{{ getCategoryName(category) }}</h3>
              <p>{{ getCategoryDescription(category) }}</p>
              <div class="category-meta">
                <span class="meta-icon">
                  <span class="i-carbon:document" />
                </span>
                <span class="meta-text">{{ category.postCount || 0 }} {{ $t('page.posts.articles') }}</span>
              </div>
            </div>
          </div>
        </div>

        <n-empty v-else :description="$t('page.categories.no_categories')" style="margin: 80px 0;">
          <template #icon>
            <span class="i-carbon:folder-blank" style="font-size: 64px;" />
          </template>
        </n-empty>
      </n-spin>
    </div>
  </div>
</template>

<style scoped lang="less">
.category-page {
  min-height: 100vh;
  background: var(--color-bg);
}

// Hero Section
.hero-section {
  background: linear-gradient(135deg, var(--color-brand) 0%, #764ba2 100%);
  padding: 60px 32px;
  margin-bottom: 40px;
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
    max-width: 1200px;
    margin: 0 auto;
    text-align: center;
    position: relative;
    z-index: 1;

    h1 {
      font-size: 48px;
      font-weight: 800;
      margin: 0 0 16px 0;
      color: #ffffff;
      text-shadow: 0 2px 8px rgba(0, 0, 0, 0.2);
      letter-spacing: -1px;
    }

    p {
      font-size: 20px;
      color: rgba(255, 255, 255, 0.95);
      margin: 0;
      font-weight: 500;
    }
  }
}

// Page Container
.page-container {
  max-width: 1200px;
  margin: 0 auto;
  padding: 0 32px 80px;
}

// Categories Grid
.categories-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
  gap: 28px;
}

.category-card {
  background: var(--color-surface);
  border-radius: 16px;
  overflow: hidden;
  cursor: pointer;
  transition: all 0.4s cubic-bezier(0.4, 0, 0.2, 1);
  border: 1px solid var(--color-border);
  height: 100%;
  display: flex;
  flex-direction: column;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.05);

  &:hover {
    transform: translateY(-12px);
    box-shadow: 0 20px 40px rgba(0, 0, 0, 0.15);
    border-color: var(--color-brand);

    .category-image {
      .image-overlay {
        opacity: 0.6;
      }

      .view-button {
        opacity: 1;
        transform: translateY(0);
      }

      img {
        transform: scale(1.1);
      }
    }

    h3 {
      color: var(--color-brand);
    }
  }

  .category-image {
    position: relative;
    width: 100%;
    height: 240px;
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
      background: rgba(0, 0, 0, 0.2);
      transition: opacity 0.3s;
      opacity: 0;
    }

    .view-button {
      position: absolute;
      bottom: 0;
      left: 0;
      right: 0;
      background: linear-gradient(to top, rgba(0, 0, 0, 0.8) 0%, transparent 100%);
      color: white;
      padding: 24px 20px;
      display: flex;
      align-items: center;
      gap: 8px;
      font-weight: 600;
      font-size: 15px;
      opacity: 0;
      transform: translateY(10px);
      transition: all 0.4s cubic-bezier(0.4, 0, 0.2, 1);

      span[class^="i-"] {
        font-size: 18px;
      }
    }
  }

  .category-content {
    padding: 24px;
    flex: 1;
    display: flex;
    flex-direction: column;

    h3 {
      font-size: 20px;
      font-weight: 700;
      margin: 0 0 12px 0;
      color: var(--color-text-primary);
      line-height: 1.4;
      transition: color 0.3s;
      letter-spacing: 0.5px;
    }

    p {
      color: var(--color-text-secondary);
      font-size: 15px;
      line-height: 1.7;
      margin: 0 0 16px 0;
      display: -webkit-box;
      -webkit-line-clamp: 2;
      -webkit-box-orient: vertical;
      overflow: hidden;
      flex: 1;
    }

    .category-meta {
      display: flex;
      align-items: center;
      gap: 8px;
      padding-top: 16px;
      border-top: 1px solid var(--color-border);
      font-size: 14px;
      color: var(--color-text-secondary);
      font-weight: 500;

      .meta-icon {
        display: flex;
        align-items: center;
        font-size: 18px;
        opacity: 0.8;
      }

      .meta-text {
        display: flex;
        align-items: center;
      }
    }
  }
}

// Responsive Design
@media (max-width: 1024px) {
  .hero-section {
    padding: 48px 24px;
  }

  .page-container {
    padding: 0 24px 60px;
  }

  .categories-grid {
    gap: 24px;
  }

  .hero-section .hero-content {
    h1 {
      font-size: 36px;
    }

    p {
      font-size: 18px;
    }
  }
}

@media (max-width: 768px) {
  .hero-section {
    padding: 40px 20px;

    .hero-content {
      h1 {
        font-size: 32px;
      }

      p {
        font-size: 16px;
      }
    }
  }

  .page-container {
    padding: 0 20px 40px;
  }

  .categories-grid {
    grid-template-columns: repeat(auto-fill, minmax(280px, 1fr));
    gap: 20px;
  }

  .category-card {
    .category-image {
      height: 200px;

      .view-button {
        padding: 20px 16px;
        font-size: 14px;
      }
    }

    .category-content {
      padding: 20px;

      h3 {
        font-size: 18px;
      }

      p {
        font-size: 14px;
      }
    }
  }
}

@media (max-width: 480px) {
  .hero-section {
    padding: 32px 16px;

    .hero-content {
      h1 {
        font-size: 28px;
      }

      p {
        font-size: 15px;
      }
    }
  }

  .page-container {
    padding: 0 16px 40px;
  }

  .categories-grid {
    grid-template-columns: 1fr;
    gap: 16px;
  }

  .category-card {
    .category-image {
      height: 180px;

      .view-button {
        padding: 16px 12px;
        font-size: 13px;
      }
    }

    .category-content {
      padding: 16px;

      h3 {
        font-size: 16px;
      }

      p {
        font-size: 13px;
      }

      .category-meta {
        font-size: 12px;
        gap: 6px;

        .meta-icon {
          font-size: 16px;
        }
      }
    }
  }
}

// Dark Mode Enhancements
html.dark {
  .category-card {
    &:hover {
      box-shadow: 0 20px 40px rgba(0, 0, 0, 0.4);
    }
  }
}
</style>

