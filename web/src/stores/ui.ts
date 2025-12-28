import { defineStore } from 'pinia'
import { ref } from 'vue'
import { useDark, useToggle } from '@vueuse/core'

export const useUiStore = defineStore('ui', () => {
  const sidebarCollapsed = ref(false)
  const isDark = useDark()
  const toggleDark = useToggle(isDark)

  function toggleSidebar() {
    sidebarCollapsed.value = !sidebarCollapsed.value
  }

  function toggleTheme() {
    toggleDark()
  }

  function initTheme() {
    // useDark handles initialization automatically
  }

  return {
    sidebarCollapsed,
    isDark,
    toggleSidebar,
    toggleTheme,
    initTheme
  }
}, {
  persist: {
    paths: ['sidebarCollapsed']
  }
})
