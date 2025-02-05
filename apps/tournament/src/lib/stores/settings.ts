import { create } from 'zustand'
import { devtools, persist } from 'zustand/middleware'
import type {} from '@redux-devtools/extension' // required for devtools typing

interface SettingsState {
  engineUrl: string
  fps: number
  showScoreboard: boolean
}

const useSettingsStore = create<SettingsState>()(
  devtools(
    // persist(
    (set) => ({
      engineUrl: 'http://localhost:5000',
      fps: 6,
      showScoreboard: true
    }),
    {
      name: 'settings-storage'
    }
    // )
  )
)

export { useSettingsStore }
