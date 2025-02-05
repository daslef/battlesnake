import fetcher from './fetcher'

export async function loadGameInfo(engineURL: string) {
  const gameInfoUrl = `${engineURL}/game`

  try {
    const response = await fetcher(gameInfoUrl)
    return await response.json()
  } catch (e: unknown) {
    console.error(e)
    throw e
  }
}

export async function setWs(engineURL: string) {
  try {
    const wsUrl = engineURL.replace(/^https:\/\//i, 'wss://').replace(/^http:\/\//i, 'ws://')

    const gameEventsUrl = `${wsUrl}/game/events`
    return new WebSocket(gameEventsUrl)
  } catch (e: unknown) {
    console.error(e)
  }
}
