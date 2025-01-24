import { useMemo } from 'react'

import { k_combinations } from '../lib'

import jsSvg from '../assets/logos/js.svg'
import digitOneSvg from '../assets/icons/digit-one.svg'
import digitTwoSvg from '../assets/icons/digit-two.svg'
import bowSvg from '../assets/icons/bow.svg'

import { teamIcons } from '../data'
import { useState } from 'react'


function BracketParticipant({ teamName, languageLogo, setScore }) {
    const handleClick = (type) => {
        setScore(score => ({ ...score, [teamName]: { ...score[teamName], [type]: score[teamName][type] + 1 } }))
    }

    return (
        <tr className="match__participant participant">
            <td className="participant__image">
                <img src={teamIcons[teamName].icon} alt="programming language" className="participant__language" />
            </td>
            <td className="participant__team">
                {teamName}
            </td>
            <td className="participant__controls participant__controls">
                <button type="button" className="participant__controls__control control" onClick={() => handleClick('gold')}>
                    <img src={digitOneSvg} alt="first place" className="control__icon control__icon--gold" />
                </button>
                <button type="button" className="participant__controls__control control" onClick={() => handleClick('silver')}>
                    <img src={digitTwoSvg} alt="second place" className="control__icon control__icon--silver" />
                </button>
                <button type="button" className="participant__controls__control control" onClick={() => handleClick('bonus')}>
                    <img src={bowSvg} alt="brave bonus" className="control__icon control__icon--bonus" />
                </button>
            </td>
        </tr>
    )
}

function BracketMatch({ teamNames, setScore, heading }) {
    const [completed, setCompleted] = useState(false)

    return (
        <table className="bracket__match match">
            <thead>
                <tr style={{ display: "flex", alignItems: "center", }}>
                    <th className='match__heading' colSpan={5}>{heading}</th>
                    <td style={{ display: "flex", alignItems: "center", gap: "8px", marginLeft: "auto" }}>
                        <label htmlFor="">Completed</label>
                        <input type="checkbox" name="" id="" onChange={() => setCompleted(completed => !completed)} />
                    </td>

                </tr>
            </thead>
            <tbody className={`${completed ? "match__body match__body--completed" : "match__body"}`}>
                {teamNames.map(team => (
                    <BracketParticipant key={`bracket_${team}`} languageLogo={jsSvg} teamName={team} setScore={setScore} />
                ))}
            </tbody>
            <tfoot>
            </tfoot>
        </table>

    )
}

export default function Brackets({ teamNames, sortedScore, setScore }) {

    const matches = useMemo(() => {
        const fields = ['11×11', '19×19']
        const combs = k_combinations(teamNames, 3).toSorted(() => Math.random() - 0.5)

        const matches = {
            threes: [],
            fives: [],
        }

        for (const field of fields) {
            for (const comb of combs) {
                matches.threes.push({
                    participants: comb,
                    field: field
                })
            }
        }


        for (const field of fields) {
            for (let i = 0; i < 5; i++) {
                matches.fives.push({
                    participants: teamNames,
                    field: field
                })
            }
        }

        return matches

    }, [])

    const finals = useMemo(() => {
        const fields = ['7×7', '11×11', '19×19']
        const matches = []

        for (const field of fields) {
            for (let i = 0; i < 3; i++) {
                matches.push({
                    participants: Object.keys(sortedScore).slice(0, 2),
                    field: field
                })
            }
        }

        return matches

    }, [sortedScore])


    return (
        <article className="brackets">
            <h2 className='brackets__heading'>Сетка</h2>
            <div className='brackets__content'>
                <section className="bracket">
                    <h3 className="bracket__heading">Групповая стадия - Тройки</h3>
                    {
                        matches.threes.map((match, ix) => (
                            <BracketMatch key={`bracket_three_${ix}_${match.participants.join("_")}`} setScore={setScore} teamNames={match.participants} heading={`Игра #${ix + 1} (поле ${match.field})`} />
                        ))
                    }
                </section>
                <section className="bracket">
                    <h3 className="bracket__heading">Групповая стадия - Квинты</h3>
                    {
                        matches.fives.map((match, ix) => (
                            <BracketMatch key={`bracket_five_${ix}_${match.participants.join("_")}`} setScore={setScore} teamNames={match.participants} heading={`Игра #${ix + 1} (поле ${match.field})`} />
                        ))
                    }
                </section>
                <section className="bracket">
                    <h3 className="bracket__heading">Суперфинал</h3>
                    {
                        finals.map((match, ix) => (
                            <BracketMatch key={`bracket_final_${ix}_${match.participants.join("_")}`} setScore={setScore} teamNames={match.participants} heading={`Игра #${ix + 1} (поле ${match.field})`} />
                        ))
                    }

                </section>
            </div>
        </article>

    )
}