import apiClient from './client'
import type { SiteConfig } from '@/types/api'

export const publicApi = {
  getSiteConfig() {
    return apiClient.get<SiteConfig>('/config/site')
  }
}
