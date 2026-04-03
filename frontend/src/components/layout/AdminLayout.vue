<script setup lang="ts">
import { ref } from 'vue'

const sidebarCollapsed = ref(false)
</script>

<template>
  <div class="admin-layout">
    <AppHeader />
    <div class="admin-body">
      <AppSidebar />
      <main :class="['admin-content', { collapsed: sidebarCollapsed }]">
        <router-view v-slot="{ Component }">
          <transition name="fade" mode="out-in">
            <component :is="Component" />
          </transition>
        </router-view>
      </main>
    </div>
    <AppFooter />
  </div>
</template>

<script lang="ts">
import AppHeader from './AppHeader.vue'
import AppSidebar from './AppSidebar.vue'
import AppFooter from './AppFooter.vue'

export default {
  components: { AppHeader, AppSidebar, AppFooter },
}
</script>

<style scoped lang="scss">
.admin-layout {
  display: flex;
  flex-direction: column;
  min-height: 100vh;
}

.admin-body {
  display: flex;
  flex: 1;
}

.admin-content {
  flex: 1;
  margin-left: var(--sidebar-width);
  padding: 24px;
  min-height: calc(100vh - var(--header-height) - var(--footer-height));
  transition: margin-left 0.3s ease;

  &.collapsed {
    margin-left: var(--sidebar-collapsed-width);
  }
}
</style>
