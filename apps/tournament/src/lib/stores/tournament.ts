import { create } from 'zustand'
import { devtools, persist } from 'zustand/middleware'
import type {} from '@redux-devtools/extension'

import jsSvg from '../../assets/logos/js.svg'
import pythonSvg from '../../assets/logos/python.svg'
import csharpSvg from '../../assets/logos/csharp.svg'

import { k_combinations } from '../'
import { Stage, Language, GameStatus, Game, Snake, Field, Participant, Score } from '../types'

interface TournamentStore {
  stage: Stage
  fields: Field[]
  games: Game[]
  participants: Participant[]

  score: Record<Participant['snakeName'], Score>
  addToScore: (participant: Participant, scoreType: keyof Score) => void
  getSortedScore: () => Record<Participant['snakeName'], Score>

  initialize: () => void
  generateBrackets: () => void
  setStage: (stage: Stage) => void
}

const useTournamentStore = create<TournamentStore>()(
  devtools(
    persist(
      (set, get) => ({
        stage: Stage.NOT_SET,
        fields: [
          { width: 8, height: 8 },
          { width: 15, height: 15 },
          { width: 25, height: 25 }
        ],
        games: [],
        participants: [
          {
            snakeName: '.conCat()',
            snakeUrl: new URL('/', 'http://10.11.21.177:6001'),
            snakeAuthor: '.conCat()',
            language: Language.Python,
            languageIcon: pythonSvg
          },
          {
            snakeName: 'Пушик',
            snakeUrl: new URL('/', 'http://10.11.21.177:6002'),
            snakeAuthor: 'Yarik',
            language: Language.JavaScript,
            languageIcon: jsSvg,
            photo: new URL('/vite.svg', 'http://localhost:5432')
          },
          {
            snakeName: 'camomile-snake',
            snakeUrl: new URL('/', 'http://10.11.21.177:6003'),
            snakeAuthor: 'Roman',
            language: Language.Python,
            languageIcon: pythonSvg
          },
          {
            snakeName: 'Генадий',
            snakeUrl: new URL('/', 'http://10.11.21.177:6004'),
            snakeAuthor: 'YuraPura',
            language: Language.JavaScript,
            languageIcon: jsSvg
          },
          {
            snakeName: 'power-shark',
            snakeUrl: new URL('/', 'http://10.11.21.177:6005'),
            snakeAuthor: 'NinthAcolite',
            language: Language.Python,
            languageIcon: pythonSvg
          },
          {
            snakeName: 'ss',
            snakeUrl: new URL('/', 'http://10.11.21.177:6006'),
            snakeAuthor: 'Никита / Марина',
            language: Language.Python,
            languageIcon: pythonSvg
          },
          {
            snakeName: 'collequesnake',
            snakeUrl: new URL('/', 'http://10.11.21.177:6006'),
            snakeAuthor: 'IT-Colleques',
            language: Language.Python,
            languageIcon: pythonSvg
          }
        ],

        score: {},

        addToScore: (participant, scoreType) => {
          set((state) => {
            const stageMultiplier = {
              [Stage.NOT_SET]: 1,
              [Stage.GROUP_THREES]: 1,
              [Stage.GROUP_FIVES]: 1.5,
              [Stage.GROUP_ALL]: 2
            }[state.stage]

            const scoreCopy = { ...state.score }
            const participantScore = scoreCopy[participant.snakeName]
            if (participantScore) {
              scoreCopy[participant.snakeName] = {
                ...scoreCopy[participant.snakeName],
                [scoreType]:
                  scoreCopy[participant.snakeName][scoreType] + 1 * (stageMultiplier ?? 1)
              }
              const { firstPlaces, secondPlaces, aggressiveBonuses } =
                scoreCopy[participant.snakeName]
              scoreCopy[participant.snakeName] = {
                ...scoreCopy[participant.snakeName],
                total: firstPlaces * 2 + secondPlaces * 1 + aggressiveBonuses * 0.001
              }
              return { score: scoreCopy }
            }
            return state
          })
        },

        getSortedScore: () =>
          Object.fromEntries(
            Object.entries(get().score).toSorted((a, b) => b[1].total - a[1].total)
          ),

        generateBrackets: () => {
          set((state) => {
            const combs_three = k_combinations(state.participants, 3).toSorted(
              () => Math.random() - 0.5
            )

            const groupGames = []

            for (const comb of combs_three) {
              groupGames.push({
                id: Date.now().toString() + Math.random().toString(),
                stage: Stage.GROUP_THREES,
                field: state.fields[0],
                status: GameStatus.NOT_PLAYED,
                gameParticipants: comb
              })
            }

            return { games: groupGames }
          })

          set((state) => {
            const combs_fives = k_combinations(state.participants, 5).toSorted(
              () => Math.random() - 0.5
            )

            const groupGames = [...state.games]

            for (const comb of combs_fives) {
              groupGames.push({
                id: Date.now().toString() + Math.random().toString(),
                stage: Stage.GROUP_FIVES,
                field: state.fields[1],
                status: GameStatus.NOT_PLAYED,
                gameParticipants: comb
              })
            }

            return { games: groupGames }
          })

          set((state) => {
            const finalMatches = [...state.games]
            for (let i = 0; i < 10; i++) {
              finalMatches.push({
                id: Date.now().toString() + Math.random().toString(),
                stage: Stage.GROUP_ALL,
                field: state.fields[2],
                status: GameStatus.NOT_PLAYED,
                gameParticipants: state.participants
              })
            }

            return { games: finalMatches }
          })
        },

        setStage: (stage) => set(() => ({ stage })),

        initialize: () => {
          if (!get().games.length) {
            get().setStage(Stage.GROUP_THREES)
            get().generateBrackets()
          }

          set((state) => ({
            score: Object.fromEntries(
              state.participants.map((participant) => [
                participant.snakeName,
                { firstPlaces: 0, secondPlaces: 0, aggressiveBonuses: 0, total: 0 }
              ])
            )
          }))
        }
      }),
      {
        name: 'tournament-storage'
      }
    )
  )
)

export { useTournamentStore }
