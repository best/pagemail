<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { adminApi, type AuditLog } from '@/api/admin'

const { t } = useI18n()

const logs = ref<AuditLog[]>([])
const loading = ref(false)
const total = ref(0)
const query = reactive({ page: 1, limit: 20 })

const fetchLogs = async () => {
  loading.value = true
  try {
    const res = await adminApi.getAuditLogs(query)
    logs.value = res.data.data
    total.value = res.data.meta.total
  } finally {
    loading.value = false
  }
}

onMounted(fetchLogs)
</script>

<template>
  <div class="audit-log">
    <h2>{{ t('admin.auditLogs.title') }}</h2>
    <el-table :data="logs" v-loading="loading" style="width: 100%">
      <el-table-column prop="created_at" :label="t('admin.auditLogs.time')" width="180">
        <template #default="{ row }">{{ new Date(row.created_at).toLocaleString() }}</template>
      </el-table-column>
      <el-table-column prop="actor_email" :label="t('admin.auditLogs.user')" width="200" />
      <el-table-column prop="action" :label="t('admin.auditLogs.action')" width="150" />
      <el-table-column prop="resource_type" :label="t('admin.auditLogs.resource')" width="120" />
      <el-table-column prop="ip_address" :label="t('admin.auditLogs.ip')" width="140" />
      <el-table-column :label="t('admin.auditLogs.details')">
        <template #default="{ row }">
          <pre class="details-json">{{ JSON.stringify(row.details, null, 2) }}</pre>
        </template>
      </el-table-column>
    </el-table>
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
  </div>
</template>

<style scoped>
.audit-log {
  padding: 20px;
}
.pagination {
  margin-top: 20px;
  display: flex;
  justify-content: center;
}
.details-json {
  margin: 0;
  font-size: 12px;
  max-height: 100px;
  overflow: auto;
}
</style>
