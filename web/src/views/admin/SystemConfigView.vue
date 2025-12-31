<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { adminApi } from '@/api/admin'
import { useSiteConfigStore } from '@/stores/siteConfig'
import { ElMessage } from 'element-plus'

const { t } = useI18n()
const siteConfigStore = useSiteConfigStore()

interface SystemConfig {
  storage_type: string
  storage_path: string
  s3_bucket?: string
  s3_region?: string
  default_formats: string[]
  max_concurrent_captures: number
}

const activeTab = ref('system')
const config = ref<SystemConfig>({
  storage_type: 'local',
  storage_path: './data/captures',
  default_formats: ['pdf'],
  max_concurrent_captures: 5
})
const siteSettings = ref({
  site_name: '',
  site_slogan: ''
})
const loading = ref(false)
const saving = ref(false)

const fetchConfig = async () => {
  loading.value = true
  try {
    const [sysRes, siteRes] = await Promise.all([
      adminApi.getSystemConfig(),
      adminApi.getSiteConfig()
    ])
    config.value = sysRes.data
    siteSettings.value = siteRes.data
  } catch {
    // handled globally
  } finally {
    loading.value = false
  }
}

const saveConfig = async () => {
  saving.value = true
  try {
    if (activeTab.value === 'system') {
      await adminApi.updateSystemConfig(config.value)
    } else {
      await adminApi.updateSiteConfig(siteSettings.value)
      siteConfigStore.updateConfig(siteSettings.value.site_name, siteSettings.value.site_slogan)
    }
    ElMessage.success(t('admin.siteConfig.configSaved'))
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
    <div class="page-header">
      <h2>{{ t(activeTab === 'system' ? 'admin.systemConfig.title' : 'admin.siteConfig.title') }}</h2>
    </div>

    <el-card>
      <el-tabs v-model="activeTab">
        <el-tab-pane :label="t('admin.systemConfig.title')" name="system">
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
              <el-button type="primary" @click="saveConfig" :loading="saving">{{ t('admin.siteConfig.saveConfig') }}</el-button>
            </el-form-item>
          </el-form>
        </el-tab-pane>

        <el-tab-pane :label="t('admin.siteConfig.title')" name="site">
          <el-form :model="siteSettings" label-position="top" style="max-width: 600px">
            <el-form-item :label="t('admin.siteConfig.siteName')" required>
              <el-input v-model="siteSettings.site_name" placeholder="Pagemail" />
            </el-form-item>
            <el-form-item :label="t('admin.siteConfig.siteSlogan')">
              <el-input v-model="siteSettings.site_slogan" :placeholder="t('admin.siteConfig.sloganHint')" />
              <div class="hint">{{ t('admin.siteConfig.sloganHint') }}</div>
            </el-form-item>
            <el-form-item>
              <el-button type="primary" @click="saveConfig" :loading="saving">{{ t('admin.siteConfig.saveConfig') }}</el-button>
            </el-form-item>
          </el-form>
        </el-tab-pane>
      </el-tabs>
    </el-card>
  </div>
</template>

<style scoped>
.system-config {
  max-width: 1200px;
  margin: 0 auto;
}
.page-header {
  margin-bottom: 20px;
}
.page-header h2 {
  margin: 0;
}
.hint {
  font-size: 12px;
  color: var(--el-text-color-secondary);
  margin-top: 4px;
}
</style>
