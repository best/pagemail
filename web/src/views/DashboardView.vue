<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { useAuthStore } from '@/stores/auth'
import { tasksApi } from '@/api/tasks'
import type { Task } from '@/types/task'
import { Document, Finished, Warning, Plus } from '@element-plus/icons-vue'

const authStore = useAuthStore()
const tasks = ref<Task[]>([])
const loading = ref(true)

const stats = computed(() => {
  const total = tasks.value.length
  const completed = tasks.value.filter(t => t.status === 'completed').length
  const pending = tasks.value.filter(t => t.status === 'pending' || t.status === 'processing').length
  return { total, completed, pending }
})

const recentTasks = computed(() => tasks.value.slice(0, 5))

const fetchTasks = async () => {
  loading.value = true
  try {
    const res = await tasksApi.listTasks({ limit: 100 })
    tasks.value = res.data.data || []
  } catch {
    tasks.value = []
  } finally {
    loading.value = false
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

onMounted(fetchTasks)
</script>

<template>
  <div class="dashboard">
    <div class="welcome">
      <h1>Welcome, {{ authStore.user?.email }}</h1>
      <p>Manage your web page captures and deliveries</p>
    </div>

    <el-row :gutter="20" class="stats">
      <el-col :span="8">
        <el-card shadow="hover">
          <div class="stat-card">
            <el-icon :size="40" color="var(--el-color-primary)"><Document /></el-icon>
            <div class="stat-info">
              <div class="stat-value">{{ stats.total }}</div>
              <div class="stat-label">Total Captures</div>
            </div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="8">
        <el-card shadow="hover">
          <div class="stat-card">
            <el-icon :size="40" color="var(--el-color-success)"><Finished /></el-icon>
            <div class="stat-info">
              <div class="stat-value">{{ stats.completed }}</div>
              <div class="stat-label">Completed</div>
            </div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="8">
        <el-card shadow="hover">
          <div class="stat-card">
            <el-icon :size="40" color="var(--el-color-warning)"><Warning /></el-icon>
            <div class="stat-info">
              <div class="stat-value">{{ stats.pending }}</div>
              <div class="stat-label">Pending</div>
            </div>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <el-card class="quick-actions">
      <template #header>
        <div class="card-header">
          <span>Quick Actions</span>
        </div>
      </template>
      <el-button type="primary" :icon="Plus" @click="$router.push('/tasks/new')">
        New Capture
      </el-button>
    </el-card>

    <el-card class="recent-tasks">
      <template #header>
        <div class="card-header">
          <span>Recent Tasks</span>
          <el-button text type="primary" @click="$router.push('/tasks')">View All</el-button>
        </div>
      </template>
      <el-table v-if="recentTasks.length > 0" :data="recentTasks" v-loading="loading">
        <el-table-column prop="url" label="URL" show-overflow-tooltip />
        <el-table-column prop="status" label="Status" width="120">
          <template #default="{ row }">
            <el-tag :type="getStatusType(row.status)" size="small">{{ row.status }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="created_at" label="Created" width="180">
          <template #default="{ row }">{{ new Date(row.created_at).toLocaleString() }}</template>
        </el-table-column>
        <el-table-column width="100" align="right">
          <template #default="{ row }">
            <el-button size="small" text type="primary" @click="$router.push(`/tasks/${row.id}`)">View</el-button>
          </template>
        </el-table-column>
      </el-table>
      <el-empty v-else-if="!loading" description="No tasks yet" />
    </el-card>
  </div>
</template>

<style scoped>
.dashboard {
  max-width: 1200px;
  margin: 0 auto;
}

.welcome {
  margin-bottom: 24px;
}

.welcome h1 {
  margin: 0;
  font-size: 1.5rem;
}

.welcome p {
  margin: 8px 0 0;
  color: var(--el-text-color-secondary);
}

.stats {
  margin-bottom: 24px;
}

.stat-card {
  display: flex;
  align-items: center;
  gap: 16px;
}

.stat-info {
  flex: 1;
}

.stat-value {
  font-size: 2rem;
  font-weight: bold;
}

.stat-label {
  color: var(--el-text-color-secondary);
}

.quick-actions {
  margin-bottom: 24px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}
</style>
