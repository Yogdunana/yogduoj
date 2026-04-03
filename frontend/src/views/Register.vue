<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { useAuthStore } from '@/stores/auth'

const router = useRouter()
const { t } = useI18n()
const authStore = useAuthStore()

const username = ref('')
const email = ref('')
const password = ref('')
const confirmPassword = ref('')
const loading = ref(false)
const errorMsg = ref('')

async function handleRegister() {
  errorMsg.value = ''

  if (!username.value) {
    errorMsg.value = t('auth.usernameRequired')
    return
  }
  if (!email.value) {
    errorMsg.value = t('auth.emailRequired')
    return
  }
  if (!password.value) {
    errorMsg.value = t('auth.passwordRequired')
    return
  }
  if (password.value.length < 6) {
    errorMsg.value = t('auth.passwordMinLength')
    return
  }
  if (password.value !== confirmPassword.value) {
    errorMsg.value = t('auth.passwordMismatch')
    return
  }

  loading.value = true
  try {
    await authStore.register(username.value, email.value, password.value)
    router.push({ name: 'Home' })
  } catch (err: unknown) {
    errorMsg.value = (err instanceof Error) ? err.message : t('errors.unknownError')
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <div class="register-page">
    <div class="register-container">
      <div class="register-card card">
        <div class="register-header">
          <h1 class="register-title">{{ t('auth.register') }}</h1>
          <p class="register-subtitle">YogduOJ</p>
        </div>

        <div v-if="errorMsg" class="error-banner">
          <span>{{ errorMsg }}</span>
        </div>

        <form class="register-form" @submit.prevent="handleRegister">
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
            <label class="form-label">{{ t('auth.email') }}</label>
            <input
              v-model="email"
              type="email"
              class="form-input"
              :placeholder="t('auth.email')"
              autocomplete="email"
            />
          </div>

          <div class="form-group">
            <label class="form-label">{{ t('auth.password') }}</label>
            <input
              v-model="password"
              type="password"
              class="form-input"
              :placeholder="t('auth.password')"
              autocomplete="new-password"
            />
          </div>

          <div class="form-group">
            <label class="form-label">{{ t('auth.confirmPassword') }}</label>
            <input
              v-model="confirmPassword"
              type="password"
              class="form-input"
              :placeholder="t('auth.confirmPassword')"
              autocomplete="new-password"
            />
          </div>

          <button type="submit" class="btn btn-primary btn-full" :disabled="loading">
            {{ loading ? t('common.loading') : t('auth.register') }}
          </button>
        </form>

        <div class="register-footer">
          <span>{{ t('auth.hasAccount') }}</span>
          <router-link to="/login">{{ t('auth.login') }}</router-link>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped lang="scss">
.register-page {
  display: flex;
  align-items: center;
  justify-content: center;
  min-height: calc(100vh - var(--header-height) - var(--footer-height) - 48px);
  padding: 24px;
}

.register-container {
  width: 100%;
  max-width: 420px;
}

.register-card {
  padding: 40px 32px;
}

.register-header {
  text-align: center;
  margin-bottom: 32px;
}

.register-title {
  font-size: 28px;
  font-weight: 700;
  color: var(--color-text);
  margin-bottom: 4px;
}

.register-subtitle {
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

.register-form {
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

.btn-full {
  width: 100%;
  padding: 12px;
  font-size: 16px;
}

.register-footer {
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
