export enum Stage {
  NOT_SET,
  GROUP_THREES,
  GROUP_FIVES,
  GROUP_ALL
}

export enum Language {
  CSharp,
  Python,
  JavaScript
}

export type Participant = {
  snakeName: string
  snakeUrl: URL
  snakeAuthor: string
  language: Language
  languageIcon: string
  photo?: URL
}

export enum GameStatus {
  NOT_PLAYED,
  IN_PROGRESS,
  COMPLETED,
  ERROR
}

export type Field = {
  width: number
  height: number
}

export type GameResult = {
  firstPlace: Participant | null
  secondPlace: Participant | null
  aggressiveBonus: Participant | null
}

export type Game = {
  id?: string
  stage: Stage
  field: Field
  status: GameStatus
  gameParticipants: Participant[]
  result?: GameResult
}

export type Score = {
  firstPlaces: number
  secondPlaces: number
  aggressiveBonuses: number
  total: number
}


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

export type GameServerObject = {
  snakes: { snakeUrl: string; snakeName: string }[]
  field: Field
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
  PAUSED,
  PLAYING
}

