<script setup lang="ts">
import { ref, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { tasksApi } from '@/api/tasks'
import type { Task } from '@/types/task'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Download, Refresh, Delete, Back } from '@element-plus/icons-vue'
import { usePolling } from '@/composables/usePolling'

const { t } = useI18n()
const route = useRoute()
const router = useRouter()
const task = ref<Task | null>(null)
const loading = ref(true)
const lastRefreshed = ref<Date>(new Date())

const fetchTask = async () => {
  const isInitialLoad = !task.value
  if (isInitialLoad) loading.value = true
  try {
    const res = await tasksApi.getTask(route.params.id as string)
    task.value = res.data
    lastRefreshed.value = new Date()
  } catch {
    // Only redirect on initial load failure, not during polling
    if (isInitialLoad) router.push('/tasks')
  } finally {
    loading.value = false
  }
}

const formatTime = (date: Date) => date.toLocaleTimeString('en-US', { hour: '2-digit', minute: '2-digit', second: '2-digit' })

const isPendingTask = computed(() => task.value?.status === 'pending' || task.value?.status === 'processing')

const { isRunning, stop: stopPolling } = usePolling(fetchTask, {
  intervalMs: 10000,
  pendingIntervalMs: 5000,
  isPending: () => isPendingTask.value
})

const handleRetry = async () => {
  if (!task.value) return
  try {
    await tasksApi.retryTask(task.value.id)
    ElMessage.success(t('taskDetail.retryInitiated'))
    fetchTask()
  } catch {
    // handled globally
  }
}

const handleDelete = () => {
  if (!task.value) return
  ElMessageBox.confirm(t('taskDetail.deleteTask'), 'Warning', { type: 'warning' })
    .then(async () => {
      stopPolling()
      await tasksApi.deleteTask(task.value!.id)
      ElMessage.success(t('tasks.deleteSuccess'))
      router.push('/tasks')
    })
    .catch(() => {})
}

const handleDownload = async (outputId: string, format: string) => {
  if (!task.value) return
  try {
    await tasksApi.downloadOutput(task.value.id, outputId, format)
  } catch {
    ElMessage.error(t('taskDetail.downloadFailed'))
  }
}

const getStatusType = (status: string): 'success' | 'danger' | 'warning' | 'info' => {
  const map: Record<string, 'success' | 'danger' | 'warning' | 'info'> = {
    completed: 'success',
    failed: 'danger',
    processing: 'warning',
    pending: 'info'
  }
  return map[status] || 'info'
}
</script>

<template>
  <div class="task-detail" v-if="task">
    <div class="header">
      <el-button :icon="Back" circle @click="$router.back()" />
      <h2>{{ t('taskDetail.title') }}</h2>
      <span v-if="!loading" class="last-updated">
        <el-icon v-if="isRunning" class="is-loading"><Refresh /></el-icon>
        {{ formatTime(lastRefreshed) }}
      </span>
    </div>

    <el-row :gutter="20">
      <el-col :xs="24" :lg="16">
        <el-card class="mb-4">
          <template #header>
            <div class="card-header">
              <span>{{ t('taskDetail.taskInfo') }}</span>
              <el-tag :type="getStatusType(task.status)">{{ task.status.toUpperCase() }}</el-tag>
            </div>
          </template>
          <el-descriptions :column="1" border>
            <el-descriptions-item :label="t('taskDetail.id')">{{ task.id }}</el-descriptions-item>
            <el-descriptions-item :label="t('tasks.url')">
              <a :href="task.url" target="_blank">{{ task.url }}</a>
            </el-descriptions-item>
            <el-descriptions-item :label="t('tasks.created')">
              {{ new Date(task.created_at).toLocaleString() }}
            </el-descriptions-item>
            <el-descriptions-item :label="t('tasks.formats')">{{ task.formats.join(', ') }}</el-descriptions-item>
            <el-descriptions-item v-if="task.error_message" :label="t('taskDetail.error')">
              <span class="text-danger">{{ task.error_message }}</span>
            </el-descriptions-item>
          </el-descriptions>

          <div class="actions mt-4">
            <el-button v-if="task.status === 'failed'" type="warning" :icon="Refresh" @click="handleRetry">
              {{ t('common.retry') }}
            </el-button>
            <el-button type="danger" :icon="Delete" @click="handleDelete">{{ t('common.delete') }}</el-button>
          </div>
        </el-card>

        <el-card v-if="task.outputs && task.outputs.length > 0" class="pm-table-card">
          <template #header><span>{{ t('taskDetail.generatedOutputs') }}</span></template>
          <el-table :data="task.outputs" stripe style="width: 100%">
            <el-table-column prop="format" :label="t('taskDetail.format')" width="100">
              <template #default="{ row }">{{ row.format.toUpperCase() }}</template>
            </el-table-column>
            <el-table-column prop="size" :label="t('taskDetail.size')">
              <template #default="{ row }">{{ (row.size / 1024).toFixed(1) }} KB</template>
            </el-table-column>
            <el-table-column align="right">
              <template #default="{ row }">
                <el-button size="small" :icon="Download" @click="handleDownload(row.id, row.format)">
                  {{ t('common.download') }}
                </el-button>
              </template>
            </el-table-column>
          </el-table>
        </el-card>
      </el-col>

      <el-col :xs="24" :lg="8">
        <el-card>
          <template #header><span>{{ t('taskDetail.deliveryHistory') }}</span></template>
          <el-timeline v-if="task.delivery_history && task.delivery_history.length">
            <el-timeline-item
              v-for="(item, index) in task.delivery_history"
              :key="index"
              :type="item.status === 'success' ? 'success' : 'danger'"
              :timestamp="new Date(item.attempt_time).toLocaleString()"
            >
              <h4>{{ item.channel }}</h4>
              <p>{{ item.status }}</p>
              <p v-if="item.error" class="text-danger">{{ item.error }}</p>
            </el-timeline-item>
          </el-timeline>
          <el-empty v-else :description="t('taskDetail.noDelivery')" />
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<style scoped>
.task-detail {
  max-width: 1200px;
  margin: 0 auto;
}
.header {
  display: flex;
  align-items: center;
  gap: 16px;
  margin-bottom: 20px;
}
.header h2 {
  margin: 0;
  flex: 1;
}
.last-updated {
  font-size: 0.75rem;
  color: var(--pm-text-muted);
  display: flex;
  align-items: center;
  gap: 6px;
}
.last-updated .is-loading {
  animation: pm-spin 1s linear infinite;
}
.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}
.text-danger {
  color: var(--el-color-danger);
}
.mt-4 {
  margin-top: 16px;
}
.mb-4 {
  margin-bottom: 16px;
}
</style>
