<script setup lang="ts">
import { ref, computed } from 'vue'
import { useAuthStore } from '@/stores/auth'
import { tasksApi } from '@/api/tasks'
import type { Task } from '@/types/task'
import { Document, Finished, Warning, Plus, Refresh } from '@element-plus/icons-vue'
import { usePolling } from '@/composables/usePolling'

const authStore = useAuthStore()
const tasks = ref<Task[]>([])
const loading = ref(true)
const loadError = ref(false)
const lastRefreshed = ref<Date>(new Date())

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
  if (tasks.value.length === 0) loading.value = true
  loadError.value = false
  try {
    const res = await tasksApi.listTasks({ limit: 100 })
    tasks.value = res.data.data || []
    lastRefreshed.value = new Date()
  } catch {
    tasks.value = []
    loadError.value = true
  } finally {
    loading.value = false
  }
}

const formatTime = (date: Date) => date.toLocaleTimeString('en-US', { hour: '2-digit', minute: '2-digit', second: '2-digit' })

const hasPendingTasks = computed(() => tasks.value.some(t => t.status === 'pending' || t.status === 'processing'))

const { isRunning } = usePolling(fetchTasks, {
  intervalMs: 10000,
  pendingIntervalMs: 5000,
  isPending: () => hasPendingTasks.value
})

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
  <div class="dashboard">
    <div class="dashboard-header anim-up">
      <div class="header-content">
        <h1><span class="gradient-text">{{ timeGreeting }}</span>, {{ userName }}</h1>
        <p>Here's an overview of your web page captures</p>
      </div>
      <el-button type="primary" size="large" :icon="Plus" @click="$router.push('/tasks/new')">
        New Capture
      </el-button>
    </div>

    <el-row :gutter="20" class="stats anim-up d1">
      <el-col :xs="24" :sm="8">
        <div class="stat-card primary">
          <div class="stat-icon">
            <el-icon :size="24"><Document /></el-icon>
          </div>
          <div class="stat-info">
            <div class="stat-label">Total Captures</div>
            <div class="stat-value">{{ stats.total }}</div>
          </div>
          <el-icon class="stat-bg-icon" aria-hidden="true"><Document /></el-icon>
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
          <el-icon class="stat-bg-icon" aria-hidden="true"><Finished /></el-icon>
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
          <el-icon class="stat-bg-icon" aria-hidden="true"><Warning /></el-icon>
        </div>
      </el-col>
    </el-row>

    <div class="section-header anim-up d2">
      <div class="title-group">
        <h2>Recent Activity</h2>
        <span v-if="!loading" class="last-updated">
          <el-icon v-if="isRunning" class="is-loading"><Refresh /></el-icon>
          {{ formatTime(lastRefreshed) }}
        </span>
      </div>
      <el-button text type="primary" @click="$router.push('/tasks')">View All</el-button>
    </div>

    <el-card class="recent-tasks anim-up d3" shadow="never">
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
/* Entry Animations */
.anim-up {
  opacity: 0;
  animation: slideUp 0.5s ease-out forwards;
}
.d1 { animation-delay: 0.1s; }
.d2 { animation-delay: 0.2s; }
.d3 { animation-delay: 0.3s; }

@keyframes slideUp {
  from { opacity: 0; transform: translateY(20px); }
  to { opacity: 1; transform: translateY(0); }
}

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

.gradient-text {
  background: linear-gradient(135deg, var(--pm-primary) 0%, var(--pm-secondary) 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
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
  background: rgba(255, 255, 255, 0.75);
  backdrop-filter: blur(12px);
  -webkit-backdrop-filter: blur(12px);
  border: 1px solid rgba(255, 255, 255, 0.5);
  border-radius: 16px;
  box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.05), 0 2px 4px -1px rgba(0, 0, 0, 0.03);
  overflow: hidden;
  transition: transform 0.3s ease, box-shadow 0.3s ease, border-color 0.3s ease;
}

html.dark .stat-card {
  background: rgba(30, 41, 59, 0.65);
  border-color: rgba(255, 255, 255, 0.08);
}

.stat-card:hover {
  transform: translateY(-6px);
  box-shadow: 0 20px 25px -5px rgba(0, 0, 0, 0.08), 0 10px 10px -5px rgba(0, 0, 0, 0.04);
  border-color: rgba(79, 70, 229, 0.3);
}

html.dark .stat-card:hover {
  box-shadow: 0 20px 25px -5px rgba(0, 0, 0, 0.3), 0 10px 10px -5px rgba(0, 0, 0, 0.2);
  border-color: rgba(99, 102, 241, 0.4);
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

.section-header .title-group {
  display: flex;
  align-items: center;
  gap: 12px;
}

.section-header h2 {
  margin: 0;
  font-size: 1.125rem;
  font-weight: 600;
  color: var(--pm-text-heading);
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

.recent-tasks {
  border: none;
  background: rgba(255, 255, 255, 0.75);
  backdrop-filter: blur(12px);
  -webkit-backdrop-filter: blur(12px);
  border: 1px solid rgba(255, 255, 255, 0.5);
  border-radius: 16px;
}

html.dark .recent-tasks {
  background: rgba(30, 41, 59, 0.65);
  border-color: rgba(255, 255, 255, 0.08);
}

:deep(.el-table) {
  --el-table-border-color: transparent;
  background: transparent;
  --el-table-tr-bg-color: transparent;
  --el-table-header-bg-color: transparent;
}

:deep(.el-table__inner-wrapper::before) {
  display: none;
}

:deep(.el-table tr) {
  transition: background-color 0.2s;
}

:deep(.el-table tr:hover > td) {
  background-color: rgba(79, 70, 229, 0.04) !important;
}

html.dark :deep(.el-table tr:hover > td) {
  background-color: rgba(99, 102, 241, 0.08) !important;
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

@media (prefers-reduced-motion: reduce) {
  .anim-up { animation: none; opacity: 1; }
}
</style>
