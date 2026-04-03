<script setup lang="ts">
import { useI18n } from 'vue-i18n'

const { t } = useI18n()

const stats = [
  { label: t('home.totalUsers'), value: '1,234', icon: '👥' },
  { label: t('home.totalProblems'), value: '256', icon: '📝' },
  { label: t('home.totalSubmissions'), value: '45,678', icon: '📊' },
]

const quickLinks = [
  { label: t('home.startPracticing'), path: '/problems', icon: '🚀' },
  { label: t('home.viewContests'), path: '/contests', icon: '🏆' },
  { label: t('nav.ctf'), path: '/ctf', icon: '🛡️' },
  { label: t('nav.ranking'), path: '/contests', icon: '📈' },
  { label: t('nav.teams'), path: '/teams', icon: '🤝' },
  { label: t('nav.help'), path: '/help', icon: '❓' },
]

const recentContests = [
  { id: 1, title: 'Weekly Contest #42', status: 'upcoming', startTime: '2026-04-10 14:00' },
  { id: 2, title: 'Monthly Challenge', status: 'running', startTime: '2026-04-01 10:00' },
  { id: 3, title: 'Beginner Contest #20', status: 'ended', startTime: '2026-03-28 09:00' },
]

const announcements = [
  { id: 1, title: 'System Maintenance Notice', date: '2026-04-01' },
  { id: 2, title: 'New CTF Challenges Available', date: '2026-03-30' },
  { id: 3, title: 'Contest Rules Update', date: '2026-03-25' },
]
</script>

<template>
  <div class="home-page">
    <!-- Hero Section -->
    <section class="hero">
      <div class="hero-bg"></div>
      <div class="hero-content">
        <h1 class="hero-title">YogduOJ</h1>
        <p class="hero-subtitle">{{ t('home.subtitle') }}</p>
        <div class="hero-actions">
          <router-link to="/problems" class="btn btn-primary btn-lg">
            {{ t('home.startPracticing') }}
          </router-link>
          <router-link to="/contests" class="btn btn-secondary btn-lg">
            {{ t('home.viewContests') }}
          </router-link>
        </div>
      </div>
    </section>

    <!-- Stats Section -->
    <section class="section">
      <div class="container">
        <div class="stats-grid">
          <div v-for="stat in stats" :key="stat.label" class="stat-card card">
            <span class="stat-icon">{{ stat.icon }}</span>
            <span class="stat-value">{{ stat.value }}</span>
            <span class="stat-label">{{ stat.label }}</span>
          </div>
        </div>
      </div>
    </section>

    <!-- Quick Links -->
    <section class="section">
      <div class="container">
        <h2 class="section-title">{{ t('home.quickLinks') }}</h2>
        <div class="quick-links-grid">
          <router-link
            v-for="link in quickLinks"
            :key="link.path"
            :to="link.path"
            class="quick-link card"
          >
            <span class="link-icon">{{ link.icon }}</span>
            <span class="link-label">{{ link.label }}</span>
          </router-link>
        </div>
      </div>
    </section>

    <!-- Recent Contests & Announcements -->
    <section class="section">
      <div class="container">
        <div class="two-column">
          <!-- Recent Contests -->
          <div class="column">
            <h2 class="section-title">{{ t('home.recentContests') }}</h2>
            <div class="card">
              <div v-for="contest in recentContests" :key="contest.id" class="list-item">
                <div class="list-item-info">
                  <span class="list-item-title">{{ contest.title }}</span>
                  <span class="list-item-meta">{{ contest.startTime }}</span>
                </div>
                <span :class="['status-badge', contest.status]">
                  {{ t(`contests.${contest.status}`) }}
                </span>
              </div>
            </div>
          </div>

          <!-- Announcements -->
          <div class="column">
            <h2 class="section-title">{{ t('home.announcements') }}</h2>
            <div class="card">
              <div v-for="item in announcements" :key="item.id" class="list-item">
                <div class="list-item-info">
                  <span class="list-item-title">{{ item.title }}</span>
                  <span class="list-item-meta">{{ item.date }}</span>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </section>
  </div>
</template>

<style scoped lang="scss">
.home-page {
  padding-bottom: 40px;
}

.hero {
  position: relative;
  display: flex;
  align-items: center;
  justify-content: center;
  min-height: 400px;
  overflow: hidden;
  margin-bottom: 40px;
}

.hero-bg {
  position: absolute;
  inset: 0;
  background: linear-gradient(135deg, #1a1a2e 0%, #16213e 30%, #0f3460 60%, #1a1a2e 100%);
  background-size: 300% 300%;
  animation: gradientShift 8s ease infinite;

  &::before {
    content: '';
    position: absolute;
    inset: 0;
    background: radial-gradient(circle at 20% 50%, rgba(0, 212, 255, 0.15) 0%, transparent 50%),
                radial-gradient(circle at 80% 50%, rgba(233, 69, 96, 0.15) 0%, transparent 50%);
  }
}

@keyframes gradientShift {
  0%, 100% { background-position: 0% 50%; }
  50% { background-position: 100% 50%; }
}

.hero-content {
  position: relative;
  text-align: center;
  z-index: 1;
}

.hero-title {
  font-size: 64px;
  font-weight: 800;
  background: linear-gradient(135deg, var(--color-primary), var(--color-accent));
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
  margin-bottom: 16px;
  letter-spacing: -2px;
}

.hero-subtitle {
  font-size: 20px;
  color: var(--color-text-secondary);
  margin-bottom: 32px;
}

.hero-actions {
  display: flex;
  gap: 16px;
  justify-content: center;
}

.btn-lg {
  padding: 12px 32px;
  font-size: 16px;
  border-radius: 12px;
}

.section {
  margin-bottom: 48px;
}

.section-title {
  font-size: 24px;
  font-weight: 700;
  color: var(--color-text);
  margin-bottom: 20px;
}

.stats-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
  gap: 16px;
}

.stat-card {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 8px;
  padding: 32px 24px;
  text-align: center;
}

.stat-icon {
  font-size: 32px;
}

.stat-value {
  font-size: 32px;
  font-weight: 700;
  color: var(--color-primary);
}

.stat-label {
  font-size: 14px;
  color: var(--color-text-secondary);
}

.quick-links-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(160px, 1fr));
  gap: 16px;
}

.quick-link {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 12px;
  padding: 24px 16px;
  text-align: center;
  text-decoration: none;
  color: var(--color-text);
  transition: all 0.2s ease;

  &:hover {
    transform: translateY(-4px);
    border-color: var(--color-primary);
    color: var(--color-primary);
  }
}

.link-icon {
  font-size: 36px;
}

.link-label {
  font-size: 14px;
  font-weight: 500;
}

.two-column {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 24px;
}

.column {
  min-width: 0;
}

.list-item {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 12px 0;
  border-bottom: 1px solid rgba(255, 255, 255, 0.06);

  &:last-child {
    border-bottom: none;
  }
}

.list-item-info {
  display: flex;
  flex-direction: column;
  gap: 4px;
  min-width: 0;
}

.list-item-title {
  font-size: 14px;
  font-weight: 500;
  color: var(--color-text);
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.list-item-meta {
  font-size: 12px;
  color: var(--color-text-secondary);
}

.status-badge {
  flex-shrink: 0;
  padding: 2px 10px;
  border-radius: 12px;
  font-size: 12px;
  font-weight: 500;

  &.upcoming {
    background-color: rgba(0, 212, 255, 0.15);
    color: var(--color-primary);
  }

  &.running {
    background-color: rgba(0, 230, 118, 0.15);
    color: var(--color-success);
  }

  &.ended {
    background-color: rgba(160, 160, 160, 0.15);
    color: var(--color-text-secondary);
  }
}

@media (max-width: 768px) {
  .hero-title {
    font-size: 40px;
  }

  .hero-subtitle {
    font-size: 16px;
  }

  .hero-actions {
    flex-direction: column;
    align-items: center;
  }

  .two-column {
    grid-template-columns: 1fr;
  }
}
</style>
