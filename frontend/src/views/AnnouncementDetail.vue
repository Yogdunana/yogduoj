<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { NCard, NButton, NTag, NSpin, useMessage } from 'naive-ui'
import { getAnnouncement } from '@/api/announcement'
import type { Announcement } from '@/types'

const props = defineProps<{ id: string }>()
const router = useRouter()
const { t } = useI18n()
const message = useMessage()

const loading = ref(true)
const announcement = ref<Announcement | null>(null)

function formatTime(dateStr: string): string {
  return new Date(dateStr).toLocaleString()
}

function goBack() {
  router.push({ name: 'AnnouncementList' })
}

async function fetchAnnouncement() {
  loading.value = true
  try {
    const res = await getAnnouncement(Number(props.id))
    announcement.value = res.data.data
  } catch (e: any) {
    message.error(e.message || t('errors.networkError'))
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  fetchAnnouncement()
})
</script>

<template>
  <div class="announcement-detail-page">
    <NSpin :show="loading">
      <template v-if="announcement">
        <!-- Back Button -->
        <div class="back-row">
          <NButton text @click="goBack" class="back-btn">
            &larr; {{ t('announcements.backToList') }}
          </NButton>
        </div>

        <!-- Announcement Card -->
        <NCard class="detail-card">
          <div class="detail-header">
            <div class="title-row">
              <NTag v-if="announcement.isPinned" type="warning" size="small" bordered>
                {{ t('announcements.pinned') }}
              </NTag>
              <h1 class="detail-title">{{ announcement.title }}</h1>
            </div>
            <div class="detail-meta">
              <span class="meta-item">
                {{ t('announcements.author') }}: {{ announcement.author?.username || '-' }}
              </span>
              <span class="meta-item">
                {{ t('announcements.createdAt') }}: {{ formatTime(announcement.createdAt) }}
              </span>
              <span v-if="announcement.updatedAt !== announcement.createdAt" class="meta-item">
                {{ t('announcements.updatedAt') }}: {{ formatTime(announcement.updatedAt) }}
              </span>
            </div>
          </div>

          <div class="divider" />

          <div class="detail-content" v-html="announcement.content" />
        </NCard>
      </template>
    </NSpin>
  </div>
</template>

<style scoped lang="scss">
.announcement-detail-page {
  padding: 24px 20px;
  max-width: 900px;
  margin: 0 auto;
}

.back-row {
  margin-bottom: 16px;
}

.back-btn {
  font-size: 14px;
  color: var(--color-text-secondary);
  padding: 0;

  &:hover {
    color: var(--color-primary);
  }
}

.detail-card {
  background-color: var(--color-bg-card);
  border: 1px solid rgba(255, 255, 255, 0.06);

  :deep(.n-card__content) {
    padding: 28px;
  }
}

.detail-header {
  margin-bottom: 20px;
}

.title-row {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-bottom: 12px;
  flex-wrap: wrap;
}

.detail-title {
  font-size: 24px;
  font-weight: 700;
  color: var(--color-text);
  margin: 0;
}

.detail-meta {
  display: flex;
  gap: 24px;
  flex-wrap: wrap;
}

.meta-item {
  font-size: 13px;
  color: var(--color-text-secondary);
}

.divider {
  height: 1px;
  background: rgba(255, 255, 255, 0.06);
  margin-bottom: 24px;
}

.detail-content {
  color: var(--color-text);
  line-height: 1.8;
  font-size: 15px;

  :deep(h1), :deep(h2), :deep(h3) {
    color: var(--color-text);
    margin-top: 24px;
    margin-bottom: 12px;
  }

  :deep(h1) {
    font-size: 22px;
  }

  :deep(h2) {
    font-size: 19px;
  }

  :deep(h3) {
    font-size: 16px;
  }

  :deep(p) {
    margin-bottom: 12px;
  }

  :deep(ul), :deep(ol) {
    padding-left: 24px;
    margin-bottom: 12px;
  }

  :deep(li) {
    margin-bottom: 4px;
  }

  :deep(pre) {
    background: rgba(0, 0, 0, 0.3);
    padding: 16px;
    border-radius: 8px;
    overflow-x: auto;
    margin: 16px 0;
  }

  :deep(code) {
    font-family: 'Courier New', monospace;
    font-size: 13px;
  }

  :deep(:not(pre) > code) {
    background: rgba(0, 0, 0, 0.2);
    padding: 2px 6px;
    border-radius: 4px;
  }

  :deep(blockquote) {
    border-left: 3px solid var(--color-primary);
    padding-left: 16px;
    margin: 16px 0;
    color: var(--color-text-secondary);
  }

  :deep(table) {
    border-collapse: collapse;
    width: 100%;
    margin: 16px 0;
  }

  :deep(th), :deep(td) {
    border: 1px solid rgba(255, 255, 255, 0.1);
    padding: 10px 14px;
    text-align: left;
  }

  :deep(th) {
    background: rgba(255, 255, 255, 0.03);
    font-weight: 600;
  }

  :deep(a) {
    color: var(--color-primary);
    text-decoration: none;

    &:hover {
      text-decoration: underline;
    }
  }

  :deep(img) {
    max-width: 100%;
    border-radius: 8px;
    margin: 12px 0;
  }
}

@media (max-width: 768px) {
  .detail-meta {
    flex-direction: column;
    gap: 4px;
  }

  .detail-card {
    :deep(.n-card__content) {
      padding: 16px;
    }
  }
}
</style>
