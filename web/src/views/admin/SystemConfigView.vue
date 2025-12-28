<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { adminApi } from '@/api/admin'
import { ElMessage } from 'element-plus'

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
    ElMessage.success('Configuration saved')
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
    <h2>System Configuration</h2>

    <el-card>
      <el-form :model="config" label-position="top" style="max-width: 600px">
        <el-form-item label="Storage Type">
          <el-radio-group v-model="config.storage_type">
            <el-radio value="local">Local Filesystem</el-radio>
            <el-radio value="s3">S3 / MinIO</el-radio>
          </el-radio-group>
        </el-form-item>

        <el-form-item v-if="config.storage_type === 'local'" label="Storage Path">
          <el-input v-model="config.storage_path" />
        </el-form-item>

        <template v-if="config.storage_type === 's3'">
          <el-form-item label="S3 Bucket">
            <el-input v-model="config.s3_bucket" />
          </el-form-item>
          <el-form-item label="S3 Region">
            <el-input v-model="config.s3_region" />
          </el-form-item>
        </template>

        <el-divider />

        <el-form-item label="Default Output Formats">
          <el-checkbox-group v-model="config.default_formats">
            <el-checkbox value="pdf">PDF</el-checkbox>
            <el-checkbox value="html">HTML</el-checkbox>
            <el-checkbox value="screenshot">Screenshot</el-checkbox>
          </el-checkbox-group>
        </el-form-item>

        <el-form-item label="Max Concurrent Captures">
          <el-input-number v-model="config.max_concurrent_captures" :min="1" :max="20" />
        </el-form-item>

        <el-form-item>
          <el-button type="primary" @click="saveConfig" :loading="saving">Save Configuration</el-button>
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
