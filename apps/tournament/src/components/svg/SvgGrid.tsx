import React from 'react'

import {
  type SvgCalcParams,
} from '../../lib/types'

import {
  svgCalcCellRect,
  svgCalcCellLabelBottom,
  svgCalcCellLabelLeft
} from '../../lib/helpers'


interface ISvgGrid {
  gridWidth: number
  gridHeight: number
  showLabels: boolean
  svgCalcParams: SvgCalcParams
}

export const SvgGrid: React.FC<ISvgGrid> = ({ gridWidth, gridHeight, showLabels, svgCalcParams }) => {
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

      {showLabels && Array(gridHeight).map((_, x) => (
        <text
          className="coordinate-label text-[0.35rem] fill-neutral-500"
          text-anchor="middle"
          transform="translate(0, 2)"
          key={`label_x_${x}`}
          {...svgCalcCellLabelBottom(svgCalcParams, { x: x, y: 0 })}
        >
          {x}
        </text>
      ))}

      {showLabels && Array(gridHeight).map((_, y) => (
        <text
          className="coordinate-label text-[0.35rem] fill-neutral-500"
          text-anchor="middle"
          transform="translate(0, 2)"
          key={`label_y_${y}`}
          {...svgCalcCellLabelLeft(svgCalcParams, { x: 0, y: y })}
        >
          {y}
        </text>
      ))}
    </g>
  )
}