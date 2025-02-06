import { forwardRef } from "react"
import { Text, Section } from '@radix-ui/themes'
import FlipMove from 'react-flip-move';

import { useTournamentStore } from "../lib/stores/tournament";

import { type Score } from '../lib/types'

import digitOneSvg from '../assets/icons/digit-one.svg'
import digitTwoSvg from '../assets/icons/digit-two.svg'
import bowSvg from '../assets/icons/bow.svg'


interface ILeaderboardParticipant {
    teamName: string,
    languageLogo: string,
    score: Score
}

const LeaderboardParticipant = forwardRef<HTMLTableRowElement, ILeaderboardParticipant>(function LeaderboardParticipant(props, ref) {
    const { teamName, languageLogo, score } = props

    return (
        <tr className="leaderboard__participant participant" ref={ref}>
            <td className="participant__image">
                <img src={languageLogo} alt="" className="participant__language" />
            </td>
            <td className="participant__team">
                {teamName}
            </td>
            <td className="participant__score participant__score--gold">
                <img src={digitOneSvg} alt="" className="participant__score__icon participant__score__icon--gold" />
                <span className="participant__score__value">✖ {score.firstPlaces} </span>
            </td>
            <td className="participant__score participant__score--silver">
                <img src={digitTwoSvg} alt="" className="participant__score__icon  participant__score__icon--gold" />
                <span className="participant__score__value">✖ {score.secondPlaces}</span>
            </td>
            <td className="participant__score participant__score--bonus">
                <img src={bowSvg} alt="" className="participant__score__icon  participant__score__icon--gold" />
                <span className="participant__score__value">✖ {score.aggressiveBonuses}</span>
            </td>
            <td className="participant__score participant__score--final">
                <span className="participant__score__value participant__score__value--final">{Math.floor(score.total)}</span>
            </td>
        </tr>
    )
})

export default function Leaderboard() {
    const score = useTournamentStore(store => store.score)
    return (
        <Section className="leaderboard">
            <table>
                <tbody className='leaderboard__content' style={{ position: "relative" }}>
                    <FlipMove staggerDurationBy="30"
                        duration={500}
                        // enterLeaveAnimation='accordianVertical'
                        // order='desc'
                        typeName={null}
                    >
                        {[...score.entries()]?.map(([participant, participantScore]) => (
                            <LeaderboardParticipant key={`leaderboard_${participant.snakeAuthor}`} languageLogo={participant.languageIcon} teamName={participant.snakeAuthor} score={participantScore} />
                        ))}
                    </FlipMove>
                </tbody>
            </table>
        </Section>
    )
}