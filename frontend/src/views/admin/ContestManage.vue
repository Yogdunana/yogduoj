<script setup lang="ts">
import { ref, reactive, computed, onMounted, h } from 'vue'
import { useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import {
  NDataTable, NInput, NSelect, NButton, NSpace, NTag, NPopconfirm, NModal, NPagination, useMessage,
} from 'naive-ui'
import { adminGetContests, adminDeleteContest } from '@/api/admin'
import type { Contest } from '@/types'

const { t } = useI18n()
const router = useRouter()
const message = useMessage()

const loading = ref(false)
const contests = ref<Contest[]>([])
const total = ref(0)

const filters = reactive({
  status: '' as string,
  type: '' as string,
  ruleType: '' as string,
  page: 1,
  pageSize: 20,
})

const statusOptions = computed(() => [
  { label: t('admin.allStatuses'), value: '' },
  { label: t('contests.upcoming'), value: 'upcoming' },
  { label: t('contests.running'), value: 'running' },
  { label: t('contests.ended'), value: 'ended' },
])

const typeOptions = computed(() => [
  { label: t('admin.allTypes'), value: '' },
  { label: 'ICPC', value: 'icpc' },
  { label: 'IOI', value: 'ioi' },
  { label: 'CTF', value: 'ctf' },
])

const ruleTypeOptions = computed(() => [
  { label: t('admin.allRules'), value: '' },
  { label: 'ACM', value: 'acm' },
  { label: 'OI', value: 'oi' },
  { label: 'Codeforces', value: 'cf' },
])

type NTagType = 'default' | 'error' | 'info' | 'primary' | 'success' | 'warning'

const contestStatusColor: Record<string, NTagType> = {
  upcoming: 'info',
  running: 'success',
  ended: 'default',
}

const typeColorMap: Record<string, NTagType> = {
  icpc: 'primary',
  ioi: 'info',
  ctf: 'warning',
}

const ruleColorMap: Record<string, NTagType> = {
  acm: 'success',
  oi: 'info',
  cf: 'warning',
}

const columns = computed(() => [
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
    title: t('contests.type'),
    key: 'type',
    width: 80,
    render(row: Contest) {
      return h(NTag, { type: typeColorMap[row.type] || 'default', size: 'small', bordered: false }, { default: () => row.type.toUpperCase() })
    },
  },
  {
    title: t('admin.ruleType'),
    key: 'ruleType',
    width: 100,
    render(row: Contest) {
      return h(NTag, { type: ruleColorMap[row.ruleType] || 'default', size: 'small', bordered: false }, { default: () => row.ruleType.toUpperCase() })
    },
  },
  {
    title: t('contests.status'),
    key: 'status',
    width: 100,
    render(row: Contest) {
      return h(NTag, { type: contestStatusColor[row.status] || 'default', size: 'small', bordered: false }, { default: () => t(`contests.${row.status}`) })
    },
  },
  {
    title: t('contests.startTime'),
    key: 'startTime',
    width: 170,
    render(row: Contest) {
      return new Date(row.startTime).toLocaleString()
    },
  },
  {
    title: t('contests.endTime'),
    key: 'endTime',
    width: 170,
    render(row: Contest) {
      return new Date(row.endTime).toLocaleString()
    },
  },
  {
    title: t('contests.participants'),
    key: 'participantCount',
    width: 100,
  },
  {
    title: t('common.actions'),
    key: 'actions',
    width: 240,
    render(row: Contest) {
      return h(NSpace, { size: 'small' }, () => [
        h(
          NButton,
          { size: 'small', type: 'primary', onClick: () => router.push({ name: 'AdminContestEdit', params: { id: row.id } }) },
          { default: () => t('common.edit') }
        ),
        h(
          NButton,
          { size: 'small', onClick: () => router.push({ name: 'ContestRanking', params: { id: row.id } }) },
          { default: () => t('contests.ranking') }
        ),
        h(
          NPopconfirm,
          { onPositiveClick: () => handleDelete(row.id) },
          {
            trigger: () => h(NButton, { size: 'small', type: 'error' }, { default: () => t('common.delete') }),
            default: () => t('admin.confirmDelete'),
          }
        ),
      ])
    },
  },
])

async function fetchContests() {
  loading.value = true
  try {
    const params: Record<string, any> = {
      page: filters.page,
      pageSize: filters.pageSize,
    }
    if (filters.status) params.status = filters.status
    if (filters.type) params.type = filters.type
    if (filters.ruleType) params.ruleType = filters.ruleType

    const res = await adminGetContests(params)
    contests.value = res.data.data.list
    total.value = res.data.data.total
  } catch (e: any) {
    message.error(e.message || t('errors.networkError'))
  } finally {
    loading.value = false
  }
}

function handlePageChange(page: number) {
  filters.page = page
  fetchContests()
}

function handleStatusChange(val: string) {
  filters.status = val
  filters.page = 1
  fetchContests()
}

function handleTypeChange(val: string) {
  filters.type = val
  filters.page = 1
  fetchContests()
}

function handleRuleTypeChange(val: string) {
  filters.ruleType = val
  filters.page = 1
  fetchContests()
}

async function handleDelete(id: number) {
  try {
    await adminDeleteContest(id)
    message.success(t('common.success'))
    fetchContests()
  } catch (e: any) {
    message.error(e.message || t('errors.networkError'))
  }
}

onMounted(() => {
  fetchContests()
})
</script>

<template>
  <div class="admin-contest-manage">
    <div class="page-header flex-between">
      <h1 class="page-title">{{ t('admin.contestManage') }}</h1>
      <NButton type="primary" @click="router.push('/admin/contests/create')">
        {{ t('admin.createContest') }}
      </NButton>
    </div>

    <!-- Filters -->
    <div class="toolbar">
      <NSpace>
        <NSelect
          :value="filters.status"
          :options="statusOptions"
          size="small"
          style="width: 140px"
          @update:value="handleStatusChange"
        />
        <NSelect
          :value="filters.type"
          :options="typeOptions"
          size="small"
          style="width: 140px"
          @update:value="handleTypeChange"
        />
        <NSelect
          :value="filters.ruleType"
          :options="ruleTypeOptions"
          size="small"
          style="width: 160px"
          @update:value="handleRuleTypeChange"
        />
      </NSpace>
    </div>

    <!-- Contest Table -->
    <NCard class="table-card">
      <NDataTable
        :columns="columns"
        :data="contests"
        :loading="loading"
        :bordered="false"
        :single-line="false"
        size="small"
        :row-key="(row: Contest) => row.id"
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
  </div>
</template>

<style scoped lang="scss">
.admin-contest-manage {
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

.pagination-wrapper {
  display: flex;
  justify-content: center;
  padding: 16px 0;
}
</style>
