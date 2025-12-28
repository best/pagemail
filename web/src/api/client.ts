import axios, { type AxiosError, type InternalAxiosRequestConfig } from 'axios'
import type { ProblemDetail } from '@/types/api'
import { ElMessage } from 'element-plus'

const apiClient = axios.create({
  baseURL: '/api',
  timeout: 30000,
  headers: {
    'Content-Type': 'application/json'
  }
})

// Token will be read from localStorage directly to avoid circular dependency
function getToken(): string | null {
  try {
    const stored = localStorage.getItem('auth')
    if (stored) {
      const parsed = JSON.parse(stored)
      return parsed.token || null
    }
  } catch {
    // ignore parse errors
  }
  return null
}

apiClient.interceptors.request.use(
  (config: InternalAxiosRequestConfig) => {
    const token = getToken()
    if (token) {
      config.headers.Authorization = `Bearer ${token}`
    }
    return config
  },
  (error) => {
    return Promise.reject(error)
  }
)

apiClient.interceptors.response.use(
  (response) => response,
  async (error: AxiosError<ProblemDetail>) => {
    const status = error.response?.status
    const problem = error.response?.data

    if (status === 401) {
      // Clear auth data and redirect
      localStorage.removeItem('auth')
      ElMessage.error('Session expired. Please login again.')
      window.location.href = '/login'
    } else if (status === 403) {
      ElMessage.error('Access denied.')
    } else if (problem?.title) {
      ElMessage.error(`${problem.title}: ${problem.detail || 'An error occurred'}`)
    } else {
      ElMessage.error(error.message || 'Network error')
    }

    return Promise.reject(error)
  }
)

export default apiClient
