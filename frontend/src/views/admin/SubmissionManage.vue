<script setup lang="ts">
import { ref, reactive, computed, onMounted, h } from 'vue'
import { useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import {
  NDataTable, NInput, NSelect, NButton, NSpace, NTag, NModal, NPagination, useMessage,
} from 'naive-ui'
import { adminGetSubmissions, adminRejudge } from '@/api/admin'
import { getSubmissionCode } from '@/api/submission'
import type { Submission } from '@/types'

const { t } = useI18n()
const router = useRouter()
const message = useMessage()

const loading = ref(false)
const submissions = ref<Submission[]>([])
const total = ref(0)

const filters = reactive({
  userId: '' as string,
  problemId: '' as string,
  contestId: '' as string,
  result: '' as string,
  language: '' as string,
  page: 1,
  pageSize: 20,
})

const codeModal = ref(false)
const codeContent = ref('')

const resultOptions = computed(() => [
  { label: t('submissions.allResults'), value: '' },
  { label: t('submissions.accepted'), value: 'accepted' },
  { label: t('submissions.wrongAnswer'), value: 'wrong_answer' },
  { label: t('submissions.timeLimitExceeded'), value: 'time_limit_exceeded' },
  { label: t('submissions.memoryLimitExceeded'), value: 'memory_limit_exceeded' },
  { label: t('submissions.runtimeError'), value: 'runtime_error' },
  { label: t('submissions.compilationError'), value: 'compilation_error' },
])

const languageOptions = computed(() => [
  { label: t('submissions.allLanguages'), value: '' },
  { label: 'C', value: 'c' },
  { label: 'C++', value: 'cpp' },
  { label: 'Java', value: 'java' },
  { label: 'Python', value: 'python' },
  { label: 'Go', value: 'go' },
  { label: 'Rust', value: 'rust' },
  { label: 'JavaScript', value: 'javascript' },
])

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

const columns = computed(() => [
  {
    title: 'ID',
    key: 'id',
    width: 80,
  },
  {
    title: t('submissions.problem'),
    key: 'problemId',
    width: 80,
    render(row: Submission) {
      return h('a', {
        class: 'link',
        onClick: () => router.push({ name: 'ProblemDetail', params: { id: row.problemId } }),
      }, `#${row.problemId}`)
    },
  },
  {
    title: t('auth.username'),
    key: 'userId',
    width: 120,
    render(row: Submission) {
      return row.user?.username || `User#${row.userId}`
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
    render(row: Submission) {
      return row.timeUsed != null ? `${row.timeUsed}ms` : '-'
    },
  },
  {
    title: t('submissions.memoryUsed'),
    key: 'memoryUsed',
    width: 100,
    render(row: Submission) {
      return row.memoryUsed != null ? `${row.memoryUsed}KB` : '-'
    },
  },
  {
    title: t('submissions.submitTime'),
    key: 'createdAt',
    width: 170,
    render(row: Submission) {
      return new Date(row.createdAt).toLocaleString()
    },
  },
  {
    title: t('common.actions'),
    key: 'actions',
    width: 160,
    render(row: Submission) {
      return h(NSpace, { size: 'small' }, () => [
        h(
          NButton,
          { size: 'small', onClick: () => handleViewCode(row.id) },
          { default: () => t('submissions.code') }
        ),
        h(
          NButton,
          { size: 'small', type: 'warning', onClick: () => handleRejudge(row.id) },
          { default: () => t('submissions.rejudge') }
        ),
      ])
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
    if (filters.userId) params.userId = filters.userId
    if (filters.problemId) params.problemId = filters.problemId
    if (filters.contestId) params.contestId = filters.contestId
    if (filters.result) params.status = filters.result
    if (filters.language) params.language = filters.language

    const res = await adminGetSubmissions(params)
    submissions.value = res.data.data.list
    total.value = res.data.data.total
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

function handleResultChange(val: string) {
  filters.result = val
  filters.page = 1
  fetchSubmissions()
}

function handleLanguageChange(val: string) {
  filters.language = val
  filters.page = 1
  fetchSubmissions()
}

function handleSearch() {
  filters.page = 1
  fetchSubmissions()
}

async function handleViewCode(id: number) {
  try {
    const res = await getSubmissionCode(id)
    codeContent.value = res.data.data.code
    codeModal.value = true
  } catch (e: any) {
    message.error(e.message || t('errors.networkError'))
  }
}

async function handleRejudge(id: number) {
  try {
    await adminRejudge(id)
    message.success(t('admin.rejudgeQueued'))
    fetchSubmissions()
  } catch (e: any) {
    message.error(e.message || t('errors.networkError'))
  }
}

onMounted(() => {
  fetchSubmissions()
})
</script>

<template>
  <div class="admin-submission-manage">
    <div class="page-header">
      <h1 class="page-title">{{ t('admin.submissionManage') }}</h1>
    </div>

    <!-- Filters -->
    <div class="toolbar">
      <NInput
        :placeholder="t('admin.userIdPlaceholder')"
        :value="filters.userId"
        clearable
        @update:value="(v: string) => { filters.userId = v }"
        style="width: 120px"
      />
      <NInput
        :placeholder="t('admin.problemIdPlaceholder')"
        :value="filters.problemId"
        clearable
        @update:value="(v: string) => { filters.problemId = v }"
        style="width: 120px"
      />
      <NSelect
        :value="filters.result"
        :options="resultOptions"
        size="small"
        style="width: 180px"
        @update:value="handleResultChange"
      />
      <NSelect
        :value="filters.language"
        :options="languageOptions"
        size="small"
        style="width: 140px"
        @update:value="handleLanguageChange"
      />
      <NButton size="small" @click="handleSearch">{{ t('common.search') }}</NButton>
    </div>

    <!-- Submission Table -->
    <NCard class="table-card">
      <NDataTable
        :columns="columns"
        :data="submissions"
        :loading="loading"
        :bordered="false"
        :single-line="false"
        size="small"
        :row-key="(row: Submission) => row.id"
        :max-height="600"
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

    <!-- Code Modal -->
    <NModal
      v-model:show="codeModal"
      preset="card"
      :title="t('submissions.sourceCode')"
      style="width: 800px"
    >
      <pre class="code-preview">{{ codeContent }}</pre>
    </NModal>
  </div>
</template>

<style scoped lang="scss">
.admin-submission-manage {
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

.toolbar {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-bottom: 16px;
  flex-wrap: wrap;
}

.table-card {
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

.link {
  color: var(--color-primary);
  cursor: pointer;
  text-decoration: none;

  &:hover {
    color: var(--color-primary-hover);
  }
}

.pagination-wrapper {
  display: flex;
  justify-content: center;
  padding: 16px 0;
}

.code-preview {
  background-color: var(--color-bg);
  padding: 16px;
  border-radius: var(--border-radius);
  font-size: 13px;
  line-height: 1.6;
  overflow-x: auto;
  max-height: 500px;
  overflow-y: auto;
}
</style>
