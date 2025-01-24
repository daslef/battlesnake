import { create } from 'zustand'
import { devtools, persist } from 'zustand/middleware'
import type { } from '@redux-devtools/extension' // required for devtools typing
import { Theme } from '../types'

interface SettingsState {
  engine: string
  fps: number
  game: string
  showControls: boolean
  showCoords: boolean
  showScoreboard: boolean
  theme: Theme
  title: string
  turn: number
  setEngine: (engineUrl: string) => void
}

const useSettingsStore = create<SettingsState>()(
  devtools(
    persist(
      (set) => ({
        engine: 'http://localhost:5000',
        fps: 6,
        game: '',
        showControls: true,
        showCoords: false,
        showScoreboard: true,
        theme: Theme.DARK,
        title: '',
        turn: 0,
        setEngine: (engineUrl) => set(() => ({ engine: engineUrl })),
      }),
      {
        name: 'settings-storage',
      },
    ),
  ),
)

export { useSettingsStore }