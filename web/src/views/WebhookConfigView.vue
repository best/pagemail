<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { webhooksApi } from '@/api/webhooks'
import type { WebhookConfig } from '@/types/settings'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Plus, Delete, Edit } from '@element-plus/icons-vue'
import type { FormInstance, FormRules } from 'element-plus'

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

const rules: FormRules = {
  name: [{ required: true, message: 'Name is required', trigger: 'blur' }],
  url: [
    { required: true, message: 'URL is required', trigger: 'blur' },
    { pattern: /^https?:\/\/.+/, message: 'Must be a valid URL', trigger: 'blur' }
  ]
}

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
      ElMessage.success(isEditing.value ? 'Webhook updated' : 'Webhook created')
      dialogVisible.value = false
      fetchWebhooks()
    } catch {
      // handled globally
    }
  })
}

const handleDelete = (id: string) => {
  ElMessageBox.confirm('Delete this webhook?', 'Warning', { type: 'warning' })
    .then(async () => {
      await webhooksApi.deleteWebhook(id)
      ElMessage.success('Deleted')
      fetchWebhooks()
    })
    .catch(() => {})
}

const handleTest = async (id: string) => {
  try {
    await webhooksApi.testWebhook(id)
    ElMessage.success('Test event sent')
  } catch {
    // handled globally
  }
}

onMounted(fetchWebhooks)
</script>

<template>
  <div class="webhook-view">
    <div class="header">
      <h2>Webhooks</h2>
      <el-button type="primary" :icon="Plus" @click="openDialog()">Add Webhook</el-button>
    </div>

    <el-card shadow="hover" class="pm-table-card">
      <el-table :data="webhooks" v-loading="loading" stripe>
      <el-table-column prop="name" label="Name" />
      <el-table-column prop="url" label="URL" show-overflow-tooltip />
      <el-table-column label="Status" width="100">
        <template #default="{ row }">
          <el-tag :type="row.is_active ? 'success' : 'info'">
            {{ row.is_active ? 'Active' : 'Inactive' }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column label="Actions" align="right" width="250">
        <template #default="{ row }">
          <el-button size="small" @click="handleTest(row.id)">Test</el-button>
          <el-button size="small" :icon="Edit" @click="openDialog(row)" />
          <el-button size="small" type="danger" :icon="Delete" @click="handleDelete(row.id)" />
        </template>
      </el-table-column>
    </el-table>
    </el-card>

    <el-dialog v-model="dialogVisible" :title="isEditing ? 'Edit Webhook' : 'New Webhook'">
      <el-form ref="formRef" :model="form" :rules="rules" label-width="120px">
        <el-form-item label="Name" prop="name">
          <el-input v-model="form.name" />
        </el-form-item>
        <el-form-item label="Payload URL" prop="url">
          <el-input v-model="form.url" placeholder="https://..." />
        </el-form-item>
        <el-form-item label="Secret" prop="secret">
          <el-input v-model="form.secret" type="password" show-password />
        </el-form-item>
        <el-form-item label="Active" prop="is_active">
          <el-switch v-model="form.is_active" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">Cancel</el-button>
        <el-button type="primary" @click="handleSubmit">Save</el-button>
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
