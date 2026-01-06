import { ref, watch, onMounted } from 'vue'
import { useAuthStore } from '@/stores/auth'
import { userApi } from '@/api/user'

// Singleton state shared across all components
const avatarUrl = ref<string>('')
let lastFetchedUrl: string | undefined = undefined

export function useAvatar() {
  const authStore = useAuthStore()

  const fetchAvatar = async () => {
    const url = authStore.user?.avatar_url

    // Dedup: skip if URL unchanged and we have a valid blob
    if (url === lastFetchedUrl && avatarUrl.value) return

    // Cleanup old URL before any state change
    if (avatarUrl.value) {
      URL.revokeObjectURL(avatarUrl.value)
      avatarUrl.value = ''
    }

    lastFetchedUrl = url

    if (!url) return

    try {
      const blob = await userApi.getAvatar(url)
      // Guard against race: if URL changed during fetch, discard result
      if (authStore.user?.avatar_url !== url) return
      avatarUrl.value = URL.createObjectURL(blob)
    } catch {
      // Reset so we can retry later
      lastFetchedUrl = undefined
    }
  }

  watch(() => authStore.user?.avatar_url, fetchAvatar)

  onMounted(() => {
    // Only fetch if not already loaded
    if (!avatarUrl.value || lastFetchedUrl !== authStore.user?.avatar_url) {
      fetchAvatar()
    }
  })

  return { avatarUrl, fetchAvatar }
}
