<script setup lang="ts">
import { ref, computed } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useI18n } from 'vue-i18n'

const router = useRouter()
const route = useRoute()
const { t } = useI18n()

const collapsed = ref(false)

const menuGroups = computed(() => [
  {
    label: t('admin.dashboard'),
    items: [
      { name: 'AdminDashboard', label: t('admin.dashboard'), path: '/admin', icon: '📊' },
    ],
  },
  {
    label: t('admin.userManage'),
    items: [
      { name: 'AdminUserManage', label: t('admin.userManage'), path: '/admin/users', icon: '👥' },
    ],
  },
  {
    label: t('admin.problemManage'),
    items: [
      { name: 'AdminProblemManage', label: t('admin.problemManage'), path: '/admin/problems', icon: '📝' },
      { name: 'AdminProblemCreate', label: t('common.create'), path: '/admin/problems/create', icon: '➕' },
    ],
  },
  {
    label: t('admin.contestManage'),
    items: [
      { name: 'AdminContestManage', label: t('admin.contestManage'), path: '/admin/contests', icon: '🏆' },
      { name: 'AdminContestCreate', label: t('common.create'), path: '/admin/contests/create', icon: '➕' },
    ],
  },
  {
    label: t('admin.submissionManage'),
    items: [
      { name: 'AdminSubmissionManage', label: t('admin.submissionManage'), path: '/admin/submissions', icon: '📋' },
      { name: 'AdminJudgeMonitor', label: t('admin.judgeMonitor'), path: '/admin/judge', icon: '⚡' },
    ],
  },
  {
    label: t('admin.announcementManage'),
    items: [
      { name: 'AdminAnnouncementManage', label: t('admin.announcementManage'), path: '/admin/announcements', icon: '📢' },
    ],
  },
  {
    label: t('common.more'),
    items: [
      { name: 'AdminAIProblemManage', label: t('admin.aiProblemManage'), path: '/admin/ai-problems', icon: '🤖' },
      { name: 'AdminImportManage', label: t('admin.importManage'), path: '/admin/import', icon: '📥' },
      { name: 'AdminCheatManage', label: t('admin.cheatManage'), path: '/admin/cheat', icon: '🔍' },
      { name: 'AdminSystemConfig', label: t('admin.systemConfig'), path: '/admin/config', icon: '⚙️' },
      { name: 'AdminDIYTemplates', label: t('admin.diyTemplates'), path: '/admin/diy-templates', icon: '🎨' },
    ],
  },
])

function isActive(path: string): boolean {
  if (path === '/admin') return route.path === '/admin'
  return route.path.startsWith(path)
}

function navigate(path: string) {
  router.push(path)
}

function toggleCollapse() {
  collapsed.value = !collapsed.value
}
</script>

<template>
  <aside :class="['app-sidebar', { collapsed }]">
    <div class="sidebar-toggle" @click="toggleCollapse">
      <span :class="['toggle-icon', { rotated: collapsed }]">◀</span>
    </div>

    <nav class="sidebar-nav">
      <div v-for="group in menuGroups" :key="group.label" class="menu-group">
        <div v-if="!collapsed" class="menu-group-label">{{ group.label }}</div>
        <router-link
          v-for="item in group.items"
          :key="item.name"
          :to="item.path"
          :class="['menu-item', { active: isActive(item.path) }]"
          :title="collapsed ? item.label : undefined"
        >
          <span class="menu-icon">{{ item.icon }}</span>
          <span v-if="!collapsed" class="menu-label">{{ item.label }}</span>
        </router-link>
      </div>
    </nav>
  </aside>
</template>

<style scoped lang="scss">
.app-sidebar {
  position: fixed;
  left: 0;
  top: var(--header-height);
  bottom: var(--footer-height);
  width: var(--sidebar-width);
  background-color: var(--color-bg-secondary);
  border-right: 1px solid rgba(255, 255, 255, 0.08);
  overflow-y: auto;
  overflow-x: hidden;
  transition: width 0.3s ease;
  z-index: 100;

  &.collapsed {
    width: var(--sidebar-collapsed-width);
  }
}

.sidebar-toggle {
  display: flex;
  align-items: center;
  justify-content: flex-end;
  padding: 12px;
  cursor: pointer;
}

.toggle-icon {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 24px;
  height: 24px;
  border-radius: 4px;
  background-color: rgba(255, 255, 255, 0.05);
  color: var(--color-text-secondary);
  font-size: 12px;
  transition: transform 0.3s ease;

  &:hover {
    background-color: rgba(255, 255, 255, 0.1);
  }

  &.rotated {
    transform: rotate(180deg);
  }
}

.sidebar-nav {
  padding: 0 8px 16px;
}

.menu-group {
  margin-bottom: 8px;
}

.menu-group-label {
  padding: 8px 12px 4px;
  font-size: 11px;
  font-weight: 600;
  text-transform: uppercase;
  color: var(--color-text-secondary);
  letter-spacing: 0.5px;
}

.menu-item {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 8px 12px;
  border-radius: var(--border-radius);
  color: var(--color-text-secondary);
  font-size: 14px;
  text-decoration: none;
  transition: all 0.2s ease;
  white-space: nowrap;
  overflow: hidden;

  &:hover {
    color: var(--color-text);
    background-color: rgba(255, 255, 255, 0.05);
  }

  &.active {
    color: var(--color-primary);
    background-color: rgba(0, 212, 255, 0.1);
  }
}

.menu-icon {
  flex-shrink: 0;
  font-size: 16px;
  width: 24px;
  text-align: center;
}

.menu-label {
  overflow: hidden;
  text-overflow: ellipsis;
}
</style>
