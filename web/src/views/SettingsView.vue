<script setup lang="ts">
import { ref } from 'vue'
import { useAuthStore } from '@/stores/auth'
import { ElMessage } from 'element-plus'
import apiClient from '@/api/client'

const authStore = useAuthStore()
const passwordForm = ref({
  current_password: '',
  new_password: '',
  confirm_password: ''
})
const loading = ref(false)

const updatePassword = async () => {
  if (passwordForm.value.new_password !== passwordForm.value.confirm_password) {
    ElMessage.error('New passwords do not match')
    return
  }

  loading.value = true
  try {
    await apiClient.put('/users/me/password', {
      current_password: passwordForm.value.current_password,
      new_password: passwordForm.value.new_password
    })
    ElMessage.success('Password updated successfully')
    passwordForm.value = { current_password: '', new_password: '', confirm_password: '' }
  } catch {
    // handled globally
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <div class="settings-view">
    <div class="header">
      <h2>Settings</h2>
    </div>

    <el-tabs type="border-card">
      <el-tab-pane label="Profile">
        <el-form label-position="top" style="max-width: 400px">
          <el-form-item label="Email">
            <el-input :model-value="authStore.user?.email" disabled />
          </el-form-item>
        </el-form>

        <el-divider />

        <h3>Change Password</h3>
        <el-form
          :model="passwordForm"
          label-position="top"
          style="max-width: 400px"
          @submit.prevent="updatePassword"
        >
          <el-form-item label="Current Password">
            <el-input v-model="passwordForm.current_password" type="password" show-password />
          </el-form-item>
          <el-form-item label="New Password">
            <el-input v-model="passwordForm.new_password" type="password" show-password />
          </el-form-item>
          <el-form-item label="Confirm New Password">
            <el-input v-model="passwordForm.confirm_password" type="password" show-password />
          </el-form-item>
          <el-button type="primary" native-type="submit" :loading="loading">
            Update Password
          </el-button>
        </el-form>
      </el-tab-pane>

      <el-tab-pane label="SMTP Configuration">
        <div class="config-link">
          <p>Manage your SMTP profiles for email delivery.</p>
          <el-button type="primary" @click="$router.push('/settings/smtp')">
            Go to SMTP Settings
          </el-button>
        </div>
      </el-tab-pane>

      <el-tab-pane label="Webhooks">
        <div class="config-link">
          <p>Manage webhooks for notifications.</p>
          <el-button type="primary" @click="$router.push('/settings/webhooks')">
            Go to Webhook Settings
          </el-button>
        </div>
      </el-tab-pane>
    </el-tabs>
  </div>
</template>

<style scoped>
.settings-view {
  max-width: 1000px;
  margin: 0 auto;
}
.config-link {
  padding: 20px;
  text-align: center;
}
</style>
