<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { useAuthStore } from '@/stores/auth'
import { NCard, NForm, NFormItem, NInput, NButton, NSpace, NSpin, useMessage } from 'naive-ui'
import type { FormRules, FormInst } from 'naive-ui'

const router = useRouter()
const { t } = useI18n()
const authStore = useAuthStore()
const message = useMessage()

const formRef = ref<FormInst | null>(null)
const username = ref('')
const email = ref('')
const password = ref('')
const confirmPassword = ref('')
const loading = ref(false)

const rules: FormRules = {
  username: [
    { required: true, message: () => t('auth.usernameRequired'), trigger: 'blur' },
    {
      min: 4,
      max: 16,
      message: () => t('auth.usernameFormat'),
      trigger: 'blur',
    },
    {
      pattern: /^[a-zA-Z0-9]+$/,
      message: () => t('auth.usernameFormat'),
      trigger: 'blur',
    },
  ],
  email: [
    { required: true, message: () => t('auth.emailRequired'), trigger: 'blur' },
    {
      type: 'email',
      message: () => t('auth.emailFormat'),
      trigger: 'blur',
    },
  ],
  password: [
    { required: true, message: () => t('auth.passwordRequired'), trigger: 'blur' },
    {
      min: 8,
      message: () => t('auth.passwordMinLength'),
      trigger: 'blur',
    },
    {
      pattern: /^(?=.*[a-zA-Z])(?=.*\d)/,
      message: () => t('auth.passwordFormat'),
      trigger: 'blur',
    },
  ],
  confirmPassword: [
    { required: true, message: () => t('auth.passwordRequired'), trigger: 'blur' },
    {
      validator: (_rule: unknown, value: string) => {
        if (value !== password.value) {
          return new Error(t('auth.passwordMismatch'))
        }
        return true
      },
      trigger: 'blur',
    },
  ],
}

async function handleRegister() {
  try {
    await formRef.value?.validate()
  } catch {
    return
  }

  loading.value = true
  try {
    await authStore.register(username.value, email.value, password.value)
    message.success(t('auth.registerSuccess'))
    router.push({ name: 'Home' })
  } catch (err: unknown) {
    message.error((err instanceof Error) ? err.message : t('errors.unknownError'))
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <div class="register-page">
    <div class="register-container">
      <NCard class="register-card">
        <div class="register-header">
          <h1 class="register-title">{{ t('auth.register') }}</h1>
          <p class="register-subtitle">YogduOJ</p>
        </div>

        <NSpin :show="loading">
          <NForm
            ref="formRef"
            :model="{ username, email, password, confirmPassword }"
            :rules="rules"
            label-placement="top"
            class="register-form"
            @submit.prevent="handleRegister"
          >
            <NFormItem :label="t('auth.username')" path="username">
              <NInput
                v-model:value="username"
                :placeholder="t('auth.username')"
                autocomplete="username"
                maxlength="16"
                show-count
                @keyup.enter="handleRegister"
              />
            </NFormItem>

            <NFormItem :label="t('auth.email')" path="email">
              <NInput
                v-model:value="email"
                :placeholder="t('auth.email')"
                autocomplete="email"
                @keyup.enter="handleRegister"
              />
            </NFormItem>

            <NFormItem :label="t('auth.password')" path="password">
              <NInput
                v-model:value="password"
                type="password"
                show-password-on="click"
                :placeholder="t('auth.password')"
                autocomplete="new-password"
                @keyup.enter="handleRegister"
              />
            </NFormItem>

            <NFormItem :label="t('auth.confirmPassword')" path="confirmPassword">
              <NInput
                v-model:value="confirmPassword"
                type="password"
                show-password-on="click"
                :placeholder="t('auth.confirmPassword')"
                autocomplete="new-password"
                @keyup.enter="handleRegister"
              />
            </NFormItem>

            <NButton
              type="primary"
              block
              size="large"
              :loading="loading"
              @click="handleRegister"
            >
              {{ t('auth.register') }}
            </NButton>
          </NForm>
        </NSpin>

        <div class="register-footer">
          <span>{{ t('auth.hasAccount') }}</span>
          <router-link to="/login">{{ t('auth.login') }}</router-link>
        </div>
      </NCard>
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

.register-form {
  display: flex;
  flex-direction: column;
  gap: 4px;
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
