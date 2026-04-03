<script setup lang="ts">
import { computed } from 'vue'
import { useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import {
  NCard,
  NGrid,
  NGridItem,
  NIcon,
  NTag,
  NSpace,
  NButton,
} from 'naive-ui'
import {
  CodeSlash,
  Bug,
  Globe,
  Key,
  Search,
  ExtensionPuzzle,
  Eye,
  ShieldCheckmark,
} from '@vicons/ionicons5'

const { t } = useI18n()
const router = useRouter()

interface CTFCategoryCard {
  key: string
  labelKey: string
  descKey: string
  icon: any
  color: string
  problemCount: number
}

const categories = computed<CTFCategoryCard[]>(() => [
  {
    key: 'reverse',
    labelKey: 'ctf.reverse',
    descKey: 'ctf.reverseDesc',
    icon: CodeSlash,
    color: '#00d4ff',
    problemCount: 24,
  },
  {
    key: 'pwn',
    labelKey: 'ctf.pwn',
    descKey: 'ctf.pwnDesc',
    icon: Bug,
    color: '#e94560',
    problemCount: 18,
  },
  {
    key: 'web',
    labelKey: 'ctf.web',
    descKey: 'ctf.webDesc',
    icon: Globe,
    color: '#00e676',
    problemCount: 32,
  },
  {
    key: 'crypto',
    labelKey: 'ctf.crypto',
    descKey: 'ctf.cryptoDesc',
    icon: Key,
    color: '#ffab00',
    problemCount: 20,
  },
  {
    key: 'forensics',
    labelKey: 'ctf.forensics',
    descKey: 'ctf.forensicsDesc',
    icon: Search,
    color: '#ab47bc',
    problemCount: 15,
  },
  {
    key: 'misc',
    labelKey: 'ctf.misc',
    descKey: 'ctf.miscDesc',
    icon: ExtensionPuzzle,
    color: '#26c6da',
    problemCount: 28,
  },
  {
    key: 'recon',
    labelKey: 'ctf.recon',
    descKey: 'ctf.reconDesc',
    icon: Eye,
    color: '#66bb6a',
    problemCount: 12,
  },
  {
    key: 'vuln-reproduce',
    labelKey: 'ctf.vulnReproduce',
    descKey: 'ctf.vulnReproduceDesc',
    icon: ShieldCheckmark,
    color: '#ff7043',
    problemCount: 10,
  },
])

const resources = computed(() => [
  { label: t('ctf.tools'), path: '/help#ctf-tools' },
  { label: t('ctf.tutorials'), path: '/help#ctf-tutorials' },
  { label: t('ctf.knowledgeBase'), path: '/help#ctf-knowledge' },
])

function goToCategory(key: string) {
  router.push(`/ctf/${key}`)
}
</script>

<template>
  <div class="ctf-practice-page">
    <!-- Hero Section -->
    <section class="ctf-hero">
      <div class="ctf-hero-bg"></div>
      <div class="ctf-hero-content">
        <div class="ctf-hero-icon">
          <NIcon size="64" color="#00d4ff">
            <ShieldCheckmark />
          </NIcon>
        </div>
        <h1 class="ctf-hero-title">CTF Practice</h1>
        <p class="ctf-hero-subtitle">{{ t('ctf.heroSubtitle') }}</p>
        <NSpace justify="center" class="ctf-hero-actions">
          <NButton type="primary" size="large" @click="router.push('/ctf/web')">
            {{ t('ctf.startChallenge') }}
          </NButton>
          <NButton size="large" ghost @click="router.push('/help#ctf')">
            {{ t('ctf.learnMore') }}
          </NButton>
        </NSpace>
      </div>
    </section>

    <!-- Category Cards Grid -->
    <section class="ctf-section">
      <h2 class="section-title">{{ t('ctf.categories') }}</h2>
      <NGrid :x-gap="20" :y-gap="20" :cols="4" responsive="screen" item-responsive>
        <NGridItem
          v-for="cat in categories"
          :key="cat.key"
          span="4 m:2 l:1"
        >
          <NCard
            class="category-card"
            hoverable
            @click="goToCategory(cat.key)"
          >
            <div class="category-card-inner">
              <div class="category-icon" :style="{ color: cat.color }">
                <NIcon size="40">
                  <component :is="cat.icon" />
                </NIcon>
              </div>
              <div class="category-info">
                <h3 class="category-name">{{ t(cat.labelKey) }}</h3>
                <p class="category-desc">{{ t(cat.descKey) }}</p>
              </div>
              <NTag :bordered="false" size="small" round :type="'info'">
                {{ cat.problemCount }} {{ t('ctf.problems') }}
              </NTag>
            </div>
          </NCard>
        </NGridItem>
      </NGrid>
    </section>

    <!-- Resources Section -->
    <section class="ctf-section">
      <h2 class="section-title">{{ t('ctf.resources') }}</h2>
      <NGrid :x-gap="20" :y-gap="20" :cols="3" responsive="screen" item-responsive>
        <NGridItem
          v-for="res in resources"
          :key="res.label"
          span="3 m:1"
        >
          <NCard class="resource-card" hoverable>
            <div class="resource-inner">
              <h3 class="resource-title">{{ res.label }}</h3>
              <NButton text type="primary" @click="router.push(res.path)">
                {{ t('common.search') }} &rarr;
              </NButton>
            </div>
          </NCard>
        </NGridItem>
      </NGrid>
    </section>
  </div>
</template>

<style scoped lang="scss">
.ctf-practice-page {
  padding-bottom: 40px;
}

.ctf-hero {
  position: relative;
  display: flex;
  align-items: center;
  justify-content: center;
  min-height: 360px;
  overflow: hidden;
  margin-bottom: 48px;
  border-radius: 12px;
}

.ctf-hero-bg {
  position: absolute;
  inset: 0;
  background: linear-gradient(135deg, #0a0a1a 0%, #1a0a2e 25%, #0f1a3e 50%, #0a2a1a 75%, #1a1a2e 100%);
  background-size: 400% 400%;
  animation: ctfGradient 12s ease infinite;

  &::before {
    content: '';
    position: absolute;
    inset: 0;
    background:
      radial-gradient(circle at 15% 30%, rgba(0, 212, 255, 0.12) 0%, transparent 40%),
      radial-gradient(circle at 85% 70%, rgba(233, 69, 96, 0.12) 0%, transparent 40%),
      radial-gradient(circle at 50% 50%, rgba(0, 230, 118, 0.06) 0%, transparent 50%);
  }

  &::after {
    content: '';
    position: absolute;
    inset: 0;
    background-image:
      linear-gradient(rgba(0, 212, 255, 0.03) 1px, transparent 1px),
      linear-gradient(90deg, rgba(0, 212, 255, 0.03) 1px, transparent 1px);
    background-size: 40px 40px;
    animation: gridMove 20s linear infinite;
  }
}

@keyframes ctfGradient {
  0%, 100% { background-position: 0% 50%; }
  25% { background-position: 100% 0%; }
  50% { background-position: 100% 100%; }
  75% { background-position: 0% 100%; }
}

@keyframes gridMove {
  0% { transform: translate(0, 0); }
  100% { transform: translate(40px, 40px); }
}

.ctf-hero-content {
  position: relative;
  text-align: center;
  z-index: 1;
  padding: 40px 20px;
}

.ctf-hero-icon {
  margin-bottom: 16px;
}

.ctf-hero-title {
  font-size: 56px;
  font-weight: 800;
  background: linear-gradient(135deg, #00d4ff, #e94560, #00e676);
  background-size: 200% 200%;
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
  animation: textShimmer 4s ease infinite;
  margin-bottom: 12px;
  letter-spacing: -1px;
}

@keyframes textShimmer {
  0%, 100% { background-position: 0% 50%; }
  50% { background-position: 100% 50%; }
}

.ctf-hero-subtitle {
  font-size: 18px;
  color: var(--color-text-secondary);
  margin-bottom: 28px;
  max-width: 600px;
  margin-left: auto;
  margin-right: auto;
}

.ctf-hero-actions {
  margin-top: 8px;
}

.ctf-section {
  margin-bottom: 48px;
}

.section-title {
  font-size: 24px;
  font-weight: 700;
  color: var(--color-text);
  margin-bottom: 20px;
}

.category-card {
  background-color: var(--color-bg-card);
  border: 1px solid rgba(255, 255, 255, 0.06);
  cursor: pointer;
  transition: all 0.3s ease;
  height: 100%;

  &:hover {
    border-color: var(--color-primary);
    transform: translateY(-4px);
    box-shadow: 0 8px 24px rgba(0, 212, 255, 0.15);
  }

  :deep(.n-card__content) {
    padding: 24px;
  }
}

.category-card-inner {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 12px;
  text-align: center;
}

.category-icon {
  width: 64px;
  height: 64px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 16px;
  background: rgba(255, 255, 255, 0.04);
}

.category-info {
  flex: 1;
}

.category-name {
  font-size: 16px;
  font-weight: 600;
  color: var(--color-text);
  margin-bottom: 6px;
}

.category-desc {
  font-size: 13px;
  color: var(--color-text-secondary);
  line-height: 1.5;
}

.resource-card {
  background-color: var(--color-bg-card);
  border: 1px solid rgba(255, 255, 255, 0.06);
  cursor: pointer;
  transition: all 0.3s ease;

  &:hover {
    border-color: var(--color-primary);
    transform: translateY(-2px);
  }

  :deep(.n-card__content) {
    padding: 20px 24px;
  }
}

.resource-inner {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.resource-title {
  font-size: 16px;
  font-weight: 600;
  color: var(--color-text);
}

@media (max-width: 768px) {
  .ctf-hero-title {
    font-size: 36px;
  }

  .ctf-hero-subtitle {
    font-size: 15px;
  }
}
</style>
