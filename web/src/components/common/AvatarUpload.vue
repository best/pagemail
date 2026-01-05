<script setup lang="ts">
import { ref, computed, onBeforeUnmount, onMounted, watch } from 'vue'
import { Plus, Check, Close } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import type { UploadFile } from 'element-plus'
import { useI18n } from 'vue-i18n'
import { userApi } from '@/api/user'
import { useAuthStore } from '@/stores/auth'
import apiClient from '@/api/client'

const { t } = useI18n()
const authStore = useAuthStore()

const loading = ref(false)
const previewUrl = ref<string>('')
const avatarBlobUrl = ref<string>('')
const selectedFile = ref<File | null>(null)

const fetchAvatar = async () => {
  if (!authStore.user?.avatar_url) return
  try {
    const response = await apiClient.get(`${authStore.user.avatar_url}`, { responseType: 'blob' })
    if (avatarBlobUrl.value) URL.revokeObjectURL(avatarBlobUrl.value)
    avatarBlobUrl.value = URL.createObjectURL(response.data)
  } catch {
    avatarBlobUrl.value = ''
  }
}

onMounted(() => fetchAvatar())
watch(() => authStore.user?.avatar_url, () => fetchAvatar())

const currentAvatar = computed(() => {
  if (previewUrl.value) return previewUrl.value
  return avatarBlobUrl.value
})

const userInitial = computed(() => authStore.user?.email?.[0]?.toUpperCase() || '?')

const handleFileChange = (file: UploadFile) => {
  if (!file.raw) return

  const allowedTypes = ['image/jpeg', 'image/png', 'image/webp']
  const isAllowedType = allowedTypes.includes(file.raw.type)
  const maxBytes = 5 * 1024 * 1024
  const isWithinSize = file.raw.size <= maxBytes

  if (!isAllowedType) {
    ElMessage.error(t('settings.avatarTypeErr'))
    return
  }
  if (!isWithinSize) {
    ElMessage.error(t('settings.avatarSizeErr'))
    return
  }

  if (previewUrl.value) {
    URL.revokeObjectURL(previewUrl.value)
  }
  previewUrl.value = URL.createObjectURL(file.raw)
  selectedFile.value = file.raw
}

const cancelUpload = () => {
  if (previewUrl.value) {
    URL.revokeObjectURL(previewUrl.value)
  }
  previewUrl.value = ''
  selectedFile.value = null
}

onBeforeUnmount(() => {
  if (previewUrl.value) URL.revokeObjectURL(previewUrl.value)
  if (avatarBlobUrl.value) URL.revokeObjectURL(avatarBlobUrl.value)
})

const uploadAvatar = async () => {
  if (!selectedFile.value) return

  loading.value = true
  try {
    const response = await userApi.uploadAvatar(selectedFile.value)
    authStore.user = response.data
    await fetchAvatar()
    cancelUpload()
    ElMessage.success(t('settings.avatarUpdated'))
  } catch {
    // handled globally
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <div class="avatar-upload-container">
    <el-upload
      class="avatar-uploader"
      action="#"
      :show-file-list="false"
      :auto-upload="false"
      :on-change="handleFileChange"
      accept="image/jpeg,image/png,image/webp"
    >
      <div class="avatar-wrapper">
        <el-avatar :size="100" :src="currentAvatar" class="user-avatar">
          {{ userInitial }}
        </el-avatar>
        <div class="avatar-overlay">
          <el-icon class="upload-icon"><Plus /></el-icon>
          <span class="upload-text">{{ t('settings.changeAvatar') }}</span>
        </div>
      </div>
    </el-upload>

    <div v-if="selectedFile" class="avatar-actions">
      <el-button :icon="Close" circle size="small" @click="cancelUpload" />
      <el-button :icon="Check" :loading="loading" circle size="small" type="primary" @click="uploadAvatar" />
    </div>
  </div>
</template>

<style scoped>
.avatar-upload-container {
  display: flex;
  flex-direction: column;
  align-items: flex-start;
  margin-bottom: 24px;
}

.avatar-uploader :deep(.el-upload) {
  border-radius: 50%;
  cursor: pointer;
  overflow: hidden;
}

.avatar-wrapper {
  position: relative;
  cursor: pointer;
}

.user-avatar {
  background: var(--el-color-primary-light-7);
  color: var(--el-color-primary);
  font-size: 36px;
  font-weight: 500;
}

.avatar-overlay {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background: rgba(0, 0, 0, 0.5);
  color: white;
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
  opacity: 0;
  transition: opacity 0.3s;
  border-radius: 50%;
}

.avatar-wrapper:hover .avatar-overlay {
  opacity: 1;
}

.upload-icon {
  font-size: 24px;
  margin-bottom: 4px;
}

.upload-text {
  font-size: 12px;
}

.avatar-actions {
  margin-top: 12px;
  display: flex;
  gap: 8px;
}
</style>
