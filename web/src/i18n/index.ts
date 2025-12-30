import { createI18n } from 'vue-i18n'
import en from './locales/en.json'
import zh from './locales/zh.json'

export type Locale = 'en' | 'zh'

export function detectBrowserLanguage(): Locale {
  const browserLang = navigator.language || 'en'
  return browserLang.startsWith('zh') ? 'zh' : 'en'
}

export function getStoredLanguage(): Locale | null {
  try {
    const stored = localStorage.getItem('ui')
    if (stored) {
      const parsed = JSON.parse(stored)
      if (parsed.language === 'zh' || parsed.language === 'en') {
        return parsed.language
      }
    }
  } catch {
    // ignore
  }
  return null
}

const i18n = createI18n({
  legacy: false,
  locale: getStoredLanguage() || detectBrowserLanguage(),
  fallbackLocale: 'en',
  messages: { en, zh }
})

export default i18n
