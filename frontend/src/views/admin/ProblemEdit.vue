<script setup lang="ts">
import { ref, reactive, computed, onMounted, h } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useI18n } from 'vue-i18n'
import {
  NForm, NFormItem, NInput, NSelect, NButton, NSpace, NCard, NInputNumber,
  NUpload, NDataTable, NTag, useMessage, type UploadFileInfo,
} from 'naive-ui'
import { getProblem } from '@/api/problem'
import { adminUpdateProblem, adminUploadTestData, adminDeleteTestData } from '@/api/problem'
import type { Problem, ProblemTestCase } from '@/types'

const props = defineProps<{ id: string }>()
const { t } = useI18n()
const router = useRouter()
const route = useRoute()
const message = useMessage()

const loading = ref(false)
const submitting = ref(false)
const problemId = computed(() => Number(props.id) || Number(route.params.id))

const form = reactive({
  title: '',
  type: 'programming' as string,
  difficulty: 'easy' as string,
  timeLimit: 1000,
  memoryLimit: 262144,
  description: '',
  inputFormat: '',
  outputFormat: '',
  hints: '',
  ctfCategory: '',
  ctfFlag: '',
  tags: [] as string[],
})

const testCases = ref<ProblemTestCase[]>([])

const typeOptions = computed(() => [
  { label: t('problems.programming'), value: 'programming' },
  { label: t('problems.algorithm'), value: 'algorithm' },
  { label: t('problems.ctf'), value: 'ctf' },
])

const difficultyOptions = computed(() => [
  { label: t('problems.easy'), value: 'easy' },
  { label: t('problems.medium'), value: 'medium' },
  { label: t('problems.hard'), value: 'hard' },
  { label: t('problems.expert'), value: 'expert' },
])

const ctfCategoryOptions = computed(() => [
  { label: 'Web', value: 'web' },
  { label: 'Crypto', value: 'crypto' },
  { label: 'Pwn', value: 'pwn' },
  { label: 'Reverse', value: 'reverse' },
  { label: 'Misc', value: 'misc' },
  { label: 'Forensics', value: 'forensics' },
  { label: 'Blockchain', value: 'blockchain' },
])

const tagOptions = computed(() => [
  { label: 'Array', value: 'Array' },
  { label: 'String', value: 'String' },
  { label: 'Hash Table', value: 'Hash Table' },
  { label: 'Tree', value: 'Tree' },
  { label: 'Graph', value: 'Graph' },
  { label: 'Dynamic Programming', value: 'Dynamic Programming' },
  { label: 'Greedy', value: 'Greedy' },
  { label: 'Binary Search', value: 'Binary Search' },
  { label: 'Math', value: 'Math' },
  { label: 'Sorting', value: 'Sorting' },
  { label: 'Stack', value: 'Stack' },
  { label: 'Queue', value: 'Queue' },
  { label: 'DFS', value: 'DFS' },
  { label: 'BFS', value: 'BFS' },
  { label: 'Bit Manipulation', value: 'Bit Manipulation' },
])

const isCtf = computed(() => form.type === 'ctf')

const testDataColumns = computed(() => [
  {
    title: '#',
    key: 'index',
    width: 60,
    render(_row: ProblemTestCase, index: number) {
      return h('span', {}, `${index + 1}`)
    },
  },
  {
    title: t('admin.inputFile'),
    key: 'input',
    ellipsis: { tooltip: true },
    render(row: ProblemTestCase) {
      return row.input || '-'
    },
  },
  {
    title: t('admin.outputFile'),
    key: 'output',
    ellipsis: { tooltip: true },
    render(row: ProblemTestCase) {
      return row.output || '-'
    },
  },
  {
    title: t('admin.isSample'),
    key: 'isSample',
    width: 90,
    render(row: ProblemTestCase) {
      return h(
        NTag,
        { type: row.isSample ? 'success' : 'default', size: 'small', bordered: false },
        { default: () => (row.isSample ? t('common.yes') : t('common.no')) }
      )
    },
  },
  {
    title: t('common.actions'),
    key: 'actions',
    width: 100,
    render(row: ProblemTestCase) {
      return h(
        NButton,
        {
          size: 'small',
          type: 'error',
          onClick: () => handleDeleteTestData(row.id),
        },
        { default: () => t('common.delete') }
      )
    },
  },
])

async function fetchProblem() {
  loading.value = true
  try {
    const res = await getProblem(problemId.value)
    const data = res.data.data
    form.title = data.title
    form.type = data.type
    form.difficulty = data.difficulty
    form.timeLimit = data.timeLimit
    form.memoryLimit = data.memoryLimit
    form.description = data.description
    form.inputFormat = data.inputFormat || ''
    form.outputFormat = data.outputFormat || ''
    form.tags = data.tags || []
    form.ctfCategory = (data as any).ctfCategory || ''
    form.ctfFlag = (data as any).ctfFlag || ''
    if (data.hints && data.hints.length > 0) {
      form.hints = data.hints.join('\n')
    }
    testCases.value = (data as any).testCases || []
  } catch (e: any) {
    message.error(e.message || t('errors.networkError'))
  } finally {
    loading.value = false
  }
}

async function handleSubmit() {
  if (!form.title.trim()) {
    message.warning(t('admin.titleRequired'))
    return
  }

  submitting.value = true
  try {
    const data: Record<string, any> = {
      title: form.title,
      type: form.type,
      difficulty: form.difficulty,
      timeLimit: form.timeLimit,
      memoryLimit: form.memoryLimit,
      description: form.description,
      inputFormat: form.inputFormat,
      outputFormat: form.outputFormat,
      tags: form.tags,
    }

    if (form.hints.trim()) {
      data.hints = form.hints.split('\n').filter(h => h.trim())
    }

    if (isCtf.value) {
      data.ctfCategory = form.ctfCategory
      data.ctfFlag = form.ctfFlag
    }

    await adminUpdateProblem(problemId.value, data as any)
    message.success(t('common.success'))
    router.push('/admin/problems')
  } catch (e: any) {
    message.error(e.message || t('errors.networkError'))
  } finally {
    submitting.value = false
  }
}

async function handleTestDataUpload({ file }: { file: UploadFileInfo }) {
  const formData = new FormData()
  formData.append('file', file.file as File)
  try {
    await adminUploadTestData(problemId.value, formData)
    message.success(t('common.success'))
    fetchProblem()
  } catch (e: any) {
    message.error(e.message || t('errors.networkError'))
  }
  return false
}

async function handleDeleteTestData(dataId: number) {
  try {
    await adminDeleteTestData(problemId.value, dataId)
    message.success(t('common.success'))
    fetchProblem()
  } catch (e: any) {
    message.error(e.message || t('errors.networkError'))
  }
}

onMounted(() => {
  fetchProblem()
})
</script>

<template>
  <div class="admin-problem-edit">
    <div class="page-header">
      <h1 class="page-title">{{ t('admin.editProblem') }} #{{ problemId }}</h1>
    </div>

    <NCard class="form-card" :loading="loading">
      <NForm label-placement="left" label-width="140">
        <NFormItem :label="t('problems.problemTitle')" required>
          <NInput
            v-model:value="form.title"
            :placeholder="t('admin.problemTitlePlaceholder')"
          />
        </NFormItem>

        <NFormItem :label="t('problems.type')">
          <NSelect
            v-model:value="form.type"
            :options="typeOptions"
            style="width: 300px"
          />
        </NFormItem>

        <NFormItem :label="t('problems.difficulty')">
          <NSelect
            v-model:value="form.difficulty"
            :options="difficultyOptions"
            style="width: 300px"
          />
        </NFormItem>

        <NFormItem :label="t('problems.timeLimit')">
          <NInputNumber
            v-model:value="form.timeLimit"
            :min="100"
            :max="30000"
            :step="100"
            style="width: 300px"
          >
            <template #suffix>ms</template>
          </NInputNumber>
        </NFormItem>

        <NFormItem :label="t('problems.memoryLimit')">
          <NInputNumber
            v-model:value="form.memoryLimit"
            :min="1024"
            :max="1048576"
            :step="1024"
            style="width: 300px"
          >
            <template #suffix>KB</template>
          </NInputNumber>
        </NFormItem>

        <NFormItem :label="t('problems.description')">
          <NInput
            v-model:value="form.description"
            type="textarea"
            :rows="12"
            :placeholder="t('admin.descriptionPlaceholder')"
          />
        </NFormItem>

        <NFormItem :label="t('problems.inputFormat')">
          <NInput
            v-model:value="form.inputFormat"
            type="textarea"
            :rows="4"
            :placeholder="t('admin.inputFormatPlaceholder')"
          />
        </NFormItem>

        <NFormItem :label="t('problems.outputFormat')">
          <NInput
            v-model:value="form.outputFormat"
            type="textarea"
            :rows="4"
            :placeholder="t('admin.outputFormatPlaceholder')"
          />
        </NFormItem>

        <NFormItem :label="t('problems.hints')">
          <NInput
            v-model:value="form.hints"
            type="textarea"
            :rows="4"
            :placeholder="t('admin.hintsPlaceholder')"
          />
        </NFormItem>

        <!-- CTF Fields -->
        <template v-if="isCtf">
          <NFormItem :label="t('admin.ctfCategory')">
            <NSelect
              v-model:value="form.ctfCategory"
              :options="ctfCategoryOptions"
              style="width: 300px"
              clearable
            />
          </NFormItem>
          <NFormItem :label="t('admin.ctfFlag')">
            <NInput
              v-model:value="form.ctfFlag"
              :placeholder="t('admin.ctfFlagPlaceholder')"
              style="width: 300px"
            />
          </NFormItem>
        </template>

        <NFormItem :label="t('problems.tags')">
          <NSelect
            v-model:value="form.tags"
            :options="tagOptions"
            multiple
            filterable
            tag
            style="width: 100%"
            :placeholder="t('admin.tagsPlaceholder')"
          />
        </NFormItem>

        <NFormItem label=" ">
          <NSpace>
            <NButton type="primary" :loading="submitting" @click="handleSubmit">
              {{ t('common.save') }}
            </NButton>
            <NButton @click="router.push('/admin/problems')">
              {{ t('common.cancel') }}
            </NButton>
          </NSpace>
        </NFormItem>
      </NForm>
    </NCard>

    <!-- Test Data Management -->
    <NCard class="section-card" :title="t('admin.testDataManage')">
      <div class="upload-section">
        <NUpload
          :custom-request="handleTestDataUpload"
          accept=".in,.out"
          :show-file-list="false"
          multiple
        >
          <NButton type="primary">{{ t('admin.uploadTestData') }}</NButton>
        </NUpload>
        <span class="upload-hint">{{ t('admin.uploadTestDataHint') }}</span>
      </div>

      <NDataTable
        :columns="testDataColumns"
        :data="testCases"
        :bordered="false"
        :single-line="false"
        size="small"
        :row-key="(row: ProblemTestCase) => row.id"
        class="mt-2"
      />
    </NCard>
  </div>
</template>

<style scoped lang="scss">
.admin-problem-edit {
  max-width: 1000px;
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

.form-card,
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
  margin-bottom: 16px;
}

.upload-hint {
  font-size: 13px;
  color: var(--color-text-secondary);
}
</style>
