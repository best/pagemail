<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { tasksApi } from '@/api/tasks'
import type { Task } from '@/types/task'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Download, Refresh, Delete, Back } from '@element-plus/icons-vue'

const route = useRoute()
const router = useRouter()
const task = ref<Task | null>(null)
const loading = ref(true)

const fetchTask = async () => {
  loading.value = true
  try {
    const res = await tasksApi.getTask(route.params.id as string)
    task.value = res.data
  } catch {
    router.push('/tasks')
  } finally {
    loading.value = false
  }
}

const handleRetry = async () => {
  if (!task.value) return
  try {
    await tasksApi.retryTask(task.value.id)
    ElMessage.success('Retry initiated')
    fetchTask()
  } catch {
    // handled globally
  }
}

const handleDelete = () => {
  if (!task.value) return
  ElMessageBox.confirm('Delete this task?', 'Warning', { type: 'warning' })
    .then(async () => {
      await tasksApi.deleteTask(task.value!.id)
      ElMessage.success('Task deleted')
      router.push('/tasks')
    })
    .catch(() => {})
}

const handleDownload = async (outputId: string, format: string) => {
  if (!task.value) return
  try {
    await tasksApi.downloadOutput(task.value.id, outputId, format)
  } catch {
    ElMessage.error('Download failed')
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

onMounted(fetchTask)
</script>

<template>
  <div class="task-detail" v-if="task">
    <div class="header">
      <el-button :icon="Back" circle @click="$router.back()" />
      <h2>Task Detail</h2>
    </div>

    <el-row :gutter="20">
      <el-col :span="16">
        <el-card class="mb-4">
          <template #header>
            <div class="card-header">
              <span>Task Information</span>
              <el-tag :type="getStatusType(task.status)">{{ task.status.toUpperCase() }}</el-tag>
            </div>
          </template>
          <el-descriptions :column="1" border>
            <el-descriptions-item label="ID">{{ task.id }}</el-descriptions-item>
            <el-descriptions-item label="URL">
              <a :href="task.url" target="_blank">{{ task.url }}</a>
            </el-descriptions-item>
            <el-descriptions-item label="Created">
              {{ new Date(task.created_at).toLocaleString() }}
            </el-descriptions-item>
            <el-descriptions-item label="Formats">{{ task.formats.join(', ') }}</el-descriptions-item>
            <el-descriptions-item v-if="task.error_message" label="Error">
              <span class="text-danger">{{ task.error_message }}</span>
            </el-descriptions-item>
          </el-descriptions>

          <div class="actions mt-4">
            <el-button v-if="task.status === 'failed'" type="warning" :icon="Refresh" @click="handleRetry">
              Retry
            </el-button>
            <el-button type="danger" :icon="Delete" @click="handleDelete">Delete</el-button>
          </div>
        </el-card>

        <el-card v-if="task.outputs && task.outputs.length > 0">
          <template #header><span>Generated Outputs</span></template>
          <el-table :data="task.outputs">
            <el-table-column prop="format" label="Format" width="100">
              <template #default="{ row }">{{ row.format.toUpperCase() }}</template>
            </el-table-column>
            <el-table-column prop="size" label="Size">
              <template #default="{ row }">{{ (row.size / 1024).toFixed(1) }} KB</template>
            </el-table-column>
            <el-table-column align="right">
              <template #default="{ row }">
                <el-button size="small" :icon="Download" @click="handleDownload(row.id, row.format)">
                  Download
                </el-button>
              </template>
            </el-table-column>
          </el-table>
        </el-card>
      </el-col>

      <el-col :span="8">
        <el-card>
          <template #header><span>Delivery History</span></template>
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
          <el-empty v-else description="No delivery attempts" />
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
