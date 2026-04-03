<script setup lang="ts">
import { ref, reactive, computed, onMounted, h } from 'vue'
import { useI18n } from 'vue-i18n'
import {
  NDataTable, NInput, NSelect, NButton, NModal, NForm, NFormItem, NTag, NSpace,
  NPopconfirm, NPagination, useMessage,
} from 'naive-ui'
import { adminGetUsers, adminUpdateUser } from '@/api/admin'
import type { User } from '@/types'

const { t } = useI18n()
const message = useMessage()

const loading = ref(false)
const users = ref<User[]>([])
const total = ref(0)

const filters = reactive({
  keyword: '',
  role: '' as string,
  status: '' as string,
  page: 1,
  pageSize: 20,
})

const resetPasswordModal = ref(false)
const resetPasswordUserId = ref<number | null>(null)
const resetPasswordForm = reactive({
  newPassword: '',
  confirmPassword: '',
})

const roleOptions = computed(() => [
  { label: t('admin.allRoles'), value: '' },
  { label: t('admin.roleUser'), value: 'user' },
  { label: t('admin.roleAdmin'), value: 'admin' },
])

const statusOptions = computed(() => [
  { label: t('admin.allStatuses'), value: '' },
  { label: t('admin.statusActive'), value: 'active' },
  { label: t('admin.statusDisabled'), value: 'disabled' },
])

type NTagType = 'default' | 'error' | 'info' | 'primary' | 'success' | 'warning'

const columns = computed(() => [
  {
    title: 'ID',
    key: 'id',
    width: 80,
  },
  {
    title: t('auth.username'),
    key: 'username',
    width: 140,
    ellipsis: { tooltip: true },
  },
  {
    title: t('auth.email'),
    key: 'email',
    width: 200,
    ellipsis: { tooltip: true },
  },
  {
    title: t('admin.role'),
    key: 'role',
    width: 100,
    render(row: User) {
      const type: NTagType = row.role === 'admin' ? 'error' : 'info'
      return h(NTag, { type, size: 'small', bordered: false }, { default: () => row.role })
    },
  },
  {
    title: t('common.status'),
    key: 'status',
    width: 100,
    render(row: User) {
      const isActive = (row as any).status !== 'disabled'
      const type: NTagType = isActive ? 'success' : 'warning'
      return h(
        NTag,
        { type, size: 'small', bordered: false },
        { default: () => isActive ? t('admin.statusActive') : t('admin.statusDisabled') }
      )
    },
  },
  {
    title: t('profile.solvedCount'),
    key: 'solved',
    width: 80,
    render(row: User) {
      return (row as any).solvedCount ?? '-'
    },
  },
  {
    title: t('profile.submissionCount'),
    key: 'submissions',
    width: 100,
    render(row: User) {
      return (row as any).submissionCount ?? '-'
    },
  },
  {
    title: t('teams.createdAt'),
    key: 'createdAt',
    width: 180,
    render(row: User) {
      return new Date(row.createdAt).toLocaleString()
    },
  },
  {
    title: t('common.actions'),
    key: 'actions',
    width: 220,
    render(row: User) {
      return h(NSpace, { size: 'small' }, () => [
        h(
          NSelect,
          {
            value: row.role,
            options: [
              { label: 'User', value: 'user' },
              { label: 'Admin', value: 'admin' },
            ],
            size: 'small',
            style: 'width: 90px',
            onUpdateValue: (val: string) => handleUpdateRole(row.id, val),
          }
        ),
        h(
          NButton,
          {
            size: 'small',
            type: (row as any).status === 'disabled' ? 'success' : 'warning',
            onClick: () => handleToggleStatus(row),
          },
          { default: () => ((row as any).status === 'disabled' ? t('admin.enable') : t('admin.disable')) }
        ),
        h(
          NButton,
          {
            size: 'small',
            onClick: () => openResetPassword(row.id),
          },
          { default: () => t('admin.resetPassword') }
        ),
      ])
    },
  },
])

async function fetchUsers() {
  loading.value = true
  try {
    const params: Record<string, any> = {
      page: filters.page,
      pageSize: filters.pageSize,
    }
    if (filters.keyword) params.keyword = filters.keyword
    if (filters.role) params.role = filters.role
    if (filters.status) params.status = filters.status

    const res = await adminGetUsers(params)
    users.value = res.data.data.items
    total.value = res.data.data.total
  } catch (e: any) {
    message.error(e.message || t('errors.networkError'))
  } finally {
    loading.value = false
  }
}

function handleSearch() {
  filters.page = 1
  fetchUsers()
}

function handleRoleChange(val: string) {
  filters.role = val
  filters.page = 1
  fetchUsers()
}

function handleStatusChange(val: string) {
  filters.status = val
  filters.page = 1
  fetchUsers()
}

function handlePageChange(page: number) {
  filters.page = page
  fetchUsers()
}

async function handleUpdateRole(userId: number, role: string) {
  try {
    await adminUpdateUser(userId, { role: role as any })
    message.success(t('common.success'))
    fetchUsers()
  } catch (e: any) {
    message.error(e.message || t('errors.networkError'))
  }
}

async function handleToggleStatus(user: User) {
  const newStatus = (user as any).status === 'disabled' ? 'active' : 'disabled'
  try {
    await adminUpdateUser(user.id, { status: newStatus } as any)
    message.success(t('common.success'))
    fetchUsers()
  } catch (e: any) {
    message.error(e.message || t('errors.networkError'))
  }
}

function openResetPassword(userId: number) {
  resetPasswordUserId.value = userId
  resetPasswordForm.newPassword = ''
  resetPasswordForm.confirmPassword = ''
  resetPasswordModal.value = true
}

async function handleResetPassword() {
  if (!resetPasswordUserId.value) return
  if (resetPasswordForm.newPassword.length < 8) {
    message.warning(t('auth.passwordMinLength'))
    return
  }
  if (resetPasswordForm.newPassword !== resetPasswordForm.confirmPassword) {
    message.warning(t('auth.passwordMismatch'))
    return
  }
  try {
    await adminUpdateUser(resetPasswordUserId.value, { password: resetPasswordForm.newPassword } as any)
    message.success(t('common.success'))
    resetPasswordModal.value = false
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
    fetchUsers()
  }, 500)
}

onMounted(() => {
  fetchUsers()
})
</script>

<template>
  <div class="admin-user-manage">
    <div class="page-header">
      <h1 class="page-title">{{ t('admin.userManage') }}</h1>
    </div>

    <!-- Filters -->
    <div class="toolbar">
      <NInput
        :placeholder="t('admin.searchUserPlaceholder')"
        :value="filters.keyword"
        clearable
        @update:value="handleKeywordInput"
        class="search-input"
      />
      <NSpace>
        <NSelect
          :value="filters.role"
          :options="roleOptions"
          size="small"
          style="width: 140px"
          @update:value="handleRoleChange"
        />
        <NSelect
          :value="filters.status"
          :options="statusOptions"
          size="small"
          style="width: 140px"
          @update:value="handleStatusChange"
        />
        <NButton size="small" @click="handleSearch">{{ t('common.search') }}</NButton>
      </NSpace>
    </div>

    <!-- User Table -->
    <NCard class="table-card">
      <NDataTable
        :columns="columns"
        :data="users"
        :loading="loading"
        :bordered="false"
        :single-line="false"
        size="small"
        :row-key="(row: User) => row.id"
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

    <!-- Reset Password Modal -->
    <NModal
      v-model:show="resetPasswordModal"
      preset="dialog"
      :title="t('admin.resetPassword')"
      positive-text="Confirm"
      negative-text="Cancel"
      @positive-click="handleResetPassword"
    >
      <NForm>
        <NFormItem :label="t('auth.newPassword')">
          <NInput
            v-model:value="resetPasswordForm.newPassword"
            type="password"
            show-password-on="click"
            :placeholder="t('auth.passwordRequired')"
          />
        </NFormItem>
        <NFormItem :label="t('auth.confirmNewPassword')">
          <NInput
            v-model:value="resetPasswordForm.confirmPassword"
            type="password"
            show-password-on="click"
            :placeholder="t('auth.confirmPassword')"
          />
        </NFormItem>
      </NForm>
    </NModal>
  </div>
</template>

<style scoped lang="scss">
.admin-user-manage {
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

.pagination-wrapper {
  display: flex;
  justify-content: center;
  padding: 16px 0;
}
</style>
