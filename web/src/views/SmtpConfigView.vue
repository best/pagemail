<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { smtpApi } from '@/api/smtp'
import type { SmtpProfile } from '@/types/settings'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Plus, Delete, Edit, Check } from '@element-plus/icons-vue'
import type { FormInstance, FormRules } from 'element-plus'

const profiles = ref<SmtpProfile[]>([])
const loading = ref(false)
const dialogVisible = ref(false)
const isEditing = ref(false)
const formRef = ref<FormInstance>()

const form = ref<Partial<SmtpProfile>>({
  name: '',
  host: '',
  port: 587,
  username: '',
  password: '',
  from_email: '',
  from_name: '',
  use_tls: true,
  is_default: false
})

const rules: FormRules = {
  name: [{ required: true, message: 'Name is required', trigger: 'blur' }],
  host: [{ required: true, message: 'Host is required', trigger: 'blur' }],
  port: [{ required: true, message: 'Port is required', trigger: 'blur' }],
  from_email: [{ required: true, message: 'From Email is required', trigger: 'blur' }]
}

const fetchProfiles = async () => {
  loading.value = true
  try {
    const res = await smtpApi.listProfiles()
    profiles.value = res.data
  } finally {
    loading.value = false
  }
}

const openDialog = (profile?: SmtpProfile) => {
  if (profile) {
    isEditing.value = true
    form.value = { ...profile, password: '' }
  } else {
    isEditing.value = false
    form.value = {
      name: '',
      host: '',
      port: 587,
      username: '',
      password: '',
      from_email: '',
      from_name: '',
      use_tls: true,
      is_default: false
    }
  }
  dialogVisible.value = true
}

const handleSubmit = async () => {
  if (!formRef.value) return
  await formRef.value.validate(async (valid) => {
    if (!valid) return
    try {
      if (isEditing.value && form.value.id) {
        await smtpApi.updateProfile(form.value.id, form.value)
      } else {
        await smtpApi.createProfile(form.value)
      }
      ElMessage.success(isEditing.value ? 'Profile updated' : 'Profile created')
      dialogVisible.value = false
      fetchProfiles()
    } catch {
      // handled globally
    }
  })
}

const handleDelete = (id: string) => {
  ElMessageBox.confirm('Are you sure you want to delete this profile?', 'Warning', {
    type: 'warning'
  }).then(async () => {
    try {
      await smtpApi.deleteProfile(id)
      ElMessage.success('Profile deleted')
      fetchProfiles()
    } catch {
      // handled globally
    }
  }).catch(() => {})
}

const handleTest = async (id: string) => {
  try {
    const { value } = await ElMessageBox.prompt('Enter email address to send test to', 'Test Connection', {
      confirmButtonText: 'Send',
      cancelButtonText: 'Cancel',
      inputPattern: /^[^\s@]+@[^\s@]+\.[^\s@]+$/,
      inputErrorMessage: 'Invalid Email'
    })
    await smtpApi.testProfile(id, value)
    ElMessage.success('Test email sent successfully')
  } catch {
    // cancelled or error
  }
}

onMounted(fetchProfiles)
</script>

<template>
  <div class="smtp-view">
    <div class="header">
      <h2>SMTP Configuration</h2>
      <el-button type="primary" :icon="Plus" @click="openDialog()">Add Profile</el-button>
    </div>

    <el-table :data="profiles" v-loading="loading">
      <el-table-column prop="name" label="Name" />
      <el-table-column prop="host" label="Host" />
      <el-table-column prop="username" label="Username" />
      <el-table-column label="Default" align="center">
        <template #default="{ row }">
          <el-icon v-if="row.is_default" color="var(--el-color-success)"><Check /></el-icon>
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

    <el-dialog v-model="dialogVisible" :title="isEditing ? 'Edit Profile' : 'New Profile'">
      <el-form ref="formRef" :model="form" :rules="rules" label-width="120px">
        <el-form-item label="Profile Name" prop="name">
          <el-input v-model="form.name" />
        </el-form-item>
        <el-form-item label="Host" prop="host">
          <el-input v-model="form.host" />
        </el-form-item>
        <el-form-item label="Port" prop="port">
          <el-input-number v-model="form.port" />
        </el-form-item>
        <el-form-item label="Username" prop="username">
          <el-input v-model="form.username" />
        </el-form-item>
        <el-form-item label="Password" prop="password">
          <el-input
            v-model="form.password"
            type="password"
            show-password
            placeholder="Leave empty to keep current"
          />
        </el-form-item>
        <el-form-item label="From Name" prop="from_name">
          <el-input v-model="form.from_name" />
        </el-form-item>
        <el-form-item label="From Email" prop="from_email">
          <el-input v-model="form.from_email" />
        </el-form-item>
        <el-form-item label="Use TLS" prop="use_tls">
          <el-switch v-model="form.use_tls" />
        </el-form-item>
        <el-form-item label="Set as Default" prop="is_default">
          <el-switch v-model="form.is_default" />
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
.smtp-view {
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
