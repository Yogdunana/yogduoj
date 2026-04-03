<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import {
  NCard,
  NGrid,
  NGridItem,
  NStatistic,
  NTag,
  NButton,
  NSpace,
  NIcon,
} from 'naive-ui'
import {
  Rocket,
  Trophy,
  ShieldCheckmark,
  People,
  TrendingUp,
  Notifications,
} from '@vicons/ionicons5'
import { listContests } from '@/api/contest'
import { listAnnouncements } from '@/api/announcement'
import type { Contest, Announcement } from '@/types'

const { t } = useI18n()
const router = useRouter()

// Animated counter logic
function useAnimatedCounter(target: number, duration = 2000) {
  const current = ref(0)
  let animationFrame: number | null = null
  let startTime: number | null = null

  function animate(timestamp: number) {
    if (!startTime) startTime = timestamp
    const progress = Math.min((timestamp - startTime) / duration, 1)
    const eased = 1 - Math.pow(1 - progress, 3) // ease-out cubic
    current.value = Math.floor(eased * target)
    if (progress < 1) {
      animationFrame = requestAnimationFrame(animate)
    }
  }

  function start() {
    current.value = 0
    startTime = null
    animationFrame = requestAnimationFrame(animate)
  }

  function stop() {
    if (animationFrame) {
      cancelAnimationFrame(animationFrame)
      animationFrame = null
    }
  }

  return { current, start, stop }
}

const usersCounter = useAnimatedCounter(1234)
const problemsCounter = useAnimatedCounter(256)
const submissionsCounter = useAnimatedCounter(45678)

const stats = computed(() => [
  {
    label: t('home.totalUsers'),
    value: usersCounter.current.value,
    icon: People,
    color: '#00d4ff',
  },
  {
    label: t('home.totalProblems'),
    value: problemsCounter.current.value,
    icon: TrendingUp,
    color: '#00e676',
  },
  {
    label: t('home.totalSubmissions'),
    value: submissionsCounter.current.value,
    icon: Rocket,
    color: '#e94560',
  },
])

const quickLinks = computed(() => [
  { label: t('home.startPracticing'), path: '/problems', icon: Rocket, color: '#00d4ff' },
  { label: t('home.viewContests'), path: '/contests', icon: Trophy, color: '#ffab00' },
  { label: t('nav.ctf'), path: '/ctf', icon: ShieldCheckmark, color: '#e94560' },
  { label: t('nav.teams'), path: '/teams', icon: People, color: '#00e676' },
  { label: t('nav.help'), path: '/help', icon: Notifications, color: '#ab47bc' },
])

// Contests data
const contests = ref<Contest[]>([])
const upcomingContests = computed(() => contests.value.filter(c => c.status === 'upcoming').slice(0, 3))
const runningContests = computed(() => contests.value.filter(c => c.status === 'running').slice(0, 3))
const recentContests = computed(() => contests.value.filter(c => c.status === 'ended').slice(0, 3))

// Announcements data
const announcements = ref<Announcement[]>([])

function formatTime(dateStr: string): string {
  const d = new Date(dateStr)
  return d.toLocaleString()
}

function formatShortDate(dateStr: string): string {
  const d = new Date(dateStr)
  return d.toLocaleDateString()
}

async function fetchData() {
  try {
    const [contestRes, announcementRes] = await Promise.all([
      listContests({ page: 1, pageSize: 10 }),
      listAnnouncements({ page: 1, pageSize: 5 }),
    ])
    contests.value = contestRes.data.data.items
    announcements.value = announcementRes.data.data.items
  } catch {
    // Silently fail - use placeholder data
  }
}

function goToContest(id: number) {
  router.push({ name: 'ContestDetail', params: { id } })
}

function goToAnnouncement(id: number) {
  router.push({ name: 'AnnouncementDetail', params: { id } })
}

onMounted(() => {
  fetchData()
  // Start counter animations after a short delay
  setTimeout(() => {
    usersCounter.start()
    problemsCounter.start()
    submissionsCounter.start()
  }, 300)
})

onUnmounted(() => {
  usersCounter.stop()
  problemsCounter.stop()
  submissionsCounter.stop()
})
</script>

<template>
  <div class="home-page">
    <!-- Hero Section with Animated Gradient -->
    <section class="hero">
      <div class="hero-bg"></div>
      <div class="hero-particles"></div>
      <div class="hero-content">
        <h1 class="hero-title">YogduOJ</h1>
        <p class="hero-subtitle">{{ t('home.subtitle') }}</p>
        <div class="hero-actions">
          <NButton type="primary" size="large" @click="router.push('/problems')">
            <template #icon>
              <NIcon><Rocket /></NIcon>
            </template>
            {{ t('home.startPracticing') }}
          </NButton>
          <NButton size="large" ghost @click="router.push('/contests')">
            <template #icon>
              <NIcon><Trophy /></NIcon>
            </template>
            {{ t('home.viewContests') }}
          </NButton>
        </div>
      </div>
    </section>

    <!-- Statistics Section with Animated Counters -->
    <section class="section">
      <div class="container">
        <NGrid :x-gap="20" :y-gap="20" :cols="3" responsive="screen" item-responsive>
          <NGridItem span="3 m:1" v-for="stat in stats" :key="stat.label">
            <NCard class="stat-card">
              <div class="stat-inner">
                <div class="stat-icon" :style="{ color: stat.color }">
                  <NIcon size="36">
                    <component :is="stat.icon" />
                  </NIcon>
                </div>
                <NStatistic :value="stat.value" class="stat-value">
                  <template #label>
                    <span class="stat-label">{{ stat.label }}</span>
                  </template>
                </NStatistic>
              </div>
            </NCard>
          </NGridItem>
        </NGrid>
      </div>
    </section>

    <!-- Quick Links Grid -->
    <section class="section">
      <div class="container">
        <h2 class="section-title">{{ t('home.quickLinks') }}</h2>
        <NGrid :x-gap="16" :y-gap="16" :cols="5" responsive="screen" item-responsive>
          <NGridItem
            v-for="link in quickLinks"
            :key="link.path"
            span="5 m:2 l:1"
          >
            <NCard class="quick-link-card" hoverable @click="router.push(link.path)">
              <div class="quick-link-inner">
                <div class="quick-link-icon" :style="{ color: link.color }">
                  <NIcon size="32">
                    <component :is="link.icon" />
                  </NIcon>
                </div>
                <span class="quick-link-label">{{ link.label }}</span>
              </div>
            </NCard>
          </NGridItem>
        </NGrid>
      </div>
    </section>

    <!-- Contests Section -->
    <section class="section">
      <div class="container">
        <h2 class="section-title">{{ t('home.recentContests') }}</h2>

        <!-- Running Contests -->
        <div v-if="runningContests.length > 0" class="contest-group">
          <h3 class="contest-group-title">
            <NTag type="success" size="small" :bordered="false">{{ t('contests.running') }}</NTag>
          </h3>
          <div class="contest-list">
            <NCard
              v-for="contest in runningContests"
              :key="contest.id"
              class="contest-card"
              hoverable
              @click="goToContest(contest.id)"
            >
              <div class="contest-card-inner">
                <div class="contest-info">
                  <span class="contest-title">{{ contest.title }}</span>
                  <span class="contest-meta">{{ formatTime(contest.startTime) }}</span>
                </div>
                <div class="contest-meta-right">
                  <NTag size="small" :bordered="false">{{ contest.participantCount }} {{ t('contests.participants') }}</NTag>
                  <NButton type="primary" size="small">{{ t('contests.enterContest') }}</NButton>
                </div>
              </div>
            </NCard>
          </div>
        </div>

        <!-- Upcoming Contests -->
        <div v-if="upcomingContests.length > 0" class="contest-group">
          <h3 class="contest-group-title">
            <NTag type="info" size="small" :bordered="false">{{ t('contests.upcoming') }}</NTag>
          </h3>
          <div class="contest-list">
            <NCard
              v-for="contest in upcomingContests"
              :key="contest.id"
              class="contest-card"
              hoverable
              @click="goToContest(contest.id)"
            >
              <div class="contest-card-inner">
                <div class="contest-info">
                  <span class="contest-title">{{ contest.title }}</span>
                  <span class="contest-meta">{{ formatTime(contest.startTime) }}</span>
                </div>
                <div class="contest-meta-right">
                  <NTag size="small" :bordered="false">{{ contest.participantCount }} {{ t('contests.participants') }}</NTag>
                  <NButton size="small">{{ t('contests.registerContest') }}</NButton>
                </div>
              </div>
            </NCard>
          </div>
        </div>

        <!-- Recent (Ended) Contests -->
        <div v-if="recentContests.length > 0" class="contest-group">
          <h3 class="contest-group-title">
            <NTag size="small" :bordered="false">{{ t('contests.ended') }}</NTag>
          </h3>
          <div class="contest-list">
            <NCard
              v-for="contest in recentContests"
              :key="contest.id"
              class="contest-card"
              hoverable
              @click="goToContest(contest.id)"
            >
              <div class="contest-card-inner">
                <div class="contest-info">
                  <span class="contest-title">{{ contest.title }}</span>
                  <span class="contest-meta">{{ formatTime(contest.startTime) }}</span>
                </div>
                <div class="contest-meta-right">
                  <NTag size="small" :bordered="false">{{ contest.participantCount }} {{ t('contests.participants') }}</NTag>
                </div>
              </div>
            </NCard>
          </div>
        </div>

        <div v-if="contests.length === 0" class="empty-state">
          <p>{{ t('contests.noContests') }}</p>
        </div>
      </div>
    </section>

    <!-- Announcements Section -->
    <section class="section">
      <div class="container">
        <h2 class="section-title">{{ t('home.announcements') }}</h2>
        <div v-if="announcements.length > 0" class="announcement-list">
          <NCard
            v-for="item in announcements"
            :key="item.id"
            class="announcement-card"
            hoverable
            @click="goToAnnouncement(item.id)"
          >
            <div class="announcement-inner">
              <div class="announcement-info">
                <div class="announcement-title-row">
                  <NTag v-if="item.isPinned" type="error" size="small" :bordered="false" class="pin-tag">
                    {{ t('announcements.pinned') }}
                  </NTag>
                  <span class="announcement-title">{{ item.title }}</span>
                </div>
                <span class="announcement-date">{{ formatShortDate(item.createdAt) }}</span>
              </div>
              <NButton text type="primary">{{ t('announcements.readMore') }} &rarr;</NButton>
            </div>
          </NCard>
        </div>
        <div v-else class="empty-state">
          <p>{{ t('announcements.noAnnouncements') }}</p>
        </div>
      </div>
    </section>

    <!-- Footer Branding -->
    <section class="section branding-section">
      <div class="container">
        <p class="branding-text">Powered by <strong>YogduOJ</strong></p>
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
  min-height: 440px;
  overflow: hidden;
  margin-bottom: 48px;
}

.hero-bg {
  position: absolute;
  inset: 0;
  background: linear-gradient(135deg, #1a1a2e 0%, #16213e 20%, #0f3460 40%, #1a0a2e 60%, #0a2a1a 80%, #1a1a2e 100%);
  background-size: 400% 400%;
  animation: heroGradient 10s ease infinite;

  &::before {
    content: '';
    position: absolute;
    inset: 0;
    background: radial-gradient(circle at 20% 50%, rgba(0, 212, 255, 0.15) 0%, transparent 50%),
                radial-gradient(circle at 80% 50%, rgba(233, 69, 96, 0.15) 0%, transparent 50%),
                radial-gradient(circle at 50% 20%, rgba(0, 230, 118, 0.08) 0%, transparent 40%);
  }
}

@keyframes heroGradient {
  0%, 100% { background-position: 0% 50%; }
  25% { background-position: 100% 0%; }
  50% { background-position: 50% 100%; }
  75% { background-position: 0% 100%; }
}

.hero-particles {
  position: absolute;
  inset: 0;
  background-image:
    radial-gradient(1px 1px at 10% 20%, rgba(0, 212, 255, 0.4) 0%, transparent 100%),
    radial-gradient(1px 1px at 30% 60%, rgba(233, 69, 96, 0.3) 0%, transparent 100%),
    radial-gradient(1px 1px at 50% 30%, rgba(0, 230, 118, 0.3) 0%, transparent 100%),
    radial-gradient(1px 1px at 70% 70%, rgba(0, 212, 255, 0.4) 0%, transparent 100%),
    radial-gradient(1px 1px at 90% 40%, rgba(255, 171, 0, 0.3) 0%, transparent 100%),
    radial-gradient(1px 1px at 20% 80%, rgba(171, 71, 188, 0.3) 0%, transparent 100%),
    radial-gradient(1px 1px at 60% 10%, rgba(0, 212, 255, 0.3) 0%, transparent 100%),
    radial-gradient(1px 1px at 80% 90%, rgba(233, 69, 96, 0.3) 0%, transparent 100%);
  animation: particleFloat 15s ease-in-out infinite;
}

@keyframes particleFloat {
  0%, 100% { opacity: 0.6; transform: translateY(0); }
  50% { opacity: 1; transform: translateY(-10px); }
}

.hero-content {
  position: relative;
  text-align: center;
  z-index: 1;
  padding: 40px 20px;
}

.hero-title {
  font-size: 72px;
  font-weight: 800;
  background: linear-gradient(135deg, var(--color-primary), var(--color-accent), var(--color-success));
  background-size: 200% 200%;
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
  animation: titleShimmer 5s ease infinite;
  margin-bottom: 16px;
  letter-spacing: -2px;
}

@keyframes titleShimmer {
  0%, 100% { background-position: 0% 50%; }
  50% { background-position: 100% 50%; }
}

.hero-subtitle {
  font-size: 20px;
  color: var(--color-text-secondary);
  margin-bottom: 36px;
}

.hero-actions {
  display: flex;
  gap: 16px;
  justify-content: center;
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

// Stat Cards
.stat-card {
  background-color: var(--color-bg-card);
  border: 1px solid rgba(255, 255, 255, 0.06);
  transition: all 0.3s ease;

  &:hover {
    border-color: var(--color-primary);
    transform: translateY(-2px);
  }

  :deep(.n-card__content) {
    padding: 24px;
  }
}

.stat-inner {
  display: flex;
  align-items: center;
  gap: 16px;
}

.stat-icon {
  width: 56px;
  height: 56px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 12px;
  background: rgba(255, 255, 255, 0.04);
  flex-shrink: 0;
}

.stat-label {
  font-size: 14px;
  color: var(--color-text-secondary);
}

.stat-value {
  :deep(.n-statistic-value__content) {
    font-size: 32px;
    font-weight: 700;
    color: var(--color-text);
  }
}

// Quick Links
.quick-link-card {
  background-color: var(--color-bg-card);
  border: 1px solid rgba(255, 255, 255, 0.06);
  cursor: pointer;
  transition: all 0.3s ease;

  &:hover {
    border-color: var(--color-primary);
    transform: translateY(-4px);
    box-shadow: 0 8px 24px rgba(0, 0, 0, 0.3);
  }

  :deep(.n-card__content) {
    padding: 24px 16px;
  }
}

.quick-link-inner {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 12px;
  text-align: center;
}

.quick-link-icon {
  width: 56px;
  height: 56px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 14px;
  background: rgba(255, 255, 255, 0.04);
}

.quick-link-label {
  font-size: 14px;
  font-weight: 500;
  color: var(--color-text);
}

// Contest Section
.contest-group {
  margin-bottom: 24px;

  &:last-child {
    margin-bottom: 0;
  }
}

.contest-group-title {
  margin-bottom: 12px;
}

.contest-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.contest-card {
  background-color: var(--color-bg-card);
  border: 1px solid rgba(255, 255, 255, 0.06);
  cursor: pointer;
  transition: all 0.2s ease;

  &:hover {
    border-color: var(--color-primary);
  }

  :deep(.n-card__content) {
    padding: 16px 20px;
  }
}

.contest-card-inner {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 16px;
}

.contest-info {
  display: flex;
  flex-direction: column;
  gap: 4px;
  min-width: 0;
}

.contest-title {
  font-size: 15px;
  font-weight: 500;
  color: var(--color-text);
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.contest-meta {
  font-size: 13px;
  color: var(--color-text-secondary);
}

.contest-meta-right {
  display: flex;
  align-items: center;
  gap: 12px;
  flex-shrink: 0;
}

// Announcements
.announcement-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.announcement-card {
  background-color: var(--color-bg-card);
  border: 1px solid rgba(255, 255, 255, 0.06);
  cursor: pointer;
  transition: all 0.2s ease;

  &:hover {
    border-color: var(--color-primary);
  }

  :deep(.n-card__content) {
    padding: 16px 20px;
  }
}

.announcement-inner {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 16px;
}

.announcement-info {
  display: flex;
  flex-direction: column;
  gap: 4px;
  min-width: 0;
}

.announcement-title-row {
  display: flex;
  align-items: center;
  gap: 8px;
}

.pin-tag {
  flex-shrink: 0;
}

.announcement-title {
  font-size: 15px;
  font-weight: 500;
  color: var(--color-text);
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.announcement-date {
  font-size: 13px;
  color: var(--color-text-secondary);
}

.empty-state {
  text-align: center;
  padding: 32px;
  color: var(--color-text-secondary);
  font-size: 14px;
}

// Branding
.branding-section {
  text-align: center;
  padding: 24px 0;
  border-top: 1px solid rgba(255, 255, 255, 0.06);
}

.branding-text {
  font-size: 14px;
  color: var(--color-text-secondary);

  strong {
    color: var(--color-primary);
    font-weight: 600;
  }
}

@media (max-width: 768px) {
  .hero {
    min-height: 320px;
  }

  .hero-title {
    font-size: 48px;
  }

  .hero-subtitle {
    font-size: 16px;
  }

  .hero-actions {
    flex-direction: column;
    align-items: center;
  }

  .contest-card-inner {
    flex-direction: column;
    align-items: flex-start;
  }

  .contest-meta-right {
    width: 100%;
    justify-content: flex-end;
  }

  .announcement-inner {
    flex-direction: column;
    align-items: flex-start;
  }
}
</style>
