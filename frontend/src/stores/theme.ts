// src/stores/theme.ts
import { defineStore } from 'pinia'
import { ref, computed } from 'vue'

export type Theme = 'light' | 'dark'

// Standalone function to get initial theme (can be used before store initialization)
export const getInitialTheme = (): Theme => {
    const saved = localStorage.getItem('theme') as Theme | null
    if (saved === 'light' || saved === 'dark') return saved

    if (window.matchMedia('(prefers-color-scheme: dark)').matches) {
        return 'dark'
    }
    return 'light'
}

// Standalone function to apply theme to DOM
export const applyThemeToDOM = (theme: Theme): void => {
    document.documentElement.setAttribute('data-theme', theme)
}

export const useThemeStore = defineStore('theme', () => {
    const theme = ref<Theme>(getInitialTheme())

    const applyTheme = (newTheme: Theme) => {
        theme.value = newTheme
        applyThemeToDOM(newTheme)
        localStorage.setItem('theme', newTheme)
    }

    const toggleTheme = () => {
        applyTheme(theme.value === 'light' ? 'dark' : 'light')
    }

    // Apply theme on init
    applyTheme(theme.value)

    return {
        theme,
        toggleTheme,
        applyTheme,
        isDark: computed(() => theme.value === 'dark'),
    }
})