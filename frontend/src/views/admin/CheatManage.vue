<script setup lang="ts">
import { ref, reactive, computed, onMounted, h } from 'vue'
import { useI18n } from 'vue-i18n'
import {
  NDataTable, NInput, NSelect, NButton, NSpace, NTag, NModal, NForm, NFormItem,
  NInputNumber, NPagination, useMessage,
} from 'naive-ui'

const { t } = useI18n()
const message = useMessage()

const loading = ref(false)
const cheatRecords = ref<any[]>([])
const total = ref(0)

const filters = reactive({
  contestId: '' as string,
  cheatType: '' as string,
  reviewStatus: '' as string,
  page: 1,
  pageSize: 20,
})

const reviewModal = ref(false)
const reviewRecord = ref<any>(null)
const reviewForm = reactive({
  status: 'confirmed' as string,
  penalty: 'warning' as string,
  note: '',
})

const cheatTypeOptions = computed(() => [
  { label: t('admin.allCheatTypes'), value: '' },
  { label: t('admin.cheatCodeSimilarity'), value: 'code_similarity' },
  { label: t('admin.cheatSubmissionTime'), value: 'submission_time' },
  { label: t('admin.cheatSameIP'), value: 'same_ip' },
])

const reviewStatusOptions = computed(() => [
  { label: t('admin.allReviewStatuses'), value: '' },
  { label: t('admin.reviewPending'), value: 'pending' },
  { label: t('admin.reviewConfirmed'), value: 'confirmed' },
  { label: t('admin.reviewDismissed'), value: 'dismissed' },
])

type NTagType = 'default' | 'error' | 'info' | 'primary' | 'success' | 'warning'

const reviewStatusColor: Record<string, NTagType> = {
  pending: 'warning',
  confirmed: 'error',
  dismissed: 'success',
}

const columns = computed(() => [
  {
    title: t('auth.username'),
    key: 'username',
    width: 120,
  },
  {
    title: t('contests.title'),
    key: 'contestTitle',
    width: 150,
    ellipsis: { tooltip: true },
  },
  {
    title: t('admin.cheatType'),
    key: 'cheatType',
    width: 140,
  },
  {
    title: t('admin.similarity'),
    key: 'similarity',
    width: 100,
    render(row: any) {
      return row.similarity != null ? `${(row.similarity * 100).toFixed(1)}%` : '-'
    },
  },
  {
    title: t('admin.cheatDetail'),
    key: 'detail',
    ellipsis: { tooltip: true },
    render(row: any) {
      return row.detail || '-'
    },
  },
  {
    title: t('admin.reviewStatus'),
    key: 'reviewStatus',
    width: 100,
    render(row: any) {
      const status = row.reviewStatus || 'pending'
      return h(
        NTag,
        { type: reviewStatusColor[status] || 'default', size: 'small', bordered: false },
        { default: () => t(`admin.review${status.charAt(0).toUpperCase() + status.slice(1)}`) }
      )
    },
  },
  {
    title: t('common.actions'),
    key: 'actions',
    width: 160,
    render(row: any) {
      return h(NSpace, { size: 'small' }, () => [
        h(
          NButton,
          { size: 'small', type: 'primary', onClick: () => openReview(row) },
          { default: () => t('admin.review') }
        ),
        h(
          NButton,
          { size: 'small', onClick: () => openReview(row) },
          { default: () => t('admin.viewDetail') }
        ),
      ])
    },
  },
])

function handleContestFilter(val: string) {
  filters.contestId = val
  filters.page = 1
}

function handleCheatTypeChange(val: string) {
  filters.cheatType = val
  filters.page = 1
}

function handleReviewStatusChange(val: string) {
  filters.reviewStatus = val
  filters.page = 1
}

function handlePageChange(page: number) {
  filters.page = page
}

function openReview(record: any) {
  reviewRecord.value = record
  reviewForm.status = 'confirmed'
  reviewForm.penalty = 'warning'
  reviewForm.note = ''
  reviewModal.value = true
}

async function handleReviewSubmit() {
  message.success(t('common.success'))
  reviewModal.value = false
}

onMounted(() => {
  // Cheat records would be fetched from API
})
</script>

<template>
  <div class="admin-cheat-manage">
    <div class="page-header">
      <h1 class="page-title">{{ t('admin.cheatManage') }}</h1>
    </div>

    <!-- Filters -->
    <div class="toolbar">
      <NInput
        :placeholder="t('admin.contestIdPlaceholder')"
        :value="filters.contestId"
        clearable
        @update:value="(v: string) => { filters.contestId = v }"
        style="width: 140px"
      />
      <NSelect
        :value="filters.cheatType"
        :options="cheatTypeOptions"
        size="small"
        style="width: 180px"
        @update:value="handleCheatTypeChange"
      />
      <NSelect
        :value="filters.reviewStatus"
        :options="reviewStatusOptions"
        size="small"
        style="width: 160px"
        @update:value="handleReviewStatusChange"
      />
    </div>

    <!-- Cheat Records Table -->
    <NCard class="table-card">
      <NDataTable
        :columns="columns"
        :data="cheatRecords"
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

    <!-- Review Modal -->
    <NModal
      v-model:show="reviewModal"
      preset="dialog"
      :title="t('admin.reviewCheatRecord')"
      positive-text="Confirm"
      negative-text="Cancel"
      style="width: 600px"
      @positive-click="handleReviewSubmit"
    >
      <div v-if="reviewRecord" class="review-detail">
        <p><strong>{{ t('auth.username') }}:</strong> {{ reviewRecord.username }}</p>
        <p><strong>{{ t('admin.cheatType') }}:</strong> {{ reviewRecord.cheatType }}</p>
        <p><strong>{{ t('admin.similarity') }}:</strong> {{ reviewRecord.similarity != null ? `${(reviewRecord.similarity * 100).toFixed(1)}%` : '-' }}</p>
        <p><strong>{{ t('admin.cheatDetail') }}:</strong> {{ reviewRecord.detail || '-' }}</p>
      </div>
      <NForm class="mt-2">
        <NFormItem :label="t('admin.reviewStatus')">
          <NSelect
            v-model:value="reviewForm.status"
            :options="[
              { label: t('admin.reviewConfirmed'), value: 'confirmed' },
              { label: t('admin.reviewDismissed'), value: 'dismissed' },
            ]"
            style="width: 100%"
          />
        </NFormItem>
        <NFormItem :label="t('admin.penalty')">
          <NSelect
            v-model:value="reviewForm.penalty"
            :options="[
              { label: t('admin.penaltyWarning'), value: 'warning' },
              { label: t('admin.penaltyScore'), value: 'score_penalty' },
              { label: t('admin.penaltyBan'), value: 'ban' },
            ]"
            style="width: 100%"
          />
        </NFormItem>
        <NFormItem :label="t('admin.reviewNote')">
          <NInput
            v-model:value="reviewForm.note"
            type="textarea"
            :rows="3"
            :placeholder="t('admin.reviewNotePlaceholder')"
          />
        </NFormItem>
      </NForm>
    </NModal>
  </div>
</template>

<style scoped lang="scss">
.admin-cheat-manage {
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

.review-detail {
  p {
    margin-bottom: 8px;
    font-size: 14px;
    color: var(--color-text);
  }
}
</style>
