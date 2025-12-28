import apiClient from './client'
import type { WebhookConfig } from '@/types/settings'

export const webhooksApi = {
  listWebhooks() {
    return apiClient.get<WebhookConfig[]>('/webhooks')
  },

  createWebhook(data: Partial<WebhookConfig>) {
    return apiClient.post<WebhookConfig>('/webhooks', data)
  },

  updateWebhook(id: string, data: Partial<WebhookConfig>) {
    return apiClient.put<WebhookConfig>(`/webhooks/${id}`, data)
  },

  deleteWebhook(id: string) {
    return apiClient.delete(`/webhooks/${id}`)
  },

  testWebhook(id: string) {
    return apiClient.post(`/webhooks/${id}/test`)
  }
}
