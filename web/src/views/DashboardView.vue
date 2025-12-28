<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { useAuthStore } from '@/stores/auth'
import { tasksApi } from '@/api/tasks'
import type { Task } from '@/types/task'
import { Document, Finished, Warning, Plus } from '@element-plus/icons-vue'

const authStore = useAuthStore()
const tasks = ref<Task[]>([])
const loading = ref(true)
const loadError = ref(false)

const stats = computed(() => {
  const total = tasks.value.length
  const completed = tasks.value.filter(t => t.status === 'completed').length
  const pending = tasks.value.filter(t => t.status === 'pending' || t.status === 'processing').length
  return { total, completed, pending }
})

const timeGreeting = computed(() => {
  const hour = new Date().getHours()
  if (hour < 12) return 'Good morning'
  if (hour < 18) return 'Good afternoon'
  return 'Good evening'
})

const userName = computed(() => {
  const email = authStore.user?.email
  return email ? email.split('@')[0] : 'User'
})

const recentTasks = computed(() => tasks.value.slice(0, 5))

const fetchTasks = async () => {
  loading.value = true
  loadError.value = false
  try {
    const res = await tasksApi.listTasks({ limit: 100 })
    tasks.value = res.data.data || []
  } catch {
    tasks.value = []
    loadError.value = true
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
    <div class="dashboard-header">
      <div class="header-content">
        <h1>{{ timeGreeting }}, {{ userName }}</h1>
        <p>Here's an overview of your web page captures</p>
      </div>
      <el-button type="primary" size="large" :icon="Plus" @click="$router.push('/tasks/new')">
        New Capture
      </el-button>
    </div>

    <el-row :gutter="20" class="stats">
      <el-col :xs="24" :sm="8">
        <div class="stat-card primary">
          <div class="stat-icon">
            <el-icon :size="24"><Document /></el-icon>
          </div>
          <div class="stat-info">
            <div class="stat-label">Total Captures</div>
            <div class="stat-value">{{ stats.total }}</div>
          </div>
          <el-icon class="stat-bg-icon"><Document /></el-icon>
        </div>
      </el-col>
      <el-col :xs="24" :sm="8">
        <div class="stat-card success">
          <div class="stat-icon">
            <el-icon :size="24"><Finished /></el-icon>
          </div>
          <div class="stat-info">
            <div class="stat-label">Completed</div>
            <div class="stat-value">{{ stats.completed }}</div>
          </div>
          <el-icon class="stat-bg-icon"><Finished /></el-icon>
        </div>
      </el-col>
      <el-col :xs="24" :sm="8">
        <div class="stat-card warning">
          <div class="stat-icon">
            <el-icon :size="24"><Warning /></el-icon>
          </div>
          <div class="stat-info">
            <div class="stat-label">Pending</div>
            <div class="stat-value">{{ stats.pending }}</div>
          </div>
          <el-icon class="stat-bg-icon"><Warning /></el-icon>
        </div>
      </el-col>
    </el-row>

    <div class="section-header">
      <h2>Recent Activity</h2>
      <el-button text type="primary" @click="$router.push('/tasks')">View All</el-button>
    </div>

    <el-card class="recent-tasks" shadow="never">
      <el-alert
        v-if="loadError"
        title="Failed to load tasks"
        type="error"
        show-icon
        :closable="false"
        style="margin-bottom: 16px"
      >
        <template #default>
          <el-button size="small" type="primary" link @click="fetchTasks">Retry</el-button>
        </template>
      </el-alert>
      <el-table v-if="recentTasks.length > 0" :data="recentTasks" v-loading="loading">
        <el-table-column prop="url" label="URL" show-overflow-tooltip />
        <el-table-column prop="status" label="Status" width="120">
          <template #default="{ row }">
            <el-tag :type="getStatusType(row.status)" effect="light" round>{{ row.status }}</el-tag>
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
      <el-empty v-else-if="!loading && !loadError" description="No tasks yet" />
    </el-card>
  </div>
</template>

<style scoped>
.dashboard {
  max-width: 1200px;
  margin: 0 auto;
  padding-bottom: 40px;
}

.dashboard-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 32px;
}

.dashboard-header h1 {
  margin: 0;
  font-size: 1.75rem;
  font-weight: 600;
  color: var(--pm-text-heading);
}

.dashboard-header p {
  margin: 8px 0 0;
  color: var(--pm-text-muted);
}

.stats {
  margin-bottom: 40px;
}

.stat-card {
  position: relative;
  display: flex;
  align-items: flex-start;
  gap: 16px;
  padding: 24px;
  background: var(--pm-bg-card);
  border-radius: 12px;
  box-shadow: var(--pm-shadow-md);
  overflow: hidden;
  transition: transform 0.2s, box-shadow 0.2s;
}

.stat-card:hover {
  transform: translateY(-4px);
  box-shadow: var(--pm-shadow-lg);
}

.stat-icon {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 48px;
  height: 48px;
  border-radius: 50%;
  flex-shrink: 0;
}

.stat-card.primary .stat-icon {
  background: var(--el-color-primary-light-9);
  color: var(--pm-primary);
}

.stat-card.success .stat-icon {
  background: var(--el-color-success-light-9);
  color: var(--el-color-success);
}

.stat-card.warning .stat-icon {
  background: var(--el-color-warning-light-9);
  color: var(--el-color-warning);
}

html.dark .stat-card.primary .stat-icon {
  background: rgba(79, 70, 229, 0.15);
  color: var(--pm-primary-light);
}

html.dark .stat-card.success .stat-icon {
  background: rgba(34, 197, 94, 0.15);
  color: #4ade80;
}

html.dark .stat-card.warning .stat-icon {
  background: rgba(234, 179, 8, 0.15);
  color: #fbbf24;
}

.stat-info {
  flex: 1;
}

.stat-label {
  color: var(--pm-text-muted);
  font-size: 0.875rem;
  margin-bottom: 4px;
}

.stat-value {
  font-size: 2.25rem;
  font-weight: 700;
  line-height: 1.2;
  color: var(--pm-text-heading);
}

.stat-bg-icon {
  position: absolute;
  right: -10px;
  bottom: -10px;
  font-size: 100px;
  opacity: 0.05;
  pointer-events: none;
}

.section-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16px;
}

.section-header h2 {
  margin: 0;
  font-size: 1.125rem;
  font-weight: 600;
  color: var(--pm-text-heading);
}

.recent-tasks {
  border: none;
  background: var(--pm-bg-card);
  border-radius: 12px;
}

:deep(.el-table) {
  --el-table-border-color: transparent;
}

:deep(.el-table__inner-wrapper::before) {
  display: none;
}

:deep(.el-table tr) {
  transition: background-color 0.2s;
}

:deep(.el-table tr:hover > td) {
  background-color: var(--el-fill-color-lighter) !important;
}

@media (max-width: 768px) {
  .dashboard-header {
    flex-direction: column;
    align-items: flex-start;
    gap: 16px;
  }

  .stat-card {
    margin-bottom: 16px;
  }
}
</style>
