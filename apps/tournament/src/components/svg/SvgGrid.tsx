import React from 'react'

import {
  type SvgCalcParams,
} from '../../lib/types'

import {
  svgCalcCellRect,
} from '../../lib/helpers'


interface ISvgGrid {
  gridWidth: number
  gridHeight: number
  svgCalcParams: SvgCalcParams
}

export const SvgGrid: React.FC<ISvgGrid> = ({ gridWidth, gridHeight, svgCalcParams }) => {
  function generateProduct() {
    const result = []
    for (let x = 0; x < gridWidth; x++) {
      for (let y = 0; y < gridHeight; y++) {
        result.push({ x, y })
      }
    }
    return result
  }
  return (
    <g>
      {generateProduct().map(({ x, y }) => {
        return (
          <rect
            id={`grid-${x}-${y}`}
            key={`grid-${x}-${y}`}
            className="grid fill-[#f1f1f1] dark:fill-[#393939]"
            {...svgCalcCellRect(svgCalcParams, { x, y })}
          />
        )
      })}
    </g>
  )
}