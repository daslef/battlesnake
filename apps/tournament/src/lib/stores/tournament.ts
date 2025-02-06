import { create } from 'zustand'
import { devtools, persist } from 'zustand/middleware'
import type {} from '@redux-devtools/extension'

import jsSvg from '../../assets/logos/js.svg'
import pythonSvg from '../../assets/logos/python.svg'
import csharpSvg from '../../assets/logos/csharp.svg'

import { k_combinations } from '../'
import { TournamentStore, Stage, Language, GameStatus, Game, Snake } from '../types'

const useTournamentStore = create<TournamentStore>()(
  devtools(
    // persist(
    (set, get) => ({
      stage: Stage.NOT_SET,
      fields: [
        { width: 8, height: 8 },
        { width: 11, height: 11 },
        { width: 19, height: 19 }
      ],
      games: [],
      participants: [
        {
          snakeName: 'yaro-snake',
          snakeUrl: new URL('/', 'http://localhost:9001'),
          snakeAuthor: 'Ярослав',
          language: Language.JavaScript,
          languageIcon: jsSvg
        },
        {
          snakeName: 'ri-snake',
          snakeUrl: new URL('/', 'http://localhost:9002'),
          snakeAuthor: 'РИ-9',
          language: Language.CSharp,
          languageIcon: csharpSvg,
          photo: new URL('/vite.svg', 'http://localhost:5432')
        },
        {
          snakeName: 'guest-snake',
          snakeUrl: new URL('/', 'https://snapepy.onrender.com'),
          snakeAuthor: 'Гость',
          language: Language.CSharp,
          languageIcon: csharpSvg
        },
        {
          snakeName: 'another-snake',
          snakeUrl: new URL('/', 'http://localhost:9005'),
          snakeAuthor: 'Anothers',
          language: Language.Python,
          languageIcon: pythonSvg
        },
        {
          snakeName: 'puppies-snake',
          snakeUrl: new URL('/', 'http://localhost:9003'),
          snakeAuthor: 'Щеночки',
          language: Language.JavaScript,
          languageIcon: jsSvg
        }
      ],

      score: new Map(),

      addToScore: (participant, scoreType) => {
        set((state) => {
          const scoreCopy = new Map(state.score)
          const participantScore = scoreCopy.get(participant)
          if (participantScore) {
            scoreCopy.set(participant, {
              ...scoreCopy.get(participant)!,
              [scoreType]: scoreCopy.get(participant)![scoreType] + 1
            })
            return { score: scoreCopy }
          }
          return state
        })
      },

      setGameResult: (game, frames) => {
        const firstPlaceSnakeName =
          frames.at(-1)?.snakes.find((participant) => !participant.isEliminated)?.name ?? null
        const firstPlaceParticipant =
          game.gameParticipants.find(
            (participant) => participant.snakeName === firstPlaceSnakeName
          ) ?? null

        const lastFrameWithTwoSnakes = frames
          .toReversed()
          .find((frame) => frame.snakes.filter((snake) => !snake.isEliminated).length === 2)

        const secondPlaceSnakeName =
          lastFrameWithTwoSnakes?.snakes.filter(
            (participant) => !participant.isEliminated && participant.name !== firstPlaceSnakeName
          )[0]?.name ?? null

        const secondPlaceParticipant =
          game.gameParticipants.find(
            (participant) => participant.snakeName === secondPlaceSnakeName
          ) ?? null

        set((state) => {
          const gameIndex = state.games.findIndex((el) => el.id === game.id)

          if (gameIndex === -1) {
            return state
          }

          const updatedGame = {
            ...state.games[gameIndex],
            result: {
              firstPlace: firstPlaceParticipant,
              secondPlace: secondPlaceParticipant,
              aggressiveBonus: null
            }
          }

          return {
            games: [
              ...state.games.slice(0, gameIndex),
              updatedGame,
              ...state.games.slice(gameIndex + 1)
            ]
          }
        })
      },

      setGameStatus: (game: Game, status: GameStatus) => {
        set((state) => {
          const gameIndex = state.games.findIndex((el) => el.id === game.id)
          return {
            games: [
              ...state.games.slice(0, gameIndex),
              { ...state.games[gameIndex], status },
              ...state.games.slice(gameIndex + 1)
            ]
          }
        })
      },

      calculateTotalScore: (participant) => {
        const participantScores = get().score.get(participant)

        if (!participantScores) {
          throw new Error('Scores not found')
        }

        const { firstPlaces, secondPlaces, aggressiveBonuses } = participantScores
        return firstPlaces * 2 + secondPlaces * 1 + aggressiveBonuses * 0.001
      },

      getSortedScore: () =>
        new Map(Object.entries(get().score).toSorted((a, b) => b[1].result - a[1].result)),

      setGames: () => {
        if (get().stage === Stage.GROUP_THREES) {
          set((state) => {
            const combs_three = k_combinations(state.participants, 3).toSorted(
              () => Math.random() - 0.5
            )

            const groupGames = []

            for (const comb of combs_three) {
              groupGames.push({
                id: Date.now().toString() + Math.random().toString(),
                stage: state.stage,
                field: { width: 8, height: 8 },
                status: GameStatus.NOT_PLAYED,
                gameParticipants: comb
              })
            }

            return { games: groupGames }
          })
        } else if (get().stage === Stage.GROUP_FIVES) {
          set((state) => {
            const combs_fives = k_combinations(state.participants, 5).toSorted(
              () => Math.random() - 0.5
            )

            const groupGames = [...state.games]

            // fives from last commit
            for (const comb of combs_fives) {
              for (const _ of Array(5)) {
                groupGames.push({
                  id: Date.now().toString() + Math.random().toString(),
                  stage: state.stage,
                  field: { width: 11, height: 11 },
                  status: GameStatus.NOT_PLAYED,
                  gameParticipants: comb
                })
              }
            }

            return { games: groupGames }
          })
        } else if (get().stage === Stage.FINALS) {
          set((state) => {
            const finalMatches = []
            for (const field of state.fields) {
              for (let i = 0; i < 3; i++) {
                finalMatches.push({
                  id: Date.now().toString() + Math.random().toString(),
                  stage: state.stage,
                  field: field,
                  status: GameStatus.NOT_PLAYED,
                  gameParticipants: [...state.getSortedScore().keys()].slice(0, 2)
                })
              }
            }

            return { games: { ...state.games, finalMatches } }
          })
        }
      },

      setStage: (stage) => set(() => ({ stage })),

      initializeScore: () =>
        set((state) => ({
          score: new Map(
            state.participants.map((participant) => [
              participant,
              { firstPlaces: 0, secondPlaces: 0, aggressiveBonuses: 0, total: 0 }
            ])
          )
        })),

      initialize: () => {
        get().setStage(Stage.GROUP_THREES)
        get().setGames()
      }
    }),
    {
      name: 'tournament-storage'
    }
  )
  // )
)

export { useTournamentStore }
