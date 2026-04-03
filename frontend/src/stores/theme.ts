import { defineStore } from 'pinia'
import { ref } from 'vue'

export const useThemeStore = defineStore('theme', () => {
  const isDark = ref(true)
  const locale = ref('zh')

  function toggleTheme() {
    isDark.value = !isDark.value
    applyTheme()
  }

  function setLocale(lang: string) {
    locale.value = lang
    localStorage.setItem('locale', lang)
  }

  function applyTheme() {
    if (isDark.value) {
      document.documentElement.setAttribute('data-theme', 'dark')
    } else {
      document.documentElement.removeAttribute('data-theme')
    }
    localStorage.setItem('theme', isDark.value ? 'dark' : 'light')
  }

  function init() {
    const storedTheme = localStorage.getItem('theme')
    const storedLocale = localStorage.getItem('locale')
    if (storedTheme !== null) {
      isDark.value = storedTheme === 'dark'
    }
    if (storedLocale) {
      locale.value = storedLocale
    }
    applyTheme()
  }

  init()

  return {
    isDark,
    locale,
    toggleTheme,
    setLocale,
  }
})
