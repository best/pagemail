<script setup lang="ts">
import { ref, reactive, onMounted, watch } from 'vue'
import { tasksApi } from '@/api/tasks'
import type { Task } from '@/types/task'
import { useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { View, Refresh, Delete, Plus } from '@element-plus/icons-vue'

const router = useRouter()
const tasks = ref<Task[]>([])
const total = ref(0)
const loading = ref(false)

const query = reactive({
  page: 1,
  limit: 10,
  status: ''
})

const fetchTasks = async () => {
  loading.value = true
  try {
    const res = await tasksApi.listTasks({
      page: query.page,
      limit: query.limit,
      status: query.status || undefined
    })
    tasks.value = res.data.data
    total.value = res.data.meta.total
  } finally {
    loading.value = false
  }
}

const handlePageChange = (page: number) => {
  query.page = page
  fetchTasks()
}

const handleRetry = async (id: string) => {
  try {
    await tasksApi.retryTask(id)
    ElMessage.success('Task queued for retry')
    fetchTasks()
  } catch {
    // handled globally
  }
}

const handleDelete = (id: string) => {
  ElMessageBox.confirm('Are you sure?', 'Warning', { type: 'warning' })
    .then(async () => {
      await tasksApi.deleteTask(id)
      ElMessage.success('Task deleted')
      fetchTasks()
    })
    .catch(() => {})
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

onMounted(fetchTasks)
watch(() => query.status, () => {
  query.page = 1
  fetchTasks()
})
</script>

<template>
  <div class="task-list">
    <div class="header">
      <h2>Tasks</h2>
      <div class="filters">
        <el-select v-model="query.status" placeholder="Filter Status" clearable style="width: 150px">
          <el-option label="Pending" value="pending" />
          <el-option label="Processing" value="processing" />
          <el-option label="Completed" value="completed" />
          <el-option label="Failed" value="failed" />
        </el-select>
        <el-button type="primary" :icon="Plus" @click="router.push('/tasks/new')">New Task</el-button>
      </div>
    </div>

    <el-table :data="tasks" v-loading="loading" style="width: 100%">
      <el-table-column label="URL" min-width="250">
        <template #default="{ row }">
          <div class="url-col">{{ row.url }}</div>
        </template>
      </el-table-column>
      <el-table-column label="Status" width="100">
        <template #default="{ row }">
          <el-tag :type="getStatusType(row.status)" size="small">{{ row.status }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column label="Formats" width="150">
        <template #default="{ row }">
          <el-tag v-for="fmt in row.formats" :key="fmt" size="small" class="mr-1">{{ fmt }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column label="Created" width="180">
        <template #default="{ row }">{{ new Date(row.created_at).toLocaleString() }}</template>
      </el-table-column>
      <el-table-column label="Actions" width="180" align="right">
        <template #default="{ row }">
          <el-button size="small" :icon="View" circle @click="router.push(`/tasks/${row.id}`)" />
          <el-button
            v-if="row.status === 'failed'"
            size="small"
            type="warning"
            :icon="Refresh"
            circle
            @click="handleRetry(row.id)"
          />
          <el-button size="small" type="danger" :icon="Delete" circle @click="handleDelete(row.id)" />
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
        @current-change="handlePageChange"
      />
    </div>
  </div>
</template>

<style scoped>
.task-list {
  max-width: 1200px;
  margin: 0 auto;
}
.header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}
.filters {
  display: flex;
  gap: 10px;
}
.pagination {
  margin-top: 20px;
  display: flex;
  justify-content: center;
}
.mr-1 {
  margin-right: 4px;
}
.url-col {
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}
</style>
