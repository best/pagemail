import apiClient from './client'
import type { PaginatedResponse } from '@/types/api'
import type { Task, TaskCreatePayload } from '@/types/task'

export const tasksApi = {
  listTasks(params: { page?: number; limit?: number; status?: string }) {
    return apiClient.get<PaginatedResponse<Task>>('/captures', { params })
  },

  createTask(data: TaskCreatePayload) {
    return apiClient.post<Task>('/captures', data)
  },

  getTask(id: string) {
    return apiClient.get<Task>(`/captures/${id}`)
  },

  retryTask(id: string) {
    return apiClient.post(`/captures/${id}/retry`)
  },

  deleteTask(id: string) {
    return apiClient.delete(`/captures/${id}`)
  },

  getTaskOutputs(taskId: string) {
    return apiClient.get(`/captures/${taskId}/outputs`)
  },

  async downloadOutput(taskId: string, outputId: string, format: string) {
    const response = await apiClient.get(`/captures/${taskId}/outputs/${outputId}/download`, {
      responseType: 'blob'
    })
    const blob = new Blob([response.data])
    const url = window.URL.createObjectURL(blob)
    const link = document.createElement('a')
    link.href = url
    link.download = `${format}.${format === 'screenshot' ? 'png' : format}`
    document.body.appendChild(link)
    link.click()
    document.body.removeChild(link)
    window.URL.revokeObjectURL(url)
  },

  async previewOutput(taskId: string, outputId: string): Promise<Blob> {
    const response = await apiClient.get(`/captures/${taskId}/outputs/${outputId}/preview`, {
      responseType: 'blob'
    })
    return response.data as Blob
  }
}
