<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { smtpApi } from '@/api/smtp'
import type { SmtpProfile } from '@/types/settings'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Plus, Delete, Edit, Check } from '@element-plus/icons-vue'
import type { FormInstance, FormRules } from 'element-plus'

const { t } = useI18n()

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

const rules = computed<FormRules>(() => ({
  name: [{ required: true, message: t('webhook.nameRequired'), trigger: 'blur' }],
  host: [{ required: true, message: t('validation.urlRequired'), trigger: 'blur' }],
  port: [{ required: true, message: t('validation.urlRequired'), trigger: 'blur' }],
  from_email: [{ required: true, message: t('validation.emailRequired'), trigger: 'blur' }]
}))

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
      ElMessage.success(isEditing.value ? t('smtp.profileUpdated') : t('smtp.profileCreated'))
      dialogVisible.value = false
      fetchProfiles()
    } catch {
      // handled globally
    }
  })
}

const handleDelete = (id: string) => {
  ElMessageBox.confirm(t('smtp.deleteConfirm'), 'Warning', {
    type: 'warning'
  }).then(async () => {
    try {
      await smtpApi.deleteProfile(id)
      ElMessage.success(t('smtp.profileDeleted'))
      fetchProfiles()
    } catch {
      // handled globally
    }
  }).catch(() => {})
}

const handleTest = async (id: string) => {
  try {
    const { value } = await ElMessageBox.prompt(t('smtp.testPrompt'), t('smtp.testTitle'), {
      confirmButtonText: t('common.confirm'),
      cancelButtonText: t('common.cancel'),
      inputPattern: /^[^\s@]+@[^\s@]+\.[^\s@]+$/,
      inputErrorMessage: t('smtp.invalidEmail')
    })
    await smtpApi.testProfile(id, value)
    ElMessage.success(t('smtp.testSuccess'))
  } catch {
    // cancelled or error
  }
}

onMounted(fetchProfiles)
</script>

<template>
  <div class="smtp-view">
    <div class="header">
      <h2>{{ t('smtp.title') }}</h2>
      <el-button type="primary" :icon="Plus" @click="openDialog()">{{ t('smtp.addProfile') }}</el-button>
    </div>

    <el-card shadow="hover" class="pm-table-card">
      <el-table :data="profiles" v-loading="loading" stripe>
      <el-table-column prop="name" :label="t('smtp.name')" />
      <el-table-column prop="host" :label="t('smtp.host')" />
      <el-table-column prop="username" :label="t('smtp.username')" />
      <el-table-column :label="t('smtp.default')" align="center">
        <template #default="{ row }">
          <el-icon v-if="row.is_default" color="var(--el-color-success)"><Check /></el-icon>
        </template>
      </el-table-column>
      <el-table-column :label="t('smtp.actions')" align="right" width="250">
        <template #default="{ row }">
          <el-button size="small" @click="handleTest(row.id)">{{ t('smtp.test') }}</el-button>
          <el-button size="small" :icon="Edit" @click="openDialog(row)" />
          <el-button size="small" type="danger" :icon="Delete" @click="handleDelete(row.id)" />
        </template>
      </el-table-column>
    </el-table>
    </el-card>

    <el-dialog v-model="dialogVisible" :title="isEditing ? t('smtp.editProfile') : t('smtp.newProfile')">
      <el-form ref="formRef" :model="form" :rules="rules" label-width="120px">
        <el-form-item :label="t('smtp.profileName')" prop="name">
          <el-input v-model="form.name" />
        </el-form-item>
        <el-form-item :label="t('smtp.host')" prop="host">
          <el-input v-model="form.host" />
        </el-form-item>
        <el-form-item :label="t('smtp.port')" prop="port">
          <el-input-number v-model="form.port" />
        </el-form-item>
        <el-form-item :label="t('smtp.username')" prop="username">
          <el-input v-model="form.username" />
        </el-form-item>
        <el-form-item :label="t('smtp.password')" prop="password">
          <el-input
            v-model="form.password"
            type="password"
            show-password
            :placeholder="t('smtp.passwordPlaceholder')"
          />
        </el-form-item>
        <el-form-item :label="t('smtp.fromName')" prop="from_name">
          <el-input v-model="form.from_name" />
        </el-form-item>
        <el-form-item :label="t('smtp.fromEmail')" prop="from_email">
          <el-input v-model="form.from_email" />
        </el-form-item>
        <el-form-item :label="t('smtp.useTls')" prop="use_tls">
          <el-switch v-model="form.use_tls" />
        </el-form-item>
        <el-form-item :label="t('smtp.setDefault')" prop="is_default">
          <el-switch v-model="form.is_default" />
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
.header h2 {
  margin: 0;
}
</style>
