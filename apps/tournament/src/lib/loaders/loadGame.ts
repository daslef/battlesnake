import ReconnectingWebSocket from 'reconnecting-websocket'
import fetcher from './fetcher'

export async function loadGameInfo(engineURL: string) {
    const gameInfoUrl = `${engineURL}/game`

    try {
        const response = await fetcher(gameInfoUrl)
        return await response.json()
    } catch (e: unknown) {
        console.error(e)
    }
}

export async function loadGameEvents(engineURL: string) {
    try {
        const wsUrl = engineURL
            .replace(/^https:\/\//i, 'wss://')
            .replace(/^http:\/\//i, 'ws://')

        const gameEventsUrl = `${wsUrl}/game/events`
        return new ReconnectingWebSocket(gameEventsUrl)
    } catch (e: unknown) {
        console.error(e)
    }

}