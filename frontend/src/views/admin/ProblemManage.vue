<script setup lang="ts">
import { ref, reactive, computed, onMounted, h } from 'vue'
import { useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import {
  NDataTable, NInput, NSelect, NButton, NSpace, NTag, NPopconfirm, NModal, NPagination, useMessage,
} from 'naive-ui'
import { adminGetProblems, adminDeleteProblem, adminUpdateProblem } from '@/api/admin'
import type { Problem } from '@/types'

const { t } = useI18n()
const router = useRouter()
const message = useMessage()

const loading = ref(false)
const problems = ref<Problem[]>([])
const total = ref(0)

const filters = reactive({
  keyword: '',
  type: '' as string,
  difficulty: '' as string,
  status: '' as string,
  source: '' as string,
  page: 1,
  pageSize: 20,
})

const statusModal = ref(false)
const statusProblemId = ref<number | null>(null)
const statusValue = ref('public')

const typeOptions = computed(() => [
  { label: t('problems.allTypes'), value: '' },
  { label: t('problems.programming'), value: 'programming' },
  { label: t('problems.algorithm'), value: 'algorithm' },
  { label: t('problems.ctf'), value: 'ctf' },
])

const difficultyOptions = computed(() => [
  { label: t('problems.allDifficulties'), value: '' },
  { label: t('problems.easy'), value: 'easy' },
  { label: t('problems.medium'), value: 'medium' },
  { label: t('problems.hard'), value: 'hard' },
  { label: t('problems.expert'), value: 'expert' },
])

const statusFilterOptions = computed(() => [
  { label: t('admin.allStatuses'), value: '' },
  { label: t('admin.statusPublic'), value: 'public' },
  { label: t('admin.statusPrivate'), value: 'private' },
  { label: t('admin.statusDisabled'), value: 'disabled' },
])

const statusOptions = computed(() => [
  { label: t('admin.statusPublic'), value: 'public' },
  { label: t('admin.statusPrivate'), value: 'private' },
  { label: t('admin.statusDisabled'), value: 'disabled' },
])

type NTagType = 'default' | 'error' | 'info' | 'primary' | 'success' | 'warning'

const difficultyColorMap: Record<string, NTagType> = {
  easy: 'success',
  medium: 'warning',
  hard: 'error',
  expert: 'default',
}

const typeColorMap: Record<string, NTagType> = {
  programming: 'info',
  algorithm: 'primary',
  ctf: 'warning',
}

const problemStatusColor: Record<string, NTagType> = {
  public: 'success',
  private: 'warning',
  disabled: 'error',
}

function getAcceptanceRate(problem: Problem): string {
  if (problem.totalSubmit === 0) return '0%'
  return ((problem.acceptedCount / problem.totalSubmit) * 100).toFixed(1) + '%'
}

const columns = computed(() => [
  {
    title: 'ID',
    key: 'id',
    width: 80,
    render(row: Problem) {
      return h('span', { class: 'problem-id' }, `#${row.id}`)
    },
  },
  {
    title: t('problems.problemTitle'),
    key: 'title',
    ellipsis: { tooltip: true },
    render(row: Problem) {
      return h(
        'a',
        {
          class: 'problem-title-link',
          onClick: () => router.push({ name: 'ProblemDetail', params: { id: row.id } }),
        },
        row.title
      )
    },
  },
  {
    title: t('problems.type'),
    key: 'type',
    width: 110,
    render(row: Problem) {
      return h(
        NTag,
        { type: typeColorMap[row.type] || 'default', size: 'small', bordered: false },
        { default: () => t(`problems.${row.type}`) }
      )
    },
  },
  {
    title: t('problems.difficulty'),
    key: 'difficulty',
    width: 100,
    render(row: Problem) {
      return h(
        NTag,
        { type: difficultyColorMap[row.difficulty] || 'default', size: 'small', bordered: false },
        { default: () => t(`problems.${row.difficulty}`) }
      )
    },
  },
  {
    title: t('common.status'),
    key: 'isPublic',
    width: 90,
    render(row: Problem) {
      const status = (row as any).status || (row.isPublic ? 'public' : 'private')
      return h(
        NTag,
        { type: problemStatusColor[status] || 'default', size: 'small', bordered: false },
        { default: () => t(`admin.status${status.charAt(0).toUpperCase() + status.slice(1)}`) }
      )
    },
  },
  {
    title: t('problems.submitCount'),
    key: 'totalSubmit',
    width: 100,
  },
  {
    title: t('problems.acceptanceRate'),
    key: 'acceptance',
    width: 100,
    render(row: Problem) {
      return getAcceptanceRate(row)
    },
  },
  {
    title: t('common.actions'),
    key: 'actions',
    width: 280,
    render(row: Problem) {
      return h(NSpace, { size: 'small' }, () => [
        h(
          NButton,
          { size: 'small', type: 'primary', onClick: () => router.push({ name: 'AdminProblemEdit', params: { id: row.id } }) },
          { default: () => t('common.edit') }
        ),
        h(
          NPopconfirm,
          { onPositiveClick: () => handleDelete(row.id) },
          {
            trigger: () => h(NButton, { size: 'small', type: 'error' }, { default: () => t('common.delete') }),
            default: () => t('admin.confirmDelete'),
          }
        ),
        h(
          NButton,
          { size: 'small', onClick: () => openStatusModal(row.id, (row as any).status || (row.isPublic ? 'public' : 'private')) },
          { default: () => t('admin.setStatus') }
        ),
      ])
    },
  },
])

async function fetchProblems() {
  loading.value = true
  try {
    const params: Record<string, any> = {
      page: filters.page,
      pageSize: filters.pageSize,
    }
    if (filters.keyword) params.keyword = filters.keyword
    if (filters.type) params.type = filters.type
    if (filters.difficulty) params.difficulty = filters.difficulty
    if (filters.status) params.status = filters.status

    const res = await adminGetProblems(params)
    problems.value = res.data.data.items
    total.value = res.data.data.total
  } catch (e: any) {
    message.error(e.message || t('errors.networkError'))
  } finally {
    loading.value = false
  }
}

function handlePageChange(page: number) {
  filters.page = page
  fetchProblems()
}

function handleTypeChange(val: string) {
  filters.type = val
  filters.page = 1
  fetchProblems()
}

function handleDifficultyChange(val: string) {
  filters.difficulty = val
  filters.page = 1
  fetchProblems()
}

function handleStatusFilterChange(val: string) {
  filters.status = val
  filters.page = 1
  fetchProblems()
}

async function handleDelete(id: number) {
  try {
    await adminDeleteProblem(id)
    message.success(t('common.success'))
    fetchProblems()
  } catch (e: any) {
    message.error(e.message || t('errors.networkError'))
  }
}

function openStatusModal(id: number, status: string) {
  statusProblemId.value = id
  statusValue.value = status
  statusModal.value = true
}

async function handleSetStatus() {
  if (!statusProblemId.value) return
  try {
    await adminUpdateProblem(statusProblemId.value, {
      isPublic: statusValue.value === 'public',
      status: statusValue.value,
    } as any)
    message.success(t('common.success'))
    statusModal.value = false
    fetchProblems()
  } catch (e: any) {
    message.error(e.message || t('errors.networkError'))
  }
}

let searchTimer: ReturnType<typeof setTimeout> | null = null
function handleKeywordInput(val: string) {
  filters.keyword = val
  if (searchTimer) clearTimeout(searchTimer)
  searchTimer = setTimeout(() => {
    filters.page = 1
    fetchProblems()
  }, 500)
}

onMounted(() => {
  fetchProblems()
})
</script>

<template>
  <div class="admin-problem-manage">
    <div class="page-header flex-between">
      <h1 class="page-title">{{ t('admin.problemManage') }}</h1>
      <NButton type="primary" @click="router.push('/admin/problems/create')">
        {{ t('admin.createProblem') }}
      </NButton>
    </div>

    <!-- Filters -->
    <div class="toolbar">
      <NInput
        :placeholder="t('admin.searchProblemPlaceholder')"
        :value="filters.keyword"
        clearable
        @update:value="handleKeywordInput"
        class="search-input"
      />
      <NSpace>
        <NSelect
          :value="filters.type"
          :options="typeOptions"
          size="small"
          style="width: 140px"
          @update:value="handleTypeChange"
        />
        <NSelect
          :value="filters.difficulty"
          :options="difficultyOptions"
          size="small"
          style="width: 140px"
          @update:value="handleDifficultyChange"
        />
        <NSelect
          :value="filters.status"
          :options="statusFilterOptions"
          size="small"
          style="width: 140px"
          @update:value="handleStatusFilterChange"
        />
      </NSpace>
    </div>

    <!-- Problem Table -->
    <NCard class="table-card">
      <NDataTable
        :columns="columns"
        :data="problems"
        :loading="loading"
        :bordered="false"
        :single-line="false"
        size="small"
        :row-key="(row: Problem) => row.id"
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

    <!-- Set Status Modal -->
    <NModal
      v-model:show="statusModal"
      preset="dialog"
      :title="t('admin.setStatus')"
      positive-text="Confirm"
      negative-text="Cancel"
      @positive-click="handleSetStatus"
    >
      <NSelect
        v-model:value="statusValue"
        :options="statusOptions"
        style="width: 100%"
      />
    </NModal>
  </div>
</template>

<style scoped lang="scss">
.admin-problem-manage {
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

.search-input {
  flex: 1;
  max-width: 400px;
  min-width: 200px;
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

.problem-id {
  color: var(--color-text-secondary);
  font-size: 13px;
  font-weight: 500;
}

.problem-title-link {
  color: var(--color-primary);
  text-decoration: none;
  font-weight: 500;
  cursor: pointer;

  &:hover {
    color: var(--color-primary-hover);
  }
}

.pagination-wrapper {
  display: flex;
  justify-content: center;
  padding: 16px 0;
}
</style>
