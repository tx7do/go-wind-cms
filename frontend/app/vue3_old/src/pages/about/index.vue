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
  background: linear-gradient(135deg, #6366f1 0%, #8b5cf6 50%, #a855f7 100%);
  color: white;
  padding: 3rem 2rem;
  text-align: center;
  position: relative;
  overflow: hidden;
  min-height: 350px;
  height: auto;
  display: flex;
  align-items: center;
  justify-content: center;

  // 背景装饰：渐变和网格
  &::before {
    content: '';
    position: absolute;
    top: 0;
    left: 0;
    right: 0;
    bottom: 0;
    background:
      radial-gradient(ellipse at 20% 50%, rgba(100, 200, 255, 0.15) 0%, transparent 50%),
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
    background-image:
      linear-gradient(rgba(255, 255, 255, 0.05) 1px, transparent 1px),
      linear-gradient(90deg, rgba(255, 255, 255, 0.05) 1px, transparent 1px);
    background-size: 50px 50px;
    animation: gridMove 20s linear infinite;
    opacity: 0.6;
    z-index: 1;
  }

  .hero-content {
    max-width: 900px;
    margin: 0 auto;
    position: relative;
    z-index: 2;
    animation: fadeInUp 0.8s ease-out;
  }

  .hero-title {
    font-size: 48px;
    font-weight: 800;
    margin: 0 0 20px 0;
    line-height: 1.2;
    letter-spacing: -1px;

    // 添加发光效果
    background: linear-gradient(
      135deg,
      #ffffff 0%,
      #f0f0ff 25%,
      #e0e0ff 50%,
      #f0f0ff 75%,
      #ffffff 100%
    );
    background-size: 200% auto;
    -webkit-background-clip: text;
    -webkit-text-fill-color: transparent;
    background-clip: text;

    text-shadow:
      0 0 40px rgba(255, 255, 255, 0.9),
      0 0 80px rgba(var(--color-primary-purple-rgb), 0.7),
      0 0 120px rgba(99, 102, 241, 0.5),
      0 6px 24px rgba(0, 0, 0, 0.5),
      0 3px 12px rgba(0, 0, 0, 0.4);

    filter:
      drop-shadow(0 0 30px rgba(var(--color-primary-purple-rgb), 0.6))
      drop-shadow(0 0 60px rgba(99, 102, 241, 0.4))
      drop-shadow(0 5px 15px rgba(0, 0, 0, 0.3));

    animation-name: slideDown, glowPulseTitle, gradientShine;
    animation-duration: 0.8s, 3s, 6s;
    animation-timing-function: cubic-bezier(0.34, 1.56, 0.64, 1), ease-in-out, linear;
    animation-iteration-count: 1, infinite, infinite;
    animation-delay: 0s, 0.5s, 0s;
  }

  .hero-subtitle {
    font-size: 20px;
    margin: 0 0 16px 0;
    opacity: 0.95;
    font-weight: 600;
    text-shadow: 0 2px 10px rgba(0, 0, 0, 0.25);
    animation: fadeInUp 0.8s ease-out 0.2s both;
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
    text-shadow: 0 2px 10px rgba(0, 0, 0, 0.25);
    animation: fadeInUp 0.8s ease-out 0.3s both;
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
    background:
      radial-gradient(ellipse at 20% 50%, rgba(100, 200, 255, 0.15) 0%, transparent 50%),
      radial-gradient(ellipse at 80% 80%, rgba(200, 100, 255, 0.12) 0%, transparent 50%),
      linear-gradient(135deg, #6366f1 0%, #8b5cf6 50%, #a855f7 100%);
  }
  50% {
    background:
      radial-gradient(ellipse at 80% 50%, rgba(100, 200, 255, 0.15) 0%, transparent 50%),
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
    filter:
      drop-shadow(0 0 25px rgba(var(--color-primary-purple-rgb), 0.5))
      drop-shadow(0 0 50px rgba(99, 102, 241, 0.3))
      drop-shadow(0 5px 15px rgba(0, 0, 0, 0.3));
  }
  50% {
    filter:
      drop-shadow(0 0 40px rgba(var(--color-primary-purple-rgb), 0.8))
      drop-shadow(0 0 80px rgba(99, 102, 241, 0.6))
      drop-shadow(0 8px 20px rgba(0, 0, 0, 0.4));
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

// Responsive Design
@media (max-width: 1024px) {
  .hero {
    padding: 3rem 2rem;
    min-height: 300px;

    .hero-title {
      font-size: 2.5rem;
      letter-spacing: -0.8px;
    }

    .hero-subtitle {
      font-size: 1.2rem;
    }

    .hero-description {
      font-size: 1rem;
    }
  }

  .about-section {
    padding: 60px 24px;

    .section-container {
      gap: 40px;
    }

    .about-content h2 {
      font-size: 36px;
    }
  }

  .features-section,
  .team-section,
  .values-section {
    padding: 60px 24px;
  }

  .section-header h2 {
    font-size: 36px;
  }

  .cta-section {
    padding: 60px 24px;

    h2 {
      font-size: 36px;
    }
  }
}

@media (max-width: 768px) {
  .hero {
    padding: 2.5rem 1.5rem;
    min-height: 280px;

    &::after {
      background-size: 30px 30px;
    }

    .hero-title {
      font-size: 2rem;
      margin-bottom: 1rem;
      letter-spacing: -0.5px;
    }

    .hero-subtitle {
      font-size: 1.1rem;
      margin-bottom: 0.875rem;
    }

    .hero-description {
      font-size: 0.95rem;
      line-height: 1.6;
    }
  }

  .about-section {
    padding: 50px 20px;

    .section-container {
      grid-template-columns: 1fr;
      gap: 40px;
    }

    .about-content {
      h2 {
        font-size: 28px;
        margin-bottom: 24px;
      }

      p {
        font-size: 15px;
        margin-bottom: 20px;
      }
    }

    .about-stats {
      grid-template-columns: repeat(2, 1fr);
      gap: 16px;

      .stat-card {
        padding: 24px 20px;

        .stat-number {
          font-size: 32px;
        }

        .stat-label {
          font-size: 13px;
        }
      }
    }
  }

  .features-section,
  .team-section,
  .values-section {
    padding: 50px 20px;

    .section-header {
      margin-bottom: 40px;

      h2 {
        font-size: 28px;
        margin-bottom: 16px;
      }

      p {
        font-size: 16px;
      }
    }

    .features-grid {
      grid-template-columns: 1fr;
      gap: 20px;

      .feature-card {
        padding: 32px 24px;

        .feature-icon {
          width: 70px;
          height: 70px;
          margin-bottom: 20px;
        }

        h3 {
          font-size: 18px;
        }

        p {
          font-size: 14px;
        }
      }
    }

    .team-grid {
      grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
      gap: 24px;

      .team-card {
        padding: 28px 20px;

        .team-avatar {
          width: 100px;
          height: 100px;
        }

        h3 {
          font-size: 18px;
        }
      }
    }

    .values-grid {
      grid-template-columns: 1fr;
      gap: 20px;

      .value-card {
        padding: 32px 24px;

        h3 {
          font-size: 18px;
        }

        p {
          font-size: 14px;
        }
      }
    }
  }

  .cta-section {
    padding: 50px 20px;

    h2 {
      font-size: 28px;
      margin-bottom: 16px;
    }

    p {
      font-size: 16px;
      margin-bottom: 32px;
    }

    .cta-buttons {
      gap: 12px;

      :deep(.n-button) {
        padding: 10px 28px;
        font-size: 15px;
      }
    }
  }
}

@media (max-width: 480px) {
  .hero {
    padding: 2rem 1rem;
    min-height: 240px;

    &::after {
      background-size: 20px 20px;
    }

    .hero-content {
      padding: 0;
    }

    .hero-title {
      font-size: 1.6rem;
      margin-bottom: 0.75rem;
      letter-spacing: -0.3px;
    }

    .hero-subtitle {
      font-size: 1rem;
      margin-bottom: 0.75rem;
    }

    .hero-description {
      font-size: 0.875rem;
      line-height: 1.5;
    }
  }

  .about-section {
    padding: 40px 16px;

    .section-container {
      gap: 32px;
      padding: 0;
    }

    .about-content {
      h2 {
        font-size: 24px;
        margin-bottom: 20px;
      }

      p {
        font-size: 14px;
        line-height: 1.7;
        margin-bottom: 16px;
      }
    }

    .about-stats {
      grid-template-columns: 1fr;
      gap: 12px;

      .stat-card {
        padding: 20px 16px;

        &:hover {
          transform: translateY(-3px);
        }

        .stat-number {
          font-size: 28px;
          margin-bottom: 8px;
        }

        .stat-label {
          font-size: 12px;
        }
      }
    }
  }

  .features-section,
  .team-section,
  .values-section {
    padding: 40px 16px;

    .section-container {
      padding: 0;
    }

    .section-header {
      margin-bottom: 32px;

      h2 {
        font-size: 24px;
        margin-bottom: 12px;
      }

      p {
        font-size: 14px;
        line-height: 1.6;
      }
    }

    .features-grid {
      grid-template-columns: 1fr;
      gap: 16px;

      .feature-card {
        padding: 28px 20px;

        &:hover {
          transform: translateY(-6px);
        }

        .feature-icon {
          width: 60px;
          height: 60px;
          margin-bottom: 16px;
          border-radius: 10px;

          :deep(svg) {
            width: 28px;
            height: 28px;
          }
        }

        h3 {
          font-size: 17px;
          margin-bottom: 12px;
        }

        p {
          font-size: 14px;
          line-height: 1.6;
        }
      }
    }

    .team-grid {
      grid-template-columns: 1fr;
      gap: 16px;

      .team-card {
        padding: 24px 20px;

        &:hover {
          transform: translateY(-6px);
        }

        .team-avatar {
          width: 90px;
          height: 90px;
          margin-bottom: 16px;
          border-width: 3px;
        }

        h3 {
          font-size: 17px;
          margin-bottom: 6px;
        }

        p {
          font-size: 13px;
        }
      }
    }

    .values-grid {
      grid-template-columns: 1fr;
      gap: 16px;

      .value-card {
        padding: 28px 20px;

        &:hover {
          transform: translateY(-4px);
        }

        h3 {
          font-size: 17px;
          margin-bottom: 12px;
        }

        p {
          font-size: 14px;
          line-height: 1.6;
        }
      }
    }
  }

  .cta-section {
    padding: 40px 16px;

    .section-container {
      padding: 0;
    }

    h2 {
      font-size: 24px;
      margin-bottom: 12px;
    }

    p {
      font-size: 14px;
      line-height: 1.6;
      margin-bottom: 28px;
    }

    .cta-buttons {
      flex-direction: column;
      gap: 12px;

      :deep(.n-button) {
        width: 100%;
        padding: 12px 24px;
        font-size: 15px;
        border-radius: 8px;
      }
    }
  }
}

// Dark Mode Enhancements
html.dark {
  .hero {
    background: linear-gradient(135deg, #6366f1 0%, #8b5cf6 50%, #a855f7 100%);
    box-shadow: 0 20px 60px rgba(99, 102, 241, 0.3);

    &::before {
      opacity: 1;
    }

    .hero-title {
      background: linear-gradient(
        135deg,
        #ffffff 0%,
        #f5f5ff 25%,
        #ebebff 50%,
        #f5f5ff 75%,
        #ffffff 100%
      );
      background-size: 200% auto;
      -webkit-background-clip: text;
      -webkit-text-fill-color: transparent;
      background-clip: text;

      text-shadow:
        0 0 50px rgba(255, 255, 255, 1),
        0 0 100px rgba(var(--color-primary-purple-rgb), 0.8),
        0 0 150px rgba(99, 102, 241, 0.6),
        0 8px 30px rgba(0, 0, 0, 0.6),
        0 4px 15px rgba(0, 0, 0, 0.4);

      filter:
        drop-shadow(0 0 35px rgba(var(--color-primary-purple-rgb), 0.7))
        drop-shadow(0 0 70px rgba(99, 102, 241, 0.5))
        drop-shadow(0 6px 20px rgba(0, 0, 0, 0.4));
    }
  }

  .cta-section {
    &::before {
      opacity: 0.1;
    }
  }
}
</style>
