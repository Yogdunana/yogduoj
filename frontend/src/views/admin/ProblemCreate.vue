<script setup lang="ts">
import { ref, reactive, computed } from 'vue'
import { useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import {
  NForm, NFormItem, NInput, NSelect, NButton, NSpace, NCard, NInputNumber, useMessage,
} from 'naive-ui'
import { adminCreateProblem } from '@/api/admin'

const { t } = useI18n()
const router = useRouter()
const message = useMessage()

const submitting = ref(false)

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

    const res = await adminCreateProblem(data as any)
    const problemId = res.data.data?.id
    message.success(t('common.success'))
    if (problemId) {
      router.push({ name: 'AdminProblemEdit', params: { id: problemId } })
    } else {
      router.push('/admin/problems')
    }
  } catch (e: any) {
    message.error(e.message || t('errors.networkError'))
  } finally {
    submitting.value = false
  }
}
</script>

<template>
  <div class="admin-problem-create">
    <div class="page-header">
      <h1 class="page-title">{{ t('admin.createProblem') }}</h1>
    </div>

    <NCard class="form-card">
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
              {{ t('common.create') }}
            </NButton>
            <NButton @click="router.back()">
              {{ t('common.cancel') }}
            </NButton>
          </NSpace>
        </NFormItem>
      </NForm>
    </NCard>
  </div>
</template>

<style scoped lang="scss">
.admin-problem-create {
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

.form-card {
  background-color: var(--color-bg-card);
  border: 1px solid rgba(255, 255, 255, 0.06);
}
</style>
