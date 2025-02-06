import { create } from 'zustand'
import { devtools, persist } from 'zustand/middleware'
import type {} from '@redux-devtools/extension'

import { loadGameInfo, setWs } from '../loaders/loadGame'
import { startPlayback, stopPlayback } from '../playback/animation'
import { engineEventToFrame, type Frame, PlaybackMode, type GameServerObject, Game } from '../types'

interface PlaybackStore {
  loadedFrames: Set<Frame>
  frames: Frame[]
  currentFrameIndex: number
  mode: PlaybackMode
  playbackError: string | null
  isLoading: boolean
  title: string
  load: (engineURL: string) => void
  sendGameInfo: (gameObject: Game) => Promise<GameServerObject>
  reset: () => void
  setCurrentFrame: (index: number) => void
  setMode: (mode: PlaybackMode) => void
  firstFrame: () => void
  lastFrame: () => void
  prevFrame: () => void
  nextFrame: () => void
  play: () => void
  pause: () => void
  togglePlayPause: () => void
  jumpToFrame: (index: number) => void
}

const usePlaybackStore = create<PlaybackStore>()(
  devtools(
    // persist(
    (set, get) => ({
      isLoading: false,
      loadedFrames: new Set(),
      frames: [],
      currentFrameIndex: 0,
      playbackError: null,
      title: '',
      mode: PlaybackMode.NEW,
      setFrames: (frames: Frame[]) => set({ frames }),
      load: async (engineURL) => {
        get().reset()
        set(() => ({ isLoading: true }))

        console.debug(`[playback] loading game`)

        try {
          const gameInfo = await loadGameInfo(engineURL)
          const ws = await setWs(engineURL)

          ws!.onopen = () => {
            console.debug('[playback] opening engine websocket')
          }

          ws!.onmessage = (message: { data: string }) => {
            const engineEvent = JSON.parse(message.data)

            if (engineEvent.Type == 'game_end') {
              console.debug('[playback] received final frame')
              get().frames[get().frames.length - 1].isFinalFrame = true
            } else if (
              engineEvent.Type == 'frame' &&
              !get().loadedFrames.has(engineEvent.Data.Turn)
            ) {
              const frame = engineEventToFrame(gameInfo, engineEvent.Data)

              set(() => ({
                loadedFrames: new Set([...get().loadedFrames, engineEvent.Data.Turn]),
                frames: [...get().frames, frame].toSorted((a: Frame, b: Frame) => a.turn - b.turn)
              }))

              if (frame.turn === 0) {
                set(() => ({
                  currentFrameIndex: frame.turn,
                  frame: frame,
                  isLoading: false,
                  mode: PlaybackMode.READY
                }))
              }
            }
          }

          ws!.onclose = () => {
            console.debug('[playback] closing engine websocket')
          }
        } catch (error) {
          console.error(error)
          set(() => ({ playbackError: (error as Error).message, isLoading: false }))
        }
      },
      sendGameInfo: async (game: Game) => {
        fetch('http://localhost:5001/new', {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json'
          },
          body: JSON.stringify({
            snakes: game.gameParticipants.map(({ snakeName, snakeUrl }) => ({
              snakeName,
              snakeUrl
            })),
            field: game.field
          })
        })
      },
      reset: () => {
        set(() => ({
          frames: [],
          loadedFrames: new Set(),
          currentFrameIndex: 0,
          playbackError: null,
          mode: PlaybackMode.NEW,
          isLoading: false,
          title: ''
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
          get().setMode(PlaybackMode.FINISHED)
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
        if (![PlaybackMode.PAUSED, PlaybackMode.READY].includes(get().mode)) {
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
        if ([PlaybackMode.PAUSED, PlaybackMode.READY].includes(get().mode)) {
          get().play()
        } else if (get().mode == PlaybackMode.PLAYING) {
          get().pause()
        }
      },
      jumpToFrame: (i: number) => {
        get().pause()
        get().setCurrentFrame(i)
      }
    }),
    {
      name: 'playback-storage'
    }
  )
)

export { usePlaybackStore }
