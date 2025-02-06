import fetcher from './fetcher'

export async function fetchCustomizationSvgDef(type: string, name: string) {
  const mediaPath = `snakes/${type}s/${name}.svg`
  try {
    const response = await fetcher(`https://media.battlesnake.com/${mediaPath}`)
    const textSVG = await response.text()
    const tempElememt = document.createElement('template')
    tempElememt.innerHTML = textSVG.trim()

    if (tempElememt.content.firstChild === null) {
      console.debug('[customizations] error loading customization, no elements found')
      return ''
    }

    const child = <HTMLElement>tempElememt.content.firstChild
    return child.innerHTML
  } catch (error: unknown) {
    throw error as Error
  }
}
