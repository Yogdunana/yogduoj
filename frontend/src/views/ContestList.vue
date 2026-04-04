<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { NCard, NGrid, NGridItem, NTag, NButton, NInput, NSelect, NSpace, NEmpty, NPagination, useMessage } from 'naive-ui'
import { listContests, signupContest } from '@/api/contest'
import { useAuthStore } from '@/stores/auth'
import type { Contest, ContestStatus, ContestType, ContestRuleType } from '@/types'

const { t } = useI18n()
const router = useRouter()
const message = useMessage()
const authStore = useAuthStore()

const loading = ref(false)
const contests = ref<Contest[]>([])
const total = ref(0)
const signupLoading = ref<number | null>(null)

const filters = reactive({
  keyword: '',
  status: '' as ContestStatus | '',
  type: '' as ContestType | '',
  category: '' as string,
  ruleType: '' as ContestRuleType | '',
  page: 1,
  pageSize: 12,
})

type NTagType = 'default' | 'error' | 'info' | 'primary' | 'success' | 'warning'

const statusOptions = computed(() => [
  { label: t('contests.allStatuses'), value: '' },
  { label: t('contests.upcoming'), value: 'upcoming' },
  { label: t('contests.running'), value: 'running' },
  { label: t('contests.ended'), value: 'ended' },
])

const typeOptions = computed(() => [
  { label: t('contests.allTypes'), value: '' },
  { label: t('contests.individual'), value: 'individual' },
  { label: t('contests.team'), value: 'team' },
])

const categoryOptions = computed(() => [
  { label: t('contests.allCategories'), value: '' },
  { label: t('contests.programming'), value: 'programming' },
  { label: t('contests.algorithm'), value: 'algorithm' },
  { label: t('contests.ctf'), value: 'ctf' },
])

const ruleTypeOptions = computed(() => [
  { label: t('contests.allRuleTypes'), value: '' },
  { label: t('contests.acm'), value: 'acm' },
  { label: t('contests.oi'), value: 'oi' },
  { label: t('contests.cf'), value: 'cf' },
])

const statusColorMap: Record<string, NTagType> = {
  upcoming: 'info',
  running: 'success',
  ended: 'default',
}

const typeColorMap: Record<string, NTagType> = {
  individual: 'primary',
  team: 'warning',
}

const categoryColorMap: Record<string, NTagType> = {
  programming: 'info',
  algorithm: 'primary',
  ctf: 'warning',
}

const ruleTypeColorMap: Record<string, NTagType> = {
  acm: 'error',
  oi: 'success',
  cf: 'info',
}

function formatTime(dateStr: string): string {
  const date = new Date(dateStr)
  return date.toLocaleString()
}

function getDuration(start: string, end: string): string {
  const startTime = new Date(start).getTime()
  const endTime = new Date(end).getTime()
  const diff = endTime - startTime
  const hours = Math.floor(diff / (1000 * 60 * 60))
  const minutes = Math.floor((diff % (1000 * 60 * 60)) / (1000 * 60))
  if (hours >= 24) {
    const days = Math.floor(hours / 24)
    const remainHours = hours % 24
    return `${days}d ${remainHours}h`
  }
  return `${hours}h ${minutes}m`
}

async function fetchContests() {
  loading.value = true
  try {
    const params: Record<string, any> = {
      page: filters.page,
      pageSize: filters.pageSize,
    }
    if (filters.keyword) params.keyword = filters.keyword
    if (filters.status) params.status = filters.status
    if (filters.type) params.type = filters.type
    if (filters.category) params.category = filters.category
    if (filters.ruleType) params.ruleType = filters.ruleType

    const res = await listContests(params)
    const data = res.data.data
    contests.value = data.list
    total.value = data.total
  } catch (e: any) {
    message.error(e.message || t('errors.networkError'))
  } finally {
    loading.value = false
  }
}

async function handleSignup(contestId: number) {
  signupLoading.value = contestId
  try {
    await signupContest(contestId)
    message.success(t('contests.signupSuccess'))
    fetchContests()
  } catch (e: any) {
    message.error(e.message || t('errors.networkError'))
  } finally {
    signupLoading.value = null
  }
}

function handlePageChange(page: number) {
  filters.page = page
  fetchContests()
}

let searchTimer: ReturnType<typeof setTimeout> | null = null
function handleKeywordInput(val: string) {
  filters.keyword = val
  if (searchTimer) clearTimeout(searchTimer)
  searchTimer = setTimeout(() => {
    filters.page = 1
    fetchContests()
  }, 500)
}

function handleFilterChange() {
  filters.page = 1
  fetchContests()
}

function goToContest(id: number) {
  router.push({ name: 'ContestDetail', params: { id } })
}

onMounted(() => {
  fetchContests()
})
</script>

<template>
  <div class="contest-list-page">
    <div class="page-header">
      <h1 class="page-title">{{ t('contests.title') }}</h1>
    </div>

    <!-- Filter Bar -->
    <div class="filter-bar">
      <NInput
        :placeholder="t('contests.searchPlaceholder')"
        :value="filters.keyword"
        clearable
        @update:value="handleKeywordInput"
        class="search-input"
      />
      <NSpace>
        <NSelect
          :value="filters.status"
          :options="statusOptions"
          size="small"
          style="width: 140px"
          @update:value="(val: string) => { filters.status = val as any; handleFilterChange() }"
        />
        <NSelect
          :value="filters.type"
          :options="typeOptions"
          size="small"
          style="width: 130px"
          @update:value="(val: string) => { filters.type = val as any; handleFilterChange() }"
        />
        <NSelect
          :value="filters.category"
          :options="categoryOptions"
          size="small"
          style="width: 140px"
          @update:value="(val: string) => { filters.category = val; handleFilterChange() }"
        />
        <NSelect
          :value="filters.ruleType"
          :options="ruleTypeOptions"
          size="small"
          style="width: 140px"
          @update:value="(val: string) => { filters.ruleType = val as any; handleFilterChange() }"
        />
      </NSpace>
    </div>

    <!-- Contest Cards Grid -->
    <div v-if="contests.length > 0" class="contest-grid">
      <NCard
        v-for="contest in contests"
        :key="contest.id"
        class="contest-card"
        hoverable
        @click="goToContest(contest.id)"
      >
        <div class="contest-card-header">
          <h3 class="contest-title">{{ contest.title }}</h3>
          <NTag :type="statusColorMap[contest.status] || 'default'" size="small" bordered>
            {{ t(`contests.${contest.status}`) }}
          </NTag>
        </div>

        <div class="contest-badges">
          <NTag :type="typeColorMap[contest.type] || 'default'" size="tiny" bordered>
            {{ t(`contests.${contest.type}`) }}
          </NTag>
          <NTag :type="categoryColorMap[contest.type] || 'default'" size="tiny" bordered>
            {{ t(`contests.${contest.type}`) }}
          </NTag>
          <NTag :type="ruleTypeColorMap[contest.ruleType] || 'default'" size="tiny" bordered>
            {{ t(`contests.${contest.ruleType}`) }}
          </NTag>
        </div>

        <div class="contest-time">
          <div class="time-row">
            <span class="time-label">{{ t('contests.startTime') }}:</span>
            <span class="time-value">{{ formatTime(contest.startTime) }}</span>
          </div>
          <div class="time-row">
            <span class="time-label">{{ t('contests.endTime') }}:</span>
            <span class="time-value">{{ formatTime(contest.endTime) }}</span>
          </div>
          <div class="time-row">
            <span class="time-label">{{ t('contests.duration') }}:</span>
            <span class="time-value">{{ getDuration(contest.startTime, contest.endTime) }}</span>
          </div>
        </div>

        <div class="contest-footer">
          <span class="participant-count">
            {{ t('contests.participantCount') }}: {{ contest.participantCount }}
          </span>
          <NSpace v-if="contest.status === 'upcoming'" @click.stop>
            <NButton
              type="primary"
              size="small"
              :loading="signupLoading === contest.id"
              @click="handleSignup(contest.id)"
            >
              {{ t('contests.signup') }}
            </NButton>
          </NSpace>
          <NTag v-else-if="contest.status === 'running'" type="success" size="small" bordered>
            {{ t('contests.enterContest') }}
          </NTag>
          <NButton v-else type="default" size="small" @click.stop="goToContest(contest.id)">
            {{ t('contests.viewRanking') }}
          </NButton>
        </div>
      </NCard>
    </div>

    <NEmpty v-else-if="!loading" :description="t('contests.noContests')" class="empty-state" />

    <!-- Pagination -->
    <div v-if="total > 0" class="pagination-wrapper">
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
.contest-list-page {
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

.filter-bar {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-bottom: 24px;
  flex-wrap: wrap;
}

.search-input {
  flex: 1;
  min-width: 200px;
  max-width: 400px;
}

.contest-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(340px, 1fr));
  gap: 20px;
  margin-bottom: 24px;
}

.contest-card {
  background-color: var(--color-bg-card);
  border: 1px solid rgba(255, 255, 255, 0.06);
  cursor: pointer;
  transition: border-color 0.2s, transform 0.15s;

  &:hover {
    border-color: var(--color-primary);
    transform: translateY(-2px);
  }

  :deep(.n-card__content) {
    padding: 20px;
  }
}

.contest-card-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
  margin-bottom: 12px;
}

.contest-title {
  font-size: 17px;
  font-weight: 600;
  color: var(--color-text);
  margin: 0;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  flex: 1;
}

.contest-badges {
  display: flex;
  gap: 6px;
  margin-bottom: 16px;
  flex-wrap: wrap;
}

.contest-time {
  margin-bottom: 16px;
}

.time-row {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 13px;
  line-height: 1.8;
}

.time-label {
  color: var(--color-text-secondary);
  white-space: nowrap;
}

.time-value {
  color: var(--color-text);
}

.contest-footer {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding-top: 12px;
  border-top: 1px solid rgba(255, 255, 255, 0.06);
}

.participant-count {
  font-size: 13px;
  color: var(--color-text-secondary);
}

.empty-state {
  margin-top: 60px;
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

  .search-input {
    max-width: none;
  }

  .contest-grid {
    grid-template-columns: 1fr;
  }
}
</style>
