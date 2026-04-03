<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import {
  NCard, NStatistic, NTag, NButton, NDataTable, NTabs, NTabPane,
  NSpace, NModal, NEmpty, NSpin, NGrid, NGridItem, useMessage,
} from 'naive-ui'
import { getContest, signupContest, withdrawContest, getContestProblems, getContestRanking } from '@/api/contest'
import { useAuthStore } from '@/stores/auth'
import type { Contest, Problem, ContestRanking, ContestStatus } from '@/types'

const props = defineProps<{ id: string }>()
const route = useRoute()
const router = useRouter()
const { t } = useI18n()
const message = useMessage()
const authStore = useAuthStore()

const loading = ref(true)
const contest = ref<Contest | null>(null)
const problems = ref<Problem[]>([])
const rankings = ref<ContestRanking[]>([])
const isSignedUp = ref(false)
const signupLoading = ref(false)
const withdrawLoading = ref(false)
const showProblemModal = ref(false)
const selectedProblem = ref<Problem | null>(null)
const activeTab = ref('problems')
const countdownTarget = ref<number>(0)
const countdownLabel = ref('')
let countdownTimer: ReturnType<typeof setInterval> | null = null
const remainingTime = ref('')

type NTagType = 'default' | 'error' | 'info' | 'primary' | 'success' | 'warning'

const statusColorMap: Record<string, NTagType> = {
  upcoming: 'info',
  running: 'success',
  ended: 'default',
}

const ruleTypeColorMap: Record<string, NTagType> = {
  acm: 'error',
  oi: 'success',
  cf: 'info',
}

const problemLabels = computed(() => {
  const labels = 'ABCDEFGHIJKLMNOPQRSTUVWXYZ'
  return problems.value.map((_, i) => labels[i] || `P${i + 1}`)
})

const problemStatusMap = ref<Record<number, 'solved' | 'attempted' | 'unsolved'>>({})

function formatTime(dateStr: string): string {
  return new Date(dateStr).toLocaleString()
}

function getDuration(start: string, end: string): string {
  const diff = new Date(end).getTime() - new Date(start).getTime()
  const hours = Math.floor(diff / (1000 * 60 * 60))
  const minutes = Math.floor((diff % (1000 * 60 * 60)) / (1000 * 60))
  if (hours >= 24) {
    const days = Math.floor(hours / 24)
    return `${days}d ${hours % 24}h ${minutes}m`
  }
  return `${hours}h ${minutes}m`
}

function updateCountdown() {
  if (!contest.value) return
  const now = Date.now()
  const start = new Date(contest.value.startTime).getTime()
  const end = new Date(contest.value.endTime).getTime()

  if (now < start) {
    countdownLabel.value = t('contests.countdownStart')
    const diff = start - now
    remainingTime.value = formatCountdown(diff)
  } else if (now < end) {
    countdownLabel.value = t('contests.countdownEnd')
    const diff = end - now
    remainingTime.value = formatCountdown(diff)
  } else {
    remainingTime.value = '00:00:00'
    if (countdownTimer) {
      clearInterval(countdownTimer)
      countdownTimer = null
    }
  }
}

function formatCountdown(ms: number): string {
  const totalSeconds = Math.floor(ms / 1000)
  const hours = Math.floor(totalSeconds / 3600)
  const minutes = Math.floor((totalSeconds % 3600) / 60)
  const seconds = totalSeconds % 60
  return `${String(hours).padStart(2, '0')}:${String(minutes).padStart(2, '0')}:${String(seconds).padStart(2, '0')}`
}

async function fetchContest() {
  loading.value = true
  try {
    const res = await getContest(Number(props.id))
    contest.value = res.data.data
    updateCountdown()
  } catch (e: any) {
    message.error(e.message || t('errors.networkError'))
  } finally {
    loading.value = false
  }
}

async function fetchProblems() {
  try {
    const res = await getContestProblems(Number(props.id))
    problems.value = res.data.data || []
  } catch {
    // problems may not be available
  }
}

async function fetchRankings() {
  try {
    const res = await getContestRanking(Number(props.id))
    rankings.value = res.data.data?.items || []
  } catch {
    // ranking may not be available
  }
}

async function handleSignup() {
  signupLoading.value = true
  try {
    await signupContest(Number(props.id))
    isSignedUp.value = true
    message.success(t('contests.signupSuccess'))
  } catch (e: any) {
    message.error(e.message || t('errors.networkError'))
  } finally {
    signupLoading.value = false
  }
}

async function handleWithdraw() {
  withdrawLoading.value = true
  try {
    await withdrawContest(Number(props.id))
    isSignedUp.value = false
    message.success(t('contests.withdrawSuccess'))
  } catch (e: any) {
    message.error(e.message || t('errors.networkError'))
  } finally {
    withdrawLoading.value = false
  }
}

function selectProblem(index: number) {
  selectedProblem.value = problems.value[index]
  showProblemModal.value = true
}

function goToRanking() {
  router.push({ name: 'ContestRanking', params: { id: props.id } })
}

const miniRankingColumns = computed(() => {
  const cols: any[] = [
    {
      title: '#',
      key: 'rank',
      width: 50,
      render(row: ContestRanking) {
        if (row.rank === 1) return h('span', { class: 'rank-badge gold' }, '1')
        if (row.rank === 2) return h('span', { class: 'rank-badge silver' }, '2')
        if (row.rank === 3) return h('span', { class: 'rank-badge bronze' }, '3')
        return h('span', {}, String(row.rank))
      },
    },
    {
      title: t('contests.user'),
      key: 'user',
      ellipsis: { tooltip: true },
      render(row: ContestRanking) {
        return row.user?.username || `User #${row.userId}`
      },
    },
  ]

  if (contest.value?.ruleType === 'acm') {
    cols.push(
      { title: t('contests.solved'), key: 'solvedCount', width: 70 },
      { title: t('contests.penalty'), key: 'penalty', width: 80, render(row: ContestRanking) { return `${row.penalty}` } },
    )
  } else if (contest.value?.ruleType === 'oi') {
    cols.push({ title: t('contests.totalScore'), key: 'totalScore', width: 80, render(row: ContestRanking) { return row.solvedCount } })
  } else {
    cols.push({ title: t('contests.score'), key: 'solvedCount', width: 80 })
  }

  return cols
})

const miniRankings = computed(() => rankings.value.slice(0, 10))

onMounted(async () => {
  await fetchContest()
  if (contest.value) {
    const now = Date.now()
    const start = new Date(contest.value.startTime).getTime()
    const end = new Date(contest.value.endTime).getTime()
    if (now >= start) {
      await fetchProblems()
      await fetchRankings()
    }
    if (now < end) {
      countdownTimer = setInterval(updateCountdown, 1000)
    }
  }
})

onUnmounted(() => {
  if (countdownTimer) {
    clearInterval(countdownTimer)
    countdownTimer = null
  }
})
</script>

<template>
  <div class="contest-detail-page">
    <NSpin :show="loading">
      <template v-if="contest">
        <!-- Contest Header -->
        <div class="contest-header">
          <div class="header-left">
            <h1 class="contest-title">{{ contest.title }}</h1>
            <div class="header-badges">
              <NTag :type="statusColorMap[contest.status] || 'default'" size="small" bordered>
                {{ t(`contests.${contest.status}`) }}
              </NTag>
              <NTag :type="ruleTypeColorMap[contest.ruleType] || 'default'" size="small" bordered>
                {{ t(`contests.${contest.ruleType}`) }}
              </NTag>
              <NTag type="info" size="small" bordered>
                {{ t(`contests.${contest.type}`) }}
              </NTag>
            </div>
          </div>
          <div class="header-right">
            <div v-if="remainingTime && contest.status !== 'ended'" class="countdown-block">
              <span class="countdown-label">{{ countdownLabel }}</span>
              <span class="countdown-value">{{ remainingTime }}</span>
            </div>
            <NStatistic :label="t('contests.participantCount')" :value="contest.participantCount" class="participant-stat" />
          </div>
        </div>

        <!-- Contest Info Bar -->
        <div class="info-bar">
          <div class="info-item">
            <span class="info-label">{{ t('contests.startTime') }}</span>
            <span class="info-value">{{ formatTime(contest.startTime) }}</span>
          </div>
          <div class="info-item">
            <span class="info-label">{{ t('contests.endTime') }}</span>
            <span class="info-value">{{ formatTime(contest.endTime) }}</span>
          </div>
          <div class="info-item">
            <span class="info-label">{{ t('contests.duration') }}</span>
            <span class="info-value">{{ getDuration(contest.startTime, contest.endTime) }}</span>
          </div>
          <div class="info-item">
            <span class="info-label">{{ t('contests.ruleType') }}</span>
            <span class="info-value">{{ t(`contests.${contest.ruleType}`) }}</span>
          </div>
        </div>

        <!-- Not signed up / Contest ended view -->
        <div v-if="contest.status === 'ended' || !isSignedUp" class="contest-ended-view">
          <NCard :title="t('contests.contestDescription')" class="desc-card">
            <div class="contest-description" v-html="contest.description" />
          </NCard>

          <div v-if="contest.status === 'ended'" class="ended-actions">
            <NButton type="primary" @click="goToRanking">
              {{ t('contests.viewRanking') }}
            </NButton>
          </div>

          <div v-else-if="!isSignedUp && contest.status === 'upcoming'" class="signup-section">
            <NButton type="primary" size="large" :loading="signupLoading" @click="handleSignup">
              {{ t('contests.signup') }}
            </NButton>
          </div>
        </div>

        <!-- Signed up and running / upcoming -->
        <div v-else class="contest-active-view">
          <NTabs v-model:value="activeTab" type="line" animated>
            <NTabPane name="problems" :tab="t('contests.problems')">
              <div v-if="problems.length > 0" class="problem-panel">
                <div
                  v-for="(problem, index) in problems"
                  :key="problem.id"
                  class="problem-item"
                  :class="{
                    solved: problemStatusMap[problem.id] === 'solved',
                    attempted: problemStatusMap[problem.id] === 'attempted',
                  }"
                  @click="selectProblem(index)"
                >
                  <span class="problem-label">{{ problemLabels[index] }}</span>
                  <span class="problem-title-text">{{ problem.title }}</span>
                  <NTag
                    v-if="problemStatusMap[problem.id] === 'solved'"
                    type="success"
                    size="tiny"
                    bordered
                  >
                    {{ t('contests.solved') }}
                  </NTag>
                  <NTag
                    v-else-if="problemStatusMap[problem.id] === 'attempted'"
                    type="warning"
                    size="tiny"
                    bordered
                  >
                    {{ t('contests.attempted') }}
                  </NTag>
                </div>
              </div>
              <NEmpty v-else :description="t('contests.noProblems')" />
            </NTabPane>

            <NTabPane name="ranking" :tab="t('contests.miniRanking')">
              <div class="ranking-section">
                <NButton type="primary" size="small" class="full-ranking-btn" @click="goToRanking">
                  {{ t('contests.fullRanking') }}
                </NButton>
                <NDataTable
                  v-if="miniRankings.length > 0"
                  :columns="miniRankingColumns"
                  :data="miniRankings"
                  :bordered="false"
                  :single-line="false"
                  size="small"
                  :row-key="(row: ContestRanking) => row.userId"
                />
                <NEmpty v-else :description="t('contests.noRanking')" />
              </div>
            </NTabPane>
          </NTabs>
        </div>
      </template>
    </NSpin>

    <!-- Problem Detail Modal -->
    <NModal
      v-model:show="showProblemModal"
      preset="card"
      :title="selectedProblem?.title || ''"
      class="problem-modal"
      style="max-width: 900px; width: 90vw;"
    >
      <template v-if="selectedProblem">
        <div class="problem-detail-content">
          <div class="problem-meta">
            <NTag type="info" size="small" bordered>
              {{ t(`problems.${selectedProblem.type}`) }}
            </NTag>
            <NTag type="warning" size="small" bordered>
              {{ t(`problems.${selectedProblem.difficulty}`) }}
            </NTag>
            <span class="time-limit">{{ t('problems.timeLimit') }}: {{ selectedProblem.timeLimit }}ms</span>
            <span class="memory-limit">{{ t('problems.memoryLimit') }}: {{ selectedProblem.memoryLimit }}MB</span>
          </div>
          <div class="problem-description" v-html="selectedProblem.description" />
        </div>
      </template>
    </NModal>
  </div>
</template>

<style scoped lang="scss">
.contest-detail-page {
  padding: 24px 20px;
  max-width: 1400px;
  margin: 0 auto;
}

.contest-header {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  gap: 24px;
  margin-bottom: 20px;
  flex-wrap: wrap;
}

.header-left {
  flex: 1;
  min-width: 0;
}

.contest-title {
  font-size: 28px;
  font-weight: 700;
  color: var(--color-text);
  margin: 0 0 12px 0;
}

.header-badges {
  display: flex;
  gap: 8px;
  flex-wrap: wrap;
}

.header-right {
  display: flex;
  flex-direction: column;
  align-items: flex-end;
  gap: 12px;
  flex-shrink: 0;
}

.countdown-block {
  text-align: right;
}

.countdown-label {
  display: block;
  font-size: 12px;
  color: var(--color-text-secondary);
  margin-bottom: 4px;
}

.countdown-value {
  font-size: 28px;
  font-weight: 700;
  font-family: 'Courier New', monospace;
  color: var(--color-primary);
}

.participant-stat {
  :deep(.n-statistic-value) {
    color: var(--color-text);
    font-size: 24px;
  }
}

.info-bar {
  display: flex;
  gap: 32px;
  padding: 16px 20px;
  background-color: var(--color-bg-card);
  border: 1px solid rgba(255, 255, 255, 0.06);
  border-radius: var(--border-radius);
  margin-bottom: 24px;
  flex-wrap: wrap;
}

.info-item {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.info-label {
  font-size: 12px;
  color: var(--color-text-secondary);
}

.info-value {
  font-size: 14px;
  color: var(--color-text);
  font-weight: 500;
}

.desc-card {
  background-color: var(--color-bg-card);
  border: 1px solid rgba(255, 255, 255, 0.06);
  margin-bottom: 24px;

  :deep(.n-card__content) {
    max-height: 500px;
    overflow-y: auto;
  }
}

.contest-description {
  color: var(--color-text);
  line-height: 1.8;
  font-size: 14px;

  :deep(h1), :deep(h2), :deep(h3) {
    color: var(--color-text);
    margin-top: 16px;
    margin-bottom: 8px;
  }

  :deep(pre) {
    background: rgba(0, 0, 0, 0.3);
    padding: 12px;
    border-radius: 6px;
    overflow-x: auto;
  }

  :deep(code) {
    font-family: 'Courier New', monospace;
    font-size: 13px;
  }

  :deep(table) {
    border-collapse: collapse;
    width: 100%;
    margin: 12px 0;
  }

  :deep(th), :deep(td) {
    border: 1px solid rgba(255, 255, 255, 0.1);
    padding: 8px 12px;
    text-align: left;
  }

  :deep(th) {
    background: rgba(255, 255, 255, 0.03);
  }
}

.ended-actions,
.signup-section {
  display: flex;
  justify-content: center;
  padding: 32px 0;
}

.contest-active-view {
  :deep(.n-tabs) {
    .n-tabs-tab {
      color: var(--color-text-secondary);
    }

    .n-tabs-tab--active {
      color: var(--color-primary);
    }
  }
}

.problem-panel {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(280px, 1fr));
  gap: 12px;
  padding: 8px 0;
}

.problem-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 14px 16px;
  background-color: var(--color-bg-card);
  border: 1px solid rgba(255, 255, 255, 0.06);
  border-radius: var(--border-radius);
  cursor: pointer;
  transition: border-color 0.2s, background-color 0.2s;

  &:hover {
    border-color: var(--color-primary);
    background-color: rgba(0, 212, 255, 0.04);
  }

  &.solved {
    border-left: 3px solid var(--color-success);
  }

  &.attempted {
    border-left: 3px solid var(--color-warning);
  }
}

.problem-label {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 32px;
  height: 32px;
  border-radius: 6px;
  background: rgba(0, 212, 255, 0.12);
  color: var(--color-primary);
  font-weight: 700;
  font-size: 14px;
  flex-shrink: 0;
}

.problem-title-text {
  flex: 1;
  font-size: 14px;
  color: var(--color-text);
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.ranking-section {
  position: relative;
}

.full-ranking-btn {
  position: absolute;
  top: -40px;
  right: 0;
  z-index: 1;
}

.rank-badge {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 24px;
  height: 24px;
  border-radius: 50%;
  font-size: 12px;
  font-weight: 700;

  &.gold {
    background: linear-gradient(135deg, #ffd700, #ffb800);
    color: #1a1a2e;
  }

  &.silver {
    background: linear-gradient(135deg, #c0c0c0, #a0a0a0);
    color: #1a1a2e;
  }

  &.bronze {
    background: linear-gradient(135deg, #cd7f32, #b87333);
    color: #1a1a2e;
  }
}

.problem-modal {
  :deep(.n-card__content) {
    max-height: 70vh;
    overflow-y: auto;
  }
}

.problem-detail-content {
  .problem-meta {
    display: flex;
    align-items: center;
    gap: 12px;
    margin-bottom: 20px;
    flex-wrap: wrap;
  }

  .time-limit,
  .memory-limit {
    font-size: 13px;
    color: var(--color-text-secondary);
  }

  .problem-description {
    color: var(--color-text);
    line-height: 1.8;
    font-size: 14px;

    :deep(h1), :deep(h2), :deep(h3) {
      color: var(--color-text);
      margin-top: 16px;
      margin-bottom: 8px;
    }

    :deep(pre) {
      background: rgba(0, 0, 0, 0.3);
      padding: 12px;
      border-radius: 6px;
      overflow-x: auto;
    }

    :deep(code) {
      font-family: 'Courier New', monospace;
      font-size: 13px;
    }

    :deep(table) {
      border-collapse: collapse;
      width: 100%;
      margin: 12px 0;
    }

    :deep(th), :deep(td) {
      border: 1px solid rgba(255, 255, 255, 0.1);
      padding: 8px 12px;
      text-align: left;
    }

    :deep(th) {
      background: rgba(255, 255, 255, 0.03);
    }
  }
}

@media (max-width: 768px) {
  .contest-header {
    flex-direction: column;
  }

  .header-right {
    align-items: flex-start;
  }

  .info-bar {
    flex-direction: column;
    gap: 16px;
  }

  .problem-panel {
    grid-template-columns: 1fr;
  }
}
</style>
