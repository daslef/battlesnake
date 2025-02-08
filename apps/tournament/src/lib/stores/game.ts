import { create } from 'zustand'
import { devtools, persist } from 'zustand/middleware'
import type {} from '@redux-devtools/extension'

import { engineEventToFrame, type Frame, PlaybackMode, type GameServerObject, Game } from '../types'

interface GameStore {
  game: Game | null
  frames: Frame[]
  isFinished: boolean
  sendGameInfo: () => Promise<GameServerObject>
  reset: () => void
  setResult: () => void
  setGameResult: (game: Game, frames: Frame[]) => void
  setIsFinished: () => void
}

const useGameStore = create<GameStore>()(
  devtools(
    // persist(
    (set, get) => ({
      isFinished: false,
      frames: [],
      game: null,

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

      setResult: () => {
        set((state) => {
          const currentGame = state.game

          if (currentGame === null) {
            return state
          }

          const firstPlaceSnakeName =
            state.frames.at(-1)?.snakes.find((participant) => !participant.isEliminated)?.name ??
            null
          const firstPlaceParticipant =
            currentGame!.gameParticipants.find(
              (participant) => participant.snakeName === firstPlaceSnakeName
            ) ?? null

          const lastFrameWithTwoSnakes = state.frames
            .toReversed()
            .find((frame) => frame.snakes.filter((snake) => !snake.isEliminated).length === 2)

          const secondPlaceSnakeName =
            lastFrameWithTwoSnakes?.snakes.filter(
              (participant) => !participant.isEliminated && participant.name !== firstPlaceSnakeName
            )[0]?.name ?? null

          const secondPlaceParticipant =
            currentGame!.gameParticipants.find(
              (participant) => participant.snakeName === secondPlaceSnakeName
            ) ?? null

          currentGame.result = {
            firstPlace: firstPlaceParticipant,
            secondPlace: secondPlaceParticipant,
            aggressiveBonus: null
          }

          return {
            game: currentGame
          }
        })
      },

      setIsFinished: () => {
        set(() => ({ isFinished: true }))
      },

      reset: () => {
        set(() => ({
          frames: [],
          isFinished: false,
          game: null
        }))
      }
    }),
    {
      name: 'game-storage'
    }
  )
)

export { useGameStore }
