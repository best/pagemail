<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { adminApi } from '@/api/admin'
import type { User } from '@/types/user'
import { ElMessage, ElMessageBox } from 'element-plus'

const { t } = useI18n()

const users = ref<User[]>([])
const loading = ref(false)
const total = ref(0)
const page = ref(1)

const fetchUsers = async () => {
  loading.value = true
  try {
    const res = await adminApi.listUsers({ page: page.value, limit: 10 })
    users.value = res.data.data
    total.value = res.data.meta.total
  } finally {
    loading.value = false
  }
}

const toggleStatus = async (user: User) => {
  try {
    await adminApi.updateUser(user.id, { is_active: !user.is_active })
    ElMessage.success(t('admin.userManagement.statusUpdated'))
    fetchUsers()
  } catch {
    // handled globally
  }
}

const handleDelete = (id: string) => {
  ElMessageBox.confirm(t('admin.userManagement.deleteConfirm'), 'Warning', { type: 'warning' })
    .then(async () => {
      await adminApi.deleteUser(id)
      ElMessage.success(t('admin.userManagement.userDeleted'))
      fetchUsers()
    })
    .catch(() => {})
}

onMounted(fetchUsers)
</script>

<template>
  <div class="user-mgmt">
    <div class="page-header">
      <h2>{{ t('admin.userManagement.title') }}</h2>
    </div>
    <el-table :data="users" v-loading="loading">
      <el-table-column prop="email" :label="t('admin.userManagement.email')" />
      <el-table-column prop="role" :label="t('admin.userManagement.role')">
        <template #default="{ row }">
          <el-tag :type="row.role === 'admin' ? 'danger' : 'info'">{{ row.role }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column :label="t('admin.userManagement.status')">
        <template #default="{ row }">
          <el-tag :type="row.is_active ? 'success' : 'info'">
            {{ row.is_active ? t('webhook.active') : t('webhook.inactive') }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column :label="t('admin.userManagement.created')" width="180">
        <template #default="{ row }">{{ new Date(row.created_at).toLocaleString() }}</template>
      </el-table-column>
      <el-table-column :label="t('admin.userManagement.actions')" align="right">
        <template #default="{ row }">
          <el-button
            size="small"
            :type="row.is_active ? 'warning' : 'success'"
            @click="toggleStatus(row)"
          >
            {{ row.is_active ? t('admin.userManagement.deactivate') : t('admin.userManagement.activate') }}
          </el-button>
          <el-button size="small" type="danger" @click="handleDelete(row.id)">{{ t('common.delete') }}</el-button>
        </template>
      </el-table-column>
    </el-table>
    <div class="pagination">
      <el-pagination
        background
        layout="prev, pager, next"
        :total="total"
        v-model:current-page="page"
        @current-change="fetchUsers"
      />
    </div>
  </div>
</template>

<style scoped>
.user-mgmt {
  max-width: 1200px;
  margin: 0 auto;
}
.page-header {
  margin-bottom: 20px;
}
.page-header h2 {
  margin: 0;
}
.pagination {
  margin-top: 20px;
  display: flex;
  justify-content: center;
}
</style>
