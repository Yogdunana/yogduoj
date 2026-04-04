<script setup lang="ts">
import { ref, reactive, computed, onMounted, h } from 'vue'
import { useI18n } from 'vue-i18n'
import {
  NDataTable, NButton, NSpace, NTag, NPopconfirm, NModal, NForm, NFormItem,
  NInput, NSwitch, NPagination, useMessage,
} from 'naive-ui'
import { adminGetAnnouncements } from '@/api/admin'
import { createAnnouncement, updateAnnouncement, deleteAnnouncement } from '@/api/announcement'
import type { Announcement } from '@/types'

const { t } = useI18n()
const message = useMessage()

const loading = ref(false)
const announcements = ref<any[]>([])
const total = ref(0)

const filters = reactive({
  page: 1,
  pageSize: 20,
})

const modalVisible = ref(false)
const editingId = ref<number | null>(null)
const form = reactive({
  title: '',
  content: '',
  isPinned: false,
})

const columns = computed(() => [
  {
    title: 'ID',
    key: 'id',
    width: 80,
  },
  {
    title: t('nav.announcements'),
    key: 'title',
    ellipsis: { tooltip: true },
  },
  {
    title: t('admin.pinned'),
    key: 'isPinned',
    width: 80,
    render(row: any) {
      return h(
        NTag,
        { type: row.isPinned ? 'success' : 'default', size: 'small', bordered: false },
        { default: () => (row.isPinned ? t('common.yes') : t('common.no')) }
      )
    },
  },
  {
    title: t('admin.createdBy'),
    key: 'authorId',
    width: 120,
    render(row: any) {
      return row.author?.username || `Admin#${row.authorId}`
    },
  },
  {
    title: t('teams.createdAt'),
    key: 'createdAt',
    width: 170,
    render(row: any) {
      return new Date(row.createdAt).toLocaleString()
    },
  },
  {
    title: t('common.actions'),
    key: 'actions',
    width: 220,
    render(row: any) {
      return h(NSpace, { size: 'small' }, () => [
        h(
          NButton,
          { size: 'small', type: 'primary', onClick: () => openEditModal(row) },
          { default: () => t('common.edit') }
        ),
        h(
          NButton,
          { size: 'small', onClick: () => handleTogglePin(row) },
          { default: () => (row.isPinned ? t('admin.unpin') : t('admin.pin')) }
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

async function fetchAnnouncements() {
  loading.value = true
  try {
    const res = await adminGetAnnouncements({ page: filters.page, pageSize: filters.pageSize })
    announcements.value = res.data.data.list
    total.value = res.data.data.total
  } catch (e: any) {
    message.error(e.message || t('errors.networkError'))
  } finally {
    loading.value = false
  }
}

function handlePageChange(page: number) {
  filters.page = page
  fetchAnnouncements()
}

function openCreateModal() {
  editingId.value = null
  form.title = ''
  form.content = ''
  form.isPinned = false
  modalVisible.value = true
}

function openEditModal(row: any) {
  editingId.value = row.id
  form.title = row.title
  form.content = row.content
  form.isPinned = row.isPinned
  modalVisible.value = true
}

async function handleSave() {
  if (!form.title.trim() || !form.content.trim()) {
    message.warning(t('admin.titleContentRequired'))
    return
  }
  try {
    const data = { title: form.title, content: form.content, isPinned: form.isPinned }
    if (editingId.value) {
      await updateAnnouncement(editingId.value, data)
    } else {
      await createAnnouncement(data)
    }
    message.success(t('common.success'))
    modalVisible.value = false
    fetchAnnouncements()
  } catch (e: any) {
    message.error(e.message || t('errors.networkError'))
  }
}

async function handleTogglePin(row: any) {
  try {
    await updateAnnouncement(row.id, { isPinned: !row.isPinned })
    message.success(t('common.success'))
    fetchAnnouncements()
  } catch (e: any) {
    message.error(e.message || t('errors.networkError'))
  }
}

async function handleDelete(id: number) {
  try {
    await deleteAnnouncement(id)
    message.success(t('common.success'))
    fetchAnnouncements()
  } catch (e: any) {
    message.error(e.message || t('errors.networkError'))
  }
}

onMounted(() => {
  fetchAnnouncements()
})
</script>

<template>
  <div class="admin-announcement-manage">
    <div class="page-header flex-between">
      <h1 class="page-title">{{ t('admin.announcementManage') }}</h1>
      <NButton type="primary" @click="openCreateModal">
        {{ t('admin.createAnnouncement') }}
      </NButton>
    </div>

    <NCard class="table-card">
      <NDataTable
        :columns="columns"
        :data="announcements"
        :loading="loading"
        :bordered="false"
        :single-line="false"
        size="small"
        :row-key="(row: any) => row.id"
        :max-height="600"
      />
    </NCard>

    <div class="pagination-wrapper">
      <NPagination
        :page="filters.page"
        :page-size="filters.pageSize"
        :item-count="total"
        @update:page="handlePageChange"
      />
    </div>

    <!-- Create/Edit Modal -->
    <NModal
      v-model:show="modalVisible"
      preset="dialog"
      :title="editingId ? t('common.edit') : t('common.create')"
      positive-text="Confirm"
      negative-text="Cancel"
      style="width: 600px"
      @positive-click="handleSave"
    >
      <NForm>
        <NFormItem :label="t('nav.announcements')" required>
          <NInput v-model:value="form.title" :placeholder="t('admin.announcementTitlePlaceholder')" />
        </NFormItem>
        <NFormItem :label="t('teams.description')" required>
          <NInput
            v-model:value="form.content"
            type="textarea"
            :rows="8"
            :placeholder="t('admin.announcementContentPlaceholder')"
          />
        </NFormItem>
        <NFormItem :label="t('admin.pinned')">
          <NSwitch v-model:value="form.isPinned" />
        </NFormItem>
      </NForm>
    </NModal>
  </div>
</template>

<style scoped lang="scss">
.admin-announcement-manage {
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
