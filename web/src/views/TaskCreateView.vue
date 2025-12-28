<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { tasksApi } from '@/api/tasks'
import { smtpApi } from '@/api/smtp'
import { webhooksApi } from '@/api/webhooks'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import type { SmtpProfile, WebhookConfig } from '@/types/settings'
import type { TaskCreatePayload } from '@/types/task'
import type { FormInstance, FormRules } from 'element-plus'

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

const rules: FormRules = {
  url: [
    { required: true, message: 'URL is required', trigger: 'blur' },
    { pattern: /^https?:\/\/.+/, message: 'Please enter a valid URL', trigger: 'blur' }
  ],
  formats: [
    { type: 'array', required: true, message: 'Please select at least one format', trigger: 'change' }
  ]
}

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
      ElMessage.success('Task created successfully')
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
      <h2>New Capture Task</h2>
    </div>

    <el-card>
      <el-form ref="formRef" :model="form" :rules="rules" label-position="top">
        <el-form-item label="Target URL" prop="url">
          <el-input v-model="form.url" placeholder="https://example.com" />
        </el-form-item>

        <el-form-item label="Output Formats" prop="formats">
          <el-checkbox-group v-model="form.formats">
            <el-checkbox value="pdf">PDF</el-checkbox>
            <el-checkbox value="html">HTML</el-checkbox>
            <el-checkbox value="screenshot">Screenshot (PNG)</el-checkbox>
          </el-checkbox-group>
        </el-form-item>

        <el-form-item label="Cookies (Optional)">
          <el-input
            v-model="form.cookies"
            type="textarea"
            :rows="3"
            placeholder="name=value; name2=value2"
          />
          <span class="hint">Format: name=value; name2=value2</span>
        </el-form-item>

        <el-divider>Delivery Options</el-divider>

        <el-form-item label="Delivery Method">
          <el-radio-group v-model="form.delivery_type">
            <el-radio value="none">None (Download only)</el-radio>
            <el-radio value="email">Email</el-radio>
            <el-radio value="webhook">Webhook</el-radio>
          </el-radio-group>
        </el-form-item>

        <el-form-item v-if="form.delivery_type === 'email'" label="Select SMTP Profile">
          <el-select v-model="form.delivery_config_id" placeholder="Select SMTP Profile">
            <el-option
              v-for="profile in smtpProfiles"
              :key="profile.id"
              :label="profile.name"
              :value="profile.id"
            />
          </el-select>
        </el-form-item>

        <el-form-item v-if="form.delivery_type === 'webhook'" label="Select Webhook">
          <el-select v-model="form.delivery_config_id" placeholder="Select Webhook">
            <el-option
              v-for="hook in webhooks"
              :key="hook.id"
              :label="hook.name"
              :value="hook.id"
            />
          </el-select>
        </el-form-item>

        <el-form-item>
          <el-button type="primary" @click="handleSubmit" :loading="loading">Create Task</el-button>
          <el-button @click="router.back()">Cancel</el-button>
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
.hint {
  font-size: 12px;
  color: var(--el-text-color-secondary);
}
</style>
