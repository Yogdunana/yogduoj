<script setup lang="ts">
import { useI18n } from 'vue-i18n'
import { NCard, NEmpty, NDataTable, NTag, NSpace } from 'naive-ui'

const { t } = useI18n()

const historyColumns = [
  { title: 'ID', key: 'id', width: 80 },
  { title: t('admin.aiProblemTitle'), key: 'title', ellipsis: { tooltip: true } },
  { title: t('admin.aiModel'), key: 'model', width: 140 },
  { title: t('admin.aiStatus'), key: 'status', width: 100 },
  { title: t('teams.createdAt'), key: 'createdAt', width: 170 },
]
</script>

<template>
  <div class="admin-ai-problem-manage">
    <div class="page-header">
      <h1 class="page-title">{{ t('admin.aiProblemManage') }}</h1>
    </div>

    <NCard class="coming-soon-card">
      <div class="coming-soon-content">
        <div class="icon-wrapper">
          <span class="ai-icon">AI</span>
        </div>
        <h2 class="coming-soon-title">{{ t('admin.comingSoon') }}</h2>
        <p class="coming-soon-desc">{{ t('admin.aiProblemDesc') }}</p>
        <ul class="feature-list">
          <li>{{ t('admin.aiFeature1') }}</li>
          <li>{{ t('admin.aiFeature2') }}</li>
          <li>{{ t('admin.aiFeature3') }}</li>
          <li>{{ t('admin.aiFeature4') }}</li>
        </ul>
      </div>
    </NCard>

    <NCard class="section-card" :title="t('admin.aiGenerationHistory')">
      <NDataTable
        :columns="historyColumns"
        :data="[]"
        :bordered="false"
        :single-line="false"
        size="small"
      >
        <template #empty>
          <NEmpty :description="t('admin.noAIHistory')" />
        </template>
      </NDataTable>
    </NCard>
  </div>
</template>

<style scoped lang="scss">
.admin-ai-problem-manage {
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

.coming-soon-card {
  background-color: var(--color-bg-card);
  border: 1px solid rgba(255, 255, 255, 0.06);
  margin-bottom: 16px;
}

.coming-soon-content {
  text-align: center;
  padding: 40px 20px;
}

.icon-wrapper {
  margin-bottom: 20px;
}

.ai-icon {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 80px;
  height: 80px;
  border-radius: 50%;
  background: linear-gradient(135deg, var(--color-primary), var(--color-accent));
  color: white;
  font-size: 28px;
  font-weight: 700;
}

.coming-soon-title {
  font-size: 24px;
  font-weight: 700;
  color: var(--color-text);
  margin-bottom: 12px;
}

.coming-soon-desc {
  font-size: 14px;
  color: var(--color-text-secondary);
  margin-bottom: 24px;
  max-width: 500px;
  margin-left: auto;
  margin-right: auto;
}

.feature-list {
  list-style: none;
  padding: 0;
  display: inline-block;
  text-align: left;

  li {
    padding: 6px 0;
    font-size: 14px;
    color: var(--color-text-secondary);

    &::before {
      content: '>';
      color: var(--color-primary);
      margin-right: 8px;
      font-weight: 700;
    }
  }
}

.section-card {
  background-color: var(--color-bg-card);
  border: 1px solid rgba(255, 255, 255, 0.06);

  :deep(.n-data-table) {
    --n-td-color: transparent;
    --n-th-color: rgba(255, 255, 255, 0.03);
    --n-border-color: rgba(255, 255, 255, 0.06);
    --n-td-text-color: var(--color-text);
    --n-th-text-color: var(--color-text-secondary);
  }
}
</style>
