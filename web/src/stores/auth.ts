import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import type { User } from '@/types/user'
import apiClient from '@/api/client'

export const useAuthStore = defineStore('auth', () => {
  const token = ref<string | null>(null)
  const refreshToken = ref<string | null>(null)
  const user = ref<User | null>(null)

  const isAuthenticated = computed(() => !!token.value)
  const isAdmin = computed(() => user.value?.role === 'admin')

  async function login(email: string, password: string) {
    const response = await apiClient.post('/auth/login', { email, password })
    token.value = response.data.access_token
    refreshToken.value = response.data.refresh_token
    user.value = response.data.user
    return response.data
  }

  async function register(email: string, password: string) {
    const response = await apiClient.post('/auth/register', { email, password })
    return response.data
  }

  function logout() {
    token.value = null
    refreshToken.value = null
    user.value = null
  }

  async function fetchProfile() {
    if (!token.value) return
    const response = await apiClient.get('/users/me')
    user.value = response.data
  }

  return {
    token,
    refreshToken,
    user,
    isAuthenticated,
    isAdmin,
    login,
    register,
    logout,
    fetchProfile
  }
}, {
  persist: {
    paths: ['token', 'refreshToken', 'user']
  }
})
