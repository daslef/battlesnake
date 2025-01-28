import { create } from 'zustand'
import { devtools, persist } from 'zustand/middleware'
import type { } from '@redux-devtools/extension'

import { loadGameInfo, loadGameEvents } from '../loaders/loadGame'
import { startPlayback, stopPlayback } from '../playback/animation'
import { engineEventToFrame, type Frame, PlaybackMode } from '../types'


interface PlaybackStore {
    loadedFrames: Set<Frame>
    frames: Frame[]
    currentFrameIndex: number
    mode: PlaybackMode
    playbackError: string | null
    isLoading: boolean
    load: (engineURL: string) => void
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
            mode: PlaybackMode.PAUSED,
            load: async (engineURL) => {
                get().reset()
                set(state => ({ ...state, loadedFrames: new Set(), isLoading: true }))

                console.debug(`[playback] loading game`)

                try {
                    const gameInfo = await loadGameInfo(engineURL)
                    const ws = await loadGameEvents(engineURL)
                    set(() => ({ isLoading: false }))

                    ws.onopen = () => {
                        console.debug('[playback] opening engine websocket')
                    }

                    ws.onmessage = (message: { data: string }) => {
                        const engineEvent = JSON.parse(message.data)

                        if (engineEvent.Data.Turn == 0) {
                            console.debug('[playback] received new game')
                        }

                        if (engineEvent.Type == 'frame' && !get().loadedFrames.has(engineEvent.Data.Turn)) {
                            const updatedLoadedFrames = [...get().loadedFrames, engineEvent.Data.Turn]
                            set(() => ({ loadedFrames: new Set(updatedLoadedFrames) }))

                            const frame = engineEventToFrame(gameInfo, engineEvent.Data)
                            const updatedFrames = [...get().frames, frame].toSorted((a: Frame, b: Frame) => a.turn - b.turn)
                            set(() => ({ frames: updatedFrames }))

                            if (frame.turn !== 0) {
                                return
                            }

                            get().setCurrentFrame(frame.turn)

                            set(() => ({
                                frame: frame,
                                mode: PlaybackMode.PAUSED,
                            }))

                        } else if (engineEvent.Type == 'game_end') {
                            console.debug('[playback] received final frame')

                            get().frames[get().frames.length - 1].isFinalFrame = true
                        }
                    }

                    ws.onclose = () => {
                        console.debug('[playback] closing engine websocket')
                    }

                } catch (error) {
                    console.error(error)
                    set(state => ({ ...state, playbackError: (error as Error).message, isLoading: false }))
                }

            },
            reset: () => {
                set((state) => ({ ...state, frames: [], currentFrameIndex: 0, playbackError: null }))
            },
            setCurrentFrame: (index) => {
                set((state) => ({ currentFrameIndex: Math.min(Math.max(index, 0), state.frames.length - 1) }))

                if (get().frames[get().currentFrameIndex].isFinalFrame && get().mode == PlaybackMode.PLAYING) {
                    stopPlayback()
                    set(() => ({ mode: PlaybackMode.FINISHED }))
                }
            },
            setMode: (mode) => {
                set((state) => ({ ...state, mode }))
            },
            firstFrame: () => {
                set(() => ({ mode: PlaybackMode.PAUSED, currentFrameIndex: 0 }))
            },
            lastFrame: () => {
                set((state) => ({ mode: PlaybackMode.PAUSED, currentFrameIndex: state.frames.length - 1 }))
            },
            prevFrame: () => {
                if (get().currentFrameIndex === 0) {
                    return
                }
                set((state) => ({ mode: PlaybackMode.PAUSED, currentFrameIndex: state.currentFrameIndex - 1 }))
            },
            nextFrame: () => {
                if (get().currentFrameIndex === get().frames.length - 1) {
                    return
                }
                set((state) => ({ mode: PlaybackMode.PAUSED, currentFrameIndex: state.currentFrameIndex + 1 }))
            },
            play: () => {
                const fps = 10
                if (get().mode !== PlaybackMode.PAUSED) {
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
                if (get().mode == PlaybackMode.PAUSED) {
                    get().play()
                } else if (get().mode == PlaybackMode.PLAYING) {
                    get().pause()
                }
            },
            jumpToFrame: (i: number) => {
                get().pause()
                get().setCurrentFrame(i)
            },
        }),
        {
            name: 'playback-storage',
        },
    ),
)

export { usePlaybackStore }

