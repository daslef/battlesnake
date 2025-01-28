import { useState, useEffect, useMemo } from 'react'

import { useThemeStore } from './lib/stores/settings'

import Leaderboard from './components/Leaderboard'
import Brackets from './components/Brackets'
import Board from './components/Board'

import { teamNames } from './data'

import './App.css'

function App() {
  const { theme } = useThemeStore()

  const [score, setScore] = useState(() => {
    const initialScore = {
      gold: 0, silver: 0, bonus: 0, result: 0
    }

    return teamNames.reduce((acc, team) => {
      return { ...acc, [team]: { ...initialScore } }
    }, {})
  })

  const sortedScore = useMemo(() => {
    return Object.fromEntries(Object.entries(score).toSorted((a, b) => b[1].result - a[1].result))
  }, [JSON.stringify(score)])

  useEffect(() => {
    setScore(score => {
      const newScore = { ...score }
      for (const team of teamNames) {
        const { gold, silver, bonus } = score?.[team]
        newScore[team].result = gold * 2 + silver * 1 + bonus * 0.001
      }
      return newScore
    })

  }, [JSON.stringify(score)])

  return (
    <main className={`app ${theme}`}>
      <button type="button" onClick={() => {
        fetch('http://localhost:5000/create')
      }}>Create!</button>
      <Board engineUrl='http://127.0.0.1:5000' />
      <Leaderboard score={sortedScore} />
      <Brackets teamNames={teamNames} setScore={setScore} sortedScore={sortedScore} />
    </main>
  )
}

export default App
