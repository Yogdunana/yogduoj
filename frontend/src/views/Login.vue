<script setup lang="ts">
import { ref } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { useAuthStore } from '@/stores/auth'
import { NCard, NForm, NFormItem, NInput, NButton, NCheckbox, NSpace, NAlert, NSpin, useMessage } from 'naive-ui'

const router = useRouter()
const route = useRoute()
const { t } = useI18n()
const authStore = useAuthStore()
const message = useMessage()

const formRef = ref()
const username = ref('')
const password = ref('')
const rememberMe = ref(false)
const loading = ref(false)

const rules = {
  username: [
    { required: true, message: () => t('auth.usernameRequired'), trigger: 'blur' },
  ],
  password: [
    { required: true, message: () => t('auth.passwordRequired'), trigger: 'blur' },
  ],
}

async function handleLogin() {
  try {
    await formRef.value?.validate()
  } catch {
    return
  }

  loading.value = true
  try {
    await authStore.login(username.value, password.value)
    const redirect = (route.query.redirect as string) || '/'
    router.push(redirect)
    message.success(t('auth.loginSuccess'))
  } catch (err: unknown) {
    message.error((err instanceof Error) ? err.message : t('errors.unknownError'))
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <div class="login-page">
    <div class="login-container">
      <NCard class="login-card">
        <div class="login-header">
          <h1 class="login-title">{{ t('auth.login') }}</h1>
          <p class="login-subtitle">YogduOJ</p>
        </div>

        <NSpin :show="loading">
          <NForm
            ref="formRef"
            :model="{ username, password }"
            :rules="rules"
            label-placement="top"
            class="login-form"
            @submit.prevent="handleLogin"
          >
            <NFormItem :label="t('auth.username')" path="username">
              <NInput
                v-model:value="username"
                :placeholder="t('auth.username')"
                autocomplete="username"
                size="large"
                @keyup.enter="handleLogin"
              />
            </NFormItem>

            <NFormItem :label="t('auth.password')" path="password">
              <NInput
                v-model:value="password"
                type="password"
                show-password-on="click"
                :placeholder="t('auth.password')"
                autocomplete="current-password"
                size="large"
                @keyup.enter="handleLogin"
              />
            </NFormItem>

            <div class="form-row">
              <NCheckbox v-model:checked="rememberMe">
                {{ t('auth.rememberMe') }}
              </NCheckbox>
              <router-link to="/forgot-password" class="forgot-link">
                {{ t('auth.forgotPassword') }}
              </router-link>
            </div>

            <NButton
              type="primary"
              block
              size="large"
              :loading="loading"
              @click="handleLogin"
            >
              {{ t('auth.login') }}
            </NButton>
          </NForm>
        </NSpin>

        <div class="login-footer">
          <span>{{ t('auth.noAccount') }}</span>
          <router-link to="/register">{{ t('auth.register') }}</router-link>
        </div>
      </NCard>
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

.login-form {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.form-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin: 8px 0 16px 0;
}

.forgot-link {
  font-size: 13px;
  color: var(--color-text-secondary);

  &:hover {
    color: var(--color-primary);
  }
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
