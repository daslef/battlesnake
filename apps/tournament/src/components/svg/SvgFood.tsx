import React from 'react'
import { type SvgCalcParams, type Point } from '../../lib/types'
import { svgCalcCellCircle } from '../../lib/helpers'

interface ISvgFood {
    svgCalcParams: SvgCalcParams
    point: Point
}

export const SvgFood: React.FC<ISvgFood> = ({ svgCalcParams, point }) => {
    const circleProps = svgCalcCellCircle(svgCalcParams, point)
    const foodRadius = (svgCalcParams.cellSize / 3.25).toFixed(2)
    return (
        <circle id={`food-${point.x}-${point.y}`} className="food fill-rose-500" r={foodRadius} {...circleProps} />
    )
}

