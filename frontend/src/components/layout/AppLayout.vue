<script setup lang="ts">
</script>

<template>
  <div class="app-layout">
    <AppHeader />
    <main class="main-content">
      <router-view v-slot="{ Component }">
        <transition name="fade" mode="out-in">
          <component :is="Component" />
        </transition>
      </router-view>
    </main>
    <AppFooter />
  </div>
</template>

<script lang="ts">
import AppHeader from './AppHeader.vue'
import AppFooter from './AppFooter.vue'

export default {
  components: { AppHeader, AppFooter },
}
</script>

<style scoped lang="scss">
.app-layout {
  display: flex;
  flex-direction: column;
  min-height: 100vh;
}

.main-content {
  flex: 1;
  padding: 24px 0;
}
</style>
