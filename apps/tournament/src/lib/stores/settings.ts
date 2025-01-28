import { create } from 'zustand'
import { devtools, persist } from 'zustand/middleware'
import type { } from '@redux-devtools/extension' // required for devtools typing
import { Theme } from '../types'

interface SettingsState {
  engineUrl: string
  fps: number
  showScoreboard: boolean
  title: string
  error: string | null
  setError: (errorMessage: string) => void
  setEngine: (engineUrl: string) => void
}

interface ThemeState {
  theme: Theme
  setTheme: (theme: Theme) => void
}

const useSettingsStore = create<SettingsState>()(
  devtools(
    // persist(
    (set) => ({
      engineUrl: 'http://localhost:5000',
      fps: 6,
      showScoreboard: true,
      error: null,
      setError: (errorMessage) => set((state) => ({ ...state, error: errorMessage })),
      setEngine: (engineUrl) => set((state) => ({ ...state, engineUrl }))
    }),
    {
      name: 'settings-storage'
    }
    // )
  )
)

const useThemeStore = create<ThemeState>()(
  devtools(
    persist(
      (set) => ({
        theme: Theme.DARK,
        setTheme: (theme) => {
          set((state) => ({ ...state, theme: theme }))
        }
      }),
      { name: 'theme-storage' }
    )
  )
)

export { useSettingsStore, useThemeStore }
