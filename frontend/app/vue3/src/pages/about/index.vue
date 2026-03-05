<script setup lang='ts'>
import { definePage } from "unplugin-vue-router/runtime"
import { computed } from 'vue'
import { $t, i18n } from '@/locales'
import XIcon from '@/plugins/xicon/XIcon.vue'

definePage({
  name: 'about',
  meta: {
    level: 2,
  },
})

const features = computed(() => {
  void i18n.global.locale.value

  return [
    {
      icon: 'carbon:document-add',
      title: $t('page.about.feature_content'),
      description: $t('page.about.feature_content_desc'),
    },
    {
      icon: 'carbon:cloud-upload',
      title: $t('page.about.feature_multi_tenant'),
      description: $t('page.about.feature_multi_tenant_desc'),
    },
    {
      icon: 'carbon:security',
      title: $t('page.about.feature_security'),
      description: $t('page.about.feature_security_desc'),
    },
    {
      icon: 'carbon:api',
      title: $t('page.about.feature_api'),
      description: $t('page.about.feature_api_desc'),
    },
    {
      icon: 'carbon:collaborate',
      title: $t('page.about.feature_collaboration'),
      description: $t('page.about.feature_collaboration_desc'),
    },
    {
      icon: 'carbon:analytics',
      title: $t('page.about.feature_analytics'),
      description: $t('page.about.feature_analytics_desc'),
    },
  ]
})

const teamMembers = computed(() => {
  return [
    {
      name: $t('page.about.team_member_1'),
      role: $t('page.about.team_role_1'),
      avatar: '/logo.png',
    },
    {
      name: $t('page.about.team_member_2'),
      role: $t('page.about.team_role_2'),
      avatar: '/logo.png',
    },
    {
      name: $t('page.about.team_member_3'),
      role: $t('page.about.team_role_3'),
      avatar: '/logo.png',
    },
  ]
})
</script>

<template>
  <div class="about-page">
    <!-- Hero Section -->
    <section class="hero">
      <div class="hero-content">
        <h1 class="hero-title">{{ $t('page.about.title') }}</h1>
        <p class="hero-subtitle">{{ $t('page.about.subtitle') }}</p>
        <p class="hero-description">{{ $t('page.about.description') }}</p>
      </div>
    </section>

    <!-- About Section -->
    <section class="about-section">
      <div class="section-container">
        <div class="about-content">
          <h2>{{ $t('page.about.about_us') }}</h2>
          <p>{{ $t('page.about.about_us_desc_1') }}</p>
          <p>{{ $t('page.about.about_us_desc_2') }}</p>
          <p>{{ $t('page.about.about_us_desc_3') }}</p>
        </div>
        <div class="about-stats">
          <div class="stat-card">
            <div class="stat-number">10K+</div>
            <div class="stat-label">{{ $t('page.about.stat_users') }}</div>
          </div>
          <div class="stat-card">
            <div class="stat-number">500+</div>
            <div class="stat-label">{{ $t('page.about.stat_projects') }}</div>
          </div>
          <div class="stat-card">
            <div class="stat-number">99.9%</div>
            <div class="stat-label">{{ $t('page.about.stat_uptime') }}</div>
          </div>
          <div class="stat-card">
            <div class="stat-number">24/7</div>
            <div class="stat-label">{{ $t('page.about.stat_support') }}</div>
          </div>
        </div>
      </div>
    </section>

    <!-- Features Section -->
    <section class="features-section">
      <div class="section-container">
        <div class="section-header">
          <h2>{{ $t('page.about.features') }}</h2>
          <p>{{ $t('page.about.features_desc') }}</p>
        </div>
        <div class="features-grid">
          <div v-for="feature in features" :key="feature.title" class="feature-card">
            <div class="feature-icon">
              <XIcon :name="feature.icon" :size="32" />
            </div>
            <h3>{{ feature.title }}</h3>
            <p>{{ feature.description }}</p>
          </div>
        </div>
      </div>
    </section>

    <!-- Team Section -->
    <section class="team-section">
      <div class="section-container">
        <div class="section-header">
          <h2>{{ $t('page.about.team') }}</h2>
          <p>{{ $t('page.about.team_desc') }}</p>
        </div>
        <div class="team-grid">
          <div v-for="member in teamMembers" :key="member.name" class="team-card">
            <div class="team-avatar">
              <img :src="member.avatar" :alt="member.name" />
            </div>
            <h3>{{ member.name }}</h3>
            <p>{{ member.role }}</p>
          </div>
        </div>
      </div>
    </section>

    <!-- Values Section -->
    <section class="values-section">
      <div class="section-container">
        <div class="section-header">
          <h2>{{ $t('page.about.values') }}</h2>
          <p>{{ $t('page.about.values_desc') }}</p>
        </div>
        <div class="values-grid">
          <div class="value-card">
            <h3>{{ $t('page.about.value_innovation') }}</h3>
            <p>{{ $t('page.about.value_innovation_desc') }}</p>
          </div>
          <div class="value-card">
            <h3>{{ $t('page.about.value_reliability') }}</h3>
            <p>{{ $t('page.about.value_reliability_desc') }}</p>
          </div>
          <div class="value-card">
            <h3>{{ $t('page.about.value_customer') }}</h3>
            <p>{{ $t('page.about.value_customer_desc') }}</p>
          </div>
        </div>
      </div>
    </section>

    <!-- CTA Section -->
    <section class="cta-section">
      <div class="section-container">
        <h2>{{ $t('page.about.cta_title') }}</h2>
        <p>{{ $t('page.about.cta_desc') }}</p>
        <div class="cta-buttons">
          <n-button type="primary" size="large" @click="$router.push('/post')">
            {{ $t('page.about.cta_explore') }}
          </n-button>
          <n-button type="default" size="large">
            {{ $t('page.about.cta_contact') }}
          </n-button>
        </div>
      </div>
    </section>
  </div>
</template>

<style scoped lang='less'>
.about-page {
  width: 100%;
  background: var(--color-bg);
  color: var(--color-text-primary);
}

.section-container {
  max-width: 1200px;
  margin: 0 auto;
  padding: 0 var(--space-4);
}

// Hero Section
.hero {
  background: linear-gradient(135deg, var(--color-brand) 0%, #764ba2 100%);
  color: white;
  padding: 80px 32px;
  text-align: center;
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
    max-width: 900px;
    margin: 0 auto;
    position: relative;
    z-index: 1;
    animation: fadeInUp 0.8s ease-out;
  }

  .hero-title {
    font-size: 56px;
    font-weight: 800;
    margin: 0 0 20px 0;
    line-height: 1.2;
    letter-spacing: -1px;
    text-shadow: 0 2px 8px rgba(0, 0, 0, 0.2);
  }

  .hero-subtitle {
    font-size: 24px;
    margin: 0 0 16px 0;
    opacity: 0.95;
    font-weight: 600;
  }

  .hero-description {
    font-size: 18px;
    opacity: 0.9;
    margin: 0;
    line-height: 1.7;
    max-width: 700px;
    margin-left: auto;
    margin-right: auto;
    font-weight: 500;
  }
}

// About Section
.about-section {
  padding: 80px 32px;
  background: var(--color-surface);

  .section-container {
    display: grid;
    grid-template-columns: 1fr 1fr;
    gap: 60px;
    align-items: center;
    max-width: 1200px;
    margin: 0 auto;
  }

  .about-content {
    h2 {
      font-size: 44px;
      margin: 0 0 32px 0;
      color: var(--color-text-primary);
      font-weight: 800;
      letter-spacing: -0.5px;
    }

    p {
      font-size: 16px;
      line-height: 1.9;
      margin: 0 0 24px 0;
      color: var(--color-text-secondary);
      font-weight: 500;

      &:last-child {
        margin-bottom: 0;
      }
    }
  }

  .about-stats {
    display: grid;
    grid-template-columns: 1fr 1fr;
    gap: 20px;

    .stat-card {
      background: var(--color-bg);
      padding: 32px 24px;
      border-radius: 12px;
      border: 1px solid var(--color-border);
      text-align: center;
      transition: all 0.4s cubic-bezier(0.4, 0, 0.2, 1);
      position: relative;
      overflow: hidden;

      &::before {
        content: '';
        position: absolute;
        top: 0;
        left: 0;
        right: 0;
        bottom: 0;
        background: linear-gradient(135deg, rgba(102, 126, 234, 0.1) 0%, transparent 100%);
        opacity: 0;
        transition: opacity 0.3s;
      }

      &:hover {
        transform: translateY(-6px);
        box-shadow: 0 12px 32px rgba(0, 0, 0, 0.12);
        border-color: var(--color-brand);

        &::before {
          opacity: 1;
        }

        .stat-number {
          color: var(--color-brand);
        }
      }

      .stat-number {
        font-size: 40px;
        font-weight: 800;
        color: var(--color-text-primary);
        margin: 0 0 12px 0;
        transition: color 0.3s;
        letter-spacing: -1px;
      }

      .stat-label {
        font-size: 14px;
        color: var(--color-text-secondary);
        margin: 0;
        font-weight: 600;
      }
    }
  }
}

// Features Section
.features-section {
  padding: 80px 32px;
  background: var(--color-bg);

  .section-header {
    text-align: center;
    margin-bottom: 60px;
    max-width: 700px;
    margin-left: auto;
    margin-right: auto;

    h2 {
      font-size: 44px;
      margin: 0 0 20px 0;
      color: var(--color-text-primary);
      font-weight: 800;
      letter-spacing: -0.5px;
    }

    p {
      font-size: 18px;
      color: var(--color-text-secondary);
      margin: 0;
      line-height: 1.7;
      font-weight: 500;
    }
  }

  .features-grid {
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(320px, 1fr));
    gap: 28px;
    max-width: 1200px;
    margin: 0 auto;

    .feature-card {
      background: var(--color-surface);
      padding: 40px 28px;
      border-radius: 16px;
      border: 1px solid var(--color-border);
      text-align: center;
      transition: all 0.4s cubic-bezier(0.4, 0, 0.2, 1);
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
        transform: translateY(-12px);
        box-shadow: 0 20px 40px rgba(0, 0, 0, 0.15);
        border-color: var(--color-brand);

        &::before {
          top: 0;
        }

        .feature-icon {
          transform: scale(1.2) rotate(-5deg);
          color: #fff;
          background: linear-gradient(135deg, var(--color-brand) 0%, #764ba2 100%);
        }

        h3 {
          color: var(--color-brand);
        }
      }

      .feature-icon {
        display: flex;
        justify-content: center;
        align-items: center;
        width: 80px;
        height: 80px;
        background: rgba(102, 126, 234, 0.12);
        border-radius: 12px;
        margin: 0 auto 24px;
        color: var(--color-brand);
        transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
        font-size: 32px;
      }

      h3 {
        font-size: 20px;
        margin: 0 0 16px 0;
        color: var(--color-text-primary);
        font-weight: 700;
        transition: color 0.3s;
        letter-spacing: 0.5px;
      }

      p {
        font-size: 15px;
        color: var(--color-text-secondary);
        line-height: 1.8;
        margin: 0;
      }
    }
  }
}

// Team Section
.team-section {
  padding: 80px 32px;
  background: var(--color-surface);

  .section-header {
    text-align: center;
    margin-bottom: 60px;
    max-width: 700px;
    margin-left: auto;
    margin-right: auto;

    h2 {
      font-size: 44px;
      margin: 0 0 20px 0;
      color: var(--color-text-primary);
      font-weight: 800;
      letter-spacing: -0.5px;
    }

    p {
      font-size: 18px;
      color: var(--color-text-secondary);
      margin: 0;
      line-height: 1.7;
      font-weight: 500;
    }
  }

  .team-grid {
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(280px, 1fr));
    gap: 32px;
    max-width: 1200px;
    margin: 0 auto;

    .team-card {
      text-align: center;
      padding: 32px 24px;
      border-radius: 16px;
      background: var(--color-bg);
      border: 1px solid var(--color-border);
      transition: all 0.4s cubic-bezier(0.4, 0, 0.2, 1);

      &:hover {
        transform: translateY(-12px);
        box-shadow: 0 20px 40px rgba(0, 0, 0, 0.15);
        border-color: var(--color-brand);

        .team-avatar {
          transform: scale(1.1);
          box-shadow: 0 12px 28px rgba(102, 126, 234, 0.3);
          border-color: var(--color-brand);
        }

        h3 {
          color: var(--color-brand);
        }
      }

      .team-avatar {
        width: 120px;
        height: 120px;
        margin: 0 auto 20px;
        border-radius: 50%;
        overflow: hidden;
        border: 4px solid var(--color-border);
        transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
        background: linear-gradient(135deg, var(--color-brand) 0%, #764ba2 100%);

        img {
          width: 100%;
          height: 100%;
          object-fit: cover;
        }
      }

      h3 {
        font-size: 20px;
        margin: 0 0 8px 0;
        color: var(--color-text-primary);
        font-weight: 700;
        transition: color 0.3s;
        letter-spacing: 0.5px;
      }

      p {
        font-size: 14px;
        color: var(--color-text-secondary);
        margin: 0;
        font-weight: 600;
      }
    }
  }
}

// Values Section
.values-section {
  padding: 80px 32px;
  background: var(--color-bg);

  .section-header {
    text-align: center;
    margin-bottom: 60px;
    max-width: 700px;
    margin-left: auto;
    margin-right: auto;

    h2 {
      font-size: 44px;
      margin: 0 0 20px 0;
      color: var(--color-text-primary);
      font-weight: 800;
      letter-spacing: -0.5px;
    }

    p {
      font-size: 18px;
      color: var(--color-text-secondary);
      margin: 0;
      line-height: 1.7;
      font-weight: 500;
    }
  }

  .values-grid {
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(320px, 1fr));
    gap: 28px;
    max-width: 1200px;
    margin: 0 auto;

    .value-card {
      background: var(--color-surface);
      padding: 40px 28px;
      border-radius: 16px;
      border: 1px solid var(--color-border);
      transition: all 0.4s cubic-bezier(0.4, 0, 0.2, 1);
      position: relative;
      overflow: hidden;

      &::before {
        content: '';
        position: absolute;
        top: 0;
        left: 0;
        right: 0;
        bottom: 0;
        background: linear-gradient(135deg, rgba(102, 126, 234, 0.05) 0%, transparent 100%);
        opacity: 0;
        transition: opacity 0.3s;
      }

      &:hover {
        transform: translateY(-8px);
        box-shadow: 0 12px 32px rgba(0, 0, 0, 0.12);
        border-color: var(--color-brand);

        &::before {
          opacity: 1;
        }

        h3 {
          color: var(--color-brand);
        }
      }

      h3 {
        font-size: 20px;
        margin: 0 0 16px 0;
        color: var(--color-text-primary);
        font-weight: 700;
        transition: color 0.3s;
        letter-spacing: 0.5px;
      }

      p {
        font-size: 15px;
        color: var(--color-text-secondary);
        line-height: 1.8;
        margin: 0;
      }
    }
  }
}

// CTA Section
.cta-section {
  background: linear-gradient(135deg, var(--color-brand) 0%, #764ba2 100%);
  color: white;
  padding: 80px 32px;
  text-align: center;
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

  .section-container {
    max-width: 800px;
    margin: 0 auto;
    position: relative;
    z-index: 1;
  }

  h2 {
    font-size: 44px;
    margin: 0 0 20px 0;
    font-weight: 800;
    letter-spacing: -0.5px;
    text-shadow: 0 2px 8px rgba(0, 0, 0, 0.2);
  }

  p {
    font-size: 18px;
    opacity: 0.95;
    margin: 0 0 40px 0;
    line-height: 1.7;
    font-weight: 500;
  }

  .cta-buttons {
    display: flex;
    gap: 16px;
    justify-content: center;
    flex-wrap: wrap;

    :deep(.n-button) {
      font-size: 16px;
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
  }
}

// Animations
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
  .hero {
    padding: 60px 24px;

    .hero-title {
      font-size: 44px;
    }

    .hero-subtitle {
      font-size: 20px;
    }

    .hero-description {
      font-size: 16px;
    }
  }

  .about-section .section-container {
    gap: 40px;
  }

  .section-header h2 {
    font-size: 36px;
  }
}

@media (max-width: 768px) {
  .hero {
    padding: 48px 20px;

    .hero-title {
      font-size: 32px;
    }

    .hero-subtitle {
      font-size: 18px;
    }

    .hero-description {
      font-size: 15px;
    }
  }

  .about-section {
    padding: 60px 20px;

    .section-container {
      grid-template-columns: 1fr;
      gap: 40px;
    }

    .about-content h2 {
      font-size: 32px;
    }
  }

  .features-section,
  .team-section,
  .values-section {
    padding: 60px 20px;

    .section-header {
      margin-bottom: 40px;

      h2 {
        font-size: 32px;
      }

      p {
        font-size: 16px;
      }
    }
  }

  .cta-section {
    padding: 60px 20px;

    h2 {
      font-size: 32px;
    }

    p {
      font-size: 16px;
    }

    .cta-buttons {
      gap: 12px;
    }
  }
}

@media (max-width: 480px) {
  .hero {
    padding: 40px 16px;

    .hero-title {
      font-size: 28px;
    }

    .hero-subtitle {
      font-size: 16px;
    }

    .hero-description {
      font-size: 14px;
    }
  }

  .about-section {
    padding: 40px 16px;

    .about-stats {
      grid-template-columns: 1fr;
    }
  }

  .features-section,
  .team-section,
  .values-section {
    padding: 40px 16px;

    .section-header {
      margin-bottom: 30px;

      h2 {
        font-size: 24px;
      }

      p {
        font-size: 14px;
      }
    }

    .features-grid,
    .team-grid,
    .values-grid {
      gap: 16px;
    }
  }

  .cta-section {
    padding: 40px 16px;

    h2 {
      font-size: 24px;
    }

    p {
      font-size: 14px;
      margin-bottom: 24px;
    }

    .cta-buttons {
      flex-direction: column;
      gap: 12px;

      :deep(.n-button) {
        width: 100%;
      }
    }
  }
}

// Dark Mode Enhancements
html.dark {
  .hero {
    &::before {
      opacity: 0.1;
    }
  }

  .cta-section {
    &::before {
      opacity: 0.1;
    }
  }
}
</style>
