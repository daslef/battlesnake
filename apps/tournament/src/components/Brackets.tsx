import { type FC, useMemo } from 'react'
import { Dialog, Box, Button, Text, Flex } from '@radix-ui/themes'
import Board from './Board'

import { usePlaybackStore } from '../lib/stores/playback'
import { useTournamentStore } from '../lib/stores/tournament'

import digitOneSvg from '../assets/icons/digit-one.svg'
import digitTwoSvg from '../assets/icons/digit-two.svg'
import bowSvg from '../assets/icons/bow.svg'

import { Stage, type Game, GameStatus, type Participant, type Score } from '../lib/types'

interface IBracketMatch {
  game: Game
  heading: string
}

interface IBracketParticipant {
  participant: Participant
  isGameCompleted: boolean
}

const BracketParticipant: FC<IBracketParticipant> = ({ participant, isGameCompleted }) => {
  const addToScore = useTournamentStore((store) => store.addToScore)

  const handleClick = (type: keyof Score) => {
    addToScore(participant, type)
  }

  return (
    <tr className="match__participant participant">
      <td className="participant__image">
        <img
          src={participant.languageIcon}
          alt="programming language"
          className="participant__language"
        />
      </td>
      <td className="participant__team">{participant.snakeAuthor}</td>
      {!isGameCompleted && (
        <td className="participant__controls participant__controls">
          <button
            type="button"
            className="participant__controls__control control"
            onClick={() => handleClick('firstPlaces')}
          >
            <img
              src={digitOneSvg}
              alt="first place"
              className="control__icon control__icon--gold"
            />
          </button>
          <button
            type="button"
            className="participant__controls__control control"
            onClick={() => handleClick('secondPlaces')}
          >
            <img
              src={digitTwoSvg}
              alt="second place"
              className="control__icon control__icon--silver"
            />
          </button>
          <button
            type="button"
            className="participant__controls__control control"
            onClick={() => handleClick('aggressiveBonuses')}
          >
            <img src={bowSvg} alt="brave bonus" className="control__icon control__icon--bonus" />
          </button>
        </td>
      )}
    </tr>
  )
}

const BracketMatch: FC<IBracketMatch> = ({ game, heading }) => {
  const load = usePlaybackStore((store) => store.load)
  const setGame = usePlaybackStore((store) => store.setGame)
  const sendGameInfo = usePlaybackStore((store) => store.sendGameInfo)
  const isGameCompleted = game.status === GameStatus.COMPLETED

  return (
    <Dialog.Root>
      <Box className="bracket__match match">
        <Flex align={'center'} justify={'between'}>
          <Text size={'4'} weight={'bold'}>
            {heading}
          </Text>
          <Dialog.Trigger>
            {isGameCompleted ? (
              <></>
            ) : (
              <Button
                size={'1'}
                color="violet"
                style={{ cursor: 'pointer' }}
                onClick={() => {
                  setGame(game)
                  sendGameInfo().then(() => {
                    load('http://127.0.0.1:5000')
                  })
                }}
              >
                ▶︎
              </Button>
            )}
          </Dialog.Trigger>
        </Flex>
        <table>
          <tbody
            className={`${isGameCompleted ? 'match__body match__body--completed' : 'match__body'}`}
          >
            {game.gameParticipants.map((participant) => (
              <BracketParticipant
                key={`bracket_${participant.snakeName}`}
                participant={participant}
                isGameCompleted={isGameCompleted}
              />
            ))}
          </tbody>
        </table>
      </Box>
      <Dialog.Content
        maxWidth="80vw"
        style={{ backgroundColor: 'var(--gray-4)' }}
        className="dialog"
      >
        <Board></Board>
      </Dialog.Content>
    </Dialog.Root>
  )
}

export default function Brackets() {
  const stage = useTournamentStore((store) => store.stage)
  const games = useTournamentStore((store) => store.games)

  return (
    <Box className="brackets">
      <div className="brackets__content">
        <section className="bracket">
          <h3 className="bracket__heading">Групповая стадия - Тройки</h3>
          {games
            .filter((game) => game.stage === Stage.GROUP_THREES)
            .map((game, ix) => (
              <BracketMatch
                key={`bracket_three__${game.id}`}
                game={game}
                heading={`Игра #${ix + 1} (поле ${Object.values(game.field).join('x')})`}
              />
            ))}
        </section>
        <section className="bracket">
          <h3 className="bracket__heading">Групповая стадия - Квинты</h3>
          {games
            .filter((game) => game.stage === Stage.GROUP_FIVES)
            .map((game, ix) => (
              <BracketMatch
                key={`bracket_five__${game.id}`}
                game={game}
                heading={`Игра #${ix + 1} (поле ${Object.values(game.field).join('x')})`}
              />
            ))}
        </section>
        {stage === Stage.FINALS && (
          <section className="bracket">
            <h3 className="bracket__heading">Суперфинал</h3>
            {games
              .filter((game) => game.stage === Stage.FINALS)
              .map((game, ix) => (
                <BracketMatch
                  key={`bracket_final__${game.id}`}
                  game={game}
                  heading={`Игра #${ix + 1} (поле ${Object.values(game.field).join('x')})`}
                />
              ))}
          </section>
        )}
      </div>
    </Box>
  )
}
