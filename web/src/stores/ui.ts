import { defineStore } from 'pinia'
import { ref, watch } from 'vue'
import { useDark, useToggle } from '@vueuse/core'
import i18n, { detectBrowserLanguage, type Locale } from '@/i18n'

export const useUiStore = defineStore('ui', () => {
  const sidebarCollapsed = ref(false)
  const language = ref<Locale>(detectBrowserLanguage())
  const isDark = useDark()
  const toggleDark = useToggle(isDark)

  watch(language, (newLang) => {
    i18n.global.locale.value = newLang
  }, { immediate: true })

  watch(sidebarCollapsed, (collapsed) => {
    document.documentElement.style.setProperty('--pm-sidebar-width', collapsed ? '64px' : '220px')
  }, { immediate: true })

  function toggleSidebar() {
    sidebarCollapsed.value = !sidebarCollapsed.value
  }

  function toggleTheme() {
    toggleDark()
  }

  function setLanguage(lang: Locale) {
    language.value = lang
  }

  function initTheme() {
    // useDark handles initialization automatically
  }

  return {
    sidebarCollapsed,
    language,
    isDark,
    toggleSidebar,
    toggleTheme,
    setLanguage,
    initTheme
  }
}, {
  persist: {
    paths: ['sidebarCollapsed', 'language']
  }
})
