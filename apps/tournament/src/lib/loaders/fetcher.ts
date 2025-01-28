export default async function fetcher(url: string, options: object = {}, retries: number = 10) {
    try {
        const response = await fetch(url, options)
        if (response.ok) {
            return response
        }
        throw new Error("Response is not OK")
    } catch (error: unknown) {
        console.error((error as Error).message)
        if (retries > 0) {
            return fetcher(url, options, retries - 1)
        }
        throw error as Error
    }
}