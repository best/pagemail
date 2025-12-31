<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { tasksApi } from '@/api/tasks'
import { smtpApi } from '@/api/smtp'
import { webhooksApi } from '@/api/webhooks'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import type { SmtpProfile, WebhookConfig } from '@/types/settings'
import type { TaskCreatePayload } from '@/types/task'
import type { FormInstance, FormRules } from 'element-plus'

const { t } = useI18n()
const router = useRouter()
const loading = ref(false)
const smtpProfiles = ref<SmtpProfile[]>([])
const webhooks = ref<WebhookConfig[]>([])
const formRef = ref<FormInstance>()

const form = reactive({
  url: '',
  formats: ['pdf'] as string[],
  cookies: '',
  delivery_type: 'none',
  delivery_config_id: ''
})

const rules = computed<FormRules>(() => ({
  url: [
    { required: true, message: t('validation.urlRequired'), trigger: 'blur' },
    { pattern: /^https?:\/\/.+/, message: t('validation.urlInvalid'), trigger: 'blur' }
  ],
  formats: [
    { type: 'array', required: true, message: t('validation.formatRequired'), trigger: 'change' }
  ]
}))

const fetchConfigs = async () => {
  try {
    const [smtpRes, webhookRes] = await Promise.all([
      smtpApi.listProfiles(),
      webhooksApi.listWebhooks()
    ])
    smtpProfiles.value = smtpRes.data
    webhooks.value = webhookRes.data
  } catch {
    // handled globally
  }
}

onMounted(fetchConfigs)

const handleSubmit = async () => {
  if (!formRef.value) return
  await formRef.value.validate(async (valid) => {
    if (!valid) return

    loading.value = true
    try {
      const payload: TaskCreatePayload = {
        url: form.url,
        formats: form.formats,
        cookies: form.cookies || undefined
      }

      if (form.delivery_type !== 'none' && form.delivery_config_id) {
        payload.delivery_config = {
          type: form.delivery_type as 'email' | 'webhook',
          id: form.delivery_config_id
        }
      }

      const res = await tasksApi.createTask(payload)
      ElMessage.success(t('tasks.createSuccess'))
      router.push({ name: 'task-detail', params: { id: res.data.id } })
    } catch {
      // handled globally
    } finally {
      loading.value = false
    }
  })
}
</script>

<template>
  <div class="create-task">
    <div class="header">
      <h2>{{ t('taskCreate.title') }}</h2>
    </div>

    <el-card>
      <el-form ref="formRef" :model="form" :rules="rules" label-position="top">
        <el-form-item :label="t('taskCreate.targetUrl')" prop="url">
          <el-input v-model="form.url" :placeholder="t('taskCreate.urlPlaceholder')" />
        </el-form-item>

        <el-form-item :label="t('taskCreate.outputFormats')" prop="formats">
          <el-checkbox-group v-model="form.formats">
            <el-checkbox value="pdf">{{ t('taskCreate.pdf') }}</el-checkbox>
            <el-checkbox value="html">{{ t('taskCreate.html') }}</el-checkbox>
            <el-checkbox value="screenshot">{{ t('taskCreate.screenshot') }}</el-checkbox>
          </el-checkbox-group>
        </el-form-item>

        <el-form-item :label="t('taskCreate.cookies')">
          <el-input
            v-model="form.cookies"
            type="textarea"
            :rows="3"
            :placeholder="t('taskCreate.cookiesPlaceholder')"
          />
          <span class="hint">{{ t('taskCreate.cookiesHint') }}</span>
        </el-form-item>

        <el-divider>{{ t('taskCreate.deliveryOptions') }}</el-divider>

        <el-form-item :label="t('taskCreate.deliveryMethod')">
          <el-radio-group v-model="form.delivery_type">
            <el-radio value="none">{{ t('taskCreate.none') }}</el-radio>
            <el-radio value="email">{{ t('taskCreate.emailDelivery') }}</el-radio>
            <el-radio value="webhook">{{ t('taskCreate.webhookDelivery') }}</el-radio>
          </el-radio-group>
        </el-form-item>

        <el-form-item v-if="form.delivery_type === 'email'" :label="t('taskCreate.selectSmtp')">
          <el-select v-model="form.delivery_config_id" :placeholder="t('taskCreate.selectSmtp')">
            <el-option
              v-for="profile in smtpProfiles"
              :key="profile.id"
              :label="profile.name"
              :value="profile.id"
            />
          </el-select>
        </el-form-item>

        <el-form-item v-if="form.delivery_type === 'webhook'" :label="t('taskCreate.selectWebhook')">
          <el-select v-model="form.delivery_config_id" :placeholder="t('taskCreate.selectWebhook')">
            <el-option
              v-for="hook in webhooks"
              :key="hook.id"
              :label="hook.name"
              :value="hook.id"
            />
          </el-select>
        </el-form-item>

        <el-form-item>
          <el-button type="primary" @click="handleSubmit" :loading="loading">{{ t('taskCreate.createTask') }}</el-button>
          <el-button @click="router.back()">{{ t('common.cancel') }}</el-button>
        </el-form-item>
      </el-form>
    </el-card>
  </div>
</template>

<style scoped>
.create-task {
  max-width: 800px;
  margin: 0 auto;
}
.header {
  margin-bottom: 20px;
}
.header h2 {
  margin: 0;
}
.hint {
  font-size: 12px;
  color: var(--el-text-color-secondary);
}
</style>
