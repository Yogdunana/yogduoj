<script setup lang="ts">
import { ref, reactive, computed, onMounted, watch } from 'vue'
import { useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { NDataTable, NInput, NSelect, NTag, NSpace, NCard, NCheckbox, NPagination, useMessage } from 'naive-ui'
import { listProblems } from '@/api/problem'
import type { Problem, ProblemType, ProblemDifficulty } from '@/types'

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
  tag: '' as string,
  sortBy: 'id' as 'id' | 'acceptanceRate' | 'submissions' | 'difficulty',
  sortOrder: 'asc' as 'asc' | 'desc',
  page: 1,
  pageSize: 20,
})

const allTags = ref<string[]>([])

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

const sortOptions = computed(() => [
  { label: t('problems.sortById'), value: 'id' },
  { label: t('problems.sortByAcceptanceRate'), value: 'acceptanceRate' },
  { label: t('problems.sortBySubmissions'), value: 'submissions' },
  { label: t('problems.sortByDifficulty'), value: 'difficulty' },
])

const tagOptions = computed(() => [
  { label: t('problems.allTags'), value: '' },
  ...allTags.value.map(tag => ({ label: tag, value: tag })),
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

const statusIconMap: Record<string, string> = {
  unsubmitted: '',
  submitted: 'pending',
  accepted: 'success',
}

function getAcceptanceRate(problem: Problem): string {
  if (problem.totalSubmit === 0) return '0%'
  return ((problem.acceptedCount / problem.totalSubmit) * 100).toFixed(1) + '%'
}

const columns = computed(() => [
  {
    title: t('problems.problemId'),
    key: 'id',
    width: 80,
    sorter: false,
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
    filter: false,
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
    title: t('problems.acceptanceRate'),
    key: 'acceptanceRate',
    width: 120,
    sorter: false,
    render(row: Problem) {
      return h('span', {}, getAcceptanceRate(row))
    },
  },
  {
    title: t('problems.submitCount'),
    key: 'totalSubmit',
    width: 100,
    sorter: false,
  },
  {
    title: t('problems.status'),
    key: 'userStatus',
    width: 80,
    render(row: Problem) {
      const status = row.userStatus || 'unsubmitted'
      if (status === 'unsubmitted') return h('span', { class: 'status-icon unsubmitted' }, '-')
      if (status === 'submitted') return h('span', { class: 'status-icon submitted' }, '...')
      return h('span', { class: 'status-icon accepted' }, 'AC')
    },
  },
])

async function fetchProblems() {
  loading.value = true
  try {
    const params: Record<string, any> = {
      page: filters.page,
      pageSize: filters.pageSize,
      sortBy: filters.sortBy,
      sortOrder: filters.sortOrder,
    }
    if (filters.keyword) params.keyword = filters.keyword
    if (filters.type) params.type = filters.type
    if (filters.difficulty) params.difficulty = filters.difficulty
    if (filters.tag) params.tag = filters.tag

    const res = await listProblems(params)
    const data = res.data.data
    problems.value = data.items
    total.value = data.total

    // Collect unique tags
    const tagSet = new Set<string>()
    data.items.forEach(p => p.tags?.forEach(tag => tagSet.add(tag)))
    allTags.value = Array.from(tagSet).sort()
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

function handleSearch() {
  filters.page = 1
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

function handleTagChange(val: string) {
  filters.tag = val
  filters.page = 1
  fetchProblems()
}

function handleSortChange(val: string) {
  filters.sortBy = val as any
  filters.page = 1
  fetchProblems()
}

function handleRowClick(row: Problem) {
  router.push({ name: 'ProblemDetail', params: { id: row.id } })
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
  <div class="problem-list-page">
    <div class="page-header">
      <h1 class="page-title">{{ t('problems.title') }}</h1>
    </div>

    <div class="problem-list-layout">
      <!-- Left Sidebar Filter -->
      <aside class="filter-sidebar">
        <NCard :title="t('common.filter')" size="small" class="filter-card">
          <div class="filter-group">
            <label class="filter-label">{{ t('problems.filterByType') }}</label>
            <NSelect
              :value="filters.type"
              :options="typeOptions"
              size="small"
              @update:value="handleTypeChange"
            />
          </div>

          <div class="filter-group">
            <label class="filter-label">{{ t('problems.filterByDifficulty') }}</label>
            <NSelect
              :value="filters.difficulty"
              :options="difficultyOptions"
              size="small"
              @update:value="handleDifficultyChange"
            />
          </div>

          <div class="filter-group">
            <label class="filter-label">{{ t('problems.filterByTag') }}</label>
            <NSelect
              :value="filters.tag"
              :options="tagOptions"
              size="small"
              @update:value="handleTagChange"
            />
          </div>
        </NCard>
      </aside>

      <!-- Main Content -->
      <main class="problem-main">
        <!-- Search and Sort Bar -->
        <div class="toolbar">
          <NInput
            :placeholder="t('problems.searchPlaceholder')"
            :value="filters.keyword"
            clearable
            @update:value="handleKeywordInput"
            class="search-input"
          />
          <NSpace>
            <NSelect
              :value="filters.sortBy"
              :options="sortOptions"
              size="small"
              style="width: 160px"
              @update:value="handleSortChange"
            />
          </NSpace>
        </div>

        <!-- Problem Table -->
        <NCard class="table-card">
          <NDataTable
            :columns="columns"
            :data="problems"
            :loading="loading"
            :row-props="(row: Problem) => ({
              style: 'cursor: pointer',
              onClick: () => handleRowClick(row),
            })"
            :bordered="false"
            :single-line="false"
            size="small"
            :row-key="(row: Problem) => row.id"
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
      </main>
    </div>
  </div>
</template>

<style scoped lang="scss">
.problem-list-page {
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

.problem-list-layout {
  display: flex;
  gap: 24px;
  align-items: flex-start;
}

.filter-sidebar {
  width: 220px;
  flex-shrink: 0;
}

.filter-card {
  background-color: var(--color-bg-card);
  border: 1px solid rgba(255, 255, 255, 0.06);
}

.filter-group {
  margin-bottom: 16px;

  &:last-child {
    margin-bottom: 0;
  }
}

.filter-label {
  display: block;
  font-size: 13px;
  font-weight: 500;
  color: var(--color-text-secondary);
  margin-bottom: 6px;
}

.problem-main {
  flex: 1;
  min-width: 0;
}

.toolbar {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-bottom: 16px;
}

.search-input {
  flex: 1;
  max-width: 400px;
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

.problem-id {
  color: var(--color-text-secondary);
  font-size: 13px;
  font-weight: 500;
}

.problem-title-link {
  color: var(--color-primary);
  text-decoration: none;
  font-weight: 500;
  transition: color 0.2s;

  &:hover {
    color: var(--color-primary-hover);
  }
}

.status-icon {
  font-size: 13px;
  font-weight: 600;

  &.unsubmitted {
    color: var(--color-text-secondary);
  }

  &.submitted {
    color: var(--color-warning);
  }

  &.accepted {
    color: var(--color-success);
  }
}

.pagination-wrapper {
  display: flex;
  justify-content: center;
  padding: 8px 0;
}

@media (max-width: 768px) {
  .problem-list-layout {
    flex-direction: column;
  }

  .filter-sidebar {
    width: 100%;
  }

  .toolbar {
    flex-direction: column;
    align-items: stretch;
  }

  .search-input {
    max-width: none;
  }
}
</style>
