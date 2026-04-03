<script setup lang="ts">
import { ref, reactive, onMounted, computed, h } from 'vue'
import { useI18n } from 'vue-i18n'
import { useAuthStore } from '@/stores/auth'
import { updateProfile, changePassword, getUserSubmissions, getUserContests } from '@/api/user'
import type { Submission, Contest, SubmissionStatus } from '@/types'
import { NCard, NDataTable, NModal, NForm, NFormItem, NInput, NButton, NTabs, NTabPane, NStatistic, NTag, NAvatar, NSpace, NAlert, NSpin, NEmpty, NIcon, useMessage } from 'naive-ui'
import { PersonOutline, MailOutline, ImageOutline, LockClosedOutline } from '@vicons/ionicons5'

const { t } = useI18n()
const authStore = useAuthStore()
const message = useMessage()

const loading = ref(false)
const submissions = ref<Submission[]>([])
const contests = ref<Contest[]>([])
const submissionsLoading = ref(false)
const contestsLoading = ref(false)
const submissionPage = ref(1)
const contestPage = ref(1)
const submissionTotal = ref(0)
const contestTotal = ref(0)

// Edit profile modal
const showEditProfile = ref(false)
const editLoading = ref(false)
const editForm = reactive({
  username: '',
  email: '',
  avatar: '',
})
const usernameChanged = ref(false)

// Change password modal
const showChangePassword = ref(false)
const passwordLoading = ref(false)
const passwordForm = reactive({
  oldPassword: '',
  newPassword: '',
  confirmNewPassword: '',
})

const user = computed(() => authStore.user)

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

function getContestTypeLabel(type: string) {
  return type.toUpperCase()
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
      return h(NTag, { size: 'small', bordered: false }, { default: () => getContestTypeLabel(row.type) })
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

async function fetchSubmissions() {
  submissionsLoading.value = true
  try {
    const res = await getUserSubmissions({ page: submissionPage.value, pageSize: 20 })
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
    const res = await getUserContests({ page: contestPage.value, pageSize: 20 })
    contests.value = res.data.data.items
    contestTotal.value = res.data.data.total
  } catch {
    // ignore
  } finally {
    contestsLoading.value = false
  }
}

function openEditProfile() {
  if (!user.value) return
  editForm.username = user.value.username
  editForm.email = user.value.email
  editForm.avatar = user.value.avatar || ''
  showEditProfile.value = true
}

async function handleUpdateProfile() {
  editLoading.value = true
  try {
    const data: Record<string, string> = { email: editForm.email }
    if (editForm.avatar) {
      data.avatar = editForm.avatar
    }
    if (!usernameChanged.value && editForm.username !== user.value?.username) {
      data.username = editForm.username
    }
    const res = await updateProfile(data)
    authStore.user = res.data.data
    localStorage.setItem('user', JSON.stringify(res.data.data))
    if (data.username) {
      usernameChanged.value = true
    }
    showEditProfile.value = false
    message.success(t('profile.profileUpdated'))
  } catch (err: unknown) {
    message.error((err instanceof Error) ? err.message : t('errors.unknownError'))
  } finally {
    editLoading.value = false
  }
}

function openChangePassword() {
  passwordForm.oldPassword = ''
  passwordForm.newPassword = ''
  passwordForm.confirmNewPassword = ''
  showChangePassword.value = true
}

async function handleChangePassword() {
  if (!passwordForm.oldPassword) {
    message.warning(t('auth.passwordRequired'))
    return
  }
  if (passwordForm.newPassword.length < 8) {
    message.warning(t('auth.passwordMinLength'))
    return
  }
  if (passwordForm.newPassword !== passwordForm.confirmNewPassword) {
    message.warning(t('auth.passwordMismatch'))
    return
  }
  passwordLoading.value = true
  try {
    await changePassword({ oldPassword: passwordForm.oldPassword, newPassword: passwordForm.newPassword })
    showChangePassword.value = false
    message.success(t('profile.passwordChanged'))
  } catch (err: unknown) {
    message.error((err instanceof Error) ? err.message : t('errors.unknownError'))
  } finally {
    passwordLoading.value = false
  }
}

onMounted(() => {
  fetchSubmissions()
  fetchContests()
})
</script>

<template>
  <div class="profile-page">
    <NSpin :show="loading">
      <!-- Profile Header -->
      <NCard class="profile-header-card">
        <div class="profile-header">
          <div class="profile-avatar-section">
            <NAvatar
              :src="user?.avatar"
              :size="80"
              round
              class="profile-avatar"
            >
              {{ user?.username?.charAt(0)?.toUpperCase() }}
            </NAvatar>
            <div class="profile-info">
              <h1 class="profile-username">{{ user?.username }}</h1>
              <p class="profile-email">{{ user?.email }}</p>
              <p class="profile-date">{{ t('profile.registrationDate') }}: {{ user?.createdAt ? formatDate(user.createdAt) : '-' }}</p>
            </div>
          </div>
          <NSpace>
            <NButton type="primary" @click="openEditProfile">
              <template #icon>
                <NIcon><PersonOutline /></NIcon>
              </template>
              {{ t('profile.editProfile') }}
            </NButton>
            <NButton @click="openChangePassword">
              <template #icon>
                <NIcon><LockClosedOutline /></NIcon>
              </template>
              {{ t('profile.changePassword') }}
            </NButton>
          </NSpace>
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

    <!-- Edit Profile Modal -->
    <NModal
      v-model:show="showEditProfile"
      preset="card"
      :title="t('profile.editProfile')"
      class="modal-card"
      style="max-width: 500px;"
    >
      <NForm label-placement="top">
        <NFormItem :label="t('profile.username')">
          <NInput
            v-model:value="editForm.username"
            :placeholder="t('auth.username')"
            :disabled="usernameChanged"
          />
        </NFormItem>
        <NAlert v-if="usernameChanged" type="warning" class="mb-1">
          {{ t('profile.usernameChanged') }}
        </NAlert>
        <NAlert v-else type="warning" class="mb-1">
          {{ t('profile.usernameWarning') }}
        </NAlert>
        <NFormItem :label="t('profile.email')">
            <NInput
              v-model:value="editForm.email"
              :placeholder="t('auth.email')"
            />
          </NFormItem>
        <NFormItem :label="t('profile.avatar')">
          <NInput
            v-model:value="editForm.avatar"
            :placeholder="t('profile.avatar')"
          />
        </NFormItem>
      </NForm>
      <template #footer>
        <NSpace justify="end">
          <NButton @click="showEditProfile = false">{{ t('common.cancel') }}</NButton>
          <NButton type="primary" :loading="editLoading" @click="handleUpdateProfile">
            {{ t('common.save') }}
          </NButton>
        </NSpace>
      </template>
    </NModal>

    <!-- Change Password Modal -->
    <NModal
      v-model:show="showChangePassword"
      preset="card"
      :title="t('profile.changePassword')"
      class="modal-card"
      style="max-width: 500px;"
    >
      <NForm label-placement="top">
        <NFormItem :label="t('auth.oldPassword')">
          <NInput
            v-model:value="passwordForm.oldPassword"
            type="password"
            show-password-on="click"
            :placeholder="t('auth.oldPassword')"
          />
        </NFormItem>
        <NFormItem :label="t('auth.newPassword')">
          <NInput
            v-model:value="passwordForm.newPassword"
            type="password"
            show-password-on="click"
            :placeholder="t('auth.newPassword')"
          />
        </NFormItem>
        <NFormItem :label="t('auth.confirmNewPassword')">
          <NInput
            v-model:value="passwordForm.confirmNewPassword"
            type="password"
            show-password-on="click"
            :placeholder="t('auth.confirmNewPassword')"
          />
        </NFormItem>
      </NForm>
      <template #footer>
        <NSpace justify="end">
          <NButton @click="showChangePassword = false">{{ t('common.cancel') }}</NButton>
          <NButton type="primary" :loading="passwordLoading" @click="handleChangePassword">
            {{ t('common.confirm') }}
          </NButton>
        </NSpace>
      </template>
    </NModal>
  </div>
</template>

<style scoped lang="scss">
.profile-page {
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

.profile-email {
  font-size: 14px;
  color: var(--color-text-secondary);
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

.mb-1 {
  margin-bottom: 8px;
}
</style>
