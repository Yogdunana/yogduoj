<script setup lang="ts">
import { ref, computed, onMounted, watch } from 'vue'
import { useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { NCard, NButton, NTag, NSpin, NDescriptions, NDescriptionsItem, useMessage } from 'naive-ui'
import { getSubmission, getSubmissionCode } from '@/api/submission'
import CodeEditor from '@/components/common/CodeEditor.vue'
import type { Submission, SubmissionStatus } from '@/types'

const props = defineProps<{ id: string }>()
const router = useRouter()
const { t } = useI18n()
const message = useMessage()

const loading = ref(true)
const submission = ref<Submission | null>(null)
const code = ref('')
const codeLoading = ref(false)

type NTagType = 'default' | 'error' | 'info' | 'primary' | 'success' | 'warning'

const statusColorMap: Record<string, NTagType> = {
  pending: 'default',
  judging: 'info',
  accepted: 'success',
  wrong_answer: 'error',
  time_limit_exceeded: 'warning',
  memory_limit_exceeded: 'warning',
  runtime_error: 'error',
  compilation_error: 'error',
  presentation_error: 'warning',
  system_error: 'error',
}

function getStatusLabel(status: SubmissionStatus): string {
  const map: Record<string, string> = {
    pending: t('submissions.pending'),
    judging: t('submissions.judging'),
    accepted: t('submissions.accepted'),
    wrong_answer: t('submissions.wrongAnswer'),
    time_limit_exceeded: t('submissions.timeLimitExceeded'),
    memory_limit_exceeded: t('submissions.memoryLimitExceeded'),
    runtime_error: t('submissions.runtimeError'),
    compilation_error: t('submissions.compilationError'),
    presentation_error: t('submissions.presentationError'),
    system_error: t('submissions.systemError'),
  }
  return map[status] || status
}

function getTestCaseStatusLabel(status: SubmissionStatus): string {
  return getStatusLabel(status)
}

const languageMap: Record<string, string> = {
  c: 'c',
  cpp: 'cpp',
  java: 'java',
  python3: 'python',
  python: 'python',
}

const mappedLanguage = computed(() => {
  if (!submission.value) return 'cpp'
  return languageMap[submission.value.language] || submission.value.language
})

const hasJudgeDetail = computed(() => {
  return submission.value?.judgeDetail && submission.value.judgeDetail.length > 0
})

const hasError = computed(() => {
  const status = submission.value?.status
  return status === 'compilation_error' || status === 'runtime_error' || status === 'system_error'
})

function formatTime(ms: number): string {
  return `${ms} ms`
}

function formatMemory(kb: number): string {
  if (kb >= 1024) return `${(kb / 1024).toFixed(1)} MB`
  return `${kb} KB`
}

function formatDateTime(dateStr: string): string {
  return new Date(dateStr).toLocaleString()
}

async function fetchSubmission() {
  loading.value = true
  try {
    const res = await getSubmission(Number(props.id))
    submission.value = res.data.data
    // If code is included in the response, use it
    if (submission.value.code) {
      code.value = submission.value.code
    } else {
      // Fetch code separately
      await fetchCode()
    }
  } catch (e: any) {
    message.error(e.message || t('errors.networkError'))
  } finally {
    loading.value = false
  }
}

async function fetchCode() {
  codeLoading.value = true
  try {
    const res = await getSubmissionCode(Number(props.id))
    code.value = res.data.data.code
  } catch {
    code.value = '// Code not available'
  } finally {
    codeLoading.value = false
  }
}

function goBack() {
  router.push({ name: 'SubmissionList' })
}

function goToProblem() {
  if (submission.value) {
    router.push({ name: 'ProblemDetail', params: { id: submission.value.problemId } })
  }
}

onMounted(() => {
  fetchSubmission()
})

watch(() => props.id, () => {
  fetchSubmission()
})
</script>

<template>
  <div class="submission-detail-page">
    <NSpin :show="loading">
      <template v-if="submission">
        <!-- Back Button -->
        <div class="page-actions">
          <NButton @click="goBack" size="small">
            {{ t('submissions.backToList') }}
          </NButton>
        </div>

        <!-- Submission Info Card -->
        <NCard :title="t('submissions.submissionInfo')" size="small" class="info-card">
          <NDescriptions bordered :column="2" label-placement="left" size="small">
            <NDescriptionsItem :label="t('submissions.submissionId')">
              #{{ submission.id }}
            </NDescriptionsItem>
            <NDescriptionsItem :label="t('submissions.problem')">
              <a
                class="problem-link"
                @click.prevent="goToProblem"
              >
                {{ submission.problem?.title || `Problem #${submission.problemId}` }}
              </a>
            </NDescriptionsItem>
            <NDescriptionsItem :label="t('submissions.status')">
              <NTag :type="statusColorMap[submission.status] || 'default'" size="small" :bordered="false">
                {{ getStatusLabel(submission.status) }}
              </NTag>
            </NDescriptionsItem>
            <NDescriptionsItem :label="t('submissions.language')">
              {{ submission.language }}
            </NDescriptionsItem>
            <NDescriptionsItem :label="t('submissions.timeUsed')">
              {{ submission.timeUsed != null ? formatTime(submission.timeUsed) : '-' }}
            </NDescriptionsItem>
            <NDescriptionsItem :label="t('submissions.memoryUsed')">
              {{ submission.memoryUsed != null ? formatMemory(submission.memoryUsed) : '-' }}
            </NDescriptionsItem>
            <NDescriptionsItem :label="t('submissions.submitTime')">
              {{ formatDateTime(submission.createdAt) }}
            </NDescriptionsItem>
            <NDescriptionsItem v-if="submission.score != null" :label="t('submissions.score')">
              <span class="score-value">
                {{ submission.score }}{{ submission.totalScore != null ? ` / ${submission.totalScore}` : '' }}
              </span>
            </NDescriptionsItem>
          </NDescriptions>
        </NCard>

        <!-- Error Message Display -->
        <NCard v-if="hasError && submission.errorMessage" :title="t('submissions.errorMessage')" size="small" class="error-card">
          <pre class="error-message">{{ submission.errorMessage }}</pre>
        </NCard>

        <!-- Judge Detail -->
        <NCard v-if="hasJudgeDetail" :title="t('submissions.judgeDetail')" size="small" class="judge-card">
          <div class="judge-detail-list">
            <div
              v-for="testCase in submission.judgeDetail"
              :key="testCase.id"
              class="test-case-item"
            >
              <div class="test-case-header">
                <span class="test-case-id">
                  {{ t('submissions.testCase') }} #{{ testCase.id }}
                </span>
                <NTag :type="statusColorMap[testCase.status] || 'default'" size="small" :bordered="false">
                  {{ getTestCaseStatusLabel(testCase.status) }}
                </NTag>
              </div>
              <div class="test-case-info">
                <span v-if="testCase.timeUsed != null" class="info-item">
                  {{ t('submissions.timeUsed') }}: {{ formatTime(testCase.timeUsed) }}
                </span>
                <span v-if="testCase.memoryUsed != null" class="info-item">
                  {{ t('submissions.memoryUsed') }}: {{ formatMemory(testCase.memoryUsed) }}
                </span>
                <span v-if="testCase.score != null" class="info-item score">
                  {{ t('submissions.score') }}: {{ testCase.score }}
                </span>
              </div>
              <div v-if="testCase.errorMessage" class="test-case-error">
                {{ testCase.errorMessage }}
              </div>
            </div>
          </div>
        </NCard>

        <!-- Source Code -->
        <NCard :title="t('submissions.sourceCode')" size="small" class="code-card">
          <NSpin :show="codeLoading">
            <CodeEditor
              :model-value="code"
              :language="mappedLanguage"
              :readonly="true"
              :height="'500px'"
            />
          </NSpin>
        </NCard>
      </template>
    </NSpin>
  </div>
</template>

<style scoped lang="scss">
.submission-detail-page {
  padding: 24px 20px;
  max-width: 1200px;
  margin: 0 auto;
}

.page-actions {
  margin-bottom: 16px;
}

.info-card,
.error-card,
.judge-card,
.code-card {
  background-color: var(--color-bg-card);
  border: 1px solid rgba(255, 255, 255, 0.06);
  margin-bottom: 16px;
}

.problem-link {
  color: var(--color-primary);
  text-decoration: none;
  cursor: pointer;
  font-weight: 500;

  &:hover {
    color: var(--color-primary-hover);
    text-decoration: underline;
  }
}

.score-value {
  font-weight: 600;
  color: var(--color-primary);
}

.error-card {
  border-color: rgba(255, 82, 82, 0.3);
}

.error-message {
  background-color: rgba(255, 82, 82, 0.08);
  border-radius: 6px;
  padding: 16px;
  font-family: 'Fira Code', 'Consolas', monospace;
  font-size: 13px;
  white-space: pre-wrap;
  word-break: break-all;
  color: var(--color-error);
  overflow-x: auto;
  line-height: 1.6;
}

.judge-detail-list {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.test-case-item {
  padding: 12px;
  border-radius: 6px;
  background-color: rgba(0, 0, 0, 0.15);
  border: 1px solid rgba(255, 255, 255, 0.04);
}

.test-case-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 8px;
}

.test-case-id {
  font-size: 13px;
  font-weight: 600;
  color: var(--color-text);
}

.test-case-info {
  display: flex;
  gap: 16px;
  flex-wrap: wrap;
}

.info-item {
  font-size: 12px;
  color: var(--color-text-secondary);

  &.score {
    font-weight: 600;
    color: var(--color-primary);
  }
}

.test-case-error {
  margin-top: 8px;
  padding: 8px;
  border-radius: 4px;
  background-color: rgba(255, 82, 82, 0.08);
  font-size: 12px;
  color: var(--color-error);
  font-family: 'Fira Code', 'Consolas', monospace;
  white-space: pre-wrap;
  word-break: break-all;
}

@media (max-width: 768px) {
  .test-case-info {
    flex-direction: column;
    gap: 4px;
  }
}
</style>
