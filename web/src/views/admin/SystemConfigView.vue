<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { adminApi } from '@/api/admin'
import { ElMessage } from 'element-plus'

const { t } = useI18n()

interface SystemConfig {
  storage_type: string
  storage_path: string
  s3_bucket?: string
  s3_region?: string
  default_formats: string[]
  max_concurrent_captures: number
}

const config = ref<SystemConfig>({
  storage_type: 'local',
  storage_path: './data/captures',
  default_formats: ['pdf'],
  max_concurrent_captures: 5
})
const loading = ref(false)
const saving = ref(false)

const fetchConfig = async () => {
  loading.value = true
  try {
    const res = await adminApi.getSystemConfig()
    config.value = res.data
  } catch {
    // handled globally
  } finally {
    loading.value = false
  }
}

const saveConfig = async () => {
  saving.value = true
  try {
    await adminApi.updateSystemConfig(config.value)
    ElMessage.success(t('admin.systemConfig.configSaved'))
  } catch {
    // handled globally
  } finally {
    saving.value = false
  }
}

onMounted(fetchConfig)
</script>

<template>
  <div class="system-config" v-loading="loading">
    <h2>{{ t('admin.systemConfig.title') }}</h2>

    <el-card>
      <el-form :model="config" label-position="top" style="max-width: 600px">
        <el-form-item :label="t('admin.systemConfig.storageType')">
          <el-radio-group v-model="config.storage_type">
            <el-radio value="local">{{ t('admin.systemConfig.local') }}</el-radio>
            <el-radio value="s3">{{ t('admin.systemConfig.s3') }}</el-radio>
          </el-radio-group>
        </el-form-item>

        <el-form-item v-if="config.storage_type === 'local'" :label="t('admin.systemConfig.storagePath')">
          <el-input v-model="config.storage_path" />
        </el-form-item>

        <template v-if="config.storage_type === 's3'">
          <el-form-item :label="t('admin.systemConfig.s3Bucket')">
            <el-input v-model="config.s3_bucket" />
          </el-form-item>
          <el-form-item :label="t('admin.systemConfig.s3Region')">
            <el-input v-model="config.s3_region" />
          </el-form-item>
        </template>

        <el-divider />

        <el-form-item :label="t('admin.systemConfig.defaultFormats')">
          <el-checkbox-group v-model="config.default_formats">
            <el-checkbox value="pdf">PDF</el-checkbox>
            <el-checkbox value="html">HTML</el-checkbox>
            <el-checkbox value="screenshot">Screenshot</el-checkbox>
          </el-checkbox-group>
        </el-form-item>

        <el-form-item :label="t('admin.systemConfig.maxConcurrent')">
          <el-input-number v-model="config.max_concurrent_captures" :min="1" :max="20" />
        </el-form-item>

        <el-form-item>
          <el-button type="primary" @click="saveConfig" :loading="saving">{{ t('admin.systemConfig.saveConfig') }}</el-button>
        </el-form-item>
      </el-form>
    </el-card>
  </div>
</template>

<style scoped>
.system-config {
  padding: 20px;
  max-width: 800px;
}
</style>
