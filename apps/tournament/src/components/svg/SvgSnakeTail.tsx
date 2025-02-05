import React, { useEffect, useState } from 'react'
import type { Snake, SvgCalcParams, } from '../../lib/types'
import { fetchCustomizationSvgDef } from '../../lib/loaders/loadCustomization'
import { svgCalcCellRect, calcDestinationWrapPosition, isAdjacentPoint } from '../../lib/helpers'

interface ISvgSnakeTail {
  snake: Snake
  svgCalcParams: SvgCalcParams
}

const DEFAULT_TAIL = '<path xmlns="http://www.w3.org/2000/svg" d="M50 0H0v100h50l50-50L50 0z"/>'


export const SvgSnakeTail: React.FC<ISvgSnakeTail> = ({ snake, svgCalcParams }) => {
  const [customTail, setCustomTail] = useState(DEFAULT_TAIL)

  function calcDrawTail(snake: Snake): boolean {
    const head = snake.body[0]
    const tail = snake.body[snake.body.length - 1]

    return head.x != tail.x || head.y != tail.y
  }

  function calcTailTransform(snake: Snake): string {
    const tail = snake.body[snake.body.length - 1]

    // Work backwards from the tail until we reach a segment that isn't stacked.
    let preTailIndex = snake.body.length - 2
    let preTail = snake.body[preTailIndex]
    while (preTail.x == tail.x && preTail.y == tail.y) {
      preTailIndex -= 1
      if (preTailIndex < 0) {
        return ''
      }
      preTail = snake.body[preTailIndex]
    }

    // If tail is wrapped we need to calcualte neck position on border
    if (!isAdjacentPoint(preTail, tail)) {
      preTail = calcDestinationWrapPosition(preTail, tail)
    }

    // Return transform based on relative location
    if (preTail.x > tail.x) {
      return 'scale(-1,1) translate(-100,0)'
    } else if (preTail.y > tail.y) {
      return 'scale(-1,1) translate(-100,0) rotate(90, 50, 50)'
    } else if (preTail.y < tail.y) {
      return 'scale(-1,1) translate(-100,0) rotate(-90, 50, 50)'
    }

    return ''
  }

  useEffect(() => {
    fetchCustomizationSvgDef("tail", snake.tail).then(tailHtml => {
      setCustomTail(tailHtml)
    })
  }, [])

  const drawTail = calcDrawTail(snake)
  const tailRectProps = svgCalcCellRect(svgCalcParams, snake.body[snake.body.length - 1])
  const tailTransform = calcTailTransform(snake)

  if (!drawTail) {
    return
  }

  return (
    <svg className="tail" viewBox="0 0 100 100" fill={snake.color} {...tailRectProps}>
      <g transform={tailTransform} dangerouslySetInnerHTML={{ __html: customTail }}>
      </g>
    </svg>

  )
}