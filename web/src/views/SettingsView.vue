<script setup lang="ts">
import { ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { useAuthStore } from '@/stores/auth'
import { ElMessage } from 'element-plus'
import apiClient from '@/api/client'

const { t } = useI18n()
const authStore = useAuthStore()
const passwordForm = ref({
  current_password: '',
  new_password: '',
  confirm_password: ''
})
const loading = ref(false)

const updatePassword = async () => {
  if (passwordForm.value.new_password !== passwordForm.value.confirm_password) {
    ElMessage.error(t('settings.passwordMismatch'))
    return
  }

  loading.value = true
  try {
    await apiClient.put('/users/me/password', {
      current_password: passwordForm.value.current_password,
      new_password: passwordForm.value.new_password
    })
    ElMessage.success(t('settings.passwordUpdated'))
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
    <div class="page-header">
      <h2>{{ t('settings.title') }}</h2>
    </div>

    <el-tabs type="border-card">
      <el-tab-pane :label="t('settings.profile')">
        <el-form label-position="top" style="max-width: 400px">
          <el-form-item :label="t('auth.email')">
            <el-input :model-value="authStore.user?.email" disabled />
          </el-form-item>
        </el-form>

        <el-divider />

        <h3>{{ t('settings.changePassword') }}</h3>
        <el-form
          :model="passwordForm"
          label-position="top"
          style="max-width: 400px"
          @submit.prevent="updatePassword"
        >
          <el-form-item :label="t('settings.currentPassword')">
            <el-input v-model="passwordForm.current_password" type="password" show-password />
          </el-form-item>
          <el-form-item :label="t('settings.newPassword')">
            <el-input v-model="passwordForm.new_password" type="password" show-password />
          </el-form-item>
          <el-form-item :label="t('settings.confirmNewPassword')">
            <el-input v-model="passwordForm.confirm_password" type="password" show-password />
          </el-form-item>
          <el-button type="primary" native-type="submit" :loading="loading">
            {{ t('settings.updatePassword') }}
          </el-button>
        </el-form>
      </el-tab-pane>

      <el-tab-pane :label="t('settings.smtpConfig')">
        <div class="config-link">
          <p>{{ t('settings.smtpDesc') }}</p>
          <el-button type="primary" @click="$router.push('/settings/smtp')">
            {{ t('settings.goToSmtp') }}
          </el-button>
        </div>
      </el-tab-pane>

      <el-tab-pane :label="t('settings.webhooks')">
        <div class="config-link">
          <p>{{ t('settings.webhooksDesc') }}</p>
          <el-button type="primary" @click="$router.push('/settings/webhooks')">
            {{ t('settings.goToWebhooks') }}
          </el-button>
        </div>
      </el-tab-pane>
    </el-tabs>
  </div>
</template>

<style scoped>
.settings-view {
  max-width: 1200px;
  margin: 0 auto;
}
.page-header {
  margin-bottom: 20px;
}
.page-header h2 {
  margin: 0;
}
.config-link {
  padding: 20px;
  text-align: center;
}
</style>
