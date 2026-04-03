<script setup lang="ts">
import { ref, onMounted, computed, h } from 'vue'
import { useI18n } from 'vue-i18n'
import { getUserPublicProfile, getUserSubmissions, getUserContests } from '@/api/user'
import type { User, Submission, Contest, SubmissionStatus } from '@/types'
import { NCard, NDataTable, NButton, NTabs, NTabPane, NStatistic, NTag, NAvatar, NSpace, NSpin, NEmpty, NIcon } from 'naive-ui'
import { PersonOutline } from '@vicons/ionicons5'

const props = defineProps<{ id: string | number }>()
const { t } = useI18n()

const loading = ref(true)
const userProfile = ref<User | null>(null)
const submissions = ref<Submission[]>([])
const contests = ref<Contest[]>([])
const submissionsLoading = ref(false)
const contestsLoading = ref(false)
const submissionTotal = ref(0)
const contestTotal = ref(0)

function formatDate(dateStr: string) {
  return new Date(dateStr).toLocaleDateString()
}

function getStatusType(status: SubmissionStatus) {
  const map: Record<string, 'success' | 'error' | 'warning' | 'info' | 'default'> = {
    accepted: 'success',
    wrong_answer: 'error',
    time_limit_exceeded: 'warning',
    memory_limit_exceeded: 'warning',
    runtime_error: 'error',
    compilation_error: 'error',
    pending: 'info',
    judging: 'info',
    presentation_error: 'warning',
  }
  return map[status] || 'default'
}

function getStatusLabel(status: SubmissionStatus) {
  const key: Record<string, string> = {
    pending: 'submissions.pending',
    judging: 'submissions.judging',
    accepted: 'submissions.accepted',
    wrong_answer: 'submissions.wrongAnswer',
    time_limit_exceeded: 'submissions.timeLimitExceeded',
    memory_limit_exceeded: 'submissions.memoryLimitExceeded',
    runtime_error: 'submissions.runtimeError',
    compilation_error: 'submissions.compilationError',
    presentation_error: 'submissions.presentationError',
  }
  return t(key[status] || 'common.noData')
}

function getContestStatusLabel(status: string) {
  const key: Record<string, string> = {
    upcoming: 'contests.upcoming',
    running: 'contests.running',
    ended: 'contests.ended',
  }
  return t(key[status] || status)
}

function getContestStatusType(status: string) {
  const map: Record<string, 'success' | 'warning' | 'info' | 'default'> = {
    upcoming: 'info',
    running: 'success',
    ended: 'default',
  }
  return map[status] || 'default'
}

const submissionColumns = computed(() => [
  {
    title: t('submissions.problem'),
    key: 'problem',
    render(row: Submission) {
      return row.problem?.title || `#${row.problemId}`
    },
  },
  {
    title: t('submissions.language'),
    key: 'language',
    width: 120,
  },
  {
    title: t('submissions.status'),
    key: 'status',
    width: 160,
    render(row: Submission) {
      return h(NTag, { type: getStatusType(row.status), size: 'small', bordered: false }, { default: () => getStatusLabel(row.status) })
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
    width: 180,
    render(row: Submission) {
      return new Date(row.createdAt).toLocaleString()
    },
  },
])

const contestColumns = computed(() => [
  {
    title: t('contests.title'),
    key: 'title',
    render(row: Contest) {
      return h('router-link', { to: `/contests/${row.id}`, style: { color: 'var(--color-primary)' } }, { default: () => row.title })
    },
  },
  {
    title: t('contests.type'),
    key: 'type',
    width: 100,
    render(row: Contest) {
      return h(NTag, { size: 'small', bordered: false }, { default: () => row.type.toUpperCase() })
    },
  },
  {
    title: 'Rule',
    key: 'ruleType',
    width: 80,
    render(row: Contest) {
      return row.ruleType?.toUpperCase()
    },
  },
  {
    title: t('contests.status'),
    key: 'status',
    width: 120,
    render(row: Contest) {
      return h(NTag, { type: getContestStatusType(row.status), size: 'small', bordered: false }, { default: () => getContestStatusLabel(row.status) })
    },
  },
  {
    title: t('contests.startTime'),
    key: 'startTime',
    width: 180,
    render(row: Contest) {
      return new Date(row.startTime).toLocaleString()
    },
  },
])

async function fetchUserProfile() {
  loading.value = true
  try {
    const res = await getUserPublicProfile(Number(props.id))
    userProfile.value = res.data.data
  } catch {
    // ignore
  } finally {
    loading.value = false
  }
}

async function fetchSubmissions() {
  submissionsLoading.value = true
  try {
    const res = await getUserSubmissions({ userId: Number(props.id), page: 1, pageSize: 20 })
    submissions.value = res.data.data.items
    submissionTotal.value = res.data.data.total
  } catch {
    // ignore
  } finally {
    submissionsLoading.value = false
  }
}

async function fetchContests() {
  contestsLoading.value = true
  try {
    const res = await getUserContests({ userId: Number(props.id), page: 1, pageSize: 20 })
    contests.value = res.data.data.items
    contestTotal.value = res.data.data.total
  } catch {
    // ignore
  } finally {
    contestsLoading.value = false
  }
}

onMounted(() => {
  fetchUserProfile()
  fetchSubmissions()
  fetchContests()
})
</script>

<template>
  <div class="public-profile-page">
    <NSpin :show="loading">
      <!-- Profile Header -->
      <NCard class="profile-header-card">
        <div class="profile-header">
          <div class="profile-avatar-section">
            <NAvatar
              :src="userProfile?.avatar"
              :size="80"
              round
              class="profile-avatar"
            >
              {{ userProfile?.username?.charAt(0)?.toUpperCase() }}
            </NAvatar>
            <div class="profile-info">
              <h1 class="profile-username">{{ userProfile?.username }}</h1>
              <p class="profile-date">{{ t('profile.registrationDate') }}: {{ userProfile?.createdAt ? formatDate(userProfile.createdAt) : '-' }}</p>
            </div>
          </div>
        </div>
      </NCard>

      <!-- Statistics -->
      <div class="stats-row">
        <NCard class="stat-card">
          <NStatistic :label="t('profile.solvedCount')" :value="0" />
        </NCard>
        <NCard class="stat-card">
          <NStatistic :label="t('profile.submissionCount')" :value="submissionTotal" />
        </NCard>
        <NCard class="stat-card">
          <NStatistic :label="t('profile.contestCount')" :value="contestTotal" />
        </NCard>
      </div>

      <!-- Tabs -->
      <NCard class="profile-tabs-card">
        <NTabs type="line" animated>
          <NTabPane :name="'submissions'" :tab="t('profile.submissionsTab')">
            <NSpin :show="submissionsLoading">
              <NDataTable
                v-if="submissions.length > 0"
                :columns="submissionColumns"
                :data="submissions"
                :bordered="false"
                :single-line="false"
                striped
                class="data-table"
              />
              <NEmpty v-else :description="t('submissions.noSubmissions')" class="empty-state" />
            </NSpin>
          </NTabPane>
          <NTabPane :name="'contests'" :tab="t('profile.contestsTab')">
            <NSpin :show="contestsLoading">
              <NDataTable
                v-if="contests.length > 0"
                :columns="contestColumns"
                :data="contests"
                :bordered="false"
                :single-line="false"
                striped
                class="data-table"
              />
              <NEmpty v-else :description="t('contests.noContests')" class="empty-state" />
            </NSpin>
          </NTabPane>
        </NTabs>
      </NCard>
    </NSpin>
  </div>
</template>

<style scoped lang="scss">
.public-profile-page {
  max-width: var(--content-max-width);
  margin: 0 auto;
  padding: 24px 20px;
}

.profile-header-card {
  margin-bottom: 20px;
}

.profile-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  flex-wrap: wrap;
  gap: 16px;
}

.profile-avatar-section {
  display: flex;
  align-items: center;
  gap: 20px;
}

.profile-avatar {
  background-color: var(--color-primary);
  color: #1a1a2e;
  font-size: 32px;
  font-weight: 700;
  flex-shrink: 0;
}

.profile-info {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.profile-username {
  font-size: 24px;
  font-weight: 700;
  color: var(--color-text);
  margin: 0;
}

.profile-date {
  font-size: 13px;
  color: var(--color-text-secondary);
  margin: 0;
}

.stats-row {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
  gap: 16px;
  margin-bottom: 20px;
}

.stat-card {
  text-align: center;
}

.profile-tabs-card {
  min-height: 300px;
}

.data-table {
  margin-top: 16px;
}

.empty-state {
  margin-top: 40px;
}
</style>
