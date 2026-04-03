<script setup lang="ts">
import { ref, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { createTeam } from '@/api/team'
import { NCard, NForm, NFormItem, NInput, NButton, NSpace, useMessage } from 'naive-ui'

const router = useRouter()
const { t } = useI18n()
const message = useMessage()

const loading = ref(false)
const form = reactive({
  name: '',
  description: '',
})

async function handleCreate() {
  if (!form.name.trim()) {
    message.warning(t('teams.teamName') + ' ' + t('auth.usernameRequired'))
    return
  }
  if (form.name.trim().length < 2) {
    message.warning(t('auth.usernameFormat'))
    return
  }

  loading.value = true
  try {
    const res = await createTeam({ name: form.name.trim(), description: form.description.trim() })
    message.success(t('teams.createSuccess'))
    router.push(`/teams/${res.data.data.id}`)
  } catch (err: unknown) {
    message.error((err instanceof Error) ? err.message : t('errors.unknownError'))
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <div class="team-create-page">
    <div class="create-container">
      <NCard class="create-card">
        <div class="create-header">
          <h1 class="create-title">{{ t('teams.createTeam') }}</h1>
        </div>

        <NForm label-placement="top" class="create-form">
          <NFormItem :label="t('teams.teamName')">
            <NInput
              v-model:value="form.name"
              :placeholder="t('teams.teamName')"
              maxlength="50"
              show-count
            />
          </NFormItem>

          <NFormItem :label="t('teams.slogan')">
            <NInput
              v-model:value="form.description"
              type="textarea"
              :placeholder="t('teams.slogan')"
              :rows="3"
              maxlength="200"
              show-count
            />
          </NFormItem>

          <NSpace class="form-actions">
            <NButton @click="router.back()">{{ t('common.cancel') }}</NButton>
            <NButton type="primary" :loading="loading" @click="handleCreate">
              {{ t('common.create') }}
            </NButton>
          </NSpace>
        </NForm>
      </NCard>
    </div>
  </div>
</template>

<style scoped lang="scss">
.team-create-page {
  max-width: var(--content-max-width);
  margin: 0 auto;
  padding: 24px 20px;
}

.create-container {
  max-width: 560px;
  margin: 0 auto;
}

.create-header {
  margin-bottom: 24px;
}

.create-title {
  font-size: 24px;
  font-weight: 700;
  color: var(--color-text);
  margin: 0;
}

.create-form {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.form-actions {
  margin-top: 16px;
  justify-content: flex-end;
}
</style>
