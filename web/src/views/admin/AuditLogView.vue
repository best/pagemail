<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { Search, Refresh, Download, View as ViewIcon } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import { adminApi, type AuditLog } from '@/api/admin'

const { t } = useI18n()

const logs = ref<AuditLog[]>([])
const loading = ref(false)
const total = ref(0)
const drawerVisible = ref(false)
const selectedLog = ref<AuditLog | null>(null)

const query = reactive({
  page: 1,
  limit: 20,
  action: '',
  actor: '',
  dateRange: [] as string[]
})

const actionOptions = [
  'user.login',
  'user.create',
  'user.update',
  'user.delete',
  'settings.update',
  'smtp.create',
  'smtp.update',
  'smtp.delete',
  'webhook.create',
  'webhook.update',
  'webhook.delete',
  'capture.create',
  'capture.delete'
]

const fetchLogs = async () => {
  loading.value = true
  try {
    const params: Record<string, unknown> = {
      page: query.page,
      limit: query.limit
    }
    if (query.action) params.action = query.action
    if (query.actor) params.actor = query.actor
    if (query.dateRange?.length === 2) {
      params.from = query.dateRange[0]
      params.to = query.dateRange[1]
    }
    const res = await adminApi.getAuditLogs(params as Parameters<typeof adminApi.getAuditLogs>[0])
    logs.value = res.data.data
    total.value = res.data.meta.total
  } catch {
    ElMessage.error(t('common.fetchError'))
  } finally {
    loading.value = false
  }
}

const handleSearch = () => {
  query.page = 1
  fetchLogs()
}

const handleReset = () => {
  query.action = ''
  query.actor = ''
  query.dateRange = []
  handleSearch()
}

const viewDetails = (row: AuditLog) => {
  selectedLog.value = row
  drawerVisible.value = true
}

const getActionType = (action: string): 'primary' | 'success' | 'warning' | 'info' | 'danger' => {
  if (action.includes('login')) return 'info'
  if (action.includes('create')) return 'success'
  if (action.includes('update')) return 'warning'
  if (action.includes('delete')) return 'danger'
  return 'primary'
}

const translateAction = (action: string): string => {
  const key = `admin.auditLogs.actionTypes.${action.replace(/\./g, '_')}`
  const translated = t(key)
  return translated === key ? action : translated
}

const handleExport = () => {
  const dataStr = JSON.stringify(logs.value, null, 2)
  const blob = new Blob([dataStr], { type: 'application/json' })
  const url = URL.createObjectURL(blob)
  const link = document.createElement('a')
  link.href = url
  link.download = `audit-logs-${new Date().toISOString().slice(0, 10)}.json`
  document.body.appendChild(link)
  link.click()
  document.body.removeChild(link)
  URL.revokeObjectURL(url)
}

onMounted(fetchLogs)
</script>

<template>
  <div class="audit-log">
    <div class="header anim-up">
      <div class="title-group">
        <h2>{{ t('admin.auditLogs.title') }}</h2>
      </div>
      <el-button :icon="Download" @click="handleExport">{{ t('common.export') }}</el-button>
    </div>

    <el-card shadow="never" class="filter-card anim-up d1">
      <div class="filters-row">
        <el-select v-model="query.action" clearable :placeholder="t('admin.auditLogs.action')" style="width: 150px">
          <el-option v-for="action in actionOptions" :key="action" :label="translateAction(action)" :value="action" />
        </el-select>
        <el-input v-model="query.actor" :placeholder="t('admin.auditLogs.user')" clearable style="width: 150px" @keyup.enter="handleSearch" />
        <el-date-picker
          v-model="query.dateRange"
          type="daterange"
          range-separator="-"
          :start-placeholder="t('common.startDate')"
          :end-placeholder="t('common.endDate')"
          value-format="YYYY-MM-DD"
          style="width: 240px"
        />
        <el-button type="primary" :icon="Search" @click="handleSearch">{{ t('common.search') }}</el-button>
        <el-button :icon="Refresh" @click="handleReset">{{ t('common.reset') }}</el-button>
      </div>
    </el-card>

    <el-card shadow="hover" class="pm-table-card anim-up d2">
      <el-table :data="logs" v-loading="loading" stripe style="width: 100%">
        <el-table-column prop="created_at" :label="t('admin.auditLogs.time')" width="180">
          <template #default="{ row }">{{ new Date(row.created_at).toLocaleString() }}</template>
        </el-table-column>
        <el-table-column prop="actor_email" :label="t('admin.auditLogs.user')" min-width="200" />
        <el-table-column prop="action" :label="t('admin.auditLogs.action')" width="140">
          <template #default="{ row }">
            <el-tag :type="getActionType(row.action)" size="small" effect="light" round>{{ translateAction(row.action) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="resource_type" :label="t('admin.auditLogs.resource')" width="100" />
        <el-table-column prop="ip_address" :label="t('admin.auditLogs.ip')" width="130" />
        <el-table-column :label="t('common.actions')" width="80" align="center">
          <template #default="{ row }">
            <el-tooltip :content="t('admin.auditLogs.viewDetails')" placement="top">
              <el-button link type="primary" :icon="ViewIcon" :aria-label="t('admin.auditLogs.viewDetails')" @click="viewDetails(row)" />
            </el-tooltip>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <div class="pagination">
      <el-pagination
        background
        layout="prev, pager, next"
        :total="total"
        :page-size="query.limit"
        v-model:current-page="query.page"
        @current-change="fetchLogs"
      />
    </div>

    <el-drawer v-model="drawerVisible" :title="t('admin.auditLogs.details')" size="450px">
      <template v-if="selectedLog">
        <el-descriptions :column="1" border class="detail-descriptions">
          <el-descriptions-item :label="t('admin.auditLogs.action')">
            <el-tag :type="getActionType(selectedLog.action)" effect="light" round>{{ translateAction(selectedLog.action) }}</el-tag>
          </el-descriptions-item>
          <el-descriptions-item :label="t('admin.auditLogs.time')">
            {{ new Date(selectedLog.created_at).toLocaleString() }}
          </el-descriptions-item>
          <el-descriptions-item :label="t('admin.auditLogs.user')">
            {{ selectedLog.actor_email }}
          </el-descriptions-item>
          <el-descriptions-item :label="t('admin.auditLogs.ip')">
            {{ selectedLog.ip_address }}
          </el-descriptions-item>
        </el-descriptions>

        <div v-if="selectedLog.details" class="detail-section">
          <h4>{{ t('admin.auditLogs.changeDetails') }}</h4>

          <div v-if="selectedLog.action === 'user.login'" class="detail-content">
            <p class="login-success">{{ t('admin.auditLogs.loginSuccess') }}</p>
            <p v-if="selectedLog.details.user_agent" class="user-agent">
              <strong>{{ t('admin.auditLogs.userAgent') }}:</strong><br>
              {{ selectedLog.details.user_agent }}
            </p>
          </div>

          <div v-else-if="selectedLog.details_type === 'change' && selectedLog.details.changes" class="detail-content">
            <div v-for="(change, idx) in (selectedLog.details.changes as Array<{field: string, old_value?: string, new_value?: string}>)" :key="idx" class="change-item">
              <strong>{{ change.field }}:</strong>
              <span v-if="change.old_value != null" class="old-value">{{ change.old_value }}</span>
              <span v-if="change.old_value != null && change.new_value != null"> → </span>
              <span v-if="change.new_value != null" class="new-value">{{ change.new_value }}</span>
            </div>
          </div>

          <div v-else class="detail-content">
            <div v-for="(val, key) in selectedLog.details" :key="key" class="detail-item">
              <strong>{{ key }}:</strong>
              <span v-if="Array.isArray(val)">{{ val.join(', ') }}</span>
              <span v-else-if="typeof val === 'object'">{{ JSON.stringify(val) }}</span>
              <span v-else>{{ val }}</span>
            </div>
          </div>
        </div>

        <div v-else class="no-details">
          {{ t('admin.auditLogs.noDetails') }}
        </div>
      </template>
    </el-drawer>
  </div>
</template>

<style scoped>
/* Entry Animations */
.anim-up {
  opacity: 0;
  animation: slideUp 0.5s ease-out forwards;
}
.d1 { animation-delay: 0.1s; }
.d2 { animation-delay: 0.2s; }

@keyframes slideUp {
  from { opacity: 0; transform: translateY(20px); }
  to { opacity: 1; transform: translateY(0); }
}

/* Layout */
.audit-log {
  max-width: 1200px;
  margin: 0 auto;
}

/* Header */
.header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.title-group h2 {
  margin: 0;
}

/* Filter Card */
.filter-card {
  margin-bottom: 20px;
}

.filters-row {
  display: flex;
  gap: 12px;
  align-items: center;
  flex-wrap: wrap;
}

/* Pagination */
.pagination {
  margin-top: 20px;
  display: flex;
  justify-content: center;
}

/* Detail Drawer */
.detail-descriptions {
  margin-bottom: 24px;
}

:deep(.detail-descriptions .el-descriptions__body) {
  background: transparent;
}

.detail-section {
  margin-top: 24px;
}

.detail-section h4 {
  margin: 0 0 16px 0;
  color: var(--el-text-color-primary);
  font-weight: 600;
  font-size: 0.875rem;
  display: flex;
  align-items: center;
  gap: 8px;
}

.detail-section h4::before {
  content: '';
  display: block;
  width: 4px;
  height: 16px;
  background: linear-gradient(180deg, var(--pm-primary), var(--pm-secondary));
  border-radius: 2px;
}

.detail-content {
  background-color: var(--el-fill-color-light);
  padding: 20px;
  border-radius: 12px;
  border: 1px solid var(--el-border-color-lighter);
}

.detail-item, .change-item {
  margin-bottom: 12px;
  padding-bottom: 12px;
  border-bottom: 1px dashed var(--el-border-color-lighter);
  word-break: break-all;
}

.detail-item:last-child, .change-item:last-child {
  margin-bottom: 0;
  padding-bottom: 0;
  border-bottom: none;
}

.old-value {
  color: var(--el-color-danger);
  text-decoration: line-through;
  padding: 2px 6px;
  background: rgba(239, 68, 68, 0.1);
  border-radius: 4px;
}

.new-value {
  color: var(--el-color-success);
  padding: 2px 6px;
  background: rgba(34, 197, 94, 0.1);
  border-radius: 4px;
}

.login-success {
  color: var(--el-color-success);
  font-weight: 600;
  margin: 0 0 12px 0;
  display: flex;
  align-items: center;
  gap: 8px;
}

.login-success::before {
  content: '✓';
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 20px;
  height: 20px;
  background: rgba(34, 197, 94, 0.15);
  border-radius: 50%;
  font-size: 12px;
}

.user-agent {
  color: var(--el-text-color-secondary);
  font-size: 12px;
  margin: 0;
  word-break: break-all;
  line-height: 1.6;
  padding: 12px;
  background: var(--el-fill-color);
  border-radius: 8px;
}

.no-details {
  text-align: center;
  color: var(--el-text-color-secondary);
  padding: 60px 0;
  font-size: 0.875rem;
}
</style>
