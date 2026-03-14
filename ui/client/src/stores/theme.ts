import { defineStore } from 'pinia'
import { ref, watch } from 'vue'

const STORAGE_KEY = 'beehive-theme'

export type ThemeName = 'light' | 'dark' | null

function loadSavedTheme(): ThemeName {
  const s = localStorage.getItem(STORAGE_KEY)
  if (s === 'dark') return 'dark'
  if (s === 'light') return 'light'
  return null
}

export const useThemeStore = defineStore('theme', () => {
  const theme = ref<ThemeName>(loadSavedTheme())

  watch(
    theme,
    (v) => {
      if (v) localStorage.setItem(STORAGE_KEY, v)
    },
    { immediate: true }
  )

  function setTheme(value: ThemeName) {
    theme.value = value
  }

  return { theme, setTheme }
})
