<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { webhooksApi } from '@/api/webhooks'
import type { WebhookConfig } from '@/types/settings'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Plus, Delete, Edit } from '@element-plus/icons-vue'
import type { FormInstance, FormRules } from 'element-plus'

const { t } = useI18n()

const webhooks = ref<WebhookConfig[]>([])
const loading = ref(false)
const dialogVisible = ref(false)
const isEditing = ref(false)
const formRef = ref<FormInstance>()

const form = ref<Partial<WebhookConfig>>({
  name: '',
  url: '',
  secret: '',
  is_active: true
})

const rules = computed<FormRules>(() => ({
  name: [{ required: true, message: t('webhook.nameRequired'), trigger: 'blur' }],
  url: [
    { required: true, message: t('webhook.urlRequired'), trigger: 'blur' },
    { pattern: /^https?:\/\/.+/, message: t('webhook.urlInvalid'), trigger: 'blur' }
  ]
}))

const fetchWebhooks = async () => {
  loading.value = true
  try {
    const res = await webhooksApi.listWebhooks()
    webhooks.value = res.data
  } finally {
    loading.value = false
  }
}

const openDialog = (webhook?: WebhookConfig) => {
  if (webhook) {
    isEditing.value = true
    form.value = { ...webhook }
  } else {
    isEditing.value = false
    form.value = { name: '', url: '', secret: '', is_active: true }
  }
  dialogVisible.value = true
}

const handleSubmit = async () => {
  if (!formRef.value) return
  await formRef.value.validate(async (valid) => {
    if (!valid) return
    try {
      if (isEditing.value && form.value.id) {
        await webhooksApi.updateWebhook(form.value.id, form.value)
      } else {
        await webhooksApi.createWebhook(form.value)
      }
      ElMessage.success(isEditing.value ? t('webhook.webhookUpdated') : t('webhook.webhookCreated'))
      dialogVisible.value = false
      fetchWebhooks()
    } catch {
      // handled globally
    }
  })
}

const handleDelete = (id: string) => {
  ElMessageBox.confirm(t('webhook.deleteConfirm'), 'Warning', { type: 'warning' })
    .then(async () => {
      await webhooksApi.deleteWebhook(id)
      ElMessage.success(t('webhook.webhookDeleted'))
      fetchWebhooks()
    })
    .catch(() => {})
}

const handleTest = async (id: string) => {
  try {
    await webhooksApi.testWebhook(id)
    ElMessage.success(t('webhook.testSuccess'))
  } catch {
    // handled globally
  }
}

onMounted(fetchWebhooks)
</script>

<template>
  <div class="webhook-view">
    <div class="header">
      <h2>{{ t('webhook.title') }}</h2>
      <el-button type="primary" :icon="Plus" @click="openDialog()">{{ t('webhook.addWebhook') }}</el-button>
    </div>

    <el-card shadow="hover" class="pm-table-card">
      <el-table :data="webhooks" v-loading="loading" stripe>
      <el-table-column prop="name" :label="t('webhook.name')" />
      <el-table-column prop="url" :label="t('webhook.url')" show-overflow-tooltip />
      <el-table-column :label="t('webhook.status')" width="100">
        <template #default="{ row }">
          <el-tag :type="row.is_active ? 'success' : 'info'">
            {{ row.is_active ? t('webhook.active') : t('webhook.inactive') }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column :label="t('webhook.actions')" align="right" width="250">
        <template #default="{ row }">
          <el-button size="small" @click="handleTest(row.id)">{{ t('webhook.test') }}</el-button>
          <el-button size="small" :icon="Edit" @click="openDialog(row)" />
          <el-button size="small" type="danger" :icon="Delete" @click="handleDelete(row.id)" />
        </template>
      </el-table-column>
    </el-table>
    </el-card>

    <el-dialog v-model="dialogVisible" :title="isEditing ? t('webhook.editWebhook') : t('webhook.newWebhook')">
      <el-form ref="formRef" :model="form" :rules="rules" label-width="120px">
        <el-form-item :label="t('webhook.name')" prop="name">
          <el-input v-model="form.name" />
        </el-form-item>
        <el-form-item :label="t('webhook.payloadUrl')" prop="url">
          <el-input v-model="form.url" placeholder="https://..." />
        </el-form-item>
        <el-form-item :label="t('webhook.secret')" prop="secret">
          <el-input v-model="form.secret" type="password" show-password />
        </el-form-item>
        <el-form-item :label="t('webhook.isActive')" prop="is_active">
          <el-switch v-model="form.is_active" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">{{ t('common.cancel') }}</el-button>
        <el-button type="primary" @click="handleSubmit">{{ t('common.save') }}</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<style scoped>
.webhook-view {
  max-width: 1200px;
  margin: 0 auto;
}
.header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}
</style>
