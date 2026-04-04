<script setup lang="ts">
import { ref, computed, onMounted, watch } from 'vue'
import { useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import {
  NDataTable,
  NSelect,
  NTag,
  NCard,
  NIcon,
  NSpace,
  NButton,
  useMessage,
} from 'naive-ui'
import { listProblems } from '@/api/problem'
import type { Problem } from '@/types'
import {
  CodeSlash,
  Bug,
  Globe,
  Key,
  Search,
  ExtensionPuzzle,
  Eye,
  ShieldCheckmark,
} from '@vicons/ionicons5'

const props = defineProps<{ category: string }>()
const router = useRouter()
const { t } = useI18n()
const message = useMessage()

const loading = ref(false)
const problems = ref<Problem[]>([])
const total = ref(0)
const difficultyFilter = ref<string>('')

const iconMap: Record<string, any> = {
  reverse: CodeSlash,
  pwn: Bug,
  web: Globe,
  crypto: Key,
  forensics: Search,
  misc: ExtensionPuzzle,
  recon: Eye,
  'vuln-reproduce': ShieldCheckmark,
}

const colorMap: Record<string, string> = {
  reverse: '#00d4ff',
  pwn: '#e94560',
  web: '#00e676',
  crypto: '#ffab00',
  forensics: '#ab47bc',
  misc: '#26c6da',
  recon: '#66bb6a',
  'vuln-reproduce': '#ff7043',
}

const categoryInfo = computed(() => {
  const key = props.category
  return {
    icon: iconMap[key] || ExtensionPuzzle,
    color: colorMap[key] || '#00d4ff',
    label: t(`ctf.${key === 'vuln-reproduce' ? 'vulnReproduce' : key}`) || key,
    description: t(`ctf.${key === 'vuln-reproduce' ? 'vulnReproduceDesc' : key + 'Desc'}`) || '',
  }
})

const difficultyOptions = computed(() => [
  { label: t('problems.allDifficulties'), value: '' },
  { label: t('problems.easy'), value: 'easy' },
  { label: t('problems.medium'), value: 'medium' },
  { label: t('problems.hard'), value: 'hard' },
  { label: t('problems.expert'), value: 'expert' },
])

type NTagType = 'default' | 'error' | 'info' | 'primary' | 'success' | 'warning'

const difficultyColorMap: Record<string, NTagType> = {
  easy: 'success',
  medium: 'warning',
  hard: 'error',
  expert: 'default',
}

const columns = computed(() => [
  {
    title: t('problems.problemId'),
    key: 'id',
    width: 80,
    render(row: Problem) {
      return h('span', { class: 'problem-id' }, `#${row.id}`)
    },
  },
  {
    title: t('problems.problemTitle'),
    key: 'title',
    ellipsis: { tooltip: true },
    render(row: Problem) {
      return h(
        'a',
        {
          class: 'problem-title-link',
          onClick: () => router.push({ name: 'ProblemDetail', params: { id: row.id } }),
        },
        row.title
      )
    },
  },
  {
    title: t('problems.difficulty'),
    key: 'difficulty',
    width: 110,
    render(row: Problem) {
      return h(
        NTag,
        { type: difficultyColorMap[row.difficulty] || 'default', size: 'small', bordered: false },
        { default: () => t(`problems.${row.difficulty}`) }
      )
    },
  },
  {
    title: t('ctf.solves'),
    key: 'acceptedCount',
    width: 100,
    render(row: Problem) {
      return h('span', {}, String(row.acceptedCount))
    },
  },
  {
    title: t('problems.status'),
    key: 'userStatus',
    width: 80,
    render(row: Problem) {
      const status = row.userStatus || 'unsubmitted'
      if (status === 'unsubmitted') return h('span', { class: 'status-icon unsubmitted' }, '-')
      if (status === 'submitted') return h('span', { class: 'status-icon submitted' }, '...')
      return h('span', { class: 'status-icon accepted' }, 'AC')
    },
  },
])

async function fetchProblems() {
  loading.value = true
  try {
    const params: Record<string, any> = {
      type: 'ctf',
      tag: props.category,
      page: 1,
      pageSize: 100,
    }
    if (difficultyFilter.value) {
      params.difficulty = difficultyFilter.value
    }
    const res = await listProblems(params)
    const data = res.data.data
    problems.value = data.list
    total.value = data.total
  } catch (e: any) {
    message.error(e.message || t('errors.networkError'))
  } finally {
    loading.value = false
  }
}

function handleDifficultyChange(val: string) {
  difficultyFilter.value = val
  fetchProblems()
}

function handleRowClick(row: Problem) {
  router.push({ name: 'ProblemDetail', params: { id: row.id } })
}

function goBack() {
  router.push('/ctf')
}

watch(() => props.category, () => {
  difficultyFilter.value = ''
  fetchProblems()
})

onMounted(() => {
  fetchProblems()
})
</script>

<template>
  <div class="ctf-category-page">
    <!-- Category Header -->
    <div class="category-header">
      <NButton text @click="goBack" class="back-btn">
        &larr; {{ t('common.back') }}
      </NButton>
      <div class="header-content">
        <div class="header-icon" :style="{ color: categoryInfo.color }">
          <NIcon size="48">
            <component :is="categoryInfo.icon" />
          </NIcon>
        </div>
        <div class="header-text">
          <h1 class="header-title">{{ categoryInfo.label }}</h1>
          <p class="header-desc">{{ categoryInfo.description }}</p>
        </div>
      </div>
    </div>

    <!-- Filter Bar -->
    <div class="filter-bar">
      <NSpace align="center">
        <span class="filter-label">{{ t('problems.filterByDifficulty') }}:</span>
        <NSelect
          :value="difficultyFilter"
          :options="difficultyOptions"
          size="small"
          style="width: 160px"
          @update:value="handleDifficultyChange"
        />
        <NTag :bordered="false" size="small" type="info">
          {{ total }} {{ t('ctf.problems') }}
        </NTag>
      </NSpace>
    </div>

    <!-- Problem Table -->
    <NCard class="table-card">
      <NDataTable
        :columns="columns"
        :data="problems"
        :loading="loading"
        :row-props="(row: Problem) => ({
          style: 'cursor: pointer',
          onClick: () => handleRowClick(row),
        })"
        :bordered="false"
        :single-line="false"
        size="small"
        :row-key="(row: Problem) => row.id"
        :max-height="600"
      />
    </NCard>
  </div>
</template>

<style scoped lang="scss">
.ctf-category-page {
  padding: 24px 20px;
  max-width: 1200px;
  margin: 0 auto;
}

.category-header {
  margin-bottom: 24px;
}

.back-btn {
  color: var(--color-text-secondary);
  margin-bottom: 16px;
  font-size: 14px;

  &:hover {
    color: var(--color-primary);
  }
}

.header-content {
  display: flex;
  align-items: center;
  gap: 20px;
}

.header-icon {
  width: 72px;
  height: 72px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 16px;
  background: rgba(255, 255, 255, 0.04);
  flex-shrink: 0;
}

.header-text {
  flex: 1;
}

.header-title {
  font-size: 28px;
  font-weight: 700;
  color: var(--color-text);
  margin-bottom: 6px;
}

.header-desc {
  font-size: 15px;
  color: var(--color-text-secondary);
  line-height: 1.5;
}

.filter-bar {
  margin-bottom: 16px;
}

.filter-label {
  font-size: 14px;
  color: var(--color-text-secondary);
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

.problem-id {
  color: var(--color-text-secondary);
  font-size: 13px;
  font-weight: 500;
}

.problem-title-link {
  color: var(--color-primary);
  text-decoration: none;
  font-weight: 500;
  transition: color 0.2s;

  &:hover {
    color: var(--color-primary-hover);
  }
}

.status-icon {
  font-size: 13px;
  font-weight: 600;

  &.unsubmitted {
    color: var(--color-text-secondary);
  }

  &.submitted {
    color: var(--color-warning);
  }

  &.accepted {
    color: var(--color-success);
  }
}

@media (max-width: 768px) {
  .header-content {
    flex-direction: column;
    align-items: flex-start;
  }

  .header-title {
    font-size: 22px;
  }
}
</style>
