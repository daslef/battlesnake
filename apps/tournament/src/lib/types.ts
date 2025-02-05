export type SvgCalcParams = {
  cellSize: number
  cellSizeHalf: number
  cellSpacing: number
  gridBorder: number
  height: number
  width: number
}

export type SvgPoint = {
  x: number
  y: number
}

export type SvgCircleProps = {
  cx: number
  cy: number
}

export type SvgRectProps = {
  x: number
  y: number
  width: number
  height: number
}

export enum Setting {
  ENGINE = 'engine',
  FPS = 'fps',
  GAME = 'game',
  SHOW_SCOREBOARD = 'showScoreboard',
  THEME = 'theme',
  TITLE = 'title'
}

export type Point = {
  x: number
  y: number
}

export type Elimination = {
  turn: number
  cause: string
  by: string
}

export type EnginePoint = { X: number; Y: number }

export type EngineSnake = {
  ID: string
  Name: string
  Author: string
  Color: string
  HeadType: string
  TailType: string
  // Frame specific
  Health: number
  Latency: number
  Body: EnginePoint[]
  Death: EngineDeath
}

export type Snake = {
  id: string
  name: string
  author: string
  color: string
  head: string
  tail: string
  health: number
  latency: number
  body: Point[]
  length: number
  elimination: Elimination | null
  // Helpers
  isEliminated: boolean
}

export type EngineGame = {
  Width: number
  Height: number
}

export type EngineGameEvent = {
  Turn: number
  Snakes: EngineSnake[]
  Food: EnginePoint[]
}

export type EngineGameInfo = {
  Turn: number
  Game: EngineGame
  Snakes: EngineSnake[]
}

export type EngineDeath = {
  Turn: number
  Cause: string
  EliminatedBy: string
}

export type Frame = {
  turn: number
  width: number
  height: number
  snakes: Snake[]
  food: Point[]
  isFinalFrame: boolean
}

export enum PlaybackMode {
  NEW,
  READY,
  PAUSED,
  PLAYING,
  FINISHED
}

export function engineEventToFrame(
  engineGameInfo: EngineGameInfo,
  engineGameEvent: EngineGameEvent
): Frame {
  return {
    turn: engineGameEvent.Turn,
    width: engineGameInfo.Game.Width,
    height: engineGameInfo.Game.Height,
    snakes: engineGameEvent.Snakes.map((engineSnake: EngineSnake) => ({
      // Fixed properties
      id: engineSnake.ID,
      name: engineSnake.Name,
      author: engineSnake.Author,
      color: engineSnake.Color,
      head: engineSnake.HeadType,
      tail: engineSnake.TailType,
      // Frame specific
      health: engineSnake.Health,
      latency: engineSnake.Latency,
      body: engineSnake.Body.map((engineCoords) => ({ x: engineCoords.X, y: engineCoords.Y })),
      length: engineSnake.Body.length,
      elimination: engineSnake.Death
        ? {
            turn: engineSnake.Death.Turn,
            cause: engineSnake.Death.Cause,
            by: engineSnake.Death.EliminatedBy
          }
        : null,
      // Helpers
      isEliminated: engineSnake.Death != null
    })),
    food: engineGameEvent.Food.map((engineCoords) => ({ x: engineCoords.X, y: engineCoords.Y })),
    isFinalFrame: false
  }
}
