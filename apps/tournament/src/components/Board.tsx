import React, { useEffect } from 'react'
import { Box } from '@radix-ui/themes'

import { type SvgCalcParams, PlaybackMode, Elimination } from '../lib/types'

import { SvgSnake } from './svg/SvgSnake'
import { SvgFood } from './svg/SvgFood'
import { SvgGrid } from './svg/SvgGrid'

import iconPlay from '../assets/icons/play.svg'
import iconPause from '../assets/icons/pause.svg'
import iconNext from '../assets/icons/chevron-right.svg'
import iconPrev from '../assets/icons/chevron-left.svg'
import iconFirst from '../assets/icons/chevron-left-double.svg'
import iconLast from '../assets/icons/chevron-right-double.svg'

import { usePlaybackStore } from '../lib/stores/playback'
import { useSettingsStore } from '../lib/stores/settings'

function GameScore() {
  function snakeIdToName(id: string) {
    for (let i = 0; i < currentFrame.snakes.length; i++) {
      if (currentFrame.snakes[i].id == id) {
        return currentFrame.snakes[i].name
      }
    }
  }

  function eliminationToString(elimination: Elimination) {
    switch (elimination.cause) {
      case 'snake-collision':
        return `Collided with body of ${snakeIdToName(elimination.by)} on Turn ${elimination.turn}`
      case 'snake-self-collision':
        return `Collided with itself on Turn ${elimination.turn}`
      case 'out-of-health':
        return `Ran out of health on Turn ${elimination.turn}`
      case 'head-collision':
        return `Lost head-to-head with ${snakeIdToName(elimination.by)} on Turn ${elimination.turn}`
      case 'wall-collision':
        return `Moved out of bounds on Turn ${elimination.turn}`
      default:
        return elimination.cause
    }
  }

  const playback = usePlaybackStore()
  const currentFrame = playback.frames[playback.currentFrameIndex]

  // We sort snakes by elimination state, then lowercase name alphabetical
  const sortedSnakes =
    currentFrame ?
      [...currentFrame.snakes].sort((a, b) => {
        if (a.isEliminated != b.isEliminated) {
          return a.isEliminated ? 1 : -1
        }
        return a.name.toLowerCase().localeCompare(b.name.toLowerCase())
      }) : []

  // if (!currentFrame) {
  //   return <Box className="gamescore">loading...</Box>
  // }

  return (
    <Box className="gamescore">
      <div className="flex flex-row font-bold text-lg">
        <div className="basis-1/2 text-right">TURN{ }</div>
        <div className="basis-1/2 pl-2">{currentFrame.turn}</div>
      </div>

      {sortedSnakes.map((snake) => (
        <div
          className={`p-2 cursor-pointer rounded-sm flex flex-col gap-2 ${snake.isEliminated ? 'eliminated' : ''}`}
          role="presentation"
          key={`gamescore-${snake.id}`}
        >
          <div className="flex flex-row font-bold">
            <p className="grow">{snake.name}</p>
            <p className="ps-4 text-right">{snake.length}</p>
          </div>
          <div className="flex flex-row text-xs">
            <p className="grow truncate">by {snake.author}</p>
            <p className="text-right">{snake.latency ? `${snake.latency}ms` : ''}</p>
          </div>

          <div className="h-4 text-xs mt-1">
            {snake.elimination ? (
              <p>{eliminationToString(snake.elimination)}</p>
            ) : (
              <div className="text-outline w-full h-full rounded-full bg-neutral-200 dark:bg-neutral-800">
                <div
                  className="transition-all h-full rounded-full text-white ps-2"
                  style={{ background: snake.color, width: `${snake.health}%` }}
                >
                  {snake.health}
                </div>
              </div>
            )}
          </div>
        </div>
      ))}
    </Box>
  )
}

function PlaybackControls() {
  const playback = usePlaybackStore()
  const disableDuringPlayback = playback.mode == PlaybackMode.PLAYING

  return (
    <div>
      <button
        className="mx-2 disabled:text-neutral-400"
        onClick={playback.firstFrame}
        disabled={disableDuringPlayback}
      >
        <img src={iconFirst} className='w-8 h-8' />
      </button>
      <button
        className="mx-2 disabled:text-neutral-400"
        onClick={playback.prevFrame}
        disabled={disableDuringPlayback}
      >
        <img src={iconPrev} className='w-8 h-8' />
      </button>

      {playback.mode == PlaybackMode.PLAYING ? (
        <button className="mx-2" onClick={playback.pause}>
          <img src={iconPause} className='w-8 h-8' />
        </button>
      ) : (
        <button className="mx-2" onClick={playback.play}>
          <img src={iconPlay} className='w-8 h-8' />
        </button>
      )}
      <button
        className="mx-2 disabled:text-neutral-400"
        onClick={playback.nextFrame}
        disabled={disableDuringPlayback}
      >
        <img src={iconNext} className='w-8 h-8' />
      </button>
      <button
        className="mx-2 disabled:text-neutral-400"
        onClick={playback.lastFrame}
        disabled={disableDuringPlayback}
      >
        <img src={iconLast} className='w-8 h-8' />
      </button>
    </div>
  )
}

function Gameboard() {
  const frames = usePlaybackStore(state => state.frames)
  const currentFrameIndex = usePlaybackStore(state => state.currentFrameIndex)

  const CELL_SIZE = 20
  const CELL_SIZE_HALF = CELL_SIZE / 2
  const CELL_SPACING = 4
  const GRID_BORDER = 10

  const svgWidth = frames.length
    ? 2 * GRID_BORDER +
    frames[0].width * CELL_SIZE +
    Math.max(frames[0].width - 1, 0) * CELL_SPACING
    : 0

  const svgHeight = frames.length
    ? 2 * GRID_BORDER +
    frames[0].height * CELL_SIZE +
    Math.max(frames[0].height - 1, 0) * CELL_SPACING
    : 0

  const svgCalcParams = {
    cellSize: CELL_SIZE,
    cellSizeHalf: CELL_SIZE_HALF,
    cellSpacing: CELL_SPACING,
    gridBorder: GRID_BORDER,
    height: svgHeight,
    width: svgWidth
  } as SvgCalcParams

  const currentFrame = frames[currentFrameIndex]

  if (!currentFrame) {
    return <p>loading...</p>
  }

  return (
    <svg className="gameboard flex-shrink" viewBox={`0 0 ${svgWidth} ${svgHeight}`}>
      <SvgGrid
        gridWidth={frames[0].width}
        gridHeight={frames[0].height}
        svgCalcParams={svgCalcParams}
      />

      {currentFrame.snakes
        .filter((snake) => snake.isEliminated)
        .map((snake) => <SvgSnake key={snake.id} snake={snake} svgCalcParams={svgCalcParams} opacity={0.1} />)}

      {currentFrame.snakes
        .filter((snake) => !snake.isEliminated)
        .map((snake) => <SvgSnake key={snake.id} snake={snake} svgCalcParams={svgCalcParams} />)}

      {currentFrame.food.map((foodElement, index) => (
        <SvgFood point={foodElement} key={`food-${index}`} svgCalcParams={svgCalcParams} />
      ))}
    </svg>
  )
}

export default function Board({ children }: { children: React.ReactNode }) {
  const settings = useSettingsStore()
  const playback = usePlaybackStore()

  if (!playback.frames.length) {
    return <Box>loading...</Box>
  }

  return (
    <div className="w-full h-full flex">
      <div className="flex flex-col grow">
        <Gameboard />
        <div className="flex justify-evenly text-xl py-2 px-6">
          <PlaybackControls />
        </div>
      </div>
      {settings.showScoreboard && (
        <div className="basis-full md:basis-[45%] order-first p-2 md:order-last gap-8 pt-2 flex flex-col">
          <Box>
            {children}
          </Box>
          <GameScore />
        </div>
      )}
    </div>
  )
}
