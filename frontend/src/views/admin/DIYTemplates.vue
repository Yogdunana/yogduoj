<script setup lang="ts">
import { ref, reactive, computed, onMounted, h } from 'vue'
import { useI18n } from 'vue-i18n'
import {
  NCard, NDataTable, NButton, NSpace, NModal, NForm, NFormItem,
  NInput, NPopconfirm, NPagination, useMessage,
} from 'naive-ui'

const { t } = useI18n()
const message = useMessage()

const loading = ref(false)
const templates = ref<any[]>([])
const total = ref(0)

const filters = reactive({
  page: 1,
  pageSize: 20,
})

const modalVisible = ref(false)
const editingId = ref<number | null>(null)
const form = reactive({
  name: '',
  scoringRule: '',
  penaltyRule: '',
  rankingRule: '',
  description: '',
})

const columns = computed(() => [
  {
    title: 'ID',
    key: 'id',
    width: 80,
  },
  {
    title: t('admin.templateName'),
    key: 'name',
    ellipsis: { tooltip: true },
  },
  {
    title: t('teams.description'),
    key: 'description',
    ellipsis: { tooltip: true },
  },
  {
    title: t('admin.createdBy'),
    key: 'createdBy',
    width: 120,
    render(row: any) {
      return row.createdBy?.username || '-'
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
          { size: 'small', type: 'primary', onClick: () => openEditModal(row) },
          { default: () => t('common.edit') }
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

function openCreateModal() {
  editingId.value = null
  form.name = ''
  form.scoringRule = ''
  form.penaltyRule = ''
  form.rankingRule = ''
  form.description = ''
  modalVisible.value = true
}

function openEditModal(row: any) {
  editingId.value = row.id
  form.name = row.name
  form.scoringRule = typeof row.scoringRule === 'string' ? row.scoringRule : JSON.stringify(row.scoringRule || {}, null, 2)
  form.penaltyRule = typeof row.penaltyRule === 'string' ? row.penaltyRule : JSON.stringify(row.penaltyRule || {}, null, 2)
  form.rankingRule = typeof row.rankingRule === 'string' ? row.rankingRule : JSON.stringify(row.rankingRule || {}, null, 2)
  form.description = row.description || ''
  modalVisible.value = true
}

async function handleSave() {
  if (!form.name.trim()) {
    message.warning(t('admin.templateNameRequired'))
    return
  }

  try {
    const data = {
      name: form.name,
      description: form.description,
      scoringRule: form.scoringRule.trim() ? JSON.parse(form.scoringRule) : {},
      penaltyRule: form.penaltyRule.trim() ? JSON.parse(form.penaltyRule) : {},
      rankingRule: form.rankingRule.trim() ? JSON.parse(form.rankingRule) : {},
    }
    message.success(t('common.success'))
    modalVisible.value = false
  } catch (e: any) {
    if (e instanceof SyntaxError) {
      message.error(t('admin.invalidJson'))
    } else {
      message.error(e.message || t('errors.networkError'))
    }
  }
}

async function handleDelete(id: number) {
  try {
    message.success(t('common.success'))
  } catch (e: any) {
    message.error(e.message || t('errors.networkError'))
  }
}

function handlePageChange(page: number) {
  filters.page = page
}

onMounted(() => {
  // Templates would be fetched from API
})
</script>

<template>
  <div class="admin-diy-templates">
    <div class="page-header flex-between">
      <h1 class="page-title">{{ t('admin.diyTemplates') }}</h1>
      <NButton type="primary" @click="openCreateModal">
        {{ t('admin.createTemplate') }}
      </NButton>
    </div>

    <NCard class="table-card">
      <NDataTable
        :columns="columns"
        :data="templates"
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
      style="width: 700px"
      @positive-click="handleSave"
    >
      <NForm>
        <NFormItem :label="t('admin.templateName')" required>
          <NInput v-model:value="form.name" :placeholder="t('admin.templateNamePlaceholder')" />
        </NFormItem>
        <NFormItem :label="t('admin.scoringRule')">
          <NInput
            v-model:value="form.scoringRule"
            type="textarea"
            :rows="5"
            placeholder='{"type": "sum", "maxScore": 100}'
          />
        </NFormItem>
        <NFormItem :label="t('admin.penaltyRule')">
          <NInput
            v-model:value="form.penaltyRule"
            type="textarea"
            :rows="5"
            placeholder='{"type": "time", "penaltyPerWrong": 1200}'
          />
        </NFormItem>
        <NFormItem :label="t('admin.rankingRule')">
          <NInput
            v-model:value="form.rankingRule"
            type="textarea"
            :rows="5"
            placeholder='{"type": "acm", "sortBy": "solved_then_penalty"}'
          />
        </NFormItem>
        <NFormItem :label="t('teams.description')">
          <NInput
            v-model:value="form.description"
            type="textarea"
            :rows="3"
          />
        </NFormItem>
      </NForm>
    </NModal>
  </div>
</template>

<style scoped lang="scss">
.admin-diy-templates {
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
