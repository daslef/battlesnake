import React, { useEffect } from 'react'
import type { Snake, SvgCalcParams } from '../../lib/types'
import { fetchCustomizationSvgDef } from '../../lib/loaders/loadCustomization'
import { svgCalcCellRect, calcDestinationWrapPosition, isAdjacentPoint } from '../../lib/helpers'


interface ISvgSnakeHead {
  snake: Snake
  svgCalcParams: SvgCalcParams
}


export const SvgSnakeHead: React.FC<ISvgSnakeHead> = ({ snake, svgCalcParams }) => {
  function calcHeadDirection(snake: Snake): string {
    const [head, neckPoint] = snake.body.slice(0, 2)

    const neck =
      isAdjacentPoint(neckPoint, head)
        ? neckPoint
        : calcDestinationWrapPosition(neckPoint, head)

    // Determine head direction based on relative position of neck and tail.
    // If neck and tail overlap, we return the default direction (right).
    if (head.x < neck.x) {
      return 'left'
    } else if (head.y > neck.y) {
      return 'up'
    } else if (head.y < neck.y) {
      return 'down'
    }
    return 'right'
  }

  function calcHeadTransform(headDirection: string): string {
    switch (headDirection) {
      case 'left':
        return 'scale(-1,1) translate(-100, 0)'
      case 'up':
        return 'rotate(-90, 50, 50)'
      case 'down':
        return 'rotate(90, 50, 50)'
      default:
        return ''
    }
  }

  function calcDrawHeadShadow(snake: Snake): boolean {
    return snake.isEliminated && snake.elimination?.cause == 'snake-self-collision'
  }

  let headSvgDef;

  const headRectProps = svgCalcCellRect(svgCalcParams, snake.body[0])
  const headDirection = calcHeadDirection(snake)
  const headTransform = calcHeadTransform(headDirection)
  const drawHeadShadow = calcDrawHeadShadow(snake)


  useEffect(() => {
    fetchCustomizationSvgDef("head", snake.head).then(html => {
      headSvgDef = html
    })
  }, [])

  if (!headSvgDef) {
    headSvgDef = `<circle fill="none" cx="12.52" cy="28.55" r="9.26"/>
        <path d="M0 100h100L56 55.39l44-39.89V.11L0 0zm12.52-80.71a9.26 9.26 0 1 1-9.26 9.26 9.26 9.26 0 0 1 9.26-9.26z"/>`
  }

  return (
    <svg
      className={`head ${headDirection} ${drawHeadShadow ? 'shadow' : ''}`}
      viewBox="0 0 100 100"
      fill={snake.color}
      {...headRectProps}
    >
      <g transform={headTransform} dangerouslySetInnerHTML={{ __html: headSvgDef }}></g>
    </svg>

  )
}