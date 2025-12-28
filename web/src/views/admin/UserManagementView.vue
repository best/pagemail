<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { adminApi } from '@/api/admin'
import type { User } from '@/types/user'
import { ElMessage, ElMessageBox } from 'element-plus'

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
    ElMessage.success('User status updated')
    fetchUsers()
  } catch {
    // handled globally
  }
}

const handleDelete = (id: string) => {
  ElMessageBox.confirm('Delete user?', 'Warning', { type: 'warning' })
    .then(async () => {
      await adminApi.deleteUser(id)
      ElMessage.success('User deleted')
      fetchUsers()
    })
    .catch(() => {})
}

onMounted(fetchUsers)
</script>

<template>
  <div class="user-mgmt">
    <h2>User Management</h2>
    <el-table :data="users" v-loading="loading">
      <el-table-column prop="email" label="Email" />
      <el-table-column prop="role" label="Role">
        <template #default="{ row }">
          <el-tag :type="row.role === 'admin' ? 'danger' : 'info'">{{ row.role }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column label="Status">
        <template #default="{ row }">
          <el-tag :type="row.is_active ? 'success' : 'info'">
            {{ row.is_active ? 'Active' : 'Inactive' }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column label="Created" width="180">
        <template #default="{ row }">{{ new Date(row.created_at).toLocaleString() }}</template>
      </el-table-column>
      <el-table-column label="Actions" align="right">
        <template #default="{ row }">
          <el-button
            size="small"
            :type="row.is_active ? 'warning' : 'success'"
            @click="toggleStatus(row)"
          >
            {{ row.is_active ? 'Deactivate' : 'Activate' }}
          </el-button>
          <el-button size="small" type="danger" @click="handleDelete(row.id)">Delete</el-button>
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
  padding: 20px;
}
.pagination {
  margin-top: 20px;
  display: flex;
  justify-content: center;
}
</style>
