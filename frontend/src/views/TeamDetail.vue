<script setup lang="ts">
import { ref, reactive, onMounted, computed, h } from 'vue'
import { useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { useAuthStore } from '@/stores/auth'
import { getTeamDetail, getTeamMembers, updateTeam, deleteTeam, leaveTeam, inviteMember, removeMember } from '@/api/team'
import type { Team, TeamMember } from '@/types'
import { NCard, NDataTable, NButton, NInput, NSpace, NSpin, NEmpty, NTag, NAvatar, NModal, NForm, NFormItem, NPopconfirm, NAlert, NIcon, useMessage } from 'naive-ui'
import { PersonAddOutline, TrashOutline, CreateOutline, LogOutOutline } from '@vicons/ionicons5'

const props = defineProps<{ id: string | number }>()
const router = useRouter()
const { t } = useI18n()
const authStore = useAuthStore()
const message = useMessage()

const loading = ref(true)
const team = ref<Team | null>(null)
const members = ref<TeamMember[]>([])
const membersLoading = ref(false)

// Edit team modal
const showEditTeam = ref(false)
const editLoading = ref(false)
const editForm = reactive({
  name: '',
  description: '',
})

// Invite member
const inviteUsername = ref('')
const inviteLoading = ref(false)

const isLeader = computed(() => {
  if (!team.value || !authStore.user) return false
  return team.value.leaderId === authStore.user.id
})

const isMember = computed(() => {
  if (!authStore.user) return false
  return members.value.some(m => m.userId === authStore.user!.id)
})

function formatDate(dateStr: string) {
  return new Date(dateStr).toLocaleDateString()
}

const memberColumns = computed(() => [
  {
    title: t('auth.username'),
    key: 'user',
    render(row: TeamMember) {
      return h('div', { style: { display: 'flex', alignItems: 'center', gap: '8px' } }, [
        h(NAvatar, {
          src: row.user?.avatar,
          size: 28,
          round: true,
        }, { default: () => row.user?.username?.charAt(0)?.toUpperCase() }),
        h('span', {}, row.user?.username || `#${row.userId}`),
      ])
    },
  },
  {
    title: t('teams.leader'),
    key: 'role',
    width: 100,
    render(row: TeamMember) {
      return h(NTag, {
        type: row.role === 'leader' ? 'warning' : 'default',
        size: 'small',
        bordered: false,
      }, { default: () => row.role === 'leader' ? t('teams.leader') : t('teams.member') })
    },
  },
  {
    title: t('teams.createdAt'),
    key: 'joinedAt',
    width: 140,
    render(row: TeamMember) {
      return formatDate(row.joinedAt)
    },
  },
  {
    title: t('common.actions'),
    key: 'actions',
    width: 120,
    render(row: TeamMember) {
      if (!isLeader.value || row.role === 'leader') return null
      return h(NPopconfirm, {
        onPositiveClick: () => handleRemoveMember(row.userId),
      }, {
        trigger: () => h(NButton, { size: 'small', type: 'error', quaternary: true }, {
          default: () => t('teams.removeMember'),
          icon: () => h(NIcon, null, { default: () => h(TrashOutline) }),
        }),
        default: () => t('teams.confirmRemove'),
      })
    },
  },
])

async function fetchTeam() {
  loading.value = true
  try {
    const res = await getTeamDetail(Number(props.id))
    team.value = res.data.data
  } catch {
    // ignore
  } finally {
    loading.value = false
  }
}

async function fetchMembers() {
  membersLoading.value = true
  try {
    const res = await getTeamMembers(Number(props.id), { page: 1, pageSize: 100 })
    members.value = res.data.data.items
  } catch {
    // ignore
  } finally {
    membersLoading.value = false
  }
}

function openEditTeam() {
  if (!team.value) return
  editForm.name = team.value.name
  editForm.description = team.value.description
  showEditTeam.value = true
}

async function handleUpdateTeam() {
  if (!editForm.name.trim()) {
    message.warning(t('teams.teamName') + ' ' + t('auth.usernameRequired'))
    return
  }
  editLoading.value = true
  try {
    const res = await updateTeam(Number(props.id), { name: editForm.name.trim(), description: editForm.description.trim() })
    team.value = res.data.data
    showEditTeam.value = false
    message.success(t('teams.updateSuccess'))
  } catch (err: unknown) {
    message.error((err instanceof Error) ? err.message : t('errors.unknownError'))
  } finally {
    editLoading.value = false
  }
}

async function handleInviteMember() {
  if (!inviteUsername.value.trim()) {
    message.warning(t('teams.inviteUsername') + ' ' + t('auth.usernameRequired'))
    return
  }
  inviteLoading.value = true
  try {
    await inviteMember(Number(props.id), Number(inviteUsername.value.trim()))
    inviteUsername.value = ''
    message.success(t('teams.inviteSuccess'))
    fetchMembers()
  } catch (err: unknown) {
    message.error((err instanceof Error) ? err.message : t('errors.unknownError'))
  } finally {
    inviteLoading.value = false
  }
}

async function handleRemoveMember(userId: number) {
  try {
    await removeMember(Number(props.id), userId)
    message.success(t('teams.removeSuccess'))
    fetchMembers()
  } catch (err: unknown) {
    message.error((err instanceof Error) ? err.message : t('errors.unknownError'))
  }
}

async function handleLeaveTeam() {
  try {
    await leaveTeam(Number(props.id))
    message.success(t('teams.leaveSuccess'))
    router.push('/teams')
  } catch (err: unknown) {
    message.error((err instanceof Error) ? err.message : t('errors.unknownError'))
  }
}

async function handleDisbandTeam() {
  try {
    await deleteTeam(Number(props.id))
    message.success(t('teams.deleteSuccess'))
    router.push('/teams')
  } catch (err: unknown) {
    message.error((err instanceof Error) ? err.message : t('errors.unknownError'))
  }
}

onMounted(() => {
  fetchTeam()
  fetchMembers()
})
</script>

<template>
  <div class="team-detail-page">
    <NSpin :show="loading">
      <!-- Team Info Header -->
      <NCard class="team-header-card">
        <div class="team-header">
          <div class="team-info">
            <h1 class="team-name">{{ team?.name }}</h1>
            <p class="team-slogan" v-if="team?.description">{{ team.description }}</p>
            <p class="team-meta">
              <span>{{ t('teams.leader') }}: {{ team?.leader?.username || '-' }}</span>
              <span class="meta-sep">|</span>
              <span>{{ t('teams.memberCount') }}: {{ team?.memberCount || members.length }}</span>
              <span class="meta-sep">|</span>
              <span>{{ t('teams.createdAt') }}: {{ team?.createdAt ? formatDate(team.createdAt) : '-' }}</span>
            </p>
          </div>
          <NSpace v-if="authStore.isLoggedIn" vertical>
            <template v-if="isLeader">
              <NSpace>
                <NButton type="primary" @click="openEditTeam">
                  <template #icon>
                    <NIcon><CreateOutline /></NIcon>
                  </template>
                  {{ t('teams.editTeam') }}
                </NButton>
                <NPopconfirm @positive-click="handleDisbandTeam">
                  <template #trigger>
                    <NButton type="error">
                      <template #icon>
                        <NIcon><TrashOutline /></NIcon>
                      </template>
                      {{ t('teams.disbandTeam') }}
                    </NButton>
                  </template>
                  {{ t('teams.confirmDisband') }}
                </NPopconfirm>
              </NSpace>
            </template>
            <NPopconfirm v-else-if="isMember" @positive-click="handleLeaveTeam">
              <template #trigger>
                <NButton type="warning">
                  <template #icon>
                    <NIcon><LogOutOutline /></NIcon>
                  </template>
                  {{ t('teams.leaveTeam') }}
                </NButton>
              </template>
              {{ t('teams.confirmLeave') }}
            </NPopconfirm>
          </NSpace>
        </div>
      </NCard>

      <!-- Invite Member (Leader only) -->
      <NCard v-if="isLeader" class="invite-card">
        <div class="invite-section">
          <h3 class="section-title">{{ t('teams.inviteMember') }}</h3>
          <NSpace align="center" class="invite-form">
            <NInput
              v-model:value="inviteUsername"
              :placeholder="t('teams.inviteUsername')"
              style="max-width: 240px;"
              @keyup.enter="handleInviteMember"
            />
            <NButton type="primary" :loading="inviteLoading" @click="handleInviteMember">
              <template #icon>
                <NIcon><PersonAddOutline /></NIcon>
              </template>
              {{ t('teams.inviteMember') }}
            </NButton>
          </NSpace>
        </div>
      </NCard>

      <!-- Members List -->
      <NCard class="members-card">
        <h3 class="section-title">{{ t('teams.members') }}</h3>
        <NSpin :show="membersLoading">
          <NDataTable
            v-if="members.length > 0"
            :columns="memberColumns"
            :data="members"
            :bordered="false"
            :single-line="false"
            striped
            class="data-table"
          />
          <NEmpty v-else :description="t('teams.noMembers')" class="empty-state" />
        </NSpin>
      </NCard>
    </NSpin>

    <!-- Edit Team Modal -->
    <NModal
      v-model:show="showEditTeam"
      preset="card"
      :title="t('teams.editTeam')"
      class="modal-card"
      style="max-width: 500px;"
    >
      <NForm label-placement="top">
        <NFormItem :label="t('teams.teamName')">
          <NInput
            v-model:value="editForm.name"
            :placeholder="t('teams.teamName')"
            maxlength="50"
          />
        </NFormItem>
        <NFormItem :label="t('teams.slogan')">
          <NInput
            v-model:value="editForm.description"
            type="textarea"
            :placeholder="t('teams.slogan')"
            :rows="3"
            maxlength="200"
          />
        </NFormItem>
      </NForm>
      <template #footer>
        <NSpace justify="end">
          <NButton @click="showEditTeam = false">{{ t('common.cancel') }}</NButton>
          <NButton type="primary" :loading="editLoading" @click="handleUpdateTeam">
            {{ t('common.save') }}
          </NButton>
        </NSpace>
      </template>
    </NModal>
  </div>
</template>

<style scoped lang="scss">
.team-detail-page {
  max-width: var(--content-max-width);
  margin: 0 auto;
  padding: 24px 20px;
}

.team-header-card {
  margin-bottom: 20px;
}

.team-header {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  flex-wrap: wrap;
  gap: 16px;
}

.team-info {
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.team-name {
  font-size: 24px;
  font-weight: 700;
  color: var(--color-text);
  margin: 0;
}

.team-slogan {
  font-size: 14px;
  color: var(--color-text-secondary);
  margin: 0;
  font-style: italic;
}

.team-meta {
  font-size: 13px;
  color: var(--color-text-secondary);
  margin: 0;
  display: flex;
  align-items: center;
  gap: 8px;
  flex-wrap: wrap;
}

.meta-sep {
  opacity: 0.4;
}

.invite-card {
  margin-bottom: 20px;
}

.section-title {
  font-size: 16px;
  font-weight: 600;
  color: var(--color-text);
  margin: 0 0 12px 0;
}

.invite-section {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.invite-form {
  flex-wrap: wrap;
}

.members-card {
  min-height: 200px;
}

.data-table {
  margin-top: 8px;
}

.empty-state {
  margin-top: 40px;
}
</style>
