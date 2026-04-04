<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useI18n } from 'vue-i18n'
import {
  NForm, NFormItem, NInput, NSelect, NButton, NSpace, NCard, NDatePicker,
  NInputNumber, NSwitch, useMessage,
} from 'naive-ui'
import { getContestDetail } from '@/api/contest'
import { adminUpdateContest } from '@/api/admin'
import { listProblems } from '@/api/problem'
import type { Problem } from '@/types'

const props = defineProps<{ id: string }>()
const { t } = useI18n()
const router = useRouter()
const route = useRoute()
const message = useMessage()

const loading = ref(false)
const submitting = ref(false)
const contestId = computed(() => Number(props.id) || Number(route.params.id))
const allProblems = ref<Problem[]>([])

const form = reactive({
  title: '',
  type: 'icpc' as string,
  category: '' as string,
  ruleType: 'acm' as string,
  description: '',
  startTime: null as number | null,
  endTime: null as number | null,
  freezeTime: null as number | null,
  maxTeamSize: 1,
  signupLimit: 0,
  allowViewOthers: true,
  showRealtimeRank: true,
  enableAIHint: false,
  problemIds: [] as number[],
  problemScores: [] as { problemId: number; score: number; label: string; displayOrder: number }[],
  diyRules: '',
})

const typeOptions = computed(() => [
  { label: 'ICPC', value: 'icpc' },
  { label: 'IOI', value: 'ioi' },
  { label: 'CTF', value: 'ctf' },
])

const ruleTypeOptions = computed(() => [
  { label: 'ACM', value: 'acm' },
  { label: 'OI', value: 'oi' },
  { label: 'Codeforces', value: 'cf' },
  { label: 'DIY', value: 'diy' },
])

const categoryOptions = computed(() => [
  { label: t('admin.categoryOfficial'), value: 'official' },
  { label: t('admin.categoryTraining'), value: 'training' },
  { label: t('admin.categoryPrivate'), value: 'private' },
])

const problemOptions = computed(() =>
  allProblems.value.map(p => ({
    label: `#${p.id} ${p.title}`,
    value: p.id,
  }))
)

const isDIY = computed(() => form.ruleType === 'diy')

async function fetchContest() {
  loading.value = true
  try {
    const res = await getContestDetail(contestId.value)
    const data = res.data.data
    form.title = data.title
    form.type = data.type
    form.ruleType = data.ruleType
    form.description = data.description
    form.startTime = new Date(data.startTime).getTime()
    form.endTime = new Date(data.endTime).getTime()
    form.problemIds = data.problemIds || []
    form.category = (data as any).category || ''
    form.maxTeamSize = (data as any).maxTeamSize || 1
    form.signupLimit = (data as any).signupLimit || 0
    form.allowViewOthers = (data as any).allowViewOthers !== false
    form.showRealtimeRank = (data as any).showRealtimeRank !== false
    form.enableAIHint = !!(data as any).enableAIHint
    form.problemScores = (data as any).problemScores || form.problemIds.map((id, i) => ({
      problemId: id, score: 100, label: String.fromCharCode(65 + i), displayOrder: i + 1,
    }))
    if ((data as any).freezeTime) {
      form.freezeTime = new Date((data as any).freezeTime).getTime()
    }
    if ((data as any).diyRules) {
      form.diyRules = typeof (data as any).diyRules === 'string'
        ? (data as any).diyRules
        : JSON.stringify((data as any).diyRules, null, 2)
    }
  } catch (e: any) {
    message.error(e.message || t('errors.networkError'))
  } finally {
    loading.value = false
  }
}

async function fetchProblems() {
  try {
    const res = await listProblems({ page: 1, pageSize: 500 })
    allProblems.value = res.data.data.list
  } catch {
    // ignore
  }
}

function handleProblemSelect(ids: number[]) {
  form.problemIds = ids
  const existing = new Map(form.problemScores.map(p => [p.problemId, p]))
  form.problemScores = ids.map((id, index) => {
    if (existing.has(id)) return existing.get(id)!
    return { problemId: id, score: 100, label: String.fromCharCode(65 + index), displayOrder: index + 1 }
  })
}

function updateProblemScore(index: number, field: 'score' | 'label' | 'displayOrder', value: any) {
  if (form.problemScores[index]) {
    (form.problemScores[index] as any)[field] = value
  }
}

async function handleSubmit() {
  if (!form.title.trim()) {
    message.warning(t('admin.titleRequired'))
    return
  }
  if (!form.startTime || !form.endTime) {
    message.warning(t('admin.timeRequired'))
    return
  }

  submitting.value = true
  try {
    const data: Record<string, any> = {
      title: form.title,
      type: form.type,
      category: form.category,
      ruleType: form.ruleType,
      description: form.description,
      startTime: new Date(form.startTime).toISOString(),
      endTime: new Date(form.endTime).toISOString(),
      maxTeamSize: form.maxTeamSize,
      signupLimit: form.signupLimit || undefined,
      allowViewOthers: form.allowViewOthers,
      showRealtimeRank: form.showRealtimeRank,
      enableAIHint: form.enableAIHint,
      problemIds: form.problemIds,
      problemScores: form.problemScores,
    }

    if (form.freezeTime) {
      data.freezeTime = new Date(form.freezeTime).toISOString()
    }

    if (isDIY.value && form.diyRules.trim()) {
      try {
        data.diyRules = JSON.parse(form.diyRules)
      } catch {
        message.error(t('admin.invalidJson'))
        submitting.value = false
        return
      }
    }

    await adminUpdateContest(contestId.value, data as any)
    message.success(t('common.success'))
    router.push('/admin/contests')
  } catch (e: any) {
    message.error(e.message || t('errors.networkError'))
  } finally {
    submitting.value = false
  }
}

onMounted(() => {
  fetchContest()
  fetchProblems()
})
</script>

<template>
  <div class="admin-contest-edit">
    <div class="page-header">
      <h1 class="page-title">{{ t('admin.editContest') }} #{{ contestId }}</h1>
    </div>

    <NCard class="form-card" :loading="loading">
      <NForm label-placement="left" label-width="160">
        <NFormItem :label="t('contests.title')" required>
          <NInput v-model:value="form.title" :placeholder="t('admin.contestTitlePlaceholder')" />
        </NFormItem>

        <NFormItem :label="t('contests.type')">
          <NSelect v-model:value="form.type" :options="typeOptions" style="width: 300px" />
        </NFormItem>

        <NFormItem :label="t('admin.category')">
          <NSelect v-model:value="form.category" :options="categoryOptions" style="width: 300px" clearable />
        </NFormItem>

        <NFormItem :label="t('admin.ruleType')">
          <NSelect v-model:value="form.ruleType" :options="ruleTypeOptions" style="width: 300px" />
        </NFormItem>

        <NFormItem :label="t('teams.description')">
          <NInput v-model:value="form.description" type="textarea" :rows="6" />
        </NFormItem>

        <NFormItem :label="t('contests.startTime')" required>
          <NDatePicker v-model:value="form.startTime" type="datetime" style="width: 300px" />
        </NFormItem>

        <NFormItem :label="t('contests.endTime')" required>
          <NDatePicker v-model:value="form.endTime" type="datetime" style="width: 300px" />
        </NFormItem>

        <NFormItem :label="t('admin.freezeTime')">
          <NDatePicker v-model:value="form.freezeTime" type="datetime" style="width: 300px" clearable />
        </NFormItem>

        <NFormItem :label="t('admin.maxTeamSize')">
          <NInputNumber v-model:value="form.maxTeamSize" :min="1" :max="10" style="width: 300px" />
        </NFormItem>

        <NFormItem :label="t('admin.signupLimit')">
          <NInputNumber v-model:value="form.signupLimit" :min="0" style="width: 300px" />
        </NFormItem>

        <NFormItem :label="t('admin.allowViewOthers')">
          <NSwitch v-model:value="form.allowViewOthers" />
        </NFormItem>

        <NFormItem :label="t('admin.showRealtimeRank')">
          <NSwitch v-model:value="form.showRealtimeRank" />
        </NFormItem>

        <NFormItem :label="t('admin.enableAIHint')">
          <NSwitch v-model:value="form.enableAIHint" />
        </NFormItem>

        <NFormItem :label="t('admin.contestProblems')">
          <NSelect
            :value="form.problemIds"
            :options="problemOptions"
            multiple
            filterable
            style="width: 100%"
            @update:value="handleProblemSelect"
          />
        </NFormItem>

        <div v-if="form.problemScores.length > 0" class="problem-scores">
          <div v-for="(ps, index) in form.problemScores" :key="ps.problemId" class="score-row">
            <span class="score-label">#{{ ps.problemId }}</span>
            <NInput
              :value="ps.label"
              size="small"
              style="width: 60px"
              @update:value="(v: string) => updateProblemScore(index, 'label', v)"
            />
            <NInputNumber
              :value="ps.score"
              size="small"
              :min="0"
              style="width: 100px"
              @update:value="(v: number | null) => updateProblemScore(index, 'score', v)"
            />
            <NInputNumber
              :value="ps.displayOrder"
              size="small"
              :min="1"
              style="width: 100px"
              @update:value="(v: number | null) => updateProblemScore(index, 'displayOrder', v)"
            />
          </div>
        </div>

        <template v-if="isDIY">
          <NFormItem :label="t('admin.diyRules')">
            <NInput
              v-model:value="form.diyRules"
              type="textarea"
              :rows="8"
              placeholder='{"scoring": "...", "penalty": "..."}'
            />
          </NFormItem>
        </template>

        <NFormItem label=" ">
          <NSpace>
            <NButton type="primary" :loading="submitting" @click="handleSubmit">
              {{ t('common.save') }}
            </NButton>
            <NButton @click="router.push('/admin/contests')">
              {{ t('common.cancel') }}
            </NButton>
          </NSpace>
        </NFormItem>
      </NForm>
    </NCard>
  </div>
</template>

<style scoped lang="scss">
.admin-contest-edit {
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

.problem-scores {
  margin-bottom: 16px;
  padding: 12px;
  background-color: rgba(255, 255, 255, 0.03);
  border-radius: var(--border-radius);
}

.score-row {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-bottom: 8px;

  &:last-child {
    margin-bottom: 0;
  }
}

.score-label {
  font-size: 13px;
  color: var(--color-text-secondary);
  min-width: 50px;
}
</style>
