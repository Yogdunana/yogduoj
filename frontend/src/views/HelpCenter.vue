<script setup lang="ts">
import { ref, computed } from 'vue'
import { useI18n } from 'vue-i18n'
import {
  NCollapse,
  NCollapseItem,
  NInput,
  NCard,
  NSpace,
  NIcon,
} from 'naive-ui'
import { Search, HelpCircle, Book, Code, Trophy, People, ShieldCheckmark, Chatbubble } from '@vicons/ionicons5'

const { t } = useI18n()

const searchQuery = ref('')

interface FAQItem {
  question: string
  answer: string
}

interface HelpSection {
  title: string
  icon: any
  items: FAQItem[]
}

const sections = computed<HelpSection[]>(() => [
  {
    title: t('help.gettingStarted'),
    icon: Book,
    items: [
      {
        question: t('help.q_register'),
        answer: t('help.a_register'),
      },
      {
        question: t('help.q_login'),
        answer: t('help.a_login'),
      },
      {
        question: t('help.q_profile'),
        answer: t('help.a_profile'),
      },
    ],
  },
  {
    title: t('help.problems'),
    icon: Code,
    items: [
      {
        question: t('help.q_browseProblems'),
        answer: t('help.a_browseProblems'),
      },
      {
        question: t('help.q_submitCode'),
        answer: t('help.a_submitCode'),
      },
      {
        question: t('help.q_understandResults'),
        answer: t('help.a_understandResults'),
      },
    ],
  },
  {
    title: t('help.contests'),
    icon: Trophy,
    items: [
      {
        question: t('help.q_participate'),
        answer: t('help.a_participate'),
      },
      {
        question: t('help.q_contestRules'),
        answer: t('help.a_contestRules'),
      },
      {
        question: t('help.q_rankings'),
        answer: t('help.a_rankings'),
      },
    ],
  },
  {
    title: t('help.teams'),
    icon: People,
    items: [
      {
        question: t('help.q_createTeam'),
        answer: t('help.a_createTeam'),
      },
      {
        question: t('help.q_joinTeam'),
        answer: t('help.a_joinTeam'),
      },
      {
        question: t('help.q_manageTeam'),
        answer: t('help.a_manageTeam'),
      },
    ],
  },
  {
    title: t('help.ctfPractice'),
    icon: ShieldCheckmark,
    items: [
      {
        question: t('help.q_ctfCategories'),
        answer: t('help.a_ctfCategories'),
      },
      {
        question: t('help.q_submitFlag'),
        answer: t('help.a_submitFlag'),
      },
    ],
  },
  {
    title: t('help.faq'),
    icon: Chatbubble,
    items: [
      {
        question: t('help.q_supportedLanguages'),
        answer: t('help.a_supportedLanguages'),
      },
      {
        question: t('help.q_timeLimit'),
        answer: t('help.a_timeLimit'),
      },
      {
        question: t('help.q_cheatPolicy'),
        answer: t('help.a_cheatPolicy'),
      },
      {
        question: t('help.q_reportBug'),
        answer: t('help.a_reportBug'),
      },
    ],
  },
])

const filteredSections = computed(() => {
  if (!searchQuery.value.trim()) return sections.value
  const query = searchQuery.value.toLowerCase()
  return sections.value
    .map(section => ({
      ...section,
      items: section.items.filter(
        item =>
          item.question.toLowerCase().includes(query) ||
          item.answer.toLowerCase().includes(query)
      ),
    }))
    .filter(section => section.items.length > 0)
})
</script>

<template>
  <div class="help-center-page">
    <!-- Page Header -->
    <div class="page-header">
      <div class="header-icon">
        <NIcon size="40" color="#00d4ff">
          <HelpCircle />
        </NIcon>
      </div>
      <h1 class="page-title">{{ t('help.title') }}</h1>
      <p class="page-subtitle">{{ t('help.subtitle') }}</p>
    </div>

    <!-- Search Bar -->
    <NCard class="search-card">
      <NSpace align="center" :size="12">
        <NIcon size="20" color="var(--color-text-secondary)">
          <Search />
        </NIcon>
        <NInput
          v-model:value="searchQuery"
          :placeholder="t('help.searchPlaceholder')"
          clearable
          size="large"
        />
      </NSpace>
    </NCard>

    <!-- Help Sections -->
    <div class="sections-list">
      <NCard
        v-for="(section, idx) in filteredSections"
        :key="idx"
        class="section-card"
      >
        <div class="section-header">
          <NIcon size="22" color="#00d4ff">
            <component :is="section.icon" />
          </NIcon>
          <h2 class="section-title">{{ section.title }}</h2>
        </div>
        <NCollapse
          v-if="section.items.length > 0"
          :default-expanded-names="searchQuery ? section.items.map((_, i) => `${idx}-${i}`) : []"
        >
          <NCollapseItem
            v-for="(item, itemIdx) in section.items"
            :key="itemIdx"
            :title="item.question"
            :name="`${idx}-${itemIdx}`"
          >
            <div class="answer-text" v-html="item.answer"></div>
          </NCollapseItem>
        </NCollapse>
        <div v-else class="no-results">
          {{ t('common.noData') }}
        </div>
      </NCard>
    </div>
  </div>
</template>

<style scoped lang="scss">
.help-center-page {
  max-width: 900px;
  margin: 0 auto;
  padding: 24px 20px 40px;
}

.page-header {
  text-align: center;
  margin-bottom: 32px;
}

.header-icon {
  margin-bottom: 12px;
}

.page-title {
  font-size: 32px;
  font-weight: 700;
  color: var(--color-text);
  margin-bottom: 8px;
}

.page-subtitle {
  font-size: 16px;
  color: var(--color-text-secondary);
}

.search-card {
  margin-bottom: 32px;
  background-color: var(--color-bg-card);
  border: 1px solid rgba(255, 255, 255, 0.06);
}

.sections-list {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.section-card {
  background-color: var(--color-bg-card);
  border: 1px solid rgba(255, 255, 255, 0.06);
}

.section-header {
  display: flex;
  align-items: center;
  gap: 10px;
  margin-bottom: 16px;
}

.section-title {
  font-size: 20px;
  font-weight: 600;
  color: var(--color-text);
}

.answer-text {
  font-size: 14px;
  color: var(--color-text-secondary);
  line-height: 1.8;

  :deep(a) {
    color: var(--color-primary);
    text-decoration: none;

    &:hover {
      text-decoration: underline;
    }
  }

  :deep(code) {
    background: rgba(0, 212, 255, 0.1);
    padding: 2px 6px;
    border-radius: 4px;
    font-size: 13px;
    color: var(--color-primary);
  }

  :deep(ul),
  :deep(ol) {
    padding-left: 20px;
    margin: 8px 0;
  }

  :deep(li) {
    margin-bottom: 4px;
  }
}

.no-results {
  text-align: center;
  padding: 16px;
  color: var(--color-text-secondary);
  font-size: 14px;
}

:deep(.n-collapse) {
  --n-title-text-color: var(--color-text);
  --n-title-font-size: 15px;
  --n-arrow-color: var(--color-text-secondary);
}

:deep(.n-collapse-item__content-inner) {
  padding-top: 4px;
}
</style>
