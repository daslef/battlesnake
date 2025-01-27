import React from 'react'
import type { Snake, SvgCalcParams } from '../../lib/types'

import { SvgSnakeBody } from './SvgSnakeBody'
import { SvgSnakeHead } from './SvgSnakeHead'
import { SvgSnakeTail } from './SvgSnakeTail'

interface ISvgSnake {
  snake: Snake
  svgCalcParams: SvgCalcParams
  opacity?: number
}

export const SvgSnake: React.FC<ISvgSnake> = ({ snake, svgCalcParams, opacity }) => {
  return (
    <g id={`snake-${snake.id}`} key={`snake-${snake.id}`} className="snake" style={{ opacity: opacity ?? 1.0 }} >
      <SvgSnakeTail svgCalcParams={svgCalcParams} snake={snake} />
      <SvgSnakeBody svgCalcParams={svgCalcParams} snake={snake} />
      <SvgSnakeHead svgCalcParams={svgCalcParams} snake={snake} />
    </g >

  )
}




