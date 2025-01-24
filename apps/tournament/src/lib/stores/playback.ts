import { create } from 'zustand'
import { devtools, persist } from 'zustand/middleware'
import type { } from '@redux-devtools/extension'

import { startPlayback, stopPlayback } from '../playback/animation'
import { engineEventToFrame, type Frame, PlaybackMode } from '../types'


interface PlaybackStore {
    loadedFrames: Set<Frame>
    frames: Frame[]
    currentFrameIndex: number
    currentFrame: Frame | null
    mode: PlaybackMode
    finalFrame: null | Frame
    playbackError: string | null
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
    onFrameLoad: (frame: Frame) => void
    onFinalFrame: (frame: Frame) => void
}

const usePlaybackStore = create<PlaybackStore>()(
    devtools(
        persist(
            (set, get) => ({
                loadedFrames: new Set(),
                frames: [],
                currentFrameIndex: 0,
                currentFrame: null,
                playbackError: null,
                finalFrame: null,
                mode: PlaybackMode.PAUSED,
                load: async (engineURL, gameID) => {
                    const wsUrl = engineURL
                        .replace(/^https:\/\//i, 'wss://')
                        .replace(/^http:\/\//i, 'ws://')

                    const gameInfoUrl = `${engineURL}/games/${gameID}`
                    const gameEventsUrl = `${wsUrl}/games/${gameID}/events`

                    console.debug(`[playback] loading game ${gameID}`)

                    set(state => ({ ...state, loadedFrames: new Set() }))

                    try {
                        const response = await fetch(gameInfoUrl)

                        if (response.status == 404) {
                            throw new Error('Game not found')
                        } else if (!response.ok) {
                            throw new Error('Error loading game')
                        }

                        const gameInfo = await response.json()
                        const ws = new WebSocket(gameEventsUrl)

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

                                get().onFrameLoad(frame)
                            } else if (engineEvent.Type == 'game_end') {
                                console.debug('[playback] received final frame')

                                get().frames[frames.length - 1].isFinalFrame = true
                                get().onFinalFrame(get().frames[frames.length - 1])
                            }
                        }

                        ws.onclose = () => {
                            console.debug('[playback] closing engine websocket')
                        }

                    } catch (e: unknown) {
                        console.error(e)
                        set(state => ({ ...state, playbackError: (e as Error).message }))
                    }
                },
                reset: () => {
                    set((state) => ({ ...state, frames: [], currentFrameIndex: 0, playbackError: null }))
                },
                setCurrentFrame: (index) => {
                    set((state) => {
                        state.currentFrameIndex = Math.min(Math.max(index, 0), state.frames.length - 1)
                        state.currentFrame = state.frames[state.currentFrameIndex]

                        if (state.currentFrame.isFinalFrame && state.mode == PlaybackMode.PLAYING) {
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
                onFrameLoad: (frame: Frame) => {
                    // Load the first frame when we see it.
                    if (frame.turn !== 0) {
                        return
                    }

                    set((state) => {
                        get().setCurrentFrame(frame.turn)
                        return {
                            ...state,
                            frame: frame,
                            mode: PlaybackMode.PAUSED,
                            finalFrame: null
                        }
                    })
                },
                onFinalFrame: (frame: Frame) => {
                    set(state => ({ ...state, finalFrame: frame, mode: PlaybackMode.FINISHED }))
                },
            }),
            {
                name: 'playback-storage',
            },
        ),
    ),
)

export { usePlaybackStore }

