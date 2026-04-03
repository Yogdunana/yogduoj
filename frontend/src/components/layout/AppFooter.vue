<script setup lang="ts">
</script>

<template>
  <footer class="app-footer">
    <div class="footer-inner container">
      <span>Powered by </span>
      <a href="https://github.com/Yogdunana" target="_blank" rel="noopener noreferrer">YogduOJ</a>
      <span>, Copyright </span>
      <a href="https://github.com/Yogdunana" target="_blank" rel="noopener noreferrer">Yogdunana-悠渡</a>
    </div>
  </footer>
</template>

<style scoped lang="scss">
.app-footer {
  height: var(--footer-height);
  display: flex;
  align-items: center;
  justify-content: center;
  background-color: var(--color-bg-secondary);
  border-top: 1px solid rgba(255, 255, 255, 0.08);
  flex-shrink: 0;
}

.footer-inner {
  display: flex;
  align-items: center;
  gap: 4px;
  font-size: 13px;
  color: var(--color-text-secondary);

  a {
    color: var(--color-primary);
    text-decoration: none;

    &:hover {
      color: var(--color-primary-hover);
    }
  }
}
</style>
