<script setup lang="ts">
withDefaults(defineProps<{
  size?: 'small' | 'medium' | 'large'
  text?: string
}>(), {
  size: 'medium',
  text: 'Loading...',
})
</script>

<template>
  <div class="loading-spinner">
    <div :class="['spinner', size]">
      <div class="spinner-ring"></div>
    </div>
    <span v-if="text" class="loading-text">{{ text }}</span>
  </div>
</template>

<style scoped lang="scss">
.loading-spinner {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 12px;
  padding: 40px;
}

.spinner {
  position: relative;
}

.spinner-ring {
  border: 3px solid rgba(255, 255, 255, 0.1);
  border-top-color: var(--color-primary);
  border-radius: 50%;
  animation: spin 0.8s linear infinite;
}

.small .spinner-ring {
  width: 24px;
  height: 24px;
  border-width: 2px;
}

.medium .spinner-ring {
  width: 40px;
  height: 40px;
  border-width: 3px;
}

.large .spinner-ring {
  width: 56px;
  height: 56px;
  border-width: 4px;
}

.loading-text {
  color: var(--color-text-secondary);
  font-size: 14px;
}

@keyframes spin {
  to {
    transform: rotate(360deg);
  }
}
</style>
