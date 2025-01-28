import { SvgCalcParams, Point, SvgPoint, SvgRectProps, SvgCircleProps } from "./types"

export function svgCalcCellCenter(params: SvgCalcParams, p: Point): SvgPoint {
  const topLeft = svgCalcCellTopLeft(params, p)
  return {
    x: topLeft.x + params.cellSizeHalf,
    y: topLeft.y + params.cellSizeHalf
  }
}

export function svgCalcCellTopLeft(params: SvgCalcParams, p: Point): SvgPoint {
  return {
    x: params.gridBorder + p.x * (params.cellSize + params.cellSpacing),
    y:
      params.height -
      (params.gridBorder + p.y * (params.cellSize + params.cellSpacing) + params.cellSize)
  }
}

export function svgCalcCellCircle(params: SvgCalcParams, p: Point): SvgCircleProps {
  const center = svgCalcCellCenter(params, p)
  return { cx: center.x, cy: center.y }
}

export function svgCalcCellRect(params: SvgCalcParams, p: Point): SvgRectProps {
  const topLeft = svgCalcCellTopLeft(params, p)
  return { x: topLeft.x, y: topLeft.y, width: params.cellSize, height: params.cellSize }
}


export function fromLocalStorage(key: string, defaultValue: boolean | number | string) {
  const val = localStorage.getItem(`setting.${key}`)
  if (val) {
    return JSON.parse(val)
  }
  return defaultValue
}

export function toLocalStorage(key: string, value: boolean | number | string) {
  localStorage.setItem(`setting.${key}`, JSON.stringify(value))
}

export function getBoolFromURL(url: URL, key: string, defaultValue: boolean): boolean {
  const val = url.searchParams.get(key)
  if (val) {
    if (val === 'true') return true
    if (val === 'false') return false
  }
  return defaultValue
}

export function getIntFromURL(url: URL, key: string, defaultValue: number): number {
  const val = url.searchParams.get(key)
  if (val) {
    const parsedVal = parseInt(val)
    if (!isNaN(parsedVal)) {
      return parsedVal
    }
  }
  return defaultValue
}

export function getStringFromURL(url: URL, key: string, defaultValue: string): string {
  return url.searchParams.get(key) ?? defaultValue
}

export function isEqualPoint(p1?: Point, p2?: Point): boolean {
  if (p1 == undefined || p2 == undefined) {
    return false
  }

  return p1.x == p2.x && p1.y == p2.y
}

export function isAdjacentPoint(p1: Point, p2: Point): boolean {
  return calcManhattan(p1, p2) == 1
}

export function calcManhattan(p1: Point, p2: Point): number {
  return Math.abs(p1.x - p2.x) + Math.abs(p1.y - p2.y)
}

export function calcSourceWrapPosition(src: Point, dst: Point): Point {
  return {
    x: src.x - Math.sign(dst.x - src.x),
    y: src.y - Math.sign(dst.y - src.y)
  }
}

export function calcDestinationWrapPosition(src: Point, dst: Point): Point {
  return {
    x: dst.x + Math.sign(dst.x - src.x),
    y: dst.y + Math.sign(dst.y - src.y)
  }
}
