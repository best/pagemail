import apiClient from './client'
import type { User } from '@/types/user'
import type { PaginatedResponse } from '@/types/api'

export const adminApi = {
  listUsers(params: { page?: number; limit?: number }) {
    return apiClient.get<PaginatedResponse<User>>('/admin/users', { params })
  },

  updateUser(id: string, data: Partial<User>) {
    return apiClient.put<User>(`/admin/users/${id}`, data)
  },

  deleteUser(id: string) {
    return apiClient.delete(`/admin/users/${id}`)
  },

  getSystemConfig() {
    return apiClient.get('/admin/system/config')
  },

  updateSystemConfig(data: Record<string, unknown>) {
    return apiClient.put('/admin/system/config', data)
  },

  getAuditLogs(params: { page?: number; limit?: number }) {
    return apiClient.get<PaginatedResponse<AuditLog>>('/admin/audit-logs', { params })
  }
}

export interface AuditLog {
  id: string
  actor_id: string
  actor_email: string
  action: string
  resource_type: string
  resource_id?: string
  details: Record<string, unknown>
  ip_address: string
  created_at: string
}
