import apiClient from './client'
import type { User } from '@/types/user'

export const userApi = {
  uploadAvatar(file: File) {
    const formData = new FormData()
    formData.append('avatar', file)
    return apiClient.put<User>('/users/me/avatar', formData)
  },

  deleteAvatar() {
    return apiClient.delete<User>('/users/me/avatar')
  }
}
