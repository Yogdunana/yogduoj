<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { NDataTable, NInput, NSelect, NTag, NSpace, NCard, NDatePicker, NPagination, useMessage } from 'naive-ui'
import { listSubmissions } from '@/api/submission'
import type { Submission, SubmissionStatus } from '@/types'

const { t } = useI18n()
const router = useRouter()
const message = useMessage()

const loading = ref(false)
const submissions = ref<Submission[]>([])
const total = ref(0)

const filters = reactive({
  problemId: null as number | null,
  status: '' as string,
  language: '' as string,
  startTime: null as number | null,
  endTime: null as number | null,
  page: 1,
  pageSize: 20,
})

const statusOptions = computed(() => [
  { label: t('submissions.allStatuses'), value: '' },
  { label: t('submissions.pending'), value: 'pending' },
  { label: t('submissions.judging'), value: 'judging' },
  { label: t('submissions.accepted'), value: 'accepted' },
  { label: t('submissions.wrongAnswer'), value: 'wrong_answer' },
  { label: t('submissions.timeLimitExceeded'), value: 'time_limit_exceeded' },
  { label: t('submissions.memoryLimitExceeded'), value: 'memory_limit_exceeded' },
  { label: t('submissions.runtimeError'), value: 'runtime_error' },
  { label: t('submissions.compilationError'), value: 'compilation_error' },
  { label: t('submissions.presentationError'), value: 'presentation_error' },
  { label: t('submissions.systemError'), value: 'system_error' },
])

const languageOptions = computed(() => [
  { label: t('submissions.allLanguages'), value: '' },
  { label: 'C', value: 'c' },
  { label: 'C++', value: 'cpp' },
  { label: 'Java', value: 'java' },
  { label: 'Python3', value: 'python3' },
])

type NTagType = 'default' | 'error' | 'info' | 'primary' | 'success' | 'warning'

const statusColorMap: Record<string, NTagType> = {
  pending: 'default',
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

function getStatusLabel(status: SubmissionStatus): string {
  const map: Record<string, string> = {
    pending: t('submissions.pending'),
    judging: t('submissions.judging'),
    accepted: t('submissions.accepted'),
    wrong_answer: t('submissions.wrongAnswer'),
    time_limit_exceeded: t('submissions.timeLimitExceeded'),
    memory_limit_exceeded: t('submissions.memoryLimitExceeded'),
    runtime_error: t('submissions.runtimeError'),
    compilation_error: t('submissions.compilationError'),
    presentation_error: t('submissions.presentationError'),
    system_error: t('submissions.systemError'),
  }
  return map[status] || status
}

function formatTime(ms: number): string {
  return `${ms} ms`
}

function formatMemory(kb: number): string {
  if (kb >= 1024) return `${(kb / 1024).toFixed(1)} MB`
  return `${kb} KB`
}

function formatDateTime(dateStr: string): string {
  return new Date(dateStr).toLocaleString()
}

const columns = computed(() => [
  {
    title: t('submissions.submissionId'),
    key: 'id',
    width: 90,
    render(row: Submission) {
      return h('span', { class: 'submission-id' }, `#${row.id}`)
    },
  },
  {
    title: t('submissions.problem'),
    key: 'problem',
    ellipsis: { tooltip: true },
    render(row: Submission) {
      const title = row.problem?.title || `Problem #${row.problemId}`
      return h(
        'a',
        {
          class: 'problem-link',
          onClick: () => router.push({ name: 'ProblemDetail', params: { id: row.problemId } }),
        },
        title
      )
    },
  },
  {
    title: t('submissions.language'),
    key: 'language',
    width: 100,
    render(row: Submission) {
      return h('span', { class: 'lang-badge' }, row.language)
    },
  },
  {
    title: t('submissions.status'),
    key: 'status',
    width: 160,
    render(row: Submission) {
      return h(
        NTag,
        { type: statusColorMap[row.status] || 'default', size: 'small', bordered: false },
        { default: () => getStatusLabel(row.status) }
      )
    },
  },
  {
    title: t('submissions.timeUsed'),
    key: 'timeUsed',
    width: 100,
    render(row: Submission) {
      return row.timeUsed != null ? formatTime(row.timeUsed) : '-'
    },
  },
  {
    title: t('submissions.memoryUsed'),
    key: 'memoryUsed',
    width: 110,
    render(row: Submission) {
      return row.memoryUsed != null ? formatMemory(row.memoryUsed) : '-'
    },
  },
  {
    title: t('submissions.submitTime'),
    key: 'createdAt',
    width: 170,
    render(row: Submission) {
      return h('span', { class: 'time-text' }, formatDateTime(row.createdAt))
    },
  },
])

async function fetchSubmissions() {
  loading.value = true
  try {
    const params: Record<string, any> = {
      page: filters.page,
      pageSize: filters.pageSize,
    }
    if (filters.problemId) params.problemId = filters.problemId
    if (filters.status) params.status = filters.status
    if (filters.language) params.language = filters.language
    if (filters.startTime) params.startTime = new Date(filters.startTime).toISOString()
    if (filters.endTime) params.endTime = new Date(filters.endTime).toISOString()

    const res = await listSubmissions(params)
    const data = res.data.data
    submissions.value = data.items
    total.value = data.total
  } catch (e: any) {
    message.error(e.message || t('errors.networkError'))
  } finally {
    loading.value = false
  }
}

function handlePageChange(page: number) {
  filters.page = page
  fetchSubmissions()
}

function handleProblemIdInput(val: string) {
  const num = parseInt(val)
  filters.problemId = isNaN(num) ? null : num
}

function handleProblemIdSearch() {
  filters.page = 1
  fetchSubmissions()
}

function handleStatusChange(val: string) {
  filters.status = val
  filters.page = 1
  fetchSubmissions()
}

function handleLanguageChange(val: string) {
  filters.language = val
  filters.page = 1
  fetchSubmissions()
}

function handleTimeRangeChange(val: number[] | null) {
  if (val && val.length === 2) {
    filters.startTime = val[0]
    filters.endTime = val[1]
  } else {
    filters.startTime = null
    filters.endTime = null
  }
  filters.page = 1
  fetchSubmissions()
}

function handleRowClick(row: Submission) {
  router.push({ name: 'SubmissionDetail', params: { id: row.id } })
}

function resetFilters() {
  filters.problemId = null
  filters.status = ''
  filters.language = ''
  filters.startTime = null
  filters.endTime = null
  filters.page = 1
  fetchSubmissions()
}

onMounted(() => {
  fetchSubmissions()
})
</script>

<template>
  <div class="submission-list-page">
    <div class="page-header">
      <h1 class="page-title">{{ t('submissions.title') }}</h1>
    </div>

    <!-- Filter Bar -->
    <NCard size="small" class="filter-card">
      <div class="filter-bar">
        <div class="filter-item">
          <label class="filter-label">{{ t('submissions.filterByProblem') }}</label>
          <NInput
            :placeholder="'#'"
            clearable
            size="small"
            style="width: 120px"
            @update:value="handleProblemIdInput"
            @keydown.enter="handleProblemIdSearch"
          />
        </div>

        <div class="filter-item">
          <label class="filter-label">{{ t('submissions.filterByResult') }}</label>
          <NSelect
            :value="filters.status"
            :options="statusOptions"
            size="small"
            style="width: 180px"
            @update:value="handleStatusChange"
          />
        </div>

        <div class="filter-item">
          <label class="filter-label">{{ t('submissions.filterByLanguage') }}</label>
          <NSelect
            :value="filters.language"
            :options="languageOptions"
            size="small"
            style="width: 120px"
            @update:value="handleLanguageChange"
          />
        </div>

        <div class="filter-item">
          <label class="filter-label">{{ t('submissions.filterByTime') }}</label>
          <NDatePicker
            :value="filters.startTime && filters.endTime ? [filters.startTime, filters.endTime] : null"
            type="datetimerange"
            size="small"
            clearable
            style="width: 340px"
            @update:value="handleTimeRangeChange"
          />
        </div>

        <div class="filter-item filter-actions">
          <NButton size="small" @click="resetFilters">
            {{ t('common.reset') }}
          </NButton>
        </div>
      </div>
    </NCard>

    <!-- Submission Table -->
    <NCard class="table-card">
      <NDataTable
        :columns="columns"
        :data="submissions"
        :loading="loading"
        :row-props="(row: Submission) => ({
          style: 'cursor: pointer',
          onClick: () => handleRowClick(row),
        })"
        :bordered="false"
        :single-line="false"
        size="small"
        :row-key="(row: Submission) => row.id"
        :max-height="600"
        virtual-scroll
      />
    </NCard>

    <!-- Pagination -->
    <div class="pagination-wrapper">
      <NPagination
        :page="filters.page"
        :page-size="filters.pageSize"
        :item-count="total"
        @update:page="handlePageChange"
      />
    </div>
  </div>
</template>

<style scoped lang="scss">
.submission-list-page {
  padding: 24px 20px;
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

.filter-card {
  background-color: var(--color-bg-card);
  border: 1px solid rgba(255, 255, 255, 0.06);
  margin-bottom: 16px;
}

.filter-bar {
  display: flex;
  align-items: flex-end;
  gap: 16px;
  flex-wrap: wrap;
}

.filter-item {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.filter-label {
  font-size: 12px;
  font-weight: 500;
  color: var(--color-text-secondary);
}

.filter-actions {
  align-self: flex-end;
}

.table-card {
  background-color: var(--color-bg-card);
  border: 1px solid rgba(255, 255, 255, 0.06);
  margin-bottom: 16px;

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

.submission-id {
  color: var(--color-text-secondary);
  font-size: 13px;
  font-weight: 500;
}

.problem-link {
  color: var(--color-primary);
  text-decoration: none;
  font-weight: 500;
  transition: color 0.2s;

  &:hover {
    color: var(--color-primary-hover);
  }
}

.lang-badge {
  display: inline-block;
  padding: 2px 8px;
  border-radius: 4px;
  background-color: rgba(0, 212, 255, 0.1);
  color: var(--color-primary);
  font-size: 12px;
  font-weight: 500;
}

.time-text {
  font-size: 13px;
  color: var(--color-text-secondary);
}

.pagination-wrapper {
  display: flex;
  justify-content: center;
  padding: 8px 0;
}

@media (max-width: 768px) {
  .filter-bar {
    flex-direction: column;
    align-items: stretch;
  }

  .filter-item {
    width: 100%;
  }
}
</style>
