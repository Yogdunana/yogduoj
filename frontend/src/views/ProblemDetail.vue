<script setup lang="ts">
import { ref, computed, onMounted, onBeforeUnmount, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import {
  NGrid, NGridItem, NCard, NButton, NSelect, NSpin, NTag,
  NCollapse, NCollapseItem, NUpload, NSpace, NInputGroup, NInput,
  useMessage, useDialog,
} from 'naive-ui'
import { getProblem } from '@/api/problem'
import { createSubmission, listSubmissions } from '@/api/submission'
import CodeEditor from '@/components/common/CodeEditor.vue'
import type { Problem, Submission, SubmissionStatus } from '@/types'

const props = defineProps<{ id: string }>()
const route = useRoute()
const router = useRouter()
const { t } = useI18n()
const message = useMessage()
const dialog = useDialog()

const loading = ref(true)
const problem = ref<Problem | null>(null)
const code = ref('')
const language = ref('cpp')
const submitting = ref(false)
const judgeStatus = ref<SubmissionStatus | null>(null)
const submissions = ref<Submission[]>([])
const flagValue = ref('')
const uploadedFile = ref<File | null>(null)
const autoSaveTimer = ref<ReturnType<typeof setTimeout> | null>(null)
const pollTimer = ref<ReturnType<typeof setInterval> | null>(null)
const lastSubmissionId = ref<number | null>(null)

const isCtf = computed(() => problem.value?.type === 'ctf')

const languageOptions = computed(() => [
  { label: 'C', value: 'c' },
  { label: 'C++', value: 'cpp' },
  { label: 'Java', value: 'java' },
  { label: 'Python3', value: 'python3' },
])

const difficultyColorMap: Record<string, NTagType> = {
  easy: 'success',
  medium: 'warning',
  hard: 'error',
  expert: 'default',
}

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

function getLocalStorageKey(): string {
  return `problem_${props.id}_${language.value}`
}

function loadCodeFromStorage() {
  const key = getLocalStorageKey()
  const stored = localStorage.getItem(key)
  if (stored) {
    code.value = stored
  } else {
    code.value = getDefaultCode(language.value)
  }
}

function getDefaultCode(lang: string): string {
  const defaults: Record<string, string> = {
    c: '#include <stdio.h>\n\nint main() {\n    // Write your code here\n    return 0;\n}\n',
    cpp: '#include <iostream>\nusing namespace std;\n\nint main() {\n    // Write your code here\n    return 0;\n}\n',
    java: 'public class Main {\n    public static void main(String[] args) {\n        // Write your code here\n    }\n}\n',
    python3: '# Write your code here\nif __name__ == "__main__":\n    pass\n',
  }
  return defaults[lang] || ''
}

function saveCodeToStorage() {
  const key = getLocalStorageKey()
  localStorage.setItem(key, code.value)
}

function handleCodeChange(val: string) {
  code.value = val
  if (autoSaveTimer.value) clearTimeout(autoSaveTimer.value)
  autoSaveTimer.value = setTimeout(() => {
    saveCodeToStorage()
  }, 2000)
}

function handleLanguageChange(val: string) {
  language.value = val
  loadCodeFromStorage()
}

async function handleSubmit() {
  if (isCtf.value) {
    await submitFlag()
    return
  }

  dialog.warning({
    title: t('common.confirm'),
    content: t('problems.submitConfirm'),
    positiveText: t('common.confirm'),
    negativeText: t('common.cancel'),
    onPositiveClick: async () => {
      await doSubmitCode()
    },
  })
}

async function submitFlag() {
  if (!flagValue.value.trim()) {
    message.warning(t('problems.flagPlaceholder'))
    return
  }
  submitting.value = true
  try {
    await createSubmission({
      problemId: Number(props.id),
      language: 'flag',
      code: flagValue.value.trim(),
    })
    message.success(t('problems.submitSuccess'))
    flagValue.value = ''
    fetchSubmissions()
  } catch (e: any) {
    message.error(e.message || t('problems.submitError'))
  } finally {
    submitting.value = false
  }
}

async function doSubmitCode() {
  if (!code.value.trim()) {
    message.warning(t('problems.code'))
    return
  }
  submitting.value = true
  judgeStatus.value = 'pending'
  try {
    saveCodeToStorage()
    const res = await createSubmission({
      problemId: Number(props.id),
      language: language.value,
      code: code.value,
    })
    const submission = res.data.data
    lastSubmissionId.value = submission.id
    judgeStatus.value = submission.status
    message.success(t('problems.submitSuccess'))
    startPolling(submission.id)
    fetchSubmissions()
  } catch (e: any) {
    message.error(e.message || t('problems.submitError'))
    judgeStatus.value = null
  } finally {
    submitting.value = false
  }
}

function startPolling(submissionId: number) {
  if (pollTimer.value) clearInterval(pollTimer.value)
  pollTimer.value = setInterval(async () => {
    try {
      const { getSubmission } = await import('@/api/submission')
      const res = await getSubmission(submissionId)
      const sub = res.data.data
      judgeStatus.value = sub.status
      if (sub.status !== 'pending' && sub.status !== 'judging') {
        if (pollTimer.value) clearInterval(pollTimer.value)
        pollTimer.value = null
        fetchSubmissions()
      }
    } catch {
      if (pollTimer.value) clearInterval(pollTimer.value)
      pollTimer.value = null
    }
  }, 2000)
}

async function fetchProblem() {
  loading.value = true
  try {
    const res = await getProblem(Number(props.id))
    problem.value = res.data.data
    loadCodeFromStorage()
  } catch (e: any) {
    message.error(e.message || t('errors.networkError'))
  } finally {
    loading.value = false
  }
}

async function fetchSubmissions() {
  try {
    const res = await listSubmissions({
      problemId: Number(props.id),
      page: 1,
      pageSize: 10,
    })
    submissions.value = res.data.data.list
  } catch {
    // silent
  }
}

function copyToClipboard(text: string) {
  navigator.clipboard.writeText(text).then(() => {
    message.success(t('problems.copied'))
  }).catch(() => {
    message.error(t('errors.unknownError'))
  })
}

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

function goToSubmission(id: number) {
  router.push({ name: 'SubmissionDetail', params: { id } })
}

onMounted(() => {
  fetchProblem()
  fetchSubmissions()
})

onBeforeUnmount(() => {
  if (pollTimer.value) clearInterval(pollTimer.value)
  if (autoSaveTimer.value) clearTimeout(autoSaveTimer.value)
})

watch(() => props.id, () => {
  if (pollTimer.value) clearInterval(pollTimer.value)
  judgeStatus.value = null
  lastSubmissionId.value = null
  fetchProblem()
  fetchSubmissions()
})
</script>

<template>
  <div class="problem-detail-page">
    <NSpin :show="loading">
      <template v-if="problem">
        <!-- Problem Header -->
        <div class="problem-header">
          <div class="problem-title-row">
            <h1 class="problem-title">{{ problem.title }}</h1>
            <NSpace>
              <NTag :type="difficultyColorMap[problem.difficulty] || 'default'" size="small" bordered>
                {{ t(`problems.${problem.difficulty}`) }}
              </NTag>
              <NTag type="info" size="small" bordered>
                {{ t(`problems.${problem.type}`) }}
              </NTag>
              <NTag v-for="tag in problem.tags" :key="tag" size="small" :bordered="false" type="default">
                {{ tag }}
              </NTag>
            </NSpace>
          </div>
          <div class="problem-meta">
            <span class="meta-item">
              {{ t('problems.timeLimit') }}: {{ problem.timeLimit }} ms
            </span>
            <span class="meta-item">
              {{ t('problems.memoryLimit') }}: {{ problem.memoryLimit }} MB
            </span>
            <span class="meta-item">
              {{ t('problems.acceptanceRate') }}:
              {{ problem.totalSubmit > 0 ? ((problem.acceptedCount / problem.totalSubmit) * 100).toFixed(1) + '%' : '0%' }}
            </span>
            <span class="meta-item">
              {{ t('problems.submitCount') }}: {{ problem.totalSubmit }}
            </span>
          </div>
        </div>

        <!-- Main Content: Problem + Editor -->
        <NGrid :x-gap="20" :cols="24" responsive="screen" item-responsive>
          <!-- Left Panel: Problem Description -->
          <NGridItem :span="24" :md="14" :lg="14">
            <div class="problem-content">
              <!-- Description -->
              <NCard :title="t('problems.description')" size="small" class="content-card">
                <div class="markdown-body" v-html="problem.description"></div>
              </NCard>

              <!-- Input/Output Format -->
              <NGrid :x-gap="16" :cols="2" responsive="screen" item-responsive>
                <NGridItem :span="2" :md="1">
                  <NCard :title="t('problems.inputFormat')" size="small" class="content-card">
                    <div class="markdown-body" v-if="problem.inputFormat" v-html="problem.inputFormat"></div>
                    <span v-else class="text-secondary">-</span>
                  </NCard>
                </NGridItem>
                <NGridItem :span="2" :md="1">
                  <NCard :title="t('problems.outputFormat')" size="small" class="content-card">
                    <div class="markdown-body" v-if="problem.outputFormat" v-html="problem.outputFormat"></div>
                    <span v-else class="text-secondary">-</span>
                  </NCard>
                </NGridItem>
              </NGrid>

              <!-- Samples -->
              <NCard :title="t('problems.samples')" size="small" class="content-card">
                <div v-if="problem.samples && problem.samples.length > 0" class="samples-list">
                  <div v-for="(sample, idx) in problem.samples" :key="idx" class="sample-item">
                    <div class="sample-header">
                      <span class="sample-label">{{ t('problems.sampleInput') }} {{ idx + 1 }}</span>
                      <NButton text size="tiny" @click="copyToClipboard(sample.input)">
                        {{ t('problems.copySample') }}
                      </NButton>
                    </div>
                    <pre class="sample-code">{{ sample.input }}</pre>

                    <div class="sample-header">
                      <span class="sample-label">{{ t('problems.sampleOutput') }} {{ idx + 1 }}</span>
                      <NButton text size="tiny" @click="copyToClipboard(sample.output)">
                        {{ t('problems.copySample') }}
                      </NButton>
                    </div>
                    <pre class="sample-code">{{ sample.output }}</pre>
                  </div>
                </div>
                <span v-else class="text-secondary">{{ t('problems.noSamples') }}</span>
              </NCard>

              <!-- Hints -->
              <NCard v-if="problem.hints && problem.hints.length > 0" size="small" class="content-card">
                <NCollapse>
                  <NCollapseItem :title="t('problems.hints')" :name="'hints'">
                    <div v-for="(hint, idx) in problem.hints" :key="idx" class="hint-item">
                      <p>{{ hint }}</p>
                    </div>
                  </NCollapseItem>
                </NCollapse>
              </NCard>

              <!-- Attachments -->
              <NCard v-if="problem.attachments && problem.attachments.length > 0" :title="t('problems.attachments')" size="small" class="content-card">
                <div class="attachments-list">
                  <a
                    v-for="(url, idx) in problem.attachments"
                    :key="idx"
                    :href="url"
                    target="_blank"
                    class="attachment-link"
                  >
                    {{ url.split('/').pop() || t('common.download') }} {{ idx + 1 }}
                  </a>
                </div>
              </NCard>
            </div>
          </NGridItem>

          <!-- Right Panel: Code Editor / CTF -->
          <NGridItem :span="24" :md="10" :lg="10">
            <div class="editor-panel">
              <NCard class="editor-card">
                <!-- Language Selector -->
                <div v-if="!isCtf" class="editor-toolbar">
                  <NSelect
                    :value="language"
                    :options="languageOptions"
                    size="small"
                    style="width: 140px"
                    @update:value="handleLanguageChange"
                  />
                  <span class="auto-save-hint">{{ t('problems.autoSaved') }}</span>
                </div>

                <!-- Code Editor -->
                <div v-if="!isCtf" class="editor-container">
                  <CodeEditor
                    :model-value="code"
                    :language="language"
                    :height="'400px'"
                    @update:model-value="handleCodeChange"
                  />
                </div>

                <!-- CTF Flag Input -->
                <div v-else class="ctf-panel">
                  <div class="flag-section">
                    <label class="flag-label">{{ t('problems.flagInput') }}</label>
                    <NInput
                      v-model:value="flagValue"
                      :placeholder="t('problems.flagPlaceholder')"
                      type="text"
                      size="large"
                    />
                  </div>

                  <div class="upload-section">
                    <label class="flag-label">{{ t('problems.fileUpload') }}</label>
                    <NUpload
                      :max="1"
                      :default-upload="false"
                      accept="*/*"
                      @before-upload="(file: any) => { uploadedFile = file.file; return false }"
                    >
                      <NButton>{{ t('problems.fileUpload') }}</NButton>
                    </NUpload>
                    <span class="upload-hint">{{ t('problems.fileUploadHint') }}</span>
                  </div>
                </div>

                <!-- Submit Button & Status -->
                <div class="submit-section">
                  <NButton
                    type="primary"
                    :loading="submitting"
                    :disabled="submitting"
                    @click="handleSubmit"
                    block
                    size="large"
                  >
                    {{ isCtf ? t('problems.flagSubmit') : t('problems.submitCode') }}
                  </NButton>

                  <div v-if="judgeStatus" class="judge-status">
                    <NTag :type="statusColorMap[judgeStatus] || 'default'" size="medium" round>
                      {{ getStatusLabel(judgeStatus) }}
                    </NTag>
                  </div>
                </div>
              </NCard>
            </div>
          </NGridItem>
        </NGrid>

        <!-- Submission History -->
        <div class="submission-history">
          <NCard :title="t('problems.submissionHistory')" size="small" class="history-card">
            <div v-if="submissions.length > 0" class="history-list">
              <div
                v-for="sub in submissions"
                :key="sub.id"
                class="history-item"
                @click="goToSubmission(sub.id)"
              >
                <span class="history-id">#{{ sub.id }}</span>
                <NTag :type="statusColorMap[sub.status] || 'default'" size="small" :bordered="false">
                  {{ getStatusLabel(sub.status) }}
                </NTag>
                <span class="history-lang">{{ sub.language }}</span>
                <span v-if="sub.timeUsed" class="history-time">{{ formatTime(sub.timeUsed) }}</span>
                <span v-if="sub.memoryUsed" class="history-memory">{{ formatMemory(sub.memoryUsed) }}</span>
                <span class="history-date">{{ formatDateTime(sub.createdAt) }}</span>
              </div>
            </div>
            <div v-else class="text-secondary" style="padding: 16px 0; text-align: center;">
              {{ t('problems.noSubmissions') }}
            </div>
          </NCard>
        </div>
      </template>
    </NSpin>
  </div>
</template>

<style scoped lang="scss">
.problem-detail-page {
  padding: 24px 20px;
  max-width: 1400px;
  margin: 0 auto;
}

.problem-header {
  margin-bottom: 24px;
}

.problem-title-row {
  display: flex;
  align-items: center;
  gap: 12px;
  flex-wrap: wrap;
  margin-bottom: 12px;
}

.problem-title {
  font-size: 24px;
  font-weight: 700;
  color: var(--color-text);
}

.problem-meta {
  display: flex;
  gap: 20px;
  flex-wrap: wrap;
}

.meta-item {
  font-size: 13px;
  color: var(--color-text-secondary);
}

.problem-content {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.content-card {
  background-color: var(--color-bg-card);
  border: 1px solid rgba(255, 255, 255, 0.06);
}

.markdown-body {
  font-size: 14px;
  line-height: 1.8;
  color: var(--color-text);

  :deep(h1), :deep(h2), :deep(h3) {
    margin-top: 16px;
    margin-bottom: 8px;
    color: var(--color-text);
  }

  :deep(pre) {
    background-color: rgba(0, 0, 0, 0.3);
    border-radius: 6px;
    padding: 12px;
    overflow-x: auto;
    font-size: 13px;
  }

  :deep(code) {
    font-family: 'Fira Code', 'Consolas', monospace;
    font-size: 13px;
  }

  :deep(p) {
    margin-bottom: 8px;
  }
}

.samples-list {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.sample-item {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.sample-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.sample-label {
  font-size: 13px;
  font-weight: 600;
  color: var(--color-text-secondary);
}

.sample-code {
  background-color: rgba(0, 0, 0, 0.3);
  border-radius: 6px;
  padding: 12px;
  font-family: 'Fira Code', 'Consolas', monospace;
  font-size: 13px;
  white-space: pre-wrap;
  word-break: break-all;
  color: var(--color-text);
  overflow-x: auto;
}

.hint-item {
  padding: 8px 0;
  border-bottom: 1px solid rgba(255, 255, 255, 0.06);

  &:last-child {
    border-bottom: none;
  }

  p {
    font-size: 14px;
    color: var(--color-text-secondary);
    line-height: 1.6;
  }
}

.attachments-list {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.attachment-link {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  padding: 6px 12px;
  background-color: rgba(0, 212, 255, 0.08);
  border-radius: 6px;
  color: var(--color-primary);
  font-size: 13px;
  text-decoration: none;
  transition: background-color 0.2s;

  &:hover {
    background-color: rgba(0, 212, 255, 0.15);
  }
}

.editor-panel {
  position: sticky;
  top: 80px;
}

.editor-card {
  background-color: var(--color-bg-card);
  border: 1px solid rgba(255, 255, 255, 0.06);
}

.editor-toolbar {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 12px;
}

.auto-save-hint {
  font-size: 12px;
  color: var(--color-text-secondary);
}

.editor-container {
  margin-bottom: 16px;
}

.ctf-panel {
  display: flex;
  flex-direction: column;
  gap: 20px;
  padding: 8px 0;
}

.flag-section {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.flag-label {
  font-size: 14px;
  font-weight: 500;
  color: var(--color-text-secondary);
}

.upload-section {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.upload-hint {
  font-size: 12px;
  color: var(--color-text-secondary);
}

.submit-section {
  margin-top: 16px;
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.judge-status {
  display: flex;
  justify-content: center;
}

.submission-history {
  margin-top: 24px;
}

.history-card {
  background-color: var(--color-bg-card);
  border: 1px solid rgba(255, 255, 255, 0.06);
}

.history-list {
  display: flex;
  flex-direction: column;
}

.history-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 10px 0;
  border-bottom: 1px solid rgba(255, 255, 255, 0.06);
  cursor: pointer;
  transition: background-color 0.2s;

  &:hover {
    background-color: rgba(0, 212, 255, 0.04);
  }

  &:last-child {
    border-bottom: none;
  }
}

.history-id {
  font-size: 13px;
  font-weight: 500;
  color: var(--color-text-secondary);
  min-width: 50px;
}

.history-lang {
  font-size: 13px;
  color: var(--color-text-secondary);
  min-width: 60px;
}

.history-time {
  font-size: 13px;
  color: var(--color-text-secondary);
  min-width: 70px;
}

.history-memory {
  font-size: 13px;
  color: var(--color-text-secondary);
  min-width: 70px;
}

.history-date {
  font-size: 12px;
  color: var(--color-text-secondary);
  margin-left: auto;
}

@media (max-width: 768px) {
  .problem-title-row {
    flex-direction: column;
    align-items: flex-start;
  }

  .problem-meta {
    flex-direction: column;
    gap: 4px;
  }

  .editor-panel {
    position: static;
  }

  .history-item {
    flex-wrap: wrap;
    gap: 6px;
  }

  .history-date {
    margin-left: 0;
    width: 100%;
  }
}
</style>
