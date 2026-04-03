<script setup lang="ts">
import { ref, onMounted, onUnmounted, h } from 'vue'
import { useI18n } from 'vue-i18n'
import {
  NCard, NGrid, NGridItem, NStatistic, NDataTable, NButton, NSpace, NTag, useMessage,
} from 'naive-ui'
import { getJudgeStatus } from '@/api/admin'
import type { Submission } from '@/types'

const { t } = useI18n()
const message = useMessage()

const loading = ref(false)
const connected = ref(false)
const queueSize = ref(0)
const runningCount = ref(0)
const poolSize = ref(0)
const recentResults = ref<any[]>([])
let refreshTimer: ReturnType<typeof setInterval> | null = null

type NTagType = 'default' | 'error' | 'info' | 'primary' | 'success' | 'warning'

const statusColorMap: Record<string, NTagType> = {
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

function formatStatus(status: string): string {
  const mapped: Record<string, string> = {
    pending: 'submissions.pending',
    judging: 'submissions.judging',
    accepted: 'submissions.accepted',
    wrong_answer: 'submissions.wrongAnswer',
    time_limit_exceeded: 'submissions.timeLimitExceeded',
    memory_limit_exceeded: 'submissions.memoryLimitExceeded',
    runtime_error: 'submissions.runtimeError',
    compilation_error: 'submissions.compilationError',
    presentation_error: 'submissions.presentationError',
    system_error: 'submissions.systemError',
  }
  return t(mapped[status] || status) || status
}

const resultColumns = [
  {
    title: 'ID',
    key: 'id',
    width: 80,
  },
  {
    title: t('submissions.problem'),
    key: 'problemId',
    width: 80,
    render(row: any) {
      return h('span', {}, `#${row.problemId}`)
    },
  },
  {
    title: t('auth.username'),
    key: 'userId',
    width: 120,
    render(row: any) {
      return row.username || `User#${row.userId}`
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
    render(row: any) {
      return h(
        NTag,
        { type: statusColorMap[row.status] || 'default', size: 'small', bordered: false },
        { default: () => formatStatus(row.status) }
      )
    },
  },
  {
    title: t('submissions.timeUsed'),
    key: 'timeUsed',
    width: 100,
    render(row: any) {
      return row.timeUsed != null ? `${row.timeUsed}ms` : '-'
    },
  },
  {
    title: t('submissions.memoryUsed'),
    key: 'memoryUsed',
    width: 100,
    render(row: any) {
      return row.memoryUsed != null ? `${row.memoryUsed}KB` : '-'
    },
  },
  {
    title: t('submissions.submitTime'),
    key: 'createdAt',
    width: 170,
    render(row: any) {
      return row.createdAt ? new Date(row.createdAt).toLocaleString() : '-'
    },
  },
]

async function fetchStatus() {
  loading.value = true
  try {
    const res = await getJudgeStatus()
    const data = res.data.data as Record<string, any> | null
    if (data) {
      connected.value = !!data.connected
      queueSize.value = data.queueSize || 0
      runningCount.value = data.runningCount || 0
      poolSize.value = data.poolSize || 0
      recentResults.value = data.recentResults || []
    }
  } catch {
    connected.value = false
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  fetchStatus()
  refreshTimer = setInterval(fetchStatus, 5000)
})

onUnmounted(() => {
  if (refreshTimer) {
    clearInterval(refreshTimer)
    refreshTimer = null
  }
})
</script>

<template>
  <div class="admin-judge-monitor">
    <div class="page-header flex-between">
      <h1 class="page-title">{{ t('admin.judgeMonitor') }}</h1>
      <NSpace>
        <NTag :type="connected ? 'success' : 'error'" size="large" bordered>
          {{ connected ? t('admin.connected') : t('admin.disconnected') }}
        </NTag>
        <NButton size="small" @click="fetchStatus">{{ t('common.refresh') }}</NButton>
      </NSpace>
    </div>

    <!-- Status Cards -->
    <NGrid :x-gap="16" :y-gap="16" :cols="3" responsive="screen" item-responsive>
      <NGridItem span="3 m:1">
        <NCard class="stat-card">
          <NStatistic :label="t('admin.queueSize')" :value="queueSize" />
        </NCard>
      </NGridItem>
      <NGridItem span="3 m:1">
        <NCard class="stat-card">
          <NStatistic :label="t('admin.runningCount')" :value="runningCount" />
        </NCard>
      </NGridItem>
      <NGridItem span="3 m:1">
        <NCard class="stat-card">
          <NStatistic :label="t('admin.poolSize')" :value="poolSize" />
        </NCard>
      </NGridItem>
    </NGrid>

    <!-- Auto-refresh hint -->
    <div class="auto-refresh-hint">
      {{ t('admin.autoRefreshHint') }}
    </div>

    <!-- Recent Judge Results -->
    <NCard class="section-card" :title="t('admin.recentJudgeResults')">
      <NDataTable
        :columns="resultColumns"
        :data="recentResults"
        :loading="loading"
        :bordered="false"
        :single-line="false"
        size="small"
        :row-key="(row: any) => row.id"
        :max-height="500"
      />
    </NCard>
  </div>
</template>

<style scoped lang="scss">
.admin-judge-monitor {
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

.auto-refresh-hint {
  margin-top: 12px;
  font-size: 13px;
  color: var(--color-text-secondary);
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
