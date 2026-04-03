<script setup lang="ts">
const props = withDefaults(defineProps<{
  total: number
  page: number
  pageSize: number
}>(), {
  total: 0,
  page: 1,
  pageSize: 20,
})

const emit = defineEmits<{
  'update:page': [value: number]
}>()

const totalPages = computed(() => Math.ceil(props.total / props.pageSize))

function goToPage(page: number) {
  if (page >= 1 && page <= totalPages.value) {
    emit('update:page', page)
  }
}

import { computed } from 'vue'
</script>

<template>
  <div v-if="totalPages > 1" class="pagination">
    <button
      class="page-btn"
      :disabled="page <= 1"
      @click="goToPage(page - 1)"
    >
      &laquo;
    </button>

    <button
      v-for="p in totalPages"
      :key="p"
      :class="['page-btn', { active: p === page }]"
      @click="goToPage(p)"
    >
      {{ p }}
    </button>

    <button
      class="page-btn"
      :disabled="page >= totalPages"
      @click="goToPage(page + 1)"
    >
      &raquo;
    </button>
  </div>
</template>

<style scoped lang="scss">
.pagination {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 4px;
  padding: 16px 0;
}

.page-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  min-width: 36px;
  height: 36px;
  padding: 0 8px;
  border: 1px solid rgba(255, 255, 255, 0.1);
  border-radius: var(--border-radius);
  background-color: var(--color-bg-card);
  color: var(--color-text);
  font-size: 14px;
  cursor: pointer;
  transition: all 0.2s ease;

  &:hover:not(:disabled):not(.active) {
    border-color: var(--color-primary);
    color: var(--color-primary);
  }

  &.active {
    background-color: var(--color-primary);
    border-color: var(--color-primary);
    color: #1a1a2e;
    font-weight: 600;
  }

  &:disabled {
    opacity: 0.3;
    cursor: not-allowed;
  }
}
</style>
