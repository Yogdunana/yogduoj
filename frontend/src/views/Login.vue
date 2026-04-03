<script setup lang="ts">
import { ref } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { useAuthStore } from '@/stores/auth'

const router = useRouter()
const route = useRoute()
const { t } = useI18n()
const authStore = useAuthStore()

const username = ref('')
const password = ref('')
const rememberMe = ref(false)
const loading = ref(false)
const errorMsg = ref('')

async function handleLogin() {
  errorMsg.value = ''

  if (!username.value) {
    errorMsg.value = t('auth.usernameRequired')
    return
  }
  if (!password.value) {
    errorMsg.value = t('auth.passwordRequired')
    return
  }

  loading.value = true
  try {
    await authStore.login(username.value, password.value)
    const redirect = (route.query.redirect as string) || '/'
    router.push(redirect)
  } catch (err: unknown) {
    errorMsg.value = (err instanceof Error) ? err.message : t('errors.unknownError')
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <div class="login-page">
    <div class="login-container">
      <div class="login-card card">
        <div class="login-header">
          <h1 class="login-title">{{ t('auth.login') }}</h1>
          <p class="login-subtitle">YogduOJ</p>
        </div>

        <form v-if="errorMsg" class="error-banner" @submit.prevent>
          <span>{{ errorMsg }}</span>
        </form>

        <form class="login-form" @submit.prevent="handleLogin">
          <div class="form-group">
            <label class="form-label">{{ t('auth.username') }}</label>
            <input
              v-model="username"
              type="text"
              class="form-input"
              :placeholder="t('auth.username')"
              autocomplete="username"
            />
          </div>

          <div class="form-group">
            <label class="form-label">{{ t('auth.password') }}</label>
            <input
              v-model="password"
              type="password"
              class="form-input"
              :placeholder="t('auth.password')"
              autocomplete="current-password"
            />
          </div>

          <div class="form-row">
            <label class="checkbox-label">
              <input v-model="rememberMe" type="checkbox" />
              <span>{{ t('auth.rememberMe') }}</span>
            </label>
          </div>

          <button type="submit" class="btn btn-primary btn-full" :disabled="loading">
            {{ loading ? t('common.loading') : t('auth.login') }}
          </button>
        </form>

        <div class="login-footer">
          <span>{{ t('auth.noAccount') }}</span>
          <router-link to="/register">{{ t('auth.register') }}</router-link>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped lang="scss">
.login-page {
  display: flex;
  align-items: center;
  justify-content: center;
  min-height: calc(100vh - var(--header-height) - var(--footer-height) - 48px);
  padding: 24px;
}

.login-container {
  width: 100%;
  max-width: 420px;
}

.login-card {
  padding: 40px 32px;
}

.login-header {
  text-align: center;
  margin-bottom: 32px;
}

.login-title {
  font-size: 28px;
  font-weight: 700;
  color: var(--color-text);
  margin-bottom: 4px;
}

.login-subtitle {
  font-size: 14px;
  color: var(--color-text-secondary);
}

.error-banner {
  padding: 12px 16px;
  margin-bottom: 20px;
  border-radius: var(--border-radius);
  background-color: rgba(255, 82, 82, 0.1);
  border: 1px solid rgba(255, 82, 82, 0.3);
  color: var(--color-error);
  font-size: 14px;
}

.login-form {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.form-group {
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.form-label {
  font-size: 14px;
  font-weight: 500;
  color: var(--color-text);
}

.form-input {
  padding: 10px 14px;
  border: 1px solid rgba(255, 255, 255, 0.1);
  border-radius: var(--border-radius);
  background-color: var(--color-bg);
  color: var(--color-text);
  font-size: 14px;
  outline: none;
  transition: border-color 0.2s ease;

  &::placeholder {
    color: var(--color-text-secondary);
  }

  &:focus {
    border-color: var(--color-primary);
  }
}

.form-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.checkbox-label {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 14px;
  color: var(--color-text-secondary);
  cursor: pointer;

  input[type="checkbox"] {
    accent-color: var(--color-primary);
  }
}

.btn-full {
  width: 100%;
  padding: 12px;
  font-size: 16px;
}

.login-footer {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  margin-top: 24px;
  font-size: 14px;
  color: var(--color-text-secondary);

  a {
    color: var(--color-primary);
    font-weight: 500;
  }
}
</style>
