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
    finalFrame: null | Frame
    playbackError: string | null
    isLoading: boolean
    load: (engineURL: string, gameID: string) => void
    reset: () => void
    setCurrentFrame: (index: number) => void
    setMode: (mode: PlaybackMode) => void
    firstFrame: () => void
    lastFrame: () => void
    prevFrame: () => void
    nextFrame: () => void
    prevEliminationFrame: () => void
    nextEliminationFrame: () => void
    play: () => void
    pause: () => void
    togglePlayPause: () => void
    jumpToFrame: (index: number) => void
}

const usePlaybackStore = create<PlaybackStore>()(
    devtools(
        persist(
            (set, get) => ({
                isLoading: false,
                loadedFrames: new Set(),
                frames: [],
                currentFrameIndex: 0,
                playbackError: null,
                finalFrame: null,
                mode: PlaybackMode.PAUSED,
                load: async (engineURL, gameID) => {
                    get().reset()
                    set(state => ({ ...state, loadedFrames: new Set(), isLoading: true }))

                    console.debug(`[playback] loading game ${gameID}`)

                    try {
                        const gameInfo = await loadGameInfo(engineURL, gameID)
                        const ws = await loadGameEvents(engineURL, gameID)
                        set(state => ({ ...state, isLoading: false }))

                        ws.onopen = () => {
                            console.debug('[playback] opening engine websocket')
                        }

                        ws.onmessage = (message: { data: string }) => {
                            const engineEvent = JSON.parse(message.data)

                            if (engineEvent.Data.Turn == 0) {
                                // loadedFrames = new Set()
                                // frames = []
                                console.debug('[playback] received new game')
                            }

                            if (engineEvent.Type == 'frame' && !get().loadedFrames.has(engineEvent.Data.Turn)) {
                                get().loadedFrames.add(engineEvent.Data.Turn)

                                const frame = engineEventToFrame(gameInfo, engineEvent.Data)
                                get().frames.push(frame)
                                get().frames.sort((a: Frame, b: Frame) => a.turn - b.turn)

                                if (frame.turn !== 0) {
                                    return
                                }

                                get().setCurrentFrame(frame.turn)

                                set((state) => {
                                    return {
                                        ...state,
                                        frame: frame,
                                        mode: PlaybackMode.PAUSED,
                                        finalFrame: null
                                    }
                                })


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
                    set((state) => {
                        state.currentFrameIndex = Math.min(Math.max(index, 0), state.frames.length - 1)

                        if (state.frames[state.currentFrameIndex].isFinalFrame && state.mode == PlaybackMode.PLAYING) {
                            stopPlayback()
                            state.mode = PlaybackMode.FINISHED
                        }

                        return state
                    })
                },
                setMode: (mode) => {
                    set((state) => ({ ...state, mode }))
                },
                firstFrame: () => {
                    set(state => {
                        if (state.mode == PlaybackMode.PAUSED) {
                            return { ...state, currentFrameIndex: 0 }
                        }
                        return state
                    })
                },
                lastFrame: () => {
                    set(state => {
                        if (state.mode == PlaybackMode.PAUSED) {
                            return { ...state, currentFrameIndex: state.frames.length - 1 }
                        }
                        return state
                    })
                },
                prevFrame: () => {
                    set(state => {
                        if (state.mode == PlaybackMode.PAUSED) {
                            return { ...state, currentFrameIndex: state.currentFrameIndex - 1 }
                        }
                        return state
                    })
                },
                nextFrame: () => {
                    set(state => {
                        if (state.mode == PlaybackMode.PAUSED) {
                            return { ...state, currentFrameIndex: state.currentFrameIndex + 1 }
                        }
                        return state
                    })
                },
                prevEliminationFrame: () => {
                    set(state => {
                        if (state.mode !== PlaybackMode.PAUSED) {
                            return state
                        }

                        for (let i = state.currentFrameIndex; i >= 0; i--) {
                            for (let s = 0; s < state.frames[i].snakes.length; s++) {
                                const snake = state.frames[i].snakes[s]
                                if (snake.elimination && snake.elimination.turn <= state.currentFrameIndex) {
                                    const newIndex = snake.elimination.turn - 1
                                    console.debug(`[playback] jump to elimination frame ${newIndex}`)
                                    state.setCurrentFrame(newIndex)
                                    break // return
                                }
                            }
                            state.firstFrame()
                        }

                        return state
                    })
                },
                nextEliminationFrame: () => {
                    set(state => {
                        if (state.mode !== PlaybackMode.PAUSED) {
                            return state
                        }

                        for (let i = state.currentFrameIndex + 2; i < state.frames.length; i++) {
                            for (let s = 0; s < state.frames[i].snakes.length; s++) {
                                const snake = state.frames[i].snakes[s]
                                if (snake.elimination && snake.elimination.turn > state.currentFrameIndex + 1) {
                                    const newIndex = snake.elimination.turn - 1
                                    console.debug(`[playback] jump to elimination frame ${newIndex}`)
                                    state.setCurrentFrame(newIndex)
                                    break // return
                                }
                            }
                        }

                        state.lastFrame()

                        return state
                    })
                },
                play: () => {
                    set(state => {
                        const fps = 10
                        if (state.mode !== PlaybackMode.PAUSED) {
                            return state
                        }
                        startPlayback(fps, () => {
                            state.setCurrentFrame(state.currentFrameIndex + 1)
                        })
                        state.setMode(PlaybackMode.PLAYING)
                        return state
                    })
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
    ),
)

export { usePlaybackStore }

