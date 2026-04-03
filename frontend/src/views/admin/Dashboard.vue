<script setup lang="ts">
import { ref, onMounted, h } from 'vue'
import { useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import {
  NCard, NGrid, NGridItem, NStatistic, NDataTable, NButton, NSpace, NTag, useMessage,
} from 'naive-ui'
import { adminGetUsers, adminGetProblems, adminGetSubmissions, adminGetContests, getJudgeStatus } from '@/api/admin'
import type { Submission, Contest } from '@/types'

const { t } = useI18n()
const router = useRouter()
const message = useMessage()

const loading = ref(false)
const stats = ref({
  totalUsers: 0,
  totalProblems: 0,
  totalSubmissions: 0,
  totalContests: 0,
  activeUsersToday: 0,
})

const recentSubmissions = ref<Submission[]>([])
const recentContests = ref<Contest[]>([])
const judgeConnected = ref(false)
const dbConnected = ref(false)

type NTagType = 'default' | 'error' | 'info' | 'primary' | 'success' | 'warning'

const submissionStatusColor: Record<string, NTagType> = {
  pending: 'warning',
  judging: 'info',
  accepted: 'success',
  wrong_answer: 'error',
  time_limit_exceeded: 'warning',
  memory_limit_exceeded: 'warning',
  runtime_error: 'error',
  compilation_error: 'error',
  presentation_error: 'warning',
  system_error: 'error',
}

const contestStatusColor: Record<string, NTagType> = {
  upcoming: 'info',
  running: 'success',
  ended: 'default',
}

const submissionColumns = [
  {
    title: 'ID',
    key: 'id',
    width: 80,
  },
  {
    title: t('submissions.problem'),
    key: 'problemId',
    width: 100,
    render(row: Submission) {
      return h('span', {}, `#${row.problemId}`)
    },
  },
  {
    title: t('submissions.language'),
    key: 'language',
    width: 100,
  },
  {
    title: t('submissions.status'),
    key: 'status',
    width: 140,
    render(row: Submission) {
      const color = submissionStatusColor[row.status] || 'default'
      return h(
        NTag,
        { type: color, size: 'small', bordered: false },
        { default: () => formatStatus(row.status) }
      )
    },
  },
  {
    title: t('submissions.timeUsed'),
    key: 'timeUsed',
    width: 100,
    render(row: Submission) {
      return row.timeUsed != null ? `${row.timeUsed}ms` : '-'
    },
  },
  {
    title: t('submissions.submitTime'),
    key: 'createdAt',
    width: 180,
    render(row: Submission) {
      return formatTime(row.createdAt)
    },
  },
]

const contestColumns = [
  {
    title: 'ID',
    key: 'id',
    width: 80,
  },
  {
    title: t('contests.title'),
    key: 'title',
    ellipsis: { tooltip: true },
  },
  {
    title: t('contests.status'),
    key: 'status',
    width: 100,
    render(row: Contest) {
      const color = contestStatusColor[row.status] || 'default'
      return h(
        NTag,
        { type: color, size: 'small', bordered: false },
        { default: () => t(`contests.${row.status}`) }
      )
    },
  },
  {
    title: t('contests.participants'),
    key: 'participantCount',
    width: 100,
  },
  {
    title: t('contests.startTime'),
    key: 'startTime',
    width: 180,
    render(row: Contest) {
      return formatTime(row.startTime)
    },
  },
]

function formatStatus(status: string): string {
  const key = status.replace(/_/g, '')
  const camelKey = key.charAt(0).toLowerCase() + key.slice(1)
  const mapped: Record<string, string> = {
    pending: 'submissions.pending',
    judging: 'submissions.judging',
    accepted: 'submissions.accepted',
    wrongAnswer: 'submissions.wrongAnswer',
    timeLimitExceeded: 'submissions.timeLimitExceeded',
    memoryLimitExceeded: 'submissions.memoryLimitExceeded',
    runtimeError: 'submissions.runtimeError',
    compilationError: 'submissions.compilationError',
    presentationError: 'submissions.presentationError',
    systemError: 'submissions.systemError',
  }
  return t(mapped[camelKey] || `submissions.${status}`) || status
}

function formatTime(time: string): string {
  if (!time) return '-'
  return new Date(time).toLocaleString()
}

async function fetchDashboard() {
  loading.value = true
  try {
    const [usersRes, problemsRes, submissionsRes, contestsRes] = await Promise.allSettled([
      adminGetUsers({ page: 1, pageSize: 1 }),
      adminGetProblems({ page: 1, pageSize: 1 }),
      adminGetSubmissions({ page: 1, pageSize: 10 }),
      adminGetContests({ page: 1, pageSize: 5 }),
    ])

    if (usersRes.status === 'fulfilled') {
      stats.value.totalUsers = usersRes.value.data.data?.total || 0
    }
    if (problemsRes.status === 'fulfilled') {
      stats.value.totalProblems = problemsRes.value.data.data?.total || 0
    }
    if (submissionsRes.status === 'fulfilled') {
      const subData = submissionsRes.value.data.data
      stats.value.totalSubmissions = subData?.total || 0
      recentSubmissions.value = subData?.items || []
    }
    if (contestsRes.status === 'fulfilled') {
      const contestData = contestsRes.value.data.data
      stats.value.totalContests = contestData?.total || 0
      recentContests.value = contestData?.items || []
    }

    try {
      const judgeRes = await getJudgeStatus()
      const judgeData = judgeRes.data.data as Record<string, any> | null
      if (judgeData) {
        judgeConnected.value = !!judgeData.connected
        dbConnected.value = !!judgeData.dbConnected
      }
    } catch {
      // Judge service may not be available
    }
  } catch (e: any) {
    message.error(e.message || t('errors.networkError'))
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  fetchDashboard()
})
</script>

<template>
  <div class="admin-dashboard">
    <div class="page-header">
      <h1 class="page-title">{{ t('admin.dashboard') }}</h1>
    </div>

    <!-- Statistics Cards -->
    <NGrid :x-gap="16" :y-gap="16" :cols="5" responsive="screen" item-responsive>
      <NGridItem span="5 m:1">
        <NCard class="stat-card">
          <NStatistic :label="t('admin.totalUsers')" :value="stats.totalUsers" />
        </NCard>
      </NGridItem>
      <NGridItem span="5 m:1">
        <NCard class="stat-card">
          <NStatistic :label="t('admin.totalProblems')" :value="stats.totalProblems" />
        </NCard>
      </NGridItem>
      <NGridItem span="5 m:1">
        <NCard class="stat-card">
          <NStatistic :label="t('admin.totalSubmissions')" :value="stats.totalSubmissions" />
        </NCard>
      </NGridItem>
      <NGridItem span="5 m:1">
        <NCard class="stat-card">
          <NStatistic :label="t('admin.totalContests')" :value="stats.totalContests" />
        </NCard>
      </NGridItem>
      <NGridItem span="5 m:1">
        <NCard class="stat-card">
          <NStatistic :label="t('admin.activeUsersToday')" :value="stats.activeUsersToday" />
        </NCard>
      </NGridItem>
    </NGrid>

    <!-- Quick Actions -->
    <NCard class="section-card" :title="t('admin.quickActions')">
      <NSpace>
        <NButton type="primary" @click="router.push('/admin/problems/create')">
          {{ t('admin.createProblem') }}
        </NButton>
        <NButton type="primary" @click="router.push('/admin/contests/create')">
          {{ t('admin.createContest') }}
        </NButton>
        <NButton type="info" @click="router.push('/admin/announcements')">
          {{ t('admin.publishAnnouncement') }}
        </NButton>
      </NSpace>
    </NCard>

    <!-- System Status -->
    <NCard class="section-card" :title="t('admin.systemStatus')">
      <NSpace>
        <NTag :type="judgeConnected ? 'success' : 'error'" size="large" bordered>
          {{ t('admin.judgeService') }}: {{ judgeConnected ? t('admin.connected') : t('admin.disconnected') }}
        </NTag>
        <NTag :type="dbConnected ? 'success' : 'error'" size="large" bordered>
          {{ t('admin.database') }}: {{ dbConnected ? t('admin.connected') : t('admin.disconnected') }}
        </NTag>
      </NSpace>
    </NCard>

    <!-- Recent Submissions -->
    <NCard class="section-card" :title="t('admin.recentSubmissions')">
      <NDataTable
        :columns="submissionColumns"
        :data="recentSubmissions"
        :loading="loading"
        :bordered="false"
        :single-line="false"
        size="small"
        :row-key="(row: Submission) => row.id"
      />
    </NCard>

    <!-- Recent Contests -->
    <NCard class="section-card" :title="t('admin.recentContests')">
      <NDataTable
        :columns="contestColumns"
        :data="recentContests"
        :loading="loading"
        :bordered="false"
        :single-line="false"
        size="small"
        :row-key="(row: Contest) => row.id"
      />
    </NCard>
  </div>
</template>

<style scoped lang="scss">
.admin-dashboard {
  max-width: 1400px;
  margin: 0 auto;
}

.page-header {
  margin-bottom: 24px;
}

.page-title {
  font-size: 28px;
  font-weight: 700;
  color: var(--color-text);
}

.stat-card {
  background-color: var(--color-bg-card);
  border: 1px solid rgba(255, 255, 255, 0.06);

  :deep(.n-statistic) {
    --n-label-font-size: 13px;
    --n-value-font-size: 28px;
    --n-label-text-color: var(--color-text-secondary);
    --n-value-text-color: var(--color-primary);
  }
}

.section-card {
  margin-top: 16px;
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
}
</style>
