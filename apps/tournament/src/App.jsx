import { useEffect } from 'react'
import { Flex } from '@radix-ui/themes'

import Leaderboard from './components/Leaderboard'
import Brackets from './components/Brackets'

import { useTournamentStore } from './lib/stores/tournament'

function App() {
  const initialize = useTournamentStore(store => store.initialize)

  useEffect(() => {
    initialize()
    // return () => { } to unset ? 
  }, [])

  return (
    <Flex style={{ justifyContent: "center", gap: "2vh", padding: "20px", alignItems: "flex-start" }}>
      <Leaderboard />
      <Brackets />
    </Flex>
  )
}

export default App
