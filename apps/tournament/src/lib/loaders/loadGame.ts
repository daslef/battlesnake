export async function loadGameInfo(engineURL: string, gameID: string) {
    const gameInfoUrl = `${engineURL}/games/${gameID}`

    try {
        const response = await fetch(gameInfoUrl)

        if (response.status == 404) {
            throw new Error('Game not found')
        } else if (!response.ok) {
            throw new Error('Error loading game')
        }
        const gameInfo = await response.json()
        return gameInfo
    } catch (e: unknown) {
        console.error(e)
        throw e as Error
    }
}

export async function loadGameEvents(engineURL: string, gameID: string) {
    try {
        const wsUrl = engineURL
            .replace(/^https:\/\//i, 'wss://')
            .replace(/^http:\/\//i, 'ws://')

        const gameEventsUrl = `${wsUrl}/games/${gameID}/events`
        const ws = new WebSocket(gameEventsUrl)
        return ws
    } catch (e: unknown) {
        console.error(e)
        throw e as Error
    }

}