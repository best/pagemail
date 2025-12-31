import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { publicApi } from '@/api/public'

export const useSiteConfigStore = defineStore('siteConfig', () => {
  const siteName = ref('Pagemail')
  const siteSlogan = ref('')
  const loaded = ref(false)
  const loading = ref(false)

  const copyright = computed(() => `© ${new Date().getFullYear()} ${siteName.value}.`)

  async function fetchConfig() {
    if (loaded.value || loading.value) return

    loading.value = true
    try {
      const response = await publicApi.getSiteConfig()
      if (response.data.site_name) {
        siteName.value = response.data.site_name
      }
      siteSlogan.value = response.data.site_slogan || ''
      loaded.value = true
    } catch (error) {
      console.error('Failed to fetch site config', error)
    } finally {
      loading.value = false
    }
  }

  function updateConfig(name: string, slogan: string) {
    siteName.value = name
    siteSlogan.value = slogan
  }

  return { siteName, siteSlogan, copyright, loaded, fetchConfig, updateConfig }
})
