<script setup lang="ts">
import { ref, reactive, onMounted, h } from 'vue'
import { useI18n } from 'vue-i18n'
import {
  NCard, NUpload, NDataTable, NButton, NSpace, NTag, NModal, NPagination, useMessage,
  type UploadFileInfo,
} from 'naive-ui'
import { importProblems } from '@/api/admin'

const { t } = useI18n()
const message = useMessage()

const loading = ref(false)
const importHistory = ref<any[]>([])
const total = ref(0)

const filters = reactive({
  page: 1,
  pageSize: 20,
})

const detailModal = ref(false)
const detailData = ref<any>(null)

type NTagType = 'default' | 'error' | 'info' | 'primary' | 'success' | 'warning'

const importStatusColor: Record<string, NTagType> = {
  success: 'success',
  failed: 'error',
  pending: 'warning',
  processing: 'info',
}

const historyColumns = computed(() => [
  {
    title: 'ID',
    key: 'id',
    width: 80,
  },
  {
    title: t('admin.sourcePlatform'),
    key: 'sourcePlatform',
    width: 120,
  },
  {
    title: t('admin.importFile'),
    key: 'fileName',
    ellipsis: { tooltip: true },
  },
  {
    title: t('admin.importedBy'),
    key: 'importedBy',
    width: 120,
    render(row: any) {
      return row.importedBy?.username || '-'
    },
  },
  {
    title: t('admin.importCount'),
    key: 'count',
    width: 80,
  },
  {
    title: t('common.status'),
    key: 'status',
    width: 100,
    render(row: any) {
      return h(
        NTag,
        { type: importStatusColor[row.status] || 'default', size: 'small', bordered: false },
        { default: () => row.status }
      )
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
    width: 80,
    render(row: any) {
      return h(
        NButton,
        { size: 'small', onClick: () => openDetail(row) },
        { default: () => t('admin.viewDetail') }
      )
    },
  },
])

async function handleUpload({ file }: { file: UploadFileInfo }) {
  loading.value = true
  try {
    await importProblems(file.file as File)
    message.success(t('admin.importStarted'))
  } catch (e: any) {
    message.error(e.message || t('errors.networkError'))
  } finally {
    loading.value = false
  }
  return false
}

function openDetail(row: any) {
  detailData.value = row
  detailModal.value = true
}

function handlePageChange(page: number) {
  filters.page = page
}

onMounted(() => {
  // Import history would be fetched from API
  // For now, showing empty state
})
</script>

<template>
  <div class="admin-import-manage">
    <div class="page-header">
      <h1 class="page-title">{{ t('admin.importManage') }}</h1>
    </div>

    <!-- Upload Section -->
    <NCard class="upload-card" :title="t('admin.uploadImportFile')">
      <div class="upload-section">
        <NUpload
          :custom-request="handleUpload"
          accept=".json"
          :show-file-list="false"
        >
          <NButton type="primary">{{ t('common.upload') }}</NButton>
        </NUpload>
        <span class="upload-hint">{{ t('admin.importFileHint') }}</span>
      </div>
    </NCard>

    <!-- Import History -->
    <NCard class="section-card" :title="t('admin.importHistory')">
      <NDataTable
        :columns="historyColumns"
        :data="importHistory"
        :loading="loading"
        :bordered="false"
        :single-line="false"
        size="small"
        :row-key="(row: any) => row.id"
        :max-height="500"
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

    <!-- Detail Modal -->
    <NModal
      v-model:show="detailModal"
      preset="card"
      :title="t('admin.importDetail')"
      style="width: 600px"
    >
      <div v-if="detailData" class="detail-content">
        <p><strong>{{ t('admin.sourcePlatform') }}:</strong> {{ detailData.sourcePlatform }}</p>
        <p><strong>{{ t('admin.importFile') }}:</strong> {{ detailData.fileName }}</p>
        <p><strong>{{ t('admin.importCount') }}:</strong> {{ detailData.count }}</p>
        <p><strong>{{ t('common.status') }}:</strong> {{ detailData.status }}</p>
        <p><strong>{{ t('teams.createdAt') }}:</strong> {{ new Date(detailData.createdAt).toLocaleString() }}</p>
      </div>
    </NModal>
  </div>
</template>

<style scoped lang="scss">
.admin-import-manage {
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

.upload-card,
.section-card {
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

.upload-section {
  display: flex;
  align-items: center;
  gap: 12px;
}

.upload-hint {
  font-size: 13px;
  color: var(--color-text-secondary);
}

.pagination-wrapper {
  display: flex;
  justify-content: center;
  padding: 16px 0;
}

.detail-content {
  p {
    margin-bottom: 8px;
    font-size: 14px;
    color: var(--color-text);
  }
}
</style>
