<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { getTeamList } from '@/api/team'
import type { Team } from '@/types'
import { NCard, NDataTable, NButton, NInput, NSpace, NSpin, NEmpty, NIcon, NTag } from 'naive-ui'
import { SearchOutline, AddOutline } from '@vicons/ionicons5'
import { useAuthStore } from '@/stores/auth'

const router = useRouter()
const { t } = useI18n()
const authStore = useAuthStore()

const loading = ref(false)
const teams = ref<Team[]>([])
const keyword = ref('')
const page = ref(1)
const pageSize = ref(20)
const total = ref(0)

function formatDate(dateStr: string) {
  return new Date(dateStr).toLocaleDateString()
}

const columns = computed(() => [
  {
    title: t('teams.teamName'),
    key: 'name',
    render(row: Team) {
      return h('a', {
        onClick: () => router.push(`/teams/${row.id}`),
        style: { color: 'var(--color-primary)', cursor: 'pointer', fontWeight: 500 },
      }, row.name)
    },
  },
  {
    title: t('teams.leader'),
    key: 'leader',
    width: 150,
    render(row: Team) {
      return row.leader?.username || '-'
    },
  },
  {
    title: t('teams.memberCount'),
    key: 'memberCount',
    width: 100,
    render(row: Team) {
      return h(NTag, { size: 'small', bordered: false, type: 'info' }, { default: () => String(row.memberCount) })
    },
  },
  {
    title: t('teams.slogan'),
    key: 'description',
    ellipsis: { tooltip: true },
  },
  {
    title: t('teams.createdAt'),
    key: 'createdAt',
    width: 140,
    render(row: Team) {
      return formatDate(row.createdAt)
    },
  },
])

async function fetchTeams() {
  loading.value = true
  try {
    const res = await getTeamList({ page: page.value, pageSize: pageSize.value, keyword: keyword.value || undefined })
    teams.value = res.data.data.list
    total.value = res.data.data.total
  } catch {
    // ignore
  } finally {
    loading.value = false
  }
}

function handleSearch() {
  page.value = 1
  fetchTeams()
}

function handlePageChange(p: number) {
  page.value = p
  fetchTeams()
}

onMounted(() => {
  fetchTeams()
})
</script>

<template>
  <div class="team-list-page">
    <div class="page-header">
      <h1 class="page-title">{{ t('teams.title') }}</h1>
      <NButton v-if="authStore.isLoggedIn" type="primary" @click="router.push('/teams/create')">
        <template #icon>
          <NIcon><AddOutline /></NIcon>
        </template>
        {{ t('teams.createTeam') }}
      </NButton>
    </div>

    <NCard class="search-card">
      <NSpace align="center">
        <NInput
          v-model:value="keyword"
          :placeholder="t('teams.searchPlaceholder')"
          clearable
          style="max-width: 320px;"
          @keyup.enter="handleSearch"
        />
        <NButton type="primary" @click="handleSearch">
          <template #icon>
            <NIcon><SearchOutline /></NIcon>
          </template>
          {{ t('common.search') }}
        </NButton>
      </NSpace>
    </NCard>

    <NCard class="table-card">
      <NSpin :show="loading">
        <NDataTable
          v-if="teams.length > 0"
          :columns="columns"
          :data="teams"
          :bordered="false"
          :single-line="false"
          striped
          :row-props="(row: Team) => ({ style: 'cursor: pointer', onClick: () => router.push(`/teams/${row.id}`) })"
        />
        <NEmpty v-else :description="t('teams.noTeams')" class="empty-state" />
      </NSpin>
    </NCard>
  </div>
</template>

<style scoped lang="scss">
.team-list-page {
  max-width: var(--content-max-width);
  margin: 0 auto;
  padding: 24px 20px;
}

.page-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 20px;
}

.page-title {
  font-size: 24px;
  font-weight: 700;
  color: var(--color-text);
  margin: 0;
}

.search-card {
  margin-bottom: 20px;
}

.table-card {
  min-height: 300px;
}

.empty-state {
  margin-top: 40px;
}
</style>
