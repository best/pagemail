import apiClient from './client'
import type { SmtpProfile } from '@/types/settings'

export const smtpApi = {
  listProfiles() {
    return apiClient.get<SmtpProfile[]>('/smtp/profiles')
  },

  createProfile(data: Partial<SmtpProfile>) {
    return apiClient.post<SmtpProfile>('/smtp/profiles', data)
  },

  updateProfile(id: string, data: Partial<SmtpProfile>) {
    return apiClient.put<SmtpProfile>(`/smtp/profiles/${id}`, data)
  },

  deleteProfile(id: string) {
    return apiClient.delete(`/smtp/profiles/${id}`)
  },

  testProfile(id: string, email: string) {
    return apiClient.post(`/smtp/profiles/${id}/test`, { email })
  }
}
