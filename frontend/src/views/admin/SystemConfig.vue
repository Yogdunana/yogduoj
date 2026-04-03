<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import {
  NCard, NTabs, NTabPane, NForm, NFormItem, NInput, NInputNumber, NSelect, NButton, NSpace, useMessage,
} from 'naive-ui'
import { getSystemConfig, updateSystemConfig } from '@/api/admin'

const { t } = useI18n()
const message = useMessage()

const loading = ref(false)
const activeTab = ref('basic')

const basicConfig = reactive({
  systemName: '',
  logoUrl: '',
  themeColor: '',
  announcement: '',
})

const judgeConfig = reactive({
  defaultTimeLimit: 1000,
  defaultMemoryLimit: 262144,
  poolSize: 4,
  supportedLanguages: '',
})

const securityConfig = reactive({
  loginFailLimit: 5,
  passwordComplexity: 'medium' as string,
  backupFrequency: 'daily' as string,
})

const aiConfig = reactive({
  provider: '',
  token: '',
})

async function fetchConfig() {
  loading.value = true
  try {
    const res = await getSystemConfig()
    const data = res.data.data || {}
    if (data.basic) Object.assign(basicConfig, data.basic)
    if (data.judge) Object.assign(judgeConfig, data.judge)
    if (data.security) Object.assign(securityConfig, data.security)
    if (data.ai) Object.assign(aiConfig, data.ai)
  } catch (e: any) {
    message.error(e.message || t('errors.networkError'))
  } finally {
    loading.value = false
  }
}

async function saveBasic() {
  try {
    await updateSystemConfig({ basic: { ...basicConfig } })
    message.success(t('common.success'))
  } catch (e: any) {
    message.error(e.message || t('errors.networkError'))
  }
}

async function saveJudge() {
  try {
    await updateSystemConfig({ judge: { ...judgeConfig } })
    message.success(t('common.success'))
  } catch (e: any) {
    message.error(e.message || t('errors.networkError'))
  }
}

async function saveSecurity() {
  try {
    await updateSystemConfig({ security: { ...securityConfig } })
    message.success(t('common.success'))
  } catch (e: any) {
    message.error(e.message || t('errors.networkError'))
  }
}

async function saveAI() {
  try {
    await updateSystemConfig({ ai: { ...aiConfig } })
    message.success(t('common.success'))
  } catch (e: any) {
    message.error(e.message || t('errors.networkError'))
  }
}

const passwordComplexityOptions = [
  { label: t('admin.complexityLow'), value: 'low' },
  { label: t('admin.complexityMedium'), value: 'medium' },
  { label: t('admin.complexityHigh'), value: 'high' },
]

const backupFrequencyOptions = [
  { label: t('admin.backupDaily'), value: 'daily' },
  { label: t('admin.backupWeekly'), value: 'weekly' },
  { label: t('admin.backupMonthly'), value: 'monthly' },
]

onMounted(() => {
  fetchConfig()
})
</script>

<template>
  <div class="admin-system-config">
    <div class="page-header">
      <h1 class="page-title">{{ t('admin.systemConfig') }}</h1>
    </div>

    <NCard class="config-card">
      <NTabs v-model:value="activeTab" type="line">
        <!-- Basic Tab -->
        <NTabPane :name="'basic'" :tab="t('admin.tabBasic')">
          <NForm label-placement="left" label-width="160" class="config-form">
            <NFormItem :label="t('admin.systemName')">
              <NInput v-model:value="basicConfig.systemName" :placeholder="t('admin.systemNamePlaceholder')" />
            </NFormItem>
            <NFormItem :label="t('admin.logoUrl')">
              <NInput v-model:value="basicConfig.logoUrl" :placeholder="t('admin.logoUrlPlaceholder')" />
            </NFormItem>
            <NFormItem :label="t('admin.themeColor')">
              <NInput v-model:value="basicConfig.themeColor" placeholder="#00d4ff" />
            </NFormItem>
            <NFormItem :label="t('nav.announcements')">
              <NInput
                v-model:value="basicConfig.announcement"
                type="textarea"
                :rows="4"
                :placeholder="t('admin.systemAnnouncementPlaceholder')"
              />
            </NFormItem>
            <NFormItem label=" ">
              <NButton type="primary" :loading="loading" @click="saveBasic">
                {{ t('common.save') }}
              </NButton>
            </NFormItem>
          </NForm>
        </NTabPane>

        <!-- Judge Tab -->
        <NTabPane :name="'judge'" :tab="t('admin.tabJudge')">
          <NForm label-placement="left" label-width="160" class="config-form">
            <NFormItem :label="t('admin.defaultTimeLimit')">
              <NInputNumber v-model:value="judgeConfig.defaultTimeLimit" :min="100" :max="30000" :step="100" style="width: 300px">
                <template #suffix>ms</template>
              </NInputNumber>
            </NFormItem>
            <NFormItem :label="t('admin.defaultMemoryLimit')">
              <NInputNumber v-model:value="judgeConfig.defaultMemoryLimit" :min="1024" :max="1048576" :step="1024" style="width: 300px">
                <template #suffix>KB</template>
              </NInputNumber>
            </NFormItem>
            <NFormItem :label="t('admin.poolSize')">
              <NInputNumber v-model:value="judgeConfig.poolSize" :min="1" :max="32" style="width: 300px" />
            </NFormItem>
            <NFormItem :label="t('admin.supportedLanguages')">
              <NInput
                v-model:value="judgeConfig.supportedLanguages"
                type="textarea"
                :rows="4"
                placeholder="c, cpp, java, python, go, rust"
              />
            </NFormItem>
            <NFormItem label=" ">
              <NButton type="primary" :loading="loading" @click="saveJudge">
                {{ t('common.save') }}
              </NButton>
            </NFormItem>
          </NForm>
        </NTabPane>

        <!-- Security Tab -->
        <NTabPane :name="'security'" :tab="t('admin.tabSecurity')">
          <NForm label-placement="left" label-width="160" class="config-form">
            <NFormItem :label="t('admin.loginFailLimit')">
              <NInputNumber v-model:value="securityConfig.loginFailLimit" :min="1" :max="20" style="width: 300px" />
            </NFormItem>
            <NFormItem :label="t('admin.passwordComplexity')">
              <NSelect v-model:value="securityConfig.passwordComplexity" :options="passwordComplexityOptions" style="width: 300px" />
            </NFormItem>
            <NFormItem :label="t('admin.backupFrequency')">
              <NSelect v-model:value="securityConfig.backupFrequency" :options="backupFrequencyOptions" style="width: 300px" />
            </NFormItem>
            <NFormItem label=" ">
              <NButton type="primary" :loading="loading" @click="saveSecurity">
                {{ t('common.save') }}
              </NButton>
            </NFormItem>
          </NForm>
        </NTabPane>

        <!-- AI Tab -->
        <NTabPane :name="'ai'" :tab="t('admin.tabAI')">
          <NForm label-placement="left" label-width="160" class="config-form">
            <NFormItem :label="t('admin.aiProvider')">
              <NInput v-model:value="aiConfig.provider" :placeholder="t('admin.aiProviderPlaceholder')" style="width: 300px" />
            </NFormItem>
            <NFormItem :label="t('admin.aiToken')">
              <NInput v-model:value="aiConfig.token" type="password" show-password-on="click" :placeholder="t('admin.aiTokenPlaceholder')" style="width: 300px" />
            </NFormItem>
            <NFormItem label=" ">
              <NButton type="primary" :loading="loading" @click="saveAI">
                {{ t('common.save') }}
              </NButton>
            </NFormItem>
          </NForm>
        </NTabPane>
      </NTabs>
    </NCard>
  </div>
</template>

<style scoped lang="scss">
.admin-system-config {
  max-width: 1000px;
  margin: 0 auto;
}

.page-header {
  margin-bottom: 24px;
}

.page-title {
  font-size: 28px;
  font-weight: 700;
  color: var(--color-text);
}

.config-card {
  background-color: var(--color-bg-card);
  border: 1px solid rgba(255, 255, 255, 0.06);
}

.config-form {
  margin-top: 16px;
}
</style>
