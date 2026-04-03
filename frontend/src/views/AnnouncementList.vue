<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { NCard, NTag, NButton, NInput, NEmpty, NPagination, NSpace, useMessage } from 'naive-ui'
import { listAnnouncements } from '@/api/announcement'
import type { Announcement } from '@/types'

const { t } = useI18n()
const router = useRouter()
const message = useMessage()

const loading = ref(false)
const announcements = ref<Announcement[]>([])
const total = ref(0)

const filters = reactive({
  keyword: '',
  page: 1,
  pageSize: 10,
})

function formatTime(dateStr: string): string {
  return new Date(dateStr).toLocaleString()
}

function getContentPreview(content: string, maxLen: number = 150): string {
  const plain = content.replace(/<[^>]*>/g, '').replace(/&[^;]+;/g, ' ')
  if (plain.length <= maxLen) return plain
  return plain.substring(0, maxLen) + '...'
}

async function fetchAnnouncements() {
  loading.value = true
  try {
    const params: Record<string, any> = {
      page: filters.page,
      pageSize: filters.pageSize,
    }
    if (filters.keyword) params.keyword = filters.keyword

    const res = await listAnnouncements(params)
    const data = res.data.data
    announcements.value = data.items
    total.value = data.total
  } catch (e: any) {
    message.error(e.message || t('errors.networkError'))
  } finally {
    loading.value = false
  }
}

function handlePageChange(page: number) {
  filters.page = page
  fetchAnnouncements()
}

let searchTimer: ReturnType<typeof setTimeout> | null = null
function handleKeywordInput(val: string) {
  filters.keyword = val
  if (searchTimer) clearTimeout(searchTimer)
  searchTimer = setTimeout(() => {
    filters.page = 1
    fetchAnnouncements()
  }, 500)
}

function goToAnnouncement(id: number) {
  router.push({ name: 'AnnouncementDetail', params: { id } })
}

onMounted(() => {
  fetchAnnouncements()
})
</script>

<template>
  <div class="announcement-list-page">
    <div class="page-header">
      <h1 class="page-title">{{ t('announcements.title') }}</h1>
    </div>

    <!-- Search -->
    <div class="toolbar">
      <NInput
        :placeholder="t('announcements.searchPlaceholder')"
        :value="filters.keyword"
        clearable
        @update:value="handleKeywordInput"
        class="search-input"
      />
    </div>

    <!-- Announcement Cards -->
    <div v-if="announcements.length > 0" class="announcement-list">
      <NCard
        v-for="item in announcements"
        :key="item.id"
        class="announcement-card"
        hoverable
        @click="goToAnnouncement(item.id)"
      >
        <div class="card-header">
          <div class="card-title-row">
            <NTag v-if="item.isPinned" type="warning" size="small" bordered class="pinned-tag">
              {{ t('announcements.pinned') }}
            </NTag>
            <h3 class="card-title">{{ item.title }}</h3>
          </div>
          <div class="card-meta">
            <span class="meta-item">
              {{ t('announcements.author') }}: {{ item.author?.username || '-' }}
            </span>
            <span class="meta-item">
              {{ t('announcements.createdAt') }}: {{ formatTime(item.createdAt) }}
            </span>
          </div>
        </div>
        <div class="card-content">
          {{ getContentPreview(item.content) }}
        </div>
        <div class="card-footer">
          <NButton text type="primary" size="small">
            {{ t('announcements.readMore') }} &rarr;
          </NButton>
        </div>
      </NCard>
    </div>

    <NEmpty v-else-if="!loading" :description="t('announcements.noAnnouncements')" class="empty-state" />

    <!-- Pagination -->
    <div v-if="total > 0" class="pagination-wrapper">
      <NPagination
        :page="filters.page"
        :page-size="filters.pageSize"
        :item-count="total"
        @update:page="handlePageChange"
      />
    </div>
  </div>
</template>

<style scoped lang="scss">
.announcement-list-page {
  padding: 24px 20px;
  max-width: 900px;
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
  margin-bottom: 20px;
}

.search-input {
  max-width: 400px;
}

.announcement-list {
  display: flex;
  flex-direction: column;
  gap: 16px;
  margin-bottom: 24px;
}

.announcement-card {
  background-color: var(--color-bg-card);
  border: 1px solid rgba(255, 255, 255, 0.06);
  cursor: pointer;
  transition: border-color 0.2s;

  &:hover {
    border-color: var(--color-primary);
  }

  :deep(.n-card__content) {
    padding: 20px;
  }
}

.card-header {
  margin-bottom: 12px;
}

.card-title-row {
  display: flex;
  align-items: center;
  gap: 10px;
  margin-bottom: 8px;
}

.pinned-tag {
  flex-shrink: 0;
}

.card-title {
  font-size: 17px;
  font-weight: 600;
  color: var(--color-text);
  margin: 0;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.card-meta {
  display: flex;
  gap: 20px;
  flex-wrap: wrap;
}

.meta-item {
  font-size: 13px;
  color: var(--color-text-secondary);
}

.card-content {
  font-size: 14px;
  color: var(--color-text-secondary);
  line-height: 1.7;
  margin-bottom: 12px;
}

.card-footer {
  display: flex;
  justify-content: flex-end;
}

.empty-state {
  margin-top: 60px;
}

.pagination-wrapper {
  display: flex;
  justify-content: center;
  padding: 8px 0;
}

@media (max-width: 768px) {
  .search-input {
    max-width: none;
  }

  .card-meta {
    flex-direction: column;
    gap: 4px;
  }
}
</style>
