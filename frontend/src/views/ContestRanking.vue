<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { NDataTable, NTag, NCard, NSpace, NButton, NEmpty, NSpin, useMessage } from 'naive-ui'
import { getContest, getContestRanking, getContestFrozenRanking } from '@/api/contest'
import { useAuthStore } from '@/stores/auth'
import type { Contest, ContestRanking, Problem } from '@/types'

const props = defineProps<{ id: string }>()
const router = useRouter()
const { t } = useI18n()
const message = useMessage()
const authStore = useAuthStore()

const loading = ref(true)
const contest = ref<Contest | null>(null)
const rankings = ref<ContestRanking[]>([])
const frozenRankings = ref<ContestRanking[]>([])
const isFrozen = ref(false)
const autoRefresh = ref(true)
let refreshTimer: ReturnType<typeof setInterval> | null = null

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
  if (!contest.value?.problems?.length) return []
  const labels = 'ABCDEFGHIJKLMNOPQRSTUVWXYZ'
  return contest.value.problems.map((_, i) => labels[i] || `P${i + 1}`)
})

const currentUserId = computed(() => authStore.user?.id)

function formatTime(dateStr: string): string {
  return new Date(dateStr).toLocaleString()
}

function formatDuration(minutes: number): string {
  if (minutes === 0) return '-'
  const h = Math.floor(minutes / 60)
  const m = minutes % 60
  if (h > 0) return `${h}:${String(m).padStart(2, '0')}`
  return `${m}m`
}

function buildColumns() {
  if (!contest.value) return []

  const cols: any[] = [
    {
      title: '#',
      key: 'rank',
      width: 60,
      sorter: false,
      render(row: ContestRanking) {
        if (row.rank === 1) return h('span', { class: 'rank-badge gold' }, '1')
        if (row.rank === 2) return h('span', { class: 'rank-badge silver' }, '2')
        if (row.rank === 3) return h('span', { class: 'rank-badge bronze' }, '3')
        return h('span', { class: 'rank-number' }, String(row.rank))
      },
    },
    {
      title: t('contests.user'),
      key: 'user',
      width: 150,
      ellipsis: { tooltip: true },
      render(row: ContestRanking) {
        return row.user?.username || `User #${row.userId}`
      },
    },
  ]

  const ruleType = contest.value.ruleType

  if (ruleType === 'acm') {
    cols.push(
      {
        title: t('contests.solved'),
        key: 'solvedCount',
        width: 80,
        sorter: (a: ContestRanking, b: ContestRanking) => a.solvedCount - b.solvedCount,
      },
      {
        title: t('contests.penalty'),
        key: 'penalty',
        width: 90,
        sorter: (a: ContestRanking, b: ContestRanking) => a.penalty - b.penalty,
        render(row: ContestRanking) {
          return formatDuration(row.penalty)
        },
      },
    )
  } else if (ruleType === 'oi') {
    cols.push({
      title: t('contests.totalScore'),
      key: 'totalScore',
      width: 100,
      sorter: (a: ContestRanking, b: ContestRanking) => a.solvedCount - b.solvedCount,
      render(row: ContestRanking) {
        return row.solvedCount
      },
    })
  } else if (ruleType === 'cf') {
    cols.push({
      title: t('contests.score'),
      key: 'score',
      width: 90,
      sorter: (a: ContestRanking, b: ContestRanking) => a.solvedCount - b.solvedCount,
      render(row: ContestRanking) {
        return row.solvedCount
      },
    })
  } else {
    // CTF
    cols.push(
      {
        title: t('contests.solved'),
        key: 'solvedCount',
        width: 80,
      },
      {
        title: t('contests.lastSolveTime'),
        key: 'lastSolve',
        width: 140,
        render(row: ContestRanking) {
          const results = Object.values(row.problemResults)
          const solved = results.filter(r => r.status === 'accepted')
          if (solved.length === 0) return '-'
          const times = solved.map(r => r.acceptedAt).filter(Boolean).sort()
          return times.length > 0 ? formatTime(times[times.length - 1]!) : '-'
        },
      },
    )
  }

  // Per-problem columns
  const problemIds = contest.value.problemIds || []
  problemIds.forEach((pid, index) => {
    const label = problemLabels.value[index] || `P${index + 1}`
    cols.push({
      title: label,
      key: `problem_${pid}`,
      width: 90,
      align: 'center' as const,
      render(row: ContestRanking) {
        const result = row.problemResults[pid]
        if (!result) return h('span', { class: 'problem-cell empty' }, '-')

        if (result.status === 'accepted') {
          if (ruleType === 'acm') {
            const attempts = result.submitCount - 1
            const timeMin = Math.floor(result.timeUsed / 60)
            return h('span', { class: 'problem-cell accepted' }, `${attempts > 0 ? attempts + '/' : ''}${timeMin}`)
          } else if (ruleType === 'oi') {
            return h('span', { class: 'problem-cell scored' }, `${result.submitCount}`)
          } else {
            return h('span', { class: 'problem-cell accepted' }, '+')
          }
        } else if (result.submitCount > 0) {
          if (ruleType === 'acm') {
            return h('span', { class: 'problem-cell wrong' }, `-${result.submitCount}`)
          } else if (ruleType === 'oi') {
            return h('span', { class: 'problem-cell partial' }, `${result.submitCount}`)
          } else {
            return h('span', { class: 'problem-cell wrong' }, `${result.submitCount}`)
          }
        }
        return h('span', { class: 'problem-cell empty' }, '-')
      },
    })
  })

  return cols
}

const columns = computed(() => buildColumns())

const rowProps = (row: ContestRanking) => {
  const isCurrentUser = row.userId === currentUserId.value
  return {
    style: isCurrentUser
      ? 'background-color: rgba(0, 212, 255, 0.08); font-weight: 600;'
      : '',
  }
}

async function fetchContest() {
  try {
    const res = await getContest(Number(props.id))
    contest.value = res.data.data
  } catch (e: any) {
    message.error(e.message || t('errors.networkError'))
  }
}

async function fetchRankings() {
  loading.value = true
  try {
    const res = await getContestRanking(Number(props.id))
    rankings.value = res.data.data?.list || []
  } catch {
    // ignore
  } finally {
    loading.value = false
  }
}

async function fetchFrozenRankings() {
  try {
    const res = await getContestFrozenRanking(Number(props.id))
    frozenRankings.value = res.data.data || []
    isFrozen.value = frozenRankings.value.length > 0
  } catch {
    isFrozen.value = false
  }
}

function toggleAutoRefresh() {
  autoRefresh.value = !autoRefresh.value
  if (autoRefresh.value) {
    startAutoRefresh()
  } else {
    stopAutoRefresh()
  }
}

function startAutoRefresh() {
  stopAutoRefresh()
  refreshTimer = setInterval(() => {
    fetchRankings()
  }, 30000)
}

function stopAutoRefresh() {
  if (refreshTimer) {
    clearInterval(refreshTimer)
    refreshTimer = null
  }
}

function goToContest() {
  router.push({ name: 'ContestDetail', params: { id: props.id } })
}

onMounted(async () => {
  await fetchContest()
  await fetchRankings()
  if (contest.value?.status === 'running') {
    await fetchFrozenRankings()
    startAutoRefresh()
  }
})

onUnmounted(() => {
  stopAutoRefresh()
})
</script>

<template>
  <div class="contest-ranking-page">
    <NSpin :show="loading">
      <template v-if="contest">
        <!-- Contest Info Header -->
        <div class="ranking-header">
          <div class="header-left">
            <NButton text @click="goToContest" class="back-btn">
              {{ t('common.back') }}
            </NButton>
            <h1 class="page-title">{{ contest.title }} - {{ t('contests.ranking') }}</h1>
            <div class="header-badges">
              <NTag :type="statusColorMap[contest.status] || 'default'" size="small" bordered>
                {{ t(`contests.${contest.status}`) }}
              </NTag>
              <NTag :type="ruleTypeColorMap[contest.ruleType] || 'default'" size="small" bordered>
                {{ t(`contests.${contest.ruleType}`) }}
              </NTag>
            </div>
          </div>
          <div class="header-right">
            <NTag v-if="isFrozen" type="warning" size="small" bordered>
              {{ t('contests.boardFrozen') }}
            </NTag>
            <NButton
              v-if="contest.status === 'running'"
              size="small"
              :type="autoRefresh ? 'primary' : 'default'"
              @click="toggleAutoRefresh"
            >
              {{ t('contests.autoRefresh') }}: {{ autoRefresh ? 'ON' : 'OFF' }}
            </NButton>
          </div>
        </div>

        <!-- Ranking Table -->
        <NCard class="ranking-card" v-if="rankings.length > 0">
          <NDataTable
            :columns="columns"
            :data="rankings"
            :bordered="false"
            :single-line="false"
            size="small"
            :row-props="rowProps"
            :row-key="(row: ContestRanking) => row.userId"
            :max-height="700"
            virtual-scroll
            :scroll-x="800"
          />
        </NCard>

        <NEmpty v-else :description="t('contests.noRanking')" class="empty-state" />
      </template>
    </NSpin>
  </div>
</template>

<style scoped lang="scss">
.contest-ranking-page {
  padding: 24px 20px;
  max-width: 1600px;
  margin: 0 auto;
}

.ranking-header {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  gap: 24px;
  margin-bottom: 24px;
  flex-wrap: wrap;
}

.header-left {
  flex: 1;
  min-width: 0;
}

.back-btn {
  font-size: 13px;
  color: var(--color-text-secondary);
  margin-bottom: 8px;
  padding: 0;

  &:hover {
    color: var(--color-primary);
  }
}

.page-title {
  font-size: 24px;
  font-weight: 700;
  color: var(--color-text);
  margin: 0 0 12px 0;
}

.header-badges {
  display: flex;
  gap: 8px;
}

.header-right {
  display: flex;
  align-items: center;
  gap: 12px;
  flex-shrink: 0;
}

.ranking-card {
  background-color: var(--color-bg-card);
  border: 1px solid rgba(255, 255, 255, 0.06);

  :deep(.n-data-table) {
    --n-td-color: transparent;
    --n-th-color: rgba(255, 255, 255, 0.03);
    --n-border-color: rgba(255, 255, 255, 0.06);
    --n-td-text-color: var(--color-text);
    --n-th-text-color: var(--color-text-secondary);
  }

  :deep(.n-data-table-tr:hover .n-data-table-td) {
    background-color: rgba(0, 212, 255, 0.04);
  }

  :deep(.n-data-table-th) {
    text-align: center;
    font-weight: 600;
    font-size: 13px;
  }

  :deep(.n-data-table-td) {
    text-align: center;
  }
}

.rank-badge {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 26px;
  height: 26px;
  border-radius: 50%;
  font-size: 13px;
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

.rank-number {
  font-weight: 600;
  color: var(--color-text-secondary);
}

.problem-cell {
  display: inline-block;
  font-size: 12px;
  font-weight: 600;
  padding: 2px 6px;
  border-radius: 4px;
  min-width: 32px;

  &.accepted {
    color: var(--color-success);
    background: rgba(0, 230, 118, 0.1);
  }

  &.scored {
    color: var(--color-primary);
    background: rgba(0, 212, 255, 0.1);
  }

  &.partial {
    color: var(--color-warning);
    background: rgba(255, 171, 0, 0.1);
  }

  &.wrong {
    color: var(--color-error);
    background: rgba(255, 82, 82, 0.1);
  }

  &.empty {
    color: var(--color-text-secondary);
    opacity: 0.5;
  }
}

.empty-state {
  margin-top: 60px;
}

@media (max-width: 768px) {
  .ranking-header {
    flex-direction: column;
  }

  .header-right {
    align-items: flex-start;
  }
}
</style>
