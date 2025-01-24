import { forwardRef } from "react"
import FlipMove from 'react-flip-move';

import jsSvg from '../assets/logos/js.svg'
import digitOneSvg from '../assets/icons/digit-one.svg'
import digitTwoSvg from '../assets/icons/digit-two.svg'
import bowSvg from '../assets/icons/bow.svg'

import { teamIcons } from '../data'


const LeaderboardParticipant = forwardRef(function LeaderboardParticipant(props, ref) {
    const { teamName, languageLogo, score } = props

    return (
        <tr className="leaderboard__participant participant" ref={ref}>
            <td className="participant__image">
                <img src={teamIcons[teamName].icon} alt="" className="participant__language" />
            </td>
            <td className="participant__team">
                {teamName}
            </td>
            <td className="participant__score participant__score--gold">
                <img src={digitOneSvg} alt="" className="participant__score__icon participant__score__icon--gold" />
                <span className="participant__score__value">✖ {score?.gold ?? ""} </span>
            </td>
            <td className="participant__score participant__score--silver">
                <img src={digitTwoSvg} alt="" className="participant__score__icon  participant__score__icon--gold" />
                <span className="participant__score__value">✖ {score?.silver ?? ""}</span>
            </td>
            <td className="participant__score participant__score--bonus">
                <img src={bowSvg} alt="" className="participant__score__icon  participant__score__icon--gold" />
                <span className="participant__score__value">✖ {score?.bonus ?? ""}</span>
            </td>
            <td className="participant__score participant__score--final">
                <span className="participant__score__value participant__score__value--final">{Math.floor(score?.result) ?? ""}</span>
            </td>
        </tr>
    )
})

export default function Leaderboard({ score }) {
    return (
        <article className="leaderboard">
            <h2 className='leaderboard__heading'>Таблица</h2>
            <table>
                <tbody className='leaderboard__content'>
                    <FlipMove staggerDurationBy="30"
                        duration={500}
                        enterLeaveAnimation='accordianVertical'
                        order='desc'
                        typeName={null}
                    >
                        {Object.entries(score).map(([team, teamScore]) => (
                            <LeaderboardParticipant key={`leaderboard_${team}`} languageLogo={jsSvg} teamName={team} score={teamScore} />
                        ))}
                    </FlipMove>
                </tbody>
            </table>
        </article>

    )
}