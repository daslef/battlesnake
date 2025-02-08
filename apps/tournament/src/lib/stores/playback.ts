import { create } from 'zustand'
import { devtools, persist } from 'zustand/middleware'
import type {} from '@redux-devtools/extension'

import { loadGameInfo, setWs } from '../loaders/loadGame'
import { startPlayback, stopPlayback } from '../playback/animation'
import {
  type EngineGameInfo,
  type EngineGameEvent,
  type EngineSnake,
  type Frame,
  type Game,
  type GameServerObject,
  PlaybackMode
} from '../types'

interface PlaybackStore {
  loadedFrames: Set<Frame>
  frames: Frame[]
  currentFrameIndex: number
  setCurrentFrame: (index: number) => void
  checkIfFirstFrame: (frame: Frame) => void

  mode: PlaybackMode
  setMode: (mode: PlaybackMode) => void

  playbackError: string | null
  isLoading: boolean
  isFinished: boolean

  game: Game | null
  setGame: (game: Game) => void

  sendGameInfo: () => Promise<GameServerObject>
  setIsFinished: () => void

  title: string
  load: (engineURL: string) => void
  reset: () => void

  firstFrame: () => void
  lastFrame: () => void
  prevFrame: () => void
  nextFrame: () => void
  play: () => void
  pause: () => void
  togglePlayPause: () => void
}

const engineEventToFrame = (
  engineGameInfo: EngineGameInfo,
  engineGameEvent: EngineGameEvent
): Frame => ({
  turn: engineGameEvent.Turn,
  width: engineGameInfo.Game.Width,
  height: engineGameInfo.Game.Height,
  snakes: engineGameEvent.Snakes.map((engineSnake: EngineSnake) => ({
    // Fixed properties
    id: engineSnake.ID,
    name: engineSnake.Name,
    author: engineSnake.Author,
    color: engineSnake.Color,
    head: engineSnake.HeadType,
    tail: engineSnake.TailType,
    // Frame specific
    health: engineSnake.Health,
    latency: engineSnake.Latency,
    body: engineSnake.Body.map((engineCoords) => ({ x: engineCoords.X, y: engineCoords.Y })),
    length: engineSnake.Body.length,
    elimination: engineSnake.Death
      ? {
          turn: engineSnake.Death.Turn,
          cause: engineSnake.Death.Cause,
          by: engineSnake.Death.EliminatedBy
        }
      : null,
    // Helpers
    isEliminated: engineSnake.Death != null
  })),
  food: engineGameEvent.Food.map((engineCoords) => ({
    x: engineCoords.X,
    y: engineCoords.Y
  })),
  isFinalFrame: false
})

const usePlaybackStore = create<PlaybackStore>()(
  devtools(
    // persist(
    (set, get) => ({
      isLoading: false,
      isFinished: false,
      game: null,
      loadedFrames: new Set(),
      frames: [],
      currentFrameIndex: 0,
      playbackError: null,
      title: '',
      mode: PlaybackMode.PAUSED,

      checkIfFirstFrame: (frame: Frame) => {
        if (frame.turn !== 0) {
          return
        }

        set(() => ({
          currentFrameIndex: frame.turn,
          frame: frame,
          isLoading: false,
          mode: PlaybackMode.PAUSED
        }))
      },

      setFrames: (frames: Frame[]) => set({ frames }),

      load: async (engineURL) => {
        get().reset()
        set(() => ({ isLoading: true }))

        try {
          const gameInfo = await loadGameInfo(engineURL)
          const ws = await setWs(engineURL)

          ws!.onmessage = (message: { data: string }) => {
            const engineEvent = JSON.parse(message.data)

            if (engineEvent.Type == 'game_end') {
              get().frames.at(-1)!.isFinalFrame = true
            } else if (
              engineEvent.Type == 'frame' &&
              !get().loadedFrames.has(engineEvent.Data.Turn)
            ) {
              const frame = engineEventToFrame(gameInfo, engineEvent.Data)

              set(() => ({
                loadedFrames: new Set([...get().loadedFrames, engineEvent.Data.Turn]),
                frames: [...get().frames, frame].toSorted((a: Frame, b: Frame) => a.turn - b.turn)
              }))

              get().checkIfFirstFrame(frame)
            }
          }
        } catch (error) {
          console.error(error)
          set(() => ({ playbackError: (error as Error).message, isLoading: false }))
        }
      },

      sendGameInfo: async () => {
        if (!get().game) {
          return Promise.reject('game does not set')
        }

        const { gameParticipants, field } = get().game!

        const snakes = gameParticipants.map(({ snakeName, snakeUrl }) => ({
          snakeName,
          snakeUrl
        }))

        fetch('http://localhost:5001/new', {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json'
          },
          body: JSON.stringify({
            snakes,
            field
          })
        })
      },

      setGame: (game) => set(() => ({ game })),

      reset: () => {
        set(() => ({
          frames: [],
          loadedFrames: new Set(),
          currentFrameIndex: 0,
          playbackError: null,
          mode: PlaybackMode.PAUSED,
          isLoading: false,
          title: '',
          isFinished: false,
          game: null
        }))
      },

      setCurrentFrame: (index) => {
        set((state) => ({
          currentFrameIndex: Math.min(Math.max(index, 0), state.frames.length - 1)
        }))

        if (
          get().frames[get().currentFrameIndex].isFinalFrame &&
          get().mode == PlaybackMode.PLAYING
        ) {
          stopPlayback()
          get().setMode(PlaybackMode.PAUSED) // TODO?
        }
      },

      setMode: (mode) => {
        set((state) => ({ ...state, mode }))
      },

      firstFrame: () => {
        get().setMode(PlaybackMode.PAUSED)
        set(() => ({ mode: PlaybackMode.PAUSED, currentFrameIndex: 0 }))
      },

      lastFrame: () => {
        set((state) => ({ mode: PlaybackMode.PAUSED, currentFrameIndex: state.frames.length - 1 }))
      },

      prevFrame: () => {
        if (get().currentFrameIndex === 0) {
          return
        }
        set((state) => ({
          mode: PlaybackMode.PAUSED,
          currentFrameIndex: state.currentFrameIndex - 1
        }))
      },

      nextFrame: () => {
        if (get().currentFrameIndex === get().frames.length - 1) {
          return
        }
        set((state) => ({
          mode: PlaybackMode.PAUSED,
          currentFrameIndex: state.currentFrameIndex + 1
        }))
      },

      play: () => {
        const fps = 10
        if (get().mode === PlaybackMode.PLAYING) {
          return
        }
        startPlayback(fps, () => {
          get().setCurrentFrame(get().currentFrameIndex + 1)
        })
        get().setMode(PlaybackMode.PLAYING)
      },

      pause: () => {
        stopPlayback()
        get().setMode(PlaybackMode.PAUSED)
      },

      togglePlayPause: () => {
        if (get().mode === PlaybackMode.PAUSED) {
          get().play()
        } else if (get().mode == PlaybackMode.PLAYING) {
          get().pause()
        }
      }
    }),
    {
      name: 'playback-storage'
    }
  )
)

export { usePlaybackStore }
